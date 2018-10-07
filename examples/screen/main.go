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

func bitmap() {
	abitMap := robotgo.CaptureScreen()
	fmt.Println("abitMap...", abitMap)
	gbit := robotgo.ToBitmap(abitMap)
	fmt.Println("bitmap...", gbit.Width)

	gbitMap := robotgo.GoCaptureScreen()
	fmt.Println("GoCaptureScreen...", gbitMap.Width)
	// fmt.Println("...", gbitmap.Width, gbitmap.BytesPerPixel)

	robotgo.SaveCapture("saveCapture.png", 10, 20, 100, 100)
}

func color() {
	// gets the pixel color at 100, 200.
	color := robotgo.GetPixelColor(100, 200)
	fmt.Println("color----", color, "-----------------")

	clo := robotgo.GetPxColor(100, 200)
	fmt.Println("color...", clo)
	clostr := robotgo.PadHex(clo)
	fmt.Println("color...", clostr)

	rgb := robotgo.RgbToHex(255, 100, 200)
	rgbstr := robotgo.PadHex(robotgo.U32ToHex(rgb))
	fmt.Println("rgb...", rgbstr)

	hex := robotgo.HexToRgb(uint32(rgb))
	fmt.Println("hex...", hex)
	hexh := robotgo.PadHex(robotgo.U8ToHex(hex))
	fmt.Println("HexToRgb...", hexh)

	// gets the pixel color at 10, 20.
	color2 := robotgo.GetPixelColor(10, 20)
	fmt.Println("color---", color2)
}

func screen() {
	////////////////////////////////////////////////////////////////////////////////
	// Read the screen
	////////////////////////////////////////////////////////////////////////////////

	bitmap()

	// gets the screen width and height
	sx, sy := robotgo.GetScreenSize()
	fmt.Println("...", sx, sy)

	color()
}

func main() {
	screen()
}
