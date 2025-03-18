//go:generate mockery --all --srcpkg harmony/internal --recursive=true --with-expecter=true
package authrouter_test

import (
	"errors"
	"testing"

	"harmony/internal/features/auth"
	router "harmony/internal/features/auth/authrouter"
	. "harmony/internal/server/testing"
	ariarole "harmony/internal/testing/aria-role"
	. "harmony/internal/testing/mocks/features/auth/authrouter_mock"
	"harmony/internal/testing/servertest"
	"harmony/internal/testing/shaman"
	. "harmony/internal/testing/shaman/predicates"

	matchers "github.com/gost-dom/browser/testing/gomega-matchers"
	"github.com/gost-dom/surgeon"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type LoginPageSuite struct {
	servertest.BrowserSuite
	loginForm LoginForm
	authMock  *MockAuthenticator
}

func (s *LoginPageSuite) SetupTest() {
	s.BrowserSuite.SetupTest()
	s.authMock = NewMockAuthenticator(s.T())
	s.Graph = surgeon.Replace[router.Authenticator](s.Graph, s.authMock)
	s.OpenWindow("/auth/login")
	s.loginForm = NewLoginForm(s.Scope)
}

func (s *LoginPageSuite) TestMissingUsername() {
	s.authMock.EXPECT().
		Authenticate(mock.Anything, "", "s3cret").
		Return(auth.Account{}, auth.ErrBadCredentials)

	s.loginForm.Password().SetAttribute("value", "s3cret")
	s.loginForm.SubmitBtn().Click()

	s.Equal("/auth/login", s.Win.Location().Href())

	s.Expect(s.loginForm.Email()).To(matchers.HaveAttribute("aria-invalid", "true"))
	s.Expect(s.loginForm.Password()).ToNot(matchers.HaveAttribute("aria-invalid", "true"))
}

func (s *LoginPageSuite) TestMissingPassword() {
	s.authMock.EXPECT().
		Authenticate(mock.Anything, "valid-user@example.com", "").
		Return(auth.Account{}, auth.ErrBadCredentials)
	s.loginForm.Email().SetAttribute("value", "valid-user@example.com")
	s.loginForm.SubmitBtn().Click()

	s.Equal("/auth/login", s.Win.Location().Href())

	s.Expect(s.loginForm.Email()).ToNot(matchers.HaveAttribute("aria-invalid", "true"))
	s.Expect(s.loginForm.Password()).To(matchers.HaveAttribute("aria-invalid", "true"))
	s.Equal("Password is required", shaman.GetDescription(s.loginForm.Password()))
}

func (s *LoginPageSuite) TestValidCredentialsRedirects() {
	s.authMock.EXPECT().
		Authenticate(mock.Anything, "valid-user@example.com", "s3cret").
		Return(auth.Account{}, nil).Once()
	s.loginForm.Email().SetAttribute("value", "valid-user@example.com")
	s.loginForm.Password().SetAttribute("value", "s3cret")
	s.loginForm.SubmitBtn().Click()

	s.Equal("/", s.Win.Location().Pathname())
}

func (s *LoginPageSuite) TestInvalidCredentials() {
	s.authMock.EXPECT().
		Authenticate(mock.Anything, "bad-user@example.com", "s3cret").
		Return(auth.Account{}, auth.ErrBadCredentials).Once()
	s.loginForm.Email().SetAttribute("value", "bad-user@example.com")
	s.loginForm.Password().SetAttribute("value", "s3cret")
	s.loginForm.SubmitBtn().Click()

	s.Equal("/auth/login", s.Win.Location().Href())

	alert := s.Get(ByRole(ariarole.Alert))

	s.Assert().Equal("Email or password did not match", alert.TextContent())

	s.Expect(s.Win.Document().ActiveElement()).To(matchers.HaveAttribute("id", "email"))
}

func (s *LoginPageSuite) TestUnexpectedError() {
	s.authMock.EXPECT().
		Authenticate(mock.Anything, mock.Anything, mock.Anything).
		Return(auth.Account{}, errors.New("Unexpected")).Once()

	s.loginForm.Email().SetAttribute("value", "valid-user@example.com")
	s.loginForm.Password().SetAttribute("value", "s3cret")
	s.loginForm.SubmitBtn().Click()

	s.Equal("/auth/login", s.Win.Location().Href())

	alert := s.Get(ByRole(ariarole.Alert))

	s.Assert().NotContains(alert.TextContent(), "Email or password did not match")
	s.Assert().Contains(alert.TextContent(), "unexpected error")
	s.Expect(s.Win.Document().ActiveElement()).To(matchers.HaveAttribute("id", "email"))
}

func TestLoginPage(t *testing.T) {
	suite.Run(t, new(LoginPageSuite))
}
