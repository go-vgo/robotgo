//go:build darwin
// +build darwin

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

package darwin

import "testing"

func TestKeyToCode(t *testing.T) {
	named := []string{
		"enter", "tab", "space", "backspace", "delete", "esc", "escape",
		"caps", "capslock", "fn",
		"up", "down", "left", "right", "home", "end",
		"shift", "ctrl", "alt", "cmd", "f1", "f12",
	}
	for _, k := range named {
		if _, _, ok := keyToCode(k); !ok {
			t.Errorf("keyToCode(%q): expected resolvable named key", k)
		}
	}

	// Single alphanumerics resolve.
	for _, k := range []string{"a", "z", "0", "9"} {
		if _, _, ok := keyToCode(k); !ok {
			t.Errorf("keyToCode(%q): expected resolvable char", k)
		}
	}

	// Uppercase letters resolve and imply the SHIFT flag.
	if _, flags, ok := keyToCode("A"); !ok || flags&kCGEventFlagMaskShift == 0 {
		t.Errorf("keyToCode(A): expected resolvable with shift flag, got ok=%v flags=0x%x", ok, flags)
	}

	if _, _, ok := keyToCode("nonexistent_key"); ok {
		t.Error("keyToCode(nonexistent_key): expected not resolvable")
	}

	// Names from the Cgo keyNames table must map to the same kVK codes.
	parity := map[string]uint16{
		"print": 105, "printscreen": 105, "right_shift": 60,
		"numpad_0": 82, "numpad_9": 92, "numpad_lock": 71,
	}
	for k, want := range parity {
		if code, _, ok := keyToCode(k); !ok || code != want {
			t.Errorf("keyToCode(%q) = (%d, %v), want %d", k, code, ok, want)
		}
	}
}

func TestMediaKeyData(t *testing.T) {
	// Same NX_KEYTYPE_* codes the Cgo backend encodes as 1000+code.
	want := map[string]int{
		"audio_vol_up": 0, "audio_vol_down": 1, "audio_mute": 7,
		"audio_play": 16, "audio_pause": 16, "audio_next": 17, "audio_prev": 18,
		"lights_mon_up": 2, "lights_mon_down": 3,
		"lights_kbd_up": 21, "lights_kbd_down": 22, "lights_kbd_toggle": 23,
	}
	for k, code := range want {
		if got, ok := mediaCodes[k]; !ok || got != code {
			t.Errorf("mediaCodes[%q] = (%d, %v), want %d", k, got, ok, code)
		}
	}

	flags, data1 := mediaKeyData(7, true)
	if flags != 0xa00 || data1 != 0x70a00 {
		t.Errorf("mediaKeyData(7, down) = (%#x, %#x), want (0xa00, 0x70a00)", flags, data1)
	}
	flags, data1 = mediaKeyData(16, false)
	if flags != 0xb00 || data1 != 0x100b00 {
		t.Errorf("mediaKeyData(16, up) = (%#x, %#x), want (0xb00, 0x100b00)", flags, data1)
	}
}

func TestWithMediaEvent(t *testing.T) {
	if !loaded {
		t.Skip("CoreGraphics not loaded")
	}
	// Build the NSEvent without posting it (posting would change the volume).
	for _, down := range []bool{true, false} {
		var cg uintptr
		if err := withMediaEvent(mediaCodes["audio_mute"], down, func(ev uintptr) { cg = ev }); err != nil {
			t.Fatalf("withMediaEvent(down=%v): %v", down, err)
		}
		if cg == 0 {
			t.Fatalf("withMediaEvent(down=%v): got nil CGEvent", down)
		}
	}
}

func TestModKeyCodes(t *testing.T) {
	// Every modifier accepted by extractModifiers must resolve to a keycode
	// so upModKeys (the upKeyArr equivalent) can key it up after a tap.
	mods := []string{"cmd", "command", "shift", "ctrl", "control", "alt", "option"}
	codes := modKeyCodes(mods)
	if len(codes) != len(mods) {
		t.Errorf("modKeyCodes(%v): got %d codes, want %d", mods, len(codes), len(mods))
	}

	// Unresolvable names are skipped, resolvable ones kept in order.
	codes = modKeyCodes([]string{"cmd", "not_a_modifier", "shift"})
	if len(codes) != 2 || codes[0] != namedCodes["cmd"] || codes[1] != namedCodes["shift"] {
		t.Errorf("modKeyCodes(cmd,not_a_modifier,shift): got %v", codes)
	}

	if got := modKeyCodes(nil); len(got) != 0 {
		t.Errorf("modKeyCodes(nil): got %v, want empty", got)
	}
}

func TestKeyTapUnknownKey(t *testing.T) {
	// An unknown key must error out before posting any event.
	if err := KeyTap("nonexistent_key", "cmd"); err == nil {
		t.Error("KeyTap(nonexistent_key): expected error")
	}
	if err := KeyToggle("nonexistent_key", "up"); err == nil {
		t.Error("KeyToggle(nonexistent_key, up): expected error")
	}
}

func TestExtractModifiers(t *testing.T) {
	tests := []struct {
		name string
		args []interface{}
		want []string
	}{
		{"no mods", []interface{}{}, nil},
		{"cmd", []interface{}{"cmd"}, []string{"cmd"}},
		{"cmd+shift", []interface{}{"cmd", "shift"}, []string{"cmd", "shift"}},
		{"mixed types", []interface{}{"ctrl", 42, true, "alt"}, []string{"ctrl", "alt"}},
		{"non-modifier string", []interface{}{"hello"}, nil},
		{"[]string slice", []interface{}{[]string{"cmd", "shift"}}, []string{"cmd", "shift"}},
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

func TestFlagsFromMods(t *testing.T) {
	if got := flagsFromMods([]string{"cmd", "shift"}); got != kCGEventFlagMaskCommand|kCGEventFlagMaskShift {
		t.Errorf("flagsFromMods(cmd,shift): got 0x%x", got)
	}
	if got := flagsFromMods(nil); got != 0 {
		t.Errorf("flagsFromMods(nil): got 0x%x, want 0", got)
	}
}

func TestExtractPid(t *testing.T) {
	tests := []struct {
		name string
		args []interface{}
		want int
	}{
		{"none", []interface{}{"cmd", "shift"}, 0},
		{"pid first", []interface{}{4321, "cmd"}, 4321},
		{"pid after mods", []interface{}{"cmd", 99}, 99},
		{"first int wins", []interface{}{7, 8}, 7},
		{"empty", []interface{}{}, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := extractPid(tt.args); got != tt.want {
				t.Errorf("extractPid(%v): got %d, want %d", tt.args, got, tt.want)
			}
		})
	}
}

func TestMouseButton(t *testing.T) {
	for _, btn := range []string{"left", "right", "center", "middle", "", "unknown"} {
		down, up, dragged, _ := mouseButton(btn)
		if down == 0 || up == 0 || dragged == 0 {
			t.Errorf("mouseButton(%q): got zero event type down=%d up=%d dragged=%d", btn, down, up, dragged)
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
	if got := CmdCtrl(); got != "cmd" {
		t.Errorf("CmdCtrl(): got %q, want %q", got, "cmd")
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

func TestWindowUnsupported(t *testing.T) {
	if err := ActiveName("nope"); err != ErrNotSupported {
		t.Errorf("ActiveName: got %v, want ErrNotSupported", err)
	}
	if GetTitle() != "" {
		t.Error("GetTitle: expected empty string")
	}
}

func TestFrameworksLoaded(t *testing.T) {
	// The system frameworks are always present on macOS; init must resolve
	// them. (Posting events may still require user-granted permissions.)
	if !loaded {
		t.Error("CoreGraphics/CoreFoundation frameworks failed to load")
	}
}

func TestScreenSize(t *testing.T) {
	// Headless/CI may report 0; just ensure it does not panic and is sane.
	w, h := GetScreenSize()
	if w < 0 || h < 0 {
		t.Errorf("GetScreenSize: got negative %dx%d", w, h)
	}
}

func TestGetPixelColor(t *testing.T) {
	// Headless / unpermissioned environments return the "000000" fallback;
	// just ensure a valid 6-char hex string comes back without panicking.
	c := GetPixelColor(1, 1)
	if len(c) != 6 {
		t.Errorf("GetPixelColor: got %q, want 6 hex chars", c)
	}
}

func TestMainDisplayID(t *testing.T) {
	if id := MainDisplayID(); id < 0 {
		t.Errorf("MainDisplayID: got negative id %d", id)
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

func TestKeyToCodePunctuationAndSpecial(t *testing.T) {
	tests := []struct {
		key   string
		code  uint16
		flags uint64
	}{
		{"-", 27, 0},
		{"=", 24, 0},
		{"/", 44, 0},
		{"`", 50, 0},
		{"!", 18, kCGEventFlagMaskShift}, // shift+1
		{"{", 33, kCGEventFlagMaskShift}, // shift+[
		{"~", 50, kCGEventFlagMaskShift}, // shift+`
		{"num0", 82, 0},
		{"num_enter", 76, 0},
	}
	for _, tt := range tests {
		code, flags, ok := keyToCode(tt.key)
		if !ok || code != tt.code || flags != tt.flags {
			t.Errorf("keyToCode(%q): got (%d,%#x,%v), want (%d,%#x,true)", tt.key, code, flags, ok, tt.code, tt.flags)
		}
	}
}

func TestScaleAndRect(t *testing.T) {
	if !loaded {
		t.Skip("CoreGraphics not loaded")
	}
	if f := ScaleF(); f < 1 || f > 4 {
		t.Errorf("ScaleF: got %v", f)
	}
	w, h := GetScreenSize()
	sw, sh := GetScaleSize()
	if sw < w || sh < h {
		t.Errorf("GetScaleSize (%d,%d) must be >= GetScreenSize (%d,%d)", sw, sh, w, h)
	}
	if r := GetScreenRect(); r.W != w || r.H != h {
		t.Errorf("GetScreenRect: got %+v, want %dx%d", r, w, h)
	}
	// Out-of-range index falls back to the main display, no panic.
	if r := GetScreenRect(99); r.W != w {
		t.Errorf("GetScreenRect(99): got %+v", r)
	}
}

func TestActiveAppAndAccess(t *testing.T) {
	CheckAccess(false)
	name, id, pid := GetActiveApp()
	if pid <= 0 {
		t.Skip("no frontmost app in this session")
	}
	// Some apps (e.g. bare executables) have no bundle id, so id is not checked.
	if name == "" {
		t.Errorf("GetActiveApp() = %q, %q, %d; want app name", name, id, pid)
	}
}
