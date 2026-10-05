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

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/jezek/xgb/xproto"
)

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
		{"caps", xkCapsLock, true},
		{"capslock", xkCapsLock, true},
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
		{"audio_random", 0x1008ff99, true},
		{"right_shift", xkShiftR, true},
		{"numpad_7", xkKP0 + 7, true},
		{"numpad_lock", xkNumLock, true},
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
		{"up dir in slice", []interface{}{[]string{"up", "ctrl"}}, []string{"ctrl"}, false, true},
		{"last dir wins", []interface{}{"ctrl", "up", []string{"down"}}, []string{"ctrl"}, true, true},
		{"up dir in nested iface", []interface{}{[]interface{}{"up", "ctrl"}}, []string{"ctrl"}, false, true},
		{"nested without dir keeps outer", []interface{}{"up", []interface{}{"ctrl"}}, []string{"ctrl"}, false, true},
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

// germanConn builds a conn around the core keyboard mapping an XKB server
// exposes for a German (de) layout, as read with xmodmap -pke under Xvfb:
// columns 0/1 are the plain/Shift levels, 2/3 repeat them, 4/5 hold the AltGr
// levels. Keycode 25 is left empty so a scratch keycode exists.
func germanConn() *conn {
	const per = 7
	rows := map[xproto.Keycode][]xproto.Keysym{
		16:  {'7', '/', '7', '/', '{', 0xbe, '{'},  // 7 slash ... braceleft seveneighths
		24:  {'q', 'Q', 'q', 'Q', '@', 0x7d9, '@'}, // q Q ... at Greek_OMEGA
		50:  {xkShiftL, 0, xkShiftL},
		52:  {'y', 'Y', 'y', 'Y', 0xbb, 0xab, 0xbb}, // y Y (z/y swap) ... guillemotright
		108: {xkISOLevel3Shift, 0, xkISOLevel3Shift},
	}
	c := &conn{minKeycode: 16, keysymsPerKeycode: per}
	for kc := xproto.Keycode(16); kc <= 108; kc++ {
		row := make([]xproto.Keysym, per)
		copy(row, rows[kc])
		c.keysyms = append(c.keysyms, row...)
	}
	c.findModKeycodes()
	c.scratch, c.scratchOK = c.findScratchKeycode()
	return c
}

// #640: on a German layout '@' lives on the AltGr level of q and '/' on the
// Shift level of 7. The lookup must report those levels so pressKeysym holds
// AltGr/Shift, instead of returning q+Shift for '@' (which types 'Q').
func TestKeysymToKeycodeGermanLayout(t *testing.T) {
	c := germanConn()
	if !c.scratchOK || c.shiftKeycode != 50 || c.level3Keycode != 108 {
		t.Fatalf("setup: scratchOK=%v shift=%d level3=%d", c.scratchOK, c.shiftKeycode, c.level3Keycode)
	}
	tests := []struct {
		ks   uint32
		kc   xproto.Keycode
		mods uint8
		ok   bool
	}{
		{'7', 16, 0, true},
		{'/', 16, modShift, true},
		{'{', 16, modLevel3, true},
		{'q', 24, 0, true},
		{'Q', 24, modShift, true},
		{'@', 24, modLevel3, true},
		{0x7d9, 24, modLevel3 | modShift, true}, // Greek_OMEGA (legacy keysym)
		{'y', 52, 0, true},
		{'z', 0, 0, false}, // not on these rows: scratch keycode
	}
	for _, tt := range tests {
		kc, mods, ok := c.keysymToKeycode(tt.ks)
		if kc != tt.kc || mods != tt.mods || ok != tt.ok {
			t.Errorf("keysymToKeycode(%#x) = (%d, mods=%#b, ok=%v), want (%d, %#b, %v)",
				tt.ks, kc, mods, ok, tt.kc, tt.mods, tt.ok)
		}
	}
	if got := c.modKeycodesFor(modLevel3 | modShift); len(got) != 2 || got[0] != 108 || got[1] != 50 {
		t.Errorf("modKeycodesFor(level3|shift) = %v, want [108 50]", got)
	}
}

// Layout with AltGr symbols but no level-3 key: '@' must not be typed as
// its bare base key 'q'; it falls back to scratch keycode. With
// scratch disabled that path returns ErrNotSupported before any event is
// sent (c.c is nil, so sending would panic).
func TestPressKeysymNoLevel3Key(t *testing.T) {
	c := germanConn()
	c.level3Keycode = 0
	c.scratchOK = false

	if c.canReach(modLevel3) || c.canReach(modLevel3|modShift) {
		t.Error("canReach(level3) = true without level-3 key")
	}
	if !c.canReach(0) || !c.canReach(modShift) {
		t.Error("canReach(plain/shift) = false")
	}
	if err := c.pressKeysym('@'); err != ErrNotSupported {
		t.Errorf("pressKeysym('@') = %v, want ErrNotSupported", err)
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

func TestProcInfo(t *testing.T) {
	pid := os.Getpid()
	name, path, got := procInfo(pid)
	if got != pid || path == "" {
		t.Fatalf("procInfo(self) = %q, %q, %d", name, path, got)
	}
	// The name comes from the exe basename, not the 15-byte comm.
	if want := filepath.Base(path); name != want {
		t.Errorf("procInfo name = %q, want %q", name, want)
	}
	if name, path, got := procInfo(-1); name != "" || path != "" || got != -1 {
		t.Errorf("procInfo(-1) = %q, %q, %d", name, path, got)
	}
}

func TestParseXftDPI(t *testing.T) {
	tests := []struct {
		db   string
		want float64
		ok   bool
	}{
		{"Xft.dpi:\t144\nXft.antialias:\t1\n", 144, true},
		{"Xcursor.size: 24\nXft.dpi: 192.5", 192.5, true},
		{"Xft.antialias: 1\n", 0, false},
		{"Xft.dpi: nope\n", 0, false},
		{"", 0, false},
	}
	for _, tt := range tests {
		got, ok := parseXftDPI(tt.db)
		if got != tt.want || ok != tt.ok {
			t.Errorf("parseXftDPI(%q) = (%v, %v), want (%v, %v)", tt.db, got, ok, tt.want, tt.ok)
		}
	}
}
