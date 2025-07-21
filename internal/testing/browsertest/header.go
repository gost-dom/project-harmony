package browsertest

import (
	"github.com/gost-dom/browser/html"
	"github.com/gost-dom/shaman"
	"github.com/gost-dom/shaman/ariarole"
	. "github.com/gost-dom/shaman/predicates"
)

type Header struct{ shaman.Scope }

// LoginBtn returns the login "button" if found, otherwise it returns nil.
func (h Header) LoginBtn() html.HTMLElement {
	return h.Find(ByRole(ariarole.Link), ByName("Login"))
}

// LogoutBtn returns the logout button if found, otherwise it returns nil.
func (h Header) LogoutBtn() html.HTMLElement {
	return h.Find(ByRole(ariarole.Button), ByName("Logout"))
}
