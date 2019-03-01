package gosseract

// #if __FreeBSD__ >= 10
// #cgo LDFLAGS: -L/usr/local/lib -llept -ltesseract
// #else
// #cgo CXXFLAGS: -std=c++0x
// #cgo LDFLAGS: -llept -ltesseract
// #endif
// #include <stdlib.h>
// #include <stdbool.h>
// #include "tessbridge.h"
import "C"
import (
	"fmt"
	"image"
	"os"
	"strings"
	"unsafe"
)

// Version returns the version of Tesseract-OCR
func Version() string {
	api := C.Create()
	defer C.Free(api)
	version := C.Version(api)
	return C.GoString(version)
}

// ClearPersistentCache clears any library-level memory caches. There are a variety of expensive-to-load constant data structures (mostly language dictionaries) that are cached globally – surviving the Init() and End() of individual TessBaseAPI's. This function allows the clearing of these caches.
func ClearPersistentCache() {
	api := C.Create()
	defer C.Free(api)
	C.ClearPersistentCache(api)
}

// Client is argument builder for tesseract::TessBaseAPI.
type Client struct {
	api C.TessBaseAPI

	// Holds a reference to the pix image to be able to destroy on client close
	// or when a new image is set
	pixImage C.PixImage

	// Trim specifies characters to trim, which would be trimed from result string.
	// As results of OCR, text often contains unnecessary characters, such as newlines, on the head/foot of string.
	// If `Trim` is set, this client will remove specified characters from the result.
	Trim bool

	// TessdataPrefix can indicate directory path to `tessdata`.
	// It is set `/usr/local/share/tessdata/` or something like that, as default.
	// TODO: Implement and test
	TessdataPrefix *string

	// Languages are languages to be detected. If not specified, it's gonna be "eng".
	Languages []string

	// Variables is just a pool to evaluate "tesseract::TessBaseAPI->SetVariable" in delay.
	// TODO: Think if it should be public, or private property.
	Variables map[SettableVariable]string

	// Config is a file path to the configuration for Tesseract
	// See http://www.sk-spell.sk.cx/tesseract-ocr-parameters-in-302-version
	// TODO: Fix link to official page
	ConfigFilePath string
}

// NewClient construct new Client. It's due to caller to Close this client.
func NewClient() *Client {
	client := &Client{
		api:       C.Create(),
		Variables: map[SettableVariable]string{},
		Trim:      true,
	}
	return client
}

// Close frees allocated API. This MUST be called for ANY client constructed by "NewClient" function.
func (client *Client) Close() (err error) {
	// defer func() {
	// 	if e := recover(); e != nil {
	// 		err = fmt.Errorf("%v", e)
	// 	}
	// }()
	C.Clear(client.api)
	C.Free(client.api)
	if client.pixImage != nil {
		C.DestroyPixImage(client.pixImage)
		client.pixImage = nil
	}
	return err
}

// SetImage sets path to image file to be processed OCR.
func (client *Client) SetImage(imagepath string) error {

	if client.api == nil {
		return fmt.Errorf("TessBaseAPI is not constructed, please use `gosseract.NewClient`")
	}
	if imagepath == "" {
		return fmt.Errorf("image path cannot be empty")
	}
	if _, err := os.Stat(imagepath); err != nil {
		return fmt.Errorf("cannot detect the stat of specified file: %v", err)
	}

	if client.pixImage != nil {
		C.DestroyPixImage(client.pixImage)
		client.pixImage = nil
	}

	p := C.CString(imagepath)
	defer C.free(unsafe.Pointer(p))

	img := C.CreatePixImageByFilePath(p)
	client.pixImage = img

	return nil
}

// SetImageFromBytes sets the image data to be processed OCR.
func (client *Client) SetImageFromBytes(data []byte) error {

	if client.api == nil {
		return fmt.Errorf("TessBaseAPI is not constructed, please use `gosseract.NewClient`")
	}
	if len(data) == 0 {
		return fmt.Errorf("image data cannot be empty")
	}

	if client.pixImage != nil {
		C.DestroyPixImage(client.pixImage)
		client.pixImage = nil
	}

	img := C.CreatePixImageFromBytes((*C.uchar)(unsafe.Pointer(&data[0])), C.int(len(data)))
	client.pixImage = img

	return nil
}

// SetLanguage sets languages to use. English as default.
func (client *Client) SetLanguage(langs ...string) error {
	if len(langs) == 0 {
		return fmt.Errorf("languages cannot be empty")
	}
	client.Languages = langs
	return nil
}

func (client *Client) DisableOutput() error {
	return client.SetVariable(DEBUG_FILE, os.DevNull)
}

// SetWhitelist sets whitelist chars.
// See official documentation for whitelist here https://github.com/tesseract-ocr/tesseract/wiki/ImproveQuality#dictionaries-word-lists-and-patterns
func (client *Client) SetWhitelist(whitelist string) error {
	return client.SetVariable(TESSEDIT_CHAR_WHITELIST, whitelist)
}

// SetBlacklist sets whitelist chars.
// See official documentation for whitelist here https://github.com/tesseract-ocr/tesseract/wiki/ImproveQuality#dictionaries-word-lists-and-patterns
func (client *Client) SetBlacklist(whitelist string) error {
	return client.SetVariable(TESSEDIT_CHAR_BLACKLIST, whitelist)
}

// SetVariable sets parameters, representing tesseract::TessBaseAPI->SetVariable.
// See official documentation here https://zdenop.github.io/tesseract-doc/classtesseract_1_1_tess_base_a_p_i.html#a2e09259c558c6d8e0f7e523cbaf5adf5
// Because `api->SetVariable` must be called after `api->Init`, this method cannot detect unexpected key for variables.
// Check `client.setVariablesToInitializedAPI` for more information.
func (client *Client) SetVariable(key SettableVariable, value string) error {
	client.Variables[key] = value
	return nil
}

// SetPageSegMode sets "Page Segmentation Mode" (PSM) to detect layout of characters.
// See official documentation for PSM here https://github.com/tesseract-ocr/tesseract/wiki/ImproveQuality#page-segmentation-method
// See https://github.com/otiai10/gosseract/issues/52 for more information.
func (client *Client) SetPageSegMode(mode PageSegMode) error {
	C.SetPageSegMode(client.api, C.int(mode))
	return nil
}

// SetConfigFile sets the file path to config file.
func (client *Client) SetConfigFile(fpath string) error {
	info, err := os.Stat(fpath)
	if err != nil {
		return err
	}
	if info.IsDir() {
		return fmt.Errorf("the specified config file path seems to be a directory")
	}
	client.ConfigFilePath = fpath
	return nil
}

// Initialize tesseract::TessBaseAPI
// TODO: add tessdata prefix
func (client *Client) init() error {

	var languages *C.char
	if len(client.Languages) != 0 {
		languages = C.CString(strings.Join(client.Languages, "+"))
	}
	defer C.free(unsafe.Pointer(languages))

	var configfile *C.char
	if _, err := os.Stat(client.ConfigFilePath); err == nil {
		configfile = C.CString(client.ConfigFilePath)
	}
	defer C.free(unsafe.Pointer(configfile))

	errbuf := [512]C.char{}
	res := C.Init(client.api, nil, languages, configfile, &errbuf[0])
	msg := C.GoString(&errbuf[0])

	if res != 0 {
		return fmt.Errorf("failed to initialize TessBaseAPI with code %d: %s", res, msg)
	}

	if err := client.setVariablesToInitializedAPI(); err != nil {
		return err
	}

	if client.pixImage == nil {
		return fmt.Errorf("PixImage is not set, use SetImage or SetImageFromBytes before Text or HOCRText")
	}
	C.SetPixImage(client.api, client.pixImage)

	return nil
}

// This method sets all the sspecified variables to TessBaseAPI structure.
// Because `api->SetVariable` must be called after `api->Init()`,
// gosseract.Client.SetVariable cannot call `api->SetVariable` directly.
// See https://zdenop.github.io/tesseract-doc/classtesseract_1_1_tess_base_a_p_i.html#a2e09259c558c6d8e0f7e523cbaf5adf5
func (client *Client) setVariablesToInitializedAPI() error {
	for key, value := range client.Variables {
		k, v := C.CString(string(key)), C.CString(value)
		defer C.free(unsafe.Pointer(k))
		defer C.free(unsafe.Pointer(v))
		res := C.SetVariable(client.api, k, v)
		if bool(res) == false {
			return fmt.Errorf("failed to set variable with key(%v) and value(%v)", key, value)
		}
	}
	return nil
}

// Text finally initialize tesseract::TessBaseAPI, execute OCR and extract text detected as string.
func (client *Client) Text() (out string, err error) {
	if err = client.init(); err != nil {
		return
	}
	out = C.GoString(C.UTF8Text(client.api))
	if client.Trim {
		out = strings.Trim(out, "\n")
	}
	return out, err
}

// HOCRText finally initialize tesseract::TessBaseAPI, execute OCR and returns hOCR text.
// See https://en.wikipedia.org/wiki/HOCR for more information of hOCR.
func (client *Client) HOCRText() (out string, err error) {
	if err = client.init(); err != nil {
		return
	}
	out = C.GoString(C.HOCRText(client.api))
	return
}

// BoundingBox contains the position, confidence and UTF8 text of the recognized word
type BoundingBox struct {
	Box        image.Rectangle
	Word       string
	Confidence float64
}

// GetBoundingBoxes returns bounding boxes for each matched word
func (client *Client) GetBoundingBoxes(level PageIteratorLevel) (out []BoundingBox, err error) {
	if client.api == nil {
		return out, fmt.Errorf("TessBaseAPI is not constructed, please use `gosseract.NewClient`")
	}
	if err = client.init(); err != nil {
		return
	}
	boxArray := C.GetBoundingBoxes(client.api, C.int(level))
	length := int(boxArray.length)
	defer C.free(unsafe.Pointer(boxArray.boxes))
	defer C.free(unsafe.Pointer(boxArray))

	for i := 0; i < length; i++ {
		// cast to bounding_box: boxes + i*sizeof(box)
		box := (*C.struct_bounding_box)(unsafe.Pointer(uintptr(unsafe.Pointer(boxArray.boxes)) + uintptr(i)*unsafe.Sizeof(C.struct_bounding_box{})))
		out = append(out, BoundingBox{
			Box:        image.Rect(int(box.x1), int(box.y1), int(box.x2), int(box.y2)),
			Word:       C.GoString(box.word),
			Confidence: float64(box.confidence),
		})
	}

	return
}
