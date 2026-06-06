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
)

// Window management via zwlr_foreign_toplevel_management_v1.
// Only available on wlroots-based compositors.

// GetTitle returns the title of the active (or specified) window.
func GetTitle(args ...int) string {
	c, err := ensureConn()
	if err != nil || c.toplevelMgr == nil {
		return ""
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	// Find the activated toplevel
	for _, info := range c.toplevels {
		if isActivated(info.states) {
			return info.title
		}
	}

	// Fallback: return first toplevel's title
	for _, info := range c.toplevels {
		return info.title
	}
	return ""
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

	c.mu.Lock()
	defer c.mu.Unlock()

	nameLower := strings.ToLower(name)
	for _, info := range c.toplevels {
		if strings.Contains(strings.ToLower(info.title), nameLower) ||
			strings.Contains(strings.ToLower(info.appId), nameLower) {
			return info.handle.Activate(c.seat)
		}
	}
	return ErrNotSupported
}

// MinWindow minimizes a window. If state is true, minimize; if false, unminimize.
func MinWindow(pid int, args ...interface{}) {
	c, err := ensureConn()
	if err != nil || c.toplevelMgr == nil {
		return
	}

	minimize := true
	if len(args) > 0 {
		if b, ok := args[0].(bool); ok {
			minimize = b
		}
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	// Find toplevel by matching (we don't have PID info from the protocol,
	// so we operate on all toplevels or the active one)
	for _, info := range c.toplevels {
		if minimize {
			_ = info.handle.SetMinimized()
		} else {
			_ = info.handle.UnsetMinimized()
		}
		return // operate on first match
	}
}

// MaxWindow maximizes a window.
func MaxWindow(pid int, args ...interface{}) {
	c, err := ensureConn()
	if err != nil || c.toplevelMgr == nil {
		return
	}

	maximize := true
	if len(args) > 0 {
		if b, ok := args[0].(bool); ok {
			maximize = b
		}
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	for _, info := range c.toplevels {
		if maximize {
			_ = info.handle.SetMaximized()
		} else {
			_ = info.handle.UnsetMaximized()
		}
		return
	}
}

// CloseWindow closes a window.
func CloseWindow(args ...int) {
	c, err := ensureConn()
	if err != nil || c.toplevelMgr == nil {
		return
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	// Close the active toplevel
	for _, info := range c.toplevels {
		if isActivated(info.states) {
			_ = info.handle.Close()
			return
		}
	}
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
