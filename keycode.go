// Copyright (c) 2016-2026 AtomAI, All rights reserved.
//
// See the COPYRIGHT file at the top-level directory of this distribution and at
// https://github.com/go-vgo/robotgo/blob/master/LICENSE
//
// Licensed under the Apache License, Version 2.0 <LICENSE-APACHE or
// http://www.apache.org/licenses/LICENSE-2.0>
//
// This file may not be copied, modified, or distributed
// except according to those terms.

package robotgo

import (
	"github.com/vcaesar/keycode"
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
	Key0      = "0"
	Key1      = "1"
	Key2      = "2"
	Key3      = "3"
	Key4      = "4"
	Key5      = "5"
	Key6      = "6"
	Key7      = "7"
	Key8      = "8"
	Key9      = "9"
	KeyGrave  = "`"
	KeyQuoter = '"'
	KeyQuote  = "'"

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
	CmdL = "cmdl" // left command
	CmdR = "cmdr" // right command
	// "command"
	Alt     = "alt"
	AltL    = "altl" // left alt
	AltR    = "altr" // right alt
	Ctrl    = "ctrl"
	CtrlL   = "ctrll" // left ctrl
	CtrlR   = "ctrlr" // right ctrl
	Control = "control"
	Shift   = "shift"
	ShiftL  = "shiftl" // left shift
	ShiftR  = "shiftr" // right shift
	// "right_shift"
	Capslock    = "capslock"
	Space       = "space"
	Print       = "print"
	Printscreen = "printscreen" // No Mac support
	Insert      = "insert"
	Menu        = "menu" // Windows only
	ScrollLock  = "scroll_lock"
	PauseBreak  = "pause_break"

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

type uMap map[string]uint16

// MouseMap robotgo hook mouse's code map
var MouseMap = keycode.MouseMap

const (
	// Mleft mouse left button
	Mleft      = "left"
	Mright     = "right"
	Center     = "center"
	WheelDown  = "wheelDown"
	WheelUp    = "wheelUp"
	WheelLeft  = "wheelLeft"
	WheelRight = "wheelRight"
)

// Keycode robotgo hook key's code map
var Keycode = keycode.Keycode

// Special is the special key map
var Special = keycode.Special
