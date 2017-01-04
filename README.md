#Robotgo
<!--<img align="right" src="https://raw.githubusercontent.com/go-vgo/robotgo/master/logo.jpg">-->
<!--[![Build Status](https://travis-ci.org/go-vgo/robotgo.svg)](https://travis-ci.org/go-vgo/robotgo)
[![codecov](https://codecov.io/gh/go-vgo/robotgo/branch/master/graph/badge.svg)](https://codecov.io/gh/go-vgo/robotgo)-->
[![Go Report Card](https://goreportcard.com/badge/github.com/go-vgo/robotgo)](https://goreportcard.com/report/github.com/go-vgo/robotgo)
[![GoDoc](https://godoc.org/github.com/go-vgo/robotgo?status.svg)](https://godoc.org/github.com/go-vgo/robotgo)
[![Release](http://github-release-version.herokuapp.com/github/go-vgo/robotgo/release.svg?style=flat)](https://github.com/go-vgo/robotgo/releases/latest)
<!--<a href="https://github.com/go-vgo/robotgo/releases"><img src="https://img.shields.io/badge/%20version%20-%206.0.0%20-blue.svg?style=flat-square" alt="Releases"></a>-->
<!--[![Join the chat at https://gitter.im/go-vgo/robotgo](https://badges.gitter.im/Join%20Chat.svg)](https://gitter.im/go-vgo/robotgo?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)-->
  
  >Golang Desktop Automation. Control the mouse, keyboard, bitmap,read the screen,Window Handle and global event listener.
  
RobotGo supports Mac, Windows, and Linux(X11).

This is a work in progress.



##[API Document](https://github.com/go-vgo/robotgo/blob/master/doc.md) &nbsp;&nbsp;&nbsp;[中文文档](https://github.com/go-vgo/robotgo/blob/master/doc_zh.md)
  Please click API Document;This is a work in progress.



##Installation:
    go get github.com/go-vgo/robotgo

  It's that easy!

##Update:
    go get -u github.com/go-vgo/robotgo  

###Requirements:

####ALL:  
    Golang
    //Gcc
    zlib & libpng (bitmap)

####For Mac OS X:
    Xcode Command Line Tools

    brew install libpng
    brew install homebrew/dupes/zlib
####For Windows:
    MinGW or other GCC

#####[zlib&libpng Windows32 GCC's Course](https://github.com/go-vgo/Mingw32)

#####[Download include zlib&libpng Windows64 GCC](https://github.com/go-vgo/Mingw)

####For everything else:
    GCC
    
    X11 with the XTest extension (also known as the Xtst library)

    Event:
    
    xcb,xkb,libxkbcommon

    

#####Ubuntu:

      sudo apt-get install libx11-dev
      sudo apt-get install libgtkglextmm-x11-dev
      sudo apt-get install libghc6-x11-dev
      sudo apt-get install libgl1-mesa-swx11-dev
      sudo apt-get install xorg-dev

      sudo apt-get install libxtst-dev libpng++-dev   

      Event:

      sudo apt-get install xcb libxcb-xkb-dev x11-xkb-utils libx11-xcb-dev libxkbcommon-x11-dev
      sudo apt-get install libxkbcommon-dev


##[Examples:](https://github.com/go-vgo/robotgo/blob/master/example/main.go)

###Mouse

```Go
package main

import (
	//. "fmt"

	"github.com/go-vgo/robotgo"
)

func main() {
  robotgo.ScrollMouse(10, "up")
  robotgo.MouseClick("left",true)
  robotgo.MoveMouseSmooth(100, 200, 1.0, 100.0)
} 
``` 

###Keyboard

```Go
package main

import (
	//. "fmt"

	"github.com/go-vgo/robotgo"
)

func main() {
  robotgo.TypeString("Hello World")
  robotgo.KeyTap("enter")
  robotgo.TypeString("en")
  robotgo.KeyTap("i", "alt", "command")
  arr := []string{"alt", "command"}
  robotgo.KeyTap("i", arr)
} 
```

###Screen

```Go
package main

import (
	. "fmt"

	"github.com/go-vgo/robotgo"
)

func main() {
  x, y := robotgo.GetMousePos()
  Println("pos:", x, y)
  color := robotgo.GetPixelColor(100, 200)
  Println("color----", color)
} 
```

###Bitmap

```Go
package main

import (
	. "fmt"

	"github.com/go-vgo/robotgo"
)

func main() {
  bit_map := robotgo.CaptureScreen(10, 20, 30, 40)
  Println("...", bit_map)

  fx, fy := robotgo.FindBitmap(bit_map)
  Println("FindBitmap------", fx, fy)

  robotgo.SaveBitmap(bit_map, "test.png")
} 
```

###Event

```Go
package main

import (
	. "fmt"

	"github.com/go-vgo/robotgo"
)

func main() {
  keve := robotgo.AddEvent("k")
  if keve == 0 {
    Println("you press...", "k")
  }

  mleft := robotgo.AddEvent("mleft")
  if mleft == 0 {
    Println("you press...", "mouse left button")
  }
} 
```

###Window

```Go
package main

import (
	. "fmt"

	"github.com/go-vgo/robotgo"
)

func main() {
  //Println("test")
  abool := robotgo.ShowAlert("test", "robotgo")
  if abool == 0 {
 	  Println("ok@@@", "ok")
  }

  title:=robotgo.GetTitle()
  Println("title@@@", title)
} 
```

##Contributors

- See [contributors page](https://github.com/go-vgo/robotgo/graphs/contributors) for full list of contributors.
- See [Contribution Guidelines](https://github.com/go-vgo/robotgo/blob/master/CONTRIBUTING.md).
