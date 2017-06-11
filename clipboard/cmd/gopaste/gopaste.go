package main

import (
	"fmt"

	"github.com/go-vgo/robotgo/clipboard"
	// "github.com/atotto/clipboard"
)

func main() {
	text, err := clipboard.ReadAll()
	if err != nil {
		panic(err)
	}

	fmt.Print(text)
}
