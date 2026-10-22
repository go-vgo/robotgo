//go:build linux
// +build linux

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

package wayland

import (
	"os"
	"testing"

	"github.com/go-vgo/robotgo/internal/wlcursor"
)

func noCursor() (int, int, error) { return 0, 0, wlcursor.ErrUnavailable }

// TestMain pins Location to the tracked fallback so tests do not depend on
// the compositor (Hyprland/KDE) they happen to run under.
func TestMain(m *testing.M) {
	cursorPos = noCursor
	os.Exit(m.Run())
}

func TestLocationRealCursor(t *testing.T) {
	fp := installFakeConn(t, &outputInfo{x: 0, y: 0, width: 1920, height: 1080})
	cursorPos = func() (int, int, error) { return 500, 600, nil }
	t.Cleanup(func() { cursorPos = noCursor })

	if x, y := Location(); x != 500 || y != 600 {
		t.Fatalf("Location: got (%d,%d), want (500,600)", x, y)
	}

	// MoveSmooth starts from the real cursor: one step lands on the target.
	if !MoveSmooth(510, 600, 1, 0) {
		t.Fatal("MoveSmooth returned false")
	}
	if len(fp.abs) != 1 || fp.abs[0] != [2]uint32{510, 600} {
		t.Fatalf("motion_absolute: got %v", fp.abs)
	}

	// Without a source, Location falls back to the resynced tracked position.
	cursorPos = noCursor
	if err := MoveRelative(5, 5); err != nil {
		t.Fatal(err)
	}
	if x, y := Location(); x != 515 || y != 605 {
		t.Fatalf("tracked Location: got (%d,%d), want (515,605)", x, y)
	}
}

// The compositor reports logical coordinates; Move works in output pixels.
func TestLocationRealCursorScaled(t *testing.T) {
	installFakeConn(t, &outputInfo{width: 3840, height: 2160, scale: 2})
	cursorPos = func() (int, int, error) { return 100, 50, nil }
	t.Cleanup(func() { cursorPos = noCursor })
	if x, y := Location(); x != 200 || y != 100 {
		t.Fatalf("Location: got (%d,%d), want (200,100)", x, y)
	}
}

// Mixed scales have no single factor: fall back to the tracked position.
func TestLocationRealCursorMixedScale(t *testing.T) {
	installFakeConn(t,
		&outputInfo{width: 3840, height: 2160, scale: 2},
		&outputInfo{x: 3840, width: 1920, height: 1080, scale: 1})
	if err := Move(10, 20); err != nil {
		t.Fatal(err)
	}
	cursorPos = func() (int, int, error) { return 100, 50, nil }
	t.Cleanup(func() { cursorPos = noCursor })
	if x, y := Location(); x != 10 || y != 20 {
		t.Fatalf("Location: got (%d,%d), want tracked (10,20)", x, y)
	}
}
