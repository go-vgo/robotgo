package main

import (
	"fmt"

	"github.com/go-vgo/robotgo"
)

func colorPicker() {
	m := robotgo.AddEvent("mleft")
	if m {
		x, y := robotgo.GetMousePos()
		fmt.Println("mouse pos: ", x, y)

		clo := robotgo.GetPixelColor(x, y)
		fmt.Println("color: #", clo)

		// clipboard
		robotgo.WriteAll("#" + clo)
	}
}

func main() {
	fmt.Println("color picker...")

	for {
		colorPicker()
	}
}
