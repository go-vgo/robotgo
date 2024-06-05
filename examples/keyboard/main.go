package main

import (
	"flag"
	"fmt"
	keyboard "github.com/go-vgo/robotgo"
	"os"
	"strings"
)

func main() {
	// Define a flag to accept the key combination from the command line
	keysPtr := flag.String("keys", "", "Key combination to be pressed, e.g., 'cmd+shift+F'")
	flag.Parse()

	// Check if a key combination was provided
	if *keysPtr == "" {
		fmt.Println("Please provide a key combination using the -keys flag.")
		os.Exit(1)
	}

	keys := strings.Split(*keysPtr, "^")
	fmt.Printf("keys: %v\n", strings.Join(keys, "^"))
	keyboard.KeyTap(keys[len(keys)-1], keys[:len(keys)-1])
	fmt.Printf("Key combination '%s' pressed.\n", *keysPtr)
}
