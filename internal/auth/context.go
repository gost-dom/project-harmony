package auth

import (
	"context"
	"harmony/internal/auth/domain"
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
func AuthenticatedUser(ctx context.Context) (acc domain.AuthenticatedAccount, ok bool) {
	ctxVal := ctx.Value(CtxKeyAuthAccount)
	if ok = (ctxVal != nil); !ok {
		return
	}

	acc, ok = ctxVal.(domain.AuthenticatedAccount)
	return
}

// WithContext is the interface for a context-bearing value, where a new value
// with a new context can be created using a WithContext function.
//
// In reality, this represents a [*net/http.Request], but callers shouldn't be
// coupled to the reques object.
type Contexter[T any] interface {
	Context() context.Context
	WithContext(context.Context) T
}

// SetAuthenticatedUser storeds an authenticated user in the request context.
func SetAuthenticatedUser[T Contexter[T]](r *T, acc domain.AuthenticatedAccount) {
	ctx := context.WithValue((*r).Context(), CtxKeyAuthAccount, acc)
	*r = (*r).WithContext(ctx)
}
