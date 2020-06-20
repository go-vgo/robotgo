package main

import (
	"github.com/go-vgo/robotgo"
)

func main() {
	bmp, free := loadBitmaps("start.png", "chest.png", "eat.png")
	defer free()

	for {
		clickBitmap(bmp["start.png"], false)
		clickBitmap(bmp["chest.png"], true)
		clickBitmap(bmp["eat.png"], false)
	}
}

func loadBitmaps(files ...string) (bitmaps map[string]robotgo.Bitmap, free func()) {
	freeFuncs := make([]func(), 0)
	bitmaps = make(map[string]robotgo.Bitmap)
	for _, f := range files {
		bitmap, freeFunc := readBitmap(f)
		bitmaps[f] = bitmap
		freeFuncs = append(freeFuncs, freeFunc)
	}

	free = func() {
		for key := range freeFuncs {
			freeFuncs[key]()
		}
	}
	return bitmaps, free
}

func readBitmap(file string) (bitmap robotgo.Bitmap, free func()) {
	cBitmap := robotgo.OpenBitmap(file)
	bitmap = robotgo.ToBitmap(cBitmap)
	free = func() {
		robotgo.FreeBitmap(cBitmap)
	}
	return bitmap, free
}

func clickBitmap(bmp robotgo.Bitmap, doubleClick bool) bool {
	fx, fy := robotgo.FindBitmap(robotgo.ToCBitmap(bmp))
	if fx != -1 && fy != -1 {
		robotgo.MoveMouse(fx, fy)
		robotgo.MouseClick("left", doubleClick)
		return true
	}

	return false
}
