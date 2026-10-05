//go:build !wayland && !win && !libei && !mac && !x11 && !purego
// +build !wayland,!win,!libei,!mac,!x11,!purego

// Copyright (c) 2016-2026 AtomAI, All rights reserved.
//
// See the COPYRIGHT file at the top-level directory of this distribution and at
// https://github.com/go-vgo/robotgo/blob/master/LICENSE
//
// Licensed under the Apache License, Version 2.0 <LICENSE-APACHE or
// http://www.apache.org/licenses/LICENSE-2.0>
//
// This file may not be copied, modified, or distributed
// except according to those terms.

package robotgo

import (
	"testing"

	"github.com/vcaesar/tt"
)

// Cgo-only surface (C bitmaps and C.MMRGBHex) that the pure-Go backends do
// not expose; the portable half of these checks lives in robotgo_test.go.

func TestColorCgo(t *testing.T) {
	requireDisplay(t)

	s := GetPixelColor(10, 10)
	c := GetPxColor(10, 10)
	tt.Equal(t, s, PadHex(c))
}

func TestImageCgo(t *testing.T) {
	requireDisplay(t)

	bit := CaptureScreen()
	if bit == nil {
		// ToImage must not crash on a failed capture (nil C bitmap).
		tt.Equal(t, 0, Width(ToImage(bit)))
		t.Skip("CaptureScreen returned nil: no Screen Recording permission")
	}
	defer FreeBitmap(bit)

	img := ToImage(bit)
	tt.NotZero(t, Width(img))
	tt.NotZero(t, Height(img))
}
