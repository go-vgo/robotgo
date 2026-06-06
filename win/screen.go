//go:build windows
// +build windows

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

package win

import (
	"fmt"
	"image"
	"unsafe"

	"github.com/tailscale/win"
)

// GetScreenSize returns the primary display's width and height in pixels.
func GetScreenSize() (int, int) {
	return int(win.GetSystemMetrics(win.SM_CXSCREEN)), int(win.GetSystemMetrics(win.SM_CYSCREEN))
}

// GetScaleSize returns the scaled screen size. The GDI capture path already
// works in physical pixels, so this returns the same value as GetScreenSize.
// Provided for robotgo API parity.
func GetScaleSize(displayId ...int) (int, int) {
	return GetScreenSize()
}

// GetScreenRect returns the primary display rectangle.
func GetScreenRect(displayId ...int) Rect {
	w, h := GetScreenSize()
	return Rect{Point: Point{X: 0, Y: 0}, Size: Size{W: w, H: h}}
}

// DisplaysNum returns the number of displays.
func DisplaysNum() int {
	n := int(win.GetSystemMetrics(win.SM_CMONITORS))
	if n <= 0 {
		return 1
	}
	return n
}

// GetPixelColor returns the pixel color at (x, y) as a 6-digit hex string.
func GetPixelColor(x, y int, displayId ...int) string {
	img, err := CaptureImg(x, y, 1, 1)
	if err != nil || img == nil {
		return "000000"
	}
	r, g, b, _ := img.At(0, 0).RGBA()
	return fmt.Sprintf("%02x%02x%02x", r>>8, g>>8, b>>8)
}

// CaptureImg captures the screen and returns an image.Image.
// Optional args: x, y, w, h (region). With no args the full primary screen
// is captured.
func CaptureImg(args ...int) (image.Image, error) {
	x, y := 0, 0
	w, h := GetScreenSize()
	if len(args) >= 4 {
		x, y, w, h = args[0], args[1], args[2], args[3]
	}
	if w <= 0 || h <= 0 {
		return nil, fmt.Errorf("robotgo: invalid capture size %dx%d", w, h)
	}

	hScreen := win.GetDC(0)
	if hScreen == 0 {
		return nil, fmt.Errorf("robotgo: GetDC failed")
	}
	defer win.ReleaseDC(0, hScreen)

	hDC := win.CreateCompatibleDC(hScreen)
	if hDC == 0 {
		return nil, fmt.Errorf("robotgo: CreateCompatibleDC failed")
	}
	defer win.DeleteDC(hDC)

	hBmp := win.CreateCompatibleBitmap(hScreen, int32(w), int32(h))
	if hBmp == 0 {
		return nil, fmt.Errorf("robotgo: CreateCompatibleBitmap failed")
	}
	defer win.DeleteObject(win.HGDIOBJ(hBmp))

	old := win.SelectObject(hDC, win.HGDIOBJ(hBmp))
	if !win.BitBlt(hDC, 0, 0, int32(w), int32(h), hScreen, int32(x), int32(y), win.SRCCOPY) {
		win.SelectObject(hDC, old)
		return nil, fmt.Errorf("robotgo: BitBlt failed")
	}
	// GetDIBits requires the bitmap not be selected into a DC.
	win.SelectObject(hDC, old)

	bi := win.BITMAPINFO{
		BmiHeader: win.BITMAPINFOHEADER{
			BiSize:        uint32(unsafe.Sizeof(win.BITMAPINFOHEADER{})),
			BiWidth:       int32(w),
			BiHeight:      -int32(h), // negative = top-down rows
			BiPlanes:      1,
			BiBitCount:    32,
			BiCompression: win.BI_RGB,
		},
	}

	buf := make([]byte, w*h*4)
	ret := win.GetDIBits(hDC, hBmp, 0, uint32(h), &buf[0], &bi, win.DIB_RGB_COLORS)
	if ret == 0 {
		return nil, fmt.Errorf("robotgo: GetDIBits failed")
	}

	// GDI 32bpp DIBs are stored as B, G, R, X bytes; convert to RGBA.
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	for i := 0; i < w*h; i++ {
		b := buf[i*4+0]
		g := buf[i*4+1]
		r := buf[i*4+2]
		img.Pix[i*4+0] = r
		img.Pix[i*4+1] = g
		img.Pix[i*4+2] = b
		img.Pix[i*4+3] = 0xff
	}
	return img, nil
}

// Capture captures the screen and returns an *image.RGBA.
func Capture(args ...int) (*image.RGBA, error) {
	img, err := CaptureImg(args...)
	if err != nil {
		return nil, err
	}
	if rgba, ok := img.(*image.RGBA); ok {
		return rgba, nil
	}
	bounds := img.Bounds()
	rgba := image.NewRGBA(bounds)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			rgba.Set(x, y, img.At(x, y))
		}
	}
	return rgba, nil
}
