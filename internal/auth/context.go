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
func UserAuthenticated(c context.Context) bool {
	acc := c.Value(CtxKeyAuthAccount)
	return acc != nil
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
