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
	"fmt"
	"reflect"
	"testing"

	"github.com/go-vgo/robotgo/pub"
	"github.com/tailscale/win"
)

func TestCharText(t *testing.T) {
	cases := map[string]string{
		"ab":       "ab",
		"a\nb":     "a\rb",
		"a\r\nb":   "a\rb",
		"a\rb\n\n": "a\rb\r\r",
	}
	for in, want := range cases {
		if got := charText(in); got != want {
			t.Errorf("charText(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestMouseMessage(t *testing.T) {
	cases := []struct {
		kind   MouseKind
		btn    string
		clicks int
		msg    uint32
		wp     uintptr
	}{
		{MouseMoved, "", 0, wmMouseMove, 0},
		{MouseDragged, "left", 0, wmMouseMove, mkLButton},
		{MouseButtonDown, "left", 1, wmLButtonDown, mkLButton},
		{MouseButtonDown, "left", 2, wmLButtonDblClk, mkLButton},
		{MouseButtonUp, "right", 1, wmRButtonUp, 0},
		{MouseButtonDown, "center", 1, wmMButtonDown, mkMButton},
	}
	for _, c := range cases {
		msg, wp, err := mouseMessage(c.kind, c.btn, c.clicks)
		if err != nil || msg != c.msg || wp != c.wp {
			t.Errorf("mouseMessage(%d,%q,%d) = 0x%x,%d,%v", c.kind, c.btn, c.clicks, msg, wp, err)
		}
	}
	if _, _, err := mouseMessage(MouseButtonDown, "bogus", 1); err == nil {
		t.Error("unknown button should fail")
	}
}

func TestPackParams(t *testing.T) {
	if got := makeLParam(10, 20); got != 20<<16|10 {
		t.Errorf("makeLParam = 0x%x", got)
	}
	if got := makeLParam(-1, 0); got != 0xffff {
		t.Errorf("makeLParam(-1,0) = 0x%x", got)
	}
	if got := wheelWParam(-120); got != 0xff880000 {
		t.Errorf("wheelWParam(-120) = 0x%x", got)
	}
}

func TestKeyParams(t *testing.T) {
	if got := keyLParam(0, false); got != 1 {
		t.Errorf("keyLParam down = 0x%x, want repeat count 1", got)
	}
	if got := keyLParam(0, true); got&0xC0000001 != 0xC0000001 {
		t.Errorf("keyLParam up = 0x%x, want transition bits", got)
	}
	got := impliedVKs(1 | 2 | 4)
	if len(got) != 3 || got[0] != 0x12 || got[1] != 0x11 || got[2] != 0x10 {
		t.Errorf("impliedVKs = %v", got)
	}
	if len(impliedVKs(0)) != 0 {
		t.Error("no implied modifiers expected")
	}
}

func TestKeyTarget(t *testing.T) {
	oldThread, oldInfo, oldChild := getWindowThreadProcessID, getGUIThreadInfo, isChildWindow
	t.Cleanup(func() {
		getWindowThreadProcessID, getGUIThreadInfo, isChildWindow = oldThread, oldInfo, oldChild
	})
	for _, tc := range []struct {
		name  string
		tid   uint32
		ok    bool
		focus win.HWND
		child bool
		want  win.HWND
	}{
		{"same window", 5, true, 101, false, 101},
		{"child", 5, true, 102, true, 102},
		{"other window on thread", 5, true, 202, false, 101},
		{"no focus", 5, true, 0, false, 101},
		{"query failed", 5, false, 102, true, 101},
		{"invalid window", 0, true, 102, true, 101},
	} {
		t.Run(tc.name, func(t *testing.T) {
			getWindowThreadProcessID = func(hwnd win.HWND, pid *uint32) uint32 {
				if hwnd != 101 || pid != nil {
					t.Fatalf("unexpected thread lookup: %d, %v", hwnd, pid)
				}
				return tc.tid
			}
			getGUIThreadInfo = func(tid uint32, info *guiThreadInfo) bool {
				if tid != 5 || info.cbSize == 0 {
					t.Fatalf("invalid GUI thread query: %d, %+v", tid, info)
				}
				info.hwndFocus = tc.focus
				return tc.ok
			}
			isChildWindow = func(parent, child win.HWND) bool {
				if parent != 101 || child != tc.focus || child == 0 {
					t.Fatalf("unexpected child lookup: %d, %d", parent, child)
				}
				return tc.child
			}
			if got := keyTarget(101); got != tc.want {
				t.Fatalf("keyTarget = %d, want %d", got, tc.want)
			}
		})
	}
}

func TestStandaloneKeyboardPosting(t *testing.T) {
	oldPost, oldThread, oldSleep := postMessageW, getWindowThreadProcessID, pub.KeySleep
	t.Cleanup(func() {
		postMessageW, getWindowThreadProcessID, pub.KeySleep = oldPost, oldThread, oldSleep
	})
	getWindowThreadProcessID = func(win.HWND, *uint32) uint32 { return 0 }
	pub.KeySleep = 0
	w := WindowInfo{ID: 101}
	type message struct {
		kind, key uintptr
	}
	postErr := errors.New("PostMessageW failed")
	// A failed press releases the keys already pressed; a release continues
	// past failures. Other calls stop at the failing message.
	pressCleanup := func(want []message, failAt int) []message {
		got := append([]message{}, want[:failAt]...)
		for i := failAt - 2; i >= 0; i-- {
			got = append(got, message{wmKeyUp, want[i].key})
		}
		return got
	}
	releaseAll := func(want []message, _ int) []message { return want }
	cases := []struct {
		name   string
		call   func() error
		want   []message
		onFail func(want []message, failAt int) []message
	}{
		{"down", func() error { return PostKeyToggle(w, "enter", "down", "ctrl", "shift") },
			[]message{{wmKeyDown, win.VK_CONTROL}, {wmKeyDown, win.VK_SHIFT}, {wmKeyDown, win.VK_RETURN}}, pressCleanup},
		{"up", func() error { return PostKeyToggle(w, "enter", "up", "ctrl", "shift") },
			[]message{{wmKeyUp, win.VK_RETURN}, {wmKeyUp, win.VK_SHIFT}, {wmKeyUp, win.VK_CONTROL}}, releaseAll},
		{"tap", func() error { return PostKeyTap(w, "enter") },
			[]message{{wmKeyDown, win.VK_RETURN}, {wmKeyUp, win.VK_RETURN}}, nil},
		{"text", func() error { return PostType(w, "a😀b") },
			[]message{{wmChar, 'a'}, {wmChar, 0xd83d}, {wmChar, 0xde00}, {wmChar, 'b'}}, nil},
		{"empty text", func() error { return PostType(w, "") }, nil, nil},
	}
	for _, tc := range cases {
		for failAt := 0; failAt <= len(tc.want); failAt++ {
			t.Run(fmt.Sprintf("%s/failAt=%d", tc.name, failAt), func(t *testing.T) {
				var got []message
				postMessageW = func(args ...uintptr) (uintptr, uintptr, error) {
					if len(args) != 4 || args[0] != uintptr(w.ID) {
						t.Fatalf("unexpected post: %v", args)
					}
					if args[3]&0xffff != 1 {
						t.Errorf("missing repeat count: %#x", args[3])
					}
					if args[1] == wmKeyUp && args[3]&0xc0000000 != 0xc0000000 {
						t.Errorf("missing key-up bits: %#x", args[3])
					}
					got = append(got, message{args[1], args[2]})
					if len(got) == failAt {
						return 0, 0, postErr
					}
					return 1, 0, postErr // LastError can be stale on success.
				}
				err := tc.call()
				want := tc.want
				if failAt > 0 {
					if tc.onFail != nil {
						want = tc.onFail(want, failAt)
					} else {
						want = want[:failAt]
					}
					if !errors.Is(err, postErr) {
						t.Errorf("error = %v, want %v", err, postErr)
					}
				} else if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if !reflect.DeepEqual(got, want) {
					t.Errorf("posted = %v, want %v", got, want)
				}
			})
		}
	}
}

func TestStandaloneWindowInfo(t *testing.T) {
	w := WindowInfo{X: 10, Y: 10, W: 5, H: 5}
	if !w.Contains(10, 10) || w.Contains(15, 10) {
		t.Error("Contains: wrong bounds")
	}
	if err := PostMouse(WindowInfo{}, MouseEvent{}); err != ErrNotFound {
		t.Errorf("PostMouse: %v", err)
	}
	if err := PostKeyToggle(WindowInfo{ID: 1}, "a", "sideways"); err == nil {
		t.Error("bad direction should fail")
	}
}
