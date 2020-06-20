package main

import (
	"fmt"

	hook "github.com/robotn/gohook"
)

// hook listen and return values using detailed examples
func add() {
	s := hook.Start()
	defer hook.End()

	ct := false
	for {
		i := <-s

		if i.Kind == hook.KeyHold && i.Rawcode == 59 {
			ct = true
		}

		if ct && i.Rawcode == 12 {
			break
		}
	}
}

// base hook example
func base() {
	evChan := hook.Start()
	defer hook.End()

	for ev := range evChan {
		fmt.Println("hook: ", ev)
	}
}

func main() {
	base()

	add()
}
