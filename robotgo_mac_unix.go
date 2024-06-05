// Copyright 2016 The go-vgo Project Developers. See the COPYRIGHT
// file at the top-level directory of this distribution and at
// https://github.com/go-vgo/robotgo/blob/master/LICENSE
//
// Licensed under the Apache License, Version 2.0 <LICENSE-APACHE or
// http://www.apache.org/licenses/LICENSE-2.0> or the MIT license
// <LICENSE-MIT or http://opensource.org/licenses/MIT>, at your
// option. This file may not be copied, modified, or distributed
// except according to those terms.

//go:build !windows
// +build !windows

package keyboard

// ScaleF get the system scale val
func ScaleF(displayId ...int) float64 {
	f := SysScale(displayId...)
	if f == 0.0 {
		f = 1.0
	}
	return f
}
