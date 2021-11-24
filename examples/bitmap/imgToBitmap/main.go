//go:build go1.16
// +build go1.16

package main

import (
	_ "embed"
	"fmt"

	"github.com/go-vgo/robotgo"
	"github.com/vcaesar/imgo"
)

//go:embed test_007.jpeg
var testPng []byte

func main() {
	bit1 := robotgo.CaptureScreen(300, 300, 100, 100)
	defer robotgo.FreeBitmap(bit1)
	robotgo.SaveBitmap(bit1, "test_003.jpeg")

	m1 := robotgo.ToImage(bit1)
	fmt.Println("m1: ", m1.Bounds())
	imgo.SaveToPNG("test_01.png", m1)

	r1 := robotgo.ToRGBA(bit1)
	fmt.Println("r1: ", r1.Pix)

	bit2 := robotgo.ToCBitmap(robotgo.ImgToBitmap(m1))
	robotgo.SaveBitmap(bit2, "test_002.jpeg")

	test()
}

func test() {
	bitmap := robotgo.CaptureScreen(10, 10, 10, 10)
	defer robotgo.FreeBitmap(bitmap)

	img := robotgo.ToImage(bitmap)
	robotgo.SavePng(img, "test_1.png")

	img1, _ := robotgo.ByteToImg(testPng)
	robotgo.SaveJpeg(img1, "test_7.jpeg")

	bit2 := robotgo.ToCBitmap(robotgo.ImgToBitmap(img))
	fx, fy := robotgo.FindBitmap(bit2)
	fmt.Println("FindBitmap------ ", fx, fy)

	arr := robotgo.FindAllBitmap(bit2)
	fmt.Println("Find every bitmap: ", arr)
	robotgo.SaveBitmap(bitmap, "test.png")
}
