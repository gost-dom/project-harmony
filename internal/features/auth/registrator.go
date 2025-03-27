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
type RegistratorInput struct {
	Email    string
	Password string
}

type Registrator struct {
	Repository AccountRepository
}

func (r Registrator) Register(ctx context.Context, input RegistratorInput) error {
	id, err := NewID()
	if err != nil {
		return err
	}
	account := Account{
		Id:    AccountID(id),
		Email: input.Email,
	}
	res := NewResult(account)
	res.AddEvent(AccountRegistered{AccountID: account.Id})
	return r.Repository.Insert(ctx, *res)
}
