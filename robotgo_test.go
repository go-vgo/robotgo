// Copyright 2016 The go-vgo Project Developers. See the COPYRIGHT
// file at the top-level directory of this distribution and at
// https://github.com/go-vgo/robotgo/blob/master/LICENSE
//
// Licensed under the Apache License, Version 2.0 <LICENSE-APACHE or
// http://www.apache.org/licenses/LICENSE-2.0> or the MIT license
// <LICENSE-MIT or http://opensource.org/licenses/MIT>, at your
// option. This file may not be copied, modified, or distributed
// except according to those terms.

// +build darwin windows

package robotgo

import (
	"fmt"
	"log"
	"runtime"
	"testing"

	"github.com/vcaesar/tt"
)

func TestGetVer(t *testing.T) {
	fmt.Println("go version: ", runtime.Version())
	ver := GetVersion()

	tt.Expect(t, Version, ver)
}

func TestMoveMouse(t *testing.T) {
	MoveMouse(20, 20)
	MilliSleep(10)
	x, y := GetMousePos()

	tt.Equal(t, 20, x)
	tt.Equal(t, 20, y)
}

func TestMoveMouseSmooth(t *testing.T) {
	MoveMouseSmooth(100, 100)
	MilliSleep(10)
	x, y := GetMousePos()

	tt.Equal(t, 100, x)
	tt.Equal(t, 100, y)
}

func TestDragMouse(t *testing.T) {
	DragMouse(500, 500)
	MilliSleep(10)
	x, y := GetMousePos()

	tt.Equal(t, 500, x)
	tt.Equal(t, 500, y)
}

func TestScrollMouse(t *testing.T) {
	ScrollMouse(120, "up")
	MilliSleep(100)

	Scroll(210, 210)
}

func TestGetScreenSize(t *testing.T) {
	x, y := GetScreenSize()
	log.Println("GetScreenSize: ", x, y)
}
