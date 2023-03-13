// Copyright 2016 The go-vgo Project Developers. See the COPYRIGHT
// file at the top-level directory of this distribution and at
// https://github.com/go-vgo/robotgo/blob/master/LICENSE
//
// Licensed under the Apache License, Version 2.0 <LICENSE-APACHE or
// http://www.apache.org/licenses/LICENSE-2.0> or the MIT license
// <LICENSE-MIT or http://opensource.org/licenses/MIT>, at your
// option. This file may not be copied, modified, or distributed
// except according to those terms.

package robotgo

/*
// #include "key/keycode.h"
#include "key/keypress_c.h"
*/
import "C"

import (
	"errors"
	"math/rand"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"unicode"
	"unsafe"

	"github.com/go-vgo/robotgo/clipboard"
)

// Defining a bunch of constants.
const (
	// KeyA define key "a"
	KeyA = "a"
	KeyB = "b"
	KeyC = "c"
	KeyD = "d"
	KeyE = "e"
	KeyF = "f"
	KeyG = "g"
	KeyH = "h"
	KeyI = "i"
	KeyJ = "j"
	KeyK = "k"
	KeyL = "l"
	KeyM = "m"
	KeyN = "n"
	KeyO = "o"
	KeyP = "p"
	KeyQ = "q"
	KeyR = "r"
	KeyS = "s"
	KeyT = "t"
	KeyU = "u"
	KeyV = "v"
	KeyW = "w"
	KeyX = "x"
	KeyY = "y"
	KeyZ = "z"
	//
	CapA = "A"
	CapB = "B"
	CapC = "C"
	CapD = "D"
	CapE = "E"
	CapF = "F"
	CapG = "G"
	CapH = "H"
	CapI = "I"
	CapJ = "J"
	CapK = "K"
	CapL = "L"
	CapM = "M"
	CapN = "N"
	CapO = "O"
	CapP = "P"
	CapQ = "Q"
	CapR = "R"
	CapS = "S"
	CapT = "T"
	CapU = "U"
	CapV = "V"
	CapW = "W"
	CapX = "X"
	CapY = "Y"
	CapZ = "Z"
	//
	Key0 = "0"
	Key1 = "1"
	Key2 = "2"
	Key3 = "3"
	Key4 = "4"
	Key5 = "5"
	Key6 = "6"
	Key7 = "7"
	Key8 = "8"
	Key9 = "9"

	// Backspace backspace key string
	Backspace = "backspace"
	Delete    = "delete"
	Enter     = "enter"
	Tab       = "tab"
	Esc       = "esc"
	Escape    = "escape"
	Up        = "up"    // Up arrow key
	Down      = "down"  // Down arrow key
	Right     = "right" // Right arrow key
	Left      = "left"  // Left arrow key
	Home      = "home"
	End       = "end"
	Pageup    = "pageup"
	Pagedown  = "pagedown"

	F1  = "f1"
	F2  = "f2"
	F3  = "f3"
	F4  = "f4"
	F5  = "f5"
	F6  = "f6"
	F7  = "f7"
	F8  = "f8"
	F9  = "f9"
	F10 = "f10"
	F11 = "f11"
	F12 = "f12"
	F13 = "f13"
	F14 = "f14"
	F15 = "f15"
	F16 = "f16"
	F17 = "f17"
	F18 = "f18"
	F19 = "f19"
	F20 = "f20"
	F21 = "f21"
	F22 = "f22"
	F23 = "f23"
	F24 = "f24"

	Cmd  = "cmd"  // is the "win" key for windows
	Lcmd = "lcmd" // left command
	Rcmd = "rcmd" // right command
	// "command"
	Alt     = "alt"
	Lalt    = "lalt" // left alt
	Ralt    = "ralt" // right alt
	Ctrl    = "ctrl"
	Lctrl   = "lctrl" // left ctrl
	Rctrl   = "rctrl" // right ctrl
	Control = "control"
	Shift   = "shift"
	Lshift  = "lshift" // left shift
	Rshift  = "rshift" // right shift
	// "right_shift"
	Capslock    = "capslock"
	Space       = "space"
	Print       = "print"
	Printscreen = "printscreen" // No Mac support
	Insert      = "insert"
	Menu        = "menu" // Windows only

	AudioMute    = "audio_mute"     // Mute the volume
	AudioVolDown = "audio_vol_down" // Lower the volume
	AudioVolUp   = "audio_vol_up"   // Increase the volume
	AudioPlay    = "audio_play"
	AudioStop    = "audio_stop"
	AudioPause   = "audio_pause"
	AudioPrev    = "audio_prev"    // Previous Track
	AudioNext    = "audio_next"    // Next Track
	AudioRewind  = "audio_rewind"  // Linux only
	AudioForward = "audio_forward" // Linux only
	AudioRepeat  = "audio_repeat"  //  Linux only
	AudioRandom  = "audio_random"  //  Linux only

	Num0    = "num0" // numpad 0
	Num1    = "num1"
	Num2    = "num2"
	Num3    = "num3"
	Num4    = "num4"
	Num5    = "num5"
	Num6    = "num6"
	Num7    = "num7"
	Num8    = "num8"
	Num9    = "num9"
	NumLock = "num_lock"

	NumDecimal = "num."
	NumPlus    = "num+"
	NumMinus   = "num-"
	NumMul     = "num*"
	NumDiv     = "num/"
	NumClear   = "num_clear"
	NumEnter   = "num_enter"
	NumEqual   = "num_equal"

	LightsMonUp     = "lights_mon_up"     // Turn up monitor brightness			No Windows support
	LightsMonDown   = "lights_mon_down"   // Turn down monitor brightness		No Windows support
	LightsKbdToggle = "lights_kbd_toggle" // Toggle keyboard backlight on/off		No Windows support
	LightsKbdUp     = "lights_kbd_up"     // Turn up keyboard backlight brightness	No Windows support
	LightsKbdDown   = "lights_kbd_down"
)

// keyNames define a map of key names to MMKeyCode
var keyNames = map[string]C.MMKeyCode{
	"backspace": C.K_BACKSPACE,
	"delete":    C.K_DELETE,
	"enter":     C.K_RETURN,
	"tab":       C.K_TAB,
	"esc":       C.K_ESCAPE,
	"escape":    C.K_ESCAPE,
	"up":        C.K_UP,
	"down":      C.K_DOWN,
	"right":     C.K_RIGHT,
	"left":      C.K_LEFT,
	"home":      C.K_HOME,
	"end":       C.K_END,
	"pageup":    C.K_PAGEUP,
	"pagedown":  C.K_PAGEDOWN,
	//
	"f1":  C.K_F1,
	"f2":  C.K_F2,
	"f3":  C.K_F3,
	"f4":  C.K_F4,
	"f5":  C.K_F5,
	"f6":  C.K_F6,
	"f7":  C.K_F7,
	"f8":  C.K_F8,
	"f9":  C.K_F9,
	"f10": C.K_F10,
	"f11": C.K_F11,
	"f12": C.K_F12,
	"f13": C.K_F13,
	"f14": C.K_F14,
	"f15": C.K_F15,
	"f16": C.K_F16,
	"f17": C.K_F17,
	"f18": C.K_F18,
	"f19": C.K_F19,
	"f20": C.K_F20,
	"f21": C.K_F21,
	"f22": C.K_F22,
	"f23": C.K_F23,
	"f24": C.K_F24,
	//
	"cmd":         C.K_META,
	"lcmd":        C.K_LMETA,
	"rcmd":        C.K_RMETA,
	"command":     C.K_META,
	"alt":         C.K_ALT,
	"lalt":        C.K_LALT,
	"ralt":        C.K_RALT,
	"ctrl":        C.K_CONTROL,
	"lctrl":       C.K_LCONTROL,
	"rctrl":       C.K_RCONTROL,
	"control":     C.K_CONTROL,
	"shift":       C.K_SHIFT,
	"lshift":      C.K_LSHIFT,
	"rshift":      C.K_RSHIFT,
	"right_shift": C.K_RSHIFT,
	"capslock":    C.K_CAPSLOCK,
	"space":       C.K_SPACE,
	"print":       C.K_PRINTSCREEN,
	"printscreen": C.K_PRINTSCREEN,
	"insert":      C.K_INSERT,
	"menu":        C.K_MENU,

	"audio_mute":     C.K_AUDIO_VOLUME_MUTE,
	"audio_vol_down": C.K_AUDIO_VOLUME_DOWN,
	"audio_vol_up":   C.K_AUDIO_VOLUME_UP,
	"audio_play":     C.K_AUDIO_PLAY,
	"audio_stop":     C.K_AUDIO_STOP,
	"audio_pause":    C.K_AUDIO_PAUSE,
	"audio_prev":     C.K_AUDIO_PREV,
	"audio_next":     C.K_AUDIO_NEXT,
	"audio_rewind":   C.K_AUDIO_REWIND,
	"audio_forward":  C.K_AUDIO_FORWARD,
	"audio_repeat":   C.K_AUDIO_REPEAT,
	"audio_random":   C.K_AUDIO_RANDOM,

	"num0":     C.K_NUMPAD_0,
	"num1":     C.K_NUMPAD_1,
	"num2":     C.K_NUMPAD_2,
	"num3":     C.K_NUMPAD_3,
	"num4":     C.K_NUMPAD_4,
	"num5":     C.K_NUMPAD_5,
	"num6":     C.K_NUMPAD_6,
	"num7":     C.K_NUMPAD_7,
	"num8":     C.K_NUMPAD_8,
	"num9":     C.K_NUMPAD_9,
	"num_lock": C.K_NUMPAD_LOCK,

	// todo: removed
	"numpad_0":    C.K_NUMPAD_0,
	"numpad_1":    C.K_NUMPAD_1,
	"numpad_2":    C.K_NUMPAD_2,
	"numpad_3":    C.K_NUMPAD_3,
	"numpad_4":    C.K_NUMPAD_4,
	"numpad_5":    C.K_NUMPAD_5,
	"numpad_6":    C.K_NUMPAD_6,
	"numpad_7":    C.K_NUMPAD_7,
	"numpad_8":    C.K_NUMPAD_8,
	"numpad_9":    C.K_NUMPAD_9,
	"numpad_lock": C.K_NUMPAD_LOCK,

	"num.":      C.K_NUMPAD_DECIMAL,
	"num+":      C.K_NUMPAD_PLUS,
	"num-":      C.K_NUMPAD_MINUS,
	"num*":      C.K_NUMPAD_MUL,
	"num/":      C.K_NUMPAD_DIV,
	"num_clear": C.K_NUMPAD_CLEAR,
	"num_enter": C.K_NUMPAD_ENTER,
	"num_equal": C.K_NUMPAD_EQUAL,

	"lights_mon_up":     C.K_LIGHTS_MON_UP,
	"lights_mon_down":   C.K_LIGHTS_MON_DOWN,
	"lights_kbd_toggle": C.K_LIGHTS_KBD_TOGGLE,
	"lights_kbd_up":     C.K_LIGHTS_KBD_UP,
	"lights_kbd_down":   C.K_LIGHTS_KBD_DOWN,

	// { NULL:              C.K_NOT_A_KEY }
}

// CmdCtrl If the operating system is macOS, return the key string "cmd",
// otherwise return the key string "ctrl
func CmdCtrl() string {
	if runtime.GOOS == "darwin" {
		return "cmd"
	}
	return "ctrl"
}

// It sends a key press and release to the active application
func tapKeyCode(code C.MMKeyCode, flags C.MMKeyFlags, pid C.uintptr) {
	C.toggleKeyCode(code, true, flags, pid)
	MilliSleep(3)
	C.toggleKeyCode(code, false, flags, pid)
}

var keyErr = errors.New("Invalid key flag specified.")

func checkKeyCodes(k string) (key C.MMKeyCode, err error) {
	if k == "" {
		return
	}

	if len(k) == 1 {
		val1 := C.CString(k)
		defer C.free(unsafe.Pointer(val1))

		key = C.keyCodeForChar(*val1)
		if key == C.K_NOT_A_KEY {
			err = keyErr
			return
		}
		return
	}

	if v, ok := keyNames[k]; ok {
		key = v
		if key == C.K_NOT_A_KEY {
			err = keyErr
			return
		}
	}
	return
}

func checkKeyFlags(f string) (flags C.MMKeyFlags) {
	m := map[string]C.MMKeyFlags{
		"alt":    C.MOD_ALT,
		"ralt":   C.MOD_ALT,
		"lalt":   C.MOD_ALT,
		"cmd":    C.MOD_META,
		"rcmd":   C.MOD_META,
		"lcmd":   C.MOD_META,
		"ctrl":   C.MOD_CONTROL,
		"rctrl":  C.MOD_CONTROL,
		"lctrl":  C.MOD_CONTROL,
		"shift":  C.MOD_SHIFT,
		"rshift": C.MOD_SHIFT,
		"lshift": C.MOD_SHIFT,
		"none":   C.MOD_NONE,
	}

	if v, ok := m[f]; ok {
		return v
	}
	return
}

func getFlagsFromValue(value []string) (flags C.MMKeyFlags) {
	if len(value) <= 0 {
		return
	}

	for i := 0; i < len(value); i++ {
		var f C.MMKeyFlags = C.MOD_NONE

		f = checkKeyFlags(value[i])
		flags = (C.MMKeyFlags)(flags | f)
	}

	return
}

func keyTaps(k string, keyArr []string, pid int) error {
	flags := getFlagsFromValue(keyArr)
	key, err := checkKeyCodes(k)
	if err != nil {
		return err
	}

	tapKeyCode(key, flags, C.uintptr(pid))
	MilliSleep(KeySleep)
	return nil
}

func keyToggles(k string, keyArr []string, pid int, args ...interface{}) error {
	if len(keyArr) <= 0 {
		keyArr = append(keyArr, "down")
	}

	down := true
	if keyArr[0] == "up" {
		down = false
	}

	if keyArr[0] == "up" || keyArr[0] == "down" {
		keyArr = keyArr[1:]
	}
	flags := getFlagsFromValue(keyArr)

	key, err := checkKeyCodes(k)
	if err != nil {
		return err
	}

	C.toggleKeyCode(key, C.bool(down), flags, C.uintptr(pid))
	if len(args) > 0 {
		MilliSleep(KeySleep)
	}
	return nil
}

/*
 __  ___  ___________    ____ .______     ______        ___      .______       _______
|  |/  / |   ____\   \  /   / |   _  \   /  __  \      /   \     |   _  \     |       \
|  '  /  |  |__   \   \/   /  |  |_)  | |  |  |  |    /  ^  \    |  |_)  |    |  .--.  |
|    <   |   __|   \_    _/   |   _  <  |  |  |  |   /  /_\  \   |      /     |  |  |  |
|  .  \  |  |____    |  |     |  |_)  | |  `--'  |  /  _____  \  |  |\  \----.|  '--'  |
|__|\__\ |_______|   |__|     |______/   \______/  /__/     \__\ | _| `._____||_______/

*/

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

// toErr it converts a C string to a Go error
func toErr(str *C.char) error {
	gstr := C.GoString(str)
	if gstr == "" {
		return nil
	}
	return errors.New(gstr)
}

// KeyTap taps the keyboard code;
//
// See keys supported:
//
//	https://github.com/go-vgo/robotgo/blob/master/docs/keys.md#keys
//
// Examples:
//
//	robotgo.KeySleep = 100 // 100 millisecond
//	robotgo.KeyTap("a")
//	robotgo.KeyTap("i", "alt", "command")
//
//	arr := []string{"alt", "command"}
//	robotgo.KeyTap("i", arr)
//
//	robotgo.KeyTap("k", pid int)
func KeyTap(key string, args ...interface{}) error {
	var keyArr []string

	if len(key) > 0 && unicode.IsUpper([]rune(key)[0]) {
		args = append(args, "shift")
	}

	key = strings.ToLower(key)
	if _, ok := Special[key]; ok {
		key = Special[key]
		if len(args) <= 0 {
			args = append(args, "shift")
		}
	}

	pid := 0
	if len(args) > 0 {
		if reflect.TypeOf(args[0]) == reflect.TypeOf(keyArr) {
			keyArr = args[0].([]string)
		} else {
			if reflect.TypeOf(args[0]) == reflect.TypeOf(pid) {
				pid = args[0].(int)
				keyArr = ToStrings(args[1:])
			} else {
				keyArr = ToStrings(args)
			}
		}
	}

	return keyTaps(key, keyArr, pid)
}

// KeyToggle toggles the keyboard, if there not have args default is "down"
//
// See keys:
//
//	https://github.com/go-vgo/robotgo/blob/master/docs/keys.md#keys
//
// Examples:
//
//	robotgo.KeyToggle("a")
//	robotgo.KeyToggle("a", "up")
//
//	robotgo.KeyToggle("a", "up", "alt", "cmd")
//	robotgo.KeyToggle("k", pid int)
func KeyToggle(key string, args ...interface{}) error {

	if len(key) > 0 && unicode.IsUpper([]rune(key)[0]) {
		args = append(args, "shift")
	}

	key = strings.ToLower(key)
	if _, ok := Special[key]; ok {
		key = Special[key]
		if len(args) <= 1 {
			args = append(args, "shift")
		}
	}

	pid := 0
	var keyArr []string
	if len(args) > 0 && reflect.TypeOf(args[0]) == reflect.TypeOf(pid) {
		pid = args[0].(int)
		keyArr = ToStrings(args[1:])
	} else {
		keyArr = ToStrings(args)
	}

	return keyToggles(key, keyArr, pid)
}

// KeyPress press key string
func KeyPress(key string, args ...interface{}) error {
	err := KeyDown(key, args...)
	if err != nil {
		return err
	}

	MilliSleep(1 + rand.Intn(3))
	return KeyUp(key, args...)
}

// KeyDown press down a key
func KeyDown(key string, args ...interface{}) error {
	return KeyToggle(key, args...)
}

// KeyUp press up a key
func KeyUp(key string, args ...interface{}) error {
	arr := []interface{}{"up"}
	arr = append(arr, args...)
	return KeyToggle(key, arr...)
}

// ReadAll read string from clipboard
func ReadAll() (string, error) {
	return clipboard.ReadAll()
}

// WriteAll write string to clipboard
func WriteAll(text string) error {
	return clipboard.WriteAll(text)
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

// UnicodeType tap the uint32 unicode
func UnicodeType(str uint32, args ...int) {
	cstr := C.uint(str)
	pid := 0
	if len(args) > 0 {
		pid = args[0]
	}

	isPid := 0
	if len(args) > 1 {
		isPid = args[1]
	}

	C.unicodeType(cstr, C.uintptr(pid), C.int8_t(isPid))
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

func inputUTF(str string) {
	cstr := C.CString(str)
	C.input_utf(cstr)

	C.free(unsafe.Pointer(cstr))
}

// TypeStr send a string (supported UTF-8)
//
// robotgo.TypeStr(string: "The string to send", int: pid, "milli_sleep time", "x11 option")
//
// Examples:
//
//	robotgo.TypeStr("abc@123, Hi galaxy, こんにちは")
//	robotgo.TypeStr("To be or not to be, this is questions.", pid int)
func TypeStr(str string, args ...int) {
	var tm, tm1 = 0, 7

	if len(args) > 1 {
		tm = args[1]
	}
	if len(args) > 2 {
		tm1 = args[2]
	}
	pid := 0
	if len(args) > 0 {
		pid = args[0]
	}

	if runtime.GOOS == "linux" {
		strUc := ToUC(str)
		for i := 0; i < len(strUc); i++ {
			ru := []rune(strUc[i])
			if len(ru) <= 1 {
				ustr := uint32(CharCodeAt(strUc[i], 0))
				UnicodeType(ustr, pid)
			} else {
				inputUTF(strUc[i])
				MilliSleep(tm1)
			}

			MilliSleep(tm)
		}
		return
	}

	for i := 0; i < len([]rune(str)); i++ {
		ustr := uint32(CharCodeAt(str, i))
		UnicodeType(ustr, pid)
		// if len(args) > 0 {
		MilliSleep(tm)
		// }
	}
	MilliSleep(KeySleep)
}

// PasteStr paste a string (support UTF-8),
// write the string to clipboard and tap `cmd + v`
func PasteStr(str string) error {
	err := clipboard.WriteAll(str)
	if err != nil {
		return err
	}

	if runtime.GOOS == "darwin" {
		return KeyTap("v", "command")
	}

	return KeyTap("v", "control")
}

// TypeStrDelay type string with delayed
// And you can use robotgo.KeySleep = 100 to delayed not this function
func TypeStrDelay(str string, delay int) {
	TypeStr(str)
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
