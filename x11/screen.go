//go:build linux
// +build linux

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

package x11

import (
	"errors"
	"image"

	"github.com/jezek/xgb/xinerama"
	"github.com/jezek/xgb/xproto"
)

// GetScreenSize returns the full virtual screen size in pixels.
func GetScreenSize() (int, int) {
	c, err := ensureConn()
	if err != nil {
		return 0, 0
	}
	s := c.xu.Setup().DefaultScreen(c.c)
	return int(s.WidthInPixels), int(s.HeightInPixels)
}

// GetScaleSize returns the screen scale size. X11 reports physical pixels, so
// this is identical to GetScreenSize.
func GetScaleSize(displayId ...int) (int, int) {
	return GetScreenSize()
}

// GetScreenRect returns the rect (x, y, w, h) of a display. With no displayId
// (or an out-of-range one) it returns the full virtual screen.
func GetScreenRect(displayId ...int) Rect {
	c, err := ensureConn()
	if err != nil {
		return Rect{}
	}

	if c.xineramaOK && len(displayId) > 0 && displayId[0] >= 0 {
		if r, ok := xineramaRect(c, displayId[0]); ok {
			return r
		}
	}

	s := c.xu.Setup().DefaultScreen(c.c)
	return Rect{Point{0, 0}, Size{int(s.WidthInPixels), int(s.HeightInPixels)}}
}

// xineramaRect returns the geometry of the given Xinerama screen index.
func xineramaRect(c *conn, idx int) (Rect, bool) {
	reply, err := xinerama.QueryScreens(c.c).Reply()
	if err != nil || reply == nil || idx >= len(reply.ScreenInfo) {
		return Rect{}, false
	}
	si := reply.ScreenInfo[idx]
	return Rect{
		Point{int(si.XOrg), int(si.YOrg)},
		Size{int(si.Width), int(si.Height)},
	}, true
}

// DisplaysNum returns the number of displays.
func DisplaysNum() int {
	c, err := ensureConn()
	if err != nil {
		return 0
	}
	if c.xineramaOK {
		if reply, err := xinerama.QueryScreens(c.c).Reply(); err == nil && reply != nil {
			if reply.Number > 0 {
				return int(reply.Number)
			}
		}
	}
	return len(c.xu.Setup().Roots)
}

// MainDisplayID returns the index of the default screen among the setup
// roots, mirroring the Cgo backend's GetMainId. Returns -1 when no X11
// connection is available.
func MainDisplayID() int {
	c, err := ensureConn()
	if err != nil {
		return -1
	}
	setup := c.xu.Setup()
	def := setup.DefaultScreen(c.c)
	for i := range setup.Roots {
		if setup.Roots[i].Root == def.Root {
			return i
		}
	}
	return -1
}

// ScaleX returns the horizontal DPI of the default screen (dots per inch),
// mirroring the Cgo backend's deprecated ScaleX. Returns 0 when no X11
// connection or physical size information is available.
func ScaleX() int {
	c, err := ensureConn()
	if err != nil {
		return 0
	}
	s := c.xu.Setup().DefaultScreen(c.c)
	if s.WidthInMillimeters == 0 {
		return 0
	}
	return int(float64(s.WidthInPixels) * 25.4 / float64(s.WidthInMillimeters))
}

// GetPixelColor returns the pixel color at (x, y) as a 6-char hex string.
func GetPixelColor(x, y int, displayId ...int) string {
	img, err := CaptureImg(x, y, 1, 1)
	if err != nil {
		return ""
	}
	r, g, b, _ := img.At(0, 0).RGBA()
	hex := (uint32(r>>8) << 16) | (uint32(g>>8) << 8) | uint32(b>>8)
	return PadHex(hex)
}

// CaptureImg captures the screen and returns an image.Image.
//
//	CaptureImg()                 // full screen
//	CaptureImg(x, y, w, h int)
func CaptureImg(args ...int) (image.Image, error) {
	return Capture(args...)
}

// Capture captures the screen and returns an *image.RGBA.
func Capture(args ...int) (*image.RGBA, error) {
	c, err := ensureConn()
	if err != nil {
		return nil, err
	}

	var x, y, w, h int
	if len(args) >= 4 {
		x, y, w, h = args[0], args[1], args[2], args[3]
	} else {
		r := GetScreenRect()
		x, y, w, h = r.X, r.Y, r.W, r.H
	}
	if w <= 0 || h <= 0 {
		return nil, errors.New("robotgo: capture size must be positive")
	}

	reply, err := xproto.GetImage(
		c.c, xproto.ImageFormatZPixmap, xproto.Drawable(c.root),
		int16(x), int16(y), uint16(w), uint16(h), 0xffffffff).Reply()
	if err != nil {
		return nil, err
	}
	if reply == nil {
		return nil, ErrNotSupported
	}

	return zpixmapToRGBA(reply.Data, w, h), nil
}

// zpixmapToRGBA converts X11 ZPixmap data (BGRX, little-endian, typically 4
// bytes per pixel) into an *image.RGBA.
func zpixmapToRGBA(data []byte, w, h int) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	if w == 0 || h == 0 {
		return img
	}

	bpp := len(data) / (w * h)
	if bpp < 3 {
		return img
	}

	di := 0
	for i := 0; i+3 < len(img.Pix) && di+2 < len(data); i += 4 {
		img.Pix[i+0] = data[di+2] // R (ZPixmap is BGR(X))
		img.Pix[i+1] = data[di+1] // G
		img.Pix[i+2] = data[di+0] // B
		img.Pix[i+3] = 255        // A (X bytes carry no usable alpha)
		di += bpp
	}
	return img
}
