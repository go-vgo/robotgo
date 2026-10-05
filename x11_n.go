//go:build linux && x11
// +build linux,x11

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

// This file wires the pure-Go (Cgo-free) X11 implementation
// (github.com/go-vgo/robotgo/x11) into the top-level robotgo package.
//
// Build it with:
//
//	go build -tags x11 ./...
//
// Under this tag the default Cgo/X11 backend (robotgo.go, key.go,
// robotgo_x11.go, ...) is excluded via `!x11` constraints, and the wrappers
// below forward to the x11 package so callers keep using the same robotgo API
// (robotgo.KeyTap, robotgo.Move, ...) with no source changes.
//
// Portable, backend-agnostic code (img.go, ps.go, keycode.go, screen.go) is
// compiled as usual and provides the image/process/keycode helpers, so those
// symbols are NOT re-exported here. In particular screen.go already provides
// Capture()/SaveCaptureGo(), so this file exposes the x11 capture as Capture1
// to avoid a redeclaration (mirroring wayland_n.go / windows_n.go).
package robotgo

import (
	"image"

	x11 "github.com/go-vgo/robotgo/x11"
)

// The shared API surface — Version, GetVersion, Sleep, MilliSleep, the
// Bitmap/Point/Size/Rect types and the DisplayID/NotPid/Scale tunables — lives
// in the build-tag-free robotgo_pub.go and is compiled for every backend, so
// it is NOT re-declared here.

// Sentinel errors (values — aliased so errors.Is works across packages).
var (
	ErrNotFound     = x11.ErrNotFound
	ErrNotSupported = x11.ErrNotSupported
	ErrNoConnection = x11.ErrNoConnection
)

// --- General ---

// Close the X11 connection and release resources.
func Close() { x11.Close() }

// --- Keyboard ---

// KeyTap tap the keyboard code.
func KeyTap(key string, args ...any) error { return x11.KeyTap(key, args...) }

// KeyToggle toggle the keyboard.
func KeyToggle(key string, args ...any) error { return x11.KeyToggle(key, args...) }

// KeyDown press down a key.
func KeyDown(key string, args ...any) error { return x11.KeyDown(key, args...) }

// KeyUp release a key.
func KeyUp(key string, args ...any) error { return x11.KeyUp(key, args...) }

// KeyPress press and release a key.
func KeyPress(key string, args ...any) error { return x11.KeyPress(key, args...) }

// Type type a string (alias of TypeStr).
func Type(str string, args ...int) int { return x11.Type(str, args...) }

// TypeStr type a string.
func TypeStr(str string, args ...int) error { return x11.TypeStr(str, args...) }

// TypeDelay, SetDelay and CmdCtrl live in robotgo_pub.go (build-tag-free) — they
// call the package-level Type/KeySleep/MouseSleep — so they are NOT re-declared
// here.

// --- Mouse ---

// Move move the mouse to (x, y).
func Move(x, y int, displayId ...int) error { return x11.Move(x, y, displayId...) }

// MoveSmooth move the mouse smoothly to (x, y).
func MoveSmooth(x, y int, args ...any) bool { return x11.MoveSmooth(x, y, args...) }

// Click click the mouse button.
func Click(args ...any) error { return x11.Click(args...) }

// Toggle toggle the mouse button.
func Toggle(key ...any) error { return x11.Toggle(key...) }

// MouseDown send a mouse down event.
func MouseDown(key ...any) error { return x11.MouseDown(key...) }

// MouseUp send a mouse up event.
func MouseUp(key ...any) error { return x11.MouseUp(key...) }

// Scroll scroll the mouse to (x, y).
func Scroll(x, y int, args ...int) error { return x11.Scroll(x, y, args...) }

// ScrollDir scroll the mouse to a direction.
func ScrollDir(x int, direction ...any) error { return x11.ScrollDir(x, direction...) }

// ScrollSmooth scroll the mouse smoothly.
func ScrollSmooth(to int, args ...int) error { return x11.ScrollSmooth(to, args...) }

// DragSmooth drag the mouse smoothly to (x, y).
func DragSmooth(x, y int, args ...any) error { return x11.DragSmooth(x, y, args...) }

// MultiClick click the mouse button count times; click is accepted for API parity.
func MultiClick(button string, count int, click ...bool) error { return clickTimes(button, count) }

// MoveClick move and click the mouse.
func MoveClick(x, y int, args ...any) error { return x11.MoveClick(x, y, args...) }

// Location get the mouse location position, return x, y.
func Location() (int, int) { return x11.Location() }

// GetMousePos get the mouse position, return x, y.
func GetMousePos() (int, int) { return x11.GetMousePos() }

// --- Screen (capture lives in the backend; img.go provides Save*/Width/Height) ---

// GetScreenSize get the screen size.
func GetScreenSize() (int, int) { return x11.GetScreenSize() }

// GetScaleSize get the screen scale size.
func GetScaleSize(displayId ...int) (int, int) { return x11.GetScaleSize(displayId...) }

// SysScale get the sys scale (physical DPI or Xft.dpi, / 96).
func SysScale(displayId ...int) float64 { return x11.ScaleF(displayId...) }

// ScaleF get the system scale val.
func ScaleF(displayId ...int) float64 { return SysScale(displayId...) }

// Scaled get the screen scaled return scale size.
func Scaled(x int, displayId ...int) int { return Scaled0(x, ScaleF(displayId...)) }

// GetScreenRect get the screen rect (x, y, w, h).
func GetScreenRect(displayId ...int) Rect {
	r := x11.GetScreenRect(displayId...)
	return Rect{Point{r.X, r.Y}, Size{r.W, r.H}}
}

// GetMainId get the main display id.
func GetMainId() int { return x11.MainDisplayID() }

// IsMain is main display.
func IsMain(displayId int) bool { return displayId == GetMainId() }

// Deprecated: use the ScaledF(),
//
// ScaleX get the primary display horizontal DPI, drop.
func ScaleX() int { return x11.ScaleX() }

// DisplaysNum get the number of displays.
func DisplaysNum() int { return x11.DisplaysNum() }

// GetPixelColor get the pixel color at (x, y), return string.
func GetPixelColor(x, y int, displayId ...int) string { return x11.GetPixelColor(x, y, displayId...) }

// CaptureImg capture the screen, return image.Image, error.
func CaptureImg(args ...int) (image.Image, error) { return x11.CaptureImg(args...) }

// Capture1 capture the screen, return *image.RGBA, error.
//
// Named Capture1 because the portable screen.go already declares Capture().
func Capture1(args ...int) (*image.RGBA, error) { return x11.Capture(args...) }

// SaveCapture capture the screen and save it to a file.
func SaveCapture(path string, args ...int) error { return x11.SaveCapture(path, args...) }

// PadHex pad a hex color value to 6 characters.
func PadHex(hex uint32) string { return x11.PadHex(hex) }

// --- Window ---

// GetTitle get the window title, return string.
func GetTitle(args ...int) string { return x11.GetTitle(args...) }

// CheckAccess check the os accessibility, p shows the system prompt.
func CheckAccess(p bool) bool { return x11.CheckAccess(p) }

// GetActiveApp get the active app name, bundle id or path, and pid.
func GetActiveApp() (string, string, int) { return x11.GetActiveApp() }

// ActiveName active the window by name.
func ActiveName(name string) error { return x11.ActiveName(name) }

// MinWindow set the window min.
func MinWindow(pid int, args ...any) error { return x11.MinWindow(pid, args...) }

// MaxWindow set the window max.
func MaxWindow(pid int, args ...any) error { return x11.MaxWindow(pid, args...) }

// CloseWindow close the window.
func CloseWindow(args ...int) error { return x11.CloseWindow(args...) }

// ActivePid active the window by pid; with args or NotPid, pid is an X window id.
func ActivePid(pid int, args ...int) error { return x11.ActivePid(pid, isHandle(args)) }

// GetBounds get the window bounds (x, y, w, h); with args or NotPid, pid is an X window id.
func GetBounds(pid int, args ...int) (int, int, int, int) { return x11.GetBounds(pid, isHandle(args)) }

// GetClient get the window client bounds (x, y, w, h); with args or NotPid, pid is an X window id.
func GetClient(pid int, args ...int) (int, int, int, int) { return x11.GetClient(pid, isHandle(args)) }

// --- Process (Pids/Process/Kill/... come from the portable ps.go) ---

// GetPid get the current process id.
func GetPid() int { return x11.GetPid() }
