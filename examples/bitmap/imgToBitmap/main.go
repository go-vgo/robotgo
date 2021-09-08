package main

import (
	"fmt"

	"github.com/go-vgo/robotgo"
	"github.com/vcaesar/imgo"
)

func main() {
	bit1 := robotgo.CaptureScreen(300, 300, 100, 100)
	robotgo.SaveBitmap(bit1, "test_003.jpeg")

	m1 := robotgo.ToImage(bit1)
	fmt.Println("m1: ", m1.Bounds())
	imgo.SaveToPNG("test_01.png", m1)

	r1 := robotgo.ToRGBA(bit1)
	fmt.Println("r1: ", r1.Pix)

	bit2 := robotgo.ToCBitmap(robotgo.ImgToBitmap(m1))
	robotgo.SaveBitmap(bit2, "test_002.jpeg")
}
