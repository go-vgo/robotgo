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
#include "win32.h"

int show_alert(const char *title, const char *msg,
	const char *defaultButton, const char *cancelButton){
	int alert = showAlert(title, msg, defaultButton, cancelButton);

	return alert;
}

intptr scalex(){
	#if defined(IS_MACOSX)
		return 0;
	#elif defined(USE_X11)
		return 0;
	#elif defined(IS_WINDOWS)
		// Get desktop dc
		HDC desktopDc = GetDC(NULL);
		// Get native resolution
		intptr horizontalDPI = GetDeviceCaps(desktopDc, LOGPIXELSX);
		// intptr verticalDPI = GetDeviceCaps(desktopDc, LOGPIXELSY);
		return horizontalDPI;
	#endif
}

intptr scaley(){
	#if defined(IS_MACOSX)
		return 0;
	#elif defined(USE_X11)
		return 0;
	#elif defined(IS_WINDOWS)
		// Get desktop dc
		HDC desktopDc = GetDC(NULL);
		// Get native resolution
		intptr verticalDPI = GetDeviceCaps(desktopDc, LOGPIXELSY);
		return verticalDPI;
	#endif
}

bool is_valid(){
	bool abool = IsValid();
	return abool;
}

// int find_window(char* name){
// 	int z = findwindow(name);
// 	return z;
// }

void min_window(uintptr pid, bool state, uintptr isHwnd){
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
		if (isHwnd == 0) {
			HWND hwnd = GetHwndByPId(pid);
			win_min(hwnd, state);
		} else {
			win_min((HWND)pid, state);
		}
	#endif
}

void max_window(uintptr pid, bool state, uintptr isHwnd){
	#if defined(IS_MACOSX)
		// return 0;
	#elif defined(USE_X11)
		XDismissErrors();
		// SetState((Window)pid, STATE_MINIMIZE, false);
		// SetState((Window)pid, STATE_MAXIMIZE, state);
	#elif defined(IS_WINDOWS)
		if (isHwnd == 0) {
			HWND hwnd = GetHwndByPId(pid);
			win_max(hwnd, state);
		} else {
			win_max((HWND)pid, state);
		}
	#endif
}

void close_window(void){
	CloseWin();
}

bool set_handle(uintptr handle){
	bool hwnd = setHandle(handle);
	return hwnd;
}

uintptr get_handle(){
	MData mData = GetActive();
	#if defined(IS_MACOSX)
		return (uintptr)mData.CgID;
	#elif defined(USE_X11)
		return (uintptr)mData.XWin;
	#elif defined(IS_WINDOWS)
		return (uintptr)mData.HWnd;
	#endif
}

uintptr bget_handle(){
	uintptr hwnd = getHandle();
	return hwnd;
}

void set_active(const MData win){
	SetActive(win);
}

void active_PID(uintptr pid, uintptr isHwnd){
	MData win;
	#if defined(IS_MACOSX)
		// Handle to a AXUIElementRef
		win.AxID = AXUIElementCreateApplication(pid);
	#elif defined(USE_X11)
		win.XWin = (Window)pid;		// Handle to an X11 window
	#elif defined(IS_WINDOWS)
		// win.HWnd = (HWND)pid;		// Handle to a window HWND
		if (isHwnd == 0) {
			win.HWnd = GetHwndByPId(pid);
		} else {
			win.HWnd = (HWND)pid;
		}
	#endif

	SetActive(win);
}

MData get_active(){
	MData mdata = GetActive();
	return mdata;
}

char* get_title(){
	char* title = GetTitle();
	// printf("title::::%s\n", title );
	return title;
}

int32 get_PID(void){
	int pid = WGetPID();
	return pid;
}
