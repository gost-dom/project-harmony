package ioc

import (
	"harmony/internal/core/corerepo"
	"harmony/internal/auth"
	"harmony/internal/auth/authrepo"
	"harmony/internal/auth/authrouter"

	"github.com/gost-dom/surgeon"
)

func Install[T any](graph *surgeon.Graph[T]) *surgeon.Graph[T] {
	graph = surgeon.Replace[authrouter.Registrator](graph, &auth.Registrator{})
	graph = surgeon.Replace[authrouter.Authenticator](graph, &auth.Authenticator{})
	graph = surgeon.Replace[authrouter.EmailValidator](graph, &auth.EmailChallengeValidator{})

	repo := &authrepo.AccountRepository{
		Connection: corerepo.DefaultConnection,
	}
	graph = surgeon.ReplaceAll(graph, repo)
	return graph
}
