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
	"unsafe"

	"github.com/vcaesar/imgo"
)

// RGBAToBitmap convert the standard image.RGBA to Bitmap
func RGBAToBitmap(r1 *image.RGBA) (bit Bitmap) {
	bit.Width = r1.Bounds().Size().X
	bit.Height = r1.Bounds().Size().Y
	bit.Bytewidth = r1.Stride

	src := ToUint8p(r1.Pix)
	bit.ImgBuf = src

	bit.BitsPixel = 32
	bit.BytesPerPixel = 32 / 8

	return
}

// ImgToBitmap convert the standard image.Image to Bitmap
func ImgToBitmap(m image.Image) (bit Bitmap) {
	bit.Width = m.Bounds().Size().X
	bit.Height = m.Bounds().Size().Y

	pix, stride, _ := imgo.EncodeImg(m)
	bit.Bytewidth = stride

	src := ToUint8p(pix)
	bit.ImgBuf = src
	//
	bit.BitsPixel = 32
	bit.BytesPerPixel = 32 / 8
	return
}

// ToUint8p convert the []uint8 to uint8 pointer
func ToUint8p(dst []uint8) *uint8 {
	src := make([]uint8, len(dst)+10)
	for i := 0; i < len(dst)-4; i += 4 {
		src[i+3] = dst[i+3]
		src[i] = dst[i+2]
		src[i+1] = dst[i+1]
		src[i+2] = dst[i]
	}

	ptr := unsafe.Pointer(&src[0])
	return (*uint8)(ptr)
}
