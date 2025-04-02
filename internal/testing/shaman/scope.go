package shaman

import (
	"fmt"
	ariarole "harmony/internal/testing/aria-role"
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
	t         *testing.T
	Container dom.ElementContainer
}

func NewScope(t *testing.T, c dom.ElementContainer) Scope {
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

// Get returns the element that matches the options. Exactly one element is
// expected to exist in the dom mathing the options. If zero, or more than one
// are found, a fatal error is generated.
func (h Scope) Get(opts ...ElementPredicate) html.HTMLElement {
	next, stop := iter.Pull(h.FindAll(opts...))
	defer stop()
	if v, ok := next(); ok {
		if _, ok := next(); ok {
			h.t.Fatalf("Multiple elements matching options: %s", predicates(opts))
		}
		return v.(html.HTMLElement)
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

// A helper to interact with "text boxes"
type TextboxRole struct {
	html.HTMLElement
}

// Write is intended to simulate the user typing in. Currently it merely sets
// the value content attribute, making it only applicable to input elements, not
// custom implementations of the textbox aria role.
func (tb TextboxRole) Write(input string) {
	tb.SetAttribute("value", input)
}
