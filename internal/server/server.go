package server

import (
	"harmony/internal/project"
	"harmony/internal/server/views"
	"net/http"
	"path/filepath"

	"github.com/a-h/templ"
)

func noCache(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// fmt.Println("Request", r.Method, r.URL.Path)
		w.Header().Add("cache-control", "no-cache")
		h.ServeHTTP(w, r)
	})
}

func staticFilesPath() string { return filepath.Join(project.Root(), "static") }

func New() http.Handler {
	component := views.Index()
	login := views.AuthLogin()
	loggedIn := false

	mux := http.NewServeMux()
	mux.Handle("GET /{$}", templ.Handler(component))
	mux.Handle("GET /auth/login/{$}", templ.Handler(login))
	mux.HandleFunc("POST /auth/login", func(w http.ResponseWriter, r *http.Request) {
		loggedIn = true
		w.Header().Add("hx-push-url", "/host")
	})
	mux.HandleFunc("GET /host/{$}", func(w http.ResponseWriter, r *http.Request) {
		if !loggedIn {
			w.Header().Add("hx-push-url", "/auth/login")
			login.Render(r.Context(), w)
		} else {
			views.HostsPage().Render(r.Context(), w)
		}
	})
	mux.Handle(
		"GET /static/",
		http.StripPrefix("/static", http.FileServer(
			http.Dir(staticFilesPath()))),
	)
	return noCache(mux)
}
