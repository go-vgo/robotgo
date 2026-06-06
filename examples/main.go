// Copyright (c) 2016-2026 AtomAI go-vgo, All rights reserved.
//
// See the COPYRIGHT file at the top-level directory of this distribution and at
// https://github.com/go-vgo/robotgo/blob/master/LICENSE
//
// Licensed under the Apache License, Version 2.0 <LICENSE-APACHE or
// http://www.apache.org/licenses/LICENSE-2.0>
//
// This file may not be copied, modified, or distributed
// except according to those terms.

package main

import (
	"fmt"

	"github.com/go-vgo/robotgo"
	// "go-vgo/robotgo"
)

func main() {
	ver := robotgo.GetVersion()
	fmt.Println("robotgo version is: ", ver)

	// Control the keyboard
	// key()

	// Control the mouse
	// mouse()

	// Read the screen
	// screen()

	// Bitmap and image processing
	// bitmap()

	// Global event listener
	// event()

	// Window Handle and progress
	// window()
}
