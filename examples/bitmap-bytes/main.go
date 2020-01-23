package main

import (
	"bytes"
	"log"

	"image/jpeg"
	"io/ioutil"

	"github.com/go-vgo/robotgo"
	"golang.org/x/image/bmp"
)

func main() {
	bitMap := robotgo.CaptureScreen()
	defer robotgo.FreeBitmap(bitMap)

	bs := robotgo.ToBitmapBytes(bitMap)
	img, err := bmp.Decode(bytes.NewReader(bs))
	if err != nil {
		log.Println("bmp.Decode err is: ", err)
		return
	}

	b := new(bytes.Buffer)
	err = jpeg.Encode(b, img, &jpeg.Options{Quality: 20})
	if err != nil {
		log.Println("jpeg.Encode err is: ", err)
		return
	}

	log.Println("bytes: ", b.Bytes())
	ioutil.WriteFile("out.jpg", b.Bytes(), 0666)
}
