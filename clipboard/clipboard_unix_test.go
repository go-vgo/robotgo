// Copyright 2013 @atotto. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
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

//go:build freebsd || linux || netbsd || openbsd || solaris || dragonfly
// +build freebsd linux netbsd openbsd solaris dragonfly

package clipboard

import (
	"bytes"
	"errors"
	"reflect"
	"testing"
)

func fakeLookPath(available ...string) func(string) (string, error) {
	return func(name string) (string, error) {
		for _, a := range available {
			if a == name {
				return "/usr/bin/" + name, nil
			}
		}
		return "", errors.New("not found")
	}
}

func fakeGetenv(vars ...string) func(string) string {
	return func(key string) string {
		for _, v := range vars {
			if v == key {
				return "1"
			}
		}
		return ""
	}
}

func TestDetect(t *testing.T) {
	wsl := []string{"powershell.exe", "clip.exe"}
	cases := []struct {
		name      string
		env       []string
		available []string
		want      *tool
	}{
		{"none", nil, nil, nil},
		{"xclip", []string{"DISPLAY"}, []string{"xclip", "xsel"}, xclipTool},
		{"xsel", []string{"DISPLAY"}, []string{"xsel"}, xselTool},
		{"x11 without DISPLAY", nil, []string{"xclip"}, xclipTool},
		{"wayland prefers wl", []string{"WAYLAND_DISPLAY", "DISPLAY"},
			[]string{"wl-copy", "wl-paste", "xclip"}, wlTool},
		{"wl ignored on x11", []string{"DISPLAY"}, []string{"wl-copy", "wl-paste", "xsel"}, xselTool},
		{"wl needs both tools", []string{"WAYLAND_DISPLAY"}, []string{"wl-copy", "xclip"}, xclipTool},
		{"termux", nil, []string{"termux-clipboard-get", "termux-clipboard-set"}, termuxTool},
		{"termux needs both tools", nil, []string{"termux-clipboard-set"}, nil},
		{"wsl", nil, wsl, wslTool},
		{"wslg prefers wayland", []string{"WAYLAND_DISPLAY"},
			append([]string{"wl-copy", "wl-paste"}, wsl...), wlTool},
		{"wsl before stale xclip", nil, append([]string{"xclip"}, wsl...), wslTool},
		{"tmux", []string{"TMUX"}, []string{"tmux", "xclip"}, tmuxTool},
		{"tmux needs TMUX", nil, []string{"tmux"}, nil},
		{"x11 before tmux", []string{"TMUX", "DISPLAY"}, []string{"tmux", "xclip"}, xclipTool},
	}

	for _, c := range cases {
		if got := detect(fakeGetenv(c.env...), fakeLookPath(c.available...)); got != c.want {
			t.Errorf("%s: got %+v, want %+v", c.name, got, c.want)
		}
	}
}

func TestEncodeUTF16LE(t *testing.T) {
	cases := map[string][]byte{
		"":    {0xff, 0xfe},
		"a\n": {0xff, 0xfe, 'a', 0, '\n', 0},
		"日":   {0xff, 0xfe, 0xe5, 0x65},
		"💩":   {0xff, 0xfe, 0x3d, 0xd8, 0xa9, 0xdc},
	}
	for in, want := range cases {
		if got := encodeUTF16LE(in); !bytes.Equal(got, want) {
			t.Errorf("encodeUTF16LE(%q) = % x, want % x", in, got, want)
		}
	}
}

func TestDecodeWSL(t *testing.T) {
	cases := map[string]string{
		"":             "",
		"x":            "x",
		"日本語\r\n":      "日本語",
		"a\r\n\r\n":    "a\r\n",
		"no newline\n": "no newline\n",
	}
	for in, want := range cases {
		if got := decodeWSL([]byte(in)); got != want {
			t.Errorf("decodeWSL(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestPrimaryFallback(t *testing.T) {
	saved, savedPrimary := current, Primary
	defer func() { current, Primary = saved, savedPrimary }()

	Primary = true
	for _, tl := range []*tool{termuxTool, wslTool, tmuxTool} {
		current = tl
		if got := getPasteCommand().Args; !reflect.DeepEqual(got, tl.paste) {
			t.Errorf("paste args: got %v, want %v", got, tl.paste)
		}
		if got := getCopyCommand().Args; !reflect.DeepEqual(got, tl.copy) {
			t.Errorf("copy args: got %v, want %v", got, tl.copy)
		}
	}
}

func TestPrimaryArgs(t *testing.T) {
	saved, savedPrimary := current, Primary
	defer func() { current, Primary = saved, savedPrimary }()

	current = xclipTool
	for _, step := range []struct {
		primary         bool
		paste, copyArgs []string
	}{
		{false, xclipTool.paste, xclipTool.copy},
		{true, xclipTool.pastePrimary, xclipTool.copyPrimary},
		// toggling back must restore the clipboard selection
		{false, xclipTool.paste, xclipTool.copy},
	} {
		Primary = step.primary
		if got := getPasteCommand().Args; !reflect.DeepEqual(got, step.paste) {
			t.Errorf("primary=%v paste args: got %v, want %v", step.primary, got, step.paste)
		}
		if got := getCopyCommand().Args; !reflect.DeepEqual(got, step.copyArgs) {
			t.Errorf("primary=%v copy args: got %v, want %v", step.primary, got, step.copyArgs)
		}
	}
}

func TestMissingCommands(t *testing.T) {
	saved, savedUnsupported := current, Unsupported
	defer func() { current, Unsupported = saved, savedUnsupported }()

	current, Unsupported = nil, false
	if _, err := readAll(); err != errMissingCommands {
		t.Errorf("readAll: got %v, want %v", err, errMissingCommands)
	}
	if err := writeAll("x"); err != errMissingCommands {
		t.Errorf("writeAll: got %v, want %v", err, errMissingCommands)
	}
}
