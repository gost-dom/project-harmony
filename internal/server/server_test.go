package server_test

import (
	"harmony/internal/features/auth/authrouter"
	. "harmony/internal/server/testing"
	. "harmony/internal/testing/domaintest"
	. "harmony/internal/testing/mocks/features/auth/authrouter_mock"
	"harmony/internal/testing/servertest"
	"harmony/internal/testing/shaman"
	"harmony/internal/testing/shaman/ariarole"
	. "harmony/internal/testing/shaman/predicates"
	"testing"

	"github.com/gost-dom/browser/html"
	"github.com/gost-dom/surgeon"
	"github.com/onsi/gomega"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type serverSuite struct {
	t *testing.T
	gomega.Gomega
	shaman.Scope
	Win html.Window
}

func initServerSuite(t *testing.T) serverSuite {
	authMock := NewMockAuthenticator(t)
	authMock.EXPECT().
		Authenticate(mock.Anything, mock.Anything, mock.Anything).
		Return(InitAuthenticatedAccount(), nil).Maybe()
	g := surgeon.Replace[authrouter.Authenticator](servertest.Graph, authMock)

	b := servertest.InitBrowser(t, g)

	win, err := b.Open("https://example.com/")
	assert.NoError(t, err)

	return serverSuite{t: t,
		Gomega: gomega.NewWithT(t),
		Scope:  shaman.WindowScope(t, win),
		Win:    win,
	}
}

func TestLoginFlow(t *testing.T) {
	s := initServerSuite(t)
	t.Run("Login button exists before login", func(t *testing.T) {
		header := s.Subscope(ByRole(ariarole.Banner))
		_, hasLoginButton := header.Query(ByRole(ariarole.Link), ByName("Login"))
		assert.True(t, hasLoginButton)
	})

	t.Run("/host redirects to /auth/login", func(t *testing.T) {
		s.Get(ByRole(ariarole.Link), ByName("Go to hosting")).Click()
		s.Win.Clock().RunAll()

		assert.Equal(t, "/auth/login", s.Win.Location().Pathname(), "Location after host")
		mainHeading := s.Get(ByH1)
		assert.Equal(t, "Login", mainHeading.TextContent())
	})

	t.Run("Performing a logn redirects back to /host", func(t *testing.T) {
		loginForm := NewLoginForm(s.Scope)
		loginForm.Email().SetAttribute("value", "valid-user@example.com")
		loginForm.Password().SetAttribute("value", "s3cret")
		loginForm.SubmitBtn().Click()

		assert.Equal(t, "/host", s.Win.Location().Pathname(), "path after login name")
		assert.Equal(t, "Host", s.Get(ByH1).TextContent(), "page heading after login")
	})

	t.Run("Login button disappears after login", func(t *testing.T) {
		header := s.Subscope(ByRole(ariarole.Banner))
		_, hasLoginButton := header.Query(ByRole(ariarole.Link), ByName("Login"))
		assert.False(t, hasLoginButton, "A login link exists in the header")

		_, hasLogoutButton := header.Query(ByRole(ariarole.Link), ByName("Logout"))
		assert.False(t, hasLogoutButton, "A logout button exists in the header")
	})
}

func TestOpeningHostDirectlyRedirects(t *testing.T) {
	s := initServerSuite(t)
	s.Win.Navigate("https://example.com/host")
	assert.Equal(t, "/auth/login", s.Win.Location().Pathname(), "Location after host")
}
