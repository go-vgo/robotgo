// Copyright 2016 The go-vgo Project Developers. See the COPYRIGHT
// file at the top-level directory of this distribution and at
// https://github.com/go-vgo/robotgo/blob/master/LICENSE
//
// Licensed under the Apache License, Version 2.0 <LICENSE-APACHE or
// http://www.apache.org/licenses/LICENSE-2.0> or the MIT license
// <LICENSE-MIT or http://opensource.org/licenses/MIT>, at your
// option. This file may not be copied, modified, or distributed
// except according to those terms.

package main

import "C"
import (
	"fmt"
)

//export showKeyCode
func showKeyCode(keyCode C.int) int {
	fmt.Println("show msg in go ",C.int(keyCode))
	//defer C.free(unsafe.Pointer(keyCode)) // will destruct in c
	return 1
}
