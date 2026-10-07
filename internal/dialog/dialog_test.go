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

package dialog

import (
	"errors"
	"slices"
	"testing"

	"github.com/vcaesar/tt"
)

func fakeLook(found ...string) func(string) (string, error) {
	return func(name string) (string, error) {
		if slices.Contains(found, name) {
			return "/usr/bin/" + name, nil
		}
		return "", errors.New("not found")
	}
}

func TestCommand(t *testing.T) {
	tt.Equal(t, []string{"/usr/bin/zenity", "--question", "--title", "T", "--text", "m",
		"--ok-label", "Yes", "--cancel-label", "No"},
		command(fakeLook("zenity", "xmessage"), "T", "m", "Yes", "No"))

	tt.Equal(t, []string{"/usr/bin/xmessage", "-center", "-title", "T",
		"-buttons", "Yes:0,No:1", "-default", "Yes", "m"},
		command(fakeLook("xmessage"), "T", "m", "Yes", "No"))

	tt.True(t, command(fakeLook(), "T", "m", "Yes", "No") == nil)
}
