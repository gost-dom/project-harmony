package corerepo

import (
	"errors"
	"fmt"
	"harmony/internal/core"
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
var ErrNotFound = fmt.Errorf("couchdb: %w", core.ErrNotFound)

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
	ID    string `json:"id"`
	Key   string `json:"key"`
	Rev   string `json:"rev"`
	Value T      `json:"value"`
}

type DocViewRow[T any] struct {
	ID  string `json:"id"`
	Key string `json:"key"`
	Rev string `json:"rev"`
	Doc T      `json:"doc"`
}

type DocsViewResult[T any] struct {
	Offset    int             `json:"offset"`
	Rows      []DocViewRow[T] `json:"rows"`
	TotalRows int             `json:"total_rows"`
}

// Values return a []T containing the Value fields of each row.
func (r DocsViewResult[T]) Docs() []T {
	res := make([]T, len(r.Rows))
	for i, r := range r.Rows {
		res[i] = r.Doc
	}
	return res
}

type ViewResult[T any] struct {
	Offset    int          `json:"offset"`
	Rows      []ViewRow[T] `json:"rows"`
	TotalRows int          `json:"total_rows"`
}

// Values return a []T containing the Value fields of each row.
func (r ViewResult[T]) Values() []T {
	res := make([]T, len(r.Rows))
	for i, r := range r.Rows {
		res[i] = r.Value
	}
	return res
}
