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

Bounds get_client(uintptr pid, uintptr isHwnd);

intptr scaleX(){
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

double sys_scale() {
	#if defined(IS_MACOSX)
	
		CGDirectDisplayID displayID = CGMainDisplayID();
		CGDisplayModeRef modeRef = CGDisplayCopyDisplayMode(displayID);

		double pixelWidth = CGDisplayModeGetPixelWidth(modeRef);
		double targetWidth = CGDisplayModeGetWidth(modeRef);
		
		return pixelWidth / targetWidth;
	#elif defined(USE_X11)
		
		double xres;
		Display *dpy;

		char *displayname = NULL;
		int scr = 0; /* Screen number */

		dpy = XOpenDisplay (displayname);
		xres = ((((double) DisplayWidth(dpy, scr)) * 25.4) /
			((double) DisplayWidthMM(dpy, scr)));

   		XCloseDisplay (dpy);

   		return xres + 0.5;
   	#elif defined(IS_WINDOWS)
   		double s = scaleX() / 96.0;
   		return s;
   	#endif
}

intptr scaleY(){
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

Bounds get_bounds(uintptr pid, uintptr isHwnd){
	// Check if the window is valid
	Bounds bounds;
	if (!IsValid()) { return bounds; }

    #if defined(IS_MACOSX)

		// Bounds bounds;
		AXValueRef axp = NULL;
		AXValueRef axs = NULL;
		AXUIElementRef AxID = AXUIElementCreateApplication(pid);

		// Determine the current point of the window
		if (AXUIElementCopyAttributeValue(
			AxID, kAXPositionAttribute, (CFTypeRef*) &axp)
			!= kAXErrorSuccess || axp == NULL){
			goto exit;
		}

		// Determine the current size of the window
		if (AXUIElementCopyAttributeValue(
			AxID, kAXSizeAttribute, (CFTypeRef*) &axs)
			!= kAXErrorSuccess || axs == NULL){
			goto exit;
		}

		CGPoint p; CGSize s;
		// Attempt to convert both values into atomic types
		if (AXValueGetValue(axp, kAXValueCGPointType, &p) &&
			AXValueGetValue(axs, kAXValueCGSizeType, &s)){
			bounds.X = p.x;
			bounds.Y = p.y;
			bounds.W = s.width;
			bounds.H = s.height;
		}
		
		// return bounds;
	exit:
		if (axp != NULL) { CFRelease(axp); }
		if (axs != NULL) { CFRelease(axs); }

		return bounds;

    #elif defined(USE_X11)

        // Ignore X errors
        XDismissErrors();
        MData win;
        win.XWin = (Window)pid;

        Bounds client = get_client(pid, isHwnd);
        Bounds frame = GetFrame(win);

        bounds.X = client.X - frame.X;
        bounds.Y = client.Y - frame.Y;
        bounds.W = client.W + frame.W;
        bounds.H = client.H + frame.H;

        return bounds;

    #elif defined(IS_WINDOWS)
        HWND hwnd;
        if (isHwnd == 0) {
            hwnd= GetHwndByPId(pid);
        } else {
            hwnd = (HWND)pid;
        }

        RECT rect = { 0 };
        GetWindowRect(hwnd, &rect);

        bounds.X = rect.left;
        bounds.Y = rect.top;
        bounds.W = rect.right - rect.left;
        bounds.H = rect.bottom - rect.top;

        return bounds;

    #endif
}

Bounds get_client(uintptr pid, uintptr isHwnd){
	// Check if the window is valid
	Bounds bounds;
	if (!IsValid()) { return bounds; }

	#if defined(IS_MACOSX)

		return get_bounds(pid, isHwnd);

	#elif defined(USE_X11)

        Display *rDisplay = XOpenDisplay(NULL);

		// Ignore X errors
		XDismissErrors();
		MData win;
        win.XWin = (Window)pid;

		// Property variables
		Window root, parent;
		Window* children;
		unsigned int count;
		int32 x = 0, y = 0;

		// Check if the window is the root
		XQueryTree(rDisplay, win.XWin,
			&root, &parent, &children, &count);
		if (children) { XFree(children); }

		// Retrieve window attributes
		XWindowAttributes attr = { 0 };
		XGetWindowAttributes(rDisplay, win.XWin, &attr);

		// Coordinates must be translated
		if (parent != attr.root){
			XTranslateCoordinates(rDisplay, win.XWin, attr.root, attr.x,
			 attr.y, &x, &y, &parent);
		}
		// Coordinates can be left alone
		else {
			x = attr.x;
			y = attr.y;
		}

		// Return resulting window bounds
		bounds.X = x;
		bounds.Y = y;
		bounds.W = attr.width;
		bounds.H = attr.height;
		return bounds;

	#elif defined(IS_WINDOWS)
		HWND hwnd;
		if (isHwnd == 0) {
			hwnd= GetHwndByPId(pid);
		} else {
			hwnd = (HWND)pid;
		}


		RECT rect = { 0 };
		GetClientRect(hwnd, &rect);

		POINT point;
		point.x = rect.left;
		point.y = rect.top;

		// Convert the client point to screen
		ClientToScreen(hwnd, &point);

		bounds.X = point.x;
		bounds.Y = point.y;
		bounds.W = rect.right - rect.left;
		bounds.H = rect.bottom - rect.top;

		return bounds;

	#endif
}