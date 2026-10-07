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
	"fmt"
	"strings"
	"unicode/utf16"
	"unicode/utf8"
	"unsafe"

	"github.com/go-vgo/robotgo/pub"
	"github.com/tailscale/win"
)

// NotPid, when true, makes the pid argument of the keyboard APIs be treated as
// a window handle (HWND) directly instead of a process id, mirroring the
// default Cgo backend's robotgo.NotPid behavior on Windows.
var NotPid bool

// Window messages used to post key events to a specific window (pid-directed
// input), mirroring the PostMessageW path in key/keypress_c.h.
const (
	wmKeyDown = 0x0100
	wmKeyUp   = 0x0101
	wmChar    = 0x0102
)

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
	Caps     = "caps"
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
	"right_shift": win.VK_RSHIFT,
	"ctrl":        win.VK_CONTROL, "control": win.VK_CONTROL,
	"ctrll": win.VK_LCONTROL, "ctrlr": win.VK_RCONTROL,
	"alt": win.VK_MENU, "altl": win.VK_LMENU, "altr": win.VK_RMENU,
	"cmd": win.VK_LWIN, "command": win.VK_LWIN, "cmdl": win.VK_LWIN, "win": win.VK_LWIN,
	"cmdr": win.VK_RWIN, "rwin": win.VK_RWIN,
	"capslock": win.VK_CAPITAL,
	"caps":     win.VK_CAPITAL,
	"print":    win.VK_SNAPSHOT, "printscreen": win.VK_SNAPSHOT,
	"menu":     win.VK_APPS,
	"num_lock": win.VK_NUMLOCK, "scroll_lock": win.VK_SCROLL,
	"pause": win.VK_PAUSE, "pause_break": win.VK_PAUSE,

	"f1": win.VK_F1, "f2": win.VK_F2, "f3": win.VK_F3, "f4": win.VK_F4,
	"f5": win.VK_F5, "f6": win.VK_F6, "f7": win.VK_F7, "f8": win.VK_F8,
	"f9": win.VK_F9, "f10": win.VK_F10, "f11": win.VK_F11, "f12": win.VK_F12,
	"f13": win.VK_F13, "f14": win.VK_F14, "f15": win.VK_F15, "f16": win.VK_F16,
	"f17": win.VK_F17, "f18": win.VK_F18, "f19": win.VK_F19, "f20": win.VK_F20,
	"f21": win.VK_F21, "f22": win.VK_F22, "f23": win.VK_F23, "f24": win.VK_F24,

	"num0": win.VK_NUMPAD0, "num1": win.VK_NUMPAD1, "num2": win.VK_NUMPAD2,
	"num3": win.VK_NUMPAD3, "num4": win.VK_NUMPAD4, "num5": win.VK_NUMPAD5,
	"num6": win.VK_NUMPAD6, "num7": win.VK_NUMPAD7, "num8": win.VK_NUMPAD8,
	"num9": win.VK_NUMPAD9,
	"num.": win.VK_DECIMAL, "num+": win.VK_ADD, "num-": win.VK_SUBTRACT,
	"num*": win.VK_MULTIPLY, "num/": win.VK_DIVIDE, "num_enter": win.VK_RETURN,
	"num_clear": win.VK_CLEAR,
	"num_equal": win.VK_OEM_PLUS, // same as the Cgo K_NUMPAD_EQUAL
	// deprecated numpad_* aliases (see keyNames in the root key.go)
	"numpad_0": win.VK_NUMPAD0, "numpad_1": win.VK_NUMPAD1, "numpad_2": win.VK_NUMPAD2,
	"numpad_3": win.VK_NUMPAD3, "numpad_4": win.VK_NUMPAD4, "numpad_5": win.VK_NUMPAD5,
	"numpad_6": win.VK_NUMPAD6, "numpad_7": win.VK_NUMPAD7, "numpad_8": win.VK_NUMPAD8,
	"numpad_9": win.VK_NUMPAD9, "numpad_lock": win.VK_NUMLOCK,

	"audio_mute": win.VK_VOLUME_MUTE, "audio_vol_down": win.VK_VOLUME_DOWN,
	"audio_vol_up": win.VK_VOLUME_UP, "audio_play": win.VK_MEDIA_PLAY_PAUSE,
	"audio_pause": win.VK_MEDIA_PLAY_PAUSE, "audio_stop": win.VK_MEDIA_STOP,
	"audio_prev": win.VK_MEDIA_PREV_TRACK, "audio_next": win.VK_MEDIA_NEXT_TRACK,
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
	win.VK_NUMLOCK: true, win.VK_SNAPSHOT: true, win.VK_PAUSE: true,
	win.VK_LWIN: true, win.VK_RWIN: true, win.VK_APPS: true,
	win.VK_DIVIDE:      true,
	win.VK_VOLUME_MUTE: true, win.VK_VOLUME_DOWN: true, win.VK_VOLUME_UP: true,
	win.VK_MEDIA_PLAY_PAUSE: true, win.VK_MEDIA_STOP: true,
	win.VK_MEDIA_PREV_TRACK: true, win.VK_MEDIA_NEXT_TRACK: true,
}

var procMapVirtualKeyW = modUser32.NewProc("MapVirtualKeyW")

// mapvkVkToVsc is the MapVirtualKey translation from VK to scan code.
const mapvkVkToVsc = 0

// scanCode returns the hardware scan code for a virtual key. Raw-input
// consumers (games, RDP, low-level hooks) ignore events whose scan code is
// zero, so it is filled in like the Cgo backend's win32KeyEvent does.
func scanCode(vk uint16) uint16 {
	r, _, _ := procMapVirtualKeyW.Call(uintptr(vk), mapvkVkToVsc)
	return uint16(r)
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
	if len(r) == 1 && r[0] <= 0xffff {
		res := win.VkKeyScan(uint16(r[0]))
		if res != -1 {
			return uint16(byte(res & 0xff)), uint8((res >> 8) & 0xff), true
		}
	}
	return 0, 0, false
}

// sendVK dispatches a single virtual-key event (down or up).
func sendVK(vk uint16, up bool) error {
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
			WScan:   scanCode(vk),
			DwFlags: flags,
		},
	}
	if win.SendInput(1, unsafe.Pointer(&in), int32(unsafe.Sizeof(in))) != 1 {
		return errSendInput
	}
	return nil
}

// sendUnicode dispatches a single UTF-16 code unit as a Unicode key event.
// It reports whether SendInput inserted the event.
func sendUnicode(u uint16, up bool) bool {
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
	return win.SendInput(1, unsafe.Pointer(&in), int32(unsafe.Sizeof(in))) == 1
}

// hwndByPid returns the first top-level window owned by pid. It mirrors the C
// GetHwndByPid helper in base/pubs.h and does not require the window to be
// visible.
func hwndByPid(pid int) win.HWND {
	var found win.HWND
	enumWindows(func(hwnd win.HWND) bool {
		if windowPid(hwnd) == pid {
			found = hwnd
			return false // stop
		}
		return true
	})
	return found
}

// keyHwnd resolves the target window for pid-directed key input. When NotPid is
// set the value is treated as an HWND directly (matching robotgo.NotPid in the
// default backend); otherwise the first window owned by that pid is returned.
func keyHwnd(pid int) win.HWND {
	if NotPid {
		return win.HWND(uintptr(pid))
	}
	return hwndByPid(pid)
}

// postKey posts a WM_KEYDOWN/WM_KEYUP message for a virtual-key code to a
// window, mirroring the PostMessageW path in key/keypress_c.h.
func postKey(hwnd win.HWND, vk uint16, up bool) error {
	msg := uintptr(wmKeyDown)
	if up {
		msg = wmKeyUp
	}
	if r, _, _ := procPostMessageW.Call(uintptr(hwnd), msg, uintptr(vk), 0); r == 0 {
		return errPostMessage
	}
	return nil
}

// postChar posts a WM_CHAR message carrying a single UTF-16 code unit to a
// window, mirroring unicodeType()'s PostMessageW path in key/keypress_c.h.
// It reports whether PostMessageW succeeded.
func postChar(hwnd win.HWND, u uint16) bool {
	r, _, _ := procPostMessageW.Call(uintptr(hwnd), uintptr(wmChar), uintptr(u), 0)
	return r != 0
}

// KeyTap taps a key (press + release). Optional trailing modifiers (strings
// or a []string) and an int pid: when a pid is supplied the key (and its
// modifiers) is posted to that
// process's window via PostMessageW, mirroring the Windows path in
// key/keypress_c.h; otherwise it is injected into the focused window via
// SendInput. Set NotPid to pass an HWND instead of a pid.
//
//	KeyTap("a")
//	KeyTap("a", "ctrl")
//	KeyTap("a", "ctrl", "shift")
//	KeyTap("a", []string{"ctrl", "shift"})
//	KeyTap("a", pid)
//	KeyTap("a", pid, "ctrl")
func KeyTap(key string, args ...interface{}) error {
	vks, _, err := toggleKeys(key, args)
	if err != nil {
		return err
	}
	send, err := keySender(extractPid(args))
	if err != nil {
		return err
	}

	// Press modifiers then the key, and release in reverse order. Only the
	// keys that were actually pressed are released after a failure, so none
	// is left stuck down and keys held by the user are not released.
	pressed, err := pressKeys(send, vks)
	pub.MilliSleep(pub.KeySleep)
	if upErr := releaseKeys(send, pressed); err == nil {
		err = upErr
	}
	return err
}

// pressKeys sends key down for vks in order and stops at the first error,
// returning the prefix that was successfully pressed.
func pressKeys(send func(vk uint16, up bool) error, vks []uint16) ([]uint16, error) {
	for i, vk := range vks {
		if err := send(vk, false); err != nil {
			return vks[:i], err
		}
	}
	return vks, nil
}

// releaseKeys sends key up for vks in reverse order and returns the first
// error.
func releaseKeys(send func(vk uint16, up bool) error, vks []uint16) error {
	var err error
	for i := len(vks) - 1; i >= 0; i-- {
		if e := send(vks[i], true); err == nil {
			err = e
		}
	}
	return err
}

// keySender returns the function delivering virtual-key events: PostMessageW
// to the pid's window (mirroring key/keypress_c.h) when pid != 0, otherwise
// SendInput to the focused window.
func keySender(pid int) (func(vk uint16, up bool) error, error) {
	if pid == 0 {
		return sendVK, nil
	}
	hwnd := keyHwnd(pid)
	if hwnd == 0 {
		return nil, ErrNotFound
	}
	return func(vk uint16, up bool) error { return postKey(hwnd, vk, up) }, nil
}

// toggleKeys resolves a key and its modifiers into virtual keys in press
// order (explicit modifiers, modifiers implied by the key such as shift for
// uppercase letters, then the key) and the direction; the last "up" or
// "down" argument wins. Virtual keys are deduplicated so aliases press once.
func toggleKeys(key string, args []interface{}) (vks []uint16, up bool, err error) {
	vk, autoMods, ok := keyToVK(key)
	if !ok {
		return nil, false, errors.New("robotgo: unknown key: " + key)
	}
	for _, s := range keyArgs(args) {
		switch s {
		case "up":
			up = true
		case "down":
			up = false
		}
	}

	// Add the modifiers implied by the key itself, deduplicating against
	// explicit ones (including left/right variants).
	mods := extractModifiers(args)
	for i, mod := range []string{"shift", "ctrl", "alt"} {
		if autoMods&(1<<i) != 0 {
			mods = appendUniqueMod(mods, mod)
		}
	}
	seen := map[uint16]bool{vk: true}
	for _, mod := range mods {
		if mvk, _, ok := keyToVK(mod); ok && !seen[mvk] {
			seen[mvk] = true
			vks = append(vks, mvk)
		}
	}
	return append(vks, vk), up, nil
}

// appendUniqueMod appends mod unless an equivalent modifier (including
// left/right variants and aliases) is already present.
func appendUniqueMod(mods []string, mod string) []string {
	for _, m := range mods {
		if m == mod || m == mod+"l" || m == mod+"r" ||
			(mod == "ctrl" && m == "control") || (mod == "shift" && m == "right_shift") {
			return mods
		}
	}
	return append(mods, mod)
}

// KeyToggle toggles a key. Default is "down"; pass "up" to release. Modifiers
// (strings or a []string) are pressed before the key and released in reverse
// order after it. An optional int pid posts the events to that process's
// window via PostMessageW, mirroring key/keypress_c.h; otherwise SendInput is
// used. Set NotPid to pass an HWND.
//
//	KeyToggle("a")
//	KeyToggle("a", "up")
//	KeyToggle("a", "down", []string{"ctrl", "shift"})
//	KeyToggle("a", "up", []string{"ctrl", "shift"})
//	KeyToggle("a", pid)
func KeyToggle(key string, args ...interface{}) error {
	vks, up, err := toggleKeys(key, args)
	if err != nil {
		return err
	}
	send, err := keySender(extractPid(args))
	if err != nil {
		return err
	}
	if up {
		return releaseKeys(send, vks)
	}
	// A partial press releases the keys pressed so far instead of leaving
	// modifiers stuck down.
	pressed, err := pressKeys(send, vks)
	if err != nil {
		releaseKeys(send, pressed) //nolint:errcheck // best-effort cleanup, the press error is reported
	}
	return err
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
// (including those not present on the current keyboard layout). An optional
// first int argument is the target pid: when non-zero each character is posted
// to that process's window as a WM_CHAR message, mirroring unicodeType() in
// key/keypress_c.h; otherwise it is injected into the focused window via
// SendInput. Set NotPid to pass an HWND instead of a pid.
//
//	Type("hello")
//	Type("hello", pid)
//
// It returns the number of characters (runes) typed, stopping at the first
// failed input event; 0 if the target window is not found.
func Type(str string, args ...int) int {
	pid := 0
	if len(args) > 0 {
		pid = args[0]
	}
	send := func(u uint16) bool { return sendUnicode(u, false) && sendUnicode(u, true) }
	if pid != 0 {
		hwnd := keyHwnd(pid)
		if hwnd == 0 {
			return 0
		}
		send = func(u uint16) bool { return postChar(hwnd, u) }
	}
	return typeRunes(str, send)
}

// typeRunes sends each rune of str as UTF-16 code units via send and returns
// the number of runes fully sent, stopping at the first failure.
func typeRunes(str string, send func(u uint16) bool) int {
	n := 0
	for _, r := range str {
		for _, u := range utf16.Encode([]rune{r}) {
			if !send(u) {
				return n
			}
		}
		n++
		pub.MilliSleep(pub.KeySleep)
	}
	return n
}

// TypeStr types a string, mirroring the robotgo API. It returns an error if
// not every character was typed.
func TypeStr(str string, args ...int) error {
	return typeErr(Type(str, args...), str)
}

// TypeDelay types a string with a per-character delay in milliseconds.
func TypeDelay(str string, delay int) error {
	old := pub.KeySleep
	pub.KeySleep = delay
	n := Type(str)
	pub.KeySleep = old
	return typeErr(n, str)
}

// typeErr reports an error when fewer than all runes of str were typed.
func typeErr(n int, str string) error {
	if total := utf8.RuneCountInString(str); n < total {
		return fmt.Errorf("robotgo: typed %d of %d characters", n, total)
	}
	return nil
}

// CmdCtrl returns "cmd" on macOS, "ctrl" elsewhere. On Windows: "ctrl".
func CmdCtrl() string {
	return "ctrl"
}

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
			"cmd", "command", "cmdl", "cmdr", "win", "rwin":
			mods = append(mods, s)
		}
	}
	return mods
}

// extractPid returns the first int argument as the target pid (0 means inject
// into the focused window via SendInput), matching the default robotgo
// backend's pid handling.
func extractPid(args []interface{}) int {
	for _, arg := range args {
		if v, ok := arg.(int); ok {
			return v
		}
	}
	return 0
}
