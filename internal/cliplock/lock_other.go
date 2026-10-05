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

//go:build !darwin && !linux && !freebsd && !netbsd && !openbsd && !dragonfly && !windows
// +build !darwin,!linux,!freebsd,!netbsd,!openbsd,!dragonfly,!windows

package cliplock

import "os"

// Platforms without flock/LockFileEx (plan9, solaris, ...) do not serialize.
func lockFile(f *os.File) error { return nil }

func unlockFile(f *os.File) error { return nil }
