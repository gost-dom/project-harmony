package ariarole

import "github.com/gost-dom/browser/dom"

// Role represents an [ARIA role]. See package documentation for more
// information of aria roles.
//
// Two special values are defined in this library:
//   - [None] represents an element has no role
//   - [PasswordText] represents password input fields
//
// PasswordText is not an official ARIA role. It is reported by Firefox's
// accessibility tools, and is helpful in locating password fields, as the
// element <input type="password" /> does not have an ARIA role.
//
// [ARIA roles]: https://developer.mozilla.org/en-US/docs/Web/Accessibility/ARIA/Roles
type Role string

const (
	// None represents an element that doesn't have a role specified.
	None   Role = ""
	Alert  Role = "alert"
	Button Role = "button"
	Form   Role = "form"
	Link   Role = "link"
	Main   Role = "main"
	Banner Role = "banner"

	// PasswordText represents the "password text" role, which isn't an official
	// ARIA role. It is reported by Firefox's accessibility tools, and helpful
	// as password fields don't actually have an official role, i.e., you cannot
	// find them as a textbox role.
	PasswordText Role = "password text"
	Textbox      Role = "textbox"
	Checkbox     Role = "checkbox"
)

var elementRoles map[string]Role = map[string]Role{
	"MAIN":   Main,
	"BUTTON": Button,
	"A":      Link,
	"FORM":   Form,
	"HEADER": Banner,
}

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
			case "checkbox":
				return Checkbox
			case "button", "submit", "reset":
				return Button
			}
			return Textbox
		}
	}
	return elementRoles[e.TagName()]
}
