package ariarole_test

import (
	"fmt"
	"testing"

	"harmony/internal/testing/shaman/ariarole"

	"github.com/gost-dom/browser/html"
)

func TestARIARoles(t *testing.T) {
	createElement := newRoleHelper().createElement

	specs := []struct {
		TagName  string
		RoleAttr string
		Want     ariarole.Role
	}{
		{TagName: "button", RoleAttr: "button", Want: ariarole.Button},
		{TagName: "", RoleAttr: "alert", Want: ariarole.Alert},
	}

	for _, spec := range specs {
		name := fmt.Sprintf("Test aria role: %s", spec.RoleAttr)
		t.Run(name, func(t *testing.T) {
			d := createElement("div")
			d.SetAttribute("role", spec.RoleAttr)
			assertRole(t, spec.Want, d)
		})
		if spec.TagName != "" {
			t.Run(
				fmt.Sprintf("Test aria role for <%s>", spec.TagName),
				func(t *testing.T) {
					d := createElement(spec.TagName)
					assertRole(t, spec.Want, d)
				})
		}
	}
}

type roleHelper struct {
	doc html.HTMLDocument
}

func newRoleHelper() roleHelper {
	win := html.NewWindow()
	doc := html.NewHTMLDocument(win)
	return roleHelper{doc}
}

func (h roleHelper) createElement(tagname string) html.HTMLElement {
	return h.doc.CreateElement(tagname).(html.HTMLElement)
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
