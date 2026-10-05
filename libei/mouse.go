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

import (
	"fmt"
	"time"
)

// Linux evdev button codes (input-event-codes.h). The RemoteDesktop portal's
// NotifyPointerButton expects these evdev codes.
const (
	btnLeft   = 0x110 // BTN_LEFT
	btnRight  = 0x111 // BTN_RIGHT
	btnMiddle = 0x112 // BTN_MIDDLE
)

// MouseSleep is the global mouse delay in milliseconds.
var MouseSleep = 0

// cornerReset is the relative delta used to park the pointer in the top-left
// corner when its position is unknown; larger than any screen.
const cornerReset = 1 << 16

// Move moves the mouse to absolute position (x, y).
//
// With a linked ScreenCast stream (see LinkScreenCast) the move is sent as
// NotifyPointerMotionAbsolute on the stream that contains (x, y) — or the one
// selected by displayId. Without a stream the portal only accepts relative
// motion, so Move sends a delta from the last tracked position; when no
// position has been tracked yet the pointer is first parked in the top-left
// corner with a large relative move so the origin is known.
func Move(x, y int, displayId ...int) error {
	_, err := move(x, y, displayId...)
	return err
}

// move implements Move and reports whether the pointer was injected exactly at
// (x, y): false when injection fails or the target was clamped into a stream.
func move(x, y int, displayId ...int) (bool, error) {
	c, err := pointerReady()
	if err != nil {
		return false, err
	}

	if len(c.streams) == 0 {
		cx, cy, ok := c.position()
		if !ok {
			if err := c.inj.pointerMotion(-cornerReset, -cornerReset); err != nil {
				return false, err
			}
			c.setPos(0, 0)
			cx, cy = 0, 0
		}
		err := c.inj.pointerMotion(float64(x-cx), float64(y-cy))
		if err == nil {
			c.setPos(x, y)
		}
		mouseDelay()
		return err == nil, err
	}

	idx, found := c.streamAt(x, y)
	if len(displayId) > 0 && displayId[0] >= 0 && displayId[0] < len(c.streams) {
		idx, found = displayId[0], true
	}
	s := c.streams[idx]
	// Map global coordinates into the stream's local space; a point outside
	// every stream (a gap between monitors) is clamped onto the chosen one
	// rather than sent as an out-of-range position the portal rejects.
	lx, ly := x-int(s.x), y-int(s.y)
	if !found {
		lx = clamp(lx, 0, int(s.width)-1)
		ly = clamp(ly, 0, int(s.height)-1)
	}
	err = c.inj.pointerMotionAbsolute(s.nodeID, float64(lx), float64(ly))
	if err == nil {
		c.setPos(lx+int(s.x), ly+int(s.y))
	}
	mouseDelay()
	return err == nil && found, err
}

// MoveRelative moves the mouse relative to its current position.
func MoveRelative(x, y int) error {
	c, err := pointerReady()
	if err != nil {
		return err
	}
	err = c.inj.pointerMotion(float64(x), float64(y))
	if err == nil {
		c.addPos(x, y)
	}
	mouseDelay()
	return err
}

// MoveSmooth moves the mouse smoothly to absolute position (x, y). Optional
// args: steps (default 20), sleep ms between steps (default 5). Returns true
// on success; false when a step fails or the target was clamped into a
// linked stream.
//
// When the current position is unknown, it jumps via Move, using a corner
// reset first if no stream is linked.
func MoveSmooth(x, y int, args ...interface{}) bool {
	reached, err := moveSmooth(x, y, args...)
	return err == nil && reached
}

// moveSmooth implements MoveSmooth, stopping at the first failed step. It
// reports whether the pointer ended exactly at (x, y).
func moveSmooth(x, y int, args ...interface{}) (bool, error) {
	c, err := pointerReady()
	if err != nil {
		return false, err
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

	sx, sy, known := c.position()
	if !known {
		return move(x, y)
	}

	for i := 1; i <= steps; i++ {
		tx := sx + (x-sx)*i/steps
		ty := sy + (y-sy)*i/steps
		cx, cy, _ := c.position()
		if tx == cx && ty == cy {
			continue
		}
		if len(c.streams) > 0 {
			err = Move(tx, ty)
		} else {
			err = MoveRelative(tx-cx, ty-cy)
		}
		if err != nil {
			return false, err
		}
		if sleepMs > 0 {
			time.Sleep(time.Duration(sleepMs) * time.Millisecond)
		}
	}
	cx, cy, _ := c.position()
	return cx == x && cy == y, nil
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
// Scroll scrolls the mouse by wheel notches. Positive y scrolls up, negative
// scrolls down; positive x scrolls left, negative scrolls right (matching
// robotgo's Cgo backend convention). Optional arg: delay ms.
func Scroll(x, y int, args ...int) error {
	c, err := pointerReady()
	if err != nil {
		return err
	}

	msDelay := 10
	if len(args) > 0 {
		msDelay = args[0]
	}
	// The portal counts positive steps as down/right.
	if y != 0 {
		if err := c.inj.pointerAxisDiscrete(axisVertical, int32(-y)); err != nil {
			return err
		}
	}
	if x != 0 {
		if err := c.inj.pointerAxisDiscrete(axisHorizontal, int32(-x)); err != nil {
			return err
		}
	}
	if msDelay > 0 {
		time.Sleep(time.Duration(msDelay) * time.Millisecond)
	}
	return nil
}

// ScrollDir scrolls in a named direction: "up", "down" (default), "left",
// "right". Any other direction is an error.
func ScrollDir(x int, direction ...interface{}) error {
	dir := "down"
	if len(direction) > 0 {
		s, ok := direction[0].(string)
		if !ok {
			return fmt.Errorf("robotgo: unknown scroll direction: %v", direction[0])
		}
		dir = s
	}
	switch dir {
	case "down":
		return Scroll(0, -x)
	case "up":
		return Scroll(0, x)
	case "left":
		return Scroll(x, 0)
	case "right":
		return Scroll(-x, 0)
	}
	return fmt.Errorf("robotgo: unknown scroll direction: %v", dir)
}

// ScrollSmooth scrolls the mouse smoothly by `to` steps, repeating `num` times
// (default 5) with `tm` ms between steps (default 100). An optional third arg
// sets the horizontal offset per step. It stops at the first failed scroll.
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
		MilliSleep(tm)
	}
	MilliSleep(MouseSleep)
	return nil
}

// DragSmooth moves the mouse smoothly while holding a button down. The button
// is released even if the move fails; the first error of the press, move and
// release is returned. A target clamped into a linked stream is not an error.
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
	_, err := moveSmooth(x, y)
	time.Sleep(50 * time.Millisecond)
	if upErr := Toggle(btn, "up"); err == nil {
		err = upErr
	}
	return err
}

// MoveClick moves to (x, y) then clicks. It does not click if the move fails.
func MoveClick(x, y int, args ...interface{}) error {
	if err := Move(x, y); err != nil {
		return err
	}
	return Click(args...)
}

// Location returns the current mouse position.
//
// The RemoteDesktop portal does not expose the real cursor position, so this
// returns the last position injected by this backend (Move, MoveRelative,
// MoveSmooth, ...). It is (0, 0) until the first move and does not follow
// movement made by the physical mouse. Relative motion alone cannot establish
// an absolute position; Location stays (0, 0) until Move establishes one.
//
// It is a pure query: it never opens a portal session (and so never shows the
// consent dialog), and it keeps answering from the last session after the
// portal closed it (#783).
func Location() (int, int) {
	connMu.Lock()
	c := globalConn
	connMu.Unlock()
	if c == nil {
		return 0, 0
	}
	x, y, _ := c.position()
	return x, y
}

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
