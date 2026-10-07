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
	"fmt"
	"math"
	"sync"
	"unsafe"

	"github.com/tailscale/win"
	"golang.org/x/sys/windows"
)

// HID-level input goes through the Interception kernel driver
// (https://github.com/oblitum/Interception): strokes are injected below the
// Win32 input stack, so they look like real device input to apps, games and
// raw-input consumers. The driver and interception.dll must be installed.

// ErrDriverNotInstalled is returned when the Interception driver is missing.
var ErrDriverNotInstalled = errors.New("robotgo: interception driver not installed or accessible")

// Interception device slots: keyboards are 1..10, mice 11..20.
const (
	hidMaxKeyboard = 10
	hidMaxMouse    = 10
	hidFirstMouse  = hidMaxKeyboard + 1
)

// Interception mouse states and flags (interception.h).
const (
	hidMouseLeftDown   = 0x001
	hidMouseLeftUp     = 0x002
	hidMouseRightDown  = 0x004
	hidMouseRightUp    = 0x008
	hidMouseMiddleDown = 0x010
	hidMouseMiddleUp   = 0x020
	hidMouseWheel      = 0x400
	hidMouseHWheel     = 0x800

	hidMoveAbsolute       = 0x001
	hidMoveVirtualDesktop = 0x002
)

// Interception key states (interception.h).
const (
	hidKeyDown = 0x00
	hidKeyUp   = 0x01
	hidKeyE0   = 0x02
	hidKeyE1   = 0x04
)

// Scan codes MapVirtualKey reports wrongly for real keyboards: Pause is
// E1 1D 45, Print Screen is E0 37 (MapVirtualKey gives SysRq 0x54).
const (
	hidScanCtrl     = 0x1D
	hidScanNumLock  = 0x45
	hidScanPrtSc    = 0x37
	hidMaxWheelStep = math.MaxInt16 / wheelDelta
)

// hidMouseStroke mirrors InterceptionMouseStroke (20 bytes).
type hidMouseStroke struct {
	State       uint16
	Flags       uint16
	Rolling     int16
	_           uint16
	X, Y        int32
	Information uint32
}

// hidKeyStroke mirrors InterceptionKeyStroke (8 bytes).
type hidKeyStroke struct {
	Code        uint16
	State       uint16
	Information uint32
}

var hid struct {
	mu       sync.Mutex
	path     string
	dll      *windows.LazyDLL
	send     *windows.LazyProc
	ctx      uintptr
	keyboard uintptr
	mouse    uintptr
}

// Seams replaced in tests.
var (
	hidSendMouse = func(s *hidMouseStroke) error { return hidSend(false, unsafe.Pointer(s)) }
	hidSendKey   = func(s *hidKeyStroke) error { return hidSend(true, unsafe.Pointer(s)) }
	hidOpen      = openHID
)

// SetHIDLibraryPath sets the interception.dll path used by the next HID
// initialization (default "interception.dll" on the DLL search path).
func SetHIDLibraryPath(path string) {
	hid.mu.Lock()
	hid.path = path
	hid.mu.Unlock()
}

// InitHID loads interception.dll and opens a driver context. It is called
// by SetBackend(BackendHID) and lazily by HID input.
func InitHID() error {
	hid.mu.Lock()
	defer hid.mu.Unlock()
	if hid.ctx != 0 {
		return nil
	}
	return hidOpen()
}

// CloseHID releases the Interception driver context.
func CloseHID() error {
	hid.mu.Lock()
	defer hid.mu.Unlock()
	if hid.ctx == 0 {
		return nil
	}
	destroy := hid.dll.NewProc("interception_destroy_context")
	if err := destroy.Find(); err != nil {
		return err
	}
	// interception_destroy_context returns void; Call's error is a stale LastError.
	destroy.Call(hid.ctx)
	hid.ctx = 0
	return nil
}

// openHID loads the DLL, creates a context and picks the first keyboard and
// mouse slot that reports a hardware id. Caller holds hid.mu.
func openHID() error {
	path := hid.path
	if path == "" {
		path = "interception.dll"
	}
	dll := windows.NewLazyDLL(path)
	if err := dll.Load(); err != nil {
		return fmt.Errorf("%w: %v", ErrDriverNotInstalled, err)
	}
	create := dll.NewProc("interception_create_context")
	send := dll.NewProc("interception_send")
	hwid := dll.NewProc("interception_get_hardware_id")
	for _, p := range []*windows.LazyProc{create, send, hwid} {
		if err := p.Find(); err != nil {
			return fmt.Errorf("%w: %v", ErrDriverNotInstalled, err)
		}
	}
	ctx, _, _ := create.Call()
	if ctx == 0 {
		return ErrDriverNotInstalled
	}

	present := func(dev uintptr) bool {
		var buf [512]byte
		n, _, _ := hwid.Call(ctx, dev, uintptr(unsafe.Pointer(&buf[0])), uintptr(len(buf)))
		return n > 0
	}
	hid.keyboard, hid.mouse = 1, hidFirstMouse
	for d := uintptr(1); d <= hidMaxKeyboard; d++ {
		if present(d) {
			hid.keyboard = d
			break
		}
	}
	for d := uintptr(hidFirstMouse); d < hidFirstMouse+hidMaxMouse; d++ {
		if present(d) {
			hid.mouse = d
			break
		}
	}
	hid.dll, hid.send, hid.ctx = dll, send, ctx
	return nil
}

func hidSend(keyboard bool, stroke unsafe.Pointer) error {
	hid.mu.Lock()
	defer hid.mu.Unlock()
	if hid.ctx == 0 {
		if err := hidOpen(); err != nil {
			return err
		}
	}
	dev := hid.mouse
	if keyboard {
		dev = hid.keyboard
	}
	// interception_send returns the number of strokes sent and sets no
	// LastError, so a zero result is reported with the device instead.
	if n, _, _ := hid.send.Call(hid.ctx, dev, uintptr(stroke), 1); n == 0 {
		return fmt.Errorf("robotgo: interception_send to device %d failed", dev)
	}
	return nil
}

// virtualScreen returns the virtual-desktop origin and size in pixels.
var virtualScreen = func() (x, y, w, h int) {
	return int(win.GetSystemMetrics(win.SM_XVIRTUALSCREEN)),
		int(win.GetSystemMetrics(win.SM_YVIRTUALSCREEN)),
		int(win.GetSystemMetrics(win.SM_CXVIRTUALSCREEN)),
		int(win.GetSystemMetrics(win.SM_CYVIRTUALSCREEN))
}

// normalize maps a screen coordinate to the 0..65535 absolute range.
func normalize(v, origin, size int) int32 {
	if size <= 1 {
		return 0
	}
	return int32((v - origin) * 65535 / (size - 1))
}

// hidMove moves the pointer to screen point (x, y) with one absolute stroke.
func hidMove(x, y int) error {
	vx, vy, vw, vh := virtualScreen()
	return hidSendMouse(&hidMouseStroke{
		Flags: hidMoveAbsolute | hidMoveVirtualDesktop,
		X:     normalize(x, vx, vw),
		Y:     normalize(y, vy, vh),
	})
}

// hidButtonState returns the Interception (down, up) states for a button.
func hidButtonState(btn string) (down, up uint16, err error) {
	switch btn {
	case "", "left":
		return hidMouseLeftDown, hidMouseLeftUp, nil
	case "right":
		return hidMouseRightDown, hidMouseRightUp, nil
	case "center", "middle":
		return hidMouseMiddleDown, hidMouseMiddleUp, nil
	}
	return 0, 0, errors.New("robotgo: unknown mouse button: " + btn)
}

func hidButton(btn string, up bool) error {
	down, release, err := hidButtonState(btn)
	if err != nil {
		return err
	}
	if up {
		down = release
	}
	return hidSendMouse(&hidMouseStroke{State: down})
}

// hidWheel scrolls by notches; positive y scrolls up, positive x scrolls left.
func hidWheel(x, y int) error {
	if err := hidRoll(hidMouseWheel, y); err != nil {
		return err
	}
	// HWHEEL is right-positive; robotgo is left-positive.
	return hidRoll(hidMouseHWheel, -x)
}

// hidRoll sends notches in strokes small enough for the int16 Rolling field.
func hidRoll(state uint16, notches int) error {
	for notches != 0 {
		step := max(-hidMaxWheelStep, min(hidMaxWheelStep, notches))
		if err := hidSendMouse(&hidMouseStroke{State: state, Rolling: int16(step * wheelDelta)}); err != nil {
			return err
		}
		notches -= step
	}
	return nil
}

// hidKey sends a virtual key as its hardware scan code sequence.
func hidKey(vk uint16, up bool) error {
	strokes, err := hidKeyStrokes(vk, up)
	if err != nil {
		return err
	}
	for i := range strokes {
		if err := hidSendKey(&strokes[i]); err != nil {
			return err
		}
	}
	return nil
}

// hidKeyStrokes returns the strokes a physical keyboard emits for vk.
func hidKeyStrokes(vk uint16, up bool) ([]hidKeyStroke, error) {
	state := uint16(hidKeyDown)
	if up {
		state = hidKeyUp
	}
	switch vk {
	case win.VK_PAUSE:
		return []hidKeyStroke{
			{Code: hidScanCtrl, State: state | hidKeyE1},
			{Code: hidScanNumLock, State: state},
		}, nil
	case win.VK_SNAPSHOT:
		return []hidKeyStroke{{Code: hidScanPrtSc, State: state | hidKeyE0}}, nil
	case win.VK_NUMLOCK:
		// Extended for SendInput, but the hardware scan code has no E0 prefix.
		return []hidKeyStroke{{Code: hidScanNumLock, State: state}}, nil
	}
	sc := scanCode(vk)
	if sc == 0 {
		return nil, ErrNotSupported
	}
	if extendedVKs[vk] {
		state |= hidKeyE0
	}
	return []hidKeyStroke{{Code: sc, State: state}}, nil
}
