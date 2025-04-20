package ioc

import (
	"harmony/internal/couchdb"
	"harmony/internal/features/auth/authrepo"
	"harmony/internal/messaging"

	"github.com/gost-dom/surgeon"
)

var Graph *surgeon.Graph[messaging.MessageHandler]

func init() {
	Graph = surgeon.BuildGraph(messaging.MessageHandler{})

	Graph.Inject(authrepo.AccountRepository{
		Connection: couchdb.DefaultConnection,
	})
}

func Handler() messaging.MessageHandler { return Graph.Instance() }
