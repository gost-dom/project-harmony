package main

import (
	"meet-the-locals/views"
	"net/http"

	"github.com/a-h/templ"
)

func noCache(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("cache-control", "no-cache")
		h.ServeHTTP(w, r)
	})
}

func main() {
	component := views.Index()
	login := views.AuthLogin()

	http.DefaultServeMux.Handle("GET /{$}", templ.Handler(component))
	http.DefaultServeMux.Handle("GET /login/{$}", templ.Handler(login))
	http.DefaultServeMux.Handle(
		"/static/",
		http.StripPrefix("/static", http.FileServer(http.Dir("static"))),
	)

	http.ListenAndServe("0.0.0.0:8081", noCache(http.DefaultServeMux))
}
