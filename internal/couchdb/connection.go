package couchdb

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"maps"
	"net/http"
	"net/url"
)

type Connection struct{ dbURL *url.URL }

// Bootstrap creates the database, as well as updates any design documents, such
// as views. Panics on unrecoverable errors, e.g., an invalid configuration.
func (c Connection) Bootstrap(ctx context.Context) error {
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

func NewCouchConnection(couchURL string) (conn Connection, err error) {
	var url *url.URL
	url, err = url.Parse(couchURL)
	conn = Connection{url}
	if err == nil {
		err = conn.Bootstrap(context.Background())
	}
	return
}

// docURL generates the full couchDB resource URL for a given document ID. The
// full URL will include both database name and credentials.
func (c Connection) docURL(id string) string {
	return c.dbURL.JoinPath(url.PathEscape(id)).String()
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
	// req, err := http.NewRequestWithContext(ctx, "PUT", url, &b)
	// if err != nil {
	// 	err = fmt.Errorf("%w: put failed: %v", ErrConn, err)
	// 	return
	// }
	// resp, err := c.do(req)
	// if err != nil {
	// 	return
	// }
	var resp *http.Response
	resp, err = c.req(ctx, "PUT", url, nil, &b)
	defer resp.Body.Close()

	switch resp.StatusCode {
	case 201:
		rev = resp.Header.Get("Etag")
	case 409:
		err = ErrConflict
	default:
		err = fmt.Errorf("couch: insert id(%s): %w", id, errUnexpectedStatusCode(resp))
		return
	}
	return
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
		err = ErrNotFound
	default:
		err = fmt.Errorf("couchdb: insert id(%s): %w", id, errUnexpectedStatusCode(resp))
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
