// Copyright (c) 2016-2026 AtomAI, All rights reserved.
// 
// See the COPYRIGHT file at the top-level directory of this distribution and at
// https://github.com/go-vgo/robotgo/blob/master/LICENSE
//
// Licensed under the Apache License, Version 2.0 <LICENSE-APACHE or
// http://www.apache.org/licenses/LICENSE-2.0> 
//
// This file may not be copied, modified, or distributed
// except according to those terms.

#include "../base/deadbeef_rand_c.h"
#include "../base/microsleep.h"
#include "keypress.h"
#include "keycode_c.h"

#include <ctype.h> /* For isupper() */
#if defined(IS_MACOSX)
	#include <ApplicationServices/ApplicationServices.h>
	#import <IOKit/hidsystem/IOHIDLib.h>
	#import <IOKit/hidsystem/ev_keymap.h>
#elif defined(USE_X11)
	#include <X11/extensions/XTest.h>
	#include <X11/XKBlib.h>
	// #include "../base/xdisplay_c.h"
#endif

/* Convenience wrappers around ugly APIs. */
#if defined(IS_WINDOWS)
	HWND GetHwndByPid(DWORD dwProcessId);

	HWND getHwnd(uintptr pid, int8_t isPid) { 
		if (isPid == 0) { 
			return GetHwndByPid((DWORD)pid);
		}
		return (HWND)pid;	
	}

	// Construct lParam for WM_KEYDOWN/WM_KEYUP messages
	LPARAM makeKeyLParam(int scanCode, BOOL isExtended, BOOL isKeyUp) {
		LPARAM lParam = 1; // repeat count = 1
		lParam |= (scanCode & 0xFF) << 16; // scan code
		if (isExtended) {
			lParam |= (1 << 24); // extended key flag
		}
		if (isKeyUp) {
			lParam |= (1 << 30); // previous key state (was down)
			lParam |= (1 << 31); // transition state (being released)
		}
		return lParam;
	}

	int WIN32_KEY_EVENT_WAIT(MMKeyCode key, DWORD flags, uintptr pid) {
		int ret = win32KeyEvent(key, flags, pid, 0); 
		Sleep(DEADBEEF_RANDRANGE(0, 1));
		return ret;
	}
#elif defined(USE_X11)
	Display *XGetMainDisplay(void);

	// X_KEYCODE_EVENT returns 0 on success, 1 if XTest rejected the event.
	int X_KEYCODE_EVENT(Display *display, KeyCode code, bool is_press) {
		Bool ok = XTestFakeKeyEvent(display, code, is_press, CurrentTime);
		XSync(display, false);
		return ok ? 0 : 1;
	}

	int X_KEY_EVENT(Display *display, MMKeyCode key, bool is_press) {
		return X_KEYCODE_EVENT(display, XKeysymToKeycode(display, key), is_press);
	}

	void X_KEY_EVENT_WAIT(Display *display, MMKeyCode key, bool is_press) {
		X_KEY_EVENT(display, key, is_press);
		microsleep(DEADBEEF_UNIFORM(0.0, 0.5));
	}

	/* Keycode and shift level that produce sym in the active XKB group, lowest
	level first: 0 plain, 1 Shift, 2 ISO_Level3_Shift (AltGr), 3 both (the
	pc/complete key types). XKeysymToKeycode alone drops the level, so on a
	German layout '@' (AltGr+q) came out as Shift+q = 'Q' and '/' (Shift+7)
	as '7' (#640). Falls back to XKeysymToKeycode for syms not on the layout. */
	KeyCode X_KEYSYM_TO_KEYCODE(Display *display, KeySym sym, int *level) {
		XkbStateRec state;
		unsigned int group = 0;
		if (XkbGetState(display, XkbUseCoreKbd, &state) == Success) {
			group = state.group;
		}
		int min = 8, max = 255;
		XDisplayKeycodes(display, &min, &max);
		for (int lv = 0; lv < 4; lv++) {
			for (int kc = min; kc <= max; kc++) {
				if (XkbKeycodeToKeysym(display, (KeyCode)kc, group, lv) == sym) {
					*level = lv;
					return (KeyCode)kc;
				}
			}
		}
		*level = 0;
		return XKeysymToKeycode(display, sym);
	}

	/* Keycode of the level-3 chooser (AltGr on pc layouts), 0 if none. */
	KeyCode X_LEVEL3_KEYCODE(Display *display) {
		KeyCode kc = XKeysymToKeycode(display, XK_ISO_Level3_Shift);
		if (kc == 0) {
			kc = XKeysymToKeycode(display, XK_Mode_switch);
		}
		return kc;
	}
#endif

#if defined(IS_MACOSX)
	int SendTo(uintptr pid, CGEventRef event) {
		if (pid != 0) {
			CGEventPostToPid(pid, event);
		} else {
			CGEventPost(kCGHIDEventTap, event);
		}
		
		CFRelease(event);
		return 0;
	}

	static io_connect_t _getAuxiliaryKeyDriver(void) {
		static mach_port_t sEventDrvrRef = 0;
		mach_port_t masterPort, service, iter;
		kern_return_t kr;

		if (!sEventDrvrRef) {
			#if __ENVIRONMENT_MAC_OS_X_VERSION_MIN_REQUIRED__ < 120000
				kr = IOMasterPort(bootstrap_port, &masterPort); // waring deprecated
			#else
				kr = IOMainPort(bootstrap_port, &masterPort);
			#endif

			assert(KERN_SUCCESS == kr);
			kr = IOServiceGetMatchingServices(masterPort, IOServiceMatching(kIOHIDSystemClass), &iter);
			assert(KERN_SUCCESS == kr);

			service = IOIteratorNext(iter);
			assert(service);
			kr = IOServiceOpen(service, mach_task_self(), kIOHIDParamConnectType, &sEventDrvrRef);
			assert(KERN_SUCCESS == kr);

			IOObjectRelease(service);
			IOObjectRelease(iter);
		}
		return sEventDrvrRef;
	}
#elif defined(IS_WINDOWS)
	int win32KeyEvent(int key, MMKeyFlags flags, uintptr pid, int8_t isPid) {
		int scan = MapVirtualKey(key & 0xff, MAPVK_VK_TO_VSC);
		BOOL isExtended = FALSE;

		/* Set the scan code for extended keys */
		switch (key){
			case VK_RCONTROL:
			case VK_SNAPSHOT: /* Print Screen */
			case VK_RMENU: /* Right Alt / Alt Gr */
			case VK_PAUSE: /* Pause / Break */
			case VK_HOME:
			case VK_UP:
			case VK_PRIOR: /* Page up */
			case VK_LEFT:
			case VK_RIGHT:
			case VK_END:
			case VK_DOWN:
			case VK_NEXT: /* 'Page Down' */
			case VK_INSERT:
			case VK_DELETE:
			case VK_LWIN:
			case VK_RWIN:
			case VK_APPS: /* Application */
			case VK_VOLUME_MUTE:
			case VK_VOLUME_DOWN:
			case VK_VOLUME_UP:
			case VK_MEDIA_NEXT_TRACK:
			case VK_MEDIA_PREV_TRACK:
			case VK_MEDIA_STOP:
			case VK_MEDIA_PLAY_PAUSE:
			case VK_BROWSER_BACK:
			case VK_BROWSER_FORWARD:
			case VK_BROWSER_REFRESH:
			case VK_BROWSER_STOP:
			case VK_BROWSER_SEARCH:
			case VK_BROWSER_FAVORITES:
			case VK_BROWSER_HOME:
			case VK_LAUNCH_MAIL:
			{
				flags |= KEYEVENTF_EXTENDEDKEY;
				isExtended = TRUE;
				break;
			}
		}

		// todo: test this
		if (pid != 0) {
			HWND hwnd = getHwnd(pid, isPid);
			if (hwnd == NULL) {
				return -5; /* Window not found */
			}

			// int down = (flags == 0 ? WM_KEYDOWN : WM_KEYUP);
			BOOL isKeyUp = (flags & KEYEVENTF_KEYUP) != 0;
			UINT down = isKeyUp ? WM_KEYUP : WM_KEYDOWN;
			LPARAM lParam = makeKeyLParam(scan, isExtended, isKeyUp);
			// SendMessage(hwnd, down, key, 0);
			UINT sent = PostMessageW(hwnd, down, key, lParam);
			return sent == 1 ? 0 : (int)GetLastError();
		}

		/* Set the scan code for keyup */
		// if ( flags & KEYEVENTF_KEYUP ) {
		// 	scan |= 0x80;
		// }
		// keybd_event(key, scan, flags, 0);
		
		INPUT keyInput;

		keyInput.type = INPUT_KEYBOARD;
		keyInput.ki.wVk = key;
		keyInput.ki.wScan = scan;
		keyInput.ki.dwFlags = flags;
		keyInput.ki.time = 0;
		keyInput.ki.dwExtraInfo = 0;
		UINT sent = SendInput(1, &keyInput, sizeof(keyInput));
		if (sent == 1) {
			return 0;
		}
		DWORD err = GetLastError();
		return err != 0 ? (int)err : MM_ERR_INPUT_BLOCKED;
	}
#endif

#if defined(IS_MACOSX)
/* isModKeyCode reports whether code is a modifier key, for which
CGEventCreateKeyboardEvent builds a kCGEventFlagsChanged event. */
static bool isModKeyCode(MMKeyCode code) {
	switch (code) {
	case K_META: case K_RMETA:
	case K_ALT: case K_RALT:
	case K_CONTROL: case K_RCONTROL:
	case K_SHIFT: case K_RSHIFT:
		return true;
	}
	return false;
}
#endif

int toggleKeyCode(MMKeyCode code, const bool down, MMKeyFlags flags, uintptr pid) {
#if defined(IS_MACOSX)
	/* The media keys all have 1000 added to them to help us detect them. */
	if (code >= 1000) {
		code = code - 1000; /* Get the real keycode. */
		NXEventData event;
		kern_return_t kr;

		IOGPoint loc = { 0, 0 };
		UInt32 evtInfo = code << 16 | (down?NX_KEYDOWN:NX_KEYUP) << 8;

		bzero(&event, sizeof(NXEventData));
		event.compound.subType = NX_SUBTYPE_AUX_CONTROL_BUTTONS;
		event.compound.misc.L[0] = evtInfo;

		io_connect_t auxDriver = _getAuxiliaryKeyDriver();
		if (auxDriver == 0) {
			return -1;
		}
		kr = IOHIDPostEvent(auxDriver, NX_SYSDEFINED, loc, &event, kNXEventDataVersion, 0, FALSE);
		// assert(KERN_SUCCESS == kr);
		if (kr != KERN_SUCCESS) {
			return kr;
		}
	} else {
		CGEventSourceRef source = MMEventSourceCreate();
		CGEventRef keyEvent = CGEventCreateKeyboardEvent(source, (CGKeyCode)code, down);
		// assert(keyEvent != NULL);
		if (keyEvent == NULL) {
			if (source != NULL) { CFRelease(source); }
			return (int)kCGErrorCannotComplete;
		}

		/* Keep the type CGEventCreateKeyboardEvent chose: modifier keycodes
		get kCGEventFlagsChanged; forcing KeyDown/KeyUp on them meant the
		release was never applied and cmd/alt/shift stayed stuck down.

		A modifier release always carries an explicit mask (flags = the
		modifiers that remain held). The auto-filled mask is read from the
		HID state asynchronously, so back-to-back releases (cmd+alt) saw a
		stale state and re-asserted the modifier just released. */
		if (flags != 0 || (!down && isModKeyCode(code))) {
			CGEventSetFlags(keyEvent, (CGEventFlags) flags);
		}

		SendTo(pid, keyEvent);
		if (source != NULL) { CFRelease(source); }
	}
	return 0;
#elif defined(IS_WINDOWS)
	const DWORD dwFlags = down ? 0 : KEYEVENTF_KEYUP;

	/* Parse modifier keys. */
	if (flags & MOD_META) { WIN32_KEY_EVENT_WAIT(K_META, dwFlags, pid); }
	if (flags & MOD_ALT) { WIN32_KEY_EVENT_WAIT(K_ALT, dwFlags, pid); }
	if (flags & MOD_CONTROL) { WIN32_KEY_EVENT_WAIT(K_CONTROL, dwFlags, pid); }
	if (flags & MOD_SHIFT) { WIN32_KEY_EVENT_WAIT(K_SHIFT, dwFlags, pid); }

	return win32KeyEvent(code, dwFlags, pid, 0);
#elif defined(USE_X11)
	Display *display = XGetMainDisplay();
	if (display == NULL) {
		return -8;
	}
	const Bool is_press = down ? True : False; /* Just to be safe. */

	/* The layout decides which modifiers the keysym needs (#640). */
	int level = 0;
	KeyCode kc = X_KEYSYM_TO_KEYCODE(display, code, &level);
	if (level & 1) { flags |= MOD_SHIFT; }
	KeyCode level3 = (level & 2) ? X_LEVEL3_KEYCODE(display) : 0;
	if ((level & 2) && level3 == 0) {
		return -9; /* AltGr level, but layout has no level-3 key */
	}

	/* Release main key while its modifiers are still held, so the
	release keeps the level its press had. */
	int ret = 0;
	if (!is_press) { ret = X_KEYCODE_EVENT(display, kc, is_press); }

	/* Parse modifier keys. */
	if (flags & MOD_META) { X_KEY_EVENT_WAIT(display, K_META, is_press); }
	if (flags & MOD_ALT) { X_KEY_EVENT_WAIT(display, K_ALT, is_press); }
	if (flags & MOD_CONTROL) { X_KEY_EVENT_WAIT(display, K_CONTROL, is_press); }
	if (flags & MOD_SHIFT) { X_KEY_EVENT_WAIT(display, K_SHIFT, is_press); }
	if (level3 != 0) {
		X_KEYCODE_EVENT(display, level3, is_press);
		microsleep(DEADBEEF_UNIFORM(0.0, 0.5));
	}

	if (is_press) { ret = X_KEYCODE_EVENT(display, kc, is_press); }
	return ret;
#endif
}

// void tapKeyCode(MMKeyCode code, MMKeyFlags flags){
// 	toggleKeyCode(code, true, flags);
// 	microsleep(5.0);
// 	toggleKeyCode(code, false, flags);
// }

int toggleKey(char c, const bool down, MMKeyFlags flags, uintptr pid) {
	MMKeyCode keyCode = keyCodeForChar(c);

	/* On X11 toggleKeyCode reads the shift level off the layout; a fixed US
	shifted-character list ("~!@#...") was what mistyped '@' and '/' on a
	German layout (#640). */
	#if !defined(USE_X11)
		if (isupper(c) && !(flags & MOD_SHIFT)) {
			flags |= MOD_SHIFT; /* Not sure if this is safe for all layouts. */
		}
	#endif

	#if defined(IS_WINDOWS)
		int modifiers = keyCode >> 8; // Pull out modifers.

		if ((modifiers & 1) != 0) { flags |= MOD_SHIFT; } // Uptdate flags from keycode modifiers.
		if ((modifiers & 2) != 0) { flags |= MOD_CONTROL; }
		if ((modifiers & 4) != 0) { flags |= MOD_ALT; }
		keyCode = keyCode & 0xff; // Mask out modifiers.
	#endif

	return toggleKeyCode(keyCode, down, flags, pid);
}

// void tapKey(char c, MMKeyFlags flags){
// 	toggleKey(c, true, flags);
// 	microsleep(5.0);
// 	toggleKey(c, false, flags);
// }

#if defined(IS_MACOSX)
	/* Post one key event carrying a UTF-16 string (1 code unit, or a
	   surrogate pair for runes above the BMP). */
	int toggleUnicode(const UniChar *ch, UniCharCount len, const bool down, uintptr pid) {
		/* This function relies on the convenient CGEventKeyboardSetUnicodeString()*/
		CGEventSourceRef source = MMEventSourceCreate();
		CGEventRef keyEvent = CGEventCreateKeyboardEvent(source, 0, down);
		if (keyEvent == NULL) {
			// fputs("Could not create keyboard event.\n", stderr);
			if (source != NULL) { CFRelease(source); }
			return (int)kCGErrorCannotComplete;
		}

		CGEventKeyboardSetUnicodeString(keyEvent, len, ch);
		SendTo(pid, keyEvent);
		if (source != NULL) { CFRelease(source); }
		return 0;
	}
#else
	#define toggleUniKey(c, down) toggleKey(c, down, MOD_NONE, 0)
#endif

#if defined(IS_MACOSX) || defined(IS_WINDOWS)
	/* Encode a code point as UTF-16 into out[2]; returns the unit count. */
	static int utf16Units(unsigned value, uint16_t out[2]) {
		if (value < 0x10000) {
			out[0] = (uint16_t)value;
			return 1;
		}
		value -= 0x10000;
		out[0] = (uint16_t)(0xD800 | (value >> 10));
		out[1] = (uint16_t)(0xDC00 | (value & 0x3FF));
		return 2;
	}
#endif

// unicode type
int unicodeType(const unsigned value, uintptr pid, int8_t isPid) {
	#if defined(IS_MACOSX)
		uint16_t units[2];
		UniCharCount len = (UniCharCount)utf16Units(value, units);

		int err = toggleUnicode((const UniChar *)units, len, true, pid);
		microsleep(5.0);
		int err1 = toggleUnicode((const UniChar *)units, len, false, pid);
		return err != 0 ? err : err1;
	#elif defined(IS_WINDOWS)
		uint16_t units[2];
		int len = utf16Units(value, units);

		if (pid != 0) {
			HWND hwnd = getHwnd(pid, isPid);
			if (hwnd == NULL) {
				/* PostMessageW(NULL, ...) would post to our own queue and
				   "succeed". */
				return ERROR_INVALID_WINDOW_HANDLE;
			}

			// SendMessage(hwnd, down, value, 0);
			int i;
			for (i = 0; i < len; i++) {
				if (!PostMessageW(hwnd, WM_CHAR, units[i], 0)) {
					DWORD err = GetLastError();
					return err != 0 ? (int)err : MM_ERR_INPUT_BLOCKED;
				}
			}
			return 0;
		}

		/* down/up per UTF-16 unit: 2 inputs for the BMP, 4 for a pair. */
		INPUT input[4];
		memset(input, 0, sizeof(input));
		UINT n = 0;
		int i;
		for (i = 0; i < len; i++) {
			input[n].type = INPUT_KEYBOARD;
			input[n].ki.wScan = units[i];
			input[n].ki.dwFlags = KEYEVENTF_UNICODE;
			n++;
		}
		for (i = 0; i < len; i++) {
			input[n].type = INPUT_KEYBOARD;
			input[n].ki.wScan = units[i];
			input[n].ki.dwFlags = KEYEVENTF_KEYUP | KEYEVENTF_UNICODE;
			n++;
		}

		UINT sent = SendInput(n, input, sizeof(INPUT));
		if (sent == n) {
			return 0;
		}
		DWORD err = GetLastError();
		return err != 0 ? (int)err : MM_ERR_INPUT_BLOCKED;
	#elif defined(USE_X11)
		int err = toggleUniKey(value, true);
		microsleep(5.0);
		int err1 = toggleUniKey(value, false);
		return err != 0 ? err : err1;
	#endif
}

#if defined(USE_X11)
	/* Type a keysym name (e.g. "U1F600") by remapping a spare keycode.
	   Returns 0 on success, 1 when there is no display, 2 for an unknown
	   keysym, 3 when the keyboard mapping is unavailable and 4 when XTEST
	   rejects the key event. */
	int input_utf(const char *utf) {
		Display *dpy = XOpenDisplay(NULL);
		if (dpy == NULL) { return 1; }
		KeySym sym = XStringToKeysym(utf);
		if (sym == NoSymbol) {
			XCloseDisplay(dpy);
			return 2;
		}
		// KeySym sym = XKeycodeToKeysym(dpy, utf);

		int min, max, numcodes;
		XDisplayKeycodes(dpy, &min, &max);
		KeySym *keysym;
		keysym = XGetKeyboardMapping(dpy, min, max-min+1, &numcodes);
		if (keysym == NULL) {
			XCloseDisplay(dpy);
			return 3;
		}
		keysym[(max-min-1)*numcodes] = sym;
		XChangeKeyboardMapping(dpy, min, numcodes, keysym, (max-min));
		XFree(keysym);
		XFlush(dpy);

		KeyCode code = XKeysymToKeycode(dpy, sym);
		int ok = code != 0
			&& XTestFakeKeyEvent(dpy, code, True, 1)
			&& XTestFakeKeyEvent(dpy, code, False, 1);

		XFlush(dpy);
		XCloseDisplay(dpy);
		return ok ? 0 : 4;
	}
#else
	int input_utf(const char *utf){
		return 0;
	}
#endif