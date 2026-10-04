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
	"runtime"
	"sync"
	"unsafe"

	"github.com/ebitengine/purego"
	"github.com/ebitengine/purego/objc"
)

var (
	appOnce sync.Once
	appOK   bool

	axIsProcessTrustedWithOptions func(options uintptr) bool
	axTrustedCheckOptionPrompt    uintptr // CFStringRef
)

// loadApp loads AppKit (NSWorkspace) and the AX trust API on first use.
func loadApp() bool {
	appOnce.Do(func() {
		if _, err := purego.Dlopen("/System/Library/Frameworks/AppKit.framework/AppKit",
			purego.RTLD_NOW|purego.RTLD_GLOBAL); err != nil {
			return
		}
		as, err := purego.Dlopen(
			"/System/Library/Frameworks/ApplicationServices.framework/ApplicationServices",
			purego.RTLD_NOW|purego.RTLD_GLOBAL)
		if err != nil {
			return
		}
		sym, err := purego.Dlsym(as, "kAXTrustedCheckOptionPrompt")
		if err != nil {
			return
		}
		fn, err := purego.Dlsym(as, "AXIsProcessTrustedWithOptions")
		if err != nil {
			return
		}
		purego.RegisterFunc(&axIsProcessTrustedWithOptions, fn)
		axTrustedCheckOptionPrompt = **(**uintptr)(unsafe.Pointer(&sym))
		appOK = true
	})
	return appOK
}

// withPool runs fn inside an NSAutoreleasePool so autoreleased objects are freed.
// The pool is thread-local, so the goroutine is pinned to its OS thread.
func withPool(fn func()) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	pool := objc.ID(objc.GetClass("NSAutoreleasePool")).Send(objc.RegisterName("new"))
	defer pool.Send(objc.RegisterName("drain"))
	fn()
}

// nsString converts an NSString to a Go string.
func nsString(s objc.ID) string {
	if s == 0 {
		return ""
	}
	p := objc.Send[*byte](s, objc.RegisterName("UTF8String"))
	if p == nil {
		return ""
	}
	n := 0
	for *(*byte)(unsafe.Add(unsafe.Pointer(p), n)) != 0 {
		n++
	}
	return string(unsafe.Slice(p, n))
}

// CheckAccess reports whether the process is trusted for Accessibility
// (needed to post input events); prompt shows the system dialog if not.
func CheckAccess(prompt bool) bool {
	if !loadApp() {
		return false
	}
	var ok bool
	withPool(func() {
		val := objc.ID(objc.GetClass("NSNumber")).Send(objc.RegisterName("numberWithBool:"), prompt)
		opts := objc.ID(objc.GetClass("NSDictionary")).Send(
			objc.RegisterName("dictionaryWithObject:forKey:"), val, objc.ID(axTrustedCheckOptionPrompt))
		ok = axIsProcessTrustedWithOptions(uintptr(opts))
	})
	return ok
}

// GetActiveApp returns the frontmost app's localized name, bundle id and pid.
func GetActiveApp() (name, bundleID string, pid int) {
	if !loadApp() {
		return "", "", 0
	}
	withPool(func() {
		ws := objc.ID(objc.GetClass("NSWorkspace")).Send(objc.RegisterName("sharedWorkspace"))
		app := ws.Send(objc.RegisterName("frontmostApplication"))
		if app == 0 {
			return
		}
		name = nsString(app.Send(objc.RegisterName("localizedName")))
		bundleID = nsString(app.Send(objc.RegisterName("bundleIdentifier")))
		pid = int(objc.Send[int32](app, objc.RegisterName("processIdentifier")))
	})
	return name, bundleID, pid
}
