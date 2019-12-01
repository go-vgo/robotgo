#pragma once
#ifndef TYPES_H
#define TYPES_H

#include "os.h"
#include "inline_keywords.h" /* For H_INLINE */
#include <stddef.h>
#include <stdint.h>

/* Some generic, cross-platform types. */

struct _MMPoint {
	size_t x;
	size_t y;
};

typedef struct _MMPoint MMPoint;

struct _MMPointInt32 {
	int32_t x;
	int32_t y;
};

typedef struct _MMPointInt32 MMPointInt32;

struct _MMSize {
	size_t width;
	size_t height;
};

typedef struct _MMSize MMSize;

struct _MMSizeInt32 {
	int32_t w;
	int32_t h;
};

typedef struct _MMSizeInt32 MMSizeInt32;


struct _MMRect {
	MMPoint origin;
	MMSize size;
};

typedef struct _MMRect MMRect;

struct _MMRectInt32 {
	MMPointInt32 origin;
	MMSizeInt32 size;
};

typedef struct _MMRectInt32 MMRectInt32;

H_INLINE MMPoint MMPointMake(size_t x, size_t y)
{
	MMPoint point;
	point.x = x;
	point.y = y;
	return point;
}

H_INLINE MMPointInt32 MMPointInt32Make(int32_t x, int32_t y)
{
	MMPointInt32 point;
	point.x = x;
	point.y = y;
	return point;
}

H_INLINE MMSize MMSizeMake(size_t width, size_t height)
{
	MMSize size;
	size.width = width;
	size.height = height;
	return size;
}

H_INLINE MMSizeInt32 MMSizeInt32Make(int32_t w, int32_t h)
{
	MMSizeInt32 size;
	size.w = w;
	size.h = h;
	return size;
}

H_INLINE MMRect MMRectMake(size_t x, size_t y, size_t width, size_t height)
{
	MMRect rect;
	rect.origin = MMPointMake(x, y);
	rect.size = MMSizeMake(width, height);
	return rect;
}

H_INLINE MMRectInt32 MMRectInt32Make(int32_t x, int32_t y, int32_t w, int32_t h)
{
	MMRectInt32 rect;
	rect.origin = MMPointInt32Make(x, y);
	rect.size = MMSizeInt32Make(w, h);
	return rect;
}

//
#define MMPointZero MMPointMake(0, 0)

#if defined(IS_MACOSX)

#define CGPointFromMMPoint(p) CGPointMake((CGFloat)(p).x, (CGFloat)(p).y)
#define MMPointFromCGPoint(p) MMPointMake((size_t)(p).x, (size_t)(p).y)

#define CGPointFromMMPointInt32(p) CGPointMake((CGFloat)(p).x, (CGFloat)(p).y)
#define MMPointInt32FromCGPoint(p) MMPointInt32Make((int32_t)(p).x, (int32_t)(p).y)

#elif defined(IS_WINDOWS)

#define MMPointFromPOINT(p) MMPointMake((size_t)p.x, (size_t)p.y)
#define MMPointInt32FromPOINT(p) MMPointInt32Make((int32_t)p.x, (int32_t)p.y)

#endif

#endif /* TYPES_H */
