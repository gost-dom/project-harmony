package ioc

import (
	"harmony/internal/server"
	"harmony/internal/server/authenticator"

	"github.com/gost-dom/surgeon"
	"github.com/quasoft/memstore"
)

var Graph *surgeon.Graph[*server.Server]

func init() {
	Graph = surgeon.BuildGraph(server.New())
	Graph.Inject(memstore.NewMemStore(
		[]byte("authkey123"),
		[]byte("enckey12341234567890123456789012"),
	))
	Graph.Inject(authenticator.New())
}

func Server() *server.Server { return Graph.Instance() }
