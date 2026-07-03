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
