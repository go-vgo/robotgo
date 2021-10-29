#include "screen.h"
//#include "../base/os.h"

#if defined(IS_MACOSX)
	#include <ApplicationServices/ApplicationServices.h>
#elif defined(USE_X11)
	#include <X11/Xlib.h>
	// #include "../base/xdisplay_c.h"
#endif

MMSizeInt32 getMainDisplaySize(void){
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
	if (GetSystemMetrics(SM_CMONITORS) == 1) {
 		return MMSizeInt32Make(
 						(int32_t)GetSystemMetrics(SM_CXSCREEN),
 		                (int32_t)GetSystemMetrics(SM_CYSCREEN));
 	} else {
 		return MMSizeInt32Make(
 						(int32_t)GetSystemMetrics(SM_CXVIRTUALSCREEN),
 		                (int32_t)GetSystemMetrics(SM_CYVIRTUALSCREEN));
 	}
#endif
}

MMRectInt32 getScreenRect(int32_t display_id) {
#if defined(IS_MACOSX)
	CGDirectDisplayID displayID = (CGDirectDisplayID) display_id;
	if (display_id == 0) {
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
					(int32_t)0,
					(int32_t)0,
					(int32_t)DisplayWidth(display, screen),
	                (int32_t)DisplayHeight(display, screen));
#elif defined(IS_WINDOWS)
	if (GetSystemMetrics(SM_CMONITORS) == 1) {
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
