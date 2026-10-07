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
	"errors"
	"reflect"
	"testing"
	"unsafe"

	"golang.org/x/sys/windows"
)

type fakeIoctl struct {
	code  uint32
	slot  int
	bytes []byte
}

// fakeDriver emulates the \\.\interceptionNN devices. Slots in present
// report a hardware id; slots in missing fail to open. Handles are real
// events so CloseHandle succeeds.
func fakeDriver(t *testing.T, present, missing map[int]bool) *[]fakeIoctl {
	t.Helper()
	oldCreate, oldIoctl := hidCreateFile, hidIoctl
	slots := map[windows.Handle]int{}
	calls := &[]fakeIoctl{}
	t.Cleanup(func() {
		if err := CloseHID(); err != nil {
			t.Errorf("CloseHID: %v", err)
		}
		hidCreateFile, hidIoctl = oldCreate, oldIoctl
	})
	hidCreateFile = func(name string) (windows.Handle, error) {
		idx, ok := deviceIndex(name)
		if !ok {
			t.Fatalf("bad device name %q", name)
		}
		if missing[idx+1] {
			return windows.InvalidHandle, windows.ERROR_FILE_NOT_FOUND
		}
		h, err := windows.CreateEvent(nil, 0, 0, nil)
		if err != nil {
			return 0, err
		}
		slots[h] = idx + 1
		return h, nil
	}
	hidIoctl = func(h windows.Handle, code uint32, in unsafe.Pointer, inLen uint32,
		out unsafe.Pointer, outLen uint32) (uint32, error) {
		slot := slots[h]
		c := fakeIoctl{code: code, slot: slot}
		if in != nil {
			c.bytes = append([]byte(nil), unsafe.Slice((*byte)(in), inLen)...)
		}
		*calls = append(*calls, c)
		switch code {
		case ioctlGetHardwareID:
			if present[slot] {
				return 4, nil
			}
			return 0, nil
		case ioctlWrite:
			return inLen, nil
		}
		return 0, nil
	}
	return calls
}

// deviceIndex parses the NN of \\.\interceptionNN.
func deviceIndex(name string) (int, bool) {
	const prefix = `\\.\interception`
	if len(name) != len(prefix)+2 || name[:len(prefix)] != prefix {
		return 0, false
	}
	return int(name[len(prefix)]-'0')*10 + int(name[len(prefix)+1]-'0'), true
}

func lastWrite(t *testing.T, calls []fakeIoctl) fakeIoctl {
	t.Helper()
	for i := len(calls) - 1; i >= 0; i-- {
		if calls[i].code == ioctlWrite {
			return calls[i]
		}
	}
	t.Fatal("no IOCTL_WRITE")
	return fakeIoctl{}
}

func TestHIDDriverWrite(t *testing.T) {
	calls := fakeDriver(t, map[int]bool{3: true, 12: true}, nil)
	if err := InitHID(); err != nil {
		t.Fatal(err)
	}
	if hid.keyboard.slot != 3 || hid.mouse.slot != 12 {
		t.Fatalf("slots = %d/%d, want 3/12", hid.keyboard.slot, hid.mouse.slot)
	}

	if err := hidSendKey(&hidKeyStroke{Code: 0x1E, State: hidKeyUp | hidKeyE0, Information: 7}); err != nil {
		t.Fatal(err)
	}
	w := lastWrite(t, *calls)
	want := hidRawKey{MakeCode: 0x1E, Flags: hidKeyUp | hidKeyE0, ExtraInformation: 7}
	if w.slot != 3 || !reflect.DeepEqual(w.bytes, rawBytes(&want)) {
		t.Errorf("key write slot %d bytes %v, want slot 3 %v", w.slot, w.bytes, rawBytes(&want))
	}

	if err := hidSendMouse(&hidMouseStroke{State: hidMouseWheel, Flags: hidMoveAbsolute, Rolling: -120, X: 5, Y: -6}); err != nil {
		t.Fatal(err)
	}
	w = lastWrite(t, *calls)
	wantM := hidRawMouse{Flags: hidMoveAbsolute, ButtonFlags: hidMouseWheel, ButtonData: 0xFF88, LastX: 5, LastY: -6}
	if w.slot != 12 || !reflect.DeepEqual(w.bytes, rawBytes(&wantM)) {
		t.Errorf("mouse write slot %d bytes %v, want slot 12 %v", w.slot, w.bytes, rawBytes(&wantM))
	}

	// Every opened device registered an event first.
	if (*calls)[0].code != ioctlSetEvent || len((*calls)[0].bytes) != int(unsafe.Sizeof([2]windows.Handle{})) {
		t.Errorf("first ioctl = %+v, want IOCTL_SET_EVENT", (*calls)[0])
	}
	if err := CloseHID(); err != nil {
		t.Fatal(err)
	}
	if hid.open {
		t.Error("CloseHID left the driver open")
	}
}

func rawBytes[T any](v *T) []byte {
	return append([]byte(nil), unsafe.Slice((*byte)(unsafe.Pointer(v)), unsafe.Sizeof(*v))...)
}

func TestHIDDriverFallbackSlot(t *testing.T) {
	// No hardware ids: use the first slot that opens.
	fakeDriver(t, nil, map[int]bool{1: true, 11: true, 12: true})
	if err := InitHID(); err != nil {
		t.Fatal(err)
	}
	if hid.keyboard.slot != 2 || hid.mouse.slot != 13 {
		t.Errorf("slots = %d/%d, want 2/13", hid.keyboard.slot, hid.mouse.slot)
	}
}

func TestHIDDriverMissing(t *testing.T) {
	missing := map[int]bool{}
	for s := 1; s <= hidMaxKeyboard+hidMaxMouse; s++ {
		missing[s] = true
	}
	fakeDriver(t, nil, missing)
	if err := InitHID(); !errors.Is(err, ErrDriverNotInstalled) {
		t.Fatalf("InitHID = %v, want ErrDriverNotInstalled", err)
	}

	// Mice missing: the opened keyboard is released and init fails.
	for s := 1; s <= hidMaxKeyboard; s++ {
		delete(missing, s)
	}
	if err := InitHID(); !errors.Is(err, ErrDriverNotInstalled) {
		t.Fatalf("InitHID = %v, want ErrDriverNotInstalled", err)
	}
	if hid.open {
		t.Error("failed InitHID must not mark the driver open")
	}
}

func TestHIDDriverShortWrite(t *testing.T) {
	fakeDriver(t, nil, nil)
	ioctl := hidIoctl
	hidIoctl = func(h windows.Handle, code uint32, in unsafe.Pointer, inLen uint32,
		out unsafe.Pointer, outLen uint32) (uint32, error) {
		if code == ioctlWrite {
			return 0, nil
		}
		return ioctl(h, code, in, inLen, out, outLen)
	}
	if err := hidSendKey(&hidKeyStroke{Code: 0x1E}); err == nil {
		t.Error("zero-byte write should fail")
	}
}
