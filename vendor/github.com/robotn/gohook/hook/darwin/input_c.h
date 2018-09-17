
#ifdef HAVE_CONFIG_H
#include <config.h>
#endif

#ifdef USE_COREFOUNDATION
	#include <CoreFoundation/CoreFoundation.h>
#endif
#ifndef USE_WEAK_IMPORT
	#include <dlfcn.h>
#endif

#include <stdbool.h>
#include "input.h"
#include "../iohook.h"
#include "../logger_c.h"

// Current dead key state.
#if defined(USE_CARBON_LEGACY) || defined(USE_COREFOUNDATION)
static UInt32 deadkey_state;
#endif

// Input source data for the keyboard.
#if defined(USE_CARBON_LEGACY)
static KeyboardLayoutRef prev_keyboard_layout = NULL;
#elif defined(USE_COREFOUNDATION)
static TISInputSourceRef prev_keyboard_layout = NULL;
#endif

#ifdef USE_WEAK_IMPORT
// Required to dynamically check for AXIsProcessTrustedWithOptions availability.
extern Boolean AXIsProcessTrustedWithOptions(CFDictionaryRef options) __attribute__((weak_import));
extern CFStringRef kAXTrustedCheckOptionPrompt __attribute__((weak_import));
#else
static Boolean (*AXIsProcessTrustedWithOptions_t)(CFDictionaryRef);
#endif

bool is_accessibility_enabled() {
	bool is_enabled = false;

	#ifdef USE_WEAK_IMPORT
	// Check and make sure assistive devices is enabled.
	if (AXIsProcessTrustedWithOptions != NULL) {
		// New accessibility API 10.9 and later.
		const void * keys[] = { kAXTrustedCheckOptionPrompt };
		const void * values[] = { kCFBooleanTrue };

		CFDictionaryRef options = CFDictionaryCreate(
				kCFAllocatorDefault,
				keys,
				values,
				sizeof(keys) / sizeof(*keys),
				&kCFCopyStringDictionaryKeyCallBacks,
				&kCFTypeDictionaryValueCallBacks);

		is_enabled = AXIsProcessTrustedWithOptions(options);
	}
	#else
	// Dynamically load the application services framework for examination.
	*(void **) (&AXIsProcessTrustedWithOptions_t) = dlsym(RTLD_DEFAULT, "AXIsProcessTrustedWithOptions");
	const char *dlError = dlerror();
	if (AXIsProcessTrustedWithOptions_t != NULL && dlError == NULL) {
		// Check for property CFStringRef kAXTrustedCheckOptionPrompt
		void ** kAXTrustedCheckOptionPrompt_t = dlsym(RTLD_DEFAULT, "kAXTrustedCheckOptionPrompt");

		dlError = dlerror();
		if (kAXTrustedCheckOptionPrompt_t != NULL && dlError == NULL) {
			// New accessibility API 10.9 and later.
			const void * keys[] = { *kAXTrustedCheckOptionPrompt_t };
			const void * values[] = { kCFBooleanTrue };

			CFDictionaryRef options = CFDictionaryCreate(
					kCFAllocatorDefault,
					keys,
					values,
					sizeof(keys) / sizeof(*keys),
					&kCFCopyStringDictionaryKeyCallBacks,
					&kCFTypeDictionaryValueCallBacks);

			is_enabled = (*AXIsProcessTrustedWithOptions_t)(options);
		}
	}
	#endif
	else {
		#ifndef USE_WEAK_IMPORT
		if (dlError != NULL) {
			// Could not load the AXIsProcessTrustedWithOptions function!
			logger(LOG_LEVEL_DEBUG, "%s [%u]: %s.\n",
					__FUNCTION__, __LINE__, dlError);
		}
		#endif

		logger(LOG_LEVEL_DEBUG, "%s [%u]: Weak import AXIsProcessTrustedWithOptions not found.\n",
					__FUNCTION__, __LINE__, dlError);

		logger(LOG_LEVEL_DEBUG, "%s [%u]: Falling back to AXAPIEnabled().\n",
				__FUNCTION__, __LINE__, dlError);

		// Old accessibility check 10.8 and older.
		is_enabled = AXAPIEnabled();
	}

	return is_enabled;
}


UniCharCount keycode_to_unicode(CGEventRef event_ref, UniChar *buffer, UniCharCount size) {
	UniCharCount count = 0;

	#if defined(USE_CARBON_LEGACY)
	KeyboardLayoutRef curr_keyboard_layout;
	void *inputData = NULL;
	if (KLGetCurrentKeyboardLayout(&curr_keyboard_layout) == noErr) {
		if (KLGetKeyboardLayoutProperty(curr_keyboard_layout, kKLuchrData, (const void **) &inputData) != noErr) {
			inputData = NULL;
		}
	}
	#elif defined(USE_COREFOUNDATION)
	CFDataRef inputData = NULL;
	if (CFEqual(CFRunLoopGetCurrent(), CFRunLoopGetMain())) {
		// NOTE The following block must execute on the main runloop,
		// Ex: CFEqual(CFRunLoopGetCurrent(), CFRunLoopGetMain()) to avoid
		// Exception detected while handling key input and TSMProcessRawKeyCode failed
		// (-192) errors.
		TISInputSourceRef curr_keyboard_layout = TISCopyCurrentKeyboardLayoutInputSource();
		if (curr_keyboard_layout != NULL && CFGetTypeID(curr_keyboard_layout) == TISInputSourceGetTypeID()) {
			CFDataRef data = (CFDataRef) TISGetInputSourceProperty(curr_keyboard_layout, kTISPropertyUnicodeKeyLayoutData);
			if (data != NULL && CFGetTypeID(data) == CFDataGetTypeID() && CFDataGetLength(data) > 0) {
				inputData = (CFDataRef) data;
			}
		}

		// Check if the keyboard layout has changed to see if the dead key state needs to be discarded.
		if (prev_keyboard_layout != NULL && curr_keyboard_layout != NULL && CFEqual(curr_keyboard_layout, prev_keyboard_layout) == false) {
			deadkey_state = 0x00;
		}

		// Release the previous keyboard layout.
		if (prev_keyboard_layout != NULL) {
			CFRelease(prev_keyboard_layout);
			prev_keyboard_layout = NULL;
		}

		// Set the previous keyboard layout to the current layout.
		if (curr_keyboard_layout != NULL) {
			prev_keyboard_layout = curr_keyboard_layout;
		}
	}
	#endif

	#if defined(USE_CARBON_LEGACY) || defined(USE_COREFOUNDATION)
	if (inputData != NULL) {
		#ifdef USE_CARBON_LEGACY
		const UCKeyboardLayout *keyboard_layout = (const UCKeyboardLayout *) inputData;
		#else
		const UCKeyboardLayout *keyboard_layout = (const UCKeyboardLayout*) CFDataGetBytePtr(inputData);
		#endif

		if (keyboard_layout != NULL) {
			//Extract keycode and modifier information.
			CGKeyCode keycode = CGEventGetIntegerValueField(event_ref, kCGKeyboardEventKeycode);
			CGEventFlags modifiers = CGEventGetFlags(event_ref);

			// Disable all command modifiers for translation.  This is required
			// so UCKeyTranslate will provide a keysym for the separate event.
			static const CGEventFlags cmd_modifiers = kCGEventFlagMaskCommand |
					kCGEventFlagMaskControl | kCGEventFlagMaskAlternate;
			modifiers &= ~cmd_modifiers;

			// I don't know why but UCKeyTranslate does not process the
			// kCGEventFlagMaskAlphaShift (A.K.A. Caps Lock Mask) correctly.
			// We need to basically turn off the mask and process the capital
			// letters after UCKeyTranslate().
			bool is_caps_lock = modifiers & kCGEventFlagMaskAlphaShift;
			modifiers &= ~kCGEventFlagMaskAlphaShift;

			// Run the translation with the saved deadkey_state.
			OSStatus status = UCKeyTranslate(
					keyboard_layout,
					keycode,
					kUCKeyActionDown, //kUCKeyActionDisplay,
					(modifiers >> 16) & 0xFF, //(modifiers >> 16) & 0xFF, || (modifiers >> 8) & 0xFF,
					LMGetKbdType(),
					kNilOptions, //kNilOptions, //kUCKeyTranslateNoDeadKeysMask
					&deadkey_state,
					size,
					&count,
					buffer);

			if (status == noErr && count > 0) {
				if (is_caps_lock) {
					// We *had* a caps lock mask so we need to convert to uppercase.
					CFMutableStringRef keytxt = CFStringCreateMutableWithExternalCharactersNoCopy(kCFAllocatorDefault, buffer, count, size, kCFAllocatorNull);
					if (keytxt != NULL) {
						CFLocaleRef locale = CFLocaleCopyCurrent();
						CFStringUppercase(keytxt, locale);
						CFRelease(locale);
						CFRelease(keytxt);
					}
					else {
						// There was an problem creating the CFMutableStringRef.
						count = 0;
					}
				}
			}
			else {
				// Make sure the buffer count is zero if an error occurred.
				count = 0;
			}
		}
	}
	#else
	CGEventKeyboardGetUnicodeString(event_ref, size, &count, buffer);
	#endif

	// The following codes should not be processed because they are invalid.
	if (count == 1) {
		switch (buffer[0]) {
			case 0x01:		// Home
			case 0x04:		// End
			case 0x05:		// Help Key
			case 0x10:		// Function Keys
			case 0x0B:		// Page Up
			case 0x0C:		// Page Down
			case 0x1F:		// Volume Up
				count = 0;
		}
	}

	return count;
}


static const uint16_t keycode_scancode_table[][2] = {
	/* idx		{ keycode,				scancode				}, */
	/*   0 */	{ VC_A,					kVK_Undefined			},	// 0x00
	/*   1 */	{ VC_S,					kVK_Escape				},	// 0x01
	/*   2 */	{ VC_D,					kVK_ANSI_1				},	// 0x02
	/*   3 */	{ VC_F,					kVK_ANSI_2				},	// 0x03
	/*   4 */	{ VC_H,					kVK_ANSI_3				},	// 0x04
	/*   5 */	{ VC_G,					kVK_ANSI_4				},	// 0x05
	/*   6 */	{ VC_Z,					kVK_ANSI_5				},	// 0x07
	/*   7 */	{ VC_X,					kVK_ANSI_6				},	// 0x08
	/*   8 */	{ VC_C,					kVK_ANSI_7				},	// 0x09
	/*   9 */	{ VC_V,					kVK_ANSI_8				},	// 0x0A
	/*  10 */	{ VC_UNDEFINED,			kVK_ANSI_9				},	// 0x0B
	/*  11 */	{ VC_B,					kVK_ANSI_0				},	// 0x0C
	/*  12 */	{ VC_Q,					kVK_ANSI_Minus			},	// 0x0D
	/*  13 */	{ VC_W,					kVK_ANSI_Equal			},	// 0x0E
	/*  14 */	{ VC_E,					kVK_Delete				},	// 0x0F
	/*  15 */	{ VC_R,					kVK_Tab					},	// 0x10
	/*  16 */	{ VC_Y,					kVK_ANSI_Q				},	// 0x11
	/*  17 */	{ VC_T,					kVK_ANSI_W				},	// 0x12
	/*  18 */	{ VC_1,					kVK_ANSI_E				},	// 0x13
	/*  19 */	{ VC_2,					kVK_ANSI_R				},	// 0x14
	/*  20 */	{ VC_3,					kVK_ANSI_T				},	// 0x15
	/*  21 */	{ VC_4,					kVK_ANSI_Y				},	// 0x16
	/*  22 */	{ VC_6,					kVK_ANSI_U				},	// 0x17
	/*  23 */	{ VC_5,					kVK_ANSI_I				},	// 0x18
	/*  24 */	{ VC_EQUALS,			kVK_ANSI_O				},	// 0x19
	/*  25 */	{ VC_9,					kVK_ANSI_P				},	// 0x19
	/*  26 */	{ VC_7,					kVK_ANSI_LeftBracket	},	// 0x1A
	/*  27 */	{ VC_MINUS,				kVK_ANSI_RightBracket	},	// 0x1B
	/*  28 */	{ VC_8,					kVK_Return				},	// 0x1C
	/*  29 */	{ VC_0,					kVK_Control				},	// 0x1D
	/*  30 */	{ VC_CLOSE_BRACKET,		kVK_ANSI_A				},	// 0x1E
	/*  31 */	{ VC_O,					kVK_ANSI_S				},	// 0x1F
	/*  32 */	{ VC_U,					kVK_ANSI_D				},	// 0x20
	/*  33 */	{ VC_OPEN_BRACKET,		kVK_ANSI_F				},	// 0x21
	/*  34 */	{ VC_I,					kVK_ANSI_G				},	// 0x22
	/*  35 */	{ VC_P,					kVK_ANSI_H				},	// 0x23
	/*  36 */	{ VC_ENTER,				kVK_ANSI_J				},	// 0x24
	/*  37 */	{ VC_L,					kVK_ANSI_K				},	// 0x25
	/*  38 */	{ VC_J,					kVK_ANSI_L				},	// 0x26
	/*  39 */	{ VC_QUOTE,				kVK_ANSI_Semicolon		},	// 0x27
	/*  40 */	{ VC_K,					kVK_ANSI_Quote			},	// 0x28
	/*  41 */	{ VC_SEMICOLON,			kVK_ANSI_Grave			},	// 0x29
	/*  42 */	{ VC_BACK_SLASH,		kVK_Shift				},	// 0x2A
	/*  43 */	{ VC_COMMA,				kVK_ANSI_Backslash		},	// 0x2B
	/*  44 */	{ VC_SLASH,				kVK_ANSI_Z				},	// 0x2C
	/*  45 */	{ VC_N,					kVK_ANSI_X				},	// 0x2D
	/*  46 */	{ VC_M,					kVK_ANSI_C				},	// 0x2E
	/*  47 */	{ VC_PERIOD,			kVK_ANSI_V				},	// 0x2F
	/*  48 */	{ VC_TAB,				kVK_ANSI_B				},	// 0x30
	/*  49 */	{ VC_SPACE,				kVK_ANSI_N				},	// 0x31
	/*  50 */	{ VC_BACKQUOTE,			kVK_ANSI_M				},	// 0x32
	/*  51 */	{ VC_BACKSPACE,			kVK_ANSI_Comma			},	// 0x33
	/*  52 */	{ VC_UNDEFINED,			kVK_ANSI_Period			},	// 0x34
	/*  53 */	{ VC_ESCAPE,			kVK_ANSI_Slash			},	// 0x35
	/*  54 */	{ VC_META_R,			kVK_RightShift			},	// 0x36
	/*  55 */	{ VC_META_L,			kVK_ANSI_KeypadMultiply	},	// 0x37
	/*  56 */	{ VC_SHIFT_L,			kVK_Option				},	// 0x38
	/*  57 */	{ VC_CAPS_LOCK,			kVK_Space				},	// 0x39
	/*  58 */	{ VC_ALT_L,				kVK_CapsLock			},	// 0x3A
	/*  59 */	{ VC_CONTROL_L,			kVK_F1					},	// 0x3B
	/*  60 */	{ VC_SHIFT_R,			kVK_F2					},	// 0x3C
	/*  61 */	{ VC_ALT_R,				kVK_F3					},	// 0x3D
	/*  62 */	{ VC_CONTROL_R,			kVK_F4					},	// 0x3E
	/*  63 */	{ VC_UNDEFINED,			kVK_F5					},	// 0x3F
	/*  64 */	{ VC_F17,				kVK_F6					},	// 0x40
	/*  65 */	{ VC_KP_SEPARATOR,		kVK_F7					},	// 0x41
	/*  66 */	{ VC_UNDEFINED,			kVK_F8					},	// 0x42
	/*  67 */	{ VC_KP_MULTIPLY,		kVK_F9					},	// 0x43
	/*  68 */	{ VC_UNDEFINED,			kVK_F10					},	// 0x44
	/*  69 */	{ VC_KP_ADD,			kVK_ANSI_KeypadClear	},	// 0x45
	/*  70 */	{ VC_UNDEFINED,			kVK_Undefined			},	// 0x46
	/*  71 */	{ VC_NUM_LOCK,			kVK_ANSI_Keypad7		},	// 0x47
	/*  72 */	{ VC_VOLUME_UP,			kVK_ANSI_Keypad8		},	// 0x48
	/*  73 */	{ VC_VOLUME_DOWN,		kVK_ANSI_Keypad9		},	// 0x49
	/*  74 */	{ VC_VOLUME_MUTE,		kVK_ANSI_KeypadMinus	},	// 0x4A
	/*  75 */	{ VC_KP_DIVIDE,			kVK_ANSI_Keypad4		},	// 0x4B
	/*  76 */	{ VC_KP_ENTER,			kVK_ANSI_Keypad5		},	// 0x4C
	/*  77 */	{ VC_UNDEFINED,			kVK_ANSI_Keypad6		},	// 0x4D
	/*  78 */	{ VC_KP_SUBTRACT,		kVK_ANSI_KeypadPlus		},	// 0x4E
	/*  79 */	{ VC_F18,				kVK_ANSI_Keypad1		},	// 0x4F
	/*  80 */	{ VC_F19,				kVK_ANSI_Keypad2		},	// 0x50
	/*  81 */	{ VC_KP_EQUALS,			kVK_ANSI_Keypad3		},	// 0x51
	/*  82 */	{ VC_KP_0,				kVK_ANSI_Keypad0		},	// 0x52
	/*  83 */	{ VC_KP_1,				kVK_ANSI_KeypadDecimal	},	// 0x53
	/*  84 */	{ VC_KP_2,				kVK_Undefined			},	// 0x54
	/*  85 */	{ VC_KP_3,				kVK_Undefined			},	// 0x55
	/*  86 */	{ VC_KP_4,				kVK_Undefined			},	// 0x56
	/*  87 */	{ VC_KP_5,				kVK_F11					},	// 0x57
	/*  88 */	{ VC_KP_6,				kVK_F12					},	// 0x58
	/*  89 */	{ VC_KP_7,				kVK_Undefined			},	// 0x59
	/*  90 */	{ VC_F20,				kVK_Undefined			},	// 0x5A
	/*  91 */	{ VC_KP_8,				kVK_F13					},	// 0x5B
	/*  92 */	{ VC_KP_9,				kVK_F14					},	// 0x5C
	/*  93 */	{ VC_YEN,				kVK_F15					},	// 0x5D
	/*  94 */	{ VC_UNDERSCORE,		kVK_Undefined			},	// 0x5E
	/*  95 */	{ VC_KP_COMMA,			kVK_Undefined			},	// 0x5F
	/*  96 */	{ VC_F5,				kVK_Undefined			},	// 0x60
	/*  97 */	{ VC_F6,				kVK_Undefined			},	// 0x61
	/*  98 */	{ VC_F7,				kVK_Undefined			},	// 0x62
	/*  99 */	{ VC_F3,				kVK_F16					},	// 0x63
	/* 100 */	{ VC_F8,				kVK_F17					},	// 0x64
	/* 101 */	{ VC_F9,				kVK_F18					},	// 0x65
	/* 102 */	{ VC_UNDEFINED,			kVK_F19					},	// 0x66
	/* 103 */	{ VC_F11,				kVK_F20					},	// 0x67
	/* 104 */	{ VC_KATAKANA,			kVK_Undefined			},	// 0x68
	/* 105 */	{ VC_F13,				kVK_Undefined			},	// 0x69
	/* 106 */	{ VC_F16,				kVK_Undefined			},	// 0x6A
	/* 107 */	{ VC_F14,				kVK_Undefined			},	// 0x6B
	/* 108 */	{ VC_UNDEFINED,			kVK_Undefined			},	// 0x6C		FIXME kVK_JIS_Eisu same as Caps Lock ?
	/* 109 */	{ VC_F10,				kVK_Undefined			},	// 0x6D
	/* 110 */	{ VC_UNDEFINED,			kVK_Undefined			},	// 0x6E
	/* 111 */	{ VC_F12,				kVK_Undefined			},	// 0x6F
	/* 112 */	{ VC_UNDEFINED,			kVK_JIS_Kana			},	// 0x70
	/* 113 */	{ VC_F15,				kVK_Undefined			},	// 0x71
	/* 114 */	{ VC_INSERT,			kVK_Undefined			},	// 0x72
	/* 115 */	{ VC_HOME,				kVK_JIS_Underscore		},	// 0x73
	/* 116 */	{ VC_PAGE_UP,			kVK_Undefined			},	// 0x74
	/* 117 */	{ VC_DELETE,			kVK_Undefined			},	// 0x75
	/* 118 */	{ VC_F4,				kVK_Undefined			},	// 0x76
	/* 119 */	{ VC_END,				kVK_Undefined			},	// 0x77
	/* 120 */	{ VC_F2,				kVK_Undefined			},	// 0x78
	/* 121 */	{ VC_PAGE_DOWN,			kVK_Undefined			},	// 0x79
	/* 122 */	{ VC_F1,				kVK_Undefined			},	// 0x7A
	/* 123 */	{ VC_LEFT,				kVK_Undefined			},	// 0x7B
	/* 124 */	{ VC_RIGHT,				kVK_Undefined			},	// 0x7C
	/* 125 */	{ VC_DOWN,				kVK_JIS_Yen				},	// 0x7D
	/* 126 */	{ VC_UP,				kVK_JIS_KeypadComma		},	// 0x7E
	/* 127 */	{ VC_UNDEFINED,			kVK_Undefined			},	// 0x7F

	//			No Offset				Offset (i & 0x007F) + 128

	/* 128 */	{ VC_UNDEFINED,			kVK_Undefined			},	// 0x80
	/* 129 */	{ VC_UNDEFINED,			kVK_Undefined			},	// 0x81
	/* 130 */	{ VC_UNDEFINED,			kVK_Undefined			},	// 0x82
	/* 131 */	{ VC_UNDEFINED,			kVK_Undefined			},	// 0x83
	/* 132 */	{ VC_UNDEFINED,			kVK_Undefined			},	// 0x84
	/* 133 */	{ VC_UNDEFINED,			kVK_Undefined			},	// 0x85
	/* 134 */	{ VC_UNDEFINED,			kVK_Undefined			},	// 0x86
	/* 135 */	{ VC_UNDEFINED,			kVK_Undefined			},	// 0x87
	/* 136 */	{ VC_UNDEFINED,			kVK_Undefined			},	// 0x88
	/* 137 */	{ VC_UNDEFINED,			kVK_Undefined			},	// 0x89
	/* 138 */	{ VC_UNDEFINED,			kVK_Undefined			},	// 0x8A
	/* 139 */	{ VC_UNDEFINED,			kVK_Undefined			},	// 0x8B
	/* 140 */	{ VC_UNDEFINED,			kVK_Undefined			},	// 0x8C
	/* 141 */	{ VC_UNDEFINED,			kVK_ANSI_KeypadEquals	},	// 0x8D
	/* 142 */	{ VC_UNDEFINED,			kVK_Undefined			},	// 0x8E
	/* 143 */	{ VC_UNDEFINED,			kVK_Undefined			},	// 0x8F
	/* 144 */	{ VC_UNDEFINED,			kVK_MEDIA_Previous		},	// 0x90
	/* 145 */	{ VC_UNDEFINED,			kVK_Undefined			},	// 0x91
	/* 146 */	{ VC_UNDEFINED,			kVK_Undefined			},	// 0x92
	/* 147 */	{ VC_UNDEFINED,			kVK_Undefined			},	// 0x93
	/* 148 */	{ VC_UNDEFINED,			kVK_Undefined			},	// 0x94
	/* 149 */	{ VC_UNDEFINED,			kVK_Undefined			},	// 0x95
	/* 150 */	{ VC_UNDEFINED,			kVK_Undefined			},	// 0x96
	/* 151 */	{ VC_UNDEFINED,			kVK_Undefined			},	// 0x97
	/* 152 */	{ VC_UNDEFINED,			kVK_Undefined			},	// 0x98
	/* 153 */	{ VC_UNDEFINED,			kVK_MEDIA_Next			},	// 0x99
	/* 154 */	{ VC_UNDEFINED,			kVK_Undefined			},	// 0x9A
	/* 155 */	{ VC_UNDEFINED,			kVK_Undefined			},	// 0x9B
	/* 156 */	{ VC_UNDEFINED,			kVK_ANSI_KeypadEnter	},	// 0x9C
	/* 157 */	{ VC_UNDEFINED,			kVK_RightControl		},	// 0x9D
	/* 158 */	{ VC_UNDEFINED,			kVK_Undefined			},	// 0x9E
	/* 159 */	{ VC_UNDEFINED,			kVK_Undefined			},	// 0x9F
	/* 160 */	{ VC_UNDEFINED,			kVK_Mute				},	// 0xA0
	/* 161 */	{ VC_UNDEFINED,			kVK_Undefined			},	// 0xA1
	/* 162 */	{ VC_UNDEFINED,			kVK_MEDIA_Play			},	// 0xA2
	/* 163 */	{ VC_UNDEFINED,			kVK_Undefined			},	// 0xA3
	/* 164 */	{ VC_UNDEFINED,			kVK_Undefined			},	// 0xA4
	/* 165 */	{ VC_UNDEFINED,			kVK_Undefined			},	// 0xA5
	/* 166 */	{ VC_UNDEFINED,			kVK_Undefined			},	// 0xA6
	/* 167 */	{ VC_UNDEFINED,			kVK_Undefined			},	// 0xA7
	/* 168 */	{ VC_UNDEFINED,			kVK_Undefined			},	// 0xA8
	/* 169 */	{ VC_UNDEFINED,			kVK_Undefined			},	// 0xA9
	/* 170 */	{ VC_UNDEFINED,			kVK_Undefined			},	// 0xAA
	/* 171 */	{ VC_UNDEFINED,			kVK_Undefined			},	// 0xAB
	/* 172 */	{ VC_UNDEFINED,			kVK_NX_Eject			},	// 0xAC
	/* 173 */	{ VC_UNDEFINED,			kVK_Undefined			},	// 0xAD
	/* 174 */	{ VC_UNDEFINED,			kVK_VolumeDown			},	// 0xAE
	/* 175 */	{ VC_UNDEFINED,			kVK_Undefined			},	// 0xAF
	/* 176 */	{ VC_UNDEFINED,			kVK_VolumeUp			},	// 0xB0
	/* 177 */	{ VC_UNDEFINED,			kVK_Undefined			},	// 0xB1
	/* 178 */	{ VC_UNDEFINED,			kVK_Undefined			},	// 0xB2
	/* 179 */	{ VC_UNDEFINED,			kVK_Undefined			},	// 0xB3
	/* 180 */	{ VC_UNDEFINED,			kVK_Undefined			},	// 0xB4
	/* 181 */	{ VC_UNDEFINED,			kVK_ANSI_KeypadDivide	},	// 0xB5
	/* 182 */	{ VC_UNDEFINED,			kVK_Undefined			},	// 0xB6
	/* 183 */	{ VC_UNDEFINED,			kVK_Undefined			},	// 0xB7
	/* 184 */	{ VC_UNDEFINED,			kVK_RightOption			},	// 0xB8
	/* 185 */	{ VC_UNDEFINED,			kVK_Undefined			},	// 0xB9
	/* 186 */	{ VC_UNDEFINED,			kVK_Undefined			},	// 0xBA
	/* 187 */	{ VC_UNDEFINED,			kVK_Undefined			},	// 0xBB
	/* 188 */	{ VC_UNDEFINED,			kVK_Undefined			},	// 0xBC
	/* 189 */	{ VC_UNDEFINED,			kVK_Undefined			},	// 0xBD
	/* 190 */	{ VC_UNDEFINED,			kVK_Undefined			},	// 0xBE
	/* 191 */	{ VC_UNDEFINED,			kVK_Undefined			},	// 0xBF
	/* 192 */	{ VC_UNDEFINED,			kVK_Undefined			},	// 0xC0
	/* 193 */	{ VC_UNDEFINED,			kVK_Undefined			},	// 0xC1
	/* 194 */	{ VC_UNDEFINED,			kVK_Undefined			},	// 0xC2
	/* 195 */	{ VC_UNDEFINED,			kVK_Undefined			},	// 0xC3
	/* 196 */	{ VC_UNDEFINED,			kVK_Undefined			},	// 0xC4
	/* 197 */	{ VC_UNDEFINED,			kVK_Undefined			},	// 0xC5
	/* 198 */	{ VC_UNDEFINED,			kVK_Undefined			},	// 0xC6
	/* 199 */	{ VC_UNDEFINED,			kVK_Home				},	// 0xC7
	/* 200 */	{ VC_UNDEFINED,			kVK_UpArrow				},	// 0xC8
	/* 201 */	{ VC_UNDEFINED,			kVK_PageUp				},	// 0xC9
	/* 202 */	{ VC_UNDEFINED,			kVK_Undefined			},	// 0xCA
	/* 203 */	{ VC_UNDEFINED,			kVK_LeftArrow			},	// 0xCB
	/* 204 */	{ VC_UNDEFINED,			kVK_Undefined			},	// 0xCC
	/* 205 */	{ VC_UNDEFINED,			kVK_RightArrow			},	// 0xCD
	/* 206 */	{ VC_UNDEFINED,			kVK_Undefined			},	// 0xCE
	/* 207 */	{ VC_UNDEFINED,			kVK_End					},	// 0xCF
	/* 208 */	{ VC_UNDEFINED,			kVK_DownArrow			},	// 0xD0
	/* 209 */	{ VC_UNDEFINED,			kVK_PageDown			},	// 0xD1
	/* 210 */	{ VC_UNDEFINED,			kVK_Help				},	// 0xD2
	/* 211 */	{ VC_UNDEFINED,			kVK_ForwardDelete		},	// 0xD3
	/* 212 */	{ VC_UNDEFINED,			kVK_Undefined			},	// 0xD4
	/* 213 */	{ VC_UNDEFINED,			kVK_Undefined			},	// 0xD5
	/* 214 */	{ VC_UNDEFINED,			kVK_Undefined			},	// 0xD6
	/* 215 */	{ VC_UNDEFINED,			kVK_Undefined			},	// 0xD7
	/* 216 */	{ VC_UNDEFINED,			kVK_Undefined			},	// 0xD8
	/* 217 */	{ VC_UNDEFINED,			kVK_Undefined			},	// 0xD9
	/* 218 */	{ VC_UNDEFINED,			kVK_Undefined			},	// 0xDA
	/* 219 */	{ VC_UNDEFINED,			kVK_Command				},	// 0xDB
	/* 220 */	{ VC_UNDEFINED,			kVK_RightCommand		},	// 0xDC
	/* 221 */	{ VC_UNDEFINED,			kVK_Undefined			},	// 0xDD
	/* 222 */	{ VC_UNDEFINED,			kVK_NX_Power			},	// 0xDE
	/* 223 */	{ VC_UNDEFINED,			kVK_Undefined			},	// 0xDF
	/* 224 */	{ VC_UNDEFINED,			kVK_Undefined			},	// 0xE0
	/* 225 */	{ VC_UNDEFINED,			kVK_Undefined			},	// 0xE1
	/* 226 */	{ VC_UNDEFINED,			kVK_Undefined			},	// 0xE2
	/* 227 */	{ VC_UNDEFINED,			kVK_Undefined			},	// 0xE3
	/* 228 */	{ VC_UNDEFINED,			kVK_Undefined			},	// 0xE4
	/* 229 */	{ VC_UNDEFINED,			kVK_Undefined			},	// 0xE5
	/* 230 */	{ VC_POWER,				kVK_Undefined			},	// 0xE6
	/* 231 */	{ VC_UNDEFINED,			kVK_Undefined			},	// 0xE7
	/* 232 */	{ VC_UNDEFINED,			kVK_Undefined			},	// 0xE8
	/* 233 */	{ VC_UNDEFINED,			kVK_Undefined			},	// 0xE9
	/* 234 */	{ VC_UNDEFINED,			kVK_Undefined			},	// 0xEA
	/* 235 */	{ VC_UNDEFINED,			kVK_Undefined			},	// 0xEB
	/* 236 */	{ VC_UNDEFINED,			kVK_Undefined			},	// 0xEC
	/* 237 */	{ VC_UNDEFINED,			kVK_Undefined			},	// 0xED
	/* 238 */	{ VC_MEDIA_EJECT,		kVK_Undefined			},	// 0xEE
	/* 239 */	{ VC_UNDEFINED,			kVK_Undefined			},	// 0xEF
	/* 240 */	{ VC_MEDIA_PLAY,		kVK_Undefined			},	// 0xF0
	/* 241 */	{ VC_MEDIA_NEXT,		kVK_Undefined			},	// 0xF1
	/* 242 */	{ VC_MEDIA_PREVIOUS,	kVK_Undefined			},	// 0xF2
	/* 243 */	{ VC_UNDEFINED,			kVK_Undefined			},	// 0xF3
	/* 244 */	{ VC_UNDEFINED,			kVK_Undefined			},	// 0xF4
	/* 245 */	{ VC_UNDEFINED,			kVK_Undefined			},	// 0xF5
	/* 246 */	{ VC_UNDEFINED,			kVK_Undefined			},	// 0xF6
	/* 247 */	{ VC_UNDEFINED,			kVK_Undefined			},	// 0xF7
	/* 248 */	{ VC_UNDEFINED,			kVK_Undefined			},	// 0xF8
	/* 249 */	{ VC_UNDEFINED,			kVK_Undefined			},	// 0xF9
	/* 250 */	{ VC_UNDEFINED,			kVK_Undefined			},	// 0xFA
	/* 251 */	{ VC_UNDEFINED,			kVK_Undefined			},	// 0xFB
	/* 252 */	{ VC_UNDEFINED,			kVK_Undefined			},	// 0xFC
	/* 253 */	{ VC_UNDEFINED,			kVK_Undefined			},	// 0xFD
	/* 254 */	{ VC_UNDEFINED,			kVK_Undefined			},	// 0xFE
	/* 255 */	{ VC_UNDEFINED,			kVK_Undefined			},	// 0xFF
};

uint16_t keycode_to_scancode(UInt64 keycode) {
	uint16_t scancode = VC_UNDEFINED;

	// Bound check 0 <= keycode < 256
	if (keycode < sizeof(keycode_scancode_table) / sizeof(keycode_scancode_table[0])) {
		scancode = keycode_scancode_table[keycode][0];
	}

	return scancode;
}

UInt64 scancode_to_keycode(uint16_t scancode) {
	UInt64 keycode = kVK_Undefined;

	// Bound check 0 <= keycode < 128
	if (scancode < 128) {
		keycode = keycode_scancode_table[scancode][1];
	}
	else {
		// Calculate the upper offset.
		unsigned short i = (scancode & 0x007F) | 0x80;

		if (i < sizeof(keycode_scancode_table) / sizeof(keycode_scancode_table[1])) {
			keycode = keycode_scancode_table[i][1];
		}
	}

	return keycode;
}

void load_input_helper() {
	#if defined(USE_CARBON_LEGACY) || defined(USE_COREFOUNDATION)
	// Start with a fresh dead key state.
	//curr_deadkey_state = 0;
	#endif
}

void unload_input_helper() {
	#if defined(USE_CARBON_LEGACY) || defined(USE_COREFOUNDATION)
	if (prev_keyboard_layout != NULL) {
		// Cleanup tracking of the previous layout.
		CFRelease(prev_keyboard_layout);
		prev_keyboard_layout = NULL;
	}
	#endif
}
