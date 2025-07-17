package browsertest

import (
	"harmony/internal/testing/shaman"
	"harmony/internal/testing/shaman/ariarole"
	. "harmony/internal/testing/shaman/predicates"

	"github.com/gost-dom/browser/html"
)

// A helper to make test code interacting with the login form more expressive.
type LoginForm struct {
	shaman.Scope
}

func NewLoginForm(s shaman.Scope) LoginForm {
	return LoginForm{s.Subscope(ByRole(ariarole.Main)).Subscope(ByRole(ariarole.Form))}
}

func (f LoginForm) Email() shaman.TextboxRole {
	return f.Textbox(ByName("Email"))
}

func (f LoginForm) Password() shaman.TextboxRole {
	return f.PasswordText(ByRole(ariarole.PasswordText), ByName("Password"))
}

func (f LoginForm) SubmitBtn() html.HTMLElement {
	return f.Get(ByRole(ariarole.Button), ByName("Sign in"))
}

type LoginPage struct {
	Page
}

func (p LoginPage) LoginForm() LoginForm {
	return NewLoginForm(p.Scope())
}
