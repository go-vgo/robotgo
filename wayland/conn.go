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
	"golang.org/x/sys/unix"

	"github.com/go-vgo/robotgo/wayland/internal/protocols/wlr_foreign_toplevel"
	"github.com/go-vgo/robotgo/wayland/internal/protocols/wlr_screencopy"
	"github.com/go-vgo/robotgo/wayland/internal/protocols/wlr_virtual_keyboard"
	"github.com/go-vgo/robotgo/wayland/internal/protocols/wlr_virtual_pointer"
)

// conn holds the singleton Wayland connection and all bound protocol objects.
type conn struct {
	// mu guards the Go-side window/output bookkeeping (toplevels, outputs).
	mu sync.Mutex

	// wl serializes every interaction with the go-wayland Context: the
	// library keeps an unsynchronized object table, so requests that create
	// or destroy proxies must never run concurrently with event dispatch.
	// dispatchLoop holds it while delivering each event, so event handlers
	// (which already run under it) call proxies directly; every other
	// goroutine goes through c.do.
	wl sync.Mutex

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
	// mods is the XKB modifier mask currently held down via the virtual
	// keyboard; it is reported to the compositor with the modifiers request.
	mods uint32

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
			c.mu.Lock()
			info.x = int32(ge.X)
			info.y = int32(ge.Y)
			c.mu.Unlock()
		})
		out.SetModeHandler(func(me client.OutputModeEvent) {
			if me.Flags&0x1 != 0 { // WL_OUTPUT_MODE_CURRENT
				c.mu.Lock()
				info.width = int32(me.Width)
				info.height = int32(me.Height)
				c.mu.Unlock()
			}
		})
		c.mu.Lock()
		c.outputs = append(c.outputs, info)
		c.mu.Unlock()

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
		// The handle is inert after closed; release it so the proxy slot and
		// server resource do not leak for every window that comes and goes.
		if err := handle.Destroy(); err != nil {
			log.Printf("robotgo: destroy toplevel handle: %v", err)
		}
	})
}

// roundtrip performs a blocking wl_display.sync roundtrip. It is only used
// during setup, before the dispatch loop starts.
func (c *conn) roundtrip() {
	if err := c.display.Roundtrip(); err != nil {
		log.Printf("robotgo: wayland roundtrip: %v", err)
	}
}

// do runs fn with exclusive access to the Wayland context. It returns
// ErrNoConnection once the connection has been closed or its dispatch loop
// has died, so callers never block on a dead compositor.
func (c *conn) do(fn func() error) error {
	if c.closed.Load() {
		return ErrNoConnection
	}
	select {
	case <-c.dispatchDone:
		return ErrNoConnection
	default:
	}
	c.wl.Lock()
	defer c.wl.Unlock()
	return fn()
}

// dispatchLoop reads Wayland events in a background goroutine. The socket read
// happens unlocked; the actual dispatch (which touches the shared object table
// and runs handlers) is serialized with c.do via c.wl.
func (c *conn) dispatchLoop() {
	defer close(c.dispatchDone)
	ctx := c.display.Context()
	for {
		dispatch := ctx.GetDispatch()
		c.wl.Lock()
		err := dispatch()
		c.wl.Unlock()
		if err != nil {
			if !c.closed.Load() {
				log.Printf("robotgo: dispatch error: %v", err)
			}
			return
		}
	}
}

// outputBounds returns the union rectangle of all outputs, which is the
// coordinate space wlr_virtual_pointer.motion_absolute maps onto when the
// pointer is not bound to a single output.
func (c *conn) outputBounds() (x, y, w, h int32) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if len(c.outputs) == 0 {
		return 0, 0, 0, 0
	}
	minX, minY := c.outputs[0].x, c.outputs[0].y
	maxX, maxY := minX, minY
	for _, o := range c.outputs {
		minX = min(minX, o.x)
		minY = min(minY, o.y)
		maxX = max(maxX, o.x+o.width)
		maxY = max(maxY, o.y+o.height)
	}
	return minX, minY, maxX - minX, maxY - minY
}

// output returns a snapshot of output idx (0 when out of range) and whether
// any output exists.
func (c *conn) output(idx int) (outputInfo, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if len(c.outputs) == 0 {
		return outputInfo{}, false
	}
	if idx < 0 || idx >= len(c.outputs) {
		idx = 0
	}
	return *c.outputs[idx], true
}

// timestamp returns a Wayland-compatible millisecond timestamp.
func timestamp() uint32 {
	return uint32(time.Now().UnixMilli())
}

// keymap is the XKB keymap installed on the virtual keyboard. It is the same
// shape `setxkbmap -print` produces: the full evdev keycode set (so every
// code in evdevKeyMap — F13-F24, media keys, pause, ... — resolves) with the
// US symbols the key table assumes.
const keymap = `xkb_keymap {
	xkb_keycodes { include "evdev+aliases(qwerty)" };
	xkb_types { include "complete" };
	xkb_compat { include "complete" };
	xkb_symbols { include "pc+us+inet(evdev)" };
	xkb_geometry { include "pc(pc105)" };
};
`

// setupKeymap shares the XKB keymap with the compositor through a sealed
// memfd, as the protocol expects (the compositor mmaps the fd read-only).
// This must be done before any key events.
func (c *conn) setupKeymap() error {
	data := append([]byte(keymap), 0) // NUL terminate

	fd, err := unix.MemfdCreate("robotgo-keymap", unix.MFD_CLOEXEC|unix.MFD_ALLOW_SEALING)
	if err != nil {
		return fmt.Errorf("memfd_create keymap: %w", err)
	}
	f := os.NewFile(uintptr(fd), "robotgo-keymap")
	defer f.Close()

	if _, err := f.Write(data); err != nil {
		return fmt.Errorf("write keymap: %w", err)
	}
	if _, err := unix.FcntlInt(uintptr(fd), unix.F_ADD_SEALS,
		unix.F_SEAL_SHRINK|unix.F_SEAL_GROW|unix.F_SEAL_WRITE|unix.F_SEAL_SEAL); err != nil {
		return fmt.Errorf("seal keymap: %w", err)
	}

	// The fd is duplicated onto the socket by SCM_RIGHTS, so it can be
	// closed once the request has been written.
	if err := c.keyboard.Keymap(1, fd, uint32(len(data))); err != nil {
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

	c.wl.Lock()
	defer c.wl.Unlock()
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
	// Closing the socket unblocks the dispatch goroutine's pending read.
	_ = c.display.Context().Close()
}
