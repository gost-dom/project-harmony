package router_test

import (
	"testing"

	"harmony/internal/testing/browsertest"
	"harmony/internal/testing/servertest"

	"github.com/gost-dom/shaman"
	"github.com/gost-dom/shaman/ariarole"
	. "github.com/gost-dom/shaman/predicates"
)

func TestPOSTLogout(t *testing.T) {
	win := servertest.InitAuthenticatedWindow(t, servertest.Graph)

	header := shaman.WindowScope(t, win).Subscope(ByRole(ariarole.Banner))
	logoutBtn := header.Get(ByRole(ariarole.Button), ByName("Logout"))
	logoutBtn.Click()

	browsertest.AssertUnauthenticated(t, win)
}
