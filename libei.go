//go:build linux && libei
// +build linux,libei

// Copyright (c) 2016-2026 AtomAI, All rights reserved.
//
// Licensed under the Apache License, Version 2.0 <LICENSE-APACHE or
// http://www.apache.org/licenses/LICENSE-2.0>

// This file wires the pure-Go libei/portal implementation
// (github.com/go-vgo/robotgo/libei) into the top-level robotgo package.
//
// Build it with:
//
//	go build -tags libei ./...
//
// Under this tag the Cgo/X11 backend (robotgo.go, key.go, screen.go, ...) and
// the wlroots wayland backend (wayland_n.go) are excluded via build
// constraints, and the wrappers below forward to the libei package so callers
// keep using the same robotgo API (robotgo.KeyTap, robotgo.Move, ...) with no
// source changes.
//
// The libei backend drives input through the freedesktop xdg-desktop-portal
// RemoteDesktop interface, which works on GNOME and KDE (unlike the wlroots
// wayland backend). It does not implement screen capture or window management;
// those report ErrNotSupported.
//
// Portable, backend-agnostic code (img.go, ps.go, keycode.go) is compiled as
// usual and provides the image/process/keycode helpers, so those symbols are
// NOT re-exported here.
package robotgo

import (
	"image"

	lb "github.com/go-vgo/robotgo/libei"
)

// The shared API surface — Version, GetVersion, Sleep, MilliSleep, the
// Bitmap/Point/Size/Rect types and the DisplayID/NotPid/Scale tunables — lives
// in the build-tag-free robotgo_pub.go and is compiled for every backend, so
// it is NOT re-declared here.

// Sentinel errors (values — aliased so errors.Is works across packages).
var (
	ErrNotSupported = lb.ErrNotSupported
	ErrNoConnection = lb.ErrNoConnection
)

// --- General ---

// Close the libei/portal connection and release resources.
func Close() { lb.Close() }

// --- Keyboard ---

// KeyTap tap the keyboard code.
func KeyTap(key string, args ...any) error { return lb.KeyTap(key, args...) }

// KeyToggle toggle the keyboard.
func KeyToggle(key string, args ...any) error { return lb.KeyToggle(key, args...) }

// KeyDown press down a key.
func KeyDown(key string, args ...any) error { return lb.KeyDown(key, args...) }

// KeyUp release a key.
func KeyUp(key string, args ...any) error { return lb.KeyUp(key, args...) }

// KeyPress press and release a key.
func KeyPress(key string, args ...any) error { return lb.KeyPress(key, args...) }

// Type type a string (alias of TypeStr).
func Type(str string, args ...int) { lb.Type(str, args...) }

// TypeStr type a string.
func TypeStr(str string, args ...int) { lb.TypeStr(str, args...) }

// TypeDelay type a string with delay.
func TypeDelay(str string, delay int) { lb.TypeDelay(str, delay) }

// SetDelay set the default typing delay.
func SetDelay(d ...int) { lb.SetDelay(d...) }

// CmdCtrl lives in robotgo_pub.go (build-tag-free), so it is NOT re-declared here.

// --- Mouse ---

// Move move the mouse to (x, y).
func Move(x, y int, displayId ...int) { lb.Move(x, y, displayId...) }

// // MoveRelative move the mouse relative to the current position.
// func MoveRelative(x, y int) { lb.MoveRelative(x, y) }

// MoveSmooth move the mouse smoothly to (x, y).
func MoveSmooth(x, y int, args ...any) bool { return lb.MoveSmooth(x, y, args...) }

// Click click the mouse button.
func Click(args ...any) error { return lb.Click(args...) }

// Toggle toggle the mouse button.
func Toggle(key ...any) error { return lb.Toggle(key...) }

// MouseDown send a mouse down event.
func MouseDown(key ...any) error { return lb.MouseDown(key...) }

// MouseUp send a mouse up event.
func MouseUp(key ...any) error { return lb.MouseUp(key...) }

// Scroll scroll the mouse to (x, y).
func Scroll(x, y int, args ...int) { lb.Scroll(x, y, args...) }

// ScrollDir scroll the mouse to a direction.
func ScrollDir(x int, direction ...any) { lb.ScrollDir(x, direction...) }

// ScrollSmooth scroll the mouse smoothly.
func ScrollSmooth(to int, args ...int) { lb.ScrollSmooth(to, args...) }

// DragSmooth drag the mouse smoothly to (x, y).
func DragSmooth(x, y int, args ...any) { lb.DragSmooth(x, y, args...) }

// MoveClick move and click the mouse.
func MoveClick(x, y int, args ...any) { lb.MoveClick(x, y, args...) }

// Location get the mouse location position, return x, y.
func Location() (int, int) { return lb.Location() }

// GetMousePos get the mouse position, return x, y.
func GetMousePos() (int, int) { return lb.GetMousePos() }

// --- Screen (not supported by the libei backend; report ErrNotSupported) ---

// GetScreenSize get the screen size.
func GetScreenSize() (int, int) { return lb.GetScreenSize() }

// GetScaleSize get the screen scale size.
func GetScaleSize(displayId ...int) (int, int) { return lb.GetScaleSize(displayId...) }

// GetScreenRect get the screen rect (x, y, w, h).
func GetScreenRect(displayId ...int) Rect {
	r := lb.GetScreenRect(displayId...)
	return Rect{Point{r.X, r.Y}, Size{r.W, r.H}}
}

// DisplaysNum get the number of displays.
func DisplaysNum() int { return lb.DisplaysNum() }

// GetPixelColor get the pixel color at (x, y), return string.
func GetPixelColor(x, y int, displayId ...int) string { return lb.GetPixelColor(x, y, displayId...) }

// CaptureImg capture the screen, return image.Image, error.
func CaptureImg(args ...int) (image.Image, error) { return lb.CaptureImg(args...) }

// Capture capture the screen, return *image.RGBA, error.
func Capture1(args ...int) (*image.RGBA, error) { return lb.Capture(args...) }

// SaveCapture capture the screen and save it to a file.
func SaveCapture(path string, args ...int) error { return lb.SaveCapture(path, args...) }

// PadHex pad a hex color value to 6 characters.
func PadHex(hex uint32) string { return lb.PadHex(hex) }

// --- Window ---

// GetTitle get the window title, return string.
func GetTitle(args ...int) string { return lb.GetTitle(args...) }

// ActiveName active the window by name.
func ActiveName(name string) error { return lb.ActiveName(name) }

// MinWindow set the window min.
func MinWindow(pid int, args ...any) { lb.MinWindow(pid, args...) }

// MaxWindow set the window max.
func MaxWindow(pid int, args ...any) { lb.MaxWindow(pid, args...) }

// CloseWindow close the window.
func CloseWindow(args ...int) { lb.CloseWindow(args...) }

// --- Process ---

// GetPid get the current process id.
func GetPid() int { return lb.GetPid() }
