//go:build linux
// +build linux

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

package libei

import "time"

// Linux evdev button codes (input-event-codes.h). The RemoteDesktop portal's
// NotifyPointerButton expects these evdev codes.
const (
	btnLeft   = 0x110 // BTN_LEFT
	btnRight  = 0x111 // BTN_RIGHT
	btnMiddle = 0x112 // BTN_MIDDLE
)

// MouseSleep is the global mouse delay in milliseconds.
var MouseSleep = 0

// Move moves the mouse to absolute position (x, y).
//
// Absolute positioning through the portal requires a linked ScreenCast stream
// (NotifyPointerMotionAbsolute). When no stream has been negotiated the portal
// rejects absolute motion, so this becomes a no-op. Use MoveRelative for
// reliable cursor movement on this backend.
func Move(x, y int, displayId ...int) {
	c, err := pointerReady()
	if err != nil || len(c.streams) == 0 {
		return
	}

	idx := 0
	if len(displayId) > 0 && displayId[0] >= 0 && displayId[0] < len(c.streams) {
		idx = displayId[0]
	}
	s := c.streams[idx]
	// Map global coordinates into the stream's local space.
	lx := float64(x - int(s.x))
	ly := float64(y - int(s.y))
	_ = c.inj.pointerMotionAbsolute(s.nodeID, lx, ly)
	mouseDelay()
}

// MoveRelative moves the mouse relative to its current position.
func MoveRelative(x, y int) {
	c, err := pointerReady()
	if err != nil {
		return
	}
	_ = c.inj.pointerMotion(float64(x), float64(y))
	mouseDelay()
}

// MoveSmooth moves the mouse smoothly using relative steps toward an offset
// (dx, dy). Because the portal does not expose the absolute cursor position,
// the arguments are treated as a relative delta and interpolated. Returns true
// on success.
func MoveSmooth(x, y int, args ...interface{}) bool {
	c, err := pointerReady()
	if err != nil {
		return false
	}

	steps := 20
	sleepMs := 5
	if len(args) >= 1 {
		if v, ok := args[0].(int); ok && v > 0 {
			steps = v
		}
	}
	if len(args) >= 2 {
		if v, ok := args[1].(int); ok {
			sleepMs = v
		}
	}

	var movedX, movedY int
	for i := 1; i <= steps; i++ {
		tx := x * i / steps
		ty := y * i / steps
		dx := tx - movedX
		dy := ty - movedY
		movedX, movedY = tx, ty
		if dx == 0 && dy == 0 {
			continue
		}
		if err := c.inj.pointerMotion(float64(dx), float64(dy)); err != nil {
			return false
		}
		if sleepMs > 0 {
			time.Sleep(time.Duration(sleepMs) * time.Millisecond)
		}
	}
	return true
}

// Click clicks a mouse button. Default is the left button. Pass a bool true to
// double-click.
func Click(args ...interface{}) error {
	c, err := pointerReady()
	if err != nil {
		return err
	}

	button := int32(btnLeft)
	double := false
	for _, arg := range args {
		switch v := arg.(type) {
		case string:
			button = resolveButton(v)
		case bool:
			double = v
		}
	}

	count := 1
	if double {
		count = 2
	}
	for i := 0; i < count; i++ {
		if err := c.inj.pointerButton(button, statePressed); err != nil {
			return err
		}
		time.Sleep(10 * time.Millisecond)
		if err := c.inj.pointerButton(button, stateReleased); err != nil {
			return err
		}
		if i < count-1 {
			time.Sleep(50 * time.Millisecond)
		}
	}
	mouseDelay()
	return nil
}

// Toggle toggles a mouse button down or up.
//
//	Toggle("left")        // press
//	Toggle("left", "up")  // release
func Toggle(key ...interface{}) error {
	c, err := pointerReady()
	if err != nil {
		return err
	}

	button := int32(btnLeft)
	state := statePressed
	for _, arg := range key {
		if v, ok := arg.(string); ok {
			switch v {
			case "up":
				state = stateReleased
			case "down":
				state = statePressed
			default:
				button = resolveButton(v)
			}
		}
	}
	return c.inj.pointerButton(button, state)
}

// MouseDown sends a mouse button down event.
func MouseDown(key ...interface{}) error {
	return Toggle(append(append([]interface{}{}, key...), "down")...)
}

// MouseUp sends a mouse button up event.
func MouseUp(key ...interface{}) error {
	return Toggle(append(append([]interface{}{}, key...), "up")...)
}

// Scroll scrolls the mouse. Positive y scrolls down, negative up; positive x
// scrolls right, negative left.
func Scroll(x, y int, args ...int) {
	c, err := pointerReady()
	if err != nil {
		return
	}

	msDelay := 10
	if len(args) > 0 {
		msDelay = args[0]
	}
	if y != 0 {
		_ = c.inj.pointerAxisDiscrete(axisVertical, int32(y))
	}
	if x != 0 {
		_ = c.inj.pointerAxisDiscrete(axisHorizontal, int32(x))
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

// ScrollSmooth scrolls the mouse smoothly by `to` steps, repeating `num` times
// (default 5) with `tm` ms between steps (default 100). An optional third arg
// sets the horizontal offset per step.
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
// The RemoteDesktop portal does not expose the global pointer position, so this
// returns (0, 0).
func Location() (int, int) { return 0, 0 }

// GetMousePos returns the current mouse position (alias of Location).
func GetMousePos() (int, int) { return Location() }

// pointerReady returns the connection if pointer injection is available.
func pointerReady() (*conn, error) {
	c, err := ensureConn()
	if err != nil {
		return nil, err
	}
	if c.inj == nil || !c.hasPointer() {
		return nil, ErrNotSupported
	}
	return c, nil
}

func resolveButton(btn string) int32 {
	switch btn {
	case "right":
		return btnRight
	case "center", "middle":
		return btnMiddle
	default:
		return btnLeft
	}
}

func mouseDelay() {
	if MouseSleep > 0 {
		time.Sleep(time.Duration(MouseSleep) * time.Millisecond)
	}
}
