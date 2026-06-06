//go:build windows
// +build windows

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

package win

import (
	"errors"
	"strings"
	"time"
	"unicode/utf16"
	"unsafe"

	"github.com/tailscale/win"
)

// KeySleep is the global keyboard delay in milliseconds.
var KeySleep = 10

// Key constants matching robotgo's API.
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

// vkMap maps robotgo named keys to Win32 virtual-key codes.
var vkMap = map[string]uint16{
	"enter": win.VK_RETURN, "return": win.VK_RETURN,
	"tab": win.VK_TAB, "space": win.VK_SPACE,
	"backspace": win.VK_BACK, "delete": win.VK_DELETE,
	"esc": win.VK_ESCAPE, "escape": win.VK_ESCAPE,
	"up": win.VK_UP, "down": win.VK_DOWN, "left": win.VK_LEFT, "right": win.VK_RIGHT,
	"home": win.VK_HOME, "end": win.VK_END,
	"pageup": win.VK_PRIOR, "pgup": win.VK_PRIOR,
	"pagedown": win.VK_NEXT, "pgdn": win.VK_NEXT,
	"insert": win.VK_INSERT,

	"shift": win.VK_SHIFT, "shiftl": win.VK_LSHIFT, "shiftr": win.VK_RSHIFT,
	"ctrl": win.VK_CONTROL, "control": win.VK_CONTROL,
	"ctrll": win.VK_LCONTROL, "ctrlr": win.VK_RCONTROL,
	"alt": win.VK_MENU, "altl": win.VK_LMENU, "altr": win.VK_RMENU,
	"cmd": win.VK_LWIN, "cmdl": win.VK_LWIN, "win": win.VK_LWIN,
	"cmdr": win.VK_RWIN, "rwin": win.VK_RWIN,
	"capslock": win.VK_CAPITAL,
	"print":    win.VK_SNAPSHOT, "printscreen": win.VK_SNAPSHOT,
	"menu":     win.VK_APPS,
	"num_lock": win.VK_NUMLOCK, "scroll_lock": win.VK_SCROLL,
	"pause": win.VK_PAUSE,

	"f1": win.VK_F1, "f2": win.VK_F2, "f3": win.VK_F3, "f4": win.VK_F4,
	"f5": win.VK_F5, "f6": win.VK_F6, "f7": win.VK_F7, "f8": win.VK_F8,
	"f9": win.VK_F9, "f10": win.VK_F10, "f11": win.VK_F11, "f12": win.VK_F12,
	"f13": win.VK_F13, "f14": win.VK_F14, "f15": win.VK_F15, "f16": win.VK_F16,
	"f17": win.VK_F17, "f18": win.VK_F18, "f19": win.VK_F19, "f20": win.VK_F20,
	"f21": win.VK_F21, "f22": win.VK_F22, "f23": win.VK_F23, "f24": win.VK_F24,
}

// extendedVKs are the virtual keys that must be sent with the
// KEYEVENTF_EXTENDEDKEY flag so apps reading the scancode/extended bit
// (games, RDP, DirectInput, left/right modifier discrimination) see them
// correctly.
var extendedVKs = map[uint16]bool{
	win.VK_RCONTROL: true, win.VK_RMENU: true,
	win.VK_INSERT: true, win.VK_DELETE: true,
	win.VK_HOME: true, win.VK_END: true,
	win.VK_PRIOR: true, win.VK_NEXT: true,
	win.VK_LEFT: true, win.VK_RIGHT: true, win.VK_UP: true, win.VK_DOWN: true,
	win.VK_NUMLOCK: true, win.VK_SNAPSHOT: true,
	win.VK_LWIN: true, win.VK_RWIN: true, win.VK_APPS: true,
}

// keyToVK resolves a robotgo key name to a Win32 virtual-key code plus the
// modifier bitmask (bit0=shift, bit1=ctrl, bit2=alt) needed to produce it.
// Named keys come from vkMap (no implied modifiers); single characters fall
// back to VkKeyScan, which covers letters, digits, and punctuation for the
// active layout and reports the required shift state in its high byte.
func keyToVK(key string) (vk uint16, mods uint8, ok bool) {
	if v, found := vkMap[strings.ToLower(key)]; found {
		return v, 0, true
	}
	r := []rune(key)
	if len(r) == 1 {
		res := win.VkKeyScan(uint16(r[0]))
		if res != -1 {
			return uint16(byte(res & 0xff)), uint8((res >> 8) & 0xff), true
		}
	}
	return 0, 0, false
}

// sendVK dispatches a single virtual-key event (down or up).
func sendVK(vk uint16, up bool) {
	flags := uint32(0)
	if extendedVKs[vk] {
		flags |= win.KEYEVENTF_EXTENDEDKEY
	}
	if up {
		flags |= win.KEYEVENTF_KEYUP
	}
	in := win.KEYBD_INPUT{
		Type: win.INPUT_KEYBOARD,
		Ki: win.KEYBDINPUT{
			WVk:     vk,
			DwFlags: flags,
		},
	}
	win.SendInput(1, unsafe.Pointer(&in), int32(unsafe.Sizeof(in)))
}

// sendUnicode dispatches a single UTF-16 code unit as a Unicode key event.
func sendUnicode(u uint16, up bool) {
	flags := uint32(win.KEYEVENTF_UNICODE)
	if up {
		flags |= win.KEYEVENTF_KEYUP
	}
	in := win.KEYBD_INPUT{
		Type: win.INPUT_KEYBOARD,
		Ki: win.KEYBDINPUT{
			WScan:   u,
			DwFlags: flags,
		},
	}
	win.SendInput(1, unsafe.Pointer(&in), int32(unsafe.Sizeof(in)))
}

// KeyTap taps a key (press + release). Optional trailing modifiers.
//
//	KeyTap("a")
//	KeyTap("a", "ctrl")
//	KeyTap("a", "ctrl", "shift")
func KeyTap(key string, args ...interface{}) error {
	modifiers := extractModifiers(args)

	vk, autoMods, ok := keyToVK(key)
	if !ok {
		return errors.New("robotgo: unknown key: " + key)
	}

	// Add the modifiers implied by the key itself (an uppercase letter or
	// shifted symbol needs SHIFT held), deduplicating against explicit ones.
	if autoMods&1 != 0 {
		modifiers = appendUniqueMod(modifiers, "shift")
	}
	if autoMods&2 != 0 {
		modifiers = appendUniqueMod(modifiers, "ctrl")
	}
	if autoMods&4 != 0 {
		modifiers = appendUniqueMod(modifiers, "alt")
	}

	// Press modifiers.
	for _, mod := range modifiers {
		if mvk, _, ok := keyToVK(mod); ok {
			sendVK(mvk, false)
		}
	}

	sendVK(vk, false)
	time.Sleep(time.Duration(KeySleep) * time.Millisecond)
	sendVK(vk, true)

	// Release modifiers in reverse order.
	for i := len(modifiers) - 1; i >= 0; i-- {
		if mvk, _, ok := keyToVK(modifiers[i]); ok {
			sendVK(mvk, true)
		}
	}
	return nil
}

// appendUniqueMod appends mod unless an equivalent modifier (ignoring an
// l/r prefix) is already present.
func appendUniqueMod(mods []string, mod string) []string {
	for _, m := range mods {
		if m == mod || strings.TrimLeft(m, "lr") == mod {
			return mods
		}
	}
	return append(mods, mod)
}

// KeyToggle toggles a key. Default is "down"; pass "up" to release.
func KeyToggle(key string, args ...interface{}) error {
	up := false
	for _, arg := range args {
		if s, ok := arg.(string); ok && s == "up" {
			up = true
		}
	}
	vk, _, ok := keyToVK(key)
	if !ok {
		return errors.New("robotgo: unknown key: " + key)
	}
	sendVK(vk, up)
	return nil
}

// KeyDown presses a key down. Extra args (e.g. modifiers) are forwarded to
// KeyToggle for API parity with the default robotgo backend.
func KeyDown(key string, args ...interface{}) error {
	return KeyToggle(key, append([]interface{}{"down"}, args...)...)
}

// KeyUp releases a key. Extra args (e.g. modifiers) are forwarded to
// KeyToggle for API parity with the default robotgo backend.
func KeyUp(key string, args ...interface{}) error {
	return KeyToggle(key, append([]interface{}{"up"}, args...)...)
}

// KeyPress presses a key (down + delay + up).
func KeyPress(key string, args ...interface{}) error {
	return KeyTap(key, args...)
}

// Type types a string using Unicode key events, supporting any character
// (including those not present on the current keyboard layout).
func Type(str string, args ...int) {
	for _, u := range utf16.Encode([]rune(str)) {
		sendUnicode(u, false)
		sendUnicode(u, true)
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

// CmdCtrl returns "cmd" on macOS, "ctrl" elsewhere. On Windows: "ctrl".
func CmdCtrl() string {
	return "ctrl"
}

func extractModifiers(args []interface{}) []string {
	var mods []string
	for _, arg := range args {
		if s, ok := arg.(string); ok {
			switch s {
			case "ctrl", "control", "ctrll", "ctrlr",
				"shift", "shiftl", "shiftr",
				"alt", "altl", "altr",
				"cmd", "cmdl", "cmdr", "win", "rwin":
				mods = append(mods, s)
			}
		}
	}
	return mods
}
