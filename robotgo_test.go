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

// Untagged on purpose: these interactive tests only use the public API that
// every backend wires (Cgo and -tags mac/win/x11/wayland/libei/purego), so
// the same file checks Move -> Location, keys, clipboard and capture on all
// of them. Cgo-only APIs (CaptureScreen, GetPxColor, ...) are covered in
// robot_info_test.go instead.

package robotgo

import (
	"os"
	"runtime"
	"strings"
	"sync"
	"testing"

	"github.com/vcaesar/tt"
)

var (
	displayOnce sync.Once
	displayOK   bool
)

// requireDisplay skips the test when no session can receive injected input:
// Linux without DISPLAY/WAYLAND_DISPLAY, a pure-Go backend with no
// connection, or a runner without Accessibility rights (Move has no effect).
// The probe is retried once: the first event a process posts to the HID tap
// can land a px or two off while the WindowServer connection settles.
func requireDisplay(t *testing.T) {
	t.Helper()
	displayOnce.Do(func() {
		if runtime.GOOS == "linux" && os.Getenv("DISPLAY") == "" && os.Getenv("WAYLAND_DISPLAY") == "" {
			return
		}
		for i := 0; i < 2 && !displayOK; i++ {
			Move(10, 10)
			MilliSleep(50)
			x, y := Location()
			displayOK = x == 10 && y == 10
		}
	})
	if !displayOK {
		t.Skip("no display session for input injection")
	}
}

// requireScreen skips only when there is no screen at all: Linux without
// DISPLAY/WAYLAND_DISPLAY, or a backend that reports a zero screen size
// (libei without a portal ScreenCast stream, wayland/darwin without a
// compositor or WindowServer connection). Screen-read APIs work without
// Accessibility or an input session, so these tests also run on
// macOS/Windows CI runners.
func requireScreen(t *testing.T) {
	t.Helper()
	if runtime.GOOS == "linux" && os.Getenv("DISPLAY") == "" && os.Getenv("WAYLAND_DISPLAY") == "" {
		t.Skip("no display for screen reads")
	}
	if w, h := GetScreenSize(); w == 0 || h == 0 {
		t.Skipf("screen size unavailable: %dx%d", w, h)
	}
}

func TestColor(t *testing.T) {
	requireScreen(t)

	s := GetPixelColor(10, 10)
	tt.IsType(t, "string", s)
	tt.Equal(t, 6, len(s))

	tt.Equal(t, "abcdef", PadHex(0xABCDEF))
	tt.Equal(t, "000123", PadHex(0x123))
}

func TestSize(t *testing.T) {
	requireScreen(t)

	x, y := GetScreenSize()
	tt.NotZero(t, x)
	tt.NotZero(t, y)

	x, y = GetScaleSize()
	tt.NotZero(t, x)
	tt.NotZero(t, y)
}

// Move then Location must agree on every backend (#783): the Wayland/libei
// ports cannot read the real cursor, so Location is the last injected point
// and polling it must not drift or re-open the session.
func TestMoveMouse(t *testing.T) {
	requireDisplay(t)

	Move(20, 20)
	MilliSleep(50)
	for i := 0; i < 3; i++ {
		x, y := Location()
		tt.Equal(t, 20, x)
		tt.Equal(t, 20, y)
		MilliSleep(20)
	}
}

func TestMoveMouseSmooth(t *testing.T) {
	requireDisplay(t)

	b := MoveSmooth(100, 100)
	MilliSleep(50)
	x, y := Location()

	tt.True(t, b)
	tt.Equal(t, 100, x)
	tt.Equal(t, 100, y)
}

func TestDragMouse(t *testing.T) {
	requireDisplay(t)

	DragSmooth(500, 500)
	MilliSleep(50)
	x, y := Location()

	tt.Equal(t, 500, x)
	tt.Equal(t, 500, y)
}

func TestScrollMouse(t *testing.T) {
	requireDisplay(t)

	ScrollDir(120, "up")
	ScrollDir(100, "right")

	Scroll(0, 120)
	MilliSleep(100)

	Scroll(210, 210)
	MilliSleep(10)
}

func TestMoveRelative(t *testing.T) {
	requireDisplay(t)

	Move(200, 200)
	MilliSleep(50)

	MoveRelative(10, -10)
	MilliSleep(50)

	x, y := Location()
	tt.Equal(t, 210, x)
	tt.Equal(t, 190, y)
}

func TestMoveSmoothRelative(t *testing.T) {
	requireDisplay(t)

	Move(200, 200)
	MilliSleep(50)

	MoveSmoothRelative(10, -10)
	MilliSleep(50)

	x, y := Location()
	tt.Equal(t, 210, x)
	tt.Equal(t, 190, y)
}

func TestMouseToggle(t *testing.T) {
	requireDisplay(t)

	e := Toggle("right")
	tt.Nil(t, e)

	e = Toggle("right", "up")
	tt.Nil(t, e)

	e = MouseDown("left")
	tt.Nil(t, e)

	e = MouseUp("left")
	tt.Nil(t, e)
}

// skipNoClipboard skips when the platform has no clipboard tool (e.g. a Linux
// CI image without xclip/xsel/wl-clipboard).
func skipNoClipboard(t *testing.T, err error) {
	t.Helper()
	if err != nil && strings.Contains(err.Error(), "no clipboard utilities") {
		t.Skipf("clipboard unavailable: %v", err)
	}
}

func TestClip(t *testing.T) {
	requireDisplay(t)

	err := WriteAll("s")
	skipNoClipboard(t, err)
	tt.Nil(t, err)

	s, e := ReadAll()
	tt.Equal(t, "s", s)
	tt.Nil(t, e)
}

func TestImage(t *testing.T) {
	requireScreen(t)

	img1, err := CaptureImg(10, 10, 20, 20)
	if err != nil {
		// e.g. macOS without the Screen Recording permission, Xvfb without
		// a usable framebuffer, Wayland without screencopy.
		t.Skipf("screen capture unavailable: %v", err)
	}
	if err := os.MkdirAll("test/tmp", 0o755); err != nil {
		t.Fatal(err)
	}
	e := Save(img1, "test/tmp/robot_img.jpeg", 50)
	tt.Nil(t, e)
	e = SavePng(img1, "test/tmp/robot_test.png")
	tt.Nil(t, e)

	// CGDisplayCreateImageForRect (Cgo and pure-Go darwin alike) returns
	// backing pixels, so the captured size is the request times the scale.
	f := ScaleF()
	if runtime.GOOS != "darwin" {
		f = 1
	}
	tt.Equal(t, int(20*f), Width(img1))
	tt.Equal(t, int(20*f), Height(img1))

	bit1 := ImgToBitmap(img1)
	tt.Equal(t, bit1.Width, Width(img1))
	tt.Equal(t, bit1.Height, Height(img1))
}

func TestPs(t *testing.T) {
	id, err := Pids()
	tt.Not(t, "[]", id)
	tt.IsType(t, "[]int", id)
	tt.Nil(t, err)

	ps, e := Process()
	tt.Not(t, "[]", ps)
	tt.IsType(t, "[]robotgo.Nps", ps)
	tt.Nil(t, e)

	// Pids()[0] may be the kernel (pid 0), which gopsutil rejects; use our own.
	self := os.Getpid()
	b, e := PidExists(self)
	tt.True(t, b)
	tt.Nil(t, e)

	n, e := FindName(self)
	tt.NotEmpty(t, n)
	tt.Nil(t, e)

	n1, e := FindNames()
	tt.Not(t, "[]", n1)
	tt.IsType(t, "[]string", n1)
	tt.Nil(t, e)

	id, err = FindIds(n1[0])
	tt.Not(t, "[]", id)
	tt.IsType(t, "[]int", id)
	tt.Nil(t, err)

	// n, e = FindPath(id[0])
	// tt.NotEmpty(t, n)
	// tt.Nil(t, e)
}

// func TestAlert(t *testing.T) {
// 	go func() {
// 		MilliSleep(200)
// 		KeyTap("enter")
// 		log.Println("tap...")
// 	}()

// 	i := Alert("t", "msg")
//	tt.True(t, i)
// }
