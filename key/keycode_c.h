#include "keycode.h"

#if defined(IS_MACOSX)

#include <CoreFoundation/CoreFoundation.h>
#include <Carbon/Carbon.h> /* For kVK_ constants, and TIS functions. */

/* Returns string representation of key, if it is printable.
 * Ownership follows the Create Rule; that is, it is the caller's
 * responsibility to release the returned object. */
CFStringRef createStringForKey(CGKeyCode keyCode);

MMKeyCode keyCodeForCharFallBack(const char c);

#elif defined(USE_X11)

/*
 * Structs to store key mappings not handled by XStringToKeysym() on some
 * Linux systems.
 */

struct XSpecialCharacterMapping {
	char name;
	MMKeyCode code;
};

struct XSpecialCharacterMapping XSpecialCharacterTable[] = {
	{'~', XK_asciitilde},
  	{'_', XK_underscore},
  	{'[', XK_bracketleft},
  	{']', XK_bracketright},
  	{'!', XK_exclam},
  	{'#', XK_numbersign},
  	{'$', XK_dollar},
  	{'%', XK_percent},
  	{'&', XK_ampersand},
  	{'*', XK_asterisk},
  	{'+', XK_plus},
  	{',', XK_comma},
  	{'-', XK_minus},
  	{'.', XK_period},
  	{'?', XK_question},
  	{'<', XK_less},
  	{'>', XK_greater},
  	{'=', XK_equal},
  	{'@', XK_at},
  	{':', XK_colon},
  	{';', XK_semicolon},
  	{'{', XK_braceleft},
  	{'}', XK_braceright},
  	{'|', XK_bar},
  	{'^', XK_asciicircum},
  	{'(', XK_parenleft},
  	{')', XK_parenright},
  	{' ', XK_space},
  	{'/', XK_slash},
	{'\\', XK_backslash},
	{'`', XK_grave},
	{'"', XK_quoteright},
  	{'\'', XK_quotedbl},
	// {'\'', XK_quoteright},
  	{'\t', XK_Tab},
  	{'\n', XK_Return}
};

#endif

MMKeyCode keyCodeForChar(const char c){
	#if defined(IS_MACOSX)
		/* OS X does not appear to have a built-in function for this, so instead we
		* have to write our own. */
		static CFMutableDictionaryRef charToCodeDict = NULL;
		CGKeyCode code;
		UniChar character = c;
		CFStringRef charStr = NULL;

		/* Generate table of keycodes and characters. */
		if (charToCodeDict == NULL) {
			size_t i;
			charToCodeDict = CFDictionaryCreateMutable(kCFAllocatorDefault,
													128,
													&kCFCopyStringDictionaryKeyCallBacks,
													NULL);
			if (charToCodeDict == NULL) return K_NOT_A_KEY;

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
		if (!CFDictionaryGetValueIfPresent(charToCodeDict, charStr,
										(const void **)&code)) {
			code = UINT16_MAX; /* Error */
		}

		CFRelease(charStr);

		// TISGetInputSourceProperty may return nil so we need fallback
		if (code == UINT16_MAX) {
			code = keyCodeForCharFallBack(c);
		}

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
		MMKeyCode code;

		char buf[2];
		buf[0] = c;
		buf[1] = '\0';

		code = XStringToKeysym(buf);
		if (code == NoSymbol) {
			/* Some special keys are apparently not handled properly by
			* XStringToKeysym() on some systems, so search for them instead in our
			* mapping table. */
			struct XSpecialCharacterMapping* xs = XSpecialCharacterTable;
			while (xs->name) {
				if (c == xs->name ) {
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

		return code;
	#endif
}


#if defined(IS_MACOSX)

CFStringRef createStringForKey(CGKeyCode keyCode){
	TISInputSourceRef currentKeyboard = TISCopyCurrentASCIICapableKeyboardInputSource();
	CFDataRef layoutData =
		(CFDataRef)TISGetInputSourceProperty(currentKeyboard,
		                          kTISPropertyUnicodeKeyLayoutData);

	if (layoutData == nil) { return 0; }

	const UCKeyboardLayout *keyboardLayout =
		(const UCKeyboardLayout *)CFDataGetBytePtr(layoutData);

	UInt32 keysDown = 0;
	UniChar chars[4];
	UniCharCount realLength;

	UCKeyTranslate(keyboardLayout,
	               keyCode,
	               kUCKeyActionDisplay,
	               0,
	               LMGetKbdType(),
	               kUCKeyTranslateNoDeadKeysBit,
	               &keysDown,
	               sizeof(chars) / sizeof(chars[0]),
	               &realLength,
	               chars);
	CFRelease(currentKeyboard);

	return CFStringCreateWithCharacters(kCFAllocatorDefault, chars, 1);
}

MMKeyCode keyCodeForCharFallBack(const char c) {
	switch (c) {
		case 'A': return kVK_ANSI_A;
		case 'B': return kVK_ANSI_B;
		case 'C': return kVK_ANSI_C;
		case 'D': return kVK_ANSI_D;
		case 'E': return kVK_ANSI_E;
		case 'F': return kVK_ANSI_F;
		case 'G': return kVK_ANSI_G;
		case 'H': return kVK_ANSI_H;
		case 'I': return kVK_ANSI_I;
		case 'J': return kVK_ANSI_J;
		case 'K': return kVK_ANSI_K;
		case 'L': return kVK_ANSI_L;
		case 'M': return kVK_ANSI_M;
		case 'N': return kVK_ANSI_N;
		case 'O': return kVK_ANSI_O;
		case 'P': return kVK_ANSI_P;
		case 'Q': return kVK_ANSI_Q;
		case 'R': return kVK_ANSI_R;
		case 'S': return kVK_ANSI_S;
		case 'T': return kVK_ANSI_T;
		case 'U': return kVK_ANSI_U;
		case 'V': return kVK_ANSI_V;
		case 'W': return kVK_ANSI_W;
		case 'X': return kVK_ANSI_X;
		case 'Y': return kVK_ANSI_Y;
		case 'Z': return kVK_ANSI_Z;


		case '0': return kVK_ANSI_0;
		case '1': return kVK_ANSI_1;
		case '2': return kVK_ANSI_2;
		case '3': return kVK_ANSI_3;
		case '4': return kVK_ANSI_4;
		case '5': return kVK_ANSI_5;
		case '6': return kVK_ANSI_6;
		case '7': return kVK_ANSI_7;
		case '8': return kVK_ANSI_8;
		case '9': return kVK_ANSI_9;
	}

	return K_NOT_A_KEY;
}

#endif
