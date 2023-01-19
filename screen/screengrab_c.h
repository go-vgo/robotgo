#include "../base/bitmap_free_c.h"
#include <stdlib.h> /* malloc() */

#if defined(IS_MACOSX)
	#include <OpenGL/OpenGL.h>
	#include <OpenGL/gl.h>
	#include <ApplicationServices/ApplicationServices.h>
#elif defined(USE_X11)
	#include <X11/Xlib.h>
	#include <X11/Xutil.h>
	#include "../base/xdisplay_c.h"
#elif defined(IS_WINDOWS)
	#include <string.h>
#endif

MMBitmapRef copyMMBitmapFromDisplayInRect(MMRectInt32 rect, int32_t display_id, int8_t isPid) {
#if defined(IS_MACOSX)
	MMBitmapRef bitmap = NULL;
	uint8_t *buffer = NULL;
	size_t bufferSize = 0;

	CGDirectDisplayID displayID = (CGDirectDisplayID) display_id;
	if (displayID == -1 || displayID == 0) {
		displayID = CGMainDisplayID();
	}

	MMPointInt32 o = rect.origin; MMSizeInt32 s = rect.size;
	CGImageRef image = CGDisplayCreateImageForRect(displayID, CGRectMake(o.x, o.y, s.w, s.h));
	if (!image) { return NULL; }

	CFDataRef imageData = CGDataProviderCopyData(CGImageGetDataProvider(image));
	if (!imageData) { return NULL; }

	bufferSize = CFDataGetLength(imageData);
	buffer = malloc(bufferSize);
	CFDataGetBytes(imageData, CFRangeMake(0, bufferSize), buffer);

	bitmap = createMMBitmap_c(buffer, 
			CGImageGetWidth(image), CGImageGetHeight(image), CGImageGetBytesPerRow(image), 
			CGImageGetBitsPerPixel(image), CGImageGetBitsPerPixel(image) / 8);

	CFRelease(imageData);
	CGImageRelease(image);

	return bitmap;
#elif defined(USE_X11)
	MMBitmapRef bitmap;
	Display *display;
	if (display_id == -1) {
		display = XOpenDisplay(NULL);
	} else {
		display = XGetMainDisplay();
	}

	MMPointInt32 o = rect.origin; MMSizeInt32 s = rect.size;
	XImage *image = XGetImage(display, XDefaultRootWindow(display), 
							(int)o.x, (int)o.y, (unsigned int)s.w, (unsigned int)s.h, AllPlanes, ZPixmap);
	XCloseDisplay(display);
	if (image == NULL) { return NULL; }

	bitmap = createMMBitmap_c((uint8_t *)image->data, 
				s.w, s.h, (size_t)image->bytes_per_line, 
				(uint8_t)image->bits_per_pixel, (uint8_t)image->bits_per_pixel / 8);
	image->data = NULL; /* Steal ownership of bitmap data so we don't have to copy it. */
	XDestroyImage(image);

	return bitmap;
#elif defined(IS_WINDOWS)
	MMBitmapRef bitmap;
	void *data;
	HDC screen = NULL, screenMem = NULL;
	HBITMAP dib;
	BITMAPINFO bi;

	int32_t x = rect.origin.x, y = rect.origin.y;
	int32_t w = rect.size.w, h = rect.size.h;

	/* Initialize bitmap info. */
	bi.bmiHeader.biSize = sizeof(bi.bmiHeader);
   	bi.bmiHeader.biWidth = (long) w;
   	bi.bmiHeader.biHeight = -(long) h; /* Non-cartesian, please */
   	bi.bmiHeader.biPlanes = 1;
   	bi.bmiHeader.biBitCount = 32;
   	bi.bmiHeader.biCompression = BI_RGB;
   	bi.bmiHeader.biSizeImage = (DWORD)(4 * w * h);
	bi.bmiHeader.biXPelsPerMeter = 0;
	bi.bmiHeader.biYPelsPerMeter = 0;
	bi.bmiHeader.biClrUsed = 0;
	bi.bmiHeader.biClrImportant = 0;

	HWND hwnd;
	if (display_id == -1 || isPid == 0) {
	// 	screen = GetDC(NULL); /* Get entire screen */
		hwnd = GetDesktopWindow();
	} else {
		hwnd = (HWND) (uintptr) display_id;
	}
	screen = GetDC(hwnd);
	
	if (screen == NULL) { return NULL; }

	// Todo: Use DXGI
	screenMem = CreateCompatibleDC(screen);
	/* Get screen data in display device context. */
   	dib = CreateDIBSection(screen, &bi, DIB_RGB_COLORS, &data, NULL, 0);

	/* Copy the data into a bitmap struct. */
	BOOL b = (screenMem == NULL) || 
		SelectObject(screenMem, dib) == NULL ||
	    !BitBlt(screenMem, (int)0, (int)0, (int)w, (int)h, screen, x, y, SRCCOPY);
	if (b) {
		/* Error copying data. */
		ReleaseDC(hwnd, screen);
		DeleteObject(dib);
		if (screenMem != NULL) { DeleteDC(screenMem); }

		return NULL;
	}

	bitmap = createMMBitmap_c(NULL, w, h, 4 * w, (uint8_t)bi.bmiHeader.biBitCount, 4);

	/* Copy the data to our pixel buffer. */
	if (bitmap != NULL) {
		bitmap->imageBuffer = malloc(bitmap->bytewidth * bitmap->height);
		memcpy(bitmap->imageBuffer, data, bitmap->bytewidth * bitmap->height);
	}

	ReleaseDC(hwnd, screen);
	DeleteObject(dib);
	DeleteDC(screenMem);

	return bitmap;
#endif
}
