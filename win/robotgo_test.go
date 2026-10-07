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

import (
	"fmt"
	"testing"
)

func TestKeyToVK(t *testing.T) {
	// Named keys must resolve.
	named := []string{
		"enter", "tab", "space", "backspace", "delete", "esc", "escape",
		"caps", "capslock",
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

	// Names from the Cgo keyNames table must map to the same virtual keys.
	parity := map[string]uint16{
		"pause_break": 0x13, // VK_PAUSE
		"right_shift": 0xA1, // VK_RSHIFT
		"num_equal":   0xBB, // VK_OEM_PLUS
		"numpad_0":    0x60, // VK_NUMPAD0
		"numpad_9":    0x69, // VK_NUMPAD9
		"numpad_lock": 0x90, // VK_NUMLOCK
	}
	for k, want := range parity {
		if vk, _, ok := keyToVK(k); !ok || vk != want {
			t.Errorf("keyToVK(%q) = (%#x, %v), want %#x", k, vk, ok, want)
		}
	}
}

func TestKeyToVKRejectsNonBMP(t *testing.T) {
	// These code points truncate to ASCII keys if cast directly to uint16.
	for _, key := range []string{"\U00010061", "\U00010041", "\U00010030", "\U00100061", "😀"} {
		if vk, mods, ok := keyToVK(key); ok || vk != 0 || mods != 0 {
			t.Errorf("keyToVK(%q) = (%#x, %d, %v), want (0, 0, false)", key, vk, mods, ok)
		}
	}
}

func TestAppendUniqueMod(t *testing.T) {
	for _, mod := range []string{"shift", "ctrl", "alt"} {
		variants := []string{mod, mod + "l", mod + "r"}
		if mod == "ctrl" {
			variants = append(variants, "control")
		}
		for _, variant := range variants {
			t.Run(variant, func(t *testing.T) {
				got := appendUniqueMod([]string{"cmd", variant}, mod)
				if len(got) != 2 || got[0] != "cmd" || got[1] != variant {
					t.Errorf("appendUniqueMod([cmd %s], %q) = %v, want unchanged modifiers", variant, mod, got)
				}
			})
		}
		t.Run(mod+"/missing", func(t *testing.T) {
			got := appendUniqueMod([]string{"cmd"}, mod)
			if len(got) != 2 || got[0] != "cmd" || got[1] != mod {
				t.Errorf("appendUniqueMod([cmd], %q) = %v, want [cmd %s]", mod, got, mod)
			}
			got = appendUniqueMod(nil, mod)
			if len(got) != 1 || got[0] != mod {
				t.Errorf("appendUniqueMod(nil, %q) = %v, want [%s]", mod, got, mod)
			}
		})
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

func TestExtractPid(t *testing.T) {
	tests := []struct {
		name string
		args []interface{}
		want int
	}{
		{"no args", []interface{}{}, 0},
		{"only modifiers", []interface{}{"ctrl", "shift"}, 0},
		{"pid first", []interface{}{1234}, 1234},
		{"pid after modifier", []interface{}{"ctrl", 4321}, 4321},
		{"first int wins", []interface{}{42, 99}, 42},
		{"mixed types", []interface{}{"alt", true, 7, "shift"}, 7},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := extractPid(tt.args); got != tt.want {
				t.Errorf("extractPid(%v): got %d, want %d", tt.args, got, tt.want)
			}
		})
	}
}

func TestKeyHwndNotPid(t *testing.T) {
	// With NotPid set, the pid value is treated as an HWND directly.
	old := NotPid
	NotPid = true
	defer func() { NotPid = old }()

	if got := keyHwnd(0x1234); uintptr(got) != 0x1234 {
		t.Errorf("keyHwnd(0x1234) with NotPid: got %#x, want %#x", uintptr(got), 0x1234)
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

func TestTypeRunesStopsOnFailure(t *testing.T) {
	ok := func(uint16) bool { return true }
	if n := typeRunes("héllo😀", ok); n != 6 {
		t.Errorf("typeRunes all ok: got %d, want 6", n)
	}
	calls := 0
	failThird := func(uint16) bool { calls++; return calls < 3 }
	if n := typeRunes("abcd", failThird); n != 2 {
		t.Errorf("typeRunes fail on 3rd: got %d, want 2", n)
	}
	// A surrogate pair counts only when both halves are sent.
	calls = 0
	failSecond := func(uint16) bool { calls++; return calls < 2 }
	if n := typeRunes("😀", failSecond); n != 0 {
		t.Errorf("typeRunes partial surrogate: got %d, want 0", n)
	}
}

func TestScrollDirInvalid(t *testing.T) {
	// An invalid direction must error out before sending any input.
	if err := ScrollDir(1, "sideways"); err == nil {
		t.Error("ScrollDir(sideways): expected error")
	}
	if err := ScrollDir(1, 3); err == nil {
		t.Error("ScrollDir(3): expected error")
	}
}

func TestTypeErr(t *testing.T) {
	if err := typeErr(2, "hé"); err != nil {
		t.Errorf("typeErr(full): got %v, want nil", err)
	}
	err := typeErr(1, "héllo")
	if err == nil || err.Error() != "robotgo: typed 1 of 5 characters" {
		t.Errorf("typeErr(partial): got %v", err)
	}
	if err := TypeStr(""); err != nil {
		t.Errorf("TypeStr(empty): got %v, want nil", err)
	}
}

// A failed press releases only the keys that were pressed so far; keys the
// user holds (never pressed by us) must not get a key up.
func TestPressKeysPartial(t *testing.T) {
	var log []string
	fail := uint16(0x42)
	send := func(vk uint16, up bool) error {
		if !up && vk == fail {
			return errSendInput
		}
		log = append(log, fmt.Sprintf("%x:%v", vk, up))
		return nil
	}
	vks := []uint16{0x10, 0x11, fail, 0x41}
	pressed, err := pressKeys(send, vks)
	if err != errSendInput {
		t.Fatalf("pressKeys: got %v, want errSendInput", err)
	}
	if len(pressed) != 2 || pressed[0] != 0x10 || pressed[1] != 0x11 {
		t.Fatalf("pressKeys prefix: got %x", pressed)
	}
	if err := releaseKeys(send, pressed); err != nil {
		t.Fatalf("releaseKeys: %v", err)
	}
	want := []string{"10:false", "11:false", "11:true", "10:true"}
	if fmt.Sprint(log) != fmt.Sprint(want) {
		t.Errorf("events: got %v, want %v", log, want)
	}

	// All presses succeed: the whole list is returned.
	log = nil
	pressed, err = pressKeys(send, []uint16{0x10, 0x41})
	if err != nil || len(pressed) != 2 {
		t.Errorf("pressKeys ok: got %x, %v", pressed, err)
	}
}

// In handle mode (Cgo NotPid / extra arg) pid is an HWND and must exist.
func TestWindowForHandle(t *testing.T) {
	const bogus = 0x7ffffff0
	if h := windowFor(0, true); h != 0 {
		t.Errorf("windowFor(0, handle): got %v, want 0", h)
	}
	if h := windowFor(bogus, true); h != 0 {
		t.Errorf("windowFor(bogus, handle): got %v, want 0", h)
	}
	if err := ActivePid(bogus, true); err != ErrNotFound {
		t.Errorf("ActivePid(bogus, handle): got %v, want ErrNotFound", err)
	}
	if err := ActivePid(0, false); err != ErrNotFound {
		t.Errorf("ActivePid(0): got %v, want ErrNotFound", err)
	}
	// A pid with no visible window must be reported, not silently ignored.
	if err := MinWindow(bogus); err != ErrNotFound {
		t.Errorf("MinWindow(bogus): got %v, want ErrNotFound", err)
	}
	if err := MaxWindow(bogus, false); err != ErrNotFound {
		t.Errorf("MaxWindow(bogus): got %v, want ErrNotFound", err)
	}
	if err := CloseWindow(bogus); err != ErrNotFound {
		t.Errorf("CloseWindow(bogus): got %v, want ErrNotFound", err)
	}
	// pid <= 0 (e.g. a failed FindIds) must not hit the foreground window.
	if err := CloseWindow(0); err != ErrNotFound {
		t.Errorf("CloseWindow(0): got %v, want ErrNotFound", err)
	}
	if err := MinWindow(0); err != ErrNotFound {
		t.Errorf("MinWindow(0): got %v, want ErrNotFound", err)
	}
	if err := MaxWindow(-1); err != ErrNotFound {
		t.Errorf("MaxWindow(-1): got %v, want ErrNotFound", err)
	}
	if x, y, w, h := GetBounds(bogus, true); x|y|w|h != 0 {
		t.Errorf("GetBounds(bogus, handle): got %d,%d,%d,%d", x, y, w, h)
	}
	if x, y, w, h := GetClient(bogus, true); x|y|w|h != 0 {
		t.Errorf("GetClient(bogus, handle): got %d,%d,%d,%d", x, y, w, h)
	}
}

func TestGetHWNDByPidInvalid(t *testing.T) {
	for _, pid := range []int{0, -1} {
		if got := GetHWNDByPid(pid); got != 0 {
			t.Errorf("GetHWNDByPid(%d): got %d, want 0", pid, got)
		}
	}
}
