package robotgo

/*
//#if defined(IS_MACOSX)
	#cgo darwin CFLAGS: -x objective-c  -Wno-deprecated-declarations -I/usr/local/opt/libpng/include -I/usr/local/opt/zlib/include
	#cgo darwin LDFLAGS: -framework Cocoa -framework OpenGL -framework IOKit -framework Carbon -framework CoreFoundation -L/usr/local/opt/libpng/lib -lpng -L/usr/local/opt/zlib/lib -lz
//#elif defined(USE_X11)
	#cgo linux CFLAGS:-I/usr/src
	#cgo linux LDFLAGS:-L/usr/src -lpng -lz -lX11 -lXtst -lX11-xcb -lxcb -lxcb-xkb -lxkbcommon -lxkbcommon-x11 -lm
//#endif
	#cgo windows LDFLAGS: -lgdi32 -luser32 -lpng -lz
//#include <AppKit/NSEvent.h>
#include "screen/goScreen.h"
#include "mouse/goMouse.h"
#include "key/goKey.h"
#include "bitmap/goBitmap.h"
#include "event/goEvent.h"
#include "window/goWindow.h"
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

func GetPixelColor(x, y int) string {
	cx := C.size_t(x)
	cy := C.size_t(y)
	color := C.aGetPixelColor(cx, cy)
	// color := C.aGetPixelColor(x, y)
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

func CaptureScreen(args ...int) C.MMBitmapRef {
	var x C.size_t
	var y C.size_t
	var w C.size_t
	var h C.size_t
	Try(func() {
		x = C.size_t(args[0])
		y = C.size_t(args[1])
		w = C.size_t(args[2])
		h = C.size_t(args[3])
	}, func(e interface{}) {
		// Println("err:::", e)
		x = 0
		y = 0
		//Get screen size.
		var displaySize C.MMSize
		displaySize = C.getMainDisplaySize()
		w = displaySize.width
		h = displaySize.height
	})

	bit := C.aCaptureScreen(x, y, w, h)
	// Println("...", bit.width)
	return bit
}

func Capture_Screen(x, y, w, h C.size_t) Bit_map {
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
func FindBitmap(args ...interface{}) (C.size_t, C.size_t) {
	var bit C.MMBitmapRef
	bit = args[0].(C.MMBitmapRef)

	var rect C.MMRect
	Try(func() {
		rect.origin.x = C.size_t(args[1].(int))
		rect.origin.y = C.size_t(args[2].(int))
		rect.size.width = C.size_t(args[3].(int))
		rect.size.height = C.size_t(args[4].(int))
	}, func(e interface{}) {
		Println("err:::", e)
		// rect.origin.x = x
		// rect.origin.y = y
		// rect.size.width = w
		// rect.size.height = h
	})

	pos := C.aFindBitmap(bit, rect)
	// Println("pos----", pos)
	return pos.x, pos.y
}

func OpenBitmap(args ...interface{}) C.MMBitmapRef {
	path := C.CString(args[0].(string))
	var mtype C.uint16_t
	Try(func() {
		mtype = C.uint16_t(args[1].(int))
	}, func(e interface{}) {
		// Println("err:::", e)
		mtype = 1
	})
	bit := C.aOpenBitmap(path, mtype)
	// Println("opening...", bit)
	return bit
	// defer C.free(unsafe.Pointer(path))
}

func SaveBitmap(args ...interface{}) {
	var mtype C.uint16_t
	Try(func() {
		mtype = C.uint16_t(args[2].(int))
	}, func(e interface{}) {
		// Println("err:::", e)
		mtype = 1
	})

	path := C.CString(args[1].(string))
	savebit := C.aSaveBitmap(args[0].(C.MMBitmapRef), path, mtype)
	Println("saved...", savebit)
	// return bit
	// defer C.free(unsafe.Pointer(path))
}

// func SaveBitmap(bit C.MMBitmapRef, gpath string, mtype C.MMImageType) {
// 	path := C.CString(gpath)
// 	savebit := C.aSaveBitmap(bit, path, mtype)
// 	Println("opening...", savebit)
// 	// return bit
// 	// defer C.free(unsafe.Pointer(path))
// }

func TostringBitmap(bit C.MMBitmapRef) *C.char {
	str_bit := C.aTostringBitmap(bit)
	// Println("...", str_bit)
	return str_bit
}

func GetPortion(bit C.MMBitmapRef, x, y, w, h C.size_t) C.MMBitmapRef {
	var rect C.MMRect
	rect.origin.x = x
	rect.origin.y = y
	rect.size.width = w
	rect.size.height = h

	pos := C.aGetPortion(bit, rect)
	return pos
}

func Convert(args ...interface{}) {
	var mtype int
	Try(func() {
		mtype = args[2].(int)
	}, func(e interface{}) {
		Println("err:::", e)
		mtype = 1
	})
	//C.CString()
	opath := args[0].(string)
	spath := args[1].(string)
	bit_map := OpenBitmap(opath)
	// Println("a----", bit_map)
	SaveBitmap(bit_map, spath, mtype)
}

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

func LEvent(aeve string) int {
	cs := C.CString(aeve)
	eve := C.aEvent(cs)
	// Println("event@@", eve)
	geve := int(eve)
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

	cbool := C.aShowAlert(atitle, amsg, adefaultButton, acancelButton)
	ibool := int(cbool)
	return ibool
}
