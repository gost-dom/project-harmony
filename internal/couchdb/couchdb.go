package couchdb

import (
	"errors"
	"fmt"
	"net/http"
)

// ErrConn indicates that an error occurred trying to communicate with CouchDB
// itself. Possible causes:
//   - Temporary condition such as a disconnected network
//   - Configuration issue, e.g., a wrong host name
var ErrConn = errors.New("couchdb: connection error")

var ErrConflict = errors.New("couchdb: conflict")
var ErrNotFound = errors.New("couchdb: not found")

type Document any

type couchDoc struct {
	ID  string `json:"_id"`
	Rev string `json:"_rev"`
	Document
}

func errUnexpectedStatusCode(resp *http.Response) error {
	return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
}
