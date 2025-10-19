package auth

import (
	"context"
	"harmony/internal/auth/domain"
	"harmony/internal/web"
	"net/http"
)

// Auth-related context values

type contextKey string

const (
	CtxKeyRewritten   contextKey = "rewritten"
	CtxKeyRewriter    contextKey = "rewriter"
	CtxKeyAuthAccount contextKey = "account"
)

// UserAuthenticated returns whether we are processing a request from an
// authenticated user.
func UserAuthenticated(c context.Context) (res bool) {
	_, res = AuthenticatedUser(c)
	return
}

// AuthenticatedUser returns the currently authenticated account. If no user is
// authenticated, acc will be a zero Account, and ok will be false.
func AuthenticatedUser(ctx context.Context) (acc domain.Account, ok bool) {
	ctxVal := ctx.Value(CtxKeyAuthAccount)
	if ok = (ctxVal != nil); !ok {
		return
	}

	acc, ok = ctxVal.(domain.Account)
	return
}

// SetAuthenticatedUser storeds an authenticated user in the request context.
func SetAuthenticatedUser(r **http.Request, acc domain.Account) {
	web.SetReqValue(r, CtxKeyAuthAccount, acc)
}
