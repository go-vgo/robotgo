//go:build darwin
// +build darwin

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

package darwin

import (
	"unsafe"

	"github.com/ebitengine/purego"
)

// CGFloat is a 64-bit C double on every 64-bit Apple platform (amd64/arm64).
type CGFloat = float64

// CGPoint mirrors the C struct CGPoint { CGFloat x, y; }.
type CGPoint struct {
	X, Y CGFloat
}

// CGSize mirrors the C struct CGSize { CGFloat width, height; }.
type CGSize struct {
	Width, Height CGFloat
}

// CGRect mirrors the C struct CGRect { CGPoint origin; CGSize size; }.
type CGRect struct {
	Origin CGPoint
	Size   CGSize
}

// CoreGraphics event tap location.
const kCGHIDEventTap = 0

// CoreGraphics mouse event types.
const (
	kCGEventLeftMouseDown     = 1
	kCGEventLeftMouseUp       = 2
	kCGEventRightMouseDown    = 3
	kCGEventRightMouseUp      = 4
	kCGEventMouseMoved        = 5
	kCGEventLeftMouseDragged  = 6
	kCGEventRightMouseDragged = 7
	kCGEventOtherMouseDown    = 25
	kCGEventOtherMouseUp      = 26
	kCGEventOtherMouseDragged = 27
)

// CoreGraphics mouse button numbers.
const (
	kCGMouseButtonLeft   = 0
	kCGMouseButtonRight  = 1
	kCGMouseButtonCenter = 2
)

// CoreGraphics event flag masks (modifier keys).
const (
	kCGEventFlagMaskShift     = 0x00020000
	kCGEventFlagMaskControl   = 0x00040000
	kCGEventFlagMaskAlternate = 0x00080000
	kCGEventFlagMaskCommand   = 0x00100000
)

// CoreGraphics scroll-wheel units and delta fields.
const (
	kCGScrollEventUnitPixel = 0
	kCGScrollEventUnitLine  = 1

	kCGScrollWheelEventDeltaAxis1 = 11 // vertical
	kCGScrollWheelEventDeltaAxis2 = 12 // horizontal
)

// kCGMouseEventClickState is the CGEventField holding the click state:
// 1 = single click, 2 = double click.
const kCGMouseEventClickState = 1

// CoreGraphics / Quartz function bindings (loaded via purego).
var (
	// Events.
	cgEventCreate              func(source uintptr) uintptr
	cgEventCreateMouseEvent    func(source uintptr, mouseType uint32, point CGPoint, button uint32) uintptr
	cgEventCreateKeyboardEvent func(source uintptr, keycode uint16, keyDown bool) uintptr
	// CGEventCreateScrollWheelEvent is variadic (the wheelCount varargs do not
	// survive purego's fixed-arity trampoline on arm64), so the fixed-arity
	// CGEventCreateScrollWheelEvent2 variant is bound instead.
	cgEventCreateScrollWheelEvent2 func(source uintptr, units uint32, wheelCount uint32, wheel1, wheel2, wheel3 int32) uintptr
	cgEventKeyboardSetUnicode      func(event uintptr, length uint64, str *uint16)
	cgEventSetFlags                func(event uintptr, flags uint64)
	cgEventSetIntegerValueField    func(event uintptr, field uint32, value int64)
	cgEventGetLocation             func(event uintptr) CGPoint
	cgEventPost                    func(tap uint32, event uintptr)
	cgEventPostToPid               func(pid uint32, event uintptr)

	// Displays.
	cgMainDisplayID       func() uint32
	cgDisplayPixelsWide   func(display uint32) uint64
	cgDisplayPixelsHigh   func(display uint32) uint64
	cgDisplayBounds       func(display uint32) CGRect
	cgGetActiveDisplayLst func(maxDisplays uint32, active *uint32, count *uint32) int32

	// Screen capture.
	cgDisplayCreateImageForRect func(display uint32, rect CGRect) uintptr
	cgImageGetWidth             func(image uintptr) uint64
	cgImageGetHeight            func(image uintptr) uint64
	cgImageGetBytesPerRow       func(image uintptr) uint64
	cgImageGetDataProvider      func(image uintptr) uintptr
	cgImageRelease              func(image uintptr)
	cgDataProviderCopyData      func(provider uintptr) uintptr

	// CoreFoundation.
	cfDataGetBytePtr func(data uintptr) unsafe.Pointer
	cfDataGetLength  func(data uintptr) int64
	cfRelease        func(ref uintptr)
)

// loaded reports whether the system frameworks were resolved successfully.
var loaded bool

func init() {
	cg, err := purego.Dlopen(
		"/System/Library/Frameworks/CoreGraphics.framework/CoreGraphics",
		purego.RTLD_NOW|purego.RTLD_GLOBAL)
	if err != nil {
		return
	}
	cf, err := purego.Dlopen(
		"/System/Library/Frameworks/CoreFoundation.framework/CoreFoundation",
		purego.RTLD_NOW|purego.RTLD_GLOBAL)
	if err != nil {
		return
	}

	// RegisterLibFunc panics on a missing symbol; recover so an unavailable
	// framework entry degrades to loaded=false instead of aborting the
	// importing program.
	defer func() {
		if r := recover(); r != nil {
			loaded = false
		}
	}()

	purego.RegisterLibFunc(&cgEventCreate, cg, "CGEventCreate")
	purego.RegisterLibFunc(&cgEventCreateMouseEvent, cg, "CGEventCreateMouseEvent")
	purego.RegisterLibFunc(&cgEventCreateKeyboardEvent, cg, "CGEventCreateKeyboardEvent")
	purego.RegisterLibFunc(&cgEventCreateScrollWheelEvent2, cg, "CGEventCreateScrollWheelEvent2")
	purego.RegisterLibFunc(&cgEventKeyboardSetUnicode, cg, "CGEventKeyboardSetUnicodeString")
	purego.RegisterLibFunc(&cgEventSetFlags, cg, "CGEventSetFlags")
	purego.RegisterLibFunc(&cgEventSetIntegerValueField, cg, "CGEventSetIntegerValueField")
	purego.RegisterLibFunc(&cgEventGetLocation, cg, "CGEventGetLocation")
	purego.RegisterLibFunc(&cgEventPost, cg, "CGEventPost")
	purego.RegisterLibFunc(&cgEventPostToPid, cg, "CGEventPostToPid")

	purego.RegisterLibFunc(&cgMainDisplayID, cg, "CGMainDisplayID")
	purego.RegisterLibFunc(&cgDisplayPixelsWide, cg, "CGDisplayPixelsWide")
	purego.RegisterLibFunc(&cgDisplayPixelsHigh, cg, "CGDisplayPixelsHigh")
	purego.RegisterLibFunc(&cgDisplayBounds, cg, "CGDisplayBounds")
	purego.RegisterLibFunc(&cgGetActiveDisplayLst, cg, "CGGetActiveDisplayList")

	purego.RegisterLibFunc(&cgDisplayCreateImageForRect, cg, "CGDisplayCreateImageForRect")
	purego.RegisterLibFunc(&cgImageGetWidth, cg, "CGImageGetWidth")
	purego.RegisterLibFunc(&cgImageGetHeight, cg, "CGImageGetHeight")
	purego.RegisterLibFunc(&cgImageGetBytesPerRow, cg, "CGImageGetBytesPerRow")
	purego.RegisterLibFunc(&cgImageGetDataProvider, cg, "CGImageGetDataProvider")
	purego.RegisterLibFunc(&cgImageRelease, cg, "CGImageRelease")
	purego.RegisterLibFunc(&cgDataProviderCopyData, cg, "CGDataProviderCopyData")

	purego.RegisterLibFunc(&cfDataGetBytePtr, cf, "CFDataGetBytePtr")
	purego.RegisterLibFunc(&cfDataGetLength, cf, "CFDataGetLength")
	purego.RegisterLibFunc(&cfRelease, cf, "CFRelease")

	loaded = true
}

// postEvent posts a CoreGraphics event to the global HID event tap and
// releases it.
func postEvent(event uintptr) {
	postEventTo(event, 0)
}

// postEventTo posts a CoreGraphics event and releases it, mirroring the C
// SendTo() helper in key/keypress_c.h: when pid > 0 the event is delivered
// to that specific process via CGEventPostToPid, otherwise (zero or an
// invalid negative pid, which would wrap when narrowed to uint32) it is
// posted to the global HID event tap via CGEventPost.
func postEventTo(event uintptr, pid int) {
	if event == 0 {
		return
	}
	if pid > 0 {
		cgEventPostToPid(uint32(pid), event)
	} else {
		cgEventPost(kCGHIDEventTap, event)
	}
	cfRelease(event)
}
