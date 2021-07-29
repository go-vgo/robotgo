package main

import (
	"log"

	"github.com/go-vgo/robotgo/clipboard"
)

func main() {
	err := clipboard.WriteAll("日本語")
	if err != nil {
		log.Println("clipboard write all error: ", err)
	}

	text, err := clipboard.ReadAll()
	if err != nil {
		log.Println("clipboard read all error: ", err)
		return
	}

	if text != "" {
		log.Println("text is: ", text)
		// Output: 日本語
	}
}
