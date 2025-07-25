package ioc

import (
	authioc "harmony/internal/auth/ioc"
	"harmony/internal/core/corerepo"
	"harmony/internal/messaging"
	mioc "harmony/internal/messaging/ioc"
	"harmony/internal/web/server"

	"github.com/gost-dom/surgeon"
)

type RootGraph struct {
	Server      *server.Server
	MessagePump messaging.MessagePump
}

var Graph *surgeon.Graph[RootGraph]

func init() {
	Graph = surgeon.BuildGraph(RootGraph{
		new(server.Server),
		messaging.MessagePump{
			MessageSource:         corerepo.DefaultMessageSource,
			DomainEventRepository: corerepo.DefaultDomainEventRepo,
			Handler:               mioc.Handler(),
		},
	})
	Graph = authioc.Install(Graph)
	if err := Graph.Validate(); err != nil {
		panic(err)
	}
}

func Root() RootGraph { return Graph.Instance() }
