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
	"fmt"

	"github.com/go-vgo/robotgo"
	// "go-vgo/robotgo"
)

func main() {
	////////////////////////////////////////////////////////////////////////////////
	// Control the keyboard
	////////////////////////////////////////////////////////////////////////////////

	// importing "Hello World"
	robotgo.TypeString("Hello World")

	// press "enter"
	robotgo.KeyTap("enter")
	robotgo.KeyTap("a", "control")
	// hide window
	robotgo.KeyTap("h", "command")

	// press "i", "alt", "command" Key combination
	robotgo.KeyTap("i", "alt", "command")
	arr := []string{"alt", "command"}
	robotgo.KeyTap("i", arr)

	// close window
	robotgo.KeyTap("w", "command")
	// minimize window
	robotgo.KeyTap("m", "command")
	robotgo.KeyTap("f1", "control")
	robotgo.KeyTap("a", "control")

	robotgo.KeyToggle("a", "down")
	robotgo.KeyToggle("a", "down", "alt")
	robotgo.KeyToggle("a", "down", "alt", "command")
	robotgo.KeyToggle("enter", "down")

	robotgo.TypeString("en")

	// write string to clipboard
	robotgo.WriteAll("测试")
	// read string from clipboard
	text, err := robotgo.ReadAll()
	if err == nil {
		fmt.Println(text)
	}
}
