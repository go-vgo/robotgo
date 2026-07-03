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

package x11

import (
	"math"
	"time"

	"github.com/jezek/xgb/xproto"
	"github.com/jezek/xgb/xtest"
)

// X11 pointer button numbers.
const (
	btnLeft       = 1
	btnMiddle     = 2
	btnRight      = 3
	btnWheelUp    = 4
	btnWheelDown  = 5
	btnWheelLeft  = 6
	btnWheelRight = 7
)

// MouseSleep is the global mouse delay in milliseconds, applied after a mouse
// action completes.
var MouseSleep = 0

// resolveButton maps a button name to its X11 button number.
func resolveButton(btn string) byte {
	switch btn {
	case "right":
		return btnRight
	case "center", "middle":
		return btnMiddle
	case "wheelUp":
		return btnWheelUp
	case "wheelDown":
		return btnWheelDown
	case "wheelLeft":
		return btnWheelLeft
	case "wheelRight":
		return btnWheelRight
	default:
		return btnLeft
	}
}

// sendButton presses or releases a pointer button via XTEST.
func (c *conn) sendButton(button byte, press bool) {
	t := byte(xproto.ButtonRelease)
	if press {
		t = byte(xproto.ButtonPress)
	}
	xtest.FakeInput(c.c, t, button, 0, c.root, 0, 0, 0)
}

// motion moves the pointer to absolute (x, y) via XTEST.
func (c *conn) motion(x, y int) {
	xtest.FakeInput(c.c, xproto.MotionNotify, 0, 0, c.root, int16(x), int16(y), 0)
	c.c.Sync()
}

// Move moves the mouse to absolute position (x, y).
func Move(x, y int, displayId ...int) {
	c, err := ensureConn()
	if err != nil {
		return
	}
	c.motion(x, y)
	mouseDelay()
}

// MoveRelative moves the mouse relative to its current position.
func MoveRelative(x, y int) {
	cx, cy := Location()
	Move(cx+x, cy+y)
}

// MoveSmooth moves the mouse smoothly to (x, y) with an ease-in-out curve.
// Optional args: MoveSmooth(x, y, steps int, sleepMs int). Returns true.
func MoveSmooth(x, y int, args ...interface{}) bool {
	c, err := ensureConn()
	if err != nil {
		return false
	}

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
		if t < 0.5 {
			t = 4 * t * t * t
		} else {
			t = 1 - math.Pow(-2*t+2, 3)/2
		}
		cx := sx + int(float64(x-sx)*t)
		cy := sy + int(float64(y-sy)*t)
		c.motion(cx, cy)
		time.Sleep(time.Duration(sleepMs) * time.Millisecond)
	}
	mouseDelay()
	return true
}

// Click clicks a mouse button. Default is the left button.
//
//	Click()             // left
//	Click("right")
//	Click("left", true) // double click
func Click(args ...interface{}) error {
	c, err := ensureConn()
	if err != nil {
		return err
	}

	button := byte(btnLeft)
	double := false
	for _, a := range args {
		switch v := a.(type) {
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
		c.sendButton(button, true)
		time.Sleep(10 * time.Millisecond)
		c.sendButton(button, false)
		c.c.Sync()
		if i < count-1 {
			time.Sleep(50 * time.Millisecond)
		}
	}
	mouseDelay()
	return nil
}

// Toggle toggles a mouse button down or up.
//
//	Toggle("left")        // down
//	Toggle("left", "up")
func Toggle(key ...interface{}) error {
	c, err := ensureConn()
	if err != nil {
		return err
	}

	button := byte(btnLeft)
	press := true
	for _, a := range key {
		if v, ok := a.(string); ok {
			switch v {
			case "up":
				press = false
			case "down":
				press = true
			default:
				button = resolveButton(v)
			}
		}
	}

	c.sendButton(button, press)
	c.c.Sync()
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

// Scroll scrolls the mouse. Positive y scrolls up, negative scrolls down;
// positive x scrolls left, negative scrolls right (matching robotgo's Cgo
// backend convention). args[0] is an optional inter-step delay in ms.
func Scroll(x, y int, args ...int) {
	c, err := ensureConn()
	if err != nil {
		return
	}

	msDelay := 10
	if len(args) > 0 {
		msDelay = args[0]
	}

	if y != 0 {
		btn := byte(btnWheelUp)
		n := y
		if y < 0 {
			btn, n = btnWheelDown, -y
		}
		c.wheel(btn, n)
	}
	if x != 0 {
		btn := byte(btnWheelLeft)
		n := x
		if x < 0 {
			btn, n = btnWheelRight, -x
		}
		c.wheel(btn, n)
	}
	c.c.Sync()
	if msDelay > 0 {
		time.Sleep(time.Duration(msDelay) * time.Millisecond)
	}
}

// wheel emits n press/release pairs of a scroll button.
func (c *conn) wheel(button byte, n int) {
	for i := 0; i < n; i++ {
		c.sendButton(button, true)
		c.sendButton(button, false)
	}
}

// ScrollDir scrolls in a named direction: "up", "down", "left", "right".
func ScrollDir(x int, direction ...interface{}) {
	d := "down"
	if len(direction) > 0 {
		if s, ok := direction[0].(string); ok {
			d = s
		}
	}
	switch d {
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

// ScrollSmooth scrolls smoothly by `to` steps, repeating `num` times (default
// 5) with `tm` ms between steps (default 100). An optional third arg sets the
// horizontal offset per step.
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

// DragSmooth moves the mouse smoothly to (x, y) while holding a button down.
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

// MoveClick moves the mouse to (x, y) then clicks.
func MoveClick(x, y int, args ...interface{}) {
	Move(x, y)
	MilliSleep(50)
	_ = Click(args...)
}

// Location returns the current mouse position.
func Location() (int, int) {
	c, err := ensureConn()
	if err != nil {
		return 0, 0
	}
	reply, err := xproto.QueryPointer(c.c, c.root).Reply()
	if err != nil || reply == nil {
		return 0, 0
	}
	return int(reply.RootX), int(reply.RootY)
}

// GetMousePos returns the current mouse position (alias of Location).
func GetMousePos() (int, int) {
	return Location()
}

func mouseDelay() {
	if MouseSleep > 0 {
		time.Sleep(time.Duration(MouseSleep) * time.Millisecond)
	}
}
