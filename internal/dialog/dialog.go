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

// Package dialog shows a modal OK/Cancel alert through a desktop helper
// program (zenity, else xmessage), shared by the pure-Go Linux backends.
package dialog

import "os/exec"

// command builds the alert command line with the first helper lookPath finds;
// nil when none is installed.
func command(lookPath func(string) (string, error), title, msg, ok, cancel string) []string {
	if p, err := lookPath("zenity"); err == nil {
		return []string{p, "--question", "--title", title, "--text", msg,
			"--ok-label", ok, "--cancel-label", cancel}
	}
	if p, err := lookPath("xmessage"); err == nil {
		return []string{p, "-center", "-title", title,
			"-buttons", ok + ":0," + cancel + ":1", "-default", ok, msg}
	}
	return nil
}

// Show shows the alert and reports whether the ok button was chosen; false
// when it was canceled or no helper program is available.
func Show(title, msg, ok, cancel string) bool {
	args := command(exec.LookPath, title, msg, ok, cancel)
	if args == nil {
		return false
	}
	// both helpers exit 0 for the ok button only
	return exec.Command(args[0], args[1:]...).Run() == nil
}
