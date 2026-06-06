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

// Screen capture and geometry are not provided by the RemoteDesktop portal
// input path. Capturing the screen requires the separate ScreenCast portal
// plus a PipeWire client, which is out of scope for this input backend. These
// functions exist for robotgo API parity and report ErrNotSupported (or zero
// values) so callers can detect the limitation.
//
// For screen capture on GNOME/KDE use the wlr-screencopy or xdg-desktop-portal
// ScreenShot path provided by the default (cgo) backend.

// GetScreenSize returns (0, 0); screen geometry is unavailable on this backend.
func GetScreenSize() (int, int) { return 0, 0 }

// GetScaleSize returns (0, 0); unavailable on this backend.
func GetScaleSize(displayId ...int) (int, int) { return 0, 0 }

// GetScreenRect returns an empty Rect; unavailable on this backend.
func GetScreenRect(displayId ...int) Rect { return Rect{} }

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
