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
	"os"
	"path/filepath"
	"strconv"
	"strings"

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

// windowFor resolves pid to a window: with isXid (the Cgo NotPid / extra arg
// mode) pid is already an X window id, otherwise it is resolved like
// targetWindow.
func (c *conn) windowFor(pid int, isXid bool) (xproto.Window, error) {
	if !isXid {
		return c.targetWindow(pid)
	}
	if pid <= 0 {
		return 0, ErrNotFound
	}
	return xproto.Window(pid), nil
}

// ActivePid activates the first window owned by pid, or the window pid
// itself when isXid.
func ActivePid(pid int, isXid bool) error {
	if pid <= 0 {
		return ErrNotFound
	}
	c, err := ensureConn()
	if err != nil {
		return err
	}
	w, err := c.windowFor(pid, isXid)
	if err != nil {
		return err
	}
	return ewmh.ActiveWindowReq(c.xu, w)
}

// clientRect returns the client area of w in root coordinates.
func (c *conn) clientRect(w xproto.Window) (x, y, width, height int, err error) {
	geom, err := xproto.GetGeometry(c.c, xproto.Drawable(w)).Reply()
	if err != nil {
		return 0, 0, 0, 0, err
	}
	tr, err := xproto.TranslateCoordinates(c.c, w, c.root, 0, 0).Reply()
	if err != nil {
		return 0, 0, 0, 0, err
	}
	if !tr.SameScreen {
		return 0, 0, 0, 0, ErrNotFound
	}
	return int(tr.DstX), int(tr.DstY), int(geom.Width), int(geom.Height), nil
}

// GetClient returns the client area (x, y, w, h) of pid's window (an X
// window id when isXid); pid <= 0 selects the active window. It returns
// zeros when no window is found.
func GetClient(pid int, isXid bool) (int, int, int, int) {
	c, err := ensureConn()
	if err != nil {
		return 0, 0, 0, 0
	}
	w, err := c.windowFor(pid, isXid)
	if err != nil {
		return 0, 0, 0, 0
	}
	x, y, width, height, err := c.clientRect(w)
	if err != nil {
		return 0, 0, 0, 0
	}
	return x, y, width, height
}

// GetBounds returns the window bounds (x, y, w, h) of pid's window including
// the window manager frame (_NET_FRAME_EXTENTS); pid is an X window id when
// isXid, pid <= 0 selects the active window. It returns zeros when no window
// is found.
func GetBounds(pid int, isXid bool) (int, int, int, int) {
	c, err := ensureConn()
	if err != nil {
		return 0, 0, 0, 0
	}
	w, err := c.windowFor(pid, isXid)
	if err != nil {
		return 0, 0, 0, 0
	}
	x, y, width, height, err := c.clientRect(w)
	if err != nil {
		return 0, 0, 0, 0
	}
	// Undecorated windows or WMs without the property have no frame.
	if ext, err := ewmh.FrameExtentsGet(c.xu, w); err == nil {
		x, y = x-ext.Left, y-ext.Top
		width, height = width+ext.Left+ext.Right, height+ext.Top+ext.Bottom
	}
	return x, y, width, height
}

// MinWindow minimizes (or restores) the window owned by pid.
//
//	MinWindow(pid)        // minimize
//	MinWindow(pid, false) // restore
func MinWindow(pid int, args ...interface{}) error {
	c, err := ensureConn()
	if err != nil {
		return err
	}
	state := true
	if len(args) > 0 {
		if v, ok := args[0].(bool); ok {
			state = v
		}
	}
	w, err := c.xidByPid(pid)
	if err != nil {
		return err
	}
	if state {
		// WM_CHANGE_STATE -> IconicState minimizes via the window manager.
		return ewmh.ClientEvent(c.xu, w, "WM_CHANGE_STATE", icccm.StateIconic)
	}
	return ewmh.ActiveWindowReq(c.xu, w)
}

// MaxWindow maximizes (or unmaximizes) the window owned by pid.
//
//	MaxWindow(pid)        // maximize
//	MaxWindow(pid, false) // unmaximize
func MaxWindow(pid int, args ...interface{}) error {
	c, err := ensureConn()
	if err != nil {
		return err
	}
	state := true
	if len(args) > 0 {
		if v, ok := args[0].(bool); ok {
			state = v
		}
	}
	w, err := c.xidByPid(pid)
	if err != nil {
		return err
	}
	action := ewmh.StateAdd
	if !state {
		action = ewmh.StateRemove
	}
	return ewmh.WmStateReqExtra(c.xu, w, action,
		"_NET_WM_STATE_MAXIMIZED_VERT", "_NET_WM_STATE_MAXIMIZED_HORZ", 2)
}

// CloseWindow closes the window. With no argument it closes the active window;
// otherwise it closes the first window owned by pid (args[0]).
func CloseWindow(args ...int) error {
	c, err := ensureConn()
	if err != nil {
		return err
	}
	pid := 0
	if len(args) > 0 {
		pid = args[0]
	}
	w, err := c.targetWindow(pid)
	if err != nil {
		return err
	}
	return ewmh.CloseWindow(c.xu, w)
}

// CheckAccess reports whether input injection is permitted. X11 has no
// accessibility gate, so it always returns true; prompt is ignored.
func CheckAccess(prompt bool) bool { return true }

// GetActiveApp returns the active window's process name, executable path and
// pid (from _NET_WM_PID and /proc). Unresolved fields are left empty.
func GetActiveApp() (string, string, int) {
	c, err := ensureConn()
	if err != nil {
		return "", "", 0
	}
	w, err := ewmh.ActiveWindowGet(c.xu)
	if err != nil {
		return "", "", 0
	}
	wmPid, err := ewmh.WmPidGet(c.xu, w)
	if err != nil || wmPid == 0 {
		return "", "", 0
	}
	return procInfo(int(wmPid))
}

// procInfo returns the name and executable path of pid read from /proc.
// The name prefers the exe basename, since comm is truncated to 15 bytes.
func procInfo(pid int) (string, string, int) {
	dir := "/proc/" + strconv.Itoa(pid)
	path, err := os.Readlink(dir + "/exe")
	if err != nil {
		path = ""
	} else if name := filepath.Base(path); name != "." && name != "/" {
		return name, path, pid
	}
	comm, err := os.ReadFile(dir + "/comm")
	if err != nil {
		return "", path, pid
	}
	return strings.TrimSpace(string(comm)), path, pid
}
