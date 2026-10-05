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

import (
	"runtime"
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

func TestAppInfo(t *testing.T) {
	name, path, pid := appInfo(0)
	tt.Equal(t, "", name)
	tt.Equal(t, "", path)
	tt.Equal(t, 0, pid)
}
