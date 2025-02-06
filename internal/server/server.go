package server

import (
	"context"
	"fmt"
	"harmony/internal/project"
	"harmony/internal/server/views"
	"net/http"
	"net/url"
	"path/filepath"

	"github.com/a-h/templ"
)

func noCache(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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

func (s *Server) GetHost(w http.ResponseWriter, r *http.Request) {
	if !s.loggedIn {
		fmtNewLocation := fmt.Sprintf("/auth/login?redirectUrl=%s", url.QueryEscape("/hosts"))
		w.Header().Add("hx-push-url", fmtNewLocation)
		views.AuthLogin("/host").Render(r.Context(), w)
	} else {
		views.HostsPage().Render(r.Context(), w)
	}
}

func (s *Server) PostAuthLogin(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	email := r.FormValue("email")
	password := r.FormValue("password")
	redirectUrl := r.FormValue("redirectUrl")
	if redirectUrl == "" {
		redirectUrl = "/"
	}
	if email == "valid-user@example.com" && password == "s3cret" {
		s.loggedIn = true
		w.Header().Add("hx-push-url", redirectUrl)
	} else {
		data := views.LoginFormData{
			Email:              "",
			Password:           "",
			InvalidCredentials: true,
		}
		if r.FormValue("email") == "" {
			data.EmailMissing = true
		}
		if r.FormValue("password") == "" {
			data.PasswordMissing = true
		}
		views.LoginForm(redirectUrl, data).Render(r.Context(), w)
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
	mux.HandleFunc("GET /auth/login/{$}", func(w http.ResponseWriter, r *http.Request) {
		redirectUrl := r.URL.Query().Get("redirectUrl")
		views.AuthLogin(redirectUrl).Render(r.Context(), w)

	})
	mux.HandleFunc("POST /auth/login", server.PostAuthLogin)
	mux.HandleFunc("GET /host/{$}", server.GetHost)
	mux.Handle(
		"GET /static/",
		http.StripPrefix("/static", http.FileServer(
			http.Dir(staticFilesPath()))),
	)
	return server
}
