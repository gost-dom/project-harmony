//go:generate mockery --all --with-expecter=true
package server_test

import (
	"fmt"
	"testing"

	"harmony/internal/server"
	"harmony/internal/server/mocks"
	. "harmony/internal/server/testing"
	ariarole "harmony/internal/testing/aria-role"
	"harmony/internal/testing/shaman"
	. "harmony/internal/testing/shaman/predicates"
	"harmony/internal/testing/shaman/sync"

	"github.com/gost-dom/browser"
	"github.com/gost-dom/browser/dom"
	"github.com/gost-dom/browser/html"
	matchers "github.com/gost-dom/browser/testing/gomega-matchers"
	"github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type LoginPageSuite struct {
	suite.Suite
	gomega.Gomega
	shaman.Scope
	sync.EventSync
	win       html.Window
	events    chan dom.Event
	loginForm LoginForm
	authMock  *mocks.Authenticator
}

func (s *LoginPageSuite) SetupTest() {
	s.Gomega = gomega.NewWithT(s.T())
	s.events = make(chan dom.Event, 100)
	serv := server.New()
	s.authMock = mocks.NewAuthenticator(s.T())
	serv.Authenticator = s.authMock
	b := browser.NewBrowserFromHandler(serv)
	win, err := b.Open("/auth/login")
	s.NoError(err)
	// Theoretically, this is setup too late, as DOMContentLoaded has already
	// fired by the time we get here. But in practice it works, as HTMX delays
	// processing with a setTimeout call.
	//
	// A future version of Gost will allow setting up synch _before_ opening the
	// page.
	//
	// Technically, you can create an empty browser, setup sync, and navigate. But
	// that opens a blank page, and a script context, which is a bit wasted.
	s.EventSync = sync.SetupEventSync(win)
	s.win = win
	s.Scope = shaman.NewScope(s.T(), win.Document())
	s.WaitFor("htmx:load")
	s.loginForm = NewLoginForm(s.Scope)
}

func (s *LoginPageSuite) TestMissingUsername() {
	s.authMock.EXPECT().
		Authenticate(mock.Anything, "", "s3cret").
		Return(server.Account{}, server.ErrBadCredentials)

	s.loginForm.Password().SetAttribute("value", "s3cret")
	s.loginForm.SubmitBtn().Click()
	s.WaitFor("htmx:afterSettle")

	s.Equal("/auth/login", s.win.Location().Href())

	s.Expect(s.loginForm.Email()).To(matchers.HaveAttribute("aria-invalid", "true"))
	s.Expect(s.loginForm.Password()).ToNot(matchers.HaveAttribute("aria-invalid", "true"))
}

func (s *LoginPageSuite) TestMissingPassword() {
	s.authMock.EXPECT().
		Authenticate(mock.Anything, "valid-user@example.com", "").
		Return(server.Account{}, server.ErrBadCredentials)
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
	s.authMock.EXPECT().
		Authenticate(mock.Anything, "valid-user@example.com", "s3cret").
		Return(server.Account{}, nil).Once()
	s.loginForm.Email().SetAttribute("value", "valid-user@example.com")
	s.loginForm.Password().SetAttribute("value", "s3cret")
	s.loginForm.SubmitBtn().Click()
	s.WaitFor("htmx:afterSettle")

	s.Equal("/", s.win.Location().Pathname())
}

func (s *LoginPageSuite) TestInvalidCredentials() {
	s.authMock.EXPECT().
		Authenticate(mock.Anything, "bad-user@example.com", "s3cret").
		Return(server.Account{}, server.ErrBadCredentials).Once()
	s.loginForm.Email().SetAttribute("value", "bad-user@example.com")
	s.loginForm.Password().SetAttribute("value", "s3cret")
	s.loginForm.SubmitBtn().Click()
	s.WaitFor("htmx:afterSettle")

	s.Expect(s.win.Location().Href()).To(gomega.Equal("/auth/login"))
	s.Equal("/auth/login", s.win.Location().Href())

	alert := s.Get(ByRole(ariarole.Alert))

	s.Expect(alert).To(matchers.HaveTextContent(gomega.ContainSubstring(
		"Email or password did not match")))
}

func TestLoginPage(t *testing.T) {
	suite.Run(t, new(LoginPageSuite))
}
