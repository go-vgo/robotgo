//go:build darwin && (mac || purego)
// +build darwin
// +build mac purego

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

// This file wires the pure-Go (Cgo-free) macOS implementation
// (github.com/go-vgo/robotgo/darwin) into the top-level robotgo package.
//
// Build it with:
//
//	go build -tags mac ./...
//
// or use the cross-platform pure-Go default backends (mac on macOS,
// win on Windows, wayland on Linux):
//
//	go build -tags purego ./...
//
// Under this tag the default Cgo backend (robotgo.go, key.go, robotgo_mac.go,
// ...) is excluded via `!mac` constraints, and the wrappers below forward to
// the darwin package so callers keep using the same robotgo API
// (robotgo.KeyTap, robotgo.Move, ...) with no source changes.
//
// The darwin backend drives input and screen capture through the
// Quartz/CoreGraphics frameworks loaded at runtime via purego (no Cgo). Window
// management is not implemented (it needs the Objective-C Accessibility/AppKit
// APIs) and reports ErrNotSupported.
//
// Portable, backend-agnostic code (img.go, ps.go, keycode.go, screen.go and
// the shared helpers in robotgo_pub.go) is compiled as usual and provides the
// image/process/keycode/clipboard helpers (Save, Width, Height, Pids, Process,
// Capture, TypeDelay, SetDelay, CmdCtrl, MoveRelative, ...), so those symbols
// are NOT re-exported here.
package robotgo

import (
	"image"

	dm "github.com/go-vgo/robotgo/darwin"
)

// The shared API surface — Version, GetVersion, Sleep, MilliSleep, the
// Bitmap/Point/Size/Rect types and the DisplayID/NotPid/Scale tunables — lives
// in the build-tag-free robotgo_pub.go and is compiled for every backend, so
// it is NOT re-declared here.

// Sentinel errors (values — aliased so errors.Is works across packages).
var (
	ErrNotFound     = dm.ErrNotFound
	ErrNotSupported = dm.ErrNotSupported
)

// --- Keyboard ---

// KeyTap tap the keyboard code.
func KeyTap(key string, args ...any) error { return dm.KeyTap(key, args...) }

// KeyToggle toggle the keyboard.
func KeyToggle(key string, args ...any) error { return dm.KeyToggle(key, args...) }

// KeyDown press down a key.
func KeyDown(key string, args ...any) error { return dm.KeyDown(key, args...) }

// KeyUp release a key.
func KeyUp(key string, args ...any) error { return dm.KeyUp(key, args...) }

// KeyPress press and release a key.
func KeyPress(key string, args ...any) error { return dm.KeyPress(key, args...) }

// Type type a string (alias of TypeStr).
func Type(str string, args ...int) { dm.Type(str, args...) }

// TypeStr type a string.
func TypeStr(str string, args ...int) { dm.TypeStr(str, args...) }

// TypeDelay, SetDelay and CmdCtrl live in robotgo_pub.go (build-tag-free), so
// they are NOT re-declared here.

// --- Mouse ---

// Move move the mouse to (x, y).
func Move(x, y int, displayId ...int) { dm.Move(x, y, displayId...) }

// // MoveRelative lives in robotgo_pub.go, so it is NOT re-declared here.

// MoveSmooth move the mouse smoothly to (x, y).
func MoveSmooth(x, y int, args ...any) bool { return dm.MoveSmooth(x, y, args...) }

// Click click the mouse button.
func Click(args ...any) error { return dm.Click(args...) }

// Toggle toggle the mouse button.
func Toggle(key ...any) error { return dm.Toggle(key...) }

// MouseDown send a mouse down event.
func MouseDown(key ...any) error { return dm.MouseDown(key...) }

// MouseUp send a mouse up event.
func MouseUp(key ...any) error { return dm.MouseUp(key...) }

// Scroll scroll the mouse to (x, y).
func Scroll(x, y int, args ...int) { dm.Scroll(x, y, args...) }

// ScrollDir scroll the mouse to a direction.
func ScrollDir(x int, direction ...any) { dm.ScrollDir(x, direction...) }

// ScrollSmooth scroll the mouse smoothly.
func ScrollSmooth(to int, args ...int) { dm.ScrollSmooth(to, args...) }

// DragSmooth drag the mouse smoothly to (x, y).
func DragSmooth(x, y int, args ...any) { dm.DragSmooth(x, y, args...) }

// MoveClick move and click the mouse.
func MoveClick(x, y int, args ...any) { dm.MoveClick(x, y, args...) }

// Location get the mouse location position, return x, y.
func Location() (int, int) { return dm.Location() }

// GetMousePos get the mouse position, return x, y.
func GetMousePos() (int, int) { return dm.GetMousePos() }

// --- Screen (the portable screen.go provides Capture/GetDisplayBounds; img.go
// provides Save/Width/Height) ---

// GetScreenSize get the screen size.
func GetScreenSize() (int, int) { return dm.GetScreenSize() }

// GetScaleSize get the screen scale size.
func GetScaleSize(displayId ...int) (int, int) { return dm.GetScaleSize(displayId...) }

// GetScreenRect get the screen rect (x, y, w, h).
func GetScreenRect(displayId ...int) Rect {
	r := dm.GetScreenRect(displayId...)
	return Rect{Point{r.X, r.Y}, Size{r.W, r.H}}
}

// GetMainId get the main display id.
func GetMainId() int { return dm.MainDisplayID() }

// IsMain is main display.
func IsMain(displayId int) bool { return displayId == GetMainId() }

// DisplaysNum get the number of displays.
func DisplaysNum() int { return dm.DisplaysNum() }

// GetPixelColor get the pixel color at (x, y), return string.
func GetPixelColor(x, y int, displayId ...int) string { return dm.GetPixelColor(x, y, displayId...) }

// CaptureImg capture the screen, return image.Image, error.
func CaptureImg(args ...int) (image.Image, error) { return dm.CaptureImg(args...) }

// Capture1 capture the screen, return *image.RGBA, error.
// (The portable screen.go owns the name Capture.)
func Capture1(args ...int) (*image.RGBA, error) { return dm.Capture(args...) }

// SaveCapture capture the screen and save it to a file.
func SaveCapture(path string, args ...int) error { return dm.SaveCapture(path, args...) }

// PadHex pad a hex color value to 6 characters.
func PadHex(hex uint32) string { return dm.PadHex(hex) }

// --- Window (not supported by the darwin backend; report ErrNotSupported) ---

// GetTitle get the window title, return string.
func GetTitle(args ...int) string { return dm.GetTitle(args...) }

// ActiveName active the window by name.
func ActiveName(name string) error { return dm.ActiveName(name) }

// MinWindow set the window min.
func MinWindow(pid int, args ...any) { dm.MinWindow(pid, args...) }

// MaxWindow set the window max.
func MaxWindow(pid int, args ...any) { dm.MaxWindow(pid, args...) }

// CloseWindow close the window.
func CloseWindow(args ...int) { dm.CloseWindow(args...) }

// --- Process (Pids/Process/Kill/... come from the portable ps.go) ---

// GetPid get the current process id.
func GetPid() int { return dm.GetPid() }
