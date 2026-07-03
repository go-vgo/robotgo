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

package x11

import (
	"github.com/jezek/xgb/xproto"
	"github.com/jezek/xgbutil/ewmh"
	"github.com/jezek/xgbutil/icccm"
)

// xidByPid returns the first managed window owned by pid (EWMH client list).
func (c *conn) xidByPid(pid int) (xproto.Window, error) {
	wins, err := ewmh.ClientListGet(c.xu)
	if err != nil {
		return 0, err
	}
	for _, w := range wins {
		if wmPid, err := ewmh.WmPidGet(c.xu, w); err == nil && uint(pid) == wmPid {
			return w, nil
		}
	}
	return 0, ErrNotFound
}

// windowName returns a window's title, preferring EWMH _NET_WM_NAME and falling
// back to ICCCM WM_NAME.
func (c *conn) windowName(w xproto.Window) string {
	if name, err := ewmh.WmNameGet(c.xu, w); err == nil && name != "" {
		return name
	}
	if name, err := icccm.WmNameGet(c.xu, w); err == nil {
		return name
	}
	return ""
}

// targetWindow resolves the window referenced by an optional pid; pid <= 0
// selects the currently active window.
func (c *conn) targetWindow(pid int) (xproto.Window, error) {
	if pid <= 0 {
		return ewmh.ActiveWindowGet(c.xu)
	}
	return c.xidByPid(pid)
}

// GetTitle returns the window title. With no argument (or pid <= 0) it returns
// the active window's title; otherwise the title of the first window for pid.
func GetTitle(args ...int) string {
	c, err := ensureConn()
	if err != nil {
		return ""
	}
	pid := 0
	if len(args) > 0 {
		pid = args[0]
	}
	w, err := c.targetWindow(pid)
	if err != nil {
		return ""
	}
	return c.windowName(w)
}

// ActiveName activates the first window whose owning process name matches.
func ActiveName(name string) error {
	c, err := ensureConn()
	if err != nil {
		return err
	}
	pids, err := FindIds(name)
	if err != nil {
		return err
	}
	if len(pids) == 0 {
		return ErrNotFound
	}
	w, err := c.xidByPid(pids[0])
	if err != nil {
		return err
	}
	return ewmh.ActiveWindowReq(c.xu, w)
}

// MinWindow minimizes (or restores) the window owned by pid.
//
//	MinWindow(pid)        // minimize
//	MinWindow(pid, false) // restore
func MinWindow(pid int, args ...interface{}) {
	c, err := ensureConn()
	if err != nil {
		return
	}
	state := true
	if len(args) > 0 {
		if v, ok := args[0].(bool); ok {
			state = v
		}
	}
	w, err := c.xidByPid(pid)
	if err != nil {
		return
	}
	if state {
		// WM_CHANGE_STATE -> IconicState minimizes via the window manager.
		_ = ewmh.ClientEvent(c.xu, w, "WM_CHANGE_STATE", icccm.StateIconic)
		return
	}
	_ = ewmh.ActiveWindowReq(c.xu, w)
}

// MaxWindow maximizes (or unmaximizes) the window owned by pid.
//
//	MaxWindow(pid)        // maximize
//	MaxWindow(pid, false) // unmaximize
func MaxWindow(pid int, args ...interface{}) {
	c, err := ensureConn()
	if err != nil {
		return
	}
	state := true
	if len(args) > 0 {
		if v, ok := args[0].(bool); ok {
			state = v
		}
	}
	w, err := c.xidByPid(pid)
	if err != nil {
		return
	}
	action := ewmh.StateAdd
	if !state {
		action = ewmh.StateRemove
	}
	_ = ewmh.WmStateReqExtra(c.xu, w, action,
		"_NET_WM_STATE_MAXIMIZED_VERT", "_NET_WM_STATE_MAXIMIZED_HORZ", 2)
}

// CloseWindow closes the window. With no argument it closes the active window;
// otherwise it closes the first window owned by pid (args[0]).
func CloseWindow(args ...int) {
	c, err := ensureConn()
	if err != nil {
		return
	}
	pid := 0
	if len(args) > 0 {
		pid = args[0]
	}
	w, err := c.targetWindow(pid)
	if err != nil {
		return
	}
	_ = ewmh.CloseWindow(c.xu, w)
}
