package main

import (
	"fmt"

	"github.com/go-vgo/robotgo"
)

func main() {
	// syscall.NewLazyDLL("user32.dll").NewProc("SetProcessDPIAware").Call()

	width, height := robotgo.GetScaleSize()
	fmt.Println("get scale screen size: ", width, height)

	bitmap := robotgo.CaptureScreen(0, 0, width, height)
	defer robotgo.FreeBitmap(bitmap)
	// bitmap.Save(bitmap, "test.png")
	robotgo.Save(robotgo.ToImage(bitmap), "test.png")

	robotgo.Scale = true
	robotgo.Move(10, 10)
	robotgo.MoveSmooth(100, 100)

	fmt.Println(robotgo.Location())

	num := robotgo.DisplaysNum()
	for i := 0; i < num; i++ {
		rect := robotgo.GetScreenRect(i)
		fmt.Println("rect: ", rect)
	}
}

func old() {
	sx := robotgo.ScaleX() // Deprecated
	s := robotgo.Scale1()  // Deprecated, use the ScaleF() function
	robotx, roboty := 35*s/100, 25*s/100
	fmt.Println("scale: ", sx, s, " pos: ", robotx, roboty)

	mx, my := robotgo.Location()
	sx, sy := mx*s/100, my*s/100

	rx, ry, rw, rh := sx, sy, robotx, roboty
	// bit1 := robotgo.CaptureScreen(10, 20, robotw, roboth)
	bit1 := robotgo.CaptureScreen(rx, ry, rw, rh)
	defer robotgo.FreeBitmap(bit1)
	// bitmap.Save(bit1, "test2.png")
	robotgo.Save(robotgo.ToImage(bit1), "test2.png")

	clo := robotgo.GetPixelColor(robotx, roboty)
	fmt.Println("GetPixelColor...", clo)
}
