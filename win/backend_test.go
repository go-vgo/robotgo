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
	"errors"
	"reflect"
	"testing"
	"unsafe"

	"github.com/go-vgo/robotgo/pub"
	"github.com/tailscale/win"
)

// fakeHID captures HID strokes instead of calling the Interception driver.
func fakeHID(t *testing.T) (mice *[]hidMouseStroke, keys *[]hidKeyStroke) {
	t.Helper()
	oldMouse, oldKey, oldOpen, oldScreen := hidSendMouse, hidSendKey, hidOpen, virtualScreen
	oldSleep := pub.KeySleep
	t.Cleanup(func() {
		hidSendMouse, hidSendKey, hidOpen, virtualScreen = oldMouse, oldKey, oldOpen, oldScreen
		pub.KeySleep = oldSleep
		hid.open = false
		if err := SetBackend(BackendSendInput); err != nil {
			t.Errorf("reset backend: %v", err)
		}
	})
	mice, keys = &[]hidMouseStroke{}, &[]hidKeyStroke{}
	hidSendMouse = func(s *hidMouseStroke) error { *mice = append(*mice, *s); return nil }
	hidSendKey = func(s *hidKeyStroke) error { *keys = append(*keys, *s); return nil }
	hidOpen = func() error { hid.open = true; return nil }
	virtualScreen = func() (int, int, int, int) { return -100, 0, 201, 101 }
	pub.KeySleep = 0
	if err := SetBackend(BackendHID); err != nil {
		t.Fatalf("SetBackend(HID): %v", err)
	}
	return mice, keys
}

func TestHIDStrokeLayout(t *testing.T) {
	if n := unsafe.Sizeof(hidMouseStroke{}); n != 20 {
		t.Errorf("mouse stroke size = %d, want 20", n)
	}
	if off := unsafe.Offsetof(hidMouseStroke{}.X); off != 8 {
		t.Errorf("mouse stroke X offset = %d, want 8", off)
	}
	if n := unsafe.Sizeof(hidKeyStroke{}); n != 8 {
		t.Errorf("key stroke size = %d, want 8", n)
	}
	if n := unsafe.Sizeof(hidRawKey{}); n != 12 {
		t.Errorf("KEYBOARD_INPUT_DATA size = %d, want 12", n)
	}
	if n := unsafe.Sizeof(hidRawMouse{}); n != 24 {
		t.Errorf("MOUSE_INPUT_DATA size = %d, want 24", n)
	}
	if off := unsafe.Offsetof(hidRawMouse{}.LastX); off != 12 {
		t.Errorf("MOUSE_INPUT_DATA LastX offset = %d, want 12", off)
	}
}

func TestNormalize(t *testing.T) {
	for _, c := range []struct {
		v, origin, size int
		want            int32
	}{
		{-100, -100, 201, 0},
		{100, -100, 201, 65535},
		{0, -100, 201, 32767},
		{5, 0, 1, 0},
		// (v-origin)*65535 overflows a 32-bit int.
		{40000, 0, 40001, 65535},
		{0, -20000, 40001, 32767},
	} {
		if got := normalize(c.v, c.origin, c.size); got != c.want {
			t.Errorf("normalize(%d,%d,%d) = %d, want %d", c.v, c.origin, c.size, got, c.want)
		}
	}
}

func TestHIDMouse(t *testing.T) {
	mice, _ := fakeHID(t)
	if err := InputMove(100, 50); err != nil {
		t.Fatal(err)
	}
	if err := InputClick("right", false); err != nil {
		t.Fatal(err)
	}
	if err := InputScroll(1, -2); err != nil {
		t.Fatal(err)
	}
	want := []hidMouseStroke{
		{Flags: hidMoveAbsolute | hidMoveVirtualDesktop, X: 65535, Y: 32767},
		{State: hidMouseRightDown},
		{State: hidMouseRightUp},
		{State: hidMouseWheel, Rolling: -240},
		{State: hidMouseHWheel, Rolling: -120},
	}
	if !reflect.DeepEqual(*mice, want) {
		t.Errorf("strokes = %+v\nwant %+v", *mice, want)
	}
	if err := InputToggle("bogus", false); err == nil {
		t.Error("unknown button should fail")
	}
}

func TestHIDKeyboard(t *testing.T) {
	_, keys := fakeHID(t)
	if err := KeyTap("delete", "ctrl"); err != nil {
		t.Fatal(err)
	}
	ctrl, del := scanCode(win.VK_CONTROL), scanCode(win.VK_DELETE)
	want := []hidKeyStroke{
		{Code: ctrl, State: hidKeyDown},
		{Code: del, State: hidKeyDown | hidKeyE0},
		{Code: del, State: hidKeyUp | hidKeyE0},
		{Code: ctrl, State: hidKeyUp},
	}
	if !reflect.DeepEqual(*keys, want) {
		t.Errorf("strokes = %+v\nwant %+v", *keys, want)
	}

	*keys = nil
	if err := InputType("A"); err != nil {
		t.Fatal(err)
	}
	shift, a := scanCode(win.VK_SHIFT), scanCode('A')
	want = []hidKeyStroke{
		{Code: shift, State: hidKeyDown}, {Code: a, State: hidKeyDown},
		{Code: a, State: hidKeyUp}, {Code: shift, State: hidKeyUp},
	}
	if !reflect.DeepEqual(*keys, want) {
		t.Errorf("type strokes = %+v\nwant %+v", *keys, want)
	}
}

func TestHIDWheelSplit(t *testing.T) {
	mice, _ := fakeHID(t)
	if err := InputScroll(-300, 600); err != nil {
		t.Fatal(err)
	}
	full := int16(hidMaxWheelStep * wheelDelta)
	want := []hidMouseStroke{
		{State: hidMouseWheel, Rolling: full},
		{State: hidMouseWheel, Rolling: full},
		{State: hidMouseWheel, Rolling: int16((600 - 2*hidMaxWheelStep) * wheelDelta)},
		{State: hidMouseHWheel, Rolling: full},
		{State: hidMouseHWheel, Rolling: int16((300 - hidMaxWheelStep) * wheelDelta)},
	}
	if !reflect.DeepEqual(*mice, want) {
		t.Errorf("strokes = %+v\nwant %+v", *mice, want)
	}
}

func TestHIDSpecialKeys(t *testing.T) {
	for _, c := range []struct {
		vk   uint16
		up   bool
		want []hidKeyStroke
	}{
		{win.VK_PAUSE, false, []hidKeyStroke{{Code: 0x1D, State: hidKeyE1}, {Code: 0x45}}},
		{win.VK_PAUSE, true, []hidKeyStroke{{Code: 0x1D, State: hidKeyUp | hidKeyE1}, {Code: 0x45, State: hidKeyUp}}},
		{win.VK_SNAPSHOT, false, []hidKeyStroke{{Code: 0x37, State: hidKeyE0}}},
		{win.VK_NUMLOCK, true, []hidKeyStroke{{Code: 0x45, State: hidKeyUp}}},
	} {
		got, err := hidKeyStrokes(c.vk, c.up)
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(got, c.want) {
			t.Errorf("vk %#x up=%v: %+v, want %+v", c.vk, c.up, got, c.want)
		}
	}

	_, keys := fakeHID(t)
	if err := KeyTap("pause"); err != nil {
		t.Skipf("pause not mapped: %v", err)
	}
	if len(*keys) != 4 {
		t.Errorf("pause tap sent %d strokes, want 4", len(*keys))
	}
}

func TestHIDSendError(t *testing.T) {
	fakeHID(t)
	sendErr := errors.New("send failed")
	hidSendKey = func(*hidKeyStroke) error { return sendErr }
	if err := KeyTap("a"); !errors.Is(err, sendErr) {
		t.Errorf("KeyTap error = %v, want %v", err, sendErr)
	}
}

func TestSetBackendHIDFailure(t *testing.T) {
	oldOpen := hidOpen
	t.Cleanup(func() { hidOpen = oldOpen })
	hidOpen = func() error { return ErrDriverNotInstalled }
	if err := SetBackend(BackendHID); !errors.Is(err, ErrDriverNotInstalled) {
		t.Fatalf("SetBackend(HID) = %v", err)
	}
	if GetBackend() != BackendSendInput {
		t.Error("failed SetBackend must keep the previous backend")
	}
	if err := SetBackend(Backend(99)); err == nil {
		t.Error("unknown backend should fail")
	}
}

func TestMessageBackend(t *testing.T) {
	oldPost, oldThread, oldFg, oldSleep := postMessageW, getWindowThreadProcessID, getForegroundWindow, pub.KeySleep
	oldDbl := classDblClks
	classDblClks = func(win.HWND) bool { return true }
	t.Cleanup(func() {
		classDblClks = oldDbl
		postMessageW, getWindowThreadProcessID, getForegroundWindow, pub.KeySleep = oldPost, oldThread, oldFg, oldSleep
		SetMessageTarget(WindowInfo{})
		if err := SetBackend(BackendSendInput); err != nil {
			t.Errorf("reset backend: %v", err)
		}
	})
	type msg struct{ hwnd, kind, wp uintptr }
	var got []msg
	postMessageW = func(args ...uintptr) (uintptr, uintptr, error) {
		got = append(got, msg{args[0], args[1], args[2]})
		return 1, 0, nil
	}
	getWindowThreadProcessID = func(win.HWND, *uint32) uint32 { return 0 }
	pub.KeySleep = 0

	if err := SetBackend(BackendMessage); err != nil {
		t.Fatal(err)
	}
	// No target and no foreground window.
	getForegroundWindow = func() win.HWND { return 0 }
	if err := InputType("x"); !errors.Is(err, ErrNotFound) {
		t.Errorf("no window: %v", err)
	}

	getForegroundWindow = func() win.HWND { return 7 }
	if err := KeyTap("enter"); err != nil {
		t.Fatal(err)
	}
	if err := InputType("hi"); err != nil {
		t.Fatal(err)
	}
	want := []msg{{7, wmKeyDown, win.VK_RETURN}, {7, wmKeyUp, win.VK_RETURN}, {7, wmChar, 'h'}, {7, wmChar, 'i'}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("keys = %v, want %v", got, want)
	}

	got = nil
	SetMessageTarget(WindowInfo{ID: 9})
	if err := InputMove(5, 6); err != nil {
		t.Fatal(err)
	}
	if err := InputToggle("left", false); err != nil {
		t.Fatal(err)
	}
	if err := InputMove(8, 9); err != nil {
		t.Fatal(err)
	}
	if err := InputToggle("left", true); err != nil {
		t.Fatal(err)
	}
	if err := InputClick("left", true); err != nil {
		t.Fatal(err)
	}
	if err := InputScroll(0, 1); err != nil {
		t.Fatal(err)
	}
	if x, y := pointer(); x != 8 || y != 9 {
		t.Errorf("virtual pointer = %d,%d", x, y)
	}
	want = []msg{
		{9, wmMouseMove, 0},
		{9, wmLButtonDown, mkLButton},
		{9, wmMouseMove, mkLButton}, // drag while held
		{9, wmLButtonUp, 0},
		{9, wmLButtonDown, mkLButton},
		{9, wmLButtonUp, 0},
		{9, wmLButtonDblClk, mkLButton},
		{9, wmLButtonUp, 0},
		{9, wmMouseWheel, wheelWParam(wheelDelta)},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("mouse = %v\nwant %v", got, want)
	}
}

func TestInputKeyToggleOrder(t *testing.T) {
	_, keys := fakeHID(t)
	if err := InputKeyToggle("a", false, "shift", "alt"); err != nil {
		t.Fatal(err)
	}
	if err := InputKeyToggle("a", true, "shift", "alt"); err != nil {
		t.Fatal(err)
	}
	sh, alt, a := scanCode(win.VK_SHIFT), scanCode(win.VK_MENU), scanCode('A')
	want := []hidKeyStroke{
		{Code: sh}, {Code: alt}, {Code: a},
		{Code: a, State: hidKeyUp}, {Code: alt, State: hidKeyUp}, {Code: sh, State: hidKeyUp},
	}
	if !reflect.DeepEqual(*keys, want) {
		t.Errorf("strokes = %+v\nwant %+v", *keys, want)
	}
	if err := InputKeyToggle("nope", false); err == nil {
		t.Error("unknown key should fail")
	}
	if err := InputKeyToggle("a", false, "nope"); err == nil {
		t.Error("unknown modifier should fail")
	}
}

// InputKeyToggle must hold the shift a key implies ("A"), and not press it
// twice when it is also given explicitly, like KeyTap and KeyToggle.
func TestInputKeyToggleImpliedShift(t *testing.T) {
	if _, mods, ok := keyToVK("A"); !ok || mods&1 == 0 {
		t.Skip("layout does not need shift for 'A'")
	}
	_, keys := fakeHID(t)
	if err := InputKeyToggle("A", false); err != nil {
		t.Fatal(err)
	}
	if err := InputKeyToggle("A", true, "shift"); err != nil {
		t.Fatal(err)
	}
	sh, a := scanCode(win.VK_SHIFT), scanCode('A')
	want := []hidKeyStroke{
		{Code: sh}, {Code: a},
		{Code: a, State: hidKeyUp}, {Code: sh, State: hidKeyUp},
	}
	if !reflect.DeepEqual(*keys, want) {
		t.Errorf("strokes = %+v\nwant %+v", *keys, want)
	}
}
