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

func mouse() {
	////////////////////////////////////////////////////////////////////////////////
	// Control the mouse
	////////////////////////////////////////////////////////////////////////////////

	// move the mouse to 100, 200
	robotgo.MoveMouse(100, 200)

	// click the left mouse button
	robotgo.MouseClick()
	// click the right mouse button
	robotgo.MouseClick("right", false)
	// double click the left mouse button
	robotgo.MouseClick("left", true)

	// scrolls the mouse either up
	robotgo.ScrollMouse(10, "up")
	robotgo.Scroll(100, 200)
	// toggles right mouse button
	robotgo.MouseToggle("down", "right")

	// smooth move the mouse to 100, 200
	robotgo.MoveMouseSmooth(100, 200)
	robotgo.MoveMouseSmooth(100, 200, 1.0, 100.0)

	// gets the mouse coordinates
	x, y := robotgo.GetMousePos()
	fmt.Println("pos:", x, y)
	if x == 456 && y == 586 {
		fmt.Println("mouse...", "586")
	}

	robotgo.MouseToggle("up")
	robotgo.MoveMouse(x, y)
	robotgo.MoveMouse(100, 200)

	for i := 0; i < 1080; i += 1000 {
		fmt.Println(i)
		robotgo.MoveMouse(800, i)
	}

}

func main() {
	mouse()
}
