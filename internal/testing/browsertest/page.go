package browsertest

import (
	"testing"

	"github.com/gost-dom/browser/html"
	"github.com/gost-dom/shaman"
	"github.com/gost-dom/shaman/ariarole"
	. "github.com/gost-dom/shaman/predicates"
	"github.com/stretchr/testify/assert"
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

func (p Page) Scope() shaman.Scope {
	return shaman.WindowScope(p.t, p.win)
}

func (p Page) AssertLoginPage() LoginPage {
	main := p.Scope().Subscope(ByRole(ariarole.Main))
	h1 := main.Get(ByH1)
	if !assert.Equal(p.t, "Login", h1.TextContent()) {
		p.t.Fatal("Not a login page")
	}
	return LoginPage{p}
}
