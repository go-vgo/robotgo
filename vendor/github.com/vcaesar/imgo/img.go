// Copyright 2018 go-vgo authors
//
// Licensed under the Apache License, Version 2.0 (the "License"): you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

/*

Package imgo get the image info
*/
package imgo

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"io/ioutil"
	"log"
	"os"

	"golang.org/x/image/bmp"
)

var (
	br, bg, bb, ba = color.Black.RGBA()
)

func IsBlack(c color.Color) bool {
	r, g, b, a := c.RGBA()

	return r == br && g == bg && b == bb && a == ba
}

// DecodeFile decodes image file
func DecodeFile(fileName string) (image.Image, string, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, "", fmt.Errorf("%s: %s", fileName, err)
	}

	img, fm, err := image.Decode(file)
	if err != nil {
		return nil, fm, fmt.Errorf("%s: %s", fileName, err)
	}

	return img, fm, nil
}

func GetSize(imagePath string) (int, int) {
	file, err := os.Open(imagePath)
	defer file.Close()
	if err != nil {
		log.Println(err)
	}

	img, _, err := image.DecodeConfig(file)
	if err != nil {
		log.Println(imagePath, err)
	}

	w := img.Width / 2
	h := img.Height / 2

	return w, h
}

// SaveToPNG create a png file with the image.Image
func SaveToPNG(path string, img image.Image) {
	f, err := os.Create(path)
	if err != nil {
		log.Println(err)
	}
	defer f.Close()

	png.Encode(f, img)
}

func ReadPNG(path string) image.Image {
	f, err := os.Open(path)
	if err != nil {
		log.Println(err)
	}

	defer f.Close()

	img, derr := png.Decode(f)
	if derr != nil {
		log.Println(derr)
	}

	return img
}

func ModifiedTime(filePath string) (int64, error) {
	f, e := os.Stat(filePath)
	if e != nil {
		return 0, e
	}

	return f.ModTime().Unix(), nil
}

func Rename(filePath, to string) error {
	return os.Rename(filePath, to)
}

func Destroy(filePath string) error {
	return os.Remove(filePath)
}

// Encode encode image to buf
func Encode(out io.Writer, subImg image.Image, fm string) error {
	switch fm {
	case "jpeg":
		return jpeg.Encode(out, subImg, nil)
	case "png":
		return png.Encode(out, subImg)
	case "gif":
		return gif.Encode(out, subImg, &gif.Options{})
	case "bmp":
		return bmp.Encode(out, subImg)
	default:
		return errors.New("ERROR FORMAT")
	}
}

func ToString(img image.Image) (result string) {
	for row := img.Bounds().Min.Y; row < img.Bounds().Max.Y; row++ {
		for col := img.Bounds().Min.X; col < img.Bounds().Max.X; col++ {
			if IsBlack(img.At(col, row)) {
				result += "."
			} else {
				result += "O"
			}
		}

		result += "\n"
	}

	return
}

func ToBytes(img image.Image, fm string) []byte {

	buf := new(bytes.Buffer)

	err := Encode(buf, img, fm)

	if err != nil {
		log.Println(err)
	}

	return buf.Bytes()
}

func ToBytesPng(img image.Image) []byte {

	buf := new(bytes.Buffer)
	err := png.Encode(buf, img)

	if err != nil {
		log.Println(err)
	}

	return buf.Bytes()
}

func ImgToBytes(path string) []byte {
	img, fm, err := DecodeFile(path)
	if err != nil {
		log.Println("to image...", err)
		return nil
	}

	return ToBytes(img, fm)
}

func PngToBytes(path string) []byte {
	img := ReadPNG(path)

	return ToBytesPng(img)
}

// Save []byte to image
func Save(path string, dist []byte) {
	ioutil.WriteFile(path, dist, 0666)
}
