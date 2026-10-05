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

package libei

import (
	"reflect"
	"testing"
)

func TestExtractModifiersSlice(t *testing.T) {
	tests := []struct {
		name string
		args []interface{}
		want []string
	}{
		{"slice", []interface{}{[]string{"ctrl", "shift"}}, []string{"ctrl", "shift"}},
		{"mixed", []interface{}{42, "alt", []string{"CTRL", "hello", "Shift"}}, []string{"alt", "ctrl", "shift"}},
		{"command alias", []interface{}{[]string{"command"}}, []string{"command"}},
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
}

func TestToggleKeys(t *testing.T) {
	const (
		f1    = 59
		ctrl  = 29
		shift = 42
		meta  = 125
		a     = 30
	)
	tests := []struct {
		name  string
		key   string
		args  []interface{}
		codes []int32
		up    bool
	}{
		{"plain down", "f1", nil, []int32{f1}, false},
		{"plain up", "f1", []interface{}{"up"}, []int32{f1}, true},
		{"slice mods", "f1", []interface{}{"down", []string{"ctrl", "shift"}}, []int32{ctrl, shift, f1}, false},
		{"direction in slice", "f1", []interface{}{1234, []string{"up", "ctrl"}}, []int32{ctrl, f1}, true},
		{"last direction wins", "f1", []interface{}{"up", "down"}, []int32{f1}, false},
		{"aliases dedup", "f1", []interface{}{[]string{"ctrl", "control", "CTRL", "cmd", "command"}}, []int32{ctrl, meta, f1}, false},
		{"key is modifier", "ctrl", []interface{}{"ctrl"}, []int32{ctrl}, false},
		{"implicit shift", "A", []interface{}{"up"}, []int32{shift, a}, true},
		{"explicit shift dedup", "A", []interface{}{"shift"}, []int32{shift, a}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			codes, up, err := toggleKeys(tt.key, tt.args)
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(codes, tt.codes) || up != tt.up {
				t.Fatalf("toggleKeys = %v, up=%v; want %v, up=%v", codes, up, tt.codes, tt.up)
			}
		})
	}

	if codes, _, err := toggleKeys("nonexistent_key", []interface{}{"ctrl"}); err == nil || codes != nil {
		t.Fatalf("unknown key: codes %v, err %v", codes, err)
	}
}
