// Copyright 2016 The cauefcr Project Developers. See the COPYRIGHT
// file at the top-level directory of this distribution and at
// https://github.com/cauefcr/robotgo/blob/master/LICENSE
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
	hook "github.com/robotn/gohook"
)

func addEvent() {
	fmt.Println("--- Please press ctrl + shift + q to stop hook ---")
	robotgo.EventHook(hook.KeyDown, []string{"q", "ctrl", "shift"}, func(e hook.Event) {
		fmt.Println("ctrl-shift-q")
		robotgo.EventEnd()
	})

	fmt.Println("--- Please press w---")
	robotgo.EventHook(hook.KeyDown, []string{"w"}, func(e hook.Event) {
		fmt.Println("w")
	})

	s := robotgo.EventStart()
	<-robotgo.EventProcess(s)
}

func addMouse() {
	fmt.Println("--- Please press left mouse button to see it's position and the right mouse button to exit ---")
	robotgo.EventHook(hook.MouseDown, []string{}, func(e hook.Event) {
		if e.Button == hook.MouseMap["left"] {
			fmt.Printf("mouse left @ %v - %v\n", e.X, e.Y)
		} else if e.Button == hook.MouseMap["right"] {
			robotgo.EventEnd()
		}
	})

	s := robotgo.EventStart()
	<-robotgo.EventProcess(s)
}

func lowLevel() {
	////////////////////////////////////////////////////////////////////////////////
	// Global event listener
	////////////////////////////////////////////////////////////////////////////////
	fmt.Println("Press q to stop event gathering")
	evChan := robotgo.EventStart()
	for e := range evChan {
		fmt.Println(e)
		if e.Keychar == 'q' {
			robotgo.EventEnd()
			// break
		}
	}
}

func main() {
	fmt.Println("test begin...")
	addEvent()

	addMouse()

	lowLevel()
}
