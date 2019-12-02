#pragma once
#ifndef MOUSE_H
#define MOUSE_H

#include "../base/os.h"
#include "../base/types.h"

#if defined(_MSC_VER)
	#include "../base/ms_stdbool.h"
#else
	#include <stdbool.h>
#endif

#ifdef __cplusplus
// #ifdefined(__cplusplus)||defined(c_plusplus)
extern "C"
{
#endif

#if defined(IS_MACOSX)

	// #include </System/Library/Frameworks/ApplicationServices.framework/Headers/ApplicationServices.h>
	#include <ApplicationServices/ApplicationServices.h>
	// #include </System/Library/Frameworks/ApplicationServices.framework/Versions/A/Headers/ApplicationServices.h>

	typedef enum  {
		LEFT_BUTTON = kCGMouseButtonLeft,
		RIGHT_BUTTON = kCGMouseButtonRight,
		CENTER_BUTTON = kCGMouseButtonCenter
	} MMMouseButton;

#elif defined(USE_X11)

	enum _MMMouseButton {
		LEFT_BUTTON = 1,
		CENTER_BUTTON = 2,
		RIGHT_BUTTON = 3
	};
	typedef unsigned int MMMouseButton;

#elif defined(IS_WINDOWS)

	enum _MMMouseButton {
		LEFT_BUTTON = 1,
		CENTER_BUTTON = 2,
		RIGHT_BUTTON = 3
	};
	typedef unsigned int MMMouseButton;

#else
	#error "No mouse button constants set for platform"
#endif

#define MMMouseButtonIsValid(button) \
	(button == LEFT_BUTTON || button == RIGHT_BUTTON || \
	 button == CENTER_BUTTON)

enum __MMMouseWheelDirection
{
	DIRECTION_DOWN 	= -1,
	DIRECTION_UP	= 1
};
typedef int MMMouseWheelDirection;

/* Immediately moves the mouse to the given point on-screen.
 * It is up to the caller to ensure that this point is within the
 * screen boundaries. */
void moveMouse(MMPointInt32 point);

/* Like moveMouse, moves the mouse to the given point on-screen, but marks
 * the event as the mouse being dragged on platforms where it is supported.
 * It is up to the caller to ensure that this point is within the screen
 * boundaries. */
void dragMouse(MMPointInt32 point, const MMMouseButton button);

/* Smoothly moves the mouse from the current position to the given point.
 * deadbeef_srand() should be called before using this function.
 *
 * Returns false if unsuccessful (i.e. a point was hit that is outside of the
 * screen boundaries), or true if successful. */
bool smoothlyMoveMouse(MMPointInt32 endPoint, double lowSpeed, double highSpeed);
// bool smoothlyMoveMouse(MMPoint point);

/* Returns the coordinates of the mouse on the current screen. */
MMPointInt32 getMousePos(void);

/* Holds down or releases the mouse with the given button in the current
 * position. */
void toggleMouse(bool down, MMMouseButton button);

/* Clicks the mouse with the given button in the current position. */
void clickMouse(MMMouseButton button);

/* Double clicks the mouse with the given button. */
void doubleClick(MMMouseButton button);

/* Scrolls the mouse in the stated direction.
 * TODO: Add a smoothly scroll mouse next. */
void scrollMouse(int scrollMagnitude, MMMouseWheelDirection scrollDirection);

//#ifdefined(__cplusplus)||defined(c_plusplus)
#ifdef __cplusplus
}
#endif

#endif /* MOUSE_H */