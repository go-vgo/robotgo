#include "MMBitmap.h"
#include <assert.h>
#include <string.h>

MMBitmapRef createMMBitmap_c(uint8_t *buffer, int32_t width, int32_t height, 
	int32_t bytewidth, uint8_t bitsPerPixel, uint8_t bytesPerPixel
) {
	MMBitmapRef bitmap = malloc(sizeof(MMBitmap));
	if (bitmap == NULL) { return NULL; }

	bitmap->imageBuffer = buffer;
	bitmap->width = width;
	bitmap->height = height;
	bitmap->bytewidth = bytewidth;
	bitmap->bitsPerPixel = bitsPerPixel;
	bitmap->bytesPerPixel = bytesPerPixel;

	return bitmap;
}

void destroyMMBitmap(MMBitmapRef bitmap) {
	assert(bitmap != NULL);

	if (bitmap->imageBuffer != NULL) {
		free(bitmap->imageBuffer);
		bitmap->imageBuffer = NULL;
	}

	free(bitmap);
}

void destroyMMBitmapBuffer(char * bitmapBuffer, void * hint) {
	if (bitmapBuffer != NULL) {
		free(bitmapBuffer);
	}
}
