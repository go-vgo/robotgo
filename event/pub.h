// Copyright 2016 The go-vgo Project Developers. See the COPYRIGHT
// file at the top-level directory of this distribution and at
// https://github.com/go-vgo/robotgo/blob/master/LICENSE
//
// Licensed under the Apache License, Version 2.0 <LICENSE-APACHE or
// http://www.apache.org/licenses/LICENSE-2.0> or the MIT license
// <LICENSE-MIT or http://opensource.org/licenses/MIT>, at your
// option. This file may not be copied, modified, or distributed
// except according to those terms.

#include "../base/os.h"

#if defined(IS_MACOSX)
	#include "../cdeps/hook/darwin/input_c.h"
	#include "../cdeps/hook/darwin/hook_c.h"
	#include "../cdeps/hook/darwin/event_c.h"
	#include "../cdeps/hook/darwin/properties_c.h"
#elif defined(USE_X11)
	//#define USE_XKBCOMMON 0
	#include "../cdeps/hook/x11/input_c.h"
	#include "../cdeps/hook/x11/hook_c.h"
	#include "../cdeps/hook/x11/event_c.h"
	#include "../cdeps/hook/x11/properties_c.h"
#elif defined(IS_WINDOWS)
	#include "../cdeps/hook/windows/input_c.h"
	#include "../cdeps/hook/windows/hook_c.h"
	#include "../cdeps/hook/windows/event_c.h"
	#include "../cdeps/hook/windows/properties_c.h"
#endif

#include <inttypes.h>
#include <stdarg.h>
#include <stdbool.h>
#include <stdio.h>
#include <string.h>
#include "../cdeps/hook/iohook.h"


int vccode[100];
int  codesz;

char *cevent;
int rrevent;
// uint16_t *cevent;
int cstatus = 1;


int stop_event();
int add_event(char *key_event);
// int allEvent(char *key_event);
int allEvent(char *key_event, int vcode[], int size);

// NOTE: The following callback executes on the same thread that hook_run() is called
// from.

struct _MEvent {
	uint8_t id;
	size_t mask;
	uint16_t keychar;
	// char *keychar;
	size_t x;
	uint8_t y;
	uint8_t bytesPerPixel;
};

typedef struct _MEvent MEvent;
// typedef MMBitmap *MMBitmapRef;

MEvent mEvent;


bool loggerProc(unsigned int level, const char *format, ...) {
	bool status = false;

	va_list args;
	switch (level) {
		#ifdef USE_DEBUG
		case LOG_LEVEL_DEBUG:
		case LOG_LEVEL_INFO:
			va_start(args, format);
			status = vfprintf(stdout, format, args) >= 0;
			va_end(args);
			break;
		#endif

		case LOG_LEVEL_WARN:
		case LOG_LEVEL_ERROR:
			va_start(args, format);
			status = vfprintf(stderr, format, args) >= 0;
			va_end(args);
			break;
	}

	return status;
}