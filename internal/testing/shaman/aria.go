package shaman

import (
	"fmt"
	"strings"

	"github.com/gost-dom/browser/dom"
)

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

func GetDescription(e dom.Element) string {
	if id, ok := e.GetAttribute("aria-describedby"); ok {
		ids := strings.Split(id, " ")
		ss := make([]string, 0, len(ids))
		for _, id := range ids {
			if e := e.OwnerDocument().GetElementById(id); e != nil {
				ss = append(ss, e.TextContent())
			}
		}
		return strings.Join(ss, "\n")
	}
	s, _ := e.GetAttribute("aria-description")
	return s
}
