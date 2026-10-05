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

	UnicodeType(uint32(' '))
	UnicodeType(uint32(' '), 0)
	UnicodeType(uint32(' '), 0, 0)
	inputUTF("space")
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
		tt.True(t, strings.Contains(err.Error(), "XTestFakeButtonEvent"), err.Error())
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
