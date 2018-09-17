
#ifdef HAVE_CONFIG_H
	#include <config.h>
#endif

#include <windows.h>
#include "../iohook.h"
#include "input.h"
// #include "logger.h"

// The handle to the DLL module pulled in DllMain on DLL_PROCESS_ATTACH.
HINSTANCE hInst;

// input_hook.c
extern void unregister_running_hooks();


// Structure for the monitor_enum_proc() callback so we can track the count.
typedef struct _screen_info {
	uint8_t count;
	screen_data *data;
} screen_info;


static BOOL CALLBACK monitor_enum_proc(HMONITOR hMonitor, HDC hdcMonitor, LPRECT lprcMonitor, LPARAM dwData) {
 	int width  = lprcMonitor->right - lprcMonitor->left;
	int height = lprcMonitor->bottom - lprcMonitor->top;
	int origin_x = lprcMonitor->left;
	int origin_y = lprcMonitor->top;

	if (width > 0 && height > 0) {
		screen_info *screens = (screen_info *) dwData;

		if (screens->data == NULL) {
			screens->data = (screen_data *) malloc(sizeof(screen_data));
		}
		else {
			screens->data = (screen_data *) realloc(screens, sizeof(screen_data) * screens->count);
		}

		screens->data[screens->count++] = (screen_data) {
				// Should monitor count start @ zero? Currently it starts at 1.
				.number = screens->count,
				.x = origin_x,
				.y = origin_y,
				.width = width,
				.height = height
			};

			logger(LOG_LEVEL_INFO,	"%s [%u]: Monitor %d: %ldx%ld (%ld, %ld)\n",
					__FUNCTION__, __LINE__, screens->count, width, height, origin_x, origin_y);
	}

	return TRUE;
}

IOHOOK_API screen_data* hook_create_screen_info(unsigned char *count) {
	// Initialize count to zero.
	*count = 0;

	// Create a simple structure to make working with monitor_enum_proc easier.
	screen_info screens = {
		.count = 0,
		.data = NULL
	};

	BOOL status = EnumDisplayMonitors(NULL, NULL, monitor_enum_proc, (LPARAM) &screens);

	if (!status || screens.count == 0) {
		// Fallback in case EnumDisplayMonitors fails.
		logger(LOG_LEVEL_INFO,	"%s [%u]: EnumDisplayMonitors failed. Fallback.\n",
				__FUNCTION__, __LINE__);

		int width  = GetSystemMetrics(SM_CXSCREEN);
		int height = GetSystemMetrics(SM_CYSCREEN);

		if (width > 0 && height > 0) {
			screens.data = (screen_data *) malloc(sizeof(screen_data));

			if (screens.data != NULL) {
				*count = 1;
				screens.data[0] = (screen_data) {
					.number = 1,
					.x = 0,
					.y = 0,
					.width = width,
					.height = height
				};
			}
		}
	} else {
		// Populate the count.
		*count = screens.count;
	}

	return screens.data;
}

IOHOOK_API long int hook_get_auto_repeat_rate() {
	long int value = -1;
	long int rate;

	if (SystemParametersInfo(SPI_GETKEYBOARDSPEED, 0, &rate, 0)) {
		logger(LOG_LEVEL_INFO,	"%s [%u]: SPI_GETKEYBOARDSPEED: %li.\n",
			__FUNCTION__, __LINE__, rate);

		value = rate;
	}

	return value;
}

IOHOOK_API long int hook_get_auto_repeat_delay() {
	long int value = -1;
	long int delay;

	if (SystemParametersInfo(SPI_GETKEYBOARDDELAY, 0, &delay, 0)) {
		logger(LOG_LEVEL_INFO,	"%s [%u]: SPI_GETKEYBOARDDELAY: %li.\n",
			__FUNCTION__, __LINE__, delay);

		value = delay;
	}

	return value;
}

IOHOOK_API long int hook_get_pointer_acceleration_multiplier() {
	long int value = -1;
	int mouse[3]; // 0-Threshold X, 1-Threshold Y and 2-Speed.

	if (SystemParametersInfo(SPI_GETMOUSE, 0, &mouse, 0)) {
		logger(LOG_LEVEL_INFO,	"%s [%u]: SPI_GETMOUSE[2]: %i.\n",
			__FUNCTION__, __LINE__, mouse[2]);

		value = mouse[2];
	}

	return value;
}

IOHOOK_API long int hook_get_pointer_acceleration_threshold() {
	long int value = -1;
	int mouse[3]; // 0-Threshold X, 1-Threshold Y and 2-Speed.

	if (SystemParametersInfo(SPI_GETMOUSE, 0, &mouse, 0)) {
		logger(LOG_LEVEL_INFO,	"%s [%u]: SPI_GETMOUSE[0]: %i.\n",
			__FUNCTION__, __LINE__, mouse[0]);
		logger(LOG_LEVEL_INFO,	"%s [%u]: SPI_GETMOUSE[1]: %i.\n",
			__FUNCTION__, __LINE__, mouse[1]);

		// Average the x and y thresholds.
		value = (mouse[0] + mouse[1]) / 2;
	}

	return value;
}

IOHOOK_API long int hook_get_pointer_sensitivity() {
	long int value = -1;
	int sensitivity;

	if (SystemParametersInfo(SPI_GETMOUSESPEED, 0, &sensitivity, 0)) {
		logger(LOG_LEVEL_INFO,	"%s [%u]: SPI_GETMOUSESPEED: %i.\n",
			__FUNCTION__, __LINE__, sensitivity);

		value = sensitivity;
	}

	return value;
}

IOHOOK_API long int hook_get_multi_click_time() {
	long int value = -1;
	UINT clicktime;

	clicktime = GetDoubleClickTime();
	logger(LOG_LEVEL_INFO,	"%s [%u]: GetDoubleClickTime: %u.\n",
			__FUNCTION__, __LINE__, (unsigned int) clicktime);

	value = (long int) clicktime;

	return value;
}

// DLL Entry point.
BOOL WINAPI DllMain(HINSTANCE hInstDLL, DWORD fdwReason, LPVOID lpReserved) {
	switch (fdwReason) {
		case DLL_PROCESS_ATTACH:
			// Save the DLL address.
			hInst = hInstDLL;

			// Initialize native input helper functions.
			load_input_helper();
			break;

		case DLL_PROCESS_DETACH:
			// Unregister any hooks that may still be installed.
			unregister_running_hooks();

			// Deinitialize native input helper functions.
			unload_input_helper();
			break;

		case DLL_THREAD_ATTACH:
		case DLL_THREAD_DETACH:
			// Do Nothing.
			break;
	}

	return TRUE;
}
