# gohook

This is a work in progress.

```Go
package main

import (
	"fmt"

	"go-vgo/robotn/gohook"
)

func main() {
	// hook.AsyncHook()
	veve := hook.AddEvent("v")
	if veve == 0 {
		fmt.Println("v...")
	}
}
```