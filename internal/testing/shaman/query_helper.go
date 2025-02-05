package shaman

import (
	"fmt"
	ariarole "harmony/internal/testing/aria-role"
	"iter"
	"testing"

	"github.com/gost-dom/browser/dom"
	"github.com/gost-dom/browser/html"
	"github.com/stretchr/testify/assert"
)

type QueryHelper struct {
	t         *testing.T
	Container dom.ElementContainer
}

func NewQueryHelper(t *testing.T) QueryHelper { return QueryHelper{t: t} }

type ElementPredicate func(dom.Element) bool

func ByRole(role ariarole.Role) ElementPredicate {
	return func(e dom.Element) bool {
		return ariarole.GetElementRole(e) == role
	}
}

// Gets the by their accessibility name of an element. I.e., an associated
// label, the value of an aria-label, or the text-content of an element
// referenced by an aria-labelledby property
func GetName(e dom.Element) string {
	// TODO: This should be exposed as IDL attributes
	if l, ok := e.GetAttribute("aria-label"); ok {
		return l
	}
	doc := e.OwnerDocument()
	if l, ok := e.GetAttribute("aria-labelledby"); ok {
		if labelElm := doc.GetElementById(l); labelElm != nil {
			return labelElm.TextContent()
		}
	}
	switch e.TagName() {
	case "INPUT":
		if id, ok := e.GetAttribute("id"); ok {
			if label, _ := doc.QuerySelector(fmt.Sprintf("label[for='%s']", id)); label != nil {
				return label.TextContent()
			}
		}
	}
	return e.TextContent()
}

// Finds elements by their accessibility name. I.e., an associated label, the
// value of an aria-label, or the text-content of an element referenced by an
// aria-labelledby property
func ByName(name string) ElementPredicate {
	return func(e dom.Element) bool { return GetName(e) == name }
}

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

func (h QueryHelper) FindAll(options ...ElementPredicate) iter.Seq[dom.Element] {
	return func(yield func(dom.Element) bool) {
		next, done := iter.Pull(h.All())
		defer done()
	loop:
		for {
			e, ok := next()
			if !ok {
				return
			}
			for _, o := range options {
				if !o(e) {
					continue loop
				}
			}
			if !yield(e) {
				return
			}
		}
	}
}

// Get returns the element that matches the options. Exactly one element is
// expected to exist in the dom mathing the options. If zero, or more than one
// are found, a fatal error is generated.
func (h QueryHelper) Get(options ...ElementPredicate) dom.Element {
	next, stop := iter.Pull(h.FindAll(options...))
	defer stop()
	if v, ok := next(); ok {
		if _, ok := next(); ok {
			h.t.Fatal("Multiple elements matching options")
		}
		return v
	}
	h.t.Fatal("No elements mathing options")
	return nil
}

func (h QueryHelper) FindLinkWithName(name string) html.HTMLAnchorElement {
	as, err := h.Container.QuerySelectorAll("a")
	t := h.t
	assert.NoError(t, err)
	var res html.HTMLAnchorElement
	for _, e := range as.All() {
		a, ok := e.(html.HTMLAnchorElement)
		if !ok {
			t.Fatalf(
				"Something very very wrong in the dom. Element was found as an 'a', but not an HTMLAnchorElement: %s",
				e,
			)
		}
		if a.TextContent() == name {
			if res != nil {
				t.Fatalf("Expected to find one anchor with name, '%s'. Found multiple", name)
			}
			res = a
		}
	}
	if res == nil {
		t.Fatalf("Expected to find one anchor with name, '%s'. Found none", name)
	}
	return res
}
