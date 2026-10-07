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
	"os"
	"path/filepath"
	"strings"
	"unicode/utf16"
	"unsafe"

	"github.com/go-vgo/robotgo/pub"
	"github.com/tailscale/win"
)

// Standalone input: mouse and keyboard events are posted straight to a target
// window (PostMessageW) without moving the user's real cursor or activating
// the window — the "app_post" style background input used by computer-use
// agents.

// WindowInfo describes an on-screen top-level app window in screen
// coordinates.
type WindowInfo struct {
	ID    uint64 `json:"id"`    // HWND
	Pid   int    `json:"pid"`   // owning process
	Owner string `json:"owner"` // executable name without ".exe"
	X     int    `json:"x"`
	Y     int    `json:"y"`
	W     int    `json:"width"`
	H     int    `json:"height"`
}

// Contains reports whether screen point (x, y) lies inside the window.
func (w WindowInfo) Contains(x, y int) bool {
	return x >= w.X && y >= w.Y && x < w.X+w.W && y < w.Y+w.H
}

// MouseKind is the type of a standalone mouse event.
type MouseKind int

// Standalone mouse event kinds.
const (
	MouseMoved      MouseKind = iota // pointer moved (hover)
	MouseButtonDown                  // button pressed
	MouseButtonUp                    // button released
	MouseDragged                     // pointer moved while a button is held
	MouseScrolled                    // wheel scrolled by DX/DY lines
)

// MouseEvent is one standalone mouse event at screen point (X, Y).
type MouseEvent struct {
	Kind   MouseKind
	Button string // "left" (default), "right", "center"
	X, Y   int
	Clicks int // click state: 1 single, 2 posts a double-click message
	// DX, DY are scroll notches; positive DY scrolls up, positive DX scrolls
	// left, matching Scroll.
	DX, DY int
}

// Window messages and key-state flags used for posted mouse input.
const (
	wmMouseMove     = 0x0200
	wmLButtonDown   = 0x0201
	wmLButtonUp     = 0x0202
	wmLButtonDblClk = 0x0203
	wmRButtonDown   = 0x0204
	wmRButtonUp     = 0x0205
	wmRButtonDblClk = 0x0206
	wmMButtonDown   = 0x0207
	wmMButtonUp     = 0x0208
	wmMButtonDblClk = 0x0209
	wmMouseWheel    = 0x020A
	wmMouseHWheel   = 0x020E

	mkLButton = 0x0001
	mkRButton = 0x0002
	mkMButton = 0x0010

	cwpSkipInvisible   = 0x0001
	cwpSkipDisabled    = 0x0002
	cwpSkipTransparent = 0x0004

	gclStyle = ^uintptr(25) // GCL_STYLE (-26)
)

var (
	procChildWindowFromPointEx = modUser32.NewProc("ChildWindowFromPointEx")
	procGetGUIThreadInfo       = modUser32.NewProc("GetGUIThreadInfo")
	getWindowThreadProcessID   = win.GetWindowThreadProcessId
	isChildWindow              = win.IsChild
	postMessageW               = procPostMessageW.Call
	procGetClassLongPtrW       = modUser32.NewProc("GetClassLongPtrW")
	procGetClassLongW          = modUser32.NewProc("GetClassLongW") // 386 has no ...Ptr export
	classDblClks               = func(hwnd win.HWND) bool {
		p := procGetClassLongPtrW
		if p.Find() != nil {
			p = procGetClassLongW
		}
		style, _, _ := p.Call(uintptr(hwnd), gclStyle)
		return style&win.CS_DBLCLKS != 0
	}
	getGUIThreadInfo = func(tid uint32, info *guiThreadInfo) bool {
		r, _, _ := procGetGUIThreadInfo.Call(uintptr(tid), uintptr(unsafe.Pointer(info)))
		return r != 0
	}
)

// guiThreadInfo is the Win32 GUITHREADINFO structure.
type guiThreadInfo struct {
	cbSize        uint32
	flags         uint32
	hwndActive    win.HWND
	hwndFocus     win.HWND
	hwndCapture   win.HWND
	hwndMenuOwner win.HWND
	hwndMoveSize  win.HWND
	hwndCaret     win.HWND
	rcCaret       win.RECT
}

// ListWindows returns visible, non-minimized top-level windows in z-order,
// front-most first. Windows of the calling process are excluded.
func ListWindows() ([]WindowInfo, error) {
	self := os.Getpid()
	names := map[int]string{}
	var out []WindowInfo
	enumWindows(func(hwnd win.HWND) bool {
		if !win.IsWindowVisible(hwnd) || win.IsIconic(hwnd) {
			return true
		}
		pid := windowPid(hwnd)
		if pid <= 0 || pid == self {
			return true
		}
		var r win.RECT
		if !win.GetWindowRect(hwnd, &r) || r.Right-r.Left <= 1 || r.Bottom-r.Top <= 1 {
			return true
		}
		name, ok := names[pid]
		if !ok {
			name, _ = FindName(pid)
			name = strings.TrimSuffix(name, filepath.Ext(name))
			names[pid] = name
		}
		out = append(out, WindowInfo{
			ID: uint64(hwnd), Pid: pid, Owner: name,
			X: int(r.Left), Y: int(r.Top),
			W: int(r.Right - r.Left), H: int(r.Bottom - r.Top),
		})
		return true
	})
	return out, nil
}

// WindowAt returns the front-most app window containing screen point (x, y).
func WindowAt(x, y int) (WindowInfo, error) {
	wins, err := ListWindows()
	if err != nil {
		return WindowInfo{}, err
	}
	for _, w := range wins {
		if w.Contains(x, y) {
			return w, nil
		}
	}
	return WindowInfo{}, ErrNotFound
}

// PostMouse posts a mouse message to the deepest child of window w under
// the event point. The real cursor does not move and the window is not
// activated.
func PostMouse(w WindowInfo, ev MouseEvent) error {
	if w.ID == 0 {
		return ErrNotFound
	}
	target := childAt(win.HWND(w.ID), ev.X, ev.Y)

	if ev.Kind == MouseScrolled {
		// Wheel messages carry screen coordinates in lParam.
		lp := makeLParam(ev.X, ev.Y)
		if err := postWheel(target, wmMouseWheel, ev.DY, lp); err != nil {
			return err
		}
		// WM_MOUSEHWHEEL is right-positive; robotgo is left-positive.
		return postWheel(target, wmMouseHWheel, -ev.DX, lp)
	}

	clicks := ev.Clicks
	if clicks == 2 && !classDblClks(target) {
		// Windows only sends DBLCLK to classes with CS_DBLCLKS; others get a
		// second plain button-down.
		clicks = 1
	}
	msg, wp, err := mouseMessage(ev.Kind, ev.Button, clicks)
	if err != nil {
		return err
	}
	p := win.POINT{X: int32(ev.X), Y: int32(ev.Y)}
	win.ScreenToClient(target, &p)
	return postMessage(target, msg, wp, makeLParam(int(p.X), int(p.Y)))
}

// PostKeyTap taps key (with optional modifiers) in window w via
// WM_KEYDOWN/WM_KEYUP messages.
func PostKeyTap(w WindowInfo, key string, mods ...string) error {
	if err := PostKeyToggle(w, key, "down", mods...); err != nil {
		return err
	}
	pub.MilliSleep(pub.KeySleep)
	return PostKeyToggle(w, key, "up", mods...)
}

// PostKeyToggle presses ("down") or releases ("up") key in window w.
// Modifiers are pressed before the key and released after it.
func PostKeyToggle(w WindowInfo, key, direction string, mods ...string) error {
	if w.ID == 0 {
		return ErrNotFound
	}
	if direction != "down" && direction != "up" {
		return errors.New(`robotgo: direction must be "down" or "up"`)
	}
	vk, implied, ok := keyToVK(key)
	if !ok {
		return errors.New("robotgo: unknown key: " + key)
	}
	mvks, err := modVKs(mods)
	if err != nil {
		return err
	}
	mvks = append(mvks, impliedVKs(implied)...)
	return postKeySeq(keyTarget(win.HWND(w.ID)), vk, mvks, direction == "up")
}

// PostType types text into window w as WM_CHAR messages.
func PostType(w WindowInfo, text string) error {
	if w.ID == 0 {
		return ErrNotFound
	}
	hwnd := keyTarget(win.HWND(w.ID))
	for _, u := range utf16.Encode([]rune(charText(text))) {
		if err := postMessage(hwnd, wmChar, uintptr(u), 1); err != nil {
			return err
		}
		pub.MilliSleep(pub.KeySleep)
	}
	return nil
}

// charText maps newlines to '\r': Enter arrives as WM_CHAR 0x0D, while 0x0A
// is Ctrl+Enter (linefeed), which edit controls ignore.
func charText(text string) string {
	return strings.ReplaceAll(strings.ReplaceAll(text, "\r\n", "\r"), "\n", "\r")
}

// keyTarget uses the thread's focused control only when it belongs to hwnd.
func keyTarget(hwnd win.HWND) win.HWND {
	tid := getWindowThreadProcessID(hwnd, nil)
	info := guiThreadInfo{}
	info.cbSize = uint32(unsafe.Sizeof(info))
	if tid != 0 && getGUIThreadInfo(tid, &info) && info.hwndFocus != 0 {
		if info.hwndFocus == hwnd || isChildWindow(hwnd, info.hwndFocus) {
			return info.hwndFocus
		}
	}
	return hwnd
}

// impliedVKs maps VkKeyScan's modifier bits (1 shift, 2 ctrl, 4 alt) to
// virtual keys, so "A" or "!" is posted with shift held.
func impliedVKs(bits uint8) []uint16 {
	var out []uint16
	if bits&4 != 0 {
		out = append(out, win.VK_MENU)
	}
	if bits&2 != 0 {
		out = append(out, win.VK_CONTROL)
	}
	if bits&1 != 0 {
		out = append(out, win.VK_SHIFT)
	}
	return out
}

// childAt descends to the deepest visible, enabled child window of h under
// screen point (x, y).
func childAt(h win.HWND, x, y int) win.HWND {
	for {
		p := win.POINT{X: int32(x), Y: int32(y)}
		win.ScreenToClient(h, &p)
		args := append([]uintptr{uintptr(h)}, pointArgs(p)...)
		args = append(args, cwpSkipInvisible|cwpSkipDisabled|cwpSkipTransparent)
		c, _, _ := procChildWindowFromPointEx.Call(args...)
		if c == 0 || win.HWND(c) == h {
			return h
		}
		h = win.HWND(c)
	}
}

// pointArgs passes a POINT by value: one register on 64-bit (X low, Y high
// 32 bits), two stack slots on 32-bit.
func pointArgs(p win.POINT) []uintptr {
	x, y := uint64(uint32(p.X)), uint64(uint32(p.Y))
	if unsafe.Sizeof(uintptr(0)) == 8 {
		return []uintptr{uintptr(x | y<<32)}
	}
	return []uintptr{uintptr(x), uintptr(y)}
}

// mouseMessage maps a standalone event kind and button to a window message
// and its wParam key-state flags.
func mouseMessage(kind MouseKind, btn string, clicks int) (uint32, uintptr, error) {
	var down, up, dbl uint32
	var mk uintptr
	switch btn {
	case "", "left":
		down, up, dbl, mk = wmLButtonDown, wmLButtonUp, wmLButtonDblClk, mkLButton
	case "right":
		down, up, dbl, mk = wmRButtonDown, wmRButtonUp, wmRButtonDblClk, mkRButton
	case "center", "middle":
		down, up, dbl, mk = wmMButtonDown, wmMButtonUp, wmMButtonDblClk, mkMButton
	default:
		return 0, 0, errors.New("robotgo: unknown mouse button: " + btn)
	}
	switch kind {
	case MouseMoved:
		return wmMouseMove, 0, nil
	case MouseDragged:
		return wmMouseMove, mk, nil
	case MouseButtonDown:
		if clicks == 2 {
			return dbl, mk, nil
		}
		return down, mk, nil
	case MouseButtonUp:
		return up, 0, nil
	}
	return 0, 0, errors.New("robotgo: unknown mouse event kind")
}

// modVKs resolves modifier names to virtual-key codes.
func modVKs(mods []string) ([]uint16, error) {
	out := make([]uint16, 0, len(mods))
	for _, m := range mods {
		vk, _, ok := keyToVK(m)
		if !ok {
			return nil, errors.New("robotgo: unknown modifier: " + m)
		}
		out = append(out, vk)
	}
	return out, nil
}

func postMessage(hwnd win.HWND, msg uint32, wp, lp uintptr) error {
	if r, _, err := postMessageW(uintptr(hwnd), uintptr(msg), wp, lp); r == 0 {
		return err
	}
	return nil
}

// makeLParam packs signed 16-bit x/y coordinates (MAKELPARAM).
func makeLParam(x, y int) uintptr {
	return uintptr(uint32(uint16(int16(x))) | uint32(uint16(int16(y)))<<16)
}

// maxWheelNotches keeps notches*WHEEL_DELTA within the int16 wParam word.
const maxWheelNotches = 0x7fff / wheelDelta

// postWheel posts notches wheel steps in chunks that fit the delta word.
func postWheel(hwnd win.HWND, msg uint32, notches int, lp uintptr) error {
	for notches != 0 {
		step := max(-maxWheelNotches, min(maxWheelNotches, notches))
		if err := postMessage(hwnd, msg, wheelWParam(step*wheelDelta), lp); err != nil {
			return err
		}
		notches -= step
	}
	return nil
}

// wheelWParam packs a signed wheel delta into the high word of wParam.
func wheelWParam(delta int) uintptr {
	return uintptr(uint32(uint16(int16(delta))) << 16)
}
