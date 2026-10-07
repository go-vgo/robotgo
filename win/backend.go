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
	"sync"
	"unicode/utf16"

	"github.com/go-vgo/robotgo/pub"
	"github.com/tailscale/win"
)

// Backend selects how the global mouse and keyboard API injects input.
type Backend int

const (
	// BackendSendInput injects input with SendInput/SetCursorPos (default).
	BackendSendInput Backend = iota
	// BackendHID injects hardware-level strokes through the Interception
	// driver. Characters missing from the keyboard layout fall back to
	// SendInput Unicode events in Type.
	BackendHID
	// BackendMessage posts window messages (PostMessageW) without moving the
	// real cursor or activating a window. Mouse events go to the message
	// target, or the window under the event point when no target is set, and
	// use a virtual pointer started at the real cursor position. Key events
	// go to the target, or the foreground window.
	BackendMessage
)

var input struct {
	mu      sync.Mutex
	backend Backend
	target  WindowInfo
	x, y    int
	posSet  bool
	held    string  // button held under BackendMessage, "" if none
	keys    keyHeld // Alt/Ctrl posted under BackendMessage
}

var getForegroundWindow = win.GetForegroundWindow

// SetBackend selects the input backend. BackendHID initializes the
// Interception driver immediately and fails if it is unavailable.
func SetBackend(b Backend) error {
	switch b {
	case BackendSendInput, BackendMessage:
	case BackendHID:
		if err := InitHID(); err != nil {
			return err
		}
	default:
		return errors.New("robotgo: unknown input backend")
	}
	input.mu.Lock()
	input.backend, input.posSet, input.held, input.keys = b, false, "", keyHeld{}
	input.mu.Unlock()
	return nil
}

// GetBackend returns the current input backend.
func GetBackend() Backend {
	input.mu.Lock()
	defer input.mu.Unlock()
	return input.backend
}

// SetMessageTarget fixes the window BackendMessage posts to. A zero
// WindowInfo restores automatic targeting.
func SetMessageTarget(w WindowInfo) {
	input.mu.Lock()
	input.target = w
	input.mu.Unlock()
}

// pointer returns the position mouse input applies to: the virtual pointer
// under BackendMessage, otherwise the real cursor.
func pointer() (int, int) {
	input.mu.Lock()
	defer input.mu.Unlock()
	if input.backend == BackendMessage && input.posSet {
		return input.x, input.y
	}
	return Location()
}

// postPointer posts ev to the message target or the window under ev.
func postPointer(ev MouseEvent) error {
	input.mu.Lock()
	w := input.target
	input.mu.Unlock()
	if w.ID == 0 {
		var err error
		if w, err = WindowAt(ev.X, ev.Y); err != nil {
			return err
		}
	}
	return PostMouse(w, ev)
}

func inputMoveTo(x, y int) error {
	switch GetBackend() {
	case BackendHID:
		return hidMove(x, y)
	case BackendMessage:
		input.mu.Lock()
		input.x, input.y, input.posSet = x, y, true
		held := input.held
		input.mu.Unlock()
		ev := MouseEvent{Kind: MouseMoved, X: x, Y: y}
		if held != "" {
			ev.Kind, ev.Button = MouseDragged, held
		}
		return postPointer(ev)
	}
	return setCursorPos(x, y)
}

// inputButton presses or releases btn; clicks 2 marks a double-click press.
func inputButton(btn string, up bool, clicks int) error {
	switch GetBackend() {
	case BackendHID:
		return hidButton(btn, up)
	case BackendMessage:
		x, y := pointer()
		ev := MouseEvent{Kind: MouseButtonDown, Button: btn, X: x, Y: y, Clicks: clicks}
		if up {
			ev.Kind = MouseButtonUp
		}
		if err := postPointer(ev); err != nil {
			return err
		}
		input.mu.Lock()
		if up {
			input.held = ""
		} else {
			input.held = btn
		}
		input.mu.Unlock()
		return nil
	}
	down, release := mouseButtonFlags(btn)
	if up {
		down = release
	}
	return sendMouseInput(down, 0, 0, 0)
}

// inputWheel scrolls by notches; positive y scrolls up, positive x left.
func inputWheel(x, y int) error {
	switch GetBackend() {
	case BackendHID:
		return hidWheel(x, y)
	case BackendMessage:
		px, py := pointer()
		return postPointer(MouseEvent{Kind: MouseScrolled, X: px, Y: py, DX: x, DY: y})
	}
	if y != 0 {
		if err := sendMouseInput(win.MOUSEEVENTF_WHEEL, uint32(int32(y*wheelDelta)), 0, 0); err != nil {
			return err
		}
	}
	if x != 0 {
		// Win32 horizontal wheel is right-positive; robotgo is left-positive.
		return sendMouseInput(win.MOUSEEVENTF_HWHEEL, uint32(int32(-x*wheelDelta)), 0, 0)
	}
	return nil
}

// keyWindow returns the window BackendMessage posts key input to.
func keyWindow() (win.HWND, error) {
	input.mu.Lock()
	hwnd := win.HWND(input.target.ID)
	input.mu.Unlock()
	if hwnd == 0 {
		hwnd = getForegroundWindow()
	}
	if hwnd == 0 {
		return 0, ErrNotFound
	}
	return keyTarget(hwnd), nil
}

func inputKey(vk uint16, up bool) error {
	switch GetBackend() {
	case BackendHID:
		return hidKey(vk, up)
	case BackendMessage:
		hwnd, err := keyWindow()
		if err != nil {
			return err
		}
		input.mu.Lock()
		defer input.mu.Unlock()
		return input.keys.post(hwnd, vk, up)
	}
	return sendVK(vk, up)
}

func keyDelay() {
	pub.MilliSleep(pub.KeySleep)
}

func inputText(str string) error {
	switch GetBackend() {
	case BackendHID:
		for _, r := range str {
			if err := hidRune(r); err != nil {
				return err
			}
			keyDelay()
		}
		return nil
	case BackendMessage:
		hwnd, err := keyWindow()
		if err != nil {
			return err
		}
		for _, u := range utf16.Encode([]rune(str)) {
			if err := postMessage(hwnd, wmChar, uintptr(u), 1); err != nil {
				return err
			}
			keyDelay()
		}
		return nil
	}
	return sendUnicodeText(str, true)
}

// sendUnicodeText types str as SendInput Unicode events, stopping at the
// first event SendInput rejects.
func sendUnicodeText(str string, delay bool) error {
	for _, u := range utf16.Encode([]rune(str)) {
		if !sendUnicode(u, false) || !sendUnicode(u, true) {
			return errSendInput
		}
		if delay {
			keyDelay()
		}
	}
	return nil
}

// hidRune types r as scan codes with the modifiers the layout needs, or as
// SendInput Unicode events when the layout has no key for it.
func hidRune(r rune) error {
	vk, bits, ok := keyToVK(string(r))
	if r == '\n' { // VkKeyScan maps '\n' to Ctrl+Enter
		vk, bits, ok = win.VK_RETURN, 0, true
	}
	if !ok {
		return sendUnicodeText(string(r), false)
	}
	vks := append(impliedVKs(bits), vk)
	pressed, err := pressKeys(hidKey, vks)
	// Release whatever was pressed, even after a failure, so no modifier
	// is left held at the driver level.
	if rerr := releaseKeys(hidKey, pressed); err == nil {
		err = rerr
	}
	return err
}
