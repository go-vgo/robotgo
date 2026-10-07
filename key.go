//go:build !wayland && !win && !libei && !mac && !x11 && !purego
// +build !wayland,!win,!libei,!mac,!x11,!purego

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

package robotgo

/*
// #include "key/keycode.h"
#include "key/keypress_c.h"
*/
import "C"

import (
	"errors"
	"fmt"
	"runtime"
	"unsafe"
)

// keyNames define a map of key names to MMKeyCode
var keyNames = map[string]C.MMKeyCode{
	"backspace": C.K_BACKSPACE,
	"delete":    C.K_DELETE,
	"enter":     C.K_RETURN,
	"tab":       C.K_TAB,
	"esc":       C.K_ESCAPE,
	"escape":    C.K_ESCAPE,
	"up":        C.K_UP,
	"down":      C.K_DOWN,
	"right":     C.K_RIGHT,
	"left":      C.K_LEFT,
	"home":      C.K_HOME,
	"end":       C.K_END,
	"pageup":    C.K_PAGEUP,
	"pagedown":  C.K_PAGEDOWN,
	//
	"fn":  C.K_Fn,
	"f1":  C.K_F1,
	"f2":  C.K_F2,
	"f3":  C.K_F3,
	"f4":  C.K_F4,
	"f5":  C.K_F5,
	"f6":  C.K_F6,
	"f7":  C.K_F7,
	"f8":  C.K_F8,
	"f9":  C.K_F9,
	"f10": C.K_F10,
	"f11": C.K_F11,
	"f12": C.K_F12,
	"f13": C.K_F13,
	"f14": C.K_F14,
	"f15": C.K_F15,
	"f16": C.K_F16,
	"f17": C.K_F17,
	"f18": C.K_F18,
	"f19": C.K_F19,
	"f20": C.K_F20,
	"f21": C.K_F21,
	"f22": C.K_F22,
	"f23": C.K_F23,
	"f24": C.K_F24,
	//
	"cmd":         C.K_META,
	"cmdl":        C.K_LMETA,
	"cmdr":        C.K_RMETA,
	"command":     C.K_META,
	"alt":         C.K_ALT,
	"altl":        C.K_LALT,
	"altr":        C.K_RALT,
	"ctrl":        C.K_CONTROL,
	"ctrll":       C.K_LCONTROL,
	"ctrlr":       C.K_RCONTROL,
	"control":     C.K_CONTROL,
	"shift":       C.K_SHIFT,
	"shiftl":      C.K_LSHIFT,
	"shiftr":      C.K_RSHIFT,
	"right_shift": C.K_RSHIFT,
	"caps":        C.K_CAPSLOCK,
	"capslock":    C.K_CAPSLOCK,
	"space":       C.K_SPACE,
	"print":       C.K_PRINTSCREEN,
	"printscreen": C.K_PRINTSCREEN,
	"insert":      C.K_INSERT,
	"menu":        C.K_MENU,
	"scroll_lock": C.K_SCROLL_LOCK,
	"pause_break": C.K_PAUSE,

	"audio_mute":     C.K_AUDIO_VOLUME_MUTE,
	"audio_vol_down": C.K_AUDIO_VOLUME_DOWN,
	"audio_vol_up":   C.K_AUDIO_VOLUME_UP,
	"audio_play":     C.K_AUDIO_PLAY,
	"audio_stop":     C.K_AUDIO_STOP,
	"audio_pause":    C.K_AUDIO_PAUSE,
	"audio_prev":     C.K_AUDIO_PREV,
	"audio_next":     C.K_AUDIO_NEXT,
	"audio_rewind":   C.K_AUDIO_REWIND,
	"audio_forward":  C.K_AUDIO_FORWARD,
	"audio_repeat":   C.K_AUDIO_REPEAT,
	"audio_random":   C.K_AUDIO_RANDOM,

	"num0":     C.K_NUMPAD_0,
	"num1":     C.K_NUMPAD_1,
	"num2":     C.K_NUMPAD_2,
	"num3":     C.K_NUMPAD_3,
	"num4":     C.K_NUMPAD_4,
	"num5":     C.K_NUMPAD_5,
	"num6":     C.K_NUMPAD_6,
	"num7":     C.K_NUMPAD_7,
	"num8":     C.K_NUMPAD_8,
	"num9":     C.K_NUMPAD_9,
	"num_lock": C.K_NUMPAD_LOCK,

	// todo: removed
	"numpad_0":    C.K_NUMPAD_0,
	"numpad_1":    C.K_NUMPAD_1,
	"numpad_2":    C.K_NUMPAD_2,
	"numpad_3":    C.K_NUMPAD_3,
	"numpad_4":    C.K_NUMPAD_4,
	"numpad_5":    C.K_NUMPAD_5,
	"numpad_6":    C.K_NUMPAD_6,
	"numpad_7":    C.K_NUMPAD_7,
	"numpad_8":    C.K_NUMPAD_8,
	"numpad_9":    C.K_NUMPAD_9,
	"numpad_lock": C.K_NUMPAD_LOCK,

	"num.":      C.K_NUMPAD_DECIMAL,
	"num+":      C.K_NUMPAD_PLUS,
	"num-":      C.K_NUMPAD_MINUS,
	"num*":      C.K_NUMPAD_MUL,
	"num/":      C.K_NUMPAD_DIV,
	"num_clear": C.K_NUMPAD_CLEAR,
	"num_enter": C.K_NUMPAD_ENTER,
	"num_equal": C.K_NUMPAD_EQUAL,

	"lights_mon_up":     C.K_LIGHTS_MON_UP,
	"lights_mon_down":   C.K_LIGHTS_MON_DOWN,
	"lights_kbd_toggle": C.K_LIGHTS_KBD_TOGGLE,
	"lights_kbd_up":     C.K_LIGHTS_KBD_UP,
	"lights_kbd_down":   C.K_LIGHTS_KBD_DOWN,

	// { NULL:              C.K_NOT_A_KEY }
}

// macCharToKeyCode maps ASCII characters to macOS virtual key codes (kVK_ANSI_*)
// This avoids calling C.keyCodeForChar which can cause SIGTRAP on macOS
var macCharToKeyCode = map[byte]C.MMKeyCode{
	'a': 0x00, 'A': 0x00,
	's': 0x01, 'S': 0x01,
	'd': 0x02, 'D': 0x02,
	'f': 0x03, 'F': 0x03,
	'h': 0x04, 'H': 0x04,
	'g': 0x05, 'G': 0x05,
	'z': 0x06, 'Z': 0x06,
	'x': 0x07, 'X': 0x07,
	'c': 0x08, 'C': 0x08,
	'v': 0x09, 'V': 0x09,
	'b': 0x0B, 'B': 0x0B,
	'q': 0x0C, 'Q': 0x0C,
	'w': 0x0D, 'W': 0x0D,
	'e': 0x0E, 'E': 0x0E,
	'r': 0x0F, 'R': 0x0F,
	'y': 0x10, 'Y': 0x10,
	't': 0x11, 'T': 0x11,
	'1': 0x12, '!': 0x12,
	'2': 0x13, '@': 0x13,
	'3': 0x14, '#': 0x14,
	'4': 0x15, '$': 0x15,
	'6': 0x16, '^': 0x16,
	'5': 0x17, '%': 0x17,
	'=': 0x18, '+': 0x18,
	'9': 0x19, '(': 0x19,
	'7': 0x1A, '&': 0x1A,
	'-': 0x1B, '_': 0x1B,
	'8': 0x1C, '*': 0x1C,
	'0': 0x1D, ')': 0x1D,
	']': 0x1E, '}': 0x1E,
	'o': 0x1F, 'O': 0x1F,
	'u': 0x20, 'U': 0x20,
	'[': 0x21, '{': 0x21,
	'i': 0x22, 'I': 0x22,
	'p': 0x23, 'P': 0x23,
	'l': 0x25, 'L': 0x25,
	'j': 0x26, 'J': 0x26,
	'\'': 0x27, '"': 0x27,
	'k': 0x28, 'K': 0x28,
	';': 0x29, ':': 0x29,
	'\\': 0x2A, '|': 0x2A,
	',': 0x2B, '<': 0x2B,
	'/': 0x2C, '?': 0x2C,
	'n': 0x2D, 'N': 0x2D,
	'm': 0x2E, 'M': 0x2E,
	'.': 0x2F, '>': 0x2F,
	'`': 0x32, '~': 0x32,
	' ': 0x31, // kVK_Space
}

// It sends a key press and release to the active application
// The key up is always sent so modifiers pressed by a failed key down are released.
func tapKeyCode(code C.MMKeyCode, flags C.MMKeyFlags, pid C.uintptr) (int, string) {
	c1 := C.toggleKeyCode(code, true, flags, pid)
	MilliSleep(3)
	c2 := C.toggleKeyCode(code, false, flags, pid)
	if c1 != 0 {
		return int(c1), "down"
	}
	return int(c2), "up"
}

var keyErr = errors.New("Invalid key flag specified.")

func checkKeyCodes(k string) (key C.MMKeyCode, err error) {
	if k == "" {
		return
	}

	if len(k) == 1 {
		// On macOS, use Go lookup table to avoid SIGTRAP in CGO
		if runtime.GOOS == "darwin" {
			c := k[0]
			if code, ok := macCharToKeyCode[c]; ok {
				key = code
				return
			}
			err = keyErr
			return
		}

		val1 := C.CString(k)
		defer C.free(unsafe.Pointer(val1))

		key = C.keyCodeForChar(*val1)
		if key == C.K_NOT_A_KEY {
			err = keyErr
			return
		}
		return
	}

	if v, ok := keyNames[k]; ok {
		key = v
		if key == C.K_NOT_A_KEY {
			err = keyErr
		}
		return
	}
	// Unknown names used to fall through as keycode 0 ("a" on macOS); report
	// them like the pure-Go backends do.
	err = keyErr
	return
}

// keyFlags maps modifier names to MMKeyFlags.
var keyFlags = map[string]C.MMKeyFlags{
	"alt":     C.MOD_ALT,
	"altr":    C.MOD_ALT,
	"altl":    C.MOD_ALT,
	"cmd":     C.MOD_META,
	"command": C.MOD_META,
	"cmdr":    C.MOD_META,
	"cmdl":    C.MOD_META,
	"ctrl":    C.MOD_CONTROL,
	"control": C.MOD_CONTROL,
	"ctrlr":   C.MOD_CONTROL,
	"ctrll":   C.MOD_CONTROL,
	"shift":   C.MOD_SHIFT,
	"shiftr":  C.MOD_SHIFT,
	"shiftl":  C.MOD_SHIFT,
	// legacy key name (keyNames) accepted as a modifier too
	"right_shift": C.MOD_SHIFT,
	"none":        C.MOD_NONE,
}

func checkKeyFlags(f string) C.MMKeyFlags {
	return keyFlags[f]
}

func getFlagsFromValue(value []string) (flags C.MMKeyFlags) {
	for _, v := range value {
		flags |= checkKeyFlags(v)
	}
	return
}

// upKeyArr releases every modifier in keyArr. On macOS each release carries
// the mask of modifiers still held (see toggleKeyCode) so back-to-back
// releases cannot re-assert a modifier from stale HID state.
func upKeyArr(keyArr []string, pid int) {
	var remaining C.MMKeyFlags
	if runtime.GOOS == "darwin" {
		remaining = getFlagsFromValue(keyArr)
	}
	for _, k := range keyArr {
		key1, err := checkKeyCodes(k)
		if err != nil {
			continue
		}
		remaining &^= checkKeyFlags(k)
		C.toggleKeyCode(key1, false, remaining, C.uintptr(pid))
	}
}

func keyTaps(k string, keyArr []string, pid int) error {
	flags := getFlagsFromValue(keyArr)
	key, err := checkKeyCodes(k)
	if err != nil {
		return err
	}

	c1, stage := tapKeyCode(key, flags, C.uintptr(pid))
	MilliSleep(KeySleep)
	upKeyArr(keyArr, pid)
	return formatClickError(c1, k, stage, 1)
}

func keyTogglesB(k string, down bool, keyArr []string, pid int) error {
	flags := getFlagsFromValue(keyArr)
	key, err := checkKeyCodes(k)
	if err != nil {
		return err
	}

	c1 := C.toggleKeyCode(key, C.bool(down), flags, C.uintptr(pid))
	MilliSleep(KeySleep)
	if !down {
		upKeyArr(keyArr, pid)
	}
	return formatClickError(int(c1), k, getDown(down), 1)
}

func keyToggles(k string, keyArr []string, pid int) error {
	down, keyArr1 := getKeyDown(keyArr)
	return keyTogglesB(k, down, keyArr1, pid)
}

/*
 __  ___  ___________    ____ .______     ______        ___      .______       _______
|  |/  / |   ____\   \  /   / |   _  \   /  __  \      /   \     |   _  \     |       \
|  '  /  |  |__   \   \/   /  |  |_)  | |  |  |  |    /  ^  \    |  |_)  |    |  .--.  |
|    <   |   __|   \_    _/   |   _  <  |  |  |  |   /  /_\  \   |      /     |  |  |  |
|  .  \  |  |____    |  |     |  |_)  | |  `--'  |  /  _____  \  |  |\  \----.|  '--'  |
|__|\__\ |_______|   |__|     |______/   \______/  /__/     \__\ | _| `._____||_______/

*/

// toErr it converts a C string to a Go error
func toErr(str *C.char) error {
	gstr := C.GoString(str)
	if gstr == "" {
		return nil
	}
	return errors.New(gstr)
}

// KeyTap taps the keyboard code;
//
// See keys supported:
//
//	https://github.com/go-vgo/robotgo/blob/master/docs/keys.md#keys
//
// Examples:
//
//	robotgo.KeySleep = 100 // 100 millisecond
//	robotgo.KeyTap("a")
//	robotgo.KeyTap("i", "alt", "command")
//
//	arr := []string{"alt", "command"}
//	robotgo.KeyTap("i", arr)
//
//	robotgo.KeyTap("k", pid int)
func KeyTap(key string, args ...interface{}) error {
	key, args = appendShift(key, 0, args...)
	pid, keyArr := getToggleArgs(args...)
	return keyTaps(key, keyArr, pid)
}

// KeyToggle toggles the keyboard, if there not have args default is "down"
//
// See keys:
//
//	https://github.com/go-vgo/robotgo/blob/master/docs/keys.md#keys
//
// Examples:
//
//	robotgo.KeyToggle("a")
//	robotgo.KeyToggle("a", "up")
//
//	robotgo.KeyToggle("a", "up", "alt", "cmd")
//	robotgo.KeyToggle("a", "up", []string{"alt", "cmd"})
//	robotgo.KeyToggle("k", pid int)
func KeyToggle(key string, args ...interface{}) error {
	key, args = appendShift(key, 1, args...)
	pid, keyArr := getToggleArgs(args...)
	return keyToggles(key, keyArr, pid)
}

// UnicodeType tap the uint32 unicode
func UnicodeType(str uint32, args ...int) error {
	cstr := C.uint(str)
	pid := 0
	if len(args) > 0 {
		pid = args[0]
	}

	isPid := 0
	if len(args) > 1 {
		isPid = args[1]
	}

	code := C.unicodeType(cstr, C.uintptr(pid), C.int8_t(isPid))
	return formatKeyError(int(code), str)
}

// inputUTFDetail describes the non-zero codes of C input_utf (X11 only).
var inputUTFDetail = map[int]string{
	1: "no X display",
	2: "unknown keysym",
	3: "keyboard mapping unavailable",
	4: "XTestFakeKeyEvent returned false",
}

func inputUTF(str string) error {
	cstr := C.CString(str)
	code := C.input_utf(cstr)

	C.free(unsafe.Pointer(cstr))
	if code != 0 {
		return fmt.Errorf("input %q failed: %s (code=%d)", str, inputUTFDetail[int(code)], int(code))
	}
	return nil
}

// formatKeyError converts a non-zero C unicode type code to an error
func formatKeyError(code int, r uint32) error {
	if code == 0 {
		return nil
	}
	if detail := codeDetail(code); detail != "" {
		return fmt.Errorf("type %q failed: %s (code=%d)", rune(r), detail, code)
	}
	return fmt.Errorf("type %q failed, code=%d", rune(r), code)
}

// TypeStr tap a string and return error
func TypeStr(str string, args ...int) error {
	_, err := typeStr(str, args...)
	return err
}

// Type type a string (supported UTF-8)
//
// robotgo.Type(string: "The string to send", int: pid, "milli_sleep time", "x11 option")
//
// Examples:
//
//	robotgo.Type("abc@123, Hi galaxy, こんにちは")
//	robotgo.Type("To be or not to be, this is questions.", pid int)
//
// It returns the count of characters typed, stops at the first failure.
func Type(str string, args ...int) int {
	n, _ := typeStr(str, args...)
	return n
}

// typeStr types str and returns the count of characters typed
// and the first error.
func typeStr(str string, args ...int) (int, error) {
	var tm, tm1 = 0, 7

	if len(args) > 1 {
		tm = args[1]
	}
	if len(args) > 2 {
		tm1 = args[2]
	}
	pid := 0
	if len(args) > 0 {
		pid = args[0]
	}
	// Windows: NotPid makes pid an HWND, like the other key APIs.
	isPid := 0
	if NotPid {
		isPid = 1
	}

	if runtime.GOOS == "linux" {
		strUc := ToUC(str)
		for i := 0; i < len(strUc); i++ {
			ru := []rune(strUc[i])
			var err error
			if len(ru) <= 1 {
				ustr := uint32(CharCodeAt(strUc[i], 0))
				err = UnicodeType(ustr, pid, isPid)
			} else {
				err = inputUTF(strUc[i])
				MilliSleep(tm1)
			}
			if err != nil {
				return i, err
			}

			MilliSleep(tm)
		}
		return len(strUc), nil
	}

	l1 := len([]rune(str))
	for i := 0; i < l1; i++ {
		ustr := uint32(CharCodeAt(str, i))
		if err := UnicodeType(ustr, pid, isPid); err != nil {
			return i, err
		}
		// if len(args) > 0 {
		MilliSleep(tm)
		// }
	}
	MilliSleep(KeySleep)
	return l1, nil
}
