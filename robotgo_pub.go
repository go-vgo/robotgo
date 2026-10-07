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

package robotgo

import (
	"errors"
	"fmt"
	"runtime"
	"time"
	"unicode/utf8"

	"github.com/go-vgo/robotgo/clipboard"
	"github.com/go-vgo/robotgo/pub"
)

const (
	// Version get the robotgo version
	Version = "v2.00.0.1658, MT. Baker!"
)

// GetVersion get the robotgo version
func GetVersion() string {
	return Version
}

var (
	// MouseSleep set the mouse default millisecond sleep time
	MouseSleep = 0
	// KeySleep set the key default millisecond sleep time
	KeySleep = 10

	// DisplayID set the screen display id
	DisplayID = -1

	// NotPid used the hwnd not pid in windows
	NotPid bool
	// Scale option the os screen scale
	Scale bool
)

// MilliSleep sleep tm milli second
func MilliSleep(tm int) {
	pub.MilliSleep(tm)
}

// Deprecated: use the MilliSleep(),
//
// MicroSleep sleep tm milliseconds (fractions allowed)
func MicroSleep(tm float64) {
	time.Sleep(time.Duration(tm * float64(time.Millisecond)))
}

// Sleep time.Sleep tm second
func Sleep(tm int) {
	pub.Sleep(tm)
}

// SetDelay sets the key and mouse delay
// robotgo.SetDelay(100) option the robotgo.KeySleep and robotgo.MouseSleep = d
func SetDelay(d ...int) {
	v := 10
	if len(d) > 0 {
		v = d[0]
	}

	KeySleep = v
	MouseSleep = v
	//
	pub.SetDelay(v)
}

// Map a map[string]interface{}
type Map map[string]interface{}

// Bitmap define the go Bitmap struct
//
// The common type conversion of bitmap:
//
//	https://github.com/go-vgo/robotgo/blob/master/docs/keys.md#type-conversion
type Bitmap struct {
	ImgBuf        *uint8
	Width, Height int

	Bytewidth     int
	BitsPixel     uint8
	BytesPerPixel uint8
}

// Point is point struct
type Point struct {
	X int
	Y int
}

// Size is size structure
type Size struct {
	W, H int
}

// Rect is rect structure
type Rect struct {
	Point
	Size
}

// Try handler(err)
func Try(fun func(), handler func(interface{})) {
	defer func() {
		if err := recover(); err != nil {
			handler(err)
		}
	}()
	fun()
}

// ToInterfaces convert []string to []interface{}
func ToInterfaces(fields []string) []interface{} {
	res := make([]interface{}, 0, len(fields))
	for _, s := range fields {
		res = append(res, s)
	}
	return res
}

// ToStrings convert []interface{} to []string
func ToStrings(fields []interface{}) []string {
	res := make([]string, 0, len(fields))
	for _, s := range fields {
		res = append(res, s.(string))
	}
	return res
}

// CmdCtrl If the operating system is macOS, return the key string "cmd",
// otherwise return the key string "ctrl"
func CmdCtrl() string {
	if runtime.GOOS == "darwin" {
		return "cmd"
	}
	return "ctrl" // Ctrl
}

// CmdV tap key command + v or control + v
func CmdV(pid ...int) error {
	pid1 := 0
	if len(pid) > 0 {
		pid1 = pid[0]
	}
	return KeyTap("v", pid1, CmdCtrl())
}

// Scaled0 return int(x * f)
func Scaled0(x int, f float64) int {
	return int(float64(x) * f)
}

// Scaled1 return int(x / f)
func Scaled1(x int, f float64) int {
	return int(float64(x) / f)
}

// Scaled get the screen scaled return scale size
func Scaled(x int, displayId ...int) int {
	return Scaled0(x, ScaleF(displayId...))
}

// MoveScale calculate the os scale factor x, y
func MoveScale(x, y int, displayId ...int) (int, int) {
	if Scale || runtime.GOOS == "windows" {
		f := ScaleF()
		// on Windows ScaleF takes a window handle, not a display id
		if runtime.GOOS != "windows" {
			f = ScaleF(displayId...)
		}
		x, y = Scaled1(x, f), Scaled1(y, f)
	}

	return x, y
}

// IsMain is main display
func IsMain(displayId int) bool {
	return displayId == GetMainId()
}

// GetLocationColor get the location pos's color
func GetLocationColor(displayId ...int) string {
	x, y := Location()
	return GetPixelColor(x, y, displayId...)
}

// MoveClick move and click the mouse
//
// robotgo.MoveClick(x, y int, button string, double bool)
//
// Examples:
//
//	robotgo.MouseSleep = 100
//	robotgo.MoveClick(10, 10)
func MoveClick(x, y int, args ...interface{}) error {
	if err := Move(x, y); err != nil {
		return err
	}
	MilliSleep(50)
	return Click(args...)
}

// MouseDown send mouse down event
func MouseDown(key ...interface{}) error {
	return Toggle(mouseToggleArgs(key, "down")...)
}

// MouseUp send mouse up event
func MouseUp(key ...interface{}) error {
	return Toggle(mouseToggleArgs(key, "up")...)
}

// mouseToggleArgs builds the Toggle args (button, dir, rest...) so dir always
// wins over a caller-supplied key[1]; button defaults to "left".
func mouseToggleArgs(key []interface{}, dir string) []interface{} {
	args := []interface{}{"left", dir}
	if len(key) > 0 {
		args[0] = key[0]
	}
	if len(key) > 2 {
		args = append(args, key[2:]...)
	}
	return args
}

// ScrollDir scroll the mouse with direction to (x, "up")
// supported: "up", "down", "left", "right"
//
// Examples:
//
//	robotgo.ScrollDir(10, "down")
//	robotgo.ScrollDir(10, "up")
func ScrollDir(x int, direction ...interface{}) error {
	d := "down"
	if len(direction) > 0 {
		s, ok := direction[0].(string)
		if !ok {
			return fmt.Errorf("unknown scroll direction: %v", direction[0])
		}
		d = s
	}

	switch d {
	case "down":
		return Scroll(0, -x)
	case "up":
		return Scroll(0, x)
	case "left":
		return Scroll(x, 0)
	case "right":
		return Scroll(-x, 0)
	}
	return fmt.Errorf("unknown scroll direction: %v", d)
}

// ScrollSmooth scroll the mouse smooth,
// default scroll 5 times and sleep 100 millisecond
//
// robotgo.ScrollSmooth(toy, num, sleep, tox)
//
// Examples:
//
//	robotgo.ScrollSmooth(-10)
//	robotgo.ScrollSmooth(-10, 6, 200, -10)
func ScrollSmooth(to int, args ...int) error {
	num := 5
	if len(args) > 0 {
		num = args[0]
	}
	tm := 100
	if len(args) > 1 {
		tm = args[1]
	}
	tox := 0
	if len(args) > 2 {
		tox = args[2]
	}

	for i := 0; i < num; i++ {
		if err := Scroll(tox, to); err != nil {
			return err
		}
		MilliSleep(tm)
	}
	MilliSleep(MouseSleep)
	return nil
}

// MoveArgs get the mouse relative args
func MoveArgs(x, y int) (int, int) {
	mx, my := Location()
	mx = mx + x
	my = my + y

	return mx, my
}

// ErrSmoothMove is returned when MoveSmooth reports a failure
var ErrSmoothMove = errors.New("robotgo: smooth move failed")

// MoveRelative move mouse with relative
func MoveRelative(x, y int) error {
	return Move(MoveArgs(x, y))
}

// MoveSmoothRelative move mouse smooth with relative
func MoveSmoothRelative(x, y int, args ...interface{}) error {
	mx, my := MoveArgs(x, y)
	if !MoveSmooth(mx, my, args...) {
		return ErrSmoothMove
	}
	return nil
}

// MovesClick move smooth and click the mouse
//
// use the `robotgo.MouseSleep = 100`
func MovesClick(x, y int, args ...interface{}) error {
	if !MoveSmooth(x, y) {
		return ErrSmoothMove
	}
	MilliSleep(50)
	return Click(args...)
}

// ScrollRelative scroll mouse with relative
//
// Examples:
//
//	robotgo.ScrollRelative(10, 10)
func ScrollRelative(x, y int, args ...int) error {
	mx, my := MoveArgs(x, y)
	return Scroll(mx, my, args...)
}

// ReadAll read string from clipboard
func ReadAll() (string, error) {
	return clipboard.ReadAll()
}

// WriteAll write string to clipboard
func WriteAll(text string) error {
	return clipboard.WriteAll(text)
}

// PasteStr paste a string
//
// Deprecated: use the Paste()
func PasteStr(str string) error {
	_, err := Paste(str)
	return err
}

// Pastes paste a string and return the pasted character count, 0 on failure
func Pastes(str string, pid ...int) int {
	n, err := Paste(str, pid...)
	if err != nil {
		return 0
	}
	return n
}

// Paste paste a string (supported UTF-8),
// write the string to clipboard and tap `cmd + v`,
// return the pasted character (rune) count
func Paste(str string, pid ...int) (int, error) {
	if err := clipboard.WriteAll(str); err != nil {
		return 0, err
	}
	if err := CmdV(pid...); err != nil {
		return 0, err
	}
	return utf8.RuneCountInString(str), nil
}

// alertArgs returns the Alert default and cancel button labels
func alertArgs(args ...string) (string, string) {
	defaultBtn, cancelBtn := "Ok", "Cancel"
	if len(args) > 0 {
		defaultBtn = args[0]
	}
	if len(args) > 1 {
		cancelBtn = args[1]
	}
	return defaultBtn, cancelBtn
}
