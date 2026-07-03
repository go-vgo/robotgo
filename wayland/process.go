//go:build linux
// +build linux

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

package wayland

import (
	"errors"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Process management — pure Go, no Wayland dependency.

// Pids returns all process IDs by reading /proc.
func Pids() ([]int, error) {
	entries, err := os.ReadDir("/proc")
	if err != nil {
		return nil, err
	}
	var pids []int
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		pid, err := strconv.Atoi(e.Name())
		if err != nil {
			continue
		}
		pids = append(pids, pid)
	}
	return pids, nil
}

// PidExists checks if a process exists.
func PidExists(pid int) (bool, error) {
	_, err := os.Stat("/proc/" + strconv.Itoa(pid))
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// FindName returns the process name for a given PID.
func FindName(pid int) (string, error) {
	data, err := os.ReadFile("/proc/" + strconv.Itoa(pid) + "/comm")
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(data)), nil
}

// FindNames returns all process names.
func FindNames() ([]string, error) {
	pids, err := Pids()
	if err != nil {
		return nil, err
	}
	var names []string
	for _, pid := range pids {
		name, err := FindName(pid)
		if err != nil {
			continue
		}
		names = append(names, name)
	}
	return names, nil
}

// FindIds finds PIDs by name substring (case-insensitive).
func FindIds(name string) ([]int, error) {
	pids, err := Pids()
	if err != nil {
		return nil, err
	}
	nameLower := strings.ToLower(name)
	var result []int
	for _, pid := range pids {
		pname, err := FindName(pid)
		if err != nil {
			continue
		}
		if strings.Contains(strings.ToLower(pname), nameLower) {
			result = append(result, pid)
		}
	}
	return result, nil
}

// FindPath returns the executable path for a PID.
func FindPath(pid int) (string, error) {
	return os.Readlink("/proc/" + strconv.Itoa(pid) + "/exe")
}

// Process returns all processes as []Nps.
func Process() ([]Nps, error) {
	pids, err := Pids()
	if err != nil {
		return nil, err
	}
	var procs []Nps
	for _, pid := range pids {
		name, err := FindName(pid)
		if err != nil {
			continue
		}
		procs = append(procs, Nps{Pid: pid, Name: name})
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

// Run runs a shell command and returns its output.
func Run(path string) ([]byte, error) {
	return exec.Command("sh", "-c", path).CombinedOutput()
}
