//go:build darwin
// +build darwin

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

package darwin

import (
	"fmt"
	"image"
	"unsafe"
)

// displayForIndex resolves an optional display index to a CGDirectDisplayID,
// defaulting to the main display.
func displayForIndex(displayId ...int) uint32 {
	if !loaded {
		return 0
	}
	if len(displayId) > 0 && displayId[0] > 0 {
		var ids [16]uint32
		var count uint32
		if cgGetActiveDisplayLst(uint32(len(ids)), &ids[0], &count) == 0 {
			idx := displayId[0]
			// count reports every active display, not only the ones that
			// fit in ids, so bound by both.
			if idx < int(count) && idx < len(ids) {
				return ids[idx]
			}
		}
	}
	return cgMainDisplayID()
}

// MainDisplayID returns the main display's CGDirectDisplayID, mirroring the
// Cgo backend's GetMainId. Returns 0 when the frameworks are unavailable.
func MainDisplayID() int {
	if !loaded {
		return 0
	}
	return int(cgMainDisplayID())
}

// GetScreenSize returns the main display's width and height in pixels.
func GetScreenSize() (int, int) {
	if !loaded {
		return 0, 0
	}
	id := cgMainDisplayID()
	return int(cgDisplayPixelsWide(id)), int(cgDisplayPixelsHigh(id))
}

// ScaleF returns the backing scale factor (pixel width / point width) of the
// selected display, mirroring the Cgo backend's sys_scale. Returns 1 when the
// frameworks are unavailable or the mode cannot be read.
func ScaleF(displayId ...int) float64 {
	if !loaded {
		return 1
	}
	mode := cgDisplayCopyDisplayMode(displayForIndex(displayId...))
	if mode == 0 {
		return 1
	}
	defer cgDisplayModeRelease(mode)
	w := cgDisplayModeGetWidth(mode)
	if w == 0 {
		return 1
	}
	return float64(cgDisplayModeGetPixelWidth(mode)) / float64(w)
}

// GetScaleSize returns the screen size in physical pixels (point size times
// the backing scale), matching robotgo.GetScaleSize on the Cgo backend.
func GetScaleSize(displayId ...int) (int, int) {
	w, h := GetScreenSize()
	f := ScaleF(displayId...)
	return int(float64(w) * f), int(float64(h) * f)
}

// GetScreenRect returns the selected display's bounds (origin in global
// display space, size in points), defaulting to the main display.
func GetScreenRect(displayId ...int) Rect {
	if !loaded {
		return Rect{}
	}
	b := cgDisplayBounds(displayForIndex(displayId...))
	return Rect{
		Point: Point{X: int(b.Origin.X), Y: int(b.Origin.Y)},
		Size:  Size{W: int(b.Size.Width), H: int(b.Size.Height)},
	}
}

// DisplaysNum returns the number of active displays.
func DisplaysNum() int {
	if !loaded {
		return 1
	}
	var count uint32
	if cgGetActiveDisplayLst(0, nil, &count) != 0 || count == 0 {
		return 1
	}
	return int(count)
}

// GetPixelColor returns the pixel color at (x, y) as a 6-digit hex string.
func GetPixelColor(x, y int, displayId ...int) string {
	img, err := CaptureImg(x, y, 1, 1, indexOf(displayId))
	if err != nil || img == nil {
		return "000000"
	}
	r, g, b, _ := img.At(img.Bounds().Min.X, img.Bounds().Min.Y).RGBA()
	return fmt.Sprintf("%02x%02x%02x", r>>8, g>>8, b>>8)
}

// indexOf returns the first display id or -1 if none was supplied.
func indexOf(displayId []int) int {
	if len(displayId) > 0 {
		return displayId[0]
	}
	return -1
}

// CaptureImg captures the screen and returns an image.Image.
// Optional args: x, y, w, h (region) and a trailing display index. With no
// args the full main display is captured.
func CaptureImg(args ...int) (image.Image, error) {
	if !loaded {
		return nil, ErrNotSupported
	}

	display := cgMainDisplayID()
	if len(args) > 4 && args[4] >= 0 {
		display = displayForIndex(args[4])
	}

	var rect CGRect
	if len(args) >= 4 {
		if args[2] <= 0 || args[3] <= 0 {
			return nil, fmt.Errorf("robotgo: invalid capture size %dx%d", args[2], args[3])
		}
		rect = CGRect{
			Origin: CGPoint{X: float64(args[0]), Y: float64(args[1])},
			Size:   CGSize{Width: float64(args[2]), Height: float64(args[3])},
		}
	} else {
		rect = cgDisplayBounds(display)
	}

	cgImage := cgDisplayCreateImageForRect(display, rect)
	if cgImage == 0 {
		return nil, fmt.Errorf("robotgo: CGDisplayCreateImageForRect failed (check Screen Recording permission)")
	}
	defer cgImageRelease(cgImage)

	w := int(cgImageGetWidth(cgImage))
	h := int(cgImageGetHeight(cgImage))
	bpr := int(cgImageGetBytesPerRow(cgImage))
	if w <= 0 || h <= 0 {
		return nil, fmt.Errorf("robotgo: empty capture %dx%d", w, h)
	}

	provider := cgImageGetDataProvider(cgImage)
	if provider == 0 {
		return nil, fmt.Errorf("robotgo: CGImageGetDataProvider failed")
	}
	data := cgDataProviderCopyData(provider)
	if data == 0 {
		return nil, fmt.Errorf("robotgo: CGDataProviderCopyData failed")
	}
	defer cfRelease(data)

	ptr := cfDataGetBytePtr(data)
	length := int(cfDataGetLength(data))
	if ptr == nil || length < bpr*h {
		return nil, fmt.Errorf("robotgo: short pixel buffer")
	}
	src := unsafe.Slice((*byte)(ptr), length)

	// CoreGraphics display images are 32-bit, host-endian, alpha skip-first:
	// in memory the bytes are B, G, R, X. Convert to RGBA.
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	for y := 0; y < h; y++ {
		row := y * bpr
		for x := 0; x < w; x++ {
			s := row + x*4
			d := img.PixOffset(x, y)
			img.Pix[d+0] = src[s+2] // R
			img.Pix[d+1] = src[s+1] // G
			img.Pix[d+2] = src[s+0] // B
			img.Pix[d+3] = 0xff     // A
		}
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
