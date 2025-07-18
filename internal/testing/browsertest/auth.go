package browsertest

import (
	"testing"

	"github.com/gost-dom/browser/html"
	"github.com/stretchr/testify/assert"
)

// Authenticated returns whether the current use has been authenticated.
// Marks the test as Fatal if the state couldn't be determined.
func Authenticated(t testing.TB, win html.Window) bool {
	hdr := NewPage(t, win).Header()
	loginVisible := hdr.LoginBtn() != nil
	logoutVisible := hdr.LogoutBtn() != nil
	bothVisible := loginVisible && logoutVisible
	noneVisible := !loginVisible && !logoutVisible
	if !assert.False(t, noneVisible, "Neither login, nor logout button is visible") ||
		!assert.False(t, bothVisible, "Both login and logout button is visible") {
		t.Fatal("Browser is in a bad state")
	}
	return logoutVisible && !loginVisible // Testing both is unnecessary, but communicates intent.
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
