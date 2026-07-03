//go:build windows
// +build windows

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

package win

import (
	"errors"
	"os"
	"os/exec"
	"strings"
	"unsafe"

	"golang.org/x/sys/windows"
)

// snapshotProcesses walks the process snapshot, calling fn for each entry.
// fn returns false to stop early.
func snapshotProcesses(fn func(e *windows.ProcessEntry32) bool) error {
	snap, err := windows.CreateToolhelp32Snapshot(windows.TH32CS_SNAPPROCESS, 0)
	if err != nil {
		return err
	}
	defer windows.CloseHandle(snap)

	var e windows.ProcessEntry32
	e.Size = uint32(unsafe.Sizeof(e))

	err = windows.Process32First(snap, &e)
	for err == nil {
		if !fn(&e) {
			return nil
		}
		err = windows.Process32Next(snap, &e)
	}
	if err == windows.ERROR_NO_MORE_FILES {
		return nil
	}
	return err
}

// Pids returns all process IDs.
func Pids() ([]int, error) {
	var pids []int
	err := snapshotProcesses(func(e *windows.ProcessEntry32) bool {
		pids = append(pids, int(e.ProcessID))
		return true
	})
	if err != nil {
		return nil, err
	}
	return pids, nil
}

// PidExists reports whether a process with the given PID exists.
func PidExists(pid int) (bool, error) {
	exists := false
	err := snapshotProcesses(func(e *windows.ProcessEntry32) bool {
		if int(e.ProcessID) == pid {
			exists = true
			return false
		}
		return true
	})
	if err != nil {
		return false, err
	}
	return exists, nil
}

// FindName returns the executable name for a given PID.
func FindName(pid int) (string, error) {
	name := ""
	found := false
	err := snapshotProcesses(func(e *windows.ProcessEntry32) bool {
		if int(e.ProcessID) == pid {
			name = windows.UTF16ToString(e.ExeFile[:])
			found = true
			return false
		}
		return true
	})
	if err != nil {
		return "", err
	}
	if !found {
		return "", ErrNotFound
	}
	return name, nil
}

// FindNames returns all process names.
func FindNames() ([]string, error) {
	var names []string
	err := snapshotProcesses(func(e *windows.ProcessEntry32) bool {
		names = append(names, windows.UTF16ToString(e.ExeFile[:]))
		return true
	})
	if err != nil {
		return nil, err
	}
	return names, nil
}

// FindIds finds PIDs whose executable name contains name (case insensitive).
func FindIds(name string) ([]int, error) {
	nameLower := strings.ToLower(name)
	var result []int
	err := snapshotProcesses(func(e *windows.ProcessEntry32) bool {
		ename := strings.ToLower(windows.UTF16ToString(e.ExeFile[:]))
		if strings.Contains(ename, nameLower) {
			result = append(result, int(e.ProcessID))
		}
		return true
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

// FindPath returns the full executable path for a PID.
func FindPath(pid int) (string, error) {
	h, err := windows.OpenProcess(windows.PROCESS_QUERY_LIMITED_INFORMATION, false, uint32(pid))
	if err != nil {
		return "", err
	}
	defer windows.CloseHandle(h)

	buf := make([]uint16, windows.MAX_PATH)
	size := uint32(len(buf))
	if err := windows.QueryFullProcessImageName(h, 0, &buf[0], &size); err != nil {
		return "", err
	}
	return windows.UTF16ToString(buf[:size]), nil
}

// Process returns all processes as []Nps.
func Process() ([]Nps, error) {
	var procs []Nps
	err := snapshotProcesses(func(e *windows.ProcessEntry32) bool {
		procs = append(procs, Nps{
			Pid:  int(e.ProcessID),
			Name: windows.UTF16ToString(e.ExeFile[:]),
		})
		return true
	})
	if err != nil {
		return nil, err
	}
	return procs, nil
}

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
	return exec.Command("cmd", "/c", path).CombinedOutput()
}
