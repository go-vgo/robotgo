
#ifdef HAVE_CONFIG_H
	#include <config.h>
#endif

#include <stdbool.h>
#include <stdio.h>
#include <stdlib.h>
#include <X11/Xlib.h>
#include <X11/Xutil.h>
#ifdef USE_XTEST
	#include <X11/extensions/XTest.h>
#endif

#include "../iohook.h"
#include "input.h"
// #include "../logger.h"

extern Display *properties_disp;

// This lookup table must be in the same order the masks are defined.
#ifdef USE_XTEST
static KeySym keymask_lookup[8] = {
	XK_Shift_L,
	XK_Control_L,
	XK_Meta_L,
	XK_Alt_L,

	XK_Shift_R,
	XK_Control_R,
	XK_Meta_R,
	XK_Alt_R
};

static unsigned int btnmask_lookup[5] = {
	MASK_BUTTON1,
	MASK_BUTTON2,
	MASK_BUTTON3,
	MASK_BUTTON4,
	MASK_BUTTON5
};
#else
// TODO Possibly relocate to input helper.
static unsigned int convert_to_native_mask(unsigned int mask) {
	unsigned int native_mask = 0x00;

	if (mask & (MASK_SHIFT))	{ native_mask |= ShiftMask;		}
	if (mask & (MASK_CTRL))		{ native_mask |= ControlMask;	}
	if (mask & (MASK_META))		{ native_mask |= Mod4Mask;		}
	if (mask & (MASK_ALT))		{ native_mask |= Mod1Mask;		}

	if (mask & MASK_BUTTON1)	{ native_mask |= Button1Mask;	}
	if (mask & MASK_BUTTON2)	{ native_mask |= Button2Mask;	}
	if (mask & MASK_BUTTON3)	{ native_mask |= Button3Mask;	}
	if (mask & MASK_BUTTON4)	{ native_mask |= Button4Mask;	}
	if (mask & MASK_BUTTON5)	{ native_mask |= Button5Mask;	}

	return native_mask;
}
#endif

static inline void post_key_event(iohook_event * const event) {
	#ifdef USE_XTEST
	// FIXME Currently ignoring EVENT_KEY_TYPED.
	if (event->type == EVENT_KEY_PRESSED) {
		XTestFakeKeyEvent(
			properties_disp,
			scancode_to_keycode(event->data.keyboard.keycode),
			True,
			0);
	}
	else if (event->type == EVENT_KEY_RELEASED) {
		XTestFakeKeyEvent(
			properties_disp,
			scancode_to_keycode(event->data.keyboard.keycode),
			False,
			0);
	}
	#else
	XKeyEvent key_event;

	key_event.serial = 0x00;
	key_event.send_event = False;
	key_event.display = properties_disp;
	key_event.time = CurrentTime;
	key_event.same_screen = True;

	unsigned int mask;
	if (!XQueryPointer(properties_disp, DefaultRootWindow(properties_disp), &(key_event.root), &(key_event.subwindow), &(key_event.x_root), &(key_event.y_root), &(key_event.x), &(key_event.y), &mask)) {
		key_event.root = DefaultRootWindow(properties_disp);
		key_event.window = key_event.root;
		key_event.subwindow = None;

		key_event.x_root = 0;
		key_event.y_root = 0;
		key_event.x = 0;
		key_event.y = 0;
	}

	key_event.state = convert_to_native_mask(event->mask);
	key_event.keycode = XKeysymToKeycode(properties_disp, scancode_to_keycode(event->data.keyboard.keycode));

	// FIXME Currently ignoring typed events.
	if (event->type == EVENT_KEY_PRESSED) {
		key_event.type = KeyPress;
		XSendEvent(properties_disp, InputFocus, False, KeyPressMask, (XEvent *) &key_event);
	}
	else if (event->type == EVENT_KEY_RELEASED) {
		key_event.type = KeyRelease;
		XSendEvent(properties_disp, InputFocus, False, KeyReleaseMask, (XEvent *) &key_event);
	}
	#endif
}

static inline void post_mouse_button_event(iohook_event * const event) {
	#ifdef USE_XTEST
	Window ret_root;
	Window ret_child;
	int root_x;
	int root_y;
	int win_x;
	int win_y;
	unsigned int mask;

	Window win_root = XDefaultRootWindow(properties_disp);
	Bool query_status = XQueryPointer(properties_disp, win_root, &ret_root, &ret_child, &root_x, &root_y, &win_x, &win_y, &mask);
	if (query_status) {
		if (event->data.mouse.x != root_x || event->data.mouse.y != root_y) {
			// Move the pointer to the specified position.
			XTestFakeMotionEvent(properties_disp, -1, event->data.mouse.x, event->data.mouse.y, 0);
		}
		else {
			query_status = False;
		}
	}

	if (event->type == EVENT_MOUSE_WHEEL) {
		// Wheel events should be the same as click events on X11.
		// type, amount and rotation
		if (event->data.wheel.rotation < 0) {
			XTestFakeButtonEvent(properties_disp, WheelUp, True, 0);
			XTestFakeButtonEvent(properties_disp, WheelUp, False, 0);
		}
		else {
			XTestFakeButtonEvent(properties_disp, WheelDown, True, 0);
			XTestFakeButtonEvent(properties_disp, WheelDown, False, 0);
		}
	}
	else if (event->type == EVENT_MOUSE_PRESSED) {
		XTestFakeButtonEvent(properties_disp, event->data.mouse.button, True, 0);
	}
	else if (event->type == EVENT_MOUSE_RELEASED) {
		XTestFakeButtonEvent(properties_disp, event->data.mouse.button, False, 0);
	}
	else if (event->type == EVENT_MOUSE_CLICKED) {
		XTestFakeButtonEvent(properties_disp, event->data.mouse.button, True, 0);
		XTestFakeButtonEvent(properties_disp, event->data.mouse.button, False, 0);
	}

	if (query_status) {
		// Move the pointer back to the original position.
		XTestFakeMotionEvent(properties_disp, -1, root_x, root_y, 0);
	}
	#else
	XButtonEvent btn_event;

	btn_event.serial = 0x00;
	btn_event.send_event = False;
	btn_event.display = properties_disp;
	btn_event.time = CurrentTime;
	btn_event.same_screen = True;

	btn_event.root = DefaultRootWindow(properties_disp);
	btn_event.window = btn_event.root;
	btn_event.subwindow = None;

	btn_event.type = 0x00;
	btn_event.state = 0x00;
	btn_event.x_root = 0;
	btn_event.y_root = 0;
	btn_event.x = 0;
	btn_event.y = 0;
	btn_event.button = 0x00;

	btn_event.state = convert_to_native_mask(event->mask);

	btn_event.x = event->data.mouse.x;
	btn_event.y = event->data.mouse.y;

	#if defined(USE_XINERAMA) || defined(USE_XRANDR)
	uint8_t screen_count;
	screen_data *screens = hook_create_screen_info(&screen_count);
	if (screen_count > 1) {
		btn_event.x += screens[0].x;
		btn_event.y += screens[0].y;
	}

	if (screens != NULL) {
		free(screens);
	}
	#endif

	// These are the same because Window == Root Window.
	btn_event.x_root = btn_event.x;
	btn_event.y_root = btn_event.y;

	if (event->type == EVENT_MOUSE_WHEEL) {
		// type, amount and rotation
		if (event->data.wheel.rotation < 0) {
			btn_event.button = WheelUp;
		}
		else {
			btn_event.button = WheelDown;
		}
	}

	if (event->type != EVENT_MOUSE_RELEASED) {
		// FIXME Where do we set event->button?
		btn_event.type = ButtonPress;
		XSendEvent(properties_disp, InputFocus, False, ButtonPressMask, (XEvent *) &btn_event);
	}

	if (event->type != EVENT_MOUSE_PRESSED) {
		btn_event.type = ButtonRelease;
		XSendEvent(properties_disp, InputFocus, False, ButtonReleaseMask, (XEvent *) &btn_event);
	}
	#endif
}

static inline void post_mouse_motion_event(iohook_event * const event) {
    #ifdef USE_XTEST
	XTestFakeMotionEvent(properties_disp, -1, event->data.mouse.x, event->data.mouse.y, 0);
    #else
	XMotionEvent mov_event;

	mov_event.serial = MotionNotify;
	mov_event.send_event = False;
	mov_event.display = properties_disp;
	mov_event.time = CurrentTime;
	mov_event.same_screen = True;
	mov_event.is_hint = NotifyNormal,
	mov_event.root = DefaultRootWindow(properties_disp);
	mov_event.window = mov_event.root;
	mov_event.subwindow = None;

	mov_event.type = 0x00;
	mov_event.state = 0x00;
	mov_event.x_root = 0;
	mov_event.y_root = 0;
	mov_event.x = 0;
	mov_event.y = 0;

	mov_event.state = convert_to_native_mask(event->mask);

	mov_event.x = event->data.mouse.x;
	mov_event.y = event->data.mouse.y;

	#if defined(USE_XINERAMA) || defined(USE_XRANDR)
	uint8_t screen_count;
	screen_data *screens = hook_create_screen_info(&screen_count);
	if (screen_count > 1) {
		mov_event.x += screens[0].x;
		mov_event.y += screens[0].y;
	}

	if (screens != NULL) {
		free(screens);
	}
	#endif

	// These are the same because Window == Root Window.
	mov_event.x_root = mov_event.x;
	mov_event.y_root = mov_event.y;

	long int event_mask = NoEventMask;
	if (event->type == EVENT_MOUSE_DRAGGED) {
		#if Button1Mask == Button1MotionMask && \
			Button2Mask == Button2MotionMask && \
			Button3Mask == Button3MotionMask && \
			Button4Mask == Button4MotionMask && \
			Button5Mask == Button5MotionMask
		// This little trick only works if Button#MotionMasks align with
		// the Button#Masks.
		event_mask = mov_event.state &
				(Button1MotionMask | Button2MotionMask |
				Button2MotionMask | Button3MotionMask | Button5MotionMask);
		#else
		// Fallback to some slightly larger...
		if (event->state & Button1Mask) {
			event_mask |= Button1MotionMask;
		}

		if (event->state & Button2Mask) {
			event_mask |= Button2MotionMask;
		}

		if (event->state & Button3Mask) {
			event_mask |= Button3MotionMask;
		}

		if (event->state & Button4Mask) {
			event_mask |= Button4MotionMask;
		}

		if (event->state & Button5Mask) {
			event_mask |= Button5MotionMask;
		}
		#endif
	}

	// NOTE x_mask = NoEventMask.
	XSendEvent(properties_disp, InputFocus, False, event_mask, (XEvent *) &mov_event);
    #endif
}

IOHOOK_API void hook_post_event(iohook_event * const event) {
	XLockDisplay(properties_disp);

	#ifdef USE_XTEST
	// XTest does not have modifier support, so we fake it by depressing the
	// appropriate modifier keys.
	unsigned int i;
	for (i = 0; i < sizeof(keymask_lookup) / sizeof(KeySym); i++) {
		if (event->mask & 1 << i) {
			XTestFakeKeyEvent(properties_disp, XKeysymToKeycode(properties_disp, keymask_lookup[i]), True, 0);
		}
	}

	unsigned int i;
	for (i = 0; i < sizeof(btnmask_lookup) / sizeof(unsigned int); i++) {
		if (event->mask & btnmask_lookup[i]) {
			XTestFakeButtonEvent(properties_disp, i + 1, True, 0);
		}
	}
	#endif

	switch (event->type) {
		case EVENT_KEY_PRESSED:
		case EVENT_KEY_RELEASED:
		case EVENT_KEY_TYPED:
			post_key_event(event);
			break;

		case EVENT_MOUSE_PRESSED:
		case EVENT_MOUSE_RELEASED:
		case EVENT_MOUSE_WHEEL:
		case EVENT_MOUSE_CLICKED:
			post_mouse_button_event(event);
			break;

		case EVENT_MOUSE_DRAGGED:
		case EVENT_MOUSE_MOVED:
			post_mouse_motion_event(event);
			break;

		case EVENT_HOOK_ENABLED:
		case EVENT_HOOK_DISABLED:
			// Ignore hook enabled / disabled events.

		default:
			// Ignore any other garbage.
			logger(LOG_LEVEL_WARN, "%s [%u]: Ignoring post event type %#X\n",
				__FUNCTION__, __LINE__, event->type);
			break;
	}

	#ifdef USE_XTEST
	// Release the previously held modifier keys used to fake the event mask.
	unsigned int i ;
	for (i= 0; i < sizeof(keymask_lookup) / sizeof(KeySym); i++) {
		if (event->mask & 1 << i) {
			XTestFakeKeyEvent(properties_disp, XKeysymToKeycode(properties_disp, keymask_lookup[i]), False, 0);
		}
	}
	unsigned int i;
	for (i = 0; i < sizeof(btnmask_lookup) / sizeof(unsigned int); i++) {
		if (event->mask & btnmask_lookup[i]) {
			XTestFakeButtonEvent(properties_disp, i + 1, False, 0);
		}
	}
	#endif

	// Don't forget to flush!
	XSync(properties_disp, True);
	XUnlockDisplay(properties_disp);
}
