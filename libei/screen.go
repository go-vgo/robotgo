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

package libei

import "image"

// Screen capture is not provided by the RemoteDesktop portal input path: it
// would require a PipeWire client on top of the linked ScreenCast stream,
// which is out of scope for this input backend. Screen geometry, however, is
// reported by the portal for every linked stream (see LinkScreenCast), so the
// size/rect helpers below work whenever a stream was negotiated and return
// zero values otherwise. Capture and pixel reads report ErrNotSupported.

// GetScreenSize returns the union size of all linked streams, or (0, 0) when
// no ScreenCast stream is linked.
func GetScreenSize() (int, int) {
	c, err := ensureConn()
	if err != nil {
		return 0, 0
	}
	r, _ := c.bounds()
	return r.W, r.H
}

// GetScaleSize returns the size of the stream selected by displayId (stream
// geometry is already in logical pixels), or the union size when no display
// is given.
func GetScaleSize(displayId ...int) (int, int) {
	if len(displayId) == 0 {
		return GetScreenSize()
	}
	r := GetScreenRect(displayId...)
	return r.W, r.H
}

// GetScreenRect returns the rectangle of the stream selected by displayId
// (default 0), or an empty Rect when no stream is linked.
func GetScreenRect(displayId ...int) Rect {
	c, err := ensureConn()
	if err != nil || len(c.streams) == 0 {
		return Rect{}
	}
	idx := 0
	if len(displayId) > 0 && displayId[0] >= 0 && displayId[0] < len(c.streams) {
		idx = displayId[0]
	}
	s := c.streams[idx]
	return Rect{Point{int(s.x), int(s.y)}, Size{int(s.width), int(s.height)}}
}

// DisplaysNum returns the number of ScreenCast streams linked to the session
// (0 unless a ScreenCast source was negotiated).
func DisplaysNum() int {
	c, err := ensureConn()
	if err != nil {
		return 0
	}
	return len(c.streams)
}

// GetPixelColor returns "000000"; screen reading is unavailable on this backend.
func GetPixelColor(x, y int, displayId ...int) string { return "000000" }

// CaptureImg is not supported by this backend.
func CaptureImg(args ...int) (image.Image, error) { return nil, ErrNotSupported }

// Capture is not supported by this backend.
func Capture(args ...int) (*image.RGBA, error) { return nil, ErrNotSupported }
