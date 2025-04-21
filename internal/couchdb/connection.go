package couchdb

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"maps"
	"net/http"
	"net/url"
	"os"
	"reflect"

	"github.com/lampctl/go-sse"
)

// Connection provides basic functionality to use CouchDB. A single instance is
// safe to use from multiple goroutines.
type Connection struct {
	dbURL       *url.URL
	initialized bool
}

// DefaultConnection is a Connection that is initialized with the default
// connection URL from environment variables.
//
// This is not guaranteed to be valid. Client code can call [AssertInitialized]
// if it depends on this value to be valid.
var DefaultConnection Connection

func (c Connection) createDB(ctx context.Context) error {
	resp, err := c.req(ctx, "PUT", c.dbURL.String(), nil, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	switch resp.StatusCode {
	case 201, 202, 412: // 412 means the database already exists.
		return nil
	case 400:
		// This is an unrecoverable error. The configured database name is not a
		// valid name in couchdb.
		panic("couchdb: invalid configuration: invalid database name")
	case 401:
		// This is an unrecoverable error. The configured credentials are wrong.
		panic("couchdb: invalid configuration: bad credentials")
	default:
		panic("couchdb: unable to bootstrap database")
	}
}

type View struct {
	Map    string `json:"map,omitempty"`
	Reduce string `json:"reduce,omitempty"`
}

type Views map[string]View
type Filters map[string]string

type designDoc struct {
	Views   Views   `json:"views,omitempty"`
	Filters Filters `json:"filters"`
}

// aggregateEventsFilter retrieves domain events stored with aggregate entities in
// couchdb.
const aggregateEventsFilter = `function(doc, req) { 
	return doc.events && doc.events.length
}`

const newEventFilter = `function(doc, req) {
	return doc._id.startsWith("domain_event:") && !doc.published_at
}`

func (c Connection) createViews(ctx context.Context) error {
	var doc designDoc = designDoc{
		Filters: Filters{
			"aggregate_events":          aggregateEventsFilter,
			"unpublished_domain_events": newEventFilter,
		},
	}
	return c.SetDesignDoc(ctx, "events", doc)
}

func (c Connection) SetDesignDoc(ctx context.Context, id string, doc designDoc) error {
	var existing designDoc
	path := fmt.Sprintf("_design/%s", id)
	rev, err := c.Get(ctx, path, &existing)
	if errors.Is(err, ErrNotFound) {
		_, err = c.Insert(ctx, path, doc)
	} else {
		if !reflect.DeepEqual(doc, existing) {
			_, err = c.Update(ctx, path, rev, doc)
		}
	}
	return err
}

type changeEventChange struct {
	Rev     string `json:"rev"`
	Deleted bool   `json:"deleted,omitempty"` // omitempty probably irrelevant, as we only read
}

// ChangeEvent is a record emitted by subscribing to /{db}/_changes server-sent
// events
type ChangeEvent struct {
	Seq     string            `json:"seq"`
	ID      string            `json:"id"`
	Changes []json.RawMessage `json:"changes"`
	Doc     json.RawMessage   `json:"doc,omitempty"` // Included if include_docs options is used
}

type changeOption func(*url.Values)

// ChangeOptViewFilter specifies to filter on documents for which the map function of view in design
// document ddoc produces a value.
func ChangeOptViewFilter(ddoc, view string) changeOption {
	return func(v *url.Values) {
		v.Set("filter", "_view")
		v.Set("view", fmt.Sprintf("%s/%s", ddoc, view))
	}
}

// ChangeOptFilter specifies to filter the events using filter function on
// design document ddoc.
func ChangeOptFilter(ddoc, filter string) changeOption {
	return func(v *url.Values) {
		v.Set("filter", fmt.Sprintf("%s/%s", ddoc, filter))
	}
}

func ChangeOptIncludeDocs() changeOption {
	return func(v *url.Values) { v.Set("include_docs", "true") }
}

func getChangeEvents(ctx context.Context, ch <-chan *sse.Event) <-chan ChangeEvent {
	res := make(chan ChangeEvent)
	go func() {
		defer close(res)
		for e := range ch {
			// Ignore heartbeat events
			if e.Data != "" {
				var cev ChangeEvent
				err := json.Unmarshal([]byte(e.Data), &cev)
				// slog.InfoContext(ctx, "couchdb: process event", "event", e.Data)
				if err != nil {
					slog.ErrorContext(ctx, "couchdb: process event", "err", err, "event", e.Data)
					continue
				}
				res <- cev
			}
		}
	}()
	return res
}

// Changes subscribe to change events from CouchDB.
func (c Connection) Changes(
	ctx context.Context,
	options ...changeOption,
) (<-chan ChangeEvent, error) {
	u := c.dbURL.JoinPath("_changes")
	q := u.Query()
	q.Set("feed", "eventsource")
	q.Set("since", "0")
	for _, o := range options {
		o(&q)
	}
	u.RawQuery = q.Encode()
	conn, err := sse.NewClientFromURL(u.String())
	if err != nil {
		return nil, err
	}

	go func() {
		<-ctx.Done()
		slog.InfoContext(ctx, "couchdb: closing event stream")
		conn.Close()
	}()
	return getChangeEvents(ctx, conn.Events), nil
}

// Bootstrap creates the database, as well as updates any design documents, such
// as views. Panics on unrecoverable errors, e.g., an invalid configuration.
func (c Connection) Bootstrap(ctx context.Context) error {
	if err := c.createDB(ctx); err != nil {
		return err
	}
	if err := c.createViews(ctx); err != nil {
		return err
	}

	return nil
}

func NewCouchConnection(couchURL string) (conn Connection, err error) {
	if couchURL == "" {
		err = errors.New("couchdb: NewCouchConnection: empty couchURL")
		return
	}
	var url *url.URL
	url, err = url.Parse(couchURL)
	conn = Connection{url, false}
	if err == nil {
		if err = conn.Bootstrap(context.Background()); err == nil {
			conn.initialized = true
		}
	}
	return
}

// docURL generates the full couchDB resource URL for a given document ID. The
// full URL will include both database name and credentials.
func (c Connection) docURL(id string) string {
	return c.dbURL.JoinPath(id).String()
}

// Insert creates a new document with the specified id. If the operation
// succeeds, the revision of the new document is returned in the rev return
// value.
func (c Connection) Insert(ctx context.Context, id string, doc any) (rev string, err error) {
	if id == "" {
		err = errors.New("couchdb: missing id")
		return
	}
	var b bytes.Buffer
	enc := json.NewEncoder(&b)
	if err = enc.Encode(doc); err != nil {
		return
	}

	url := c.docURL(id)
	var resp *http.Response
	if resp, err = c.req(ctx, "PUT", url, nil, &b); err != nil {
		return
	}
	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.ErrorContext(ctx, "couchdb: insert: reading response body", "err", err)
	}

	switch resp.StatusCode {
	case 201:
		etag := resp.Header.Get("Etag")
		err = json.Unmarshal([]byte(etag), &rev)
		if err != nil {
			err = fmt.Errorf("couchdb: Insert: unable to parse etag \"%s\" : %w", etag, err)
		}
	case 409:
		err = ErrConflict
	default:
		slog.ErrorContext(
			ctx, "couchdb: insert failed",
			"status", resp.StatusCode,
			"resp", string(respBody),
		)
		err = fmt.Errorf("couchdb: insert id(%s): %w", id, errUnexpectedStatusCode(resp))
		return
	}
	return
}

func (c Connection) RawPost(ctx context.Context, path string, body any) (*http.Response, error) {
	var reader io.Reader
	reader, ok := body.(io.Reader)
	if !ok {
		var b bytes.Buffer
		enc := json.NewEncoder(&b)
		if err := enc.Encode(body); err != nil {
			return nil, err
		}
		reader = &b
	}

	var header = make(http.Header)
	header.Add("Content-Type", "application/json")
	return c.req(ctx, "POST", c.docURL(path), header, reader)
}

func (c Connection) Get(ctx context.Context, id string, doc any) (rev string, err error) {
	var resp *http.Response
	if resp, err = c.req(ctx, "GET", c.docURL(id), nil, nil); err != nil {
		return
	}
	defer resp.Body.Close()
	switch resp.StatusCode {
	case 200:
		var bodyBytes []byte
		bodyBytes, err = io.ReadAll(resp.Body)
		if err != nil {
			return
		}
		cd := couchDoc{}
		if err = json.Unmarshal(bodyBytes, &cd); err == nil {
			err = json.Unmarshal(bodyBytes, &doc)
		}
		rev = cd.Rev
	case 404:
		err = fmt.Errorf("%w: %s", ErrNotFound, id)
	default:
		err = fmt.Errorf("couchdb: get(%s): %w", id, errUnexpectedStatusCode(resp))
	}
	return
}

func (c Connection) GetPath(path string, q url.Values, doc any) (rev string, err error) {
	var resp *http.Response
	u := c.dbURL.JoinPath(path)
	u.RawQuery = q.Encode()
	if resp, err = http.Get(u.String()); err != nil {
		return
	}
	defer resp.Body.Close()
	switch resp.StatusCode {
	case 200:
		var bodyBytes []byte
		bodyBytes, err = io.ReadAll(resp.Body)
		if err != nil {
			return
		}
		cd := couchDoc{}
		if err = json.Unmarshal(bodyBytes, &cd); err == nil {
			err = json.Unmarshal(bodyBytes, &doc)
		}
		rev = cd.Rev
	case 404:
		err = fmt.Errorf("%w: %s", ErrNotFound, path)
	default:
		err = fmt.Errorf("couchdb: get(%s): %w", path, errUnexpectedStatusCode(resp))
	}
	return
}

// Update updates the document in the database. If successful, it will return
// the updated revision of the document. If there is a conflict, it will return
// ErrConflict
func (c Connection) Update(
	ctx context.Context,
	id, oldRev string,
	doc any,
) (newRev string, err error) {
	var (
		b bytes.Buffer
		// req  *http.Request
		resp *http.Response
	)
	enc := json.NewEncoder(&b)
	if err = enc.Encode(doc); err != nil {
		return
	}
	var header = make(http.Header)
	header.Add("If-Match", oldRev)
	if resp, err = c.req(ctx, "PUT", c.docURL(id), header, &b); err != nil {
		return
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case 201:
		var bodyBytes []byte
		bodyBytes, err = io.ReadAll(resp.Body)
		if err != nil {
			return
		}
		cd := couchDoc{}
		if err = json.Unmarshal(bodyBytes, &cd); err == nil {
			err = json.Unmarshal(bodyBytes, &doc)
		}
		newRev = cd.Rev
	case 409:
		err = ErrConflict
	default:
		err = fmt.Errorf("couchdb: update %s: unexpected http status code: %d", id, resp.StatusCode)
	}
	return
}

// req wraps http NewRequest and Client.Do method, but converts the error to an ErrConn. This
// makes the caller able to distinguish between:
//   - The couchdb server responded with an unexpected status code.
//   - The call failed, and no response could be retrieved from couch DB
func (c Connection) req(
	ctx context.Context,
	method, url string,
	headers http.Header,
	body io.Reader,
) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		err = fmt.Errorf("%w: %w: %v", ErrConn, ErrRequest, err)
		return nil, err
	}
	if headers != nil {
		maps.Copy(req.Header, headers)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		err = fmt.Errorf("%w: %v", ErrConn, err)
	}
	return resp, err
}

// AssertInitialized verifies that a DefaultConnection exists. The function
// panics if no defualt connection was initialized.
func AssertInitialized() {
	if !DefaultConnection.initialized {
		panic("couchdb: DefaultConnection not initialized")
	}
}

func init() {
	couchURL := os.Getenv("COUCHDB_URL")
	if couchURL == "" {
		// Default local test DB
		couchURL = "http://admin:password@localhost:5984/harmony-test"
	}
	conn, err := NewCouchConnection(couchURL)
	if err == nil {
		DefaultConnection = conn
	} else {
		slog.Error("couchdb: Error initializing", "err", err)
	}
}
