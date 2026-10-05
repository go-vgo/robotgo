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

// Untagged on purpose: headless tests for the portable helpers in
// robotgo_pub.go and ps.go, run by the Cgo and every pure-Go backend job.

package robotgo

import (
	"os"
	"slices"
	"strings"
	"testing"

	"github.com/vcaesar/tt"
)

func TestGetVersionPub(t *testing.T) {
	tt.Equal(t, Version, GetVersion())
	tt.True(t, strings.HasPrefix(GetVersion(), "v2."))
}

func TestTry(t *testing.T) {
	var got any
	Try(func() { panic("boom") }, func(e any) { got = e })
	tt.Equal(t, "boom", got)

	called, ran := false, false
	Try(func() { ran = true }, func(any) { called = true })
	tt.True(t, ran)
	tt.False(t, called)
}

func TestScaled01(t *testing.T) {
	tests := []struct {
		x        int
		f        float64
		mul, div int
	}{
		{100, 1, 100, 100},
		{100, 2, 200, 50},
		{101, 1.5, 151, 67},
		{-10, 2, -20, -5},
		{0, 1.25, 0, 0},
	}
	for _, c := range tests {
		tt.Equal(t, c.mul, Scaled0(c.x, c.f))
		tt.Equal(t, c.div, Scaled1(c.x, c.f))
	}
}

func TestRectEmbedding(t *testing.T) {
	r := Rect{Point{1, 2}, Size{3, 4}}
	tt.Equal(t, 1, r.X)
	tt.Equal(t, 2, r.Y)
	tt.Equal(t, 3, r.W)
	tt.Equal(t, 4, r.H)
}

func TestToStringsPanicsOnNonString(t *testing.T) {
	var got any
	Try(func() { ToStrings([]any{"a", 1}) }, func(e any) { got = e })
	tt.NotNil(t, got)
}

func TestProcessSelf(t *testing.T) {
	self := os.Getpid()

	ids, err := Pids()
	tt.Nil(t, err)
	tt.True(t, slices.Contains(ids, self))

	procs, err := Process()
	tt.Nil(t, err)
	tt.True(t, slices.ContainsFunc(procs, func(p Nps) bool { return p.Pid == self }))

	name, err := FindName(self)
	tt.Nil(t, err)
	tt.NotEmpty(t, name)

	names, err := FindNames()
	tt.Nil(t, err)
	tt.True(t, slices.Contains(names, name))

	found, err := FindIds(name)
	tt.Nil(t, err)
	tt.True(t, slices.Contains(found, self))

	path, err := FindPath(self)
	tt.Nil(t, err)
	tt.NotEmpty(t, path)
}

func TestPidExistsUnknown(t *testing.T) {
	// Pids are positive; a huge one is never allocated.
	ok, err := PidExists(1 << 30)
	tt.Nil(t, err)
	tt.False(t, ok)
}
