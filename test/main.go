package main

import (
	. "fmt"

	"github.com/go-vgo/robotgo"
)

func arobotgo() {
	x, y := robotgo.GetMousePos()
	Println("pos:", x, y)

	Println(robotgo.GetPixelColor(x, y))

	color := robotgo.GetPixelColor(100, 200)
	Println("color@@@", color)

	robotgo.TypeString("Hello World")
	robotgo.KeyTap("f1", "control")
	// robotgo.KeyTap("enter")
	// robotgo.KeyToggle("enter", "down")
	robotgo.TypeString("en")

	// robotgo.MouseClick()
	robotgo.ScrollMouse(10, "up")
}

func main() {
	arobotgo()
}
