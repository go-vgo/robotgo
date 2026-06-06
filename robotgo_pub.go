package robotgo

import (
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/go-vgo/robotgo/clipboard"
)

const (
	// Version get the robotgo version
	Version = "v2.00.0.1658, MT. Baker!"
)

// GetVersion get the robotgo version
func GetVersion() string {
	return Version
}

var (
	// MouseSleep set the mouse default millisecond sleep time
	MouseSleep = 0
	// KeySleep set the key default millisecond sleep time
	KeySleep = 10

	// DisplayID set the screen display id
	DisplayID = -1

	// NotPid used the hwnd not pid in windows
	NotPid bool
	// Scale option the os screen scale
	Scale bool
)

// Bitmap define the go Bitmap struct
//
// The common type conversion of bitmap:
//
//	https://github.com/go-vgo/robotgo/blob/master/docs/keys.md#type-conversion
type Bitmap struct {
	ImgBuf        *uint8
	Width, Height int

	Bytewidth     int
	BitsPixel     uint8
	BytesPerPixel uint8
}

// Point is point struct
type Point struct {
	X int
	Y int
}

// Size is size structure
type Size struct {
	W, H int
}

// Rect is rect structure
type Rect struct {
	Point
	Size
}

// Try handler(err)
func Try(fun func(), handler func(interface{})) {
	defer func() {
		if err := recover(); err != nil {
			handler(err)
		}
	}()
	fun()
}

// MilliSleep sleep tm milli second
func MilliSleep(tm int) {
	time.Sleep(time.Duration(tm) * time.Millisecond)
}

// Sleep time.Sleep tm second
func Sleep(tm int) {
	time.Sleep(time.Duration(tm) * time.Second)
}

// ToInterfaces convert []string to []interface{}
func ToInterfaces(fields []string) []interface{} {
	res := make([]interface{}, 0, len(fields))
	for _, s := range fields {
		res = append(res, s)
	}
	return res
}

// ToStrings convert []interface{} to []string
func ToStrings(fields []interface{}) []string {
	res := make([]string, 0, len(fields))
	for _, s := range fields {
		res = append(res, s.(string))
	}
	return res
}

// CmdCtrl If the operating system is macOS, return the key string "cmd",
// otherwise return the key string "ctrl"
func CmdCtrl() string {
	if runtime.GOOS == "darwin" {
		return "cmd"
	}
	return "ctrl" // Ctrl
}

// CmdV tap key command + v or control + v
func CmdV() error {
	return KeyTap("v", CmdCtrl())
}

// Scaled0 return int(x * f)
func Scaled0(x int, f float64) int {
	return int(float64(x) * f)
}

// Scaled1 return int(x / f)
func Scaled1(x int, f float64) int {
	return int(float64(x) / f)
}

// MoveArgs get the mouse relative args
func MoveArgs(x, y int) (int, int) {
	mx, my := Location()
	mx = mx + x
	my = my + y

	return mx, my
}

// MoveRelative move mouse with relative
func MoveRelative(x, y int) {
	Move(MoveArgs(x, y))
}

// MoveSmoothRelative move mouse smooth with relative
func MoveSmoothRelative(x, y int, args ...interface{}) {
	mx, my := MoveArgs(x, y)
	MoveSmooth(mx, my, args...)
}

// MovesClick move smooth and click the mouse
//
// use the `robotgo.MouseSleep = 100`
func MovesClick(x, y int, args ...interface{}) {
	MoveSmooth(x, y)
	MilliSleep(50)
	Click(args...)
}

// ScrollRelative scroll mouse with relative
//
// Examples:
//
//	robotgo.ScrollRelative(10, 10)
func ScrollRelative(x, y int, args ...int) {
	mx, my := MoveArgs(x, y)
	Scroll(mx, my, args...)
}

// CharCodeAt char code at utf-8
func CharCodeAt(s string, n int) rune {
	i := 0
	for _, r := range s {
		if i == n {
			return r
		}
		i++
	}

	return 0
}

// ToUC trans string to unicode []string
func ToUC(text string) []string {
	var uc []string

	for _, r := range text {
		textQ := strconv.QuoteToASCII(string(r))
		textUnQ := textQ[1 : len(textQ)-1]

		st := strings.Replace(textUnQ, "\\u", "U", -1)
		if st == "\\\\" {
			st = "\\"
		}
		if st == `\"` {
			st = `"`
		}
		uc = append(uc, st)
	}

	return uc
}

// ReadAll read string from clipboard
func ReadAll() (string, error) {
	return clipboard.ReadAll()
}

// WriteAll write string to clipboard
func WriteAll(text string) error {
	return clipboard.WriteAll(text)
}

// PasteStr paste a string
//
// Deprecated: use the Paste()
func PasteStr(str string) error {
	return Paste(str)
}

// Paste paste a string (supported UTF-8),
// write the string to clipboard and tap `cmd + v`
func Paste(str string) error {
	err := clipboard.WriteAll(str)
	if err != nil {
		return err
	}
	return CmdV()
}

// TypeStrDelay type string width delay
//
// Deprecated: use the TypeDelay()
func TypeStrDelay(str string, delay int) {
	TypeDelay(str, delay)
}

// TypeDelay type string with delayed
// And you can use robotgo.KeySleep = 100 to delayed not this function
func TypeDelay(str string, delay int) {
	Type(str)
	MilliSleep(delay)
}

// SetDelay sets the key and mouse delay
// robotgo.SetDelay(100) option the robotgo.KeySleep and robotgo.MouseSleep = d
func SetDelay(d ...int) {
	v := 10
	if len(d) > 0 {
		v = d[0]
	}

	KeySleep = v
	MouseSleep = v
}
