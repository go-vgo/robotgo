// Copyright 2013 @atotto. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
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
