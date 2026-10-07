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
	"syscall"

	"github.com/tailscale/win"
)

// Alert shows a modal MessageBox and reports whether OK was chosen; like the
// Cgo backend it uses the stock OK/Cancel buttons (OK only when cancel is
// ""), so the ok label is ignored.
func Alert(title, msg, ok, cancel string) bool {
	t, err := syscall.UTF16PtrFromString(title)
	if err != nil {
		return false
	}
	m, err := syscall.UTF16PtrFromString(msg)
	if err != nil {
		return false
	}
	var style uint32 = win.MB_OKCANCEL
	if cancel == "" {
		style = win.MB_OK
	}
	return win.MessageBox(0, m, t, style) == win.IDOK
}

// GetHWNDByPid returns the first visible top-level window owned by pid,
// 0 when there is none.
func GetHWNDByPid(pid int) int { return int(pidWindow(pid)) }
