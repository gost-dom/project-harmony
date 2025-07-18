package gosthttp

import (
	"context"
	"net/http"
)

// RewriterMiddleware injects a "rewriter" into the request context. See
// [Rewrite] for more information
func RewriterMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var rewriter http.Handler = http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				if a, b := r.Context().Value("rewritten").(bool); a && b {
					w.WriteHeader(500)
					return
				}
				r = r.WithContext(context.WithValue(r.Context(), "rewritten", true))
				h.ServeHTTP(w, r)
			})
		h.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), "rewriter", rewriter)))
	})
}

// Rewrite reprocesses the request as if it was a GET request to a different
// path; allowing e.g., a POST request to return the body.
//
// This is not recommended for normal POST requests; but useful for HTMX handled
// forms where the response body should be included with the response.
func Rewrite(w http.ResponseWriter, r *http.Request, path string, query string) {
	rewriter, ok := r.Context().Value("rewriter").(http.Handler)
	if !ok {
		w.WriteHeader(500)
		return
	}
	r.Method = "GET"
	r.URL.Path = path
	r.URL.RawQuery = query
	rewriter.ServeHTTP(w, r)
}
