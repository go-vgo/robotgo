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

Please make sure Golang, GCC, zlib and libpng is installed correctly before installing RobotGo;

See Requirements:
	https://github.com/go-vgo/robotgo#requirements

Installation:
	go get -u github.com/go-vgo/robotgo
*/
package robotgo

/*
//#if defined(IS_MACOSX)
	#cgo darwin CFLAGS: -x objective-c  -Wno-deprecated-declarations -I/usr/local/opt/libpng/include -I/usr/local/opt/zlib/include
	#cgo darwin LDFLAGS: -framework Cocoa -framework OpenGL -framework IOKit -framework Carbon -framework CoreFoundation 
//#elif defined(USE_X11)
	// drop -std=c11
	#cgo linux CFLAGS:-I/usr/src
	#cgo linux LDFLAGS:-L/usr/src  -lX11 -lXtst -lX11-xcb -lxcb -lxcb-xkb -lxkbcommon -lxkbcommon-x11 -lm
//#endif
	#cgo windows LDFLAGS: -lgdi32 -luser32 
//#include <AppKit/NSEvent.h>
#include "screen/goScreen.h"
#include "mouse/goMouse.h"
#include "key/goKey.h"
// // // // // #include "bitmap/goBitmap.h"
#include "event/goEvent.h"
#include "window/goWindow.h"
*/
import "C"

import (
	// "fmt"
	"reflect"
	"runtime"
	"strings"
	"unsafe"
	// "syscall"

	"github.com/go-vgo/robotgo/clipboard"
	"github.com/shirou/gopsutil/process"
)

const (
	version string = "v0.46.0.400, Pyrenees Mountains!"
)

// GetVersion get version
func GetVersion() string {
	return version
}

/*
      _______.  ______ .______       _______  _______ .__   __.
    /       | /      ||   _  \     |   ____||   ____||  \ |  |
   |   (----`|  ,----'|  |_)  |    |  |__   |  |__   |   \|  |
    \   \    |  |     |      /     |   __|  |   __|  |  . `  |
.----)   |   |  `----.|  |\  \----.|  |____ |  |____ |  |\   |
|_______/     \______|| _| `._____||_______||_______||__| \__|
*/

// Bitmap is Bitmap struct
type Bitmap struct {
	ImageBuffer   *uint8
	Width         int
	Height        int
	Bytewidth     int
	BitsPerPixel  uint8
	BytesPerPixel uint8
}

// GetPixelColor get pixel color
func GetPixelColor(x, y int) string {
	cx := C.size_t(x)
	cy := C.size_t(y)
	color := C.aGetPixelColor(cx, cy)
	// color := C.aGetPixelColor(x, y)
	gcolor := C.GoString(color)
	defer C.free(unsafe.Pointer(color))
	return gcolor
}

// GetScreenSize get screen size
func GetScreenSize() (int, int) {
	size := C.aGetScreenSize()
	// fmt.Println("...", size, size.width)
	return int(size.width), int(size.height)
}

// SetXDisplayName set XDisplay name
func SetXDisplayName(name string) string {
	cname := C.CString(name)
	str := C.aSetXDisplayName(cname)
	gstr := C.GoString(str)
	defer C.free(unsafe.Pointer(cname))
	return gstr
}

// GetXDisplayName get XDisplay name
func GetXDisplayName() string {
	name := C.aGetXDisplayName()
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
		//Get screen size.
		var displaySize C.MMSize
		displaySize = C.getMainDisplaySize()
		w = displaySize.width
		h = displaySize.height
	})

	bit := C.aCaptureScreen(x, y, w, h)
	// fmt.Println("...", bit.width)
	return bit
}

// BCaptureScreen capture the screen and return bitmap(go struct)
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

	bit := C.aCaptureScreen(x, y, w, h)
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

/*
.___  ___.   ______    __    __       _______. _______
|   \/   |  /  __  \  |  |  |  |     /       ||   ____|
|  \  /  | |  |  |  | |  |  |  |    |   (----`|  |__
|  |\/|  | |  |  |  | |  |  |  |     \   \    |   __|
|  |  |  | |  `--'  | |  `--'  | .----)   |   |  |____
|__|  |__|  \______/   \______/  |_______/    |_______|

*/

// MPoint is MPoint struct
type MPoint struct {
	x int
	y int
}

// MoveMouse move the mouse
func MoveMouse(x, y int) {
	//C.size_t  int
	cx := C.size_t(x)
	cy := C.size_t(y)
	C.aMoveMouse(cx, cy)
}

// Move move the mouse
func Move(x, y int) {
	cx := C.size_t(x)
	cy := C.size_t(y)
	C.aMoveMouse(cx, cy)
}

// DragMouse drag the mouse
func DragMouse(x, y int) {
	cx := C.size_t(x)
	cy := C.size_t(y)
	C.aDragMouse(cx, cy)
}

// Drag drag the mouse
func Drag(x, y int) {
	cx := C.size_t(x)
	cy := C.size_t(y)
	C.aDragMouse(cx, cy)
}

// MoveMouseSmooth move the mouse smooth
func MoveMouseSmooth(x, y int, args ...float64) {
	cx := C.size_t(x)
	cy := C.size_t(y)

	var (
		low  C.double
		high C.double
	)

	if len(args) > 1 {
		low = C.double(args[0])
		high = C.double(args[1])
	} else {
		low = 5.0
		high = 500.0
	}

	C.aMoveMouseSmooth(cx, cy, low, high)
}

// MoveSmooth move the mouse smooth
func MoveSmooth(x, y int, args ...float64) {
	cx := C.size_t(x)
	cy := C.size_t(y)

	var (
		low  C.double
		high C.double
	)

	if len(args) > 1 {
		low = C.double(args[0])
		high = C.double(args[1])
	} else {
		low = 5.0
		high = 500.0
	}

	C.aMoveMouseSmooth(cx, cy, low, high)
}

// GetMousePos get mouse portion
func GetMousePos() (int, int) {
	pos := C.aGetMousePos()
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
	C.aMouseClick(button, double)
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
	C.aMouseClick(button, double)
}

// MoveClick move and click the mouse
func MoveClick(x, y int, args ...interface{}) {
	MoveMouse(x, y)
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
	C.aMouseToggle(down, button)
	defer C.free(unsafe.Pointer(down))
}

// SetMouseDelay set mouse delay
func SetMouseDelay(x int) {
	cx := C.size_t(x)
	C.aSetMouseDelay(cx)
}

// ScrollMouse scroll the mouse
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

// Try handler(err)
func Try(fun func(), handler func(interface{})) {
	defer func() {
		if err := recover(); err != nil {
			handler(err)
		}
	}()
	fun()
}

// KeyTap tap the keyboard;
//
// See keys:
//	https://github.com/go-vgo/robotgo/blob/master/docs/keys.md
func KeyTap(args ...interface{}) {
	var (
		akey   string
		akeyt  string
		keyarr []string
		num    int
	)
	// var ckeyarr []*C.char
	ckeyarr := make([](*_Ctype_char), 0)

	Try(func() {
		if reflect.TypeOf(args[1]) == reflect.TypeOf(keyarr) {

			keyarr = args[1].([]string)

			num = len(keyarr)

			for i := 0; i < num; i++ {
				ckeyarr = append(ckeyarr, (*C.char)(unsafe.Pointer(C.CString(keyarr[i]))))
			}

		} else {
			akey = args[1].(string)

			Try(func() {
				akeyt = args[2].(string)

			}, func(e interface{}) {
				// fmt.Println("err:::", e)
				akeyt = "null"
			})
		}

	}, func(e interface{}) {
		// fmt.Println("err:::", e)
		akey = "null"
		keyarr = []string{"null"}
	})
	// }()

	zkey := C.CString(args[0].(string))

	if akey == "" && len(keyarr) != 0 {
		C.aKey_Tap(zkey, (**_Ctype_char)(unsafe.Pointer(&ckeyarr[0])), C.int(num))
	} else {
		// zkey := C.CString(args[0])
		amod := C.CString(akey)
		amodt := C.CString(akeyt)

		C.aKeyTap(zkey, amod, amodt)

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
		adown  string
		amkey  string
		amkeyt string
	)

	Try(func() {
		adown = args[1]
		Try(func() {
			amkey = args[2]
			Try(func() {
				amkeyt = args[3]
			}, func(e interface{}) {
				// fmt.Println("err:::", e)
				amkeyt = "null"
			})
		}, func(e interface{}) {
			// fmt.Println("err:::", e)
			amkey = "null"
		})
	}, func(e interface{}) {
		// fmt.Println("err:::", e)
		adown = "null"
	})

	ckey := C.CString(args[0])
	cadown := C.CString(adown)
	camkey := C.CString(amkey)
	camkeyt := C.CString(amkeyt)
	// defer func() {
	str := C.aKeyToggle(ckey, cadown, camkey, camkeyt)
	// fmt.Println(str)
	// C.aKeyToggle(ckey, cadown, camkey, camkeyt)
	// }()
	defer C.free(unsafe.Pointer(ckey))
	defer C.free(unsafe.Pointer(cadown))
	defer C.free(unsafe.Pointer(camkey))
	defer C.free(unsafe.Pointer(camkeyt))

	return C.GoString(str)
}

// TypeString type string
func TypeString(x string) {
	cx := C.CString(x)
	C.aTypeString(cx)
	defer C.free(unsafe.Pointer(cx))
}

// TypeStr type string, support UTF-8
func TypeStr(str string) {
	clipboard.WriteAll(str)
	if runtime.GOOS == "darwin" {
		KeyTap("v", "command")
	} else {
		KeyTap("v", "control")
	}
}

// ReadAll read string from clipboard
func ReadAll() (string, error) {
	return clipboard.ReadAll()
}

// WriteAll write string to clipboard
func WriteAll(text string) {
	clipboard.WriteAll(text)
}

// TypeStrDelay type string delayed
func TypeStrDelay(x string, y int) {
	cx := C.CString(x)
	cy := C.size_t(y)
	C.aTypeStringDelayed(cx, cy)
	defer C.free(unsafe.Pointer(cx))
}

// TypeStringDelayed type string delayed, Wno-deprecated
func TypeStringDelayed(x string, y int) {
	cx := C.CString(x)
	cy := C.size_t(y)
	C.aTypeStringDelayed(cx, cy)
	defer C.free(unsafe.Pointer(cx))
}

// SetKeyDelay set keyboard delay
func SetKeyDelay(x int) {
	C.aSetKeyboardDelay(C.size_t(x))
}

// SetKeyboardDelay set keyboard delay, Wno-deprecated
func SetKeyboardDelay(x int) {
	C.aSetKeyboardDelay(C.size_t(x))
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
 ___________    ____  _______ .__   __. .___________.
|   ____\   \  /   / |   ____||  \ |  | |           |
|  |__   \   \/   /  |  |__   |   \|  | `---|  |----`
|   __|   \      /   |   __|  |  . `  |     |  |
|  |____   \    /    |  |____ |  |\   |     |  |
|_______|   \__/     |_______||__| \__|     |__|
*/

// Map a map
type Map map[string]interface{}

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
	}

	var cs *C.char
	var keve string

	if len(aeve) > 1 && len(aeve) < 4 {
		keve = keycode[aeve].(string)
		cs = C.CString(keve)
	} else {
		cs = C.CString(aeve)
	}

	// cs := C.CString(aeve)
	eve := C.aEvent(cs)
	// fmt.Println("event@@", eve)
	geve := int(eve)
	defer C.free(unsafe.Pointer(cs))
	return geve
}

// StopEvent stop event listener
func StopEvent() {
	C.aStop()
}

// LEvent add event listener, Wno-deprecated
func LEvent(aeve string) int {
	cs := C.CString(aeve)
	eve := C.aEvent(cs)
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

	cbool := C.aShowAlert(atitle, amsg, adefaultButton, acancelButton)
	ibool := int(cbool)
	defer C.free(unsafe.Pointer(atitle))
	defer C.free(unsafe.Pointer(amsg))
	defer C.free(unsafe.Pointer(adefaultButton))
	defer C.free(unsafe.Pointer(acancelButton))
	return ibool
}

// IsValid valid the window
func IsValid() bool {
	abool := C.aIsValid()
	gbool := bool(abool)
	// fmt.Println("bool---------", gbool)
	return gbool
}

// SetActive set the window active
func SetActive(win C.MData) {
	C.aSetActive(win)
}

// GetActive get the active window
func GetActive() C.MData {
	mdata := C.aGetActive()
	// fmt.Println("active----", mdata)
	return mdata
}

// CloseWindow close the window
func CloseWindow() {
	C.aCloseWindow()
}

// SetHandle set the window handle
func SetHandle(hwnd int) {
	chwnd := C.uintptr(hwnd)
	C.aSetHandle(chwnd)
}

// GetHandle get the window handle
func GetHandle() int {
	hwnd := C.aGetHandle()
	ghwnd := int(hwnd)
	// fmt.Println("gethwnd---", ghwnd)
	return ghwnd
}

// GetBHandle get the window handle
func GetBHandle() int {
	hwnd := C.bGetHandle()
	ghwnd := int(hwnd)
	//fmt.Println("gethwnd---", ghwnd)
	return ghwnd
}

// GetTitle get the window title
func GetTitle() string {
	title := C.aGetTitle()
	gtittle := C.GoString(title)
	// fmt.Println("title...", gtittle)
	return gtittle
}

// GetPID get the process id
func GetPID() int {
	pid := C.aGetPID()
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

// ActivePID window active by PID
func ActivePID(pid int32) {
	C.active_PID(C.uintptr(pid))
}
