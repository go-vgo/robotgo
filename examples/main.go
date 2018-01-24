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

	robotgo.TypeString("留给真爱你的人")
	robotgo.TypeStr("留给真爱你的人")
	ustr := uint32(robotgo.CharCodeAt("留给真爱你的人", 0))
	robotgo.UnicodeType(ustr)

	robotgo.PasteStr("粘贴字符串, paste")

	// press "enter"
	robotgo.KeyTap("enter")
	robotgo.KeyTap("a", "control")
	// hide window
	robotgo.KeyTap("h", "command")
	robotgo.KeyTap("h", "command", 12)

	// press "i", "alt", "command" Key combination
	robotgo.KeyTap("i", "alt", "command")
	robotgo.KeyTap("i", "alt", "command", 11)
	arr := []string{"alt", "command"}
	robotgo.KeyTap("i", arr)
	robotgo.KeyTap("i", arr, 12)

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

	gbitMap := robotgo.GoCaptureScreen()
	fmt.Println("GoCaptureScreen...", gbitMap.Width)
	// fmt.Println("...", gbitmap.Width, gbitmap.BytesPerPixel)

	robotgo.SaveCapture("saveCapture.png", 10, 20, 100, 100)

	// gets the screen width and height
	sx, sy := robotgo.GetScreenSize()
	fmt.Println("...", sx, sy)

	// gets the pixel color at 100, 200.
	color := robotgo.GetPixelColor(100, 200)
	fmt.Println("color----", color, "-----------------")

	clo := robotgo.GetPxColor(100, 200)
	fmt.Println("color...", clo)
	clostr := robotgo.PadHex(clo)
	fmt.Println("color...", clostr)

	rgb := robotgo.RgbToHex(255, 100, 200)
	rgbstr := robotgo.PadHex(robotgo.U32ToHex(rgb))
	fmt.Println("rgb...", rgbstr)

	hex := robotgo.HexToRgb(uint32(rgb))
	fmt.Println("hex...", hex)
	hexh := robotgo.PadHex(robotgo.U8ToHex(hex))
	fmt.Println("HexToRgb...", hexh)

	// gets the pixel color at 10, 20.
	color2 := robotgo.GetPixelColor(10, 20)
	fmt.Println("color---", color2)

}

func bitmap() {
	////////////////////////////////////////////////////////////////////////////////
	// Bitmap
	////////////////////////////////////////////////////////////////////////////////

	// gets all of the screen
	abitMap := robotgo.CaptureScreen()
	fmt.Println("abitMap...", abitMap)

	// gets part of the screen
	bitmap := robotgo.CaptureScreen(100, 200, 30, 40)
	fmt.Println("CaptureScreen...", bitmap)

	gbit := robotgo.ToBitmap(bitmap)
	fmt.Println("go bitmap", gbit, gbit.Width)

	// searches for needle in bitmap
	fx, fy := robotgo.FindBit(bitmap)
	fmt.Println("FindBitmap------", fx, fy)

	color := robotgo.GetColor(bitmap, 1, 2)
	fmt.Println("color...", color)
	cx, cy := robotgo.FindColor(bitmap, robotgo.CHex(color), 1.0)
	fmt.Println("pos...", cx, cy)
	cx, cy = robotgo.FindColor(bitmap, 0xAADCDC)
	fmt.Println("pos...", cx, cy)
	cx, cy = robotgo.FindColorCS(388, 179, 300, 300, 0xAADCDC)
	fmt.Println("pos...", cx, cy)

	cnt := robotgo.CountColor(bitmap, 0xAADCDC)
	fmt.Println("count...", cnt)
	cnt1 := robotgo.CountColorCS(10, 20, 30, 40, 0xAADCDC)
	fmt.Println("count...", cnt1)

	count := robotgo.CountBitmap(abitMap, bitmap)
	fmt.Println("count...", count)

	bit := robotgo.CaptureScreen(1, 2, 40, 40)
	fmt.Println("CaptureScreen...", bit)
	fx, fy = robotgo.FindBitmap(bit)
	fmt.Println("FindBitmap------", fx, fy)
	fx, fy = robotgo.FindBitmap(bit, bitmap)
	fmt.Println("FindBitmap------", fx, fy)

	// bitmap := robotgo.CaptureScreen(10, 20, 30, 40)
	abool := robotgo.PointInBounds(bitmap, 1, 2)
	fmt.Println("point in bounds...", abool)

	// returns new bitmap object created from a portion of another
	bitpos := robotgo.GetPortion(bitmap, 10, 10, 11, 10)
	fmt.Println(bitpos)

	// creates bitmap from string by bitmap
	bitstr := robotgo.TostringBitmap(bitmap)
	fmt.Println("bitstr...", bitstr)

	// sbitmap := robotgo.BitmapFromstring(bitstr, 2)
	// fmt.Println("...", sbitmap)

	// saves image to absolute filepath in the given format
	robotgo.SaveBitmap(bitmap, "test.png")
	robotgo.SaveBitmap(bitmap, "test31.tif", 1)

	// convert image
	robotgo.Convert("test.png", "test.tif")

	// open image bitmap
	openbit := robotgo.OpenBitmap("test.tif")
	fmt.Println("openBitmap...", openbit)

	fx, fy = robotgo.FindBitmap(openbit)
	fmt.Println("FindBitmap------", fx, fy)
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
		if len(fpid) > 0 {
			robotgo.ActivePID(fpid[0])
		}
	}

	robotgo.ActiveName("chrome")

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
	ver := robotgo.GetVersion()
	fmt.Println("robotgo version", ver)

	// Control the keyboard
	key()
	// Control the mouse
	mouse()
	// Read the screen
	screen()
	// Bitmap and image processing
	bitmap()
	// Global event listener
	event()
	// Window Handle and progress
	window()
}
