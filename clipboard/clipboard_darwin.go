// Copyright 2013 @atotto. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build darwin
// +build darwin

package clipboard

import (
	"os"
	"os/exec"
	"strings"
)

var (
	pasteCmdArgs = "pbpaste"
	copyCmdArgs  = "pbcopy"
)

// utf8Env forces pbcopy/pbpaste to use UTF-8, they fall back to a legacy
// encoding when the process has no UTF-8 locale (e.g. launched from a GUI)
func utf8Env() []string {
	return append(os.Environ(), "LANG=en_US.UTF-8")
}

func getPasteCommand() *exec.Cmd {
	cmd := exec.Command(pasteCmdArgs)
	cmd.Env = utf8Env()
	return cmd
}

func getCopyCommand() *exec.Cmd {
	cmd := exec.Command(copyCmdArgs)
	cmd.Env = utf8Env()
	return cmd
}

func readAll() (string, error) {
	pasteCmd := getPasteCommand()
	out, err := pasteCmd.Output()
	if err != nil {
		return "", err
	}
	return string(out), nil
}

func writeAll(text string) error {
	copyCmd := getCopyCommand()
	copyCmd.Stdin = strings.NewReader(text)
	return copyCmd.Run()
}
