
#ifdef HAVE_CONFIG_H
	#include <config.h>
#endif

#include <ApplicationServices/ApplicationServices.h>
#include <stdbool.h>
#include <stdio.h>
#include <stdlib.h>
#include "../iohook.h"

// #include "../logger_c.h"
#include "input.h"

// TODO Possibly relocate to input helper.
static inline CGEventFlags get_key_event_mask(iohook_event * const event) {
	CGEventFlags native_mask = 0x00;

	if (event->mask & (MASK_SHIFT))	{ native_mask |= kCGEventFlagMaskShift;		}
	if (event->mask & (MASK_CTRL))	{ native_mask |= kCGEventFlagMaskControl;	}
	if (event->mask & (MASK_META))	{ native_mask |= kCGEventFlagMaskControl;	}
	if (event->mask & (MASK_ALT))	{ native_mask |= kCGEventFlagMaskAlternate;	}

	if (event->type == EVENT_KEY_PRESSED || event->type == EVENT_KEY_RELEASED || event->type == EVENT_KEY_TYPED) {
		switch (event->data.keyboard.keycode) {
			case VC_KP_0:
			case VC_KP_1:
			case VC_KP_2:
			case VC_KP_3:
			case VC_KP_4:
			case VC_KP_5:
			case VC_KP_6:
			case VC_KP_7:
			case VC_KP_8:
			case VC_KP_9:

			case VC_NUM_LOCK:
			case VC_KP_ENTER:
			case VC_KP_MULTIPLY:
			case VC_KP_ADD:
			case VC_KP_SEPARATOR:
			case VC_KP_SUBTRACT:
			case VC_KP_DIVIDE:
			case VC_KP_COMMA:
				native_mask |= kCGEventFlagMaskNumericPad;
				break;
		}
	}

	return native_mask;
}

static inline void post_key_event(iohook_event * const event) {
	bool is_pressed = event->type == EVENT_KEY_PRESSED;

	CGEventSourceRef src = CGEventSourceCreate(kCGEventSourceStateHIDSystemState);
	CGEventRef cg_event = CGEventCreateKeyboardEvent(src,
		(CGKeyCode) scancode_to_keycode(event->data.keyboard.keycode),
		is_pressed);

	CGEventSetFlags(cg_event, get_key_event_mask(event));
	CGEventPost(kCGHIDEventTap, cg_event);	// kCGSessionEventTap also works.
	CFRelease(cg_event);
	CFRelease(src);
}

static inline void post_mouse_button_event(iohook_event * const event, bool is_pressed) {
	CGMouseButton mouse_button;
	CGEventType mouse_type;
	if (event->data.mouse.button == MOUSE_BUTTON1) {
		if (is_pressed) {
			mouse_type = kCGEventLeftMouseDown;
		}
		else {
			mouse_type = kCGEventLeftMouseUp;
		}
		mouse_button = kCGMouseButtonLeft;
	}
	else if (event->data.mouse.button == MOUSE_BUTTON2) {
		if (is_pressed) {
			mouse_type = kCGEventRightMouseDown;
		}
		else {
			mouse_type = kCGEventRightMouseUp;
		}
		mouse_button = kCGMouseButtonRight;
	}
	else {
		if (is_pressed) {
			mouse_type = kCGEventOtherMouseDown;
		}
		else {
			mouse_type = kCGEventOtherMouseUp;
		}
        mouse_button = event->data.mouse.button - 1;
	}

	CGEventSourceRef src = CGEventSourceCreate(kCGEventSourceStateHIDSystemState);
	CGEventRef cg_event = CGEventCreateMouseEvent(src,
		mouse_type,
		CGPointMake(
			(CGFloat) event->data.mouse.x,
			(CGFloat) event->data.mouse.y
		),
        mouse_button
	);
	CGEventPost(kCGHIDEventTap, cg_event);	// kCGSessionEventTap also works.
	CFRelease(cg_event);
	CFRelease(src);
}

static inline void post_mouse_wheel_event(iohook_event * const event) {
	// FIXME Should I create a source event with the coords?
	// It seems to automagically use the current location of the cursor.
	// Two options: Query the mouse, move it to x/y, scroll, then move back
	// OR disable x/y for scroll events on Windows & X11.
	CGScrollEventUnit scroll_unit;
	if (event->data.wheel.type == WHEEL_BLOCK_SCROLL) {
		// Scrolling data is line-based.
		scroll_unit = kCGScrollEventUnitLine;
	}
	else {
		// Scrolling data is pixel-based.
		scroll_unit = kCGScrollEventUnitPixel;
	}

	CGEventSourceRef src = CGEventSourceCreate(kCGEventSourceStateHIDSystemState);
	CGEventRef cg_event = CGEventCreateScrollWheelEvent(src,
		kCGScrollEventUnitLine,
		// TODO Currently only support 1 wheel axis.
		(CGWheelCount) 1, // 1 for Y-only, 2 for Y-X, 3 for Y-X-Z
		event->data.wheel.amount * event->data.wheel.rotation);

	CGEventPost(kCGHIDEventTap, cg_event);	// kCGSessionEventTap also works.
	CFRelease(cg_event);
	CFRelease(src);
}

static inline void post_mouse_motion_event(iohook_event * const event) {
	CGEventSourceRef src = CGEventSourceCreate(kCGEventSourceStateHIDSystemState);
	CGEventRef cg_event;
	if (event->mask >> 8 == 0x00) {
		// No mouse flags.
		cg_event = CGEventCreateMouseEvent(src,
			kCGEventMouseMoved,
			CGPointMake(
				(CGFloat) event->data.mouse.x,
				(CGFloat) event->data.mouse.y
			),
			0
		);
	}
	else if (event->mask & MASK_BUTTON1) {
		cg_event = CGEventCreateMouseEvent(src,
			kCGEventLeftMouseDragged,
			CGPointMake(
				(CGFloat) event->data.mouse.x,
				(CGFloat) event->data.mouse.y
			),
			kCGMouseButtonLeft
		);
	}
	else if (event->mask & MASK_BUTTON2) {
		cg_event = CGEventCreateMouseEvent(src,
			kCGEventRightMouseDragged,
			CGPointMake(
				(CGFloat) event->data.mouse.x,
				(CGFloat) event->data.mouse.y
			),
			kCGMouseButtonRight
		);
	}
	else {
		cg_event = CGEventCreateMouseEvent(src,
			kCGEventOtherMouseDragged,
			CGPointMake(
				(CGFloat) event->data.mouse.x,
				(CGFloat) event->data.mouse.y
			),
			(event->mask >> 8) - 1
		);
	}

	// kCGSessionEventTap also works.
	CGEventPost(kCGHIDEventTap, cg_event);
	CFRelease(cg_event);
	CFRelease(src);
}

IOHOOK_API void hook_post_event(iohook_event * const event) {
	switch (event->type) {
		case EVENT_KEY_PRESSED:
		case EVENT_KEY_RELEASED:
			post_key_event(event);
			break;


		case EVENT_MOUSE_PRESSED:
			post_mouse_button_event(event, true);
			break;

		case EVENT_MOUSE_RELEASED:
			post_mouse_button_event(event, false);
			break;

		case EVENT_MOUSE_CLICKED:
			post_mouse_button_event(event, true);
			post_mouse_button_event(event, false);
			break;

		case EVENT_MOUSE_WHEEL:
            post_mouse_wheel_event(event);
			break;


		case EVENT_MOUSE_MOVED:
		case EVENT_MOUSE_DRAGGED:
			post_mouse_motion_event(event);
			break;


		case EVENT_KEY_TYPED:
			// FIXME Ignoreing EVENT_KEY_TYPED events.

		case EVENT_HOOK_ENABLED:
		case EVENT_HOOK_DISABLED:
			// Ignore hook enabled / disabled events.

		default:
			// Ignore any other garbage.
			logger(LOG_LEVEL_WARN, "%s [%u]: Ignoring post event type %#X\n",
					__FUNCTION__, __LINE__, event->type);
			break;
	}
}
