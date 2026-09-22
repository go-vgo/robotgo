// Copyright (c) 2016-2026 AtomAI, All rights reserved.
//
// See COPYRIGHT file at top-level directory of this distribution and at
// https://github.com/go-vgo/robotgo/blob/master/LICENSE
//
// Licensed under Apache License, Version 2.0 <LICENSE-APACHE or
// http://www.apache.org/licenses/LICENSE-2.0>
//
// This file may not be copied, modified, or distributed
// except according to those terms.

package robotgo

import (
	"bytes"
	"testing"
)

func TestToRGBAGoSwapsBGRA(t *testing.T) {
	// 2x1 BGRA bitmap: pixel0 = (B1,G2,R3,A4), pixel1 = (B5,G6,R7,A8).
	buf := []uint8{1, 2, 3, 4, 5, 6, 7, 8}
	bmp := Bitmap{ImgBuf: &buf[0], Width: 2, Height: 1, Bytewidth: 8}

	img := ToRGBAGo(bmp)
	want := []uint8{3, 2, 1, 4, 7, 6, 5, 8}
	if !bytes.Equal(img.Pix, want) {
		t.Fatalf("Pix = %v, want %v", img.Pix, want)
	}
	if img.Stride != 8 || img.Rect.Dx() != 2 || img.Rect.Dy() != 1 {
		t.Fatalf("stride/rect = %d/%v", img.Stride, img.Rect)
	}
}
