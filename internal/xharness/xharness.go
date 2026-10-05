//go:build linux && cgo
// +build linux,cgo

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

// Package xharness is a test-only helper: it owns a focused X11 window and
// reads back the text that window receives, through XLookupString, so a test
// sees exactly what an application would be given for injected key events
// (layout, Shift and AltGr levels applied by the server).
package xharness

/*
#cgo LDFLAGS: -lX11
#include <string.h>
#include <X11/Xlib.h>
#include <X11/Xutil.h>

static Display *hdpy;
static Window hwin;

static int harnessOpen(void) {
	hdpy = XOpenDisplay(NULL);
	if (hdpy == NULL) {
		return 0;
	}
	hwin = XCreateSimpleWindow(hdpy, DefaultRootWindow(hdpy), 0, 0, 200, 100, 0, 0, 0);
	XSelectInput(hdpy, hwin, KeyPressMask | StructureNotifyMask);
	XMapWindow(hdpy, hwin);
	XEvent ev;
	do {
		XNextEvent(hdpy, &ev);
	} while (ev.type != MapNotify);
	XSetInputFocus(hdpy, hwin, RevertToParent, CurrentTime);
	XSync(hdpy, False);
	return 1;
}

static int harnessRead(char *buf, int n) {
	int len = 0;
	XEvent ev;
	XSync(hdpy, False);
	while (XCheckWindowEvent(hdpy, hwin, KeyPressMask, &ev)) {
		KeySym ks;
		char tmp[8];
		int k = XLookupString(&ev.xkey, tmp, sizeof tmp, &ks, NULL);
		if (k > 0 && len + k < n) {
			memcpy(buf + len, tmp, k);
			len += k;
		}
	}
	buf[len] = 0;
	return len;
}

static void harnessClose(void) {
	if (hdpy != NULL) {
		XDestroyWindow(hdpy, hwin);
		XCloseDisplay(hdpy);
		hdpy = NULL;
	}
}
*/
import "C"

import (
	"errors"
	"unsafe"
)

// Open maps a window on $DISPLAY and gives it the input focus.
func Open() error {
	if C.harnessOpen() == 0 {
		return errors.New("xharness: cannot open X display")
	}
	return nil
}

// Read drains the KeyPress events the window received into text.
func Read() string {
	buf := make([]byte, 256)
	n := C.harnessRead((*C.char)(unsafe.Pointer(&buf[0])), C.int(len(buf)))
	return string(buf[:n])
}

// Close destroys the window and the display connection.
func Close() { C.harnessClose() }
