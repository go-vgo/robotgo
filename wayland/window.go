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
	"encoding/binary"
	"strings"

	"github.com/go-vgo/robotgo/wayland/internal/protocols/wlr_foreign_toplevel"
)

// Window management via zwlr_foreign_toplevel_management_v1.
// Only available on wlroots-based compositors.

// activeToplevel returns the activated toplevel, falling back to any toplevel
// when none is activated. The lookup runs under c.mu; the returned handle is
// used only after the lock is released (requests need c.wl, which the
// dispatch goroutine holds while it updates c.toplevels under c.mu).
func (c *conn) activeToplevel() *toplevelInfo {
	c.mu.Lock()
	defer c.mu.Unlock()
	var any *toplevelInfo
	for _, info := range c.toplevels {
		if isActivated(info.states) {
			return info
		}
		any = info
	}
	return any
}

// GetTitle returns the title of the active (or specified) window.
func GetTitle(args ...int) string {
	c, err := ensureConn()
	if err != nil || c.toplevelMgr == nil {
		return ""
	}
	info := c.activeToplevel()
	if info == nil {
		return ""
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	return info.title
}

// ActiveName activates a window by matching process/app name.
func ActiveName(name string) error {
	c, err := ensureConn()
	if err != nil {
		return err
	}
	if c.toplevelMgr == nil || c.seat == nil {
		return ErrNotSupported
	}

	nameLower := strings.ToLower(name)
	var handle *wlr_foreign_toplevel.ZwlrForeignToplevelHandleV1
	c.mu.Lock()
	for _, info := range c.toplevels {
		if strings.Contains(strings.ToLower(info.title), nameLower) ||
			strings.Contains(strings.ToLower(info.appId), nameLower) {
			handle = info.handle
			break
		}
	}
	c.mu.Unlock()
	if handle == nil {
		return ErrNotSupported
	}
	return c.do(func() error { return handle.Activate(c.seat) })
}

// withActive runs fn on the active toplevel's handle, if any.
func withActive(fn func(h *wlr_foreign_toplevel.ZwlrForeignToplevelHandleV1) error) {
	c, err := ensureConn()
	if err != nil || c.toplevelMgr == nil {
		return
	}
	info := c.activeToplevel()
	if info == nil {
		return
	}
	_ = c.do(func() error { return fn(info.handle) })
}

// MinWindow minimizes the active window. If the first arg is false, unminimize.
// The pid is accepted for API parity; the protocol exposes no pid mapping.
func MinWindow(pid int, args ...interface{}) {
	minimize := true
	if len(args) > 0 {
		if b, ok := args[0].(bool); ok {
			minimize = b
		}
	}
	withActive(func(h *wlr_foreign_toplevel.ZwlrForeignToplevelHandleV1) error {
		if minimize {
			return h.SetMinimized()
		}
		return h.UnsetMinimized()
	})
}

// MaxWindow maximizes the active window. If the first arg is false, unmaximize.
// The pid is accepted for API parity; the protocol exposes no pid mapping.
func MaxWindow(pid int, args ...interface{}) {
	maximize := true
	if len(args) > 0 {
		if b, ok := args[0].(bool); ok {
			maximize = b
		}
	}
	withActive(func(h *wlr_foreign_toplevel.ZwlrForeignToplevelHandleV1) error {
		if maximize {
			return h.SetMaximized()
		}
		return h.UnsetMaximized()
	})
}

// CloseWindow closes the active window.
func CloseWindow(args ...int) {
	withActive(func(h *wlr_foreign_toplevel.ZwlrForeignToplevelHandleV1) error {
		return h.Close()
	})
}

// isActivated checks if the toplevel state array contains the "activated" state.
// The state is an array of uint32 values, where 2 = activated.
func isActivated(states []byte) bool {
	for i := 0; i+3 < len(states); i += 4 {
		state := binary.LittleEndian.Uint32(states[i : i+4])
		if state == 2 { // activated
			return true
		}
	}
	return false
}
