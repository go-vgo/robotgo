// Copyright 2016 The go-vgo Project Developers. See the COPYRIGHT
// file at the top-level directory of this distribution and at
// https://github.com/go-vgo/robotgo/blob/master/LICENSE
//
// Licensed under the Apache License, Version 2.0 <LICENSE-APACHE or
// http://www.apache.org/licenses/LICENSE-2.0> or the MIT license
// <LICENSE-MIT or http://opensource.org/licenses/MIT>, at your
// option. This file may not be copied, modified, or distributed
// except according to those terms.

#include "../base/types.h"
#include "../base/rgb.h"
#include "screengrab_c.h"
#include "screen_c.h"
// #include "../MMBitmap_c.h"

void padHex(MMRGBHex color, char* hex){
	// Length needs to be 7 because snprintf includes a terminating null.
	// Use %06x to pad hex value with leading 0s.
	snprintf(hex, 7, "%06x", color);
}

char* pad_hex(MMRGBHex color){
	char hex[7];
	padHex(color, hex);
	// destroyMMBitmap(bitmap);

	char* str = (char*)calloc(100, sizeof(char*));
    if (str) { strcpy(str, hex); }

	return str;
}

static uint8_t rgb[3];

uint8_t* color_hex_to_rgb(uint32_t h){
	rgb[0] = RED_FROM_HEX(h);
	rgb[1] = GREEN_FROM_HEX(h);
	rgb[2] = BLUE_FROM_HEX(h);
	return rgb;
}

uint32_t color_rgb_to_hex(uint8_t r, uint8_t g, uint8_t b){
	return RGB_TO_HEX(r, g, b);
}

MMRGBHex get_px_color(int32_t x, int32_t y){
	MMBitmapRef bitmap;
	MMRGBHex color;

	if (!pointVisibleOnMainDisplay(MMPointInt32Make(x, y))){
		return color;
	}

	bitmap = copyMMBitmapFromDisplayInRect(MMRectInt32Make(x, y, 1, 1));
	// bitmap = MMRectMake(x, y, 1, 1);

	color = MMRGBHexAtPoint(bitmap, 0, 0);
	destroyMMBitmap(bitmap);

	return color;
}

char* get_pixel_color(int32_t x, int32_t y){
	MMRGBHex color = get_px_color(x, y);

	char* s = pad_hex(color);
	return s;
}

MMSizeInt32 get_screen_size(){
	// Get display size.
	MMSizeInt32 displaySize = getMainDisplaySize();
	return displaySize;
}

char* set_XDisplay_name(char* name){
	#if defined(USE_X11)
		setXDisplay(name);
		return "success";
	#else
		return "SetXDisplayName is only supported on Linux";
	#endif
}

char* get_XDisplay_name(){
	#if defined(USE_X11)
		const char* display = getXDisplay();
		char* sd = (char*)calloc(100, sizeof(char*));

		if (sd) { strcpy(sd, display); }
		return sd;
	#else
		return "GetXDisplayName is only supported on Linux";
	#endif
}

// capture_screen capture screen
MMBitmapRef capture_screen(int32_t x, int32_t y, int32_t w, int32_t h){
	// if (){
	// 	x = 0;
	// 	y = 0;
	// 	// Get screen size.
	// 	MMSize displaySize = getMainDisplaySize();
	// 	w = displaySize.width;
	// 	h = displaySize.height;
	// }
	MMBitmapRef bitmap = copyMMBitmapFromDisplayInRect(MMRectInt32Make(x, y, w, h));
	// printf("%s\n", bitmap);
	return bitmap;
}

