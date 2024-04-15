// Copyright 2016 The go-vgo Project Developers. See the COPYRIGHT
// file at the top-level directory of this distribution and at
// https://github.com/go-vgo/robotgo/blob/master/LICENSE
//
// Licensed under the Apache License, Version 2.0 <LICENSE-APACHE or
// http://www.apache.org/licenses/LICENSE-2.0> or the MIT license
// <LICENSE-MIT or http://opensource.org/licenses/MIT>, at your
// option. This file may not be copied, modified, or distributed
// except according to those terms.

package robotgo

import (
	"image"

	"github.com/kbinani/screenshot"
)

// GetDisplayBounds gets the display screen bounds
func GetDisplayBounds(i int) (x, y, w, h int) {
	bs := screenshot.GetDisplayBounds(i)
	return bs.Min.X, bs.Min.Y, bs.Dx(), bs.Dy()
}

// GetDisplayRect gets the display rect
func GetDisplayRect(i int) Rect {
	x, y, w, h := GetDisplayBounds(i)
	return Rect{
		Point{X: x, Y: y},
		Size{W: w, H: h}}
}

// Capture capture the screenshot
func Capture(args ...int) (*image.RGBA, error) {
	displayId := 0
	if DisplayID != -1 {
		displayId = DisplayID
	}

	if len(args) > 4 {
		displayId = args[4]
	}

	var x, y, w, h int
	if len(args) > 3 {
		x, y, w, h = args[0], args[1], args[2], args[3]
	} else {
		x, y, w, h = GetDisplayBounds(displayId)
	}

	return screenshot.Capture(x, y, w, h)
}
