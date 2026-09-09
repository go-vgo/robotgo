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
	"testing"
)

// --- Pure Go tests (run anywhere, no Wayland needed) ---

func TestKeyToEvdev(t *testing.T) {
	tests := []struct {
		key  string
		code uint32
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
		want int
	}{
		{"left", btnLeft},
		{"right", btnRight},
		{"center", btnMiddle},
		{"middle", btnMiddle},
		{"", btnLeft},
		{"unknown", btnLeft},
	}

	for _, tt := range tests {
		got := resolveButton(tt.btn)
		if got != tt.want {
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

// TestPidIgnored documents that an int pid argument is accepted for API parity
// but ignored: it must not be mistaken for a modifier (Wayland injects into the
// focused surface, mirroring the X11 path in key/keypress_c.h).
func TestPidIgnored(t *testing.T) {
	got := extractModifiers([]interface{}{"ctrl", 1234, "shift"})
	want := []string{"ctrl", "shift"}
	if len(got) != len(want) {
		t.Fatalf("extractModifiers dropped/added entries: got %v, want %v", got, want)
	}
	for i := range got {
		if got[i] != want[i] {
			t.Errorf("extractModifiers[%d]: got %q, want %q", i, got[i], want[i])
		}
	}
}

// TestReleaseKeys verifies the upKeyArr-equivalent helper: modifiers are
// keyed up in reverse order, every code gets a release attempt even after a
// failure, and errors are propagated.
func TestReleaseKeys(t *testing.T) {
	var got []uint32
	err := releaseKeys([]uint32{29, 42, 56}, func(code uint32) error {
		got = append(got, code)
		return nil
	})
	if err != nil {
		t.Fatalf("releaseKeys: unexpected error %v", err)
	}
	want := []uint32{56, 42, 29}
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
	failOn := uint32(42)
	err = releaseKeys([]uint32{29, 42, 56}, func(code uint32) error {
		got = append(got, code)
		if code == failOn {
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

	if err := releaseKeys(nil, func(uint32) error { return ErrNotSupported }); err != nil {
		t.Errorf("releaseKeys(nil): got %v, want nil", err)
	}
}

func TestShiftedChars(t *testing.T) {
	// Verify all shifted chars map to valid evdev keys
	for ch, baseKey := range shiftedChars {
		_, ok := keyToEvdev(baseKey)
		if !ok {
			t.Errorf("shiftedChars[%q] = %q, but %q has no evdev mapping", string(ch), baseKey, baseKey)
		}
	}
}

func TestIsActivated(t *testing.T) {
	tests := []struct {
		name   string
		states []byte
		want   bool
	}{
		{"empty", nil, false},
		{"activated only", []byte{2, 0, 0, 0}, true},
		{"maximized then activated", []byte{0, 0, 0, 0, 2, 0, 0, 0}, true},
		{"maximized only", []byte{0, 0, 0, 0}, false},
		{"minimized", []byte{1, 0, 0, 0}, false},
		{"short data", []byte{2, 0}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isActivated(tt.states)
			if got != tt.want {
				t.Errorf("isActivated(%v): got %v, want %v", tt.states, got, tt.want)
			}
		})
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
		got := PadHex(tt.hex)
		if got != tt.want {
			t.Errorf("PadHex(0x%x): got %q, want %q", tt.hex, got, tt.want)
		}
	}
}

func TestGetVersion(t *testing.T) {
	v := GetVersion()
	if v == "" {
		t.Error("GetVersion() returned empty string")
	}
}

func TestCmdCtrl(t *testing.T) {
	got := CmdCtrl()
	if got != "ctrl" {
		t.Errorf("CmdCtrl(): got %q, want %q", got, "ctrl")
	}
}

func TestTypes(t *testing.T) {
	// Verify types exist and are constructible
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
		t.Skipf("Pids() error (may not be on Linux): %v", err)
	}
	if len(pids) == 0 {
		t.Error("Pids() returned empty list")
	}
	// PID 1 should always exist
	found := false
	for _, pid := range pids {
		if pid == 1 {
			found = true
			break
		}
	}
	if !found {
		t.Error("Pids() didn't contain PID 1")
	}
}

func TestPidExists(t *testing.T) {
	exists, err := PidExists(1)
	if err != nil {
		t.Skipf("PidExists(1) error: %v", err)
	}
	if !exists {
		t.Error("PidExists(1) returned false")
	}

	exists, err = PidExists(99999999)
	if err != nil {
		t.Skipf("PidExists(99999999) error: %v", err)
	}
	if exists {
		t.Error("PidExists(99999999) returned true")
	}
}

func TestFindName(t *testing.T) {
	name, err := FindName(1)
	if err != nil {
		t.Skipf("FindName(1) error: %v", err)
	}
	if name == "" {
		t.Error("FindName(1) returned empty string")
	}
}

func TestGetPid(t *testing.T) {
	pid := GetPid()
	if pid <= 0 {
		t.Errorf("GetPid() returned %d", pid)
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

func TestModMask(t *testing.T) {
	tests := map[string]uint32{
		"shift": modShift, "shiftr": modShift,
		"ctrl": modControl, "ctrlr": modControl,
		"alt": modAlt, "altr": modAlt,
		"cmd": modSuper, "cmdr": modSuper,
	}
	for name, want := range tests {
		code, _ := keyToEvdev(name)
		if got := modMask[code]; got != want {
			t.Errorf("modMask[%s]: got %d, want %d", name, got, want)
		}
	}
	if _, isMod := modMask[30]; isMod {
		t.Error("KEY_A must not be a modifier")
	}
}

func TestDecodeShm(t *testing.T) {
	// 2x2, stride 12 (one padding pixel per row), rows: top then bottom.
	px := func(a, b, c, d byte) []byte { return []byte{a, b, c, d} }
	row := func(p1, p2 []byte) []byte { return append(append(append([]byte{}, p1...), p2...), 0, 0, 0, 0) }
	top := row(px(1, 2, 3, 0x80), px(4, 5, 6, 0x80))
	bottom := row(px(7, 8, 9, 0x80), px(10, 11, 12, 0x80))
	data := append(append([]byte{}, top...), bottom...)

	// XRGB: bytes are B,G,R,X -> opaque, swapped.
	img := decodeShm(data, shmXRGB8888, 2, 2, 12, false)
	if c := img.RGBAAt(0, 0); c.R != 3 || c.G != 2 || c.B != 1 || c.A != 0xff {
		t.Errorf("XRGB (0,0): got %+v", c)
	}
	// ARGB keeps alpha.
	if c := decodeShm(data, shmARGB8888, 2, 2, 12, false).RGBAAt(1, 0); c.R != 6 || c.A != 0x80 {
		t.Errorf("ARGB (1,0): got %+v", c)
	}
	// XBGR: bytes are R,G,B,X -> no swap.
	if c := decodeShm(data, shmXBGR8888, 2, 2, 12, false).RGBAAt(0, 1); c.R != 7 || c.B != 9 || c.A != 0xff {
		t.Errorf("XBGR (0,1): got %+v", c)
	}
	// y_invert flips rows.
	if c := decodeShm(data, shmABGR8888, 2, 2, 12, true).RGBAAt(0, 0); c.R != 7 || c.A != 0x80 {
		t.Errorf("y_invert (0,0): got %+v", c)
	}
}

func TestOutputBounds(t *testing.T) {
	c := &conn{outputs: []*outputInfo{
		{x: 0, y: 0, width: 1920, height: 1080},
		{x: 1920, y: -200, width: 1280, height: 720},
	}}
	x, y, w, h := c.outputBounds()
	if x != 0 || y != -200 || w != 3200 || h != 1280 {
		t.Errorf("outputBounds: got (%d,%d,%d,%d)", x, y, w, h)
	}
	if o, ok := c.output(5); !ok || o.width != 1920 {
		t.Errorf("output(5) must fall back to output 0, got %+v %v", o, ok)
	}
	if _, ok := (&conn{}).output(0); ok {
		t.Error("output on empty conn must report !ok")
	}
}

func TestActiveToplevel(t *testing.T) {
	c := &conn{toplevels: map[uint32]*toplevelInfo{}}
	if c.activeToplevel() != nil {
		t.Error("empty toplevels must yield nil")
	}
	inactive := &toplevelInfo{title: "bg"}
	active := &toplevelInfo{title: "fg", states: []byte{2, 0, 0, 0}}
	c.toplevels[1] = inactive
	c.toplevels[2] = active
	for i := 0; i < 20; i++ { // map order is random; must be stable
		if got := c.activeToplevel(); got != active {
			t.Fatalf("activeToplevel: got %q, want fg", got.title)
		}
	}
}

func TestDoRejectsClosedConn(t *testing.T) {
	c := &conn{dispatchDone: make(chan struct{})}
	c.closed.Store(true)
	if err := c.do(func() error { return nil }); err != ErrNoConnection {
		t.Errorf("closed conn: got %v, want ErrNoConnection", err)
	}
	c2 := &conn{dispatchDone: make(chan struct{})}
	close(c2.dispatchDone)
	if err := c2.do(func() error { return nil }); err != ErrNoConnection {
		t.Errorf("dead dispatch loop: got %v, want ErrNoConnection", err)
	}
	c3 := &conn{dispatchDone: make(chan struct{})}
	if err := c3.do(func() error { return nil }); err != nil {
		t.Errorf("live conn: got %v", err)
	}
}
