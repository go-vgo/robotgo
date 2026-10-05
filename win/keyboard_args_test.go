//go:build windows
// +build windows

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

package win

import (
	"reflect"
	"testing"

	"github.com/tailscale/win"
)

func TestExtractModifiersSlice(t *testing.T) {
	tests := []struct {
		name string
		args []interface{}
		want []string
	}{
		{"slice", []interface{}{[]string{"ctrl", "shift"}}, []string{"ctrl", "shift"}},
		{"mixed", []interface{}{42, "alt", []string{"CTRL", "hello", "SHIFT"}}, []string{"alt", "ctrl", "shift"}},
		{"command alias", []interface{}{[]string{"command", "cmdr", "rwin"}}, []string{"command", "cmdr", "rwin"}},
		{"right_shift alias", []interface{}{"Right_Shift"}, []string{"right_shift"}},
		{"direction only", []interface{}{[]string{"up", "down"}}, nil},
		{"empty slice", []interface{}{[]string{}}, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := extractModifiers(tt.args); !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("extractModifiers(%v) = %v, want %v", tt.args, got, tt.want)
			}
		})
	}
	if vk, _, ok := keyToVK("command"); !ok || vk != win.VK_LWIN {
		t.Fatalf("command alias = %#x, %v; want VK_LWIN", vk, ok)
	}
}

func TestToggleKeys(t *testing.T) {
	tests := []struct {
		name string
		key  string
		args []interface{}
		vks  []uint16
		up   bool
	}{
		{"plain down", "f1", nil, []uint16{win.VK_F1}, false},
		{"plain up", "f1", []interface{}{"up"}, []uint16{win.VK_F1}, true},
		{"slice mods", "f1", []interface{}{"down", []string{"ctrl", "shift"}}, []uint16{win.VK_CONTROL, win.VK_SHIFT, win.VK_F1}, false},
		{"direction in slice with pid", "f1", []interface{}{123, []string{"up", "ctrl"}}, []uint16{win.VK_CONTROL, win.VK_F1}, true},
		{"last direction wins", "f1", []interface{}{"up", "down"}, []uint16{win.VK_F1}, false},
		{"aliases dedup", "f1", []interface{}{[]string{"ctrl", "control", "CTRL", "cmd", "command"}}, []uint16{win.VK_CONTROL, win.VK_LWIN, win.VK_F1}, false},
		{"key is modifier", "ctrl", []interface{}{"ctrl"}, []uint16{win.VK_CONTROL}, false},
		{"side-specific", "f1", []interface{}{"ctrll", "ctrlr"}, []uint16{win.VK_LCONTROL, win.VK_RCONTROL, win.VK_F1}, false},
		{"implicit shift", "A", []interface{}{"up"}, []uint16{win.VK_SHIFT, 0x41}, true},
		{"explicit side shift kept", "A", []interface{}{"shiftl"}, []uint16{win.VK_LSHIFT, 0x41}, false},
		{"right_shift satisfies implied", "A", []interface{}{"right_shift"}, []uint16{win.VK_RSHIFT, 0x41}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vks, up, err := toggleKeys(tt.key, tt.args)
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(vks, tt.vks) || up != tt.up {
				t.Fatalf("toggleKeys = %#v, up=%v; want %#v, up=%v", vks, up, tt.vks, tt.up)
			}
		})
	}
}

func TestToggleKeysRejectsUnknownKey(t *testing.T) {
	if vks, _, err := toggleKeys("nonexistent_key", []interface{}{"ctrl"}); err == nil || vks != nil {
		t.Fatalf("unknown key: vks %v, err %v", vks, err)
	}
	// These reject the key before any input backend can emit events.
	for _, fn := range []func(string, ...interface{}) error{KeyTap, KeyToggle, KeyDown, KeyUp} {
		if err := fn("nonexistent_key", []string{"ctrl", "shift"}); err == nil {
			t.Fatal("expected unknown-key error")
		}
	}
}
