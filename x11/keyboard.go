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
	"errors"
	"fmt"
	"time"
	"unicode/utf8"

	"github.com/go-vgo/robotgo/pub"
	"github.com/jezek/xgb/xproto"
)

// keyDelay is the time spent between a key press and its release.
const keyDelay = 5 * time.Millisecond

// Modifier bits returned by keysymToKeycode.
const (
	modShift  uint8 = 1 << 0
	modLevel3 uint8 = 1 << 1 // ISO_Level3_Shift / Mode_switch (AltGr)
)

// modKeycodesFor lists the keycodes to hold for a modifier mask, skipping
// modifiers the layout has no key for.
func (c *conn) modKeycodesFor(mods uint8) []xproto.Keycode {
	var out []xproto.Keycode
	if mods&modLevel3 != 0 && c.level3Keycode != 0 {
		out = append(out, c.level3Keycode)
	}
	if mods&modShift != 0 && c.shiftKeycode != 0 {
		out = append(out, c.shiftKeycode)
	}
	return out
}

// canReach reports whether layout has keys for modifier mask: AltGr
// level is unreachable when no key selects level 3.
func (c *conn) canReach(mods uint8) bool {
	return mods&modLevel3 == 0 || c.level3Keycode != 0
}

// sendKeycode generates a press or release for the given keycode via XTEST.
func (c *conn) sendKeycode(kc xproto.Keycode, press bool) error {
	t := byte(xproto.KeyRelease)
	if press {
		t = byte(xproto.KeyPress)
	}
	return c.fakeInput(t, byte(kc), 0, 0)
}

// pressKeycodes presses kcs in order. If a press fails, the keys already
// pressed are released so none stays held, and the errors are returned.
func (c *conn) pressKeycodes(kcs []xproto.Keycode) error {
	for i, kc := range kcs {
		if err := c.sendKeycode(kc, true); err != nil {
			return errors.Join(err, c.releaseKeycodes(kcs[:i]))
		}
	}
	return nil
}

// releaseKeycodes releases kcs in reverse order. It keeps going past failures
// so every key gets a release attempt, and returns the errors joined.
func (c *conn) releaseKeycodes(kcs []xproto.Keycode) error {
	var errs []error
	for i := len(kcs) - 1; i >= 0; i-- {
		if err := c.sendKeycode(kcs[i], false); err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}

// pressKeysym presses and releases a keysym, holding the modifiers (Shift,
// AltGr) its level on the layout needs (#640).
func (c *conn) pressKeysym(ks uint32) error {
	kc, mods, ok := c.keysymToKeycode(ks)
	if ok && c.canReach(mods) {
		keys := append(c.modKeycodesFor(mods), kc)
		if err := c.pressKeycodes(keys); err != nil {
			return err
		}
		time.Sleep(keyDelay)
		return c.releaseKeycodes(keys)
	}
	// Not in the current layout, or only on an AltGr level the layout has
	// no key for: type it through the scratch keycode.
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

	err := c.sendKeycode(c.scratch, true)
	if err == nil {
		time.Sleep(keyDelay)
		err = c.sendKeycode(c.scratch, false)
	}

	// Restore the scratch keycode to NoSymbol so we leave the map as we found it.
	for i := range syms {
		syms[i] = 0
	}
	rerr := xproto.ChangeKeyboardMappingChecked(c.c, 1, c.scratch, per, syms).Check()
	c.sync()
	return errors.Join(err, rerr)
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
	add := func(s string) {
		switch s {
		case "up":
			down, hasDir = false, true
		case "down":
			down, hasDir = true, true
		default:
			mods = append(mods, s)
		}
	}
	for _, a := range args {
		switch v := a.(type) {
		case string:
			add(v)
		case []string:
			for _, s := range v {
				add(s)
			}
		case []interface{}:
			m, d, ok := extractMods(v)
			mods = append(mods, m...)
			if ok {
				down, hasDir = d, true
			}
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
	if err := c.loadKeymap(); err != nil {
		return err
	}

	ks, ok := keyKeysym(key)
	if !ok {
		return ErrNotFound
	}

	mods, _, _ := extractMods(args)
	mkc := c.modKeycodes(mods)
	if err := c.pressKeycodes(mkc); err != nil {
		return err
	}

	perr := c.pressKeysym(ks)
	// Release the modifiers even if the key failed, so none stays held.
	rerr := c.releaseKeycodes(mkc)

	keySleep()
	return errors.Join(perr, rerr)
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
	if err := c.loadKeymap(); err != nil {
		return err
	}

	ks, ok := keyKeysym(key)
	if !ok {
		return ErrNotFound
	}
	kc, lvl, ok := c.keysymToKeycode(ks)
	if !ok {
		return ErrNotFound
	}
	if !c.canReach(lvl) {
		return ErrNotSupported
	}
	levelMods := c.modKeycodesFor(lvl)

	mods, down, _ := extractMods(args)
	mkc := c.modKeycodes(mods)

	// Press order: modifiers, level modifiers, key; release is the reverse.
	keys := append(append(mkc, levelMods...), kc)
	if down {
		err = c.pressKeycodes(keys)
	} else {
		err = c.releaseKeycodes(keys)
	}
	keySleep()
	return err
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

// Type types a string of (possibly Unicode) characters and returns the
// number of characters (runes) typed; it stops at the first failed key event.
func Type(str string, args ...int) int {
	c, err := ensureConn()
	if err != nil {
		return 0
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	if err := c.loadKeymap(); err != nil {
		return 0
	}

	n := 0
	for _, r := range str {
		if err := c.pressKeysym(runeKeysym(r)); err != nil {
			return n
		}
		n++
		pub.MilliSleep(pub.KeySleep)
	}
	return n
}

// TypeStr types a string (alias of Type). It returns an error if not every
// character was typed.
func TypeStr(str string, args ...int) error {
	return typeErr(Type(str, args...), str)
}

// TypeDelay types a string then sleeps for delay milliseconds. It returns an
// error if not every character was typed.
func TypeDelay(str string, delay int) error {
	n := Type(str)
	pub.MilliSleep(delay)
	return typeErr(n, str)
}

// typeErr reports an error when fewer than all runes of str were typed.
func typeErr(n int, str string) error {
	if total := utf8.RuneCountInString(str); n < total {
		return fmt.Errorf("robotgo: typed %d of %d characters", n, total)
	}
	return nil
}

func keySleep() {
	pub.MilliSleep(pub.KeySleep)
}
