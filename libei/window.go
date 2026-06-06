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

// Window management is not exposed by the RemoteDesktop portal. The portal's
// security model deliberately prevents clients from enumerating, focusing, or
// resizing other applications' windows. These functions exist for robotgo API
// parity and are no-ops / report ErrNotSupported.

// GetTitle returns an empty string; window titles are unavailable on this backend.
func GetTitle(args ...int) string { return "" }

// ActiveName is not supported by this backend.
func ActiveName(name string) error { return ErrNotSupported }

// MinWindow is a no-op on this backend.
func MinWindow(pid int, args ...interface{}) {}

// MaxWindow is a no-op on this backend.
func MaxWindow(pid int, args ...interface{}) {}

// CloseWindow is a no-op on this backend.
func CloseWindow(args ...int) {}
