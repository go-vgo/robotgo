// Copyright 2016 The go-vgo Project Developers. See the COPYRIGHT
// file at the top-level directory of this distribution and at
// https://github.com/go-vgo/robotgo/blob/master/LICENSE
//
// Licensed under the Apache License, Version 2.0 <LICENSE-APACHE or
// http://www.apache.org/licenses/LICENSE-2.0> or the MIT license
// <LICENSE-MIT or http://opensource.org/licenses/MIT>, at your
// option. This file may not be copied, modified, or distributed
// except according to those terms.

// +build darwin windows

package robotgo

// ActivePID active the window by PID,
// If args[0] > 0 on the Windows platform via a window handle to active
func ActivePID(pid int32, args ...int) error {
	var hwnd int
	if len(args) > 0 {
		hwnd = args[0]
	}

	internalActive(pid, hwnd)
	return nil
}
