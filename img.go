// Copyright 2016 The go-vgo Project Developers. See the COPYRIGHT
// file at the top-level directory of this distribution and at
// https://github.com/go-vgo/robotgo/blob/master/LICENSE
//
// Licensed under the Apache License, Version 2.0 <LICENSE-APACHE or
// http://www.apache.org/licenses/LICENSE-2.0> or the MIT license
// <LICENSE-MIT or http://opensource.org/licenses/MIT>, at your
// option. This file may not be copied, modified, or distributed
// except according to those terms.

package robotgo

import (
	"image"
	"os/exec"
	"unsafe"

	"github.com/vcaesar/imgo"
)

// DecodeImg decode the image to image.Image and return
func DecodeImg(path string) (image.Image, string, error) {
	return imgo.DecodeFile(path)
}

// OpenImg open the image return []byte
func OpenImg(path string) ([]byte, error) {
	return imgo.ImgToBytes(path)
}

// Read read the file return image.Image
func Read(path string) (image.Image, error) {
	return imgo.Read(path)
}

// Save create a image file with the image.Image
func Save(img image.Image, path string, quality ...int) error {
	return imgo.Save(path, img, quality...)
}

// SaveImg save the image by []byte
func SaveImg(b []byte, path string) error {
	return imgo.SaveByte(path, b)
}

// SavePng save the image by image.Image
func SavePng(img image.Image, path string) error {
	return imgo.SaveToPNG(path, img)
}

// SaveJpeg save the image by image.Image
func SaveJpeg(img image.Image, path string, quality ...int) error {
	return imgo.SaveToJpeg(path, img, quality...)
}

// ToByteImg convert image.Image to []byte
func ToByteImg(img image.Image, fm ...string) []byte {
	return imgo.ToByte(img, fm...)
}

// ToStringImg convert image.Image to string
func ToStringImg(img image.Image, fm ...string) string {
	return string(ToByteImg(img, fm...))
}

// StrToImg convert base64 string to image.Image
func StrToImg(data string) (image.Image, error) {
	return imgo.StrToImg(data)
}

// ByteToImg convert []byte to image.Image
func ByteToImg(b []byte) (image.Image, error) {
	return imgo.ByteToImg(b)
}

// ImgSize get the file image size
func ImgSize(path string) (int, int, error) {
	return imgo.GetSize(path)
}

// Width return the image.Image width
func Width(img image.Image) int {
	return img.Bounds().Max.X
}

// Height return the image.Image height
func Height(img image.Image) int {
	return img.Bounds().Max.Y
}

// RGBAToBitmap convert the standard image.RGBA to Bitmap
func RGBAToBitmap(r1 *image.RGBA) (bit Bitmap) {
	bit.Width = r1.Bounds().Size().X
	bit.Height = r1.Bounds().Size().Y
	bit.Bytewidth = r1.Stride

	src := ToUint8p(r1.Pix)
	bit.ImgBuf = src

	bit.BitsPixel = 32
	bit.BytesPerPixel = 32 / 8

	return
}

// ImgToBitmap convert the standard image.Image to Bitmap
func ImgToBitmap(m image.Image) (bit Bitmap) {
	bit.Width = m.Bounds().Size().X
	bit.Height = m.Bounds().Size().Y

	pix, stride, _ := imgo.EncodeImg(m)
	bit.Bytewidth = stride

	src := ToUint8p(pix)
	bit.ImgBuf = src
	//
	bit.BitsPixel = 32
	bit.BytesPerPixel = 32 / 8
	return
}

// ToUint8p convert the []uint8 to uint8 pointer
func ToUint8p(dst []uint8) *uint8 {
	src := make([]uint8, len(dst)+10)
	for i := 0; i <= len(dst)-4; i += 4 {
		src[i+3] = dst[i+3]
		src[i] = dst[i+2]
		src[i+1] = dst[i+1]
		src[i+2] = dst[i]
	}

	ptr := unsafe.Pointer(&src[0])
	return (*uint8)(ptr)
}

// ToRGBAGo convert Bitmap to standard image.RGBA
func ToRGBAGo(bmp1 Bitmap) *image.RGBA {
	img1 := image.NewRGBA(image.Rect(0, 0, bmp1.Width, bmp1.Height))
	img1.Pix = make([]uint8, bmp1.Bytewidth*bmp1.Height)

	copyToVUint8A(img1.Pix, bmp1.ImgBuf)
	img1.Stride = bmp1.Bytewidth
	return img1
}

func val(p *uint8, n int) uint8 {
	addr := uintptr(unsafe.Pointer(p))
	addr += uintptr(n)
	p1 := (*uint8)(unsafe.Pointer(addr))
	return *p1
}

func copyToVUint8A(dst []uint8, src *uint8) {
	for i := 0; i <= len(dst)-4; i += 4 {
		dst[i] = val(src, i+2)
		dst[i+1] = val(src, i+1)
		dst[i+2] = val(src, i)
		dst[i+3] = val(src, i+3)
	}
}

// GetText get the image text by tesseract ocr
//
// robotgo.GetText(imgPath, lang string)
func GetText(imgPath string, args ...string) (string, error) {
	var lang = "eng"

	if len(args) > 0 {
		lang = args[0]
		if lang == "zh" {
			lang = "chi_sim"
		}
	}

	body, err := exec.Command("tesseract", imgPath,
		"stdout", "-l", lang).Output()
	if err != nil {
		return "", err
	}
	return string(body), nil
}
