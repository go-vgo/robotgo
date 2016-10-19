#include "../base/types.h"
#include "mouse_init.h"

//Global delays.
int mouseDelay = 10;
// int keyboardDelay = 10;


// int CheckMouseButton(const char * const b, MMMouseButton * const button){
// 	if (!button) return -1;

// 	if (strcmp(b, "left") == 0)
// 	{
// 		*button = LEFT_BUTTON;
// 	}
// 	else if (strcmp(b, "right") == 0)
// 	{
// 		*button = RIGHT_BUTTON;
// 	}
// 	else if (strcmp(b, "middle") == 0)
// 	{
// 		*button = CENTER_BUTTON;
// 	}
// 	else
// 	{
// 		return -2;
// 	}

// 	return 0;
// }

int amoveMouse(size_t x, size_t y){
	MMPoint point;
	//int x =103;
	//int y =104;
	point = MMPointMake(x, y);
	moveMouse(point);

	return 0;
}

int adragMouse(size_t x, size_t y){
	// const size_t x=10;
	// const size_t y=20;
	MMMouseButton button = LEFT_BUTTON;

	MMPoint point;
	point = MMPointMake(x, y);
	dragMouse(point, button);
	microsleep(mouseDelay);

	// printf("%s\n","gyp-----");
	return 0;
}

int amoveMouseSmooth(size_t x, size_t y){
	MMPoint point;
	point = MMPointMake(x, y);
	smoothlyMoveMouse(point);
	microsleep(mouseDelay);

	return 0;

}

MMPoint agetMousePos(){
	MMPoint pos = getMousePos();

	//Return object with .x and .y.
	// printf("%zu\n%zu\n", pos.x, pos.y );
	return pos;
}

int amouseClick(){
	MMMouseButton button = LEFT_BUTTON;
	bool doubleC = false;

	if (!doubleC){
		clickMouse(button);
	}else{
		doubleClick(button);
	}

	microsleep(mouseDelay);

	return 0;
}

int amouseToggle(){
	MMMouseButton button = LEFT_BUTTON;
	bool down = false;

	return 0;
}

int asetMouseDelay(size_t val){
	// int val=10;
	mouseDelay = val;

	return 0;
}

int ascrollMouse(size_t scrollMagnitude,char *s){
	// int scrollMagnitude = 20;

	MMMouseWheelDirection scrollDirection;

	if (strcmp(s, "up") == 0){
		scrollDirection = DIRECTION_UP;
	}else if (strcmp(s, "down") == 0){
			scrollDirection = DIRECTION_DOWN;
	}else{
		// return "Invalid scroll direction specified.";
		return 1;
	}

	scrollMouse(scrollMagnitude, scrollDirection);
	microsleep(mouseDelay);

	return 0;
}
