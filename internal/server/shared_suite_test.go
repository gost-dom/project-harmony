package server_test

import (
	"context"
	"harmony/internal/server"
	"harmony/internal/testing/shaman"
	"time"

	"github.com/gost-dom/browser"
	"github.com/gost-dom/browser/dom"
	"github.com/gost-dom/browser/html"
	"github.com/gost-dom/surgeon"
	"github.com/onsi/gomega"
	"github.com/stretchr/testify/suite"
)

// GomegaSuite is a specialized [suite.Suite] that add [gomega.Gomega] assertion
// semantics to the test suite.
//
// This can provide more expressive assertions when combined with custom
// mathers.
type GomegaSuite struct {
	suite.Suite
	gomega.Gomega
}

func (s *GomegaSuite) SetupTest() {
	s.Gomega = gomega.NewWithT(s.T())
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
	GomegaSuite
	shaman.Scope
	graph     *surgeon.Graph[*server.Server]
	win       html.Window
	ctx       context.Context
	cancelCtx context.CancelFunc
}

func (s *BrowserSuite) SetupTest() {
	s.GomegaSuite.SetupTest()
	s.graph = graph
	s.ctx, s.cancelCtx = context.WithTimeout(context.Background(), time.Millisecond*100)
}

func (s *BrowserSuite) OpenWindow(path string) html.Window {
	if s.win != nil {
		panic("BrowserSuite: This suite does not support opening multiple windows pr. test case")
	}
	serv := s.graph.Instance()
	b := browser.NewBrowserFromHandler(serv)

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
	win, err := b.Open("about:blank")
	s.Assert().NoError(err)
	err = win.Navigate(path)
	s.Assert().NoError(err)
	s.win = win
	s.Scope = shaman.NewScope(s.T(), s.win.Document())
	return win
}

func (s *BrowserSuite) TearDownTest() {
	s.win = nil
	s.cancelCtx()
}
