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

// Package cliplock serializes system clipboard access between test
// processes: go test ./... runs packages in parallel and they all share
// one clipboard.
package cliplock

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

// Lock blocks until tb holds the clipboard lock; it is released when tb
// finishes. The OS drops the lock if the process dies. It is not reentrant:
// calling Lock again before tb finishes (e.g. in a subtest) deadlocks.
func Lock(tb testing.TB) {
	tb.Helper()
	path, err := lockPath()
	if err != nil {
		tb.Fatal(err)
	}
	unlock, err := acquire(path)
	if err != nil {
		tb.Fatal(err)
	}
	tb.Cleanup(func() {
		if err := unlock(); err != nil {
			tb.Error(err)
		}
	})
}

// lockPath returns the lock file in a private per-user directory: in a
// shared /tmp another user could pre-create the file and hold the lock.
func lockPath() (string, error) {
	base, err := os.UserCacheDir()
	if err != nil {
		return "", err
	}
	dir := filepath.Join(base, "robotgo")
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return "", err
	}
	// MkdirAll keeps the mode of an existing directory.
	if err := os.Chmod(dir, 0o700); err != nil {
		return "", err
	}
	return filepath.Join(dir, "clipboard.lock"), nil
}

func acquire(path string) (func() error, error) {
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0o600)
	if err != nil {
		return nil, err
	}
	if err := lockFile(f); err != nil {
		return nil, errors.Join(err, f.Close())
	}
	return func() error { return errors.Join(unlockFile(f), f.Close()) }, nil
}
