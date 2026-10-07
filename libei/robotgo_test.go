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
	"errors"
	"slices"
	"testing"

	"github.com/godbus/dbus/v5"
)

// --- Pure Go tests (run anywhere, no portal/D-Bus needed) ---

func TestKeyToEvdev(t *testing.T) {
	tests := []struct {
		key  string
		code int32
		ok   bool
	}{
		{"a", 30, true},
		{"z", 44, true},
		{"enter", 28, true},
		{"escape", 1, true},
		{"esc", 1, true},
		{"caps", 58, true},
		{"capslock", 58, true},
		{"f1", 59, true},
		{"f12", 88, true},
		{"shift", 42, true},
		{"shiftl", 42, true},
		{"shiftr", 54, true},
		{"ctrl", 29, true},
		{"alt", 56, true},
		{"space", 57, true},
		{"tab", 15, true},
		{"backspace", 14, true},
		{"delete", 111, true},
		{"up", 103, true},
		{"down", 108, true},
		{"left", 105, true},
		{"right", 106, true},
		{"home", 102, true},
		{"end", 107, true},
		{"num+", 78, true},
		{"num-", 74, true},
		{"num*", 55, true},
		{"num/", 98, true},
		{"num_enter", 96, true},
		{"num.", 83, true},
		{"num0", 82, true},
		{"num_equal", 117, true},
		{"numpad_5", 76, true},
		{"numpad_lock", 69, true},
		{"pause_break", 119, true},
		{"right_shift", 54, true},
		{"audio_rewind", 168, true},
		{"audio_forward", 208, true},
		{"audio_repeat", 439, true},
		{"audio_random", 410, true},
		{"lights_mon_up", 225, true},
		{"lights_mon_down", 224, true},
		{"lights_kbd_toggle", 228, true},
		{"lights_kbd_up", 230, true},
		{"lights_kbd_down", 229, true},
		{"nonexistent_key", 0, false},
		{"", 0, false},
	}

	for _, tt := range tests {
		code, ok := keyToEvdev(tt.key)
		if ok != tt.ok {
			t.Errorf("keyToEvdev(%q): got ok=%v, want ok=%v", tt.key, ok, tt.ok)
		}
		if ok && code != tt.code {
			t.Errorf("keyToEvdev(%q): got code=%d, want code=%d", tt.key, code, tt.code)
		}
	}
}

func TestResolveButton(t *testing.T) {
	tests := []struct {
		btn  string
		want int32
	}{
		{"left", btnLeft},
		{"right", btnRight},
		{"center", btnMiddle},
		{"middle", btnMiddle},
		{"", btnLeft},
		{"unknown", btnLeft},
	}
	for _, tt := range tests {
		if got := resolveButton(tt.btn); got != tt.want {
			t.Errorf("resolveButton(%q): got %d, want %d", tt.btn, got, tt.want)
		}
	}
}

func TestExtractModifiers(t *testing.T) {
	tests := []struct {
		name string
		args []interface{}
		want []string
	}{
		{"no mods", []interface{}{}, nil},
		{"ctrl", []interface{}{"ctrl"}, []string{"ctrl"}},
		{"ctrl+shift", []interface{}{"ctrl", "shift"}, []string{"ctrl", "shift"}},
		{"mixed types", []interface{}{"ctrl", 42, true, "alt"}, []string{"ctrl", "alt"}},
		{"non-modifier string", []interface{}{"hello"}, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := extractModifiers(tt.args)
			if len(got) != len(tt.want) {
				t.Errorf("extractModifiers(%v): got %v, want %v", tt.args, got, tt.want)
				return
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Errorf("extractModifiers(%v)[%d]: got %q, want %q", tt.args, i, got[i], tt.want[i])
				}
			}
		})
	}
}

// TestReleaseKeys verifies the upKeyArr-equivalent helper: modifiers are
// keyed up in reverse order, every code gets a release attempt even after a
// failure, and errors are propagated.
func TestReleaseKeys(t *testing.T) {
	var got []int32
	err := releaseKeys([]int32{29, 42, 56}, func(code int32) error {
		got = append(got, code)
		return nil
	})
	if err != nil {
		t.Fatalf("releaseKeys: unexpected error %v", err)
	}
	want := []int32{56, 42, 29}
	if len(got) != len(want) {
		t.Fatalf("releaseKeys: got %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("releaseKeys[%d]: got %d, want %d (reverse order)", i, got[i], want[i])
		}
	}

	// A failing release must not stop the remaining releases.
	got = nil
	err = releaseKeys([]int32{29, 42, 56}, func(code int32) error {
		got = append(got, code)
		if code == 42 {
			return ErrNotSupported
		}
		return nil
	})
	if err == nil {
		t.Error("releaseKeys: expected error to propagate")
	}
	if len(got) != 3 {
		t.Errorf("releaseKeys: released %v, want all 3 despite failure", got)
	}

	if err := releaseKeys(nil, func(int32) error { return ErrNotSupported }); err != nil {
		t.Errorf("releaseKeys(nil): got %v, want nil", err)
	}
}

func TestRuneToKeysym(t *testing.T) {
	tests := []struct {
		r    rune
		want int32
	}{
		{'a', 0x61},
		{'A', 0x41},
		{'0', 0x30},
		{' ', 0x20},
		{'~', 0x7e},
		{'é', 0xe9},                // Latin-1: keysym == codepoint
		{'€', 0x20ac | 0x01000000}, // beyond Latin-1: Unicode keysym range
		{'😀', 0x1f600 | 0x01000000},
		{'\n', 0xff0d}, // control chars map to function keysyms
		{'\r', 0xff0d},
		{'\t', 0xff09},
		{'\b', 0xff08},
		{0x1b, 0xff1b},
	}
	for _, tt := range tests {
		if got := runeToKeysym(tt.r); got != tt.want {
			t.Errorf("runeToKeysym(%q): got 0x%x, want 0x%x", tt.r, got, tt.want)
		}
	}
}

func TestEscapedSender(t *testing.T) {
	c := &conn{uniqueName: ":1.42"}
	if got := c.escapedSender(); got != "1_42" {
		t.Errorf("escapedSender(:1.42): got %q, want %q", got, "1_42")
	}
	c2 := &conn{uniqueName: ":1.2345"}
	if got := c2.escapedSender(); got != "1_2345" {
		t.Errorf("escapedSender(:1.2345): got %q, want %q", got, "1_2345")
	}
}

func TestNextTokenUnique(t *testing.T) {
	a := nextToken("req")
	b := nextToken("req")
	if a == b {
		t.Errorf("nextToken returned duplicate tokens: %q", a)
	}
}

func TestDeviceCaps(t *testing.T) {
	c := &conn{devices: deviceKeyboard | devicePointer}
	if !c.hasKeyboard() || !c.hasPointer() {
		t.Errorf("expected keyboard+pointer granted, got devices=%d", c.devices)
	}
	c2 := &conn{devices: devicePointer}
	if c2.hasKeyboard() {
		t.Error("hasKeyboard() true when only pointer granted")
	}
	if !c2.hasPointer() {
		t.Error("hasPointer() false when pointer granted")
	}
}

func TestParseStreams(t *testing.T) {
	// streams: a(ua{sv}) — one node with position (10,20) and size (640,480).
	raw := [][]interface{}{
		{
			uint32(7),
			map[string]dbus.Variant{
				"position": dbus.MakeVariant([]interface{}{int32(10), int32(20)}),
				"size":     dbus.MakeVariant([]interface{}{int32(640), int32(480)}),
			},
		},
	}
	got := parseStreams(dbus.MakeVariant(raw))
	if len(got) != 1 {
		t.Fatalf("parseStreams: got %d streams, want 1", len(got))
	}
	s := got[0]
	if s.nodeID != 7 || s.x != 10 || s.y != 20 || s.width != 640 || s.height != 480 {
		t.Errorf("parseStreams: got %+v", s)
	}

	// Non-stream variant should yield nil without panicking.
	if got := parseStreams(dbus.MakeVariant("not-a-stream")); got != nil {
		t.Errorf("parseStreams(garbage): got %+v, want nil", got)
	}
}

func TestPadHex(t *testing.T) {
	tests := []struct {
		hex  uint32
		want string
	}{
		{0x000000, "000000"},
		{0xFF0000, "ff0000"},
		{0x00FF00, "00ff00"},
		{0x0000FF, "0000ff"},
		{0xABCDEF, "abcdef"},
		{0x123, "000123"},
	}
	for _, tt := range tests {
		if got := PadHex(tt.hex); got != tt.want {
			t.Errorf("PadHex(0x%x): got %q, want %q", tt.hex, got, tt.want)
		}
	}
}

func TestGetVersion(t *testing.T) {
	if GetVersion() == "" {
		t.Error("GetVersion() returned empty string")
	}
}

func TestCmdCtrl(t *testing.T) {
	if got := CmdCtrl(); got != "ctrl" {
		t.Errorf("CmdCtrl(): got %q, want %q", got, "ctrl")
	}
}

// fakeInjector records injected pointer events; used to test position
// tracking without a portal.
type fakeInjector struct {
	abs     []absCall
	rel     []relCall
	axis    []axisCall
	buttons []buttonCall
	syms    []int32 // keyboardKeysym calls (press and release)
	codes   []int32 // keyboardKeycode calls
	err     error
}

type absCall struct {
	stream uint32
	x, y   float64
}

type relCall struct{ dx, dy float64 }

type axisCall struct {
	axis  uint32
	steps int32
}

type buttonCall struct {
	button int32
	state  uint32
}

func (f *fakeInjector) keyboardKeycode(c int32, _ uint32) error {
	f.codes = append(f.codes, c)
	return f.err
}
func (f *fakeInjector) keyboardKeysym(s int32, _ uint32) error {
	f.syms = append(f.syms, s)
	return f.err
}
func (f *fakeInjector) pointerButton(b int32, s uint32) error {
	f.buttons = append(f.buttons, buttonCall{b, s})
	return f.err
}
func (f *fakeInjector) pointerAxisDiscrete(axis uint32, steps int32) error {
	f.axis = append(f.axis, axisCall{axis, steps})
	return f.err
}
func (f *fakeInjector) pointerMotion(dx, dy float64) error {
	f.rel = append(f.rel, relCall{dx, dy})
	return f.err
}
func (f *fakeInjector) pointerMotionAbsolute(s uint32, x, y float64) error {
	f.abs = append(f.abs, absCall{s, x, y})
	return f.err
}

// installFakeConn installs a healthy conn backed by fakeInjector as the global
// connection so the public API never touches D-Bus during tests.
func installFakeConn(t *testing.T, streams ...stream) *fakeInjector {
	t.Helper()
	inj := &fakeInjector{}
	connMu.Lock()
	prev := globalConn
	globalConn = &conn{devices: deviceKeyboard | devicePointer, streams: streams, inj: inj}
	connMu.Unlock()
	t.Cleanup(func() {
		connMu.Lock()
		globalConn = prev
		connMu.Unlock()
	})
	return inj
}

// Location is a pure query (#783): it must not negotiate a portal session
// when there is none, and must keep reporting the last injected position
// after the portal closed the session.
func TestLocationIsReadOnly(t *testing.T) {
	connMu.Lock()
	prev := globalConn
	globalConn = nil
	connMu.Unlock()
	t.Cleanup(func() {
		connMu.Lock()
		globalConn = prev
		connMu.Unlock()
	})

	if x, y := Location(); x != 0 || y != 0 {
		t.Fatalf("Location without session: got (%d,%d), want (0,0)", x, y)
	}
	connMu.Lock()
	if globalConn != nil {
		connMu.Unlock()
		t.Fatal("Location opened a portal session")
	}
	globalConn = &conn{posX: 40, posY: 50, posKnown: true, closed: true}
	connMu.Unlock()

	if x, y := Location(); x != 40 || y != 50 {
		t.Fatalf("Location after session closed: got (%d,%d), want (40,50)", x, y)
	}
}

func TestMoveAbsoluteWithStreams(t *testing.T) {
	inj := installFakeConn(t,
		stream{nodeID: 1, x: 0, y: 0, width: 1920, height: 1080},
		stream{nodeID: 2, x: 1920, y: 0, width: 1280, height: 720},
	)

	if x, y := Location(); x != 0 || y != 0 {
		t.Fatalf("Location before move: got (%d,%d), want (0,0)", x, y)
	}

	if err := Move(100, 200); err != nil {
		t.Fatal(err)
	}
	if x, y := Location(); x != 100 || y != 200 {
		t.Errorf("Location after Move: got (%d,%d), want (100,200)", x, y)
	}

	// Point on the second monitor must be mapped into that stream's space.
	if err := Move(2000, 50); err != nil {
		t.Fatal(err)
	}
	if len(inj.abs) != 2 {
		t.Fatalf("got %d absolute calls, want 2", len(inj.abs))
	}
	if got := inj.abs[1]; got.stream != 2 || got.x != 80 || got.y != 50 {
		t.Errorf("second Move: got %+v, want stream 2 at (80,50)", got)
	}
	if x, y := Location(); x != 2000 || y != 50 {
		t.Errorf("Location: got (%d,%d), want (2000,50)", x, y)
	}

	// Explicit displayId wins over the containing-stream lookup.
	if err := Move(10, 10, 1); err != nil {
		t.Fatal(err)
	}
	if got := inj.abs[2]; got.stream != 2 || got.x != -1910 {
		t.Errorf("Move with displayId=1: got %+v", got)
	}
	if len(inj.rel) != 0 {
		t.Errorf("unexpected relative calls: %+v", inj.rel)
	}
}

func TestMoveRelativeTracksPosition(t *testing.T) {
	inj := installFakeConn(t, stream{nodeID: 1, width: 800, height: 600})

	if err := Move(100, 100); err != nil {
		t.Fatal(err)
	}
	if err := MoveRelative(50, -30); err != nil {
		t.Fatal(err)
	}
	if x, y := Location(); x != 150 || y != 70 {
		t.Errorf("Location: got (%d,%d), want (150,70)", x, y)
	}
	// Clamped to the stream bounds.
	if err := MoveRelative(10000, 10000); err != nil {
		t.Fatal(err)
	}
	if x, y := Location(); x != 799 || y != 599 {
		t.Errorf("Location clamped: got (%d,%d), want (799,599)", x, y)
	}
	if len(inj.rel) != 2 {
		t.Errorf("got %d relative calls, want 2", len(inj.rel))
	}
}

func TestMoveWithoutStreamsFallsBackToRelative(t *testing.T) {
	inj := installFakeConn(t)

	// Unknown position: Move first parks the pointer in the corner, then
	// moves by the delta from (0,0).
	if err := Move(100, 100); err != nil {
		t.Fatal(err)
	}
	if len(inj.abs) != 0 {
		t.Fatalf("unexpected absolute calls: %v", inj.abs)
	}
	want := []relCall{{-cornerReset, -cornerReset}, {100, 100}}
	if len(inj.rel) != 2 || inj.rel[0] != want[0] || inj.rel[1] != want[1] {
		t.Fatalf("Move with unknown position: got %+v, want %+v", inj.rel, want)
	}
	if x, y := Location(); x != 100 || y != 100 {
		t.Errorf("Location: got (%d,%d), want (100,100)", x, y)
	}

	// Once a position is known, Move works as a plain delta.
	if err := MoveRelative(10, 20); err != nil {
		t.Fatal(err)
	}
	if err := Move(100, 100); err != nil {
		t.Fatal(err)
	}
	if len(inj.rel) != 4 || inj.rel[3] != (relCall{-10, -20}) {
		t.Errorf("Move fallback: got %+v, want delta (-10,-20)", inj.rel)
	}
	if x, y := Location(); x != 100 || y != 100 {
		t.Errorf("Location: got (%d,%d), want (100,100)", x, y)
	}

	// MoveSmooth with an unknown start and no stream jumps via Move.
	inj2 := installFakeConn(t)
	if !MoveSmooth(50, 50, 2, 0) {
		t.Error("MoveSmooth with unknown position must succeed via corner reset")
	}
	if x, y := Location(); x != 50 || y != 50 {
		t.Errorf("Location after MoveSmooth: got (%d,%d), want (50,50)", x, y)
	}
	if len(inj2.rel) != 2 {
		t.Errorf("MoveSmooth jump: got %d relative calls, want 2", len(inj2.rel))
	}
}

func TestMoveIntoMonitorGapClamps(t *testing.T) {
	inj := installFakeConn(t,
		stream{nodeID: 1, x: 0, y: 0, width: 1000, height: 1000},
		stream{nodeID: 2, x: 1500, y: 0, width: 1000, height: 1000},
	)
	// (1200, 50) lies in the gap: clamp onto stream 0 instead of sending
	// an out-of-range local coordinate.
	if err := Move(1200, 50); err != nil {
		t.Fatal(err)
	}
	if len(inj.abs) != 1 || inj.abs[0] != (absCall{1, 999, 50}) {
		t.Errorf("gap Move: got %+v, want stream 1 at (999,50)", inj.abs)
	}
	if x, y := Location(); x != 999 || y != 50 {
		t.Errorf("Location: got (%d,%d), want (999,50)", x, y)
	}
}

func TestScrollSignConvention(t *testing.T) {
	inj := installFakeConn(t)
	// robotgo: positive y = up, positive x = left; the portal counts
	// positive steps as down/right.
	if err := Scroll(2, 3, 0); err != nil {
		t.Fatal(err)
	}
	if err := ScrollDir(1, "down"); err != nil {
		t.Fatal(err)
	}
	if err := ScrollDir(1, "right"); err != nil {
		t.Fatal(err)
	}
	want := []axisCall{
		{axisVertical, -3}, {axisHorizontal, -2},
		{axisVertical, 1}, {axisHorizontal, 1},
	}
	if len(inj.axis) != len(want) {
		t.Fatalf("got %d axis calls, want %d: %+v", len(inj.axis), len(want), inj.axis)
	}
	for i := range want {
		if inj.axis[i] != want[i] {
			t.Errorf("axis call %d: got %+v, want %+v", i, inj.axis[i], want[i])
		}
	}
}

func TestResolveKeyUppercase(t *testing.T) {
	code, shift, ok := resolveKey("A")
	if !ok || !shift || code != 30 {
		t.Errorf("resolveKey(A): got (%d,%v,%v), want (30,true,true)", code, shift, ok)
	}
	if _, shift, ok := resolveKey("a"); !ok || shift {
		t.Error("resolveKey(a) must not imply shift")
	}
	if _, _, ok := resolveKey("É"); ok {
		t.Error("resolveKey(É) must be unknown")
	}
}

func TestMoveSmoothAbsoluteTarget(t *testing.T) {
	inj := installFakeConn(t, stream{nodeID: 1, width: 1000, height: 1000})

	if err := Move(0, 0); err != nil {
		t.Fatal(err)
	}
	if !MoveSmooth(100, 50, 4, 0) {
		t.Fatal("MoveSmooth returned false")
	}
	if x, y := Location(); x != 100 || y != 50 {
		t.Errorf("Location: got (%d,%d), want (100,50)", x, y)
	}
	want := []absCall{{1, 0, 0}, {1, 25, 12}, {1, 50, 25}, {1, 75, 37}, {1, 100, 50}}
	if len(inj.abs) != len(want) {
		t.Fatalf("got %d absolute calls, want %d: %+v", len(inj.abs), len(want), inj.abs)
	}
	for i := range want {
		if inj.abs[i] != want[i] {
			t.Errorf("step %d: got %+v, want %+v", i, inj.abs[i], want[i])
		}
	}

	// Unknown start + stream: single jump.
	inj2 := installFakeConn(t, stream{nodeID: 1, width: 1000, height: 1000})
	if !MoveSmooth(300, 300, 5, 0) {
		t.Fatal("MoveSmooth jump returned false")
	}
	if len(inj2.abs) != 1 {
		t.Errorf("jump: got %d absolute calls, want 1", len(inj2.abs))
	}
}

func TestMoveRelativeUnknownPosition(t *testing.T) {
	for _, linked := range []bool{false, true} {
		t.Run(map[bool]string{false: "relative fallback", true: "linked stream"}[linked], func(t *testing.T) {
			var streams []stream
			if linked {
				streams = []stream{{nodeID: 1, width: 800, height: 600}}
			}
			inj := installFakeConn(t, streams...)
			if err := MoveRelative(10, 20); err != nil {
				t.Fatal(err)
			}
			if len(inj.rel) != 1 || inj.rel[0] != (relCall{10, 20}) {
				t.Fatalf("relative injection: got %+v", inj.rel)
			}
			if x, y, known := globalConn.position(); known || x != 0 || y != 0 {
				t.Errorf("relative motion cannot establish an absolute position: (%d,%d), known=%v", x, y, known)
			}
			if err := Move(100, 100); err != nil {
				t.Fatal(err)
			}
			if !linked && (len(inj.rel) != 3 || inj.rel[1] != (relCall{-cornerReset, -cornerReset}) || inj.rel[2] != (relCall{100, 100})) {
				t.Errorf("Move must still reset the unknown origin: %+v", inj.rel)
			}
			if x, y := Location(); x != 100 || y != 100 {
				t.Errorf("Location after absolute move: (%d,%d)", x, y)
			}
		})
	}
}

type failTargetInjector struct{ fakeInjector }

func (f *failTargetInjector) pointerMotion(dx, dy float64) error {
	if len(f.rel) == 1 {
		f.err = ErrNotSupported
	}
	return f.fakeInjector.pointerMotion(dx, dy)
}

func TestMoveSmoothUnknownStartFailure(t *testing.T) {
	installFakeConn(t)
	inj := &failTargetInjector{}
	globalConn.inj = inj
	if MoveSmooth(50, 50, 2, 0) {
		t.Error("MoveSmooth reported success after target injection failed")
	}
	if len(inj.rel) != 2 {
		t.Fatalf("got %d calls, want reset then target", len(inj.rel))
	}
	if x, y, known := globalConn.position(); !known || x != 0 || y != 0 {
		t.Errorf("successful reset position lost: (%d,%d), known=%v", x, y, known)
	}
}

// Target (0,0) coincides with the reset origin, so the tracked position alone
// cannot reveal that the target injection failed.
func TestMoveSmoothUnknownStartOriginFailure(t *testing.T) {
	installFakeConn(t)
	inj := &failTargetInjector{}
	globalConn.inj = inj
	if MoveSmooth(0, 0, 2, 0) {
		t.Error("MoveSmooth reported success after target injection to (0,0) failed")
	}
	if len(inj.rel) != 2 {
		t.Fatalf("got %d calls, want reset then target", len(inj.rel))
	}
}

func TestMoveSmoothUnknownStartClampedTarget(t *testing.T) {
	installFakeConn(t, stream{nodeID: 1, width: 100, height: 100})
	if MoveSmooth(200, 200, 2, 0) {
		t.Error("MoveSmooth reported reaching a target outside the stream")
	}
	if x, y := Location(); x != 99 || y != 99 {
		t.Errorf("clamped position: (%d,%d), want (99,99)", x, y)
	}
}

func TestMoveInjectErrorKeepsPosition(t *testing.T) {
	inj := installFakeConn(t, stream{nodeID: 1, width: 1000, height: 1000})
	if err := Move(10, 10); err != nil {
		t.Fatal(err)
	}
	inj.err = ErrNotSupported
	if err := Move(500, 500); !errors.Is(err, ErrNotSupported) {
		t.Errorf("failed Move: got %v, want ErrNotSupported", err)
	}
	if err := MoveRelative(5, 5); !errors.Is(err, ErrNotSupported) {
		t.Errorf("failed MoveRelative: got %v, want ErrNotSupported", err)
	}
	if x, y := Location(); x != 10 || y != 10 {
		t.Errorf("Location after failed moves: got (%d,%d), want (10,10)", x, y)
	}
}

// Pointer injection failures must be returned instead of dropped, and
// ScrollSmooth/MoveClick must stop at the first one.
func TestPointerErrorsPropagate(t *testing.T) {
	inj := installFakeConn(t, stream{nodeID: 1, width: 1000, height: 1000})
	inj.err = ErrNotSupported

	if err := Scroll(0, 1); !errors.Is(err, ErrNotSupported) {
		t.Errorf("Scroll: got %v, want ErrNotSupported", err)
	}
	if err := ScrollDir(1, "left"); !errors.Is(err, ErrNotSupported) {
		t.Errorf("ScrollDir: got %v, want ErrNotSupported", err)
	}
	inj.axis = nil
	if err := ScrollSmooth(1, 3, 0); !errors.Is(err, ErrNotSupported) {
		t.Errorf("ScrollSmooth: got %v, want ErrNotSupported", err)
	}
	if len(inj.axis) != 1 {
		t.Errorf("ScrollSmooth scrolled %d times, want 1 (stop at the first error)", len(inj.axis))
	}
	if err := MoveClick(10, 10); !errors.Is(err, ErrNotSupported) {
		t.Errorf("MoveClick: got %v, want ErrNotSupported", err)
	}
	if len(inj.buttons) != 0 {
		t.Errorf("MoveClick clicked after a failed move: %+v", inj.buttons)
	}

	// A failed press makes DragSmooth return before moving or releasing.
	inj.abs = nil
	if err := DragSmooth(20, 20); !errors.Is(err, ErrNotSupported) {
		t.Errorf("DragSmooth: got %v, want ErrNotSupported", err)
	}
	if len(inj.buttons) != 1 || len(inj.abs) != 0 {
		t.Errorf("DragSmooth after failed press: buttons %+v, moves %+v", inj.buttons, inj.abs)
	}

	// Without a granted pointer device nothing is injected.
	inj.err = nil
	globalConn.devices = deviceKeyboard
	if err := Move(1, 1); !errors.Is(err, ErrNotSupported) {
		t.Errorf("Move without pointer device: got %v, want ErrNotSupported", err)
	}
}

// stepFailInjector fails only the second absolute motion.
type stepFailInjector struct{ fakeInjector }

func (f *stepFailInjector) pointerMotionAbsolute(s uint32, x, y float64) error {
	f.abs = append(f.abs, absCall{s, x, y})
	if len(f.abs) == 2 {
		return ErrNotSupported
	}
	return nil
}

// A failed intermediate step must fail MoveSmooth even though later steps
// would still reach the target.
func TestMoveSmoothStopsOnStepError(t *testing.T) {
	installFakeConn(t, stream{nodeID: 1, width: 1000, height: 1000})
	inj := &stepFailInjector{}
	globalConn.inj = inj
	if err := Move(0, 0); err != nil {
		t.Fatal(err)
	}
	if MoveSmooth(100, 50, 4, 0) {
		t.Error("MoveSmooth reported success after a failed step")
	}
	if len(inj.abs) != 2 {
		t.Errorf("MoveSmooth kept moving after a failed step: %d calls, want 2", len(inj.abs))
	}
}

var errMoveFailed = errors.New("move failed")

// moveFailInjector fails every pointer motion while buttons still work.
type moveFailInjector struct{ fakeInjector }

func (f *moveFailInjector) pointerMotion(float64, float64) error { return errMoveFailed }
func (f *moveFailInjector) pointerMotionAbsolute(uint32, float64, float64) error {
	return errMoveFailed
}

// DragSmooth must release the button even when the move fails, and return
// the move error.
func TestDragSmoothReleasesOnMoveError(t *testing.T) {
	installFakeConn(t)
	inj := &moveFailInjector{}
	globalConn.inj = inj
	if err := DragSmooth(100, 100, "right"); !errors.Is(err, errMoveFailed) {
		t.Errorf("DragSmooth: got %v, want the move error", err)
	}
	want := []buttonCall{{btnRight, statePressed}, {btnRight, stateReleased}}
	if !slices.Equal(inj.buttons, want) {
		t.Errorf("buttons: got %+v, want press then release %+v", inj.buttons, want)
	}
}

// ScrollDir must reject directions other than up/down/left/right instead of
// silently doing nothing; no direction means down.
func TestScrollDirUnknown(t *testing.T) {
	inj := installFakeConn(t)
	tests := []struct {
		dir  interface{}
		want string
	}{
		{"diagonal", "robotgo: unknown scroll direction: diagonal"},
		{"Up", "robotgo: unknown scroll direction: Up"},
		{42, "robotgo: unknown scroll direction: 42"},
		{nil, "robotgo: unknown scroll direction: <nil>"},
	}
	for _, tt := range tests {
		if err := ScrollDir(1, tt.dir); err == nil || err.Error() != tt.want {
			t.Errorf("ScrollDir(1, %#v) = %v, want %q", tt.dir, err, tt.want)
		}
	}
	if len(inj.axis) != 0 {
		t.Errorf("unknown directions scrolled: %+v", inj.axis)
	}
	if err := ScrollDir(2); err != nil {
		t.Fatal(err)
	}
	if len(inj.axis) != 1 || inj.axis[0] != (axisCall{axisVertical, 2}) {
		t.Errorf("ScrollDir(2) default: got %+v, want one down step of 2", inj.axis)
	}
}

// keysymFailInjector accepts n keysym events, then fails.
type keysymFailInjector struct {
	fakeInjector
	n int
}

func (f *keysymFailInjector) keyboardKeysym(s int32, state uint32) error {
	if len(f.syms) == f.n {
		return ErrNotSupported
	}
	return f.fakeInjector.keyboardKeysym(s, state)
}

// TypeStr/TypeDelay must report how many characters (runes) were typed when
// typing stops early, and TypeDelay must restore KeySleep.
func TestTypeStrPartial(t *testing.T) {
	installFakeConn(t)
	old := KeySleep
	t.Cleanup(func() { KeySleep = old })
	KeySleep = 0

	// 4 keysym events = press+release of the first 2 runes.
	globalConn.inj = &keysymFailInjector{n: 4}
	if err := TypeStr("hé€lo"); err == nil || err.Error() != "robotgo: typed 2 of 5 characters" {
		t.Errorf("TypeStr: got %v, want typed 2 of 5", err)
	}

	KeySleep = 7
	globalConn.inj = &keysymFailInjector{n: 2}
	if err := TypeDelay("abc", 0); err == nil || err.Error() != "robotgo: typed 1 of 3 characters" {
		t.Errorf("TypeDelay: got %v, want typed 1 of 3", err)
	}
	if KeySleep != 7 {
		t.Errorf("TypeDelay left KeySleep = %d, want 7", KeySleep)
	}

	globalConn.inj = &fakeInjector{}
	if err := TypeStr("hé€lo"); err != nil {
		t.Errorf("TypeStr: got %v, want nil", err)
	}
}

func TestScreenGeometryFromStreams(t *testing.T) {
	installFakeConn(t,
		stream{nodeID: 1, x: 0, y: 0, width: 1920, height: 1080},
		stream{nodeID: 2, x: 1920, y: -200, width: 1280, height: 720},
	)
	if w, h := GetScreenSize(); w != 3200 || h != 1280 {
		t.Errorf("GetScreenSize: got (%d,%d), want (3200,1280)", w, h)
	}
	if n := DisplaysNum(); n != 2 {
		t.Errorf("DisplaysNum: got %d, want 2", n)
	}
	r := GetScreenRect(1)
	if r.X != 1920 || r.Y != -200 || r.W != 1280 || r.H != 720 {
		t.Errorf("GetScreenRect(1): got %+v", r)
	}
	if r := GetScreenRect(7); r != GetScreenRect(0) {
		t.Errorf("GetScreenRect(out of range): got %+v, want display 0", r)
	}
}

func TestCloseMarksUnhealthy(t *testing.T) {
	installFakeConn(t)
	c := globalConn
	if !c.healthy() {
		t.Fatal("fresh conn must be healthy")
	}
	Close()
	if c.healthy() {
		t.Error("closed conn must not be healthy")
	}
	if globalConn != nil {
		t.Error("Close must drop the global conn")
	}
}

func TestUnsupportedSurface(t *testing.T) {
	installFakeConn(t)
	// Screen capture and window management are intentionally unsupported.
	if _, err := CaptureImg(); err != ErrNotSupported {
		t.Errorf("CaptureImg: got %v, want ErrNotSupported", err)
	}
	if _, err := Capture(); err != ErrNotSupported {
		t.Errorf("Capture: got %v, want ErrNotSupported", err)
	}
	if err := ActiveName("x"); err != ErrNotSupported {
		t.Errorf("ActiveName: got %v, want ErrNotSupported", err)
	}
	if err := ActivePid(1); err != ErrNotSupported {
		t.Errorf("ActivePid: got %v, want ErrNotSupported", err)
	}
	if err := MinWindow(1); err != ErrNotSupported {
		t.Errorf("MinWindow: got %v, want ErrNotSupported", err)
	}
	if err := MaxWindow(1); err != ErrNotSupported {
		t.Errorf("MaxWindow: got %v, want ErrNotSupported", err)
	}
	if err := CloseWindow(); err != ErrNotSupported {
		t.Errorf("CloseWindow: got %v, want ErrNotSupported", err)
	}
	if x, y, w, h := GetBounds(1); x != 0 || y != 0 || w != 0 || h != 0 {
		t.Errorf("GetBounds: got (%d,%d,%d,%d), want zeros", x, y, w, h)
	}
	if x, y, w, h := GetClient(1); x != 0 || y != 0 || w != 0 || h != 0 {
		t.Errorf("GetClient: got (%d,%d,%d,%d), want zeros", x, y, w, h)
	}
	if got := GetPixelColor(0, 0); got != "000000" {
		t.Errorf("GetPixelColor: got %q, want 000000", got)
	}
	if w, h := GetScreenSize(); w != 0 || h != 0 {
		t.Errorf("GetScreenSize without streams: got (%d,%d), want (0,0)", w, h)
	}
}

// #640: Type must send every rune as a keysym (layout independent), never as
// a US-layout keycode, so '@' and '/' come out right on a German layout.
func TestTypeUsesKeysyms(t *testing.T) {
	inj := installFakeConn(t)
	old := KeySleep
	KeySleep = 0
	t.Cleanup(func() { KeySleep = old })

	const text = `test@example.org/ <>|{}\~7Qzy"'`
	if n := Type(text); n != len(text) {
		t.Fatalf("Type returned %d, want %d", n, len(text))
	}
	if len(inj.codes) != 0 {
		t.Errorf("Type sent %d raw keycodes, want 0: %v", len(inj.codes), inj.codes)
	}
	var want []int32
	for _, r := range text {
		ks := runeToKeysym(r)
		want = append(want, ks, ks) // press, release
	}
	if len(inj.syms) != len(want) {
		t.Fatalf("got %d keysym events, want %d", len(inj.syms), len(want))
	}
	for i := range want {
		if inj.syms[i] != want[i] {
			t.Fatalf("keysym event %d = 0x%x, want 0x%x", i, inj.syms[i], want[i])
		}
	}
}

func TestTypes(t *testing.T) {
	p := Point{X: 1, Y: 2}
	if p.X != 1 || p.Y != 2 {
		t.Errorf("Point: got %+v", p)
	}
	s := Size{W: 100, H: 200}
	if s.W != 100 || s.H != 200 {
		t.Errorf("Size: got %+v", s)
	}
	r := Rect{Point: p, Size: s}
	if r.X != 1 || r.W != 100 {
		t.Errorf("Rect: got %+v", r)
	}
	n := Nps{Pid: 42, Name: "test"}
	if n.Pid != 42 || n.Name != "test" {
		t.Errorf("Nps: got %+v", n)
	}
}

// --- Process tests (work on any Linux) ---

func TestPids(t *testing.T) {
	pids, err := Pids()
	if err != nil {
		t.Skipf("Pids() error: %v", err)
	}
	if len(pids) == 0 {
		t.Error("Pids() returned empty list")
	}
}

func TestGetPid(t *testing.T) {
	if pid := GetPid(); pid <= 0 {
		t.Errorf("GetPid() returned %d", pid)
	}
}

func TestMainDisplayID(t *testing.T) {
	if got := MainDisplayID(); got != 0 {
		t.Errorf("MainDisplayID: got %d, want 0", got)
	}
}
