package browsertest

import (
	"harmony/internal/testing/shaman"
	"harmony/internal/testing/shaman/ariarole"
	. "harmony/internal/testing/shaman/predicates"
	"testing"

	"github.com/gost-dom/browser/html"
)

// Page represents the standard page layout
type Page struct {
	t   testing.TB
	win html.Window
}

func NewPage(t testing.TB, win html.Window) Page { return Page{t, win} }

func (p Page) Header() Header {
	scope := shaman.
		WindowScope(p.t, p.win).
		Subscope(ByRole(ariarole.Banner))
	return Header{scope}
}

type Header struct{ shaman.Scope }

// LoginBtn returns the login "button" if found, otherwise it returns nil.
func (h Header) LoginBtn() html.HTMLElement {
	return h.Find(ByRole(ariarole.Link), ByName("Login"))
}

// LogoutBtn returns the logout button if found, otherwise it returns nil.
func (h Header) LogoutBtn() html.HTMLElement {
	return h.Find(ByRole(ariarole.Button), ByName("Logout"))
}
