// Copyright 2016 The go-vgo Project Developers. See the COPYRIGHT
// file at the top-level directory of this distribution and at
// https://github.com/go-vgo/robotgo/blob/master/LICENSE
//
// Licensed under the Apache License, Version 2.0 <LICENSE-APACHE or
// http://www.apache.org/licenses/LICENSE-2.0> or the MIT license
// <LICENSE-MIT or http://opensource.org/licenses/MIT>, at your
// option. This file may not be copied, modified, or distributed
// except according to those terms.

#include "alert_c.h"
#include "window.h"
#include "win_sys.h"

void min_window(uintptr pid, bool state, int8_t isPid){
	#if defined(IS_MACOSX)
		// return 0;
		AXUIElementRef axID = AXUIElementCreateApplication(pid);
		AXUIElementSetAttributeValue(axID, kAXMinimizedAttribute, 
										state ? kCFBooleanTrue : kCFBooleanFalse);
	#elif defined(USE_X11)
		// Ignore X errors
		XDismissErrors();
		// SetState((Window)pid, STATE_MINIMIZE, state);
	#elif defined(IS_WINDOWS)
        HWND hwnd = getHwnd(pid, isPid);
		win_min(hwnd, state);
	#endif
}

void max_window(uintptr pid, bool state, int8_t isPid){
	#if defined(IS_MACOSX)
		// return 0;
	#elif defined(USE_X11)
		XDismissErrors();
		// SetState((Window)pid, STATE_MINIMIZE, false);
		// SetState((Window)pid, STATE_MAXIMIZE, state);
	#elif defined(IS_WINDOWS)
        HWND hwnd = getHwnd(pid, isPid);
		win_max(hwnd, state);
	#endif
}

uintptr get_handle(){
	MData mData = get_active();

	#if defined(IS_MACOSX)
		return (uintptr)mData.CgID;
	#elif defined(USE_X11)
		return (uintptr)mData.XWin;
	#elif defined(IS_WINDOWS)
		return (uintptr)mData.HWnd;
	#endif
}

uintptr b_get_handle() {
	#if defined(IS_MACOSX)
		return (uintptr)pub_mData.CgID;
	#elif defined(USE_X11)
		return (uintptr)pub_mData.XWin;
	#elif defined(IS_WINDOWS)
		return (uintptr)pub_mData.HWnd;
	#endif
}

void active_PID(uintptr pid, int8_t isPid){
	MData win = set_handle_pid(pid, isPid);
	set_active(win);
}
