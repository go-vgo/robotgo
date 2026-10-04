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
