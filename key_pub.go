// Copyright (c) 2016-2026 AtomAI, All rights reserved.
//
// See COPYRIGHT file at top-level directory of this distribution and at
// https://github.com/go-vgo/robotgo/blob/master/LICENSE
//
// Licensed under Apache License, Version 2.0 <LICENSE-APACHE or
// http://www.apache.org/licenses/LICENSE-2.0>
//
// This file may not be copied, modified, or distributed
// except according to those terms.

package robotgo

import (
	"strconv"
	"strings"
	"unicode"
)

// KeyPress press key string
//
// It dispatches to the backend keyPress (Cgo: down, short delay, up;
// pure-Go: the backend's atomic KeyTap) so each keeps its timing and locking.
func KeyPress(key string, args ...interface{}) error {
	return keyPress(key, args...)
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

// TypeStrDelay type string width delay
//
// Deprecated: use the TypeDelay()
func TypeStrDelay(str string, delay int) error {
	return TypeDelay(str, delay)
}

// TypeDelay type string with delayed
// And you can use robotgo.KeySleep = 100 to delayed not this function
func TypeDelay(str string, delay int) error {
	err := TypeStr(str)
	MilliSleep(delay)
	return err
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
//
// Runes outside ASCII become the Xlib keysym name ("U4e16", "U1F600").
func ToUC(text string) []string {
	var uc []string

	for _, r := range text {
		textQ := strconv.QuoteToASCII(string(r))
		textUnQ := textQ[1 : len(textQ)-1]

		// QuoteToASCII spells BMP runes as \uXXXX and the rest as \UXXXXXXXX;
		// Xlib wants a plain U prefix for both.
		st := strings.Replace(textUnQ, "\\u", "U", -1)
		st = strings.Replace(st, "\\U", "U", -1)
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

// getToggleArgs splits args into the pid (the first int, at any position, so
// KeyUp("a", pid) keeps it after the prepended "up") and the key array;
// string and []string args are flattened in order, other types are skipped.
func getToggleArgs(args ...any) (pid int, keyArr []string) {
	hasPid := false
	for _, arg := range args {
		switch v := arg.(type) {
		case int:
			if !hasPid {
				pid, hasPid = v, true
			}
		case string:
			keyArr = append(keyArr, v)
		case []string:
			keyArr = append(keyArr, v...)
		}
	}
	return
}

// getKeyDown strips a leading "up"/"down" direction from keyArr and reports
// whether it is a key down (the default).
func getKeyDown(keyArr []string) (bool, []string) {
	if len(keyArr) == 0 {
		return true, keyArr
	}
	switch keyArr[0] {
	case "up":
		return false, keyArr[1:]
	case "down":
		return true, keyArr[1:]
	}
	return true, keyArr
}

// getDown returns "down" or "up".
func getDown(down bool) string {
	if down {
		return "down"
	}
	return "up"
}

// appendShift lowers k and appends "shift" to args for a single upper-case
// char, or for a shifted symbol ("!" -> "1") when args has at most n
// entries (n counts a leading direction arg).
func appendShift(k string, n int, args ...any) (string, []any) {
	// only a single upper-case char implies shift; "Enter" is just "enter"
	if r := []rune(k); len(r) == 1 && unicode.IsUpper(r[0]) {
		args = append(args, "shift")
	}

	k = strings.ToLower(k)
	if v, ok := Special[k]; ok {
		k = v
		if len(args) <= n {
			args = append(args, "shift")
		}
	}

	return k, args
}
