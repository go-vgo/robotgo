// Copyright (c) 2016-2025 AtomAI, All rights reserved.
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

	void X_KEY_EVENT(Display *display, MMKeyCode key, bool is_press) {
		XTestFakeKeyEvent(display, XKeysymToKeycode(display, key), is_press, CurrentTime); 
		XSync(display, false);
	}

	void X_KEY_EVENT_WAIT(Display *display, MMKeyCode key, bool is_press) {
		X_KEY_EVENT(display, key, is_press);
		microsleep(DEADBEEF_UNIFORM(0.0, 0.5));
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
			#if __ENVIRONMENT_MAC_OS_X_VERSION_MIN_REQUIRED__ <= 12000
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
		return sent == 1 ? 0 : (int)GetLastError();
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
		CGEventSourceRef source = CGEventSourceCreate(kCGEventSourceStateHIDSystemState);
		CGEventRef keyEvent = CGEventCreateKeyboardEvent(source, (CGKeyCode)code, down);
		// assert(keyEvent != NULL);
		if (keyEvent == NULL) {
			CFRelease(source);
			return (int)kCGErrorCannotComplete;
		}

		CGEventSetType(keyEvent, down ? kCGEventKeyDown : kCGEventKeyUp);
		if (flags != 0) {
			CGEventSetFlags(keyEvent, (CGEventFlags) flags);
		}
		
		SendTo(pid, keyEvent);
		CFRelease(source);
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

	/* Parse modifier keys. */
	if (flags & MOD_META) { X_KEY_EVENT_WAIT(display, K_META, is_press); }
	if (flags & MOD_ALT) { X_KEY_EVENT_WAIT(display, K_ALT, is_press); }
	if (flags & MOD_CONTROL) { X_KEY_EVENT_WAIT(display, K_CONTROL, is_press); }
	if (flags & MOD_SHIFT) { X_KEY_EVENT_WAIT(display, K_SHIFT, is_press); }

	return X_KEY_EVENT(display, code, is_press);
#endif
}

// void tapKeyCode(MMKeyCode code, MMKeyFlags flags){
// 	toggleKeyCode(code, true, flags);
// 	microsleep(5.0);
// 	toggleKeyCode(code, false, flags);
// }

#if defined(USE_X11)
	bool toUpper(char c) {
		if (isupper(c)) {
			return true;
		}

		char *special = "~!@#$%^&*()_+{}|:\"<>?";
		while (*special) {
			if (*special == c) {
				return true;
			}
			special++;
		}
		return false;
	}
#endif

int toggleKey(char c, const bool down, MMKeyFlags flags, uintptr pid) {
	MMKeyCode keyCode = keyCodeForChar(c);

	#if defined(USE_X11)
		if (toUpper(c) && !(flags & MOD_SHIFT)) {
			flags |= MOD_SHIFT; /* Not sure if this is safe for all layouts. */
		}
	#else
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
	int toggleUnicode(UniChar ch, const bool down, uintptr pid) {
		/* This function relies on the convenient CGEventKeyboardSetUnicodeString()*/
		CGEventSourceRef source = CGEventSourceCreate(kCGEventSourceStateHIDSystemState);
		CGEventRef keyEvent = CGEventCreateKeyboardEvent(source, 0, down);
		if (keyEvent == NULL) {
			// fputs("Could not create keyboard event.\n", stderr);
			CFRelease(source);
			return (int)kCGErrorCannotComplete;
		}

		CGEventKeyboardSetUnicodeString(keyEvent, 1, &ch);
		SendTo(pid, keyEvent);
		CFRelease(source);
		return 0;
	}
#else
	#define toggleUniKey(c, down) toggleKey(c, down, MOD_NONE, 0)
#endif

// unicode type
int unicodeType(const unsigned value, uintptr pid, int8_t isPid) {
	#if defined(IS_MACOSX)
		UniChar ch = (UniChar)value; // Convert to unsigned char

		toggleUnicode(ch, true, pid);
		microsleep(5.0);
		return toggleUnicode(ch, false, pid);
	#elif defined(IS_WINDOWS)
		if (pid != 0) {
			HWND hwnd = getHwnd(pid, isPid);

			// SendMessage(hwnd, down, value, 0);
			UINT sent = PostMessageW(hwnd, WM_CHAR, value, 0);
			return sent == 1 ? 0 : (int)GetLastError();
		}

		INPUT input[2];
        memset(input, 0, sizeof(input));

        input[0].type = INPUT_KEYBOARD;
  		input[0].ki.wVk = 0;
  		input[0].ki.wScan = value;
  		input[0].ki.dwFlags = 0x4; // KEYEVENTF_UNICODE;

  		input[1].type = INPUT_KEYBOARD;
  		input[1].ki.wVk = 0;
  		input[1].ki.wScan = value;
  		input[1].ki.dwFlags = KEYEVENTF_KEYUP | 0x4; // KEYEVENTF_UNICODE;

  		UINT sent = SendInput(2, input, sizeof(INPUT));
		return sent == 1 ? 0 : (int)GetLastError();
	#elif defined(USE_X11)
		toggleUniKey(value, true);
		microsleep(5.0);
		return toggleUniKey(value, false);	
	#endif
}

#if defined(USE_X11)
	int input_utf(const char *utf) {
		Display *dpy = XOpenDisplay(NULL);
		KeySym sym = XStringToKeysym(utf);
		// KeySym sym = XKeycodeToKeysym(dpy, utf);

		int min, max, numcodes;
		XDisplayKeycodes(dpy, &min, &max);
		KeySym *keysym;
		keysym = XGetKeyboardMapping(dpy, min, max-min+1, &numcodes);
		keysym[(max-min-1)*numcodes] = sym;
		XChangeKeyboardMapping(dpy, min, numcodes, keysym, (max-min));
		XFree(keysym);
		XFlush(dpy);

		KeyCode code = XKeysymToKeycode(dpy, sym);
		XTestFakeKeyEvent(dpy, code, True, 1);
		XTestFakeKeyEvent(dpy, code, False, 1);

		XFlush(dpy);
		XCloseDisplay(dpy);
		return 0;
	}
#else
	int input_utf(const char *utf){
		return 0;
	}
#endif