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
// raw-input consumers. Only the driver must be installed: robotgo talks to
// its \\.\interceptionNN devices with DeviceIoControl, without interception.dll.

// ErrDriverNotInstalled is returned when the Interception driver is missing.
var ErrDriverNotInstalled = errors.New("robotgo: interception driver not installed or accessible")

// Interception device slots: keyboards are 1..10, mice 11..20.
const (
	hidMaxKeyboard = 10
	hidMaxMouse    = 10
	hidFirstMouse  = hidMaxKeyboard + 1
)

// Interception driver IOCTLs:
// CTL_CODE(FILE_DEVICE_UNKNOWN, fn, METHOD_BUFFERED, FILE_ANY_ACCESS).
const (
	ioctlSetEvent      = 0x222040 // fn 0x810
	ioctlWrite         = 0x222080 // fn 0x820
	ioctlGetHardwareID = 0x222200 // fn 0x880
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

// hidRawKey mirrors KEYBOARD_INPUT_DATA (ntddkbd.h), the driver's input.
type hidRawKey struct {
	UnitID, MakeCode, Flags, Reserved uint16
	ExtraInformation                  uint32
}

// hidRawMouse mirrors MOUSE_INPUT_DATA (ntddmou.h), the driver's input.
type hidRawMouse struct {
	UnitID, Flags, ButtonFlags, ButtonData uint16
	RawButtons                             uint32
	LastX, LastY                           int32
	ExtraInformation                       uint32
}

// raw converts s the way interception_send does.
func (s *hidKeyStroke) raw() hidRawKey {
	return hidRawKey{MakeCode: s.Code, Flags: s.State, ExtraInformation: s.Information}
}

// raw converts s the way interception_send does.
func (s *hidMouseStroke) raw() hidRawMouse {
	return hidRawMouse{
		Flags: s.Flags, ButtonFlags: s.State, ButtonData: uint16(s.Rolling),
		LastX: s.X, LastY: s.Y, ExtraInformation: s.Information,
	}
}

// hidDevice is an open \\.\interceptionNN handle and its registered event.
type hidDevice struct {
	slot          int
	handle, event windows.Handle
}

func (d hidDevice) close() error {
	return errors.Join(windows.CloseHandle(d.event), windows.CloseHandle(d.handle))
}

var hid struct {
	mu              sync.Mutex
	open            bool
	keyboard, mouse hidDevice
}

// Seams replaced in tests.
var (
	hidSendMouse = func(s *hidMouseStroke) error {
		r := s.raw()
		return hidWrite(false, unsafe.Pointer(&r), uint32(unsafe.Sizeof(r)))
	}
	hidSendKey = func(s *hidKeyStroke) error {
		r := s.raw()
		return hidWrite(true, unsafe.Pointer(&r), uint32(unsafe.Sizeof(r)))
	}
	hidOpen       = openHID
	hidCreateFile = func(name string) (windows.Handle, error) {
		p, err := windows.UTF16PtrFromString(name)
		if err != nil {
			return 0, err
		}
		return windows.CreateFile(p, windows.GENERIC_READ, 0, nil, windows.OPEN_EXISTING, 0, 0)
	}
	hidIoctl = func(h windows.Handle, code uint32, in unsafe.Pointer, inLen uint32,
		out unsafe.Pointer, outLen uint32) (uint32, error) {
		var n uint32
		err := windows.DeviceIoControl(h, code, (*byte)(in), inLen, (*byte)(out), outLen, &n, nil)
		return n, err
	}
)

// SetHIDLibraryPath is a no-op kept for compatibility.
//
// Deprecated: the HID backend talks to the Interception driver directly and
// no longer loads interception.dll.
func SetHIDLibraryPath(path string) {}

// InitHID opens the Interception driver devices. It is called by
// SetBackend(BackendHID) and lazily by HID input.
func InitHID() error {
	hid.mu.Lock()
	defer hid.mu.Unlock()
	if hid.open {
		return nil
	}
	return hidOpen()
}

// CloseHID releases the Interception driver devices.
func CloseHID() error {
	hid.mu.Lock()
	defer hid.mu.Unlock()
	if !hid.open {
		return nil
	}
	err := errors.Join(hid.keyboard.close(), hid.mouse.close())
	hid.open, hid.keyboard, hid.mouse = false, hidDevice{}, hidDevice{}
	return err
}

// openHID opens the first keyboard and mouse slot that reports a hardware
// id. Caller holds hid.mu.
func openHID() error {
	kb, err := openHIDDevice(1, hidMaxKeyboard)
	if err != nil {
		return err
	}
	ms, err := openHIDDevice(hidFirstMouse, hidMaxMouse)
	if err != nil {
		return errors.Join(err, kb.close())
	}
	hid.open, hid.keyboard, hid.mouse = true, kb, ms
	return nil
}

// openHIDDevice opens the first of n slots from first that reports a
// hardware id, or else the first slot that opens.
func openHIDDevice(first, n int) (hidDevice, error) {
	fallback := 0
	var openErr error
	for slot := first; slot < first+n; slot++ {
		d, err := openHIDSlot(slot)
		if err != nil {
			openErr = err
			continue
		}
		if hidPresent(d) {
			return d, nil
		}
		if fallback == 0 {
			fallback = slot
		}
		if err := d.close(); err != nil {
			return hidDevice{}, err
		}
	}
	if fallback == 0 {
		return hidDevice{}, openErr
	}
	return openHIDSlot(fallback)
}

// openHIDSlot opens device slot (1-based) and registers its input event,
// like interception_create_context.
func openHIDSlot(slot int) (hidDevice, error) {
	h, err := hidCreateFile(fmt.Sprintf(`\\.\interception%02d`, slot-1))
	if err != nil {
		return hidDevice{}, fmt.Errorf("%w: %v", ErrDriverNotInstalled, err)
	}
	ev, err := windows.CreateEvent(nil, 1, 0, nil)
	if err != nil {
		return hidDevice{}, errors.Join(err, windows.CloseHandle(h))
	}
	d := hidDevice{slot: slot, handle: h, event: ev}
	reg := [2]windows.Handle{ev}
	if _, err := hidIoctl(h, ioctlSetEvent, unsafe.Pointer(&reg), uint32(unsafe.Sizeof(reg)), nil, 0); err != nil {
		return hidDevice{}, errors.Join(fmt.Errorf("robotgo: interception set event on device %d: %w", slot, err), d.close())
	}
	return d, nil
}

// hidPresent reports whether a real device is attached to d's slot.
func hidPresent(d hidDevice) bool {
	var buf [512]byte
	n, err := hidIoctl(d.handle, ioctlGetHardwareID, nil, 0, unsafe.Pointer(&buf[0]), uint32(len(buf)))
	return err == nil && n > 0
}

// hidWrite sends one raw stroke of size bytes to the keyboard or mouse.
func hidWrite(keyboard bool, raw unsafe.Pointer, size uint32) error {
	hid.mu.Lock()
	defer hid.mu.Unlock()
	if !hid.open {
		if err := hidOpen(); err != nil {
			return err
		}
	}
	d := hid.mouse
	if keyboard {
		d = hid.keyboard
	}
	n, err := hidIoctl(d.handle, ioctlWrite, raw, size, nil, 0)
	if err != nil {
		return fmt.Errorf("robotgo: interception write to device %d: %w", d.slot, err)
	}
	// The driver reports the bytes it consumed; less than one stroke is a failure.
	if n < size {
		return fmt.Errorf("robotgo: interception write to device %d failed", d.slot)
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
