#include "keycode.h"

#if defined(IS_MACOSX)
	#include <CoreFoundation/CoreFoundation.h>
	#include <Carbon/Carbon.h> /* For kVK_ constants, and TIS functions. */

	/* Returns string representation of key, if it is printable. 
	Ownership follows the Create Rule; 
	it is the caller's responsibility to release the returned object. */
	CFStringRef createStringForKey(CGKeyCode keyCode);
#endif

MMKeyCode keyCodeForChar(const char c) {
	#if defined(IS_MACOSX)
		/* OS X does not appear to have a built-in function for this, 
		so instead it to write our own. */
		static CFMutableDictionaryRef charToCodeDict = NULL;
		CGKeyCode code;
		UniChar character = c;
		CFStringRef charStr = NULL;

		/* Generate table of keycodes and characters. */
		if (charToCodeDict == NULL) {
			size_t i;
			charToCodeDict = CFDictionaryCreateMutable(kCFAllocatorDefault, 128,
				&kCFCopyStringDictionaryKeyCallBacks, NULL);
			if (charToCodeDict == NULL) { return K_NOT_A_KEY; }

			/* Loop through every keycode (0 - 127) to find its current mapping. */
			for (i = 0; i < 128; ++i) {
				CFStringRef string = createStringForKey((CGKeyCode)i);
				if (string != NULL) {
					CFDictionaryAddValue(charToCodeDict, string, (const void *)i);
					CFRelease(string);
				}
			}
		}

		charStr = CFStringCreateWithCharacters(kCFAllocatorDefault, &character, 1);
		/* Our values may be NULL (0), so we need to use this function. */
		if (!CFDictionaryGetValueIfPresent(charToCodeDict, charStr, (const void **)&code)) {
			code = UINT16_MAX; /* Error */
		}
		CFRelease(charStr);

		// TISGetInputSourceProperty may return nil so we need fallback
		if (code == UINT16_MAX) {
			return K_NOT_A_KEY;
		}

		return (MMKeyCode)code;
	#elif defined(IS_WINDOWS)
		MMKeyCode code;
		code = VkKeyScan(c);
		if (code == 0xFFFF) {
			return K_NOT_A_KEY;
		}

		return code;
	#elif defined(USE_X11)
		char buf[2];
		buf[0] = c;
		buf[1] = '\0';

		MMKeyCode code = XStringToKeysym(buf);
		if (code == NoSymbol) {
			/* Some special keys are apparently not handled properly */
			struct XSpecialCharacterMapping* xs = XSpecialCharacterTable;
			while (xs->name) {
				if (c == xs->name) {
					code = xs->code;
					// 
					break;
				}
				xs++;
			}
		}

		if (code == NoSymbol) {
			return K_NOT_A_KEY;
		}

		// x11 key bug
		if (c == 60) {
			code = 44;
		}
		return code;
	#endif
}

#if defined(IS_MACOSX)
	CFStringRef createStringForKey(CGKeyCode keyCode){
		TISInputSourceRef currentKeyboard = TISCopyCurrentASCIICapableKeyboardInputSource();
		CFDataRef layoutData = (CFDataRef) TISGetInputSourceProperty(
			currentKeyboard, kTISPropertyUnicodeKeyLayoutData);

		if (layoutData == nil) { return 0; }

		const UCKeyboardLayout *keyboardLayout = (const UCKeyboardLayout *) CFDataGetBytePtr(layoutData);
		UInt32 keysDown = 0;
		UniChar chars[4];
		UniCharCount realLength;

		UCKeyTranslate(keyboardLayout, keyCode, kUCKeyActionDisplay, 0, LMGetKbdType(),
					kUCKeyTranslateNoDeadKeysBit, &keysDown,
					sizeof(chars) / sizeof(chars[0]), &realLength, chars);
		CFRelease(currentKeyboard);

		return CFStringCreateWithCharacters(kCFAllocatorDefault, chars, 1);
	}
#endif
