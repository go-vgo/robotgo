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
	"time"

	"github.com/jezek/xgb/xproto"
	"github.com/jezek/xgb/xtest"
)

// KeySleep is the global keyboard delay in milliseconds, applied after a key
// action completes.
var KeySleep = 10

// keyDelay is the time spent between a key press and its release.
const keyDelay = 5 * time.Millisecond

// sendKeycode generates a press or release for the given keycode via XTEST.
func (c *conn) sendKeycode(kc xproto.Keycode, press bool) {
	t := byte(xproto.KeyRelease)
	if press {
		t = byte(xproto.KeyPress)
	}
	xtest.FakeInput(c.c, t, byte(kc), 0, c.root, 0, 0, 0)
}

// pressKeysym presses and releases a keysym, holding Shift if the keysym sits
// in the shifted column of its keycode.
func (c *conn) pressKeysym(ks uint32) error {
	kc, shift, ok := c.keysymToKeycode(ks)
	if ok {
		if shift && c.shiftKeycode != 0 {
			c.sendKeycode(c.shiftKeycode, true)
		}
		c.sendKeycode(kc, true)
		time.Sleep(keyDelay)
		c.sendKeycode(kc, false)
		if shift && c.shiftKeycode != 0 {
			c.sendKeycode(c.shiftKeycode, false)
		}
		c.c.Sync()
		return nil
	}
	// Not in the current layout — type it through the scratch keycode.
	return c.pressScratch(ks)
}

// pressScratch temporarily remaps the spare keycode to the keysym, taps it, and
// restores it. This is how arbitrary Unicode characters are typed.
func (c *conn) pressScratch(ks uint32) error {
	if !c.scratchOK {
		return ErrNotSupported
	}
	per := c.keysymsPerKeycode
	syms := make([]xproto.Keysym, per)
	for i := range syms {
		syms[i] = xproto.Keysym(ks)
	}

	if err := xproto.ChangeKeyboardMappingChecked(
		c.c, 1, c.scratch, per, syms).Check(); err != nil {
		return err
	}
	c.sync()

	c.sendKeycode(c.scratch, true)
	time.Sleep(keyDelay)
	c.sendKeycode(c.scratch, false)
	c.c.Sync()

	// Restore the scratch keycode to NoSymbol so we leave the map as we found it.
	for i := range syms {
		syms[i] = 0
	}
	_ = xproto.ChangeKeyboardMappingChecked(c.c, 1, c.scratch, per, syms).Check()
	c.sync()
	return nil
}

// modKeycodes resolves modifier names ("ctrl", "shift", "alt", "cmd", ...) to
// keycodes, skipping any that cannot be mapped.
func (c *conn) modKeycodes(mods []string) []xproto.Keycode {
	var out []xproto.Keycode
	for _, m := range mods {
		ks, ok := keyKeysym(m)
		if !ok {
			continue
		}
		if kc, _, ok := c.keysymToKeycode(ks); ok {
			out = append(out, kc)
		}
	}
	return out
}

// extractMods flattens KeyTap/KeyToggle variadic args into a list of modifier
// names. It accepts string, []string and []interface{} forms, and recognizes a
// leading "up"/"down" direction which is returned separately.
func extractMods(args []interface{}) (mods []string, down bool, hasDir bool) {
	down = true
	for _, a := range args {
		switch v := a.(type) {
		case string:
			if v == "up" {
				down, hasDir = false, true
			} else if v == "down" {
				down, hasDir = true, true
			} else {
				mods = append(mods, v)
			}
		case []string:
			mods = append(mods, v...)
		case []interface{}:
			m, _, _ := extractMods(v)
			mods = append(mods, m...)
		}
	}
	return mods, down, hasDir
}

// KeyTap taps a key, optionally with modifiers.
//
//	KeyTap("a")
//	KeyTap("c", "ctrl")
//	KeyTap("a", "ctrl", "shift")
//	KeyTap("a", []string{"ctrl", "shift"})
func KeyTap(key string, args ...interface{}) error {
	c, err := ensureConn()
	if err != nil {
		return err
	}
	c.mu.Lock()
	defer c.mu.Unlock()

	ks, ok := keyKeysym(key)
	if !ok {
		return ErrNotFound
	}

	mods, _, _ := extractMods(args)
	mkc := c.modKeycodes(mods)

	for _, kc := range mkc {
		c.sendKeycode(kc, true)
	}
	c.c.Sync()

	perr := c.pressKeysym(ks)

	for i := len(mkc) - 1; i >= 0; i-- {
		c.sendKeycode(mkc[i], false)
	}
	c.c.Sync()

	keySleep()
	return perr
}

// KeyPress is an alias of KeyTap.
func KeyPress(key string, args ...interface{}) error {
	return KeyTap(key, args...)
}

// KeyToggle presses or releases a key (and optional held modifiers).
//
//	KeyToggle("a", "down")
//	KeyToggle("a", "up")
//	KeyToggle("ctrl", "down")
func KeyToggle(key string, args ...interface{}) error {
	c, err := ensureConn()
	if err != nil {
		return err
	}
	c.mu.Lock()
	defer c.mu.Unlock()

	ks, ok := keyKeysym(key)
	if !ok {
		return ErrNotFound
	}
	kc, shift, ok := c.keysymToKeycode(ks)
	if !ok {
		return ErrNotFound
	}

	mods, down, _ := extractMods(args)
	mkc := c.modKeycodes(mods)

	if down {
		for _, m := range mkc {
			c.sendKeycode(m, true)
		}
		if shift && c.shiftKeycode != 0 {
			c.sendKeycode(c.shiftKeycode, true)
		}
		c.sendKeycode(kc, true)
	} else {
		c.sendKeycode(kc, false)
		if shift && c.shiftKeycode != 0 {
			c.sendKeycode(c.shiftKeycode, false)
		}
		for i := len(mkc) - 1; i >= 0; i-- {
			c.sendKeycode(mkc[i], false)
		}
	}
	c.c.Sync()
	keySleep()
	return nil
}

// KeyDown presses a key down (and holds it).
func KeyDown(key string, args ...interface{}) error {
	args = append([]interface{}{}, args...)
	args = append(args, "down")
	return KeyToggle(key, args...)
}

// KeyUp releases a previously held key.
func KeyUp(key string, args ...interface{}) error {
	args = append([]interface{}{}, args...)
	args = append(args, "up")
	return KeyToggle(key, args...)
}

// Type types a string (alias of TypeStr).
func Type(str string, args ...int) {
	TypeStr(str, args...)
}

// TypeStr types a string of (possibly Unicode) characters.
func TypeStr(str string, args ...int) {
	c, err := ensureConn()
	if err != nil {
		return
	}
	c.mu.Lock()
	defer c.mu.Unlock()

	for _, r := range str {
		_ = c.pressKeysym(runeKeysym(r))
		if KeySleep > 0 {
			time.Sleep(time.Duration(KeySleep) * time.Millisecond)
		}
	}
}

// TypeDelay types a string then sleeps for delay milliseconds.
func TypeDelay(str string, delay int) {
	Type(str)
	MilliSleep(delay)
}

// SetDelay sets the default keyboard and mouse delay (default 10).
func SetDelay(d ...int) {
	v := 10
	if len(d) > 0 {
		v = d[0]
	}
	KeySleep = v
	MouseSleep = v
}

func keySleep() {
	if KeySleep > 0 {
		time.Sleep(time.Duration(KeySleep) * time.Millisecond)
	}
}
