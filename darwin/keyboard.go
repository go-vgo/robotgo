//go:build darwin
// +build darwin

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

package darwin

import (
	"errors"
	"strings"
	"time"
	"unicode/utf16"
)

// KeySleep is the global keyboard delay in milliseconds.
var KeySleep = 10

// letterCodes maps a-z to their macOS ANSI virtual key codes (kVK_ANSI_*).
var letterCodes = map[rune]uint16{
	'a': 0, 'b': 11, 'c': 8, 'd': 2, 'e': 14, 'f': 3, 'g': 5, 'h': 4,
	'i': 34, 'j': 38, 'k': 40, 'l': 37, 'm': 46, 'n': 45, 'o': 31, 'p': 35,
	'q': 12, 'r': 15, 's': 1, 't': 17, 'u': 32, 'v': 9, 'w': 13, 'x': 7,
	'y': 16, 'z': 6,
}

// digitCodes maps 0-9 to their macOS virtual key codes.
var digitCodes = map[rune]uint16{
	'0': 29, '1': 18, '2': 19, '3': 20, '4': 21,
	'5': 23, '6': 22, '7': 26, '8': 28, '9': 25,
}

// namedCodes maps robotgo named keys to macOS virtual key codes.
var namedCodes = map[string]uint16{
	"enter": 36, "return": 36, "tab": 48, "space": 49,
	"backspace": 51, "delete": 117, "del": 117, "forwarddelete": 117,
	"esc": 53, "escape": 53,
	"up": 126, "down": 125, "left": 123, "right": 124,
	"home": 115, "end": 119, "pageup": 116, "pagedown": 121,
	"capslock": 57, "help": 114, "insert": 114,

	"cmd": 55, "command": 55, "lcmd": 55, "cmdl": 55,
	"rcmd": 54, "cmdr": 54,
	"shift": 56, "shiftl": 56, "lshift": 56, "shiftr": 60, "rshift": 60,
	"ctrl": 59, "control": 59, "ctrll": 59, "ctrlr": 62,
	"alt": 58, "option": 58, "altl": 58, "altr": 61,

	"f1": 122, "f2": 120, "f3": 99, "f4": 118, "f5": 96, "f6": 97,
	"f7": 98, "f8": 100, "f9": 101, "f10": 109, "f11": 103, "f12": 111,
	"f13": 105, "f14": 107, "f15": 113, "f16": 106, "f17": 64, "f18": 79,
	"f19": 80, "f20": 90,
}

// modifierFlags maps a modifier name to its CGEventFlags mask.
var modifierFlags = map[string]uint64{
	"shift": kCGEventFlagMaskShift, "shiftl": kCGEventFlagMaskShift,
	"lshift": kCGEventFlagMaskShift, "shiftr": kCGEventFlagMaskShift,
	"rshift": kCGEventFlagMaskShift,
	"ctrl":   kCGEventFlagMaskControl, "control": kCGEventFlagMaskControl,
	"ctrll": kCGEventFlagMaskControl, "ctrlr": kCGEventFlagMaskControl,
	"alt": kCGEventFlagMaskAlternate, "option": kCGEventFlagMaskAlternate,
	"altl": kCGEventFlagMaskAlternate, "altr": kCGEventFlagMaskAlternate,
	"cmd": kCGEventFlagMaskCommand, "command": kCGEventFlagMaskCommand,
	"cmdl": kCGEventFlagMaskCommand, "lcmd": kCGEventFlagMaskCommand,
	"cmdr": kCGEventFlagMaskCommand, "rcmd": kCGEventFlagMaskCommand,
	"win": kCGEventFlagMaskCommand,
}

// keyToCode resolves a robotgo key name to a macOS virtual key code plus a
// flag mask implied by the key itself (an uppercase letter needs SHIFT). The
// final bool reports whether the key was resolved.
func keyToCode(key string) (code uint16, flags uint64, ok bool) {
	if v, found := namedCodes[strings.ToLower(key)]; found {
		return v, 0, true
	}
	r := []rune(key)
	if len(r) == 1 {
		c := r[0]
		if v, found := letterCodes[c]; found {
			return v, 0, true
		}
		if v, found := digitCodes[c]; found {
			return v, 0, true
		}
		// Uppercase letter -> base key + SHIFT.
		lower := []rune(strings.ToLower(string(c)))
		if len(lower) == 1 {
			if v, found := letterCodes[lower[0]]; found {
				return v, kCGEventFlagMaskShift, true
			}
		}
	}
	return 0, 0, false
}

// modKeyCodes resolves modifier names to their virtual key codes, skipping
// any that cannot be mapped.
func modKeyCodes(mods []string) []uint16 {
	out := make([]uint16, 0, len(mods))
	for _, m := range mods {
		if code, _, ok := keyToCode(m); ok {
			out = append(out, code)
		}
	}
	return out
}

// upModKeys posts a key-up for every modifier keycode in reverse order,
// mirroring the C backend's upKeyArr(): after a tap (or an "up" toggle) the
// modifier keys are explicitly keyed up so none can be left stuck down.
func upModKeys(mods []string, pid int) {
	codes := modKeyCodes(mods)
	for i := len(codes) - 1; i >= 0; i-- {
		sendKeyCode(codes[i], false, 0, pid)
	}
}

// flagsFromMods folds a slice of modifier names into a combined flag mask.
func flagsFromMods(mods []string) uint64 {
	var flags uint64
	for _, m := range mods {
		if f, ok := modifierFlags[strings.ToLower(m)]; ok {
			flags |= f
		}
	}
	return flags
}

// sendKeyCode posts a single keyboard event (down or up) with the given flags
// to the target pid (0 posts to the global HID event tap).
func sendKeyCode(code uint16, down bool, flags uint64, pid int) {
	if !loaded {
		return
	}
	ev := cgEventCreateKeyboardEvent(0, code, down)
	if ev == 0 {
		return
	}
	if flags != 0 {
		cgEventSetFlags(ev, flags)
	}
	postEventTo(ev, pid)
}

// sendUnicode posts a key event carrying a single rune as a Unicode string to
// the target pid (0 posts to the global HID event tap), supporting any
// character regardless of the active keyboard layout.
func sendUnicode(r rune, down bool, pid int) {
	if !loaded {
		return
	}
	ev := cgEventCreateKeyboardEvent(0, 0, down)
	if ev == 0 {
		return
	}
	units := utf16.Encode([]rune{r})
	if len(units) > 0 {
		cgEventKeyboardSetUnicode(ev, uint64(len(units)), &units[0])
	}
	postEventTo(ev, pid)
}

// KeyTap taps a key (press + release). Trailing arguments may be modifier
// names (or a []string of them) and an int pid: when a pid is supplied the
// event is delivered to that process via CGEventPostToPid, mirroring the C
// SendTo() helper; otherwise it is posted to the global HID event tap.
//
//	KeyTap("a")
//	KeyTap("a", "cmd")
//	KeyTap("a", "cmd", "shift")
//	KeyTap("a", pid)
//	KeyTap("a", pid, "cmd")
func KeyTap(key string, args ...interface{}) error {
	// fmt.Println("darwin----------------")
	pid := extractPid(args)
	mods := extractModifiers(args)
	flags := flagsFromMods(mods)

	code, autoFlags, ok := keyToCode(key)
	if !ok {
		return errors.New("robotgo: unknown key: " + key)
	}
	flags |= autoFlags

	sendKeyCode(code, true, flags, pid)
	time.Sleep(time.Duration(KeySleep) * time.Millisecond)
	sendKeyCode(code, false, flags, pid)
	// upKeyArr equivalent: explicitly key-up each modifier keycode so none
	// is left stuck down (flags alone never generate modifier key events).
	upModKeys(mods, pid)
	return nil
}

// KeyToggle toggles a key. Default is "down"; pass "up" to release. Trailing
// arguments may be modifier names (or a []string of them) and an int pid:
// when a pid is supplied the event is delivered to that process via
// CGEventPostToPid, mirroring the C SendTo() helper.
//
//	KeyToggle("a")
//	KeyToggle("a", "up")
//	KeyToggle("a", "up", "cmd")
//	KeyToggle("a", pid)
func KeyToggle(key string, args ...interface{}) error {
	up := false
	pid := extractPid(args)
	var mods []string
	toggle := func(s string) bool {
		switch s {
		case "up":
			up = true
			return true
		case "down":
			up = false
			return true
		}
		return false
	}
	for _, arg := range args {
		switch v := arg.(type) {
		case string:
			if !toggle(v) {
				if _, isMod := modifierFlags[strings.ToLower(v)]; isMod {
					mods = append(mods, v)
				}
			}
		case []string:
			for _, s := range v {
				if !toggle(s) {
					if _, isMod := modifierFlags[strings.ToLower(s)]; isMod {
						mods = append(mods, s)
					}
				}
			}
		}
	}

	code, autoFlags, ok := keyToCode(key)
	if !ok {
		return errors.New("robotgo: unknown key: " + key)
	}
	sendKeyCode(code, !up, flagsFromMods(mods)|autoFlags, pid)
	if up {
		// Mirror the C backend's keyTogglesB: releasing a key also keys up
		// its modifiers (upKeyArr) so none is left stuck down.
		upModKeys(mods, pid)
	}
	return nil
}

// KeyDown presses a key down.
func KeyDown(key string, args ...interface{}) error {
	return KeyToggle(key, append([]interface{}{"down"}, args...)...)
}

// KeyUp releases a key.
func KeyUp(key string, args ...interface{}) error {
	return KeyToggle(key, append([]interface{}{"up"}, args...)...)
}

// KeyPress presses and releases a key (alias of KeyTap).
func KeyPress(key string, args ...interface{}) error {
	return KeyTap(key, args...)
}

// Type types a string using Unicode key events, supporting any character
// (including those not on the active keyboard layout). An optional first int
// argument is the target pid: when non-zero the events are delivered to that
// process via CGEventPostToPid, mirroring the C SendTo() helper.
//
//	Type("hello")
//	Type("hello", pid)
func Type(str string, args ...int) {
	pid := 0
	if len(args) > 0 {
		pid = args[0]
	}
	for _, r := range str {
		sendUnicode(r, true, pid)
		sendUnicode(r, false, pid)
		if KeySleep > 0 {
			time.Sleep(time.Duration(KeySleep) * time.Millisecond)
		}
	}
}

// TypeStr types a string. Alias of Type, mirroring the robotgo API.
func TypeStr(str string, args ...int) {
	Type(str, args...)
}

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

// CmdCtrl returns "cmd" on macOS.
func CmdCtrl() string {
	return "cmd"
}

// extractModifiers picks the modifier names out of a variadic argument list,
// expanding any []string entries. Non-modifier strings, ints and other types
// are ignored.
func extractModifiers(args []interface{}) []string {
	var mods []string
	add := func(s string) {
		if _, isMod := modifierFlags[strings.ToLower(s)]; isMod {
			mods = append(mods, s)
		}
	}
	for _, arg := range args {
		switch v := arg.(type) {
		case string:
			add(v)
		case []string:
			for _, s := range v {
				add(s)
			}
		}
	}
	return mods
}

// extractPid returns the first int argument as the target pid for
// CGEventPostToPid (0 means post to the global HID event tap), matching the
// default robotgo backend's pid handling.
func extractPid(args []interface{}) int {
	for _, arg := range args {
		if v, ok := arg.(int); ok {
			return v
		}
	}
	return 0
}
