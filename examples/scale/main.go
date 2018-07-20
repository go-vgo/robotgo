package main

import (
	"fmt"

	"github.com/go-vgo/robotgo"
)

func main() {
	width, height := robotgo.GetScaleSize()
	fmt.Println("GetScreenSize: ", width, height)

	bitmap := robotgo.CaptureScreen(0, 0, width, height)
	robotgo.SaveBitmap(bitmap, "test.png")

	s := robotgo.Scale()
	robotx := 35 * s / 100
	roboty := 25 * s / 100
	bit1 := robotgo.CaptureScreen(0, 0, robotx, roboty)
	robotgo.SaveBitmap(bit1, "test2.png")

	clo := robotgo.GetPixelColor(robotx, roboty)
	fmt.Println("GetPixelColor...", clo)
}
