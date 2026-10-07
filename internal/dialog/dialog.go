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

import (
	"os/exec"
	"strings"
)

// xmessageEscape escapes the xmessage -buttons separators "," and ":" and
// its escape char "\\", so a label never splits or remaps the exit codes.
var xmessageEscape = strings.NewReplacer(`\`, `\\`, `,`, `\,`, `:`, `\:`).Replace

// command builds the alert command line with the first helper lookPath finds;
// an empty cancel omits the cancel button, nil when no helper is installed.
// Labels and the message are kept as plain text: zenity gets --no-markup and
// xmessage reads msg from stdin (see Show) so a leading "-" is not an option.
func command(lookPath func(string) (string, error), title, msg, ok, cancel string) []string {
	if p, err := lookPath("zenity"); err == nil {
		if cancel == "" {
			return []string{p, "--info", "--no-markup", "--title", title, "--text", msg,
				"--ok-label", ok}
		}
		return []string{p, "--question", "--no-markup", "--title", title, "--text", msg,
			"--ok-label", ok, "--cancel-label", cancel}
	}
	if p, err := lookPath("xmessage"); err == nil {
		buttons := xmessageEscape(ok) + ":0"
		if cancel != "" {
			buttons += "," + xmessageEscape(cancel) + ":1"
		}
		// -default matches the unescaped button name; -geometry as the old Cgo Alert
		return []string{p, "-center", "-title", title,
			"-buttons", buttons, "-default", ok, "-geometry", "400x200", "-file", "-"}
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
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdin = strings.NewReader(msg) // xmessage -file -; zenity ignores it
	// both helpers exit 0 for the ok button only
	return cmd.Run() == nil
}
