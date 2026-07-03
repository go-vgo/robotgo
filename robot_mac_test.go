//go:build darwin && mac
// +build darwin,mac

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

package robotgo_test

import (
	"testing"

	"github.com/go-vgo/robotgo"
)

// TestGetVerMac checks the version string on the pure-Go mac backend.
func TestGetVerMac(t *testing.T) {
	if robotgo.GetVersion() != robotgo.Version {
		t.Fatalf("GetVersion: got %q, want %q", robotgo.GetVersion(), robotgo.Version)
	}
}

// TestColorMac exercises the color path (GetPixelColor/PadHex) on the
// pure-Go mac backend; robotgo_test.go's TestColor is excluded by -tags mac.
func TestColorMac(t *testing.T) {
	// Headless / unpermissioned environments return the "000000" fallback;
	// just ensure a valid 6-char hex string comes back without panicking.
	c := robotgo.GetPixelColor(1, 1)
	if len(c) != 6 {
		t.Fatalf("GetPixelColor: got %q, want 6 hex chars", c)
	}

	if got := robotgo.PadHex(0xABCDEF); got != "abcdef" {
		t.Fatalf("PadHex(0xABCDEF): got %q, want %q", got, "abcdef")
	}
	if got := robotgo.PadHex(0x123); got != "000123" {
		t.Fatalf("PadHex(0x123): got %q, want %q", got, "000123")
	}
}

// TestGetMainIdMac checks the main display id wrapper on the mac backend.
func TestGetMainIdMac(t *testing.T) {
	id := robotgo.GetMainId()
	if id < 0 {
		t.Fatalf("GetMainId: got negative id %d", id)
	}
	if !robotgo.IsMain(id) {
		t.Fatalf("IsMain(%d): expected true", id)
	}
}
