package server_test

import (
	"fmt"
	"harmony/server"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/stroiman/go-dom/browser"
	"github.com/stroiman/go-dom/browser/dom"
	"github.com/stroiman/go-dom/browser/html"
)

func TestCanServe(t *testing.T) {
	b := browser.NewBrowserFromHandler(server.New())
	w, err := b.Open("/")
	if err != nil {
		t.Fatal(err)
	}
	h1, err := w.Document().Body().QuerySelector("h1")
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "Project harmony", h1.TextContent())
}

type NavigateToLoginSuite struct {
	suite.Suite
	win html.Window
}

func (s *NavigateToLoginSuite) SetupTest() {
	var err error
	b := browser.NewBrowserFromHandler(server.New())
	s.win, err = b.Open("/")
	assert.NoError(s.T(), err)
}

func (s *NavigateToLoginSuite) TestClickLoginLink() {
	as, err := s.win.Document().QuerySelectorAll("a")
	assert.NoError(s.T(), err)
	for _, a := range as.All() {
		fmt.Println("Found link", a.TextContent())
		if a.TextContent() == "Login" {
			fmt.Println("Click!!", a.(dom.Element).OuterHTML())
			a.(html.HTMLElement).Click()
			break
		}
	}
	assert.Equal(s.T(), "/auth/login", s.win.Location().Pathname())
}

func TestNavigateToLogin(t *testing.T) {
	suite.Run(t, new(NavigateToLoginSuite))
}
