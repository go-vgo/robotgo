//go:build windows
// +build windows

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

package win

import "testing"

func TestKeyToVK(t *testing.T) {
	// Named keys must resolve.
	named := []string{
		"enter", "tab", "space", "backspace", "delete", "esc", "escape",
		"up", "down", "left", "right", "home", "end",
		"shift", "ctrl", "alt", "f1", "f12",
	}
	for _, k := range named {
		if _, _, ok := keyToVK(k); !ok {
			t.Errorf("keyToVK(%q): expected resolvable named key", k)
		}
	}

	// Single alphanumerics resolve via VkKeyScan.
	for _, k := range []string{"a", "z", "0", "9"} {
		if _, _, ok := keyToVK(k); !ok {
			t.Errorf("keyToVK(%q): expected resolvable char", k)
		}
	}

	if _, _, ok := keyToVK("nonexistent_key"); ok {
		t.Error("keyToVK(nonexistent_key): expected not resolvable")
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

func TestMouseButtonFlags(t *testing.T) {
	tests := []string{"left", "right", "center", "middle", "", "unknown"}
	for _, btn := range tests {
		down, up := mouseButtonFlags(btn)
		if down == 0 || up == 0 {
			t.Errorf("mouseButtonFlags(%q): got zero flags down=%d up=%d", btn, down, up)
		}
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
