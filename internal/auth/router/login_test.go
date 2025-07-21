//go:generate mockery --all --srcpkg harmony/internal --recursive=true --with-expecter=true
package router_test

import (
	"errors"
	"testing"

	"harmony/internal/auth"
	"harmony/internal/auth/authdomain"
	"harmony/internal/auth/authdomain/password"
	"harmony/internal/auth/router"
	. "harmony/internal/testing/browsertest"
	. "harmony/internal/testing/domaintest"
	"harmony/internal/testing/mocks/auth/router_mock"
	"harmony/internal/testing/servertest"

	. "github.com/gost-dom/browser/testing/gomega-matchers"
	"github.com/gost-dom/shaman"
	"github.com/gost-dom/shaman/ariarole"
	. "github.com/gost-dom/shaman/predicates"
	"github.com/gost-dom/surgeon"
	"github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

func matchPassword(pw string) any {
	return mock.MatchedBy(password.Parse(pw).Equals)
}

type LoginPageSuite struct {
	servertest.BrowserSuite
	loginForm LoginForm
	authMock  *router_mock.MockAuthenticator
}

func (s *LoginPageSuite) SetupTest() {
	s.BrowserSuite.SetupTest()
	s.authMock = router_mock.NewMockAuthenticator(s.T())
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

	s.Expect(s.loginForm.Email()).To(HaveAttribute("aria-invalid", "true"))
	s.Expect(s.loginForm.Password()).ToNot(HaveAttribute("aria-invalid", "true"))
}

func (s *LoginPageSuite) TestMissingPassword() {
	s.authMock.EXPECT().
		Authenticate(mock.Anything, "valid-user@example.com", matchPassword("")).
		Return(authdomain.AuthenticatedAccount{}, auth.ErrBadCredentials)
	s.loginForm.Email().SetAttribute("value", "valid-user@example.com")
	s.loginForm.SubmitBtn().Click()

	s.Equal("/auth/login", s.Win.Location().Pathname())

	s.Expect(s.loginForm.Email()).ToNot(HaveAttribute("aria-invalid", "true"))
	s.Expect(s.loginForm.Password()).To(HaveAttribute("aria-invalid", "true"))
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
	s.AllowErrorLogs()
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
	s.AllowErrorLogs()
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

	s.Expect(s.Win.Document().ActiveElement()).To(HaveAttribute("id", "email"))
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

	s.Expect(alert).
		ToNot(HaveTextContent(gomega.ContainSubstring("Email or password did not match")))
	s.Expect(alert).To(HaveTextContent(gomega.ContainSubstring("unexpected error")))
	s.Expect(s.Win.Document().ActiveElement()).To(HaveAttribute("id", "email"))
}

func TestLoginPage(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(LoginPageSuite))
}
