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

package wayland

import (
	"encoding/binary"
	"io"
	"net"
	"path/filepath"
	"testing"
	"time"

	"github.com/vcaesar/go-wayland/client"

	"github.com/go-vgo/robotgo/wayland/internal/protocols/wlr_virtual_keyboard"
)

func TestSendKeyPairedModifiers(t *testing.T) {
	path := filepath.Join(t.TempDir(), "wayland")
	listener, err := net.ListenUnix("unix", &net.UnixAddr{Name: path, Net: "unix"})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := listener.Close(); err != nil {
			t.Error(err)
		}
	})
	display, err := client.Connect(path)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := display.Context().Close(); err != nil {
			t.Error(err)
		}
	})
	peer, err := listener.AcceptUnix()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := peer.Close(); err != nil {
			t.Error(err)
		}
	})
	if err := peer.SetDeadline(time.Now().Add(5 * time.Second)); err != nil {
		t.Fatal(err)
	}
	c := &conn{keyboard: wlr_virtual_keyboard.NewZwpVirtualKeyboardV1(display.Context())}
	for _, pair := range []struct {
		name              string
		left, right, mask uint32
	}{
		{"shift", 42, 54, modShift},
		{"control", 29, 97, modControl},
		{"alt", 56, 100, modAlt},
		{"super", 125, 126, modSuper},
	} {
		t.Run(pair.name, func(t *testing.T) {
			for _, event := range []struct{ code, state, want uint32 }{
				{pair.left, keyStatePressed, pair.mask},
				{pair.left, keyStatePressed, pair.mask},
				{pair.right, keyStatePressed, pair.mask},
				{pair.left, keyStateReleased, pair.mask},
				{pair.left, keyStateReleased, pair.mask},
				{pair.right, keyStateReleased, 0},
			} {
				if err := c.sendKey(123, event.code, event.state); err != nil {
					t.Fatal(err)
				}
				var wire [44]byte // key request followed by modifiers request
				if _, err := io.ReadFull(peer, wire[:]); err != nil {
					t.Fatal(err)
				}
				if got := binary.NativeEndian.Uint32(wire[28:32]); got != event.want {
					t.Fatalf("code %d state %d: modifiers = %#x, want %#x", event.code, event.state, got, event.want)
				}
			}
		})
	}
}
