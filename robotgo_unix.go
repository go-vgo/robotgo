// Copyright 2016 The go-vgo Project Developers. See the COPYRIGHT
// file at the top-level directory of this distribution and at
// https://github.com/go-vgo/robotgo/blob/master/LICENSE
//
// Licensed under the Apache License, Version 2.0 <LICENSE-APACHE or
// http://www.apache.org/licenses/LICENSE-2.0> or the MIT license
// <LICENSE-MIT or http://opensource.org/licenses/MIT>, at your
// option. This file may not be copied, modified, or distributed
// except according to those terms.

// +build !darwin,!windows

package robotgo

import (
	"errors"
	"log"

	"github.com/robotn/xgb/xproto"
	"github.com/robotn/xgbutil"
	"github.com/robotn/xgbutil/ewmh"
)

var xu *xgbutil.XUtil

// GetBounds get the window bounds
func GetBounds(pid int32, args ...int) (int, int, int, int) {
	var hwnd int
	if len(args) > 0 {
		hwnd = args[0]

		return internalGetBounds(pid, hwnd)
	}

	xid, err := GetXId(xu, pid)
	if err != nil {
		log.Println("GetXid from Pid errors is: ", err)
		return 0, 0, 0, 0
	}

	return internalGetBounds(int32(xid), hwnd)
}

// internalGetTitle get the window title
func internalGetTitle(pid int32, args ...int32) string {
	var hwnd int32
	if len(args) > 0 {
		hwnd = args[0]

		return cgetTitle(pid, hwnd)
	}

	xid, err := GetXId(xu, pid)
	if err != nil {
		log.Println("GetXid from Pid errors is: ", err)
		return ""
	}

	return cgetTitle(int32(xid), hwnd)
}

// ActivePIDC active the window by PID,
// If args[0] > 0 on the unix platform via a xid to active
func ActivePIDC(pid int32, args ...int) {
	var hwnd int
	if len(args) > 0 {
		hwnd = args[0]

		internalActive(pid, hwnd)
		return
	}

	xid, err := GetXId(xu, pid)
	if err != nil {
		log.Println("GetXid from Pid errors is: ", err)
		return
	}

	internalActive(int32(xid), hwnd)
}

// ActivePID active the window by PID,
//
// If args[0] > 0 on the Windows platform via a window handle to active,
// If args[0] > 0 on the unix platform via a xid to active
func ActivePID(pid int32, args ...int) error {
	if xu == nil {
		var err error
		xu, err = xgbutil.NewConn()
		if err != nil {
			return err
		}
	}

	if len(args) > 0 {
		err := ewmh.ActiveWindowReq(xu, xproto.Window(pid))
		if err != nil {
			return err
		}

		return nil
	}

	// get the xid from pid
	xid, err := GetXidFromPid(xu, pid)
	if err != nil {
		return err
	}

	err = ewmh.ActiveWindowReq(xu, xid)
	if err != nil {
		return err
	}

	return nil
}

// GetXId get the xid return window and error
func GetXId(xu *xgbutil.XUtil, pid int32) (xproto.Window, error) {
	if xu == nil {
		var err error
		xu, err = xgbutil.NewConn()
		if err != nil {
			// log.Println("xgbutil.NewConn errors is: ", err)
			return 0, err
		}
	}

	xid, err := GetXidFromPid(xu, pid)
	return xid, err
}

// GetXidFromPid get the xide from pid
func GetXidFromPid(xu *xgbutil.XUtil, pid int32) (xproto.Window, error) {
	windows, err := ewmh.ClientListGet(xu)
	if err != nil {
		return 0, err
	}

	for _, window := range windows {
		wmPid, err := ewmh.WmPidGet(xu, window)
		if err != nil {
			return 0, err
		}

		if uint(pid) == wmPid {
			return window, nil
		}
	}

	return 0, errors.New("failed to find a window with a matching pid.")
}
