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

package libei

import (
	"os"
	"testing"

	"github.com/go-vgo/robotgo/internal/wlcursor"
)

func noCursor() (int, int, error) { return 0, 0, wlcursor.ErrUnavailable }

// TestMain pins Location to the tracked fallback so tests do not depend on
// the compositor (KDE/Hyprland) they happen to run under.
func TestMain(m *testing.M) {
	cursorPos = noCursor
	os.Exit(m.Run())
}

func stubCursor(t *testing.T, x, y int) {
	t.Helper()
	cursorPos = func() (int, int, error) { return x, y, nil }
	t.Cleanup(func() { cursorPos = noCursor })
}

func TestLocationRealCursor(t *testing.T) {
	inj := installFakeConn(t)
	stubCursor(t, 300, 400)
	if x, y := Location(); x != 300 || y != 400 {
		t.Fatalf("Location: got (%d,%d), want (300,400)", x, y)
	}

	// The real position resyncs tracking: a relative-only Move sends the
	// delta from it, without the corner reset.
	cursorPos = noCursor
	if x, y := Location(); x != 300 || y != 400 {
		t.Fatalf("tracked Location: got (%d,%d), want (300,400)", x, y)
	}
	if err := Move(310, 420); err != nil {
		t.Fatal(err)
	}
	if len(inj.rel) != 1 || inj.rel[0] != (relCall{10, 20}) {
		t.Fatalf("relative motion: got %v, want [{10 20}]", inj.rel)
	}
}

func TestLocationRealCursorNoSession(t *testing.T) {
	connMu.Lock()
	prev := globalConn
	globalConn = nil
	connMu.Unlock()
	t.Cleanup(func() {
		connMu.Lock()
		globalConn = prev
		connMu.Unlock()
	})

	stubCursor(t, 12, 34)
	if x, y := Location(); x != 12 || y != 34 {
		t.Fatalf("Location: got (%d,%d), want (12,34)", x, y)
	}
	connMu.Lock()
	defer connMu.Unlock()
	if globalConn != nil {
		t.Fatal("Location opened a portal session")
	}
}

// Without a stream, Move resyncs from the real cursor first, so the delta
// starts where the pointer really is (no stale delta, no corner reset).
func TestMoveSyncsRealCursor(t *testing.T) {
	inj := installFakeConn(t)
	stubCursor(t, 100, 100)
	if err := Move(130, 140); err != nil {
		t.Fatal(err)
	}
	if len(inj.rel) != 1 || inj.rel[0] != (relCall{30, 40}) {
		t.Fatalf("relative motion: got %v, want [{30 40}]", inj.rel)
	}
}
