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
	"math"
	"strconv"
	"strings"
	"sync"
	"unsafe"

	"github.com/ebitengine/purego"
	"github.com/ebitengine/purego/objc"
)

// Window management goes through the Accessibility (AXUIElement) API, so it
// needs the Accessibility permission (see CheckAccess). App lookups use the
// live CGWindowList / AX state rather than NSWorkspace, whose app list is
// only refreshed by a running main run loop (which a Go program lacks).

const (
	kAXErrorSuccess     = 0
	kAXValueCGPointType = 1
	kAXValueCGSizeType  = 2

	// axTimeout bounds every AX call, so a hung app can not block for the
	// 6s system default.
	axTimeout = 1.0

	// NSApplicationActivateIgnoringOtherApps (a no-op since macOS 14, where
	// AXFrontmost below does the work).
	nsActivateIgnoringOtherApps = 1 << 1
)

var errActivate = errors.New("robotgo: failed to activate the application")

var (
	axOnce sync.Once
	axOK   bool

	axUIElementCreateSystemWide    func() uintptr
	axUIElementCreateApplication   func(pid int32) uintptr
	axUIElementGetPid              func(elem uintptr, pid *int32) int32
	axUIElementSetMessagingTimeout func(elem uintptr, seconds float32) int32
	axUIElementCopyAttributeValue  func(elem, attr uintptr, value *uintptr) int32
	axUIElementSetAttributeValue   func(elem, attr, value uintptr) int32
	axUIElementPerformAction       func(elem, action uintptr) int32
	axValueGetValue                func(value uintptr, typ uint32, out unsafe.Pointer) bool
	axValueTypeID                  uint64

	cfBooleanTrue, cfBooleanFalse uintptr

	// axSystem is the system-wide element (never released).
	axSystem uintptr

	// AX attribute / action names (CFStringRef, created once, never released).
	axFocusedApplication, axFrontmost, axFocusedWindow, axMainWindow,
	axWindows, axTitle, axPosition, axSize, axMinimized, axFullScreen,
	axCloseButton, axPress, axRaise uintptr
)

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
		purego.RegisterLibFunc(&axUIElementCreateSystemWide, as, "AXUIElementCreateSystemWide")
		purego.RegisterLibFunc(&axUIElementCreateApplication, as, "AXUIElementCreateApplication")
		purego.RegisterLibFunc(&axUIElementGetPid, as, "AXUIElementGetPid")
		purego.RegisterLibFunc(&axUIElementSetMessagingTimeout, as, "AXUIElementSetMessagingTimeout")
		purego.RegisterLibFunc(&axUIElementCopyAttributeValue, as, "AXUIElementCopyAttributeValue")
		purego.RegisterLibFunc(&axUIElementSetAttributeValue, as, "AXUIElementSetAttributeValue")
		purego.RegisterLibFunc(&axUIElementPerformAction, as, "AXUIElementPerformAction")
		purego.RegisterLibFunc(&axValueGetValue, as, "AXValueGetValue")
		var valueTypeID func() uint64
		purego.RegisterLibFunc(&valueTypeID, as, "AXValueGetTypeID")
		axValueTypeID = valueTypeID()

		cfBooleanTrue = **(**uintptr)(unsafe.Pointer(&symTrue))
		cfBooleanFalse = **(**uintptr)(unsafe.Pointer(&symFalse))

		axSystem = axUIElementCreateSystemWide()
		if axSystem == 0 {
			return
		}
		// On the system-wide element this sets the global timeout.
		axUIElementSetMessagingTimeout(axSystem, axTimeout)

		cfStr := func(s string) uintptr { return cfStringCreateWithCString(0, s, cfStringEncodingUTF8) }
		axFocusedApplication = cfStr("AXFocusedApplication")
		axFrontmost = cfStr("AXFrontmost")
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

// isType reports whether the CF object ref is non-nil and of type id.
func isType(ref uintptr, id uint64) bool {
	return ref != 0 && cfGetTypeID(ref) == id
}

// cfGoString converts a CFStringRef to a Go string; non-strings yield "".
func cfGoString(ref uintptr) string {
	if !isType(ref, cfStringTypeID) {
		return ""
	}
	size := cfStringGetMaxSizeForEnc(cfStringGetLength(ref), cfStringEncodingUTF8) + 1
	if size <= 1 {
		return ""
	}
	buf := make([]byte, size)
	if !cfStringGetCString(ref, &buf[0], size, cfStringEncodingUTF8) {
		return ""
	}
	n := 0
	for n < len(buf) && buf[n] != 0 {
		n++
	}
	return string(buf[:n])
}

// cfInt reads a CFNumber as an int64.
func cfInt(ref uintptr) (int64, bool) {
	if !isType(ref, cfNumberTypeID) {
		return 0, false
	}
	var v int64
	ok := cfNumberGetValue(ref, cfNumberSInt64Type, unsafe.Pointer(&v))
	return v, ok
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

// validPid reports whether pid fits a pid_t and names a process.
func validPid(pid int) bool {
	return pid > 0 && pid <= math.MaxInt32
}

// owner is the app owning a normal (layer 0) window.
type owner struct {
	pid  int
	name string
}

// windowOwners lists the distinct owners of layer-0 windows in front-to-back
// order. It reads the live window server state and needs no permission.
func windowOwners(option uint32) []owner {
	if !loaded {
		return nil
	}
	arr := cgWindowListCopyWindowInfo(option|kCGWindowListExcludeDesktopElements, 0)
	if arr == 0 {
		return nil
	}
	defer cfRelease(arr)
	if !isType(arr, cfArrayTypeID) {
		return nil
	}

	var res []owner
	seen := map[int]bool{}
	for i, n := int64(0), cfArrayGetCount(arr); i < n; i++ {
		d := cfArrayGetValueAtIndex(arr, i)
		if !isType(d, cfDictionaryTypeID) {
			continue
		}
		if layer, ok := cfInt(cfDictionaryGetValue(d, kCGWindowLayer)); !ok || layer != 0 {
			continue
		}
		pid, ok := cfInt(cfDictionaryGetValue(d, kCGWindowOwnerPID))
		if !ok || pid <= 0 || seen[int(pid)] {
			continue
		}
		seen[int(pid)] = true
		res = append(res, owner{int(pid), cfGoString(cfDictionaryGetValue(d, kCGWindowOwnerName))})
	}
	return res
}

// frontmostPid returns the pid of the focused app: from AX when trusted,
// else the owner of the frontmost on-screen window.
func frontmostPid() int {
	if loadAX() {
		var app uintptr
		if axUIElementCopyAttributeValue(axSystem, axFocusedApplication, &app) == kAXErrorSuccess && app != 0 {
			var pid int32
			code := axUIElementGetPid(app, &pid)
			cfRelease(app)
			if code == kAXErrorSuccess && pid > 0 {
				return int(pid)
			}
		}
	}
	if owners := windowOwners(kCGWindowListOptionOnScreenOnly); len(owners) > 0 {
		return owners[0].pid
	}
	return 0
}

// copyWindows returns app's AXWindows array (retained), or 0.
func copyWindows(app uintptr) uintptr {
	var arr uintptr
	if axUIElementCopyAttributeValue(app, axWindows, &arr) != kAXErrorSuccess || arr == 0 {
		return 0
	}
	if !isType(arr, cfArrayTypeID) {
		cfRelease(arr)
		return 0
	}
	return arr
}

// isMinimized reports whether the AX window win is minimized.
func isMinimized(win uintptr) bool {
	var v uintptr
	if axUIElementCopyAttributeValue(win, axMinimized, &v) != kAXErrorSuccess || v == 0 {
		return false
	}
	defer cfRelease(v)
	return v == cfBooleanTrue
}

// appWindow returns pid's focused, main or first window (retained; the
// caller releases it), or 0 when the app has no reachable window. With
// minimized it returns the first minimized window instead, which is never
// focused or main.
func appWindow(pid int, minimized bool) uintptr {
	if !validPid(pid) {
		return 0
	}
	app := axUIElementCreateApplication(int32(pid))
	if app == 0 {
		return 0
	}
	defer cfRelease(app)

	if !minimized {
		for _, attr := range []uintptr{axFocusedWindow, axMainWindow} {
			var w uintptr
			if axUIElementCopyAttributeValue(app, attr, &w) == kAXErrorSuccess && w != 0 {
				return w
			}
		}
	}

	arr := copyWindows(app)
	if arr == 0 {
		return 0
	}
	defer cfRelease(arr)
	for i, n := int64(0), cfArrayGetCount(arr); i < n; i++ {
		w := cfArrayGetValueAtIndex(arr, i)
		if w != 0 && (!minimized || isMinimized(w)) {
			return cfRetain(w)
		}
	}
	return 0
}

// withWindow runs fn on the window of pid (its first minimized window when
// minimized); pid <= 0 selects the frontmost app.
func withWindow(pid int, minimized bool, fn func(win uintptr) error) error {
	if !loadAX() {
		return ErrNotSupported
	}
	if pid <= 0 {
		pid = frontmostPid()
	}
	if !validPid(pid) {
		return ErrNotFound
	}
	win := appWindow(pid, minimized)
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
		if args[0] <= 0 {
			return ""
		}
		pid = args[0]
	}
	var title string
	err := withWindow(pid, false, func(win uintptr) error {
		var v uintptr
		if axUIElementCopyAttributeValue(win, axTitle, &v) != kAXErrorSuccess || v == 0 {
			return ErrNotFound
		}
		defer cfRelease(v)
		title = cfGoString(v)
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
	err := withWindow(pid, false, func(win uintptr) error {
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
	return isType(v, axValueTypeID) && axValueGetValue(v, typ, out)
}

// ActivePid brings the application with pid to the foreground, restoring
// and raising its window.
func ActivePid(pid int) error {
	if !validPid(pid) {
		return ErrNotFound
	}
	running, activated := false, false
	if loadApp() {
		withPool(func() {
			app := objc.ID(objc.GetClass("NSRunningApplication")).Send(
				objc.RegisterName("runningApplicationWithProcessIdentifier:"), int32(pid))
			if app == 0 {
				return
			}
			running = true
			activated = objc.Send[bool](app, objc.RegisterName("activateWithOptions:"),
				uint(nsActivateIgnoringOtherApps))
		})
	}

	// AXFrontmost also works on macOS 14+, where activation requests from a
	// non-active process may be ignored; it needs the Accessibility permission.
	if loadAX() {
		app := axUIElementCreateApplication(int32(pid))
		if app != 0 {
			if axUIElementSetAttributeValue(app, axFrontmost, cfBooleanTrue) == kAXErrorSuccess {
				running, activated = true, true
			}
			cfRelease(app)
		}
		if win := appWindow(pid, false); win != 0 {
			if isMinimized(win) {
				axUIElementSetAttributeValue(win, axMinimized, cfBooleanFalse)
			}
			axUIElementPerformAction(win, axRaise)
			cfRelease(win)
		}
	}

	switch {
	case !running:
		return ErrNotFound
	case !activated:
		return errActivate
	}
	return nil
}

// matchOwner picks the app whose name equals, else starts with, else
// contains name (case insensitive), so "code" prefers "Code" over "Xcode".
func matchOwner(owners []owner, name string) int {
	lower := strings.ToLower(name)
	if lower == "" {
		return 0
	}
	for _, match := range []func(s string) bool{
		func(s string) bool { return s == lower },
		func(s string) bool { return strings.HasPrefix(s, lower) },
		func(s string) bool { return strings.Contains(s, lower) },
	} {
		for _, o := range owners {
			if match(strings.ToLower(o.name)) {
				return o.pid
			}
		}
	}
	return 0
}

// ActiveName activates the app whose name best matches name (exact, then
// prefix, then substring; case insensitive) among apps with a window.
func ActiveName(name string) error {
	pid := matchOwner(windowOwners(kCGWindowListOptionAll), name)
	if pid == 0 {
		return ErrNotFound
	}
	return ActivePid(pid)
}

// setWindowAttr sets a boolean attribute on pid's window.
func setWindowAttr(pid int, minimized bool, attr uintptr, state bool, op string) error {
	if pid <= 0 {
		// Unlike reads, a missing pid never targets the frontmost window.
		return ErrNotFound
	}
	return withWindow(pid, minimized, func(win uintptr) error {
		return axError(axUIElementSetAttributeValue(win, attr, cfBool(state)), op)
	})
}

// minWindow minimizes, or restores pid's minimized window when state is false.
func minWindow(pid int, state bool) error {
	return setWindowAttr(pid, !state, axMinimized, state, "minimize")
}

// MinWindow minimizes (or restores, if the bool arg is false) pid's window.
func MinWindow(pid int, args ...interface{}) {
	if err := minWindow(pid, boolArg(args, true)); err != nil {
		return
	}
}

// MaxWindow enters (or exits, if the bool arg is false) native full screen
// for pid's window.
func MaxWindow(pid int, args ...interface{}) {
	if err := setWindowAttr(pid, false, axFullScreen, boolArg(args, true), "full screen"); err != nil {
		return
	}
}

// closeWindowPid closes pid's window; pid <= 0 closes nothing.
func closeWindowPid(pid int) error {
	if pid <= 0 {
		return ErrNotFound
	}
	return withWindow(pid, false, pressClose)
}

// CloseWindow closes the frontmost app's window (no argument), or the window
// of the pid given as the first argument, by pressing its close button.
func CloseWindow(args ...int) {
	var err error
	if len(args) == 0 {
		err = withWindow(0, false, pressClose)
	} else {
		err = closeWindowPid(args[0])
	}
	if err != nil {
		return
	}
}

func pressClose(win uintptr) error {
	var btn uintptr
	if axUIElementCopyAttributeValue(win, axCloseButton, &btn) != kAXErrorSuccess || btn == 0 {
		return ErrNotFound
	}
	defer cfRelease(btn)
	return axError(axUIElementPerformAction(btn, axPress), "close")
}
