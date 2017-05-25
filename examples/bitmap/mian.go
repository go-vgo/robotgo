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

	// Gets all of the screen
	abitMap := robotgo.CaptureScreen()
	fmt.Println("abitMap...", abitMap)

	// Gets part of the screen
	bitmap := robotgo.CaptureScreen(100, 200, 30, 40)
	fmt.Println("CaptureScreen...", bitmap)

	// Searches for needle in bitmap
	fx, fy := robotgo.FindBitmap(bitmap)
	fmt.Println("FindBitmap------", fx, fy)

	// Returns new bitmap object created from a portion of another
	bitpos := robotgo.GetPortion(bitmap, 10, 10, 11, 10)
	fmt.Println(bitpos)

	// Creates bitmap from string by bitmap
	bitstr := robotgo.TostringBitmap(bitmap)
	fmt.Println("bitstr...", bitstr)

	// sbitmap := robotgo.BitmapFromstring(bitstr, 2)
	// fmt.Println("...", sbitmap)

	// Saves image to absolute filepath in the given format
	robotgo.SaveBitmap(bitmap, "test.png")
	robotgo.SaveBitmap(bitmap, "test31.tif", 1)
	robotgo.Convert("test.png", "test.tif") // Convert image

	openbit := robotgo.OpenBitmap("test.tif") // open image bitmap
	fmt.Println("openBitmap...", openbit)
}
