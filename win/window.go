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
	"path/filepath"
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
	procIsWindow     = modUser32.NewProc("IsWindow")
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

// windowFor resolves pid to a window: with isHandle (the Cgo NotPid / extra
// arg mode) pid is an HWND and must name an existing window, otherwise it is
// resolved like targetWindow.
func windowFor(pid int, isHandle bool) win.HWND {
	if !isHandle {
		return targetWindow(pid)
	}
	if pid == 0 {
		return 0
	}
	if r, _, _ := procIsWindow.Call(uintptr(pid)); r == 0 {
		return 0
	}
	return win.HWND(pid)
}

// ActivePid restores (if minimized) and brings to the foreground the first
// visible window owned by pid, or the window pid itself when isHandle.
func ActivePid(pid int, isHandle bool) error {
	if pid <= 0 {
		return ErrNotFound
	}
	hwnd := windowFor(pid, isHandle)
	if hwnd == 0 {
		return ErrNotFound
	}
	if win.IsIconic(hwnd) {
		win.ShowWindow(hwnd, win.SW_RESTORE)
	}
	if !win.SetForegroundWindow(hwnd) {
		return errActivate
	}
	return nil
}

// GetBounds returns the window rect (x, y, w, h) of pid's window (an HWND
// when isHandle); pid <= 0 selects the foreground window. It returns zeros
// when no window is found.
func GetBounds(pid int, isHandle bool) (int, int, int, int) {
	hwnd := windowFor(pid, isHandle)
	var r win.RECT
	if hwnd == 0 || !win.GetWindowRect(hwnd, &r) {
		return 0, 0, 0, 0
	}
	return int(r.Left), int(r.Top), int(r.Right - r.Left), int(r.Bottom - r.Top)
}

// GetClient returns the client area (x, y, w, h) of pid's window in screen
// coordinates (pid is an HWND when isHandle); pid <= 0 selects the
// foreground window.
func GetClient(pid int, isHandle bool) (int, int, int, int) {
	hwnd := windowFor(pid, isHandle)
	var r win.RECT
	if hwnd == 0 || !win.GetClientRect(hwnd, &r) {
		return 0, 0, 0, 0
	}
	p := win.POINT{}
	if !win.ClientToScreen(hwnd, &p) {
		return 0, 0, 0, 0
	}
	return int(p.X), int(p.Y), int(r.Right - r.Left), int(r.Bottom - r.Top)
}

// MinWindow minimizes (or restores, if the bool arg is false) a window.
func MinWindow(pid int, args ...interface{}) error {
	hwnd := targetWindow(pid)
	if hwnd == 0 {
		return ErrNotFound
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
	return nil
}

// MaxWindow maximizes (or restores, if the bool arg is false) a window.
func MaxWindow(pid int, args ...interface{}) error {
	hwnd := targetWindow(pid)
	if hwnd == 0 {
		return ErrNotFound
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
	return nil
}

// CloseWindow closes a window by posting WM_CLOSE. With no args the
// foreground window is closed; the first arg may specify a pid.
func CloseWindow(args ...int) error {
	pid := 0
	if len(args) > 0 {
		pid = args[0]
	}
	hwnd := targetWindow(pid)
	if hwnd == 0 {
		return ErrNotFound
	}
	if r, _, err := procPostMessageW.Call(uintptr(hwnd), uintptr(wmClose), 0, 0); r == 0 {
		return err
	}
	return nil
}

// CheckAccess reports whether input injection is permitted. Windows has no
// accessibility gate, so it always returns true; prompt is ignored.
func CheckAccess(prompt bool) bool { return true }

// GetActiveApp returns the foreground app's executable name, full
// executable path and pid. Fields that can not be resolved are left empty.
func GetActiveApp() (string, string, int) {
	hwnd := win.GetForegroundWindow()
	if hwnd == 0 {
		return "", "", 0
	}
	pid := windowPid(hwnd)
	if pid == 0 {
		return "", "", 0
	}
	path, err := processPath(pid)
	if err != nil {
		return "", "", pid
	}
	return filepath.Base(path), path, pid
}

// processPath returns the full executable path of pid.
func processPath(pid int) (string, error) {
	h, err := windows.OpenProcess(windows.PROCESS_QUERY_LIMITED_INFORMATION, false, uint32(pid))
	if err != nil {
		return "", err
	}
	defer windows.CloseHandle(h)

	buf := make([]uint16, windows.MAX_LONG_PATH)
	size := uint32(len(buf))
	if err := windows.QueryFullProcessImageName(h, 0, &buf[0], &size); err != nil {
		return "", err
	}
	return windows.UTF16ToString(buf[:size]), nil
}
