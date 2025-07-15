package ariarole_test

import (
	ariarole "harmony/internal/testing/aria-role"
	"testing"

	"github.com/gost-dom/browser/html"
)

func TestAriaRoleButton(t *testing.T) {
	btn := createElement("button")
	assertRole(t, ariarole.Button, btn)

	d := createElement("div")
	d.SetAttribute("role", "button")
	assertRole(t, ariarole.Button, d)
}

func createElement(tagname string) html.HTMLElement {
	win := html.NewWindow()
	doc := html.NewHTMLDocument(win)
	return doc.CreateElement(tagname).(html.HTMLElement)
}

func assertRole(t testing.TB, want ariarole.Role, e html.HTMLElement) {
	t.Helper()
	got := ariarole.GetElementRole(e)
	if got != want {
		t.Errorf(
			"expected ARIA role: %s, got: %s\nElement: %s",
			want, got, e.OuterHTML(),
		)
	}
}
