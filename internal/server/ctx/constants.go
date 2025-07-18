package serverctx

import (
	"context"
	"harmony/internal/features/auth/authdomain"
)

type ContextKey string

const (
	AuthAccount ContextKey = "auth:account"
)

func IsLoggedIn(c context.Context) bool {
	acc := c.Value(AuthAccount)
	return acc != nil
}

func SetUser(c context.Context, acc *authdomain.Account) context.Context {
	return context.WithValue(c, AuthAccount, acc)
}

func GetUser(c context.Context) *authdomain.Account {
	acc, _ := c.Value(AuthAccount).(*authdomain.Account)
	return acc
}
