#Robotgo
  
  >Golang Desktop Automation. Control the mouse, keyboard, and read the screen.
  
RobotGo supports Mac, Windows, and Linux.

This is a work in progress.

##Installation:
  go get github.com/go-vgo/robotgo

  It's that easy!


##Examples:

###Mouse

```Go
package main

import (
	. "fmt"

	"github.com/go-vgo/robotgo"
)

func main() {
  robotgo.ScrollMouse(10, "up")
} 
``` 
