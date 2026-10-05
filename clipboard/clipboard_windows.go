// Copyright 2013 @atotto. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build windows
// +build windows

package clipboard

import (
	"errors"
	"runtime"
	"syscall"
	"time"
	"unsafe"

	"golang.org/x/sys/windows"
)

const (
	cfUnicodetext = 13
	gmemMoveable  = 0x0002
)

var (
	user32                     = windows.NewLazySystemDLL("user32.dll")
	isClipboardFormatAvailable = user32.NewProc("IsClipboardFormatAvailable")
	openClipboard              = user32.NewProc("OpenClipboard")
	closeClipboard             = user32.NewProc("CloseClipboard")
	emptyClipboard             = user32.NewProc("EmptyClipboard")
	getClipboardData           = user32.NewProc("GetClipboardData")
	setClipboardData           = user32.NewProc("SetClipboardData")

	kernel32     = windows.NewLazySystemDLL("kernel32.dll")
	globalAlloc  = kernel32.NewProc("GlobalAlloc")
	globalFree   = kernel32.NewProc("GlobalFree")
	globalLock   = kernel32.NewProc("GlobalLock")
	globalUnlock = kernel32.NewProc("GlobalUnlock")
	globalSize   = kernel32.NewProc("GlobalSize")
	setLastError = kernel32.NewProc("SetLastError")

	errUnknown = errors.New("clipboard: unknown windows error")
)

// winErr avoids returning a non-nil "operation completed successfully" error
func winErr(err error) error {
	if err == nil || err == syscall.Errno(0) {
		return errUnknown
	}
	return err
}

// waitOpenClipboard opens the clipboard, waiting for up to a second to do so.
// The caller must hold runtime.LockOSThread until the clipboard is closed.
func waitOpenClipboard() error {
	limit := time.Now().Add(time.Second)
	for {
		r, _, err := openClipboard.Call(0)
		if r != 0 {
			return nil
		}
		if time.Now().After(limit) {
			return winErr(err)
		}
		time.Sleep(time.Millisecond)
	}
}

// utf16Ptr converts a GlobalLock address, which lives outside the Go heap
func utf16Ptr(addr uintptr) *uint16 {
	return *(**uint16)(unsafe.Pointer(&addr))
}

// unlock calls GlobalUnlock, a zero result with no error code means
// the memory object is no longer locked, which is success. The last error
// is cleared first: GetClipboardData can leave a stale code (such as
// ERROR_CLIPBOARD_NOT_OPEN) behind that would otherwise fail a good unlock.
func unlock(h uintptr) error {
	setLastError.Call(0)
	r, _, err := globalUnlock.Call(h)
	if r == 0 && err != syscall.Errno(0) {
		return err
	}
	return nil
}

// readTimeout bounds how long readAll retries a failed read.
var readTimeout = time.Second

// readAll retries the whole open-read-close cycle: right after a write,
// clipboard listeners (clipboard history, RDP) can briefly take over and
// GetClipboardData fails with ERROR_CLIPBOARD_NOT_OPEN.
func readAll() (string, error) {
	return retryRead(readOnce)
}

func retryRead(read func() (string, error)) (string, error) {
	limit := time.Now().Add(readTimeout)
	for {
		text, err := read()
		if err == nil || time.Now().After(limit) {
			return text, err
		}
		time.Sleep(10 * time.Millisecond)
	}
}

func readOnce() (string, error) {
	// OpenClipboard and CloseClipboard must run on the same OS thread
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	if err := waitOpenClipboard(); err != nil {
		return "", err
	}
	defer closeClipboard.Call()

	if r, _, _ := isClipboardFormatAvailable.Call(cfUnicodetext); r == 0 {
		return "", nil
	}

	h, _, err := getClipboardData.Call(cfUnicodetext)
	if h == 0 {
		return "", winErr(err)
	}

	// valid CF_UNICODETEXT data holds at least the UTF-16 terminator
	size, _, err := globalSize.Call(h)
	if size == 0 {
		return "", winErr(err)
	}

	l, _, err := globalLock.Call(h)
	if l == 0 {
		return "", winErr(err)
	}

	buf := unsafe.Slice(utf16Ptr(l), size/2)
	text := windows.UTF16ToString(buf)

	if err := unlock(h); err != nil {
		return "", err
	}
	return text, nil
}

func writeAll(text string) error {
	data, err := windows.UTF16FromString(text)
	if err != nil {
		return err
	}

	// "If the hMem parameter identifies a memory object, the object must have
	// been allocated using the function with the GMEM_MOVEABLE flag."
	size := uintptr(len(data)) * unsafe.Sizeof(data[0])
	h, _, err := globalAlloc.Call(gmemMoveable, size)
	if h == 0 {
		return winErr(err)
	}
	defer func() {
		if h != 0 {
			globalFree.Call(h)
		}
	}()

	l, _, err := globalLock.Call(h)
	if l == 0 {
		return winErr(err)
	}
	copy(unsafe.Slice(utf16Ptr(l), len(data)), data)
	if err := unlock(h); err != nil {
		return err
	}

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	if err := waitOpenClipboard(); err != nil {
		return err
	}
	defer closeClipboard.Call()

	if r, _, err := emptyClipboard.Call(); r == 0 {
		return winErr(err)
	}

	if r, _, err := setClipboardData.Call(cfUnicodetext, h); r == 0 {
		return winErr(err)
	}

	h = 0 // the system owns the memory now, suppress deferred cleanup
	return nil
}
