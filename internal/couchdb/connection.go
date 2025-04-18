package couchdb

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"harmony/internal/domain"
	"io"
	"log/slog"
	"maps"
	"net/http"
	"net/url"
	"os"

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

type view struct {
	Map      string `json:"map,omitempty"`
	Reduce   string `json:"reduce,omitempty"`
	ReReduce string `json:"rereduce,omitempty"`
}

type views map[string]view
type filters map[string]string

type designDoc struct {
	Views   views   `json:"views"`
	Filters filters `json:"filters"`
	updated bool    `json:"-"`
}

func (d *designDoc) setView(name, mapFn string) {
	if d.Views == nil {
		d.Views = make(views)
	}
	view, ok := d.Views[name]
	if ok {
		if view.Map == mapFn {
			return
		}
	}

	view.Map = mapFn
	d.Views[name] = view
	d.updated = true
}

func (d *designDoc) setFilter(name, fn string) {
	if d.Filters == nil {
		d.Filters = make(filters)
	}
	if filter, ok := d.Filters[name]; !ok || filter != fn {
		d.Filters[name] = fn
		d.updated = true
	}
}

const mapUnpublishedEvents = `function(doc) { 
	if (doc.events) { 
		for (const e of doc.events) { 
			emit(e.id, e) 
		} 
	} 
}`

const newEventFilter = `function(doc, req) {
	return doc._id.startsWith("domain_event:")
}`

func updateEventsDesignDoc(doc *designDoc) {
	doc.setView("unpublished_events", mapUnpublishedEvents)
	doc.setFilter("domain_events", newEventFilter)
}

func (c Connection) createViews(ctx context.Context) error {
	var doc designDoc
	rev, err := c.Get("_design/events", &doc)
	if err == ErrNotFound {
		doc = designDoc{}
		updateEventsDesignDoc(&doc)
		_, err = c.Insert(ctx, "_design/events", doc)
	} else {
		updateEventsDesignDoc(&doc)
		if doc.updated {
			_, err = c.Update(ctx, "_design/events", rev, doc)
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

type DocumentWithEvents[T any] struct {
	ID       string         `json:"_id,omitempty"`
	Rev      string         `json:"_rev,omitempty"`
	Document T              `json:"doc"`
	Events   []domain.Event `json:"events,omitempty"`
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
	q.Set("since", "now")
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

func getNewEntityEvents(
	ctx context.Context,
	ch <-chan ChangeEvent,
) <-chan DocumentWithEvents[json.RawMessage] {
	res := make(chan DocumentWithEvents[json.RawMessage])
	go func() {
		defer close(res)
		for changeEvent := range ch {
			var doc DocumentWithEvents[json.RawMessage]
			err := json.Unmarshal(changeEvent.Doc, &doc)
			if err != nil {
				slog.ErrorContext(ctx, "couchdb: process event document", "err", err)
				continue
			}
			res <- doc
		}
	}()
	return res
}

func (c Connection) processNewEntity(ctx context.Context, doc DocumentWithEvents[json.RawMessage]) {
	for _, domainEvent := range doc.Events {
		_, err := c.Insert(ctx, "domain_event:"+string(domainEvent.ID), domainEvent)
		if err != nil && !errors.Is(err, ErrConflict) {
			slog.ErrorContext(ctx, "couchdb: insert domain event", "err", err)
			return
		}
	}
	doc.Events = nil
	_, err := c.Update(ctx, doc.ID, doc.Rev, doc)
	if err != nil {
		slog.ErrorContext(ctx, "couchdb: process event", "err", err)
		return
	}
}

func (c Connection) processNewDomainEvents(ctx context.Context) (err error) {
	ch, err := c.Changes(
		ctx,
		ChangeOptViewFilter("events", "unpublished_events"),
		ChangeOptIncludeDocs(),
	)
	if err != nil {
		return
	}
	go func() {
		for doc := range getNewEntityEvents(ctx, ch) {
			c.processNewEntity(ctx, doc)
		}
	}()
	return nil
}

// domainEventsOfChangeEvents takes a channel of CouchDB change events, assumed
// to contain new domain event documents, and transforms it to a channel of
// [domain.Event]
func domainEventsOfChangeEvents(ctx context.Context, ch <-chan ChangeEvent) <-chan domain.Event {
	cha := make(chan domain.Event)
	go func() {
		defer close(cha)
		for changeEvent := range ch {
			var ev domain.Event
			err := json.Unmarshal(changeEvent.Doc, &ev)
			if err != nil {
				slog.ErrorContext(ctx, "couchdb: process event", "err", err)
				continue
			}
			cha <- ev
		}
	}()
	return cha
}

func (c Connection) processUnpublishedDomainEvents(
	ctx context.Context,
) (<-chan domain.Event, error) {
	ch, err := c.Changes(
		ctx,
		ChangeOptFilter("events", "domain_events"),
		ChangeOptIncludeDocs(),
	)
	if err != nil {
		return nil, err
	}
	return domainEventsOfChangeEvents(ctx, ch), nil
}

func (c Connection) StartListener(
	ctx context.Context,
) (ch <-chan domain.Event, err error) {
	err1 := c.processNewDomainEvents(ctx)
	ch, err2 := c.processUnpublishedDomainEvents(ctx)
	err = errors.Join(err1, err2)
	return
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
		rev = resp.Header.Get("Etag")
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

func (c Connection) Get(id string, doc any) (rev string, err error) {
	var resp *http.Response
	if resp, err = http.Get(c.docURL(id)); err != nil {
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
		err = fmt.Errorf("couchdb: unexpected http status code: %d", resp.StatusCode)
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
		couchURL = "http://admin:password@localhost:5984/harmony"
	}
	conn, err := NewCouchConnection(couchURL)
	if err == nil {
		DefaultConnection = conn
	} else {
		slog.Error("couchdb: Error initializing", "err", err)
	}
}
