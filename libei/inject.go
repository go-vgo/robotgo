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

package libei

import (
	"fmt"

	"github.com/godbus/dbus/v5"
)

// notifyInjector drives input through the RemoteDesktop portal Notify* methods.
// Each call is a synchronous D-Bus method call on the session object; no
// Request/Response round-trip is involved for these (only the setup handshake
// in conn.go uses that pattern).
type notifyInjector struct {
	c *conn
}

// emptyOpts is the (a{sv}) options argument the Notify* methods require. The
// portal does not define any options for them today.
func emptyOpts() map[string]dbus.Variant { return map[string]dbus.Variant{} }

func (n *notifyInjector) call(method string, args ...interface{}) error {
	full := append([]interface{}{n.c.sessionHandle, emptyOpts()}, args...)
	call := n.c.obj.Call(ifaceRemoteDesktop+"."+method, 0, full...)
	if call.Err != nil {
		return fmt.Errorf("robotgo/libei: %s: %w", method, call.Err)
	}
	return nil
}

func (n *notifyInjector) keyboardKeycode(code int32, state uint32) error {
	return n.call("NotifyKeyboardKeycode", code, state)
}

func (n *notifyInjector) keyboardKeysym(sym int32, state uint32) error {
	return n.call("NotifyKeyboardKeysym", sym, state)
}

func (n *notifyInjector) pointerMotion(dx, dy float64) error {
	return n.call("NotifyPointerMotion", dx, dy)
}

func (n *notifyInjector) pointerButton(button int32, state uint32) error {
	return n.call("NotifyPointerButton", button, state)
}

func (n *notifyInjector) pointerAxisDiscrete(axis uint32, steps int32) error {
	return n.call("NotifyPointerAxisDiscrete", axis, steps)
}

func (n *notifyInjector) pointerMotionAbsolute(streamNode uint32, x, y float64) error {
	return n.call("NotifyPointerMotionAbsolute", streamNode, x, y)
}
