//go:build linux && cgo && (x11 || (!wayland && !libei && !purego))
// +build linux
// +build cgo
// +build x11 !wayland,!libei,!purego

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

// X11 end-to-end check for Type (#640): the text is typed into a window owned
// by the test and read back through XLookupString (internal/xharness), so it
// sees exactly what an application would. The harness needs Xlib, so this
// runs with cgo on, against both the Cgo backend and -tags x11.

package robotgo_test

import (
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/go-vgo/robotgo"
	"github.com/go-vgo/robotgo/internal/xharness"
)

// xkbLayout returns the "layout" and "variant" fields of setxkbmap -query.
func xkbLayout(t *testing.T) (layout, variant string) {
	t.Helper()
	out, err := exec.Command("setxkbmap", "-query").Output()
	if err != nil {
		t.Skipf("setxkbmap -query: %v", err)
	}
	for _, line := range strings.Split(string(out), "\n") {
		k, v, _ := strings.Cut(line, ":")
		switch strings.TrimSpace(k) {
		case "layout":
			layout = strings.TrimSpace(v)
		case "variant":
			variant = strings.TrimSpace(v)
		}
	}
	return layout, variant
}

func setXkbLayout(t *testing.T, layout, variant string) {
	t.Helper()
	args := []string{"-layout", layout}
	if variant != "" {
		args = append(args, "-variant", variant)
	}
	if out, err := exec.Command("setxkbmap", args...).CombinedOutput(); err != nil {
		t.Skipf("setxkbmap %v: %v: %s", args, err, out)
	}
}

// TestTypeLayout types every XKB level of a German layout: AltGr (@ | { } \
// ~), Shift (/ < > " ' Q), plain keys, and z/y which swap places (#640: '@'
// used to come out as 'Q' and '/' as '7'). The layout is only switched on CI
// (CI env set, setxkbmap present); elsewhere the current layout is used.
func TestTypeLayout(t *testing.T) {
	if os.Getenv("DISPLAY") == "" {
		t.Skip("no DISPLAY")
	}
	const text = `test@example.org/ <>|{}\~7Qzy"'`

	layouts := []string{""}
	if _, err := exec.LookPath("setxkbmap"); err == nil && os.Getenv("CI") != "" {
		layouts = append(layouts, "de", "us")
	}
	for _, layout := range layouts {
		name := layout
		if name == "" {
			name = "current"
		}
		t.Run(name, func(t *testing.T) {
			if layout != "" {
				prev, variant := xkbLayout(t)
				setXkbLayout(t, layout, "")
				t.Cleanup(func() { setXkbLayout(t, prev, variant) })
			}
			if err := xharness.Open(); err != nil {
				t.Skip(err)
			}
			defer xharness.Close()

			robotgo.Type(text)
			robotgo.MilliSleep(100)

			got := xharness.Read()
			if got == "" && os.Getenv("CI") == "" {
				t.Skip("no key events reached the test window (focus denied?)")
			}
			if got != text {
				t.Errorf("Type(%q) was received as %q", text, got)
			}
		})
	}
}
