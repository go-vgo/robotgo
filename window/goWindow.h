// Copyright 2016-2017 The go-vgo Project Developers. See the COPYRIGHT
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

int show_alert(const char *title, const char *msg, const char *defaultButton,
              const char *cancelButton){
	int alert = showAlert(title, msg, defaultButton, cancelButton);

	return alert;
}

bool is_valid(){
	bool abool = IsValid();
	return abool;
}

// int find_window(char* name){
// 	int z = findwindow(name);
// 	return z;
// }

void close_window(void){
	CloseWin();
}

bool set_handle(uintptr handle){
	bool hwnd = setHandle(handle);
	return hwnd;
}

uintptr get_handle(){
	uintptr hwnd = getHandle();
	return hwnd;
}

uintptr bget_handle(){
	MData mData = GetActive();
	#if defined(IS_MACOSX)
		return (uintptr)mData.CgID;
	#elif defined(USE_X11)
		return (uintptr)mData.XWin;
	#elif defined(IS_WINDOWS)
		return (uintptr)mData.HWnd;
	#endif
}

void set_active(const MData win){
	SetActive(win);
}

#if defined(IS_WINDOWS)
	typedef struct{
	    HWND hWnd;
	    DWORD dwPid;
	}WNDINFO;

	BOOL CALLBACK EnumWindowsProc(HWND hWnd, LPARAM lParam){
	    WNDINFO* pInfo = (WNDINFO*)lParam;
	    DWORD dwProcessId = 0;
	    GetWindowThreadProcessId(hWnd, &dwProcessId);

	    if(dwProcessId == pInfo->dwPid){
	        pInfo->hWnd = hWnd;
	        return FALSE;
	    }
	    return TRUE;
	}

	HWND GetHwndByPId(DWORD dwProcessId){
	    WNDINFO info = {0};
	    info.hWnd = NULL;
	    info.dwPid = dwProcessId;
	    EnumWindows(EnumWindowsProc, (LPARAM)&info);
	    // printf("%d\n", info.hWnd);
	    return info.hWnd;
	}
#endif

void active_PID(uintptr pid){
	MData win;
	#if defined(IS_MACOSX)
		// Handle to a AXUIElementRef
		win.AxID = AXUIElementCreateApplication(pid);
	#elif defined(USE_X11)
		win.XWin = (Window)pid;		// Handle to an X11 window
	#elif defined(IS_WINDOWS)
		// win.HWnd = (HWND)pid;		// Handle to a window HWND
		win.HWnd = GetHwndByPId(pid);
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
