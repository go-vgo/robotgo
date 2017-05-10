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
#include "screengrab_c.h"
#include "screen_c.h"
// #include "../MMBitmap_c.h"

void padHex(MMRGBHex color, char* hex){
	//Length needs to be 7 because snprintf includes a terminating null.
	//Use %06x to pad hex value with leading 0s.
	snprintf(hex, 7, "%06x", color);
}


char* aGetPixelColor(size_t x, size_t y){
	MMBitmapRef bitmap;
	MMRGBHex color;

	if (!pointVisibleOnMainDisplay(MMPointMake(x, y))){
		// return 1;
		return "screen's dimensions.";
	}

	bitmap = copyMMBitmapFromDisplayInRect(MMRectMake(x, y, 1, 1));
	// bitmap = MMRectMake(x, y, 1, 1);

	color = MMRGBHexAtPoint(bitmap, 0, 0);

	char hex[7];

	padHex(color, hex);

	destroyMMBitmap(bitmap);

	// printf("%s\n", hex);

	// return 0;

	char* s = (char*)calloc(100, sizeof(char*));
    if(s)strcpy(s, hex);

	return s;
}

MMSize aGetScreenSize(){
	//Get display size.
	MMSize displaySize = getMainDisplaySize();
	return displaySize;
}

char* aSetXDisplayName(char* name){
	#if defined(USE_X11)
	setXDisplay(name);
	return "success";
	#else
	return "setXDisplayName is only supported on Linux";
	#endif
}

char* aGetXDisplayName(){
	#if defined(USE_X11)
	const char* display = getXDisplay();
	char* sd = (char*)calloc(100, sizeof(char*));
    if(sd)strcpy(sd, display);

	return sd;
	#else
	return "getXDisplayName is only supported on Linux";
	#endif
}

MMBitmapRef aCaptureScreen(size_t x, size_t y, size_t w, size_t h){
	// if (){
	// 	x = 0;
	// 	y = 0;

	// 	//Get screen size.
	// 	MMSize displaySize = getMainDisplaySize();
	// 	w = displaySize.width;
	// 	h = displaySize.height;
	// }

	MMBitmapRef bitmap = copyMMBitmapFromDisplayInRect(MMRectMake(x, y, w, h));
	// printf("%s\n", bitmap);

	return bitmap;
}

