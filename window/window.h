// Copyright 2016 The go-vgo Project Developers. See the COPYRIGHT
// file at the top-level directory of this distribution and at
// https://github.com/go-vgo/robotgo/blob/master/LICENSE
//
// Licensed under the Apache License, Version 2.0 <LICENSE-APACHE or
// http://www.apache.org/licenses/LICENSE-2.0> or the MIT license
// <LICENSE-MIT or http://opensource.org/licenses/MIT>, at your
// option. This file may not be copied, modified, or distributed
// except according to those terms.

#include "pub.h"

bool setHandle(uintptr handle);
bool is_valid();
bool IsAxEnabled(bool options);

MData get_active(void);
void initWindow(uintptr handle);
char* get_title_by_hand(MData m_data);
void close_window_by_Id(MData m_data);

// int findwindow()
uintptr initHandle = 0;

void initWindow(uintptr handle){
#if defined(IS_MACOSX)
	pub_mData.CgID = 0;
	pub_mData.AxID = 0;
#elif defined(USE_X11)
	Display *rDisplay = XOpenDisplay(NULL);
	// If atoms loaded
	if (WM_PID == None) {
		// Load all necessary atom properties
		if (rDisplay != NULL) {LoadAtoms();}
	}

	pub_mData.XWin = 0;
	XCloseDisplay(rDisplay);
#elif defined(IS_WINDOWS)
	pub_mData.HWnd = 0;
#endif
	setHandle(handle);
}

bool Is64Bit() {
	#ifdef RobotGo_64
		return true;
	#endif

	return false;
}

MData set_handle_pid(uintptr pid, int8_t isPid){
	MData win;
	#if defined(IS_MACOSX)
		// Handle to a AXUIElementRef
		win.AxID = AXUIElementCreateApplication(pid);
	#elif defined(USE_X11)
		win.XWin = (Window)pid;  // Handle to an X11 window
	#elif defined(IS_WINDOWS)
		// win.HWnd = (HWND)pid;		// Handle to a window HWND
        win.HWnd = getHwnd(pid, isPid);
	#endif

	return win;
}

void set_handle_pid_mData(uintptr pid, int8_t isPid){
	MData win = set_handle_pid(pid, isPid);
	pub_mData = win;
}

bool is_valid() {
	initWindow(initHandle);
	if (!IsAxEnabled(true)) {
		printf("%s\n", "Window: Accessibility API is disabled! "
		"Failed to enable access for assistive devices. \n");
	}
	MData actdata = get_active();

#if defined(IS_MACOSX)
	pub_mData.CgID = actdata.CgID;
	pub_mData.AxID = actdata.AxID;
	if (pub_mData.CgID == 0 || pub_mData.AxID == 0) { return false; }

	CFTypeRef r = NULL;
	// Attempt to get the window role
	if (AXUIElementCopyAttributeValue(pub_mData.AxID, kAXRoleAttribute, &r) == kAXErrorSuccess && r){
			CFRelease(r);
			return true;
	}

	return false;
#elif defined(USE_X11)
	pub_mData.XWin = actdata.XWin;
	if (pub_mData.XWin == 0) { return false; }

	Display *rDisplay = XOpenDisplay(NULL);
	// Check for a valid X-Window display
	if (rDisplay == NULL) { return false; }

	// Ignore X errors
	XDismissErrors();

	// Get the window PID property
	void* result = GetWindowProperty(pub_mData, WM_PID,NULL);
	if (result == NULL) {
		XCloseDisplay(rDisplay);
		return false;
	}

	// Free result and return true
	XFree(result);
	XCloseDisplay(rDisplay);

	return true;
#elif defined(IS_WINDOWS)
	pub_mData.HWnd = actdata.HWnd;
	if (pub_mData.HWnd == 0) {
		return false;
	}

	return IsWindow(pub_mData.HWnd) != 0;
#endif
}

bool IsAxEnabled(bool options){
#if defined(IS_MACOSX)
	// Statically load all required functions one time
	static dispatch_once_t once; dispatch_once (&once,
	^{
		// Open the framework
		void* handle = dlopen("/System/Library/Frameworks/Application" 
			"Services.framework/ApplicationServices", RTLD_LAZY);

		// Validate the handle
		if (handle != NULL) {
			*(void**) (&gAXIsProcessTrustedWithOptions) = dlsym (handle, "AXIsProcessTrustedWithOptions");
			gkAXTrustedCheckOptionPrompt = (CFStringRef*) dlsym (handle, "kAXTrustedCheckOptionPrompt");
		}
	});

	// Check for new OSX 10.9 function
	if (gAXIsProcessTrustedWithOptions) {
		// Check whether to show prompt
		CFBooleanRef displayPrompt = options ? kCFBooleanTrue : kCFBooleanFalse;

		// Convert display prompt value into a dictionary
		const void* k[] = { *gkAXTrustedCheckOptionPrompt };
		const void* v[] = { displayPrompt };
		CFDictionaryRef o = CFDictionaryCreate(NULL, k, v, 1, NULL, NULL);

		// Determine whether the process is actually trusted
		bool result = (*gAXIsProcessTrustedWithOptions)(o);
		// Free memory
		CFRelease(o);
		return result;
	} else {
		// Ignore deprecated warnings
		#pragma clang diagnostic push
		#pragma clang diagnostic ignored "-Wdeprecated-declarations"

		// Check whether we have accessibility access
		return AXAPIEnabled() || AXIsProcessTrusted();
		#pragma clang diagnostic pop
	}
#elif defined(USE_X11)
	return true;
#elif defined(IS_WINDOWS)
	return true;
#endif
}

// int
bool setHandle(uintptr handle){
#if defined(IS_MACOSX)
	// Release the AX element
	if (pub_mData.AxID != NULL) {
		CFRelease(pub_mData.AxID);
	}

	// Reset both values
	pub_mData.CgID = 0;
	pub_mData.AxID = 0;

	if (handle == 0) {
		// return 0;
		return true;
	}

	// Retrieve the window element
	CGWindowID cgID = (CGWindowID)handle;
	AXUIElementRef axID = GetUIElement(cgID);
	if (axID != NULL){
		pub_mData.CgID = cgID;
		pub_mData.AxID = axID;
		// return 0;
		return true;
	}

	// return 1;
	return false;
#elif defined(USE_X11)
	pub_mData.XWin = (Window)handle;
	if (handle == 0) {
		return true;
	}

	if (is_valid()) {
		return true;
	}

	pub_mData.XWin = 0;
	return false;
#elif defined(IS_WINDOWS)
	pub_mData.HWnd = (HWND)handle;
	if (handle == 0) {
		return true;
	}

	if (is_valid()) {
		return true;
	}

	pub_mData.HWnd = 0;
	return false;
#endif
}

bool IsTopMost(void){
	// Check the window validity
	if (!is_valid()) { return false; }
#if defined(IS_MACOSX)
	return false; // WARNING: Unavailable
#elif defined(USE_X11)
	// Ignore X errors
	// XDismissErrors ();
	// return GetState (mData.XWin, STATE_TOPMOST);
#elif defined(IS_WINDOWS)
	return (GetWindowLongPtr(pub_mData.HWnd, GWL_EXSTYLE) & WS_EX_TOPMOST) != 0;
#endif
}

bool IsMinimized(void){
	// Check the window validity
	if (!is_valid()) { return false; }
#if defined(IS_MACOSX)
	CFBooleanRef data = NULL;
	// Determine whether the window is minimized
	if (AXUIElementCopyAttributeValue(pub_mData.AxID, kAXMinimizedAttribute, 
	(CFTypeRef*) &data) == kAXErrorSuccess && data != NULL) {
		// Convert resulting data into a bool
		bool result = CFBooleanGetValue(data);
		CFRelease(data);
		return result;
	}

	return false;
#elif defined(USE_X11)
	// Ignore X errors
	// XDismissErrors();
	// return GetState(mData.XWin, STATE_MINIMIZE);
#elif defined(IS_WINDOWS)
	return (GetWindowLongPtr(pub_mData.HWnd, GWL_STYLE) & WS_MINIMIZE) != 0;
#endif
}

//////
bool IsMaximized(void){
	// Check the window validity
	if (!is_valid()) { return false; }
#if defined(IS_MACOSX)
	return false; // WARNING: Unavailable
#elif defined(USE_X11)
	// Ignore X errors
	// XDismissErrors();
	// return GetState(mData.XWin, STATE_MAXIMIZE);
#elif defined(IS_WINDOWS)
	return (GetWindowLongPtr(pub_mData.HWnd, GWL_STYLE) & WS_MAXIMIZE) != 0;
#endif
}

void set_active(const MData win) {
	// Check if the window is valid
	if (!is_valid()) { return; }
#if defined(IS_MACOSX)
	// Attempt to raise the specified window object
	if (AXUIElementPerformAction(win.AxID, kAXRaiseAction) != kAXErrorSuccess) {
		pid_t pid = 0;
		// Attempt to retrieve the PID of the window
		if (AXUIElementGetPid(win.AxID, &pid) != kAXErrorSuccess || !pid) { return; }

		// Ignore deprecated warnings
		#pragma clang diagnostic push
		#pragma clang diagnostic ignored "-Wdeprecated-declarations"

		ProcessSerialNumber psn;
		// Attempt to retrieve the process psn
		if (GetProcessForPID(pid, &psn) == 0) {
			// Gracefully activate process
			SetFrontProcessWithOptions(&psn, kSetFrontProcessFrontWindowOnly);
		}

		#pragma clang diagnostic pop
	}
#elif defined(USE_X11)
	// Ignore X errors
	XDismissErrors();

	// Go to the specified window's desktop
	SetDesktopForWindow(win);
	Display *rDisplay = XOpenDisplay(NULL);
	// Check the atom value
	if (WM_ACTIVE != None) {
		// Retrieve the screen number
		XWindowAttributes attr = { 0 };
		XGetWindowAttributes(rDisplay, win.XWin, &attr);
		int s = XScreenNumberOfScreen(attr.screen);

		// Prepare an event
		XClientMessageEvent e = { 0 };
		e.window = win.XWin;
		e.format = 32;
		e.message_type = WM_ACTIVE;
		e.display = rDisplay;
		e.type = ClientMessage;
		e.data.l[0] = 2;
		e.data.l[1] = CurrentTime;

		// Send the message
		XSendEvent(rDisplay, XRootWindow(rDisplay, s), False,
			SubstructureNotifyMask | SubstructureRedirectMask,
			(XEvent*) &e);
	} else {
		// Attempt to raise the specified window
		XRaiseWindow(rDisplay, win.XWin);
		// Set the specified window's input focus
		XSetInputFocus(rDisplay, win.XWin, RevertToParent, CurrentTime);
	}
	XCloseDisplay(rDisplay);
#elif defined(IS_WINDOWS)
	if (IsMinimized()) {
		ShowWindow(win.HWnd, SW_RESTORE);
	}

	SetForegroundWindow(win.HWnd);
#endif
}

MData get_active(void) {
#if defined(IS_MACOSX)
	MData result;
	// Ignore deprecated warnings
	#pragma clang diagnostic push
	#pragma clang diagnostic ignored "-Wdeprecated-declarations"

	ProcessSerialNumber psn; pid_t pid;
	// Attempt to retrieve the front process
	if (GetFrontProcess(&psn) != 0 || GetProcessPID(&psn, &pid) != 0) {
		return result;
	}

	#pragma clang diagnostic pop

	// Create accessibility object using focused PID
	AXUIElementRef focused = AXUIElementCreateApplication(pid);
	if (focused == NULL) { return result; } // Verify

	AXUIElementRef element;
	CGWindowID win = 0;
	// Retrieve the currently focused window
	if (AXUIElementCopyAttributeValue(focused, kAXFocusedWindowAttribute, (CFTypeRef*) &element) 
		== kAXErrorSuccess && element) {

		// Use undocumented API to get WID
		if (_AXUIElementGetWindow(element, &win) == kAXErrorSuccess && win) {
			// Manually set internals
			result.CgID = win;
			result.AxID = element;
		} else {
			CFRelease(element);
		}
	} else {
		result.CgID = win;
		result.AxID = element;
	}
	CFRelease(focused);

	return result;
#elif defined(USE_X11)
	MData result;
	Display *rDisplay = XOpenDisplay(NULL);
	// Check X-Window display
	if (WM_ACTIVE == None || rDisplay == NULL) {
		return result;
	}

	// Ignore X errors
	XDismissErrors();

	// Get the current active window
	result.XWin = XDefaultRootWindow(rDisplay);
	void* active = GetWindowProperty(result, WM_ACTIVE, NULL);

	// Check result value
	if (active != NULL) {
		// Extract window from the result
		long window = *((long*)active);
		XFree(active);

		if (window != 0) {
			// Set and return the foreground window
			result.XWin = (Window)window;
			XCloseDisplay(rDisplay);
			return result;
		}
	}

	// Use input focus instead
	Window window = None;
	int revert = RevertToNone;
	XGetInputFocus(rDisplay, &window, &revert);
	XCloseDisplay(rDisplay);

	// Return foreground window
	result.XWin = window;
	return result;
#elif defined(IS_WINDOWS)
	// Attempt to get the foreground window multiple times in case
	MData result;

	uint8_t times = 0;
	while (++times < 20) {
		HWND handle;
		handle = GetForegroundWindow();
		if (handle != NULL) {
			// mData.HWnd = (uintptr) handle;
			result.HWnd = (HWND)handle;
			return result;
		}
		Sleep(20);
	}

	return result;
#endif
}

void SetTopMost(bool state){
	// Check window validity
	if (!is_valid()) { return; }
#if defined(IS_MACOSX)
	// WARNING: Unavailable
#elif defined(USE_X11)
	// Ignore X errors
	// XDismissErrors();
	// SetState(pub_mData.XWin, STATE_TOPMOST, state);
#elif defined(IS_WINDOWS)
	SetWindowPos(pub_mData.HWnd, state ? HWND_TOPMOST : HWND_NOTOPMOST,
		0, 0, 0, 0, SWP_NOMOVE | SWP_NOSIZE);
#endif
}

void close_main_window () {
   // Check if the window is valid
	if (!is_valid()) { return; }

	close_window_by_Id(pub_mData);
}

void close_window_by_PId(uintptr pid, int8_t isPid){
	MData win = set_handle_pid(pid, isPid);
	close_window_by_Id(win);
}

// CloseWindow
void close_window_by_Id(MData m_data){
	// Check window validity
	if (!is_valid()) { return; }
#if defined(IS_MACOSX)
	AXUIElementRef b = NULL;
	// Retrieve the close button of this window
	if (AXUIElementCopyAttributeValue(m_data.AxID, kAXCloseButtonAttribute, (CFTypeRef*) &b) 
		== kAXErrorSuccess && b != NULL) {
		// Simulate button press on the close button
		AXUIElementPerformAction(b, kAXPressAction);
		CFRelease(b);
	}
#elif defined(USE_X11)
	Display *rDisplay = XOpenDisplay(NULL);
	// Ignore X errors
	XDismissErrors();

	// Close the window
	XDestroyWindow(rDisplay, m_data.XWin);
	XCloseDisplay(rDisplay);
#elif defined(IS_WINDOWS)
	PostMessage(m_data.HWnd, WM_CLOSE, 0, 0);
#endif
}

char* get_main_title(){
	// Check if the window is valid
	if (!is_valid()) { return "is_valid failed."; }

	return get_title_by_hand(pub_mData);
}

char* get_title_by_pid(uintptr pid, int8_t isPid){
	MData win = set_handle_pid(pid, isPid);
	return get_title_by_hand(win);
}

char* named(void *result) {
	char *name = (char*)calloc(strlen(result)+1, sizeof(char*));
	char *rptr = (char*)result;
	char *nptr = name;
	while (*rptr) {
		*nptr = *rptr;
		nptr++;
		rptr++;
	}
	*nptr = '\0';

	return name;
}

char* get_title_by_hand(MData m_data){
	// Check if the window is valid
	if (!is_valid()) { return "is_valid failed."; }
#if defined(IS_MACOSX)
	CFStringRef data = NULL;
	// Determine the current title of the window
	if (AXUIElementCopyAttributeValue(m_data.AxID, kAXTitleAttribute, (CFTypeRef*) &data) == 
			kAXErrorSuccess && data != NULL) {
		char conv[512];
		// Convert result to a C-String
		CFStringGetCString(data, conv, 512, kCFStringEncodingUTF8);
		CFRelease(data);

		char* s = (char*)calloc(100, sizeof(char*));
		if (s) { strcpy(s, conv); }
		// return (char *)&conv;
		return s;
	}

	return "";
#elif defined(USE_X11)
	void* result;
	// Ignore X errors
	XDismissErrors();

	// Get window title (UTF-8)
	result = GetWindowProperty(m_data, WM_NAME, NULL);
	// Check result value
	if (result != NULL) {
		// Convert result to a string
		char* name = named(result);
		XFree(result);

		if (name != NULL) { return name; }
	}

	// Get window title (ASCII)
	result = GetWindowProperty(m_data, XA_WM_NAME, NULL);
	// Check result value
	if (result != NULL) {
		// Convert result to a string
		char* name = named(result);
		XFree(result);

		return name;
	}

	return "";
#elif defined(IS_WINDOWS)
	if (GetWindowText(m_data.HWnd, m_data.Title, 512) > 0){
		char* name = m_data.Title;

		char* str = (char*)calloc(100, sizeof(char*));
		if (str) { strcpy(str, name); }
		return str;
	}

	return "";
#endif
}

int32_t get_PID(void) {
	// Check window validity
	if (!is_valid()) { return 0; }
#if defined(IS_MACOSX)
	pid_t pid = 0;
	// Attempt to retrieve the window pid
	if (AXUIElementGetPid(pub_mData.AxID, &pid)== kAXErrorSuccess) {
		return pid;
	}
	return 0;
#elif defined(USE_X11)
	// Ignore X errors
	XDismissErrors();

	// Get the window PID
	long* result = (long*)GetWindowProperty(pub_mData, WM_PID,NULL);
	// Check result and convert it
	if (result == NULL) { return 0; }
	
	int32_t pid = (int32_t) *result;
	XFree(result);
	return pid;
#elif defined(IS_WINDOWS)
	DWORD id = 0;
	GetWindowThreadProcessId(pub_mData.HWnd, &id);
	return id;
#endif
}
