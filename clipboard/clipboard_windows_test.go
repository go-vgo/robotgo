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
	"errors"
	"runtime"
	"testing"
	"time"

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

// A transient failure (ERROR_CLIPBOARD_NOT_OPEN right after a write on CI)
// is retried until a read succeeds.
func TestRetryReadTransient(t *testing.T) {
	calls := 0
	text, err := retryRead(func() (string, error) {
		calls++
		if calls < 3 {
			return "", windows.ERROR_CLIPBOARD_NOT_OPEN
		}
		return "s", nil
	})
	if err != nil || text != "s" || calls != 3 {
		t.Fatalf("retryRead = %q, %v after %d calls", text, err, calls)
	}
}

// A persistent failure is reported once readTimeout elapses.
func TestRetryReadTimeout(t *testing.T) {
	old := readTimeout
	readTimeout = 30 * time.Millisecond
	defer func() { readTimeout = old }()

	start := time.Now()
	calls := 0
	_, err := retryRead(func() (string, error) {
		calls++
		return "", windows.ERROR_CLIPBOARD_NOT_OPEN
	})
	if !errors.Is(err, windows.ERROR_CLIPBOARD_NOT_OPEN) || calls < 2 {
		t.Fatalf("retryRead = %v after %d calls", err, calls)
	}
	if d := time.Since(start); d > time.Second {
		t.Fatalf("retryRead took %v", d)
	}
}
