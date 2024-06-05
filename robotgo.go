// Copyright 2016 The go-vgo Project Developers. See the COPYRIGHT
// file at the top-level directory of this distribution and at
// https://github.com/go-vgo/robotgo/blob/master/LICENSE
//
// Licensed under the Apache License, Version 2.0 <LICENSE-APACHE or
// http://www.apache.org/licenses/LICENSE-2.0> or the MIT license
// <LICENSE-MIT or http://opensource.org/licenses/MIT>, at your
// option. This file may not be copied, modified, or distributed
// except according to those terms.

/*
Package robotgo Go native cross-platform system automation.

Please make sure Golang, GCC is installed correctly before installing RobotGo;

See Requirements:

	https://github.com/go-vgo/robotgo#requirements

Installation:

With Go module support (Go 1.11+), just import:

	import "github.com/go-vgo/robotgo"

Otherwise, to install the robotgo package, run the command:

	go get -u github.com/go-vgo/robotgo
*/
package keyboard

/*
#cgo darwin CFLAGS: -x objective-c -Wno-deprecated-declarations
#cgo darwin LDFLAGS: -framework Cocoa -framework OpenGL -framework IOKit
#cgo darwin LDFLAGS: -framework Carbon -framework CoreFoundation

#cgo linux CFLAGS: -I/usr/src
#cgo linux LDFLAGS: -L/usr/src -lm -lX11 -lXtst

#cgo windows LDFLAGS: -lgdi32 -luser32
//
#include "screen/goScreen.h"
#include "mouse/mouse_c.h"
#include "window/goWindow.h"
*/
import "C"

import (
	"image"
	"runtime"
	"time"
	"unsafe"

	"github.com/vcaesar/tt"
)

const (
	// Version get the robotgo version
	Version = "v1.00.0.1189, MT. Baker!"
)

// GetVersion get the robotgo version
func GetVersion() string {
	return Version
}

var (
	// MouseSleep set the mouse default millisecond sleep time
	MouseSleep = 0
	// KeySleep set the key default millisecond sleep time
	KeySleep = 0

	// DisplayID set the screen display id
	DisplayID = -1

	// NotPid used the hwnd not pid in windows
	NotPid bool
	// Scale option the os screen scale
	Scale bool
)

type (
	// Map a map[string]interface{}
	Map map[string]interface{}
	// CHex define CHex as c rgb Hex type (C.MMRGBHex)
	CHex C.MMRGBHex
	// CBitmap define CBitmap as C.MMBitmapRef type
	CBitmap C.MMBitmapRef
)

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

// MilliSleep sleep tm milli second
func MilliSleep(tm int) {
	time.Sleep(time.Duration(tm) * time.Millisecond)
}

// Sleep time.Sleep tm second
func Sleep(tm int) {
	time.Sleep(time.Duration(tm) * time.Second)
}

// Deprecated: use the MilliSleep(),
//
// MicroSleep time C.microsleep(tm)
func MicroSleep(tm float64) {
	C.microsleep(C.double(tm))
}

// GoString trans C.char to string
func GoString(char *C.char) string {
	return C.GoString(char)
}

/*
      _______.  ______ .______       _______  _______ .__   __.
    /       | /      ||   _  \     |   ____||   ____||  \ |  |
   |   (----`|  ,----'|  |_)  |    |  |__   |  |__   |   \|  |
    \   \    |  |     |      /     |   __|  |   __|  |  . `  |
.----)   |   |  `----.|  |\  \----.|  |____ |  |____ |  |\   |
|_______/     \______|| _| `._____||_______||_______||__| \__|
*/

// ToMMRGBHex trans CHex to C.MMRGBHex
func ToMMRGBHex(hex CHex) C.MMRGBHex {
	return C.MMRGBHex(hex)
}

// UintToHex trans uint32 to robotgo.CHex
func UintToHex(u uint32) CHex {
	hex := U32ToHex(C.uint32_t(u))
	return CHex(hex)
}

// U32ToHex trans C.uint32_t to C.MMRGBHex
func U32ToHex(hex C.uint32_t) C.MMRGBHex {
	return C.MMRGBHex(hex)
}

// U8ToHex trans *C.uint8_t to C.MMRGBHex
func U8ToHex(hex *C.uint8_t) C.MMRGBHex {
	return C.MMRGBHex(*hex)
}

// PadHex trans C.MMRGBHex to string
func PadHex(hex C.MMRGBHex) string {
	color := C.pad_hex(hex)
	gcolor := C.GoString(color)
	C.free(unsafe.Pointer(color))

	return gcolor
}

// PadHexs trans CHex to string
func PadHexs(hex CHex) string {
	return PadHex(C.MMRGBHex(hex))
}

// HexToRgb trans hex to rgb
func HexToRgb(hex uint32) *C.uint8_t {
	return C.color_hex_to_rgb(C.uint32_t(hex))
}

// RgbToHex trans rgb to hex
func RgbToHex(r, g, b uint8) C.uint32_t {
	return C.color_rgb_to_hex(C.uint8_t(r), C.uint8_t(g), C.uint8_t(b))
}

// GetPxColor get the pixel color return C.MMRGBHex
func GetPxColor(x, y int, displayId ...int) C.MMRGBHex {
	cx := C.int32_t(x)
	cy := C.int32_t(y)

	display := displayIdx(displayId...)
	color := C.get_px_color(cx, cy, C.int32_t(display))
	return color
}

// GetPixelColor get the pixel color return string
func GetPixelColor(x, y int, displayId ...int) string {
	return PadHex(GetPxColor(x, y, displayId...))
}

// GetLocationColor get the location pos's color
func GetLocationColor(displayId ...int) string {
	x, y := Location()
	return GetPixelColor(x, y, displayId...)
}

// IsMain is main display
func IsMain(displayId int) bool {
	return displayId == GetMainId()
}

func displayIdx(id ...int) int {
	display := -1
	if DisplayID != -1 {
		display = DisplayID
	}
	if len(id) > 0 {
		display = id[0]
	}

	return display
}

func getNumDisplays() int {
	return int(C.get_num_displays())
}

// GetHWNDByPid get the hwnd by pid
func GetHWNDByPid(pid int) int {
	return int(C.get_hwnd_by_pid(C.uintptr(pid)))
}

// SysScale get the sys scale
func SysScale(displayId ...int) float64 {
	display := displayIdx(displayId...)
	s := C.sys_scale(C.int32_t(display))
	return float64(s)
}

// Scaled get the screen scaled return scale size
func Scaled(x int, displayId ...int) int {
	f := ScaleF(displayId...)
	return Scaled0(x, f)
}

// Scaled0 return int(x * f)
func Scaled0(x int, f float64) int {
	return int(float64(x) * f)
}

// Scaled1 return int(x / f)
func Scaled1(x int, f float64) int {
	return int(float64(x) / f)
}

// GetScreenSize get the screen size
func GetScreenSize() (int, int) {
	size := C.getMainDisplaySize()
	return int(size.w), int(size.h)
}

// GetScreenRect get the screen rect (x, y, w, h)
func GetScreenRect(displayId ...int) Rect {
	display := -1
	if len(displayId) > 0 {
		display = displayId[0]
	}

	rect := C.getScreenRect(C.int32_t(display))
	x, y, w, h := int(rect.origin.x), int(rect.origin.y),
		int(rect.size.w), int(rect.size.h)

	if runtime.GOOS == "windows" {
		// f := ScaleF(displayId...)
		f := ScaleF()
		x, y, w, h = Scaled0(x, f), Scaled0(y, f), Scaled0(w, f), Scaled0(h, f)
	}
	return Rect{
		Point{X: x, Y: y},
		Size{W: w, H: h},
	}
}

// GetScaleSize get the screen scale size
func GetScaleSize(displayId ...int) (int, int) {
	x, y := GetScreenSize()
	f := ScaleF(displayId...)
	return int(float64(x) * f), int(float64(y) * f)
}

// CaptureScreen capture the screen return bitmap(c struct),
// use `defer robotgo.FreeBitmap(bitmap)` to free the bitmap
//
// robotgo.CaptureScreen(x, y, w, h int)
func CaptureScreen(args ...int) CBitmap {
	var x, y, w, h C.int32_t
	displayId := -1
	if DisplayID != -1 {
		displayId = DisplayID
	}

	if len(args) > 4 {
		displayId = args[4]
	}

	if len(args) > 3 {
		x = C.int32_t(args[0])
		y = C.int32_t(args[1])
		w = C.int32_t(args[2])
		h = C.int32_t(args[3])
	} else {
		// Get the main screen rect.
		rect := GetScreenRect(displayId)
		if runtime.GOOS == "windows" {
			x = C.int32_t(rect.X)
			y = C.int32_t(rect.Y)
		}

		w = C.int32_t(rect.W)
		h = C.int32_t(rect.H)
	}

	isPid := 0
	if NotPid || len(args) > 5 {
		isPid = 1
	}

	bit := C.capture_screen(x, y, w, h, C.int32_t(displayId), C.int8_t(isPid))
	return CBitmap(bit)
}

// CaptureGo capture the screen and return bitmap(go struct)
func CaptureGo(args ...int) Bitmap {
	bit := CaptureScreen(args...)
	defer FreeBitmap(bit)

	return ToBitmap(bit)
}

// CaptureImg capture the screen and return image.Image
func CaptureImg(args ...int) image.Image {
	bit := CaptureScreen(args...)
	defer FreeBitmap(bit)

	return ToImage(bit)
}

// FreeBitmap free and dealloc the C bitmap
func FreeBitmap(bitmap CBitmap) {
	// C.destroyMMBitmap(bitmap)
	C.bitmap_dealloc(C.MMBitmapRef(bitmap))
}

// FreeBitmapArr free and dealloc the C bitmap array
func FreeBitmapArr(bit ...CBitmap) {
	for i := 0; i < len(bit); i++ {
		FreeBitmap(bit[i])
	}
}

// ToMMBitmapRef trans CBitmap to C.MMBitmapRef
func ToMMBitmapRef(bit CBitmap) C.MMBitmapRef {
	return C.MMBitmapRef(bit)
}

// ToBitmap trans C.MMBitmapRef to Bitmap
func ToBitmap(bit CBitmap) Bitmap {
	bitmap := Bitmap{
		ImgBuf:        (*uint8)(bit.imageBuffer),
		Width:         int(bit.width),
		Height:        int(bit.height),
		Bytewidth:     int(bit.bytewidth),
		BitsPixel:     uint8(bit.bitsPerPixel),
		BytesPerPixel: uint8(bit.bytesPerPixel),
	}

	return bitmap
}

// ToCBitmap trans Bitmap to C.MMBitmapRef
func ToCBitmap(bit Bitmap) CBitmap {
	cbitmap := C.createMMBitmap_c(
		(*C.uint8_t)(bit.ImgBuf),
		C.int32_t(bit.Width),
		C.int32_t(bit.Height),
		C.int32_t(bit.Bytewidth),
		C.uint8_t(bit.BitsPixel),
		C.uint8_t(bit.BytesPerPixel),
	)

	return CBitmap(cbitmap)
}

// ToImage convert C.MMBitmapRef to standard image.Image
func ToImage(bit CBitmap) image.Image {
	return ToRGBA(bit)
}

// ToRGBA convert C.MMBitmapRef to standard image.RGBA
func ToRGBA(bit CBitmap) *image.RGBA {
	bmp1 := ToBitmap(bit)
	return ToRGBAGo(bmp1)
}

// ImgToCBitmap trans image.Image to CBitmap
func ImgToCBitmap(img image.Image) CBitmap {
	return ToCBitmap(ImgToBitmap(img))
}

// ByteToCBitmap trans []byte to CBitmap
func ByteToCBitmap(by []byte) CBitmap {
	img, _ := ByteToImg(by)
	return ImgToCBitmap(img)
}

// SetXDisplayName set XDisplay name (Linux)
func SetXDisplayName(name string) error {
	cname := C.CString(name)
	str := C.set_XDisplay_name(cname)
	C.free(unsafe.Pointer(cname))

	return toErr(str)
}

// GetXDisplayName get XDisplay name (Linux)
func GetXDisplayName() string {
	name := C.get_XDisplay_name()
	gname := C.GoString(name)
	C.free(unsafe.Pointer(name))

	return gname
}

// Deprecated: use the ScaledF(),
//
// ScaleX get the primary display horizontal DPI scale factor, drop
func ScaleX() int {
	return int(C.scaleX())
}

/*
.___  ___.   ______    __    __       _______. _______
|   \/   |  /  __  \  |  |  |  |     /       ||   ____|
|  \  /  | |  |  |  | |  |  |  |    |   (----`|  |__
|  |\/|  | |  |  |  | |  |  |  |     \   \    |   __|
|  |  |  | |  `--'  | |  `--'  | .----)   |   |  |____
|__|  |__|  \______/   \______/  |_______/    |_______|

*/

// CheckMouse check the mouse button
func CheckMouse(btn string) C.MMMouseButton {
	// button = args[0].(C.MMMouseButton)
	m1 := map[string]C.MMMouseButton{
		"left":       C.LEFT_BUTTON,
		"center":     C.CENTER_BUTTON,
		"right":      C.RIGHT_BUTTON,
		"wheelDown":  C.WheelDown,
		"wheelUp":    C.WheelUp,
		"wheelLeft":  C.WheelLeft,
		"wheelRight": C.WheelRight,
	}
	if v, ok := m1[btn]; ok {
		return v
	}

	return C.LEFT_BUTTON
}

// MoveScale calculate the os scale factor x, y
func MoveScale(x, y int, displayId ...int) (int, int) {
	if Scale && runtime.GOOS == "windows" {
		f := ScaleF()
		x, y = Scaled1(x, f), Scaled1(y, f)
	}

	return x, y
}

// Move move the mouse to (x, y)
//
// Examples:
//
//	robotgo.MouseSleep = 100  // 100 millisecond
//	robotgo.Move(10, 10)
func Move(x, y int, displayId ...int) {
	x, y = MoveScale(x, y, displayId...)

	cx := C.int32_t(x)
	cy := C.int32_t(y)
	C.moveMouse(C.MMPointInt32Make(cx, cy))

	MilliSleep(MouseSleep)
}

// Deprecated: use the DragSmooth(),
//
// Drag drag the mouse to (x, y),
// It's not valid now, use the DragSmooth()
func Drag(x, y int, args ...string) {
	x, y = MoveScale(x, y)

	var button C.MMMouseButton = C.LEFT_BUTTON
	cx := C.int32_t(x)
	cy := C.int32_t(y)

	if len(args) > 0 {
		button = CheckMouse(args[0])
	}

	C.dragMouse(C.MMPointInt32Make(cx, cy), button)
	MilliSleep(MouseSleep)
}

// DragSmooth drag the mouse like smooth to (x, y)
//
// Examples:
//
//	robotgo.DragSmooth(10, 10)
func DragSmooth(x, y int, args ...interface{}) {
	Toggle("left")
	MilliSleep(50)
	MoveSmooth(x, y, args...)
	Toggle("left", "up")
}

// MoveSmooth move the mouse smooth,
// moves mouse to x, y human like, with the mouse button up.
//
// robotgo.MoveSmooth(x, y int, low, high float64, mouseDelay int)
//
// Examples:
//
//	robotgo.MoveSmooth(10, 10)
//	robotgo.MoveSmooth(10, 10, 1.0, 2.0)
func MoveSmooth(x, y int, args ...interface{}) bool {
	// if runtime.GOOS == "windows" {
	// 	f := ScaleF()
	// 	x, y = Scaled0(x, f), Scaled0(y, f)
	// }
	x, y = MoveScale(x, y)

	cx := C.int32_t(x)
	cy := C.int32_t(y)

	var (
		mouseDelay = 1
		low        C.double
		high       C.double
	)

	if len(args) > 2 {
		mouseDelay = args[2].(int)
	}

	if len(args) > 1 {
		low = C.double(args[0].(float64))
		high = C.double(args[1].(float64))
	} else {
		low = 1.0
		high = 3.0
	}

	cbool := C.smoothlyMoveMouse(C.MMPointInt32Make(cx, cy), low, high)
	MilliSleep(MouseSleep + mouseDelay)

	return bool(cbool)
}

// MoveArgs get the mouse relative args
func MoveArgs(x, y int) (int, int) {
	mx, my := Location()
	mx = mx + x
	my = my + y

	return mx, my
}

// MoveRelative move mouse with relative
func MoveRelative(x, y int) {
	Move(MoveArgs(x, y))
}

// MoveSmoothRelative move mouse smooth with relative
func MoveSmoothRelative(x, y int, args ...interface{}) {
	mx, my := MoveArgs(x, y)
	MoveSmooth(mx, my, args...)
}

// Location get the mouse location position return x, y
func Location() (int, int) {
	pos := C.location()
	x := int(pos.x)
	y := int(pos.y)

	if runtime.GOOS == "windows" {
		f := ScaleF()
		x, y = Scaled0(x, f), Scaled0(y, f)
	}

	return x, y
}

// Click click the mouse button
//
// robotgo.Click(button string, double bool)
//
// Examples:
//
//	robotgo.Click() // default is left button
//	robotgo.Click("right")
//	robotgo.Click("wheelLeft")
func Click(args ...interface{}) {
	var (
		button C.MMMouseButton = C.LEFT_BUTTON
		double bool
	)

	if len(args) > 0 {
		button = CheckMouse(args[0].(string))
	}

	if len(args) > 1 {
		double = args[1].(bool)
	}

	if !double {
		C.clickMouse(button)
	} else {
		C.doubleClick(button)
	}

	MilliSleep(MouseSleep)
}

// MoveClick move and click the mouse
//
// robotgo.MoveClick(x, y int, button string, double bool)
//
// Examples:
//
//	robotgo.MouseSleep = 100
//	robotgo.MoveClick(10, 10)
func MoveClick(x, y int, args ...interface{}) {
	Move(x, y)
	MilliSleep(50)
	Click(args...)
}

// MovesClick move smooth and click the mouse
//
// use the `robotgo.MouseSleep = 100`
func MovesClick(x, y int, args ...interface{}) {
	MoveSmooth(x, y)
	MilliSleep(50)
	Click(args...)
}

// Toggle toggle the mouse, support button:
//
//		"left", "center", "right",
//	 "wheelDown", "wheelUp", "wheelLeft", "wheelRight"
//
// Examples:
//
//	robotgo.Toggle("left") // default is down
//	robotgo.Toggle("left", "up")
func Toggle(key ...interface{}) error {
	var button C.MMMouseButton = C.LEFT_BUTTON
	if len(key) > 0 {
		button = CheckMouse(key[0].(string))
	}

	down := true
	if len(key) > 1 && key[1].(string) == "up" {
		down = false
	}
	C.toggleMouse(C.bool(down), button)
	if len(key) > 2 {
		MilliSleep(MouseSleep)
	}

	return nil
}

// MouseDown send mouse down event
func MouseDown(key ...interface{}) error {
	return Toggle(key...)
}

// MouseUp send mouse up event
func MouseUp(key ...interface{}) error {
	if len(key) <= 0 {
		key = append(key, "left")
	}
	return Toggle(append(key, "up")...)
}

// Scroll scroll the mouse to (x, y)
//
// robotgo.Scroll(x, y, msDelay int)
//
// Examples:
//
//	robotgo.Scroll(10, 10)
func Scroll(x, y int, args ...int) {
	var msDelay = 10
	if len(args) > 0 {
		msDelay = args[0]
	}

	cx := C.int(x)
	cy := C.int(y)

	C.scrollMouseXY(cx, cy)
	MilliSleep(MouseSleep + msDelay)
}

// ScrollDir scroll the mouse with direction to (x, "up")
// supported: "up", "down", "left", "right"
//
// Examples:
//
//	robotgo.ScrollDir(10, "down")
//	robotgo.ScrollDir(10, "up")
func ScrollDir(x int, direction ...interface{}) {
	d := "down"
	if len(direction) > 0 {
		d = direction[0].(string)
	}

	if d == "down" {
		Scroll(0, -x)
	}
	if d == "up" {
		Scroll(0, x)
	}

	if d == "left" {
		Scroll(x, 0)
	}
	if d == "right" {
		Scroll(-x, 0)
	}
	// MilliSleep(MouseSleep)
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
func ScrollSmooth(to int, args ...int) {
	i := 0
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

	for {
		Scroll(tox, to)
		MilliSleep(tm)
		i++
		if i == num {
			break
		}
	}
	MilliSleep(MouseSleep)
}

// ScrollRelative scroll mouse with relative
//
// Examples:
//
//	robotgo.ScrollRelative(10, 10)
func ScrollRelative(x, y int, args ...int) {
	mx, my := MoveArgs(x, y)
	Scroll(mx, my, args...)
}

/*
____    __    ____  __  .__   __.  _______   ______   ____    __    ____
\   \  /  \  /   / |  | |  \ |  | |       \ /  __  \  \   \  /  \  /   /
 \   \/    \/   /  |  | |   \|  | |  .--.  |  |  |  |  \   \/    \/   /
  \            /   |  | |  . `  | |  |  |  |  |  |  |   \            /
   \    /\    /    |  | |  |\   | |  '--'  |  `--'  |    \    /\    /
    \__/  \__/     |__| |__| \__| |_______/ \______/      \__/  \__/

*/

func alertArgs(args ...string) (string, string) {
	var (
		defaultBtn = "Ok"
		cancelBtn  = "Cancel"
	)

	if len(args) > 0 {
		defaultBtn = args[0]
	}

	if len(args) > 1 {
		cancelBtn = args[1]
	}

	return defaultBtn, cancelBtn
}

func showAlert(title, msg string, args ...string) bool {
	defaultBtn, cancelBtn := alertArgs(args...)

	cTitle := C.CString(title)
	cMsg := C.CString(msg)
	defaultButton := C.CString(defaultBtn)
	cancelButton := C.CString(cancelBtn)

	cbool := C.showAlert(cTitle, cMsg, defaultButton, cancelButton)
	ibool := int(cbool)

	C.free(unsafe.Pointer(cTitle))
	C.free(unsafe.Pointer(cMsg))
	C.free(unsafe.Pointer(defaultButton))
	C.free(unsafe.Pointer(cancelButton))

	return ibool == 0
}

// IsValid valid the window
func IsValid() bool {
	abool := C.is_valid()
	gbool := bool(abool)
	return gbool
}

// SetActive set the window active
func SetActive(win C.MData) {
	C.set_active(win)
}

// GetActive get the active window
func GetActive() C.MData {
	mdata := C.get_active()
	// fmt.Println("active----", mdata)
	return mdata
}

// MinWindow set the window min
func MinWindow(pid int, args ...interface{}) {
	var (
		state = true
		isPid int
	)

	if len(args) > 0 {
		state = args[0].(bool)
	}
	if len(args) > 1 || NotPid {
		isPid = 1
	}

	C.min_window(C.uintptr(pid), C.bool(state), C.int8_t(isPid))
}

// MaxWindow set the window max
func MaxWindow(pid int, args ...interface{}) {
	var (
		state = true
		isPid int
	)

	if len(args) > 0 {
		state = args[0].(bool)
	}
	if len(args) > 1 || NotPid {
		isPid = 1
	}

	C.max_window(C.uintptr(pid), C.bool(state), C.int8_t(isPid))
}

// CloseWindow close the window
func CloseWindow(args ...int) {
	if len(args) <= 0 {
		C.close_main_window()
		return
	}

	var pid, isPid int
	if len(args) > 0 {
		pid = args[0]
	}
	if len(args) > 1 || NotPid {
		isPid = 1
	}

	C.close_window_by_PId(C.uintptr(pid), C.int8_t(isPid))
}

// SetHandle set the window handle
func SetHandle(hwnd int) {
	chwnd := C.uintptr(hwnd)
	C.setHandle(chwnd)
}

// SetHandlePid set the window handle by pid
func SetHandlePid(pid int, args ...int) {
	var isPid int
	if len(args) > 0 || NotPid {
		isPid = 1
	}

	C.set_handle_pid_mData(C.uintptr(pid), C.int8_t(isPid))
}

// GetHandPid get handle mdata by pid
func GetHandPid(pid int, args ...int) C.MData {
	var isPid int
	if len(args) > 0 || NotPid {
		isPid = 1
	}

	return C.set_handle_pid(C.uintptr(pid), C.int8_t(isPid))
}

// GetHandle get the window handle
func GetHandle() int {
	hwnd := C.get_handle()
	ghwnd := int(hwnd)
	// fmt.Println("gethwnd---", ghwnd)
	return ghwnd
}

// Deprecated: use the GetHandle(),
//
// # GetBHandle get the window handle, Wno-deprecated
//
// This function will be removed in version v1.0.0
func GetBHandle() int {
	tt.Drop("GetBHandle", "GetHandle")
	hwnd := C.b_get_handle()
	ghwnd := int(hwnd)
	//fmt.Println("gethwnd---", ghwnd)
	return ghwnd
}

func cgetTitle(pid, isPid int) string {
	title := C.get_title_by_pid(C.uintptr(pid), C.int8_t(isPid))
	gtitle := C.GoString(title)

	return gtitle
}

// GetTitle get the window title return string
//
// Examples:
//
//	fmt.Println(robotgo.GetTitle())
//
//	ids, _ := robotgo.FindIds()
//	robotgo.GetTitle(ids[0])
func GetTitle(args ...int) string {
	if len(args) <= 0 {
		title := C.get_main_title()
		gtitle := C.GoString(title)
		return gtitle
	}

	if len(args) > 1 {
		return internalGetTitle(args[0], args[1])
	}

	return internalGetTitle(args[0])
}

// GetPid get the process id return int32
func GetPid() int {
	pid := C.get_PID()
	return int(pid)
}

// internalGetBounds get the window bounds
func internalGetBounds(pid, isPid int) (int, int, int, int) {
	bounds := C.get_bounds(C.uintptr(pid), C.int8_t(isPid))
	return int(bounds.X), int(bounds.Y), int(bounds.W), int(bounds.H)
}

// internalGetClient get the window client bounds
func internalGetClient(pid, isPid int) (int, int, int, int) {
	bounds := C.get_client(C.uintptr(pid), C.int8_t(isPid))
	return int(bounds.X), int(bounds.Y), int(bounds.W), int(bounds.H)
}

// Is64Bit determine whether the sys is 64bit
func Is64Bit() bool {
	b := C.Is64Bit()
	return bool(b)
}

func internalActive(pid, isPid int) {
	C.active_PID(C.uintptr(pid), C.int8_t(isPid))
}

// ActivePid active the window by Pid,
// If args[0] > 0 on the Windows platform via a window handle to active
// func ActivePid(pid int32, args ...int) {
// 	var isPid int
// 	if len(args) > 0 {
// 		isPid = args[0]
// 	}

// 	C.active_PID(C.uintptr(pid), C.uintptr(isPid))
// }

// ActiveName active the window by name
//
// Examples:
//
//	robotgo.ActiveName("chrome")
func ActiveName(name string) error {
	pids, err := FindIds(name)
	if err == nil && len(pids) > 0 {
		return ActivePid(pids[0])
	}

	return err
}
