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
	"fmt"
	"log"

	"github.com/go-vgo/robotgo"
	"github.com/vcaesar/imgo"
	// "go-vgo/robotgo"
)

func toBitmap(bmp robotgo.CBitmap) {
	bitmap := robotgo.ToMMBitmapRef(bmp)

	gbit := robotgo.ToBitmap(bitmap)
	fmt.Println("go bitmap", gbit, gbit.Width)

	cbit := robotgo.ToCBitmap(gbit)
	// defer robotgo.FreeBitmap(cbit)
	log.Println("cbit == bitmap: ", cbit == bitmap)
	robotgo.SaveBitmap(cbit, "tocbitmap.png")
}

func findColor(bmp robotgo.CBitmap) {
	bitmap := robotgo.ToMMBitmapRef(bmp)

	// find the color in bitmap
	color := robotgo.GetColor(bitmap, 1, 2)
	fmt.Println("color...", color)
	cx, cy := robotgo.FindColor(robotgo.CHex(color), bitmap, 1.0)
	fmt.Println("pos...", cx, cy)
	cx, cy = robotgo.FindColor(robotgo.CHex(color))
	fmt.Println("pos...", cx, cy)

	cx, cy = robotgo.FindColor(0xAADCDC, bitmap)
	fmt.Println("pos...", cx, cy)
	cx, cy = robotgo.FindColor(0xAADCDC, nil, 0.1)
	fmt.Println("pos...", cx, cy)

	cx, cy = robotgo.FindColorCS(0xAADCDC, 388, 179, 300, 300)
	fmt.Println("pos...", cx, cy)

	cnt := robotgo.CountColor(0xAADCDC, bitmap)
	fmt.Println("count...", cnt)
	cnt1 := robotgo.CountColorCS(0xAADCDC, 10, 20, 30, 40)
	fmt.Println("count...", cnt1)
}

func bitmapString(bmp robotgo.CBitmap) {
	bitmap := robotgo.ToMMBitmapRef(bmp)

	// creates bitmap from string by bitmap
	bitstr := robotgo.TostringBitmap(bitmap)
	fmt.Println("bitstr...", bitstr)

	// sbitmap := robotgo.BitmapFromstring(bitstr, 2)
	// fmt.Println("...", sbitmap)

	// sbitmap := robotgo.BitmapStr(bitstr)
	sbitmap := robotgo.BitmapFromStr(bitstr)
	fmt.Println("bitmap str...", sbitmap)
	robotgo.SaveBitmap(sbitmap, "teststr.png")
}

func bitmapTool(bmp robotgo.CBitmap) {
	bitmap := robotgo.ToMMBitmapRef(bmp)

	// bitmap := robotgo.CaptureScreen(10, 20, 30, 40)
	abool := robotgo.PointInBounds(bitmap, 1, 2)
	fmt.Println("point in bounds...", abool)

	// returns new bitmap object created from a portion of another
	bitpos := robotgo.GetPortion(bitmap, 10, 10, 11, 10)
	fmt.Println(bitpos)

	// saves image to absolute filepath in the given format
	robotgo.SaveBitmap(bitmap, "test.png")
	robotgo.SaveBitmap(bitmap, "test31.tif", 1)
}

func decode() {
	img, name, err := robotgo.DecodeImg("test.png")
	if err != nil {
		log.Println("decode image ", err)
	}
	fmt.Println("decode test.png", img, name)

	byt := robotgo.OpenImg("test.png")
	imgo.Save("test2.png", byt)

	w, h := robotgo.GetImgSize("test.png")
	fmt.Println("image width and hight ", w, h)
	w, h = imgo.GetSize("test.png")
	fmt.Println("image width and hight ", w, h)

	// convert image
	robotgo.Convert("test.png", "test.tif")
}

func bitmapTest(bmp robotgo.CBitmap) {
	bitmap := robotgo.ToMMBitmapRef(bmp)

	bit := robotgo.CaptureScreen(1, 2, 40, 40)
	defer robotgo.FreeBitmap(bit)
	fmt.Println("CaptureScreen...", bit)

	// searches for needle in bitmap
	fx, fy := robotgo.FindBitmap(bit, bitmap)
	fmt.Println("FindBitmap------", fx, fy)

	// fx, fy := robotgo.FindBit(bitmap)
	// fmt.Println("FindBitmap------", fx, fy)

	fx, fy = robotgo.FindBitmap(bit)
	fmt.Println("FindBitmap------", fx, fy)
}

func findBitmap(bmp robotgo.CBitmap) {
	fx, fy := robotgo.FindBitmap(robotgo.ToMMBitmapRef(bmp))
	fmt.Println("findBitmap: ", fx, fy)

	fx, fy = robotgo.FindCBitmap(bmp)
	fmt.Println("findCBitmap: ", fx, fy)
	fx, fy = robotgo.FindCBitmap(bmp, nil, 0.1)
	fmt.Println("findCBitmap: ", fx, fy)

	// open image bitmap
	openbit := robotgo.OpenBitmap("test.tif")
	fmt.Println("openBitmap...", openbit)

	fx, fy = robotgo.FindBitmap(openbit)
	fmt.Println("FindBitmap------", fx, fy)

	fx, fy = robotgo.FindPic("test.tif")
	fmt.Println("FindPic------", fx, fy)
}

func bitmap() {
	////////////////////////////////////////////////////////////////////////////////
	// Bitmap
	////////////////////////////////////////////////////////////////////////////////

	// gets all of the screen
	abitMap := robotgo.CaptureScreen()
	fmt.Println("abitMap...", abitMap)

	// gets part of the screen
	bitmap := robotgo.CaptureScreen(100, 200, 30, 30)
	defer robotgo.FreeBitmap(bitmap)
	fmt.Println("CaptureScreen...", bitmap)

	cbit := robotgo.CBitmap(bitmap)
	toBitmap(cbit)

	findColor(cbit)

	count := robotgo.CountBitmap(abitMap, bitmap)
	fmt.Println("count...", count)

	bitmapTest(cbit)
	findBitmap(cbit)

	bitmapString(cbit)
	bitmapTool(cbit)

	decode()

	// free the bitmap
	robotgo.FreeBitmap(abitMap)
	// robotgo.FreeBitmap(bitmap)
}

func main() {
	bitmap()
}
