package server_test

import (
	"harmony/internal/server"
	"net/http"
	"testing"

	"github.com/gost-dom/browser"
	"github.com/gost-dom/browser/dom"
	"github.com/gost-dom/browser/html"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

func TestCanServe(t *testing.T) {
	b := browser.NewBrowserFromHandler(server.New())
	w, err := b.Open("/")
	if err != nil {
		t.Fatal(err)
	}
	h1, err := w.Document().Body().QuerySelector("h1")
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "Project Harmony", h1.TextContent())
}

type NavigateToLoginSuite struct {
	suite.Suite
	win html.Window
}

type RequestRecorder struct {
	Handler  http.Handler
	Requests []*http.Request
}

func (r *RequestRecorder) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.Requests = append(r.Requests, req)
	r.Handler.ServeHTTP(w, req)
}

func (r RequestRecorder) requestCount() int { return len(r.Requests) }

func (s NavigateToLoginSuite) Q() QueryHelper {
	return QueryHelper{s.T(), s.win.Document()}
}

func (s *NavigateToLoginSuite) SetupTest() {
	var err error
	b := browser.NewBrowserFromHandler(server.New())
	s.win, err = b.Open("/")
	assert.NoError(s.T(), err)
}

func (s *NavigateToLoginSuite) TestClickLoginLink() {
	loginLink := s.Q().FindLinkWithName("Login")
	loginLink.Click()
	// We should be on the login path
	assert.Equal(s.T(), "/auth/login", s.win.Location().Pathname())
	// TODO: Verify that the window doesn't navigate
}

type QueryHelper struct {
	T         *testing.T
	Container dom.ElementContainer
}

func (h QueryHelper) FindLinkWithName(name string) html.HTMLAnchorElement {
	as, err := h.Container.QuerySelectorAll("a")
	assert.NoError(h.T, err)
	var res *html.HTMLAnchorElement
	for _, e := range as.All() {
		a, ok := e.(html.HTMLAnchorElement)
		if !ok {
			h.T.Fatalf(
				"Something very very wrong in the dom. Element was found as an 'a', but not an HTMLAnchorElement: %s",
				e,
			)
		}
		if a.TextContent() == name {
			if res != nil {
				h.T.Fatalf("Expected to find one anchor with name, '%s'. Found multiple", name)
			}
			res = &a
		}
	}
	if res == nil {
		h.T.Fatalf("Expected to find one anchor with name, '%s'. Found none", name)
	}
	return *res
}

func TestNavigateToLogin(t *testing.T) {
	suite.Run(t, new(NavigateToLoginSuite))
}
