#include "keypress.h"
// #include "../base/deadbeef_rand_c.h"
#include "../base/microsleep.h"

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
	#define WIN32_KEY_EVENT_WAIT(key, flags) \
		(win32KeyEvent(key, flags), Sleep(DEADBEEF_RANDRANGE(0, 1)))
#elif defined(USE_X11)
	#define X_KEY_EVENT(display, key, is_press) \
		(XTestFakeKeyEvent(display, \
		                   XKeysymToKeycode(display, key), \
		                   is_press, CurrentTime), \
		 XSync(display, false))
	#define X_KEY_EVENT_WAIT(display, key, is_press) \
		(X_KEY_EVENT(display, key, is_press), \
		 microsleep(DEADBEEF_UNIFORM(0.0, 0.5)))
#endif

#if defined(IS_MACOSX)
static io_connect_t _getAuxiliaryKeyDriver(void){
	static mach_port_t sEventDrvrRef = 0;
	mach_port_t masterPort, service, iter;
	kern_return_t kr;

	if (!sEventDrvrRef) {
		kr = IOMasterPort( bootstrap_port, &masterPort );
		assert(KERN_SUCCESS == kr);
		kr = IOServiceGetMatchingServices(masterPort, 
			IOServiceMatching( kIOHIDSystemClass), &iter );
		assert(KERN_SUCCESS == kr);

		service = IOIteratorNext( iter );
		assert(service);
		
		kr = IOServiceOpen(service, 
			mach_task_self(), kIOHIDParamConnectType, &sEventDrvrRef );
		assert(KERN_SUCCESS == kr);

		IOObjectRelease(service);
		IOObjectRelease(iter);
	}
	return sEventDrvrRef;
}
#endif

#if defined(IS_WINDOWS)
void win32KeyEvent(int key, MMKeyFlags flags){
	int scan = MapVirtualKey(key & 0xff, MAPVK_VK_TO_VSC);

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
			break;
		}
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
	SendInput(1, &keyInput, sizeof(keyInput));
}
#endif

void toggleKeyCode(MMKeyCode code, const bool down, MMKeyFlags flags){
#if defined(IS_MACOSX)
	/* The media keys all have 1000 added to them to help us detect them. */
	if (code >= 1000) {
		code = code - 1000; /* Get the real keycode. */
		NXEventData   event;
		kern_return_t kr;

		IOGPoint loc = { 0, 0 };
		UInt32 evtInfo = code << 16 | (down?NX_KEYDOWN:NX_KEYUP) << 8;

		bzero(&event, sizeof(NXEventData));
		event.compound.subType = NX_SUBTYPE_AUX_CONTROL_BUTTONS;
		event.compound.misc.L[0] = evtInfo;

		kr = IOHIDPostEvent( _getAuxiliaryKeyDriver(), 
			NX_SYSDEFINED, loc, &event, kNXEventDataVersion, 0, FALSE );
		assert( KERN_SUCCESS == kr );
	} else {
		CGEventRef keyEvent = CGEventCreateKeyboardEvent(NULL,
		                            			(CGKeyCode)code, down);
		assert(keyEvent != NULL);

		CGEventSetType(keyEvent, down ? kCGEventKeyDown : kCGEventKeyUp);
		// CGEventSetFlags(keyEvent, flags);
		CGEventSetFlags(keyEvent, (int) flags);
		CGEventPost(kCGSessionEventTap, keyEvent);
		CFRelease(keyEvent);
	}
#elif defined(IS_WINDOWS)
	const DWORD dwFlags = down ? 0 : KEYEVENTF_KEYUP;

	/* Parse modifier keys. */
	if (flags & MOD_META) WIN32_KEY_EVENT_WAIT(K_META, dwFlags);
	if (flags & MOD_ALT) WIN32_KEY_EVENT_WAIT(K_ALT, dwFlags);
	if (flags & MOD_CONTROL) WIN32_KEY_EVENT_WAIT(K_CONTROL, dwFlags);
	if (flags & MOD_SHIFT) WIN32_KEY_EVENT_WAIT(K_SHIFT, dwFlags);

	win32KeyEvent(code, dwFlags);
#elif defined(USE_X11)
	Display *display = XGetMainDisplay();
	const Bool is_press = down ? True : False; /* Just to be safe. */

	/* Parse modifier keys. */
	if (flags & MOD_META) X_KEY_EVENT_WAIT(display, K_META, is_press);
	if (flags & MOD_ALT) X_KEY_EVENT_WAIT(display, K_ALT, is_press);
	if (flags & MOD_CONTROL) X_KEY_EVENT_WAIT(display, K_CONTROL, is_press);
	if (flags & MOD_SHIFT) X_KEY_EVENT_WAIT(display, K_SHIFT, is_press);

	X_KEY_EVENT(display, code, is_press);
#endif
}

void tapKeyCode(MMKeyCode code, MMKeyFlags flags){
	toggleKeyCode(code, true, flags);
	toggleKeyCode(code, false, flags);
}

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

void toggleKey(char c, const bool down, MMKeyFlags flags){
	MMKeyCode keyCode = keyCodeForChar(c);

	//Prevent unused variable warning for Mac and Linux.
	#if defined(IS_WINDOWS)
		int modifiers;
	#endif

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
		modifiers = keyCode >> 8; // Pull out modifers.
		if ((modifiers & 1) != 0) flags |= MOD_SHIFT; // Uptdate flags from keycode modifiers.
		if ((modifiers & 2) != 0) flags |= MOD_CONTROL;
		if ((modifiers & 4) != 0) flags |= MOD_ALT;
		keyCode = keyCode & 0xff; // Mask out modifiers.
	#endif

	toggleKeyCode(keyCode, down, flags);
}

void tapKey(char c, MMKeyFlags flags){
	toggleKey(c, true, flags);
	toggleKey(c, false, flags);
}

#if defined(IS_MACOSX)
void toggleUnicode(UniChar ch, const bool down){
	/* This function relies on the convenient
	 * CGEventKeyboardSetUnicodeString(), which allows us to not have to
	 * convert characters to a keycode, but does not support adding modifier
	 * flags. It is therefore only used in typeStringDelayed()
	 * -- if you need modifier keys, use the above functions instead. */
	CGEventRef keyEvent = CGEventCreateKeyboardEvent(NULL, 0, down);
	if (keyEvent == NULL) {
		fputs("Could not create keyboard event.\n", stderr);
		return;
	}

	CGEventKeyboardSetUnicodeString(keyEvent, 1, &ch);

	CGEventPost(kCGSessionEventTap, keyEvent);
	CFRelease(keyEvent);
}
#endif

#if defined(USE_X11)
	#define toggleUniKey(c, down) toggleKey(c, down, MOD_NONE)

	int input_utf(const char *utf) {
		Display *dpy;
		dpy = XOpenDisplay(NULL);

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
#endif
#if !defined(USE_X11)
	int input_utf(const char *utf){
		return 0;
	}
#endif

// unicode type
void unicodeType(const unsigned value){
	#if defined(IS_MACOSX)
		UniChar ch = (UniChar)value; // Convert to unsigned char

		toggleUnicode(ch, true);
		toggleUnicode(ch, false);
	#elif defined(IS_WINDOWS)
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

  		SendInput(2, input, sizeof(INPUT));
	#elif defined(USE_X11)
		toggleUniKey(value, true);
		toggleUniKey(value, false);	
	#endif
}

void typeStringDelayed(const char *str, const unsigned cpm){
	
	/* Characters per second */
	const double cps = (double)cpm / 60.0;

	/* Average milli-seconds per character */
	const double mspc = (cps == 0.0) ? 0.0 : 1000.0 / cps;

	unsigned long n;
	unsigned short c;
	unsigned short c1;
	unsigned short c2;
	unsigned short c3;

	while (*str != '\0') {
		c = *str++;

		// warning, the following utf8 decoder
		// doesn't perform validation
		if (c <= 0x7F) {
			// 0xxxxxxx one byte
			n = c;
		} else if ((c & 0xE0) == 0xC0)  {
			// 110xxxxx two bytes
			c1 = (*str++) & 0x3F;
			n = ((c & 0x1F) << 6) | c1;
		} else if ((c & 0xF0) == 0xE0) {
			// 1110xxxx three bytes
			c1 = (*str++) & 0x3F;
			c2 = (*str++) & 0x3F;
			n = ((c & 0x0F) << 12) | (c1 << 6) | c2;
		} else if ((c & 0xF8) == 0xF0) {
			// 11110xxx four bytes
			c1 = (*str++) & 0x3F;
			c2 = (*str++) & 0x3F;
			c3 = (*str++) & 0x3F;
			n = ((c & 0x07) << 18) | (c1 << 12) | (c2 << 6) | c3;
		}

		unicodeType(n);

		if (mspc > 0) {
			microsleep(mspc + (DEADBEEF_UNIFORM(0.0, 0.5)));
		}
	}
}