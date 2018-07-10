// +build darwin windows

package robotgo

// ActivePID active the window by PID,
// If args[0] > 0 on the Windows platform via a window handle to active
func ActivePID(pid int32, args ...int) error {
	var hwnd int
	if len(args) > 0 {
		hwnd = args[0]
	}

	internalActive(pid, hwnd)
	return nil
}
