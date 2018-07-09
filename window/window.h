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
// #include <stdlib.h>
#include "process.h"

struct _MData{
#if defined(IS_MACOSX)
	CGWindowID		CgID;		// Handle to a CGWindowID
	AXUIElementRef	AxID;		// Handle to a AXUIElementRef
#elif defined(USE_X11)
	Window		XWin;		// Handle to an X11 window
#elif defined(IS_WINDOWS)
	HWND			HWnd;		// Handle to a window HWND
	TCHAR 	Title[512];
#endif
};

typedef struct _MData MData;

MData mData;


bool setHandle(uintptr handle);
bool IsValid();
bool IsAxEnabled(bool options);
MData GetActive(void);
void initWindow();

#if defined(IS_MACOSX)

	static Boolean(*gAXIsProcessTrustedWithOptions) (CFDictionaryRef);

	static CFStringRef* gkAXTrustedCheckOptionPrompt;

	AXError _AXUIElementGetWindow(AXUIElementRef, CGWindowID* out);

	static AXUIElementRef GetUIElement(CGWindowID win){
		intptr pid = 0;
		// double_t pid = 0;

		// Create array storing window
		CGWindowID window[1] = { win };
		CFArrayRef wlist = CFArrayCreate(NULL,
				(const void**)window, 1, NULL);

		// Get window info
		CFArrayRef info = CGWindowListCreateDescriptionFromArray(wlist);
		CFRelease(wlist);

		// Check whether the resulting array is populated
		if (info != NULL && CFArrayGetCount(info) > 0){
			// Retrieve description from info array
			CFDictionaryRef desc = (CFDictionaryRef)CFArrayGetValueAtIndex(info, 0);

			// Get window PID
			CFNumberRef data =(CFNumberRef)
				CFDictionaryGetValue(desc, kCGWindowOwnerPID);

			if (data != NULL){
				CFNumberGetValue(data, kCFNumberIntType, &pid);
			}

			// Return result
			CFRelease(info);
		}

		// Check if PID was retrieved
		if (pid <= 0) {return NULL;}

		// Create an accessibility object using retrieved PID
		AXUIElementRef application = AXUIElementCreateApplication(pid);

		if (application == 0) {return NULL;}

		CFArrayRef windows = NULL;
		// Get all windows associated with the app
		AXUIElementCopyAttributeValues(application,
			kAXWindowsAttribute, 0, 1024, &windows);

		// Reference to resulting value
		AXUIElementRef result = NULL;

		if (windows != NULL) {
			int count = CFArrayGetCount(windows);
			// Loop all windows in the process
			for (CFIndex i = 0; i < count; ++i){
				// Get the element at the index
				AXUIElementRef element = (AXUIElementRef)
					CFArrayGetValueAtIndex(windows, i);

				CGWindowID temp = 0;
				// Use undocumented API to get WindowID
				_AXUIElementGetWindow(element, &temp);

				// Check results
				if (temp == win) {
					// Retain element
					CFRetain(element);
					result = element;
					break;
				}
			}

			CFRelease(windows);
		}

		CFRelease(application);
		return result;
	}
#elif defined(USE_X11)

	// Error Handling

	typedef int (*XErrorHandler) (Display*, XErrorEvent*);

	static int XHandleError(Display* dp, XErrorEvent* e) { return 0; }

		XErrorHandler mOld;

		void XDismissErrors (void) {
			Display *rDisplay = XOpenDisplay(NULL);
			// Save old handler and dismiss errors
			mOld = XSetErrorHandler(XHandleError);
			// Flush output buffer
			XSync(rDisplay, False);

			// Reinstate old handler
			XSetErrorHandler(mOld);
		}

	// Definitions

	struct Hints{
		unsigned long Flags;
		unsigned long Funcs;
		unsigned long Decorations;
		  signed long Mode;
		unsigned long Stat;
	};

	static Atom WM_STATE	= None;
	static Atom WM_ABOVE	= None;
	static Atom WM_HIDDEN	= None;
	static Atom WM_HMAX		= None;
	static Atom WM_VMAX		= None;

	static Atom WM_DESKTOP	= None;
	static Atom WM_CURDESK	= None;

	static Atom WM_NAME		= None;
	static Atom WM_UTF8		= None;
	static Atom WM_PID		= None;
	static Atom WM_ACTIVE	= None;
	static Atom WM_HINTS	= None;
	static Atom WM_EXTENTS	= None;

	////////////////////////////////////////////////////////////////////////////////

	static void LoadAtoms (void){
		Display *rDisplay = XOpenDisplay(NULL);
		WM_STATE   = XInternAtom(rDisplay, "_NET_WM_STATE",                True);
		WM_ABOVE   = XInternAtom(rDisplay, "_NET_WM_STATE_ABOVE",          True);
		WM_HIDDEN  = XInternAtom(rDisplay, "_NET_WM_STATE_HIDDEN",         True);
		WM_HMAX    = XInternAtom(rDisplay, "_NET_WM_STATE_MAXIMIZED_HORZ", True);
		WM_VMAX    = XInternAtom(rDisplay, "_NET_WM_STATE_MAXIMIZED_VERT", True);

		WM_DESKTOP = XInternAtom(rDisplay, "_NET_WM_DESKTOP",              True);
		WM_CURDESK = XInternAtom(rDisplay, "_NET_CURRENT_DESKTOP",         True);

		WM_NAME    = XInternAtom(rDisplay, "_NET_WM_NAME",                 True);
		WM_UTF8    = XInternAtom(rDisplay, "UTF8_STRING",                  True);
		WM_PID     = XInternAtom(rDisplay, "_NET_WM_PID",                  True);
		WM_ACTIVE  = XInternAtom(rDisplay, "_NET_ACTIVE_WINDOW",           True);
		WM_HINTS   = XInternAtom(rDisplay, "_MOTIF_WM_HINTS",              True);
		WM_EXTENTS = XInternAtom(rDisplay, "_NET_FRAME_EXTENTS",           True);
	}



	// Functions
	static void* GetWindowProperty(MData win, Atom atom, uint32* items){
		// Property variables
		Atom type; int format;
		unsigned long  nItems;
		unsigned long  bAfter;
		unsigned char* result = NULL;

		Display *rDisplay = XOpenDisplay(NULL);
		// Check the atom
		if (atom != None) {
			// Retrieve and validate the specified property
			if (!XGetWindowProperty(rDisplay, win.XWin, atom, 0,
				BUFSIZ, False, AnyPropertyType, &type, &format,
				&nItems, &bAfter, &result) && result && nItems) {

				// Copy items result
				if (items != NULL) {
					*items = (uint32) nItems;
				}

				return result;
			}
		}

		// Reset the items result if valid
		if (items != NULL) {*items = 0;}

		// Free the result if it got allocated
		if (result != NULL) {
			XFree (result);
		}

		return NULL;
	}

	//////

	#define STATE_TOPMOST  0
	#define STATE_MINIMIZE 1
	#define STATE_MAXIMIZE 2


	//////
	static void SetDesktopForWindow(MData win){
		Display *rDisplay = XOpenDisplay(NULL);
		// Validate every atom that we want to use
		if (WM_DESKTOP != None && WM_CURDESK != None) {
			// Get desktop property
			long* desktop = (long*)GetWindowProperty(win, WM_DESKTOP,NULL);

			// Check result value
			if (desktop != NULL) {
				// Retrieve the screen number
				XWindowAttributes attr = { 0 };
				XGetWindowAttributes(rDisplay, win.XWin, &attr);
				int s = XScreenNumberOfScreen(attr.screen);
				Window root = XRootWindow(rDisplay, s);

				// Prepare an event
				XClientMessageEvent e = { 0 };
				e.window = root; e.format = 32;
				e.message_type = WM_CURDESK;
				e.display = rDisplay;
				e.type = ClientMessage;
				e.data.l[0] = *desktop;
				e.data.l[1] = CurrentTime;

				// Send the message
				XSendEvent(rDisplay,
					root, False, SubstructureNotifyMask |
					SubstructureRedirectMask, (XEvent*) &e);

				XFree(desktop);
			}
		}
	}


#elif defined(IS_WINDOWS)
	//
#endif

//int findwindow()

uintptr initHandle = 0;

void initWindow(uintptr handle){
#if defined(IS_MACOSX)

	mData.CgID = 0;
	mData.AxID = 0;

#elif defined(USE_X11)
	Display *rDisplay = XOpenDisplay(NULL);
	// If atoms loaded
	if (WM_PID == None) {
		// Load all necessary atom properties
		if (rDisplay != NULL) {LoadAtoms();}
	}

	mData.XWin = 0;

#elif defined(IS_WINDOWS)
	mData.HWnd = 0;
#endif

	setHandle(handle);
}

bool IsValid(){
	initWindow(initHandle);
	if (!IsAxEnabled(true)) {
		printf("%s\n", "Window:Accessibility API is disabled!\n"
		"Failed to enable access for assistive devices.");
	}
	MData actdata = GetActive();

#if defined(IS_MACOSX)

	mData.CgID = actdata.CgID;
	mData.AxID = actdata.AxID;

	if (mData.CgID == 0 || mData.AxID == 0)return false;

	CFTypeRef r = NULL;

	// Attempt to get the window role
	if (AXUIElementCopyAttributeValue(mData.AxID,
		kAXRoleAttribute,&r) == kAXErrorSuccess && r){
			CFRelease (r);
			return true;
	}

	return false;

#elif defined(USE_X11)
	mData.XWin = actdata.XWin;
	if (mData.XWin == 0) {return false;}

	Display *rDisplay = XOpenDisplay(NULL);
	// Check for a valid X-Window display
	if (rDisplay == NULL) {return false;}

	// Ignore X errors
	XDismissErrors();

	// Get the window PID property
	void* result = GetWindowProperty(mData, WM_PID,NULL);
	if (result == NULL) {return false;}

	// Free result and return true
	XFree(result);
	return true;

#elif defined(IS_WINDOWS)
	mData.HWnd = actdata.HWnd;

	if (mData.HWnd == 0) {
		return false;
	}

	return IsWindow(mData.HWnd) != 0;

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
			*(void**) (&gAXIsProcessTrustedWithOptions) =
				dlsym (handle, "AXIsProcessTrustedWithOptions");

			gkAXTrustedCheckOptionPrompt = (CFStringRef*)
				dlsym (handle, "kAXTrustedCheckOptionPrompt");
		}
	});

	// Check for new OSX 10.9 function
	if (gAXIsProcessTrustedWithOptions) {
		// Check whether to show prompt
		CFBooleanRef displayPrompt = options ?
			kCFBooleanTrue : kCFBooleanFalse;

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
	if (mData.AxID != NULL){
		CFRelease(mData.AxID);
	}

	// Reset both values
	mData.CgID = 0;
	mData.AxID = 0;

	if (handle == 0){
		return 0;
		// return true;
	}

	// Retrieve the window element
	CGWindowID cgID = (CGWindowID)handle;
	AXUIElementRef axID = GetUIElement(cgID);

	if (axID != NULL){
		mData.CgID = cgID;
		mData.AxID = axID;
		// return 0;
		return true;
	}

	// return 1;
	return false;

#elif defined(USE_X11)

	mData.XWin = (Window)handle;

	if (handle == 0){
		return true;
	}

	if (IsValid()){
		return true;
	}

	mData.XWin = 0;
	return false;
#elif defined(IS_WINDOWS)

	mData.HWnd = (HWND)handle;

	if (handle == 0) {
		return true;
	}

	if (IsValid()) {
		return true;
	}

	mData.HWnd = 0;
	return false;

#endif
}

// uint32 uintptr
uintptr getHandle() {
	#if defined(IS_MACOSX)
		return (uintptr)mData.CgID;
	#elif defined(USE_X11)
		return (uintptr)mData.XWin;
	#elif defined(IS_WINDOWS)
		return (uintptr)mData.HWnd;
	#endif
}

bool IsTopMost(void){
	// Check the window validity
	if (!IsValid()) {return false;}
#if defined(IS_MACOSX)

	return false; // WARNING: Unavailable

#elif defined(USE_X11)

	// Ignore X errors
	// XDismissErrors ();
	// return GetState (mData.XWin, STATE_TOPMOST);

#elif defined(IS_WINDOWS)

	return (GetWindowLongPtr(mData.HWnd, GWL_EXSTYLE)
		& WS_EX_TOPMOST) != 0;

#endif
}

bool IsMinimized(void){
	// Check the window validity
	if (!IsValid()) {return false;}
#if defined(IS_MACOSX)

	CFBooleanRef data = NULL;

	// Determine whether the window is minimized
	if (AXUIElementCopyAttributeValue(mData.AxID,
		kAXMinimizedAttribute, (CFTypeRef*) &data)
		== kAXErrorSuccess && data != NULL) {
		// Convert resulting data into a bool
		bool result = CFBooleanGetValue(data);
		CFRelease(data);
		return result;
	}

	return false;

#elif defined(USE_X11)

	// Ignore X errors
	// XDismissErrors();
	// return GetState (mData.XWin, STATE_MINIMIZE);

#elif defined(IS_WINDOWS)

	return (GetWindowLongPtr(mData.HWnd, GWL_STYLE)
		& WS_MINIMIZE) != 0;

#endif
}

//////

bool IsMaximized(void){
	// Check the window validity
	if (!IsValid()) {return false;}
#if defined(IS_MACOSX)

	return false; // WARNING: Unavailable

#elif defined(USE_X11)

	// Ignore X errors
	// XDismissErrors();
	// return GetState (mData.XWin, STATE_MAXIMIZE);

#elif defined(IS_WINDOWS)

	return (GetWindowLongPtr(mData.HWnd, GWL_STYLE)
		& WS_MAXIMIZE) != 0;

#endif
}

void SetActive(const MData win){
	// Check if the window is valid
	if (!IsValid()) { return; }
#if defined(IS_MACOSX)

	// Attempt to raise the specified window object
	if (AXUIElementPerformAction(win.AxID, kAXRaiseAction)
		!= kAXErrorSuccess) {
		pid_t pid = 0;
		// Attempt to retrieve the PID of the window
		if (AXUIElementGetPid(win.AxID, &pid)
					!= kAXErrorSuccess || !pid) {return;}

		// Ignore deprecated warnings
		#pragma clang diagnostic push
		#pragma clang diagnostic ignored "-Wdeprecated-declarations"

		// NOTE: Until Apple actually removes
		// these functions, there's no real
		// reason to switch to the NS* flavor

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
		XSetInputFocus(rDisplay, win.XWin,
						RevertToParent, CurrentTime);
	}

#elif defined(IS_WINDOWS)

	if (IsMinimized()) {
		ShowWindow(win.HWnd, SW_RESTORE);
	}

	SetForegroundWindow(win.HWnd);

#endif
}

MData GetActive(void){
#if defined(IS_MACOSX)

	MData result;
	// Ignore deprecated warnings
	#pragma clang diagnostic push
	#pragma clang diagnostic ignored "-Wdeprecated-declarations"

	ProcessSerialNumber psn; pid_t pid;
	// Attempt to retrieve the front process
	if (GetFrontProcess(&psn) != 0 ||
		GetProcessPID(&psn, &pid) != 0) {
		return result;
	}

	#pragma clang diagnostic pop

	// Create accessibility object using focused PID
	AXUIElementRef focused = AXUIElementCreateApplication(pid);
	if (focused == NULL) {return result; }// Verify

	AXUIElementRef element;
	// Retrieve the currently focused window
	if (AXUIElementCopyAttributeValue(focused,
		kAXFocusedWindowAttribute, (CFTypeRef*)
		&element) == kAXErrorSuccess && element) {

		CGWindowID win = 0;
		// Use undocumented API to get WID
		if (_AXUIElementGetWindow(element, &win)
			== kAXErrorSuccess && win) {
			// Manually set internals
			result.CgID = win;
			result.AxID = element;
		}
		// Something went wrong
		else {
			CFRelease(element);
		}
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
	void* active = GetWindowProperty(result,WM_ACTIVE,NULL);

	// Check result value
	if (active != NULL) {
		// Extract window from the result
		long window = *((long*)active);
		XFree(active);

		if (window != 0) {
			// Set and return the foreground window
			result.XWin = (Window)window;
			return result;
		}
	}

	// Use input focus instead
	Window window = None;
	int revert = RevertToNone;
	XGetInputFocus(rDisplay, &window, &revert);

	// Return foreground window
	result.XWin = window;
	return result;

#elif defined(IS_WINDOWS)

	// Attempt to get the foreground window multiple times in case
	MData result;

	uint8 times = 0;
	while (++times < 20) {
		HWND handle;
		handle = GetForegroundWindow();
		if (handle != NULL) {
			// mData.HWnd = (uintptr) handle;
			result.HWnd = (HWND)handle;
			return result;
		}
		Sleep (20);
	}

	return result;

#endif
}


void SetTopMost(bool state){
	// Check window validity
	if (!IsValid()) {return;}
#if defined(IS_MACOSX)

	// WARNING: Unavailable

#elif defined(USE_X11)

	// Ignore X errors
	// XDismissErrors();
	// SetState(mData.XWin, STATE_TOPMOST, state);

#elif defined(IS_WINDOWS)

	SetWindowPos(mData.HWnd,
		state ? HWND_TOPMOST : HWND_NOTOPMOST,
		0, 0, 0, 0, SWP_NOMOVE | SWP_NOSIZE);

#endif
}

// CloseWindow
void CloseWin(void){
	// Check window validity
	if (!IsValid()) {return;}
#if defined(IS_MACOSX)
	AXUIElementRef b = NULL;

	// Retrieve the close button of this window
	if (AXUIElementCopyAttributeValue(mData.AxID,
		kAXCloseButtonAttribute, (CFTypeRef*) &b)
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
	XDestroyWindow(rDisplay, mData.XWin);

#elif defined(IS_WINDOWS)

	PostMessage(mData.HWnd, WM_CLOSE, 0, 0);

#endif
}

char *GetTitle(){
	// Check if the window is valid
	if (!IsValid()) {return "IsValid failed.";}

#if defined(IS_MACOSX)

	CFStringRef data = NULL;

	// Determine the current title of the window
	if (AXUIElementCopyAttributeValue(mData.AxID,
		kAXTitleAttribute, (CFTypeRef*) &data)
		== kAXErrorSuccess && data != NULL) {
		char conv[512];
		// Convert result to a C-String
		CFStringGetCString(data, conv,
			512, kCFStringEncodingUTF8);
		CFRelease(data);
		char* s = (char*)calloc(100, sizeof(char*));
    	if (s) { strcpy(s,conv); }
		// return (char *)&conv;
		return s;
	}

	return "";

#elif defined(USE_X11)

	void* result;
	// Ignore X errors
	XDismissErrors();

	// Get window title (UTF-8)
	result = GetWindowProperty(mData, WM_NAME,NULL);

	// Check result value
	if (result != NULL) {
		// Convert result to a string
		char *name = (char*)result;
		XFree(result);

		if (name != NULL) {return name;}
	}

	// Get window title (ASCII)
	result = GetWindowProperty(mData, XA_WM_NAME,NULL);

	// Check result value
	if (result != NULL) {
		// Convert result to a string
		char *name = (char*)result;
		XFree(result);
		return name;
	}

	return "";

#elif defined(IS_WINDOWS)

	return GetWindowText
		(mData.HWnd, mData.Title, 512) > 0 ? mData.Title : "";
	// return GetWindowText
	// 	(mData.HWnd, name, 512) > 0 ?
	// 	_UTF8Encode(name) : "null";

#endif
}

int32 WGetPID(void){
	// Check window validity
	if (!IsValid()) {return 0;}

#if defined(IS_MACOSX)

	pid_t pid = 0;
	// Attempt to retrieve the window pid
	if (AXUIElementGetPid(mData.AxID, &pid)
			== kAXErrorSuccess) {
				return pid;
			}

	return 0;

#elif defined(USE_X11)

	// Ignore X errors
	XDismissErrors();

	// Get the window PID
	long* result = (long*)GetWindowProperty(mData, WM_PID,NULL);

	// Check result and convert it
	if (result == NULL) {return 0;}
	int32 pid = (int32) *result;
	XFree(result);
	return pid;

#elif defined(IS_WINDOWS)

	DWORD  id = 0;
	GetWindowThreadProcessId(mData.HWnd, &id);
	return id;

#endif
}

