// Copyright 2016 The go-vgo Project Developers. See the COPYRIGHT
// file at the top-level directory of this distribution and at
// https://github.com/go-vgo/robotgo/blob/master/LICENSE
//
// Licensed under the Apache License, Version 2.0 <LICENSE-APACHE or
// http://www.apache.org/licenses/LICENSE-2.0> or the MIT license
// <LICENSE-MIT or http://opensource.org/licenses/MIT>, at your
// option. This file may not be copied, modified, or distributed
// except according to those terms.

// #include "../base/os.h"

#if defined(IS_WINDOWS)
	typedef struct{
	    HWND hWnd;
	    DWORD dwPid;
	}WNDINFO;

	BOOL CALLBACK EnumWindowsProc(HWND hWnd, LPARAM lParam){
	    WNDINFO* pInfo = (WNDINFO*)lParam;
	    DWORD dwProcessId = 0;
	    GetWindowThreadProcessId(hWnd, &dwProcessId);

	    if (dwProcessId == pInfo->dwPid) {
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

    // window
    void win_min(HWND hwnd, bool state){
        if (state) {
            ShowWindow(hwnd, SW_MINIMIZE);
        } else {
            ShowWindow(hwnd, SW_RESTORE);
        }
    }

    void win_max(HWND hwnd, bool state){
        if (state) {
            ShowWindow(hwnd, SW_MAXIMIZE);
        } else {
            ShowWindow(hwnd, SW_RESTORE);
        }
    }
#endif