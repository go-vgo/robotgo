//go:build darwin
// +build darwin

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

package darwin

import (
	"errors"

	"github.com/ebitengine/purego/objc"
)

// mediaCodes maps robotgo media and brightness key names to NX_KEYTYPE_*
// codes (IOKit ev_keymap.h); keycode.h encodes the same values as 1000+code.
var mediaCodes = map[string]int{
	"audio_vol_up": 0, "audio_vol_down": 1,
	"lights_mon_up": 2, "lights_mon_down": 3,
	"audio_mute": 7,
	"audio_play": 16, "audio_pause": 16,
	"audio_next": 17, "audio_prev": 18,
	"lights_kbd_up": 21, "lights_kbd_down": 22, "lights_kbd_toggle": 23,
}

const (
	nsEventTypeSystemDefined   = 14
	nxSubtypeAuxControlButtons = 8
	nxKeyDown                  = 0xa
	nxKeyUp                    = 0xb
)

var errMediaKey = errors.New("robotgo: media key event could not be posted")

type nsPoint struct{ X, Y float64 }

// mediaKeyData returns the NSEvent modifierFlags and data1 for an
// aux-control-button event, matching the Cgo backend's evtInfo layout.
func mediaKeyData(code int, down bool) (flags uint, data1 int) {
	state := nxKeyUp
	if down {
		state = nxKeyDown
	}
	return uint(state << 8), code<<16 | state<<8
}

// sendMediaKey posts a system-defined media key event to the HID tap, the
// NSEvent equivalent of the Cgo backend's IOHIDPostEvent(NX_SYSDEFINED).
func sendMediaKey(code int, down bool) error {
	return withMediaEvent(code, down, func(cg uintptr) {
		cgEventPost(kCGHIDEventTap, cg)
	})
}

// withMediaEvent builds the media key CGEvent and passes it to post while the
// owning NSEvent is still alive.
func withMediaEvent(code int, down bool, post func(cg uintptr)) error {
	if !loaded || !loadApp() {
		return errMediaKey
	}
	posted := false
	withPool(func() {
		flags, data1 := mediaKeyData(code, down)
		ev := objc.ID(objc.GetClass("NSEvent")).Send(objc.RegisterName(
			"otherEventWithType:location:modifierFlags:timestamp:windowNumber:context:subtype:data1:data2:"),
			uint(nsEventTypeSystemDefined), nsPoint{}, flags, float64(0), 0, objc.ID(0),
			int16(nxSubtypeAuxControlButtons), data1, -1)
		if ev == 0 {
			return
		}
		// The CGEvent is owned by the autoreleased NSEvent; do not release it.
		cg := objc.Send[uintptr](ev, objc.RegisterName("CGEvent"))
		if cg == 0 {
			return
		}
		post(cg)
		posted = true
	})
	if !posted {
		return errMediaKey
	}
	return nil
}
