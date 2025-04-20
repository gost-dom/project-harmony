package ioc

import (
	"harmony/internal/couchdb"
	"harmony/internal/features/auth"
	"harmony/internal/features/auth/authrepo"
	"harmony/internal/features/auth/authrouter"

	"github.com/gost-dom/surgeon"
)

// Install configures the dependency injection graph with default authentication and account repository implementations.
// 
// It replaces the authrouter.Registrator dependency with a new auth.Registrator instance, and the auth.AccountRepository dependency with an authrepo.AccountRepository using the default CouchDB connection.
// Returns the updated dependency injection graph.
func Install[T any](graph *surgeon.Graph[T]) *surgeon.Graph[T] {
	graph = surgeon.Replace[authrouter.Registrator](graph, &auth.Registrator{})
	graph = surgeon.Replace[auth.AccountRepository](graph, &authrepo.AccountRepository{
		Connection: couchdb.DefaultConnection,
	})
	return graph
}
