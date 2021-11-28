// Copyright 2016 The go-vgo Project Developers. See the COPYRIGHT
// file at the top-level directory of this distribution and at
// https://github.com/go-vgo/robotgo/blob/master/LICENSE
//
// Licensed under the Apache License, Version 2.0 <LICENSE-APACHE or
// http://www.apache.org/licenses/LICENSE-2.0> or the MIT license
// <LICENSE-MIT or http://opensource.org/licenses/MIT>, at your
// option. This file may not be copied, modified, or distributed
// except according to those terms.

package robotgo

/*
#cgo darwin,amd64 LDFLAGS:-L${SRCDIR}/cdeps/mac/amd -lpng -lz
#cgo darwin,arm64 LDFLAGS:-L${SRCDIR}/cdeps/mac/m1 -lpng -lz
//
#cgo linux LDFLAGS: -L/usr/src -lpng -lz
//
#cgo windows,amd64 LDFLAGS: -L${SRCDIR}/cdeps/win/amd/win64 -lpng -lz
#cgo windows,386 LDFLAGS: -L${SRCDIR}/cdeps/win/amd/win32 -lpng -lz
#cgo windows,arm64 LDFLAGS:-L${SRCDIR}/cdeps/win/arm -lpng -lz
//
//#include "screen/goScreen.h"
#include "bitmap/goBitmap.h"
*/
import "C"

import (
	"unsafe"

	"github.com/vcaesar/tt"
)

/*
.______    __  .___________..___  ___.      ___      .______
|   _  \  |  | |           ||   \/   |     /   \     |   _  \
|  |_)  | |  | `---|  |----`|  \  /  |    /  ^  \    |  |_)  |
|   _  <  |  |     |  |     |  |\/|  |   /  /_\  \   |   ___/
|  |_)  | |  |     |  |     |  |  |  |  /  _____  \  |  |
|______/  |__|     |__|     |__|  |__| /__/     \__\ | _|
*/

// SaveCapture capture screen and save
func SaveCapture(spath string, args ...int) string {
	bit := CaptureScreen(args...)

	err := SaveBitmap(bit, spath)
	FreeBitmap(bit)
	return err
}

// FreeBitmapArr free and dealloc the C bitmap array
func FreeBitmapArr(bit ...C.MMBitmapRef) {
	for i := 0; i < len(bit); i++ {
		FreeBitmap(bit[i])
	}
}

// ToCBitmap trans Bitmap to C.MMBitmapRef
func ToCBitmap(bit Bitmap) C.MMBitmapRef {
	cbitmap := C.createMMBitmap(
		(*C.uint8_t)(bit.ImgBuf),
		C.size_t(bit.Width),
		C.size_t(bit.Height),
		C.size_t(bit.Bytewidth),
		C.uint8_t(bit.BitsPixel),
		C.uint8_t(bit.BytesPerPixel),
	)

	return cbitmap
}

// ToMMBitmapRef trans CBitmap to C.MMBitmapRef
func ToMMBitmapRef(bit CBitmap) C.MMBitmapRef {
	return C.MMBitmapRef(bit)
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
//	robotgo.FindBitmap(bitmap, source_bitamp C.MMBitmapRef, tolerance float64)
//
// 	|tolerance| should be in the range 0.0f - 1.0f, denoting how closely the
// 	colors in the bitmaps need to match, with 0 being exact and 1 being any.
//
// This method only automatically free the internal bitmap,
// use `defer robotgo.FreeBitmap(bit)` to free the bitmap
func FindBitmap(bit C.MMBitmapRef, args ...interface{}) (int, int) {
	var (
		sbit      C.MMBitmapRef
		tolerance = 0.01
	)

	if len(args) > 0 && args[0] != nil {
		sbit = args[0].(C.MMBitmapRef)
	} else {
		sbit = CaptureScreen()
	}

	if len(args) > 1 {
		tolerance = args[1].(float64)
	}

	fx, fy := internalFindBitmap(bit, sbit, tolerance)
	// FreeBitmap(bit)
	if len(args) <= 0 || (len(args) > 0 && args[0] == nil) {
		FreeBitmap(sbit)
	}

	return fx, fy
}

// FindPic finding the image by path
//
//	robotgo.FindPic(path string, source_bitamp C.MMBitmapRef, tolerance float64)
//
// This method only automatically free the internal bitmap,
// use `defer robotgo.FreeBitmap(bit)` to free the bitmap
func FindPic(path string, args ...interface{}) (int, int) {
	var (
		sbit      C.MMBitmapRef
		tolerance = 0.01
	)

	openbit := OpenBitmap(path)

	if len(args) > 0 && args[0] != nil {
		sbit = args[0].(C.MMBitmapRef)
	} else {
		sbit = CaptureScreen()
	}

	if len(args) > 1 {
		tolerance = args[1].(float64)
	}

	fx, fy := internalFindBitmap(openbit, sbit, tolerance)
	FreeBitmap(openbit)
	if len(args) <= 0 || (len(args) > 0 && args[0] == nil) {
		FreeBitmap(sbit)
	}

	return fx, fy
}

// FreeMMPointArr free MMPoint array
func FreeMMPointArr(pointArray C.MMPointArrayRef) {
	C.destroyMMPointArray(pointArray)
}

// Deprecated: use the FindAllBitmap(),
//
// FindEveryBitmap find the every bitmap, same with the FindAllBitmap()
func FindEveryBitmap(bit C.MMBitmapRef, args ...interface{}) []Point {
	return FindAllBitmap(bit, args...)
}

// FindAllBitmap find the all bitmap
func FindAllBitmap(bit C.MMBitmapRef, args ...interface{}) (posArr []Point) {
	var (
		sbit      C.MMBitmapRef
		tolerance C.float = 0.01
		lpos      C.MMPoint
	)

	if len(args) > 0 && args[0] != nil {
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
	if len(args) <= 0 || (len(args) > 0 && args[0] == nil) {
		FreeBitmap(sbit)
	}
	if pos == nil {
		return
	}
	defer FreeMMPointArr(pos)

	cSize := pos.count
	cArray := pos.array
	gSlice := (*[(1 << 28) - 1]C.MMPoint)(unsafe.Pointer(cArray))[:cSize:cSize]
	for i := 0; i < len(gSlice); i++ {
		posArr = append(posArr, Point{
			X: int(gSlice[i].x),
			Y: int(gSlice[i].y),
		})
	}

	// fmt.Println("pos----", pos)
	return
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

// Deprecated: use the BitmapFromStr(),
//
// BitmapStr bitmap from string
func BitmapStr(str string) C.MMBitmapRef {
	return BitmapFromStr(str)
}

// BitmapFromStr read bitmap from the string
func BitmapFromStr(str string) C.MMBitmapRef {
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
func Convert(opath, spath string, args ...int) string {
	var mtype = 1
	if len(args) > 0 {
		mtype = args[0]
	}

	// C.CString()
	bitmap := OpenBitmap(opath)
	// fmt.Println("a----", bit_map)
	return SaveBitmap(bitmap, spath, mtype)
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

// Deprecated: CopyBitpb copy bitmap to pasteboard, Wno-deprecated
//
// This function will be removed in version v1.0.0
func CopyBitpb(bitmap C.MMBitmapRef) bool {
	tt.Drop("CopyBitpb", "CopyBitPB")
	return CopyBitPB(bitmap)
}

// DeepCopyBit deep copy bitmap
func DeepCopyBit(bitmap C.MMBitmapRef) C.MMBitmapRef {
	bit := C.bitmap_deepcopy(bitmap)
	return bit
}

// GetColor get the bitmap color
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

	if len(args) > 0 && args[0] != nil {
		bitmap = args[0].(C.MMBitmapRef)
	} else {
		bitmap = CaptureScreen()
	}

	if len(args) > 1 {
		tolerance = C.float(args[1].(float64))
	}

	pos := C.bitmap_find_color(bitmap, C.MMRGBHex(color), tolerance)
	if len(args) <= 0 || (len(args) > 0 && args[0] == nil) {
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

// Deprecated: use the FindAllColor(),
//
// FindEveryColor find the every color, same with the FindAllColor()
func FindEveryColor(color CHex, args ...interface{}) []Point {
	return FindAllColor(color, args...)
}

// FindAllColor find the all color
func FindAllColor(color CHex, args ...interface{}) (posArr []Point) {
	var (
		bitmap    C.MMBitmapRef
		tolerance C.float = 0.01
		lpos      C.MMPoint
	)

	if len(args) > 0 && args[0] != nil {
		bitmap = args[0].(C.MMBitmapRef)
	} else {
		bitmap = CaptureScreen()
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

	pos := C.bitmap_find_every_color(bitmap, C.MMRGBHex(color), tolerance, &lpos)
	if len(args) <= 0 || (len(args) > 0 && args[0] == nil) {
		FreeBitmap(bitmap)
	}

	if pos == nil {
		return
	}
	defer FreeMMPointArr(pos)

	cSize := pos.count
	cArray := pos.array
	gSlice := (*[(1 << 28) - 1]C.MMPoint)(unsafe.Pointer(cArray))[:cSize:cSize]
	for i := 0; i < len(gSlice); i++ {
		posArr = append(posArr, Point{
			X: int(gSlice[i].x),
			Y: int(gSlice[i].y),
		})
	}

	return
}

// CountColor count bitmap color
func CountColor(color CHex, args ...interface{}) int {
	var (
		tolerance C.float = 0.01
		bitmap    C.MMBitmapRef
	)

	if len(args) > 0 && args[0] != nil {
		bitmap = args[0].(C.MMBitmapRef)
	} else {
		bitmap = CaptureScreen()
	}

	if len(args) > 1 {
		tolerance = C.float(args[1].(float64))
	}

	count := C.bitmap_count_of_color(bitmap, C.MMRGBHex(color), tolerance)
	if len(args) <= 0 || (len(args) > 0 && args[0] == nil) {
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
