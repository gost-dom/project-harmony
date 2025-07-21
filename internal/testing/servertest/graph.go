package servertest

import (
	"harmony/internal/web/server"
	"harmony/internal/web/server/ioc"

	"github.com/gorilla/sessions"
	"github.com/gost-dom/surgeon"
	"github.com/quasoft/memstore"
)

func init() {
	Graph = ioc.Graph
	Graph = surgeon.Replace[sessions.Store](
		Graph, memstore.NewMemStore(
			[]byte("authkey123"),
			[]byte("enckey12341234567890123456789012"),
		))
}

type ServerGraph = *surgeon.Graph[*server.Server]

// Graph is a base dependency graph for testing. The SessionStore has been
// permanently replaced with an in-memory session store.
//
// Note, sessions are kept in memory as all tests would use the same store, but
// that shouldn't be a problem
var Graph *surgeon.Graph[*server.Server]
