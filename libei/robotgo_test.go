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

func TestUnsupportedSurface(t *testing.T) {
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
		t.Errorf("GetScreenSize: got (%d,%d), want (0,0)", w, h)
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
