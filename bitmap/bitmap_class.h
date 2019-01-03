#pragma once
// #ifndef BITMAP_CLASS_H
// #define BITMAP_CLASS_H

#include "../base/MMBitmap.h"

/* This file defines the class "Bitmap" for dealing with raw bitmaps. */
struct _BitmapObject {
	MMBitmapRef bitmap;
	MMPoint point; /* For iterator */
};

typedef struct _BitmapObject BitmapObject;

// extern PyTypeObject Bitmap_Type;

/* Returns a newly-initialized BitmapObject from the given MMBitmap.
 * The reference to |bitmap| is "stolen"; i.e., only the pointer is copied, and
 * the reponsibility for free()'ing the buffer is given to the |BitmapObject|.
 *
 * Remember to call PyType_Ready() before using this for the first time! */
BitmapObject BitmapObject_FromMMBitmap(MMBitmapRef bitmap);

// #endif /* PY_BITMAP_CLASS_H */