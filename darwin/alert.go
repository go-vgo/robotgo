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
	"sync"

	"github.com/ebitengine/purego"
)

const (
	kCFUserNotificationNoteAlertLevel  = 1
	kCFUserNotificationDefaultResponse = 0
)

var (
	alertOnce sync.Once
	// CFUserNotificationDisplayAlert, the call the Cgo backend uses.
	cfUserNotificationDisplayAlert func(timeout float64, flags uint64,
		iconURL, soundURL, localizationURL, header, message,
		defaultBtn, alternateBtn, otherBtn uintptr, response *uint64) int32
)

// loadAlert resolves CFUserNotificationDisplayAlert lazily so a missing
// symbol only disables Alert.
func loadAlert() bool {
	alertOnce.Do(func() {
		if !loaded {
			return
		}
		cf, err := purego.Dlopen(
			"/System/Library/Frameworks/CoreFoundation.framework/CoreFoundation",
			purego.RTLD_NOW|purego.RTLD_GLOBAL)
		if err != nil {
			return
		}
		if _, err := purego.Dlsym(cf, "CFUserNotificationDisplayAlert"); err != nil {
			return
		}
		purego.RegisterLibFunc(&cfUserNotificationDisplayAlert, cf, "CFUserNotificationDisplayAlert")
	})
	return cfUserNotificationDisplayAlert != nil
}

// cfStr creates a CFString (release with cfRelease); 0 for "" so the button
// is omitted.
func cfStr(s string) uintptr {
	if s == "" {
		return 0
	}
	return cfStringCreateWithCString(0, s, cfStringEncodingUTF8)
}

// Alert shows a modal alert with the ok and cancel buttons and reports
// whether ok was chosen; false on failure.
func Alert(title, msg, ok, cancel string) bool {
	if !loadAlert() {
		return false
	}
	refs := []uintptr{cfStr(title), cfStr(msg), cfStr(ok), cfStr(cancel)}
	defer func() {
		for _, r := range refs {
			if r != 0 {
				cfRelease(r)
			}
		}
	}()

	var resp uint64
	if cfUserNotificationDisplayAlert(0, kCFUserNotificationNoteAlertLevel,
		0, 0, 0, refs[0], refs[1], refs[2], refs[3], 0, &resp) != 0 {
		return false
	}
	return resp&0x3 == kCFUserNotificationDefaultResponse
}
