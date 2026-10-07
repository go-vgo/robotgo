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

package main

import (
	"errors"
	"fmt"

	"github.com/go-vgo/robotgo"
	// "go-vgo/robotgo"
)

func move() {
	robotgo.MouseSleep = 100
	robotgo.Move(100, 200)
	robotgo.MoveRelative(10, -200)

	// move the mouse to 100, 200
	robotgo.Move(100, 200)

	// drag mouse with smooth
	robotgo.DragSmooth(10, 10)
	robotgo.DragSmooth(100, 200, 1.0, 100.0)

	// smooth move the mouse to 100, 200
	robotgo.MoveSmooth(100, 200)
	robotgo.MoveSmooth(100, 200, 1.0, 100.0)
	robotgo.MoveSmoothRelative(10, -100, 1.0, 30.0)

	for i := 0; i < 1080; i += 1000 {
		fmt.Println("i: ", i)
		// MoveMouse(800, i)
		robotgo.Move(800, i)
	}
}

func click() {

	// click the left mouse button
	robotgo.Click()

	// click the right mouse button
	robotgo.Click("right", false)

	// double click the left mouse button
	robotgo.Click("left", true)
}

func get() {
	// gets the mouse coordinates
	x, y := robotgo.Location()
	fmt.Println("pos:", x, y)
	if x == 456 && y == 586 {
		fmt.Println("mouse...", "586")
	}

	robotgo.Move(x, y)
}

func toggleAndScroll() {
	// scrolls the mouse either up
	robotgo.ScrollDir(10, "up")
	robotgo.ScrollDir(10, "right")

	robotgo.Scroll(100, 10)
	robotgo.Scroll(0, -10)

	robotgo.Toggle("left")
	robotgo.Toggle("left", "up")

	// toggles the right mouse button
	robotgo.Toggle("right")
	robotgo.Toggle("right", "up")
}

// checkErr prints an error, telling unsupported operations apart
func checkErr(name string, err error) {
	switch {
	case err == nil:
	case errors.Is(err, robotgo.ErrNotSupported):
		fmt.Println(name, "is not supported by this backend")
	default:
		fmt.Println(name, "error:", err)
	}
}

func mouseErr() {
	// the mouse APIs return an error, check it instead of ignoring it
	checkErr("robotgo.Move", robotgo.Move(100, 200))

	// invalid arguments are reported as errors
	checkErr("robotgo.ScrollDir", robotgo.ScrollDir(10, "forward"))

	// backends differ: Cgo returns an error (second arg must be a bool),
	// pure-Go returns nil and does a single left click ("double" is read
	// as an unknown button name, which falls back to "left")
	checkErr("robotgo.Click", robotgo.Click("left", "double"))

	// MoveSmooth returns a bool, MoveSmoothRelative wraps it as ErrSmoothMove
	err := robotgo.MoveSmoothRelative(10, -10)
	if errors.Is(err, robotgo.ErrSmoothMove) {
		fmt.Println("smooth move failed:", err)
	}

	// release the button even when the drag fails
	if err := robotgo.MouseDown("left"); err != nil {
		checkErr("robotgo.MouseDown", err)
		return
	}
	defer func() {
		checkErr("robotgo.MouseUp", robotgo.MouseUp("left"))
	}()
	checkErr("robotgo.MoveRelative", robotgo.MoveRelative(20, 20))
}

func mouse() {
	////////////////////////////////////////////////////////////////////////////////
	// Control the mouse
	////////////////////////////////////////////////////////////////////////////////

	move()

	click()

	get()

	toggleAndScroll()

	mouseErr()
}

func main() {
	mouse()
}
