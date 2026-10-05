// Copyright (c) 2016-2026 AtomAI, All rights reserved.
//
// See COPYRIGHT file at top-level directory of this distribution and at
// https://github.com/go-vgo/robotgo/blob/master/LICENSE
//
// Licensed under Apache License, Version 2.0 <LICENSE-APACHE or
// http://www.apache.org/licenses/LICENSE-2.0>
//
// This file may not be copied, modified, or distributed
// except according to those terms.

package robotgo

import (
	"bytes"
	"encoding/base64"
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"strconv"
	"testing"
)

func TestToRGBAGoSwapsBGRA(t *testing.T) {
	// 2x1 BGRA bitmap: pixel0 = (B1,G2,R3,A4), pixel1 = (B5,G6,R7,A8).
	buf := []uint8{1, 2, 3, 4, 5, 6, 7, 8}
	bmp := Bitmap{ImgBuf: &buf[0], Width: 2, Height: 1, Bytewidth: 8}

	img := ToRGBAGo(bmp)
	want := []uint8{3, 2, 1, 4, 7, 6, 5, 8}
	if !bytes.Equal(img.Pix, want) {
		t.Fatalf("Pix = %v, want %v", img.Pix, want)
	}
	if img.Stride != 8 || img.Rect.Dx() != 2 || img.Rect.Dy() != 1 {
		t.Fatalf("stride/rect = %d/%v", img.Stride, img.Rect)
	}
}

// testRGBA returns a w x h opaque image with distinct per-pixel values
// (opaque, so PNG's non-premultiplied storage round-trips exactly).
func testRGBA(w, h int) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			img.Set(x, y, color.NRGBA{R: uint8(10 * x), G: uint8(20 * y), B: uint8(x + y), A: 255})
		}
	}
	return img
}

func TestRGBAToBitmapRoundTrip(t *testing.T) {
	src := testRGBA(3, 2)
	bit := RGBAToBitmap(src)
	if bit.Width != 3 || bit.Height != 2 || bit.Bytewidth != src.Stride {
		t.Fatalf("bitmap size = %dx%d stride %d", bit.Width, bit.Height, bit.Bytewidth)
	}
	if bit.BitsPixel != 32 || bit.BytesPerPixel != 4 {
		t.Fatalf("bitmap depth = %d/%d", bit.BitsPixel, bit.BytesPerPixel)
	}
	// RGBAToBitmap stores BGRA; first pixel R and B must be swapped.
	if val(bit.ImgBuf, 0) != src.Pix[2] || val(bit.ImgBuf, 2) != src.Pix[0] {
		t.Fatalf("bitmap not BGRA: % x vs % x", []uint8{val(bit.ImgBuf, 0), val(bit.ImgBuf, 2)}, src.Pix[:4])
	}

	got := ToRGBAGo(bit)
	if !bytes.Equal(got.Pix, src.Pix) {
		t.Fatalf("round trip Pix = %v, want %v", got.Pix, src.Pix)
	}
}

func TestToUint8pSwapsAndPads(t *testing.T) {
	p := ToUint8p([]uint8{1, 2, 3, 4})
	for i, want := range []uint8{3, 2, 1, 4} {
		if got := val(p, i); got != want {
			t.Fatalf("byte %d = %d, want %d", i, got, want)
		}
	}
}

func TestImgToBitmap(t *testing.T) {
	bit := ImgToBitmap(testRGBA(4, 3))
	if bit.Width != 4 || bit.Height != 3 || bit.ImgBuf == nil {
		t.Fatalf("bitmap = %+v", bit)
	}
	if bit.Bytewidth < 4*3 || bit.BitsPixel != 32 || bit.BytesPerPixel != 4 {
		t.Fatalf("bitmap layout = %+v", bit)
	}
}

func TestWidthHeight(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, 5, 7))
	if Width(img) != 5 || Height(img) != 7 {
		t.Fatalf("size = %dx%d", Width(img), Height(img))
	}
}

func TestSaveReadPng(t *testing.T) {
	src := testRGBA(4, 3)
	dir := t.TempDir()

	for _, save := range []struct {
		name string
		fn   func(string) error
	}{
		{"save.png", func(p string) error { return Save(src, p) }},
		{"savepng.png", func(p string) error { return SavePng(src, p) }},
	} {
		path := filepath.Join(dir, save.name)
		if err := save.fn(path); err != nil {
			t.Fatalf("%s: %v", save.name, err)
		}

		img, err := Read(path)
		if err != nil {
			t.Fatalf("Read %s: %v", save.name, err)
		}
		assertSamePixels(t, src, img)

		img, fm, err := DecodeImg(path)
		if err != nil || fm != "png" {
			t.Fatalf("DecodeImg %s = %q, %v", save.name, fm, err)
		}
		assertSamePixels(t, src, img)
	}
}

func TestSaveJpeg(t *testing.T) {
	src := testRGBA(8, 6)
	dir := t.TempDir()

	for i, save := range []func(string) error{
		func(p string) error { return Save(src, p, 80) },
		func(p string) error { return SaveJpeg(src, p, 50) },
	} {
		path := filepath.Join(dir, strconv.Itoa(i)+".jpeg")
		if err := save(path); err != nil {
			t.Fatalf("save %s: %v", path, err)
		}
		img, fm, err := DecodeImg(path)
		if err != nil || fm != "jpeg" {
			t.Fatalf("DecodeImg %s = %q, %v", path, fm, err)
		}
		if b := img.Bounds(); b.Dx() != 8 || b.Dy() != 6 {
			t.Fatalf("jpeg bounds = %v", b)
		}
	}
}

func TestOpenSaveImgBytes(t *testing.T) {
	dir := t.TempDir()
	src := testRGBA(3, 3)
	path := filepath.Join(dir, "src.png")
	if err := SavePng(src, path); err != nil {
		t.Fatal(err)
	}

	b, err := OpenImg(path)
	if err != nil || len(b) == 0 {
		t.Fatalf("OpenImg = %d bytes, %v", len(b), err)
	}
	img, err := ByteToImg(b)
	if err != nil {
		t.Fatalf("ByteToImg: %v", err)
	}
	assertSamePixels(t, src, img)

	out := filepath.Join(dir, "copy.png")
	if err := SaveImg(b, out); err != nil {
		t.Fatalf("SaveImg: %v", err)
	}
	got, err := os.ReadFile(out)
	if err != nil || !bytes.Equal(got, b) {
		t.Fatalf("SaveImg wrote %d bytes, want %d (%v)", len(got), len(b), err)
	}
}

// ToByteImg / ToStringImg return base64-encoded image data, which
// StrToImg decodes.
func TestToByteStrImg(t *testing.T) {
	src := testRGBA(2, 2)

	b := ToByteImg(src, "png")
	if len(b) == 0 {
		t.Fatal("ToByteImg returned no data")
	}
	s := ToStringImg(src, "png")
	if s != string(b) {
		t.Fatal("ToStringImg differs from ToByteImg")
	}

	img, err := StrToImg(s)
	if err != nil {
		t.Fatalf("StrToImg: %v", err)
	}
	assertSamePixels(t, src, img)

	// Default format is jpeg.
	dec := base64.NewDecoder(base64.StdEncoding, bytes.NewReader(ToByteImg(src)))
	if _, fm, err := image.Decode(dec); err != nil || fm != "jpeg" {
		t.Fatalf("default ToByteImg format = %q, %v", fm, err)
	}

	if ToByteImg(src, "webp") != nil {
		t.Error("ToByteImg unsupported format: want nil")
	}
}

func TestByteToImg(t *testing.T) {
	src := testRGBA(3, 2)
	var buf bytes.Buffer
	if err := png.Encode(&buf, src); err != nil {
		t.Fatal(err)
	}
	img, err := ByteToImg(buf.Bytes())
	if err != nil {
		t.Fatalf("ByteToImg: %v", err)
	}
	assertSamePixels(t, src, img)
}

func TestImgErrors(t *testing.T) {
	missing := filepath.Join(t.TempDir(), "missing.png")
	if _, err := Read(missing); err == nil {
		t.Error("Read missing file: want error")
	}
	if _, _, err := DecodeImg(missing); err == nil {
		t.Error("DecodeImg missing file: want error")
	}
	if _, err := OpenImg(missing); err == nil {
		t.Error("OpenImg missing file: want error")
	}
	if _, _, err := ImgSize(missing); err == nil {
		t.Error("ImgSize missing file: want error")
	}
	if _, err := ByteToImg([]byte("not an image")); err == nil {
		t.Error("ByteToImg garbage: want error")
	}
	if _, err := StrToImg("!!not base64!!"); err == nil {
		t.Error("StrToImg garbage: want error")
	}
	if err := Save(testRGBA(1, 1), filepath.Join(missing, "x.png")); err == nil {
		t.Error("Save into missing dir: want error")
	}
}

func assertSamePixels(t *testing.T, want, got image.Image) {
	t.Helper()
	if want.Bounds() != got.Bounds() {
		t.Fatalf("bounds = %v, want %v", got.Bounds(), want.Bounds())
	}
	b := want.Bounds()
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			wr, wg, wb, wa := want.At(x, y).RGBA()
			gr, gg, gb, ga := got.At(x, y).RGBA()
			if wr != gr || wg != gg || wb != gb || wa != ga {
				t.Fatalf("pixel (%d,%d) = %v, want %v", x, y, got.At(x, y), want.At(x, y))
			}
		}
	}
}
