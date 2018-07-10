package xgbutil

/*
types.go contains several types used in the XUtil structure. In an ideal world,
they would be defined in their appropriate packages, but must be defined here
(and exported) for use in some sub-packages. (Namely, xevent, keybind and
mousebind.)
*/

import (
	"github.com/BurntSushi/xgb"
	"github.com/BurntSushi/xgb/xproto"
)

// Callback is an interface that should be implemented by event callback
// functions. Namely, to assign a function to a particular event/window
// combination, simply define a function with type 'SomeEventFun' (pre-defined
// in xevent/callback.go), and call the 'Connect' method.
// The 'Run' method is used inside the Main event loop, and shouldn't be used
// by the user.
// Also, it is perfectly legitimate to connect to events that don't specify
// a window (like MappingNotify and KeymapNotify). In this case, simply
// use 'xgbutil.NoWindow' as the window id.
//
// Example to respond to ConfigureNotify events on window 0x1
//
//	xevent.ConfigureNotifyFun(
//		func(X *xgbutil.XUtil, e xevent.ConfigureNotifyEvent) {
//			fmt.Printf("(%d, %d) %dx%d\n", e.X, e.Y, e.Width, e.Height)
//		}).Connect(X, 0x1)
type Callback interface {
	// Connect modifies XUtil's state to attach an event handler to a
	// particular event.
	Connect(xu *XUtil, win xproto.Window)

	// Run is exported for use in the xevent package but should not be
	// used by the user. (It is used to run the callback function in the
	// main event loop.)
	Run(xu *XUtil, ev interface{})
}

// CallbackHook works similarly to the more general Callback, but it is
// for hooks into the main xevent loop. As such it does not get attached
// to a window.
type CallbackHook interface {
	// Connect connects this hook to the main loop of the passed XUtil
	// instance.
	Connect(xu *XUtil)

	// Run is exported for use in the xevent package, but should not be
	// used by the user.  It should return true if it's ok to process
	// the event as usual, or false if it should be suppressed.
	Run(xu *XUtil, ev interface{}) bool
}

// CallbackKey works similarly to the more general Callback, but it adds
// parameters specific to key bindings.
type CallbackKey interface {
	// Connect modifies XUtil's state to attach an event handler to a
	// particular key press. If grab is true, connect will request a passive
	// grab.
	Connect(xu *XUtil, win xproto.Window, keyStr string, grab bool) error

	// Run is exported for use in the keybind package but should not be
	// used by the user. (It is used to run the callback function in the
	// main event loop.
	Run(xu *XUtil, ev interface{})
}

// CallbackMouse works similarly to the more general Callback, but it adds
// parameters specific to mouse bindings.
type CallbackMouse interface {
	// Connect modifies XUtil's state to attach an event handler to a
	// particular button press.
	// If sync is true, the grab will be synchronous. (This will require a
	// call to xproto.AllowEvents in response, otherwise no further events
	// will be processed and your program will lock.)
	// If grab is true, connect will request a passive grab.
	Connect(xu *XUtil, win xproto.Window, buttonStr string,
		sync bool, grab bool) error

	// Run is exported for use in the mousebind package but should not be
	// used by the user. (It is used to run the callback function in the
	// main event loop.)
	Run(xu *XUtil, ev interface{})
}

// KeyKey is the type of the key in the map of keybindings.
// It essentially represents the tuple
// (event type, window id, modifier, keycode).
// It is exported for use in the keybind package. It should not be used.
type KeyKey struct {
	Evtype int
	Win    xproto.Window
	Mod    uint16
	Code   xproto.Keycode
}

// KeyString is the type of a key binding string used to connect to particular
// key combinations. A list of all such key strings is maintained in order to
// rebind keys when the keyboard mapping has been changed.
type KeyString struct {
	Str      string
	Callback CallbackKey
	Evtype   int
	Win      xproto.Window
	Grab     bool
}

// MouseKey is the type of the key in the map of mouse bindings.
// It essentially represents the tuple
// (event type, window id, modifier, button).
// It is exported for use in the mousebind package. It should not be used.
type MouseKey struct {
	Evtype int
	Win    xproto.Window
	Mod    uint16
	Button xproto.Button
}

// KeyboardMapping embeds a keyboard mapping reply from XGB.
// It should be retrieved using keybind.KeyMapGet, if necessary.
// xgbutil tries quite hard to absolve you from ever having to use this.
// A keyboard mapping is a table that maps keycodes to one or more keysyms.
type KeyboardMapping struct {
	*xproto.GetKeyboardMappingReply
}

// ModifierMapping embeds a modifier mapping reply from XGB.
// It should be retrieved using keybind.ModMapGet, if necessary.
// xgbutil tries quite hard to absolve you from ever having to use this.
// A modifier mapping is a table that maps modifiers to one or more keycodes.
type ModifierMapping struct {
	*xproto.GetModifierMappingReply
}

// ErrorHandlerFun is the type of function required to handle errors that
// come in through the main event loop.
// For example, to set a new error handler, use:
//
//	xevent.ErrorHandlerSet(xgbutil.ErrorHandlerFun(
//		func(err xgb.Error) {
//			// do something with err
//		}))
type ErrorHandlerFun func(err xgb.Error)

// EventOrError is a struct that contains either an event value or an error
// value. It is an error to contain both. Containing neither indicates an
// error too.
// This is exported for use in the xevent package. You shouldn't have any
// direct contact with values of this type, unless you need to inspect the
// queue directly with xevent.Peek.
type EventOrError struct {
	Event xgb.Event
	Err   xgb.Error
}

// MouseDragFun is the kind of function used on each dragging step
// and at the end of a drag.
type MouseDragFun func(xu *XUtil, rootX, rootY, eventX, eventY int)

// MouseDragBeginFun is the kind of function used to initialize a drag.
// The difference between this and MouseDragFun is that the begin function
// returns a bool (of whether or not to cancel the drag) and an X resource
// identifier corresponding to a cursor.
type MouseDragBeginFun func(xu *XUtil, rootX, rootY,
	eventX, eventY int) (bool, xproto.Cursor)
