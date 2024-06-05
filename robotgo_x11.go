// Copyright 2016 The go-vgo Project Developers. See the COPYRIGHT
// file at the top-level directory of this distribution and at
// https://github.com/go-vgo/robotgo/blob/master/LICENSE
//
// Licensed under the Apache License, Version 2.0 <LICENSE-APACHE or
// http://www.apache.org/licenses/LICENSE-2.0> or the MIT license
// <LICENSE-MIT or http://opensource.org/licenses/MIT>, at your
// option. This file may not be copied, modified, or distributed
// except according to those terms.

//go:build !darwin && !windows
// +build !darwin,!windows

package keyboard

import (
	"errors"
	"log"

	"github.com/robotn/xgb"
	"github.com/robotn/xgb/xinerama"
	"github.com/robotn/xgb/xproto"
	"github.com/robotn/xgbutil"
	"github.com/robotn/xgbutil/ewmh"
)

var xu *xgbutil.XUtil

// GetBounds get the window bounds
func GetBounds(pid int, args ...int) (int, int, int, int) {
	var isPid int
	if len(args) > 0 || NotPid {
		isPid = 1
		return internalGetBounds(pid, isPid)
	}

	xid, err := GetXid(xu, pid)
	if err != nil {
		log.Println("Get Xid from Pid errors is: ", err)
		return 0, 0, 0, 0
	}

	return internalGetBounds(int(xid), isPid)
}

// GetClient get the window client bounds
func GetClient(pid int, args ...int) (int, int, int, int) {
	var isPid int
	if len(args) > 0 || NotPid {
		isPid = 1
		return internalGetClient(pid, isPid)
	}

	xid, err := GetXid(xu, pid)
	if err != nil {
		log.Println("Get Xid from Pid errors is: ", err)
		return 0, 0, 0, 0
	}

	return internalGetClient(int(xid), isPid)
}

// internalGetTitle get the window title
func internalGetTitle(pid int, args ...int) string {
	var isPid int
	if len(args) > 0 || NotPid {
		isPid = 1
		return cgetTitle(pid, isPid)
	}

	xid, err := GetXid(xu, pid)
	if err != nil {
		log.Println("Get Xid from Pid errors is: ", err)
		return ""
	}

	return cgetTitle(int(xid), isPid)
}

// ActivePidC active the window by Pid,
// If args[0] > 0 on the unix platform via a xid to active
func ActivePidC(pid int, args ...int) error {
	var isPid int
	if len(args) > 0 || NotPid {
		isPid = 1
		internalActive(pid, isPid)
		return nil
	}

	xid, err := GetXid(xu, pid)
	if err != nil {
		log.Println("Get Xid from Pid errors is: ", err)
		return err
	}

	internalActive(int(xid), isPid)
	return nil
}

// ActivePid active the window by Pid,
//
// If args[0] > 0 on the Windows platform via a window handle to active,
// If args[0] > 0 on the unix platform via a xid to active
func ActivePid(pid int, args ...int) error {
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

// GetXid get the xid return window and error
func GetXid(xu *xgbutil.XUtil, pid int) (xproto.Window, error) {
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

// GetXidFromPid get the xid from pid
func GetXidFromPid(xu *xgbutil.XUtil, pid int) (xproto.Window, error) {
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

// DisplaysNum get the count of displays
func DisplaysNum() int {
	c, err := xgb.NewConn()
	if err != nil {
		return 0
	}
	defer c.Close()

	err = xinerama.Init(c)
	if err != nil {
		return 0
	}

	reply, err := xinerama.QueryScreens(c).Reply()
	if err != nil {
		return 0
	}

	return int(reply.Number)
}

// GetMainId get the main display id
func GetMainId() int {
	conn, err := xgb.NewConn()
	if err != nil {
		return -1
	}

	setup := xproto.Setup(conn)
	defaultScreen := setup.DefaultScreen(conn)
	id := -1
	for i, screen := range setup.Roots {
		if defaultScreen.Root == screen.Root {
			id = i
			break
		}
	}
	return id
}

// Alert show a alert window
// Displays alert with the attributes.
// If cancel button is not given, only the default button is displayed
//
// Examples:
//
//	robotgo.Alert("hi", "window", "ok", "cancel")
func Alert(title, msg string, args ...string) bool {
	defaultBtn, cancelBtn := alertArgs(args...)
	c := `xmessage -center ` + msg +
		` -title ` + title + ` -buttons ` + defaultBtn + ":0,"
	if cancelBtn != "" {
		c += cancelBtn + ":1"
	}
	c += ` -default ` + defaultBtn
	c += ` -geometry 400x200`

	out, err := Run(c)
	if err != nil {
		// fmt.Println("Alert: ", err, ". ", string(out))
		return false
	}

	if string(out) == "1" {
		return false
	}
	return true
}
