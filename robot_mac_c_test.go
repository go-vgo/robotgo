//go:build darwin && !mac && !purego
// +build darwin,!mac,!purego

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

// Synthetic events must never suppress the user's own mouse/keyboard: every
// Cgo event source (MMEventSourceCreate in base/types.h) has the local-events
// suppression interval cleared; the CoreGraphics default is 0.25s per event,
// which freezes the physical mouse for as long as events keep coming.
func TestEventSourceNoLocalSuppression(t *testing.T) {
	def := eventSourceSuppression(false)
	if def < 0 {
		t.Skip("CGEventSourceCreate returned nil")
	}
	tt.True(t, def > 0, "default interval")
	tt.Equal(t, 0.0, eventSourceSuppression(true))
}

// CGEventCreateKeyboardEvent picks kCGEventFlagsChanged for modifier keycodes
// and KeyDown/KeyUp for everything else. toggleKeyCode relies on that and must
// not force the type (which left cmd/alt/shift stuck down). Create-only, no
// events are posted.
func TestKeyboardEventTypeForModifiers(t *testing.T) {
	const keyDown, keyUp, flagsChanged = 10, 11, 12

	if keyboardEventType(0, true) < 0 {
		t.Skip("CGEventCreateKeyboardEvent returned nil")
	}
	for _, name := range []string{KeyA, Enter, Space, F1, Num0} {
		code, err := checkKeyCodes(name)
		tt.Nil(t, err)
		tt.Equal(t, keyDown, keyboardEventType(int(code), true), name)
		tt.Equal(t, keyUp, keyboardEventType(int(code), false), name)
	}
	for _, name := range []string{Cmd, CmdR, Alt, AltR, Ctrl, CtrlR, Shift, ShiftR, Capslock, Fn} {
		code, err := checkKeyCodes(name)
		tt.Nil(t, err)
		tt.Equal(t, flagsChanged, keyboardEventType(int(code), true), name)
		tt.Equal(t, flagsChanged, keyboardEventType(int(code), false), name)
	}
}
