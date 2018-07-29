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

	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/ewmh"
)

var xu *xgbutil.XUtil

// ActivePIDC active the window by PID,
// If args[0] > 0 on the unix platform via a xid to active
func ActivePIDC(pid int32, args ...int) {
	var hwnd int
	if len(args) > 0 {
		hwnd = args[0]

		internalActive(pid, hwnd)
		return
	}

	if xu == nil {
		var err error
		xu, err = xgbutil.NewConn()
		if err != nil {
			log.Println("xgbutil.NewConn errors is: ", err)
			return
		}
	}

	xid, err := getXidFromPid(xu, pid)
	if err != nil {
		log.Println("getXidFromPid errors is: ", err)
		return
	}

	internalActive(int32(xid), hwnd)
}

// ActivePID active the window by PID,
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

	xid, err := getXidFromPid(xu, pid)
	if err != nil {
		return err
	}

	err = ewmh.ActiveWindowReq(xu, xid)
	if err != nil {
		return err
	}

	return nil
}

func getXidFromPid(xu *xgbutil.XUtil, pid int32) (xproto.Window, error) {
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

	return 0, errors.New("failed to find a window with a matching pid")
}
