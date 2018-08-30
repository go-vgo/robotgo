#if __FreeBSD__ >= 10
#include "/usr/local/include/tesseract/baseapi.h"
#include "/usr/local/include/leptonica/allheaders.h"
#else
#include <tesseract/baseapi.h>
#include <leptonica/allheaders.h>
#endif

#include "tessbridge.h"

TessBaseAPI Create() {
  tesseract::TessBaseAPI * api = new tesseract::TessBaseAPI();
  return (void*)api;
}

void Free(TessBaseAPI a) {
  tesseract::TessBaseAPI * api = (tesseract::TessBaseAPI*)a;
  api->End();
  delete api;
}

void Clear(TessBaseAPI a) {
  tesseract::TessBaseAPI * api = (tesseract::TessBaseAPI*)a;
  api->Clear();
}

void ClearPersistentCache(TessBaseAPI a) {
  tesseract::TessBaseAPI * api = (tesseract::TessBaseAPI*)a;
  api->ClearPersistentCache();
}

int Init(TessBaseAPI a, char* tessdataprefix, char* languages) {
  tesseract::TessBaseAPI * api = (tesseract::TessBaseAPI*)a;
  return api->Init(tessdataprefix, languages);
}

int Init(TessBaseAPI a, char* tessdataprefix, char* languages, char* configfilepath) {
  tesseract::TessBaseAPI * api = (tesseract::TessBaseAPI*)a;
  if (configfilepath != NULL) {
    char *configs[]={configfilepath};
    int configs_size = 1;
    return api->Init(tessdataprefix, languages, tesseract::OEM_DEFAULT, configs, configs_size, NULL, NULL, false);
  } else {
    return api->Init(tessdataprefix, languages);
  }
}

bool SetVariable(TessBaseAPI a, char* name, char* value) {
  tesseract::TessBaseAPI * api = (tesseract::TessBaseAPI*)a;
  return api->SetVariable(name, value);
}

void SetPixImage(TessBaseAPI a, PixImage pix) {
  tesseract::TessBaseAPI * api = (tesseract::TessBaseAPI*)a;
  Pix *image = (Pix*) pix;
  api->SetImage(image);
  if (api->GetSourceYResolution() < 70) {
    api->SetSourceResolution(70);
  }
}

void SetPageSegMode(TessBaseAPI a, int m) {
  tesseract::TessBaseAPI * api = (tesseract::TessBaseAPI*)a;
  tesseract::PageSegMode mode = (tesseract::PageSegMode)m;
  api->SetPageSegMode(mode);
}

int GetPageSegMode(TessBaseAPI a) {
  tesseract::TessBaseAPI * api = (tesseract::TessBaseAPI*)a;
  return api->GetPageSegMode();
}

char* UTF8Text(TessBaseAPI a) {
  tesseract::TessBaseAPI * api = (tesseract::TessBaseAPI*)a;
  return api->GetUTF8Text();
}

char* HOCRText(TessBaseAPI a) {
  tesseract::TessBaseAPI * api = (tesseract::TessBaseAPI*)a;
  return api->GetHOCRText(0);
}

const char* Version(TessBaseAPI a) {
  tesseract::TessBaseAPI * api = (tesseract::TessBaseAPI*)a;
  const char* v = api->Version();
  return v;
}

PixImage CreatePixImageByFilePath(char* imagepath) {
  Pix *image = pixRead(imagepath);
  return (void*)image;
}

PixImage CreatePixImageFromBytes(unsigned char* data, int size) {
  Pix *image = pixReadMem(data, (size_t)size);
  return (void*)image;
}


void DestroyPixImage(PixImage pix){
  Pix *img = (Pix*) pix;
  pixDestroy(&img);
}
