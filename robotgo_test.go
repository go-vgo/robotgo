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
	MoveMouse(10, 10)
	Sleep(1)
	x, y := GetMousePos()

	tt.Equal(t, 10, x)
	tt.Equal(t, 10, y)
}

func TestMoveMouseSmooth(t *testing.T) {
	MoveMouseSmooth(100, 100)
	Sleep(1)
	x, y := GetMousePos()

	tt.Equal(t, 100, x)
	tt.Equal(t, 100, y)
}
