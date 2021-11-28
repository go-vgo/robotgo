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
	"reflect"

	"github.com/go-vgo/robotgo"
)

func aRobotgo() {
	abool := robotgo.ShowAlert("test", "robotgo")
	if abool {
		fmt.Println("ok@@@", "ok")
	}

	x, y := robotgo.GetMousePos()
	fmt.Println("pos:", x, y)

	robotgo.Move(x, y)
	robotgo.Move(100, 200)

	robotgo.Toggle("left")
	robotgo.Toggle("left", "up")

	for i := 0; i < 1080; i += 1000 {
		fmt.Println(i)
		robotgo.Move(800, i)
	}

	fmt.Println(robotgo.GetPixelColor(x, y))

	color := robotgo.GetPixelColor(100, 200)
	fmt.Println("color@@@", color)

	robotgo.TypeStr("Hello World")
	// robotgo.KeyTap("a", "control")
	robotgo.KeyTap("f1", "control")
	// robotgo.KeyTap("enter")
	// robotgo.KeyToggle("enter", "down")
	robotgo.TypeStr("en")

	abitmap := robotgo.CaptureScreen()
	fmt.Println("all...", abitmap)

	bitmap := robotgo.CaptureScreen(10, 20, 30, 40)
	fmt.Println("...", bitmap)

	fx, fy := robotgo.FindBitmap(bitmap)
	fmt.Println("FindBitmap------", fx, fy)

	robotgo.SaveBitmap(bitmap, "test.png", 1)

	var bitmapTest robotgo.Bitmap
	bitTest := robotgo.OpenBitmap("test.png")
	bitmapTest = robotgo.ToBitmap(bitTest)
	fmt.Println("...type", reflect.TypeOf(bitTest), reflect.TypeOf(bitmapTest))

	// robotgo.MouseClick()
	robotgo.Scroll(0, 10)
}

func main() {
	aRobotgo()
}
