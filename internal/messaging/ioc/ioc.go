package ioc

import (
	"harmony/internal/couchdb"
	"harmony/internal/features/auth/authrepo"
	"harmony/internal/messaging"

	"github.com/gost-dom/surgeon"
)

var Graph *surgeon.Graph[messaging.MessageHandler]

// init initializes the dependency injection graph for MessageHandler with required dependencies.
func init() {
	Graph = surgeon.BuildGraph(messaging.MessageHandler{})

	Graph.Inject(authrepo.AccountRepository{
		Connection: couchdb.DefaultConnection,
	})
}

// Handler returns a fully initialized MessageHandler instance from the dependency injection graph.
func Handler() messaging.MessageHandler { return Graph.Instance() }
