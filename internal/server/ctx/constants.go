package serverctx

import (
	"context"
	"harmony/internal/features/auth/authdomain"
	"net/http"
)

type ContextKey string

const (
	AuthAccount ContextKey = "auth:account"

	ServerRewritten    ContextKey = "server:rewritten"
	ServerRewriter     ContextKey = "server:rewriter"
	ServerCSRFTokenSrc ContextKey = "server:csrf:token-source"
	ServerReqID        ContextKey = "server:"
)

func IsLoggedIn(c context.Context) bool {
	acc := c.Value(AuthAccount)
	return acc != nil
}

func SetUser(r **http.Request, acc *authdomain.Account) {
	SetReqValue(r, AuthAccount, acc)
}

func GetUser(c context.Context) *authdomain.Account {
	acc, _ := c.Value(AuthAccount).(*authdomain.Account)
	return acc
}

func ReqValue[T any](r *http.Request, key ContextKey) (res T, ok bool) {
	v := r.Context().Value(key)
	res, ok = v.(T)
	return
}

func SetReqValue(r **http.Request, key ContextKey, v any) {
	ctx := context.WithValue((*r).Context(), key, v)
	*r = (*r).WithContext(ctx)
}
