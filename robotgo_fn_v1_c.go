//go:build !wayland && !win && !libei && !mac && !x11 && !purego
// +build !wayland,!win,!libei,!mac,!x11,!purego

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

// Cgo-only v1 aliases (they need ScaleX); the portable v1 aliases are in
// robotgo_fn_v1.go.

// Deprecated: use the ScaledF(),
//
// Scale1 get the screen scale (only windows old), drop
func Scale1() int {
	dpi := map[int]int{
		0: 100,
		// DPI Scaling Level
		96:  100,
		120: 125,
		144: 150,
		168: 175,
		192: 200,
		216: 225,
		// Custom DPI
		240: 250,
		288: 300,
		384: 400,
		480: 500,
	}

	x := ScaleX()
	return dpi[x]
}

// Deprecated: use the ScaledF(),
//
// Scale0 return ScaleX() / 0.96, drop
func Scale0() int {
	return int(float64(ScaleX()) / 0.96)
}

// Deprecated: use the ScaledF(),
//
// Mul mul the scale, drop
func Mul(x int) int {
	s := Scale1()
	return x * s / 100
}
