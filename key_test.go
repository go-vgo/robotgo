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

	// appendShift appends "shift" after a []string; it must be kept.
	key, args := appendShift("A", 0, []string{"ctrl"})
	_, arr = getToggleArgs(args...)
	tt.Equal(t, "a", key)
	tt.Equal(t, []string{"ctrl", "shift"}, arr)

	tt.Equal(t, C.MMKeyFlags(C.MOD_CONTROL|C.MOD_SHIFT), getFlagsFromValue([]string{"ctrl", "shiftr", "x"}))
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
