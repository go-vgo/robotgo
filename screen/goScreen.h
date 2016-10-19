#include "../base/types.h"
#include "screengrab_init.h"
#include "screen_init.h"
// #include "../MMBitmap_init.h"

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

	char* s=(char*)calloc(100,sizeof(char*));
    if(s)strcpy(s,hex);

	return s;
}

MMSize aGetScreenSize(){
	//Get display size.
	MMSize displaySize = getMainDisplaySize();
	return displaySize;
}

char* aGetXDisplayName(){
	#if defined(USE_X11)
	const char* display = getXDisplay();
	char* sd=(char*)calloc(100,sizeof(char*));
    if(sd)strcpy(sd,display);

	return sd;
	#else
	return "getXDisplayName is only supported on Linux";
	#endif
}

char* aSetXDisplayName(char* name){
	#if defined(USE_X11)
	setXDisplay(name);
	return "success";
	#else
	return "setXDisplayName is only supported on Linux";
	#endif
}

MMBitmapRef aCaptureScreen(int x,int y,int w,int h){
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

