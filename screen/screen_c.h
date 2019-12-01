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
	return MMSizeInt32Make((int32_t)CGDisplayPixelsWide(displayID),
	                  (int32_t)CGDisplayPixelsHigh(displayID));
#elif defined(USE_X11)
	Display *display = XGetMainDisplay();
	const int screen = DefaultScreen(display);

	return MMSizeInt32Make((int32_t)DisplayWidth(display, screen),
	                  (int32_t)DisplayHeight(display, screen));
#elif defined(IS_WINDOWS)
	if (GetSystemMetrics(SM_CMONITORS) == 1) {
 		return MMSizeInt32Make((int32_t)GetSystemMetrics(SM_CXSCREEN),
 		                  (int32_t)GetSystemMetrics(SM_CYSCREEN));
 	} else 	{
 		return MMSizeInt32Make((int32_t)GetSystemMetrics(SM_CXVIRTUALSCREEN),
 		                  (int32_t)GetSystemMetrics(SM_CYVIRTUALSCREEN));
 	}
#endif
}

bool pointVisibleOnMainDisplay(MMPointInt32 point){
	MMSizeInt32 displaySize = getMainDisplaySize();
	return point.x < displaySize.w && point.y < displaySize.h;
}
