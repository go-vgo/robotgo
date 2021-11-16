// Copyright 2016 The go-vgo Project Developers. See the COPYRIGHT
// file at the top-level directory of this distribution and at
// https://github.com/go-vgo/robotgo/blob/master/LICENSE
//
// Licensed under the Apache License, Version 2.0 <LICENSE-APACHE or
// http://www.apache.org/licenses/LICENSE-2.0> or the MIT license
// <LICENSE-MIT or http://opensource.org/licenses/MIT>, at your
// option. This file may not be copied, modified, or distributed
// except according to those terms.

//go:build windows
// +build windows

package robotgo

import (
	"syscall"
	"unsafe"

	"github.com/lxn/win"
)

// FindWindow find window hwnd by name
func FindWindow(name string) win.HWND {
	hwnd := win.FindWindow(nil, syscall.StringToUTF16Ptr(name))
	return hwnd
}

// GetHWND get foreground window hwnd
func GetHWND() win.HWND {
	hwnd := win.GetForegroundWindow()
	return hwnd
}

// SendInput send n input event
func SendInput(nInputs uint32, pInputs unsafe.Pointer, cbSize int32) uint32 {
	return win.SendInput(nInputs, pInputs, cbSize)
}

// SendMsg send a message with hwnd
func SendMsg(hwnd win.HWND, msg uint32, wParam, lParam uintptr) uintptr {
	return win.SendMessage(hwnd, msg, wParam, lParam)
}

// SetActiveWindow set window active with hwnd
func SetActiveWindow(hwnd win.HWND) win.HWND {
	return win.SetActiveWindow(hwnd)
}

// SetFocus set window focus with hwnd
func SetFocus(hwnd win.HWND) win.HWND {
	return win.SetFocus(hwnd)
}

// ScaleF get the system scale val
func ScaleF() float64 {
	f := float64(GetMainDPI()) / 96.0
	if f == 0.0 {
		f = 1.0
	}
	return f
}

// GetMainDPI get the display dpi
func GetMainDPI() int {
	return int(GetDPI(GetHWND()))
}

// GetDPI get the window dpi
func GetDPI(hwnd win.HWND) uint32 {
	return win.GetDpiForWindow(hwnd)
}

// GetSysDPI get the system metrics dpi
func GetSysDPI(idx int32, dpi uint32) int32 {
	return win.GetSystemMetricsForDpi(idx, dpi)
}
