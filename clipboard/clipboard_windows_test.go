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

//go:build windows
// +build windows

package clipboard

import (
	"runtime"
	"testing"

	"golang.org/x/sys/windows"
)

// A stale thread error left by an earlier Win32 call must not turn a
// successful GlobalUnlock into a failure (ReadAll returned
// "Thread does not have a clipboard open." on CI).
func TestUnlockIgnoresStaleLastError(t *testing.T) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	h, _, err := globalAlloc.Call(gmemMoveable, 2)
	if h == 0 {
		t.Fatal(err)
	}
	defer globalFree.Call(h)

	if l, _, err := globalLock.Call(h); l == 0 {
		t.Fatal(err)
	}

	setLastError.Call(uintptr(windows.ERROR_CLIPBOARD_NOT_OPEN))
	if err := unlock(h); err != nil {
		t.Fatalf("unlock: %v", err)
	}
}
