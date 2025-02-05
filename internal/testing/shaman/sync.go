package shaman

import (
	"context"
	"testing"

	"github.com/gost-dom/browser/dom"
)

type Sync struct {
	Target dom.EventTarget
	t      *testing.T
	ctx    context.Context
}

func NewSync(ctx context.Context, t *testing.T) Sync {
	return Sync{t: t, ctx: ctx}
}

// WaitFor executes function f, and waits for the window w dispatches an event of
// type e,
func (s Sync) WaitFor(e string, f func()) {
	c := make(chan struct{})
	l := dom.NewEventHandlerFunc(func(dom.Event) error {
		go func() { close(c) }()
		return nil
	})
	s.Target.AddEventListener(e, l)
	defer s.Target.RemoveEventListener(e, l)
	f()
	select {
	case <-c:
		return
	case <-s.ctx.Done():
		s.t.Fatalf("Timeout waiting for event: %s", e)
	}
}
