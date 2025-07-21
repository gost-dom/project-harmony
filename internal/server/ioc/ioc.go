package ioc

import (
	"harmony/internal/core/corerepo"
	authioc "harmony/internal/features/auth/ioc"
	"harmony/internal/server"
	"harmony/internal/server/sessionstore"

	"github.com/gost-dom/surgeon"
)

var Graph *surgeon.Graph[*server.Server]

func init() {
	Graph = surgeon.BuildGraph(server.New())
	Graph.Inject(sessionstore.NewCouchDBStore(
		&corerepo.DefaultConnection,
		[][]byte{
			[]byte("authkey123"),
			[]byte("enckey12341234567890123456789012"),
		},
	))

	Graph = authioc.Install(Graph)
	if err := Graph.Validate(); err != nil {
		panic(err)
	}
}

func Server() *server.Server { return Graph.Instance() }
