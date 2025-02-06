package server

import (
	"context"
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

type Account struct{}

type Authenticator interface {
	Authenticate(context.Context, string, string) *Account
}

type Server struct {
	http.Handler
	loggedIn bool
}

var login = views.AuthLogin()

func (s *Server) GetHost(w http.ResponseWriter, r *http.Request) {
	if !s.loggedIn {
		w.Header().Add("hx-push-url", "/auth/login")
		login.Render(r.Context(), w)
	} else {
		views.HostsPage().Render(r.Context(), w)
	}
}

func (s *Server) PostAuthLogin(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	if r.FormValue("email") == "valid-user@example.com" {
		s.loggedIn = true
		w.Header().Add("hx-push-url", "/host")
	} else {
		views.LoginForm("", views.LoginFormData{
			Email:              "",
			Password:           "",
			InvalidCredentials: true,
		}).Render(r.Context(), w)
	}
}

func New() http.Handler {
	component := views.Index()

	mux := http.NewServeMux()
	server := &Server{
		noCache(mux),
		false,
	}
	mux.Handle("GET /{$}", templ.Handler(component))
	mux.Handle("GET /auth/login/{$}", templ.Handler(login))
	mux.HandleFunc("POST /auth/login", server.PostAuthLogin)
	mux.HandleFunc("GET /host/{$}", server.GetHost)
	mux.Handle(
		"GET /static/",
		http.StripPrefix("/static", http.FileServer(
			http.Dir(staticFilesPath()))),
	)
	return server
}
