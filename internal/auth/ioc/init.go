package ioc

import (
	"harmony/internal/auth"
	"harmony/internal/auth/authrepo"
	"harmony/internal/auth/authrouter"
	"harmony/internal/auth/sessionstore"
	"harmony/internal/core/corerepo"

	"github.com/gost-dom/surgeon"
)

func Install[T any](graph *surgeon.Graph[T]) *surgeon.Graph[T] {
	graph = surgeon.Replace[authrouter.Registrator](graph, &auth.Registrator{})
	graph = surgeon.Replace[authrouter.Authenticator](graph, &auth.Authenticator{})
	graph = surgeon.Replace[authrouter.EmailValidator](graph, &auth.EmailChallengeValidator{})

	graph.Inject(sessionstore.NewCouchDBStore(
		&corerepo.DefaultConnection,
		[][]byte{
			[]byte("authkey123"),
			[]byte("enckey12341234567890123456789012"),
		},
	))
	repo := &authrepo.AccountRepository{
		Connection: corerepo.DefaultConnection,
	}
	graph = surgeon.ReplaceAll(graph, repo)
	return graph
}
