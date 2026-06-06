//go:build windows
// +build windows

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

package win

import (
	"math"
	"time"
	"unsafe"

	"github.com/tailscale/win"
)

// WheelDelta is one notch of mouse-wheel movement (WHEEL_DELTA).
const wheelDelta = 120

// MouseSleep is the global mouse delay in milliseconds.
var MouseSleep = 0

// sendMouseInput dispatches a single synthesized mouse event.
func sendMouseInput(flags, mouseData uint32, dx, dy int32) {
	in := win.MOUSE_INPUT{
		Type: win.INPUT_MOUSE,
		Mi: win.MOUSEINPUT{
			Dx:        dx,
			Dy:        dy,
			MouseData: mouseData,
			DwFlags:   flags,
		},
	}
	win.SendInput(1, unsafe.Pointer(&in), int32(unsafe.Sizeof(in)))
}

// mouseButtonFlags returns the (down, up) MOUSEEVENTF flags for a button name.
func mouseButtonFlags(btn string) (down, up uint32) {
	switch btn {
	case "right":
		return win.MOUSEEVENTF_RIGHTDOWN, win.MOUSEEVENTF_RIGHTUP
	case "center", "middle":
		return win.MOUSEEVENTF_MIDDLEDOWN, win.MOUSEEVENTF_MIDDLEUP
	default: // "left"
		return win.MOUSEEVENTF_LEFTDOWN, win.MOUSEEVENTF_LEFTUP
	}
}

// Move moves the mouse to absolute position (x, y).
// The optional displayId is accepted for API parity but ignored on Windows,
// where coordinates are in the unified virtual-desktop space.
func Move(x, y int, displayId ...int) {
	win.SetCursorPos(int32(x), int32(y))
	mouseDelay()
}

// MoveRelative moves the mouse relative to its current position.
func MoveRelative(x, y int) {
	cx, cy := Location()
	win.SetCursorPos(int32(cx+x), int32(cy+y))
	mouseDelay()
}

// MoveSmooth moves the mouse smoothly to (x, y) with an ease-in-out curve.
// Optional args: steps (int), sleepMs (int). Returns true on success.
func MoveSmooth(x, y int, args ...interface{}) bool {
	steps := 20
	sleepMs := 5
	if len(args) >= 1 {
		if v, ok := args[0].(int); ok {
			steps = v
		}
	}
	if len(args) >= 2 {
		if v, ok := args[1].(int); ok {
			sleepMs = v
		}
	}
	if steps < 1 {
		steps = 1
	}

	sx, sy := Location()
	for i := 1; i <= steps; i++ {
		t := float64(i) / float64(steps)
		// Ease-in-out cubic.
		if t < 0.5 {
			t = 4 * t * t * t
		} else {
			t = 1 - math.Pow(-2*t+2, 3)/2
		}
		cx := float64(sx) + float64(x-sx)*t
		cy := float64(sy) + float64(y-sy)*t
		win.SetCursorPos(int32(math.Round(cx)), int32(math.Round(cy)))
		time.Sleep(time.Duration(sleepMs) * time.Millisecond)
	}
	return true
}

// Click clicks a mouse button. Default is the left button.
// Use "left", "right", or "center". Pass true for a double-click.
func Click(args ...interface{}) error {
	button := "left"
	double := false
	for _, arg := range args {
		switch v := arg.(type) {
		case string:
			button = v
		case bool:
			double = v
		}
	}

	down, up := mouseButtonFlags(button)
	count := 1
	if double {
		count = 2
	}
	for i := 0; i < count; i++ {
		sendMouseInput(down, 0, 0, 0)
		sendMouseInput(up, 0, 0, 0)
		if i < count-1 {
			time.Sleep(50 * time.Millisecond)
		}
	}
	mouseDelay()
	return nil
}

// Toggle toggles a mouse button down or up.
//
//	Toggle("left")        // press down
//	Toggle("left", "up")  // release
func Toggle(key ...interface{}) error {
	button := "left"
	up := false
	for _, arg := range key {
		if s, ok := arg.(string); ok {
			switch s {
			case "up":
				up = true
			case "down":
				up = false
			default:
				button = s
			}
		}
	}

	downFlag, upFlag := mouseButtonFlags(button)
	if up {
		sendMouseInput(upFlag, 0, 0, 0)
	} else {
		sendMouseInput(downFlag, 0, 0, 0)
	}
	return nil
}

// MouseDown sends a mouse button down event.
func MouseDown(key ...interface{}) error {
	args := append([]interface{}{}, key...)
	args = append(args, "down")
	return Toggle(args...)
}

// MouseUp sends a mouse button up event.
func MouseUp(key ...interface{}) error {
	args := append([]interface{}{}, key...)
	args = append(args, "up")
	return Toggle(args...)
}

// Scroll scrolls the mouse. Positive y scrolls down, negative scrolls up;
// positive x scrolls right, negative scrolls left. Optional arg: delay ms.
func Scroll(x, y int, args ...int) {
	msDelay := 10
	if len(args) > 0 {
		msDelay = args[0]
	}

	if y != 0 {
		// Win32 wheel: positive delta scrolls up, so negate for down-positive.
		sendMouseInput(win.MOUSEEVENTF_WHEEL, uint32(int32(-y*wheelDelta)), 0, 0)
	}
	if x != 0 {
		sendMouseInput(win.MOUSEEVENTF_HWHEEL, uint32(int32(x*wheelDelta)), 0, 0)
	}
	if msDelay > 0 {
		time.Sleep(time.Duration(msDelay) * time.Millisecond)
	}
}

// ScrollDir scrolls in a named direction: "up", "down", "left", "right".
func ScrollDir(x int, direction ...interface{}) {
	dir := "down"
	if len(direction) > 0 {
		if s, ok := direction[0].(string); ok {
			dir = s
		}
	}
	switch dir {
	case "up":
		Scroll(0, -x)
	case "down":
		Scroll(0, x)
	case "left":
		Scroll(-x, 0)
	case "right":
		Scroll(x, 0)
	}
}

// ScrollSmooth scrolls the mouse smoothly by `to` steps, repeating `num`
// times (default 5) with `tm` ms between steps (default 100). An optional
// third arg sets the horizontal offset per step.
func ScrollSmooth(to int, args ...int) {
	num := 5
	if len(args) > 0 {
		num = args[0]
	}
	tm := 100
	if len(args) > 1 {
		tm = args[1]
	}
	tox := 0
	if len(args) > 2 {
		tox = args[2]
	}

	for i := 0; i < num; i++ {
		Scroll(tox, to)
		MilliSleep(tm)
	}
	MilliSleep(MouseSleep)
}

// DragSmooth moves the mouse smoothly while holding a button down.
func DragSmooth(x, y int, args ...interface{}) {
	btn := "left"
	if len(args) > 0 {
		if s, ok := args[0].(string); ok {
			btn = s
		}
	}
	_ = Toggle(btn, "down")
	time.Sleep(50 * time.Millisecond)
	MoveSmooth(x, y)
	time.Sleep(50 * time.Millisecond)
	_ = Toggle(btn, "up")
}

// MoveClick moves to (x, y) then clicks.
func MoveClick(x, y int, args ...interface{}) {
	Move(x, y)
	_ = Click(args...)
}

// Location returns the current mouse position.
func Location() (int, int) {
	var p win.POINT
	if !win.GetCursorPos(&p) {
		return 0, 0
	}
	return int(p.X), int(p.Y)
}

// GetMousePos returns the current mouse position. Alias of Location.
func GetMousePos() (int, int) {
	return Location()
}

func mouseDelay() {
	if MouseSleep > 0 {
		time.Sleep(time.Duration(MouseSleep) * time.Millisecond)
	}
}
