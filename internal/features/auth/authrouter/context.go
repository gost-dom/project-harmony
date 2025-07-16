// Contains context-aware operations, like store the logged in user.

package authrouter

import (
	"context"
	"harmony/internal/features/auth/authdomain"
	serverctx "harmony/internal/server/ctx"
)

func setAuth(ctx context.Context, acc authdomain.AuthenticatedAccount) context.Context {
	return serverctx.SetUser(ctx, acc.Account)
}
