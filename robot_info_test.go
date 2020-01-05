// Copyright 2016 The go-vgo Project Developers. See the COPYRIGHT
// file at the top-level directory of this distribution and at
// https://github.com/go-vgo/robotgo/blob/master/LICENSE
//
// Licensed under the Apache License, Version 2.0 <LICENSE-APACHE or
// http://www.apache.org/licenses/LICENSE-2.0> or the MIT license
// <LICENSE-MIT or http://opensource.org/licenses/MIT>, at your
// option. This file may not be copied, modified, or distributed
// except according to those terms.

package robotgo_test

import (
	"fmt"
	"log"
	"runtime"
	"testing"

	"github.com/go-vgo/robotgo"
	"github.com/vcaesar/tt"
)

func TestGetVer(t *testing.T) {
	fmt.Println("go version: ", runtime.Version())
	ver := robotgo.GetVersion()

	tt.Expect(t, robotgo.Version, ver)
}

func TestGetScreenSize(t *testing.T) {
	x, y := robotgo.GetScreenSize()
	log.Println("GetScreenSize: ", x, y)
}

func TestGetSysScale(t *testing.T) {
	s := robotgo.SysScale()
	log.Println("SysScale: ", s)
}
