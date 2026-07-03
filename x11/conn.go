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
	"fmt"
	"sync"

	"github.com/jezek/xgb"
	"github.com/jezek/xgb/xinerama"
	"github.com/jezek/xgb/xproto"
	"github.com/jezek/xgb/xtest"
	"github.com/jezek/xgbutil"
)

// conn holds the singleton X11 connection plus the bits we cache from it: the
// root window, the keyboard mapping (for keysym -> keycode resolution) and a
// spare ("scratch") keycode used to type arbitrary Unicode characters.
type conn struct {
	mu sync.Mutex

	xu   *xgbutil.XUtil
	c    *xgb.Conn
	root xproto.Window

	// keyboard mapping cache
	minKeycode        xproto.Keycode
	keysymsPerKeycode byte
	keysyms           []xproto.Keysym

	shiftKeycode xproto.Keycode // keycode that produces Shift_L (0 if none)
	scratch      xproto.Keycode // spare keycode used for Unicode typing
	scratchOK    bool

	xineramaOK bool
}

var (
	globalConn *conn
	connMu     sync.Mutex
)

// ensureConn lazily establishes the global X11 connection, re-establishing it
// if it was never created or has since been closed. A mutex (instead of
// sync.Once) keeps the backend recoverable after a failed connect or Close.
func ensureConn() (*conn, error) {
	connMu.Lock()
	defer connMu.Unlock()

	if globalConn != nil {
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

// newConn opens the X display, initializes XTEST and caches the keyboard map.
func newConn() (*conn, error) {
	xu, err := xgbutil.NewConn()
	if err != nil {
		return nil, fmt.Errorf("robotgo: connect to X11: %w", err)
	}

	c := &conn{
		xu:   xu,
		c:    xu.Conn(),
		root: xu.RootWin(),
	}

	// XTEST is required for synthetic input.
	if err := xtest.Init(c.c); err != nil {
		c.c.Close()
		return nil, fmt.Errorf("robotgo: %w: XTEST extension (%v)", ErrNotSupported, err)
	}

	// Xinerama is optional; used for per-monitor geometry.
	if err := xinerama.Init(c.c); err == nil {
		c.xineramaOK = true
	}

	if err := c.loadKeymap(); err != nil {
		c.c.Close()
		return nil, err
	}

	return c, nil
}

// loadKeymap fetches the full keyboard mapping and locates the Shift keycode
// and a spare keycode for Unicode typing.
func (c *conn) loadKeymap() error {
	setup := c.xu.Setup()
	c.minKeycode = setup.MinKeycode
	count := byte(setup.MaxKeycode-setup.MinKeycode) + 1

	reply, err := xproto.GetKeyboardMapping(c.c, c.minKeycode, count).Reply()
	if err != nil {
		return fmt.Errorf("robotgo: get keyboard mapping: %w", err)
	}
	c.keysymsPerKeycode = reply.KeysymsPerKeycode
	c.keysyms = reply.Keysyms

	if kc, _, ok := c.keysymToKeycode(xkShiftL); ok {
		c.shiftKeycode = kc
	}
	c.scratch, c.scratchOK = c.findScratchKeycode()
	return nil
}

// sync forces the server to process queued requests (used after remapping a
// keycode, before generating events for it).
func (c *conn) sync() {
	c.c.Sync()
}

// Close shuts down the X11 connection. After Close, a subsequent call into the
// backend re-establishes a fresh connection.
func Close() {
	connMu.Lock()
	c := globalConn
	globalConn = nil
	connMu.Unlock()

	if c == nil {
		return
	}
	// Take the per-connection lock so an in-flight operation (KeyTap,
	// TypeStr, ...) finishes before the underlying transport is closed.
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.c != nil {
		c.c.Close()
	}
}
