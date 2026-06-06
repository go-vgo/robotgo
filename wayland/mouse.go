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
)

// MouseSleep is the global mouse delay in milliseconds.
var MouseSleep = 0

// Move moves the mouse to absolute position (x, y).
// An optional displayId selects which output to position against.
func Move(x, y int, displayId ...int) {
	c, err := ensureConn()
	if err != nil || c.pointer == nil {
		return
	}

	// Select the output to use for the absolute-positioning extent.
	idx := 0
	if len(displayId) > 0 && displayId[0] >= 0 && displayId[0] < len(c.outputs) {
		idx = displayId[0]
	}

	var xExtent, yExtent uint32 = 1920, 1080
	if idx < len(c.outputs) {
		o := c.outputs[idx]
		if o.width > 0 && o.height > 0 {
			xExtent = uint32(o.width)
			yExtent = uint32(o.height)
		}
	}

	_ = c.pointer.MotionAbsolute(timestamp(), clampExtent(x, xExtent), clampExtent(y, yExtent), xExtent, yExtent)
	_ = c.pointer.Frame()
	mouseDelay()
}

// MoveRelative moves the mouse relative to its current position.
func MoveRelative(x, y int) {
	c, err := ensureConn()
	if err != nil || c.pointer == nil {
		return
	}

	_ = c.pointer.Motion(timestamp(), float64(x), float64(y))
	_ = c.pointer.Frame()
	mouseDelay()
}

// MoveSmooth moves the mouse smoothly to (x, y) with a human-like curve.
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

	var xExtent, yExtent uint32 = 1920, 1080
	if len(c.outputs) > 0 {
		o := c.outputs[0]
		if o.width > 0 && o.height > 0 {
			xExtent = uint32(o.width)
			yExtent = uint32(o.height)
		}
	}

	// Smooth interpolation using ease-in-out
	for i := 1; i <= steps; i++ {
		t := float64(i) / float64(steps)
		// Ease-in-out cubic
		if t < 0.5 {
			t = 4 * t * t * t
		} else {
			t = 1 - math.Pow(-2*t+2, 3)/2
		}

		cx := clampExtent(int(float64(x)*t), xExtent)
		cy := clampExtent(int(float64(y)*t), yExtent)
		_ = c.pointer.MotionAbsolute(timestamp(), cx, cy, xExtent, yExtent)
		_ = c.pointer.Frame()
		time.Sleep(time.Duration(sleepMs) * time.Millisecond)
	}
	return true
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
		if err := c.pointer.Button(ts, uint32(button), buttonPressed); err != nil {
			return err
		}
		_ = c.pointer.Frame()
		time.Sleep(10 * time.Millisecond)
		if err := c.pointer.Button(ts+10, uint32(button), buttonReleased); err != nil {
			return err
		}
		_ = c.pointer.Frame()
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

	if err := c.pointer.Button(timestamp(), uint32(button), state); err != nil {
		return err
	}
	return c.pointer.Frame()
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

// Scroll scrolls the mouse. Positive y scrolls down, negative scrolls up.
func Scroll(x, y int, args ...int) {
	c, err := ensureConn()
	if err != nil || c.pointer == nil {
		return
	}

	msDelay := 10
	if len(args) > 0 {
		msDelay = args[0]
	}

	ts := timestamp()
	if y != 0 {
		_ = c.pointer.Axis(ts, axisVerticalScroll, float64(y)*15.0)
	}
	if x != 0 {
		_ = c.pointer.Axis(ts, axisHorizontalScroll, float64(x)*15.0)
	}
	_ = c.pointer.Frame()
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
// NOTE: Wayland does not expose global pointer position.
// This returns (0, 0) as a stub — position tracking would require
// a RemoteDesktop portal session.
func Location() (int, int) {
	return 0, 0
}

// GetMousePos returns the current mouse position.
// It is an alias of Location, mirroring the robotgo API.
// NOTE: Wayland does not expose global pointer position; returns (0, 0).
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
