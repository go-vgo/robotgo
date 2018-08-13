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

struct _Bounds{
	int32		X;				// Top left X coordinate
	int32		Y;				// Top left Y coordinate
	int32		W;				// Total bounds width
	int32		H;				// Total bounds height
};

typedef struct _Bounds Bounds;

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

	static Bounds GetFrame(MData win){
		Bounds frame;
		// Retrieve frame bounds
		if (WM_EXTENTS != None) {
			long* result; uint32 nItems = 0;
			// Get the window extents property
			result = (long*) GetWindowProperty(win, WM_EXTENTS, &nItems);

			// Verify the results
			if (result != NULL) {
				if (nItems == 4) {
					frame.X = (int32) result[0];
					frame.Y = (int32) result[2];
					frame.W = (int32) result[0] + (int32) result[1];
					frame.H =  (int32) result[2] + (int32) result[3];
				}

				XFree(result);
			}
		}

		return frame;
	}


#elif defined(IS_WINDOWS)
	//
#endif