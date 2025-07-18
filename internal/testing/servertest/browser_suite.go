package servertest

import (
	"context"
	"harmony/internal/features/auth/authrouter"
	"harmony/internal/server"
	"harmony/internal/testing/browsertest"
	"harmony/internal/testing/domaintest"
	"harmony/internal/testing/htest"
	"harmony/internal/testing/mocks/features/auth/authrouter_mock"
	"harmony/internal/testing/shaman"
	"log/slog"
	"net/http/cookiejar"
	"testing"
	"time"

	"github.com/gost-dom/browser"
	"github.com/gost-dom/browser/html"
	"github.com/gost-dom/browser/testing/gosttest"
	"github.com/gost-dom/surgeon"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// CookieJar wraps cookiejar.Jar to provide additional functionality
// for testing, particularly for CSRF protection scenarios where
// browser session state needs to be manipulated.
type CookieJar struct {
	*cookiejar.Jar
}

// Clear all cookies. Useful for testing CSRF handling.
func (j *CookieJar) Clear() {
	var err error
	j.Jar, err = cookiejar.New(nil)
	if err != nil {
		// This should never happen. cookiejar.New(...) will always return a nil
		// error (in it's current implementation). And a nil value will still
		panic("Unexpected error creating a cookie jar")
	}
}

func NewCookieJar() *CookieJar {
	var result CookieJar
	result.Clear()
	return &result
}

// BrowserSuite is an extension to [suite.Suite] that adds common behaviour for
// interacting with the website using a browser.
//
// Important: The intended use case is verification of behaviour of pages in the
// application, and only _one window_ should be used. Opening a second window in
// the same test case will panic.
type BrowserSuite struct {
	htest.GomegaSuite
	shaman.Scope
	CookieJar  *CookieJar
	Graph      *surgeon.Graph[*server.Server]
	Browser    *browser.Browser
	Win        html.Window
	Ctx        context.Context
	CancelCtx  context.CancelFunc
	logHandler *TestingLogHandler
}

func (s *BrowserSuite) SetupTest() {
	s.Graph = Graph
	s.Ctx, s.CancelCtx = context.WithTimeout(s.T().Context(), time.Millisecond*100)
	s.logHandler = &TestingLogHandler{TB: s.T()}
}

func (s *BrowserSuite) TearDownTest() {
	s.CancelCtx()
	s.Win = nil
	s.Browser = nil
	s.CookieJar = nil
}

type TestingLogHandler struct {
	testing.TB
	allowErrors bool
}

func (l TestingLogHandler) Enabled(_ context.Context, lvl slog.Level) bool {
	return lvl >= slog.LevelInfo
}
func (l TestingLogHandler) Handle(_ context.Context, r slog.Record) error {
	l.TB.Helper()
	if r.Level < slog.LevelError || l.allowErrors {
		l.TB.Logf("%v: %s", r.Level, r.Message)
	} else {
		l.TB.Errorf("%v: %s", r.Level, r.Message)
	}
	return nil
}

func (l TestingLogHandler) WithAttrs(attrs []slog.Attr) slog.Handler { return l }
func (l TestingLogHandler) WithGroup(name string) slog.Handler       { return l }

// AllowErrorLogs will allow Error log levels without automatically failing a
// test.
func (s *BrowserSuite) AllowErrorLogs() {
	s.logHandler.allowErrors = true
}

func InitAuthenticatedWindow(t testing.TB, g ServerGraph) html.Window {
	authMock := authrouter_mock.NewMockAuthenticator(t)
	g = surgeon.Replace[authrouter.Authenticator](g, authMock)
	acc := domaintest.InitAuthenticatedAccount()
	authMock.EXPECT().
		Authenticate(mock.Anything, mock.Anything, mock.Anything).
		Return(acc, nil)
	b := InitBrowser(t, g)
	win, err := b.Open("https://example.com/auth/login")
	assert.NoError(t, err, "error opening login page")

	lp := browsertest.NewPage(t, win).AssertLoginPage()
	form := lp.LoginForm()
	form.Email().Write("valid@example.com")
	form.Password().Write("validpassword")
	form.SubmitBtn().Click()
	return win
}

// InitBrowser creates a new Gost-DOM browser connected to the HTTP server
// obtained from a [surgeon.Graph].
func InitBrowser(t testing.TB, g ServerGraph) *browser.Browser {
	l := gosttest.NewTestingLogger(t, gosttest.AllowErrors())
	b := browser.New(
		browser.WithHandler(g.Instance()),
		browser.WithLogger(l),
	)
	return b
}

func (s *BrowserSuite) OpenWindow(path string) html.Window {
	if s.Win != nil {
		panic("BrowserSuite: This suite does not support opening multiple windows pr. test case")
	}
	s.Browser = InitBrowser(s.T(), s.Graph)
	s.CookieJar = NewCookieJar()
	s.Browser.Client.Jar = s.CookieJar

	win, err := s.Browser.Open(path)
	s.Assert().NoError(err)
	s.Win = win
	s.Scope = shaman.WindowScope(s.T(), s.Win)
	return win
}
