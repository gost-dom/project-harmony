package authrouter

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"net/mail"

	"harmony/internal/features/auth"
	"harmony/internal/features/auth/authdomain"
	"harmony/internal/features/auth/authdomain/password"
	"harmony/internal/features/auth/authrouter/views"

	"github.com/a-h/templ"
	"github.com/gorilla/schema"
)

var decoder = schema.NewDecoder()

func init() {
	decoder.IgnoreUnknownKeys(true)
}

type Authenticator interface {
	Authenticate(
		context.Context,
		string,
		password.Password,
	) (authdomain.AuthenticatedAccount, error)
}

type Registrator interface {
	Register(ctx context.Context, input auth.RegistratorInput) error
}

type AuthRouter struct {
	*http.ServeMux
	Authenticator  Authenticator
	Registrator    Registrator
	SessionManager SessionManager
}

func (s *AuthRouter) PostRegister(w http.ResponseWriter, r *http.Request) {
	// This is a crappy implementation. But I can't be bothered to improve any
	// more right now.
	r.ParseForm()
	data := struct {
		Fullname         string `schema:"fullname,required"`
		Email            string `schema:"email,required"`
		DisplayName      string `schema:"displayname"`
		Password         string `schema:"password"`
		TermsOfUse       bool   `schema:"terms-of-use,required"`
		NewsletterSignup bool   `schema:"newsletter-signup"`
	}{}
	var formData views.RegisterFormData
	err := decoder.Decode(&data, r.PostForm)
	var registerInput auth.RegistratorInput
	if err == nil {
		if registerInput.Email, err = mail.ParseAddress(data.Email); err != nil {
			formData.Email.Value = data.Email
			formData.Email.Errors = []string{"Must be a valid email address"}
		}
	}
	if err == nil {
		registerInput.Password = password.Parse(data.Password)
	}
	if data.Email == "" {
		formData.Email.Errors = []string{"Must be filled out"}
	}
	if err == nil {
		registerInput.Name = data.Fullname
		registerInput.DisplayName = data.DisplayName
		registerInput.NewsletterSignup = data.NewsletterSignup
		err = s.Registrator.Register(r.Context(), registerInput)
	}
	if err != nil {
		slog.Error("error", "err", err)
		formData.Fullname = data.Fullname
		formData.DisplayName = data.DisplayName
		formData.TermsOfUse = data.TermsOfUse
		formData.TermsOfUseMissing = !data.TermsOfUse
		formData.NewsletterSignup = data.NewsletterSignup
		views.RegisterFormContents(formData).Render(r.Context(), w)
		return
	}

	w.Header().Add("hx-push-url", "./validate-email")
	w.Header().Add("hx-Retarget", "body")
	views.ValidateEmailPage(views.ValidateEmailForm{
		EmailAddress: data.Email,
	}).Render(r.Context(), w)
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
	r.HandleFunc("GET /register", func(w http.ResponseWriter, r *http.Request) {
		views.Register(views.RegisterFormData{}).Render(r.Context(), w)
	})
	r.HandleFunc("POST /register", r.PostRegister)
	r.Handle(
		"GET /validate-email",
		templ.Handler(views.ValidateEmailPage(views.ValidateEmailForm{})),
	)
}

func (*AuthRouter) RenderLogin(w http.ResponseWriter, r *http.Request) {
	views.Login("/host", views.LoginFormData{}).Render(r.Context(), w)
}
