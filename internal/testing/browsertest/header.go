package browsertest

import (
	"harmony/internal/testing/shaman"
	"harmony/internal/testing/shaman/ariarole"
	. "harmony/internal/testing/shaman/predicates"

	"github.com/gost-dom/browser/html"
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
