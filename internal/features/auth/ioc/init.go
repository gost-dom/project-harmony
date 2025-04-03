package ioc

import (
	"context"
	"harmony/internal/features/auth"
	"harmony/internal/features/auth/authrouter"

	"github.com/gost-dom/surgeon"
)

func Install[T any](graph *surgeon.Graph[T]) *surgeon.Graph[T] {
	graph = surgeon.Replace[authrouter.Registrator](graph, &auth.Registrator{})
	return graph
	//surgeon.Replace[auth.AccountRepository](graph, dummyAccRepo{})
}

type dummyAccRepo struct{}

func (r dummyAccRepo) Insert(context.Context, auth.AccountUseCaseResult) error {
	return nil
}
