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
	Link         Role = "link"
	Button       Role = "button"
	Textbox      Role = "textbox"
	PasswordText Role = "password text"
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
	case "BUTTON":
		return Button
	case "A":
		return Link
	}
	return None
}
