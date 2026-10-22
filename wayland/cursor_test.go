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

// Fractional scaling: wl_output.scale is rounded up (2), but the xdg_output
// logical size (1920x1080 for a 2880x1620 mode) gives the real 1.5 factor.
func TestLocationRealCursorFractional(t *testing.T) {
	installFakeConn(t, &outputInfo{width: 2880, height: 1620, scale: 2,
		logicalW: 1920, logicalH: 1080})
	cursorPos = func() (int, int, error) { return 101, 50, nil }
	t.Cleanup(func() { cursorPos = noCursor })
	if x, y := Location(); x != 152 || y != 75 {
		t.Fatalf("Location: got (%d,%d), want (152,75)", x, y)
	}
}

// Same fractional ratio on every output is still one exact factor.
func TestLocationRealCursorFractionalMulti(t *testing.T) {
	installFakeConn(t,
		&outputInfo{width: 2880, height: 1620, logicalW: 1920, logicalH: 1080},
		&outputInfo{x: 1920, width: 3840, height: 2160, logicalW: 2560, logicalH: 1440})
	cursorPos = func() (int, int, error) { return 2000, 100, nil }
	t.Cleanup(func() { cursorPos = noCursor })
	if x, y := Location(); x != 3000 || y != 150 {
		t.Fatalf("Location: got (%d,%d), want (3000,150)", x, y)
	}
}

// wlroots sends 0,0 in wl_output.geometry; the xdg_output logical position
// is the real origin and is scaled into pixel space like the cursor, so warp
// and Location agree on outputs left of / above the primary one.
func TestLocationRealCursorLogicalOrigin(t *testing.T) {
	fp := installFakeConn(t,
		&outputInfo{width: 3840, height: 2160, scale: 2,
			logicalW: 1920, logicalH: 1080, hasLogicalPos: true},
		&outputInfo{width: 3840, height: 2160, scale: 2,
			logicalW: 1920, logicalH: 1080,
			logicalX: -1920, logicalY: -540, hasLogicalPos: true})

	x, y, w, h := globalConn.outputBounds()
	if x != -3840 || y != -1080 || w != 7680 || h != 3240 {
		t.Fatalf("outputBounds: got (%d,%d,%d,%d), want (-3840,-1080,7680,3240)", x, y, w, h)
	}
	if r := GetScreenRect(1); r.X != -3840 || r.Y != -1080 || r.W != 3840 || r.H != 2160 {
		t.Fatalf("GetScreenRect(1): got %+v", r)
	}

	cursorPos = func() (int, int, error) { return -100, -10, nil }
	t.Cleanup(func() { cursorPos = noCursor })
	if x, y := Location(); x != -200 || y != -20 {
		t.Fatalf("Location: got (%d,%d), want (-200,-20)", x, y)
	}
	if !MoveSmooth(-190, -20, 1, 0) {
		t.Fatal("MoveSmooth returned false")
	}
	if len(fp.abs) != 1 || fp.abs[0] != [2]uint32{3650, 1060} {
		t.Fatalf("motion_absolute: got %v, want [[3650 1060]]", fp.abs)
	}
}

// A rotated output (logical size transposed from the mode) or differing
// fractional ratios have no exact factor: fall back to the tracked position.
func TestLocationRealCursorInexact(t *testing.T) {
	cases := map[string][]*outputInfo{
		"rotated": {{width: 1920, height: 1080, logicalW: 1080, logicalH: 1920}},
		"mixed": {
			{width: 2880, height: 1620, logicalW: 1920, logicalH: 1080},
			{x: 1920, width: 2560, height: 1440, logicalW: 2048, logicalH: 1152},
		},
	}
	for name, outs := range cases {
		t.Run(name, func(t *testing.T) {
			installFakeConn(t, outs...)
			if err := Move(10, 20); err != nil {
				t.Fatal(err)
			}
			cursorPos = func() (int, int, error) { return 100, 50, nil }
			t.Cleanup(func() { cursorPos = noCursor })
			if x, y := Location(); x != 10 || y != 20 {
				t.Fatalf("Location: got (%d,%d), want tracked (10,20)", x, y)
			}
		})
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
