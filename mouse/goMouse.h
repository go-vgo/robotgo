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
#include "mouse_c.h"

// Global delays.
int mouseDelay = 0;
// int keyboardDelay = 10;


// int CheckMouseButton(const char * const b,
// 	MMMouseButton * const button){
// 	if (!button) return -1;

// 	if (strcmp(b, "left") == 0) {
// 		*button = LEFT_BUTTON;
// 	}
// 	else if (strcmp(b, "right") == 0) {
// 		*button = RIGHT_BUTTON;
// 	}
// 	else if (strcmp(b, "middle") == 0) {
// 		*button = CENTER_BUTTON;
// 	} else {
// 		return -2;
// 	}

// 	return 0;
// }

int move_mouse(int32_t x, int32_t y){
	MMPointInt32 point;
	// int x = 103;
	// int y = 104;
	point = MMPointInt32Make(x, y);
	moveMouse(point);

	return 0;
}


int drag_mouse(int32_t x, int32_t y, MMMouseButton button){
	// const size_t x = 10;
	// const size_t y = 20;
	// MMMouseButton button = LEFT_BUTTON;

	MMPointInt32 point;
	point = MMPointInt32Make(x, y);
	dragMouse(point, button);
	microsleep(mouseDelay);

	// printf("%s\n", "gyp-----");
	return 0;
}

bool move_mouse_smooth(int32_t x, int32_t y, double lowSpeed,
	double highSpeed, int msDelay){
	MMPointInt32 point;
	point = MMPointInt32Make(x, y);
	
	bool cbool = smoothlyMoveMouse(point, lowSpeed, highSpeed);
	microsleep(msDelay);

	return cbool;
}

MMPointInt32 get_mouse_pos(){
	MMPointInt32 pos = getMousePos();

	// Return object with .x and .y.
	// printf("%zu\n%zu\n", pos.x, pos.y );
	return pos;
}

int mouse_click(MMMouseButton button, bool doubleC){
	// MMMouseButton button = LEFT_BUTTON;
	// bool doubleC = false;
	if (!doubleC) {
		clickMouse(button);
	} else {
		doubleClick(button);
	}

	microsleep(mouseDelay);

	return 0;
}

int mouse_toggle(char* d, MMMouseButton button){
	// MMMouseButton button = LEFT_BUTTON;
	bool down = false;
	if (strcmp(d, "down") == 0) {
		down = true;
	} else if (strcmp(d, "up") == 0) {
		down = false;
	} else {
		return 1;
	}

	toggleMouse(down, button);
	microsleep(mouseDelay);

	return 0;
}

int set_mouse_delay(size_t val){
	// int val = 10;
	mouseDelay = val;

	return 0;
}

int scroll(int x, int y, int msDelay){
	scrollMouseXY(x, y);
	microsleep(msDelay);

	return 0;
}

int scroll_mouse(size_t scrollMagnitude, char *s){
	// int scrollMagnitude = 20;
	MMMouseWheelDirection scrollDirection;

	if (strcmp(s, "up") == 0) {
		scrollDirection = DIRECTION_UP;
	} else if (strcmp(s, "down") == 0) {
		scrollDirection = DIRECTION_DOWN;
	} else {
		// return "Invalid scroll direction specified.";
		return 1;
	}

	scrollMouse(scrollMagnitude, scrollDirection);
	microsleep(mouseDelay);

	return 0;
}
