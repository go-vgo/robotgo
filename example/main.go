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
	// . "fmt"
	"fmt"

	"github.com/go-vgo/robotgo"
	// "go-vgo/robotgo"
)

func main() {
	//Control the keyboard
	robotgo.TypeString("Hello World") //importing "Hello World"
	robotgo.KeyTap("enter")           //Press "enter"
	robotgo.KeyTap("a", "control")
	robotgo.KeyTap("h", "command") //Hide window

	robotgo.KeyTap("i", "alt", "command")
	//Press "i", "alt", "command" Key combination
	arr := []string{"alt", "command"}
	robotgo.KeyTap("i", arr)

	robotgo.KeyTap("w", "command") //close window
	robotgo.KeyTap("m", "command") //minimize window
	robotgo.KeyTap("f1", "control")
	robotgo.KeyTap("a", "control")
	robotgo.KeyToggle("a", "down")
	robotgo.KeyToggle("a", "down", "alt")
	robotgo.KeyToggle("a", "down", "alt", "command")
	robotgo.KeyToggle("enter", "down")
	robotgo.TypeString("en")

	//Control the mouse
	robotgo.MoveMouse(100, 200)          // Move the mouse to 100, 200
	robotgo.MouseClick()                 //Click the left mouse button
	robotgo.MouseClick("right", false)   //Click the right mouse button
	robotgo.MouseClick("left", true)     //double click the left mouse button
	robotgo.ScrollMouse(10, "up")        //Scrolls the mouse either up
	robotgo.MouseToggle("down", "right") //Toggles right mouse button
	robotgo.MoveMouseSmooth(100, 200)    //Smooth move the mouse to 100, 200
	robotgo.MoveMouseSmooth(100, 200, 1.0, 100.0)
	x, y := robotgo.GetMousePos() //Gets the mouse coordinates
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

	//read the screen
	// robotgo.CaptureScreen()
	// bit_map := robotgo.CaptureScreen()
	// Println("CaptureScreen...", bit_map)
	// gbit_map := robotgo.Capture_Screen()
	// gbitMap := robotgo.Capture_Screen()
	gbitMap := robotgo.BCaptureScreen()
	fmt.Println("Capture_Screen...", gbitMap.Width)

	sx, sy := robotgo.GetScreenSize()
	//Gets the screen width and height
	fmt.Println("...", sx, sy)

	color := robotgo.GetPixelColor(100, 200)
	//Gets the pixel color at 100, 200.
	fmt.Println("color----", color, "-----------------")

	color2 := robotgo.GetPixelColor(10, 20)
	//Gets the pixel color at 10, 20.
	fmt.Println("color---", color2)

	// Bitmap
	abitMap := robotgo.CaptureScreen()
	//Gets all of the screen
	fmt.Println("a...", abitMap)

	bitmap := robotgo.CaptureScreen(100, 200, 30, 40)
	//Gets part of the screen
	fmt.Println("CaptureScreen...", bitmap)
	// Println("...", bit_map.Width, bit_map.BytesPerPixel)

	fx, fy := robotgo.FindBitmap(bitmap)
	//Searches for needle in bitmap
	fmt.Println("FindBitmap------", fx, fy)

	bitpos := robotgo.GetPortion(bitmap, 10, 10, 11, 10)
	//Returns new bitmap object created from a portion of another
	fmt.Println(bitpos)

	bitstr := robotgo.TostringBitmap(bitmap)
	//Creates bitmap from string by bit_map
	fmt.Println("bit_str...", bitstr)

	// sbit_map := robotgo.BitmapFromstring(bit_str, 2)
	// Println("...", sbit_map)

	robotgo.SaveBitmap(bitmap, "test.png")
	//Saves image to absolute filepath in the given format
	robotgo.SaveBitmap(bitmap, "test31.tif", 1)
	robotgo.Convert("test.png", "test.tif")
	//Convert image

	// open_bit := robotgo.OpenBitmap("test.tif")
	openbit := robotgo.OpenBitmap("test.tif")
	// open image bitmap
	fmt.Println("open...", openbit)

	//global event listener
	fmt.Println("---please press v---")
	eve := robotgo.AddEvent("v")

	if eve == 0 {
		fmt.Println("---you press v---", "v")
	}

	fmt.Println("---please press k---")
	keve := robotgo.AddEvent("k")
	if keve == 0 {
		fmt.Println("---you press k---", "k")
	}

	fmt.Println("---please press f1---")
	feve := robotgo.AddEvent("f1")
	if feve == 0 {
		fmt.Println("you press...", "f1")
	}

	fmt.Println("---please press left mouse button---")
	mleft := robotgo.AddEvent("mleft")
	if mleft == 0 {
		fmt.Println("---you press left mouse button---", "mleft")
	}

	// mright := robotgo.AddEvent("mright")
	// if mright == 0 {
	// 	Println("---you press right mouse button---", "mright")
	// }

	// robotgo.LStop()

	//Window Handle
	abool := robotgo.ShowAlert("hello", "robotgo") //Show Alert Window
	if abool == 0 {
		fmt.Println("ok@@@", "ok")
	}
	robotgo.ShowAlert("hello", "robotgo", "Ok", "Cancel")
	// robotgo.GetPID()
	mdata := robotgo.GetActive() //Get current Window Active
	hwnd := robotgo.GetHandle()  //Get current Window Handle
	fmt.Println("hwnd---", hwnd)
	title := robotgo.GetTitle() //Get current Window title
	fmt.Println("title-----", title)
	robotgo.CloseWindow()    //close current Window
	robotgo.SetActive(mdata) //set Window Active
}
