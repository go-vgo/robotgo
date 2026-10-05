//go:build wayland || win || libei || mac || x11 || purego
// +build wayland win libei mac x11 purego

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

package robotgo

import (
	"strconv"
	"testing"
)

func TestIs64Bit(t *testing.T) {
	if Is64Bit() != (strconv.IntSize == 64) {
		t.Errorf("Is64Bit: got %v, want %v", Is64Bit(), strconv.IntSize == 64)
	}
}

// A non-positive count must not post any click.
func TestMultiClickNoop(t *testing.T) {
	for _, n := range []int{0, -1} {
		if err := MultiClick("left", n); err != nil {
			t.Errorf("MultiClick(%d): got %v, want nil", n, err)
		}
		if err := clickTimes("left", n); err != nil {
			t.Errorf("clickTimes(%d): got %v, want nil", n, err)
		}
	}
}
