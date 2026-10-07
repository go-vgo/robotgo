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

type postedMsg struct{ hwnd, msg, wp, lp uintptr }

// capturePosts records PostMessageW calls; fail makes the n-th post fail.
func capturePosts(t *testing.T, fail int) *[]postedMsg {
	t.Helper()
	oldPost, oldThread, oldSleep := postMessageW, getWindowThreadProcessID, pub.KeySleep
	t.Cleanup(func() { postMessageW, getWindowThreadProcessID, pub.KeySleep = oldPost, oldThread, oldSleep })
	getWindowThreadProcessID = func(win.HWND, *uint32) uint32 { return 0 }
	pub.KeySleep = 0
	got := &[]postedMsg{}
	postMessageW = func(a ...uintptr) (uintptr, uintptr, error) {
		*got = append(*got, postedMsg{a[0], a[1], a[2], a[3]})
		if len(*got) == fail {
			return 0, 0, errors.New("post failed")
		}
		return 1, 0, nil
	}
	return got
}

func msgs(p []postedMsg) (out [][2]uintptr) {
	for _, m := range p {
		out = append(out, [2]uintptr{m.msg, m.wp})
	}
	return out
}

func TestPostAltUsesSysKey(t *testing.T) {
	got := capturePosts(t, 0)
	if err := PostKeyTap(WindowInfo{ID: 7}, "f4", "alt"); err != nil {
		t.Fatal(err)
	}
	want := [][2]uintptr{
		{wmSysKeyDown, win.VK_MENU},
		{wmSysKeyDown, win.VK_F4},
		{wmSysKeyDown + 1, win.VK_F4},
		{wmKeyUp, win.VK_MENU}, // Alt's own release is a plain key-up
	}
	if !reflect.DeepEqual(msgs(*got), want) {
		t.Fatalf("posted %v, want %v", msgs(*got), want)
	}
	if (*got)[1].lp&(1<<29) == 0 {
		t.Errorf("sys key lParam %#x missing context bit", (*got)[1].lp)
	}

	// Ctrl+Alt (AltGr) is not a system key once Ctrl is held.
	*got = nil
	if err := PostKeyToggle(WindowInfo{ID: 7}, "e", "down", "alt", "ctrl"); err != nil {
		t.Fatal(err)
	}
	want = [][2]uintptr{{wmSysKeyDown, win.VK_MENU}, {wmKeyDown, win.VK_CONTROL}, {wmKeyDown, 'E'}}
	if !reflect.DeepEqual(msgs(*got), want) {
		t.Errorf("AltGr posted %v, want %v", msgs(*got), want)
	}
}

func TestPostKeyGenericVKAndF10(t *testing.T) {
	got := capturePosts(t, 0)
	if err := PostKeyTap(WindowInfo{ID: 7}, "f10", "ctrll"); err != nil {
		t.Fatal(err)
	}
	want := [][2]uintptr{
		{wmKeyDown, win.VK_CONTROL}, // generic VK in wParam
		{wmSysKeyDown, win.VK_F10},  // F10 is a system key
		{wmSysKeyDown + 1, win.VK_F10},
		{wmKeyUp, win.VK_CONTROL},
	}
	if !reflect.DeepEqual(msgs(*got), want) {
		t.Fatalf("posted %v, want %v", msgs(*got), want)
	}
	if (*got)[1].lp&(1<<29) != 0 {
		t.Errorf("F10 without Alt must not set the context bit: %#x", (*got)[1].lp)
	}
}

func TestKeyLParamSpecialScanCodes(t *testing.T) {
	if lp := keyLParam(win.VK_SNAPSHOT, false); lp != 1|0x37<<16|1<<24 {
		t.Errorf("VK_SNAPSHOT lParam %#x, want E0 37", lp)
	}
	if keyLParam(win.VK_PAUSE, false)&(1<<24) != 0 {
		t.Error("VK_PAUSE must not set the extended bit (E0 45 is NumLock)")
	}
}

func TestPostWheelChunks(t *testing.T) {
	got := capturePosts(t, 0)
	if err := postWheel(7, wmMouseWheel, 300, 0); err != nil {
		t.Fatal(err)
	}
	total := 0
	for _, m := range *got {
		d := int(int16(m.wp >> 16))
		if d <= 0 {
			t.Fatalf("wheel delta %d overflowed int16", d)
		}
		total += d
	}
	if total != 300*wheelDelta {
		t.Errorf("total delta %d, want %d", total, 300*wheelDelta)
	}
}

func TestKeyLParamExtendedBit(t *testing.T) {
	if keyLParam(win.VK_DELETE, false)&(1<<24) == 0 {
		t.Error("VK_DELETE must set the extended bit")
	}
	if keyLParam('A', false)&(1<<24) != 0 {
		t.Error("'A' must not set the extended bit")
	}
}

func TestPidKeyTapPostsAndReportsErrors(t *testing.T) {
	oldNotPid := NotPid
	t.Cleanup(func() { NotPid = oldNotPid })
	NotPid = true

	got := capturePosts(t, 0)
	if err := KeyTap("delete", 0x55, "ctrl"); err != nil {
		t.Fatal(err)
	}
	want := [][2]uintptr{
		{wmKeyDown, win.VK_CONTROL}, {wmKeyDown, win.VK_DELETE},
		{wmKeyUp, win.VK_DELETE}, {wmKeyUp, win.VK_CONTROL},
	}
	if !reflect.DeepEqual(msgs(*got), want) {
		t.Fatalf("posted %v, want %v", msgs(*got), want)
	}
	for _, m := range *got {
		if m.hwnd != 0x55 {
			t.Errorf("posted to %#x, want 0x55", m.hwnd)
		}
	}
	// Key-down must not be posted as key-up (the old flags==0 C bug).
	if (*got)[1].lp&0xC0000000 != 0 || (*got)[2].lp&0xC0000000 != 0xC0000000 {
		t.Errorf("transition bits wrong: down %#x up %#x", (*got)[1].lp, (*got)[2].lp)
	}

	capturePosts(t, 2)
	if err := KeyTap("a", 0x55); err == nil {
		t.Error("post failure must be returned")
	}
	if err := KeyToggle("a", 0x55); err != nil {
		t.Errorf("KeyToggle: %v", err)
	}
}

func TestHwndByPidPrefersMainWindow(t *testing.T) {
	oldEnum, oldPid, oldMain := enumTopWindows, pidOfWindow, isMainWindow
	t.Cleanup(func() { enumTopWindows, pidOfWindow, isMainWindow = oldEnum, oldPid, oldMain })
	pids := map[win.HWND]int{1: 9, 2: 7, 3: 7, 4: 7}
	main := map[win.HWND]bool{1: true, 4: true}
	enumTopWindows = func(cb func(win.HWND) bool) {
		for h := win.HWND(1); h <= 4; h++ {
			if !cb(h) {
				return
			}
		}
	}
	pidOfWindow = func(h win.HWND) int { return pids[h] }
	isMainWindow = func(h win.HWND) bool { return main[h] }

	if got := hwndByPid(7); got != 4 {
		t.Errorf("hwndByPid(7) = %d, want main window 4 over hidden 2", got)
	}
	main[4] = false
	if got := hwndByPid(7); got != 2 {
		t.Errorf("fallback = %d, want first match 2", got)
	}
	if got := hwndByPid(5); got != 0 {
		t.Errorf("unknown pid = %d, want 0", got)
	}
}

func TestDoubleClickNeedsDblClksClass(t *testing.T) {
	got := capturePosts(t, 0)
	oldDbl := classDblClks
	t.Cleanup(func() { classDblClks = oldDbl })
	ev := MouseEvent{Kind: MouseButtonDown, Button: "left", Clicks: 2}

	classDblClks = func(win.HWND) bool { return false }
	if err := PostMouse(WindowInfo{ID: 9}, ev); err != nil {
		t.Fatal(err)
	}
	classDblClks = func(win.HWND) bool { return true }
	if err := PostMouse(WindowInfo{ID: 9}, ev); err != nil {
		t.Fatal(err)
	}
	if len(*got) != 2 || (*got)[0].msg != wmLButtonDown || (*got)[1].msg != wmLButtonDblClk {
		t.Errorf("posted %v", msgs(*got))
	}
}

func TestPointArgs(t *testing.T) {
	a := pointArgs(win.POINT{X: -1, Y: 2})
	if unsafe.Sizeof(uintptr(0)) == 8 {
		if len(a) != 1 || uint64(a[0]) != 0x2_FFFFFFFF {
			t.Errorf("64-bit pointArgs = %#x", a)
		}
		return
	}
	if len(a) != 2 || a[0] != 0xFFFFFFFF || a[1] != 2 {
		t.Errorf("32-bit pointArgs = %#x", a)
	}
}
