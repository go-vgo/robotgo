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
/*
#cgo LDFLAGS: -L ../event/
#include <stdio.h>
#include <stdlib.h>
#include "../event/goEvent.h"
*/
import "C"

func event() {
	////////////////////////////////////////////////////////////////////////////////
	// Global event listener
	////////////////////////////////////////////////////////////////////////////////

	C.add_event_listener()
	// stop AddEvent
	// StopEvent()
}

func main() {
	event()
}
