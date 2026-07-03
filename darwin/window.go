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

// Window management on macOS requires the Accessibility (AXUIElement) and
// AppKit APIs, which are only reachable through the Objective-C runtime and
// are out of scope for this pure-Go CoreGraphics backend. These calls report
// ErrNotSupported (or empty results), mirroring the libei backend.

// GetTitle returns the active window title. Not supported by this backend.
func GetTitle(args ...int) string {
	return ""
}

// ActiveName brings the first window whose title contains name to the
// foreground. Not supported by this backend.
func ActiveName(name string) error {
	return ErrNotSupported
}

// MinWindow minimizes a window. Not supported by this backend.
func MinWindow(pid int, args ...interface{}) {}

// MaxWindow maximizes a window. Not supported by this backend.
func MaxWindow(pid int, args ...interface{}) {}

// CloseWindow closes a window. Not supported by this backend.
func CloseWindow(args ...int) {}
