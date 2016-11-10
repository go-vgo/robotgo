// Copyright 2016 The go-vgo Project Developers. See the COPYRIGHT
// file at the top-level directory of this distribution and at
// http://www.
//
// Licensed under the Apache License, Version 2.0 <LICENSE-APACHE or
// http://www.apache.org/licenses/LICENSE-2.0> or the MIT license
// <LICENSE-MIT or http://opensource.org/licenses/MIT>, at your
// option. This file may not be copied, modified, or distributed
// except according to those terms.

package main

import (
	. "fmt"

	"github.com/go-vgo/robotgo"
)

func arobotgo() {
	x, y := robotgo.GetMousePos()
	Println("pos:", x, y)

	Println(robotgo.GetPixelColor(x, y))

	color := robotgo.GetPixelColor(100, 200)
	Println("color@@@", color)

	robotgo.TypeString("Hello World")
	// robotgo.KeyTap("a", "control")
	robotgo.KeyTap("f1", "control")
	// robotgo.KeyTap("enter")
	// robotgo.KeyToggle("enter", "down")
	robotgo.TypeString("en")

	bit_map := robotgo.CaptureScreen(10, 20, 30, 40)
	Println("...", bit_map)

	// robotgo.MouseClick()
	robotgo.ScrollMouse(10, "up")
}

func main() {
	arobotgo()
}
