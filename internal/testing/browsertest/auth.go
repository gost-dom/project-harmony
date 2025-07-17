package browsertest

import (
	"harmony/internal/testing/shaman"
	"harmony/internal/testing/shaman/ariarole"
	. "harmony/internal/testing/shaman/predicates"
	"testing"

	"github.com/gost-dom/browser/html"
	"github.com/stretchr/testify/assert"
)

// Authenticated returns whether the current use has been authenticated.
// Marks the test as Fatal if the state couldn't be determined.
func Authenticated(t testing.TB, win html.Window) bool {
	win.History().PushState(html.EMPTY_STATE, "/")
	header := shaman.WindowScope(t, win).Subscope(shaman.ByRole(ariarole.Banner))
	loginBtn := header.Find(ByRole(ariarole.Link), ByName("Login"))
	logoutBtn := header.Find(ByRole(ariarole.Button), ByName("Logout"))
	if !assert.False(t,
		loginBtn == nil && logoutBtn == nil,
		"Neither login, nor logout button is visible",
	) {
		t.Fatal("Browser is in a bad state")
	}
	if !assert.False(t,
		loginBtn != nil && logoutBtn != nil,
		"Both login and logout button is visible",
	) {
		t.Fatal("Browser is in a bad state")
	}
	return logoutBtn != nil && loginBtn == nil
}

// AssertAuthenticated asserts that the user has been authenticated.
//
// Marks the test as Fatal if the state couldn't be determined.
func AssertAuthenticated(t testing.TB, win html.Window) bool {
	return assert.True(t, Authenticated(t, win), "User has been authenticated")
}

// AssertUnauthenticated asserts that the user has not been authenticated (or
// logged out again)
//
// Marks the test as Fatal if the state couldn't be determined.
func AssertUnauthenticated(t testing.TB, win html.Window) bool {
	return assert.False(t, Authenticated(t, win), "User has been authenticated")
}
