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

	"github.com/go-vgo/robotgo/pub"
)

func TestIs64Bit(t *testing.T) {
	if Is64Bit() != (strconv.IntSize == 64) {
		t.Errorf("Is64Bit: got %v, want %v", Is64Bit(), strconv.IntSize == 64)
	}
}

// Invalid code points must error before any key event is posted.
func TestUnicodeTypeInvalid(t *testing.T) {
	for _, r := range []uint32{0xD800, 0x110000} {
		if err := UnicodeType(r); err == nil {
			t.Errorf("UnicodeType(%#x): want error", r)
		}
	}
}

// Handle mode follows Cgo: any extra arg or NotPid.
func TestIsHandle(t *testing.T) {
	old := NotPid
	defer func() { NotPid = old }()
	NotPid = false
	if isHandle(nil) || !isHandle([]int{1}) {
		t.Error("isHandle: wrong result with NotPid=false")
	}
	NotPid = true
	if !isHandle(nil) {
		t.Error("isHandle: want true with NotPid=true")
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

func TestRgbToHex(t *testing.T) {
	if got := RgbToHex(0x12, 0xab, 0xff); got != 0x12abff {
		t.Errorf("RgbToHex: got %#x, want 0x12abff", got)
	}
	if got := PadHex(RgbToHex(0, 0, 1)); got != "000001" {
		t.Errorf("PadHex(RgbToHex(0, 0, 1)): got %q", got)
	}
}

// TestDelaySync checks init copied the robotgo_pub.go config into pub;
// SetDelay keeping both in sync is covered by TestSetDelay.
func TestDelaySync(t *testing.T) {
	if pub.DisplayID != DisplayID {
		t.Fatalf("pub.DisplayID = %d, want %d", pub.DisplayID, DisplayID)
	}
}

// KeyPress must dispatch to the backend's atomic KeyTap (via keyPress),
// not the shared down/sleep/up sequence: an unknown key fails exactly
// like KeyTap, before any key-down is posted.
func TestKeyPressUsesBackendTap(t *testing.T) {
	const bad = "no_such_key"
	errTap, errPress := KeyTap(bad), KeyPress(bad)
	if errTap == nil || errPress == nil || errTap.Error() != errPress.Error() {
		t.Fatalf("KeyPress(%q) = %v, want KeyTap error %v", bad, errPress, errTap)
	}
}
