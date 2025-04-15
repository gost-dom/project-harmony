package ioc

import (
	"harmony/internal/features/auth"
	"harmony/internal/features/auth/authrouter"

	"github.com/gost-dom/surgeon"
)

func Install[T any](graph *surgeon.Graph[T]) *surgeon.Graph[T] {
	graph = surgeon.Replace[authrouter.Registrator](graph, &auth.Registrator{})
	return graph
}
