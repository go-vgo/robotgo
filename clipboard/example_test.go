package clipboard_test

import (
	"fmt"

	"github.com/go-vgo/robotgo/clipboard"
	// "github.com/atotto/clipboard"
)

func Example() {
	clipboard.WriteAll("日本語")
	text, err := clipboard.ReadAll()
	if err == nil {
		fmt.Println(text)
	}

	// Output:
	// 日本語
}
