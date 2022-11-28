#pragma once
#ifndef KEYPRESS_H
#define KEYPRESS_H

#include <stdlib.h>
#include "../base/os.h"
#include "../base/types.h"

#include "keycode.h"
#include <stdbool.h>

#if defined(IS_MACOSX)
	typedef enum {
		MOD_NONE = 0,
		MOD_META = kCGEventFlagMaskCommand,
		MOD_ALT = kCGEventFlagMaskAlternate,
		MOD_CONTROL = kCGEventFlagMaskControl,
		MOD_SHIFT = kCGEventFlagMaskShift
	} MMKeyFlags;
#elif defined(USE_X11)
	enum _MMKeyFlags {
		MOD_NONE = 0,
		MOD_META = Mod4Mask,
		MOD_ALT = Mod1Mask,
		MOD_CONTROL = ControlMask,
		MOD_SHIFT = ShiftMask
	};
	typedef unsigned int MMKeyFlags;
#elif defined(IS_WINDOWS)
	enum _MMKeyFlags {
		MOD_NONE = 0,
		/* These are already defined by the Win32 API */
		/* MOD_ALT = 0,
		MOD_CONTROL = 0,
		MOD_SHIFT = 0, */
		MOD_META = MOD_WIN
	};
	typedef unsigned int MMKeyFlags;
#endif

#if defined(IS_WINDOWS)
	/* Send win32 key event for given key. */
	void win32KeyEvent(int key, MMKeyFlags flags, uintptr pid, int8_t isPid);
#endif

#endif /* KEYPRESS_H */
