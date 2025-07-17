package authrouter_test

import (
	"log/slog"
	"testing"

	router "harmony/internal/features/auth/authrouter"
	. "harmony/internal/server/testing"
	"harmony/internal/testing/domaintest"
	. "harmony/internal/testing/mocks/features/auth/authrouter_mock"
	"harmony/internal/testing/servertest"
	"harmony/internal/testing/shaman"
	"harmony/internal/testing/shaman/ariarole"

	. "harmony/internal/testing/shaman/predicates"

	"github.com/gost-dom/browser/testing/gosttest"
	"github.com/gost-dom/surgeon"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func useTestLogger(t testing.TB) (cleanup func()) {
	tmp := slog.Default()

	l := gosttest.NewTestingLogger(t, gosttest.AllowErrors())
	slog.SetDefault(l)

	return func() { slog.SetDefault(tmp) }
}

func TestPOSTLogout(t *testing.T) {
	cleanup := useTestLogger(t)
	defer cleanup()

	authMock := NewMockAuthenticator(t)
	g := surgeon.Replace[router.Authenticator](servertest.Graph, authMock)
	acc := domaintest.InitAuthenticatedAccount()
	authMock.EXPECT().
		Authenticate(mock.Anything, mock.Anything, mock.Anything).
		Return(acc, nil)
	b := servertest.InitBrowser(t, g)
	win, err := b.Open("https://example.com/auth/login")
	assert.NoError(t, err, "error opening login page")
	mainScope := shaman.WindowScope(t, win)
	form := NewLoginForm(mainScope)
	form.Email().SetAttribute("value", "valid@example.com")
	form.Password().SetAttribute("value", "validpassword")
	form.SubmitBtn().Click()

	header := shaman.WindowScope(t, win).Subscope(ByRole(ariarole.Banner))
	logoutBtn := header.Get(shaman.ByRole(ariarole.Button), shaman.ByName("Logout"))
	logoutBtn.Click()

	assert.NoError(t, win.Navigate("https://example.com/host"))
	header = shaman.WindowScope(t, win).Subscope(shaman.ByRole(ariarole.Banner))
	loginBtn := header.Find(ByRole(ariarole.Link), ByName("Login"))
	logoutBtn = header.Find(ByRole(ariarole.Button), ByName("Logout"))
	userLoggedIn := logoutBtn != nil && loginBtn == nil

	assert.False(t, userLoggedIn, "user logged in after clicking logout")
}
