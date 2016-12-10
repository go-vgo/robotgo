// Copyright 2016 The go-vgo Project Developers.

// All Rights Reserved. See the COPYRIGHT file at
// the top-level directory of this distribution and at
// https://github.com/go-vgo/robotgo/blob/master/LICENSE
//
// Licensed under the Apache License, Version 2.0 <LICENSE-APACHE or
// http://www.apache.org/licenses/LICENSE-2.0> or the MIT license
// <LICENSE-MIT or http://opensource.org/licenses/MIT>, at your
// option. This file may not be copied, modified, or distributed
// except according to those terms.

package robotgo

/*
//#if defined(IS_MACOSX)
	#cgo darwin CFLAGS: -x objective-c  -Wno-deprecated-declarations
	#cgo darwin LDFLAGS: -framework Cocoa -framework OpenGL -framework IOKit -framework Carbon -framework CoreFoundation
//#elif defined(USE_X11)
	#cgo linux CFLAGS:-I/usr/src
	#cgo linux LDFLAGS:-L/usr/src -lpng -lz -lX11 -lXtst -lm
//#endif
	#cgo windows LDFLAGS: -lgdi32 -luser32
//#include <AppKit/NSEvent.h>
#include "screen/goScreen.h"
#include "mouse/goMouse.h"
#include "key/goKey.h"
//#include "bitmap/goBitmap.h"
//#include "event/goEvent.h"
//#include "window/goWindow.h"
*/
import "C"

import (
	. "fmt"
	"unsafe"
	// "runtime"
	// "syscall"
)

/*
      _______.  ______ .______       _______  _______ .__   __.
    /       | /      ||   _  \     |   ____||   ____||  \ |  |
   |   (----`|  ,----'|  |_)  |    |  |__   |  |__   |   \|  |
    \   \    |  |     |      /     |   __|  |   __|  |  . `  |
.----)   |   |  `----.|  |\  \----.|  |____ |  |____ |  |\   |
|_______/     \______|| _| `._____||_______||_______||__| \__|
*/

type Bit_map struct {
	ImageBuffer   *C.uint8_t
	Width         C.size_t
	Height        C.size_t
	Bytewidth     C.size_t
	BitsPerPixel  C.uint8_t
	BytesPerPixel C.uint8_t
}

func GetPixelColor(x, y C.size_t) string {
	color := C.aGetPixelColor(x, y)
	gcolor := C.GoString(color)
	defer C.free(unsafe.Pointer(color))
	return gcolor
}

func GetScreenSize() (C.size_t, C.size_t) {
	size := C.aGetScreenSize()
	// Println("...", size, size.width)
	return size.width, size.height
}

func GetXDisplayName() string {
	name := C.aGetXDisplayName()
	gname := C.GoString(name)
	defer C.free(unsafe.Pointer(name))
	return gname
}

func SetXDisplayName(name string) string {
	cname := C.CString(name)
	str := C.aSetXDisplayName(cname)
	gstr := C.GoString(str)
	return gstr
}

func CaptureScreen(x, y, w, h C.int) C.MMBitmapRef {
	bit := C.aCaptureScreen(x, y, w, h)
	// Println("...", bit.width)
	return bit
}

func Capture_Screen(x, y, w, h C.int) Bit_map {
	bit := C.aCaptureScreen(x, y, w, h)
	// Println("...", bit)
	bit_map := Bit_map{
		ImageBuffer:   bit.imageBuffer,
		Width:         bit.width,
		Height:        bit.height,
		Bytewidth:     bit.bytewidth,
		BitsPerPixel:  bit.bitsPerPixel,
		BytesPerPixel: bit.bytesPerPixel,
	}

	return bit_map
}

/*
.___  ___.   ______    __    __       _______. _______
|   \/   |  /  __  \  |  |  |  |     /       ||   ____|
|  \  /  | |  |  |  | |  |  |  |    |   (----`|  |__
|  |\/|  | |  |  |  | |  |  |  |     \   \    |   __|
|  |  |  | |  `--'  | |  `--'  | .----)   |   |  |____
|__|  |__|  \______/   \______/  |_______/    |_______|

*/

type MPoint struct {
	x int
	y int
}

//C.size_t  int
func MoveMouse(x, y int) {
	cx := C.size_t(x)
	cy := C.size_t(y)
	C.aMoveMouse(cx, cy)
}

func DragMouse(x, y int) {
	cx := C.size_t(x)
	cy := C.size_t(y)
	C.aDragMouse(cx, cy)
}

func MoveMouseSmooth(x, y int) {
	cx := C.size_t(x)
	cy := C.size_t(y)
	C.aMoveMouseSmooth(cx, cy)
}

func GetMousePos() (int, int) {
	pos := C.aGetMousePos()
	// Println("pos:###", pos, pos.x, pos.y)
	x := int(pos.x)
	y := int(pos.y)
	// return pos.x, pos.y
	return x, y
}

func MouseClick() {
	C.aMouseClick()
}

func MouseToggle(args ...interface{}) {
	var button C.MMMouseButton
	Try(func() {
		button = args[1].(C.MMMouseButton)
	}, func(e interface{}) {
		// Println("err:::", e)
		button = C.LEFT_BUTTON
	})
	down := C.CString(args[0].(string))
	C.aMouseToggle(down, button)
}

func SetMouseDelay(x int) {
	cx := C.size_t(x)
	C.aSetMouseDelay(cx)
}

func ScrollMouse(x int, y string) {
	cx := C.size_t(x)
	z := C.CString(y)
	C.aScrollMouse(cx, z)
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
func Try(fun func(), handler func(interface{})) {
	defer func() {
		if err := recover(); err != nil {
			handler(err)
		}
	}()
	fun()
}

func KeyTap(args ...string) {
	var apara string
	Try(func() {
		apara = args[1]
	}, func(e interface{}) {
		// Println("err:::", e)
		apara = "null"
	})

	zkey := C.CString(args[0])
	amod := C.CString(apara)
	// defer func() {
	C.aKeyTap(zkey, amod)
	// }()

	defer C.free(unsafe.Pointer(zkey))
	defer C.free(unsafe.Pointer(amod))
}

func KeyToggle(args ...string) {
	var apara string
	Try(func() {
		apara = args[1]
	}, func(e interface{}) {
		// Println("err:::", e)
		apara = "null"
	})

	zkey := C.CString(args[0])
	amod := C.CString(apara)
	// defer func() {
	str := C.aKeyToggle(zkey, amod)
	Println(str)
	// }()
	defer C.free(unsafe.Pointer(zkey))
	defer C.free(unsafe.Pointer(amod))
}

func TypeString(x string) {
	cx := C.CString(x)
	C.aTypeString(cx)
	defer C.free(unsafe.Pointer(cx))
}

func TypeStringDelayed(x string, y C.size_t) {
	cx := C.CString(x)
	C.aTypeStringDelayed(cx, y)
	defer C.free(unsafe.Pointer(cx))
}

func SetKeyboardDelay(x C.size_t) {
	C.aSetKeyboardDelay(x)
}

/*
.______    __  .___________..___  ___.      ___      .______
|   _  \  |  | |           ||   \/   |     /   \     |   _  \
|  |_)  | |  | `---|  |----`|  \  /  |    /  ^  \    |  |_)  |
|   _  <  |  |     |  |     |  |\/|  |   /  /_\  \   |   ___/
|  |_)  | |  |     |  |     |  |  |  |  /  _____  \  |  |
|______/  |__|     |__|     |__|  |__| /__/     \__\ | _|
*/

/*
____    __    ____  __  .__   __.  _______   ______   ____    __    ____
\   \  /  \  /   / |  | |  \ |  | |       \ /  __  \  \   \  /  \  /   /
 \   \/    \/   /  |  | |   \|  | |  .--.  |  |  |  |  \   \/    \/   /
  \            /   |  | |  . `  | |  |  |  |  |  |  |   \            /
   \    /\    /    |  | |  |\   | |  '--'  |  `--'  |    \    /\    /
    \__/  \__/     |__| |__| \__| |_______/ \______/      \__/  \__/

*/

/*
------------ ---    ---  ------------ ----    ---- ------------
************ ***    ***  ************ *****   **** ************
----         ---    ---  ----         ------  ---- ------------
************ ***    ***  ************ ************     ****
------------ ---    ---  ------------ ------------     ----
****          ********   ****         ****  ******     ****
------------   ------    ------------ ----   -----     ----
************    ****     ************ ****    ****     ****

*/
