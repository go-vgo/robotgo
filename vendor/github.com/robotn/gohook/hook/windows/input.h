/* Copyright (C) 2006-2017 Alexander Barker.
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Lesser General Public License as published
 * by the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
*/

/***********************************************************************
	Input
 ***********************************************************************/

#ifndef _included_input_helper
#define _included_input_helper

#include <limits.h>
#include <windows.h>

#ifndef LPFN_ISWOW64PROCESS
typedef BOOL (WINAPI *LPFN_ISWOW64PROCESS) (HANDLE, PBOOL);
#endif

typedef void* (CALLBACK *KbdLayerDescriptor) (VOID);

#define CAPLOK			0x01
#define WCH_NONE		0xF000
#define WCH_DEAD		0xF001

#ifndef WM_MOUSEHWHEEL
#define WM_MOUSEHWHEEL	0x020E
#endif

typedef struct _VK_TO_WCHARS {
	BYTE VirtualKey;
	BYTE Attributes;
	WCHAR wch[];
} VK_TO_WCHARS, *PVK_TO_WCHARS;

typedef struct _LIGATURE {
	BYTE VirtualKey;
	WORD ModificationNumber;
	WCHAR wch[];
} LIGATURE, *PLIGATURE;

typedef struct _VK_TO_BIT {
	BYTE Vk;
	BYTE ModBits;
} VK_TO_BIT, *PVK_TO_BIT;

typedef struct _MODIFIERS {
	PVK_TO_BIT pVkToBit;				// __ptr64
	WORD wMaxModBits;
	BYTE ModNumber[];
} MODIFIERS, *PMODIFIERS;

typedef struct _VSC_VK {
	BYTE Vsc;
	USHORT Vk;
} VSC_VK, *PVSC_VK;

typedef struct _VK_TO_WCHAR_TABLE {
	PVK_TO_WCHARS pVkToWchars;			// __ptr64
	BYTE nModifications;
	BYTE cbSize;
} VK_TO_WCHAR_TABLE, *PVK_TO_WCHAR_TABLE;

typedef struct _DEADKEY {
	DWORD dwBoth;
	WCHAR wchComposed;
	USHORT uFlags;
} DEADKEY, *PDEADKEY;

typedef struct _VSC_LPWSTR {
	BYTE vsc;
	WCHAR *pwsz;						// __ptr64
} VSC_LPWSTR, *PVSC_LPWSTR;

typedef struct tagKbdLayer {
	PMODIFIERS pCharModifiers;			// __ptr64
	PVK_TO_WCHAR_TABLE pVkToWcharTable;	// __ptr64
	PDEADKEY pDeadKey;					// __ptr64
	PVSC_LPWSTR pKeyNames;				// __ptr64
	PVSC_LPWSTR pKeyNamesExt;			// __ptr64
	WCHAR **pKeyNamesDead;				// __ptr64
	USHORT *pusVSCtoVK;					// __ptr64
	BYTE bMaxVSCtoVK;
	PVSC_VK pVSCtoVK_E0;				// __ptr64
	PVSC_VK pVSCtoVK_E1;				// __ptr64
	DWORD fLocaleFlags;
	BYTE nLgMax;
	BYTE cbLgEntry;
	PLIGATURE pLigature;				// __ptr64
	DWORD dwType;
	DWORD dwSubType;
} KBDTABLES, *PKBDTABLES;				// __ptr64


extern SIZE_T keycode_to_unicode(DWORD keycode, PWCHAR buffer, SIZE_T size);

//extern DWORD unicode_to_keycode(wchar_t unicode);

extern unsigned short keycode_to_scancode(DWORD vk_code, DWORD flags);

extern DWORD scancode_to_keycode(unsigned short scancode);

// Initialize the locale list and wow64 pointer size.
extern int load_input_helper();

// Cleanup the initialized locales.
extern int unload_input_helper();

#endif
