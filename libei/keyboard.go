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
	"time"
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

// KeyTap taps a key (press + release), optionally with modifiers.
//
//	KeyTap("a")
//	KeyTap("a", "ctrl")
//	KeyTap("a", "ctrl", "shift")
func KeyTap(key string, args ...interface{}) error {
	c, err := keyboardReady()
	if err != nil {
		return err
	}

	code, ok := keyToEvdev(key)
	if !ok {
		return errors.New("robotgo/libei: unknown key: " + key)
	}

	// Press modifiers, remembering the ones that actually went down so they
	// are always released (the upKeyArr behavior of the C backend), even when
	// a later injection fails.
	var pressed []int32
	upMods := func() error {
		return releaseKeys(pressed, func(mc int32) error {
			return c.inj.keyboardKeycode(mc, stateReleased)
		})
	}
	for _, mod := range extractModifiers(args) {
		mc, ok := keyToEvdev(mod)
		if !ok {
			continue
		}
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

// KeyToggle toggles a key down or up. Default is "down".
//
//	KeyToggle("a")        // press
//	KeyToggle("a", "up")  // release
func KeyToggle(key string, args ...interface{}) error {
	c, err := keyboardReady()
	if err != nil {
		return err
	}

	state := statePressed
	for _, arg := range args {
		if s, ok := arg.(string); ok && s == "up" {
			state = stateReleased
		}
	}

	code, ok := keyToEvdev(key)
	if !ok {
		return errors.New("robotgo/libei: unknown key: " + key)
	}
	return c.inj.keyboardKeycode(code, state)
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
// bookkeeping.
func Type(str string, args ...int) {
	c, err := keyboardReady()
	if err != nil {
		return
	}
	for _, r := range str {
		sym := runeToKeysym(r)
		if err := c.inj.keyboardKeysym(sym, statePressed); err != nil {
			return
		}
		time.Sleep(time.Duration(KeySleep) * time.Millisecond)
		if err := c.inj.keyboardKeysym(sym, stateReleased); err != nil {
			return
		}
	}
}

// TypeStr types a string (alias of Type).
func TypeStr(str string, args ...int) { Type(str, args...) }

// TypeDelay types a string with a per-character delay in milliseconds.
func TypeDelay(str string, delay int) {
	old := KeySleep
	KeySleep = delay
	Type(str)
	KeySleep = old
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

// extractModifiers pulls modifier key names out of variadic args.
func extractModifiers(args []interface{}) []string {
	var mods []string
	for _, arg := range args {
		if s, ok := arg.(string); ok {
			switch s {
			case "ctrl", "control", "ctrll", "ctrlr":
				mods = append(mods, s)
			case "shift", "shiftl", "shiftr":
				mods = append(mods, s)
			case "alt", "altl", "altr":
				mods = append(mods, s)
			case "cmd", "cmdl", "cmdr":
				mods = append(mods, s)
			}
		}
	}
	return mods
}
