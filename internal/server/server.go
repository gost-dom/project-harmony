package server

import (
	"fmt"
	"harmony/internal/project"
	"harmony/internal/server/views"
	"net/http"
	"path/filepath"

	"github.com/a-h/templ"
)

func noCache(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Request", r.URL.Path)
		w.Header().Add("cache-control", "no-cache")
		h.ServeHTTP(w, r)
	})
}

func staticFilesPath() string {
	root := project.Root()

	return filepath.Join(root, "static")

}

func New() http.Handler {
	component := views.Index()
	login := views.AuthLogin()

	mux := http.NewServeMux()
	mux.Handle("GET /{$}", templ.Handler(component))
	mux.Handle("GET /auth/login/{$}", templ.Handler(login))
	mux.HandleFunc("GET /host/{$}", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("hx-push-url", "/auth/login")
		w.WriteHeader(200)
	})
	mux.Handle(
		"GET /static/",
		http.StripPrefix("/static", http.FileServer(
			http.Dir(staticFilesPath()))),
	)
	return noCache(mux)
}
