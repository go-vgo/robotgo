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

func key() {
	////////////////////////////////////////////////////////////////////////////////
	// Control the keyboard
	////////////////////////////////////////////////////////////////////////////////

	// importing "Hello World"
	robotgo.TypeString("Hello World")

	// press "enter"
	robotgo.KeyTap("enter")
	robotgo.KeyTap("a", "control")
	// hide window
	robotgo.KeyTap("h", "command")

	// press "i", "alt", "command" Key combination
	robotgo.KeyTap("i", "alt", "command")
	arr := []string{"alt", "command"}
	robotgo.KeyTap("i", arr)

	// close window
	robotgo.KeyTap("w", "command")
	// minimize window
	robotgo.KeyTap("m", "command")
	robotgo.KeyTap("f1", "control")
	robotgo.KeyTap("a", "control")

	robotgo.KeyToggle("a", "down")
	robotgo.KeyToggle("a", "down", "alt")
	robotgo.KeyToggle("a", "down", "alt", "command")
	robotgo.KeyToggle("enter", "down")

	robotgo.TypeString("en")

	// write string to clipboard
	robotgo.WriteAll("测试")
	// read string from clipboard
	text, err := robotgo.ReadAll()
	if err == nil {
		fmt.Println(text)
	}
}

func mouse() {
	////////////////////////////////////////////////////////////////////////////////
	// Control the mouse
	////////////////////////////////////////////////////////////////////////////////

	// move the mouse to 100, 200
	robotgo.MoveMouse(100, 200)

	// click the left mouse button
	robotgo.MouseClick()
	// click the right mouse button
	robotgo.MouseClick("right", false)
	// double click the left mouse button
	robotgo.MouseClick("left", true)

	// scrolls the mouse either up
	robotgo.ScrollMouse(10, "up")
	// toggles right mouse button
	robotgo.MouseToggle("down", "right")

	// smooth move the mouse to 100, 200
	robotgo.MoveMouseSmooth(100, 200)
	robotgo.MoveMouseSmooth(100, 200, 1.0, 100.0)

	// gets the mouse coordinates
	x, y := robotgo.GetMousePos()
	fmt.Println("pos:", x, y)
	if x == 456 && y == 586 {
		fmt.Println("mouse...", "586")
	}

	robotgo.MouseToggle("up")
	robotgo.MoveMouse(x, y)
	robotgo.MoveMouse(100, 200)

	for i := 0; i < 1080; i += 1000 {
		fmt.Println(i)
		robotgo.MoveMouse(800, i)
	}

}

func screen() {
	////////////////////////////////////////////////////////////////////////////////
	// Read the screen
	////////////////////////////////////////////////////////////////////////////////

	abitMap := robotgo.CaptureScreen()
	fmt.Println("abitMap...", abitMap)

	gbitMap := robotgo.BCaptureScreen()
	fmt.Println("BCaptureScreen...", gbitMap.Width)
	// fmt.Println("...", gbitmap.Width, gbitmap.BytesPerPixel)

	// gets the screen width and height
	sx, sy := robotgo.GetScreenSize()
	fmt.Println("...", sx, sy)

	// gets the pixel color at 100, 200.
	color := robotgo.GetPixelColor(100, 200)
	fmt.Println("color----", color, "-----------------")

	// gets the pixel color at 10, 20.
	color2 := robotgo.GetPixelColor(10, 20)
	fmt.Println("color---", color2)

}



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

func window() {
	////////////////////////////////////////////////////////////////////////////////
	// Window Handle
	////////////////////////////////////////////////////////////////////////////////

	// show Alert Window
	abool := robotgo.ShowAlert("hello", "robotgo")
	if abool == 0 {
		fmt.Println("ok@@@", "ok")
	}
	robotgo.ShowAlert("hello", "robotgo", "Ok", "Cancel")

	// get the current process id
	pid := robotgo.GetPID()
	fmt.Println("pid----", pid)

	// get current Window Active
	mdata := robotgo.GetActive()

	// get current Window Handle
	hwnd := robotgo.GetHandle()
	fmt.Println("hwnd---", hwnd)

	// get current Window Handle
	bhwnd := robotgo.GetBHandle()
	fmt.Println("bhwnd---", bhwnd)

	// get current Window title
	title := robotgo.GetTitle()
	fmt.Println("title-----", title)

	// set Window Active
	robotgo.SetActive(mdata)

	// find the process id by the process name
	fpid, err := robotgo.FindIds("Google")
	if err == nil {
		fmt.Println("pids...", fpid)
	}

	// determine whether the process exists
	isExist, err := robotgo.PidExists(100)
	if err == nil {
		fmt.Println("pid exists is", isExist)
	}

	// get the all process id
	pids, err := robotgo.Pids()
	if err == nil {
		fmt.Println("pids: ", pids)
	}

	// find the process name by the process id
	name, err := robotgo.FindName(100)
	if err == nil {
		fmt.Println("name: ", name)
	}

	// find the all process name
	names, err := robotgo.FindNames()
	if err == nil {
		fmt.Println("name: ", names)
	}

	// get the all process struct
	ps, err := robotgo.Process()
	if err == nil {
		fmt.Println("process: ", ps)
	}

	// close current Window
	robotgo.CloseWindow()
}

func main() {
	// Control the keyboard
	key()
	// Control the mouse
	mouse()
	// Read the screen
	screen()
	
	// Global event listener
	event()
	// Window Handle and progress
	window()
}
