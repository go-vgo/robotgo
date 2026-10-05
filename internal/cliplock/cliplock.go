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
// finishes. The OS drops the lock if the process dies.
func Lock(tb testing.TB) {
	tb.Helper()
	unlock, err := acquire(filepath.Join(os.TempDir(), "robotgo-clipboard.lock"))
	if err != nil {
		tb.Fatal(err)
	}
	tb.Cleanup(func() {
		if err := unlock(); err != nil {
			tb.Error(err)
		}
	})
}

func acquire(path string) (func() error, error) {
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0o666)
	if err != nil {
		return nil, err
	}
	if err := lockFile(f); err != nil {
		return nil, errors.Join(err, f.Close())
	}
	return func() error { return errors.Join(unlockFile(f), f.Close()) }, nil
}
