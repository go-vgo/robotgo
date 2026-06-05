// Copyright (c) 2016-2025 AtomAI, All rights reserved.
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

/*
Keys are supported:
	"A-Z a-z 0-9"

	"backspace"
	"delete"
	"enter"
	"tab"
	"esc"
	"escape"
	"up"		Up arrow key
	"down"		Down arrow key
	"right"		Right arrow key
	"left"		Left arrow key
	"home"
	"end"
	"pageup"
	"pagedown"

	"f1"
	"f2"
	"f3"
	"f4"
	"f5"
	"f6"
	"f7"
	"f8"
	"f9"
	"f10"
	"f11"
	"f12"
	"f13"
	"f14"
	"f15"
	"f16"
	"f17"
	"f18"
	"f19"
	"f20"
	"f21"
	"f22"
	"f23"
	"f24"

	"cmd"		this is the "win" key for windows
	"lcmd"		left command
	"rcmd"		right command
	// "command"
	"alt"
	"lalt"		left alt
	"ralt"		right alt
	"ctrl"
	"lctrl"		left ctrl
	"rctrl"		right ctrl
	"control"
	"shift"
	"lshift"	left shift
	"rshift"	right shift
	// "right_shift"
	"capslock"
	"space"
	"print"
	"printscreen"      // No Mac support
	"insert"
	"menu"				Windows only

	"audio_mute"		Mute the volume
	"audio_vol_down"	Lower the volume
	"audio_vol_up"		Increase the volume
	"audio_play"
	"audio_stop"
	"audio_pause"
	"audio_prev"		Previous Track
	"audio_next"		Next Track
	"audio_rewind"      Linux only
	"audio_forward"     Linux only
	"audio_repeat"      Linux only
	"audio_random"      Linux only


	"num0"
	"num1"
	"num2"
	"num3"
	"num4"
	"num5"
	"num6"
	"num7"
	"num8"
	"num9"
	"num_lock"

	"scroll_lock"
	"pause_break"

	"num."
	"num+"
	"num-"
	"num*"
	"num/"
	"num_clear"
	"num_enter"
	"num_equal"

	"lights_mon_up"		 Turn up monitor brightness					No Windows support
	"lights_mon_down"	 Turn down monitor brightness				No Windows support
	"lights_kbd_toggle"	 Toggle keyboard backlight on/off			No Windows support
	"lights_kbd_up"		 Turn up keyboard backlight brightness		No Windows support
	"lights_kbd_down"	 Turn down keyboard backlight brightness	No Windows support
*/

/*
## Type Conversion

|     | type conversion	    |  func
|-----|---------------------|----------------------
|	*	| robotgo.Bitmap -> robotgo.CBitmap | robotgo.ToCBitmap()
|		| robotgo.Bitmap -> *image.RGBA | robotgo.ToRGBAGo()
|	*	| robotgo.CBitmap -> C.MMBitmapRef | robotgo.ToMMBitmapRef()
|		| robotgo.CBitmap -> robotgo.Bitmap | robotgo.ToBitmap()
|		| robotgo.CBitmap -> image.Image | robotgo.ToImage()
|		| robotgo.CBitmap -> *image.RGBA | robotgo.ToRGBA()
|	*	| C.MMBitmapRef -> robotgo.CBitmap | robotgo.CBitmap()
|	*	| image.Image -> robotgo.Bitmap | robotgo.ImgToBitmap()
|		| image.Image -> robotgo.CBitmap | robotgo.ImgToCBitmap()
|		| image.Image -> []byte | robotgo.ToByteImg()
|		| image.Image -> string | robotgo.ToStringImg()
|	*	| *image.RGBA -> robotgo.Bitmap | robotgo.RGBAToBitmap()
|	*	| []byte -> image.Image | robotgo.ByteToImg()
|		| []byte-> robotgo.CBitmap | robotgo.ByteToCBitmap()
|		| []byte -> string | string()
|	*	| string -> image.Image | robotgo.StrToImg()
|		| string -> byte | []byte()
*/
