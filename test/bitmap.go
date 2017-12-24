package main

import (
	"fmt"
	"reflect"

	"github.com/go-vgo/robotgo"
)

func main() {

	var bitmap robotgo.Bitmap
	bit := robotgo.OpenBitmap("./test.bmp", 2)
	bitmap = robotgo.ToBitmap(bit)
	fmt.Println("...type", reflect.TypeOf(bit), reflect.TypeOf(bitmap))

}
