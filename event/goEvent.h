// Copyright 2016 The go-vgo Project Developers. See the COPYRIGHT
// file at the top-level directory of this distribution and at
// https://github.com/go-vgo/robotgo/blob/master/LICENSE
//
// Licensed under the Apache License, Version 2.0 <LICENSE-APACHE or
// http://www.apache.org/licenses/LICENSE-2.0> or the MIT license
// <LICENSE-MIT or http://opensource.org/licenses/MIT>, at your
// option. This file may not be copied, modified, or distributed
// except according to those terms.
//
// The hook directory link from the https://github.com/robotn/gohook/
// hook, you need to follow the relevant agreement and LICENSE.
// See the LICENSE file at the top-level directory of this distribution and at
// https://github.com/robotn/gohook/blob/master/LICENSE

#ifdef HAVE_CONFIG_H
	#include <config.h>
#endif

#include "pub.h"

void dispatch_proc(iohook_event * const event) {

    int keycode;

	switch (event->type) {
		case EVENT_KEY_PRESSED:
			// If the escape key is pressed, naturally terminate the program.
			if (event->data.keyboard.keycode == VC_ESCAPE) {

			}
		case EVENT_KEY_RELEASED:

            keycode = (int) event->data.keyboard.keycode;
            //printf("EVENT_KEY_RELEASED:%d\n", keycode);
            showKeyCode(keycode);
            //callback(keycode);

			break;

		case EVENT_KEY_TYPED:

            keycode = (int) event->data.keyboard.keycode;
            //printf("EVENT_KEY_TYPED:%d\n", keycode);

			break;

		default:
			break;
	}
}

int add_event_listener() {

    printf("start C add_event\n");


	// Set the logger callback for library output.
	printf("start C hookSetlogger\n");

	hookSetlogger(&loggerProc);

    printf("start C hook_set_dispatch_proc\n");
	// Set the event callback for IOhook events.
	hook_set_dispatch_proc(&dispatch_proc);

	printf("start C hook_run\n");
	// Start the hook and block.
	// NOTE If EVENT_HOOK_ENABLED was delivered, the status will always succeed.
	int status = hook_run();

    printf("hook_run status:%d\n", status);

	switch (status) {
		case IOHOOK_SUCCESS:
			// Everything is ok.
			break;

		// System level errors.
		case IOHOOK_ERROR_OUT_OF_MEMORY:
			loggerProc(LOG_LEVEL_ERROR, "Failed to allocate memory. (%#X)", status);
			break;


		// X11 specific errors.
		case IOHOOK_ERROR_X_OPEN_DISPLAY:
			loggerProc(LOG_LEVEL_ERROR, "Failed to open X11 display. (%#X)", status);
			break;

		case IOHOOK_ERROR_X_RECORD_NOT_FOUND:
			loggerProc(LOG_LEVEL_ERROR, "Unable to locate XRecord extension. (%#X)", status);
			break;

		case IOHOOK_ERROR_X_RECORD_ALLOC_RANGE:
			loggerProc(LOG_LEVEL_ERROR, "Unable to allocate XRecord range. (%#X)", status);
			break;

		case IOHOOK_ERROR_X_RECORD_CREATE_CONTEXT:
			loggerProc(LOG_LEVEL_ERROR, "Unable to allocate XRecord context. (%#X)", status);
			break;

		case IOHOOK_ERROR_X_RECORD_ENABLE_CONTEXT:
			loggerProc(LOG_LEVEL_ERROR, "Failed to enable XRecord context. (%#X)", status);
			break;


		// Windows specific errors.
		case IOHOOK_ERROR_SET_WINDOWS_HOOK_EX:
			loggerProc(LOG_LEVEL_ERROR, "Failed to register low level windows hook. (%#X)", status);
			break;


		// Darwin specific errors.
		case IOHOOK_ERROR_AXAPI_DISABLED:
			loggerProc(LOG_LEVEL_ERROR, "Failed to enable access for assistive devices. (%#X)", status);
			break;

		case IOHOOK_ERROR_CREATE_EVENT_PORT:
			loggerProc(LOG_LEVEL_ERROR, "Failed to create apple event port. (%#X)", status);
			break;

		case IOHOOK_ERROR_CREATE_RUN_LOOP_SOURCE:
			loggerProc(LOG_LEVEL_ERROR, "Failed to create apple run loop source. (%#X)", status);
			break;

		case IOHOOK_ERROR_GET_RUNLOOP:
			loggerProc(LOG_LEVEL_ERROR, "Failed to acquire apple run loop. (%#X)", status);
			break;

		case IOHOOK_ERROR_CREATE_OBSERVER:
			loggerProc(LOG_LEVEL_ERROR, "Failed to create apple run loop observer. (%#X)", status);
			break;

		// Default error.
		case IOHOOK_FAILURE:
		default:
			loggerProc(LOG_LEVEL_ERROR, "An unknown hook error occurred. (%#X)", status);
			break;
	}

	return cstatus;
}

int stop_event(){
	int status = hook_stop();
	switch (status) {
		// System level errors.
		case IOHOOK_ERROR_OUT_OF_MEMORY:
			loggerProc(LOG_LEVEL_ERROR, "Failed to allocate memory. (%#X)", status);
			break;

		case IOHOOK_ERROR_X_RECORD_GET_CONTEXT:
			// NOTE This is the only platform specific error that occurs on hook_stop().
			loggerProc(LOG_LEVEL_ERROR, "Failed to get XRecord context. (%#X)", status);
			break;

				// Default error.
		case IOHOOK_FAILURE:
			default:
			// loggerProc(LOG_LEVEL_ERROR, "An unknown hook error occurred. (%#X)", status);
			break;
	}

	return status;
}