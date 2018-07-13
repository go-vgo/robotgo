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

func key() {
	////////////////////////////////////////////////////////////////////////////////
	// Control the keyboard
	////////////////////////////////////////////////////////////////////////////////

	// importing "Hello World"
	robotgo.TypeString("Hello World")

	robotgo.TypeString("留给真爱你的人")
	robotgo.MicroSleep(1)

	robotgo.TypeStr("所以, 你好, 再见")
	robotgo.Sleep(1)

	ustr := uint32(robotgo.CharCodeAt("所以, 你好, 再见", 0))
	robotgo.UnicodeType(ustr)

	robotgo.PasteStr(" 粘贴字符串, paste")

	// press "enter"
	robotgo.KeyTap("enter")
	robotgo.KeyTap("a", "control")
	// hide window
	robotgo.KeyTap("h", "command")
	robotgo.KeyTap("h", "command", 12)

	// press "i", "alt", "command" Key combination
	robotgo.KeyTap("i", "alt", "command")
	robotgo.KeyTap("i", "alt", "command", 11)
	arr := []string{"alt", "command"}
	robotgo.KeyTap("i", arr)
	robotgo.KeyTap("i", arr, 12)

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

func main() {
	key()
}
