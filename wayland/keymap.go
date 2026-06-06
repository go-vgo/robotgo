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

// keyToEvdev maps robotgo key names to Linux evdev key codes.
// Evdev codes are used directly by zwp_virtual_keyboard_v1.
// The virtual keyboard uses XKB keycodes = evdev + 8, but
// the protocol expects raw evdev codes.
func keyToEvdev(key string) (uint32, bool) {
	code, ok := evdevKeyMap[key]
	return code, ok
}

// evdevKeyMap maps robotgo key strings to Linux evdev keycodes.
// Reference: linux/input-event-codes.h
var evdevKeyMap = map[string]uint32{
	// Row 0: Escape
	"esc": 1, "escape": 1,

	// Row 1: Number row
	"1": 2, "2": 3, "3": 4, "4": 5, "5": 6,
	"6": 7, "7": 8, "8": 9, "9": 10, "0": 11,
	"-": 12, "=": 13, "backspace": 14,

	// Row 2: QWERTY
	"tab": 15,
	"q":   16, "w": 17, "e": 18, "r": 19, "t": 20,
	"y": 21, "u": 22, "i": 23, "o": 24, "p": 25,
	"[": 26, "]": 27, "enter": 28, "return": 28,

	// Row 3: ASDF
	"a": 30, "s": 31, "d": 32, "f": 33, "g": 34,
	"h": 35, "j": 36, "k": 37, "l": 38,
	";": 39, "'": 40, "`": 41,

	// Row 4: ZXCV
	"\\": 43,
	"z":  44, "x": 45, "c": 46, "v": 47, "b": 48,
	"n": 49, "m": 50,
	",": 51, ".": 52, "/": 53,

	// Modifiers
	"shiftl": 42, "shiftr": 54, "shift": 42,
	"ctrll": 29, "ctrlr": 97, "ctrl": 29, "control": 29,
	"altl": 56, "altr": 100, "alt": 56,
	"cmd": 125, "cmdl": 125, "cmdr": 126, // KEY_LEFTMETA / KEY_RIGHTMETA
	"space":    57,
	"capslock": 58,

	// Function keys
	"f1": 59, "f2": 60, "f3": 61, "f4": 62, "f5": 63, "f6": 64,
	"f7": 65, "f8": 66, "f9": 67, "f10": 68, "f11": 87, "f12": 88,
	"f13": 183, "f14": 184, "f15": 185, "f16": 186, "f17": 187, "f18": 188,
	"f19": 189, "f20": 190, "f21": 191, "f22": 192, "f23": 193, "f24": 194,

	// Navigation
	"home": 102, "up": 103, "pageup": 104, "pgup": 104,
	"left": 105, "right": 106,
	"end": 107, "down": 108, "pagedown": 109, "pgdn": 109,
	"insert": 110, "delete": 111,

	// Numpad
	"num_lock": 69,
	"num/":     98, "num*": 55, "num-": 74, "num+": 78, "num_enter": 96,
	"num.": 83, "num_clear": 69,
	"num0": 82, "num1": 79, "num2": 80, "num3": 81,
	"num4": 75, "num5": 76, "num6": 77,
	"num7": 71, "num8": 72, "num9": 73,

	// Special
	"print": 99, "printscreen": 99,
	"scroll_lock": 70,
	"pause":       119,
	"menu":        127,

	// Media keys
	"audio_mute": 113, "audio_vol_down": 114, "audio_vol_up": 115,
	"audio_play": 164, "audio_stop": 166, "audio_pause": 164,
	"audio_prev": 165, "audio_next": 163,
}
