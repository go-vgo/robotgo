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
	// typing "Hello World"
	robotgo.TypeStr("Hello World!", 0, 1)
	robotgo.KeySleep = 100
	robotgo.TypeStr("だんしゃり")

	robotgo.TypeStr("Hi galaxy, hi stars, hi MT.Rainier, hi sea. こんにちは世界.")
	robotgo.TypeStr("So, hi, bye! 你好, 再见!")
	robotgo.Sleep(1)

	robotgo.TypeStr("Hi, Seattle space needle, Golden gate bridge, One world trade center.")
	robotgo.MilliSleep(100)

	ustr := uint32(robotgo.CharCodeAt("So, hi, bye!", 0))
	robotgo.UnicodeType(ustr)

	err := robotgo.PasteStr("paste string")
	fmt.Println("PasteStr: ", err)
}

func keyTap() {
	// press "enter"
	robotgo.KeyTap("enter")
	robotgo.KeyTap(robotgo.Enter)
	robotgo.KeySleep = 200
	robotgo.KeyTap("a")
	robotgo.MilliSleep(100)
	robotgo.KeyTap("a", "ctrl")

	// hide window
	err := robotgo.KeyTap("h", "cmd")
	if err != nil {
		fmt.Println("robotgo.KeyTap run error is: ", err)
	}

	robotgo.KeyTap("h", "cmd")

	// press "i", "alt", "command" Key combination
	robotgo.KeyTap(robotgo.KeyI, robotgo.Alt, robotgo.Cmd)
	robotgo.KeyTap("i", "alt", "cmd")

	arr := []string{"alt", "cmd"}
	robotgo.KeyTap("i", arr)
	robotgo.KeyTap("i", arr)

	robotgo.KeyTap("i", "cmd", " alt", "shift")

	// close window
	robotgo.KeyTap("w", "cmd")

	// minimize window
	robotgo.KeyTap("m", "cmd")

	robotgo.KeyTap("f1", "ctrl")
	robotgo.KeyTap("a", "control")
}

func special() {
	robotgo.TypeStr("{}")
	robotgo.KeyTap("[", "]")

	robotgo.KeyToggle("(")
	robotgo.KeyToggle("(", "up")
}

func keyToggle() {
	// robotgo.KeySleep = 150
	robotgo.KeyToggle(robotgo.KeyA)
	robotgo.KeyToggle("a", "down", "alt")
	robotgo.Sleep(1)

	robotgo.KeyToggle("a", "up", "alt", "cmd")
	robotgo.MilliSleep(100)
	robotgo.KeyToggle("q", "up", "alt", "cmd", "shift")

	err := robotgo.KeyToggle(robotgo.Enter)
	if err != nil {
		fmt.Println("robotgo.KeyToggle run error is: ", err)
	}
}

func cilp() {
	// robotgo.TypeStr("en")

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
	fmt.Println("text: ", text)
}

func key() {
	////////////////////////////////////////////////////////////////////////////////
	// Control the keyboard
	////////////////////////////////////////////////////////////////////////////////

	typeStr()
	special()

	keyTap()
	keyToggle()

	cilp()
}

func main() {
	key()
}
