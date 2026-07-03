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

// Package darwin provides a pure-Go macOS implementation of the robotgo API
// for desktop automation: mouse, keyboard and screen capture.
//
// Unlike the upstream go-vgo/robotgo, this package uses NO Cgo. It is built
// entirely on top of the Quartz/CoreGraphics and CoreFoundation system
// frameworks, dynamically loaded at runtime through
// github.com/ebitengine/purego, so it cross-compiles with CGO_ENABLED=0.
//
// The exported API mirrors the Windows and Wayland sibling packages
// (github.com/go-vgo/robotgo/win, github.com/go-vgo/robotgo/wayland) so
// callers can swap implementations per platform with a build tag.
//
// Most input and screen-capture APIs silently fail (or return empty results)
// unless the host process has been granted Accessibility and Screen Recording
// permissions in System Settings → Privacy & Security. Window management
// (GetTitle/ActiveName/MinWindow/MaxWindow/CloseWindow) requires the
// Accessibility/AppKit APIs that are not reachable without Objective-C, so
// those calls report ErrNotSupported.
package darwin
