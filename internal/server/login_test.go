package server_test

import (
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
	win       html.Window
	events    chan dom.Event
	loginForm LoginForm
}

func (s *LoginPageSuite) SetupTest() {
	s.Gomega = gomega.NewWithT(s.T())
	s.events = make(chan dom.Event, 100)
	b := browser.NewBrowserFromHandler(server.New())
	win, err := b.Open("/auth/login")
	win.SetCatchAllHandler(dom.NewEventHandlerFunc(func(e dom.Event) error {
		select {
		case s.events <- e:
		default:
			panic("Event buffer full")
		}
		return nil
	}))
	s.NoError(err)
	s.win = win
	s.QueryHelper = shaman.NewQueryHelper(s.T())
	s.QueryHelper.Container = win.Document()
	s.WaitFor("htmx:load")
	s.loginForm = LoginForm{s.Scope(shaman.ByRole(ariarole.Form))}
}

func (s *LoginPageSuite) WaitFor(event string) {
	for e := range s.events {
		if e.Type() == event {
			return
		}
	}
}

func (s *LoginPageSuite) TestMissingUsername() {}
func (s *LoginPageSuite) TestMissingPassword() {}

func (s *LoginPageSuite) TestValidCredentialsRedirects() {
	s.loginForm.Email().SetAttribute("value", "valid-user@example.com")
	s.loginForm.Password().SetAttribute("value", "s3cret")
	s.loginForm.SubmitBtn().Click()

	s.Equal("/host", s.win.Location().Pathname())
}

func (s *LoginPageSuite) TestInvalidCredentials() {
	s.loginForm.Email().SetAttribute("value", "bad-user@example.com")
	s.loginForm.Password().SetAttribute("value", "s3cret")
	s.loginForm.SubmitBtn().Click()

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
