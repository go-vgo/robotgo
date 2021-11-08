package main

import (
	"fmt"

	"github.com/go-vgo/robotgo"
)

func colorPicker() {
	// click the left mouse button to get the value
	m := robotgo.AddEvent("mleft")
	if m {
		x, y := robotgo.GetMousePos()
		fmt.Println("mouse pos: ", x, y)

		clo := robotgo.GetPixelColor(x, y)
		fmt.Println("color: #", clo)

		// clipboard
		s1 := fmt.Sprint(x, ", ", y) + ": " + "#" + clo
		err := robotgo.WriteAll(s1)
		if err != nil {
			fmt.Println("clipboard err: ", err)
		}
	}
}

func main() {
	fmt.Println("color picker: ")
	fmt.Println("click the left mouse button to get the value.")
	for {
		colorPicker()
	}
}
