
#ifdef HAVE_CONFIG_H
	#include <config.h>
#endif

#include <stdbool.h>
#include <stddef.h>
#include <stdint.h>
#include <stdio.h>
#include <string.h>
#include <windows.h>

#include "../iohook.h"
#include "../logger_c.h"
#include "input.h"

static const uint16_t keycode_scancode_table[][2] = {
	/* idx		{ vk_code,				scancode				}, */
	/*   0 */	{ VC_UNDEFINED,			0x0000					},	// 0x00
	/*   1 */	{ MOUSE_BUTTON1,		VK_ESCAPE				},	// 0x01
	/*   2 */	{ MOUSE_BUTTON2,		0x0031					},	// 0x02
	/*   3 */	{ VC_UNDEFINED,			0x0032					},	// 0x03 VK_CANCEL
	/*   4 */	{ MOUSE_BUTTON3,		0x0033					},	// 0x04
	/*   5 */	{ MOUSE_BUTTON4,		0x0034					},	// 0x05
	/*   6 */	{ MOUSE_BUTTON5,		0x0035					},	// 0x06
	/*   7 */	{ VC_UNDEFINED,			0x0036					},	// 0x07							Undefined
	/*   8 */	{ VC_BACKSPACE,			0x0037					},	// 0x08 VK_BACK
	/*   9 */	{ VC_TAB,				0x0038					},	// 0x09 VK_TAB
	/*  10 */	{ VC_UNDEFINED,			0x0039					},	// 0x0A							Reserved
	/*  11 */	{ VC_UNDEFINED,			0x0030					},	// 0x0B							Reserved
	/*  12 */	{ VC_CLEAR,				VK_OEM_MINUS			},	// 0x0C VK_CLEAR
	/*  13 */	{ VC_ENTER,				VK_OEM_PLUS				},	// 0x0D VK_RETURN
	/*  14 */	{ VC_UNDEFINED,			VK_BACK					},	// 0x0E							Undefined
	/*  15 */	{ VC_UNDEFINED,			VK_TAB					},	// 0x0F							Undefined
	/*  16 */	{ VC_SHIFT_L,			0x0051					},	// 0x10 VK_SHIFT
	/*  17 */	{ VC_CONTROL_L,			0x0057					},	// 0x11 VK_CONTROL
	/*  18 */	{ VC_ALT_L,				0x0045					},	// 0x12 VK_MENU					ALT key
	/*  19 */	{ VC_PAUSE,				0x0052					},	// 0x13 VK_PAUSE
	/*  20 */	{ VC_CAPS_LOCK,			0x0054					},	// 0x14 VK_CAPITAL				CAPS LOCK key
	/*  21 */	{ VC_KATAKANA,			0x0059					},	// 0x15 VK_KANA					IME Kana mode
	/*  22 */	{ VC_UNDEFINED,			0x0055					},	// 0x16							Undefined
	/*  23 */	{ VC_UNDEFINED,			0x0049					},	// 0x17 VK_JUNJA				IME Junja mode
	/*  24 */	{ VC_UNDEFINED,			0x004F					},	// 0x18 VK_FINAL
	/*  25 */	{ VC_KANJI,				0x0050					},	// 0x19 VK_KANJI / VK_HANJA		IME Kanji / Hanja mode
	/*  26 */	{ VC_UNDEFINED,			0x00DB					},	// 0x1A Undefined
	/*  27 */	{ VC_ESCAPE,			0x00DD					},	// 0x1B	VK_ESCAPE				ESC key
	/*  28 */	{ VC_UNDEFINED,			VK_RETURN				},	// 0x1C VK_CONVERT				IME convert// 0x1C
	/*  29 */	{ VC_UNDEFINED,			VK_LCONTROL				},	// 0x1D VK_NONCONVERT			IME nonconvert
	/*  30 */	{ VC_UNDEFINED,			0x0041					},	// 0x1E VK_ACCEPT				IME accept
	/*  31 */	{ VC_UNDEFINED,			0x0053					},	// 0x1F VK_MODECHANGE			IME mode change request
	/*  32 */	{ VC_SPACE,				0x0044					},	// 0x20 VK_SPACE				SPACEBAR
	/*  33 */	{ VC_PAGE_UP,			0x0046					},	// 0x21 VK_PRIOR				PAGE UP key
	/*  34 */	{ VC_PAGE_DOWN,			0x0047					},	// 0x22 VK_NEXT					PAGE DOWN key
	/*  35 */	{ VC_END,				0x0048					},	// 0x23 VK_END					END key
	/*  36 */	{ VC_HOME,				0x004A					},	// 0x24 VK_HOME					HOME key
	/*  37 */	{ VC_LEFT,				0x004B					},	// 0x25 VK_LEFT					LEFT ARROW key
	/*  38 */	{ VC_UP,				0x004C					},	// 0x26 VK_UP					UP ARROW key
	/*  39 */	{ VC_RIGHT,				VK_OEM_1				},	// 0x27 VK_RIGHT				RIGHT ARROW key
	/*  40 */	{ VC_DOWN,				VK_OEM_7				},	// 0x28 VK_DOWN					DOWN ARROW key
	/*  41 */	{ VC_UNDEFINED,			VK_OEM_3				},	// 0x29 VK_SELECT				SELECT key
	/*  42 */	{ VC_UNDEFINED,			VK_LSHIFT				},	// 0x2A VK_PRINT				PRINT key
	/*  43 */	{ VC_UNDEFINED,			VK_OEM_5				},	// 0x2B VK_EXECUTE				EXECUTE key
	/*  44 */	{ VC_PRINTSCREEN,		0x005A					},	// 0x2C VK_SNAPSHOT				PRINT SCREEN key
	/*  45 */	{ VC_INSERT,			0x0058					},	// 0x2D VK_INSERT				INS key
	/*  46 */	{ VC_DELETE,			0x0043					},	// 0x2E VK_DELETE				DEL key
	/*  47 */	{ VC_UNDEFINED,			0x0056					},	// 0x2F VK_HELP					HELP key
	/*  48 */	{ VC_0,					0x0042					},	// 0x30							0 key
	/*  49 */	{ VC_1,					0x004E					},	// 0x31							1 key
	/*  50 */	{ VC_2,					0x004D					},	// 0x32							2 key
	/*  51 */	{ VC_3,					VK_OEM_COMMA			},	// 0x33							3 key
	/*  52 */	{ VC_4,					VK_OEM_PERIOD			},	// 0x34							4 key
	/*  53 */	{ VC_5,					VK_OEM_2				},	// 0x35							5 key
	/*  54 */	{ VC_6,					VK_RSHIFT				},	// 0x36							6 key
	/*  55 */	{ VC_7,					VK_MULTIPLY				},	// 0x37							7 key
	/*  56 */	{ VC_8,					VK_LMENU				},	// 0x38							8 key
	/*  57 */	{ VC_9,					VK_SPACE				},	// 0x39							9 key
	/*  58 */	{ VC_UNDEFINED,			VK_CAPITAL				},	// 0x3A							Undefined
	/*  59 */	{ VC_UNDEFINED,			VK_F1					},	// 0x3B							Undefined
	/*  60 */	{ VC_UNDEFINED,			VK_F2					},	// 0x3C							Undefined
	/*  61 */	{ VC_UNDEFINED,			VK_F3					},	// 0x3D							Undefined
	/*  62 */	{ VC_UNDEFINED,			VK_F4					},	// 0x3E							Undefined
	/*  63 */	{ VC_UNDEFINED,			VK_F5					},	// 0x3F							Undefined
	/*  64 */	{ VC_UNDEFINED,			VK_F6					},	// 0x40							Undefined
	/*  65 */	{ VC_A,					VK_F7					},	// 0x41							A key
	/*  66 */	{ VC_B,					VK_F8					},	// 0x42							B key
	/*  67 */	{ VC_C,					VK_F9					},	// 0x43							C key
	/*  68 */	{ VC_D,					VK_F10					},	// 0x44							D key
	/*  69 */	{ VC_E,					VK_NUMLOCK				},	// 0x45							E key
	/*  70 */	{ VC_F,					VK_SCROLL				},	// 0x46							F key
	/*  71 */	{ VC_G,					VK_NUMPAD7				},	// 0x47							G key
	/*  72 */	{ VC_H,					VK_NUMPAD8				},	// 0x48							H key
	/*  73 */	{ VC_I,					VK_NUMPAD9				},	// 0x49							I key
	/*  74 */	{ VC_J,					VK_SUBTRACT				},	// 0x4A							J key
	/*  75 */	{ VC_K,					VK_NUMPAD4				},	// 0x4B							K key
	/*  76 */	{ VC_L,					VK_NUMPAD5				},	// 0x4C							L key
	/*  77 */	{ VC_M,					VK_NUMPAD6				},	// 0x4D							M key
	/*  78 */	{ VC_N,					VK_ADD					},	// 0x4E							N key
	/*  79 */	{ VC_O,					VK_NUMPAD1				},	// 0x4F							O key
	/*  80 */	{ VC_P,					VK_NUMPAD2				},	// 0x50							P key
	/*  81 */	{ VC_Q,					VK_NUMPAD3				},	// 0x51							Q key
	/*  82 */	{ VC_R,					VK_NUMPAD0				},	// 0x52							R key
	/*  83 */	{ VC_S,					VK_DECIMAL				},	// 0x53							S key
	/*  84 */	{ VC_T,					0x0000					},	// 0x54							T key
	/*  85 */	{ VC_U,					0x0000					},	// 0x55							U key
	/*  86 */	{ VC_V,					0x0000					},	// 0x56							V key
	/*  87 */	{ VC_W,					VK_F11					},	// 0x57							W key
	/*  88 */	{ VC_X,					VK_F12					},	// 0x58							X key
	/*  89 */	{ VC_Y,					0x0000					},	// 0x59							Y key
	/*  90 */	{ VC_Z,					0x0000					},	// 0x5A							Z key
	/*  91 */	{ VC_META_L,			VK_F13					},	// 0x5B VK_LWIN 				Left Windows key (Natural keyboard)
	/*  92 */	{ VC_META_R,			VK_F14					},	// 0x5C VK_RWIN					Right Windows key (Natural keyboard)
	/*  93 */	{ VC_CONTEXT_MENU,		VK_F15					},	// 0x5D VK_APPS					Applications key (Natural keyboard)
	/*  94 */	{ VC_UNDEFINED,			0x0000					},	// 0x5E Reserved
	/*  95 */	{ VC_SLEEP,				0x0000					},	// 0x5F VK_SLEEP				Computer Sleep key
	/*  96 */	{ VC_KP_0,				0x0000					},	// 0x60 VK_NUMPAD0				Numeric keypad 0 key
	/*  97 */	{ VC_KP_1,				0x0000					},	// 0x61 VK_NUMPAD1				Numeric keypad 1 key
	/*  98 */	{ VC_KP_2,				0x0000					},	// 0x62 VK_NUMPAD2				Numeric keypad 2 key
	/*  99 */	{ VC_KP_3,				VK_F16					},	// 0x63 VK_NUMPAD3				Numeric keypad 3 key
	/* 100 */	{ VC_KP_4,				VK_F17					},	// 0x64 VK_NUMPAD4				Numeric keypad 4 key
	/* 101 */	{ VC_KP_5,				VK_F18					},	// 0x65 VK_NUMPAD5				Numeric keypad 5 key
	/* 102 */	{ VC_KP_6,				VK_F19					},	// 0x66 VK_NUMPAD6				Numeric keypad 6 key
	/* 103 */	{ VC_KP_7,				VK_F20					},	// 0x67 VK_NUMPAD7				Numeric keypad 7 key
	/* 104 */	{ VC_KP_8,				VK_F21					},	// 0x68 VK_NUMPAD8				Numeric keypad 8 key
	/* 105 */	{ VC_KP_9,				VK_F22					},	// 0x69 VK_NUMPAD9				Numeric keypad 9 key
	/* 106 */	{ VC_KP_MULTIPLY,		VK_F23					},	// 0x6A VK_MULTIPLY				Multiply key
	/* 107 */	{ VC_KP_ADD,			VK_F24					},	// 0x6B VK_ADD					Add key
	/* 108 */	{ VC_UNDEFINED,			0x0000					},	// 0x6C VK_SEPARATOR			Separator key
	/* 109 */	{ VC_KP_SUBTRACT,		0x0000					},	// 0x6D VK_SUBTRACT				Subtract key
	/* 110 */	{ VC_KP_SEPARATOR,		0x0000					},	// 0x6E VK_DECIMAL				Decimal key
	/* 111 */	{ VC_KP_DIVIDE,			0x0000					},	// 0x6F VK_DIVIDE				Divide key
	/* 112 */	{ VC_F1,				VK_KANA					},	// 0x70 VK_F1					F1 key
	/* 113 */	{ VC_F2,				0x0000					},	// 0x71 VK_F2					F2 key
	/* 114 */	{ VC_F3,				0x0000					},	// 0x72 VK_F3					F3 key
	/* 115 */	{ VC_F4,				0x0000					},	// 0x73 VK_F4					F4 key
	/* 116 */	{ VC_F5,				0x0000					},	// 0x74 VK_F5					F5 key
	/* 117 */	{ VC_F6,				0x0000					},	// 0x75 VK_F6					F6 key
	/* 118 */	{ VC_F7,				0x0000					},	// 0x76 VK_F7					F7 key
	/* 119 */	{ VC_F8,				0x0000					},	// 0x77 VK_F8					F8 key
	/* 120 */	{ VC_F9,				0x0000					},	// 0x78 VK_F9					F9 key
	/* 121 */	{ VC_F10,				VK_KANJI				},	// 0x79 VK_F10					F10 key
	/* 122 */	{ VC_F11,				0x0000					},	// 0x7A VK_F11					F11 key
	/* 123 */	{ VC_F12,				0x0000					},	// 0x7B VK_F12					F12 key
	/* 124 */	{ VC_F13,				0x0000					},	// 0x7C VK_F13					F13 key
	/* 125 */	{ VC_F14,				VK_OEM_8				},	// 0x7D VK_F14					F14 key
	/* 126 */	{ VC_F15,				0x0000					},	// 0x7E VK_F15					F15 key
	/* 127 */	{ VC_F16,				0x0000					},	// 0x7F VK_F16					F16 key

	//			No Offset				Offset (i & 0x007F) | 0x80

	/* 128 */	{ VC_F17,				0x0000					},	// 0x80 VK_F17					F17 key
	/* 129 */	{ VC_F18,				0x0000					},	// 0x81 VK_F18					F18 key
	/* 130 */	{ VC_F19,				0x0000					},	// 0x82 VK_F19					F19 key
	/* 131 */	{ VC_F20,				0x0000					},	// 0x83 VK_F20					F20 key
	/* 132 */	{ VC_F21,				0x0000					},	// 0x84 VK_F21					F21 key
	/* 133 */	{ VC_F22,				0x0000					},	// 0x85 VK_F22					F22 key
	/* 134 */	{ VC_F23,				0x0000					},	// 0x86 VK_F23					F23 key
	/* 135 */	{ VC_F24,				0x0000					},	// 0x87 VK_F24					F24 key
	/* 136 */	{ VC_UNDEFINED,			0x0000					},	// 0x88							Unassigned
	/* 137 */	{ VC_UNDEFINED,			0x0000					},	// 0x89							Unassigned
	/* 138 */	{ VC_UNDEFINED,			0x0000					},	// 0x8A							Unassigned
	/* 139 */	{ VC_UNDEFINED,			0x0000					},	// 0x8B							Unassigned
	/* 140 */	{ VC_UNDEFINED,			0x0000					},	// 0x8C							Unassigned
	/* 141 */	{ VC_UNDEFINED,			0x0000					},	// 0x8D							Unassigned
	/* 142 */	{ VC_UNDEFINED,			0x0000					},	// 0x8E							Unassigned
	/* 143 */	{ VC_UNDEFINED,			0x0000					},	// 0x8F							Unassigned
	/* 144 */	{ VC_NUM_LOCK,			VK_MEDIA_PREV_TRACK		},	// 0x90 VK_NUMLOCK				NUM LOCK key
	/* 145 */	{ VC_SCROLL_LOCK,		0x0000					},	// 0x91 VK_SCROLL				SCROLL LOCK key
	/* 146 */	{ VC_UNDEFINED,			0x0000					},	// 0x92							OEM specific
	/* 147 */	{ VC_UNDEFINED,			0x0000					},	// 0x93							OEM specific
	/* 148 */	{ VC_UNDEFINED,			0x0000					},	// 0x94							OEM specific
	/* 149 */	{ VC_UNDEFINED,			0x0000					},	// 0x95							OEM specific
	/* 150 */	{ VC_UNDEFINED,			0x0000					},	// 0x96							OEM specific
	/* 151 */	{ VC_UNDEFINED,			0x0000					},	// 0x97							Unassigned
	/* 152 */	{ VC_UNDEFINED,			0x0000					},	// 0x98							Unassigned
	/* 153 */	{ VC_UNDEFINED,			VK_MEDIA_NEXT_TRACK		},	// 0x99							Unassigned
	/* 154 */	{ VC_UNDEFINED,			0x0000					},	// 0x9A							Unassigned
	/* 155 */	{ VC_UNDEFINED,			0x0000					},	// 0x9B							Unassigned
	/* 156 */	{ VC_UNDEFINED,			0x0000					},	// 0x9C							Unassigned
	/* 157 */	{ VC_UNDEFINED,			VK_RCONTROL				},	// 0x9D							Unassigned
	/* 158 */	{ VC_UNDEFINED,			0x0000					},	// 0x9E							Unassigned
	/* 159 */	{ VC_UNDEFINED,			0x0000					},	// 0x9F							Unassigned
	/* 160 */	{ VC_SHIFT_L,			VK_VOLUME_MUTE			},	// 0xA0 VK_LSHIFT				Left SHIFT key
	/* 161 */	{ VC_SHIFT_R,			VK_LAUNCH_APP2			},	// 0xA1 VK_RSHIFT				Right SHIFT key
	/* 162 */	{ VC_CONTROL_L,			VK_MEDIA_PLAY_PAUSE		},	// 0xA2 VK_LCONTROL				Left CONTROL key
	/* 163 */	{ VC_CONTROL_R,			0x0000					},	// 0xA3 VK_RCONTROL				Right CONTROL key
	/* 164 */	{ VC_ALT_L,				VK_MEDIA_STOP			},	// 0xA4 VK_LMENU				Left MENU key
	/* 165 */	{ VC_ALT_R,				0x0000					},	// 0xA5 VK_RMENU				Right MENU key
	/* 166 */	{ VC_BROWSER_BACK,		0x0000					},	// 0xA6 VK_BROWSER_BACK			Browser Back key
	/* 167 */	{ VC_BROWSER_FORWARD,	0x0000					},	// 0xA7 VK_BROWSER_FORWARD		Browser Forward key
	/* 168 */	{ VC_BROWSER_REFRESH,	0x0000					},	// 0xA8 VK_BROWSER_REFRESH		Browser Refresh key
	/* 169 */	{ VC_BROWSER_STOP,		0x0000					},	// 0xA9 VK_BROWSER_STOP			Browser Stop key
	/* 170 */	{ VC_BROWSER_SEARCH,	0x0000					},	// 0xAA VK_BROWSER_SEARCH		Browser Search key
	/* 171 */	{ VC_BROWSER_FAVORITES,	0x0000					},	// 0xAB VK_BROWSER_FAVORITES	Browser Favorites key
	/* 172 */	{ VC_BROWSER_HOME,		0x0000					},	// 0xAC VK_BROWSER_HOME			Browser Start and Home key
	/* 173 */	{ VC_VOLUME_MUTE,		0x0000					},	// 0xAD VK_VOLUME_MUTE			Volume Mute key
	/* 174 */	{ VC_VOLUME_DOWN,		VK_VOLUME_DOWN			},	// 0xAE VK_VOLUME_DOWN			Volume Down key
	/* 175 */	{ VC_VOLUME_UP,			0x0000					},	// 0xAF VK_VOLUME_UP			Volume Up key
	/* 176 */	{ VC_MEDIA_NEXT,		VK_VOLUME_UP			},	// 0xB0 VK_MEDIA_NEXT_TRACK		Next Track key
	/* 177 */	{ VC_MEDIA_PREVIOUS,	0x0000					},	// 0xB1 VK_MEDIA_PREV_TRACK		Previous Track key
	/* 178 */	{ VC_MEDIA_STOP,		VK_BROWSER_HOME			},	// 0xB2 VK_MEDIA_STOP			Stop Media key
	/* 179 */	{ VC_MEDIA_PLAY,		0x0000					},	// 0xB3 VK_MEDIA_PLAY_PAUSE		Play/Pause Media key
	/* 180 */	{ VC_UNDEFINED,			0x0000					},	// 0xB4 VK_LAUNCH_MAIL			Start Mail key
	/* 181 */	{ VC_MEDIA_SELECT,		VK_DIVIDE				},	// 0xB5 VK_LAUNCH_MEDIA_SELECT	Select Media key
	/* 182 */	{ VC_APP_MAIL,			0x0000					},	// 0xB6 VK_LAUNCH_APP1			Start Application 1 key
	/* 183 */	{ VC_APP_CALCULATOR,	VK_SNAPSHOT				},	// 0xB7 VK_LAUNCH_APP2			Start Application 2 key
	/* 184 */	{ VC_UNDEFINED,			VK_RMENU				},	// 0xB8							Reserved
	/* 185 */	{ VC_UNDEFINED,			0x0000					},	// 0xB9							Reserved
	/* 186 */	{ VC_SEMICOLON,			0x0000					},	// 0xBA VK_OEM_1				Varies by keyboard. For the US standard keyboard, the ';:' key
	/* 187 */	{ VC_EQUALS,			0x0000					},	// 0xBB VK_OEM_PLUS				For any country/region, the '+' key
	/* 188 */	{ VC_COMMA,				0x00E6					},	// 0xBC VK_OEM_COMMA			For any country/region, the ',' key
	/* 189 */	{ VC_MINUS,				0x0000					},	// 0xBD VK_OEM_MINUS			For any country/region, the '-' key
	/* 190 */	{ VC_PERIOD,			0x0000					},	// 0xBE VK_OEM_PERIOD			For any country/region, the '.' key
	/* 191 */	{ VC_SLASH,				0x0000					},	// 0xBF VK_OEM_2				Varies by keyboard. For the US standard keyboard, the '/?' key
	/* 192 */	{ VC_BACKQUOTE,			0x0000					},	// 0xC0 VK_OEM_3				Varies by keyboard. For the US standard keyboard, the '`~' key
	/* 193 */	{ VC_UNDEFINED,			0x0000					},	// 0xC1							Reserved
	/* 194 */	{ VC_UNDEFINED,			0x0000					},	// 0xC2							Reserved
	/* 195 */	{ VC_UNDEFINED,			0x0000					},	// 0xC3							Reserved
	/* 196 */	{ VC_UNDEFINED,			0x0000					},	// 0xC4							Reserved
	/* 197 */	{ VC_UNDEFINED,			VK_PAUSE				},	// 0xC5							Reserved
	/* 198 */	{ VC_UNDEFINED,			0x0000					},	// 0xC6							Reserved
	/* 199 */	{ VC_UNDEFINED,			VK_HOME					},	// 0xC7							Reserved
	/* 200 */	{ VC_UNDEFINED,			VK_UP					},	// 0xC8							Reserved
	/* 201 */	{ VC_UNDEFINED,			VK_PRIOR				},	// 0xC9							Reserved
	/* 202 */	{ VC_UNDEFINED,			0x0000					},	// 0xCA							Reserved
	/* 203 */	{ VC_UNDEFINED,			VK_LEFT					},	// 0xCB							Reserved
	/* 204 */	{ VC_UNDEFINED,			VK_CLEAR				},	// 0xCC							Reserved
	/* 205 */	{ VC_UNDEFINED,			VK_RIGHT				},	// 0xCD							Reserved
	/* 206 */	{ VC_UNDEFINED,			0x0000					},	// 0xCE							Reserved
	/* 207 */	{ VC_UNDEFINED,			VK_END					},	// 0xCF							Reserved
	/* 208 */	{ VC_UNDEFINED,			VK_DOWN					},	// 0xD0							Reserved
	/* 209 */	{ VC_UNDEFINED,			VK_NEXT					},	// 0xD1							Reserved
	/* 210 */	{ VC_UNDEFINED,			VK_INSERT				},	// 0xD2							Reserved
	/* 211 */	{ VC_UNDEFINED,			VK_DELETE				},	// 0xD3							Reserved
	/* 212 */	{ VC_UNDEFINED,			0x0000					},	// 0xD4							Reserved
	/* 213 */	{ VC_UNDEFINED,			0x0000					},	// 0xD5							Reserved
	/* 214 */	{ VC_UNDEFINED,			0x0000					},	// 0xD6							Reserved
	/* 215 */	{ VC_UNDEFINED,			0x0000					},	// 0xD7							Reserved
	/* 216 */	{ VC_UNDEFINED,			0x0000					},	// 0xD8							Unassigned
	/* 217 */	{ VC_UNDEFINED,			0x0000					},	// 0xD9							Unassigned
	/* 218 */	{ VC_UNDEFINED,			0x0000					},	// 0xDA							Unassigned
	/* 219 */	{ VC_OPEN_BRACKET,		VK_LWIN					},	// 0xDB VK_OEM_4				Varies by keyboard. For the US standard keyboard, the '[{' key
	/* 220 */	{ VC_BACK_SLASH,		VK_RWIN					},	// 0xDC VK_OEM_5				Varies by keyboard. For the US standard keyboard, the '\|' key
	/* 221 */	{ VC_CLOSE_BRACKET,		VK_APPS					},	// 0xDD VK_OEM_6				Varies by keyboard. For the US standard keyboard, the ']}' key
	/* 222 */	{ VC_QUOTE,				0x0000					},	// 0xDE VK_OEM_7				Varies by keyboard. For the US standard keyboard, the 'single-quote/double-quote' key
	/* 223 */	{ VC_YEN,				VK_SLEEP				},	// 0xDF VK_OEM_8				Varies by keyboard.
	/* 224 */	{ VC_UNDEFINED,			0x0000					},	// 0xE0							Reserved
	/* 225 */	{ VC_UNDEFINED,			0x0000					},	// 0xE1							OEM specific
	/* 226 */	{ VC_UNDEFINED,			0x0000					},	// 0xE2 VK_OEM_102				Either the angle bracket key or the backslash key on the RT 102-key keyboard
	/* 227 */	{ VC_UNDEFINED,			0x0000					},	// 0xE3							OEM specific
	/* 228 */	{ VC_UNDEFINED,			0x00E5					},	// 0xE4	VC_APP_PICTURES 		OEM specific
	/* 229 */	{ VC_APP_PICTURES,		VK_BROWSER_SEARCH		},	// 0xE5 VK_PROCESSKEY			IME PROCESS key
	/* 230 */	{ VC_APP_MUSIC,			VK_BROWSER_FAVORITES	},	// 0xE6							OEM specific
	/* 231 */	{ VC_UNDEFINED,			VK_BROWSER_REFRESH		},	// 0xE7 VK_PACKET				Used to pass Unicode characters as if they were keystrokes. The VK_PACKET key is the low word of a 32-bit Virtual Key value used for non-keyboard input methods.
	/* 232 */	{ VC_UNDEFINED,			VK_BROWSER_STOP			},	// 0xE8							Unassigned
	/* 233 */	{ VC_UNDEFINED,			VK_BROWSER_FORWARD		},	// 0xE9							OEM specific
	/* 234 */	{ VC_UNDEFINED,			VK_BROWSER_BACK			},	// 0xEA							OEM specific
	/* 235 */	{ VC_UNDEFINED,			0x0000					},	// 0xEB							OEM specific
	/* 236 */	{ VC_UNDEFINED,			VK_LAUNCH_APP1			},	// 0xEC							OEM specific
	/* 237 */	{ VC_UNDEFINED,			VK_LAUNCH_MEDIA_SELECT	},	// 0xED							OEM specific
	/* 238 */	{ VC_UNDEFINED,			0x0000					},	// 0xEE							OEM specific
	/* 239 */	{ VC_UNDEFINED,			0x0000					},	// 0xEF							OEM specific
	/* 240 */	{ VC_UNDEFINED,			0x0000					},	// 0xF0							OEM specific
	/* 241 */	{ VC_UNDEFINED,			0x0000					},	// 0xF1							OEM specific
	/* 242 */	{ VC_UNDEFINED,			0x0000					},	// 0xF2							OEM specific
	/* 243 */	{ VC_UNDEFINED,			0x0000					},	// 0xF3							OEM specific
	/* 244 */	{ VC_UNDEFINED,			0x0000					},	// 0xF4							OEM specific
	/* 245 */	{ VC_UNDEFINED,			0x0000					},	// 0xF5							OEM specific
	/* 246 */	{ VC_UNDEFINED,			0x0000					},	// 0xF6 VK_ATTN					Attn key
	/* 247 */	{ VC_UNDEFINED,			0x0000					},	// 0xF7 VK_CRSEL				CrSel key
	/* 248 */	{ VC_UNDEFINED,			0x0000					},	// 0xF8 VK_EXSEL				ExSel key
	/* 249 */	{ VC_UNDEFINED,			0x0000					},	// 0xF9 VK_EREOF				Erase EOF key
	/* 250 */	{ VC_UNDEFINED,			0x0000					},	// 0xFA VK_PLAY					Play key
	/* 251 */	{ VC_UNDEFINED,			0x0000					},	// 0xFB VK_ZOOM					Zoom key
	/* 252 */	{ VC_UNDEFINED,			0x0000					},	// 0xFC VK_NONAME				Reserved
	/* 253 */	{ VC_UNDEFINED,			0x0000					},	// 0xFD
	/* 254 */	{ VC_CLEAR,				0x0000					},	// 0xFE VK_OEM_CLEAR			Clear key
	/* 255 */	{ VC_UNDEFINED,			0x0000					}	// 0xFE							Unassigned
};

unsigned short keycode_to_scancode(DWORD vk_code, DWORD flags) {
	unsigned short scancode = VC_UNDEFINED;

	// Check the vk_code is in range.
	// NOTE vk_code >= 0 is assumed because DWORD is unsigned.
	if (vk_code < sizeof(keycode_scancode_table) / sizeof(keycode_scancode_table[0])) {
		scancode = keycode_scancode_table[vk_code][0];

		if (flags & LLKHF_EXTENDED) {
			logger(LOG_LEVEL_WARN,	"%s [%u]: EXTD2, vk_code %li\n",
					__FUNCTION__, __LINE__, vk_code);

			switch (vk_code) {
				case VK_PRIOR:
				case VK_NEXT:
				case VK_END:
				case VK_HOME:
				case VK_LEFT:
				case VK_UP:
				case VK_RIGHT:
				case VK_DOWN:

				case VK_INSERT:
				case VK_DELETE:
					scancode |= 0xEE00;
					break;

				case VK_RETURN:
					scancode |= 0x0E00;
					break;
			}
		}
		else {
						// logger(LOG_LEVEL_WARN,	"%s [%u]: Test2, vk_code %li\n",
            			// 		__FUNCTION__, __LINE__, vk_code);
		}
	}

	return scancode;
}

DWORD scancode_to_keycode(unsigned short scancode) {
	unsigned short keycode = 0x0000;

	// Check the vk_code is in range.
	// NOTE vk_code >= 0 is assumed because the scancode is unsigned.
	if (scancode < 128) {
		keycode = keycode_scancode_table[scancode][1];
	}
	else {
		// Calculate the upper offset based on the lower half of the scancode + 128.
		unsigned short int i = (scancode & 0x007F) | 0x80;

		if (i < sizeof(keycode_scancode_table) / sizeof(keycode_scancode_table[1])) {
			keycode = keycode_scancode_table[i][1];
		}
	}

	return keycode;
}


/************************************************************************/

// Structure and pointers for the keyboard locale cache.
typedef struct _KeyboardLocale {
	HKL id;									// Locale ID
	HINSTANCE library;						// Keyboard DLL instance.
	PVK_TO_BIT pVkToBit;					// Pointers struct arrays.
	PVK_TO_WCHAR_TABLE pVkToWcharTable;
	PDEADKEY pDeadKey;
	struct _KeyboardLocale* next;
} KeyboardLocale;

static KeyboardLocale* locale_first = NULL;
static KeyboardLocale* locale_current = NULL;
static WCHAR deadChar = WCH_NONE;

// Amount of pointer padding to apply for Wow64 instances.
static unsigned short int ptr_padding = 0;

#if defined(_WIN32) && !defined(_WIN64)
// Small function to check and see if we are executing under Wow64.
static BOOL is_wow64() {
	BOOL status = FALSE;

	LPFN_ISWOW64PROCESS pIsWow64Process = (LPFN_ISWOW64PROCESS)
			GetProcAddress(GetModuleHandle("kernel32"), "IsWow64Process");

	if (pIsWow64Process != NULL) {
		HANDLE current_proc = GetCurrentProcess();

		if (!pIsWow64Process(current_proc, &status)) {
			status = FALSE;

			logger(LOG_LEVEL_DEBUG,	"%s [%u]: pIsWow64Process(%#p, %#p) failed!\n",
					__FUNCTION__, __LINE__, current_proc, &status);
		}
	}

	return status;
}
#endif

// Locate the DLL that contains the current keyboard layout.
static int get_keyboard_layout_file(char *layoutFile, DWORD bufferSize) {
	int status = IOHOOK_FAILURE;
	HKEY hKey;
	DWORD varType = REG_SZ;

	char kbdName[KL_NAMELENGTH];
	if (GetKeyboardLayoutName(kbdName)) {
		char kbdKeyPath[51 + KL_NAMELENGTH];
		snprintf(kbdKeyPath, 51 + KL_NAMELENGTH, "SYSTEM\\CurrentControlSet\\Control\\Keyboard Layouts\\%s", kbdName);

		if (RegOpenKeyEx(HKEY_LOCAL_MACHINE, (LPCTSTR) kbdKeyPath, 0, KEY_QUERY_VALUE, &hKey) == ERROR_SUCCESS) {
			if (RegQueryValueEx(hKey, "Layout File", NULL, &varType, (LPBYTE) layoutFile, &bufferSize) == ERROR_SUCCESS) {
				RegCloseKey(hKey);
				status = IOHOOK_SUCCESS;
			}
		}
	}

	return status;
}

static int refresh_locale_list() {
	int count = 0;

	// Get the number of layouts the user has activated.
	int hkl_size = GetKeyboardLayoutList(0, NULL);
	if (hkl_size > 0) {
		logger(LOG_LEVEL_INFO,	"%s [%u]: GetKeyboardLayoutList(0, NULL) found %i layouts.\n",
				__FUNCTION__, __LINE__, hkl_size);

		// Get the thread id that currently has focus for our default.
		DWORD focus_pid = GetWindowThreadProcessId(GetForegroundWindow(), NULL);
		HKL hlk_focus = GetKeyboardLayout(focus_pid);
		HKL hlk_default = GetKeyboardLayout(0);
		HKL *hkl_list = malloc(sizeof(HKL) * hkl_size);

		int new_size = GetKeyboardLayoutList(hkl_size, hkl_list);
		if (new_size > 0) {
			if (new_size != hkl_size) {
				logger(LOG_LEVEL_WARN,	"%s [%u]: Locale size mismatch!  "
						"Expected %i, received %i!\n",
						__FUNCTION__, __LINE__, hkl_size, new_size);
			}
			else {
				logger(LOG_LEVEL_INFO,	"%s [%u]: Received %i locales.\n",
						__FUNCTION__, __LINE__, new_size);
			}

			KeyboardLocale* locale_previous = NULL;
			KeyboardLocale* locale_item = locale_first;

			// Go though the linked list and remove KeyboardLocale's that are
			// no longer loaded.
			while (locale_item != NULL) {
				// Check to see if the old HKL is in the new list.
				bool is_loaded = false;
				int i;
				for (i = 0; i < new_size && !is_loaded; i++) {
					if (locale_item->id == hkl_list[i]) {
						// Flag and jump out of the loop.
						hkl_list[i] = NULL;
						is_loaded = true;
					}
				}


				if (is_loaded) {
					logger(LOG_LEVEL_DEBUG,	"%s [%u]: Found locale ID %#p in the cache.\n",
							__FUNCTION__, __LINE__, locale_item->id);

					// Set the previous local to the current locale.
					locale_previous = locale_item;

					// Check and see if the locale is our current active locale.
					if (locale_item->id == hlk_focus) {
						locale_current = locale_item;
					}

					count++;
				}
				else {
					logger(LOG_LEVEL_DEBUG,	"%s [%u]: Removing locale ID %#p from the cache.\n",
							__FUNCTION__, __LINE__, locale_item->id);

					// If the old id is not in the new list, remove it.
					locale_previous->next = locale_item->next;

					// Make sure the locale_current points NULL or something valid.
					if (locale_item == locale_current) {
						locale_current = NULL;
					}

					// Free the memory used by locale_item;
					free(locale_item);

					// Set the item to the pervious item to guarantee a next.
					locale_item = locale_previous;
				}

				// Iterate to the next linked list item.
				locale_item = locale_item->next;
			}


			// Insert anything new into the linked list.
			int i;
			for (i = 0; i < new_size; i++) {
				// Check to see if the item was already in the list.
				if (hkl_list[i] != NULL) {
					// Set the active keyboard layout for this thread to the HKL.
					ActivateKeyboardLayout(hkl_list[i], 0x00);

					// Try to pull the current keyboard layout DLL from the registry.
					char layoutFile[MAX_PATH];
					if (get_keyboard_layout_file(layoutFile, sizeof(layoutFile)) == IOHOOK_SUCCESS) {
						// You can't trust the %SYSPATH%, look it up manually.
						char systemDirectory[MAX_PATH];
						if (GetSystemDirectory(systemDirectory, MAX_PATH) != 0) {
							char kbdLayoutFilePath[MAX_PATH];
							snprintf(kbdLayoutFilePath, MAX_PATH, "%s\\%s", systemDirectory, layoutFile);

							logger(LOG_LEVEL_DEBUG,	"%s [%u]: Loading layout for %#p: %s.\n",
									__FUNCTION__, __LINE__, hkl_list[i], layoutFile);

							// Create the new locale item.
							locale_item = malloc(sizeof(KeyboardLocale));
							locale_item->id = hkl_list[i];
							locale_item->library = LoadLibrary(kbdLayoutFilePath);

							// Get the function pointer from the library to get the keyboard layer descriptor.
							KbdLayerDescriptor pKbdLayerDescriptor = (KbdLayerDescriptor) GetProcAddress(locale_item->library, "KbdLayerDescriptor");
							if (pKbdLayerDescriptor != NULL) {
								PKBDTABLES pKbd = pKbdLayerDescriptor();

								// Store the memory address of the following 3 structures.
								BYTE *base = (BYTE *) pKbd;

								// First element of each structure, no offset adjustment needed.
								locale_item->pVkToBit = pKbd->pCharModifiers->pVkToBit;

								// Second element of pKbd, +4 byte offset on wow64.
								locale_item->pVkToWcharTable = *((PVK_TO_WCHAR_TABLE *) (base + offsetof(KBDTABLES, pVkToWcharTable) + ptr_padding));

								// Third element of pKbd, +8 byte offset on wow64.
								locale_item->pDeadKey = *((PDEADKEY *) (base + offsetof(KBDTABLES, pDeadKey) + (ptr_padding * 2)));

								// This will always be added to the end of the list.
								locale_item->next = NULL;

								// Insert the item into the linked list.
								if (locale_previous == NULL) {
									// If nothing came before, the list is empty.
									locale_first = locale_item;
								}
								else {
									// Append the new locale to the end of the list.
									locale_previous->next = locale_item;
								}

								// Check and see if the locale is our current active locale.
								if (locale_item->id == hlk_focus) {
									locale_current = locale_item;
								}

								// Set the pervious locale item to the new one.
								locale_previous = locale_item;

								count++;
							}
							else {
								logger(LOG_LEVEL_ERROR,
										"%s [%u]: GetProcAddress() failed for KbdLayerDescriptor!\n",
										__FUNCTION__, __LINE__);

								FreeLibrary(locale_item->library);
								free(locale_item);
								locale_item = NULL;
							}
						}
						else {
							logger(LOG_LEVEL_ERROR,
									"%s [%u]: GetSystemDirectory() failed!\n",
									__FUNCTION__, __LINE__);
						}
					}
					else {
						logger(LOG_LEVEL_ERROR,
								"%s [%u]: Could not find keyboard map for locale %#p!\n",
								__FUNCTION__, __LINE__, hkl_list[i]);
					}
				}
			}
		}
		else {
			logger(LOG_LEVEL_ERROR,
					"%s [%u]: GetKeyboardLayoutList() failed!\n",
					__FUNCTION__, __LINE__);

			// TODO Try and recover by using the current layout.
			// Hint: Use locale_id instead of hkl_list[i] in the loop above.
		}

		free(hkl_list);
		ActivateKeyboardLayout(hlk_default, 0x00);
	}

	return count;
}

SIZE_T keycode_to_unicode(DWORD keycode, PWCHAR buffer, SIZE_T size) {
	// Get the thread id that currently has focus and ask for its current
	// locale.
	DWORD focus_pid = GetWindowThreadProcessId(GetForegroundWindow(), NULL);
	HKL locale_id = GetKeyboardLayout(focus_pid);

	// If the current Locale is not the new locale, search the linked list.
	if (locale_current == NULL || locale_current->id != locale_id) {
		locale_current = NULL;
		KeyboardLocale* locale_item = locale_first;

		// Search the linked list...
		while (locale_item != NULL && locale_item->id != locale_id) {
			locale_item = locale_item->next;
		}

		// You may already be a winner!
		if (locale_item != NULL && locale_item->id != locale_id) {
			logger(LOG_LEVEL_INFO,
					"%s [%u]: Activating keyboard layout %#p.\n",
					__FUNCTION__, __LINE__, locale_item->id);

			// Switch the current locale.
			locale_current = locale_item;
			locale_item = NULL;

			// If they layout changes the dead key state needs to be reset.
			// This is consistent with the way Windows handles locale changes.
			deadChar = WCH_NONE;
		}
		else {
			logger(LOG_LEVEL_DEBUG,
					"%s [%u]: Refreshing locale cache.\n",
					__FUNCTION__, __LINE__);

			refresh_locale_list();
		}
	}

	// Initialize to empty.
	SIZE_T charCount = 0;
	// buffer[i] = WCH_NONE;

	// Check and make sure the Unicode helper was loaded.
	if (locale_current != NULL) {
		logger(LOG_LEVEL_INFO,
				"%s [%u]: Using keyboard layout %#p.\n",
				__FUNCTION__, __LINE__, locale_current->id);

		int mod = 0;

		int capsLock = (GetKeyState(VK_CAPITAL) & 0x01);

		PVK_TO_BIT pVkToBit = locale_current->pVkToBit;
		PVK_TO_WCHAR_TABLE pVkToWcharTable = locale_current->pVkToWcharTable;
		PDEADKEY pDeadKey = locale_current->pDeadKey;

		/* Loop over the modifier keys for this locale and determine what is
		 * currently depressed.  Because this is only a structure of two
		 * bytes, we don't need to worry about the structure padding of __ptr64
		 * offsets on Wow64.
		 */
		bool is_shift = false, is_ctrl = false, is_alt = false;
		int i;
		for (i = 0; pVkToBit[i].Vk != 0; i++) {
			short state = GetAsyncKeyState(pVkToBit[i].Vk);

			// Check to see if the most significant bit is active.
			if (state & ~SHRT_MAX) {
				if (pVkToBit[i].Vk == VK_SHIFT) {
					is_shift = true;
				}
				else if (pVkToBit[i].Vk == VK_CONTROL) {
					is_ctrl = true;
				}
				else if (pVkToBit[i].Vk == VK_MENU) {
					is_alt = true;
				}
			}
		}

		// Check the Shift modifier.
		if (is_shift) {
			mod = 1;
		}

		// Check for the AltGr modifier.
		if (is_ctrl && is_alt) {
			mod += 3;
		}

		// Default 32 bit structure size should be 6 bytes (4 for the pointer and 2
		// additional byte fields) that are padded out to 8 bytes by the compiler.
		unsigned short sizeVkToWcharTable = sizeof(VK_TO_WCHAR_TABLE);
		#if defined(_WIN32) && !defined(_WIN64)
		if (is_wow64()) {
			// If we are running under Wow64 the size of the first pointer will be
			// 8 bringing the total size to 10 bytes padded out to 16.
			sizeVkToWcharTable = (sizeVkToWcharTable + ptr_padding + 7) & -8;
		}
		#endif

		BYTE *ptrCurrentVkToWcharTable = (BYTE *) pVkToWcharTable;

		int cbSize, n;
		do {
			// cbSize is used to calculate n, and n is used for the size of pVkToWchars[j].wch[n]
			cbSize = *(ptrCurrentVkToWcharTable + offsetof(VK_TO_WCHAR_TABLE, cbSize) + ptr_padding);
			n = (cbSize - 2) / 2;

			// Same as VK_TO_WCHARS pVkToWchars[] = pVkToWcharTable[i].pVkToWchars
			PVK_TO_WCHARS pVkToWchars = (PVK_TO_WCHARS) ((PVK_TO_WCHAR_TABLE) ptrCurrentVkToWcharTable)->pVkToWchars;

			if (pVkToWchars != NULL && mod < n) {
				// pVkToWchars[j].VirtualKey
				BYTE *pCurrentVkToWchars = (BYTE *) pVkToWchars;

				do {
					if (((PVK_TO_WCHARS) pCurrentVkToWchars)->VirtualKey == keycode) {
						if ((((PVK_TO_WCHARS) pCurrentVkToWchars)->Attributes == CAPLOK) && capsLock) {
							if (is_shift && mod > 0) {
								mod -= 1;
							}
							else {
								mod += 1;
							}
						}

						// Set the initial unicode char.
						WCHAR unicode = ((PVK_TO_WCHARS) pCurrentVkToWchars)->wch[mod];

						// Increment the pCurrentVkToWchars by the size of wch[n].
						pCurrentVkToWchars += sizeof(VK_TO_WCHARS) + (sizeof(WCHAR) * n);


						if (unicode == WCH_DEAD) {
							// The current unicode char is a dead key...
							if (deadChar == WCH_NONE) {
								// No previous dead key was set so cache the next
								// wchar so we know what to do next time its pressed.
								deadChar = ((PVK_TO_WCHARS) pCurrentVkToWchars)->wch[mod];
							}
							else {
								if (size >= 2) {
									// Received a second dead key.
									memset(buffer, deadChar, 2);
									//buffer[0] = deadChar;
									//buffer[1] = deadChar;

									deadChar = WCH_NONE;
									charCount = 2;
								}
							}
						}
						else if (unicode != WCH_NONE) {
							// We are not WCH_NONE or WCH_DEAD
							if (size >= 1) {
								buffer[0] = unicode;
								charCount = 1;
							}
						}

						break;
					}
					else {
						// Add sizeof WCHAR because we are really an array of WCHAR[n] not WCHAR[]
						pCurrentVkToWchars += sizeof(VK_TO_WCHARS) + (sizeof(WCHAR) * n);
					}
				} while ( ((PVK_TO_WCHARS) pCurrentVkToWchars)->VirtualKey != 0 );
			}

			// This is effectively the same as: ptrCurrentVkToWcharTable = pVkToWcharTable[++i];
			ptrCurrentVkToWcharTable += sizeVkToWcharTable;
		} while (cbSize != 0);


		// If the current local has a dead key set.
		if (deadChar != WCH_NONE) {
			// Loop over the pDeadKey lookup table for the locale.
			int i;
			for (i = 0; pDeadKey[i].dwBoth != 0; i++) {
				WCHAR baseChar = (WCHAR) pDeadKey[i].dwBoth;
				WCHAR diacritic = (WCHAR) (pDeadKey[i].dwBoth >> 16);

				// If we locate an extended dead char, set it.
				if (size >= 1 && baseChar == buffer[0] && diacritic == deadChar) {
					deadChar = WCH_NONE;

					if (charCount <= size) {
						memset(buffer, (WCHAR) pDeadKey[i].wchComposed, charCount);
						//buffer[i] = (WCHAR) pDeadKey[i].wchComposed;
					}
				}
			}
		}
	}

	return charCount;
}

int load_input_helper() {
	int count = 0;

	#if defined(_WIN32) && !defined(_WIN64)
	if (is_wow64()) {
		ptr_padding = sizeof(void *);
	}
	#endif

	count = refresh_locale_list();

	logger(LOG_LEVEL_INFO,
			"%s [%u]: refresh_locale_list() found %i locale(s).\n",
			__FUNCTION__, __LINE__, count);

	return count;
}

// This returns the number of locales that were removed.
int unload_input_helper() {
	int count = 0;

	// Cleanup and free memory from the old list.
	KeyboardLocale* locale_item = locale_first;
	while (locale_item != NULL) {
		// Remove the first item from the linked list.
		FreeLibrary(locale_item->library);
		locale_first = locale_item->next;
		free(locale_item);
		locale_item = locale_first;

		count++;
	}

	// Reset the current local.
	locale_current = NULL;

	return count;
}
