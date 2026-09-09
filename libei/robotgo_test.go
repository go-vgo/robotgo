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
	abs  []absCall
	rel  []relCall
	axis []axisCall
	err  error
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

func (f *fakeInjector) keyboardKeycode(int32, uint32) error { return nil }
func (f *fakeInjector) keyboardKeysym(int32, uint32) error  { return nil }
func (f *fakeInjector) pointerButton(int32, uint32) error   { return nil }
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

func TestMoveAbsoluteWithStreams(t *testing.T) {
	inj := installFakeConn(t,
		stream{nodeID: 1, x: 0, y: 0, width: 1920, height: 1080},
		stream{nodeID: 2, x: 1920, y: 0, width: 1280, height: 720},
	)

	if x, y := Location(); x != 0 || y != 0 {
		t.Fatalf("Location before move: got (%d,%d), want (0,0)", x, y)
	}

	Move(100, 200)
	if x, y := Location(); x != 100 || y != 200 {
		t.Errorf("Location after Move: got (%d,%d), want (100,200)", x, y)
	}

	// Point on the second monitor must be mapped into that stream's space.
	Move(2000, 50)
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
	Move(10, 10, 1)
	if got := inj.abs[2]; got.stream != 2 || got.x != -1910 {
		t.Errorf("Move with displayId=1: got %+v", got)
	}
	if len(inj.rel) != 0 {
		t.Errorf("unexpected relative calls: %+v", inj.rel)
	}
}

func TestMoveRelativeTracksPosition(t *testing.T) {
	inj := installFakeConn(t, stream{nodeID: 1, width: 800, height: 600})

	Move(100, 100)
	MoveRelative(50, -30)
	if x, y := Location(); x != 150 || y != 70 {
		t.Errorf("Location: got (%d,%d), want (150,70)", x, y)
	}
	// Clamped to the stream bounds.
	MoveRelative(10000, 10000)
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
	Move(100, 100)
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
	MoveRelative(10, 20)
	Move(100, 100)
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
	Move(1200, 50)
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
	Scroll(2, 3, 0)
	ScrollDir(1, "down")
	ScrollDir(1, "right")
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

	Move(0, 0)
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

func TestMoveInjectErrorKeepsPosition(t *testing.T) {
	inj := installFakeConn(t, stream{nodeID: 1, width: 1000, height: 1000})
	Move(10, 10)
	inj.err = ErrNotSupported
	Move(500, 500)
	MoveRelative(5, 5)
	if x, y := Location(); x != 10 || y != 10 {
		t.Errorf("Location after failed moves: got (%d,%d), want (10,10)", x, y)
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
	if got := GetPixelColor(0, 0); got != "000000" {
		t.Errorf("GetPixelColor: got %q, want 000000", got)
	}
	if w, h := GetScreenSize(); w != 0 || h != 0 {
		t.Errorf("GetScreenSize without streams: got (%d,%d), want (0,0)", w, h)
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
