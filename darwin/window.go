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
	"errors"
	"strconv"
	"strings"
	"sync"
	"unsafe"

	"github.com/ebitengine/purego"
	"github.com/ebitengine/purego/objc"
)

// Window management goes through the Accessibility (AXUIElement) API, so it
// needs the Accessibility permission (see CheckAccess); app activation uses
// AppKit's NSRunningApplication.

const (
	kAXErrorSuccess     = 0
	kAXValueCGPointType = 1
	kAXValueCGSizeType  = 2

	// NSApplicationActivateIgnoringOtherApps
	nsActivateIgnoringOtherApps = 1 << 1
	// NSApplicationActivationPolicyRegular: apps that appear in the Dock.
	nsActivationPolicyRegular = 0
)

var errActivate = errors.New("robotgo: failed to activate the application")

var (
	axOnce sync.Once
	axOK   bool

	axUIElementCreateApplication  func(pid int32) uintptr
	axUIElementCopyAttributeValue func(elem, attr uintptr, value *uintptr) int32
	axUIElementSetAttributeValue  func(elem, attr, value uintptr) int32
	axUIElementPerformAction      func(elem, action uintptr) int32
	axValueGetValue               func(value uintptr, typ uint32, out unsafe.Pointer) bool

	cfRetain                      func(ref uintptr) uintptr
	cfArrayGetCount               func(arr uintptr) int64
	cfArrayGetValueAtIndex        func(arr uintptr, idx int64) uintptr
	cfStringCreateWithCStr        func(alloc uintptr, cstr string, encoding uint32) uintptr
	cfBooleanTrue, cfBooleanFalse uintptr

	// AX attribute / action names (CFStringRef, created once, never released).
	axFocusedWindow, axMainWindow, axWindows, axTitle, axPosition, axSize,
	axMinimized, axFullScreen, axCloseButton, axPress, axRaise uintptr
)

// kCFStringEncodingUTF8
const cfStringEncodingUTF8 = 0x08000100

// loadAX resolves the AXUIElement API on first use.
func loadAX() bool {
	axOnce.Do(func() {
		if !loaded || !loadApp() {
			return
		}
		as, err := purego.Dlopen(
			"/System/Library/Frameworks/ApplicationServices.framework/ApplicationServices",
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
		symTrue, err := purego.Dlsym(cf, "kCFBooleanTrue")
		if err != nil {
			return
		}
		symFalse, err := purego.Dlsym(cf, "kCFBooleanFalse")
		if err != nil {
			return
		}

		// RegisterLibFunc panics on a missing symbol.
		defer func() {
			if r := recover(); r != nil {
				axOK = false
			}
		}()
		purego.RegisterLibFunc(&axUIElementCreateApplication, as, "AXUIElementCreateApplication")
		purego.RegisterLibFunc(&axUIElementCopyAttributeValue, as, "AXUIElementCopyAttributeValue")
		purego.RegisterLibFunc(&axUIElementSetAttributeValue, as, "AXUIElementSetAttributeValue")
		purego.RegisterLibFunc(&axUIElementPerformAction, as, "AXUIElementPerformAction")
		purego.RegisterLibFunc(&axValueGetValue, as, "AXValueGetValue")
		purego.RegisterLibFunc(&cfRetain, cf, "CFRetain")
		purego.RegisterLibFunc(&cfArrayGetCount, cf, "CFArrayGetCount")
		purego.RegisterLibFunc(&cfArrayGetValueAtIndex, cf, "CFArrayGetValueAtIndex")
		purego.RegisterLibFunc(&cfStringCreateWithCStr, cf, "CFStringCreateWithCString")

		cfBooleanTrue = **(**uintptr)(unsafe.Pointer(&symTrue))
		cfBooleanFalse = **(**uintptr)(unsafe.Pointer(&symFalse))

		cfStr := func(s string) uintptr { return cfStringCreateWithCStr(0, s, cfStringEncodingUTF8) }
		axFocusedWindow = cfStr("AXFocusedWindow")
		axMainWindow = cfStr("AXMainWindow")
		axWindows = cfStr("AXWindows")
		axTitle = cfStr("AXTitle")
		axPosition = cfStr("AXPosition")
		axSize = cfStr("AXSize")
		axMinimized = cfStr("AXMinimized")
		axFullScreen = cfStr("AXFullScreen")
		axCloseButton = cfStr("AXCloseButton")
		axPress = cfStr("AXPress")
		axRaise = cfStr("AXRaise")
		axOK = true
	})
	return axOK
}

// cfBool returns kCFBooleanTrue or kCFBooleanFalse.
func cfBool(b bool) uintptr {
	if b {
		return cfBooleanTrue
	}
	return cfBooleanFalse
}

// boolArg returns args[0] when it is a bool, else def.
func boolArg(args []interface{}, def bool) bool {
	if len(args) > 0 {
		if b, ok := args[0].(bool); ok {
			return b
		}
	}
	return def
}

// appWindow returns pid's focused, main or first window (retained; the
// caller releases it), or 0 when the app has no reachable window.
func appWindow(pid int) uintptr {
	app := axUIElementCreateApplication(int32(pid))
	if app == 0 {
		return 0
	}
	defer cfRelease(app)

	for _, attr := range []uintptr{axFocusedWindow, axMainWindow} {
		var w uintptr
		if axUIElementCopyAttributeValue(app, attr, &w) == kAXErrorSuccess && w != 0 {
			return w
		}
	}

	// Minimized windows are neither focused nor main, but are listed here.
	var arr uintptr
	if axUIElementCopyAttributeValue(app, axWindows, &arr) != kAXErrorSuccess || arr == 0 {
		return 0
	}
	defer cfRelease(arr)
	if cfArrayGetCount(arr) < 1 {
		return 0
	}
	return cfRetain(cfArrayGetValueAtIndex(arr, 0))
}

// withWindow runs fn on the window of pid; pid <= 0 selects the frontmost app.
func withWindow(pid int, fn func(win uintptr) error) error {
	if !loadAX() {
		return ErrNotSupported
	}
	if pid <= 0 {
		_, _, pid = GetActiveApp()
	}
	if pid <= 0 {
		return ErrNotFound
	}
	win := appWindow(pid)
	if win == 0 {
		return ErrNotFound
	}
	defer cfRelease(win)
	return fn(win)
}

// axError converts a non-zero AXError to an error.
func axError(code int32, op string) error {
	if code == kAXErrorSuccess {
		return nil
	}
	return errors.New("robotgo: " + op + " failed, AXError " + strconv.Itoa(int(code)))
}

// GetTitle returns the title of the frontmost app's window, or of the window
// of the pid given as the first argument. It returns "" when the window is
// not reachable (e.g. without the Accessibility permission).
func GetTitle(args ...int) string {
	pid := 0
	if len(args) > 0 {
		pid = args[0]
	}
	var title string
	err := withWindow(pid, func(win uintptr) error {
		var v uintptr
		if axUIElementCopyAttributeValue(win, axTitle, &v) != kAXErrorSuccess || v == 0 {
			return ErrNotFound
		}
		defer cfRelease(v)
		// CFStringRef is toll-free bridged to NSString.
		withPool(func() { title = nsString(objc.ID(v)) })
		return nil
	})
	if err != nil {
		return ""
	}
	return title
}

// GetBounds returns the window bounds (x, y, w, h) of pid's window;
// pid <= 0 selects the frontmost app. It returns zeros when not reachable.
func GetBounds(pid int) (x, y, w, h int) {
	err := withWindow(pid, func(win uintptr) error {
		var p CGPoint
		var s CGSize
		if !axValue(win, axPosition, kAXValueCGPointType, unsafe.Pointer(&p)) ||
			!axValue(win, axSize, kAXValueCGSizeType, unsafe.Pointer(&s)) {
			return ErrNotFound
		}
		x, y, w, h = int(p.X), int(p.Y), int(s.Width), int(s.Height)
		return nil
	})
	if err != nil {
		return 0, 0, 0, 0
	}
	return x, y, w, h
}

// axValue reads an AXValue attribute of elem into out.
func axValue(elem, attr uintptr, typ uint32, out unsafe.Pointer) bool {
	var v uintptr
	if axUIElementCopyAttributeValue(elem, attr, &v) != kAXErrorSuccess || v == 0 {
		return false
	}
	defer cfRelease(v)
	return axValueGetValue(v, typ, out)
}

// ActivePid brings the application with pid to the foreground and raises
// its window.
func ActivePid(pid int) error {
	if pid <= 0 {
		return ErrNotFound
	}
	if !loadApp() {
		return ErrNotSupported
	}
	err := ErrNotFound
	withPool(func() {
		app := objc.ID(objc.GetClass("NSRunningApplication")).Send(
			objc.RegisterName("runningApplicationWithProcessIdentifier:"), int32(pid))
		if app == 0 {
			return
		}
		if !objc.Send[bool](app, objc.RegisterName("activateWithOptions:"), uint(nsActivateIgnoringOtherApps)) {
			err = errActivate
			return
		}
		err = nil
	})
	if err != nil {
		return err
	}
	// Raising the window is best effort: it needs the Accessibility permission,
	// while activating the app does not.
	if loadAX() {
		if win := appWindow(pid); win != 0 {
			axUIElementPerformAction(win, axRaise)
			cfRelease(win)
		}
	}
	return nil
}

// ActiveName activates the first regular (Dock) application whose name
// contains name (case insensitive).
func ActiveName(name string) error {
	if !loadApp() {
		return ErrNotSupported
	}
	lower := strings.ToLower(name)
	pid := 0
	withPool(func() {
		ws := objc.ID(objc.GetClass("NSWorkspace")).Send(objc.RegisterName("sharedWorkspace"))
		apps := ws.Send(objc.RegisterName("runningApplications"))
		n := objc.Send[uint](apps, objc.RegisterName("count"))
		for i := uint(0); i < n; i++ {
			app := apps.Send(objc.RegisterName("objectAtIndex:"), i)
			if objc.Send[int](app, objc.RegisterName("activationPolicy")) != nsActivationPolicyRegular {
				continue
			}
			appName := nsString(app.Send(objc.RegisterName("localizedName")))
			if strings.Contains(strings.ToLower(appName), lower) {
				pid = int(objc.Send[int32](app, objc.RegisterName("processIdentifier")))
				return
			}
		}
	})
	if pid == 0 {
		return ErrNotFound
	}
	return ActivePid(pid)
}

// setWindowAttr sets a boolean attribute on pid's window.
func setWindowAttr(pid int, attr uintptr, state bool, op string) error {
	return withWindow(pid, func(win uintptr) error {
		return axError(axUIElementSetAttributeValue(win, attr, cfBool(state)), op)
	})
}

// MinWindow minimizes (or restores, if the bool arg is false) pid's window;
// pid <= 0 selects the frontmost app.
func MinWindow(pid int, args ...interface{}) {
	if err := setWindowAttr(pid, axMinimized, boolArg(args, true), "minimize"); err != nil {
		return
	}
}

// MaxWindow enters (or exits, if the bool arg is false) full screen for pid's
// window; pid <= 0 selects the frontmost app.
func MaxWindow(pid int, args ...interface{}) {
	if err := setWindowAttr(pid, axFullScreen, boolArg(args, true), "full screen"); err != nil {
		return
	}
}

// CloseWindow closes the frontmost app's window, or the window of the pid
// given as the first argument, by pressing its close button.
func CloseWindow(args ...int) {
	pid := 0
	if len(args) > 0 {
		pid = args[0]
	}
	if err := withWindow(pid, closeWindow); err != nil {
		return
	}
}

func closeWindow(win uintptr) error {
	var btn uintptr
	if axUIElementCopyAttributeValue(win, axCloseButton, &btn) != kAXErrorSuccess || btn == 0 {
		return ErrNotFound
	}
	defer cfRelease(btn)
	return axError(axUIElementPerformAction(btn, axPress), "close")
}
