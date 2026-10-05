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

//go:build darwin && !mac && !purego
// +build darwin,!mac,!purego

package robotgo

/*
#include "base/types.h"
#include <CoreGraphics/CoreGraphics.h>

// Interval a source created by MMEventSourceCreate (or, with mm == 0, a
// plain HID-state source) suppresses the user's own input for after each
// posted event; -1 if no source could be created.
static double eventSourceSuppression(int mm) {
	CGEventSourceRef src = mm ? MMEventSourceCreate()
		: CGEventSourceCreate(kCGEventSourceStateHIDSystemState);
	if (src == NULL) { return -1; }
	double v = CGEventSourceGetLocalEventsSuppressionInterval(src);
	CFRelease(src);
	return v;
}
*/
import "C"

// GetMainId get the main display id
func GetMainId() int {
	return int(C.CGMainDisplayID())
}

// eventSourceSuppression reports the local-events suppression interval of a
// Cgo event source (see MMEventSourceCreate in base/types.h).
func eventSourceSuppression(mm bool) float64 {
	v := 0
	if mm {
		v = 1
	}
	return float64(C.eventSourceSuppression(C.int(v)))
}
