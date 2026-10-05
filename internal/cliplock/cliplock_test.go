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

//go:build darwin || linux || freebsd || netbsd || openbsd || dragonfly || windows
// +build darwin linux freebsd netbsd openbsd dragonfly windows

package cliplock

import (
	"path/filepath"
	"testing"
	"time"
)

// A second holder (separate file handle, as in another process) must wait
// until the first releases the lock.
func TestAcquireExclusive(t *testing.T) {
	path := filepath.Join(t.TempDir(), "clip.lock")
	unlock, err := acquire(path)
	if err != nil {
		t.Fatal(err)
	}

	got := make(chan func() error)
	go func() {
		u, err := acquire(path)
		if err != nil {
			t.Error(err)
			u = func() error { return nil }
		}
		got <- u
	}()

	select {
	case <-got:
		t.Fatal("second acquire did not block while the lock was held")
	case <-time.After(200 * time.Millisecond):
	}

	if err := unlock(); err != nil {
		t.Fatal(err)
	}
	select {
	case u := <-got:
		if err := u(); err != nil {
			t.Fatal(err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("second acquire still blocked after unlock")
	}
}

func TestLockCleanup(t *testing.T) {
	t.Run("hold", func(t *testing.T) { Lock(t) })
	// Released by the subtest's cleanup, so this must not deadlock.
	Lock(t)
}
