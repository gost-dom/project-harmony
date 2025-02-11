package server_test

import (
	"harmony/internal/server"
	"harmony/internal/server/mocks"
	. "harmony/internal/server/testing"
	ariarole "harmony/internal/testing/aria-role"
	"harmony/internal/testing/shaman"
	. "harmony/internal/testing/shaman/predicates"
	"harmony/internal/testing/shaman/sync"
	"net/http"
	"testing"

	"github.com/gost-dom/browser"
	"github.com/gost-dom/browser/dom"
	"github.com/gost-dom/browser/html"
	"github.com/samber/do"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

func init() {
	// slog.SetLogLoggerLevel(slog.LevelWarn)
	// logger.SetDefault(slog.Default())
}

func TestCanServe(t *testing.T) {
	s := server.New()
	authMock := mocks.NewAuthenticator(t)
	authMock.EXPECT().
		Authenticate(mock.Anything, mock.Anything, mock.Anything).
		Return(server.Account{}, nil).Maybe()
	s.AuthRouter.Authenticator = authMock
	b := browser.NewBrowserFromHandler(s)
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
	shaman.Scope
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

func (s *NavigateToLoginSuite) Q() shaman.Scope {
	q := shaman.NewScope(s.T(), s.win.Document())
	return q
}

func (s *NavigateToLoginSuite) SetupTest() {
	var err error
	// serv := server.New()
	authMock := mocks.NewAuthenticator(s.T())
	authMock.EXPECT().
		Authenticate(mock.Anything, mock.Anything, mock.Anything).
		Return(server.Account{}, nil).Maybe()

	injector := server.Injector.Clone()
	do.OverrideValue[server.Authenticator](injector, authMock)
	serv := do.MustInvoke[*server.Server](injector)
	b := browser.NewBrowserFromHandler(serv)
	s.win, err = b.Open("http://localhost:1234/")
	s.Scope = shaman.NewScope(s.T(), s.win.Document())
	s.Sync = sync.SetupEventSync(s.win)
	s.Scope.Container = s.win.Document()

	s.Sync.WaitFor("htmx:load")
	s.NoError(err)
}

func (s *NavigateToLoginSuite) TestClickLoginLink() {
	s.Q().Get(ByRole(ariarole.Link), ByName("Login")).Click()
	s.Equal("/auth/login", s.win.Location().Pathname())
	mainHeading := getMainHeading(s.T(), s.win)
	s.Equal("Login", mainHeading.TextContent())
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

	loginForm := NewLoginForm(s.Scope)
	loginForm.Email().SetAttribute("value", "valid-user@example.com")
	loginForm.Password().SetAttribute("value", "s3cret")
	loginForm.SubmitBtn().Click()
	s.Sync.WaitFor("htmx:afterSettle")

	s.Equal("/host", s.win.Location().Pathname())
}

func TestNavigateToLogin(t *testing.T) {
	suite.Run(t, new(NavigateToLoginSuite))
}
