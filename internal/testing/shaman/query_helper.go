package shaman

import (
	"fmt"
	ariarole "harmony/internal/testing/aria-role"
	"iter"
	"strings"
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

type ElementPredicate interface{ IsMatch(dom.Element) bool }

type ByRole ariarole.Role

func (r ByRole) IsMatch(e dom.Element) bool {
	return ariarole.GetElementRole(e) == ariarole.Role(r)
}
func (r ByRole) String() string { return fmt.Sprintf("By role: %s", string(r)) }

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
type ByName string

func (n ByName) IsMatch(e dom.Element) bool { return GetName(e) == string(n) }

func (n ByName) String() string { return fmt.Sprintf("By accessibility name: %s", string(n)) }

type Options []ElementPredicate

func (o Options) IsMatch(e dom.Element) bool {
	for _, o := range o {
		if !o.IsMatch(e) {
			return false
		}
	}
	return true
}

func (o Options) String() string {
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
	opt := Options(options)
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
func (h QueryHelper) Get(options ...ElementPredicate) dom.Element {
	next, stop := iter.Pull(h.FindAll(options...))
	defer stop()
	if v, ok := next(); ok {
		if _, ok := next(); ok {
			h.t.Fatalf("Multiple elements matching options: %s", Options(options))
		}
		return v
	}
	h.t.Fatalf("No elements mathing options: %s", Options(options))
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
