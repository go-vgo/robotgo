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
	robotgo.TypeStr("Hello World!", 1.0)

	robotgo.TypeStr("だんしゃり")
	robotgo.MicroSleep(10.2)

	robotgo.TypeStr("Hi galaxy. こんにちは世界.")
	robotgo.Sleep(2)
	// robotgo.TypeString("So, hi, bye!")
	robotgo.MilliSleep(100)

	ustr := uint32(robotgo.CharCodeAt("So, hi, bye!", 0))
	robotgo.UnicodeType(ustr)

	robotgo.PasteStr("paste string")
}

func keyTap() {
	// press "enter"
	robotgo.KeyTap("enter")
	robotgo.KeyTap("a")
	robotgo.MilliSleep(100)
	robotgo.KeyTap("a", "ctrl")

	// hide window
	err := robotgo.KeyTap("h", "cmd")
	if err != "" {
		fmt.Println("robotgo.KeyTap run error is: ", err)
	}

	robotgo.KeyTap("h", "cmd", 12)

	// press "i", "alt", "command" Key combination
	robotgo.KeyTap("i", "alt", "command")
	robotgo.KeyTap("i", "alt", "cmd", 11)

	arr := []string{"alt", "cmd"}
	robotgo.KeyTap("i", arr)
	robotgo.KeyTap("i", arr, 12)

	robotgo.KeyTap("i", "cmd", " alt", "shift")

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
	robotgo.Sleep(1)

	robotgo.KeyToggle("a", "up", "alt", "cmd")
	robotgo.MilliSleep(100)
	robotgo.KeyToggle("q", "up", "alt", "cmd", "shift")

	err := robotgo.KeyToggle("enter", "down")
	if err != "" {
		fmt.Println("robotgo.KeyToggle run error is: ", err)
	}
}

func cilp() {
	// robotgo.TypeString("en")

	// write string to clipboard
	e := robotgo.WriteAll("テストする")
	if e != nil {
		fmt.Println("robotgo.WriteAll err is: ", e)
	}

	// read string from clipboard
	text, err := robotgo.ReadAll()
	if err != nil {
		fmt.Println("robotgo.ReadAll err is: ", err)
	}
	fmt.Println(text)
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
