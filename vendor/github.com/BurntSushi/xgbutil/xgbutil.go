package xgbutil

import (
	"log"
	"os"
	"sync"

	"github.com/BurntSushi/xgb"
	"github.com/BurntSushi/xgb/xinerama"
	"github.com/BurntSushi/xgb/xproto"
)

// Logger is used through xgbutil when messages need to be emitted to stderr.
var Logger = log.New(os.Stderr, "[xgbutil] ", log.Lshortfile)

// The current maximum request size. I think we can expand this with
// BigReq, but it probably isn't worth it at the moment.
const MaxReqSize = (1 << 16) * 4

// An XUtil represents the state of xgbutil. It keeps track of the current
// X connection, the root window, event callbacks, key/mouse bindings, etc.
// Regrettably, many of the members are exported, even though they should not
// be used directly by the user. They are exported for use in sub-packages.
// (Namely, xevent, keybind and mousebind.) In fact, there should *never*
// be a reason to access any members of an XUtil value directly. Any
// interaction with an XUtil value should be through its methods.
type XUtil struct {
	// conn is the XGB connection object used to issue protocol requests.
	conn *xgb.Conn

	// Quit can be set to true, and the main event loop will finish processing
	// the current event, and gracefully quit afterwards.
	// This is exported for use in the xevent package. Please us xevent.Quit
	// to set this value.
	Quit bool // when true, the main event loop will stop gracefully

	// setup contains all the setup information retrieved at connection time.
	setup *xproto.SetupInfo

	// screen is a simple alias to the default screen info.
	screen *xproto.ScreenInfo

	// root is an alias to the default root window.
	root xproto.Window

	// Atoms is a cache of atom names to resource identifiers. This minimizes
	// round trips to the X server, since atom identifiers never change.
	// It is exported for use in the xprop package. It should not be used.
	Atoms    map[string]xproto.Atom
	AtomsLck *sync.RWMutex

	// AtomNames is a cache just like 'atoms', but in the reverse direction.
	// It is exported for use in the xprop package. It should not be used.
	AtomNames    map[xproto.Atom]string
	AtomNamesLck *sync.RWMutex

	// Evqueue is the queue that stores the results of xgb.WaitForEvent.
	// Namely, each value is either an Event *or* an Error.
	// It is exported for use in the xevent package. Do not use it.
	// If you need to interact with the event queue, please use the functions
	// available in the xevent package: Dequeue, DequeueAt, QueueEmpty
	// and QueuePeek.
	Evqueue    []EventOrError
	EvqueueLck *sync.RWMutex

	// Callbacks is a map of event numbers to a map of window identifiers
	// to callback functions.
	// This is the data structure that stores all callback functions, where
	// a callback function is always attached to a (event, window) tuple.
	// It is exported for use in the xevent package. Do not use it.
	Callbacks    map[int]map[xproto.Window][]Callback
	CallbacksLck *sync.RWMutex

	// Hooks are called by the XEvent main loop before processing the event
	// itself. These are meant for instances when it's not possible / easy
	// to use the normal Hook system. You should not modify this yourself.
	Hooks    []CallbackHook
	HooksLck *sync.RWMutex

	// eventTime is the last time recorded by an event. It is automatically
	// updated if xgbutil's main event loop is used.
	eventTime xproto.Timestamp

	// Keymap corresponds to xgbutil's current conception of the keyboard
	// mapping. It is automatically kept up-to-date if xgbutil's event loop
	// is used.
	// It is exported for use in the keybind package. It should not be
	// accessed directly. Instead, use keybind.KeyMapGet.
	Keymap *KeyboardMapping

	// Modmap corresponds to xgbutil's current conception of the modifier key
	// mapping. It is automatically kept up-to-date if xgbutil's event loop
	// is used.
	// It is exported for use in the keybind package. It should not be
	// accessed directly. Instead, use keybind.ModMapGet.
	Modmap *ModifierMapping

	// KeyRedirect corresponds to a window identifier that, when set,
	// automatically receives *all* keyboard events. This is a sort-of
	// synthetic grab and is helpful in avoiding race conditions.
	// It is exported for use in the xevent and keybind packages. Do not use
	// it directly. To redirect key events, please use xevent.RedirectKeyEvents.
	KeyRedirect xproto.Window

	// Keybinds is the data structure storing all callbacks for key bindings.
	// This is extremely similar to the general notion of event callbacks,
	// but adds extra support to make handling key bindings easier. (Like
	// specifying human readable key sequences to bind to.)
	// KeyBindKey is a struct representing the 4-tuple
	// (event-type, window-id, modifiers, keycode).
	// It is exported for use in the keybind package. Do not access it directly.
	Keybinds    map[KeyKey][]CallbackKey
	KeybindsLck *sync.RWMutex

	// Keygrabs is a frequency count of the number of callbacks associated
	// with a particular KeyBindKey. This is necessary because we can only
	// grab a particular key *once*, but we may want to attach several callbacks
	// to a single keypress.
	// It is exported for use in the keybind package. Do not access it directly.
	Keygrabs map[KeyKey]int

	// Keystrings is a list of all key strings used to connect keybindings.
	// They are used to rebuild key grabs when the keyboard mapping is updated.
	// It is exported for use in the keybind package. Do not access it directly.
	Keystrings []KeyString

	// Mousebinds is the data structure storing all callbacks for mouse
	// bindings.This is extremely similar to the general notion of event
	// callbacks,but adds extra support to make handling mouse bindings easier.
	// (Like specifying human readable mouse sequences to bind to.)
	// MouseBindKey is a struct representing the 4-tuple
	// (event-type, window-id, modifiers, button).
	// It is exported for use in the mousebind package. Do not use it.
	Mousebinds    map[MouseKey][]CallbackMouse
	MousebindsLck *sync.RWMutex

	// Mousegrabs is a frequency count of the number of callbacks associated
	// with a particular MouseBindKey. This is necessary because we can only
	// grab a particular mouse button *once*, but we may want to attach
	// several callbacks to a single button press.
	// It is exported for use in the mousebind package. Do not use it.
	Mousegrabs map[MouseKey]int

	// InMouseDrag is true if a drag is currently in progress.
	// It is exported for use in the mousebind package. Do not use it.
	InMouseDrag bool

	// MouseDragStep is the function executed for each step (i.e., pointer
	// movement) in the current mouse drag. Note that this is nil when a drag
	// is not in progress.
	// It is exported for use in the mousebind package. Do not use it.
	MouseDragStepFun MouseDragFun

	// MouseDragEnd is the function executed at the end of the current
	// mouse drag. This is nil when a drag is not in progress.
	// It is exported for use in the mousebind package. Do not use it.
	MouseDragEndFun MouseDragFun

	// gc is a general purpose graphics context; used to paint images.
	// Since we don't do any real X drawing, we don't really care about the
	// particulars of our graphics context.
	gc xproto.Gcontext

	// dummy is a dummy window used for mouse/key GRABs.
	// Basically, whenever a grab is instituted, mouse and key events are
	// redirected to the dummy the window.
	dummy xproto.Window

	// ErrorHandler is the function that handles errors *in the event loop*.
	// By default, it simply emits them to stderr.
	// It is exported for use in the xevent package. To set the default error
	// handler, please use xevent.ErrorHandlerSet.
	ErrorHandler ErrorHandlerFun
}

// NewConn connects to the X server using the DISPLAY environment variable
// and creates a new XUtil. Most environments have the DISPLAY environment
// variable set, so this is probably what you want to use to connect to X.
func NewConn() (*XUtil, error) {
	return NewConnDisplay("")
}

// NewConnDisplay connects to the X server and creates a new XUtil.
// If 'display' is empty, the DISPLAY environment variable is used. Otherwise
// there are several different display formats supported:
//
//	NewConn(":1") -> net.Dial("unix", "", "/tmp/.X11-unix/X1")
//	NewConn("/tmp/launch-12/:0") -> net.Dial("unix", "", "/tmp/launch-12/:0")
//	NewConn("hostname:2.1") -> net.Dial("tcp", "", "hostname:6002")
//	NewConn("tcp/hostname:1.0") -> net.Dial("tcp", "", "hostname:6001")
func NewConnDisplay(display string) (*XUtil, error) {
	c, err := xgb.NewConnDisplay(display)

	if err != nil {
		return nil, err
	}

	return NewConnXgb(c)

}

// NewConnXgb use the specific xgb.Conn to create a new XUtil.
//
//  NewConn, NewConnDisplay are wrapper of this function.
func NewConnXgb(c *xgb.Conn) (*XUtil, error) {
	setup := xproto.Setup(c)
	screen := setup.DefaultScreen(c)

	// Initialize our central struct that stores everything.
	xu := &XUtil{
		conn:             c,
		Quit:             false,
		Evqueue:          make([]EventOrError, 0, 1000),
		EvqueueLck:       &sync.RWMutex{},
		setup:            setup,
		screen:           screen,
		root:             screen.Root,
		eventTime:        xproto.Timestamp(0), // last event time
		Atoms:            make(map[string]xproto.Atom, 50),
		AtomsLck:         &sync.RWMutex{},
		AtomNames:        make(map[xproto.Atom]string, 50),
		AtomNamesLck:     &sync.RWMutex{},
		Callbacks:        make(map[int]map[xproto.Window][]Callback, 33),
		CallbacksLck:     &sync.RWMutex{},
		Hooks:            make([]CallbackHook, 0),
		HooksLck:         &sync.RWMutex{},
		Keymap:           nil, // we don't have anything yet
		Modmap:           nil,
		KeyRedirect:      0,
		Keybinds:         make(map[KeyKey][]CallbackKey, 10),
		KeybindsLck:      &sync.RWMutex{},
		Keygrabs:         make(map[KeyKey]int, 10),
		Keystrings:       make([]KeyString, 0, 10),
		Mousebinds:       make(map[MouseKey][]CallbackMouse, 10),
		MousebindsLck:    &sync.RWMutex{},
		Mousegrabs:       make(map[MouseKey]int, 10),
		InMouseDrag:      false,
		MouseDragStepFun: nil,
		MouseDragEndFun:  nil,
		ErrorHandler:     func(err xgb.Error) { Logger.Println(err) },
	}

	var err error = nil
	// Create a general purpose graphics context
	xu.gc, err = xproto.NewGcontextId(xu.conn)
	if err != nil {
		return nil, err
	}
	xproto.CreateGC(xu.conn, xu.gc, xproto.Drawable(xu.root),
		xproto.GcForeground, []uint32{xu.screen.WhitePixel})

	// Create a dummy window
	xu.dummy, err = xproto.NewWindowId(xu.conn)
	if err != nil {
		return nil, err
	}
	xproto.CreateWindow(xu.conn, xu.Screen().RootDepth, xu.dummy, xu.RootWin(),
		-1000, -1000, 1, 1, 0,
		xproto.WindowClassInputOutput, xu.Screen().RootVisual,
		xproto.CwEventMask|xproto.CwOverrideRedirect,
		[]uint32{1, xproto.EventMaskPropertyChange})
	xproto.MapWindow(xu.conn, xu.dummy)

	// Register the Xinerama extension... because it doesn't cost much.
	err = xinerama.Init(xu.conn)

	// If we can't register Xinerama, that's okay. Output something
	// and move on.
	if err != nil {
		Logger.Printf("WARNING: %s\n", err)
		Logger.Printf("MESSAGE: The 'xinerama' package cannot be used " +
			"because the XINERAMA extension could not be loaded.")
	}

	return xu, nil
}

// Conn returns the xgb connection object.
func (xu *XUtil) Conn() *xgb.Conn {
	return xu.conn
}

// ExtInitialized returns true if an extension has been initialized.
// This is useful for determining whether an extension is available or not.
func (xu *XUtil) ExtInitialized(extName string) bool {
	_, ok := xu.Conn().Extensions[extName]
	return ok
}

// Sync forces XGB to catch up with all events/requests and synchronize.
// This is done by issuing a benign round trip request to X.
func (xu *XUtil) Sync() {
	xproto.GetInputFocus(xu.Conn()).Reply()
}

// Setup returns the setup information retrieved during connection time.
func (xu *XUtil) Setup() *xproto.SetupInfo {
	return xu.setup
}

// Screen returns the default screen
func (xu *XUtil) Screen() *xproto.ScreenInfo {
	return xu.screen
}

// RootWin returns the current root window.
func (xu *XUtil) RootWin() xproto.Window {
	return xu.root
}

// RootWinSet will change the current root window to the one provided.
// N.B. This probably shouldn't be used unless you're desperately trying
// to support multiple X screens. (This is *not* the same as Xinerama/RandR or
// TwinView. All of those have a single root window.)
func (xu *XUtil) RootWinSet(root xproto.Window) {
	xu.root = root
}

// TimeGet gets the most recent time seen by an event.
func (xu *XUtil) TimeGet() xproto.Timestamp {
	return xu.eventTime
}

// TimeSet sets the most recent time seen by an event.
func (xu *XUtil) TimeSet(t xproto.Timestamp) {
	xu.eventTime = t
}

// GC gets a general purpose graphics context that is typically used to simply
// paint images.
func (xu *XUtil) GC() xproto.Gcontext {
	return xu.gc
}

// Dummy gets the id of the dummy window.
func (xu *XUtil) Dummy() xproto.Window {
	return xu.dummy
}

// Grabs the server. Everything becomes synchronous.
func (xu *XUtil) Grab() {
	xproto.GrabServer(xu.Conn())
}

// Ungrabs the server.
func (xu *XUtil) Ungrab() {
	xproto.UngrabServer(xu.Conn())
}
