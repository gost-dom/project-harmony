package servertest

import (
	"context"
	"harmony/internal/server"
	"harmony/internal/testing/htest"
	"harmony/internal/testing/shaman"
	"net/http/cookiejar"
	"time"

	"github.com/gost-dom/browser"
	"github.com/gost-dom/browser/html"
	"github.com/gost-dom/surgeon"
)

type CookieJar struct {
	*cookiejar.Jar
}

func (j *CookieJar) Clear() {
	j.Jar, _ = cookiejar.New(nil)
}

func NewCookieJar() *CookieJar {
	jar, _ := cookiejar.New(nil)
	return &CookieJar{jar}
}

// BrowserSuite is an extension to [suite.Suite] that adds common behaviour for
// interacting with the website using a browser.
//   - Cloning a [do.Injector] to allow easy mock replacement
//   - Synchronising to events
//   - Creating a browser on top of the root HTTP handler and open a page
//
// Important: The intended use case is verification of behaviour of pages in the
// application, and only _one window_ should be used. Opening a second window in
// the same test case will panic.
type BrowserSuite struct {
	htest.GomegaSuite
	shaman.Scope
	CookieJar *CookieJar
	Graph     *surgeon.Graph[*server.Server]
	Browser   *browser.Browser
	Win       html.Window
	Ctx       context.Context
	CancelCtx context.CancelFunc
}

func (s *BrowserSuite) SetupTest() {
	s.GomegaSuite.SetupTest()
	s.Graph = graph
	s.Ctx, s.CancelCtx = context.WithTimeout(context.Background(), time.Millisecond*100)
}

func (s *BrowserSuite) OpenWindow(path string) html.Window {
	if s.Win != nil {
		panic("BrowserSuite: This suite does not support opening multiple windows pr. test case")
	}
	serv := s.Graph.Instance()
	s.Browser = browser.New(browser.WithHandler(serv))
	s.CookieJar = NewCookieJar()
	s.Browser.Client.Jar = s.CookieJar

	// Opening "about:blank" is a bit of a hack to allow adding the event sync
	// _before_ loading an HTMX page; as by the time `Open` or `Navigate`
	// returns, the `DOMContentLoaded` event has already fired.
	//
	// In practice, you _can_ add this after `Navigate(path)`, as HTMX defers
	// initialization by a `setTimeout` call. But this was explicitly created
	// like this for the sake of a more stable example for other scenarios.
	//
	// A future version of Gost will support adding the event handler _before_
	// navigating, to avoid the overhead creating a new script context; although
	// it should be minimal.
	win, err := s.Browser.Open("about:blank")
	s.Assert().NoError(err)
	err = win.Navigate(path)
	s.Assert().NoError(err)
	s.Win = win
	s.Scope = shaman.NewScope(s.T(), s.Win.Document())
	return win
}

func (s *BrowserSuite) TearDownTest() {
	s.Win = nil
	s.CancelCtx()
}
