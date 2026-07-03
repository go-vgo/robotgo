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

import "errors"

// ErrNotFound is returned when a target window or process cannot be located.
var ErrNotFound = errors.New("robotgo: window not found")

// ErrNotSupported is returned when an operation is not supported by the
// running X server (for example a missing XTEST extension).
var ErrNotSupported = errors.New("robotgo: operation not supported")

// ErrNoConnection is returned when the X11 connection is not established.
var ErrNoConnection = errors.New("robotgo: x11 connection not established")
