// Copyright (c) 2016-2026 AtomAI, All rights reserved.
//
// See COPYRIGHT file at top-level directory of this distribution and at
// https://github.com/go-vgo/robotgo/blob/master/LICENSE
//
// Licensed under Apache License, Version 2.0 <LICENSE-APACHE or
// http://www.apache.org/licenses/LICENSE-2.0>
//
// This file may not be copied, modified, or distributed
// except according to those terms.

package robotgo

import (
	"testing"

	"github.com/vcaesar/tt"
)

func TestGetKeyDown(t *testing.T) {
	down, mods := getKeyDown(nil)
	tt.True(t, down)
	tt.Equal(t, 0, len(mods))

	down, mods = getKeyDown([]string{"down"})
	tt.True(t, down)
	tt.Equal(t, 0, len(mods))

	down, mods = getKeyDown([]string{"up"})
	tt.False(t, down)
	tt.Equal(t, 0, len(mods))

	down, mods = getKeyDown([]string{"down", "ctrl", "shift"})
	tt.True(t, down)
	tt.Equal(t, []string{"ctrl", "shift"}, mods)

	// only the first element is a direction
	down, mods = getKeyDown([]string{"ctrl", "up"})
	tt.True(t, down)
	tt.Equal(t, []string{"ctrl", "up"}, mods)

	tt.Equal(t, "down", getDown(true))
	tt.Equal(t, "up", getDown(false))
}

func TestAppendShift(t *testing.T) {
	k, args := appendShift("a", 0)
	tt.Equal(t, "a", k)
	tt.Equal(t, 0, len(args))

	// upper case adds shift and lowers the key
	k, args = appendShift("A", 0)
	tt.Equal(t, "a", k)
	tt.Equal(t, []any{"shift"}, args)

	// named keys are lowered without shift
	k, args = appendShift("Enter", 0, "ctrl")
	tt.Equal(t, "enter", k)
	tt.Equal(t, []any{"ctrl"}, args)

	// shifted symbols become base key + shift
	k, args = appendShift("!", 0)
	tt.Equal(t, "1", k)
	tt.Equal(t, []any{"shift"}, args)

	// KeyToggle("!", "up"): n=1 accounts for the direction arg
	k, args = appendShift("!", 1, "up")
	tt.Equal(t, "1", k)
	tt.Equal(t, []any{"up", "shift"}, args)

	// explicit modifiers beyond n are trusted as-is
	k, args = appendShift("!", 0, "ctrl")
	tt.Equal(t, "1", k)
	tt.Equal(t, []any{"ctrl"}, args)

	k, args = appendShift("", 0)
	tt.Equal(t, "", k)
	tt.Equal(t, 0, len(args))
}

func TestToggleArgs(t *testing.T) {
	pid, arr := getToggleArgs(1.5, true, nil, "ctrl", int64(9), 7)
	tt.Equal(t, 7, pid)
	tt.Equal(t, []string{"ctrl"}, arr)

	pid, arr = getToggleArgs()
	tt.Equal(t, 0, pid)
	tt.Equal(t, 0, len(arr))

	// empty array contributes nothing
	pid, arr = getToggleArgs([]string{}, "alt")
	tt.Equal(t, 0, pid)
	tt.Equal(t, []string{"alt"}, arr)

	// KeyUp("a", pid, "ctrl") prepends "up"; the pid must still be found
	pid, arr = getToggleArgs("up", 123, "ctrl")
	tt.Equal(t, 123, pid)
	tt.Equal(t, []string{"up", "ctrl"}, arr)

	// only the first int is the pid
	pid, arr = getToggleArgs([]string{"ctrl"}, 123, 456)
	tt.Equal(t, 123, pid)
	tt.Equal(t, []string{"ctrl"}, arr)
}
