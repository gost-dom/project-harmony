package shaman

import (
	"fmt"
	ariarole "harmony/internal/testing/aria-role"

	"github.com/gost-dom/browser/dom"
)

// ByName is an [ElementMatcher] that matches elements by their accessibility
// name. E.g., the element's label, or text content. The label can be
//   - An associated label element
//   - The value of an aria-label property
//   - The text content of an element referenced by an aria-labelledby property.
//
// This is called "name" not "label", as the term in ARIA is
// "accessibility name", which is why this is called name, not label.
type ByName string

func (n ByName) IsMatch(e dom.Element) bool { return GetName(e) == string(n) }

func (n ByName) String() string { return fmt.Sprintf("By accessibility name: %s", string(n)) }

// An [ElementPredicate] that matches elements by their [ARIA role].
//
// [ARIA role]: https://developer.mozilla.org/en-US/docs/Web/Accessibility/ARIA/Roles
type ByRole ariarole.Role

func (r ByRole) IsMatch(
	e dom.Element,
) bool {
	return ariarole.GetElementRole(e) == ariarole.Role(r)
}

func (r ByRole) String() string { return fmt.Sprintf("By role: %s", string(r)) }
