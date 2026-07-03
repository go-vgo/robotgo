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

//go:build darwin && !mac
// +build darwin,!mac

package robotgo

/*
#include <CoreGraphics/CoreGraphics.h>
*/
import "C"

// GetMainId get the main display id
func GetMainId() int {
	return int(C.CGMainDisplayID())
}
