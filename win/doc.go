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

// Package win provides a pure-Go Windows implementation of the robotgo API
// for desktop automation: mouse, keyboard, screen capture, window management
// and process enumeration.
//
// Unlike the upstream go-vgo/robotgo, this package uses NO Cgo. It is built
// entirely on top of the Win32 API bindings in github.com/tailscale/win and
// golang.org/x/sys/windows, so it cross-compiles with CGO_ENABLED=0.
//
// The exported API mirrors the Wayland sibling package
// (github.com/go-vgo/robotgo/wayland) so callers can swap
// implementations per platform with a build tag.
package win
