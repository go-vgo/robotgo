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
	"errors"
	"fmt"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"
)

// KeySleep is the global keyboard delay in milliseconds (between press and
// release in KeyTap, and between characters in Type).
var KeySleep = 10

// Key name constants matching robotgo's API.
const (
	KeyA = "a"
	KeyB = "b"
	KeyC = "c"
	KeyD = "d"
	KeyE = "e"
	KeyF = "f"
	KeyG = "g"
	KeyH = "h"
	KeyI = "i"
	KeyJ = "j"
	KeyK = "k"
	KeyL = "l"
	KeyM = "m"
	KeyN = "n"
	KeyO = "o"
	KeyP = "p"
	KeyQ = "q"
	KeyR = "r"
	KeyS = "s"
	KeyT = "t"
	KeyU = "u"
	KeyV = "v"
	KeyW = "w"
	KeyX = "x"
	KeyY = "y"
	KeyZ = "z"

	Backspace = "backspace"
	Delete    = "delete"
	Enter     = "enter"
	Tab       = "tab"
	Esc       = "esc"
	Escape    = "escape"
	Up        = "up"
	Down      = "down"
	Right     = "right"
	Left      = "left"
	Home      = "home"
	End       = "end"
	Pageup    = "pageup"
	Pagedown  = "pagedown"

	F1  = "f1"
	F2  = "f2"
	F3  = "f3"
	F4  = "f4"
	F5  = "f5"
	F6  = "f6"
	F7  = "f7"
	F8  = "f8"
	F9  = "f9"
	F10 = "f10"
	F11 = "f11"
	F12 = "f12"

	Shift    = "shift"
	Ctrl     = "ctrl"
	Alt      = "alt"
	Cmd      = "cmd"
	ShiftL   = "shiftl"
	ShiftR   = "shiftr"
	CtrlL    = "ctrll"
	CtrlR    = "ctrlr"
	AltL     = "altl"
	AltR     = "altr"
	Space    = "space"
	Capslock = "capslock"
	Caps     = "caps"
	Print    = "print"
	Insert   = "insert"
	Menu     = "menu"
)

// keyboardReady returns the connection if keyboard injection is available.
func keyboardReady() (*conn, error) {
	c, err := ensureConn()
	if err != nil {
		return nil, err
	}
	if c.inj == nil || !c.hasKeyboard() {
		return nil, ErrNotSupported
	}
	return c, nil
}

// resolveKey maps a robotgo key name to its evdev code, treating a single
// uppercase letter as the lowercase key plus shift (like the Cgo backend).
func resolveKey(key string) (code int32, shift bool, ok bool) {
	if code, ok = keyToEvdev(key); ok {
		return code, false, true
	}
	if r := []rune(key); len(r) == 1 && unicode.IsUpper(r[0]) {
		code, ok = keyToEvdev(string(unicode.ToLower(r[0])))
		return code, true, ok
	}
	return 0, false, false
}

// KeyTap taps a key (press + release), optionally with modifiers.
//
//	KeyTap("a")
//	KeyTap("a", "ctrl")
//	KeyTap("a", "ctrl", "shift")
//	KeyTap("a", []string{"ctrl", "shift"})
func KeyTap(key string, args ...interface{}) error {
	c, err := keyboardReady()
	if err != nil {
		return err
	}

	codes, _, err := toggleKeys(key, args)
	if err != nil {
		return err
	}
	code := codes[len(codes)-1]

	// Press modifiers, remembering the ones that actually went down so they
	// are always released (the upKeyArr behavior of the C backend), even when
	// a later injection fails.
	var pressed []int32
	upMods := func() error {
		return releaseKeys(pressed, func(mc int32) error {
			return c.inj.keyboardKeycode(mc, stateReleased)
		})
	}
	for _, mc := range codes[:len(codes)-1] {
		if err := c.inj.keyboardKeycode(mc, statePressed); err != nil {
			return errors.Join(err, upMods())
		}
		pressed = append(pressed, mc)
	}

	// Press + release the key.
	if err := c.inj.keyboardKeycode(code, statePressed); err != nil {
		return errors.Join(err, upMods())
	}
	time.Sleep(time.Duration(KeySleep) * time.Millisecond)
	err = c.inj.keyboardKeycode(code, stateReleased)

	// Release modifiers in reverse order (upKeyArr) even if the key release
	// failed, so no modifier is left stuck down.
	return errors.Join(err, upMods())
}

// releaseKeys keys up the given evdev codes in reverse order via send,
// mirroring the C backend's upKeyArr(). It keeps going past failures so every
// key gets a release attempt, and returns the errors it hit (joined).
func releaseKeys(codes []int32, send func(code int32) error) error {
	var errs []error
	for i := len(codes) - 1; i >= 0; i-- {
		if err := send(codes[i]); err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}

// toggleKeys resolves a KeyTap/KeyToggle call into evdev codes in press order
// (modifiers, implied shift, then the key) and its direction; the last "up"
// or "down" argument wins. Codes are deduplicated so aliases press once, and
// the implied shift is skipped when any shift variant is already given.
func toggleKeys(key string, args []interface{}) (codes []int32, up bool, err error) {
	code, shift, ok := resolveKey(key)
	if !ok {
		return nil, false, errors.New("robotgo/libei: unknown key: " + key)
	}
	for _, s := range keyArgs(args) {
		switch s {
		case "up":
			up = true
		case "down":
			up = false
		}
	}

	mods := extractModifiers(args)
	if shift && !hasShift(mods) {
		mods = append(mods, "shift")
	}
	seen := map[int32]bool{code: true}
	for _, mod := range mods {
		if mc, ok := keyToEvdev(mod); ok && !seen[mc] {
			seen[mc] = true
			codes = append(codes, mc)
		}
	}
	return append(codes, code), up, nil
}

// KeyToggle toggles a key down or up. Default is "down". Modifiers (strings
// or a []string) are pressed before the key and released in reverse order
// after it.
//
//	KeyToggle("a")        // press
//	KeyToggle("a", "up")  // release
//	KeyToggle("a", "down", []string{"ctrl", "shift"})
//	KeyToggle("a", "up", []string{"ctrl", "shift"})
func KeyToggle(key string, args ...interface{}) error {
	c, err := keyboardReady()
	if err != nil {
		return err
	}

	codes, up, err := toggleKeys(key, args)
	if err != nil {
		return err
	}
	release := func(code int32) error {
		return c.inj.keyboardKeycode(code, stateReleased)
	}
	if up {
		return releaseKeys(codes, release)
	}
	for i, code := range codes {
		if err := c.inj.keyboardKeycode(code, statePressed); err != nil {
			return errors.Join(err, releaseKeys(codes[:i], release))
		}
	}
	return nil
}

// KeyDown presses a key down. Extra args are forwarded to KeyToggle for API
// parity with the other backends.
func KeyDown(key string, args ...interface{}) error {
	return KeyToggle(key, append([]interface{}{"down"}, args...)...)
}

// KeyUp releases a key. Extra args are forwarded to KeyToggle for API parity
// with the other backends.
func KeyUp(key string, args ...interface{}) error {
	return KeyToggle(key, append([]interface{}{"up"}, args...)...)
}

// KeyPress presses and releases a key (alias of KeyTap).
func KeyPress(key string, args ...interface{}) error { return KeyTap(key, args...) }

// Type types a string. Each rune is sent as an X11 keysym via
// NotifyKeyboardKeysym, so it is layout independent and needs no shift
// bookkeeping. It returns the number of characters (runes) typed and stops
// at the first failed key event.
func Type(str string, args ...int) int {
	c, err := keyboardReady()
	if err != nil {
		return 0
	}
	n := 0
	for _, r := range str {
		sym := runeToKeysym(r)
		if err := c.inj.keyboardKeysym(sym, statePressed); err != nil {
			return n
		}
		time.Sleep(time.Duration(KeySleep) * time.Millisecond)
		if err := c.inj.keyboardKeysym(sym, stateReleased); err != nil {
			return n
		}
		n++
	}
	return n
}

// TypeStr types a string (alias of Type). It returns an error if not every
// character was typed.
func TypeStr(str string, args ...int) error { return typeErr(Type(str, args...), str) }

// TypeDelay types a string with a per-character delay in milliseconds. It
// returns an error if not every character was typed.
func TypeDelay(str string, delay int) error {
	old := KeySleep
	KeySleep = delay
	n := Type(str)
	KeySleep = old
	return typeErr(n, str)
}

// typeErr reports an error when fewer than all runes of str were typed.
func typeErr(n int, str string) error {
	if total := utf8.RuneCountInString(str); n < total {
		return fmt.Errorf("robotgo: typed %d of %d characters", n, total)
	}
	return nil
}

// SetDelay sets both KeySleep and MouseSleep.
func SetDelay(d ...int) {
	delay := 10
	if len(d) > 0 {
		delay = d[0]
	}
	KeySleep = delay
	MouseSleep = delay
}

// CmdCtrl returns "ctrl" on Linux (mirrors robotgo's cross-platform helper).
func CmdCtrl() string { return "ctrl" }

// keyArgs flattens the string and []string arguments of KeyTap/KeyToggle,
// skipping other types (such as an int pid).
func keyArgs(args []interface{}) []string {
	var out []string
	for _, arg := range args {
		switch v := arg.(type) {
		case string:
			out = append(out, v)
		case []string:
			out = append(out, v...)
		}
	}
	return out
}

// hasShift reports whether mods (lowercased) contains any shift variant
// (shift, shiftl, shiftr, right_shift).
func hasShift(mods []string) bool {
	for _, m := range mods {
		if strings.Contains(m, "shift") {
			return true
		}
	}
	return false
}

// extractModifiers picks the (case-insensitive) modifier names out of the
// variadic args, expanding []string entries; results are lowercased.
func extractModifiers(args []interface{}) []string {
	var mods []string
	for _, s := range keyArgs(args) {
		s = strings.ToLower(s)
		switch s {
		case "ctrl", "control", "ctrll", "ctrlr",
			"shift", "shiftl", "shiftr", "right_shift",
			"alt", "altl", "altr",
			"cmd", "command", "cmdl", "cmdr":
			mods = append(mods, s)
		}
	}
	return mods
}
