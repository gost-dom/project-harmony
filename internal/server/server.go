package server

import (
	"context"
	"fmt"

	"harmony/internal/features/auth"
	"harmony/internal/project"
	"harmony/internal/server/views"
	"net/http"
	"net/url"
	"path/filepath"

	"github.com/a-h/templ"
	"github.com/gorilla/sessions"
)

const sessionCookieName = "accountId"

func noCache(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("cache-control", "no-cache")
		h.ServeHTTP(w, r)
	})
}

func staticFilesPath() string { return filepath.Join(project.Root(), "static") }

type Authenticator interface {
	Authenticate(context.Context, string, string) (auth.Account, error)
}

type SessionManager struct {
	SessionStore sessions.Store
}

func (m *SessionManager) LoggedInUser(r *http.Request) *auth.Account {
	session, _ := m.SessionStore.Get(r, sessionNameAuth)
	if id, ok := session.Values[sessionCookieName]; ok {
		result := new(auth.Account)
		if strId, ok := id.(string); ok {
			result.Id = auth.AccountId(strId)
			return result
		}
	}
	return nil
}

type Server struct {
	http.Handler
	SessionManager SessionManager
	AuthRouter     *AuthRouter
}

type sessionName string

const (
	sessionNameAuth = "auth"
)

func (s *Server) GetHost(w http.ResponseWriter, r *http.Request) {
	if account := s.SessionManager.LoggedInUser(r); account != nil {
		views.HostsPage().Render(r.Context(), w)
		return
	}
	// Not authenticated; show login page
	fmtNewLocation := fmt.Sprintf("/auth/login?redirectUrl=%s", url.QueryEscape("/hosts"))
	w.Header().Add("hx-push-url", fmtNewLocation)
	views.AuthLogin("/host", views.LoginFormData{}).Render(r.Context(), w)
}

type AuthRouter struct {
	*http.ServeMux
	Authenticator Authenticator
	SessionStore  sessions.Store
}

func (s *AuthRouter) PostAuthLogin(w http.ResponseWriter, r *http.Request) {
	// Get a session. We're ignoring the error resulted from decoding an
	// existing session: Get() always returns a session, even if empty.

	session, _ := s.SessionStore.Get(r, sessionNameAuth)
	r.ParseForm()
	email := r.FormValue("email")
	password := r.FormValue("password")
	redirectUrl := r.FormValue("redirectUrl")
	if redirectUrl == "" {
		redirectUrl = "/"
	}
	if account, err := s.Authenticator.Authenticate(r.Context(), email, password); err == nil {
		session.Values[sessionCookieName] = string(account.Id)
		// TODO: Handle err
		session.Save(r, w)
		w.Header().Add("hx-push-url", redirectUrl)
	} else {
		data := views.LoginFormData{
			Email:              email,
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

func (r *AuthRouter) Init() {
	r.ServeMux = http.NewServeMux()
	r.HandleFunc("GET /login", func(w http.ResponseWriter, r *http.Request) {
		redirectUrl := r.URL.Query().Get("redirectUrl")
		views.AuthLogin(redirectUrl, views.LoginFormData{}).Render(r.Context(), w)

	})
	r.HandleFunc("POST /login", r.PostAuthLogin)
}

func (s *Server) Init() {
	mux := http.NewServeMux()
	mux.Handle("/auth/", http.StripPrefix("/auth", s.AuthRouter))
	mux.Handle("GET /{$}", templ.Handler(views.Index()))
	mux.HandleFunc("GET /host/{$}", s.GetHost)
	mux.Handle(
		"GET /static/",
		http.StripPrefix("/static", http.FileServer(
			http.Dir(staticFilesPath()))),
	)
	s.Handler = noCache(mux)
}

func NewAuthRouter() *AuthRouter {
	res := &AuthRouter{}
	res.Init()
	return res
}

func New() *Server {
	res := &Server{
		AuthRouter: NewAuthRouter(),
	}
	res.Init()
	return res
}
