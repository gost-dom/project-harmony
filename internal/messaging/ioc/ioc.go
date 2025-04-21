package ioc

import (
	"harmony/internal/core/corerepo"
	"harmony/internal/couchdb"
	"harmony/internal/features/auth/authrepo"
	"harmony/internal/messaging"

	"github.com/gost-dom/surgeon"
)

var Graph *surgeon.Graph[messaging.MessageHandler]

func init() {
	handler := messaging.NewMessageHandler()
	// handler := messaging.MessageHandler{
	// 	corerepo.DefaultDomainEventRepo,
	// 	auth.EmailValidator{Repository: authrepo.DefaultAccountRepository},
	// }
	Graph = surgeon.BuildGraph(*handler) // messaging.MessageHandler{})

	Graph.Inject(authrepo.AccountRepository{Connection: couchdb.DefaultConnection})
	Graph.Inject(corerepo.DefaultDomainEventRepo)
}

func Handler() messaging.MessageHandler { return Graph.Instance() }
