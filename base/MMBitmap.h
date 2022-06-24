#pragma once
#ifndef MMBITMAP_H
#define MMBITMAP_H

#include "types.h"
#include "rgb.h"
#include <assert.h>
#include <stdint.h>

struct _MMBitmap {
	uint8_t *imageBuffer;  /* Pixels stored in Quad I format; */
	int32_t width;          /* Never 0, unless image is NULL. */
	int32_t height;         /* Never 0, unless image is NULL. */
	
	int32_t bytewidth;      /* The aligned width (width + padding). */
	uint8_t bitsPerPixel;  /* Should be either 24 or 32. */
	uint8_t bytesPerPixel; /* For convenience; should be bitsPerPixel / 8. */
};

typedef struct _MMBitmap MMBitmap;
typedef MMBitmap *MMBitmapRef;

#define MMBitmapPointInBounds(image, p) ((p).x < (image)->width && (p).y < (image)->height)

/* Get pointer to pixel of MMBitmapRef. No bounds checking is performed */
#define MMRGBColorRefAtPoint(image, x, y) \
			(MMRGBColor *)(assert(MMBitmapPointInBounds(image, MMPointInt32Make(x, y))), \
	        ((image)->imageBuffer) + (((image)->bytewidth * (y)) + ((x) * (image)->bytesPerPixel)))

/* Dereference pixel of MMBitmapRef. Again, no bounds checking is performed. */
#define MMRGBColorAtPoint(image, x, y) *MMRGBColorRefAtPoint(image, x, y)

/* Hex/integer value of color at point. */
#define MMRGBHexAtPoint(image, x, y) hexFromMMRGB(MMRGBColorAtPoint(image, x, y))

#endif /* MMBITMAP_H */