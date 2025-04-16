// Package couchtest contains helper code for interacting with a couch test db
package couchtest

import (
	"context"
	"fmt"
	"harmony/internal/couchdb"
	"testing"
)

// testWrapper is a helper function for code that is deisnged as test code can
// can potentially run outside the scope of a test. This allows the code to
// have errors reported to the test framework, but panic when executing outside.
type testWrapper struct{ t testing.TB }

func (t testWrapper) Error(args ...any) {
	if t.t == nil {
		panic(fmt.Sprintln(args...))
	}
	t.t.Error(args...)
}

func (t testWrapper) Errorf(format string, args ...any) {
	if t.t == nil {
		panic(fmt.Sprintf(format, args...))
	}
	t.t.Errorf(format, args...)
}

type row struct {
	ID    string `json:"id"`
	Value struct {
		Rev string `json:"rev"`
	} `json:"value"`
}

type allDocs struct {
	Offset    int   `json:"offset"`
	Rows      []row `json:"rows"`
	TotalRows int   `json:"total_rows"`
}

type DeleteDoc struct {
	ID      string `json:"_id"`
	Rev     string `json:"_rev"`
	Deleted bool   `json:"_deleted"`
}

type BulkDocs struct {
	Docs []DeleteDoc `json:"docs"`
}

type couchOption func(*CouchHelper)

func WithT(t testing.TB) couchOption { return func(c *CouchHelper) { c.t.t = t } }
func WithConnection(c couchdb.Connection) couchOption {
	return func(ch *CouchHelper) { ch.Connection = c }
}

type CouchHelper struct {
	t          testWrapper
	Connection couchdb.Connection
}

func NewCouchHelper(opts ...couchOption) CouchHelper {
	res := CouchHelper{Connection: couchdb.DefaultConnection}
	for _, o := range opts {
		o(&res)
	}
	return res
}

func (h CouchHelper) DeleteAllDocs() {
	conn := h.Connection
	var docs allDocs
	_, err := conn.Get("_all_docs", &docs)
	if err != nil {
		h.t.Errorf("couchdbtest: cannot initialize: %v", err)
	}
	rows := docs.Rows
	var deleteDoc BulkDocs
	deleteDoc.Docs = make([]DeleteDoc, len(rows))
	for i, d := range rows {
		deleteDoc.Docs[i] = DeleteDoc{ID: d.ID, Rev: d.Value.Rev,
			Deleted: true,
		}
	}
	resp, err := conn.RawPost(context.Background(), "_bulk_docs", deleteDoc)
	if err != nil {
		h.t.Errorf("couchdbtest: cannot delete: %v", err)
	}
	defer resp.Body.Close()
	// b, err := io.ReadAll(resp.Body)
	// if err != nil {
	// 	h.t.Errorf("couchdbtest: DeleteAllDocs: cannot read body: %v", err)
	// }
	switch resp.StatusCode {
	case 200, 201:
		return
	default:
		h.t.Errorf(
			"couchdbtest: DeleteAllDocs: unexpected status code: %d",
			resp.StatusCode)
	}
}
