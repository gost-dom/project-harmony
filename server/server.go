package server

import (
	"harmony/views"
	"net/http"

	"github.com/a-h/templ"
)

func noCache(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("cache-control", "no-cache")
		h.ServeHTTP(w, r)
	})
}

func New() http.Handler {
	component := views.Index()
	login := views.AuthLogin()

	mux := http.NewServeMux()
	mux.Handle("GET /{$}", templ.Handler(component))
	mux.Handle("GET /auth/login/{$}", templ.Handler(login))
	mux.Handle(
		"/static/",
		http.StripPrefix("/static", http.FileServer(http.Dir("static"))),
	)
	return noCache(mux)
}
