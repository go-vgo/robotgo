/* Copyright (C) 2006-2017 Alexander Barker.
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Lesser General Public License as published
 * by the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
*/

#ifndef _included_input_helper
#define _included_input_helper

#include <ApplicationServices/ApplicationServices.h>
#include <Carbon/Carbon.h>	// For HIToolbox kVK_ keycodes and TIS funcitons.
#ifdef USE_IOKIT
#include <IOKit/hidsystem/ev_keymap.h>
#endif
#include <stdbool.h>


#ifndef USE_IOKIT
// Some of the system key codes that are needed from IOKit.
#define NX_KEYTYPE_SOUND_UP				0x00
#define NX_KEYTYPE_SOUND_DOWN			0x01
#define	NX_KEYTYPE_MUTE					0x07

/* Display controls...
#define NX_KEYTYPE_BRIGHTNESS_UP		0x02
#define NX_KEYTYPE_BRIGHTNESS_DOWN		0x03
#define NX_KEYTYPE_CONTRAST_UP			0x0B
#define NX_KEYTYPE_CONTRAST_DOWN		0x0C
#define NX_KEYTYPE_ILLUMINATION_UP		0x15
#define NX_KEYTYPE_ILLUMINATION_DOWN	0x16
#define NX_KEYTYPE_ILLUMINATION_TOGGLE	0x17
*/

#define NX_KEYTYPE_CAPS_LOCK			0x04
//#define NX_KEYTYPE_HELP				0x05
#define NX_POWER_KEY					0x06

#define NX_KEYTYPE_EJECT				0x0E
#define NX_KEYTYPE_PLAY					0x10
#define NX_KEYTYPE_NEXT					0x12
#define NX_KEYTYPE_PREVIOUS				0x13

/* There is no official fast-forward or rewind scan code support.*/
#define NX_KEYTYPE_FAST					0x14
#define NX_KEYTYPE_REWIND				0x15

#endif

// These virtual key codes do not appear to be defined anywhere by Apple.
#define kVK_NX_Power					0xE0 | NX_POWER_KEY			/* 0xE6 */
#define kVK_NX_Eject					0xE0 | NX_KEYTYPE_EJECT		/* 0xEE */

#define kVK_MEDIA_Play					0xE0 | NX_KEYTYPE_PLAY		/* 0xF0 */
#define kVK_MEDIA_Next					0xE0 | NX_KEYTYPE_NEXT		/* 0xF1 */
#define kVK_MEDIA_Previous				0xE0 | NX_KEYTYPE_PREVIOUS	/* 0xF2 */

#define kVK_RightCommand				0x36
#define kVK_ContextMenu					0x6E	// AKA kMenuPowerGlyph
#define kVK_Undefined					0xFF

// These button codes do not appear to be defined anywhere by Apple.
#define kVK_LBUTTON						kCGMouseButtonLeft
#define kVK_RBUTTON						kCGMouseButtonRight
#define kVK_MBUTTON						kCGMouseButtonCenter
#define kVK_XBUTTON1					3
#define kVK_XBUTTON2					4

// These button masks do not appear to be defined anywhere by Apple.
#define kCGEventFlagMaskButtonLeft		1 << 0
#define kCGEventFlagMaskButtonRight		1 << 1
#define kCGEventFlagMaskButtonCenter	1 << 2
#define kCGEventFlagMaskXButton1		1 << 3
#define kCGEventFlagMaskXButton2		1 << 4


/* Check for access to Apples accessibility API.
 */
extern bool is_accessibility_enabled();

/* Converts an OSX key code and event mask to the appropriate Unicode character
 * representation.
 */
extern UniCharCount keycode_to_unicode(CGEventRef event_ref, UniChar *buffer, UniCharCount size);

/* Converts an OSX keycode to the appropriate IOHook scancode constant.
 */
extern uint16_t keycode_to_scancode(UInt64 keycode);

/* Converts a IOHook scancode constant to the appropriate OSX keycode.
 */
extern UInt64 scancode_to_keycode(uint16_t keycode);


/* Initialize items required for KeyCodeToKeySym() and KeySymToUnicode()
 * functionality.  This method is called by OnLibraryLoad() and may need to be
 * called in combination with UnloadInputHelper() if the native keyboard layout
 * is changed.
 */
extern void load_input_helper();

/* De-initialize items required for KeyCodeToKeySym() and KeySymToUnicode()
 * functionality.  This method is called by OnLibraryUnload() and may need to be
 * called in combination with LoadInputHelper() if the native keyboard layout
 * is changed.
 */
extern void unload_input_helper();

#endif
