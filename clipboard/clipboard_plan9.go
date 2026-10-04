// Copyright 2013 @atotto. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build plan9
// +build plan9

package clipboard

import (
	"io"
	"os"
)

const snarf = "/dev/snarf"

func readAll() (string, error) {
	f, err := os.Open(snarf)
	if err != nil {
		return "", err
	}
	defer f.Close()

	out, err := io.ReadAll(f)
	if err != nil {
		return "", err
	}
	return string(out), nil
}

func writeAll(text string) error {
	f, err := os.OpenFile(snarf, os.O_WRONLY|os.O_TRUNC, 0)
	if err != nil {
		return err
	}

	if _, err := f.WriteString(text); err != nil {
		f.Close()
		return err
	}
	return f.Close()
}
