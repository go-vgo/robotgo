//go:build windows
// +build windows

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

package win

import (
	"strings"
	"sync"
	"syscall"

	"github.com/tailscale/win"
	"golang.org/x/sys/windows"
)

// Win32 window message and helpers not exposed by github.com/tailscale/win.
const wmClose = 0x0010

var (
	modUser32        = windows.NewLazySystemDLL("user32.dll")
	procEnumWindows  = modUser32.NewProc("EnumWindows")
	procPostMessageW = modUser32.NewProc("PostMessageW")
)

// enumState guards the visitor used by the single, permanently-registered
// EnumWindows callback. syscall.NewCallback allocates a callback slot that is
// never released, so we register exactly one callback and swap the visitor
// under a lock instead of allocating a new callback per call.
var enumState struct {
	mu      sync.Mutex
	visitor func(hwnd win.HWND) bool
}

var enumCallback = syscall.NewCallback(func(hwnd uintptr, lparam uintptr) uintptr {
	if enumState.visitor != nil && enumState.visitor(win.HWND(hwnd)) {
		return 1 // continue
	}
	return 0 // stop
})

// enumWindows iterates over all top-level windows. The callback returns
// false to stop enumeration early.
func enumWindows(cb func(hwnd win.HWND) bool) {
	enumState.mu.Lock()
	defer enumState.mu.Unlock()

	enumState.visitor = cb
	defer func() { enumState.visitor = nil }()
	procEnumWindows.Call(enumCallback, 0)
}

// windowTitle returns the text/title of a window.
func windowTitle(hwnd win.HWND) string {
	n := win.GetWindowTextLength(hwnd)
	if n <= 0 {
		return ""
	}
	buf := make([]uint16, n+1)
	win.GetWindowText(hwnd, &buf[0], int32(len(buf)))
	return windows.UTF16ToString(buf)
}

// windowPid returns the process ID that owns a window.
func windowPid(hwnd win.HWND) int {
	var pid uint32
	win.GetWindowThreadProcessId(hwnd, &pid)
	return int(pid)
}

// targetWindow resolves the window to operate on. A pid <= 0 selects the
// current foreground window; otherwise the first visible window owned by
// that pid is returned.
func targetWindow(pid int) win.HWND {
	if pid <= 0 {
		return win.GetForegroundWindow()
	}
	var found win.HWND
	enumWindows(func(hwnd win.HWND) bool {
		if win.IsWindowVisible(hwnd) && windowPid(hwnd) == pid {
			found = hwnd
			return false // stop
		}
		return true
	})
	return found
}

// GetTitle returns the title of the foreground window.
// Arguments are accepted for API parity but ignored.
func GetTitle(args ...int) string {
	pid := 0
	if len(args) > 0 {
		pid = args[0]
	}
	hwnd := targetWindow(pid)
	if hwnd == 0 {
		return ""
	}
	return windowTitle(hwnd)
}

// ActiveName brings the first window whose title contains name (case
// insensitive) to the foreground.
func ActiveName(name string) error {
	nameLower := strings.ToLower(name)
	var target win.HWND
	enumWindows(func(hwnd win.HWND) bool {
		if !win.IsWindowVisible(hwnd) {
			return true
		}
		if strings.Contains(strings.ToLower(windowTitle(hwnd)), nameLower) {
			target = hwnd
			return false
		}
		return true
	})
	if target == 0 {
		return ErrNotFound
	}
	win.SetForegroundWindow(target)
	return nil
}

// MinWindow minimizes (or restores, if the bool arg is false) a window.
func MinWindow(pid int, args ...interface{}) {
	hwnd := targetWindow(pid)
	if hwnd == 0 {
		return
	}
	minimize := true
	if len(args) > 0 {
		if b, ok := args[0].(bool); ok {
			minimize = b
		}
	}
	if minimize {
		win.ShowWindow(hwnd, win.SW_MINIMIZE)
	} else {
		win.ShowWindow(hwnd, win.SW_RESTORE)
	}
}

// MaxWindow maximizes (or restores, if the bool arg is false) a window.
func MaxWindow(pid int, args ...interface{}) {
	hwnd := targetWindow(pid)
	if hwnd == 0 {
		return
	}
	maximize := true
	if len(args) > 0 {
		if b, ok := args[0].(bool); ok {
			maximize = b
		}
	}
	if maximize {
		win.ShowWindow(hwnd, win.SW_MAXIMIZE)
	} else {
		win.ShowWindow(hwnd, win.SW_RESTORE)
	}
}

// CloseWindow closes a window by posting WM_CLOSE. With no args the
// foreground window is closed; the first arg may specify a pid.
func CloseWindow(args ...int) {
	pid := 0
	if len(args) > 0 {
		pid = args[0]
	}
	hwnd := targetWindow(pid)
	if hwnd == 0 {
		return
	}
	procPostMessageW.Call(uintptr(hwnd), uintptr(wmClose), 0, 0)
}
