package auth

import "context"

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
