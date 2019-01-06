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

func typeStr() {
	// importing "Hello World"
	robotgo.TypeStr("Hello World")

	robotgo.TypeStr("留给真爱你的人")
	robotgo.MicroSleep(10.2)

	robotgo.TypeStr("山达尔星新星军团, galaxy. こんにちは世界.")
	robotgo.Sleep(2)
	robotgo.TypeString("所以, 你好, 再见")
	robotgo.MilliSleep(100)

	ustr := uint32(robotgo.CharCodeAt("所以, 你好, 再见", 0))
	robotgo.UnicodeType(ustr)

	robotgo.PasteStr(" 粘贴字符串, paste")
}

func keyTap() {
	// press "enter"
	robotgo.KeyTap("enter")
	robotgo.KeyTap("a")
	robotgo.KeyTap("a", "ctrl")

	// hide window
	err := robotgo.KeyTap("h", "cmd")
	if err == "" {
		fmt.Println("robotgo.KeyTap run error is nil.")
	}

	robotgo.KeyTap("h", "cmd", 12)

	// press "i", "alt", "command" Key combination
	robotgo.KeyTap("i", "alt", "command")
	robotgo.KeyTap("i", "alt", "command", 11)

	arr := []string{"alt", "cmd"}
	robotgo.KeyTap("i", arr)
	robotgo.KeyTap("i", arr, 12)

	// close window
	robotgo.KeyTap("w", "cmd")

	// minimize window
	robotgo.KeyTap("m", "cmd")

	robotgo.KeyTap("f1", "ctrl")
	robotgo.KeyTap("a", "control")
}

func keyToggle() {
	robotgo.KeyToggle("a", "down")
	robotgo.KeyToggle("a", "down", "alt")
	robotgo.KeyToggle("a", "down", "alt", "cmd")

	err := robotgo.KeyToggle("enter", "down")
	if err == "" {
		fmt.Println("robotgo.KeyToggle run error is nil.")
	}
}

func cilp() {
	robotgo.TypeString("en")

	// write string to clipboard
	robotgo.WriteAll("测试")
	// read string from clipboard
	text, err := robotgo.ReadAll()
	if err == nil {
		fmt.Println(text)
	}
}

func key() {
	////////////////////////////////////////////////////////////////////////////////
	// Control the keyboard
	////////////////////////////////////////////////////////////////////////////////

	typeStr()

	keyTap()
	keyToggle()

	cilp()
}

func main() {
	key()
}
