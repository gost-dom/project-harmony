package authrouter

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"net/mail"

	"harmony/internal/auth"
	"harmony/internal/auth/authdomain"
	"harmony/internal/auth/authdomain/password"
	"harmony/internal/auth/authrouter/views"
	serverctx "harmony/internal/web/server/ctx"

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

type EmailValidator interface {
	Validate(
		ctx context.Context,
		input auth.ValidateEmailInput,
	) (authdomain.AuthenticatedAccount, error)
}

type AuthRouter struct {
	*http.ServeMux
	Authenticator  Authenticator
	Registrator    Registrator
	SessionManager SessionManager
	EmailValidator EmailValidator
}

func (s *AuthRouter) PostRegister(w http.ResponseWriter, r *http.Request) {
	// Consider improving this implementation for more clearly separate
	// - Parsing and validating input
	// - Call use case on valid input
	// - Rerender form in invalid input, or use case error
	// - Redirect to email validation page on success
	// TODO: Provide good error messages for users for all scenarios
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

		serverctx.SetReqValue(&r, auth.CtxKeyAuthAccount, account.Account)
		w.Header().Add("hx-push-url", redirectUrl)
		w.Header().Add("hx-retarget", "body")
		rewrite(w, r, redirectUrl, "")
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

// Init implements interface [surgeon.Initer].
func (r *AuthRouter) Init() {
	r.ServeMux = http.NewServeMux()
	r.HandleFunc("GET /login", func(w http.ResponseWriter, r *http.Request) {
		redirectUrl := r.URL.Query().Get("redirectUrl")
		views.Login(redirectUrl, views.LoginFormData{}).Render(r.Context(), w)
	})
	r.HandleFunc("POST /login", r.PostAuthLogin)
	r.HandleFunc("POST /logout", r.postLogout)
	r.HandleFunc("GET /register", func(w http.ResponseWriter, r *http.Request) {
		views.Register(views.RegisterFormData{}).Render(r.Context(), w)
	})
	r.HandleFunc("POST /register", r.PostRegister)
	r.HandleFunc("GET /validate-email",
		func(w http.ResponseWriter, r *http.Request) {
			views.ValidateEmailPage(views.ValidateEmailForm{
				EmailAddress: r.URL.Query().Get("email"),
			}).Render(r.Context(), w)
		},
	)
	r.HandleFunc("POST /validate-email", r.postValidateEmail)
}

func (router *AuthRouter) postLogout(w http.ResponseWriter, r *http.Request) {
	if err := router.SessionManager.Logout(w, r); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/", 303)
}

func (router *AuthRouter) postValidateEmail(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	email, _ := mail.ParseAddress(r.FormValue("email"))
	code := r.FormValue("challenge-response")
	account, err := router.EmailValidator.Validate(r.Context(),
		auth.ValidateEmailInput{
			Email: email,
			Code:  authdomain.EmailValidationCode(code),
		})
	if err != nil {
		w.Header().Add("hx-retarget", "#validation-error-container")
		w.Header().Add("hx-swap", "innerHTML")
		if errors.Is(err, auth.ErrBadChallengeResponse) {
			views.InvalidCodeError().Render(r.Context(), w)
		} else {
			views.UnexpectedError().Render(r.Context(), w)
		}
		return
	}
	if err := router.SessionManager.SetAccount(w, r, account); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Add("hx-push-url", "/host")
	w.Header().Add("hx-retarget", "body")
	rewrite(w, r, "/host", "")
}

func (*AuthRouter) RenderHost(w http.ResponseWriter, r *http.Request) {
	views.Login("/host", views.LoginFormData{}).Render(r.Context(), w)
}

func New() *AuthRouter {
	r := &AuthRouter{}
	r.Init()
	return r
}
