package authrouter

import (
	"harmony/internal/auth"
	"harmony/internal/web"
	"log/slog"
	"net/http"
)

// RewriterMiddleware injects a "rewriter" into the request context. See
// [rewrite] for more information
func RewriterMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var rewriter http.Handler = http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				if a, b := r.Context().Value("rewritten").(bool); a && b {
					w.WriteHeader(500)
					return
				}
				web.SetReqValue(&r, auth.CtxKeyRewritten, true)
				h.ServeHTTP(w, r)
			})
		web.SetReqValue(&r, auth.CtxKeyRewriter, rewriter)
		h.ServeHTTP(w, r)
	})
}

// rewrite reprocesses the request as if it was a GET request to a different
// path; allowing e.g., a POST request to return the body.
//
// This is not recommended for normal POST requests; but useful for HTMX handled
// forms where the response body should be included with the response.
func rewrite(w http.ResponseWriter, r *http.Request, path string, query string) {
	slog.DebugContext(r.Context(), "Rewrite URL", "path", path)
	rewriter, ok := web.ReqValue[http.Handler](r, auth.CtxKeyRewriter)
	if !ok {
		w.WriteHeader(500)
		return
	}
	r.Method = "GET"
	r.URL.Path = path
	r.URL.RawQuery = query
	rewriter.ServeHTTP(w, r)
}
