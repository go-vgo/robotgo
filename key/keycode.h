#pragma once
#ifndef KEYCODE_H
#define KEYCODE_H

#include "../base/os.h"

#ifdef __cplusplus
extern "C"
{
#endif

#if defined(IS_MACOSX)

#include <Carbon/Carbon.h> /* Really only need <HIToolbox/Events.h> */
#include <ApplicationServices/ApplicationServices.h>
#import <IOKit/hidsystem/ev_keymap.h>

enum _MMKeyCode {
	K_NOT_A_KEY = 9999,
	K_BACKSPACE = kVK_Delete,
	K_DELETE = kVK_ForwardDelete,
	K_RETURN = kVK_Return,
	K_TAB = kVK_Tab,
	K_ESCAPE = kVK_Escape,
	K_UP = kVK_UpArrow,
	K_DOWN = kVK_DownArrow,
	K_RIGHT = kVK_RightArrow,
	K_LEFT = kVK_LeftArrow,
	K_HOME = kVK_Home,
	K_END = kVK_End,
	K_PAGEUP = kVK_PageUp,
	K_PAGEDOWN = kVK_PageDown,
	K_F1 = kVK_F1,
	K_F2 = kVK_F2,
	K_F3 = kVK_F3,
	K_F4 = kVK_F4,
	K_F5 = kVK_F5,
	K_F6 = kVK_F6,
	K_F7 = kVK_F7,
	K_F8 = kVK_F8,
	K_F9 = kVK_F9,
	K_F10 = kVK_F10,
	K_F11 = kVK_F11,
	K_F12 = kVK_F12,
	K_F13 = kVK_F13,
	K_F14 = kVK_F14,
	K_F15 = kVK_F15,
	K_F16 = kVK_F16,
	K_F17 = kVK_F17,
	K_F18 = kVK_F18,
	K_F19 = kVK_F19,
	K_F20 = kVK_F20,
	K_F21 = K_NOT_A_KEY,
	K_F22 = K_NOT_A_KEY,
	K_F23 = K_NOT_A_KEY,
	K_F24 = K_NOT_A_KEY,
	K_META = kVK_Command,
	K_ALT = kVK_Option,
	K_CONTROL = kVK_Control,
	K_SHIFT = kVK_Shift,
	K_RIGHTSHIFT = kVK_RightShift,
	K_CAPSLOCK = kVK_CapsLock,
	K_SPACE = kVK_Space,
	K_INSERT = K_NOT_A_KEY,
	K_PRINTSCREEN = K_NOT_A_KEY,

	K_NUMPAD_0 = kVK_ANSI_Keypad0,
	K_NUMPAD_1 = kVK_ANSI_Keypad1,
	K_NUMPAD_2 = kVK_ANSI_Keypad2,
	K_NUMPAD_3 = kVK_ANSI_Keypad3,
	K_NUMPAD_4 = kVK_ANSI_Keypad4,
	K_NUMPAD_5 = kVK_ANSI_Keypad5,
	K_NUMPAD_6 = kVK_ANSI_Keypad6,
	K_NUMPAD_7 = kVK_ANSI_Keypad7,
	K_NUMPAD_8 = kVK_ANSI_Keypad8,
	K_NUMPAD_9 = kVK_ANSI_Keypad9,

	K_AUDIO_VOLUME_MUTE = 1007,
	K_AUDIO_VOLUME_DOWN = 1001,
	K_AUDIO_VOLUME_UP = 1000,
	K_AUDIO_PLAY = 1016,
	K_AUDIO_STOP = K_NOT_A_KEY,
	K_AUDIO_PAUSE = 1016,
	K_AUDIO_PREV = 1018,
	K_AUDIO_NEXT = 1017,
	K_AUDIO_REWIND = K_NOT_A_KEY,
	K_AUDIO_FORWARD = K_NOT_A_KEY,
	K_AUDIO_REPEAT = K_NOT_A_KEY,
	K_AUDIO_RANDOM = K_NOT_A_KEY,

	K_LIGHTS_MON_UP = 1002,
	K_LIGHTS_MON_DOWN = 1003,
	K_LIGHTS_KBD_TOGGLE = 1023,
	K_LIGHTS_KBD_UP = 1021,
	K_LIGHTS_KBD_DOWN = 1022
};

typedef CGKeyCode MMKeyCode;

#elif defined(USE_X11)

#include <X11/Xutil.h>
#include <X11/XF86keysym.h>

enum _MMKeyCode {
	K_NOT_A_KEY = 9999,
	K_BACKSPACE = XK_BackSpace,
	K_DELETE = XK_Delete,
	K_RETURN = XK_Return,
	K_TAB = XK_Tab,
	K_ESCAPE = XK_Escape,
	K_UP = XK_Up,
	K_DOWN = XK_Down,
	K_RIGHT = XK_Right,
	K_LEFT = XK_Left,
	K_HOME = XK_Home,
	K_END = XK_End,
	K_PAGEUP = XK_Page_Up,
	K_PAGEDOWN = XK_Page_Down,
	K_F1 = XK_F1,
	K_F2 = XK_F2,
	K_F3 = XK_F3,
	K_F4 = XK_F4,
	K_F5 = XK_F5,
	K_F6 = XK_F6,
	K_F7 = XK_F7,
	K_F8 = XK_F8,
	K_F9 = XK_F9,
	K_F10 = XK_F10,
	K_F11 = XK_F11,
	K_F12 = XK_F12,
	K_F13 = XK_F13,
	K_F14 = XK_F14,
	K_F15 = XK_F15,
	K_F16 = XK_F16,
	K_F17 = XK_F17,
	K_F18 = XK_F18,
	K_F19 = XK_F19,
	K_F20 = XK_F20,
	K_F21 = XK_F21,
	K_F22 = XK_F22,
	K_F23 = XK_F23,
	K_F24 = XK_F24,
	K_META = XK_Super_L,
	K_ALT = XK_Alt_L,
	K_CONTROL = XK_Control_L,
	K_SHIFT = XK_Shift_L,
	K_RIGHTSHIFT = XK_Shift_R,
	K_CAPSLOCK = XK_Shift_Lock,
	K_SPACE = XK_space,
	K_INSERT = XK_Insert,
	K_PRINTSCREEN = XK_Print,

	K_NUMPAD_0 = K_NOT_A_KEY,
	K_NUMPAD_1 = K_NOT_A_KEY,
	K_NUMPAD_2 = K_NOT_A_KEY,
	K_NUMPAD_3 = K_NOT_A_KEY,
	K_NUMPAD_4 = K_NOT_A_KEY,
	K_NUMPAD_5 = K_NOT_A_KEY,
	K_NUMPAD_6 = K_NOT_A_KEY,
	K_NUMPAD_7 = K_NOT_A_KEY,
	K_NUMPAD_8 = K_NOT_A_KEY,
	K_NUMPAD_9 = K_NOT_A_KEY,

	K_AUDIO_VOLUME_MUTE = XF86XK_AudioMute,
	K_AUDIO_VOLUME_DOWN = XF86XK_AudioLowerVolume,
	K_AUDIO_VOLUME_UP = XF86XK_AudioRaiseVolume,
	K_AUDIO_PLAY = XF86XK_AudioPlay,
	K_AUDIO_STOP = XF86XK_AudioStop,
	K_AUDIO_PAUSE = XF86XK_AudioPause,
	K_AUDIO_PREV = XF86XK_AudioPrev,
	K_AUDIO_NEXT = XF86XK_AudioNext,
	K_AUDIO_REWIND = XF86XK_AudioRewind,
	K_AUDIO_FORWARD = XF86XK_AudioForward,
	K_AUDIO_REPEAT = XF86XK_AudioRepeat,
	K_AUDIO_RANDOM = XF86XK_AudioRandomPlay,

	K_LIGHTS_MON_UP = XF86XK_MonBrightnessUp,
	K_LIGHTS_MON_DOWN = XF86XK_MonBrightnessDown,
	K_LIGHTS_KBD_TOGGLE = XF86XK_KbdLightOnOff,
	K_LIGHTS_KBD_UP = XF86XK_KbdBrightnessUp,
	K_LIGHTS_KBD_DOWN = XF86XK_KbdBrightnessDown
};

typedef KeySym MMKeyCode;

#elif defined(IS_WINDOWS)

enum _MMKeyCode {
	K_NOT_A_KEY = 9999,
	K_BACKSPACE = VK_BACK,
	K_DELETE = VK_DELETE,
	K_RETURN = VK_RETURN,
	K_TAB = VK_TAB,
	K_ESCAPE = VK_ESCAPE,
	K_UP = VK_UP,
	K_DOWN = VK_DOWN,
	K_RIGHT = VK_RIGHT,
	K_LEFT = VK_LEFT,
	K_HOME = VK_HOME,
	K_END = VK_END,
	K_PAGEUP = VK_PRIOR,
	K_PAGEDOWN = VK_NEXT,
	K_F1 = VK_F1,
	K_F2 = VK_F2,
	K_F3 = VK_F3,
	K_F4 = VK_F4,
	K_F5 = VK_F5,
	K_F6 = VK_F6,
	K_F7 = VK_F7,
	K_F8 = VK_F8,
	K_F9 = VK_F9,
	K_F10 = VK_F10,
	K_F11 = VK_F11,
	K_F12 = VK_F12,
	K_F13 = VK_F13,
	K_F14 = VK_F14,
	K_F15 = VK_F15,
	K_F16 = VK_F16,
	K_F17 = VK_F17,
	K_F18 = VK_F18,
	K_F19 = VK_F19,
	K_F20 = VK_F20,
	K_F21 = VK_F21,
	K_F22 = VK_F22,
	K_F23 = VK_F23,
	K_F24 = VK_F24,
	K_META = VK_LWIN,
	K_CONTROL = VK_CONTROL,
	K_SHIFT = VK_SHIFT,
	K_RIGHTSHIFT = VK_RSHIFT,
	K_ALT = VK_MENU,
	K_CAPSLOCK = VK_CAPITAL,
	K_SPACE = VK_SPACE,
	K_PRINTSCREEN = VK_SNAPSHOT,
	K_INSERT = VK_INSERT,

	K_NUMPAD_0 = VK_NUMPAD0,
	K_NUMPAD_1 = VK_NUMPAD1,
	K_NUMPAD_2 = VK_NUMPAD2,
	K_NUMPAD_3 = VK_NUMPAD3,
	K_NUMPAD_4 = VK_NUMPAD4,
	K_NUMPAD_5 = VK_NUMPAD5,
	K_NUMPAD_6 = VK_NUMPAD6,
	K_NUMPAD_7 = VK_NUMPAD7,
	K_NUMPAD_8 = VK_NUMPAD8,
	K_NUMPAD_9 = VK_NUMPAD9,

	K_AUDIO_VOLUME_MUTE = VK_VOLUME_MUTE,
	K_AUDIO_VOLUME_DOWN = VK_VOLUME_DOWN,
	K_AUDIO_VOLUME_UP = VK_VOLUME_UP,
	K_AUDIO_PLAY = VK_MEDIA_PLAY_PAUSE,
	K_AUDIO_STOP = VK_MEDIA_STOP,
	K_AUDIO_PAUSE = VK_MEDIA_PLAY_PAUSE,
	K_AUDIO_PREV = VK_MEDIA_PREV_TRACK,
	K_AUDIO_NEXT = VK_MEDIA_NEXT_TRACK,
	K_AUDIO_REWIND = K_NOT_A_KEY,
	K_AUDIO_FORWARD = K_NOT_A_KEY,
	K_AUDIO_REPEAT = K_NOT_A_KEY,
	K_AUDIO_RANDOM = K_NOT_A_KEY,

	K_LIGHTS_MON_UP = K_NOT_A_KEY,
	K_LIGHTS_MON_DOWN = K_NOT_A_KEY,
	K_LIGHTS_KBD_TOGGLE = K_NOT_A_KEY,
	K_LIGHTS_KBD_UP = K_NOT_A_KEY,
	K_LIGHTS_KBD_DOWN = K_NOT_A_KEY
};

typedef int MMKeyCode;

#endif

/* Returns the keyCode corresponding to the current keyboard layout for the
 * given ASCII character. */
MMKeyCode keyCodeForChar(const char c);

#endif /* KEYCODE_H */

#ifdef __cplusplus
}
#endif
