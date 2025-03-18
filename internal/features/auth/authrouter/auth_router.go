package authrouter

import (
	"context"
	"harmony/internal/features/auth"
	"harmony/internal/features/auth/authrouter/views"
	"net/http"

	"github.com/gorilla/sessions"
)

type Authenticator interface {
	Authenticate(context.Context, string, string) (auth.Account, error)
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
		views.Login(redirectUrl, views.LoginFormData{}).Render(r.Context(), w)

	})
	r.HandleFunc("POST /login", r.PostAuthLogin)
}

func (*AuthRouter) RenderLogin(w http.ResponseWriter, r *http.Request) {
	views.Login("/host", views.LoginFormData{}).Render(r.Context(), w)
}
