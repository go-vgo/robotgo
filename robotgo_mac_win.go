// Copyright 2016 The go-vgo Project Developers. See the COPYRIGHT
// file at the top-level directory of this distribution and at
// https://github.com/go-vgo/robotgo/blob/master/LICENSE
//
// Licensed under the Apache License, Version 2.0 <LICENSE-APACHE or
// http://www.apache.org/licenses/LICENSE-2.0> or the MIT license
// <LICENSE-MIT or http://opensource.org/licenses/MIT>, at your
// option. This file may not be copied, modified, or distributed
// except according to those terms.

//go:build darwin || windows
// +build darwin windows

package robotgo

// GetBounds get the window bounds
func GetBounds(pid int, args ...int) (int, int, int, int) {
	var isPid int
	if len(args) > 0 || NotPid {
		isPid = 1
	}

	return internalGetBounds(pid, isPid)
}

// GetClient get the window client bounds
func GetClient(pid int, args ...int) (int, int, int, int) {
	var isPid int
	if len(args) > 0 || NotPid {
		isPid = 1
	}

	return internalGetClient(pid, isPid)
}

// internalGetTitle get the window title
func internalGetTitle(pid int, args ...int) string {
	var isPid int
	if len(args) > 0 || NotPid {
		isPid = 1
	}
	gtitle := cgetTitle(pid, isPid)

	return gtitle
}

// ActivePid active the window by PID,
//
// If args[0] > 0 on the Windows platform via a window handle to active
//
// Examples:
//
//	ids, _ := robotgo.FindIds()
//	robotgo.ActivePid(ids[0])
func ActivePid(pid int, args ...int) error {
	var isPid int
	if len(args) > 0 || NotPid {
		isPid = 1
	}

	internalActive(pid, isPid)
	return nil
}

// DisplaysNum get the count of displays
func DisplaysNum() int {
	return getNumDisplays()
}

// Alert show a alert window
// Displays alert with the attributes.
// If cancel button is not given, only the default button is displayed
//
// Examples:
//
//	robotgo.Alert("hi", "window", "ok", "cancel")
func Alert(title, msg string, args ...string) bool {
	return showAlert(title, msg, args...)
}
