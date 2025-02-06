package server_test

import (
	"harmony/internal/server"
	ariarole "harmony/internal/testing/aria-role"
	"harmony/internal/testing/shaman"
	. "harmony/internal/testing/shaman/predicates"
	"harmony/internal/testing/shaman/sync"
	"net/http"
	"testing"

	"github.com/gost-dom/browser"
	"github.com/gost-dom/browser/dom"
	"github.com/gost-dom/browser/html"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

func init() {
	// slog.SetLogLoggerLevel(slog.LevelWarn)
	// logger.SetDefault(slog.Default())
}

func TestCanServe(t *testing.T) {
	b := browser.NewBrowserFromHandler(server.New())
	w, err := b.Open("http://localhost:1234/") // host is imaginary - just need to exist
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
	Sync sync.EventSync
	shaman.QueryHelper
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

func (s *NavigateToLoginSuite) Q() shaman.QueryHelper {
	q := shaman.NewQueryHelper(s.T())
	q.Container = s.win.Document()
	return q
}

func (s *NavigateToLoginSuite) SetupTest() {
	var err error
	s.QueryHelper = shaman.NewQueryHelper(s.T())
	b := browser.NewBrowserFromHandler(server.New())
	s.win, err = b.Open("http://localhost:1234/")
	s.Sync = sync.SetupEventSync(s.win)
	s.QueryHelper.Container = s.win.Document()

	s.Sync.WaitFor("htmx:load")
	s.NoError(err)
}

func (s *NavigateToLoginSuite) TestClickLoginLink() {
	s.Q().Get(ByRole(ariarole.Link), ByName("Login")).Click()
	s.Equal("/auth/login", s.win.Location().Pathname())
	mainHeading := getMainHeading(s.T(), s.win)
	s.Equal("Login", mainHeading.TextContent())
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
	s.Get(ByRole(ariarole.Link), ByName("Go to hosting")).Click()
	s.Sync.WaitFor("htmx:afterSettle")
	s.Equal("/auth/login", s.win.Location().Pathname(), "Location after host")
	mainHeading := getMainHeading(s.T(), s.win)
	s.Equal("Login", mainHeading.TextContent())

	s.Get(ByRole(ariarole.Textbox), ByName("Email")).
		SetAttribute("value", "valid-user@example.com")

	s.Get(ByRole(ariarole.PasswordText), ByName("Password")).
		SetAttribute("value", "s3cret")

	s.Get(ByRole(ariarole.Button), ByName("Sign in")).Click()
	s.Sync.WaitFor("htmx:afterSettle")
	s.Equal("/host", s.win.Location().Pathname())
}

func TestNavigateToLogin(t *testing.T) {
	suite.Run(t, new(NavigateToLoginSuite))
}
