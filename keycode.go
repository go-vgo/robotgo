// Copyright 2016 The go-vgo Project Developers. See the COPYRIGHT
// file at the top-level directory of this distribution and at
// https://github.com/go-vgo/robotgo/blob/master/LICENSE
//
// Licensed under the Apache License, Version 2.0 <LICENSE-APACHE or
// http://www.apache.org/licenses/LICENSE-2.0> or the MIT license
// <LICENSE-MIT or http://opensource.org/licenses/MIT>, at your
// option. This file may not be copied, modified, or distributed
// except according to those terms.

package robotgo

import (
	"github.com/vcaesar/keycode"
)

type uMap map[string]uint16

// MouseMap robotgo hook mouse's code map
var MouseMap = keycode.MouseMap

const (
	// Mleft mouse left button
	Mleft      = "left"
	Mright     = "right"
	Center     = "center"
	WheelDown  = "wheelDown"
	WheelUp    = "wheelUp"
	WheelLeft  = "wheelLeft"
	WheelRight = "wheelRight"
)

// Keycode robotgo hook key's code map
var Keycode = keycode.Keycode

// Special is the special key map
var Special = keycode.Special
