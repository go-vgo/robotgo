// Copyright 2016-2017 The go-vgo Project Developers. See the COPYRIGHT
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
	go get -u github.com/go-vgo/robotgo
*/
package robotgo

/*
//#if defined(IS_MACOSX)
	#cgo darwin CFLAGS: -x objective-c  -Wno-deprecated-declarations
	#cgo darwin LDFLAGS: -framework Cocoa -framework OpenGL -framework IOKit
	#cgo darwin LDFLAGS: -framework Carbon -framework CoreFoundation
	#cgo darwin LDFLAGS:-L${SRCDIR}/cdeps/mac -lpng -lz
//#elif defined(USE_X11)
	// Drop -std=c11
	#cgo linux CFLAGS: -I/usr/src
	#cgo linux LDFLAGS: -L/usr/src -lpng -lz -lX11 -lXtst -lX11-xcb -lxcb
	#cgo linux LDFLAGS: -lxcb-xkb -lxkbcommon -lxkbcommon-x11 -lm
//#endif
	// #cgo windows LDFLAGS: -lgdi32 -luser32 -lpng -lz
	#cgo windows LDFLAGS: -lgdi32 -luser32
	#cgo windows,amd64 LDFLAGS: -L${SRCDIR}/cdeps/win64 -lpng -lz
	#cgo windows,386 LDFLAGS: -L${SRCDIR}/cdeps/win32 -lpng -lz
// #include <AppKit/NSEvent.h>
#include "screen/goScreen.h"
#include "mouse/goMouse.h"
#include "key/goKey.h"
#include "bitmap/goBitmap.h"
#include "event/goEvent.h"
#include "window/goWindow.h"
*/
import "C"

import (
	// "fmt"
	"os"
	"reflect"
	"runtime"
	"strings"
	"time"
	"unsafe"
	// "syscall"

	"github.com/go-vgo/robotgo/clipboard"
	"github.com/shirou/gopsutil/process"
)

const (
	version string = "v0.48.0.494, Mount Cook!"
)

type (
	// Map a map
	Map map[string]interface{}
	// CHex c rgb Hex type
	CHex C.MMRGBHex
	// CBitmap c bitmap type C.MMBitmapRef
	CBitmap C.MMBitmapRef
)

// Bitmap is Bitmap struct
type Bitmap struct {
	ImageBuffer   *uint8
	Width         int
	Height        int
	Bytewidth     int
	BitsPerPixel  uint8
	BytesPerPixel uint8
}

// MPoint is MPoint struct
type MPoint struct {
	x int
	y int
}

// GetVersion get version
func GetVersion() string {
	return version
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

// Sleep time.Sleep
func Sleep(tm float64) {
	time.Sleep(time.Duration(tm) * time.Second)
}

// GoString teans C.char to string
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

// U32ToHex trans C.uint32_t to C.MMRGBHex
func U32ToHex(hex C.uint32_t) C.MMRGBHex {
	return C.MMRGBHex(hex)
}

// U8ToHex teans *C.uint8_t to C.MMRGBHex
func U8ToHex(hex *C.uint8_t) C.MMRGBHex {
	return C.MMRGBHex(*hex)
}

// PadHex trans C.MMRGBHex to string
func PadHex(hex C.MMRGBHex) string {
	color := C.pad_hex(hex)
	gcolor := C.GoString(color)
	defer C.free(unsafe.Pointer(color))

	return gcolor
}

// HexToRgb trans hex to rgb
func HexToRgb(hex uint32) *C.uint8_t {
	return C.color_hex_to_rgb(C.uint32_t(hex))
}

// RgbToHex trans rgb to hex
func RgbToHex(r, g, b uint8) C.uint32_t {
	return C.color_rgb_to_hex(C.uint8_t(r), C.uint8_t(g), C.uint8_t(b))
}

// GetPxColor get pixel color return C.MMRGBHex
func GetPxColor(x, y int) C.MMRGBHex {
	cx := C.size_t(x)
	cy := C.size_t(y)

	color := C.get_px_color(cx, cy)
	return color
}

// GetPixelColor get pixel color return string
func GetPixelColor(x, y int) string {
	cx := C.size_t(x)
	cy := C.size_t(y)

	color := C.get_pixel_color(cx, cy)
	gcolor := C.GoString(color)
	defer C.free(unsafe.Pointer(color))

	return gcolor
}

// GetScreenSize get screen size
func GetScreenSize() (int, int) {
	size := C.get_screen_size()
	// fmt.Println("...", size, size.width)
	return int(size.width), int(size.height)
}

// SetXDisplayName set XDisplay name
func SetXDisplayName(name string) string {
	cname := C.CString(name)
	str := C.set_XDisplay_name(cname)
	gstr := C.GoString(str)
	defer C.free(unsafe.Pointer(cname))

	return gstr
}

// GetXDisplayName get XDisplay name
func GetXDisplayName() string {
	name := C.get_XDisplay_name()
	gname := C.GoString(name)
	defer C.free(unsafe.Pointer(name))

	return gname
}

// CaptureScreen capture the screen return bitmap(c struct)
func CaptureScreen(args ...int) C.MMBitmapRef {
	var (
		x C.size_t
		y C.size_t
		w C.size_t
		h C.size_t
	)

	Try(func() {
		x = C.size_t(args[0])
		y = C.size_t(args[1])
		w = C.size_t(args[2])
		h = C.size_t(args[3])
	}, func(e interface{}) {
		// fmt.Println("err:::", e)
		x = 0
		y = 0
		// Get screen size.
		var displaySize C.MMSize
		displaySize = C.getMainDisplaySize()
		w = displaySize.width
		h = displaySize.height
	})

	bit := C.capture_screen(x, y, w, h)
	// fmt.Println("...", bit.width)
	return bit
}

// GoCaptureScreen capture the screen and return bitmap(go struct)
func GoCaptureScreen(args ...int) Bitmap {
	var (
		x C.size_t
		y C.size_t
		w C.size_t
		h C.size_t
	)

	Try(func() {
		x = C.size_t(args[0])
		y = C.size_t(args[1])
		w = C.size_t(args[2])
		h = C.size_t(args[3])
	}, func(e interface{}) {
		// fmt.Println("err:::", e)
		x = 0
		y = 0
		//Get screen size.
		var displaySize C.MMSize
		displaySize = C.getMainDisplaySize()
		w = displaySize.width
		h = displaySize.height
	})

	bit := C.capture_screen(x, y, w, h)
	// fmt.Println("...", bit)
	bitmap := Bitmap{
		ImageBuffer:   (*uint8)(bit.imageBuffer),
		Width:         int(bit.width),
		Height:        int(bit.height),
		Bytewidth:     int(bit.bytewidth),
		BitsPerPixel:  uint8(bit.bitsPerPixel),
		BytesPerPixel: uint8(bit.bytesPerPixel),
	}

	return bitmap
}

// BCaptureScreen capture the screen and return bitmap(go struct), Wno-deprecated
func BCaptureScreen(args ...int) Bitmap {
	var (
		x C.size_t
		y C.size_t
		w C.size_t
		h C.size_t
	)

	Try(func() {
		x = C.size_t(args[0])
		y = C.size_t(args[1])
		w = C.size_t(args[2])
		h = C.size_t(args[3])
	}, func(e interface{}) {
		// fmt.Println("err:::", e)
		x = 0
		y = 0
		//Get screen size.
		var displaySize C.MMSize
		displaySize = C.getMainDisplaySize()
		w = displaySize.width
		h = displaySize.height
	})

	bit := C.capture_screen(x, y, w, h)
	// fmt.Println("...", bit)
	bitmap := Bitmap{
		ImageBuffer:   (*uint8)(bit.imageBuffer),
		Width:         int(bit.width),
		Height:        int(bit.height),
		Bytewidth:     int(bit.bytewidth),
		BitsPerPixel:  uint8(bit.bitsPerPixel),
		BytesPerPixel: uint8(bit.bytesPerPixel),
	}

	return bitmap
}

// SaveCapture capture screen and save
func SaveCapture(spath string, args ...int) {
	var bit C.MMBitmapRef
	if len(args) > 3 {
		var (
			x = args[0]
			y = args[1]
			w = args[2]
			h = args[3]
		)

		bit = CaptureScreen(x, y, w, h)
	} else {
		bit = CaptureScreen()
	}

	SaveBitmap(bit, spath)
}

/*
.___  ___.   ______    __    __       _______. _______
|   \/   |  /  __  \  |  |  |  |     /       ||   ____|
|  \  /  | |  |  |  | |  |  |  |    |   (----`|  |__
|  |\/|  | |  |  |  | |  |  |  |     \   \    |   __|
|  |  |  | |  `--'  | |  `--'  | .----)   |   |  |____
|__|  |__|  \______/   \______/  |_______/    |_______|

*/

// MoveMouse move the mouse
func MoveMouse(x, y int) {
	// C.size_t  int
	cx := C.size_t(x)
	cy := C.size_t(y)
	C.move_mouse(cx, cy)
}

// Move move the mouse
func Move(x, y int) {
	cx := C.size_t(x)
	cy := C.size_t(y)
	C.move_mouse(cx, cy)
}

// DragMouse drag the mouse
func DragMouse(x, y int) {
	cx := C.size_t(x)
	cy := C.size_t(y)
	C.drag_mouse(cx, cy)
}

// Drag drag the mouse
func Drag(x, y int) {
	cx := C.size_t(x)
	cy := C.size_t(y)
	C.drag_mouse(cx, cy)
}

// MoveMouseSmooth move the mouse smooth
func MoveMouseSmooth(x, y int, args ...interface{}) {
	cx := C.size_t(x)
	cy := C.size_t(y)

	var (
		mouseDelay = 10
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
		low = 5.0
		high = 500.0
	}

	C.move_mouse_smooth(cx, cy, low, high, C.int(mouseDelay))
}

// MoveSmooth move the mouse smooth
func MoveSmooth(x, y int, args ...interface{}) {
	cx := C.size_t(x)
	cy := C.size_t(y)

	var (
		mouseDelay = 10
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
		low = 5.0
		high = 500.0
	}

	C.move_mouse_smooth(cx, cy, low, high, C.int(mouseDelay))
}

// GetMousePos get mouse portion
func GetMousePos() (int, int) {
	pos := C.get_mousePos()
	// fmt.Println("pos:###", pos, pos.x, pos.y)
	x := int(pos.x)
	y := int(pos.y)
	// return pos.x, pos.y
	return x, y
}

// MouseClick click the mouse
func MouseClick(args ...interface{}) {
	var (
		button C.MMMouseButton
		double C.bool
	)

	Try(func() {
		// button = args[0].(C.MMMouseButton)
		if args[0].(string) == "left" {
			button = C.LEFT_BUTTON
		}
		if args[0].(string) == "center" {
			button = C.CENTER_BUTTON
		}
		if args[0].(string) == "right" {
			button = C.RIGHT_BUTTON
		}
		double = C.bool(args[1].(bool))
	}, func(e interface{}) {
		// fmt.Println("err:::", e)
		button = C.LEFT_BUTTON
		double = false
	})

	C.mouse_click(button, double)
}

// Click click the mouse
func Click(args ...interface{}) {
	var (
		button C.MMMouseButton
		double C.bool
	)

	Try(func() {
		// button = args[0].(C.MMMouseButton)
		if args[0].(string) == "left" {
			button = C.LEFT_BUTTON
		}
		if args[0].(string) == "center" {
			button = C.CENTER_BUTTON
		}
		if args[0].(string) == "right" {
			button = C.RIGHT_BUTTON
		}
		double = C.bool(args[1].(bool))
	}, func(e interface{}) {
		// fmt.Println("err:::", e)
		button = C.LEFT_BUTTON
		double = false
	})

	C.mouse_click(button, double)
}

// MoveClick move and click the mouse
func MoveClick(x, y int, args ...interface{}) {
	MoveMouse(x, y)
	MouseClick(args)
}

// MovesClick move smooth and click the mouse
func MovesClick(x, y int, args ...interface{}) {
	MoveSmooth(x, y)
	MouseClick(args)
}

// MouseToggle toggle the mouse
func MouseToggle(args ...interface{}) {
	var button C.MMMouseButton

	Try(func() {
		// button = args[1].(C.MMMouseButton)
		if args[1].(string) == "left" {
			button = C.LEFT_BUTTON
		}
		if args[1].(string) == "center" {
			button = C.CENTER_BUTTON
		}
		if args[1].(string) == "right" {
			button = C.RIGHT_BUTTON
		}
	}, func(e interface{}) {
		// fmt.Println("err:::", e)
		button = C.LEFT_BUTTON
	})

	down := C.CString(args[0].(string))
	C.mouse_toggle(down, button)
	defer C.free(unsafe.Pointer(down))
}

// SetMouseDelay set mouse delay
func SetMouseDelay(x int) {
	cx := C.size_t(x)
	C.set_mouse_delay(cx)
}

// ScrollMouse scroll the mouse
func ScrollMouse(x int, y string) {
	cx := C.size_t(x)
	z := C.CString(y)
	C.scroll_mouse(cx, z)
	defer C.free(unsafe.Pointer(z))
}

/*
 __  ___  ___________    ____ .______     ______        ___      .______       _______
|  |/  / |   ____\   \  /   / |   _  \   /  __  \      /   \     |   _  \     |       \
|  '  /  |  |__   \   \/   /  |  |_)  | |  |  |  |    /  ^  \    |  |_)  |    |  .--.  |
|    <   |   __|   \_    _/   |   _  <  |  |  |  |   /  /_\  \   |      /     |  |  |  |
|  .  \  |  |____    |  |     |  |_)  | |  `--'  |  /  _____  \  |  |\  \----.|  '--'  |
|__|\__\ |_______|   |__|     |______/   \______/  /__/     \__\ | _| `._____||_______/

*/

// KeyTap tap the keyboard;
//
// See keys:
//	https://github.com/go-vgo/robotgo/blob/master/docs/keys.md
func KeyTap(args ...interface{}) {
	var (
		akey     string
		keyT     = "null"
		keyArr   []string
		num      int
		keyDelay = 10
	)
	// var ckeyArr []*C.char
	ckeyArr := make([](*_Ctype_char), 0)

	Try(func() {
		if reflect.TypeOf(args[1]) == reflect.TypeOf(keyArr) {

			keyArr = args[1].([]string)

			num = len(keyArr)

			for i := 0; i < num; i++ {
				ckeyArr = append(ckeyArr, (*C.char)(unsafe.Pointer(C.CString(keyArr[i]))))
			}

			if len(args) > 2 {
				keyDelay = args[2].(int)
			}
		} else {
			akey = args[1].(string)

			if len(args) > 2 {
				if reflect.TypeOf(args[2]) == reflect.TypeOf(akey) {
					keyT = args[2].(string)
				} else {
					keyDelay = args[2].(int)
				}
			}
		}

	}, func(e interface{}) {
		// fmt.Println("err:::", e)
		akey = "null"
		keyArr = []string{"null"}
	})
	// }()

	zkey := C.CString(args[0].(string))

	if akey == "" && len(keyArr) != 0 {
		C.key_Tap(zkey, (**_Ctype_char)(unsafe.Pointer(&ckeyArr[0])),
			C.int(num), C.int(keyDelay))
	} else {
		// zkey := C.CString(args[0])
		amod := C.CString(akey)
		amodt := C.CString(keyT)

		C.key_tap(zkey, amod, amodt, C.int(keyDelay))

		defer C.free(unsafe.Pointer(amod))
		defer C.free(unsafe.Pointer(amodt))
	}

	defer C.free(unsafe.Pointer(zkey))
}

// KeyToggle toggle the keyboard
//
// See keys:
//	https://github.com/go-vgo/robotgo/blob/master/docs/keys.md
func KeyToggle(args ...string) string {
	var (
		adown string
		amkey string
		mKeyT string
		// keyDelay = 10
	)

	Try(func() {
		adown = args[1]
		if len(args) > 2 {
			amkey = args[2]

			Try(func() {
				mKeyT = args[3]
			}, func(e interface{}) {
				// fmt.Println("err:::", e)
				mKeyT = "null"
			})
		} else {
			amkey = "null"
		}
	}, func(e interface{}) {
		// fmt.Println("err:::", e)
		adown = "null"
	})

	ckey := C.CString(args[0])
	cadown := C.CString(adown)
	camkey := C.CString(amkey)
	cmKeyT := C.CString(mKeyT)
	// defer func() {
	str := C.key_toggle(ckey, cadown, camkey, cmKeyT)
	// str := C.key_Toggle(ckey, cadown, camkey, cmKeyT, C.int(keyDelay))
	// fmt.Println(str)
	// }()
	defer C.free(unsafe.Pointer(ckey))
	defer C.free(unsafe.Pointer(cadown))
	defer C.free(unsafe.Pointer(camkey))
	defer C.free(unsafe.Pointer(cmKeyT))

	return C.GoString(str)
}

// ReadAll read string from clipboard
func ReadAll() (string, error) {
	return clipboard.ReadAll()
}

// WriteAll write string to clipboard
func WriteAll(text string) {
	clipboard.WriteAll(text)
}

// CharCodeAt char code at utf-8
func CharCodeAt(s string, n int) rune {
	i := 0
	for _, r := range s {
		if i == n {
			return r
		}
		i++
	}

	return 0
}

// TypeStr type string, support UTF-8
func TypeStr(str string) {
	for i := 0; i < len([]rune(str)); i++ {
		ustr := uint32(CharCodeAt(str, i))
		UnicodeType(ustr)
	}
}

// UnicodeType tap uint32 unicode
func UnicodeType(str uint32) {
	cstr := C.uint(str)
	C.unicodeType(cstr)
}

// TypeString type string, support unicode
func TypeString(x string) {
	cx := C.CString(x)
	C.type_string(cx)
	defer C.free(unsafe.Pointer(cx))
}

// PasteStr paste string, support UTF-8
func PasteStr(str string) {
	clipboard.WriteAll(str)
	if runtime.GOOS == "darwin" {
		KeyTap("v", "command")
	} else {
		KeyTap("v", "control")
	}
}

// TypeStrDelay type string delayed
func TypeStrDelay(str string, y int) {
	cstr := C.CString(str)
	cy := C.size_t(y)
	C.type_string_delayed(cstr, cy)
	defer C.free(unsafe.Pointer(cstr))
}

// TypeStringDelayed type string delayed, Wno-deprecated
func TypeStringDelayed(x string, y int) {
	cx := C.CString(x)
	cy := C.size_t(y)
	C.type_string_delayed(cx, cy)
	defer C.free(unsafe.Pointer(cx))
}

// SetKeyDelay set keyboard delay
func SetKeyDelay(x int) {
	C.set_keyboard_delay(C.size_t(x))
}

// SetKeyboardDelay set keyboard delay, Wno-deprecated
func SetKeyboardDelay(x int) {
	C.set_keyboard_delay(C.size_t(x))
}

/*
.______    __  .___________..___  ___.      ___      .______
|   _  \  |  | |           ||   \/   |     /   \     |   _  \
|  |_)  | |  | `---|  |----`|  \  /  |    /  ^  \    |  |_)  |
|   _  <  |  |     |  |     |  |\/|  |   /  /_\  \   |   ___/
|  |_)  | |  |     |  |     |  |  |  |  /  _____  \  |  |
|______/  |__|     |__|     |__|  |__| /__/     \__\ | _|
*/

// FindBitmap find the bitmap
func FindBitmap(args ...interface{}) (int, int) {
	var (
		bit       C.MMBitmapRef
		sbit      C.MMBitmapRef
		tolerance C.float
	)

	bit = args[0].(C.MMBitmapRef)
	if len(args) > 1 {
		sbit = args[1].(C.MMBitmapRef)
	} else {
		sbit = CaptureScreen()
	}

	if len(args) > 2 {
		tolerance = C.float(args[2].(float64))
	} else {
		tolerance = 0.5
	}

	pos := C.find_bitmap(bit, sbit, tolerance)
	// fmt.Println("pos----", pos)
	return int(pos.x), int(pos.y)
}

// FindEveryBitmap find the every bitmap
func FindEveryBitmap(args ...interface{}) (int, int) {
	var (
		bit       C.MMBitmapRef
		sbit      C.MMBitmapRef
		tolerance C.float
		lpos      C.MMPoint
	)

	bit = args[0].(C.MMBitmapRef)
	if len(args) > 1 {
		sbit = args[1].(C.MMBitmapRef)
	} else {
		sbit = CaptureScreen()
	}

	if len(args) > 2 {
		tolerance = C.float(args[2].(float64))
	} else {
		tolerance = 0.5
	}

	if len(args) > 3 {
		lpos.x = C.size_t(args[3].(int))
		lpos.y = 0
	} else {
		lpos.x = 0
		lpos.y = 0
	}

	if len(args) > 4 {
		lpos.x = C.size_t(args[3].(int))
		lpos.y = C.size_t(args[4].(int))
	}

	pos := C.find_every_bitmap(bit, sbit, tolerance, &lpos)
	// fmt.Println("pos----", pos)
	return int(pos.x), int(pos.y)
}

// CountBitmap count of the bitmap
func CountBitmap(bitmap C.MMBitmapRef, sbit C.MMBitmapRef, args ...float32) int {
	var tolerance C.float
	if len(args) > 0 {
		tolerance = C.float(args[0])
	} else {
		tolerance = 0.5
	}

	count := C.count_of_bitmap(bitmap, sbit, tolerance)
	return int(count)
}

// FindBit find the bitmap, Wno-deprecated
func FindBit(args ...interface{}) (int, int) {
	var bit C.MMBitmapRef
	bit = args[0].(C.MMBitmapRef)

	var rect C.MMRect
	Try(func() {
		rect.origin.x = C.size_t(args[1].(int))
		rect.origin.y = C.size_t(args[2].(int))
		rect.size.width = C.size_t(args[3].(int))
		rect.size.height = C.size_t(args[4].(int))
	}, func(e interface{}) {
		// fmt.Println("err:::", e)
		// rect.origin.x = x
		// rect.origin.y = y
		// rect.size.width = w
		// rect.size.height = h
	})

	pos := C.aFindBitmap(bit, rect)
	// fmt.Println("pos----", pos)
	return int(pos.x), int(pos.y)
}

// BitmapClick find the bitmap and click
func BitmapClick(bitmap C.MMBitmapRef, args ...interface{}) {
	x, y := FindBitmap(bitmap)
	MovesClick(x, y, args)
}

// PointInBounds bitmap point in bounds
func PointInBounds(bitmap C.MMBitmapRef, x, y int) bool {
	var point C.MMPoint
	point.x = C.size_t(x)
	point.y = C.size_t(y)
	cbool := C.point_in_bounds(bitmap, point)

	return bool(cbool)
}

// OpenBitmap open the bitmap
func OpenBitmap(args ...interface{}) C.MMBitmapRef {
	path := C.CString(args[0].(string))
	var mtype C.uint16_t

	Try(func() {
		mtype = C.uint16_t(args[1].(int))
	}, func(e interface{}) {
		// fmt.Println("err:::", e)
		mtype = 1
	})

	bit := C.bitmap_open(path, mtype)
	defer C.free(unsafe.Pointer(path))
	// fmt.Println("opening...", bit)
	return bit
	// defer C.free(unsafe.Pointer(path))
}

// BitmapStr bitmap from string
func BitmapStr(str string) C.MMBitmapRef {
	cs := C.CString(str)
	bit := C.bitmap_from_string(cs)
	defer C.free(unsafe.Pointer(cs))

	return bit
}

// SaveBitmap save the bitmap
func SaveBitmap(args ...interface{}) string {
	var mtype C.uint16_t
	Try(func() {
		mtype = C.uint16_t(args[2].(int))
	}, func(e interface{}) {
		// fmt.Println("err:::", e)
		mtype = 1
	})

	path := C.CString(args[1].(string))
	savebit := C.bitmap_save(args[0].(C.MMBitmapRef), path, mtype)
	// fmt.Println("saved...", savebit)
	// return bit
	defer C.free(unsafe.Pointer(path))

	return C.GoString(savebit)
}

// func SaveBitmap(bit C.MMBitmapRef, gpath string, mtype C.MMImageType) {
// 	path := C.CString(gpath)
// 	savebit := C.aSaveBitmap(bit, path, mtype)
// 	fmt.Println("saving...", savebit)
// 	// return bit
// 	// defer C.free(unsafe.Pointer(path))
// }

// TostringBitmap tostring bitmap to string
func TostringBitmap(bit C.MMBitmapRef) string {
	strBit := C.tostring_bitmap(bit)
	return C.GoString(strBit)
}

// TocharBitmap tostring bitmap to C.char
func TocharBitmap(bit C.MMBitmapRef) *C.char {
	strBit := C.tostring_bitmap(bit)
	return strBit
}

// ToBitmap trans C.MMBitmapRef to Bitmap
func ToBitmap(bit C.MMBitmapRef) Bitmap {
	bitmap := Bitmap{
		ImageBuffer:   (*uint8)(bit.imageBuffer),
		Width:         int(bit.width),
		Height:        int(bit.height),
		Bytewidth:     int(bit.bytewidth),
		BitsPerPixel:  uint8(bit.bitsPerPixel),
		BytesPerPixel: uint8(bit.bytesPerPixel),
	}

	return bitmap
}

// ToMMBitmapRef trans CBitmap to C.MMBitmapRef
func ToMMBitmapRef(bit CBitmap) C.MMBitmapRef {
	return C.MMBitmapRef(bit)
}

// GetPortion get bitmap portion
func GetPortion(bit C.MMBitmapRef, x, y, w, h C.size_t) C.MMBitmapRef {
	var rect C.MMRect
	rect.origin.x = x
	rect.origin.y = y
	rect.size.width = w
	rect.size.height = h

	pos := C.get_portion(bit, rect)
	return pos
}

// Convert convert bitmap
func Convert(args ...interface{}) {
	var mtype int
	Try(func() {
		mtype = args[2].(int)
	}, func(e interface{}) {
		// fmt.Println("err:::", e)
		mtype = 1
	})

	// C.CString()
	opath := args[0].(string)
	spath := args[1].(string)
	bitmap := OpenBitmap(opath)
	// fmt.Println("a----", bit_map)
	SaveBitmap(bitmap, spath, mtype)
}

// FreeBitmap free and dealloc bitmap
func FreeBitmap(bitmap C.MMBitmapRef) {
	// C.destroyMMBitmap(bitmap)
	C.bitmap_dealloc(bitmap)
}

// ReadBitmap returns false and sets error if |bitmap| is NULL
func ReadBitmap(bitmap C.MMBitmapRef) bool {
	abool := C.bitmap_ready(bitmap)
	gbool := bool(abool)
	return gbool
}

// CopyBitpb copy bitmap to pasteboard
func CopyBitpb(bitmap C.MMBitmapRef) bool {
	abool := C.bitmap_copy_to_pboard(bitmap)
	gbool := bool(abool)
	return gbool
}

// DeepCopyBit deep copy bitmap
func DeepCopyBit(bitmap C.MMBitmapRef) C.MMBitmapRef {
	bit := C.bitmap_deepcopy(bitmap)
	return bit
}

// GetColor get bitmap color
func GetColor(bitmap C.MMBitmapRef, x, y int) C.MMRGBHex {
	color := C.bitmap_get_color(bitmap, C.size_t(x), C.size_t(y))

	return color
}

// FindColor find bitmap color
func FindColor(bitmap C.MMBitmapRef, color CHex, args ...float32) (int, int) {
	var tolerance C.float

	if len(args) > 0 {
		tolerance = C.float(args[0])
	} else {
		tolerance = 0.5
	}

	pos := C.bitmap_find_color(bitmap, C.MMRGBHex(color), tolerance)
	x := int(pos.x)
	y := int(pos.y)

	return x, y
}

// FindColorCS findcolor by CaptureScreen
func FindColorCS(x, y, w, h int, color CHex, args ...float32) (int, int) {
	var tolerance float32

	if len(args) > 0 {
		tolerance = args[0]
	} else {
		tolerance = 0.5
	}

	bitmap := CaptureScreen(x, y, w, h)
	rx, ry := FindColor(bitmap, color, tolerance)
	return rx, ry
}

// CountColor count bitmap color
func CountColor(bitmap C.MMBitmapRef, color CHex, args ...float32) int {
	var tolerance C.float

	if len(args) > 0 {
		tolerance = C.float(args[0])
	} else {
		tolerance = 0.5
	}

	count := C.bitmap_count_of_color(bitmap, C.MMRGBHex(color), tolerance)

	return int(count)
}

// CountColorCS count bitmap color by CaptureScreen
func CountColorCS(x, y, w, h int, color CHex, args ...float32) int {
	var tolerance float32

	if len(args) > 0 {
		tolerance = args[0]
	} else {
		tolerance = 0.5
	}

	bitmap := CaptureScreen(x, y, w, h)
	rx := CountColor(bitmap, color, tolerance)

	return rx
}

/*
 ___________    ____  _______ .__   __. .___________.
|   ____\   \  /   / |   ____||  \ |  | |           |
|  |__   \   \/   /  |  |__   |   \|  | `---|  |----`
|   __|   \      /   |   __|  |  . `  |     |  |
|  |____   \    /    |  |____ |  |\   |     |  |
|_______|   \__/     |_______||__| \__|     |__|
*/

// AddEvent add event listener
func AddEvent(aeve string) int {
	keycode := Map{
		"f1":  "59",
		"f2":  "60",
		"f3":  "61",
		"f4":  "62",
		"f5":  "63",
		"f6":  "64",
		"f7":  "65",
		"f8":  "66",
		"f9":  "67",
		"f10": "68",
		"f11": "69",
		"f12": "70",
		// more
		"tab":     "15",
		"ctrl":    "29",
		"control": "29",
		"alt":     "56",
		"shift":   "42",
		"enter":   "28",
		"command": "3675",
	}

	var (
		cs   *C.char
		keve string
		mArr = []string{"mleft", "mright", "wheelDown",
			"wheelUp", "wheelLeft", "wheelRight"}
		mouseBool bool
	)

	for i := 0; i < len(mArr); i++ {
		if aeve == mArr[i] {
			mouseBool = true
		}
	}

	if len(aeve) > 1 && !mouseBool {
		keve = keycode[aeve].(string)
		cs = C.CString(keve)
	} else {
		cs = C.CString(aeve)
	}

	// cs := C.CString(aeve)
	eve := C.add_event(cs)
	// fmt.Println("event@@", eve)
	geve := int(eve)
	defer C.free(unsafe.Pointer(cs))

	return geve
}

// StopEvent stop event listener
func StopEvent() {
	C.stop_event()
}

// LEvent add event listener, Wno-deprecated
func LEvent(aeve string) int {
	cs := C.CString(aeve)
	eve := C.add_event(cs)
	// fmt.Println("event@@", eve)
	geve := int(eve)
	defer C.free(unsafe.Pointer(cs))

	return geve
}

/*
____    __    ____  __  .__   __.  _______   ______   ____    __    ____
\   \  /  \  /   / |  | |  \ |  | |       \ /  __  \  \   \  /  \  /   /
 \   \/    \/   /  |  | |   \|  | |  .--.  |  |  |  |  \   \/    \/   /
  \            /   |  | |  . `  | |  |  |  |  |  |  |   \            /
   \    /\    /    |  | |  |\   | |  '--'  |  `--'  |    \    /\    /
    \__/  \__/     |__| |__| \__| |_______/ \______/      \__/  \__/

*/

// ShowAlert show a alert window
func ShowAlert(title, msg string, args ...string) int {
	var (
		// title         string
		// msg           string
		defaultButton string
		cancelButton  string
	)

	Try(func() {
		// title = args[0]
		// msg = args[1]
		defaultButton = args[0]
		cancelButton = args[1]
	}, func(e interface{}) {
		defaultButton = "Ok"
		cancelButton = "Cancel"
	})

	atitle := C.CString(title)
	amsg := C.CString(msg)
	adefaultButton := C.CString(defaultButton)
	acancelButton := C.CString(cancelButton)

	cbool := C.show_alert(atitle, amsg, adefaultButton, acancelButton)
	ibool := int(cbool)

	defer C.free(unsafe.Pointer(atitle))
	defer C.free(unsafe.Pointer(amsg))
	defer C.free(unsafe.Pointer(adefaultButton))
	defer C.free(unsafe.Pointer(acancelButton))

	return ibool
}

// IsValid valid the window
func IsValid() bool {
	abool := C.is_valid()
	gbool := bool(abool)
	// fmt.Println("bool---------", gbool)
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

// CloseWindow close the window
func CloseWindow() {
	C.close_window()
}

// SetHandle set the window handle
func SetHandle(hwnd int) {
	chwnd := C.uintptr(hwnd)
	C.set_handle(chwnd)
}

// GetHandle get the window handle
func GetHandle() int {
	hwnd := C.get_handle()
	ghwnd := int(hwnd)
	// fmt.Println("gethwnd---", ghwnd)
	return ghwnd
}

// GetBHandle get the window handle
func GetBHandle() int {
	hwnd := C.bget_handle()
	ghwnd := int(hwnd)
	//fmt.Println("gethwnd---", ghwnd)
	return ghwnd
}

// GetTitle get the window title
func GetTitle() string {
	title := C.get_title()
	gtittle := C.GoString(title)
	// fmt.Println("title...", gtittle)
	return gtittle
}

// GetPID get the process id
func GetPID() int {
	pid := C.get_PID()
	return int(pid)
}

// Pids get the all process id
func Pids() ([]int32, error) {
	var ret []int32
	pid, err := process.Pids()
	if err != nil {
		return ret, err
	}

	return pid, err
}

// PidExists determine whether the process exists
func PidExists(pid int32) (bool, error) {
	abool, err := process.PidExists(pid)

	return abool, err
}

// Nps process struct
type Nps struct {
	Pid  int32
	Name string
}

// Process get the all process struct
func Process() ([]Nps, error) {
	var npsArr []Nps

	pid, err := process.Pids()

	if err != nil {
		return npsArr, err
	}

	for i := 0; i < len(pid); i++ {
		nps, err := process.NewProcess(pid[i])
		if err != nil {
			return npsArr, err
		}
		names, err := nps.Name()
		if err != nil {
			return npsArr, err
		}

		np := Nps{
			pid[i],
			names,
		}

		npsArr = append(npsArr, np)
	}

	return npsArr, err
}

// FindName find the process name by the process id
func FindName(pid int32) (string, error) {
	nps, err := process.NewProcess(pid)
	if err != nil {
		return "", err
	}

	names, err := nps.Name()
	if err != nil {
		return "", err
	}

	return names, err
}

// FindNames find the all process name
func FindNames() ([]string, error) {
	var strArr []string
	pid, err := process.Pids()

	if err != nil {
		return strArr, err
	}

	for i := 0; i < len(pid); i++ {
		nps, err := process.NewProcess(pid[i])
		if err != nil {
			return strArr, err
		}
		names, err := nps.Name()
		if err != nil {
			return strArr, err
		}

		strArr = append(strArr, names)
		return strArr, err

	}

	return strArr, err
}

// FindIds find the process id by the process name
func FindIds(name string) ([]int32, error) {
	var pids []int32
	nps, err := Process()
	if err != nil {
		return pids, err
	}

	for i := 0; i < len(nps); i++ {
		psname := strings.ToLower(nps[i].Name)
		name = strings.ToLower(name)
		abool := strings.Contains(psname, name)
		if abool {
			pids = append(pids, nps[i].Pid)
		}
	}

	return pids, err
}

// ActivePID active window by PID
func ActivePID(pid int32) {
	C.active_PID(C.uintptr(pid))
}

// ActiveName active window by name
func ActiveName(name string) error {
	pids, err := FindIds(name)
	if err == nil && len(pids) > 0 {
		ActivePID(pids[0])
	}

	return err
}

// Kill kill the process by PID
func Kill(pid int) error {
	ps := os.Process{Pid: pid}
	return ps.Kill()
}
