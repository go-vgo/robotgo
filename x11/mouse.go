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
	"fmt"
	"math"
	"time"

	"github.com/go-vgo/robotgo/pub"
	"github.com/jezek/xgb/xproto"
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
func (c *conn) sendButton(button byte, press bool) error {
	t := byte(xproto.ButtonRelease)
	if press {
		t = byte(xproto.ButtonPress)
	}
	return c.fakeInput(t, button, 0, 0)
}

// motion moves the pointer to absolute (x, y) via XTEST.
func (c *conn) motion(x, y int) error {
	return c.fakeInput(xproto.MotionNotify, 0, int16(x), int16(y))
}

// Move moves the mouse to absolute position (x, y).
func Move(x, y int, displayId ...int) error {
	c, err := ensureConn()
	if err != nil {
		return err
	}
	err = c.motion(x, y)
	mouseDelay()
	return err
}

// MoveRelative moves the mouse relative to its current position.
func MoveRelative(x, y int) error {
	cx, cy := Location()
	return Move(cx+x, cy+y)
}

// MoveSmooth moves the mouse smoothly to (x, y) with an ease-in-out curve.
// Optional args: MoveSmooth(x, y, steps int, sleepMs int). Returns false if a
// step fails.
func MoveSmooth(x, y int, args ...interface{}) bool {
	return moveSmooth(x, y, args...) == nil
}

// moveSmooth implements MoveSmooth, stopping at the first failed step.
func moveSmooth(x, y int, args ...interface{}) error {
	c, err := ensureConn()
	if err != nil {
		return err
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
		if err := c.motion(cx, cy); err != nil {
			return err
		}
		pub.MilliSleep(sleepMs)
	}
	mouseDelay()
	return nil
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
		if err := c.sendButton(button, true); err != nil {
			return err
		}
		time.Sleep(10 * time.Millisecond)
		if err := c.sendButton(button, false); err != nil {
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

	return c.sendButton(button, press)
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
func Scroll(x, y int, args ...int) error {
	c, err := ensureConn()
	if err != nil {
		return err
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
		if err := c.wheel(btn, n); err != nil {
			return err
		}
	}
	if x != 0 {
		btn := byte(btnWheelLeft)
		n := x
		if x < 0 {
			btn, n = btnWheelRight, -x
		}
		if err := c.wheel(btn, n); err != nil {
			return err
		}
	}
	pub.MilliSleep(msDelay)
	return nil
}

// wheel emits n press/release pairs of a scroll button.
func (c *conn) wheel(button byte, n int) error {
	for i := 0; i < n; i++ {
		if err := c.sendButton(button, true); err != nil {
			return err
		}
		if err := c.sendButton(button, false); err != nil {
			return err
		}
	}
	return nil
}

// ScrollDir scrolls in a named direction: "up", "down" (default), "left",
// "right". Any other direction is an error.
func ScrollDir(x int, direction ...interface{}) error {
	d := "down"
	if len(direction) > 0 {
		s, ok := direction[0].(string)
		if !ok {
			return fmt.Errorf("robotgo: unknown scroll direction: %v", direction[0])
		}
		d = s
	}
	switch d {
	case "down":
		return Scroll(0, -x)
	case "up":
		return Scroll(0, x)
	case "left":
		return Scroll(x, 0)
	case "right":
		return Scroll(-x, 0)
	}
	return fmt.Errorf("robotgo: unknown scroll direction: %v", d)
}

// ScrollSmooth scrolls smoothly by `to` steps, repeating `num` times (default
// 5) with `tm` ms between steps (default 100). An optional third arg sets the
// horizontal offset per step. It stops at the first failed scroll.
func ScrollSmooth(to int, args ...int) error {
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
		if err := Scroll(tox, to); err != nil {
			return err
		}
		pub.MilliSleep(tm)
	}
	pub.MilliSleep(pub.MouseSleep)
	return nil
}

// DragSmooth moves the mouse smoothly to (x, y) while holding a button down.
// The button is released even if the move fails; the first error of the
// press, move and release is returned.
func DragSmooth(x, y int, args ...interface{}) error {
	btn := "left"
	if len(args) > 0 {
		if s, ok := args[0].(string); ok {
			btn = s
		}
	}
	if err := Toggle(btn, "down"); err != nil {
		return err
	}
	time.Sleep(50 * time.Millisecond)
	err := moveSmooth(x, y)
	time.Sleep(50 * time.Millisecond)
	if upErr := Toggle(btn, "up"); err == nil {
		err = upErr
	}
	return err
}

// MoveClick moves the mouse to (x, y) then clicks. It does not click if the
// move fails.
func MoveClick(x, y int, args ...interface{}) error {
	if err := Move(x, y); err != nil {
		return err
	}
	pub.MilliSleep(50)
	return Click(args...)
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
	pub.MilliSleep(pub.MouseSleep)
}
