// Copyright 2016 The go-vgo Project Developers. See the COPYRIGHT
// file at the top-level directory of this distribution and at
// https://github.com/go-vgo/robotgo/blob/master/LICENSE
//
// Licensed under the Apache License, Version 2.0 <LICENSE-APACHE or
// http://www.apache.org/licenses/LICENSE-2.0> or the MIT license
// <LICENSE-MIT or http://opensource.org/licenses/MIT>, at your
// option. This file may not be copied, modified, or distributed
// except according to those terms.

#ifndef dispatch_proc_h
#define dispatch_proc_h

// #include "pub.h"
// #include "../chan/eb_chan.h"

void dispatch_proc(iohook_event * const event) {
    if (!sending) { return; }

	// leaking memory? hope not
    char* buffer = calloc(200, sizeof(char));

	switch (event->type) {
	    case EVENT_HOOK_ENABLED:
	    case EVENT_HOOK_DISABLED:
	        sprintf(buffer,
			"{\"id\":%i,\"time\":%" PRIu64 ",\"mask\":%hu,\"reserved\":%hu}",
	        event->type, event->time, event->mask,event->reserved);
	    break;	// send it?
		case EVENT_KEY_PRESSED:
		case EVENT_KEY_RELEASED:
		case EVENT_KEY_TYPED:
           sprintf(buffer,
                "{\"id\":%i,\"time\":%" PRIu64 ",\"mask\":%hu,\"reserved\":%hu,\"keycode\":%hu,\"rawcode\":%hu,\"keychar\":%d}",
                event->type, event->time, event->mask,event->reserved,
                event->data.keyboard.keycode,
                event->data.keyboard.rawcode,
                event->data.keyboard.keychar);
            break;
		case EVENT_MOUSE_PRESSED:
		case EVENT_MOUSE_RELEASED:
		case EVENT_MOUSE_CLICKED:
		case EVENT_MOUSE_MOVED:
		case EVENT_MOUSE_DRAGGED:
			sprintf(buffer,
				"{\"id\":%i,\"time\":%" PRIu64 ",\"mask\":%hu,\"reserved\":%hu,\"x\":%hd,\"y\":%hd,\"button\":%u,\"clicks\":%u}",
				event->type, event->time, event->mask,event->reserved,
				event->data.mouse.x,
				event->data.mouse.y,
				event->data.mouse.button,
				event->data.mouse.clicks);
			break;
		case EVENT_MOUSE_WHEEL:
			sprintf(buffer,
				"{\"id\":%i,\"time\":%" PRIu64 ",\"mask\":%hu,\"reserved\":%hu,\"clicks\":%hu,\"x\":%hd,\"y\":%hd,\"type\":%d,\"ammount\":%hu,\"rotation\":%d,\"direction\":%d}",
				event->type, event->time, event->mask, event->reserved,
				event->data.wheel.clicks,
				event->data.wheel.x,
				event->data.wheel.y,
   				event->data.wheel.type,
   				event->data.wheel.amount,
   				event->data.wheel.rotation,
   				event->data.wheel.direction);
			break;
		default:
		    fprintf(stderr,"\nError on file: %s, unusual event->type: %i\n",__FILE__,event->type);
			return;
	}
	
	// to-do remove this for
	int i;
	for (i = 0; i < 5; i++) {
        switch (eb_chan_try_send(events, buffer)) { // never block the hook callback
            case eb_chan_res_ok:
				i=5;
				break;
            default:
				if (i == 4) { // let's not leak memory
					free(buffer);
				}
				continue;
        }
    }

	// fprintf(stdout, "----%s\n",	 buffer);
}

void dispatch_proc_end(iohook_event * const event) {
	char buffer[256] = { 0 };
	size_t length = snprintf(buffer, sizeof(buffer),
			"id=%i,when=%" PRIu64 ",mask=0x%X",
			event->type, event->time, event->mask);

	switch (event->type) {
		case EVENT_KEY_PRESSED:
			// If the escape key is pressed, naturally terminate the program.
			if (event->data.keyboard.keycode == VC_ESCAPE) {
				// int status = hook_stop();
				// switch (status) {
				// 	// System level errors.
				// 	case IOHOOK_ERROR_OUT_OF_MEMORY:
				// 		loggerProc(LOG_LEVEL_ERROR, "Failed to allocate memory. (%#X)", status);
				// 		break;

				// 	case IOHOOK_ERROR_X_RECORD_GET_CONTEXT:
				// 		// NOTE This is the only platform specific error that occurs on hook_stop().
				// 		loggerProc(LOG_LEVEL_ERROR, "Failed to get XRecord context. (%#X)", status);
				// 		break;

				// 	// Default error.
				// 	case IOHOOK_FAILURE:
				// 	default:
				// 		loggerProc(LOG_LEVEL_ERROR, "An unknown hook error occurred. (%#X)", status);
				// 		break;
				// }
			}
		case EVENT_KEY_RELEASED:
			snprintf(buffer + length, sizeof(buffer) - length,
				",keycode=%u,rawcode=0x%X",
				event->data.keyboard.keycode, event->data.keyboard.rawcode);
				int key_code = (uint16_t) event->data.keyboard.keycode;

				if (event->data.keyboard.keycode == VC_ESCAPE
					&& atoi(cevent) == 11) {
					int stopEvent = stop_event();
					// printf("stop_event%d\n", stopEvent);
					cstatus = 0;
				}

				// printf("atoi(str)---%d\n", atoi(cevent));
				if (key_code == atoi(cevent)) {
					int stopEvent = stop_event();
					// printf("%d\n", stopEvent);
					cstatus = 0;
				}
			break;

		case EVENT_KEY_TYPED:
			snprintf(buffer + length, sizeof(buffer) - length,
				",keychar=%lc,rawcode=%u",
				(uint16_t) event->data.keyboard.keychar,
				event->data.keyboard.rawcode);
				
				#ifdef  WE_REALLY_WANT_A_POINTER
					char *buf = malloc (6);
				#else
					char buf[6];
				#endif

					sprintf(buf, "%lc", (uint16_t) event->data.keyboard.keychar);

				#ifdef WE_REALLY_WANT_A_POINTER
					free (buf);
				#endif

				if (strcmp(buf, cevent) == 0) {
					int stopEvent = stop_event();
					// printf("%d\n", stopEvent);
					cstatus = 0;
				}
				// return (char*) event->data.keyboard.keychar;
			break;

		case EVENT_MOUSE_PRESSED:
		case EVENT_MOUSE_RELEASED:
		case EVENT_MOUSE_CLICKED:
		case EVENT_MOUSE_MOVED:
		case EVENT_MOUSE_DRAGGED:
			snprintf(buffer + length, sizeof(buffer) - length,
				",x=%i,y=%i,button=%i,clicks=%i",
				event->data.mouse.x, event->data.mouse.y,
				event->data.mouse.button, event->data.mouse.clicks);

				int abutton = event->data.mouse.button;
				int aclicks = event->data.mouse.clicks;
				int amouse = -1;

				if (strcmp(cevent, "mleft") == 0) {
					amouse = 1;
				}
				if (strcmp(cevent, "mright") == 0) {
					amouse = 2;
				}
				if (strcmp(cevent, "center") == 0) {
					amouse = 3;
				}
				if (strcmp(cevent, "wheelDown") == 0) {
					amouse = 4;
				}
				if (strcmp(cevent, "wheelUp") == 0) {
					amouse = 5;
				}
				if (strcmp(cevent, "wheelLeft") == 0) {
					amouse = 6;
				}
				if (strcmp(cevent, "wheelRight") == 0) {
					amouse = 7;
				}
				if (abutton == amouse && aclicks == 1) {
					int stopEvent = stop_event();
					cstatus = 0;
				}

			break;

		case EVENT_MOUSE_WHEEL:
			snprintf(buffer + length, sizeof(buffer) - length,
				",type=%i,amount=%i,rotation=%i",
				event->data.wheel.type, event->data.wheel.amount,
				event->data.wheel.rotation);
			break;

		default:
			break;
	}

	// fprintf(stdout, "----%s\n",	 buffer);
}

#endif