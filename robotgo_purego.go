//go:build wayland || win || libei || mac || x11 || purego
// +build wayland win libei mac x11 purego

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

// This file holds the robotgo APIs shared by every pure-Go backend: they are
// built only from the functions each wiring file (darwin.go, windows_n.go,
// x11_n.go, wayland_n.go, libei.go) forwards, so no backend redeclares them.
package robotgo

import (
	"errors"
	"strconv"
	"unicode/utf8"
)

// GetLocationColor get the color of the pixel at the mouse location
func GetLocationColor(displayId ...int) string {
	x, y := Location()
	return GetPixelColor(x, y, displayId...)
}

// CaptureGo capture the screen and return bitmap(go struct);
// a failed capture yields an empty Bitmap
func CaptureGo(args ...int) Bitmap {
	img, err := CaptureImg(args...)
	if err != nil {
		return Bitmap{}
	}
	return ImgToBitmap(img)
}

// Is64Bit determine whether the process is 64bit
func Is64Bit() bool {
	return strconv.IntSize == 64
}

// UnicodeType tap uint32 unicode,
// the optional args[0] is the target pid where the backend supports it
func UnicodeType(str uint32, args ...int) error {
	if !utf8.ValidRune(rune(str)) {
		return errors.New("robotgo: invalid unicode code point " + strconv.FormatUint(uint64(str), 16))
	}
	if len(args) > 1 {
		args = args[:1]
	}
	return TypeStr(string(rune(str)), args...)
}

// isHandle reports whether a window API's pid is a native window handle
// (HWND / X11 window id): set by any extra arg or NotPid, as in Cgo.
func isHandle(args []int) bool {
	return len(args) > 0 || NotPid
}

// clickTimes clicks button count times, stopping at the first error.
func clickTimes(button string, count int) error {
	for i := 0; i < count; i++ {
		if i > 0 {
			MilliSleep(50)
		}
		if err := Click(button); err != nil {
			return err
		}
	}
	return nil
}
