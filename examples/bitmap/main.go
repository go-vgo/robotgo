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
	// Bitmap
	////////////////////////////////////////////////////////////////////////////////

	// gets all of the screen
	abitMap := robotgo.CaptureScreen()
	fmt.Println("abitMap...", abitMap)

	// gets part of the screen
	bitmap := robotgo.CaptureScreen(100, 200, 30, 40)
	fmt.Println("CaptureScreen...", bitmap)

	gbit := robotgo.ToBitmap(bitmap)
	fmt.Println("go bitmap", gbit, gbit.Width)

	// searches for needle in bitmap
	fx, fy := robotgo.FindBit(bitmap)
	fmt.Println("FindBitmap------", fx, fy)

	color := robotgo.GetColor(bitmap, 1, 2)
	fmt.Println("color...", color)
	cx, cy := robotgo.FindColor(bitmap, robotgo.CHex(color), 1.0)
	fmt.Println("pos...", cx, cy)
	cx, cy = robotgo.FindColor(bitmap, 0xAADCDC)
	fmt.Println("pos...", cx, cy)
	cx, cy = robotgo.FindColorCS(388, 179, 300, 300, 0xAADCDC)
	fmt.Println("pos...", cx, cy)

	cnt := robotgo.CountColor(bitmap, 0xAADCDC)
	fmt.Println("count...", cnt)
	cnt1 := robotgo.CountColorCS(10, 20, 30, 40, 0xAADCDC)
	fmt.Println("count...", cnt1)

	count := robotgo.CountBitmap(abitMap, bitmap)
	fmt.Println("count...", count)

	bit := robotgo.CaptureScreen(1, 2, 40, 40)
	fmt.Println("CaptureScreen...", bit)
	fx, fy = robotgo.FindBitmap(bit)
	fmt.Println("FindBitmap------", fx, fy)
	fx, fy = robotgo.FindBitmap(bit, bitmap)
	fmt.Println("FindBitmap------", fx, fy)

	abool := robotgo.PointInBounds(bitmap, 1, 2)
	fmt.Println("point in bounds...", abool)

	// returns new bitmap object created from a portion of another
	bitpos := robotgo.GetPortion(bitmap, 10, 10, 11, 10)
	fmt.Println(bitpos)

	// creates bitmap from string by bitmap
	bitstr := robotgo.TostringBitmap(bitmap)
	fmt.Println("bitstr...", bitstr)

	// sbitmap := robotgo.BitmapFromstring(bitstr, 2)
	// fmt.Println("...", sbitmap)

	// saves image to absolute filepath in the given format
	robotgo.SaveBitmap(bitmap, "test.png")
	robotgo.SaveBitmap(bitmap, "test31.tif", 1)

	// convert image
	robotgo.Convert("test.png", "test.tif")

	// open image bitmap
	openbit := robotgo.OpenBitmap("test.tif")
	fmt.Println("openBitmap...", openbit)

	fx, fy = robotgo.FindBitmap(openbit)
	fmt.Println("FindBitmap------", fx, fy)
}
