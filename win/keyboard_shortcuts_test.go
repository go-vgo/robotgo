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

func TestShortcutModifierArguments(t *testing.T) {
	tests := []struct {
		name string
		args []interface{}
		want []string
	}{
		{"slice", []interface{}{[]string{"ctrl", "shift"}}, []string{"ctrl", "shift"}},
		{"mixed", []interface{}{42, "alt", []string{"CTRL", "hello", "SHIFT"}}, []string{"alt", "ctrl", "shift"}},
		{"command alias", []interface{}{[]string{"command", "cmdr", "rwin"}}, []string{"command", "cmdr", "rwin"}},
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

func TestHeldShortcutEvents(t *testing.T) {
	tests := []struct {
		name string
		key  string
		args []interface{}
		want []keyEvent
	}{
		{"plain down", "f1", nil, []keyEvent{{vk: win.VK_F1}}},
		{"plain up", "f1", []interface{}{"up"}, []keyEvent{{vk: win.VK_F1, up: true}}},
		{"modifiers down", "f1", []interface{}{"down", []string{"ctrl", "shift"}}, []keyEvent{
			{vk: win.VK_CONTROL}, {vk: win.VK_SHIFT}, {vk: win.VK_F1},
		}},
		{"modifiers up", "f1", []interface{}{"up", []string{"ctrl", "shift"}}, []keyEvent{
			{vk: win.VK_F1, up: true}, {vk: win.VK_SHIFT, up: true}, {vk: win.VK_CONTROL, up: true},
		}},
		{"direction in slice with pid", "f1", []interface{}{123, []string{"up", "ctrl"}}, []keyEvent{
			{vk: win.VK_F1, up: true}, {vk: win.VK_CONTROL, up: true},
		}},
		{"aliases deduplicated", "f1", []interface{}{[]string{"ctrl", "control", "CTRL", "cmd", "command"}}, []keyEvent{
			{vk: win.VK_CONTROL}, {vk: win.VK_LWIN}, {vk: win.VK_F1},
		}},
		{"key is modifier", "ctrl", []interface{}{"ctrl"}, []keyEvent{{vk: win.VK_CONTROL}}},
		{"side-specific modifiers", "f1", []interface{}{"up", "ctrll", "ctrlr"}, []keyEvent{
			{vk: win.VK_F1, up: true}, {vk: win.VK_RCONTROL, up: true}, {vk: win.VK_LCONTROL, up: true},
		}},
		{"implicit shift", "A", nil, []keyEvent{{vk: win.VK_SHIFT}, {vk: 0x41}}},
		{"explicit shift is preserved", "A", []interface{}{"shiftl"}, []keyEvent{{vk: win.VK_LSHIFT}, {vk: 0x41}}},
		{"implicit shift released", "A", []interface{}{"up"}, []keyEvent{{vk: 0x41, up: true}, {vk: win.VK_SHIFT, up: true}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := toggleKeyEvents(tt.key, tt.args)
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("events = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func TestHeldShortcutRejectsUnknownKey(t *testing.T) {
	if events, err := toggleKeyEvents("nonexistent_key", []interface{}{"ctrl"}); err == nil || len(events) != 0 {
		t.Fatalf("unknown key produced events %v, error %v", events, err)
	}
	// These calls reject the key before either input backend can emit events.
	for _, toggle := range []func(string, ...interface{}) error{KeyToggle, KeyDown, KeyUp} {
		if err := toggle("nonexistent_key", []string{"ctrl", "shift"}); err == nil {
			t.Fatal("expected unknown-key error")
		}
	}
}
