// +build windows

package robotgo

import (
	"syscall"

	"github.com/lxn/win"
)

// FindWindow find window hwnd by name
func FindWindow(str string) win.HWND {
	hwnd := win.FindWindow(nil, syscall.StringToUTF16Ptr(str))

	return hwnd
}
