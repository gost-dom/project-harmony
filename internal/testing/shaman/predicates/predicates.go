package predicates

import (
	"fmt"
	"harmony/internal/testing/shaman"
	"harmony/internal/testing/shaman/ariarole"

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

func (n ByName) IsMatch(e dom.Element) bool { return shaman.GetName(e) == string(n) }

func (n ByName) String() string { return fmt.Sprintf("By accessibility name: %s", string(n)) }

// An ElementPredicate is a type that checks if an element matches certain
// criteria, and is used to fine elements in the dom. E.g., finding the input
// element with the label "email".
type ElementPredicate interface{ IsMatch(dom.Element) bool }

type ElementPredicateFunc func(dom.Element) bool

func (f ElementPredicateFunc) IsMatch(e dom.Element) bool { return f(e) }

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

type ByH1Predicate struct{}

// Predicate to find THE <h1> element on the page.
//
// Every page should have exactly one H1 element, and it has a special meaning
// as the page title. (The <title> in the document <head> is for browser title
// and bookmarks, and will often include prefixes, such as site name or other,
// which is why the page itself should have a title).
var ByH1 = ByH1Predicate{}

func (s ByH1Predicate) IsMatch(e dom.Element) bool { return e.TagName() == "H1" }

func (s ByH1Predicate) String() string { return "Main heading (<h1>)" }
