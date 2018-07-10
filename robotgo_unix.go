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

// ActivePID makes the window of the PID the active window
func ActivePID(pid int32) error {
	if xu == nil {
		var err error
		xu, err = xgbutil.NewConn()
		if err != nil {
			return err
		}
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
