package server_test

import (
	"harmony/internal/features/auth/authdomain"
	"harmony/internal/features/auth/authrouter"
	. "harmony/internal/server/testing"
	ariarole "harmony/internal/testing/aria-role"
	. "harmony/internal/testing/mocks/features/auth/authrouter_mock"
	"harmony/internal/testing/servertest"
	. "harmony/internal/testing/shaman/predicates"
	"testing"

	"github.com/gost-dom/surgeon"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type NavigateToLoginSuite struct{ servertest.BrowserSuite }

func (s *NavigateToLoginSuite) SetupTest() {
	s.BrowserSuite.SetupTest()
	// In this scenario, authentication always succeed. Specific tests for the
	// login page exercise different aspects
	authMock := NewMockAuthenticator(s.T())
	authMock.EXPECT().
		Authenticate(mock.Anything, mock.Anything, mock.Anything).
		Return(authdomain.Account{}, nil).Maybe()
	s.Graph = surgeon.Replace[authrouter.Authenticator](s.Graph, authMock)

	s.OpenWindow("https://example.com/")
	s.Win.Clock().RunAll()
}

func (s *NavigateToLoginSuite) TestLoginFlow() {
	s.Get(ByRole(ariarole.Link), ByName("Go to hosting")).Click()
	s.Win.Clock().RunAll()

	s.Equal("/auth/login", s.Win.Location().Pathname(), "Location after host")
	mainHeading := s.Get(ByH1)
	s.Equal("Login", mainHeading.TextContent())

	loginForm := NewLoginForm(s.Scope)
	loginForm.Email().SetAttribute("value", "valid-user@example.com")
	loginForm.Password().SetAttribute("value", "s3cret")
	loginForm.SubmitBtn().Click()

	s.Equal("/host", s.Win.Location().Pathname())
}

func TestNavigateToLogin(t *testing.T) {
	suite.Run(t, new(NavigateToLoginSuite))
}
