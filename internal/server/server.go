package server

import (
	"context"
	"encoding/gob"
	"errors"
	"fmt"
	"harmony/internal/project"
	"harmony/internal/server/views"
	"net/http"
	"net/url"
	"path/filepath"

	"github.com/a-h/templ"
	"github.com/gorilla/sessions"
	"github.com/quasoft/memstore"
)

func noCache(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("cache-control", "no-cache")
		h.ServeHTTP(w, r)
	})
}

func staticFilesPath() string { return filepath.Join(project.Root(), "static") }

type AccountId string

type Account struct{ Id AccountId }

type authenticator struct{}

func (a authenticator) Authenticate(
	ctx context.Context,
	username string,
	password string,
) (account Account, err error) {
	if username == "valid-user@example.com" && password == "s3cret" {
		account = Account{}
	} else {
		err = ErrBadCredentials
	}
	return
}

type Authenticator interface {
	Authenticate(context.Context, string, string) (Account, error)
}

var ErrBadCredentials = errors.New("authenticate: Bad credentials")

type SessionManager struct {
	sessionStore sessions.Store
}

type Server struct {
	http.Handler
	SessionManager SessionManager
	AuthRouter     *AuthRouter
	sessionStore   sessions.Store
}

type sessionName string

const (
	sessionNameAuth = "auth"
)

func init() {
	gob.Register(AccountId(""))
}

func (s *Server) GetHost(w http.ResponseWriter, r *http.Request) {
	session, _ := s.sessionStore.Get(r, sessionNameAuth)
	if _, ok := session.Values["accountId"]; !ok {
		fmtNewLocation := fmt.Sprintf("/auth/login?redirectUrl=%s", url.QueryEscape("/hosts"))
		w.Header().Add("hx-push-url", fmtNewLocation)
		views.AuthLogin("/host").Render(r.Context(), w)
	} else {
		views.HostsPage().Render(r.Context(), w)
	}
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
		session.Values["accountId"] = account.Id
		// TODO: Handle err
		session.Save(r, w)
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

func NewAuthRouter(store sessions.Store, auth Authenticator) *AuthRouter {
	result := &AuthRouter{
		ServeMux:      http.NewServeMux(),
		Authenticator: auth,
		SessionStore:  store,
	}
	result.HandleFunc("GET /login", func(w http.ResponseWriter, r *http.Request) {
		redirectUrl := r.URL.Query().Get("redirectUrl")
		views.AuthLogin(redirectUrl).Render(r.Context(), w)

	})
	result.HandleFunc("POST /login", result.PostAuthLogin)
	return result
}

func New() *Server {
	component := views.Index()

	sessionStore := memstore.NewMemStore(
		[]byte("authkey123"),
		[]byte("enckey12341234567890123456789012"),
	)
	mux := http.NewServeMux()
	server := &Server{
		AuthRouter:     NewAuthRouter(sessionStore, authenticator{}),
		Handler:        noCache(mux),
		SessionManager: SessionManager{sessionStore},
		sessionStore:   sessionStore,
	}
	mux.Handle("/auth/", http.StripPrefix("/auth", server.AuthRouter))
	mux.Handle("GET /{$}", templ.Handler(component))
	mux.HandleFunc("GET /host/{$}", server.GetHost)
	mux.Handle(
		"GET /static/",
		http.StripPrefix("/static", http.FileServer(
			http.Dir(staticFilesPath()))),
	)
	return server
}
