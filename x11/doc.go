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

// Package x11 is a pure-Go (Cgo-free) X11 backend for robotgo.
//
// It mirrors the robotgo API surface — keyboard, mouse, screen, window and
// process — without any C bindings, talking to the X server over the wire with
// github.com/jezek/xgb (XTEST for input, core protocol for screen capture and
// pointer queries) and github.com/jezek/xgbutil (EWMH/ICCCM for window
// management).
//
// It is wired into the top-level robotgo package by x11_n.go under the `x11`
// build tag:
//
//	go build -tags x11 ./...
//
// Under that tag the default Cgo X11 backend (robotgo.go, key.go,
// robotgo_x11.go, ...) is excluded, and the top-level robotgo functions forward
// to this package, so callers keep using the same API (robotgo.KeyTap,
// robotgo.Move, ...) with no source changes.
package x11
