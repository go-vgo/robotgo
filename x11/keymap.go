//go:build linux
// +build linux

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

package x11

import "github.com/jezek/xgb/xproto"

// X11 keysym values (from <X11/keysymdef.h>). Only the named/special keys are
// listed; printable ASCII / Latin-1 keysyms equal their Unicode code point and
// are handled by runeKeysym.
const (
	xkBackSpace = 0xff08
	xkTab       = 0xff09
	xkReturn    = 0xff0d
	xkEscape    = 0xff1b
	xkDelete    = 0xffff
	xkHome      = 0xff50
	xkLeft      = 0xff51
	xkUp        = 0xff52
	xkRight     = 0xff53
	xkDown      = 0xff54
	xkPageUp    = 0xff55
	xkPageDown  = 0xff56
	xkEnd       = 0xff57
	xkInsert    = 0xff63
	xkSpace     = 0x0020
	xkMenu      = 0xff67
	xkPrint     = 0xff61
	xkScrollLck = 0xff14
	xkPause     = 0xff13
	xkNumLock   = 0xff7f
	xkCapsLock  = 0xffe5

	// modifiers
	xkShiftL   = 0xffe1
	xkShiftR   = 0xffe2
	xkControlL = 0xffe3
	xkControlR = 0xffe4
	xkAltL     = 0xffe9
	xkAltR     = 0xffea
	xkSuperL   = 0xffeb
	xkSuperR   = 0xffec

	// keypad
	xkKP0        = 0xffb0
	xkKPDecimal  = 0xffae
	xkKPAdd      = 0xffab
	xkKPSubtract = 0xffad
	xkKPMultiply = 0xffaa
	xkKPDivide   = 0xffaf
	xkKPEnter    = 0xff8d
	xkKPEqual    = 0xffbd

	// function keys (F1 = xkF1, Fn = xkF1 + n - 1)
	xkF1 = 0xffbe

	// XF86 multimedia keys
	xkAudioLowerVolume = 0x1008ff11
	xkAudioMute        = 0x1008ff12
	xkAudioRaiseVolume = 0x1008ff13
	xkAudioPlay        = 0x1008ff14
	xkAudioStop        = 0x1008ff15
	xkAudioPrev        = 0x1008ff16
	xkAudioNext        = 0x1008ff17
	xkAudioRewind      = 0x1008ff3e
	xkAudioForward     = 0x1008ff97
	xkAudioRepeat      = 0x1008ff3d
	xkMonBrightnessUp  = 0x1008ff02
	xkMonBrightnessDn  = 0x1008ff03
	xkKbdLightOnOff    = 0x1008ff04
	xkKbdBrightnessUp  = 0x1008ff05
	xkKbdBrightnessDn  = 0x1008ff06
)

// specialKeysyms maps robotgo key names (see keycode.go in the root package) to
// X11 keysyms. Single-character keys are resolved by runeKeysym instead.
var specialKeysyms = map[string]uint32{
	"backspace":   xkBackSpace,
	"delete":      xkDelete,
	"enter":       xkReturn,
	"return":      xkReturn,
	"tab":         xkTab,
	"esc":         xkEscape,
	"escape":      xkEscape,
	"up":          xkUp,
	"down":        xkDown,
	"right":       xkRight,
	"left":        xkLeft,
	"home":        xkHome,
	"end":         xkEnd,
	"pageup":      xkPageUp,
	"pagedown":    xkPageDown,
	"space":       xkSpace,
	"capslock":    xkCapsLock,
	"print":       xkPrint,
	"printscreen": xkPrint,
	"insert":      xkInsert,
	"menu":        xkMenu,
	"scroll_lock": xkScrollLck,
	"pause_break": xkPause,

	// modifiers
	"cmd":     xkSuperL,
	"cmdl":    xkSuperL,
	"cmdr":    xkSuperR,
	"command": xkSuperL,
	"alt":     xkAltL,
	"altl":    xkAltL,
	"altr":    xkAltR,
	"ctrl":    xkControlL,
	"ctrll":   xkControlL,
	"ctrlr":   xkControlR,
	"control": xkControlL,
	"shift":   xkShiftL,
	"shiftl":  xkShiftL,
	"shiftr":  xkShiftR,

	// audio / media
	"audio_mute":     xkAudioMute,
	"audio_vol_down": xkAudioLowerVolume,
	"audio_vol_up":   xkAudioRaiseVolume,
	"audio_play":     xkAudioPlay,
	"audio_stop":     xkAudioStop,
	"audio_pause":    xkAudioPlay,
	"audio_prev":     xkAudioPrev,
	"audio_next":     xkAudioNext,
	"audio_rewind":   xkAudioRewind,
	"audio_forward":  xkAudioForward,
	"audio_repeat":   xkAudioRepeat,

	// brightness / backlight
	"lights_mon_up":     xkMonBrightnessUp,
	"lights_mon_down":   xkMonBrightnessDn,
	"lights_kbd_toggle": xkKbdLightOnOff,
	"lights_kbd_up":     xkKbdBrightnessUp,
	"lights_kbd_down":   xkKbdBrightnessDn,

	// numpad
	"num0":      xkKP0 + 0,
	"num1":      xkKP0 + 1,
	"num2":      xkKP0 + 2,
	"num3":      xkKP0 + 3,
	"num4":      xkKP0 + 4,
	"num5":      xkKP0 + 5,
	"num6":      xkKP0 + 6,
	"num7":      xkKP0 + 7,
	"num8":      xkKP0 + 8,
	"num9":      xkKP0 + 9,
	"num.":      xkKPDecimal,
	"num+":      xkKPAdd,
	"num-":      xkKPSubtract,
	"num*":      xkKPMultiply,
	"num/":      xkKPDivide,
	"num_enter": xkKPEnter,
	"num_equal": xkKPEqual,
	"num_lock":  xkNumLock,
}

// keyKeysym resolves a robotgo key name to an X11 keysym.
func keyKeysym(name string) (uint32, bool) {
	if name == "" {
		return 0, false
	}
	if ks, ok := specialKeysyms[name]; ok {
		return ks, true
	}
	// f1..f35
	if (name[0] == 'f' || name[0] == 'F') && len(name) <= 3 {
		if n := atoiSafe(name[1:]); n >= 1 && n <= 35 {
			return uint32(xkF1) + uint32(n-1), true
		}
	}
	// single rune (letters, digits, punctuation)
	r := []rune(name)
	if len(r) == 1 {
		return runeKeysym(r[0]), true
	}
	return 0, false
}

// runeKeysym maps a Unicode rune to its X11 keysym. Latin-1 (<= 0xff) maps
// directly; everything else uses the X11 Unicode keysym range.
func runeKeysym(r rune) uint32 {
	if r <= 0xff {
		return uint32(r)
	}
	return 0x01000000 | uint32(r)
}

// keysymToKeycode looks up the keycode (and whether Shift is required) that
// produces the given keysym in the cached keyboard mapping.
func (c *conn) keysymToKeycode(ks uint32) (xproto.Keycode, bool, bool) {
	per := int(c.keysymsPerKeycode)
	if per <= 0 {
		return 0, false, false
	}
	for i := 0; i*per+per <= len(c.keysyms); i++ {
		col0 := uint32(c.keysyms[i*per])
		if col0 == ks {
			return c.minKeycode + xproto.Keycode(i), false, true
		}
		if per > 1 && uint32(c.keysyms[i*per+1]) == ks {
			return c.minKeycode + xproto.Keycode(i), true, true
		}
	}
	return 0, false, false
}

// findScratchKeycode returns a keycode whose every column is NoSymbol (0); such
// a keycode is safe to remap temporarily for Unicode typing.
func (c *conn) findScratchKeycode() (xproto.Keycode, bool) {
	per := int(c.keysymsPerKeycode)
	if per <= 0 {
		return 0, false
	}
	for i := 0; i*per+per <= len(c.keysyms); i++ {
		empty := true
		for j := 0; j < per; j++ {
			if c.keysyms[i*per+j] != 0 {
				empty = false
				break
			}
		}
		if empty {
			return c.minKeycode + xproto.Keycode(i), true
		}
	}
	return 0, false
}

// atoiSafe parses a small non-negative integer, returning -1 on any error.
func atoiSafe(s string) int {
	if s == "" {
		return -1
	}
	n := 0
	for _, ch := range s {
		if ch < '0' || ch > '9' {
			return -1
		}
		n = n*10 + int(ch-'0')
	}
	return n
}
