package ioc

import (
	"harmony/internal/couchdb"
	"harmony/internal/features/auth"
	"harmony/internal/features/auth/authrepo"
	"harmony/internal/features/auth/authrouter"

	"github.com/gost-dom/surgeon"
)

func Install[T any](graph *surgeon.Graph[T]) *surgeon.Graph[T] {
	repo := &authrepo.AccountRepository{
		Connection: couchdb.DefaultConnection,
	}
	graph = surgeon.Replace[authrouter.Registrator](graph, &auth.Registrator{})
	graph = surgeon.Replace[auth.AccountInserter](graph, repo)
	graph = surgeon.Replace[authrouter.Authenticator](graph, &auth.Authenticator{})
	graph = surgeon.Replace[authrouter.EmailValidator](
		graph,
		&auth.EmailChallengeValidator{Repository: repo},
	)
	graph = surgeon.Replace[auth.AccountEmailFinder](graph, repo)
	return graph
}
