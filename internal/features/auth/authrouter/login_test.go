//go:generate mockery --all --srcpkg harmony/internal --recursive=true --with-expecter=true
package authrouter_test

import (
	"errors"
	"testing"

	"harmony/internal/features/auth"
	"harmony/internal/features/auth/authdomain"
	"harmony/internal/features/auth/authdomain/password"
	router "harmony/internal/features/auth/authrouter"
	. "harmony/internal/server/testing"
	ariarole "harmony/internal/testing/aria-role"
	. "harmony/internal/testing/domaintest"
	. "harmony/internal/testing/mocks/features/auth/authrouter_mock"
	"harmony/internal/testing/servertest"
	"harmony/internal/testing/shaman"
	. "harmony/internal/testing/shaman/predicates"

	matchers "github.com/gost-dom/browser/testing/gomega-matchers"
	"github.com/gost-dom/surgeon"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

func matchPassword(pw string) any {
	return mock.MatchedBy(password.Parse(pw).Equals)
}

type LoginPageSuite struct {
	servertest.BrowserSuite
	loginForm LoginForm
	authMock  *MockAuthenticator
}

func (s *LoginPageSuite) SetupTest() {
	s.BrowserSuite.SetupTest()
	s.authMock = NewMockAuthenticator(s.T())
	s.Graph = surgeon.Replace[router.Authenticator](s.Graph, s.authMock)
	s.OpenWindow("https://example.com/auth/login")
	s.loginForm = NewLoginForm(s.Scope)
}

func (s *LoginPageSuite) TestMissingUsername() {
	s.authMock.EXPECT().
		Authenticate(mock.Anything, "", matchPassword("s3cret")).
		Return(authdomain.AuthenticatedAccount{}, auth.ErrBadCredentials)

	s.loginForm.Password().SetAttribute("value", "s3cret")
	s.loginForm.SubmitBtn().Click()

	s.Equal("/auth/login", s.Win.Location().Pathname())

	s.Expect(s.loginForm.Email()).To(matchers.HaveAttribute("aria-invalid", "true"))
	s.Expect(s.loginForm.Password()).ToNot(matchers.HaveAttribute("aria-invalid", "true"))
}

func (s *LoginPageSuite) TestMissingPassword() {
	s.authMock.EXPECT().
		Authenticate(mock.Anything, "valid-user@example.com", matchPassword("")).
		Return(authdomain.AuthenticatedAccount{}, auth.ErrBadCredentials)
	s.loginForm.Email().SetAttribute("value", "valid-user@example.com")
	s.loginForm.SubmitBtn().Click()

	s.Equal("/auth/login", s.Win.Location().Pathname())

	s.Expect(s.loginForm.Email()).ToNot(matchers.HaveAttribute("aria-invalid", "true"))
	s.Expect(s.loginForm.Password()).To(matchers.HaveAttribute("aria-invalid", "true"))
	s.Equal("Password is required", shaman.GetDescription(s.loginForm.Password()))
}

func (s *LoginPageSuite) TestValidCredentialsRedirects() {
	s.authMock.EXPECT().
		Authenticate(mock.Anything, "valid-user@example.com", matchPassword("s3cret")).
		Return(InitAuthenticatedAccount(), nil).Once()
	s.loginForm.Email().SetAttribute("value", "valid-user@example.com")
	s.loginForm.Password().SetAttribute("value", "s3cret")
	s.loginForm.SubmitBtn().Click()

	s.Equal("/", s.Win.Location().Pathname())
}

func (s *LoginPageSuite) TestCSRFHandling() {
	s.authMock.EXPECT().
		Authenticate(mock.Anything, "valid-user@example.com", matchPassword("s3cret")).
		Return(InitAuthenticatedAccount(), nil).Maybe()

	s.CookieJar.Clear()
	s.loginForm.Email().SetAttribute("value", "valid-user@example.com")
	s.loginForm.Password().SetAttribute("value", "s3cret")
	s.loginForm.SubmitBtn().Click()

	s.Equal("/auth/login", s.Win.Location().Pathname())
}

func (s *LoginPageSuite) TestCSRFWithMultipleWindows() {
	s.authMock.EXPECT().
		Authenticate(mock.Anything, "valid-user@example.com", matchPassword("s3cret")).
		Return(InitAuthenticatedAccount(), nil).Once()

	_, err := s.Browser.Open("https://example.com/")
	s.Assert().NoError(err)

	s.loginForm.Email().SetAttribute("value", "valid-user@example.com")
	s.loginForm.Password().SetAttribute("value", "s3cret")
	s.loginForm.SubmitBtn().Click()

	s.Equal("/", s.Win.Location().Pathname())
}

func (s *LoginPageSuite) TestInvalidCredentials() {
	s.authMock.EXPECT().
		Authenticate(mock.Anything, "bad-user@example.com", matchPassword("s3cret")).
		Return(authdomain.AuthenticatedAccount{}, auth.ErrBadCredentials).Once()
	s.loginForm.Email().SetAttribute("value", "bad-user@example.com")
	s.loginForm.Password().SetAttribute("value", "s3cret")
	s.loginForm.SubmitBtn().Click()

	s.Equal("/auth/login", s.Win.Location().Pathname())

	alert := s.Get(ByRole(ariarole.Alert))

	s.Assert().Equal("Email or password did not match", alert.TextContent())

	s.Expect(s.Win.Document().ActiveElement()).To(matchers.HaveAttribute("id", "email"))
}

func (s *LoginPageSuite) TestUnexpectedError() {
	s.authMock.EXPECT().
		Authenticate(mock.Anything, mock.Anything, mock.Anything).
		Return(authdomain.AuthenticatedAccount{}, errors.New("Unexpected")).Once()

	s.loginForm.Email().SetAttribute("value", "valid-user@example.com")
	s.loginForm.Password().SetAttribute("value", "s3cret")
	s.loginForm.SubmitBtn().Click()

	s.Equal("/auth/login", s.Win.Location().Pathname())

	alert := s.Get(ByRole(ariarole.Alert))

	s.Assert().NotContains(alert.TextContent(), "Email or password did not match")
	s.Assert().Contains(alert.TextContent(), "unexpected error")
	s.Expect(s.Win.Document().ActiveElement()).To(matchers.HaveAttribute("id", "email"))
}

func TestLoginPage(t *testing.T) {
	suite.Run(t, new(LoginPageSuite))
}
