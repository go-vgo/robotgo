#if defined(IS_WINDOWS)
    bool CALLBACK MonitorEnumProc(HMONITOR hMonitor, HDC hdcMonitor, LPRECT lprcMonitor, LPARAM dwData) {
        uint32_t *count = (uint32_t*)dwData;
        (*count)++;
        return true;
    }
   
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

	HWND GetHwndByPId(DWORD dwProcessId) {
	    WNDINFO info = {0};
	    info.hWnd = NULL;
	    info.dwPid = dwProcessId;
	    EnumWindows(EnumWindowsProc, (LPARAM)&info);
	    // printf("%d\n", info.hWnd);
	    return info.hWnd;
	}

#endif