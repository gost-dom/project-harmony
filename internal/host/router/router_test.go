package router_test

import (
	"harmony/internal/testing/servertest"
	"strings"
	"testing"

	"github.com/gost-dom/shaman"
	"github.com/gost-dom/shaman/ariarole"
	"github.com/stretchr/testify/assert"
)

func TestHostPage(t *testing.T) {
	win := servertest.InitAuthenticatedWindow(t, servertest.Graph)
	win.Navigate("/host")
	scope := shaman.WindowScope(t, win)

	t.Run("Heading", func(t *testing.T) {
		scope.Find(
			shaman.ByRole(ariarole.Link),
			shaman.ByName("New Service"),
		).Click()
		h1 := scope.Get(shaman.ByH1)
		assert.Equal(t, "New Service", strings.TrimSpace(h1.TextContent()))
	})
}

func TestNewServicePage(t *testing.T) {
	win := servertest.InitAuthenticatedWindow(t, servertest.Graph)
	assert.NoError(t, win.Navigate("/host/services/new"))
	h1 := shaman.WindowScope(t, win).Get(shaman.ByH1)
	assert.Equal(t, "New Service", strings.TrimSpace(h1.TextContent()))
}
