// #include "os.h"
#if defined(IS_MACOSX)
	#include <CoreFoundation/CoreFoundation.h>
#endif

#if defined(IS_MACOSX)
	CFStringRef CFStringCreateWithUTF8String(const char *title) {
		if (title == NULL) { return NULL; }
		return CFStringCreateWithCString(NULL, title, kCFStringEncodingUTF8);
	}
#endif

int showAlert(const char *title, const char *msg, 
		const char *defaultButton, const char *cancelButton) {
	#if defined(IS_MACOSX)
		CFStringRef alertHeader = CFStringCreateWithUTF8String(title);
		CFStringRef alertMessage = CFStringCreateWithUTF8String(msg);
		CFStringRef defaultButtonTitle = CFStringCreateWithUTF8String(defaultButton);
		CFStringRef cancelButtonTitle = CFStringCreateWithUTF8String(cancelButton);
		CFOptionFlags responseFlags;
		
		SInt32 err = CFUserNotificationDisplayAlert(
			0.0, kCFUserNotificationNoteAlertLevel, NULL, NULL, NULL, alertHeader, alertMessage,
			defaultButtonTitle, cancelButtonTitle, NULL, &responseFlags);
												
		if (alertHeader != NULL) CFRelease(alertHeader);
		if (alertMessage != NULL) CFRelease(alertMessage);
		if (defaultButtonTitle != NULL) CFRelease(defaultButtonTitle);
		if (cancelButtonTitle != NULL) CFRelease(cancelButtonTitle);

		if (err != 0) { return -1; }
		return (responseFlags == kCFUserNotificationDefaultResponse) ? 0 : 1;
	#elif defined(USE_X11)
		return 0;
	#else
		/* TODO: Display custom buttons instead of the pre-defined "OK" and "Cancel". */
		int response = MessageBox(NULL, msg, title,
								(cancelButton == NULL) ? MB_OK : MB_OKCANCEL );
		return (response == IDOK) ? 0 : 1;
	#endif
}

