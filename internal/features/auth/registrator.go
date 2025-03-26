package auth

import "context"

type Entity[T any] interface{ ID() T }

type AccountRegistered struct {
	AccountID
}

type AccountUseCaseResult = UseCaseResult[Account, AccountID]

type AccountRepository interface {
	Insert(context.Context, AccountUseCaseResult) error
}
