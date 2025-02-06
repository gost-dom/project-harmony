package sync

import (
	"github.com/gost-dom/browser/dom"
	"github.com/gost-dom/browser/html"
)

// EventSync is a helper to wait for the window to dispatch certain events. This
// helps syncrhonizing test with HTMX events, so as to not proceed with the
// test, until HTMX has settled.
//
// Test code can call [EventSync.WaitFor] to wait for an event. Events must be waited
// for in the order they are dispatched,
//
// Good events to wait for (this is just meant as a starting point)
//   - htmx:load - when HTMX has loaded, and processed all relevant nodes
//   - htmx:afterSettle - when has swapped innerHTML
//
// Note: It appears that HTMX dispatches different events depending on the
// hx-swap value, so keep that in mind.
//
// Tip: On the fron-end, you can call `htmx.logAll()` to see which events it
// emits.
type EventSync struct {
	events chan dom.Event
}

// Creates a new Sync, and listen to events from an [html.Window].
func SetupEventSync(w html.Window) (res EventSync) {
	// I'm sure there's a better way to do this, but we want to put the events
	// into a queue without blocking the sender, but a listener will block until
	// the event they are interested in
	res.events = make(chan dom.Event, 100)
	w.SetCatchAllHandler(dom.NewEventHandlerFunc(func(e dom.Event) error {
		select {
		case res.events <- e:
		default:
			panic("Event buffer full")
		}
		return nil
	}))
	return
}

// Waits for an event with a specific [dom.Event.Type] to be dispatched. This
// must be called in the order the events are dispatched. Any event before the
// one listening for will be discarded.
//
// For HTMX events, there would normally be causality, i.e., htmx:beforeSwap is
// dispatched before htmx:afterSwap, so you can call this twice in the right
// order.
//
// If you need to wait for two events with no causality, you will need to use to
// syncers. But please write a message in the [discussions], as it's quite
// possible the sync would need some kind of fork for this scenario to work
// properly
//
// [discussions]: https://github.com/orgs/gost-dom/discussions
func (s EventSync) WaitFor(type_ string) dom.Event {
	for e := range s.events {
		if e.Type() == type_ {
			return e
		}
	}
	return nil
}
