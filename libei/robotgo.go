//go:build linux
// +build linux

// Copyright (c) 2016-2026 AtomAI, All rights reserved.
//
// See the COPYRIGHT file at the top-level directory of this distribution and at
// https://github.com/go-vgo/robotgo/blob/master/LICENSE
//
// Licensed under the Apache License, Version 2.0 <LICENSE-APACHE or
// http://www.apache.org/licenses/LICENSE-2.0>
//
// This file may not be copied, modified, or distributed
// except according to those terms.

package libei

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Version is the robotgo-libei backend version string.
const Version = "v0.1.0-libei"

// Point represents a 2D point.
type Point struct {
	X, Y int
}

// Size represents a 2D size.
type Size struct {
	W, H int
}

// Rect represents a rectangle.
type Rect struct {
	Point
	Size
}

// Nps represents a process.
type Nps struct {
	Pid  int
	Name string
}

// GetVersion returns the robotgo version.
func GetVersion() string { return Version }

// Sleep sleeps for tm seconds.
func Sleep(tm int) { time.Sleep(time.Duration(tm) * time.Second) }

// MilliSleep sleeps for tm milliseconds.
func MilliSleep(tm int) { time.Sleep(time.Duration(tm) * time.Millisecond) }

// --- Image helpers (pure Go, no portal dependency) ---

// Save saves an image to a file. Format is determined by the extension.
func Save(img image.Image, path string, quality ...int) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	switch strings.ToLower(filepath.Ext(path)) {
	case ".jpg", ".jpeg":
		q := 90
		if len(quality) > 0 {
			q = quality[0]
		}
		return jpeg.Encode(f, img, &jpeg.Options{Quality: q})
	default:
		return png.Encode(f, img)
	}
}

// SavePng saves an image as PNG.
func SavePng(img image.Image, path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return png.Encode(f, img)
}

// SaveJpeg saves an image as JPEG.
func SaveJpeg(img image.Image, path string, quality ...int) error {
	q := 90
	if len(quality) > 0 {
		q = quality[0]
	}
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return jpeg.Encode(f, img, &jpeg.Options{Quality: q})
}

// Width returns the width of an image.
func Width(img image.Image) int { return img.Bounds().Dx() }

// Height returns the height of an image.
func Height(img image.Image) int { return img.Bounds().Dy() }

// PadHex pads a hex color value to 6 characters.
func PadHex(hex uint32) string { return fmt.Sprintf("%06x", hex) }

// SaveCapture captures the screen and saves it to a file.
// Screen capture is not implemented by this backend (see screen.go).
func SaveCapture(path string, args ...int) error {
	img, err := CaptureImg(args...)
	if err != nil {
		return err
	}
	return Save(img, path)
}
