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

//go:build darwin
// +build darwin

package clipboard

import (
	"slices"
	"testing"

	"github.com/go-vgo/robotgo/internal/cliplock"
)

func TestUTF8EnvOverridesLocale(t *testing.T) {
	t.Setenv("LANG", "C")
	t.Setenv("LC_CTYPE", "C")
	t.Setenv("LC_ALL", "C")

	for _, env := range [][]string{getPasteCommand().Environ(), getCopyCommand().Environ()} {
		for _, kv := range []string{"LANG=en_US.UTF-8", "LC_ALL=en_US.UTF-8"} {
			if !slices.Contains(env, kv) {
				t.Errorf("env missing %q", kv)
			}
		}
		if slices.Contains(env, "LC_ALL=C") {
			t.Error("parent LC_ALL=C not overridden")
		}
	}
}

func TestRoundTripCLocale(t *testing.T) {
	cliplock.Lock(t)
	t.Setenv("LC_ALL", "C")
	const want = "日本語 💩"
	if err := writeAll(want); err != nil {
		t.Fatal(err)
	}
	got, err := readAll()
	if err != nil {
		t.Fatal(err)
	}
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}
