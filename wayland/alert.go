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

import "github.com/go-vgo/robotgo/internal/dialog"

// Alert shows a modal alert with the ok and cancel buttons (via zenity, else
// xmessage) and reports whether ok was chosen; false when no helper exists.
func Alert(title, msg, ok, cancel string) bool {
	return dialog.Show(title, msg, ok, cancel)
}
