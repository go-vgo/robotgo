package main

import (
	"io"
	"os"

	"github.com/go-vgo/robotgo/clipboard"
)

func main() {
	out, err := io.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}

	if err := clipboard.WriteAll(string(out)); err != nil {
		panic(err)
	}
}
