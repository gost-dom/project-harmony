package serverctx

import (
	"context"
	"harmony/internal/features/auth/authdomain"
)

const (
	AuthAccount = "auth:account"
)

func IsLoggedIn(c context.Context) bool {
	acc := c.Value(AuthAccount)
	return acc != nil
}

func SetUser(c context.Context, acc *authdomain.Account) context.Context {
	return context.WithValue(c, AuthAccount, acc)
}
