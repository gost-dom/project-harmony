package server_test

import (
	"fmt"
	"harmony/server"
	"testing"

	"github.com/stroiman/go-dom/browser"
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
	fmt.Print(h1.OuterHTML())
	t.Fatal("Foo")
}
