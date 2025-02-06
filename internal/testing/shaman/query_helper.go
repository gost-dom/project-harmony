package shaman

import (
	"fmt"
	"iter"
	"strings"
	"testing"

	"github.com/gost-dom/browser/dom"
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
// with [QueryHelper.FindAll] or [QueryHelper.Get].
//
// It is better to create a type implementing both [ElementPredicate] AND
// [fmt.Stringer], as it allows for better error messages when expected elements
// cannot be found.
//
// This type is exposed for the sake of easier prototyping of test code.
type ElementPredicateFunc func(dom.Element) bool

func (f ElementPredicateFunc) IsMatch(e dom.Element) bool { return f(e) }

// options treats multiple options as one, simplifying the search for multiple
// options, as well as stringifying multiple options.
type options []ElementPredicate

func (o options) IsMatch(e dom.Element) bool {
	for _, o := range o {
		if !o.IsMatch(e) {
			return false
		}
	}
	return true
}

func (o options) String() string {
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

type QueryHelper struct {
	t         *testing.T
	Container dom.ElementContainer
}

func NewQueryHelper(t *testing.T) QueryHelper { return QueryHelper{t: t} }

func (h QueryHelper) All() iter.Seq[dom.Element] {
	return func(yield func(dom.Element) bool) {
		for _, c := range h.Container.ChildNodes().All() {
			if e, ok := c.(dom.Element); ok {
				if !yield(e) {
					return
				}
				for cc := range (QueryHelper{h.t, e}).All() {
					if !yield(cc) {
						return
					}
				}
			}
		}
	}
}

func (h QueryHelper) FindAll(opts ...ElementPredicate) iter.Seq[dom.Element] {
	opt := options(opts)
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
func (h QueryHelper) Get(opts ...ElementPredicate) dom.Element {
	next, stop := iter.Pull(h.FindAll(opts...))
	defer stop()
	if v, ok := next(); ok {
		if _, ok := next(); ok {
			h.t.Fatalf("Multiple elements matching options: %s", options(opts))
		}
		return v
	}
	h.t.Fatalf("No elements mathing options: %s", options(opts))
	return nil
}

func (h QueryHelper) Scope(opts ...ElementPredicate) QueryHelper {
	r := NewQueryHelper(h.t)
	r.Container = h.Get(opts...)
	return r
}
