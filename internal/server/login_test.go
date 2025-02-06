package server_test

import (
	"fmt"
	"harmony/internal/server"
	ariarole "harmony/internal/testing/aria-role"
	"harmony/internal/testing/shaman"
	"testing"

	"github.com/gost-dom/browser"
	"github.com/gost-dom/browser/dom"
	"github.com/gost-dom/browser/html"
	matchers "github.com/gost-dom/browser/testing/gomega-matchers"
	"github.com/onsi/gomega"
	"github.com/stretchr/testify/suite"
)

type LoginPageSuite struct {
	suite.Suite
	gomega.Gomega
	shaman.QueryHelper
	Sync
	win       html.Window
	events    chan dom.Event
	loginForm LoginForm
}

// Sync is a helper to wait for the window to dispatch certain events. This
// helps syncrhonizing test with HTMX events, so as to not proceed with the
// test, until HTMX has settled.
//
// Test code can call [Sync.WaitFor] to wait for an event. Events must be waited
// for in the order they are dispatched,
//
// Good events to wait for (this is just meant as a starting point)
//   - htmx:load - when HTMX has loaded, and processed all relevant nodes
//   - htmx:afterSettle - when has swapped innerHTML
//
// Note: It appears that HTMX dispatches different events depending on the
// hx-swap value, so keep that in mind.
//
// Tip: On the fron-end, you can call `htmx.logAll()` to see which events it
// emits.
type Sync struct {
	events chan dom.Event
}

// Waits for an event with a specific [dom.Event.Type] to be dispatched. This
// must be called in the order the events are dispatched. Any event before the
// one listening for will be discarded.
//
// For HTMX events, there would normally be causality, i.e., htmx:beforeSwap is
// dispatched before htmx:afterSwap, so you can call this twice in the right
// order.
//
// If you need to wait for two events with no causality, you will need to use to
// syncers. But please write a message in the [discussions], as it's quite
// possible the sync would need some kind of fork for this scenario to work
// properly
//
// [discussions]: https://github.com/orgs/gost-dom/discussions
func (s Sync) WaitFor(type_ string) dom.Event {
	for e := range s.events {
		if e.Type() == type_ {
			return e
		}
	}
	return nil
}

func SetupSync(w html.Window) (res Sync) {
	// I'm sure there's a better way to do this, but we want to put the events
	// into a queue without blocking the sender, but a listener will block until
	// the event they are interested in
	res.events = make(chan dom.Event, 100)
	w.SetCatchAllHandler(dom.NewEventHandlerFunc(func(e dom.Event) error {
		select {
		case res.events <- e:
		default:
			panic("Event buffer full")
		}
		return nil
	}))
	return
}

func (s *LoginPageSuite) SetupTest() {
	s.Gomega = gomega.NewWithT(s.T())
	s.events = make(chan dom.Event, 100)
	b := browser.NewBrowserFromHandler(server.New())
	win, err := b.Open("/auth/login")
	// Theoretically, this is setup too late, as DOMContentLoaded has already
	// fired by the time we get here. But in practice it works, as HTMX delays
	// processing with a setTimeout call.
	//
	// A future version of Gost will allow setting up synch _before_ opening the
	// page.
	//
	// Technically, you can create an empty browser, setup sync, and navigate. But
	// that opens a blank page, and a script context, which is a bit wasted.
	s.Sync = SetupSync(win)
	s.NoError(err)
	s.win = win
	s.QueryHelper = shaman.NewQueryHelper(s.T())
	s.QueryHelper.Container = win.Document()
	s.WaitFor("htmx:load")
	s.loginForm = LoginForm{s.Scope(shaman.ByRole(ariarole.Form))}
}

func (s *LoginPageSuite) TestMissingUsername() {
	s.loginForm.Password().SetAttribute("value", "s3cret")
	s.loginForm.SubmitBtn().Click()
	s.WaitFor("htmx:afterSettle")

	s.Equal("/auth/login", s.win.Location().Href())

	s.Expect(s.loginForm.Email()).To(matchers.HaveAttribute("aria-invalid", "true"))
	s.Expect(s.loginForm.Password()).ToNot(matchers.HaveAttribute("aria-invalid", "true"))
}

func (s *LoginPageSuite) TestMissingPassword() {
	s.loginForm.Email().SetAttribute("value", "valid-user@example.com")
	s.loginForm.SubmitBtn().Click()
	s.WaitFor("htmx:afterSettle")

	s.Equal("/auth/login", s.win.Location().Href())

	s.Expect(s.loginForm.Email()).ToNot(matchers.HaveAttribute("aria-invalid", "true"))
	s.Expect(s.loginForm.Password()).To(matchers.HaveAttribute("aria-invalid", "true"))
	fmt.Println(s.loginForm.Container.(dom.Element).OuterHTML())
	s.Equal("Password is required", shaman.GetDescription(s.loginForm.Password()))
}

func (s *LoginPageSuite) TestValidCredentialsRedirects() {
	s.loginForm.Email().SetAttribute("value", "valid-user@example.com")
	s.loginForm.Password().SetAttribute("value", "s3cret")
	s.loginForm.SubmitBtn().Click()
	s.WaitFor("htmx:afterSettle")

	s.Equal("/host", s.win.Location().Pathname())
}

func (s *LoginPageSuite) TestInvalidCredentials() {
	s.loginForm.Email().SetAttribute("value", "bad-user@example.com")
	s.loginForm.Password().SetAttribute("value", "s3cret")
	s.loginForm.SubmitBtn().Click()
	s.WaitFor("htmx:afterSettle")

	s.Expect(s.win.Location().Href()).To(gomega.Equal("/auth/login"))
	s.Equal("/auth/login", s.win.Location().Href())

	alert := s.Get(shaman.ByRole(ariarole.Alert))

	s.Expect(alert).To(matchers.HaveTextContent(gomega.ContainSubstring(
		"Email or password did not match")))
}

func TestLoginPage(t *testing.T) {
	suite.Run(t, new(LoginPageSuite))
}

type LoginForm struct {
	shaman.QueryHelper
}

func (f LoginForm) Email() dom.Element {
	return f.Get(shaman.ByRole(ariarole.Textbox), shaman.ByName("Email"))
}

func (f LoginForm) Password() dom.Element {
	return f.Get(shaman.ByRole(ariarole.PasswordText), shaman.ByName("Password"))
}

func (f LoginForm) SubmitBtn() dom.Element {
	return f.Get(shaman.ByRole(ariarole.Button), shaman.ByName("Sign in"))
}
