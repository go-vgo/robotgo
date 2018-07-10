package ewmh

import (
	"fmt"

	"github.com/BurntSushi/xgbutil"
)

// GetEwmhWM uses the EWMH spec to find if a conforming window manager
// is currently running or not. If it is, then its name will be returned.
// Otherwise, an error will be returned explaining why one couldn't be found.
func GetEwmhWM(xu *xgbutil.XUtil) (string, error) {
	childCheck, err := SupportingWmCheckGet(xu, xu.RootWin())
	if err != nil {
		return "", fmt.Errorf("GetEwmhWM: Failed because: %s", err)
	}

	childCheck2, err := SupportingWmCheckGet(xu, childCheck)
	if err != nil {
		return "", fmt.Errorf("GetEwmhWM: Failed because: %s", err)
	}

	if childCheck != childCheck2 {
		return "", fmt.Errorf(
			"GetEwmhWM: _NET_SUPPORTING_WM_CHECK value on the root window "+
				"(%x) does not match _NET_SUPPORTING_WM_CHECK value "+
				"on the child window (%x).", childCheck, childCheck2)
	}

	return WmNameGet(xu, childCheck)
}
