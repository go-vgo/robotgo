// Copyright 2016-2017 The go-vgo Project Developers. See the COPYRIGHT
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

func main() {
	////////////////////////////////////////////////////////////////////////////////
	// Window Handle
	////////////////////////////////////////////////////////////////////////////////
	abool := robotgo.ShowAlert("hello", "robotgo") // Show Alert Window
	if abool == 0 {
		fmt.Println("ok@@@", "ok")
	}
	robotgo.ShowAlert("hello", "robotgo", "Ok", "Cancel")

	pid := robotgo.GetPID() // Get the current process id
	fmt.Println("pid----", pid)

	mdata := robotgo.GetActive() // Get current Window Active

	hwnd := robotgo.GetHandle() // Get current Window Handle
	fmt.Println("hwnd---", hwnd)

	bhwnd := robotgo.GetBHandle() // Get current Window Handle
	fmt.Println("bhwnd---", bhwnd)

	title := robotgo.GetTitle() // Get current Window title
	fmt.Println("title-----", title)

	robotgo.CloseWindow()    // close current Window
	robotgo.SetActive(mdata) // set Window Active
}
