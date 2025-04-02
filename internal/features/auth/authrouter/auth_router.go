package authrouter

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"harmony/internal/features/auth"
	"harmony/internal/features/auth/authdomain"
	"harmony/internal/features/auth/authdomain/password"
	"harmony/internal/features/auth/authrouter/views"

	"github.com/a-h/templ"
)

type Authenticator interface {
	Authenticate(
		context.Context,
		string,
		password.Password,
	) (authdomain.AuthenticatedAccount, error)
}

type AuthRouter struct {
	*http.ServeMux
	Authenticator  Authenticator
	SessionManager SessionManager
}

func (s *AuthRouter) PostRegister(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("hx-push-url", "./validate-email")
	w.Header().Add("hx-Retarget", "body")
	fmt.Println("REGISTER!")
	views.ValidateEmailPage().Render(r.Context(), w)
}

func (s *AuthRouter) PostAuthLogin(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	email := r.FormValue("email")
	pw := r.FormValue("password")
	redirectUrl := r.FormValue("redirectUrl")
	if redirectUrl == "" {
		redirectUrl = "/"
	}
	if account, err := s.Authenticator.Authenticate(r.Context(), email, password.Parse(pw)); err == nil {
		if err := s.SessionManager.SetAccount(w, r, account); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
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
	r.Handle("GET /register", templ.Handler(views.Register()))
	r.HandleFunc("POST /register", r.PostRegister)
	r.Handle("GET /validate-email", templ.Handler(views.ValidateEmailPage()))
}

func (*AuthRouter) RenderLogin(w http.ResponseWriter, r *http.Request) {
	views.Login("/host", views.LoginFormData{}).Render(r.Context(), w)
}
