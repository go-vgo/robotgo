package hook

/*
#cgo darwin CFLAGS: -x objective-c  -Wno-deprecated-declarations
#cgo darwin LDFLAGS: -framework Cocoa 
#cgo linux CFLAGS:-I/usr/src
#cgo linux LDFLAGS: -L/usr/src -lX11 -lXtst
#cgo linux LDFLAGS: -lX11-xcb -lxcb -lxcb-xkb -lxkbcommon -lxkbcommon-x11
#cgo windows LDFLAGS: -lgdi32 -luser32

#include "event/goEvent.h"
// #include "event/hook_async.h"
*/
import "C"

import(
	// 	"fmt"
	"unsafe"
)

// AddEvent add event listener
func AddEvent(key string) int {
	cs := C.CString(key)
	
	eve := C.add_event(cs)
	geve := int(eve)

	defer C.free(unsafe.Pointer(cs))
	return geve
}

// StopEvent stop event listener
func StopEvent() {
	C.stop_event()
}