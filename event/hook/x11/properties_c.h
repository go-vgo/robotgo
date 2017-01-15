
#ifdef HAVE_CONFIG_H
#include <config.h>
#endif

#include <stdbool.h>
#include <stdio.h>
#include <stdlib.h>
#include <X11/Xlib.h>
#ifdef USE_XKB
	#include <X11/XKBlib.h>
#endif
#ifdef USE_XF86MISC
	#include <X11/extensions/xf86misc.h>
	#include <X11/extensions/xf86mscstr.h>
#endif
#if defined(USE_XINERAMA) && !defined(USE_XRANDR)
	#include <X11/extensions/Xinerama.h>
#elif defined(USE_XRANDR)
#include <pthread.h>
	#include <X11/extensions/Xrandr.h>
#endif
#ifdef USE_XT
	#include <X11/Intrinsic.h>

	static XtAppContext xt_context;
	static Display *xt_disp;
#endif

#include "../iohook.h"
#include "input.h"
// #include "../logger.h"

Display *properties_disp;

#ifdef USE_XRANDR
static pthread_mutex_t xrandr_mutex = PTHREAD_MUTEX_INITIALIZER;
static XRRScreenResources *xrandr_resources = NULL;

static void settings_cleanup_proc(void *arg) {
	if (pthread_mutex_trylock(&xrandr_mutex) == 0) {
		if (xrandr_resources != NULL) {
			XRRFreeScreenResources(xrandr_resources);
			xrandr_resources = NULL;
		}

		if (arg != NULL) {
			XCloseDisplay((Display *) arg);
			arg = NULL;
		}

		pthread_mutex_unlock(&xrandr_mutex);
	}
}

static void *settings_thread_proc(void *arg) {
	Display *settings_disp = XOpenDisplay(XDisplayName(NULL));;
	if (settings_disp != NULL) {
		logger(LOG_LEVEL_DEBUG,	"%s [%u]: %s\n",
				__FUNCTION__, __LINE__, "XOpenDisplay success.");

		pthread_cleanup_push(settings_cleanup_proc, settings_disp);

		int event_base = 0;
		int error_base = 0;
		if (XRRQueryExtension(settings_disp, &event_base, &error_base)) {
			Window root = XDefaultRootWindow(settings_disp);
			unsigned long event_mask = RRScreenChangeNotifyMask;
			XRRSelectInput(settings_disp, root, event_mask);

			XEvent ev;

			while(settings_disp != NULL) {
				XNextEvent(settings_disp, &ev);

				if (ev.type == event_base + RRScreenChangeNotifyMask) {
					logger(LOG_LEVEL_DEBUG,	"%s [%u]: Received XRRScreenChangeNotifyEvent.\n",
							__FUNCTION__, __LINE__);

					pthread_mutex_lock(&xrandr_mutex);
					if (xrandr_resources != NULL) {
						XRRFreeScreenResources(xrandr_resources);
					}

					xrandr_resources = XRRGetScreenResources(settings_disp, root);
					if (xrandr_resources == NULL) {
						logger(LOG_LEVEL_WARN,	"%s [%u]: XRandR could not get screen resources!\n",
								__FUNCTION__, __LINE__);
					}
					pthread_mutex_unlock(&xrandr_mutex);
				}
				else {
					logger(LOG_LEVEL_WARN,	"%s [%u]: XRandR is not currently available!\n",
							__FUNCTION__, __LINE__);
				}
			}
		}

		// Execute the thread cleanup handler.
		pthread_cleanup_pop(1);

	}
	else {
		logger(LOG_LEVEL_ERROR,	"%s [%u]: XOpenDisplay failure!\n",
				__FUNCTION__, __LINE__);
	}

	return NULL;
}
#endif

IOHOOK_API screen_data* hook_create_screen_info(unsigned char *count) {
	*count = 0;
	screen_data *screens = NULL;

	#if defined(USE_XINERAMA) && !defined(USE_XRANDR)
	if (XineramaIsActive(properties_disp)) {
		int xine_count = 0;
		XineramaScreenInfo *xine_info = XineramaQueryScreens(properties_disp, &xine_count);

		if (xine_info != NULL) {
			if (xine_count > UINT8_MAX) {
				*count = UINT8_MAX;

				logger(LOG_LEVEL_WARN, "%s [%u]: Screen count overflow detected!\n",
						__FUNCTION__, __LINE__);
			}
			else {
				*count = (uint8_t) xine_count;
			}

			screens = malloc(sizeof(screen_data) * xine_count);

			if (screens != NULL) {
				int i;
				for (i = 0; i < xine_count; i++) {
					screens[i] = (screen_data) {
						.number = xine_info[i].screen_number,
						.x = xine_info[i].x_org,
						.y = xine_info[i].y_org,
						.width = xine_info[i].width,
						.height = xine_info[i].height
					};
				}
			}

			XFree(xine_info);
		}
	}
	#elif defined(USE_XRANDR)
	pthread_mutex_lock(&xrandr_mutex);
	if (xrandr_resources != NULL) {
		int xrandr_count = xrandr_resources->ncrtc;
		if (xrandr_count > UINT8_MAX) {
			*count = UINT8_MAX;

			logger(LOG_LEVEL_WARN, "%s [%u]: Screen count overflow detected!\n",
					__FUNCTION__, __LINE__);
		}
		else {
			*count = (uint8_t) xrandr_count;
		}

		screens = malloc(sizeof(screen_data) * xrandr_count);

		if (screens != NULL) {
			int i;
			for (i = 0; i < xrandr_count; i++) {
				XRRCrtcInfo *crtc_info = XRRGetCrtcInfo(properties_disp, xrandr_resources, xrandr_resources->crtcs[i]);

				if (crtc_info != NULL) {
					screens[i] = (screen_data) {
						.number = i + 1,
						.x = crtc_info->x,
						.y = crtc_info->y,
						.width = crtc_info->width,
						.height = crtc_info->height
					};

					XRRFreeCrtcInfo(crtc_info);
				}
				else {
					logger(LOG_LEVEL_WARN,	"%s [%u]: XRandr failed to return crtc information! (%#X)\n",
							__FUNCTION__, __LINE__, xrandr_resources->crtcs[i]);
				}
			}
		}
	}
	pthread_mutex_unlock(&xrandr_mutex);
	#else
	Screen* default_screen = DefaultScreenOfDisplay(properties_disp);

	if (default_screen->width > 0 && default_screen->height > 0) {
		screens = malloc(sizeof(screen_data));

		if (screens != NULL) {
			*count = 1;
			screens[0] = (screen_data) {
				.number = 1,
				.x = 0,
				.y = 0,
				.width = default_screen->width,
				.height = default_screen->height
			};
		}
	}
	#endif

	return screens;
}

IOHOOK_API long int hook_get_auto_repeat_rate() {
	bool successful = false;
	long int value = -1;
	unsigned int delay = 0, rate = 0;

	// Check and make sure we could connect to the x server.
	if (properties_disp != NULL) {
		#ifdef USE_XKB
		// Attempt to acquire the keyboard auto repeat rate using the XKB extension.
		if (!successful) {
			successful = XkbGetAutoRepeatRate(properties_disp, XkbUseCoreKbd, &delay, &rate);

			if (successful) {
				logger(LOG_LEVEL_INFO,	"%s [%u]: XkbGetAutoRepeatRate: %u.\n",
						__FUNCTION__, __LINE__, rate);
			}
		}
		#endif

		#ifdef USE_XF86MISC
		// Fallback to the XF86 Misc extension if available and other efforts failed.
		if (!successful) {
			XF86MiscKbdSettings kb_info;
			successful = (bool) XF86MiscGetKbdSettings(properties_disp, &kb_info);
			if (successful) {
				logger(LOG_LEVEL_INFO,	"%s [%u]: XF86MiscGetKbdSettings: %i.\n",
						__FUNCTION__, __LINE__, kbdinfo.rate);

				delay = (unsigned int) kbdinfo.delay;
				rate = (unsigned int) kbdinfo.rate;
			}
		}
		#endif
	}
	else {
		logger(LOG_LEVEL_ERROR,	"%s [%u]: %s\n",
				__FUNCTION__, __LINE__, "XOpenDisplay failure!");
	}

	if (successful) {
		value = (long int) rate;
	}

	return value;
}

IOHOOK_API long int hook_get_auto_repeat_delay() {
	bool successful = false;
	long int value = -1;
	unsigned int delay = 0, rate = 0;

	// Check and make sure we could connect to the x server.
	if (properties_disp != NULL) {
		#ifdef USE_XKB
		// Attempt to acquire the keyboard auto repeat rate using the XKB extension.
		if (!successful) {
			successful = XkbGetAutoRepeatRate(properties_disp, XkbUseCoreKbd, &delay, &rate);

			if (successful) {
				logger(LOG_LEVEL_INFO,	"%s [%u]: XkbGetAutoRepeatRate: %u.\n",
						__FUNCTION__, __LINE__, delay);
			}
		}
		#endif

		#ifdef USE_XF86MISC
		// Fallback to the XF86 Misc extension if available and other efforts failed.
		if (!successful) {
			XF86MiscKbdSettings kb_info;
			successful = (bool) XF86MiscGetKbdSettings(properties_disp, &kb_info);
			if (successful) {
				logger(LOG_LEVEL_INFO,	"%s [%u]: XF86MiscGetKbdSettings: %i.\n",
						__FUNCTION__, __LINE__, kbdinfo.delay);

				delay = (unsigned int) kbdinfo.delay;
				rate = (unsigned int) kbdinfo.rate;
			}
		}
		#endif
	}
	else {
		logger(LOG_LEVEL_ERROR,	"%s [%u]: %s\n",
				__FUNCTION__, __LINE__, "XOpenDisplay failure!");
	}

	if (successful) {
		value = (long int) delay;
	}

	return value;
}

IOHOOK_API long int hook_get_pointer_acceleration_multiplier() {
	long int value = -1;
	int accel_numerator, accel_denominator, threshold;

	// Check and make sure we could connect to the x server.
	if (properties_disp != NULL) {
		XGetPointerControl(properties_disp, &accel_numerator, &accel_denominator, &threshold);
		if (accel_denominator >= 0) {
			logger(LOG_LEVEL_INFO,	"%s [%u]: XGetPointerControl: %i.\n",
					__FUNCTION__, __LINE__, accel_denominator);

			value = (long int) accel_denominator;
		}
	}
	else {
		logger(LOG_LEVEL_ERROR,	"%s [%u]: %s\n",
				__FUNCTION__, __LINE__, "XOpenDisplay failure!");
	}

	return value;
}

IOHOOK_API long int hook_get_pointer_acceleration_threshold() {
	long int value = -1;
	int accel_numerator, accel_denominator, threshold;

	// Check and make sure we could connect to the x server.
	if (properties_disp != NULL) {
		XGetPointerControl(properties_disp, &accel_numerator, &accel_denominator, &threshold);
		if (threshold >= 0) {
			logger(LOG_LEVEL_INFO,	"%s [%u]: XGetPointerControl: %i.\n",
					__FUNCTION__, __LINE__, threshold);

			value = (long int) threshold;
		}
	}
	else {
		logger(LOG_LEVEL_ERROR,	"%s [%u]: %s\n",
				__FUNCTION__, __LINE__, "XOpenDisplay failure!");
	}

	return value;
}

IOHOOK_API long int hook_get_pointer_sensitivity() {
	long int value = -1;
	int accel_numerator, accel_denominator, threshold;

	// Check and make sure we could connect to the x server.
	if (properties_disp != NULL) {
		XGetPointerControl(properties_disp, &accel_numerator, &accel_denominator, &threshold);
		if (accel_numerator >= 0) {
			logger(LOG_LEVEL_INFO,	"%s [%u]: XGetPointerControl: %i.\n",
					__FUNCTION__, __LINE__, accel_numerator);

			value = (long int) accel_numerator;
		}
	}
	else {
		logger(LOG_LEVEL_ERROR,	"%s [%u]: %s\n",
				__FUNCTION__, __LINE__, "XOpenDisplay failure!");
	}

	return value;
}

IOHOOK_API long int hook_get_multi_click_time() {
	long int value = 200;
	int click_time;
	bool successful = false;

	#ifdef USE_XT
	// Check and make sure we could connect to the x server.
	if (xt_disp != NULL) {
		// Try and use the Xt extention to get the current multi-click.
		if (!successful) {
			// Fall back to the X Toolkit extension if available and other efforts failed.
			click_time = XtGetMultiClickTime(xt_disp);
			if (click_time >= 0) {
				logger(LOG_LEVEL_INFO,	"%s [%u]: XtGetMultiClickTime: %i.\n",
						__FUNCTION__, __LINE__, click_time);

				successful = true;
			}
		}
	}
	else {
		logger(LOG_LEVEL_ERROR,	"%s [%u]: %s\n",
				__FUNCTION__, __LINE__, "XOpenDisplay failure!");
	}
	#endif

	// Check and make sure we could connect to the x server.
	if (properties_disp != NULL) {
		// Try and acquire the multi-click time from the user defined X defaults.
		if (!successful) {
			char *xprop = XGetDefault(properties_disp, "*", "multiClickTime");
			if (xprop != NULL && sscanf(xprop, "%4i", &click_time) != EOF) {
				logger(LOG_LEVEL_INFO,	"%s [%u]: X default 'multiClickTime' property: %i.\n",
						__FUNCTION__, __LINE__, click_time);

				successful = true;
			}
		}

		if (!successful) {
			char *xprop = XGetDefault(properties_disp, "OpenWindows", "MultiClickTimeout");
			if (xprop != NULL && sscanf(xprop, "%4i", &click_time) != EOF) {
				logger(LOG_LEVEL_INFO,	"%s [%u]: X default 'MultiClickTimeout' property: %i.\n",
						__FUNCTION__, __LINE__, click_time);

				successful = true;
			}
		}
	}
	else {
		logger(LOG_LEVEL_ERROR,	"%s [%u]: %s\n",
				__FUNCTION__, __LINE__, "XOpenDisplay failure!");
	}

	if (successful) {
		value = (long int) click_time;
	}

	return value;
}

// Create a shared object constructor.
__attribute__ ((constructor))
void on_library_load() {
	// Make sure we are initialized for threading.
	XInitThreads();

	// Open local display.
	properties_disp = XOpenDisplay(XDisplayName(NULL));
	if (properties_disp == NULL) {
		logger(LOG_LEVEL_ERROR,	"%s [%u]: %s\n",
				__FUNCTION__, __LINE__, "XOpenDisplay failure!");
	}
	else {
		logger(LOG_LEVEL_DEBUG,	"%s [%u]: %s\n",
				__FUNCTION__, __LINE__, "XOpenDisplay success.");
	}

	#ifdef USE_XRANDR
	// Create the thread attribute.
	pthread_attr_t settings_thread_attr;
	pthread_attr_init(&settings_thread_attr);

	pthread_t settings_thread_id;
	if (pthread_create(&settings_thread_id, &settings_thread_attr, settings_thread_proc, NULL) == 0) {
		logger(LOG_LEVEL_DEBUG,	"%s [%u]: Successfully created settings thread.\n",
				__FUNCTION__, __LINE__);
	}
	else {
		logger(LOG_LEVEL_ERROR,	"%s [%u]: Failed to create settings thread!\n",
				__FUNCTION__, __LINE__);
	}

	// Make sure the thread attribute is removed.
	pthread_attr_destroy(&settings_thread_attr);
	#endif

	#ifdef USE_XT
	XtToolkitInitialize();
	xt_context = XtCreateApplicationContext();

	int argc = 0;
	char ** argv = { NULL };
	xt_disp = XtOpenDisplay(xt_context, NULL, "IOHook", "libIOhook", NULL, 0, &argc, argv);
	#endif

	// Initialize.
	load_input_helper(properties_disp);
}

// Create a shared object destructor.
__attribute__ ((destructor))
void on_library_unload() {
	// Disable the event hook.
	//hook_stop();

	// Cleanup.
	unload_input_helper();

	#ifdef USE_XT
	XtCloseDisplay(xt_disp);
	XtDestroyApplicationContext(xt_context);
	#endif

	// Destroy the native displays.
	if (properties_disp != NULL) {
		XCloseDisplay(properties_disp);
		properties_disp = NULL;
	}
}
