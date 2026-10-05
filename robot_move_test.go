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

// Untagged: compiled for the Cgo backend and for every pure-Go backend
// (-tags mac/win/x11/wayland/libei/purego), so Move -> Location is checked
// through the public robotgo API whichever backend is wired in (#783).

package robotgo_test

import (
	"os"
	"runtime"
	"testing"

	"github.com/go-vgo/robotgo"
)

// skipHeadless skips when no display session can receive pointer input: on
// Linux without X11/Wayland (the pure-Go backends then have no connection
// and Move is a no-op), and when robotgo.Move did not take effect at all
// (e.g. macOS CI runner without Accessibility, or a Wayland portal that
// refused the session).
func skipHeadless(t *testing.T) {
	t.Helper()
	if runtime.GOOS == "linux" && os.Getenv("DISPLAY") == "" && os.Getenv("WAYLAND_DISPLAY") == "" {
		t.Skip("no DISPLAY / WAYLAND_DISPLAY")
	}
	robotgo.Move(10, 10)
	robotgo.MilliSleep(50)
	if x, y := robotgo.Location(); x != 10 || y != 10 {
		t.Skipf("pointer injection unavailable in this session: Move(10,10) -> Location() = (%d,%d)", x, y)
	}
}

func TestMoveThenLocation(t *testing.T) {
	skipHeadless(t)

	robotgo.Move(20, 20)
	robotgo.MilliSleep(50)
	if x, y := robotgo.Location(); x != 20 || y != 20 {
		t.Fatalf("Location after Move(20,20): got (%d,%d)", x, y)
	}

	// Polling Location must keep returning the injected position (#783:
	// libei/wayland used to drift back to (0,0) or re-open the session).
	for i := 0; i < 3; i++ {
		robotgo.MilliSleep(20)
		if x, y := robotgo.Location(); x != 20 || y != 20 {
			t.Fatalf("Location poll %d: got (%d,%d), want (20,20)", i, x, y)
		}
	}

	robotgo.Move(200, 200)
	robotgo.MilliSleep(50)
	robotgo.MoveRelative(10, -10)
	robotgo.MilliSleep(50)
	if x, y := robotgo.Location(); x != 210 || y != 190 {
		t.Fatalf("Location after MoveRelative(10,-10): got (%d,%d), want (210,190)", x, y)
	}

	if !robotgo.MoveSmooth(100, 100) {
		t.Fatal("MoveSmooth(100,100) returned false")
	}
	robotgo.MilliSleep(50)
	if x, y := robotgo.Location(); x != 100 || y != 100 {
		t.Fatalf("Location after MoveSmooth(100,100): got (%d,%d)", x, y)
	}
}
