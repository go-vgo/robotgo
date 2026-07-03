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

package x11

import "testing"

// --- Pure Go tests (run anywhere, no X server needed) ---

func TestKeyKeysym(t *testing.T) {
	tests := []struct {
		key string
		ks  uint32
		ok  bool
	}{
		{"a", 0x61, true},
		{"A", 0x41, true},
		{"0", 0x30, true},
		{"enter", xkReturn, true},
		{"esc", xkEscape, true},
		{"escape", xkEscape, true},
		{"f1", xkF1, true},
		{"f12", xkF1 + 11, true},
		{"f24", xkF1 + 23, true},
		{"shift", xkShiftL, true},
		{"shiftr", xkShiftR, true},
		{"ctrl", xkControlL, true},
		{"alt", xkAltL, true},
		{"cmd", xkSuperL, true},
		{"space", xkSpace, true},
		{"tab", xkTab, true},
		{"backspace", xkBackSpace, true},
		{"delete", xkDelete, true},
		{"up", xkUp, true},
		{"num5", xkKP0 + 5, true},
		{"audio_mute", xkAudioMute, true},
		{"", 0, false},
		{"no_such_key", 0, false},
	}
	for _, tt := range tests {
		ks, ok := keyKeysym(tt.key)
		if ok != tt.ok {
			t.Errorf("keyKeysym(%q): got ok=%v, want %v", tt.key, ok, tt.ok)
			continue
		}
		if ok && ks != tt.ks {
			t.Errorf("keyKeysym(%q): got 0x%x, want 0x%x", tt.key, ks, tt.ks)
		}
	}
}

func TestRuneKeysym(t *testing.T) {
	tests := []struct {
		r  rune
		ks uint32
	}{
		{'a', 0x61},
		{'Z', 0x5a},
		{' ', 0x20},
		{'é', 0xe9},       // Latin-1 direct
		{'€', 0x010020ac}, // Unicode keysym range
		{'中', 0x01004e2d}, // Unicode keysym range
	}
	for _, tt := range tests {
		if got := runeKeysym(tt.r); got != tt.ks {
			t.Errorf("runeKeysym(%q): got 0x%x, want 0x%x", tt.r, got, tt.ks)
		}
	}
}

func TestResolveButton(t *testing.T) {
	tests := []struct {
		btn  string
		want byte
	}{
		{"left", btnLeft},
		{"right", btnRight},
		{"center", btnMiddle},
		{"middle", btnMiddle},
		{"wheelUp", btnWheelUp},
		{"wheelDown", btnWheelDown},
		{"", btnLeft},
		{"unknown", btnLeft},
	}
	for _, tt := range tests {
		if got := resolveButton(tt.btn); got != tt.want {
			t.Errorf("resolveButton(%q): got %d, want %d", tt.btn, got, tt.want)
		}
	}
}

func TestExtractMods(t *testing.T) {
	tests := []struct {
		name   string
		args   []interface{}
		want   []string
		down   bool
		hasDir bool
	}{
		{"empty", nil, nil, true, false},
		{"ctrl", []interface{}{"ctrl"}, []string{"ctrl"}, true, false},
		{"ctrl+shift", []interface{}{"ctrl", "shift"}, []string{"ctrl", "shift"}, true, false},
		{"slice", []interface{}{[]string{"ctrl", "alt"}}, []string{"ctrl", "alt"}, true, false},
		{"up dir", []interface{}{"up"}, nil, false, true},
		{"down dir + mod", []interface{}{"down", "ctrl"}, []string{"ctrl"}, true, true},
		{"ignore ints", []interface{}{"ctrl", 42, "shift"}, []string{"ctrl", "shift"}, true, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mods, down, hasDir := extractMods(tt.args)
			if len(mods) != len(tt.want) {
				t.Fatalf("mods: got %v, want %v", mods, tt.want)
			}
			for i := range mods {
				if mods[i] != tt.want[i] {
					t.Errorf("mods[%d]: got %q, want %q", i, mods[i], tt.want[i])
				}
			}
			if down != tt.down {
				t.Errorf("down: got %v, want %v", down, tt.down)
			}
			if hasDir != tt.hasDir {
				t.Errorf("hasDir: got %v, want %v", hasDir, tt.hasDir)
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
		if got := PadHex(tt.hex); got != tt.want {
			t.Errorf("PadHex(0x%x): got %q, want %q", tt.hex, got, tt.want)
		}
	}
}

func TestAtoiSafe(t *testing.T) {
	tests := []struct {
		s    string
		want int
	}{
		{"0", 0},
		{"12", 12},
		{"007", 7},
		{"", -1},
		{"a1", -1},
		{"1a", -1},
	}
	for _, tt := range tests {
		if got := atoiSafe(tt.s); got != tt.want {
			t.Errorf("atoiSafe(%q): got %d, want %d", tt.s, got, tt.want)
		}
	}
}

func TestZpixmapToRGBA(t *testing.T) {
	// 2x1 image, 4 bytes/pixel BGRX: pixel0 = red, pixel1 = blue.
	data := []byte{
		0x00, 0x00, 0xff, 0x00, // B=0 G=0 R=255 -> red
		0xff, 0x00, 0x00, 0x00, // B=255 G=0 R=0 -> blue
	}
	img := zpixmapToRGBA(data, 2, 1)
	r, g, b, a := img.At(0, 0).RGBA()
	if r>>8 != 0xff || g>>8 != 0 || b>>8 != 0 || a>>8 != 0xff {
		t.Errorf("pixel0: got rgba(%d,%d,%d,%d)", r>>8, g>>8, b>>8, a>>8)
	}
	r, g, b, _ = img.At(1, 0).RGBA()
	if r>>8 != 0 || g>>8 != 0 || b>>8 != 0xff {
		t.Errorf("pixel1: got rgb(%d,%d,%d)", r>>8, g>>8, b>>8)
	}
}

func TestTypes(t *testing.T) {
	p := Point{X: 1, Y: 2}
	s := Size{W: 100, H: 200}
	r := Rect{Point: p, Size: s}
	if r.X != 1 || r.W != 100 {
		t.Errorf("Rect: got %+v", r)
	}
	n := Nps{Pid: 42, Name: "test"}
	if n.Pid != 42 || n.Name != "test" {
		t.Errorf("Nps: got %+v", n)
	}
}

func TestGetVersion(t *testing.T) {
	if GetVersion() == "" {
		t.Error("GetVersion() returned empty string")
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
	if GetPid() <= 0 {
		t.Error("GetPid() returned non-positive pid")
	}
}

func TestKillInvalidPid(t *testing.T) {
	// pid 0 / negative would signal the whole process group; must error.
	if err := Kill(0); err == nil {
		t.Error("Kill(0): expected error")
	}
	if err := Kill(-1); err == nil {
		t.Error("Kill(-1): expected error")
	}
}
