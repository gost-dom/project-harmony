package ioc

import (
	authioc "harmony/internal/auth/ioc"
	"harmony/internal/web/server"

	"github.com/gost-dom/surgeon"
)

var Graph *surgeon.Graph[*server.Server]

func init() {
	Graph = surgeon.BuildGraph(server.New())
	Graph = authioc.Install(Graph)
	if err := Graph.Validate(); err != nil {
		panic(err)
	}
}

func Server() *server.Server { return Graph.Instance() }
