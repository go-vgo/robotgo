// Copyright 2016-2017 The go-vgo Project Developers. See the COPYRIGHT
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
	"github.com/go-vgo/robotgo"
	// "go-vgo/robotgo"
)

func main() {
	////////////////////////////////////////////////////////////////////////////////
	// Control the keyboard
	////////////////////////////////////////////////////////////////////////////////
	robotgo.TypeString("Hello World") // Importing "Hello World"

	robotgo.KeyTap("enter") // Press "enter"
	robotgo.KeyTap("a", "control")
	robotgo.KeyTap("h", "command") // Hide window

	// Press "i", "alt", "command" Key combination
	robotgo.KeyTap("i", "alt", "command")
	arr := []string{"alt", "command"}
	robotgo.KeyTap("i", arr)

	robotgo.KeyTap("w", "command") // close window
	robotgo.KeyTap("m", "command") // minimize window
	robotgo.KeyTap("f1", "control")
	robotgo.KeyTap("a", "control")

	robotgo.KeyToggle("a", "down")
	robotgo.KeyToggle("a", "down", "alt")
	robotgo.KeyToggle("a", "down", "alt", "command")
	robotgo.KeyToggle("enter", "down")

	robotgo.TypeString("en")
}
