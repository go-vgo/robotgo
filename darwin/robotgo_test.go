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

import (
	"testing"

	"github.com/ebitengine/purego"
)

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
		{"right_shift alias", []interface{}{"right_shift"}, []string{"right_shift"}},
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
	if got := flagsFromMods([]string{"right_shift"}); got != kCGEventFlagMaskShift {
		t.Errorf("flagsFromMods(right_shift): got 0x%x", got)
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

// noPid is a pid no running process can have (above kern.maxproc).
const noPid = 0x7ffffff0

func TestWindowNotFound(t *testing.T) {
	if err := ActiveName("robotgo-no-such-app-0x7ff"); err != ErrNotFound {
		t.Errorf("ActiveName: got %v, want ErrNotFound", err)
	}
	if err := ActivePid(noPid); err != ErrNotFound {
		t.Errorf("ActivePid: got %v, want ErrNotFound", err)
	}
	if err := ActivePid(0); err != ErrNotFound {
		t.Errorf("ActivePid(0): got %v, want ErrNotFound", err)
	}
	if got := GetTitle(noPid); got != "" {
		t.Errorf("GetTitle: got %q, want empty", got)
	}
	if x, y, w, h := GetBounds(noPid); x != 0 || y != 0 || w != 0 || h != 0 {
		t.Errorf("GetBounds: got %d,%d,%d,%d, want zeros", x, y, w, h)
	}
	if err := withWindow(noPid, false, func(uintptr) error { return nil }); err != ErrNotFound {
		t.Errorf("withWindow: got %v, want ErrNotFound", err)
	}
	// pid_t overflow must not wrap to a real pid (1<<32+1 -> 1).
	if err := ActivePid(1<<32 + 1); err != ErrNotFound {
		t.Errorf("ActivePid(overflow): got %v, want ErrNotFound", err)
	}
	// Writes never fall back to the frontmost window for a missing pid.
	if err := minWindow(0, true); err != ErrNotFound {
		t.Errorf("minWindow(0): got %v, want ErrNotFound", err)
	}
	if err := closeWindowPid(0); err != ErrNotFound {
		t.Errorf("closeWindowPid(0): got %v, want ErrNotFound", err)
	}
	if got := GetTitle(0); got != "" {
		t.Errorf("GetTitle(0): got %q, want empty", got)
	}
	// Must not panic or touch any real window.
	if err := MinWindow(noPid); err == nil {
		t.Error("MinWindow(noPid): got nil, want error")
	}
	if err := MaxWindow(noPid, false); err == nil {
		t.Error("MaxWindow(noPid): got nil, want error")
	}
	if err := CloseWindow(noPid); err == nil {
		t.Error("CloseWindow(noPid): got nil, want error")
	}
	if err := MinWindow(0); err != ErrNotFound {
		t.Errorf("MinWindow(0): got %v, want ErrNotFound", err)
	}
	if err := CloseWindow(0); err != ErrNotFound {
		t.Errorf("CloseWindow(0): got %v, want ErrNotFound", err)
	}
}

func TestLoadAX(t *testing.T) {
	if !loadAX() {
		t.Fatal("AX API failed to load")
	}
	if cfBool(true) == 0 || cfBool(false) == 0 || cfBool(true) == cfBool(false) {
		t.Error("cfBool: kCFBooleanTrue/False not resolved")
	}
	if axSystem == 0 || axValueTypeID == 0 {
		t.Error("system-wide element / AXValue type id not resolved")
	}
	for i, s := range []uintptr{axFocusedApplication, axFrontmost, axFocusedWindow, axMainWindow,
		axWindows, axTitle, axPosition, axSize, axMinimized, axFullScreen, axCloseButton, axPress, axRaise} {
		if s == 0 {
			t.Errorf("AX attribute %d not created", i)
		}
	}
}

func TestBoolArgAndAXError(t *testing.T) {
	if !boolArg(nil, true) || boolArg([]interface{}{false}, true) || !boolArg([]interface{}{"x"}, true) {
		t.Error("boolArg: unexpected result")
	}
	if axError(kAXErrorSuccess, "op") != nil {
		t.Error("axError(success): want nil")
	}
	if err := axError(-25204, "close"); err == nil || err.Error() != "robotgo: close failed, AXError -25204" {
		t.Errorf("axError: got %v", err)
	}
}

func TestCFHelpers(t *testing.T) {
	for _, s := range []string{"", "AXTitle", "héllo 世界"} {
		ref := cfStringCreateWithCString(0, s, cfStringEncodingUTF8)
		if got := cfGoString(ref); got != s {
			t.Errorf("cfGoString(%q) = %q", s, got)
		}
		// A string is not a number.
		if _, ok := cfInt(ref); ok {
			t.Errorf("cfInt(string %q): want !ok", s)
		}
		cfRelease(ref)
	}
	if cfGoString(0) != "" || cfGoString(cfBooleanTrue) != "" {
		t.Error("cfGoString(non-string): want empty")
	}
	for _, p := range []struct {
		pid  int
		want bool
	}{{0, false}, {-1, false}, {1, true}, {1<<31 - 1, true}, {1 << 31, false}} {
		if validPid(p.pid) != p.want {
			t.Errorf("validPid(%d) != %v", p.pid, p.want)
		}
	}
}

func TestMatchOwner(t *testing.T) {
	owners := []owner{{1, "Xcode"}, {2, "Code"}, {3, "Visual Studio Code Helper"}, {4, "Safari"}}
	for _, c := range []struct {
		name string
		want int
	}{{"code", 2}, {"CODE", 2}, {"xc", 1}, {"saf", 4}, {"helper", 3}, {"", 0}, {"nope", 0}} {
		if got := matchOwner(owners, c.name); got != c.want {
			t.Errorf("matchOwner(%q) = %d, want %d", c.name, got, c.want)
		}
	}
}

// The window list is live and needs no permission.
func TestWindowOwnersLive(t *testing.T) {
	owners := windowOwners(kCGWindowListOptionAll)
	seen := map[int]bool{}
	for _, o := range owners {
		if o.pid <= 0 || seen[o.pid] {
			t.Errorf("windowOwners: bad or duplicate pid %d", o.pid)
		}
		seen[o.pid] = true
	}
	if len(owners) == 0 {
		t.Skip("no windows in this session")
	}
	if pid := frontmostPid(); pid < 0 {
		t.Errorf("frontmostPid: got %d", pid)
	}
}

func TestMultiClickZero(t *testing.T) {
	if err := MultiClick("left", 0); err != nil {
		t.Errorf("MultiClick(0): got %v, want nil", err)
	}
}

func TestFrameworksLoaded(t *testing.T) {
	// The system frameworks are always present on macOS; init must resolve
	// them. (Posting events may still require user-granted permissions.)
	if !loaded {
		t.Error("CoreGraphics/CoreFoundation frameworks failed to load")
	}
}

// Synthetic events must never suppress the user's own mouse/keyboard: the
// source handed to create has the local-events suppression interval cleared
// (the CoreGraphics default is 0.25s per posted event).
func TestWithSourceNoLocalSuppression(t *testing.T) {
	if !loaded {
		t.Skip("CoreGraphics not loaded")
	}
	cg, err := purego.Dlopen(
		"/System/Library/Frameworks/CoreGraphics.framework/CoreGraphics", purego.RTLD_NOW)
	if err != nil {
		t.Fatal(err)
	}
	var getInterval func(source uintptr) float64
	purego.RegisterLibFunc(&getInterval, cg, "CGEventSourceGetLocalEventsSuppressionInterval")

	// sanity: a fresh source carries the 0.25s default this fix removes
	src := cgEventSourceCreate(kCGEventSourceStateHIDSystemState)
	if src == 0 {
		t.Skip("CGEventSourceCreate returned nil")
	}
	if got := getInterval(src); got <= 0 {
		t.Errorf("default suppression interval = %v, want > 0", got)
	}
	cfRelease(src)

	called := false
	withSource(func(source uintptr) uintptr {
		called = true
		if got := getInterval(source); got != 0 {
			t.Errorf("withSource suppression interval = %v, want 0", got)
		}
		return 0
	})
	if !called {
		t.Error("withSource did not invoke create")
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

func TestScrollDirInvalid(t *testing.T) {
	// An invalid direction must error out before posting any event.
	if err := ScrollDir(1, "sideways"); err == nil {
		t.Error("ScrollDir(sideways): expected error")
	}
	if err := ScrollDir(1, 3); err == nil {
		t.Error("ScrollDir(3): expected error")
	}
}

func TestTypeErr(t *testing.T) {
	if err := typeErr(3, "hél"); err != nil {
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

func TestAlertLoad(t *testing.T) {
	if !loaded {
		t.Skip("frameworks not loaded")
	}
	if !loadAlert() {
		t.Fatal("CFUserNotificationDisplayAlert not resolved")
	}
	if cfStr("") != 0 {
		t.Error(`cfStr(""): want 0 so the button is omitted`)
	}
	s := cfStr("ok")
	if s == 0 {
		t.Fatal(`cfStr("ok"): got 0`)
	}
	cfRelease(s)
}
