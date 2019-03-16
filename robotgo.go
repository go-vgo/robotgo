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
	#cgo linux LDFLAGS: -L/usr/src -lpng -lz -lX11 -lXtst -lm
	// #cgo linux LDFLAGS: -lX11-xcb -lxcb -lxcb-xkb -lxkbcommon -lxkbcommon-x11
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
//#include "event/goEvent.h"
#include "window/goWindow.h"
*/
import "C"

import (
	// "fmt"
	"image"
	"os"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"time"
	"unsafe"
	// "syscall"
	"os/exec"

	"github.com/go-vgo/robotgo/clipboard"
	"github.com/robotn/gohook"
	"github.com/shirou/gopsutil/process"
	"github.com/vcaesar/imgo"
)

const (
	// Version get the robotgo version
	Version = "v0.80.0.845, Sierra Nevada!"
)

// GetVersion get the robotgo version
func GetVersion() string {
	return Version
}

type (
	// Map a map[string]interface{}
	Map map[string]interface{}
	// CHex define CHex as c rgb Hex type (C.MMRGBHex)
	CHex C.MMRGBHex
	// CBitmap define CBitmap as C.MMBitmapRef type
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

// MicroSleep time C.microsleep(tm)
func MicroSleep(tm float64) {
	C.microsleep(C.double(tm))
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
	C.free(unsafe.Pointer(color))

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
	C.free(unsafe.Pointer(color))

	return gcolor
}

// ScaleX get primary display horizontal DPI scale factor
func ScaleX() int {
	return int(C.scale_x())
}

// ScaleY get primary display vertical DPI scale factor
func ScaleY() int {
	return int(C.scale_y())
}

// GetScreenSize get the screen size
func GetScreenSize() (int, int) {
	size := C.get_screen_size()
	// fmt.Println("...", size, size.width)
	return int(size.width), int(size.height)
}

// Scale get the screen scale
func Scale() int {
	dpi := map[int]int{
		0: 100,
		// DPI Scaling Level
		96:  100,
		120: 125,
		144: 150,
		168: 175,
		192: 200,
		216: 225,
		// Custom DPI
		240: 250,
		288: 300,
		384: 400,
		480: 500,
	}

	x := ScaleX()
	return dpi[x]
}

// Mul mul the scale
func Mul(x int) int {
	s := Scale()
	return x * s / 100
}

// GetScaleSize get the screen scale size
func GetScaleSize() (int, int) {
	x, y := GetScreenSize()
	s := Scale()
	return x * s / 100, y * s / 100
}

// SetXDisplayName set XDisplay name
func SetXDisplayName(name string) string {
	cname := C.CString(name)
	str := C.set_XDisplay_name(cname)
	gstr := C.GoString(str)
	C.free(unsafe.Pointer(cname))

	return gstr
}

// GetXDisplayName get XDisplay name
func GetXDisplayName() string {
	name := C.get_XDisplay_name()
	gname := C.GoString(name)
	C.free(unsafe.Pointer(name))

	return gname
}

// CaptureScreen capture the screen return bitmap(c struct),
// use `defer robotgo.FreeBitmap(bitmap)` to free the bitmap
//
// robotgo.CaptureScreen(x, y, w, h int)
func CaptureScreen(args ...int) C.MMBitmapRef {
	var x, y, w, h C.size_t

	if len(args) > 3 {
		x = C.size_t(args[0])
		y = C.size_t(args[1])
		w = C.size_t(args[2])
		h = C.size_t(args[3])
	} else {
		x = 0
		y = 0
		// Get screen size.
		var displaySize C.MMSize
		displaySize = C.getMainDisplaySize()
		w = displaySize.width
		h = displaySize.height
	}

	bit := C.capture_screen(x, y, w, h)
	return bit
}

// GoCaptureScreen capture the screen and return bitmap(go struct)
func GoCaptureScreen(args ...int) Bitmap {
	bit := CaptureScreen(args...)
	defer FreeBitmap(bit)

	return ToBitmap(bit)
}

// SaveCapture capture screen and save
func SaveCapture(spath string, args ...int) {
	bit := CaptureScreen(args...)

	SaveBitmap(bit, spath)
	FreeBitmap(bit)
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
	if btn == "left" {
		return C.LEFT_BUTTON
	}

	if btn == "center" {
		return C.CENTER_BUTTON
	}

	if btn == "right" {
		return C.RIGHT_BUTTON
	}

	return C.LEFT_BUTTON
}

// MoveMouse move the mouse
func MoveMouse(x, y int) {
	// C.size_t  int
	Move(x, y)
}

// Move move the mouse
func Move(x, y int) {
	cx := C.size_t(x)
	cy := C.size_t(y)
	C.move_mouse(cx, cy)
}

// DragMouse drag the mouse
func DragMouse(x, y int, args ...string) {
	Drag(x, y, args...)
}

// Drag drag the mouse
func Drag(x, y int, args ...string) {
	var button C.MMMouseButton = C.LEFT_BUTTON

	cx := C.size_t(x)
	cy := C.size_t(y)

	if len(args) > 0 {
		button = CheckMouse(args[0])
	}

	C.drag_mouse(cx, cy, button)
}

// MoveMouseSmooth move the mouse smooth,
// moves mouse to x, y human like, with the mouse button up.
func MoveMouseSmooth(x, y int, args ...interface{}) bool {
	return MoveSmooth(x, y, args...)
}

// MoveSmooth move the mouse smooth,
// moves mouse to x, y human like, with the mouse button up.
//
// robotgo.MoveSmooth(x, y int, low, high float64, mouseDelay int)
func MoveSmooth(x, y int, args ...interface{}) bool {
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
		low = 1.0
		high = 3.0
	}

	cbool := C.move_mouse_smooth(cx, cy, low, high, C.int(mouseDelay))

	return bool(cbool)
}

// GetMousePos get mouse's portion
func GetMousePos() (int, int) {
	pos := C.get_mouse_pos()

	x := int(pos.x)
	y := int(pos.y)

	return x, y
}

// MouseClick click the mouse
//
// robotgo.MouseClick(button string, double bool)
func MouseClick(args ...interface{}) {
	Click(args...)
}

// Click click the mouse
//
// robotgo.Click(button string, double bool)
func Click(args ...interface{}) {
	var (
		button C.MMMouseButton = C.LEFT_BUTTON
		double C.bool
	)

	if len(args) > 0 {
		button = CheckMouse(args[0].(string))
	}

	if len(args) > 1 {
		double = C.bool(args[1].(bool))
	}

	C.mouse_click(button, double)
}

// MoveClick move and click the mouse
//
// robotgo.MoveClick(x, y int, button string, double bool)
func MoveClick(x, y int, args ...interface{}) {
	MoveMouse(x, y)
	MouseClick(args...)
}

// MovesClick move smooth and click the mouse
func MovesClick(x, y int, args ...interface{}) {
	MoveSmooth(x, y)
	MouseClick(args...)
}

// MouseToggle toggle the mouse
func MouseToggle(togKey string, args ...interface{}) {
	var button C.MMMouseButton = C.LEFT_BUTTON

	if len(args) > 0 {
		button = CheckMouse(args[0].(string))
	}

	down := C.CString(togKey)
	C.mouse_toggle(down, button)

	C.free(unsafe.Pointer(down))
}

// SetMouseDelay set mouse delay
func SetMouseDelay(delay int) {
	cdelay := C.size_t(delay)
	C.set_mouse_delay(cdelay)
}

// ScrollMouse scroll the mouse
func ScrollMouse(x int, direction string) {
	cx := C.size_t(x)
	cy := C.CString(direction)
	C.scroll_mouse(cx, cy)

	C.free(unsafe.Pointer(cy))
}

// Scroll scroll the mouse with x, y
//
// robotgo.Scroll(x, y, msDelay int)
func Scroll(x, y int, args ...int) {
	var msDelay = 10
	if len(args) > 0 {
		msDelay = args[0]
	}

	cx := C.int(x)
	cy := C.int(y)
	cz := C.int(msDelay)

	C.scroll(cx, cy, cz)
}

/*
 __  ___  ___________    ____ .______     ______        ___      .______       _______
|  |/  / |   ____\   \  /   / |   _  \   /  __  \      /   \     |   _  \     |       \
|  '  /  |  |__   \   \/   /  |  |_)  | |  |  |  |    /  ^  \    |  |_)  |    |  .--.  |
|    <   |   __|   \_    _/   |   _  <  |  |  |  |   /  /_\  \   |      /     |  |  |  |
|  .  \  |  |____    |  |     |  |_)  | |  `--'  |  /  _____  \  |  |\  \----.|  '--'  |
|__|\__\ |_______|   |__|     |______/   \______/  /__/     \__\ | _| `._____||_______/

*/

// KeyTap tap the keyboard code;
//
// See keys:
//	https://github.com/go-vgo/robotgo/blob/master/docs/keys.md
func KeyTap(tapKey string, args ...interface{}) string {
	var (
		akey     string
		keyT     = "null"
		keyArr   []string
		num      int
		keyDelay = 10
	)
	// var ckeyArr []*C.char
	ckeyArr := make([](*C.char), 0)

	if len(args) > 0 {
		if reflect.TypeOf(args[0]) == reflect.TypeOf(keyArr) {

			keyArr = args[0].([]string)
			num = len(keyArr)

			for i := 0; i < num; i++ {
				ckeyArr = append(ckeyArr, (*C.char)(unsafe.Pointer(C.CString(keyArr[i]))))
			}

			if len(args) > 1 {
				keyDelay = args[1].(int)
			}
		} else {
			akey = args[0].(string)

			if len(args) > 1 {
				if reflect.TypeOf(args[1]) == reflect.TypeOf(akey) {
					keyT = args[1].(string)
				} else {
					keyDelay = args[1].(int)
				}
			}
		}

	} else {
		akey = "null"
		keyArr = []string{"null"}
	}

	// zkey := C.CString(args[0])
	zkey := C.CString(tapKey)
	defer C.free(unsafe.Pointer(zkey))

	if akey == "" && len(keyArr) != 0 {
		str := C.key_Taps(zkey, (**C.char)(unsafe.Pointer(&ckeyArr[0])),
			C.int(num), C.int(keyDelay))

		return C.GoString(str)
	}

	amod := C.CString(akey)
	amodt := C.CString(keyT)

	str := C.key_tap(zkey, amod, amodt, C.int(keyDelay))

	C.free(unsafe.Pointer(amod))
	C.free(unsafe.Pointer(amodt))

	return C.GoString(str)
}

// KeyToggle toggle the keyboard
//
// See keys:
//	https://github.com/go-vgo/robotgo/blob/master/docs/keys.md
func KeyToggle(args ...string) string {
	var (
		down, mKey, mKeyT = "null", "null", "null"
		// keyDelay = 10
	)

	if len(args) > 1 {
		down = args[1]

		if len(args) > 2 {
			mKey = args[2]

			if len(args) > 3 {
				mKeyT = args[3]
			}
		}
	}

	ckey := C.CString(args[0])
	cdown := C.CString(down)
	cmKey := C.CString(mKey)
	cmKeyT := C.CString(mKeyT)

	str := C.key_toggle(ckey, cdown, cmKey, cmKeyT)
	// str := C.key_Toggle(ckey, cdown, cmKey, cmKeyT, C.int(keyDelay))

	C.free(unsafe.Pointer(ckey))
	C.free(unsafe.Pointer(cdown))
	C.free(unsafe.Pointer(cmKey))
	C.free(unsafe.Pointer(cmKeyT))

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

func toUC(text string) []string {
	var uc []string

	textQuoted := strconv.QuoteToASCII(text)
	textUnquoted := textQuoted[1 : len(textQuoted)-1]

	strUnicodev := strings.Split(textUnquoted, "\\u")
	for i := 1; i < len(strUnicodev); i++ {
		uc = append(uc, "U"+strUnicodev[i])
	}

	return uc
}

func inputUTF(str string) {
	cstr := C.CString(str)
	C.input_utf(cstr)

	C.free(unsafe.Pointer(cstr))
}

// TypeStr send a string, support UTF-8
//
// robotgo.TypeStr(string: The string to send, float64: microsleep time)
func TypeStr(str string, args ...float64) {
	var tm = 7.0
	if len(args) > 0 {
		tm = args[0]
	}

	if runtime.GOOS == "linux" {
		strUc := toUC(str)
		for i := 0; i < len(strUc); i++ {
			inputUTF(strUc[i])
			MicroSleep(tm)
		}

		return
	}

	for i := 0; i < len([]rune(str)); i++ {
		ustr := uint32(CharCodeAt(str, i))
		UnicodeType(ustr)

		// if len(args) > 0 {
		// 	MicroSleep(tm)
		// }
	}
}

// UnicodeType tap uint32 unicode
func UnicodeType(str uint32) {
	cstr := C.uint(str)
	C.unicodeType(cstr)
}

// TypeString send a string, support unicode
// TypeStr(string: The string to send)
func TypeString(str string) {
	cstr := C.CString(str)
	C.type_string(cstr)

	C.free(unsafe.Pointer(cstr))
}

// PasteStr paste a string, support UTF-8
func PasteStr(str string) {
	clipboard.WriteAll(str)
	if runtime.GOOS == "darwin" {
		KeyTap("v", "command")
	} else {
		KeyTap("v", "control")
	}
}

// TypeStrDelay type string delayed
func TypeStrDelay(str string, delay int) {
	cstr := C.CString(str)
	cdelay := C.size_t(delay)
	C.type_string_delayed(cstr, cdelay)

	C.free(unsafe.Pointer(cstr))
}

// TypeStringDelayed type string delayed, Wno-deprecated
func TypeStringDelayed(str string, delay int) {
	TypeStrDelay(str, delay)
}

// SetKeyDelay set keyboard delay
func SetKeyDelay(delay int) {
	C.set_keyboard_delay(C.size_t(delay))
}

// SetKeyboardDelay set keyboard delay, Wno-deprecated
func SetKeyboardDelay(delay int) {
	C.set_keyboard_delay(C.size_t(delay))
}

/*
.______    __  .___________..___  ___.      ___      .______
|   _  \  |  | |           ||   \/   |     /   \     |   _  \
|  |_)  | |  | `---|  |----`|  \  /  |    /  ^  \    |  |_)  |
|   _  <  |  |     |  |     |  |\/|  |   /  /_\  \   |   ___/
|  |_)  | |  |     |  |     |  |  |  |  /  _____  \  |  |
|______/  |__|     |__|     |__|  |__| /__/     \__\ | _|
*/

// GetText get the image text by tesseract ocr
//
// robotgo.GetText(imgPath, lang string)
func GetText(imgPath string, args ...string) (string, error) {
	var lang = "eng"

	if len(args) > 0 {
		lang = args[0]
		if lang == "zh" {
			lang = "chi_sim"
		}
	}

	body, err := exec.Command("tesseract", imgPath,
		"stdout", "-l", lang).Output()
	if err != nil {
		return "", err
	}
	return string(body), nil
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

// ToCBitmap trans Bitmap to C.MMBitmapRef
func ToCBitmap(bit Bitmap) C.MMBitmapRef {
	cbitmap := C.createMMBitmap(
		(*C.uint8_t)(bit.ImageBuffer),
		C.size_t(bit.Width),
		C.size_t(bit.Height),
		C.size_t(bit.Bytewidth),
		C.uint8_t(bit.BitsPerPixel),
		C.uint8_t(bit.BytesPerPixel),
	)

	return cbitmap
}

// ToBitmapBytes saves Bitmap to bitmap format in bytes
func ToBitmapBytes(bit C.MMBitmapRef) []byte {
	var len C.size_t
	ptr := C.saveMMBitmapAsBytes(bit, &len)
	if int(len) < 0 {
		return nil
	}
	bs := C.GoBytes(unsafe.Pointer(ptr), C.int(len))
	C.free(unsafe.Pointer(ptr))
	return bs
}


// ToMMBitmapRef trans CBitmap to C.MMBitmapRef
func ToMMBitmapRef(bit CBitmap) C.MMBitmapRef {
	return C.MMBitmapRef(bit)
}

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

func internalFindBitmap(bit, sbit C.MMBitmapRef, tolerance float64) (int, int) {
	pos := C.find_bitmap(bit, sbit, C.float(tolerance))
	// fmt.Println("pos----", pos)
	return int(pos.x), int(pos.y)
}

// FindCBitmap find bitmap's pos by CBitmap
func FindCBitmap(bmp CBitmap, args ...interface{}) (int, int) {
	return FindBitmap(ToMMBitmapRef(bmp), args...)
}

// FindBitmap find the bitmap's pos
//
//	robotgo.FindBitmap(bitmap, subbitamp C.MMBitmapRef, tolerance float64)
//
//  |tolerance| should be in the range 0.0f - 1.0f, denoting how closely the
//  colors in the bitmaps need to match, with 0 being exact and 1 being any.
//
// This method only automatically free the internal bitmap,
// use `defer robotgo.FreeBitmap(bit)` to free the bitmap
func FindBitmap(bit C.MMBitmapRef, args ...interface{}) (int, int) {
	var (
		sbit      C.MMBitmapRef
		tolerance = 0.01
	)

	if len(args) > 0 {
		sbit = args[0].(C.MMBitmapRef)
	} else {
		sbit = CaptureScreen()
	}

	if len(args) > 1 {
		tolerance = args[1].(float64)
	}

	fx, fy := internalFindBitmap(bit, sbit, tolerance)
	// FreeBitmap(bit)
	if len(args) <= 0 {
		FreeBitmap(sbit)
	}

	return fx, fy
}

// FindPic finding the image by path
//
//	robotgo.FindPic(path string, subbitamp C.MMBitmapRef, tolerance float64)
//
// This method only automatically free the internal bitmap,
// use `defer robotgo.FreeBitmap(bit)` to free the bitmap
func FindPic(path string, args ...interface{}) (int, int) {
	var (
		sbit      C.MMBitmapRef
		tolerance = 0.01
	)

	openbit := OpenBitmap(path)

	if len(args) > 0 {
		sbit = args[0].(C.MMBitmapRef)
	} else {
		sbit = CaptureScreen()
	}

	if len(args) > 1 {
		tolerance = args[1].(float64)
	}

	fx, fy := internalFindBitmap(openbit, sbit, tolerance)
	FreeBitmap(openbit)
	if len(args) <= 0 {
		FreeBitmap(sbit)
	}

	return fx, fy
}

// FindEveryBitmap find the every bitmap
func FindEveryBitmap(bit C.MMBitmapRef, args ...interface{}) (int, int) {
	var (
		sbit      C.MMBitmapRef
		tolerance C.float = 0.01
		lpos      C.MMPoint
	)

	if len(args) > 0 {
		sbit = args[0].(C.MMBitmapRef)
	} else {
		sbit = CaptureScreen()
	}

	if len(args) > 1 {
		tolerance = C.float(args[1].(float64))
	}

	if len(args) > 2 {
		lpos.x = C.size_t(args[2].(int))
		lpos.y = 0
	} else {
		lpos.x = 0
		lpos.y = 0
	}

	if len(args) > 3 {
		lpos.x = C.size_t(args[2].(int))
		lpos.y = C.size_t(args[3].(int))
	}

	pos := C.find_every_bitmap(bit, sbit, tolerance, &lpos)
	// FreeBitmap(bit)
	if len(args) <= 0 {
		FreeBitmap(sbit)
	}

	// fmt.Println("pos----", pos)
	return int(pos.x), int(pos.y)
}

// CountBitmap count of the bitmap
func CountBitmap(bitmap, sbit C.MMBitmapRef, args ...float32) int {
	var tolerance C.float = 0.01
	if len(args) > 0 {
		tolerance = C.float(args[0])
	}

	count := C.count_of_bitmap(bitmap, sbit, tolerance)
	return int(count)
}

// BitmapClick find the bitmap and click
func BitmapClick(bitmap C.MMBitmapRef, args ...interface{}) {
	x, y := FindBitmap(bitmap)
	MovesClick(x, y, args...)
}

// PointInBounds bitmap point in bounds
func PointInBounds(bitmap C.MMBitmapRef, x, y int) bool {
	var point C.MMPoint
	point.x = C.size_t(x)
	point.y = C.size_t(y)
	cbool := C.point_in_bounds(bitmap, point)

	return bool(cbool)
}

// OpenBitmap open the bitmap return C.MMBitmapRef
//
// robotgo.OpenBitmap(path string, type int)
func OpenBitmap(gpath string, args ...int) C.MMBitmapRef {
	path := C.CString(gpath)
	var mtype C.uint16_t = 1

	if len(args) > 0 {
		mtype = C.uint16_t(args[0])
	}

	bit := C.bitmap_open(path, mtype)
	C.free(unsafe.Pointer(path))

	return bit
}

// DecodeImg decode the image to image.Image and return
func DecodeImg(path string) (image.Image, string, error) {
	return imgo.DecodeFile(path)
}

// OpenImg open the image return []byte
func OpenImg(path string) []byte {
	return imgo.ImgToBytes(path)
}

// BitmapStr bitmap from string
func BitmapStr(str string) C.MMBitmapRef {
	cs := C.CString(str)
	bit := C.bitmap_from_string(cs)
	C.free(unsafe.Pointer(cs))

	return bit
}

// SaveBitmap save the bitmap to image
//
// robotgo.SaveBimap(bitmap C.MMBitmapRef, path string, type int)
func SaveBitmap(bitmap C.MMBitmapRef, gpath string, args ...int) string {
	var mtype C.uint16_t = 1
	if len(args) > 0 {
		mtype = C.uint16_t(args[0])
	}

	path := C.CString(gpath)
	saveBit := C.bitmap_save(bitmap, path, mtype)
	C.free(unsafe.Pointer(path))

	return C.GoString(saveBit)
}

// GetPortion get bitmap portion
func GetPortion(bit C.MMBitmapRef, x, y, w, h int) C.MMBitmapRef {
	var rect C.MMRect
	rect.origin.x = C.size_t(x)
	rect.origin.y = C.size_t(y)
	rect.size.width = C.size_t(w)
	rect.size.height = C.size_t(h)

	pos := C.get_portion(bit, rect)
	return pos
}

// Convert convert the bitmap
//
// robotgo.Convert(opath, spath string, type int)
func Convert(opath, spath string, args ...int) {
	var mtype = 1
	if len(args) > 0 {
		mtype = args[0]
	}

	// C.CString()
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

// CopyBitPB copy bitmap to pasteboard
func CopyBitPB(bitmap C.MMBitmapRef) bool {
	abool := C.bitmap_copy_to_pboard(bitmap)
	gbool := bool(abool)

	return gbool
}

// CopyBitpb copy bitmap to pasteboard, Wno-deprecated
func CopyBitpb(bitmap C.MMBitmapRef) bool {
	return CopyBitPB(bitmap)
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

// GetColors get bitmap color retrun string
func GetColors(bitmap C.MMBitmapRef, x, y int) string {
	clo := GetColor(bitmap, x, y)

	return PadHex(clo)
}

// FindColor find bitmap color
//
// robotgo.FindColor(color CHex, bitmap C.MMBitmapRef, tolerance float)
func FindColor(color CHex, args ...interface{}) (int, int) {
	var (
		tolerance C.float = 0.01
		bitmap    C.MMBitmapRef
	)

	if len(args) > 0 {
		bitmap = args[0].(C.MMBitmapRef)
	} else {
		bitmap = CaptureScreen()
	}

	if len(args) > 1 {
		tolerance = C.float(args[1].(float64))
	}

	pos := C.bitmap_find_color(bitmap, C.MMRGBHex(color), tolerance)
	if len(args) <= 0 {
		FreeBitmap(bitmap)
	}

	x := int(pos.x)
	y := int(pos.y)

	return x, y
}

// FindColorCS findcolor by CaptureScreen
func FindColorCS(color CHex, x, y, w, h int, args ...float64) (int, int) {
	var tolerance = 0.01

	if len(args) > 0 {
		tolerance = args[0]
	}

	bitmap := CaptureScreen(x, y, w, h)
	rx, ry := FindColor(color, bitmap, tolerance)
	FreeBitmap(bitmap)

	return rx, ry
}

// CountColor count bitmap color
func CountColor(color CHex, args ...interface{}) int {
	var (
		tolerance C.float = 0.01
		bitmap    C.MMBitmapRef
	)

	if len(args) > 0 {
		bitmap = args[0].(C.MMBitmapRef)
	} else {
		bitmap = CaptureScreen()
	}

	if len(args) > 1 {
		tolerance = C.float(args[1].(float64))
	}

	count := C.bitmap_count_of_color(bitmap, C.MMRGBHex(color), tolerance)
	if len(args) <= 0 {
		FreeBitmap(bitmap)
	}

	return int(count)
}

// CountColorCS count bitmap color by CaptureScreen
func CountColorCS(color CHex, x, y, w, h int, args ...float64) int {
	var tolerance = 0.01

	if len(args) > 0 {
		tolerance = args[0]
	}

	bitmap := CaptureScreen(x, y, w, h)
	rx := CountColor(color, bitmap, tolerance)
	FreeBitmap(bitmap)

	return rx
}

// GetImgSize get the image size
func GetImgSize(imgPath string) (int, int) {
	bitmap := OpenBitmap(imgPath)
	gbit := ToBitmap(bitmap)

	w := gbit.Width / 2
	h := gbit.Height / 2
	FreeBitmap(bitmap)

	return w, h
}

/*
 ___________    ____  _______ .__   __. .___________.
|   ____\   \  /   / |   ____||  \ |  | |           |
|  |__   \   \/   /  |  |__   |   \|  | `---|  |----`
|   __|   \      /   |   __|  |  . `  |     |  |
|  |____   \    /    |  |____ |  |\   |     |  |
|_______|   \__/     |_______||__| \__|     |__|
*/

// AddEvent add event listener,
//
// parameters for the string type,
// the keyboard corresponding key parameters,
//
// mouse arguments: mleft, center, mright, wheelDown, wheelUp,
// wheelLeft, wheelRight.
func AddEvent(key string) bool {
	var (
		// cs   *C.char
		mArr = []string{"mleft", "center", "mright", "wheelDown",
			"wheelUp", "wheelLeft", "wheelRight"}
		mouseBool bool
	)

	for i := 0; i < len(mArr); i++ {
		if key == mArr[i] {
			mouseBool = true
		}
	}

	if len(key) > 1 && !mouseBool {
		key = strconv.Itoa(keycode[key].(int))
	}

	geve := hook.AddEvent(key)
	// defer C.free(unsafe.Pointer(cs))

	if geve == 0 {
		return true
	}

	return false
}

// StopEvent stop event listener
func StopEvent() {
	hook.StopEvent()
}

// Start start global event hook
// return event channel
func Start() chan hook.Event {
	return hook.Start()
}

// End removes global event hook
func End() {
	hook.End()
}

// AddEvents add global event hook
func AddEvents(key string, arr ...string) bool {
	s := hook.Start()
	// defer hook.End()

	ct := false
	k := 0
	for {
		e := <-s

		l := len(arr)
		if l > 0 {
			for i := 0; i < l; i++ {
				ukey := uint16(keycode[arr[i]].(int))

				if e.Kind == hook.KeyHold && e.Keycode == ukey {
					k++
				}

				if k == l {
					ct = true
				}

				if e.Kind == hook.KeyUp && e.Keycode == ukey {
					k--
					time.Sleep(10 * time.Microsecond)
					ct = false
				}
			}
		} else {
			ct = true
		}

		if ct && e.Kind == hook.KeyUp && e.Keycode == uint16(keycode[key].(int)) {
			hook.End()
			// k = 0
			break
		}
	}

	return true
}

// AddMouse add mouse event hook
func AddMouse(btn string, x ...int16) bool {
	s := hook.Start()
	ukey := mouseMap[btn]

	ct := false
	for {
		e := <-s

		if len(x) > 1 {
			if e.Kind == hook.MouseMove && e.X == x[0] && e.Y == x[1] {
				ct = true
			}
		} else {
			ct = true
		}

		if ct && e.Kind == hook.MouseDown && int(e.Button) == ukey {
			hook.End()
			break
		}
	}

	return true
}

// AddMousePos add listen mouse event pos hook
func AddMousePos(x, y int16) bool {
	s := hook.Start()

	for {
		e := <-s
		if e.Kind == hook.MouseMove && e.X == x && e.Y == y {
			hook.End()
			break
		}
	}

	return true
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
		defaultButton = "Ok"
		cancelButton  = "Cancel"
	)

	if len(args) > 0 {
		// title = args[0]
		// msg = args[1]
		defaultButton = args[0]
	}

	if len(args) > 1 {
		cancelButton = args[1]
	}

	atitle := C.CString(title)
	amsg := C.CString(msg)
	adefaultButton := C.CString(defaultButton)
	acancelButton := C.CString(cancelButton)

	cbool := C.show_alert(atitle, amsg, adefaultButton, acancelButton)
	ibool := int(cbool)

	C.free(unsafe.Pointer(atitle))
	C.free(unsafe.Pointer(amsg))
	C.free(unsafe.Pointer(adefaultButton))
	C.free(unsafe.Pointer(acancelButton))

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

// MinWindow set the window min
func MinWindow(pid int32, args ...interface{}) {
	var (
		state = true
		hwnd  int
	)

	if len(args) > 0 {
		state = args[0].(bool)
	}
	if len(args) > 1 {
		hwnd = args[1].(int)
	}

	C.min_window(C.uintptr(pid), C.bool(state), C.uintptr(hwnd))
}

// MaxWindow set the window max
func MaxWindow(pid int32, args ...interface{}) {
	var (
		state = true
		hwnd  int
	)

	if len(args) > 0 {
		state = args[0].(bool)
	}
	if len(args) > 1 {
		hwnd = args[1].(int)
	}

	C.max_window(C.uintptr(pid), C.bool(state), C.uintptr(hwnd))
}

// CloseWindow close the window
func CloseWindow(args ...int32) {
	if len(args) <= 0 {
		C.close_main_window()
		return
	}

	var hwnd, isHwnd int32
	if len(args) > 0 {
		hwnd = args[0]
	}
	if len(args) > 1 {
		isHwnd = args[1]
	}

	C.close_window(C.uintptr(hwnd), C.uintptr(isHwnd))
}

// SetHandle set the window handle
func SetHandle(hwnd int) {
	chwnd := C.uintptr(hwnd)
	C.set_handle(chwnd)
}

// SetHandlePid set the window handle by pid
func SetHandlePid(pid int32, args ...int32) {
	var isHwnd int32
	if len(args) > 0 {
		isHwnd = args[0]
	}

	C.set_handle_pid_mData(C.uintptr(pid), C.uintptr(isHwnd))
}

// GetHandPid get handle mdata by pid
func GetHandPid(pid int32, args ...int32) C.MData {
	var isHwnd int32
	if len(args) > 0 {
		isHwnd = args[0]
	}

	return C.set_handle_pid(C.uintptr(pid), C.uintptr(isHwnd))
}

// GetHandle get the window handle
func GetHandle() int {
	hwnd := C.get_handle()
	ghwnd := int(hwnd)
	// fmt.Println("gethwnd---", ghwnd)
	return ghwnd
}

// GetBHandle get the window handle, Wno-deprecated
func GetBHandle() int {
	hwnd := C.bget_handle()
	ghwnd := int(hwnd)
	//fmt.Println("gethwnd---", ghwnd)
	return ghwnd
}

// GetTitle get the window title
func GetTitle(args ...int32) string {
	if len(args) <= 0 {
		title := C.get_main_title()
		gtittle := C.GoString(title)
		return gtittle
	}

	var hwnd, isHwnd int32
	if len(args) > 0 {
		hwnd = args[0]
	}
	if len(args) > 1 {
		isHwnd = args[1]
	}
	title := C.get_title(C.uintptr(hwnd), C.uintptr(isHwnd))
	gtittle := C.GoString(title)

	return gtittle
}

// GetPID get the process id
func GetPID() int32 {
	pid := C.get_PID()
	return int32(pid)
}

// internalGetBounds get the window bounds
func internalGetBounds(pid int32, hwnd int) (int, int, int, int) {
	bounds := C.get_bounds(C.uintptr(pid), C.uintptr(hwnd))
	return int(bounds.X), int(bounds.Y), int(bounds.W), int(bounds.H)
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

// FindIds finds the all processes named with a subset
// of "name" (case insensitive),
// return matched IDs.
func FindIds(name string) ([]int32, error) {
	var pids []int32
	nps, err := Process()
	if err != nil {
		return pids, err
	}

	name = strings.ToLower(name)
	for i := 0; i < len(nps); i++ {
		psname := strings.ToLower(nps[i].Name)
		abool := strings.Contains(psname, name)
		if abool {
			pids = append(pids, nps[i].Pid)
		}
	}

	return pids, err
}

func internalActive(pid int32, hwnd int) {
	C.active_PID(C.uintptr(pid), C.uintptr(hwnd))
}

// ActivePID active the window by PID,
// If args[0] > 0 on the Windows platform via a window handle to active
// func ActivePID(pid int32, args ...int) {
// 	var hwnd int
// 	if len(args) > 0 {
// 		hwnd = args[0]
// 	}

// 	C.active_PID(C.uintptr(pid), C.uintptr(hwnd))
// }

// ActiveName active window by name
func ActiveName(name string) error {
	pids, err := FindIds(name)
	if err == nil && len(pids) > 0 {
		ActivePID(pids[0])
	}

	return err
}

// Kill kill the process by PID
func Kill(pid int32) error {
	ps := os.Process{Pid: int(pid)}
	return ps.Kill()
}
