//go:build linux
// +build linux

// Copyright (c) 2016-2026 AtomAI, All rights reserved.
//
// See the COPYRIGHT file at the top-level directory of this distribution and at
// https://github.com/go-vgo/robotgo/blob/master/LICENSE
//
// Licensed under the Apache License, Version 2.0 <LICENSE-APACHE or
// http://www.apache.org/licenses/LICENSE-2.0>
//
// This file may not be copied, modified, or distributed
// except according to those terms.

package wayland

import (
	"errors"
	"fmt"
	"log"
	"os"
	"sync"
	"sync/atomic"
	"time"

	"github.com/vcaesar/go-wayland/client"

	"github.com/go-vgo/robotgo/wayland/internal/protocols/wlr_foreign_toplevel"
	"github.com/go-vgo/robotgo/wayland/internal/protocols/wlr_screencopy"
	"github.com/go-vgo/robotgo/wayland/internal/protocols/wlr_virtual_keyboard"
	"github.com/go-vgo/robotgo/wayland/internal/protocols/wlr_virtual_pointer"
)

// conn holds the singleton Wayland connection and all bound protocol objects.
type conn struct {
	mu sync.Mutex

	display  *client.Display
	registry *client.Registry

	seat    *client.Seat
	shm     *client.Shm
	outputs []*outputInfo

	pointerManager  *wlr_virtual_pointer.ZwlrVirtualPointerManagerV1
	pointer         *wlr_virtual_pointer.ZwlrVirtualPointerV1
	keyboardManager *wlr_virtual_keyboard.ZwpVirtualKeyboardManagerV1
	keyboard        *wlr_virtual_keyboard.ZwpVirtualKeyboardV1
	screencopyMgr   *wlr_screencopy.ZwlrScreencopyManagerV1
	toplevelMgr     *wlr_foreign_toplevel.ZwlrForeignToplevelManagerV1

	toplevels map[uint32]*toplevelInfo

	keymapSet bool

	// dispatch loop
	dispatchDone chan struct{}
	closed       atomic.Bool
}

// outputInfo tracks wl_output geometry.
type outputInfo struct {
	output *client.Output
	x, y   int32
	width  int32
	height int32
	name   string
}

// toplevelInfo tracks a foreign toplevel handle.
type toplevelInfo struct {
	handle *wlr_foreign_toplevel.ZwlrForeignToplevelHandleV1
	title  string
	appId  string
	states []byte
}

var (
	globalConn *conn
	connMu     sync.Mutex
)

// ErrNotSupported is returned when a required Wayland protocol is not available.
var ErrNotSupported = errors.New("robotgo: required wayland protocol not supported by compositor")

// ErrNoConnection is returned when the Wayland connection is not established.
var ErrNoConnection = errors.New("robotgo: wayland connection not established")

// ensureConn lazily initializes the global Wayland connection, re-establishing
// it if it was never created or has since been closed. A mutex (instead of
// sync.Once) keeps the backend recoverable after a failed connect or Close.
func ensureConn() (*conn, error) {
	connMu.Lock()
	defer connMu.Unlock()

	if globalConn != nil && !globalConn.closed.Load() {
		return globalConn, nil
	}

	c, err := newConn()
	if err != nil {
		globalConn = nil
		return nil, err
	}
	globalConn = c
	return globalConn, nil
}

// newConn creates a new Wayland connection and binds all required protocols.
func newConn() (*conn, error) {
	c := &conn{
		toplevels:    make(map[uint32]*toplevelInfo),
		dispatchDone: make(chan struct{}),
	}

	display, err := client.Connect("")
	if err != nil {
		return nil, fmt.Errorf("robotgo: connect to wayland: %w", err)
	}
	c.display = display

	registry, err := display.GetRegistry()
	if err != nil {
		_ = display.Context().Close()
		return nil, fmt.Errorf("robotgo: get registry: %w", err)
	}
	c.registry = registry

	registry.SetGlobalHandler(func(e client.RegistryGlobalEvent) {
		c.handleGlobal(e)
	})

	// Two roundtrips: first to get globals, second to get events from bound objects
	c.roundtrip()
	c.roundtrip()

	// Create virtual pointer if manager is available. Only retain the proxy
	// on success; on failure leave c.pointer nil so we never keep a
	// half-initialized object around.
	if c.pointerManager != nil && c.seat != nil {
		pointer, perr := c.pointerManager.CreateVirtualPointer(c.seat)
		if perr != nil {
			log.Printf("robotgo: create virtual pointer: %v", perr)
		} else {
			c.pointer = pointer
		}
	}

	// Create virtual keyboard if manager is available (same retain-on-success
	// rule as the pointer above).
	if c.keyboardManager != nil && c.seat != nil {
		keyboard, kerr := c.keyboardManager.CreateVirtualKeyboard(c.seat)
		if kerr != nil {
			log.Printf("robotgo: create virtual keyboard: %v", kerr)
		} else {
			c.keyboard = keyboard
		}
	}

	// Set up XKB keymap for virtual keyboard
	if c.keyboard != nil {
		if err := c.setupKeymap(); err != nil {
			log.Printf("robotgo: setup keymap: %v", err)
		}
	}

	// Start dispatch loop in background
	go c.dispatchLoop()

	return c, nil
}

// handleGlobal binds protocol globals as they are advertised.
func (c *conn) handleGlobal(e client.RegistryGlobalEvent) {
	switch e.Interface {
	case "wl_seat":
		if c.seat == nil {
			c.seat = client.NewSeat(c.display.Context())
			if err := c.registry.Bind(e.Name, e.Interface, e.Version, c.seat); err != nil {
				log.Printf("robotgo: bind wl_seat: %v", err)
			}
		}
	case "wl_shm":
		if c.shm == nil {
			c.shm = client.NewShm(c.display.Context())
			if err := c.registry.Bind(e.Name, e.Interface, e.Version, c.shm); err != nil {
				log.Printf("robotgo: bind wl_shm: %v", err)
			}
		}
	case "wl_output":
		out := client.NewOutput(c.display.Context())
		if err := c.registry.Bind(e.Name, e.Interface, e.Version, out); err != nil {
			log.Printf("robotgo: bind wl_output: %v", err)
			return
		}
		info := &outputInfo{output: out}
		out.SetGeometryHandler(func(ge client.OutputGeometryEvent) {
			info.x = int32(ge.X)
			info.y = int32(ge.Y)
		})
		out.SetModeHandler(func(me client.OutputModeEvent) {
			if me.Flags&0x1 != 0 { // WL_OUTPUT_MODE_CURRENT
				info.width = int32(me.Width)
				info.height = int32(me.Height)
			}
		})
		c.outputs = append(c.outputs, info)

	case wlr_virtual_pointer.ZwlrVirtualPointerManagerV1InterfaceName:
		c.pointerManager = wlr_virtual_pointer.NewZwlrVirtualPointerManagerV1(c.display.Context())
		if err := c.registry.Bind(e.Name, e.Interface, e.Version, c.pointerManager); err != nil {
			log.Printf("robotgo: bind virtual pointer manager: %v", err)
		}

	case wlr_virtual_keyboard.ZwpVirtualKeyboardManagerV1InterfaceName:
		c.keyboardManager = wlr_virtual_keyboard.NewZwpVirtualKeyboardManagerV1(c.display.Context())
		if err := c.registry.Bind(e.Name, e.Interface, e.Version, c.keyboardManager); err != nil {
			log.Printf("robotgo: bind virtual keyboard manager: %v", err)
		}

	case wlr_screencopy.ZwlrScreencopyManagerV1InterfaceName:
		c.screencopyMgr = wlr_screencopy.NewZwlrScreencopyManagerV1(c.display.Context())
		if err := c.registry.Bind(e.Name, e.Interface, e.Version, c.screencopyMgr); err != nil {
			log.Printf("robotgo: bind screencopy manager: %v", err)
		}

	case wlr_foreign_toplevel.ZwlrForeignToplevelManagerV1InterfaceName:
		c.toplevelMgr = wlr_foreign_toplevel.NewZwlrForeignToplevelManagerV1(c.display.Context())
		if err := c.registry.Bind(e.Name, e.Interface, e.Version, c.toplevelMgr); err != nil {
			log.Printf("robotgo: bind foreign toplevel manager: %v", err)
			return
		}
		c.toplevelMgr.SetToplevelHandler(func(te wlr_foreign_toplevel.ZwlrForeignToplevelManagerV1ToplevelEvent) {
			c.handleNewToplevel(te.Toplevel)
		})
	}
}

// handleNewToplevel sets up event handlers for a newly discovered toplevel.
func (c *conn) handleNewToplevel(handle *wlr_foreign_toplevel.ZwlrForeignToplevelHandleV1) {
	c.mu.Lock()
	info := &toplevelInfo{handle: handle}
	c.toplevels[handle.ID()] = info
	c.mu.Unlock()

	handle.SetTitleHandler(func(e wlr_foreign_toplevel.ZwlrForeignToplevelHandleV1TitleEvent) {
		c.mu.Lock()
		info.title = e.Title
		c.mu.Unlock()
	})
	handle.SetAppIdHandler(func(e wlr_foreign_toplevel.ZwlrForeignToplevelHandleV1AppIdEvent) {
		c.mu.Lock()
		info.appId = e.AppId
		c.mu.Unlock()
	})
	handle.SetStateHandler(func(e wlr_foreign_toplevel.ZwlrForeignToplevelHandleV1StateEvent) {
		c.mu.Lock()
		info.states = e.State
		c.mu.Unlock()
	})
	handle.SetClosedHandler(func(_ wlr_foreign_toplevel.ZwlrForeignToplevelHandleV1ClosedEvent) {
		c.mu.Lock()
		delete(c.toplevels, handle.ID())
		c.mu.Unlock()
	})
}

// roundtrip performs a blocking wl_display.sync roundtrip.
func (c *conn) roundtrip() {
	display := c.display
	display.Roundtrip()
}

// dispatchLoop runs the Wayland event dispatch loop in a background goroutine.
func (c *conn) dispatchLoop() {
	defer close(c.dispatchDone)
	for {
		if err := c.display.Context().Dispatch(); err != nil {
			if c.closed.Load() {
				return
			}
			log.Printf("robotgo: dispatch error: %v", err)
			return
		}
	}
}

// timestamp returns a Wayland-compatible millisecond timestamp.
func timestamp() uint32 {
	return uint32(time.Now().UnixMilli())
}

// setupKeymap writes a minimal XKB keymap to a temp file and sends it
// to the virtual keyboard. This must be done before any key events.
func (c *conn) setupKeymap() error {
	// Minimal but complete XKB keymap that maps evdev keycodes directly.
	// This keymap covers the standard 104-key US layout.
	keymap := `xkb_keymap {
	xkb_keycodes "evdev" {
		minimum = 8;
		maximum = 255;
		<ESC> = 9;
		<AE01> = 10; <AE02> = 11; <AE03> = 12; <AE04> = 13;
		<AE05> = 14; <AE06> = 15; <AE07> = 16; <AE08> = 17;
		<AE09> = 18; <AE10> = 19; <AE11> = 20; <AE12> = 21;
		<BKSP> = 22; <TAB> = 23;
		<AD01> = 24; <AD02> = 25; <AD03> = 26; <AD04> = 27;
		<AD05> = 28; <AD06> = 29; <AD07> = 30; <AD08> = 31;
		<AD09> = 32; <AD10> = 33; <AD11> = 34; <AD12> = 35;
		<RTRN> = 36; <LCTL> = 37;
		<AC01> = 38; <AC02> = 39; <AC03> = 40; <AC04> = 41;
		<AC05> = 42; <AC06> = 43; <AC07> = 44; <AC08> = 45;
		<AC09> = 46; <AC10> = 47; <AC11> = 48; <TLDE> = 49;
		<LFSH> = 50; <BKSL> = 51;
		<AB01> = 52; <AB02> = 53; <AB03> = 54; <AB04> = 55;
		<AB05> = 56; <AB06> = 57; <AB07> = 58; <AB08> = 59;
		<AB09> = 60; <AB10> = 61; <RTSH> = 62;
		<KPMU> = 63; <LALT> = 64; <SPCE> = 65; <CAPS> = 66;
		<FK01> = 67; <FK02> = 68; <FK03> = 69; <FK04> = 70;
		<FK05> = 71; <FK06> = 72; <FK07> = 73; <FK08> = 74;
		<FK09> = 75; <FK10> = 76;
		<NMLK> = 77; <SCLK> = 78;
		<KP7> = 79; <KP8> = 80; <KP9> = 81; <KPSU> = 82;
		<KP4> = 83; <KP5> = 84; <KP6> = 85; <KPAD> = 86;
		<KP1> = 87; <KP2> = 88; <KP3> = 89; <KP0> = 90;
		<KPDL> = 91;
		<FK11> = 95; <FK12> = 96;
		<KPEN> = 104; <RCTL> = 105; <KPDV> = 106; <PRSC> = 107;
		<RALT> = 108; <HOME> = 110; <UP> = 111; <PGUP> = 112;
		<LEFT> = 113; <RGHT> = 114; <END> = 115; <DOWN> = 116;
		<PGDN> = 117; <INS> = 118; <DELE> = 119;
		<LWIN> = 133; <RWIN> = 134; <MENU> = 135;
	};
	xkb_types "complete" { include "complete" };
	xkb_compat "complete" { include "complete" };
	xkb_symbols "us" { include "pc+us+inet(evdev)" };
};
`
	keymapBytes := []byte(keymap)
	keymapBytes = append(keymapBytes, 0) // NUL terminate

	dir := os.Getenv("XDG_RUNTIME_DIR")
	if dir == "" {
		dir = os.TempDir()
	}

	f, err := os.CreateTemp(dir, "robotgo-keymap-*")
	if err != nil {
		return fmt.Errorf("create keymap temp file: %w", err)
	}
	defer f.Close()
	defer os.Remove(f.Name())

	if _, err := f.Write(keymapBytes); err != nil {
		return fmt.Errorf("write keymap: %w", err)
	}

	// The fd must remain valid when sent to the compositor
	if err := c.keyboard.Keymap(1, int(f.Fd()), uint32(len(keymapBytes))); err != nil {
		return fmt.Errorf("send keymap: %w", err)
	}

	c.keymapSet = true
	return nil
}

// Close shuts down the Wayland connection. After Close, a subsequent call into
// the backend re-establishes a fresh connection.
func Close() {
	connMu.Lock()
	defer connMu.Unlock()
	if globalConn == nil {
		return
	}
	c := globalConn
	globalConn = nil
	c.closed.Store(true)

	if c.pointer != nil {
		_ = c.pointer.Destroy()
	}
	if c.keyboard != nil {
		_ = c.keyboard.Destroy()
	}
	if c.pointerManager != nil {
		_ = c.pointerManager.Destroy()
	}
	if c.screencopyMgr != nil {
		_ = c.screencopyMgr.Destroy()
	}
	if c.toplevelMgr != nil {
		_ = c.toplevelMgr.Stop()
	}
	_ = c.display.Context().Close()
}
