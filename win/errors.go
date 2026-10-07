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

// errSendInput is returned when SendInput does not insert an input event.
var errSendInput = errors.New("robotgo: SendInput failed")

// // errPostMessage is returned when PostMessageW fails.
// var errPostMessage = errors.New("robotgo: PostMessageW failed")

// errSetCursorPos is returned when SetCursorPos fails.
var errSetCursorPos = errors.New("robotgo: SetCursorPos failed")

// errActivate is returned when Windows refuses the foreground change.
var errActivate = errors.New("robotgo: failed to activate the window")
