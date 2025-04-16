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
//   - Bad request - this shouldn't happen at runtime, as it's a code bug ...
var ErrConn = errors.New("couchdb: connection error")

// ErrRequest indicates that an error occurred trying to build an HTTP request.
// This shouldn't occur at runtime - it could only be the result of a bug in
// code.
var ErrRequest = errors.New("request error")

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

type ViewRow[T any] struct {
	Key   string `json:"id"`
	Rev   string `json:"rev"`
	Value T      `json:"value"`
}

type ViewResult[T any] struct {
	Offset    int          `json:"offset"`
	Rows      []ViewRow[T] `json:"rows"`
	TotalRows int          `json:"total_rows"`
}
