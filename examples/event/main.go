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

import (
	"fmt"

	"github.com/go-vgo/robotgo"
	// "go-vgo/robotgo"
)

func event() {
	////////////////////////////////////////////////////////////////////////////////
	// Global event listener
	////////////////////////////////////////////////////////////////////////////////

	fmt.Println("--- Please press v---")
	eve := robotgo.AddEvent("v")

	if eve == 0 {
		fmt.Println("--- You press v---", "v")
	}

	fmt.Println("--- Please press k---")
	keve := robotgo.AddEvent("k")
	if keve == 0 {
		fmt.Println("--- You press k---", "k")
	}

	fmt.Println("--- Please press f1---")
	feve := robotgo.AddEvent("f1")
	if feve == 0 {
		fmt.Println("You press...", "f1")
	}

	fmt.Println("--- Please press left mouse button---")
	mleft := robotgo.AddEvent("mleft")
	if mleft == 0 {
		fmt.Println("--- You press left mouse button---", "mleft")
	}

	mright := robotgo.AddEvent("mright")
	if mright == 0 {
		fmt.Println("--- You press right mouse button---", "mright")
	}

	// stop AddEvent
	// robotgo.StopEvent()
}

func main() {
	event()
}
