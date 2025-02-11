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
	"github.com/samber/do"
)

const sessionCookieName = "accountId"

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

func (a *authenticator) Authenticate(
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

func (m *SessionManager) LoggedInUser(r *http.Request) *Account {
	session, _ := m.sessionStore.Get(r, sessionNameAuth)
	if id, ok := session.Values[sessionCookieName]; ok {
		result := new(Account)
		if strId, ok := id.(string); ok {
			result.Id = AccountId(strId)
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
	views.AuthLogin("/host").Render(r.Context(), w)
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

var Injector *do.Injector = do.New()

func NewServer(
	sessionStore sessions.Store,
	sessionManager SessionManager,
	authRouter *AuthRouter,
) *Server {
	component := views.Index()

	mux := http.NewServeMux()
	server := &Server{
		AuthRouter:     authRouter, //NewAuthRouter(sessionStore, authenticator{}),
		Handler:        noCache(mux),
		SessionManager: sessionManager,
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

func init() {
	gob.Register(AccountId(""))
	do.ProvideValue[Authenticator](Injector, &authenticator{})
	do.Provide(Injector, func(i *do.Injector) (*AuthRouter, error) {
		return NewAuthRouter(
			do.MustInvoke[sessions.Store](i),
			do.MustInvoke[Authenticator](i),
		), nil
	})
	do.Provide(Injector, func(i *do.Injector) (sessions.Store, error) {
		return memstore.NewMemStore(
			[]byte("authkey123"),
			[]byte("enckey12341234567890123456789012"),
		), nil
	})
	do.Provide(Injector, func(i *do.Injector) (*Server, error) {
		return NewServer(
			do.MustInvoke[sessions.Store](i),
			do.MustInvoke[SessionManager](i),
			do.MustInvoke[*AuthRouter](i),
		), nil
	})
	do.Provide(Injector, func(i *do.Injector) (SessionManager, error) {
		sessionStore := do.MustInvoke[sessions.Store](i)
		return SessionManager{sessionStore}, nil
	})
}

func New() *Server {
	return do.MustInvoke[*Server](Injector)
}
