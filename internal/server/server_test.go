package server_test

import (
	"fmt"
	"harmony/internal/server"
	"log/slog"
	"net/http"
	"testing"

	"github.com/gost-dom/browser"
	"github.com/gost-dom/browser/dom"
	"github.com/gost-dom/browser/html"
	"github.com/gost-dom/browser/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

func init() {
	logger.SetDefault(slog.Default())
}

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
	fmt.Println("*** BOOOOOOOOOOOOOOOOOOOOOOH")
	c := make(chan struct{})
	s.win.AddEventListener("htmx:afterOnLoad", dom.NewEventHandlerFunc(func(e dom.Event) error {
		fmt.Println("EVENT!!!!")
		c <- struct{}{}
		return nil
	}))

	loginLink := s.Q().FindLinkWithName("Login")
	loginLink.Click()
	// We should be on the login path
	assert.Equal(s.T(), "/auth/login", s.win.Location().Pathname())
	// fmt.Println(s.win.Document().Body().OuterHTML())
	mainHeading := getMainHeading(s.T(), s.win)
	assert.Equal(s.T(), "Login", mainHeading.TextContent())
	// TODO: Verify that the window doesn't navigate

}

func getMainHeading(t *testing.T, w html.Window) dom.Element {
	ee, err := w.Document().QuerySelectorAll("h1")
	assert.NoError(t, err)
	assert.Equal(t, 1, ee.Length(), "Expected exactly one <h1> element in the document")
	e, ok := ee.Item(0).(dom.Element)
	assert.True(t, ok, "The found <h1> was expected to be e dom Element")
	return e
}

func (s *NavigateToLoginSuite) TestLoginFlow() {
	c := make(chan struct{})
	loginLink := s.Q().FindLinkWithName("Go to hosting")
	s.win.AddEventListener("htmx:load", dom.NewEventHandlerFunc(func(e dom.Event) error {
		go func() { c <- struct{}{} }()
		return nil
	}))

	loginLink.Click()
	<-c
	// We should be on the login path
	assert.Equal(s.T(), "/auth/login", s.win.Location().Pathname(), "Location after host")
	mainHeading := getMainHeading(s.T(), s.win)
	assert.Equal(s.T(), "Login", mainHeading.TextContent())
	// TODO: Verify that the window doesn't navigate

	fmt.Println("TEST IS DONE DONE DONE\n\n----")
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
