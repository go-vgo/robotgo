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

import "errors"

// ErrNotFound is returned when no matching window is found.
var ErrNotFound = errors.New("robotgo: window not found")

// ErrNotSupported is returned when an operation is not supported.
var ErrNotSupported = errors.New("robotgo: operation not supported")
