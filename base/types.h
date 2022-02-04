#pragma once
#ifndef TYPES_H
#define TYPES_H

#include "os.h"
#include "inline_keywords.h" /* For H_INLINE */
#include <stddef.h>
#include <stdint.h>
#include <stdbool.h>

/* Some generic, cross-platform types. */
#ifdef RobotGo_64
	typedef int64_t			intptr;	
	typedef uint64_t		uintptr;	
#else
	typedef int32_t			 intptr;
	typedef uint32_t		 uintptr;	// Unsigned pointer integer
#endif

struct _MMPointInt32 {
	int32_t x;
	int32_t y;
};
typedef struct _MMPointInt32 MMPointInt32;

struct _MMSizeInt32 {
	int32_t w;
	int32_t h;
};
typedef struct _MMSizeInt32 MMSizeInt32;

struct _MMRectInt32 {
	MMPointInt32 origin;
	MMSizeInt32 size;
};
typedef struct _MMRectInt32 MMRectInt32;

H_INLINE MMPointInt32 MMPointInt32Make(int32_t x, int32_t y) {
	MMPointInt32 point;
	point.x = x;
	point.y = y;
	return point;
}

H_INLINE MMSizeInt32 MMSizeInt32Make(int32_t w, int32_t h) {
	MMSizeInt32 size;
	size.w = w;
	size.h = h;
	return size;
}

H_INLINE MMRectInt32 MMRectInt32Make(int32_t x, int32_t y, int32_t w, int32_t h) {
	MMRectInt32 rect;
	rect.origin = MMPointInt32Make(x, y);
	rect.size = MMSizeInt32Make(w, h);
	return rect;
}

#define MMPointZero MMPointInt32Make(0, 0)

#if defined(IS_MACOSX)
	#define CGPointFromMMPointInt32(p) CGPointMake((CGFloat)(p).x, (CGFloat)(p).y)
	#define MMPointInt32FromCGPoint(p) MMPointInt32Make((int32_t)(p).x, (int32_t)(p).y)
#elif defined(IS_WINDOWS)
	#define MMPointInt32FromPOINT(p) MMPointInt32Make((int32_t)p.x, (int32_t)p.y)
#endif

#endif /* TYPES_H */
