// Copyright 2016 The go-vgo Project Developers. See the COPYRIGHT
// file at the top-level directory of this distribution and at
// https://github.com/go-vgo/robotgo/blob/master/LICENSE
//
// Licensed under the Apache License, Version 2.0 <LICENSE-APACHE or
// http://www.apache.org/licenses/LICENSE-2.0> or the MIT license
// <LICENSE-MIT or http://opensource.org/licenses/MIT>, at your
// option. This file may not be copied, modified, or distributed
// except according to those terms.

package hook

/*
#cgo darwin CFLAGS: -x objective-c  -Wno-deprecated-declarations
#cgo darwin LDFLAGS: -framework Cocoa

#cgo linux CFLAGS:-I/usr/src
#cgo linux LDFLAGS: -L/usr/src -lX11 -lXtst
#cgo linux LDFLAGS: -lX11-xcb -lxcb -lxcb-xkb -lxkbcommon -lxkbcommon-x11
//#cgo windows LDFLAGS: -lgdi32 -luser32

#include "event/goEvent.h"
*/
import "C"

import (
	"fmt"
	"sync"
	"time"
	"unsafe"
)

const (
	// HookEnabled honk enable status
	HookEnabled  = 1 // iota
	HookDisabled = 2

	KeyDown = 3
	KeyHold = 4
	KeyUp   = 5

	MouseUp    = 6
	MouseHold  = 7
	MouseDown  = 8
	MouseMove  = 9
	MouseDrag  = 10
	MouseWheel = 11

	FakeEvent = 12
	// Keychar could be v
	CharUndefined = 0xFFFF
	WheelUp       = -1
	WheelDown     = 1
)

// Event Holds a system event
// If it's a Keyboard event the relevant fields are:
// Mask, Keycode, Rawcode, and Keychar,
// Keychar is probably what you want.
// If it's a Mouse event the relevant fields are:
// Button, Clicks, X, Y, Amount, Rotation and Direction
type Event struct {
	Kind     uint8 `json:"id"`
	When     time.Time
	Mask     uint16 `json:"mask"`
	Reserved uint16 `json:"reserved"`

	Keycode uint16 `json:"keycode"`
	Rawcode uint16 `json:"rawcode"`
	Keychar rune   `json:"keychar"`

	Button uint16 `json:"button"`
	Clicks uint16 `json:"clicks"`

	X int16 `json:"x"`
	Y int16 `json:"y"`

	Amount    uint16 `json:"amount"`
	Rotation  int32  `json:"rotation"`
	Direction uint8  `json:"direction"`
}

var (
	ev      = make(chan Event, 1024)
	asyncon = false

	lck sync.RWMutex
)

// String return hook kind string
func (e Event) String() string {
	switch e.Kind {
	case HookEnabled:
		return fmt.Sprintf("%v - Event: {Kind: HookEnabled}", e.When)
	case HookDisabled:
		return fmt.Sprintf("%v - Event: {Kind: HookDisabled}", e.When)
	case KeyUp:
		return fmt.Sprintf("%v - Event: {Kind: KeyUp, Rawcode: %v, Keychar: %v}",
			e.When, e.Rawcode, e.Keychar)
	case KeyHold:
		return fmt.Sprintf(
			"%v - Event: {Kind: KeyHold, Rawcode: %v, Keychar: %v}",
			e.When, e.Rawcode, e.Keychar)
	case KeyDown:
		return fmt.Sprintf(
			"%v - Event: {Kind: KeyDown, Rawcode: %v, Keychar: %v}",
			e.When, e.Rawcode, e.Keychar)
	case MouseUp:
		return fmt.Sprintf(
			"%v - Event: {Kind: MouseUp, Button: %v, X: %v, Y: %v, Clicks: %v}",
			e.When, e.Button, e.X, e.Y, e.Clicks)
	case MouseHold:
		return fmt.Sprintf(
			"%v - Event: {Kind: MouseHold, Button: %v, X: %v, Y: %v, Clicks: %v}",
			e.When, e.Button, e.X, e.Y, e.Clicks)
	case MouseDown:
		return fmt.Sprintf(
			"%v - Event: {Kind: MouseDown, Button: %v, X: %v, Y: %v, Clicks: %v}",
			e.When, e.Button, e.X, e.Y, e.Clicks)
	case MouseMove:
		return fmt.Sprintf(
			"%v - Event: {Kind: MouseMove, Button: %v, X: %v, Y: %v, Clicks: %v}",
			e.When, e.Button, e.X, e.Y, e.Clicks)
	case MouseDrag:
		return fmt.Sprintf(
			"%v - Event: {Kind: MouseDrag, Button: %v, X: %v, Y: %v, Clicks: %v}",
			e.When, e.Button, e.X, e.Y, e.Clicks)
	case MouseWheel:
		return fmt.Sprintf(
			"%v - Event: {Kind: MouseWheel, Amount: %v, Rotation: %v, Direction: %v}",
			e.When, e.Amount, e.Rotation, e.Direction)
	case FakeEvent:
		return fmt.Sprintf("%v - Event: {Kind: FakeEvent}", e.When)
	}

	return "Unknown event, contact the mantainers"
}

// RawcodetoKeychar rawcode to keychar
func RawcodetoKeychar(r uint16) string {
	lck.RLock()
	defer lck.RUnlock()

	return raw2key[r]
}

// KeychartoRawcode key char to rawcode
func KeychartoRawcode(kc string) uint16 {
	return keytoraw[kc]
}

// Start Adds global event hook to OS
// returns event channel
func Start() chan Event {
	asyncon = true
	go C.start_ev()

	go func() {
		for {
			C.pollEv()
			time.Sleep(time.Millisecond * 50)

			// todo: find smallest time that does not destroy the cpu utilization
			if !asyncon {
				return
			}
		}
	}()

	return ev
}

// End removes global event hook
func End() {
	C.endPoll()
	C.stop_event()

	for len(ev) != 0 {
		<-ev
	}

	asyncon = false
}

// AddEvent add event listener
func AddEvent(key string) int {
	cs := C.CString(key)
	defer C.free(unsafe.Pointer(cs))

	eve := C.add_event(cs)
	geve := int(eve)

	return geve
}

// StopEvent stop event listener
func StopEvent() {
	C.stop_event()
}
