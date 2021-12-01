// Copyright 2016 The go-vgo Project Developers. See the COPYRIGHT
// file at the top-level directory of this distribution and at
// https://github.com/go-vgo/robotgo/blob/master/LICENSE
//
// Licensed under the Apache License, Version 2.0 <LICENSE-APACHE or
// http://www.apache.org/licenses/LICENSE-2.0> or the MIT license
// <LICENSE-MIT or http://opensource.org/licenses/MIT>, at your
// option. This file may not be copied, modified, or distributed
// except according to those terms.

//go:build darwin || windows
// +build darwin windows

package robotgo

import (
	"fmt"
	"testing"

	"github.com/vcaesar/tt"
)

func TestColor(t *testing.T) {
	s := GetPixelColor(10, 10)
	tt.IsType(t, "string", s)
	tt.NotEmpty(t, s)

	c := GetPxColor(10, 10)
	s1 := PadHex(c)
	tt.Equal(t, s, s1)
}

func TestSize(t *testing.T) {
	x, y := GetScreenSize()
	tt.NotZero(t, x)
	tt.NotZero(t, y)

	x, y = GetScaleSize()
	tt.NotZero(t, x)
	tt.NotZero(t, y)
}

func TestMoveMouse(t *testing.T) {
	Move(20, 20)
	MilliSleep(50)
	x, y := GetMousePos()

	tt.Equal(t, 20, x)
	tt.Equal(t, 20, y)
}

func TestMoveMouseSmooth(t *testing.T) {
	b := MoveSmooth(100, 100)
	MilliSleep(50)
	x, y := GetMousePos()

	tt.True(t, b)
	tt.Equal(t, 100, x)
	tt.Equal(t, 100, y)
}

func TestDragMouse(t *testing.T) {
	DragMouse(500, 500)
	MilliSleep(50)
	x, y := GetMousePos()

	tt.Equal(t, 500, x)
	tt.Equal(t, 500, y)
}

func TestScrollMouse(t *testing.T) {
	// ScrollMouse(120, "up")
	Scroll(0, 120)
	MilliSleep(100)

	Scroll(210, 210)
	MilliSleep(10)
}

func TestMoveRelative(t *testing.T) {
	Move(200, 200)
	MilliSleep(50)

	MoveRelative(10, -10)
	MilliSleep(50)

	x, y := GetMousePos()
	tt.Equal(t, 210, x)
	tt.Equal(t, 190, y)
}

func TestMoveSmoothRelative(t *testing.T) {
	Move(200, 200)
	MilliSleep(50)

	MoveSmoothRelative(10, -10)
	MilliSleep(50)

	x, y := GetMousePos()
	tt.Equal(t, 210, x)
	tt.Equal(t, 190, y)
}

func TestMouseToggle(t *testing.T) {
	e := Toggle("right")
	tt.Zero(t, e)

	e = Toggle("right", "up")
	tt.Zero(t, e)
}

func TestKey(t *testing.T) {
	e := KeyTap("v", "cmd")
	tt.Empty(t, e)

	e = KeyTap("enter")
	tt.Empty(t, e)

	e = KeyToggle("v", "up")
	tt.Empty(t, e)
}

func TestClip(t *testing.T) {
	err := WriteAll("s")
	tt.Nil(t, err)

	s, e := ReadAll()
	tt.Equal(t, "s", s)
	tt.Nil(t, e)
}

func TestTypeStr(t *testing.T) {
	c := CharCodeAt("s", 0)
	tt.Equal(t, 115, c)

	e := PasteStr("s")
	tt.Empty(t, e)

	s1 := "abc\\\\cd/s@世界"
	uc := ToUC(s1)
	tt.Equal(t, "[a b c \\ \\ c d / s @ U4e16 U754c]", uc)
}

func TestKeyCode(t *testing.T) {
	m := MouseMap["left"]
	tt.Equal(t, 1, m)

	k := Keycode["1"]
	tt.Equal(t, 2, k)

	s := Special["+"]
	tt.Equal(t, "=", s)
}

func TestBitmap(t *testing.T) {
	bit := CaptureScreen()
	tt.NotNil(t, bit)
	e := SaveBitmap(bit, "robot_test.png")
	tt.Empty(t, e)

	bit0 := CaptureScreen(10, 10, 20, 20)
	x, y := FindBitmap(bit0)
	fmt.Println("Find bitmap: ", x, y)

	arr := FindAllBitmap(bit0)
	fmt.Println("Find all bitmap:", arr)
	tt.Equal(t, 1, len(arr))

	c1 := CHex(0xAADCDC)
	x, y = FindColor(c1)
	fmt.Println("Find color: ", x, y)
	arr = FindAllColor(c1)
	fmt.Println("Find all color: ", arr)

	img := ToImage(bit)
	err := SavePng(img, "robot_img.png")
	tt.Nil(t, err)

	bit1 := OpenBitmap("robot_test.png")
	b := tt.TypeOf(bit, bit1)
	tt.True(t, b)
	tt.NotNil(t, bit1)
}

func TestPs(t *testing.T) {
	id, err := Pids()
	tt.Not(t, "[]", id)
	tt.IsType(t, "[]int32", id)
	tt.Nil(t, err)

	ps, e := Process()
	tt.Not(t, "[]", ps)
	tt.IsType(t, "[]robotgo.Nps", ps)
	tt.Nil(t, e)

	b, e := PidExists(id[0])
	tt.Bool(t, b)
	tt.Nil(t, e)

	n, e := FindName(id[0])
	tt.NotEmpty(t, n)
	tt.Nil(t, e)

	n1, e := FindNames()
	tt.Not(t, "[]", n1)
	tt.IsType(t, "[]string", n1)
	tt.Nil(t, e)

	id, err = FindIds(n1[0])
	tt.Not(t, "[]", id)
	tt.IsType(t, "[]int32", id)
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

// 	i := ShowAlert("t", "msg")
// 	tt.Zero(t, i)
// }
