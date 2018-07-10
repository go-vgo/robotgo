package xevent

/*
xevent/eventloop.go contains code that implements a main X event loop.

Namely, it provides facilities to read new events into xevent's event queue,
run a normal main event loop and run a main event loop that pings a channel
each time an event is about to be dequeued. The latter facility allows one to
easily include other input sources for processing in a program's main event
loop.
*/

import (
	"github.com/BurntSushi/xgb/shape"
	"github.com/BurntSushi/xgb/xproto"

	"github.com/BurntSushi/xgbutil"
)

// Read reads one or more events and queues them in XUtil.
// If 'block' is True, then call 'WaitForEvent' before sucking up
// all events that have been queued by XGB.
func Read(xu *xgbutil.XUtil, block bool) {
	if block {
		ev, err := xu.Conn().WaitForEvent()
		if ev == nil && err == nil {
			xgbutil.Logger.Fatal("BUG: Could not read an event or an error.")
		}
		Enqueue(xu, ev, err)
	}

	// Clean up anything that's in the queue
	for {
		ev, err := xu.Conn().PollForEvent()

		// No events left...
		if ev == nil && err == nil {
			break
		}

		// We're good, queue it up
		Enqueue(xu, ev, err)
	}
}

// Main starts the main X event loop. It will read events and call appropriate
// callback functions.
// N.B. If you have multiple X connections in the same program, you should be
// able to run this in different goroutines concurrently. However, only
// *one* of these should run for *each* connection.
func Main(xu *xgbutil.XUtil) {
	mainEventLoop(xu, nil, nil, nil)
}

// MainPing starts the main X event loop, and returns three "ping" channels:
// the first is pinged before an event is dequeued, the second is pinged
// after all callbacks for a particular event have been called and the last
// is pinged when the event loop stops (e.g., after a call to xevent.Quit).
// pingAfter channel.
//
// This is useful if your event loop needs to draw from other sources. e.g.,
//
//	pingBefore, pingAfter, pingQuit := xevent.MainPing()
//	for {
//		select {
//		case <-pingBefore:
//			// Wait for event processing to finish.
//			<-pingAfter
//		case val <- someOtherChannel:
//			// do some work with val
//		case <-pingQuit:
//			fmt.Printf("xevent loop has quit")
//			return
//		}
//	}
//
// Note that an unbuffered channel is returned, which implies that any work
// done with 'val' will delay further X event processing.
//
// A complete example using MainPing can be found in the examples directory in
// the xgbutil package under the name multiple-source-event-loop.
func MainPing(xu *xgbutil.XUtil) (chan struct{}, chan struct{}, chan struct{}) {
	pingBefore := make(chan struct{}, 0)
	pingAfter := make(chan struct{}, 0)
	pingQuit := make(chan struct{}, 0)
	go func() {
		mainEventLoop(xu, pingBefore, pingAfter, pingQuit)
	}()
	return pingBefore, pingAfter, pingQuit
}

// mainEventLoop runs the main event loop with an optional ping channel.
func mainEventLoop(xu *xgbutil.XUtil,
	pingBefore, pingAfter, pingQuit chan struct{}) {
	for {
		if Quitting(xu) {
			if pingQuit != nil {
				pingQuit <- struct{}{}
			}
			break
		}

		// Gobble up as many events as possible (into the queue).
		// If there are no events, we block.
		Read(xu, true)

		// Now process every event/error in the queue.
		processEventQueue(xu, pingBefore, pingAfter)
	}
}

// processEventQueue processes every item in the event/error queue.
func processEventQueue(xu *xgbutil.XUtil, pingBefore, pingAfter chan struct{}) {
	for !Empty(xu) {
		if Quitting(xu) {
			return
		}

		// We send the ping *before* the next event is dequeued.
		// This is so the queue doesn't present a misrepresentation of which
		// events haven't been processed yet.
		if pingBefore != nil && pingAfter != nil {
			pingBefore <- struct{}{}
		}
		ev, err := Dequeue(xu)

		// If we gobbled up an error, send it to the error event handler
		// and move on the next event/error.
		if err != nil {
			ErrorHandlerGet(xu)(err)
			if pingBefore != nil && pingAfter != nil {
				pingAfter <- struct{}{}
			}
			continue
		}

		// We know there isn't an error. If there isn't an event either,
		// then there's a bug somewhere.
		if ev == nil {
			xgbutil.Logger.Fatal("BUG: Expected an event but got nil.")
		}

		hooks := getHooks(xu)
		for _, hook := range hooks {
			if !hook.Run(xu, ev) {
				goto END
			}
		}

		switch event := ev.(type) {
		case xproto.KeyPressEvent:
			e := KeyPressEvent{&event}

			// If we're redirecting key events, this is the place to do it!
			if wid := RedirectKeyGet(xu); wid > 0 {
				e.Event = wid
			}

			xu.TimeSet(e.Time)
			runCallbacks(xu, e, KeyPress, e.Event)
		case xproto.KeyReleaseEvent:
			e := KeyReleaseEvent{&event}

			// If we're redirecting key events, this is the place to do it!
			if wid := RedirectKeyGet(xu); wid > 0 {
				e.Event = wid
			}

			xu.TimeSet(e.Time)
			runCallbacks(xu, e, KeyRelease, e.Event)
		case xproto.ButtonPressEvent:
			e := ButtonPressEvent{&event}
			xu.TimeSet(e.Time)
			runCallbacks(xu, e, ButtonPress, e.Event)
		case xproto.ButtonReleaseEvent:
			e := ButtonReleaseEvent{&event}
			xu.TimeSet(e.Time)
			runCallbacks(xu, e, ButtonRelease, e.Event)
		case xproto.MotionNotifyEvent:
			e := MotionNotifyEvent{&event}
			xu.TimeSet(e.Time)
			runCallbacks(xu, e, MotionNotify, e.Event)
		case xproto.EnterNotifyEvent:
			e := EnterNotifyEvent{&event}
			xu.TimeSet(e.Time)
			runCallbacks(xu, e, EnterNotify, e.Event)
		case xproto.LeaveNotifyEvent:
			e := LeaveNotifyEvent{&event}
			xu.TimeSet(e.Time)
			runCallbacks(xu, e, LeaveNotify, e.Event)
		case xproto.FocusInEvent:
			e := FocusInEvent{&event}
			runCallbacks(xu, e, FocusIn, e.Event)
		case xproto.FocusOutEvent:
			e := FocusOutEvent{&event}
			runCallbacks(xu, e, FocusOut, e.Event)
		case xproto.KeymapNotifyEvent:
			e := KeymapNotifyEvent{&event}
			runCallbacks(xu, e, KeymapNotify, NoWindow)
		case xproto.ExposeEvent:
			e := ExposeEvent{&event}
			runCallbacks(xu, e, Expose, e.Window)
		case xproto.GraphicsExposureEvent:
			e := GraphicsExposureEvent{&event}
			runCallbacks(xu, e, GraphicsExposure, xproto.Window(e.Drawable))
		case xproto.NoExposureEvent:
			e := NoExposureEvent{&event}
			runCallbacks(xu, e, NoExposure, xproto.Window(e.Drawable))
		case xproto.VisibilityNotifyEvent:
			e := VisibilityNotifyEvent{&event}
			runCallbacks(xu, e, VisibilityNotify, e.Window)
		case xproto.CreateNotifyEvent:
			e := CreateNotifyEvent{&event}
			runCallbacks(xu, e, CreateNotify, e.Parent)
		case xproto.DestroyNotifyEvent:
			e := DestroyNotifyEvent{&event}
			runCallbacks(xu, e, DestroyNotify, e.Window)
		case xproto.UnmapNotifyEvent:
			e := UnmapNotifyEvent{&event}
			runCallbacks(xu, e, UnmapNotify, e.Window)
		case xproto.MapNotifyEvent:
			e := MapNotifyEvent{&event}
			runCallbacks(xu, e, MapNotify, e.Event)
		case xproto.MapRequestEvent:
			e := MapRequestEvent{&event}
			runCallbacks(xu, e, MapRequest, e.Window)
			runCallbacks(xu, e, MapRequest, e.Parent)
		case xproto.ReparentNotifyEvent:
			e := ReparentNotifyEvent{&event}
			runCallbacks(xu, e, ReparentNotify, e.Window)
		case xproto.ConfigureNotifyEvent:
			e := ConfigureNotifyEvent{&event}
			runCallbacks(xu, e, ConfigureNotify, e.Window)
		case xproto.ConfigureRequestEvent:
			e := ConfigureRequestEvent{&event}
			runCallbacks(xu, e, ConfigureRequest, e.Window)
			runCallbacks(xu, e, ConfigureRequest, e.Parent)
		case xproto.GravityNotifyEvent:
			e := GravityNotifyEvent{&event}
			runCallbacks(xu, e, GravityNotify, e.Window)
		case xproto.ResizeRequestEvent:
			e := ResizeRequestEvent{&event}
			runCallbacks(xu, e, ResizeRequest, e.Window)
		case xproto.CirculateNotifyEvent:
			e := CirculateNotifyEvent{&event}
			runCallbacks(xu, e, CirculateNotify, e.Window)
		case xproto.CirculateRequestEvent:
			e := CirculateRequestEvent{&event}
			runCallbacks(xu, e, CirculateRequest, e.Window)
		case xproto.PropertyNotifyEvent:
			e := PropertyNotifyEvent{&event}
			xu.TimeSet(e.Time)
			runCallbacks(xu, e, PropertyNotify, e.Window)
		case xproto.SelectionClearEvent:
			e := SelectionClearEvent{&event}
			xu.TimeSet(e.Time)
			runCallbacks(xu, e, SelectionClear, e.Owner)
		case xproto.SelectionRequestEvent:
			e := SelectionRequestEvent{&event}
			xu.TimeSet(e.Time)
			runCallbacks(xu, e, SelectionRequest, e.Requestor)
		case xproto.SelectionNotifyEvent:
			e := SelectionNotifyEvent{&event}
			xu.TimeSet(e.Time)
			runCallbacks(xu, e, SelectionNotify, e.Requestor)
		case xproto.ColormapNotifyEvent:
			e := ColormapNotifyEvent{&event}
			runCallbacks(xu, e, ColormapNotify, e.Window)
		case xproto.ClientMessageEvent:
			e := ClientMessageEvent{&event}
			runCallbacks(xu, e, ClientMessage, e.Window)
		case xproto.MappingNotifyEvent:
			e := MappingNotifyEvent{&event}
			runCallbacks(xu, e, MappingNotify, NoWindow)
		case shape.NotifyEvent:
			e := ShapeNotifyEvent{&event}
			runCallbacks(xu, e, ShapeNotify, e.AffectedWindow)
		default:
			if event != nil {
				xgbutil.Logger.Printf("ERROR: UNSUPPORTED EVENT TYPE: %T",
					event)
			}
		}

	END:

		if pingBefore != nil && pingAfter != nil {
			pingAfter <- struct{}{}
		}
	}
}
