// Copyright 2016 The go-vgo Project Developers. See the COPYRIGHT
// file at the top-level directory of this distribution and at
// https://github.com/go-vgo/robotgo/blob/master/LICENSE
//
// Licensed under the Apache License, Version 2.0 <LICENSE-APACHE or
// http://www.apache.org/licenses/LICENSE-2.0> or the MIT license
// <LICENSE-MIT or http://opensource.org/licenses/MIT>, at your
// option. This file may not be copied, modified, or distributed
// except according to those terms.

package robotgo

import ps "github.com/vcaesar/gops"

// Nps process struct
type Nps struct {
	Pid  int
	Name string
}

// Pids get the all process id
func Pids() ([]int, error) {
	return ps.Pids()
}

// PidExists determine whether the process exists
func PidExists(pid int) (bool, error) {
	return ps.PidExists(pid)
}

// Process get the all process struct
func Process() ([]Nps, error) {
	var npsArr []Nps
	nps, err := ps.Process()
	for i := 0; i < len(nps); i++ {
		np := Nps{
			nps[i].Pid,
			nps[i].Name,
		}

		npsArr = append(npsArr, np)
	}

	return npsArr, err
}

// FindName find the process name by the process id
func FindName(pid int) (string, error) {
	return ps.FindName(pid)
}

// FindNames find the all process name
func FindNames() ([]string, error) {
	return ps.FindNames()
}

// FindIds finds the all processes named with a subset
// of "name" (case insensitive),
// return matched IDs.
func FindIds(name string) ([]int, error) {
	return ps.FindIds(name)
}

// FindPath find the process path by the process pid
func FindPath(pid int) (string, error) {
	return ps.FindPath(pid)
}

// Run run a cmd shell
func Run(path string) ([]byte, error) {
	return ps.Run(path)
}

// Kill kill the process by PID
func Kill(pid int) error {
	return ps.Kill(pid)
}
