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

// Untagged on purpose: key tests that only use the public API every
// backend wires (Cgo and -tags mac/win/x11/wayland/libei/purego) plus the
// portable helpers in robotgo_pub.go / keycode.go / ps.go. Cgo-internal
// key checks live in key_c_test.go.

package robotgo

import (
	"os"
	"runtime"
	"testing"

	"github.com/vcaesar/tt"
)

func TestKeyConstants(t *testing.T) {
	letters := []string{
		KeyA, KeyB, KeyC, KeyD, KeyE, KeyF, KeyG, KeyH, KeyI, KeyJ, KeyK, KeyL, KeyM,
		KeyN, KeyO, KeyP, KeyQ, KeyR, KeyS, KeyT, KeyU, KeyV, KeyW, KeyX, KeyY, KeyZ,
	}
	for i, k := range letters {
		tt.Equal(t, string(rune('a'+i)), k)
	}
	tt.Equal(t, "A", CapA)
	tt.Equal(t, "Z", CapZ)
	tt.Equal(t, "0", Key0)
	tt.Equal(t, "9", Key9)

	// aliases name the same key
	tt.Equal(t, "esc", Esc)
	tt.Equal(t, "escape", Escape)
	tt.Equal(t, "caps", Caps)
	tt.Equal(t, "capslock", Capslock)
	tt.Equal(t, "ctrl", Ctrl)
	tt.Equal(t, "control", Control)
	tt.Equal(t, "cmd", Cmd)
	tt.Equal(t, "enter", Enter)
	tt.Equal(t, "num_enter", NumEnter)
	tt.Equal(t, "scroll_lock", ScrollLock)
	tt.Equal(t, "pause_break", PauseBreak)
}

func TestKeyCode(t *testing.T) {
	m := MouseMap["left"]
	tt.Equal(t, 1, m)

	k := Keycode["1"]
	tt.Equal(t, 2, k)

	s := Special["+"]
	tt.Equal(t, "=", s)

	tt.Equal(t, "0", Key0)
	tt.Equal(t, "a", KeyA)

	// every shifted symbol maps to a single-char base key with a hook code
	for sym, base := range Special {
		tt.Equal(t, 1, len(sym), sym)
		tt.Equal(t, 1, len(base), sym)
		_, ok := Keycode[base]
		tt.True(t, ok, base)
	}
}

func TestCmdCtrl(t *testing.T) {
	if runtime.GOOS == "darwin" {
		tt.Equal(t, Cmd, CmdCtrl())
	} else {
		tt.Equal(t, Ctrl, CmdCtrl())
	}
}

func TestToInterfacesStrings(t *testing.T) {
	in := []string{"alt", "cmd"}
	out := ToInterfaces(in)
	tt.Equal(t, 2, len(out))
	tt.Equal(t, "alt", out[0])
	tt.Equal(t, in, ToStrings(out))

	tt.Equal(t, 0, len(ToInterfaces(nil)))
	tt.Equal(t, 0, len(ToStrings(nil)))
}

func TestCharCodeAt(t *testing.T) {
	tt.Equal(t, 115, CharCodeAt("s", 0))
	tt.Equal(t, 'c', CharCodeAt("abc", 2))
	// index counts runes, not bytes
	tt.Equal(t, '界', CharCodeAt("世界", 1))
	tt.Equal(t, 0, CharCodeAt("abc", 3))
	tt.Equal(t, 0, CharCodeAt("", 0))
}

func TestToUC(t *testing.T) {
	uc := ToUC("abc\\\\cd/s@世界")
	tt.Equal(t, "[a b c \\ \\ c d / s @ U4e16 U754c]", uc)

	tt.Equal(t, []string{`"`}, ToUC(`"`))
	tt.Equal(t, []string{"\\"}, ToUC("\\"))
	tt.Equal(t, 0, len(ToUC("")))
}

func TestSetDelay(t *testing.T) {
	k, m := KeySleep, MouseSleep
	defer func() { KeySleep, MouseSleep = k, m }()

	SetDelay()
	tt.Equal(t, 10, KeySleep)
	tt.Equal(t, 10, MouseSleep)

	SetDelay(25)
	tt.Equal(t, 25, KeySleep)
	tt.Equal(t, 25, MouseSleep)
}

// Empty input never touches the backend and reports zero typed runes.
func TestTypeEmpty(t *testing.T) {
	tt.Equal(t, 0, Type(""))
	tt.Equal(t, 0, Type("", 0))
	tt.Nil(t, TypeStr(""))
	tt.Nil(t, TypeDelay("", 0))
	tt.Nil(t, TypeStrDelay("", 0))
}

// Unknown key names must fail on every backend, with or without a session.
func TestKeyInvalidName(t *testing.T) {
	const bad = "no_such_key"
	tt.NotNil(t, KeyTap(bad))
	tt.NotNil(t, KeyTap(bad, "ctrl"))
	tt.NotNil(t, KeyToggle(bad))
	tt.NotNil(t, KeyToggle(bad, "up"))
	tt.NotNil(t, KeyDown(bad))
	tt.NotNil(t, KeyUp(bad))
	tt.NotNil(t, KeyPress(bad))
}

func TestKey(t *testing.T) {
	requireDisplay(t)

	e := KeyTap("v", CmdCtrl())
	tt.Nil(t, e)

	e = KeyTap("enter")
	tt.Nil(t, e)

	e = KeyToggle("v", "up")
	tt.Nil(t, e)

	e = KeyDown("a")
	tt.Nil(t, e)
	e = KeyUp("a")
	tt.Nil(t, e)

	// array and pid argument forms, same on Cgo and pure-Go backends
	e = KeyTap("i", []string{"alt", CmdCtrl()})
	tt.Nil(t, e)
	e = KeyToggle("a", "down", []string{"alt", CmdCtrl()})
	tt.Nil(t, e)
	e = KeyToggle("a", "up", []string{"alt", CmdCtrl()})
	tt.Nil(t, e)
	e = KeyToggle("a", "up", "alt", CmdCtrl())
	tt.Nil(t, e)
	e = KeyToggle("v", "up", 0)
	tt.Nil(t, e)

	// Not on Mac keyboards: keycode.h maps them to K_NOT_A_KEY there.
	e = KeyTap(ScrollLock)
	if runtime.GOOS == "darwin" {
		tt.NotNil(t, e)
	} else {
		tt.Nil(t, e)
	}

	e = KeyTap(PauseBreak)
	if runtime.GOOS == "darwin" {
		tt.NotNil(t, e)
	} else {
		tt.Nil(t, e)
	}

	e = KeyTap("nonexistent_key")
	tt.NotNil(t, e)
}

func TestKeyPress(t *testing.T) {
	requireDisplay(t)

	tt.Nil(t, KeyPress("a"))
	tt.Nil(t, KeyPress("enter", 0))
	tt.Nil(t, KeyPress("a", "shift"))
	// KeyToggle without a direction defaults to down
	tt.Nil(t, KeyToggle("shift"))
	tt.Nil(t, KeyToggle("shift", "up"))
}

func TestCmdV(t *testing.T) {
	requireDisplay(t)

	tt.Nil(t, CmdV())
	tt.Nil(t, CmdV(0))
}

func TestType(t *testing.T) {
	requireDisplay(t)

	tt.Equal(t, 2, Type("ab"))
	tt.Equal(t, 1, Type("c", 0))
	tt.Nil(t, TypeStr("d"))
	tt.Nil(t, TypeDelay("e", 1))
}

func TestTypeStr(t *testing.T) {
	requireDisplay(t)
	e := PasteStr("s")
	skipNoClipboard(t, e)
	tt.Nil(t, e)

	l, e := Paste("世界")
	tt.Nil(t, e)
	tt.Equal(t, 2, l)
}

// KeyTap with the pid argument (pid-targeted delivery on macOS/Windows,
// accepted and ignored on X11/Wayland/libei).
func TestKeyTapPid(t *testing.T) {
	requireDisplay(t)
	if runtime.GOOS == "windows" {
		t.Skip("posting to another process's window is intrusive on CI")
	}

	e := KeyTap("v", os.Getpid(), CmdCtrl())
	if e != nil && runtime.GOOS == "darwin" {
		// A test binary has no window to receive the event.
		t.Skipf("pid delivery: %v", e)
	}
	tt.Nil(t, e)
}

func TestAppInfo(t *testing.T) {
	name, path, pid := appInfo(0)
	tt.Equal(t, "", name)
	tt.Equal(t, "", path)
	tt.Equal(t, 0, pid)

	name, _, pid = appInfo(-1)
	tt.Equal(t, "", name)
	tt.Equal(t, 0, pid)

	self := os.Getpid()
	name, _, pid = appInfo(self)
	tt.Equal(t, self, pid)
	tt.NotEqual(t, "", name)
}
