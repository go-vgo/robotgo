// Copyright 2016 The go-vgo Project Developers. See the COPYRIGHT
// file at the top-level directory of this distribution and at
// https://github.com/go-vgo/robotgo/blob/master/LICENSE
//
// Licensed under the Apache License, Version 2.0 <LICENSE-APACHE or
// http://www.apache.org/licenses/LICENSE-2.0> or the MIT license
// <LICENSE-MIT or http://opensource.org/licenses/MIT>, at your
// option. This file may not be copied, modified, or distributed
// except according to those terms.

package main

import (
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

func mouse() {
	////////////////////////////////////////////////////////////////////////////////
	// Control the mouse
	////////////////////////////////////////////////////////////////////////////////

	move()

	click()

	get()

	toggleAndScroll()
}

func main() {
	mouse()
}
