
#ifdef HAVE_CONFIG_H
	#include <config.h>
#endif

#ifdef USE_CARBON_LEGACY
	#include <Carbon/Carbon.h>
#endif
#ifdef USE_COREFOUNDATION
	#include <CoreFoundation/CoreFoundation.h>
#endif
#ifdef USE_IOKIT
	#include <IOKit/hidsystem/IOHIDLib.h>
	#include <IOKit/hidsystem/IOHIDParameter.h>
#endif

#include <stdbool.h>
#include "../iohook.h"
#include "input.h"
// #include "../logger_c.h"

IOHOOK_API screen_data* hook_create_screen_info(unsigned char *count) {
	CGError status = kCGErrorFailure;
	screen_data* screens = NULL;

	// Initialize count to zero.
	*count = 0;

	// Allocate memory to hold each display id.  We will just allocate our MAX
	// because its only about 1K of memory.
	// TODO This can probably be realistically cut to something like 16 or 32....
	// If you have more than 32 monitors, send me a picture and make a donation ;)
	CGDirectDisplayID *display_ids = malloc(sizeof(CGDirectDisplayID) * UCHAR_MAX);
	if (display_ids != NULL) {
		// NOTE Pass UCHAR_MAX to make sure uint32_t doesn't overflow uint8_t.
		// TOOD Test/Check whether CGGetOnlineDisplayList is more suitable...
		status = CGGetActiveDisplayList(UCHAR_MAX, display_ids, (uint32_t *) count);

		// If there is no error and at least one monitor.
		if (status == kCGErrorSuccess && *count > 0) {
			logger(LOG_LEVEL_INFO,	"%s [%u]: CGGetActiveDisplayList: %li.\n",
					__FUNCTION__, __LINE__, *count);

			// Allocate memory for the number of screens found.
			screens = malloc(sizeof(screen_data) * (*count));
			if (screens != NULL) {
				uint8_t i;
				for (i = 0; i < *count; i++) {
					//size_t width = CGDisplayPixelsWide(display_ids[i]);
					//size_t height = CGDisplayPixelsHigh(display_ids[i]);
					CGRect boundsDisp = CGDisplayBounds(display_ids[i]);
					if (boundsDisp.size.width > 0 && boundsDisp.size.height > 0) {
						screens[i] = (screen_data) {
							.number = i + 1,
							//TODO: make sure we follow the same convention for the origin
							//in all other platform implementations (upper-left)
							//TODO: document the approach with examples in order to show different
							//cases -> different resolutions (secondary monitors origin might be
							//negative)
							.x = boundsDisp.origin.x,
							.y = boundsDisp.origin.y,
							.width = boundsDisp.size.width,
							.height = boundsDisp.size.height
						};
				}
				}
			}
		}
		else {
			logger(LOG_LEVEL_INFO,	"%s [%u]: multiple_get_screen_info failed: %ld. Fallback.\n",
					__FUNCTION__, __LINE__, status);

			size_t width = CGDisplayPixelsWide(CGMainDisplayID());
			size_t height = CGDisplayPixelsHigh(CGMainDisplayID());

			if (width > 0 && height > 0) {
				screens = malloc(sizeof(screen_data));

				if (screens != NULL) {
					*count = 1;
					screens[0] = (screen_data) {
						.number = 1,
						.x = 0,
						.y = 0,
						.width = width,
						.height = height
					};
				}
			}
		}

		// Free the id's after we are done.
		free(display_ids);
	}

	return screens;
}

/*
 * Apple's documentation is not very good.  I was finally able to find this
 * information after many hours of googling.  Value is the slider value in the
 * system preferences. That value * 15 is the rate in MS.  66 / the value is the
 * chars per second rate.
 *
 * Value	MS		Char/Sec
 *
 * 1		15		66		* Out of standard range *
 *
 * 2		30		33
 * 6		90		11
 * 12		180		5.5
 * 30		450		2.2
 * 60		900		1.1
 * 90		1350	0.73
 * 120		1800	0.55
 *
 * V = MS / 15
 * V = 66 / CharSec
 *
 * MS = V * 15
 * MS = (66 / CharSec) * 15
 *
 * CharSec = 66 / V
 * CharSec = 66 / (MS / 15)
 */

IOHOOK_API long int hook_get_auto_repeat_rate() {
	#if defined USE_IOKIT || defined USE_COREFOUNDATION || defined USE_CARBON_LEGACY
	bool successful = false;
	SInt64 rate;
	#endif

	long int value = -1;

	#ifdef USE_IOKIT
	if (!successful) {
		io_service_t service = IOServiceGetMatchingService(kIOMasterPortDefault, IOServiceMatching(kIOHIDSystemClass));
		if (service) {
			kern_return_t kren_ret = kIOReturnError;
			io_connect_t connection;

			kren_ret = IOServiceOpen(service, mach_task_self(), kIOHIDParamConnectType, &connection);
			if (kren_ret == kIOReturnSuccess) {
				IOByteCount size = sizeof(rate);

				kren_ret = IOHIDGetParameter(connection, CFSTR(kIOHIDKeyRepeatKey), (IOByteCount) sizeof(rate), &rate, &size);
				if (kren_ret == kIOReturnSuccess) {
					/* This is in some undefined unit of time that if we happen
					 * to multiply by 900 gives us the time in milliseconds. We
					 * add 0.5 to the result so that when we cast to long we
					 * actually get a rounded result.  Saves the math.h depend.
					 *
					 *    33,333,333.0 / 1000.0 / 1000.0 / 1000.0 == 0.033333333	* Fast *
					 *   100,000,000.0 / 1000.0 / 1000.0 / 1000.0 == 0.1
  					 *   200,000,000.0 / 1000.0 / 1000.0 / 1000.0 == 0.2
  					 *   500,000,000.0 / 1000.0 / 1000.0 / 1000.0 == 0.5
					 * 1,000,000,000.0 / 1000.0 / 1000.0 / 1000.0 == 1
					 * 1,500,000,000.0 / 1000.0 / 1000.0 / 1000.0 == 1.5
					 * 2,000,000,000.0 / 1000.0 / 1000.0 / 1000.0 == 2				* Slow *
					 */
					value = (long) (900.0 * ((double) rate) / 1000.0 / 1000.0 / 1000.0 + 0.5);
					successful = true;

					logger(LOG_LEVEL_INFO,	"%s [%u]: IOHIDGetParameter: %li.\n",
							__FUNCTION__, __LINE__, value);
				}
			}
		}
	}
	#endif

	#ifdef USE_COREFOUNDATION
	if (!successful) {
		CFTypeRef pref_val = CFPreferencesCopyValue(CFSTR("KeyRepeat"), kCFPreferencesAnyApplication, kCFPreferencesCurrentUser, kCFPreferencesAnyHost);
		if (pref_val != NULL && CFGetTypeID(pref_val) == CFNumberGetTypeID()) {
			if (CFNumberGetValue((CFNumberRef) pref_val, kCFNumberSInt32Type, &rate)) {
				// This is the slider value, we must multiply by 15 to convert to milliseconds.
				value = (long) rate * 15;
				successful = true;

				logger(LOG_LEVEL_INFO,	"%s [%u]: CFPreferencesCopyValue: %li.\n",
						__FUNCTION__, __LINE__, value);
			}
		}
	}
	#endif

	#ifdef USE_CARBON_LEGACY
	if (!successful) {
		// Apple documentation states that value is in 'ticks'. I am not sure
		// what that means, but it looks a lot like the arbitrary slider value.
		rate = LMGetKeyRepThresh();
		if (rate > -1) {
			/* This is the slider value, we must multiply by 15 to convert to
			 * milliseconds.
			 */
			value = (long) rate * 15;
			successful = true;

			logger(LOG_LEVEL_INFO,	"%s [%u]: LMGetKeyRepThresh: %li.\n",
					__FUNCTION__, __LINE__, value);
		}
	}
	#endif

	return value;
}

IOHOOK_API long int hook_get_auto_repeat_delay() {
	#if defined USE_IOKIT || defined USE_COREFOUNDATION || defined USE_CARBON_LEGACY
	bool successful = false;
	SInt64 delay;
	#endif

	long int value = -1;

	#ifdef USE_IOKIT
	if (!successful) {
		io_service_t service = IOServiceGetMatchingService(kIOMasterPortDefault, IOServiceMatching(kIOHIDSystemClass));
		if (service) {
			kern_return_t kren_ret = kIOReturnError;
			io_connect_t connection;

			kren_ret = IOServiceOpen(service, mach_task_self(), kIOHIDParamConnectType, &connection);
			if (kren_ret == kIOReturnSuccess) {
				IOByteCount size = sizeof(delay);

				kren_ret = IOHIDGetParameter(connection, CFSTR(kIOHIDInitialKeyRepeatKey), (IOByteCount) sizeof(delay), &delay, &size);
				if (kren_ret == kIOReturnSuccess) {
					/* This is in some undefined unit of time that if we happen
					 * to multiply by 900 gives us the time in milliseconds. We
					 * add 0.5 to the result so that when we cast to long we
					 * actually get a rounded result.  Saves the math.h depend.
					 *
					 *    33,333,333.0 / 1000.0 / 1000.0 / 1000.0 == 0.033333333	* Fast *
					 *   100,000,000.0 / 1000.0 / 1000.0 / 1000.0 == 0.1
  					 *   200,000,000.0 / 1000.0 / 1000.0 / 1000.0 == 0.2
  					 *   500,000,000.0 / 1000.0 / 1000.0 / 1000.0 == 0.5
					 * 1,000,000,000.0 / 1000.0 / 1000.0 / 1000.0 == 1
					 * 1,500,000,000.0 / 1000.0 / 1000.0 / 1000.0 == 1.5
					 * 2,000,000,000.0 / 1000.0 / 1000.0 / 1000.0 == 2				* Slow *
					 */
					value = (long) (900.0 * ((double) delay) / 1000.0 / 1000.0 / 1000.0 + 0.5);
					successful = true;

					logger(LOG_LEVEL_INFO,	"%s [%u]: IOHIDGetParameter: %li.\n",
							__FUNCTION__, __LINE__, value);
				}
			}
		}
	}
	#endif

	#ifdef USE_COREFOUNDATION
	if (!successful) {
		CFTypeRef pref_val = CFPreferencesCopyValue(CFSTR("InitialKeyRepeat"), kCFPreferencesAnyApplication, kCFPreferencesCurrentUser, kCFPreferencesAnyHost);
		if (pref_val != NULL && CFGetTypeID(pref_val) == CFNumberGetTypeID()) {
			if (CFNumberGetValue((CFNumberRef) pref_val, kCFNumberSInt32Type, &delay)) {
				// This is the slider value, we must multiply by 15 to convert to
				// milliseconds.
				value = (long) delay * 15;
				successful = true;

				logger(LOG_LEVEL_INFO,	"%s [%u]: CFPreferencesCopyValue: %li.\n",
						__FUNCTION__, __LINE__, value);
			}
		}
	}
	#endif

	#ifdef USE_CARBON_LEGACY
	if (!successful) {
		// Apple documentation states that value is in 'ticks'. I am not sure
		// what that means, but it looks a lot like the arbitrary slider value.
		delay = LMGetKeyThresh();
		if (delay > -1) {
			// This is the slider value, we must multiply by 15 to convert to
			// milliseconds.
			value = (long) delay * 15;
			successful = true;

			logger(LOG_LEVEL_INFO,	"%s [%u]: LMGetKeyThresh: %li.\n",
					__FUNCTION__, __LINE__, value);
		}
	}
	#endif

	return value;
}

IOHOOK_API long int hook_get_pointer_acceleration_multiplier() {
	#if defined USE_IOKIT || defined USE_COREFOUNDATION
	bool successful = false;
	double multiplier;
	#endif

	long int value = -1;

	#ifdef USE_IOKIT
	if (!successful) {
		io_service_t service = IOServiceGetMatchingService(kIOMasterPortDefault, IOServiceMatching(kIOHIDSystemClass));

		if (service) {
			kern_return_t kren_ret = kIOReturnError;
			io_connect_t connection;

			kren_ret = IOServiceOpen(service, mach_task_self(), kIOHIDParamConnectType, &connection);
			if (kren_ret == kIOReturnSuccess) {
				// IOByteCount size = sizeof(multiplier);

				kren_ret = IOHIDGetAccelerationWithKey(connection, CFSTR(kIOHIDMouseAccelerationType), &multiplier);
				if (kren_ret == kIOReturnSuccess) {
					// Calculate the greatest common factor.

					unsigned long denominator = 1000000, d = denominator;
					unsigned long numerator = multiplier * denominator, gcf = numerator;

					while (d != 0) {
						unsigned long i = gcf % d;
						gcf = d;
						d = i;
					}

					value = denominator / gcf;
					successful = true;

					logger(LOG_LEVEL_INFO,	"%s [%u]: IOHIDGetAccelerationWithKey: %li.\n",
							__FUNCTION__, __LINE__, value);
				}
			}
		}
	}
	#endif

	#ifdef USE_COREFOUNDATION
	if (!successful) {
		CFTypeRef pref_val = CFPreferencesCopyValue(CFSTR("com.apple.mouse.scaling"), kCFPreferencesAnyApplication, kCFPreferencesCurrentUser, kCFPreferencesAnyHost);
		if (pref_val != NULL && CFGetTypeID(pref_val) == CFNumberGetTypeID()) {
			if (CFNumberGetValue((CFNumberRef) pref_val, kCFNumberSInt32Type, &multiplier)) {
				value = (long) multiplier;

				logger(LOG_LEVEL_INFO,	"%s [%u]: CFPreferencesCopyValue: %li.\n",
						__FUNCTION__, __LINE__, value);
			}
		}
	}
	#endif

	return value;
}

IOHOOK_API long int hook_get_pointer_acceleration_threshold() {
	#if defined USE_COREFOUNDATION
	bool successful = false;
	SInt32 threshold;
	#endif

	long int value = -1;

	#ifdef USE_COREFOUNDATION
	if (!successful) {
		CFTypeRef pref_val = CFPreferencesCopyValue(CFSTR("mouseDriverMaxSpeed"), CFSTR("com.apple.universalaccess"), kCFPreferencesCurrentUser, kCFPreferencesAnyHost);
		if (pref_val != NULL && CFGetTypeID(pref_val) == CFNumberGetTypeID()) {
			if (CFNumberGetValue((CFNumberRef) pref_val, kCFNumberSInt32Type, &threshold)) {
				value = (long) threshold;

				logger(LOG_LEVEL_INFO,	"%s [%u]: CFPreferencesCopyValue: %li.\n",
						__FUNCTION__, __LINE__, value);
			}
		}
	}
	#endif

	return value;
}

IOHOOK_API long int hook_get_pointer_sensitivity() {
	#ifdef USE_IOKIT
	bool successful = false;
	double sensitivity;
	#endif

	long int value = -1;

	#ifdef USE_IOKIT
	if (!successful) {
		io_service_t service = IOServiceGetMatchingService(kIOMasterPortDefault, IOServiceMatching(kIOHIDSystemClass));

		if (service) {
			kern_return_t kren_ret = kIOReturnError;
			io_connect_t connection;

			kren_ret = IOServiceOpen(service, mach_task_self(), kIOHIDParamConnectType, &connection);
			if (kren_ret == kIOReturnSuccess) {
				// IOByteCount size = sizeof(multiplier);

				kren_ret = IOHIDGetAccelerationWithKey(connection, CFSTR(kIOHIDMouseAccelerationType), &sensitivity);
				if (kren_ret == kIOReturnSuccess) {
					// Calculate the greatest common factor.

					unsigned long denominator = 1000000, d = denominator;
					unsigned long numerator = sensitivity * denominator, gcf = numerator;

					while (d != 0) {
						unsigned long i = gcf % d;
						gcf = d;
						d = i;
					}

					value = numerator / gcf;
					successful = true;

					logger(LOG_LEVEL_INFO,	"%s [%u]: IOHIDGetAccelerationWithKey: %li.\n",
							__FUNCTION__, __LINE__, value);
				}
			}
		}
	}
	#endif

	return value;
}

IOHOOK_API long int hook_get_multi_click_time() {
	#if defined USE_IOKIT || defined USE_COREFOUNDATION || defined USE_CARBON_LEGACY
	bool successful = false;
	#if defined USE_IOKIT || defined USE_CARBON_LEGACY
	// This needs to be defined only if we have USE_IOKIT or USE_CARBON_LEGACY.
	SInt64 time;
	#endif
	#endif

	long int value = -1;

	#ifdef USE_IOKIT
	if (!successful) {
		io_service_t service = IOServiceGetMatchingService(kIOMasterPortDefault, IOServiceMatching(kIOHIDSystemClass));
		if (service) {
			kern_return_t kren_ret = kIOReturnError;
			io_connect_t connection;

			kren_ret = IOServiceOpen(service, mach_task_self(), kIOHIDParamConnectType, &connection);
			if (kren_ret == kIOReturnSuccess) {
				IOByteCount size = sizeof(time);

				kren_ret = IOHIDGetParameter(connection, CFSTR(kIOHIDClickTimeKey), (IOByteCount) sizeof(time), &time, &size);
				if (kren_ret == kIOReturnSuccess) {
					/* This is in some undefined unit of time that if we happen
					 * to multiply by 900 gives us the time in milliseconds. We
					 * add 0.5 to the result so that when we cast to long we
					 * actually get a rounded result.  Saves the math.h depend.
					 */
					value = (long) (900.0 * ((double) time) / 1000.0 / 1000.0 / 1000.0 + 0.5);
					successful = true;

					logger(LOG_LEVEL_INFO,	"%s [%u]: IOHIDGetParameter: %li.\n",
							__FUNCTION__, __LINE__, value);
				}
			}
		}
	}
	#endif

	#ifdef USE_COREFOUNDATION
	if (!successful) {
		Float32 clicktime;
		CFTypeRef pref_val = CFPreferencesCopyValue(CFSTR("com.apple.mouse.doubleClickThreshold"), kCFPreferencesAnyApplication, kCFPreferencesCurrentUser, kCFPreferencesAnyHost);
		if (pref_val != NULL && CFGetTypeID(pref_val) == CFNumberGetTypeID()) {
			if (CFNumberGetValue((CFNumberRef) pref_val, kCFNumberFloat32Type, &clicktime)) {
				/* This is in some undefined unit of time that if we happen
				 * to multiply by 900 gives us the time in milliseconds.  It is
				 * completely possible that this value is in seconds and should be
				 * multiplied by 1000 but because IOKit values are undocumented and
				 * I have no idea what a Carbon 'tick' is so there really is no way
				 * to confirm this.
				 */
				value = (long) (clicktime * 900);

				logger(LOG_LEVEL_INFO,	"%s [%u]: CFPreferencesCopyValue: %li.\n",
						__FUNCTION__, __LINE__, value);
			}
		}
	}
	#endif

	#ifdef USE_CARBON_LEGACY
	if (!successful) {
		// Apple documentation states that value is in 'ticks'. I am not sure
		// what that means, but it looks a lot like the arbitrary slider value.
		time = GetDblTime();
		if (time > -1) {
			// This is the slider value, we must multiply by 15 to convert to
			// milliseconds.
			value = (long) time * 15;
			successful = true;

			logger(LOG_LEVEL_INFO,	"%s [%u]: GetDblTime: %li.\n",
					__FUNCTION__, __LINE__, value);
		}
	}
	#endif

	return value;
}


// Create a shared object constructor.
__attribute__ ((constructor))
void on_library_load() {
	// Initialize Native Input Functions.
	load_input_helper();
}

// Create a shared object destructor.
__attribute__ ((destructor))
void on_library_unload() {
	// Disable the event hook.
	//hook_stop();

	// Cleanup native input functions.
	unload_input_helper();
}
