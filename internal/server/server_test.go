package server_test

import (
	"harmony/internal/server"
	"harmony/internal/server/mocks"
	. "harmony/internal/server/testing"
	ariarole "harmony/internal/testing/aria-role"
	. "harmony/internal/testing/shaman/predicates"
	"testing"

	"github.com/gost-dom/surgeon"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

func init() {
	// slog.SetLogLoggerLevel(slog.LevelWarn)
	// logger.SetDefault(slog.Default())
	graph = surgeon.BuildGraph(server.New(), surgeon.PackagePrefixScope("harmony"))
}

var graph *surgeon.Graph[*server.Server]

type NavigateToLoginSuite struct{ BrowserSuite }

func (s *NavigateToLoginSuite) SetupTest() {
	s.BrowserSuite.SetupTest()
	// In this scenario, authentication always succeed. Specific tests for the
	// login page exercise different aspects
	authMock := mocks.NewAuthenticator(s.T())
	authMock.EXPECT().
		Authenticate(mock.Anything, mock.Anything, mock.Anything).
		Return(server.Account{}, nil).Maybe()
	s.graph = surgeon.Replace[server.Authenticator](s.graph, authMock)

	s.OpenWindow("http://localhost:1234/")
	s.WaitFor("htmx:load")
	s.win.Clock().RunAll()
}

func (s *NavigateToLoginSuite) TestLoginFlow() {
	s.Get(ByRole(ariarole.Link), ByName("Go to hosting")).Click()
	s.win.Clock().RunAll()

	s.Equal("/auth/login", s.win.Location().Pathname(), "Location after host")
	mainHeading := s.Get(ByH1)
	s.Equal("Login", mainHeading.TextContent())

	loginForm := NewLoginForm(s.Scope)
	loginForm.Email().SetAttribute("value", "valid-user@example.com")
	loginForm.Password().SetAttribute("value", "s3cret")
	loginForm.SubmitBtn().Click()

	s.Equal("/host", s.win.Location().Pathname())
}

func TestNavigateToLogin(t *testing.T) {
	suite.Run(t, new(NavigateToLoginSuite))
}
