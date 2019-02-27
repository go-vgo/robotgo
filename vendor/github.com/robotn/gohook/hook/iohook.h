/* Copyright (C) 2006-2017 Alexander Barker.
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Lesser General Public License as published
 * by the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
*/

#ifndef __IOHOOK_H
#define __IOHOOK_H

// #include "../../base/os.h"
#include <stdarg.h>
#include <stdbool.h>
#include <stdint.h>

/* Begin Error Codes */
#define IOHOOK_SUCCESS							0x00
#define IOHOOK_FAILURE							0x01

// System level errors.
#define IOHOOK_ERROR_OUT_OF_MEMORY				0x02

// Unix specific errors.
#define IOHOOK_ERROR_X_OPEN_DISPLAY			0x20
#define IOHOOK_ERROR_X_RECORD_NOT_FOUND		0x21
#define IOHOOK_ERROR_X_RECORD_ALLOC_RANGE		0x22
#define IOHOOK_ERROR_X_RECORD_CREATE_CONTEXT	0x23
#define IOHOOK_ERROR_X_RECORD_ENABLE_CONTEXT	0x24
#define IOHOOK_ERROR_X_RECORD_GET_CONTEXT		0x25

// Windows specific errors.
#define IOHOOK_ERROR_SET_WINDOWS_HOOK_EX		0x30
#define IOHOOK_ERROR_GET_MODULE_HANDLE			0x31

// Darwin specific errors.
#define IOHOOK_ERROR_AXAPI_DISABLED			0x40
#define IOHOOK_ERROR_CREATE_EVENT_PORT			0x41
#define IOHOOK_ERROR_CREATE_RUN_LOOP_SOURCE	0x42
#define IOHOOK_ERROR_GET_RUNLOOP				0x43
#define IOHOOK_ERROR_CREATE_OBSERVER			0x44
/* End Error Codes */

/* Begin Log Levels and Function Prototype */
typedef enum _log_level {
	LOG_LEVEL_DEBUG = 1,
	LOG_LEVEL_INFO,
	LOG_LEVEL_WARN,
	LOG_LEVEL_ERROR
} log_level;

// Logger callback function prototype.
typedef bool (*logger_t)(unsigned int, const char *, ...);
/* End Log Levels and Function Prototype */

/* Begin Virtual Event Types and Data Structures */
typedef enum _event_type {
	EVENT_HOOK_ENABLED = 1,
	EVENT_HOOK_DISABLED,
	EVENT_KEY_TYPED,
	EVENT_KEY_PRESSED,
	EVENT_KEY_RELEASED,
	EVENT_MOUSE_CLICKED,
	EVENT_MOUSE_PRESSED,
	EVENT_MOUSE_RELEASED,
	EVENT_MOUSE_MOVED,
	EVENT_MOUSE_DRAGGED,
	EVENT_MOUSE_WHEEL
} event_type;

typedef struct _screen_data {
	uint8_t number;
	int16_t x;
	int16_t y;
	uint16_t width;
	uint16_t height;
} screen_data;

typedef struct _keyboard_event_data {
	uint16_t keycode;
	uint16_t rawcode;
	// uint16_t keychar;
	uint32_t keychar;
	// char *keychar;
} keyboard_event_data,
		key_pressed_event_data,
		key_released_event_data,
		key_typed_event_data;

typedef struct _mouse_event_data {
	uint16_t button;
	uint16_t clicks;
	int16_t x;
	int16_t y;
} mouse_event_data,
		mouse_pressed_event_data,
		mouse_released_event_data,
		mouse_clicked_event_data;

typedef struct _mouse_wheel_event_data {
	uint16_t clicks;
	int16_t x;
	int16_t y;
	uint8_t type;
	uint16_t amount;
	int32_t rotation;
	// int16_t rotation;
	uint8_t direction;
} mouse_wheel_event_data;

typedef struct _iohook_event {
	event_type type;
	uint64_t time;
	uint16_t mask;
	uint16_t reserved;
	union {
		keyboard_event_data keyboard;
		mouse_event_data mouse;
		mouse_wheel_event_data wheel;
	} data;
} iohook_event;

typedef void (*dispatcher_t)(iohook_event *const);
/* End Virtual Event Types and Data Structures */


/* Begin Virtual Key Codes */
#define VC_ESCAPE								0x0001

// Begin Function Keys
#define VC_F1									0x003B
#define VC_F2									0x003C
#define VC_F3									0x003D
#define VC_F4									0x003E
#define VC_F5									0x003F
#define VC_F6									0x0040
#define VC_F7									0x0041
#define VC_F8									0x0042
#define VC_F9									0x0043
#define VC_F10									0x0044
#define VC_F11									0x0057
#define VC_F12									0x0058

#define VC_F13									0x005B
#define VC_F14									0x005C
#define VC_F15									0x005D
#define VC_F16									0x0063
#define VC_F17									0x0064
#define VC_F18									0x0065
#define VC_F19									0x0066
#define VC_F20									0x0067
#define VC_F21									0x0068
#define VC_F22									0x0069
#define VC_F23									0x006A
#define VC_F24									0x006B
// End Function Keys


// Begin Alphanumeric Zone
#define VC_BACKQUOTE							0x0029

#define VC_1									0x0002
#define VC_2									0x0003
#define VC_3									0x0004
#define VC_4									0x0005
#define VC_5									0x0006
#define VC_6									0x0007
#define VC_7									0x0008
#define VC_8									0x0009
#define VC_9									0x000A
#define VC_0									0x000B

#define VC_MINUS								0x000C	// '-'
#define VC_EQUALS								0x000D	// '='
#define VC_BACKSPACE							0x000E

#define VC_TAB									0x000F
#define VC_CAPS_LOCK							0x003A

#define VC_A									0x001E
#define VC_B									0x0030
#define VC_C									0x002E
#define VC_D									0x0020
#define VC_E									0x0012
#define VC_F									0x0021
#define VC_G									0x0022
#define VC_H									0x0023
#define VC_I									0x0017
#define VC_J									0x0024
#define VC_K									0x0025
#define VC_L									0x0026
#define VC_M									0x0032
#define VC_N									0x0031
#define VC_O									0x0018
#define VC_P									0x0019
#define VC_Q									0x0010
#define VC_R									0x0013
#define VC_S									0x001F
#define VC_T									0x0014
#define VC_U									0x0016
#define VC_V									0x002F
#define VC_W									0x0011
#define VC_X									0x002D
#define VC_Y									0x0015
#define VC_Z									0x002C

#define VC_OPEN_BRACKET							0x001A	// '['
#define VC_CLOSE_BRACKET						0x001B	// ']'
#define VC_BACK_SLASH							0x002B	// '\'

#define VC_SEMICOLON							0x0027	// ';'
#define VC_QUOTE								0x0028
#define VC_ENTER								0x001C

#define VC_COMMA								0x0033	// ','
#define VC_PERIOD								0x0034	// '.'
#define VC_SLASH								0x0035	// '/'

#define VC_SPACE								0x0039
// End Alphanumeric Zone


#define VC_PRINTSCREEN							0x0E37
#define VC_SCROLL_LOCK							0x0046
#define VC_PAUSE								0x0E45


// Begin Edit Key Zone
#define VC_INSERT								0x0E52
#define VC_DELETE								0x0E53
#define VC_HOME									0x0E47
#define VC_END									0x0E4F
#define VC_PAGE_UP								0x0E49
#define VC_PAGE_DOWN							0x0E51
// End Edit Key Zone


// Begin Cursor Key Zone
#define VC_UP									0xE048
#define VC_LEFT									0xE04B
#define VC_CLEAR								0xE04C
#define VC_RIGHT								0xE04D
#define VC_DOWN									0xE050
// End Cursor Key Zone


// Begin Numeric Zone
#define VC_NUM_LOCK								0x0045
#define VC_KP_DIVIDE							0x0E35
#define VC_KP_MULTIPLY							0x0037
#define VC_KP_SUBTRACT							0x004A
#define VC_KP_EQUALS							0x0E0D
#define VC_KP_ADD								0x004E
#define VC_KP_ENTER								0x0E1C
#define VC_KP_SEPARATOR							0x0053

#define VC_KP_1									0x004F
#define VC_KP_2									0x0050
#define VC_KP_3									0x0051
#define VC_KP_4									0x004B
#define VC_KP_5									0x004C
#define VC_KP_6									0x004D
#define VC_KP_7									0x0047
#define VC_KP_8									0x0048
#define VC_KP_9									0x0049
#define VC_KP_0									0x0052

#define VC_KP_END								0xEE00 | VC_KP_1
#define VC_KP_DOWN								0xEE00 | VC_KP_2
#define VC_KP_PAGE_DOWN							0xEE00 | VC_KP_3
#define VC_KP_LEFT								0xEE00 | VC_KP_4
#define VC_KP_CLEAR								0xEE00 | VC_KP_5
#define VC_KP_RIGHT								0xEE00 | VC_KP_6
#define VC_KP_HOME								0xEE00 | VC_KP_7
#define VC_KP_UP								0xEE00 | VC_KP_8
#define VC_KP_PAGE_UP							0xEE00 | VC_KP_9
#define VC_KP_INSERT							0xEE00 | VC_KP_0
#define VC_KP_DELETE							0xEE00 | VC_KP_SEPARATOR
// End Numeric Zone


// Begin Modifier and Control Keys
#define VC_SHIFT_L								0x002A
#define VC_SHIFT_R								0x0036
#define VC_CONTROL_L							0x001D
#define VC_CONTROL_R							0x0E1D
#define VC_ALT_L								0x0038	// Option or Alt Key
#define VC_ALT_R								0x0E38	// Option or Alt Key
#define VC_META_L								0x0E5B	// Windows or Command Key
#define VC_META_R								0x0E5C	// Windows or Command Key
#define VC_CONTEXT_MENU							0x0E5D
// End Modifier and Control Keys


// Begin Media Control Keys
#define VC_POWER								0xE05E
#define VC_SLEEP								0xE05F
#define VC_WAKE									0xE063

#define VC_MEDIA_PLAY							0xE022
#define VC_MEDIA_STOP							0xE024
#define VC_MEDIA_PREVIOUS						0xE010
#define VC_MEDIA_NEXT							0xE019
#define VC_MEDIA_SELECT							0xE06D
#define VC_MEDIA_EJECT							0xE02C

#define VC_VOLUME_MUTE							0xE020
#define VC_VOLUME_UP							0xE030
#define VC_VOLUME_DOWN							0xE02E

#define VC_APP_MAIL								0xE06C
#define VC_APP_CALCULATOR						0xE021
#define VC_APP_MUSIC							0xE03C
#define VC_APP_PICTURES							0xE064

#define VC_BROWSER_SEARCH						0xE065
#define VC_BROWSER_HOME							0xE032
#define VC_BROWSER_BACK							0xE06A
#define VC_BROWSER_FORWARD						0xE069
#define VC_BROWSER_STOP							0xE068
#define VC_BROWSER_REFRESH						0xE067
#define VC_BROWSER_FAVORITES					0xE066
// End Media Control Keys

// Begin Japanese Language Keys
#define VC_KATAKANA								0x0070
#define VC_UNDERSCORE							0x0073
#define VC_FURIGANA								0x0077
#define VC_KANJI								0x0079
#define VC_HIRAGANA								0x007B
#define VC_YEN									0x007D
#define VC_KP_COMMA								0x007E
// End Japanese Language Keys

// Begin Sun keyboards
#define VC_SUN_HELP								0xFF75

#define VC_SUN_STOP								0xFF78
#define VC_SUN_PROPS							0xFF76
#define VC_SUN_FRONT							0xFF77
#define VC_SUN_OPEN								0xFF74
#define VC_SUN_FIND								0xFF7E
#define VC_SUN_AGAIN							0xFF79
#define VC_SUN_UNDO								0xFF7A
#define VC_SUN_COPY								0xFF7C
#define VC_SUN_INSERT							0xFF7D
#define VC_SUN_CUT								0xFF7B
// End Sun keyboards

#define VC_UNDEFINED							0x0000	// KeyCode Unknown

#define CHAR_UNDEFINED							0xFFFF	// CharCode Unknown
/* End Virtual Key Codes */


/* Begin Virtual Modifier Masks */
#define MASK_SHIFT_L							1 << 0
#define MASK_CTRL_L								1 << 1
#define MASK_META_L								1 << 2
#define MASK_ALT_L								1 << 3

#define MASK_SHIFT_R							1 << 4
#define MASK_CTRL_R								1 << 5
#define MASK_META_R								1 << 6
#define MASK_ALT_R								1 << 7

#define MASK_SHIFT								MASK_SHIFT_L | MASK_SHIFT_R
#define MASK_CTRL								MASK_CTRL_L  | MASK_CTRL_R
#define MASK_META								MASK_META_L  | MASK_META_R
#define MASK_ALT								MASK_ALT_L   | MASK_ALT_R

#define MASK_BUTTON1							1 << 8
#define MASK_BUTTON2							1 << 9
#define MASK_BUTTON3							1 << 10
#define MASK_BUTTON4							1 << 11
#define MASK_BUTTON5							1 << 12

#define MASK_NUM_LOCK							1 << 13
#define MASK_CAPS_LOCK							1 << 14
#define MASK_SCROLL_LOCK						1 << 15
/* End Virtual Modifier Masks */


/* Begin Virtual Mouse Buttons */
#define MOUSE_NOBUTTON							0	// Any Button
#define MOUSE_BUTTON1							1	// Left Button
#define MOUSE_BUTTON2							2	// Right Button
#define MOUSE_BUTTON3							3	// Middle Button
#define MOUSE_BUTTON4							4	// Extra Mouse Button
#define MOUSE_BUTTON5							5	// Extra Mouse Button

#define WHEEL_UNIT_SCROLL						1
#define WHEEL_BLOCK_SCROLL						2

#define WHEEL_VERTICAL_DIRECTION				3
#define WHEEL_HORIZONTAL_DIRECTION              4
/* End Virtual Mouse Buttons */


#ifdef _WIN32
#define IOHOOK_API __declspec(dllexport)
#else
#define IOHOOK_API
#endif

#ifdef __cplusplus
extern "C" {
#endif

	// Set the logger callback functions.
	IOHOOK_API void hook_set_logger_proc(logger_t logger_proc);

	// Send a virtual event back to the system.
	IOHOOK_API void hook_post_event(iohook_event * const event);

	// Set the event callback function.
	IOHOOK_API void hook_set_dispatch_proc(dispatcher_t dispatch_proc);

	// Insert the event hook.
	IOHOOK_API int hook_run();

	// Withdraw the event hook.
	IOHOOK_API int hook_stop();

	// Retrieves an array of screen data for each available monitor.
	IOHOOK_API screen_data* hook_create_screen_info(unsigned char *count);

	// Retrieves the keyboard auto repeat rate.
	IOHOOK_API long int hook_get_auto_repeat_rate();

	// Retrieves the keyboard auto repeat delay.
	IOHOOK_API long int hook_get_auto_repeat_delay();

	// Retrieves the mouse acceleration multiplier.
	IOHOOK_API long int hook_get_pointer_acceleration_multiplier();

	// Retrieves the mouse acceleration threshold.
	IOHOOK_API long int hook_get_pointer_acceleration_threshold();

	// Retrieves the mouse sensitivity.
	IOHOOK_API long int hook_get_pointer_sensitivity();

	// Retrieves the double/triple click interval.
	IOHOOK_API long int hook_get_multi_click_time();

#ifdef __cplusplus
}
#endif

#endif
