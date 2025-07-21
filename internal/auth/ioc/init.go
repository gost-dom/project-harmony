package ioc

import (
	"harmony/internal/auth"
	"harmony/internal/auth/repo"
	"harmony/internal/auth/router"
	"harmony/internal/auth/sessionstore"
	"harmony/internal/core/corerepo"
	"os"

	"github.com/gost-dom/surgeon"
)

func Install[T any](graph *surgeon.Graph[T]) *surgeon.Graph[T] {
	graph = surgeon.Replace[router.Registrator](graph, &auth.Registrator{})
	graph = surgeon.Replace[router.Authenticator](graph, &auth.Authenticator{})
	graph = surgeon.Replace[router.EmailValidator](graph, &auth.EmailChallengeValidator{})

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
	repo := &repo.AccountRepository{
		Connection: corerepo.DefaultConnection,
	}
	graph = surgeon.ReplaceAll(graph, repo)
	return graph
}
