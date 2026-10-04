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

func TestLinuxNumpadKeyCodes(t *testing.T) {
	if runtime.GOOS != "linux" {
		t.Skip("X11 keysyms are specific to the Linux backend")
	}

	tests := []struct {
		key  string
		want int
	}{
		{key: "num+", want: 0xffab},      // XK_KP_Add
		{key: "num-", want: 0xffad},      // XK_KP_Subtract
		{key: "num*", want: 0xffaa},      // XK_KP_Multiply
		{key: "num/", want: 0xffaf},      // XK_KP_Divide
		{key: "num_enter", want: 0xff8d}, // XK_KP_Enter
	}
	for _, test := range tests {
		t.Run(test.key, func(t *testing.T) {
			got, err := checkKeyCodes(test.key)
			if err != nil {
				t.Fatal(err)
			}
			if int(got) != test.want {
				t.Errorf("checkKeyCodes(%q) = %#x, want %#x", test.key, int(got), test.want)
			}
		})
	}
}

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
