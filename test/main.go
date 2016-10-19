package main

import (
	. "fmt"

	"github.com/go-vgo/robotgo"
)

func arobotgo() {
	x, y := robotgo.GetMousePos()
	Println("pos:", x, y)

	robotgo.MoveMouse(100, 200)

	Println(robotgo.GetPixelColor(x, y))

	color := robotgo.GetPixelColor(100, 200)
	Println("color@@@", color)

	robotgo.TypeString("Hello World")
	// robotgo.KeyTap("a", "control")
	robotgo.KeyTap("f1", "control")
	// robotgo.KeyTap("enter")
	// robotgo.KeyToggle("enter", "down")
	robotgo.TypeString("en")

	bit_map := robotgo.CaptureScreen(10, 20, 30, 40)
	Println("...", bit_map)

	// robotgo.MouseClick()
	robotgo.ScrollMouse(10, "up")
}

func main() {
	arobotgo()
}
