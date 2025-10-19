package router

import (
	"harmony/internal/auth"
	"harmony/internal/web"
	"net/http"
)

// Middlewares contain middlewares necessary for the authentication flow.
//
// - Checking session cookies, adding an AuthenticatedAccount to context, if applicable.
// - Rewriter; used internally to "rewrite" responses, e.g., for HTMX auth flow
type Middlewares struct {
	SessionManager SessionManager
}

// SessionAuth retrieves the logged in user from the session and
// writes it to the request context.
func (s Middlewares) SessionAuth(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if account := s.SessionManager.LoggedInUser(r); account != nil {
			auth.SetAuthenticatedUser(&r, *account)
		}
		h.ServeHTTP(w, r)
	})
}

func (m Middlewares) Apply(h http.Handler) http.Handler {
	return web.ApplyMiddlewares(h,
		web.MiddlewareFunc(RewriterMiddleware),
		web.MiddlewareFunc(m.SessionAuth),
	)
}
