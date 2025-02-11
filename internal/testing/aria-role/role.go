package ariarole

import "github.com/gost-dom/browser/dom"

// The Role represents an [ARIA role]. The contstants equal to their values in
// the spec, but with the special value [None] that represents an element has no
// role.
//
// [ARIA role]: https://developer.mozilla.org/en-US/docs/Web/Accessibility/ARIA/Roles
type Role string

const (
	// None represents an element that doesn't have a role specified.
	None         Role = ""
	Alert        Role = "alert"
	Button       Role = "button"
	Form         Role = "form"
	Link         Role = "link"
	Main         Role = "main"
	PasswordText Role = "password text"
	Textbox      Role = "textbox"
)

func GetElementRole(e dom.Element) Role {
	if r, ok := e.GetAttribute("role"); ok {
		// TODO: check validity of r
		return Role(r)
	}
	switch e.TagName() {
	case "INPUT":
		if t, ok := e.GetAttribute("type"); ok {
			switch t {
			case "password":
				return PasswordText

			case "button":
				fallthrough
			case "submit":
				fallthrough
			case "reset":
				return Button
			}
			return Textbox
		}
	case "MAIN":
		return Main
	case "BUTTON":
		return Button
	case "A":
		return Link
	case "FORM":
		return Form
	}
	return None
}
