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

// Cgo-only: exercises the C-backed key/mouse internals (checkKeyCodes,
// keyFlags, tapKeyCode, formatClickError, ...). Portable key tests that run
// on every backend live in key_test.go.

package robotgo

import (
	"runtime"
	"strings"
	"syscall"
	"testing"

	"github.com/vcaesar/tt"
)

func TestKeyAliases(t *testing.T) {
	aliases := map[string]string{
		Escape:   Esc,
		Capslock: Caps,
		Control:  Ctrl,
	}
	for alias, key := range aliases {
		want, err := checkKeyCodes(key)
		tt.Nil(t, err)
		got, err := checkKeyCodes(alias)
		tt.Nil(t, err)
		tt.Equal(t, want, got)
	}

	// Unknown names must error instead of resolving to keycode 0.
	_, err := checkKeyCodes("nonexistent_key")
	tt.Equal(t, keyErr, err)
	_, err = checkKeyCodes("")
	tt.Nil(t, err)
}

// Every documented key name must resolve to a real keycode; names that are
// only placeholders on a platform are mapped to K_NOT_A_KEY and must error.
func TestKeyNamesResolve(t *testing.T) {
	for name := range keyNames {
		code, err := checkKeyCodes(name)
		if err != nil {
			tt.Equal(t, keyErr, err, name)
			continue
		}
		tt.NotZero(t, int(code), name)
	}

	// common keys must be valid on every platform
	for _, name := range []string{
		Enter, Tab, Backspace, Delete, Space, Up, Down, Left, Right,
		Home, End, Pageup, Pagedown, F1, F12, Cmd, Alt, Ctrl, Shift,
	} {
		_, err := checkKeyCodes(name)
		tt.Nil(t, err, name)
	}

	// left/right variants and long names share the base modifier code
	same := [][]string{
		{"cmd", "command"},
		{"ctrl", "control"},
		{"shiftr", "right_shift"},
		{"esc", "escape"},
		{"print", "printscreen"},
		{"num0", "numpad_0"},
	}
	for _, pair := range same {
		a, err := checkKeyCodes(pair[0])
		tt.Nil(t, err)
		b, err := checkKeyCodes(pair[1])
		tt.Nil(t, err)
		tt.Equal(t, a, b, pair[1])
	}
}

// Single characters go through keyCodeForChar (or the macOS table). Only the
// macOS table is a physical keycode where "A" and "a" share an entry; X11
// returns the keysym and VkKeyScan carries the shift state in the high byte,
// so there the case/shift variants differ and only resolution is checked.
func TestCharKeyCodes(t *testing.T) {
	sameKey := runtime.GOOS == "darwin"

	for c := 'a'; c <= 'z'; c++ {
		lower, err := checkKeyCodes(string(c))
		tt.Nil(t, err, string(c))

		upper, err := checkKeyCodes(strings.ToUpper(string(c)))
		tt.Nil(t, err)
		if sameKey {
			tt.Equal(t, lower, upper, string(c))
		}
	}

	for c := '0'; c <= '9'; c++ {
		_, err := checkKeyCodes(string(c))
		tt.Nil(t, err, string(c))
	}

	// shifted symbols map to the same physical key as their base
	for sym, base := range Special {
		s, err := checkKeyCodes(sym)
		tt.Nil(t, err, sym)
		b, err := checkKeyCodes(base)
		tt.Nil(t, err, base)
		if sameKey {
			tt.Equal(t, b, s, sym)
		}
	}

	space, err := checkKeyCodes(" ")
	tt.Nil(t, err)
	named, err := checkKeyCodes(Space)
	tt.Nil(t, err)
	tt.Equal(t, named, space)

	if runtime.GOOS == "darwin" {
		// the Go table must reject chars it does not know instead of
		// falling through to the Cgo path
		_, err = checkKeyCodes("é")
		tt.Equal(t, keyErr, err)
	}
}

func TestKeyFlags(t *testing.T) {
	none := checkKeyFlags("none")
	tt.Equal(t, none, getFlagsFromValue(nil))
	tt.Equal(t, none, getFlagsFromValue([]string{"a", "enter"}))

	groups := map[string][]string{
		"alt":   {"altl", "altr"},
		"cmd":   {"command", "cmdl", "cmdr"},
		"ctrl":  {"control", "ctrll", "ctrlr"},
		"shift": {"shiftl", "shiftr", "right_shift"},
	}
	var all []string
	for base, variants := range groups {
		flag := checkKeyFlags(base)
		tt.NotEqual(t, none, flag, base)
		for _, v := range variants {
			tt.Equal(t, flag, checkKeyFlags(v), v)
		}
		all = append(all, base)
	}

	// the four modifiers are distinct bits
	combined := getFlagsFromValue(all)
	for _, base := range all {
		tt.Equal(t, combined, combined|checkKeyFlags(base), base)
		tt.NotEqual(t, combined, combined&^checkKeyFlags(base), base)
	}
	// same flag twice is idempotent
	tt.Equal(t, checkKeyFlags("ctrl"), getFlagsFromValue([]string{"ctrl", "ctrll", "control"}))
}

func TestGetKeyDown(t *testing.T) {
	down, mods := getKeyDown(nil)
	tt.True(t, down)
	tt.Equal(t, 0, len(mods))

	down, mods = getKeyDown([]string{"down"})
	tt.True(t, down)
	tt.Equal(t, 0, len(mods))

	down, mods = getKeyDown([]string{"up"})
	tt.False(t, down)
	tt.Equal(t, 0, len(mods))

	down, mods = getKeyDown([]string{"down", "ctrl", "shift"})
	tt.True(t, down)
	tt.Equal(t, []string{"ctrl", "shift"}, mods)

	// only the first element is a direction
	down, mods = getKeyDown([]string{"ctrl", "up"})
	tt.True(t, down)
	tt.Equal(t, []string{"ctrl", "up"}, mods)

	tt.Equal(t, "down", getDown(true))
	tt.Equal(t, "up", getDown(false))
}

func TestAppendShift(t *testing.T) {
	key, args := appendShift("a", 0)
	tt.Equal(t, "a", key)
	tt.Equal(t, 0, len(args))

	// upper case adds shift and lowers the key
	key, args = appendShift("A", 0)
	tt.Equal(t, "a", key)
	tt.Equal(t, []interface{}{"shift"}, args)

	// named keys are lowered without shift
	key, args = appendShift("Enter", 0, "ctrl")
	tt.Equal(t, "enter", key)
	tt.Equal(t, []interface{}{"ctrl"}, args)

	// shifted symbols become base key + shift
	key, args = appendShift("!", 0)
	tt.Equal(t, "1", key)
	tt.Equal(t, []interface{}{"shift"}, args)

	// KeyToggle("!", "up"): len1=1 accounts for the direction arg
	key, args = appendShift("!", 1, "up")
	tt.Equal(t, "1", key)
	tt.Equal(t, []interface{}{"up", "shift"}, args)

	// explicit modifiers beyond len1 are trusted as-is
	key, args = appendShift("!", 0, "ctrl")
	tt.Equal(t, "1", key)
	tt.Equal(t, []interface{}{"ctrl"}, args)

	key, args = appendShift("", 0)
	tt.Equal(t, "", key)
	tt.Equal(t, 0, len(args))
}

func TestGetToggleArgsSkipsUnknownTypes(t *testing.T) {
	pid, arr := getToggleArgs(1.5, true, nil, "ctrl", int64(9), 7)
	tt.Equal(t, 7, pid)
	tt.Equal(t, []string{"ctrl"}, arr)

	pid, arr = getToggleArgs()
	tt.Equal(t, 0, pid)
	tt.Equal(t, 0, len(arr))

	// empty array contributes nothing
	pid, arr = getToggleArgs([]string{}, "alt")
	tt.Equal(t, 0, pid)
	tt.Equal(t, []string{"alt"}, arr)
}

// Unknown keys must come back as keyErr from the internal tap/toggle paths
// before any C call; modifiers in the array do not change that.
func TestKeyErrBeforeC(t *testing.T) {
	const bad = "no_such_key"
	tt.Equal(t, keyErr, keyTaps(bad, nil, 0))
	tt.Equal(t, keyErr, keyTaps(bad, []string{"ctrl", "shift"}, 0))
	tt.Equal(t, keyErr, keyToggles(bad, nil, 0))
	tt.Equal(t, keyErr, keyToggles(bad, []string{"up", "alt"}, 0))
	tt.Equal(t, keyErr, keyTogglesB(bad, true, nil, 0))
	tt.Equal(t, keyErr, keyTogglesB(bad, false, []string{"cmd"}, 123))

	// public wrappers surface the same sentinel
	tt.Equal(t, keyErr, KeyTap(bad))
	tt.Equal(t, keyErr, KeyToggle(bad, "up"))
	tt.Equal(t, keyErr, KeyDown(bad, "ctrl"))
	tt.Equal(t, keyErr, KeyUp(bad))
	tt.Equal(t, keyErr, KeyPress(bad))
}

// Cgo key injection paths that need a session: tapKeyCode reports the stage
// of the failed toggle, upKeyArr/keyTaps skip unknown modifier names.
func TestKeyInjectCgo(t *testing.T) {
	requireDisplay(t)

	code, err := checkKeyCodes(Shift)
	tt.Nil(t, err)
	c, stage := tapKeyCode(code, checkKeyFlags("none"), 0)
	tt.Equal(t, 0, c)
	tt.Equal(t, "up", stage)

	upKeyArr([]string{"no_such_key", Shift}, 0)
	tt.Nil(t, keyTaps(Shift, []string{"no_such_key"}, 0))
	tt.Nil(t, keyTogglesB(Shift, true, nil, 0))
	tt.Nil(t, keyTogglesB(Shift, false, []string{"no_such_key"}, 0))
	tt.Nil(t, keyToggles(Shift, []string{"down"}, 0))
	tt.Nil(t, keyToggles(Shift, []string{"up"}, 0))
}

// UnicodeType and inputUTF are Cgo-only (the pure-Go backends type via
// Type); inputUTF is a no-op stub outside X11 and must not crash anywhere.
func TestUnicodeTypeCgo(t *testing.T) {
	requireDisplay(t)

	tt.Nil(t, UnicodeType(uint32(' ')))
	tt.Nil(t, UnicodeType(uint32(' '), 0))
	tt.Nil(t, UnicodeType(uint32(' '), 0, 0))
	tt.Nil(t, inputUTF("space"))
}

func TestInputUTFX11(t *testing.T) {
	if runtime.GOOS != "linux" {
		t.Skip("inputUTF uses Xlib only on Linux")
	}
	requireDisplay(t)

	for _, sym := range ToUC("世😀") {
		if err := inputUTF(sym); err != nil {
			t.Errorf("inputUTF(%q): %v", sym, err)
		}
	}
	for _, sym := range []string{"", "no_such_keysym", `\U0001f600`} {
		if err := inputUTF(sym); err == nil {
			t.Errorf("inputUTF(%q) succeeded for an invalid keysym", sym)
		}
	}
}

// Window ops must report failure for a pid with no window instead of
// silently doing nothing.
func TestWindowOpsBogusPid(t *testing.T) {
	const bogus = 0x7ffffff0
	if err := MinWindow(bogus); err == nil {
		t.Error("MinWindow(bogus): got nil, want error")
	}
	if err := MaxWindow(bogus, false); err == nil {
		t.Error("MaxWindow(bogus): got nil, want error")
	}
	if err := CloseWindow(bogus); err == nil {
		t.Error("CloseWindow(bogus): got nil, want error")
	}
}

func TestTypeStrInvalidKeysymX11(t *testing.T) {
	if runtime.GOOS != "linux" {
		t.Skip("TypeStr uses Xlib keysyms only on Linux")
	}
	requireDisplay(t)

	// ToUC leaves the newline as an escape that Xlib cannot resolve.
	n, err := typeStr("a\nb", 0, 0, 0)
	if n != 1 || err == nil {
		t.Fatalf("typeStr returned (%d, %v), want (1, error)", n, err)
	}
	if err := TypeStr("\n", 0, 0, 0); err == nil {
		t.Error("TypeStr succeeded for an invalid keysym")
	}
	if n := Type("\n", 0, 0, 0); n != 0 {
		t.Errorf("Type returned %d, want 0", n)
	}
}

func TestCheckMouse(t *testing.T) {
	for _, btn := range []string{Mleft, Mright, Center, WheelDown, WheelUp, WheelLeft, WheelRight} {
		tt.Equal(t, btn, MouseButtonString(CheckMouse(btn)))
	}
	// unknown names default to the left button
	tt.Equal(t, CheckMouse(Mleft), CheckMouse("unknown"))
	tt.NotEqual(t, CheckMouse(Mleft), CheckMouse(Mright))
}

func TestFormatClickErrorButton(t *testing.T) {
	err := formatClickError(1, CheckMouse(Mright), "down", 2)
	tt.NotNil(t, err)
	tt.True(t, strings.Contains(err.Error(), "click down failed (right, count=2)"), err.Error())
	tt.True(t, strings.Contains(err.Error(), "code=1"), err.Error())

	switch runtime.GOOS {
	case "darwin":
		err = formatClickError(1000, "a", "up", 1)
		tt.True(t, strings.Contains(err.Error(), "kCGErrorFailure"), err.Error())
		// unknown CG codes still report the code
		err = formatClickError(42, "a", "up", 1)
		tt.True(t, strings.HasSuffix(err.Error(), "code=42"), err.Error())
	case "linux":
		err = formatClickError(1, "a", "up", 1)
		tt.True(t, strings.Contains(err.Error(), "XTest request failed"), err.Error())
		err = formatClickError(2, "a", "up", 1)
		tt.True(t, strings.Contains(err.Error(), "no X display"), err.Error())
	case "windows":
		err = formatClickError(5, "a", "up", 1)
		tt.True(t, strings.Contains(err.Error(), syscall.Errno(5).Error()), err.Error())
	}
}

func TestAlertArgs(t *testing.T) {
	ok, cancel := alertArgs()
	tt.Equal(t, "Ok", ok)
	tt.Equal(t, "Cancel", cancel)

	ok, cancel = alertArgs("Yes")
	tt.Equal(t, "Yes", ok)
	tt.Equal(t, "Cancel", cancel)

	ok, cancel = alertArgs("Yes", "No", "ignored")
	tt.Equal(t, "Yes", ok)
	tt.Equal(t, "No", cancel)
}

func TestDisplayIdx(t *testing.T) {
	saved := DisplayID
	defer func() { DisplayID = saved }()

	DisplayID = -1
	tt.Equal(t, -1, displayIdx())
	tt.Equal(t, 2, displayIdx(2))

	// the global is the default, an explicit id still wins
	DisplayID = 1
	tt.Equal(t, 1, displayIdx())
	tt.Equal(t, 0, displayIdx(0))
}

func TestLinuxNumpadKeyCodes(t *testing.T) {
	if runtime.GOOS != "linux" {
		t.Skip("X11 keysyms are specific to the Linux backend")
	}

	keys := map[string]int{
		"num+":      0xffab, // XK_KP_Add
		"num-":      0xffad, // XK_KP_Subtract
		"num*":      0xffaa, // XK_KP_Multiply
		"num/":      0xffaf, // XK_KP_Divide
		"num_enter": 0xff8d, // XK_KP_Enter
		"num_equal": 0xffbd, // XK_KP_Equal
	}
	for key, want := range keys {
		got, err := checkKeyCodes(key)
		tt.Nil(t, err)
		tt.Equal(t, want, int(got))
	}
}

func TestGetToggleArgs(t *testing.T) {
	pid, arr := getToggleArgs([]string{"alt", "cmd"})
	tt.Equal(t, 0, pid)
	tt.Equal(t, []string{"alt", "cmd"}, arr)

	pid, arr = getToggleArgs("up", []string{"alt", "cmd"})
	tt.Equal(t, 0, pid)
	tt.Equal(t, []string{"up", "alt", "cmd"}, arr)

	pid, arr = getToggleArgs(123, []string{"ctrl"}, "shift")
	tt.Equal(t, 123, pid)
	tt.Equal(t, []string{"ctrl", "shift"}, arr)

	// KeyUp("a", pid, "ctrl") prepends "up"; the pid must still be found.
	pid, arr = getToggleArgs("up", 123, "ctrl")
	tt.Equal(t, 123, pid)
	tt.Equal(t, []string{"up", "ctrl"}, arr)

	// Only the first int is the pid.
	pid, arr = getToggleArgs([]string{"ctrl"}, 123, 456)
	tt.Equal(t, 123, pid)
	tt.Equal(t, []string{"ctrl"}, arr)

	// appendShift appends "shift" after a []string; it must be kept.
	key, args := appendShift("A", 0, []string{"ctrl"})
	_, arr = getToggleArgs(args...)
	tt.Equal(t, "a", key)
	tt.Equal(t, []string{"ctrl", "shift"}, arr)

	want := checkKeyFlags("ctrl") | checkKeyFlags("shift")
	tt.Equal(t, want, getFlagsFromValue([]string{"ctrl", "shiftr", "x"}))
	tt.Equal(t, checkKeyFlags("none"), checkKeyFlags("x"))
	tt.Equal(t, checkKeyFlags("shift"), checkKeyFlags("right_shift"))
}

// KeyToggle("a", "up", []string{"alt", "cmd"}): the direction is stripped and
// the array becomes the modifier flags (this panicked before the flattening).
func TestKeyToggleArrayArgs(t *testing.T) {
	key, args := appendShift("a", 1, "up", []string{"alt", "cmd"})
	tt.Equal(t, "a", key)
	pid, arr := getToggleArgs(args...)
	tt.Equal(t, 0, pid)

	down, mods := getKeyDown(arr)
	tt.False(t, down)
	tt.Equal(t, []string{"alt", "cmd"}, mods)
	tt.Equal(t, checkKeyFlags("alt")|checkKeyFlags("cmd"), getFlagsFromValue(mods))

	// no direction defaults to down and keeps the array intact
	down, mods = getKeyDown([]string{"alt", "cmd"})
	tt.True(t, down)
	tt.Equal(t, []string{"alt", "cmd"}, mods)
}

func TestFormatClickErrorKey(t *testing.T) {
	tt.Nil(t, formatClickError(0, "a", "down", 1))
	tt.NotNil(t, formatClickError(5, "a", "up", 1))
}

func TestFormatMouseError(t *testing.T) {
	tt.Nil(t, formatMouseError(0, "move"))

	err := formatMouseError(42, "scroll")
	tt.NotNil(t, err)
	tt.True(t, strings.HasPrefix(err.Error(), "mouse scroll failed"), err.Error())
	tt.True(t, strings.Contains(err.Error(), "code=42"), err.Error())

	switch runtime.GOOS {
	case "darwin":
		err = formatMouseError(1004, "move")
		tt.True(t, strings.Contains(err.Error(), "kCGErrorCannotComplete"), err.Error())
	case "windows":
		// GetLastError codes are described via syscall.Errno
		tt.True(t, strings.Contains(err.Error(), syscall.Errno(42).Error()), err.Error())
	}
}

func TestFormatKeyError(t *testing.T) {
	tt.Nil(t, formatKeyError(0, 'a'))

	err := formatKeyError(42, 'é')
	tt.NotNil(t, err)
	tt.True(t, strings.HasPrefix(err.Error(), `type 'é' failed`), err.Error())
	tt.True(t, strings.Contains(err.Error(), "code=42"), err.Error())
}

// codeDetail only describes codes it knows; everything else is empty so the
// callers fall back to the bare "code=N" form.
func TestCodeDetail(t *testing.T) {
	switch runtime.GOOS {
	case "windows":
		tt.Equal(t, syscall.Errno(5).Error(), codeDetail(5))
		// MM_ERR_INPUT_BLOCKED: SendInput failed without a last error.
		tt.Equal(t, "input blocked (UIPI or secure desktop)", codeDetail(-1))
		// win32KeyEvent: no window for the pid.
		tt.Equal(t, "window not found", codeDetail(-5))
		tt.Equal(t, "", codeDetail(-42))
	case "darwin":
		tt.Equal(t, "kCGErrorCannotComplete", codeDetail(1004))
		tt.Equal(t, "", codeDetail(42))
	default:
		tt.Equal(t, "XTest request failed", codeDetail(1))
		tt.Equal(t, "no X display", codeDetail(2))
		tt.Equal(t, "no X display", codeDetail(-8))
		tt.Equal(t, "", codeDetail(42))
	}
}

// Wrong smooth-move args must be rejected instead of panicking, otherwise
// DragSmooth would leave the button held down.
func TestSmoothArgs(t *testing.T) {
	low, high, delay, err := smoothArgs(nil)
	tt.Nil(t, err)
	tt.Equal(t, 1.0, low)
	tt.Equal(t, 3.0, high)
	tt.Equal(t, 1, delay)

	low, high, delay, err = smoothArgs([]interface{}{2.0, 5.0, 7})
	tt.Nil(t, err)
	tt.Equal(t, 2.0, low)
	tt.Equal(t, 5.0, high)
	tt.Equal(t, 7, delay)

	_, _, _, err = smoothArgs([]interface{}{"right", 5.0})
	tt.NotNil(t, err)
	_, _, _, err = smoothArgs([]interface{}{2.0, 5.0, "7"})
	tt.NotNil(t, err)
	// DragSmooth/MoveSmooth report the bad args without touching the mouse.
	tt.NotNil(t, DragSmooth(1, 1, "right", 5.0))
	tt.False(t, MoveSmooth(1, 1, "right", 5.0))
}

func TestInputUTFDetail(t *testing.T) {
	for code := 1; code <= 4; code++ {
		tt.NotEqual(t, "", inputUTFDetail[code])
	}
}

// An invalid direction must error out before any scroll event is posted.
func TestScrollDirInvalid(t *testing.T) {
	err := ScrollDir(1, "sideways")
	tt.NotNil(t, err)
	tt.Equal(t, "unknown scroll direction: sideways", err.Error())

	err = ScrollDir(1, 3)
	tt.NotNil(t, err)
	tt.Equal(t, "unknown scroll direction: 3", err.Error())
}
