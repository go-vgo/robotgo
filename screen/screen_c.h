//#include "../base/os.h"

#if defined(IS_MACOSX)
	#include <ApplicationServices/ApplicationServices.h>
#elif defined(USE_X11)
	#include <X11/Xlib.h>
	#include <X11/Xresource.h>
	// #include "../base/xdisplay_c.h"
#endif

intptr scaleX();

double sys_scale(int32_t display_id) {
	#if defined(IS_MACOSX)
		CGDirectDisplayID displayID = (CGDirectDisplayID) display_id;
		if (displayID == -1) {
			displayID = CGMainDisplayID();
		}
		
		CGDisplayModeRef modeRef = CGDisplayCopyDisplayMode(displayID);
		double pixelWidth = CGDisplayModeGetPixelWidth(modeRef);
		double targetWidth = CGDisplayModeGetWidth(modeRef);
	
		return pixelWidth / targetWidth;
	#elif defined(USE_X11)
		Display *dpy = XOpenDisplay(NULL);

		int scr = 0; /* Screen number */
		double xres = ((((double) DisplayWidth(dpy, scr)) * 25.4) /
			((double) DisplayWidthMM(dpy, scr)));

		char *rms = XResourceManagerString(dpy);
		if (rms) {
			XrmDatabase db = XrmGetStringDatabase(rms);
			if (db) {
				XrmValue value;
				char *type = NULL;

				if (XrmGetResource(db, "Xft.dpi", "String", &type, &value)) {
					if (value.addr) {
						xres = atof(value.addr);
					}
				}

				XrmDestroyDatabase(db);
			}
		}
		XCloseDisplay (dpy);

		return xres / 96.0;
   	#elif defined(IS_WINDOWS)
   		double s = scaleX() / 96.0;
   		return s;
   	#endif
}

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
		return horizontalDPI;
	#endif
}

MMSizeInt32 getMainDisplaySize(void) {
#if defined(IS_MACOSX)
	CGDirectDisplayID displayID = CGMainDisplayID();
	CGRect displayRect = CGDisplayBounds(displayID);
	
	CGSize size = displayRect.size;
	return MMSizeInt32Make((int32_t)size.width, (int32_t)size.height);
#elif defined(USE_X11)
	Display *display = XGetMainDisplay();
	const int screen = DefaultScreen(display);

	return MMSizeInt32Make(
						(int32_t)DisplayWidth(display, screen),
	                	(int32_t)DisplayHeight(display, screen));
#elif defined(IS_WINDOWS)
	return MMSizeInt32Make(
 						(int32_t)GetSystemMetrics(SM_CXSCREEN),
 		                (int32_t)GetSystemMetrics(SM_CYSCREEN));
#endif
}

MMRectInt32 getScreenRect(int32_t display_id) {
#if defined(IS_MACOSX)
	CGDirectDisplayID displayID = (CGDirectDisplayID) display_id;
	if (display_id == -1) {
		displayID = CGMainDisplayID();
	}
	CGRect displayRect = CGDisplayBounds(displayID);

	CGPoint point = displayRect.origin;
	CGSize size = displayRect.size;
	return MMRectInt32Make(
		(int32_t)point.x, (int32_t)point.y,
		(int32_t)size.width, (int32_t)size.height);
#elif defined(USE_X11)
	Display *display = XGetMainDisplay();
	const int screen = DefaultScreen(display);

	return MMRectInt32Make(
					(int32_t)0, (int32_t)0,
					(int32_t)DisplayWidth(display, screen),
	                (int32_t)DisplayHeight(display, screen));
#elif defined(IS_WINDOWS)
	if (GetSystemMetrics(SM_CMONITORS) == 1 
			|| display_id == -1 || display_id == 0) {
 		return MMRectInt32Make(
						(int32_t)0,
						(int32_t)0,
			 			(int32_t)GetSystemMetrics(SM_CXSCREEN),
 		                (int32_t)GetSystemMetrics(SM_CYSCREEN));
 	} else {
 		return MMRectInt32Make(
			 			(int32_t)GetSystemMetrics(SM_XVIRTUALSCREEN),
						(int32_t)GetSystemMetrics(SM_YVIRTUALSCREEN),
						(int32_t)GetSystemMetrics(SM_CXVIRTUALSCREEN),
 		                (int32_t)GetSystemMetrics(SM_CYVIRTUALSCREEN));
 	}
#endif
}

bool pointVisibleOnMainDisplay(MMPointInt32 point){
	MMSizeInt32 displaySize = getMainDisplaySize();
	return point.x < displaySize.w && point.y < displaySize.h;
}
