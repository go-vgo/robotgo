
#ifdef HAVE_CONFIG_H
#include <config.h>
#endif

#include <stdio.h>
#include "../iohook.h"
#include <windows.h>

// #include "logger.h"
#include "input.h"

// Some buggy versions of MinGW and MSys do not include these constants in winuser.h.
#ifndef MAPVK_VK_TO_VSC
#define MAPVK_VK_TO_VSC			0
#define MAPVK_VSC_TO_VK			1
#define MAPVK_VK_TO_CHAR		2
#define MAPVK_VSC_TO_VK_EX		3
#endif
// Some buggy versions of MinGW and MSys only define this value for Windows
// versions >= 0x0600 (Windows Vista) when it should be 0x0500 (Windows 2000).
#ifndef MAPVK_VK_TO_VSC_EX
#define MAPVK_VK_TO_VSC_EX		4
#endif

#ifndef KEYEVENTF_SCANCODE
#define KEYEVENTF_EXTENDEDKEY	0x0001
#define KEYEVENTF_KEYUP			0x0002
#define	KEYEVENTF_UNICODE		0x0004
#define KEYEVENTF_SCANCODE		0x0008
#endif

#ifndef KEYEVENTF_KEYDOWN
#define KEYEVENTF_KEYDOWN		0x0000
#endif

#define MAX_WINDOWS_COORD_VALUE 65535

static UINT keymask_lookup[8] = {
	VK_LSHIFT,
	VK_LCONTROL,
	VK_LWIN,
	VK_LMENU,

	VK_RSHIFT,
	VK_RCONTROL,
	VK_RWIN,
	VK_RMENU
};

IOHOOK_API void hook_post_event(iohook_event * const event) {
	//FIXME implement multiple monitor support
	uint16_t screen_width   = GetSystemMetrics( SM_CXSCREEN );
	uint16_t screen_height  = GetSystemMetrics( SM_CYSCREEN );

	unsigned char events_size = 0, events_max = 28;
	INPUT *events = malloc(sizeof(INPUT) * events_max);

	if (event->mask & (MASK_SHIFT | MASK_CTRL | MASK_META | MASK_ALT)) {
		unsigned int i;
		for (i = 0; i < sizeof(keymask_lookup) / sizeof(UINT); i++) {
			if (event->mask & 1 << i) {
				events[events_size].type = INPUT_KEYBOARD;
				events[events_size].ki.wVk = keymask_lookup[i];
				events[events_size].ki.dwFlags = KEYEVENTF_KEYDOWN;
				events[events_size].ki.time = 0; // Use current system time.
				events_size++;
			}
		}
	}

	if (event->mask & (MASK_BUTTON1 | MASK_BUTTON2 | MASK_BUTTON3 | MASK_BUTTON4 | MASK_BUTTON5)) {
		events[events_size].type = INPUT_MOUSE;
		events[events_size].mi.dx = 0;	// Relative mouse movement due to
		events[events_size].mi.dy = 0;	// MOUSEEVENTF_ABSOLUTE not being set.
		events[events_size].mi.mouseData = 0x00;
		events[events_size].mi.time = 0; // Use current system time.

		if (event->mask & MASK_BUTTON1) {
			events[events_size].mi.mouseData |= MOUSEEVENTF_LEFTDOWN;
		}

		if (event->mask & MASK_BUTTON2) {
			events[events_size].mi.mouseData |= MOUSEEVENTF_RIGHTDOWN;
		}

		if (event->mask & MASK_BUTTON3) {
			events[events_size].mi.mouseData |= MOUSEEVENTF_MIDDLEDOWN;
		}

		if (event->mask & MASK_BUTTON4) {
			events[events_size].mi.mouseData = XBUTTON1;
			events[events_size].mi.mouseData |= MOUSEEVENTF_XDOWN;
		}

		if (event->mask & MASK_BUTTON5) {
			events[events_size].mi.mouseData = XBUTTON2;
			events[events_size].mi.dwFlags |= MOUSEEVENTF_XDOWN;
		}

		events_size++;
	}


	switch (event->type) {
		case EVENT_KEY_PRESSED:
			events[events_size].ki.wVk = scancode_to_keycode(event->data.keyboard.keycode);
			if (events[events_size].ki.wVk != 0x0000) {
				events[events_size].type = INPUT_KEYBOARD;
				events[events_size].ki.dwFlags = KEYEVENTF_KEYDOWN; // |= KEYEVENTF_SCANCODE;
				events[events_size].ki.wScan = 0; // event->data.keyboard.keycode;
				events[events_size].ki.time = 0; // GetSystemTime()
				events_size++;
			}
			else {
				logger(LOG_LEVEL_INFO, "%s [%u]: Unable to lookup scancode: %li\n",
						__FUNCTION__, __LINE__,
						event->data.keyboard.keycode);
			}
			break;

		case EVENT_KEY_RELEASED:
			events[events_size].ki.wVk = scancode_to_keycode(event->data.keyboard.keycode);
			if (events[events_size].ki.wVk != 0x0000) {
				events[events_size].type = INPUT_KEYBOARD;
				events[events_size].ki.dwFlags = KEYEVENTF_KEYUP; // |= KEYEVENTF_SCANCODE;
				events[events_size].ki.wVk = scancode_to_keycode(event->data.keyboard.keycode);
				events[events_size].ki.wScan = 0; // event->data.keyboard.keycode;
				events[events_size].ki.time = 0; // GetSystemTime()
				events_size++;
			}
			else {
				logger(LOG_LEVEL_INFO, "%s [%u]: Unable to lookup scancode: %li\n",
						__FUNCTION__, __LINE__,
						event->data.keyboard.keycode);
			}
			break;


		case EVENT_MOUSE_PRESSED:
			events[events_size].type = INPUT_MOUSE;
			events[events_size].mi.dwFlags = MOUSEEVENTF_XDOWN;

			switch (event->data.mouse.button) {
				case MOUSE_BUTTON1:
					events[events_size].mi.dwFlags = MOUSEEVENTF_LEFTDOWN;
					break;

				case MOUSE_BUTTON2:
					events[events_size].mi.dwFlags = MOUSEEVENTF_RIGHTDOWN;
					break;

				case MOUSE_BUTTON3:
					events[events_size].mi.dwFlags = MOUSEEVENTF_MIDDLEDOWN;
					break;

				case MOUSE_BUTTON4:
					events[events_size].mi.mouseData = XBUTTON1;
					break;

			   case MOUSE_BUTTON5:
					events[events_size].mi.mouseData = XBUTTON2;
					break;

				default:
					// Extra buttons.
					if (event->data.mouse.button > 3) {
						events[events_size].mi.mouseData = event->data.mouse.button - 3;
					}
			}

			events[events_size].mi.dx = event->data.mouse.x * (MAX_WINDOWS_COORD_VALUE / screen_width) + 1;
			events[events_size].mi.dy = event->data.mouse.y * (MAX_WINDOWS_COORD_VALUE / screen_height) + 1;

			events[events_size].mi.dwFlags |= MOUSEEVENTF_ABSOLUTE | MOUSEEVENTF_MOVE;
			events[events_size].mi.time = 0; // GetSystemTime()

			events_size++;
			break;

		case EVENT_MOUSE_RELEASED:
			events[events_size].type = INPUT_MOUSE;
			events[events_size].mi.dwFlags = MOUSEEVENTF_XUP;

			switch (event->data.mouse.button) {
				case MOUSE_BUTTON1:
					events[events_size].mi.dwFlags = MOUSEEVENTF_LEFTUP;
					break;

				case MOUSE_BUTTON2:
					events[events_size].mi.dwFlags = MOUSEEVENTF_RIGHTUP;
					break;

				case MOUSE_BUTTON3:
					events[events_size].mi.dwFlags = MOUSEEVENTF_MIDDLEUP;
					break;

				case MOUSE_BUTTON4:
					events[events_size].mi.mouseData = XBUTTON1;
					break;

			   case MOUSE_BUTTON5:
					events[events_size].mi.mouseData = XBUTTON2;
					break;

				default:
					// Extra buttons.
					if (event->data.mouse.button > 3) {
						events[events_size].mi.mouseData = event->data.mouse.button - 3;
					}
			}

			events[events_size].mi.dx = event->data.mouse.x * (MAX_WINDOWS_COORD_VALUE / screen_width) + 1;
			events[events_size].mi.dy = event->data.mouse.y * (MAX_WINDOWS_COORD_VALUE / screen_height) + 1;

			events[events_size].mi.dwFlags |= MOUSEEVENTF_ABSOLUTE | MOUSEEVENTF_MOVE;
			events[events_size].mi.time = 0; // GetSystemTime()
			events_size++;
			break;


		case EVENT_MOUSE_WHEEL:
			events[events_size].type = INPUT_MOUSE;
			events[events_size].mi.dwFlags = MOUSEEVENTF_WHEEL;

			// type, amount and rotation?
			events[events_size].mi.mouseData = event->data.wheel.amount * event->data.wheel.rotation * WHEEL_DELTA;

			events[events_size].mi.dx = event->data.wheel.x * (MAX_WINDOWS_COORD_VALUE / screen_width) + 1;
			events[events_size].mi.dy = event->data.wheel.y * (MAX_WINDOWS_COORD_VALUE / screen_height) + 1;

			events[events_size].mi.dwFlags |= MOUSEEVENTF_ABSOLUTE | MOUSEEVENTF_MOVE;
			events[events_size].mi.time = 0; // GetSystemTime()
			events_size++;
			break;


		case EVENT_MOUSE_DRAGGED:
			// The button masks are all applied with the modifier masks.

		case EVENT_MOUSE_MOVED:
			events[events_size].type = INPUT_MOUSE;
			events[events_size].mi.dwFlags = MOUSEEVENTF_MOVE;

			events[events_size].mi.dx = event->data.mouse.x * (MAX_WINDOWS_COORD_VALUE / screen_width) + 1;
			events[events_size].mi.dy = event->data.mouse.y * (MAX_WINDOWS_COORD_VALUE / screen_height) + 1;

			events[events_size].mi.dwFlags |= MOUSEEVENTF_ABSOLUTE | MOUSEEVENTF_MOVE;
			events[events_size].mi.time = 0; // GetSystemTime()
			events_size++;
			break;


		case EVENT_MOUSE_CLICKED:
		case EVENT_KEY_TYPED:
			// Ignore clicked and typed events.

		case EVENT_HOOK_ENABLED:
		case EVENT_HOOK_DISABLED:
			// Ignore hook enabled / disabled events.

		default:
			// Ignore any other garbage.
			logger(LOG_LEVEL_WARN, "%s [%u]: Ignoring post event type %#X\n",
					__FUNCTION__, __LINE__, event->type);
			break;
	}

	// Release the previously held modifier keys used to fake the event mask.
	if (event->mask & (MASK_SHIFT | MASK_CTRL | MASK_META | MASK_ALT)) {
		unsigned int i;
		for (i = 0; i < sizeof(keymask_lookup) / sizeof(UINT); i++) {
			if (event->mask & 1 << i) {
				events[events_size].type = INPUT_KEYBOARD;
				events[events_size].ki.wVk = keymask_lookup[i];
				events[events_size].ki.dwFlags = KEYEVENTF_KEYUP;
				events[events_size].ki.time = 0; // Use current system time.
				events_size++;
			}
		}
	}

	if (event->mask & (MASK_BUTTON1 | MASK_BUTTON2 | MASK_BUTTON3 | MASK_BUTTON4 | MASK_BUTTON5)) {
		events[events_size].type = INPUT_MOUSE;
		events[events_size].mi.dx = 0;	// Relative mouse movement due to
		events[events_size].mi.dy = 0;	// MOUSEEVENTF_ABSOLUTE not being set.
		events[events_size].mi.mouseData = 0x00;
		events[events_size].mi.time = 0; // Use current system time.

		// If dwFlags does not contain MOUSEEVENTF_WHEEL, MOUSEEVENTF_XDOWN, or MOUSEEVENTF_XUP,
		// then mouseData should be zero.
		// http://msdn.microsoft.com/en-us/library/windows/desktop/ms646273%28v=vs.85%29.aspx
		if (event->mask & MASK_BUTTON1) {
			events[events_size].mi.dwFlags |= MOUSEEVENTF_LEFTUP;
		}

		if (event->mask & MASK_BUTTON2) {
			events[events_size].mi.dwFlags |= MOUSEEVENTF_RIGHTUP;
		}

		if (event->mask & MASK_BUTTON3) {
			events[events_size].mi.dwFlags |= MOUSEEVENTF_MIDDLEUP;
		}

		if (event->mask & MASK_BUTTON4) {
			events[events_size].mi.mouseData = XBUTTON1;
			events[events_size].mi.dwFlags |= MOUSEEVENTF_XUP;
		}

		if (event->mask & MASK_BUTTON5) {
			events[events_size].mi.mouseData = XBUTTON2;
			events[events_size].mi.dwFlags |= MOUSEEVENTF_XUP;
		}

		events_size++;
	}

	// Create the key release input
	// memcpy(key_events + 1, key_events, sizeof(INPUT));
	// key_events[1].ki.dwFlags |= KEYEVENTF_KEYUP;

	if (! SendInput(events_size, events, sizeof(INPUT)) ) {
		logger(LOG_LEVEL_ERROR, "%s [%u]: SendInput() failed! (%#lX)\n",
				__FUNCTION__, __LINE__, (unsigned long) GetLastError());
	}

	free(events);
}
