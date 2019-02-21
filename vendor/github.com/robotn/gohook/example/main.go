package main

import (
	"fmt"

	"github.com/robotn/gohook"
)

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

func base() {
	EvChan := hook.Start()
	defer hook.End()

	for ev := range EvChan {
		fmt.Println(ev)
	}
}

func main() {
	base()

	add()
}
