#pragma once
#ifndef SCREENGRAB_H
#define SCREENGRAB_H

#include "../base/types.h"
// #include "../base/MMBitmap_c.h"
#include "../base/bitmap_free_c.h"

#ifdef __cplusplus
extern "C"
{
#endif

/* Returns a raw bitmap of screengrab of the display (to be destroyed()'d by
 * caller), or NULL on error. */
MMBitmapRef copyMMBitmapFromDisplayInRect(MMRectInt32 rect, int32_t display_id);

#ifdef __cplusplus
}
#endif

#endif /* SCREENGRAB_H */
