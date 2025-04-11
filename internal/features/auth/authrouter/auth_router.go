package authrouter

import (
	"context"
	"errors"
	"harmony/internal/features/auth"
	"harmony/internal/features/auth/authrouter/views"
	"net/http"
)

type Authenticator interface {
	Authenticate(context.Context, string, string) (auth.Account, error)
}

type AuthRouter struct {
	*http.ServeMux
	Authenticator  Authenticator
	SessionManager SessionManager
}

func (s *AuthRouter) PostAuthLogin(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	email := r.FormValue("email")
	password := r.FormValue("password")
	redirectUrl := r.FormValue("redirectUrl")
	if redirectUrl == "" {
		redirectUrl = "/"
	}
	if account, err := s.Authenticator.Authenticate(r.Context(), email, password); err == nil {
		s.SessionManager.SetAccount(w, r, account)
		w.Header().Add("hx-push-url", redirectUrl)
	} else {
		authError := errors.Is(err, auth.ErrBadCredentials)
		data := views.LoginFormData{
			Email:              email,
			Password:           "",
			InvalidCredentials: authError,
			UnexpectedError:    !authError,
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
