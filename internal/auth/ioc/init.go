package ioc

import (
	"harmony/internal/auth"
	"harmony/internal/auth/authrepo"
	"harmony/internal/auth/authrouter"
	"harmony/internal/auth/sessionstore"
	"harmony/internal/core/corerepo"
	"os"

	"github.com/gost-dom/surgeon"
)

func Install[T any](graph *surgeon.Graph[T]) *surgeon.Graph[T] {
	graph = surgeon.Replace[authrouter.Registrator](graph, &auth.Registrator{})
	graph = surgeon.Replace[authrouter.Authenticator](graph, &auth.Authenticator{})
	graph = surgeon.Replace[authrouter.EmailValidator](graph, &auth.EmailChallengeValidator{})

	authKey := os.Getenv("SESSION_AUTH_KEY")
	encKey := os.Getenv("SESSION_ENC_KEY")
	if authKey == "" && encKey == "" {
		// Fallback values for development.
		authKey = "authkey1234"
		encKey = "enckey12341234567890123456789012"
	}

	graph.Inject(sessionstore.NewCouchDBStore(
		&corerepo.DefaultConnection,
		[][]byte{
			[]byte(authKey),
			[]byte(encKey),
		},
	))
	repo := &authrepo.AccountRepository{
		Connection: corerepo.DefaultConnection,
	}
	graph = surgeon.ReplaceAll(graph, repo)
	return graph
}
