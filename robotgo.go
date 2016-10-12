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
*/
import "C"

import (
	. "fmt"
	"unsafe"
	// "runtime"
	// "syscall"
)

/*
 __  __
|  \/  | ___  _   _ ___  ___
| |\/| |/ _ \| | | / __|/ _ \
| |  | | (_) | |_| \__ \  __/
|_|  |_|\___/ \__,_|___/\___|

*/

type MPoint struct {
	x int
	y int
}

//C.size_t  int
func MoveMouse(x, y C.int) {
	C.amoveMouse(x, y)
}

func DragMouse(x, y C.int) {
	C.adragMouse(x, y)
}

func MoveMouseSmooth(x, y C.int) {
	C.amoveMouseSmooth(x, y)
}

func GetMousePos() (C.size_t, C.size_t) {
	pos := C.agetMousePos()
	// Println("pos:###", pos, pos.x, pos.y)
	return pos.x, pos.y
}

func MouseClick() {
	C.amouseClick()
}

func MouseToggle() {
	C.amouseToggle()
}

func SetMouseDelay(x C.int) {
	C.asetMouseDelay(x)
}

func ScrollMouse(x C.int, y string) {
	z := C.CString(y)
	C.ascrollMouse(x, z)
	defer C.free(unsafe.Pointer(z))
}

/*
 _  __          _                         _
| |/ /___ _   _| |__   ___   __ _ _ __ __| |
| ' // _ \ | | | '_ \ / _ \ / _` | '__/ _` |
| . \  __/ |_| | |_) | (_) | (_| | | | (_| |
|_|\_\___|\__, |_.__/ \___/ \__,_|_|  \__,_|
		  |___/
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
	C.akeyTap(zkey, amod)
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
	str := C.akeyToggle(zkey, amod)
	Println(str)
	// }()
	defer C.free(unsafe.Pointer(zkey))
	defer C.free(unsafe.Pointer(amod))
}

func TypeString(x string) {
	cx := C.CString(x)
	C.atypeString(cx)
	defer C.free(unsafe.Pointer(cx))
}

func TypeStringDelayed(x string, y C.size_t) {
	cx := C.CString(x)
	C.atypeStringDelayed(cx, y)
	defer C.free(unsafe.Pointer(cx))
}

func SetKeyboardDelay(x C.size_t) {
	C.asetKeyboardDelay(x)
}

/*
  ____
 / ___|  ___ _ __ ___  ___ _ __
 \___ \ / __| '__/ _ \/ _ \ '_ \
  ___) | (__| | |  __/  __/ | | |
 |____/ \___|_|  \___|\___|_| |_|

*/

func GetPixelColor(x, y C.size_t) string {
	color := C.agetPixelColor(x, y)
	gcolor := C.GoString(color)
	defer C.free(unsafe.Pointer(color))
	return gcolor
}

func GetScreenSize() (C.size_t, C.size_t) {
	size := C.agetScreenSize()
	// Println("...", size, size.width)
	return size.width, size.height
}

func GetXDisplayName() string {
	name := C.agetXDisplayName()
	gname := C.GoString(name)
	defer C.free(unsafe.Pointer(name))
	return gname
}

func SetXDisplayName(name string) string {
	cname := C.CString(name)
	str := C.asetXDisplayName(cname)
	gstr := C.GoString(str)
	return gstr
}

func CaptureScreen(x, y, w, h C.int) {
	bit := C.acaptureScreen(x, y, w, h)
	Println("...", bit)
}
