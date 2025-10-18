package servertest

import (
	"harmony/cmd/server/ioc"
	"harmony/internal/web/server"

	"github.com/gorilla/sessions"
	"github.com/gost-dom/surgeon"
	"github.com/quasoft/memstore"
)

func NewMemStore() sessions.Store {
	return memstore.NewMemStore(
		[]byte("authkey123"),
		[]byte("enckey12341234567890123456789012"),
	)
}

func init() {
	root := ioc.Graph.Instance()
	Graph = surgeon.BuildGraph(root.Server)
	Graph = surgeon.Replace(Graph, NewMemStore())
}

type ServerGraph = *surgeon.Graph[*server.Server]

// Graph is a base dependency graph for testing. The SessionStore has been
// permanently replaced with an in-memory session store.
//
// Note, sessions are kept in memory as all tests would use the same store, but
// that shouldn't be a problem
var Graph *surgeon.Graph[*server.Server]
