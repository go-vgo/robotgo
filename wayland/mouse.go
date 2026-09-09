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

package wayland

import (
	"math"
	"sync"
	"time"
)

// Linux evdev button codes
const (
	btnLeft   = 0x110 // BTN_LEFT
	btnRight  = 0x111 // BTN_RIGHT
	btnMiddle = 0x112 // BTN_MIDDLE
)

// Wayland pointer button states
const (
	buttonReleased = 0
	buttonPressed  = 1
)

// Wayland pointer axis types
const (
	axisVerticalScroll   = 0
	axisHorizontalScroll = 1

	// WL_POINTER_AXIS_SOURCE_WHEEL; one wheel notch is 15 axis units.
	axisSourceWheel = 0
	wheelStep       = 15.0
)

// MouseSleep is the global mouse delay in milliseconds.
var MouseSleep = 0

// Last pointer position injected by this backend; see Location.
var (
	posMu      sync.Mutex
	posX, posY int
)

func setPos(x, y int) {
	posMu.Lock()
	posX, posY = x, y
	posMu.Unlock()
}

// warp sends one absolute motion + frame. motion_absolute is normalized
// against the whole output layout (that is what wlroots maps a virtual
// pointer onto), so the extent is the union of all outputs and the
// coordinates are made relative to its origin.
func (c *conn) warp(x, y int) {
	ox, oy, ow, oh := c.outputBounds()
	if ow <= 0 || oh <= 0 {
		ox, oy, ow, oh = 0, 0, 1920, 1080
	}
	cx := clampExtent(x-int(ox), uint32(ow))
	cy := clampExtent(y-int(oy), uint32(oh))
	err := c.do(func() error {
		if err := c.pointer.MotionAbsolute(timestamp(), cx, cy, uint32(ow), uint32(oh)); err != nil {
			return err
		}
		return c.pointer.Frame()
	})
	if err == nil {
		setPos(int(cx)+int(ox), int(cy)+int(oy))
	}
}

// Move moves the mouse to absolute position (x, y) in layout coordinates.
// The optional displayId is accepted for API parity; outputs share one
// layout on Wayland so it does not change the target.
func Move(x, y int, displayId ...int) {
	c, err := ensureConn()
	if err != nil || c.pointer == nil {
		return
	}
	c.warp(x, y)
	mouseDelay()
}

// MoveRelative moves the mouse relative to its current position.
func MoveRelative(x, y int) {
	c, err := ensureConn()
	if err != nil || c.pointer == nil {
		return
	}

	err = c.do(func() error {
		if err := c.pointer.Motion(timestamp(), float64(x), float64(y)); err != nil {
			return err
		}
		return c.pointer.Frame()
	})
	if err == nil {
		posMu.Lock()
		posX += x
		posY += y
		posMu.Unlock()
	}
	mouseDelay()
}

// MoveSmooth moves the mouse smoothly from the last injected position to
// (x, y) with an ease-in-out curve. Optional args: steps (int), sleepMs (int).
// Returns true on success.
func MoveSmooth(x, y int, args ...interface{}) bool {
	c, err := ensureConn()
	if err != nil || c.pointer == nil {
		return false
	}

	// Default parameters
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
	// Smooth interpolation using ease-in-out
	for i := 1; i <= steps; i++ {
		t := float64(i) / float64(steps)
		// Ease-in-out cubic
		if t < 0.5 {
			t = 4 * t * t * t
		} else {
			t = 1 - math.Pow(-2*t+2, 3)/2
		}

		cx := float64(sx) + float64(x-sx)*t
		cy := float64(sy) + float64(y-sy)*t
		c.warp(int(math.Round(cx)), int(math.Round(cy)))
		time.Sleep(time.Duration(sleepMs) * time.Millisecond)
	}
	mouseDelay()
	return true
}

// button sends one button event + frame under the connection lock.
func (c *conn) button(ts uint32, button int, state uint32) error {
	return c.do(func() error {
		if err := c.pointer.Button(ts, uint32(button), state); err != nil {
			return err
		}
		return c.pointer.Frame()
	})
}

// Click clicks a mouse button. Default is left button.
// Use "left", "right", or "center".
func Click(args ...interface{}) error {
	c, err := ensureConn()
	if err != nil {
		return err
	}
	if c.pointer == nil {
		return ErrNotSupported
	}

	button := btnLeft
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
		ts := timestamp()
		if err := c.button(ts, button, buttonPressed); err != nil {
			return err
		}
		time.Sleep(10 * time.Millisecond)
		if err := c.button(ts+10, button, buttonReleased); err != nil {
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
// Toggle("left") or Toggle("left", "up")
func Toggle(key ...interface{}) error {
	c, err := ensureConn()
	if err != nil {
		return err
	}
	if c.pointer == nil {
		return ErrNotSupported
	}

	button := btnLeft
	state := uint32(buttonPressed)

	for _, arg := range key {
		switch v := arg.(type) {
		case string:
			switch v {
			case "up":
				state = buttonReleased
			case "down":
				state = buttonPressed
			default:
				button = resolveButton(v)
			}
		}
	}

	return c.button(timestamp(), button, state)
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

// Scroll scrolls the mouse by wheel notches. Positive y scrolls up, negative
// scrolls down; positive x scrolls left, negative scrolls right (matching
// robotgo's Cgo backend convention). Optional arg: delay ms.
func Scroll(x, y int, args ...int) {
	c, err := ensureConn()
	if err != nil || c.pointer == nil {
		return
	}

	msDelay := 10
	if len(args) > 0 {
		msDelay = args[0]
	}

	// wl_pointer axis values are positive towards down/right, so robotgo's
	// up/left-positive notches are negated. axis_source must precede the
	// axis events and axis_discrete carries the notch count that clients
	// which only handle discrete (wheel) scrolling rely on.
	ts := timestamp()
	_ = c.do(func() error {
		if err := c.pointer.AxisSource(axisSourceWheel); err != nil {
			return err
		}
		if y != 0 {
			if err := c.pointer.AxisDiscrete(ts, axisVerticalScroll, float64(-y)*wheelStep, int32(-y)); err != nil {
				return err
			}
		}
		if x != 0 {
			if err := c.pointer.AxisDiscrete(ts, axisHorizontalScroll, float64(-x)*wheelStep, int32(-x)); err != nil {
				return err
			}
		}
		return c.pointer.Frame()
	})
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
	case "down":
		Scroll(0, -x)
	case "up":
		Scroll(0, x)
	case "left":
		Scroll(x, 0)
	case "right":
		Scroll(-x, 0)
	}
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
// NOTE: Wayland does not expose the global pointer position, so this returns
// the last position injected by this backend (Move, MoveRelative, MoveSmooth).
// It is (0, 0) until the first move and does not follow the physical mouse.
func Location() (int, int) {
	posMu.Lock()
	defer posMu.Unlock()
	return posX, posY
}

// GetMousePos returns the current mouse position.
// It is an alias of Location, mirroring the robotgo API.
func GetMousePos() (int, int) {
	return Location()
}

// ScrollSmooth scrolls the mouse smoothly by `to` steps, repeating `num`
// times (default 5) with `tm` ms between steps (default 100). An optional
// third arg sets the horizontal offset per step.
//
//	robotgo.ScrollSmooth(10)
//	robotgo.ScrollSmooth(10, 6, 50, 1)
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

// clampExtent clamps a coordinate to the valid [0, extent] range and returns it
// as a uint32, avoiding the wraparound that a direct uint32(negative) conversion
// would cause for off-screen / negative inputs.
func clampExtent(v int, extent uint32) uint32 {
	if v < 0 {
		return 0
	}
	if uint32(v) > extent {
		return extent
	}
	return uint32(v)
}

func resolveButton(btn string) int {
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
