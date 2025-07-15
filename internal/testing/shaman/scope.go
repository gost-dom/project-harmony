package shaman

import (
	"fmt"
	"harmony/internal/testing/shaman/ariarole"
	"iter"
	"strings"
	"testing"

	"github.com/gost-dom/browser/dom"
	"github.com/gost-dom/browser/html"
)

// An ElementPredicate is a type that checks if an element matches certain
// criteria, and is used to fine elements in the dom. E.g., finding the input
// element with the label "email".
//
// An implementation of ElementPredicate should also implement [fmt.Stringer],
// describing what the predicate is looking for. This provides better error
// messages for failed queries.
type ElementPredicate interface{ IsMatch(dom.Element) bool }

// An ElementPredicateFunc wraps a single function as a predicate to be used
// with [Scope.FindAll] or [Scope.Get].
//
// This type is intended for quick prototyping of test code.
//
// It is strongly suggested to create a new type for predicates that also implements
// [fmt.Stringer].
//
// See also [ElementPredicate]
type ElementPredicateFunc func(dom.Element) bool

func (f ElementPredicateFunc) IsMatch(e dom.Element) bool { return f(e) }

// predicates treats multiple predicates as one, simplifying the search for multiple
// predicates, as well as stringifying multiple predicates.
type predicates []ElementPredicate

func (o predicates) IsMatch(e dom.Element) bool {
	for _, o := range o {
		if !o.IsMatch(e) {
			return false
		}
	}
	return true
}

func (o predicates) String() string {
	names := make([]string, len(o))
	for i, o := range o {
		if s, ok := o.(fmt.Stringer); ok {
			names[i] = s.String()
		} else {
			names[i] = "Unknown predicate. No String()"
		}
	}
	return strings.Join(names, ", ")
}

// Scope represents a subset of a page, and can be used to find elements withing
// that scope.
type Scope struct {
	t         testing.TB
	Container dom.ElementContainer
}

func NewScope(t testing.TB, c dom.ElementContainer) Scope {
	return Scope{t: t, Container: c}
}

// All returns an iterator over all elements in scope.
func (h Scope) All() iter.Seq[dom.Element] {
	return func(yield func(dom.Element) bool) {
		for _, c := range h.Container.ChildNodes().All() {
			if e, ok := c.(dom.Element); ok {
				if !yield(e) {
					return
				}
				for cc := range (Scope{h.t, e}).All() {
					if !yield(cc) {
						return
					}
				}
			}
		}
	}
}

func (h Scope) FindAll(opts ...ElementPredicate) iter.Seq[dom.Element] {
	opt := predicates(opts)
	return func(yield func(dom.Element) bool) {
		next, done := iter.Pull(h.All())
		defer done()
		for {
			e, ok := next()
			if !ok {
				return
			}
			if opt.IsMatch(e) {
				if !yield(e) {
					return
				}
			}
		}
	}
}

// Find returns an element that matches the options. At most one element is
// expected to exist in the dom mathing the options. If more than one
// is found, a fatal error is generated.
func (h Scope) Find(opts ...ElementPredicate) html.HTMLElement {
	h.t.Helper()
	next, stop := iter.Pull(h.FindAll(opts...))
	defer stop()
	if v, ok := next(); ok {
		if v2, ok := next(); ok {
			h.t.Fatalf(
				"At least two elements match options: %s\n1st match: %s\n2nd match: %s",
				predicates(opts),
				v.OuterHTML(),
				v2.OuterHTML(),
			)
		}
		return v.(html.HTMLElement)
	}
	return nil
}

// Get returns the element that matches the options. Exactly one element is
// expected to exist in the dom mathing the options. If zero, or more than one
// are found, a fatal error is generated.
func (h Scope) Get(opts ...ElementPredicate) html.HTMLElement {
	h.t.Helper()
	if res := h.Find(opts...); res != nil {
		return res
	}
	h.t.Fatalf("No elements mathing options: %s", predicates(opts))
	return nil
}

func (h Scope) Subscope(opts ...ElementPredicate) Scope {
	return NewScope(h.t, h.Get(opts...))
}

func (s Scope) Textbox(opts ...ElementPredicate) TextboxRole {
	opts = append(opts, ByRole(ariarole.Textbox))
	return TextboxRole{s.Get(opts...)}
}

func (s Scope) Checkbox(opts ...ElementPredicate) CheckboxRole {
	opts = append(opts, ByRole(ariarole.Checkbox))
	return CheckboxRole{s.Get(opts...)}
}

func (s Scope) PasswordText(opts ...ElementPredicate) TextboxRole {
	opts = append(opts, ByRole(ariarole.PasswordText))
	return TextboxRole{s.Get(opts...)}
}

// A helper to interact with "text boxes"
type TextboxRole struct {
	html.HTMLElement
}

func (tb TextboxRole) Value() string {
	v, _ := tb.HTMLElement.GetAttribute("value")
	return v
}

// Write is intended to simulate the user typing in. Currently it merely sets
// the value content attribute, making it only applicable to input elements, not
// custom implementations of the textbox aria role.
func (tb TextboxRole) Write(input string) { tb.SetAttribute("value", input) }

func (tb TextboxRole) Clear() { tb.SetAttribute("value", "") }

func (tb TextboxRole) ARIADescription() string {
	return GetDescription(tb)
}

type CheckboxRole struct {
	html.HTMLElement
}

func (cb CheckboxRole) Check()   { cb.setChecked(true) }
func (cb CheckboxRole) Uncheck() { cb.setChecked(false) }

func (cb CheckboxRole) setChecked(val bool) {
	input, ok := cb.HTMLElement.(html.HTMLInputElement)
	if !ok {
		// To support generic checkbox roles, the approach should probably be to
		// check for the presence of the aria-checked content attribute, and
		// call Click() on the element if it has the wrong state.
		panic("CheckboxRole.Check/Uncheck: only input elements are supported")
	}
	input.SetChecked(val)
}
