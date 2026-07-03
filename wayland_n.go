//go:build linux && wayland && !libei
// +build linux,wayland,!libei

// Copyright (c) 2016-2026 AtomAI, All rights reserved.
//
// Licensed under the Apache License, Version 2.0 <LICENSE-APACHE or
// http://www.apache.org/licenses/LICENSE-2.0>

// This file wires the pure-Go Wayland implementation
// (github.com/go-vgo/robotgo/wayland) into the top-level robotgo package.
//
// Build it with:
//
//	go build -tags wayland ./...
//
// Under this tag the Cgo/X11 backend (robotgo.go, key.go, screen.go, ...) is
// excluded via `!wayland` constraints, and the wrappers below forward to the
// wayland package so callers keep using the same robotgo API
// (robotgo.KeyTap, robotgo.Move, ...) with no source changes.
//
// Portable, backend-agnostic code (img.go, ps.go, keycode.go) is compiled as
// usual and provides the image/process/keycode helpers, so those symbols are
// NOT re-exported here.
package robotgo

import (
	"image"

	wl "github.com/go-vgo/robotgo/wayland"
)

// The shared API surface — Version, GetVersion, Sleep, MilliSleep, the
// Bitmap/Point/Size/Rect types and the DisplayID/NotPid/Scale tunables — lives
// in the build-tag-free robotgo_pub.go and is compiled for every backend, so
// it is NOT re-declared here.

// Sentinel errors (values — aliased so errors.Is works across packages).
var (
	ErrNotSupported = wl.ErrNotSupported
	ErrNoConnection = wl.ErrNoConnection
)

// --- General ---

// Close the wayland connection and release resources.
func Close() { wl.Close() }

// --- Keyboard ---

// KeyTap tap the keyboard code.
func KeyTap(key string, args ...any) error { return wl.KeyTap(key, args...) }

// KeyToggle toggle the keyboard.
func KeyToggle(key string, args ...any) error { return wl.KeyToggle(key, args...) }

// KeyDown press down a key.
func KeyDown(key string, args ...any) error { return wl.KeyDown(key, args...) }

// KeyUp release a key.
func KeyUp(key string, args ...any) error { return wl.KeyUp(key, args...) }

// KeyPress press and release a key.
func KeyPress(key string, args ...any) error { return wl.KeyPress(key, args...) }

// Type type a string (alias of TypeStr).
func Type(str string, args ...int) { wl.Type(str, args...) }

// TypeStr type a string.
func TypeStr(str string, args ...int) { wl.TypeStr(str, args...) }

// TypeDelay, SetDelay and CmdCtrl live in robotgo_pub.go (build-tag-free), so
// they are NOT re-declared here.

// --- Mouse ---

// Move move the mouse to (x, y).
func Move(x, y int, displayId ...int) { wl.Move(x, y, displayId...) }

// // MoveRelative move the mouse relative to the current position.
// func MoveRelative(x, y int) { wl.MoveRelative(x, y) }

// MoveSmooth move the mouse smoothly to (x, y).
func MoveSmooth(x, y int, args ...any) bool { return wl.MoveSmooth(x, y, args...) }

// Click click the mouse button.
func Click(args ...any) error { return wl.Click(args...) }

// Toggle toggle the mouse button.
func Toggle(key ...any) error { return wl.Toggle(key...) }

// MouseDown send a mouse down event.
func MouseDown(key ...any) error { return wl.MouseDown(key...) }

// MouseUp send a mouse up event.
func MouseUp(key ...any) error { return wl.MouseUp(key...) }

// Scroll scroll the mouse to (x, y).
func Scroll(x, y int, args ...int) { wl.Scroll(x, y, args...) }

// ScrollDir scroll the mouse to a direction.
func ScrollDir(x int, direction ...any) { wl.ScrollDir(x, direction...) }

// ScrollSmooth scroll the mouse smoothly.
func ScrollSmooth(to int, args ...int) { wl.ScrollSmooth(to, args...) }

// DragSmooth drag the mouse smoothly to (x, y).
func DragSmooth(x, y int, args ...any) { wl.DragSmooth(x, y, args...) }

// MoveClick move and click the mouse.
func MoveClick(x, y int, args ...any) { wl.MoveClick(x, y, args...) }

// Location get the mouse location position, return x, y.
func Location() (int, int) { return wl.Location() }

// GetMousePos get the mouse position, return x, y.
func GetMousePos() (int, int) { return wl.GetMousePos() }

// --- Screen (capture lives in the backend; img.go provides Save*/Width/Height) ---

// GetScreenSize get the screen size.
func GetScreenSize() (int, int) { return wl.GetScreenSize() }

// GetScaleSize get the screen scale size.
func GetScaleSize(displayId ...int) (int, int) { return wl.GetScaleSize(displayId...) }

// GetScreenRect get the screen rect (x, y, w, h).
func GetScreenRect(displayId ...int) Rect {
	r := wl.GetScreenRect(displayId...)
	return Rect{Point{r.X, r.Y}, Size{r.W, r.H}}
}

// DisplaysNum get the number of displays.
func DisplaysNum() int { return wl.DisplaysNum() }

// GetPixelColor get the pixel color at (x, y), return string.
func GetPixelColor(x, y int, displayId ...int) string { return wl.GetPixelColor(x, y, displayId...) }

// CaptureImg capture the screen, return image.Image, error.
func CaptureImg(args ...int) (image.Image, error) { return wl.CaptureImg(args...) }

// Capture capture the screen, return *image.RGBA, error.
func Capture1(args ...int) (*image.RGBA, error) { return wl.Capture(args...) }

// SaveCapture capture the screen and save it to a file.
func SaveCapture(path string, args ...int) error { return wl.SaveCapture(path, args...) }

// PadHex pad a hex color value to 6 characters.
func PadHex(hex uint32) string { return wl.PadHex(hex) }

// --- Window ---

// GetTitle get the window title, return string.
func GetTitle(args ...int) string { return wl.GetTitle(args...) }

// ActiveName active the window by name.
func ActiveName(name string) error { return wl.ActiveName(name) }

// MinWindow set the window min.
func MinWindow(pid int, args ...any) { wl.MinWindow(pid, args...) }

// MaxWindow set the window max.
func MaxWindow(pid int, args ...any) { wl.MaxWindow(pid, args...) }

// CloseWindow close the window.
func CloseWindow(args ...int) { wl.CloseWindow(args...) }

// --- Process (Pids/Process/Kill/... come from the portable ps.go) ---

// GetPid get the current process id.
func GetPid() int { return wl.GetPid() }
