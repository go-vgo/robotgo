//go:build darwin
// +build darwin

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

package darwin

import (
	"errors"
	"os"
	"os/exec"
)

// GetPid returns the current process's PID.
func GetPid() int {
	return os.Getpid()
}

// Kill kills a process by PID.
func Kill(pid int) error {
	if pid <= 0 {
		// A zero or negative pid would signal the whole process group.
		return errors.New("robotgo: invalid pid")
	}
	p, err := os.FindProcess(pid)
	if err != nil {
		return err
	}
	return p.Kill()
}

// Run runs a shell command and returns its combined output.
func Run(path string) ([]byte, error) {
	return exec.Command("/bin/sh", "-c", path).CombinedOutput()
}
