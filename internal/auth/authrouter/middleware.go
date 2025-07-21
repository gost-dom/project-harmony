package authrouter

import (
	"fmt"
	serverctx "harmony/internal/web/server/ctx"
	"net/http"
	"net/url"
)

const (
	PathAuthLogin = "/auth/login"
)

// RequireAuth is a middleware that will only render the inner handler if the
// user has been authenticated. Otherwise, it sends the user to the login page.
// If the request is an HTMX request, the login page is sent in the response,
// otherwise, an HTTP redirect response is returned to the user.
func RequireAuth(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if serverctx.IsLoggedIn(r.Context()) {
			h.ServeHTTP(w, r)
			return
		}

		query := fmt.Sprintf("redirectUrl=%s", url.QueryEscape(r.URL.Path))
		newURL := fmt.Sprintf("%s?%s", PathAuthLogin, query)
		if r.Header.Get("HX-Request") == "" {
			http.Redirect(w, r, newURL, 303)
		} else {
			w.Header().Add("hx-replace-url", newURL)
			rewrite(w, r, "/auth/login", query)
		}
	})
}
