
#ifdef HAVE_CONFIG_H
	#include <config.h>
#endif

#define USE_XKB 0
#define USE_XKBCOMMON 0
#include <inttypes.h>
#include <limits.h>
#ifdef USE_XRECORD_ASYNC
	#include <pthread.h>
#endif
#include <stdint.h>

#include <X11/keysym.h>
#include <X11/Xlibint.h>
#include <X11/Xlib.h>
#include <X11/extensions/record.h>
// #ifdef USE_XKB
#include <xcb/xkb.h>
#include <X11/XKBlib.h>
// #endif
#if defined(USE_XINERAMA) && !defined(USE_XRANDR)
	#include <X11/extensions/Xinerama.h>
#elif defined(USE_XRANDR)
	#include <X11/extensions/Xrandr.h>
#else
// TODO We may need to fallback to the xf86vm extension for things like TwinView.
// #pragma message("*** Warning: Xinerama or XRandR support is required to produce cross-platform mouse coordinates for multi-head configurations!")
// #pragma message("... Assuming single-head display.")
#endif

#include "../iohook.h"
// #include "../logger.h"
#include "input.h"

// Thread and hook handles.
#ifdef USE_XRECORD_ASYNC
static bool running;

static pthread_cond_t hook_xrecord_cond = PTHREAD_COND_INITIALIZER;
static pthread_mutex_t hook_xrecord_mutex = PTHREAD_MUTEX_INITIALIZER;
#endif

typedef struct _hook_info {
	struct _data {
		Display *display;
		XRecordRange *range;
	} data;
	struct _ctrl {
		Display *display;
		XRecordContext context;
	} ctrl;
	struct _input {
		#ifdef USE_XKBCOMMON
		xcb_connection_t *connection;
		struct xkb_context *context;
    	#endif
		uint16_t mask;
		struct _mouse {
			bool is_dragged;
			struct _click {
				unsigned short int count;
				long int time;
				unsigned short int button;
			} click;
		} mouse;
	} input;
} hook_info;
static hook_info *hook;

// For this struct, refer to libxnee, requires Xlibint.h
typedef union {
	unsigned char		type;
	xEvent				event;
	xResourceReq		req;
	xGenericReply		reply;
	xError				error;
	xConnSetupPrefix	setup;
} XRecordDatum;

#if defined(USE_XKBCOMMON)
//struct xkb_keymap *keymap;
//struct xkb_state *state = xkb_state_new(keymap);
static struct xkb_state *state = NULL;
#endif

// Virtual event pointer.
static iohook_event event;

// Event dispatch callback.
static dispatcher_t dispatcher = NULL;

IOHOOK_API void hook_set_dispatch_proc(dispatcher_t dispatch_proc) {
	logger(LOG_LEVEL_DEBUG,	"%s [%u]: Setting new dispatch callback to %#p.\n",
			__FUNCTION__, __LINE__, dispatch_proc);

	dispatcher = dispatch_proc;
}

// Send out an event if a dispatcher was set.
static inline void dispatch_event(iohook_event *const event) {
	if (dispatcher != NULL) {
		logger(LOG_LEVEL_DEBUG,	"%s [%u]: Dispatching event type %u.\n",
				__FUNCTION__, __LINE__, event->type);

		dispatcher(event);
	}
	else {
		logger(LOG_LEVEL_WARN,	"%s [%u]: No dispatch callback set!\n",
				__FUNCTION__, __LINE__);
	}
}

// Set the native modifier mask for future events.
static inline void set_modifier_mask(uint16_t mask) {
	hook->input.mask |= mask;
}

// Unset the native modifier mask for future events.
static inline void unset_modifier_mask(uint16_t mask) {
	hook->input.mask &= ~mask;
}

// Get the current native modifier mask state.
static inline uint16_t get_modifiers() {
	return hook->input.mask;
}

// Initialize the modifier lock masks.
static void initialize_locks() {
	#ifdef USE_XKBCOMMON

	if (xkb_state_led_name_is_active(state, XKB_LED_NAME_CAPS)) {
		set_modifier_mask(MASK_CAPS_LOCK);
	}
	else {
		unset_modifier_mask(MASK_CAPS_LOCK);
	}

	if (xkb_state_led_name_is_active(state, XKB_LED_NAME_NUM)) {
		set_modifier_mask(MASK_NUM_LOCK);
	}
	else {
		unset_modifier_mask(MASK_NUM_LOCK);
	}

	if (xkb_state_led_name_is_active(state, XKB_LED_NAME_SCROLL)) {
		set_modifier_mask(MASK_SCROLL_LOCK);
	}
	else {
		unset_modifier_mask(MASK_SCROLL_LOCK);
	}
	#else
	unsigned int led_mask = 0x00;
	if (XkbGetIndicatorState(hook->ctrl.display, XkbUseCoreKbd, &led_mask) == Success) {
		if (led_mask & 0x01) {
			set_modifier_mask(MASK_CAPS_LOCK);
		}
		else {
			unset_modifier_mask(MASK_CAPS_LOCK);
		}

		if (led_mask & 0x02) {
			set_modifier_mask(MASK_NUM_LOCK);
		}
		else {
			unset_modifier_mask(MASK_NUM_LOCK);
		}

		if (led_mask & 0x04) {
			set_modifier_mask(MASK_SCROLL_LOCK);
		}
		else {
			unset_modifier_mask(MASK_SCROLL_LOCK);
		}
	}
	else {
		logger(LOG_LEVEL_WARN, "%s [%u]: XkbGetIndicatorState failed to get current led mask!\n",
				__FUNCTION__, __LINE__);
	}
	#endif
}

// Initialize the modifier mask to the current modifiers.
static void initialize_modifiers() {
	hook->input.mask = 0x0000;

	KeyCode keycode;
	char keymap[32];
	XQueryKeymap(hook->ctrl.display, keymap);

  	Window unused_win;
    int unused_int;
	unsigned int mask;
	if (XQueryPointer(hook->ctrl.display, DefaultRootWindow(hook->ctrl.display), &unused_win, &unused_win, &unused_int, &unused_int, &unused_int, &unused_int, &mask)) {
		if (mask & ShiftMask) {
			keycode = XKeysymToKeycode(hook->ctrl.display, XK_Shift_L);
			if (keymap[keycode / 8] & (1 << (keycode % 8))) { set_modifier_mask(MASK_SHIFT_L);	}
			keycode = XKeysymToKeycode(hook->ctrl.display, XK_Shift_R);
			if (keymap[keycode / 8] & (1 << (keycode % 8))) { set_modifier_mask(MASK_SHIFT_R);	}
		}
		if (mask & ControlMask) {
			keycode = XKeysymToKeycode(hook->ctrl.display, XK_Control_L);
			if (keymap[keycode / 8] & (1 << (keycode % 8))) { set_modifier_mask(MASK_CTRL_L);	}
			keycode = XKeysymToKeycode(hook->ctrl.display, XK_Control_R);
			if (keymap[keycode / 8] & (1 << (keycode % 8))) { set_modifier_mask(MASK_CTRL_R);	}
		}
		if (mask & Mod1Mask) {
			keycode = XKeysymToKeycode(hook->ctrl.display, XK_Alt_L);
			if (keymap[keycode / 8] & (1 << (keycode % 8))) { set_modifier_mask(MASK_ALT_L);	}
			keycode = XKeysymToKeycode(hook->ctrl.display, XK_Alt_R);
			if (keymap[keycode / 8] & (1 << (keycode % 8))) { set_modifier_mask(MASK_ALT_R);	}
		}
		if (mask & Mod4Mask) {
			keycode = XKeysymToKeycode(hook->ctrl.display, XK_Super_L);
			if (keymap[keycode / 8] & (1 << (keycode % 8))) { set_modifier_mask(MASK_META_L);	}
			keycode = XKeysymToKeycode(hook->ctrl.display, XK_Super_R);
			if (keymap[keycode / 8] & (1 << (keycode % 8))) { set_modifier_mask(MASK_META_R);	}
		}

		if (mask & Button1Mask)	{ set_modifier_mask(MASK_BUTTON1);	}
		if (mask & Button2Mask)	{ set_modifier_mask(MASK_BUTTON2);	}
		if (mask & Button3Mask)	{ set_modifier_mask(MASK_BUTTON3);	}
		if (mask & Button4Mask)	{ set_modifier_mask(MASK_BUTTON4);	}
		if (mask & Button5Mask)	{ set_modifier_mask(MASK_BUTTON5);	}
	}
	else {
		logger(LOG_LEVEL_WARN, "%s [%u]: XQueryPointer failed to get current modifiers!\n",
				__FUNCTION__, __LINE__);

		keycode = XKeysymToKeycode(hook->ctrl.display, XK_Shift_L);
		if (keymap[keycode / 8] & (1 << (keycode % 8))) { set_modifier_mask(MASK_SHIFT_L);	}
		keycode = XKeysymToKeycode(hook->ctrl.display, XK_Shift_R);
		if (keymap[keycode / 8] & (1 << (keycode % 8))) { set_modifier_mask(MASK_SHIFT_R);	}
		keycode = XKeysymToKeycode(hook->ctrl.display, XK_Control_L);
		if (keymap[keycode / 8] & (1 << (keycode % 8))) { set_modifier_mask(MASK_CTRL_L);	}
		keycode = XKeysymToKeycode(hook->ctrl.display, XK_Control_R);
		if (keymap[keycode / 8] & (1 << (keycode % 8))) { set_modifier_mask(MASK_CTRL_R);	}
		keycode = XKeysymToKeycode(hook->ctrl.display, XK_Alt_L);
		if (keymap[keycode / 8] & (1 << (keycode % 8))) { set_modifier_mask(MASK_ALT_L);	}
		keycode = XKeysymToKeycode(hook->ctrl.display, XK_Alt_R);
		if (keymap[keycode / 8] & (1 << (keycode % 8))) { set_modifier_mask(MASK_ALT_R);	}
		keycode = XKeysymToKeycode(hook->ctrl.display, XK_Super_L);
		if (keymap[keycode / 8] & (1 << (keycode % 8))) { set_modifier_mask(MASK_META_L);	}
		keycode = XKeysymToKeycode(hook->ctrl.display, XK_Super_R);
		if (keymap[keycode / 8] & (1 << (keycode % 8))) { set_modifier_mask(MASK_META_R);	}
	}

	initialize_locks();
}

void hook_event_proc(XPointer closeure, XRecordInterceptData *recorded_data) {
	uint64_t timestamp = (uint64_t) recorded_data->server_time;

	if (recorded_data->category == XRecordStartOfData) {
		// Populate the hook start event.
		event.time = timestamp;
		event.reserved = 0x00;

		event.type = EVENT_HOOK_ENABLED;
		event.mask = 0x00;

		// Fire the hook start event.
		dispatch_event(&event);
	}
	else if (recorded_data->category == XRecordEndOfData) {
		// Populate the hook stop event.
		event.time = timestamp;
		event.reserved = 0x00;

		event.type = EVENT_HOOK_DISABLED;
		event.mask = 0x00;

		// Fire the hook stop event.
		dispatch_event(&event);
	}
	else if (recorded_data->category == XRecordFromServer || recorded_data->category == XRecordFromClient) {
		// Get XRecord data.
		XRecordDatum *data = (XRecordDatum *) recorded_data->data;

		if (data->type == KeyPress) {
			// The X11 KeyCode associated with this event.
			KeyCode keycode = (KeyCode) data->event.u.u.detail;
            KeySym keysym = 0x00;
			#if defined(USE_XKBCOMMON)
		   	if (state != NULL) {
				keysym = xkb_state_key_get_one_sym(state, keycode);
			}
			#else
			keysym = keycode_to_keysym(keycode, data->event.u.keyButtonPointer.state);
			#endif

			unsigned short int scancode = keycode_to_scancode(keycode);

			// TODO If you have a better suggestion for this ugly, let me know.
			if		(scancode == VC_SHIFT_L)		{ set_modifier_mask(MASK_SHIFT_L);		}
			else if (scancode == VC_SHIFT_R)		{ set_modifier_mask(MASK_SHIFT_R);		}
			else if (scancode == VC_CONTROL_L)		{ set_modifier_mask(MASK_CTRL_L);		}
			else if (scancode == VC_CONTROL_R)		{ set_modifier_mask(MASK_CTRL_R);		}
			else if (scancode == VC_ALT_L)			{ set_modifier_mask(MASK_ALT_L);		}
			else if (scancode == VC_ALT_R)			{ set_modifier_mask(MASK_ALT_R);		}
			else if (scancode == VC_META_L)			{ set_modifier_mask(MASK_META_L);		}
			else if (scancode == VC_META_R)			{ set_modifier_mask(MASK_META_R);		}
			xkb_state_update_key(state, keycode, XKB_KEY_DOWN);
			initialize_locks();

			if ((get_modifiers() & MASK_NUM_LOCK) == 0) {
                switch (scancode) {
					case VC_KP_SEPARATOR:
					case VC_KP_1:
					case VC_KP_2:
					case VC_KP_3:
					case VC_KP_4:
					case VC_KP_5:
					case VC_KP_6:
					case VC_KP_7:
					case VC_KP_8:
					case VC_KP_0:
					case VC_KP_9:
						scancode |= 0xEE00;
						break;
                }
			}

			// Populate key pressed event.
			event.time = timestamp;
			event.reserved = 0x00;

			event.type = EVENT_KEY_PRESSED;
			event.mask = get_modifiers();

			event.data.keyboard.keycode = scancode;
			event.data.keyboard.rawcode = keysym;
			event.data.keyboard.keychar = CHAR_UNDEFINED;

			logger(LOG_LEVEL_INFO,	"%s [%u]: Key %#X pressed. (%#X)\n",
					__FUNCTION__, __LINE__, event.data.keyboard.keycode, event.data.keyboard.rawcode);

			// Fire key pressed event.
			dispatch_event(&event);

			// If the pressed event was not consumed...
			if (event.reserved ^ 0x01) {
				uint16_t buffer[2];
			    size_t count =  0;

				// Check to make sure the key is printable.
				#ifdef USE_XKBCOMMON
				if (state != NULL) {
					count = keycode_to_unicode(state, keycode, buffer, sizeof(buffer) / sizeof(uint16_t));
				}
				#else
				count = keysym_to_unicode(keysym, buffer, sizeof(buffer) / sizeof(uint16_t));
				#endif

				unsigned int i; 
				for (i = 0; i < count; i++) {
					// Populate key typed event.
					event.time = timestamp;
					event.reserved = 0x00;

					event.type = EVENT_KEY_TYPED;
					event.mask = get_modifiers();

					event.data.keyboard.keycode = VC_UNDEFINED;
					event.data.keyboard.rawcode = keysym;
					event.data.keyboard.keychar = buffer[i];

					logger(LOG_LEVEL_INFO,	"%s [%u]: Key %#X typed. (%lc)\n",
							__FUNCTION__, __LINE__, event.data.keyboard.keycode, (uint16_t) event.data.keyboard.keychar);

					// Fire key typed event.
					dispatch_event(&event);
				}
			}
		}
		else if (data->type == KeyRelease) {
			// The X11 KeyCode associated with this event.
			KeyCode keycode = (KeyCode) data->event.u.u.detail;
			KeySym keysym = 0x00;
			#ifdef USE_XKBCOMMON
			if (state != NULL) {
				keysym = xkb_state_key_get_one_sym(state, keycode);
			}
			#else
			keysym = keycode_to_keysym(keycode, data->event.u.keyButtonPointer.state);
			#endif

			unsigned short int scancode = keycode_to_scancode(keycode);

			// TODO If you have a better suggestion for this ugly, let me know.
			if		(scancode == VC_SHIFT_L)		{ unset_modifier_mask(MASK_SHIFT_L);		}
			else if (scancode == VC_SHIFT_R)		{ unset_modifier_mask(MASK_SHIFT_R);		}
			else if (scancode == VC_CONTROL_L)		{ unset_modifier_mask(MASK_CTRL_L);			}
			else if (scancode == VC_CONTROL_R)		{ unset_modifier_mask(MASK_CTRL_R);			}
			else if (scancode == VC_ALT_L)			{ unset_modifier_mask(MASK_ALT_L);			}
			else if (scancode == VC_ALT_R)			{ unset_modifier_mask(MASK_ALT_R);			}
			else if (scancode == VC_META_L)			{ unset_modifier_mask(MASK_META_L);			}
			else if (scancode == VC_META_R)			{ unset_modifier_mask(MASK_META_R);			}
			xkb_state_update_key(state, keycode, XKB_KEY_UP);
			initialize_locks();

			if ((get_modifiers() & MASK_NUM_LOCK) == 0) {
                switch (scancode) {
					case VC_KP_SEPARATOR:
					case VC_KP_1:
					case VC_KP_2:
					case VC_KP_3:
					case VC_KP_4:
					case VC_KP_5:
					case VC_KP_6:
					case VC_KP_7:
					case VC_KP_8:
					case VC_KP_0:
					case VC_KP_9:
						scancode |= 0xEE00;
						break;
                }
			}

			// Populate key released event.
			event.time = timestamp;
			event.reserved = 0x00;

			event.type = EVENT_KEY_RELEASED;
			event.mask = get_modifiers();

			event.data.keyboard.keycode = scancode;
			event.data.keyboard.rawcode = keysym;
			event.data.keyboard.keychar = CHAR_UNDEFINED;

			logger(LOG_LEVEL_INFO, "%s [%u]: Key %#X released. (%#X)\n",
					__FUNCTION__, __LINE__, event.data.keyboard.keycode, event.data.keyboard.rawcode);

			// Fire key released event.
			dispatch_event(&event);
		}
		else if (data->type == ButtonPress) {
			// X11 handles wheel events as button events.
			if (data->event.u.u.detail == WheelUp || data->event.u.u.detail == WheelDown
					|| data->event.u.u.detail == WheelLeft || data->event.u.u.detail == WheelRight) {

				// Reset the click count and previous button.
				hook->input.mouse.click.count = 1;
				hook->input.mouse.click.button = MOUSE_NOBUTTON;

				/* Scroll wheel release events.
				 * Scroll type: WHEEL_UNIT_SCROLL
				 * Scroll amount: 3 unit increments per notch
				 * Units to scroll: 3 unit increments
				 * Vertical unit increment: 15 pixels
				 */

				// Populate mouse wheel event.
				event.time = timestamp;
				event.reserved = 0x00;

				event.type = EVENT_MOUSE_WHEEL;
				event.mask = get_modifiers();

				event.data.wheel.clicks = hook->input.mouse.click.count;
				event.data.wheel.x = data->event.u.keyButtonPointer.rootX;
				event.data.wheel.y = data->event.u.keyButtonPointer.rootY;

				#if defined(USE_XINERAMA) || defined(USE_XRANDR)
				uint8_t count;
				screen_data *screens = hook_create_screen_info(&count);
				if (count > 1) {
					event.data.wheel.x -= screens[0].x;
					event.data.wheel.y -= screens[0].y;
				}

				if (screens != NULL) {
					free(screens);
				}
				#endif

				/* X11 does not have an API call for acquiring the mouse scroll type.  This
				 * maybe part of the XInput2 (XI2) extention but I will wont know until it
				 * is available on my platform.  For the time being we will just use the
				 * unit scroll value.
				 */
				event.data.wheel.type = WHEEL_UNIT_SCROLL;

				/* Some scroll wheel properties are available via the new XInput2 (XI2)
				 * extension.  Unfortunately the extension is not available on my
				 * development platform at this time.  For the time being we will just
				 * use the Windows default value of 3.
				 */
				event.data.wheel.amount = 3;

				if (data->event.u.u.detail == WheelUp || data->event.u.u.detail == WheelLeft) {
					// Wheel Rotated Up and Away.
					event.data.wheel.rotation = -1;
				}
				else { // data->event.u.u.detail == WheelDown
					// Wheel Rotated Down and Towards.
					event.data.wheel.rotation = 1;
				}

				if (data->event.u.u.detail == WheelUp || data->event.u.u.detail == WheelDown) {
					// Wheel Rotated Up or Down.
					event.data.wheel.direction = WHEEL_VERTICAL_DIRECTION;
				}
				else { // data->event.u.u.detail == WheelLeft || data->event.u.u.detail == WheelRight
					// Wheel Rotated Left or Right.
					event.data.wheel.direction = WHEEL_HORIZONTAL_DIRECTION;
				}

				logger(LOG_LEVEL_INFO,	"%s [%u]: Mouse wheel type %u, rotated %i units in the %u direction at %u, %u.\n",
						__FUNCTION__, __LINE__, event.data.wheel.type,
						event.data.wheel.amount * event.data.wheel.rotation,
                        event.data.wheel.direction,
						event.data.wheel.x, event.data.wheel.y);

				// Fire mouse wheel event.
				dispatch_event(&event);
			}
			else {
				/* This information is all static for X11, its up to the WM to
				 * decide how to interpret the wheel events.
				 */
				uint16_t button = MOUSE_NOBUTTON;
				switch (data->event.u.u.detail) {
					// FIXME This should use a lookup table to handle button remapping.
					case Button1:
						button = MOUSE_BUTTON1;
						set_modifier_mask(MASK_BUTTON1);
						break;

					case Button2:
						button = MOUSE_BUTTON2;
						set_modifier_mask(MASK_BUTTON2);
						break;

					case Button3:
						button = MOUSE_BUTTON3;
						set_modifier_mask(MASK_BUTTON3);
						break;

					case XButton1:
						button = MOUSE_BUTTON4;
						set_modifier_mask(MASK_BUTTON5);
						break;

					case XButton2:
						button = MOUSE_BUTTON5;
						set_modifier_mask(MASK_BUTTON5);
						break;

					default:
						// Do not set modifier masks past button MASK_BUTTON5.
						break;
				}


				// Track the number of clicks, the button must match the previous button.
				if (button == hook->input.mouse.click.button && (long int) (timestamp - hook->input.mouse.click.time) <= hook_get_multi_click_time()) {
					if (hook->input.mouse.click.count < USHRT_MAX) {
						hook->input.mouse.click.count++;
					}
					else {
						logger(LOG_LEVEL_WARN, "%s [%u]: Click count overflow detected!\n",
								__FUNCTION__, __LINE__);
					}
				}
				else {
					// Reset the click count.
					hook->input.mouse.click.count = 1;

					// Set the previous button.
					hook->input.mouse.click.button = button;
				}

				// Save this events time to calculate the hook->input.mouse.click.count.
				hook->input.mouse.click.time = timestamp;


				// Populate mouse pressed event.
				event.time = timestamp;
				event.reserved = 0x00;

				event.type = EVENT_MOUSE_PRESSED;
				event.mask = get_modifiers();

				event.data.mouse.button = button;
				event.data.mouse.clicks = hook->input.mouse.click.count;
				event.data.mouse.x = data->event.u.keyButtonPointer.rootX;
				event.data.mouse.y = data->event.u.keyButtonPointer.rootY;

				#if defined(USE_XINERAMA) || defined(USE_XRANDR)
				uint8_t count;
				screen_data *screens = hook_create_screen_info(&count);
				if (count > 1) {
					event.data.mouse.x -= screens[0].x;
					event.data.mouse.y -= screens[0].y;
				}

				if (screens != NULL) {
					free(screens);
				}
				#endif

				logger(LOG_LEVEL_INFO,	"%s [%u]: Button %u  pressed %u time(s). (%u, %u)\n",
						__FUNCTION__, __LINE__, event.data.mouse.button, event.data.mouse.clicks,
						event.data.mouse.x, event.data.mouse.y);

				// Fire mouse pressed event.
				dispatch_event(&event);
			}
		}
		else if (data->type == ButtonRelease) {
			// X11 handles wheel events as button events.
			if (data->event.u.u.detail != WheelUp && data->event.u.u.detail != WheelDown) {
				/* This information is all static for X11, its up to the WM to
				 * decide how to interpret the wheel events.
				 */
				uint16_t button = MOUSE_NOBUTTON;
				switch (data->event.u.u.detail) {
					// FIXME This should use a lookup table to handle button remapping.
					case Button1:
						button = MOUSE_BUTTON1;
						unset_modifier_mask(MASK_BUTTON1);
						break;

					case Button2:
						button = MOUSE_BUTTON2;
						unset_modifier_mask(MASK_BUTTON2);
						break;

					case Button3:
						button = MOUSE_BUTTON3;
						unset_modifier_mask(MASK_BUTTON3);
						break;

					case XButton1:
						button = MOUSE_BUTTON4;
						unset_modifier_mask(MASK_BUTTON5);
						break;

					case XButton2:
						button = MOUSE_BUTTON5;
						unset_modifier_mask(MASK_BUTTON5);
						break;

					default:
						// Do not set modifier masks past button MASK_BUTTON5.
						break;
				}

				// Populate mouse released event.
				event.time = timestamp;
				event.reserved = 0x00;

				event.type = EVENT_MOUSE_RELEASED;
				event.mask = get_modifiers();

				event.data.mouse.button = button;
				event.data.mouse.clicks = hook->input.mouse.click.count;
				event.data.mouse.x = data->event.u.keyButtonPointer.rootX;
				event.data.mouse.y = data->event.u.keyButtonPointer.rootY;

				#if defined(USE_XINERAMA) || defined(USE_XRANDR)
				uint8_t count;
				screen_data *screens = hook_create_screen_info(&count);
				if (count > 1) {
					event.data.mouse.x -= screens[0].x;
					event.data.mouse.y -= screens[0].y;
				}

				if (screens != NULL) {
					free(screens);
				}
				#endif

				logger(LOG_LEVEL_INFO,	"%s [%u]: Button %u released %u time(s). (%u, %u)\n",
						__FUNCTION__, __LINE__, event.data.mouse.button,
						event.data.mouse.clicks,
						event.data.mouse.x, event.data.mouse.y);

				// Fire mouse released event.
				dispatch_event(&event);

				// If the pressed event was not consumed...
				if (event.reserved ^ 0x01 && hook->input.mouse.is_dragged != true) {
					// Populate mouse clicked event.
					event.time = timestamp;
					event.reserved = 0x00;

					event.type = EVENT_MOUSE_CLICKED;
					event.mask = get_modifiers();

					event.data.mouse.button = button;
					event.data.mouse.clicks = hook->input.mouse.click.count;
					event.data.mouse.x = data->event.u.keyButtonPointer.rootX;
					event.data.mouse.y = data->event.u.keyButtonPointer.rootY;

					#if defined(USE_XINERAMA) || defined(USE_XRANDR)
					uint8_t count;
					screen_data *screens = hook_create_screen_info(&count);
					if (count > 1) {
						event.data.mouse.x -= screens[0].x;
						event.data.mouse.y -= screens[0].y;
					}

					if (screens != NULL) {
						free(screens);
					}
					#endif

					logger(LOG_LEVEL_INFO,	"%s [%u]: Button %u clicked %u time(s). (%u, %u)\n",
							__FUNCTION__, __LINE__, event.data.mouse.button,
							event.data.mouse.clicks,
							event.data.mouse.x, event.data.mouse.y);

					// Fire mouse clicked event.
					dispatch_event(&event);
				}

				// Reset the number of clicks.
				if (button == hook->input.mouse.click.button && (long int) (event.time - hook->input.mouse.click.time) > hook_get_multi_click_time()) {
					// Reset the click count.
					hook->input.mouse.click.count = 0;
				}
			}
		}
		else if (data->type == MotionNotify) {
			// Reset the click count.
			if (hook->input.mouse.click.count != 0 && (long int) (timestamp - hook->input.mouse.click.time) > hook_get_multi_click_time()) {
				hook->input.mouse.click.count = 0;
			}

			// Populate mouse move event.
			event.time = timestamp;
			event.reserved = 0x00;

			event.mask = get_modifiers();

			// Check the upper half of virtual modifiers for non-zero
			// values and set the mouse dragged flag.
			hook->input.mouse.is_dragged = (event.mask >> 8 > 0);
			if (hook->input.mouse.is_dragged) {
				// Create Mouse Dragged event.
				event.type = EVENT_MOUSE_DRAGGED;
			}
			else {
				// Create a Mouse Moved event.
				event.type = EVENT_MOUSE_MOVED;
			}

			event.data.mouse.button = MOUSE_NOBUTTON;
			event.data.mouse.clicks = hook->input.mouse.click.count;
			event.data.mouse.x = data->event.u.keyButtonPointer.rootX;
			event.data.mouse.y = data->event.u.keyButtonPointer.rootY;

			#if defined(USE_XINERAMA) || defined(USE_XRANDR)
			uint8_t count;
			screen_data *screens = hook_create_screen_info(&count);
			if (count > 1) {
				event.data.mouse.x -= screens[0].x;
				event.data.mouse.y -= screens[0].y;
			}

			if (screens != NULL) {
				free(screens);
			}
			#endif

			logger(LOG_LEVEL_INFO,	"%s [%u]: Mouse %s to %i, %i. (%#X)\n",
					__FUNCTION__, __LINE__, hook->input.mouse.is_dragged ? "dragged" : "moved",
					event.data.mouse.x, event.data.mouse.y, event.mask);

			// Fire mouse move event.
			dispatch_event(&event);
		}
		else {
			// In theory this *should* never execute.
			logger(LOG_LEVEL_DEBUG,	"%s [%u]: Unhandled X11 event: %#X.\n",
					__FUNCTION__, __LINE__, (unsigned int) data->type);
		}
	}
	else {
		logger(LOG_LEVEL_WARN,	"%s [%u]: Unhandled X11 hook category! (%#X)\n",
				__FUNCTION__, __LINE__, recorded_data->category);
	}

	// TODO There is no way to consume the XRecord event.

	XRecordFreeData(recorded_data);
}


static inline bool enable_key_repeate() {
	// Attempt to setup detectable autorepeat.
	// NOTE: is_auto_repeat is NOT stdbool!
	Bool is_auto_repeat = False;
	#ifdef USE_XKB
	// Enable detectable auto-repeat.
	XkbSetDetectableAutoRepeat(hook->ctrl.display, True, &is_auto_repeat);
	#else
	XAutoRepeatOn(hook->ctrl.display);

	XKeyboardState kb_state;
	XGetKeyboardControl(hook->ctrl.display, &kb_state);

	is_auto_repeat = (kb_state.global_auto_repeat == AutoRepeatModeOn);
	#endif

	return is_auto_repeat;
}


static inline int xrecord_block() {
	int status = IOHOOK_FAILURE;

	// Save the data display associated with this hook so it is passed to each event.
	//XPointer closeure = (XPointer) (ctrl_display);
	XPointer closeure = NULL;

	#ifdef USE_XRECORD_ASYNC
	// Async requires that we loop so that our thread does not return.
	if (XRecordEnableContextAsync(hook->data.display, context, hook_event_proc, closeure) != 0) {
		// Time in MS to sleep the runloop.
		int timesleep = 100;

		// Allow the thread loop to block.
		pthread_mutex_lock(&hook_xrecord_mutex);
		running = true;

		do {
			// Unlock the mutex from the previous iteration.
			pthread_mutex_unlock(&hook_xrecord_mutex);

			XRecordProcessReplies(hook->data.display);

			// Prevent 100% CPU utilization.
			struct timeval tv;
			gettimeofday(&tv, NULL);

			struct timespec ts;
			ts.tv_sec = time(NULL) + timesleep / 1000;
			ts.tv_nsec = tv.tv_usec * 1000 + 1000 * 1000 * (timesleep % 1000);
			ts.tv_sec += ts.tv_nsec / (1000 * 1000 * 1000);
			ts.tv_nsec %= (1000 * 1000 * 1000);

			pthread_mutex_lock(&hook_xrecord_mutex);
			pthread_cond_timedwait(&hook_xrecord_cond, &hook_xrecord_mutex, &ts);
		} while (running);

		// Unlock after loop exit.
		pthread_mutex_unlock(&hook_xrecord_mutex);

		// Set the exit status.
		status = NULL;
	}
	#else
	// Sync blocks until XRecordDisableContext() is called.
	if (XRecordEnableContext(hook->data.display, hook->ctrl.context, hook_event_proc, closeure) != 0) {
		status = IOHOOK_SUCCESS;
	}
	#endif
	else {
		logger(LOG_LEVEL_ERROR,	"%s [%u]: XRecordEnableContext failure!\n",
			__FUNCTION__, __LINE__);

		#ifdef USE_XRECORD_ASYNC
		// Reset the running state.
		pthread_mutex_lock(&hook_xrecord_mutex);
		running = false;
		pthread_mutex_unlock(&hook_xrecord_mutex);
		#endif

		// Set the exit status.
		status = IOHOOK_ERROR_X_RECORD_ENABLE_CONTEXT;
	}

	return status;
}

static int xrecord_alloc() {
	int status = IOHOOK_FAILURE;

	// Make sure the data display is synchronized to prevent late event delivery!
	// See Bug 42356 for more information.
	// https://bugs.freedesktop.org/show_bug.cgi?id=42356#c4
	XSynchronize(hook->data.display, True);

	// Setup XRecord range.
	XRecordClientSpec clients = XRecordAllClients;

	hook->data.range = XRecordAllocRange();
	if (hook->data.range != NULL) {
		logger(LOG_LEVEL_DEBUG,	"%s [%u]: XRecordAllocRange successful.\n",
				__FUNCTION__, __LINE__);

		hook->data.range->device_events.first = KeyPress;
		hook->data.range->device_events.last = MotionNotify;

		// Note that the documentation for this function is incorrect,
		// hook->data.display should be used!
		// See: http://www.x.org/releases/X11R7.6/doc/libXtst/recordlib.txt
		hook->ctrl.context = XRecordCreateContext(hook->data.display, XRecordFromServerTime, &clients, 1, &hook->data.range, 1);
		if (hook->ctrl.context != 0) {
			logger(LOG_LEVEL_DEBUG,	"%s [%u]: XRecordCreateContext successful.\n",
					__FUNCTION__, __LINE__);

			// Block until hook_stop() is called.
			status = xrecord_block();

			// Free up the context if it was set.
			XRecordFreeContext(hook->data.display, hook->ctrl.context);
		}
		else {
			logger(LOG_LEVEL_ERROR,	"%s [%u]: XRecordCreateContext failure!\n",
					__FUNCTION__, __LINE__);

			// Set the exit status.
			status = IOHOOK_ERROR_X_RECORD_CREATE_CONTEXT;
		}

		// Free the XRecord range.
		XFree(hook->data.range);
	}
	else {
		logger(LOG_LEVEL_ERROR,	"%s [%u]: XRecordAllocRange failure!\n",
				__FUNCTION__, __LINE__);

		// Set the exit status.
		status = IOHOOK_ERROR_X_RECORD_ALLOC_RANGE;
	}

	return status;
}

static int xrecord_query() {
	int status = IOHOOK_FAILURE;

	// Check to make sure XRecord is installed and enabled.
	int major, minor;
	if (XRecordQueryVersion(hook->ctrl.display, &major, &minor) != 0) {
		logger(LOG_LEVEL_INFO,	"%s [%u]: XRecord version: %i.%i.\n",
				__FUNCTION__, __LINE__, major, minor);

		status = xrecord_alloc();
	}
	else {
		logger(LOG_LEVEL_ERROR,	"%s [%u]: XRecord is not currently available!\n",
				__FUNCTION__, __LINE__);

		status = IOHOOK_ERROR_X_RECORD_NOT_FOUND;
	}

	return status;
}

static int xrecord_start() {
	int status = IOHOOK_FAILURE;

	// Open the control display for XRecord.
	hook->ctrl.display = XOpenDisplay(NULL);

	// Open a data display for XRecord.
	// NOTE This display must be opened on the same thread as XRecord.
	hook->data.display = XOpenDisplay(NULL);
	if (hook->ctrl.display != NULL && hook->data.display != NULL) {
		logger(LOG_LEVEL_DEBUG,	"%s [%u]: XOpenDisplay successful.\n",
				__FUNCTION__, __LINE__);

		bool is_auto_repeat = enable_key_repeate();
		if (is_auto_repeat) {
			logger(LOG_LEVEL_DEBUG,	"%s [%u]: Successfully enabled detectable autorepeat.\n",
					__FUNCTION__, __LINE__);
		}
		else {
			logger(LOG_LEVEL_WARN,	"%s [%u]: Could not enable detectable auto-repeat!\n",
					__FUNCTION__, __LINE__);
		}

		 #if defined(USE_XKBCOMMON)
		// Open XCB Connection
		hook->input.connection = XGetXCBConnection(hook->ctrl.display);
		int xcb_status = xcb_connection_has_error(hook->input.connection);
		if (xcb_status <= 0) {
			// Initialize xkbcommon context.
			struct xkb_context *context = xkb_context_new(XKB_CONTEXT_NO_FLAGS);

			if (context != NULL) {
				hook->input.context = xkb_context_ref(context);
			}
			else {
				logger(LOG_LEVEL_ERROR,	"%s [%u]: xkb_context_new failure!\n",
						__FUNCTION__, __LINE__);
			}
		}
		else {
			logger(LOG_LEVEL_ERROR,	"%s [%u]: xcb_connect failure! (%d)\n",
					__FUNCTION__, __LINE__, xcb_status);
		}
		#endif

		#ifdef USE_XKBCOMMON
		state = create_xkb_state(hook->input.context, hook->input.connection);
		#endif

		// Initialize starting modifiers.
		initialize_modifiers();

		status = xrecord_query();

		#ifdef USE_XKBCOMMON
		if (state != NULL) {
			destroy_xkb_state(state);
		}

		if (hook->input.context != NULL) {
			xkb_context_unref(hook->input.context);
			hook->input.context = NULL;
		}

		if (hook->input.connection != NULL) {
			// xcb_disconnect(hook->input.connection);
			hook->input.connection = NULL;
		}
		#endif
	}
	else {
		logger(LOG_LEVEL_ERROR,	"%s [%u]: XOpenDisplay failure!\n",
				__FUNCTION__, __LINE__);

		status = IOHOOK_ERROR_X_OPEN_DISPLAY;
	}

	// Close down the XRecord data display.
	if (hook->data.display != NULL) {
		XCloseDisplay(hook->data.display);
		hook->data.display = NULL;
	}

	// Close down the XRecord control display.
	if (hook->ctrl.display) {
		XCloseDisplay(hook->ctrl.display);
		hook->ctrl.display = NULL;
	}

	return status;
}

IOHOOK_API int hook_run() {
	int status = IOHOOK_FAILURE;

	// Hook data for future cleanup.
	hook = malloc(sizeof(hook_info));
	if (hook != NULL) {
		hook->input.mask = 0x0000;
		hook->input.mouse.is_dragged = false;
		hook->input.mouse.click.count = 0;
		hook->input.mouse.click.time = 0;
		hook->input.mouse.click.button = MOUSE_NOBUTTON;

		status = xrecord_start();

		// Free data associated with this hook.
		free(hook);
		hook = NULL;
	}
	else {
		logger(LOG_LEVEL_ERROR,	"%s [%u]: Failed to allocate memory for hook structure!\n",
				__FUNCTION__, __LINE__);

		status = IOHOOK_ERROR_OUT_OF_MEMORY;
	}

	logger(LOG_LEVEL_DEBUG,	"%s [%u]: Something, something, something, complete.\n",
			__FUNCTION__, __LINE__);

	return status;
}

IOHOOK_API int hook_stop() {
	int status = IOHOOK_FAILURE;

	if (hook != NULL && hook->ctrl.display != NULL && hook->ctrl.context != 0) {
		// We need to make sure the context is still valid.
		XRecordState *state = malloc(sizeof(XRecordState));
		if (state != NULL) {
			if (XRecordGetContext(hook->ctrl.display, hook->ctrl.context, &state) != 0) {
				// Try to exit the thread naturally.
				if (state->enabled && XRecordDisableContext(hook->ctrl.display, hook->ctrl.context) != 0) {
					#ifdef USE_XRECORD_ASYNC
					pthread_mutex_lock(&hook_xrecord_mutex);
					running = false;
					pthread_cond_signal(&hook_xrecord_cond);
					pthread_mutex_unlock(&hook_xrecord_mutex);
					#endif

					// See Bug 42356 for more information.
					// https://bugs.freedesktop.org/show_bug.cgi?id=42356#c4
					//XFlush(hook->ctrl.display);
					XSync(hook->ctrl.display, False);
					if (hook->ctrl.display) {
						XCloseDisplay(hook->ctrl.display);
						hook->ctrl.display = NULL;
					}

					status = IOHOOK_SUCCESS;
				}
			}
			else {
				logger(LOG_LEVEL_ERROR,	"%s [%u]: XRecordGetContext failure!\n",
						__FUNCTION__, __LINE__);

				status = IOHOOK_ERROR_X_RECORD_GET_CONTEXT;
			}

			free(state);
		}
		else {
			logger(LOG_LEVEL_ERROR,	"%s [%u]: Failed to allocate memory for XRecordState!\n",
					__FUNCTION__, __LINE__);

			status = IOHOOK_ERROR_OUT_OF_MEMORY;
		}

		return status;
	}

	logger(LOG_LEVEL_DEBUG,	"%s [%u]: Status: %#X.\n",
			__FUNCTION__, __LINE__, status);

	return status;
}
