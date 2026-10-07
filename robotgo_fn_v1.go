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

import "github.com/vcaesar/tt"

// Deprecated: use the Move(),
//
// MoveMouse move the mouse
func MoveMouse(x, y int) {
	Move(x, y) //nolint:errcheck // v1 signature has no error result
}

// Deprecated: use the DragSmooth(),
//
// DragMouse drag the mouse to (x, y),
// It's same with the DragSmooth() now
func DragMouse(x, y int, args ...interface{}) {
	DragSmooth(x, y, args...) //nolint:errcheck // v1 signature has no error result
}

// Deprecated: use the MoveSmooth(),
//
// MoveMouseSmooth move the mouse smooth,
// moves mouse to x, y human like, with the mouse button up.
func MoveMouseSmooth(x, y int, args ...interface{}) bool {
	return MoveSmooth(x, y, args...)
}

// Deprecated: use the function Location()
//
// GetMousePos get the mouse's position return x, y
func GetMousePos() (int, int) {
	return Location()
}

// Deprecated: use the Click(),
//
// # MouseClick click the mouse
//
// robotgo.MouseClick(button string, double bool)
func MouseClick(args ...interface{}) {
	Click(args...) //nolint:errcheck // v1 signature has no error result
}

// Deprecated: use the TypeStr(),
//
// # TypeStringDelayed type string delayed, Wno-deprecated
//
// This function will be removed in version v1.0.0
func TypeStringDelayed(str string, delay int) {
	tt.Drop("TypeStringDelayed", "TypeStrDelay")
	TypeStrDelay(str, delay) //nolint:errcheck // v1 signature has no error result
}
