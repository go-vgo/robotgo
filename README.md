#Robotgo
<!--<img align="right" src="https://raw.githubusercontent.com/go-vgo/robotgo/master/logo.jpg">-->
<!--[![Build Status](https://travis-ci.org/go-vgo/robotgo.svg)](https://travis-ci.org/go-vgo/robotgo)
[![codecov](https://codecov.io/gh/go-vgo/robotgo/branch/master/graph/badge.svg)](https://codecov.io/gh/go-vgo/robotgo)-->
<!--<a href="https://circleci.com/gh/go-vgo/robotgo/tree/dev"><img src="https://img.shields.io/circleci/project/go-vgo/robotgo/dev.svg" alt="Build Status"></a>-->
[![CircleCI Status](https://circleci.com/gh/go-vgo/robotgo.svg?style=shield)](https://circleci.com/gh/go-vgo/robotgo)
[![Go Report Card](https://goreportcard.com/badge/github.com/go-vgo/robotgo)](https://goreportcard.com/report/github.com/go-vgo/robotgo)
[![GoDoc](https://godoc.org/github.com/go-vgo/robotgo?status.svg)](https://godoc.org/github.com/go-vgo/robotgo)
[![Release](https://github-release-version.herokuapp.com/github/go-vgo/robotgo/release.svg?style=flat)](https://github.com/go-vgo/robotgo/releases/latest)
[![Join the chat at https://gitter.im/go-vgo/robotgo](https://badges.gitter.im/Join%20Chat.svg)](https://gitter.im/go-vgo/robotgo?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)
<!--<a href="https://github.com/go-vgo/robotgo/releases"><img src="https://img.shields.io/badge/%20version%20-%206.0.0%20-blue.svg?style=flat-square" alt="Releases"></a>-->
  
  >Golang Desktop Automation. Control the mouse, keyboard, bitmap, read the screen,   Window Handle and global event listener.
  
RobotGo supports Mac, Windows, and Linux(X11).

This is a work in progress.

[简体中文](https://github.com/go-vgo/robotgo/blob/master/README_zh.md)

##Contents
- [Docs](#docs)
- [Requirements](#requirements)
- [Installation](#installation)
- [Update](#update)
- [Examples](#examples)
- [Future](#future)
- [Contributors](#contributors)

##Docs
  - [API Document](https://github.com/go-vgo/robotgo/blob/master/docs/doc.md) &nbsp;&nbsp;&nbsp;
  - [中文文档](https://github.com/go-vgo/robotgo/blob/master/docs/doc_zh.md)
  - [GoDoc](https://godoc.org/github.com/go-vgo/robotgo)

##Requirements:

Now, Please make sure Golang, GCC, zlib and libpng is installed correctly before installing RobotGo.

###ALL: 
``` 
Golang
GCC
zlib & libpng (bitmap)
```
####For Mac OS X:
    Xcode Command Line Tools
```
brew install libpng
brew install homebrew/dupes/zlib
```    
####For Windows:
```
MinGW or other GCC

zlib & libpng (bitmap need it.)
```
#####[zlib&libpng Windows32 GCC's Course](https://github.com/go-vgo/Mingw32)

#####[Download include zlib&libpng Windows64 GCC](https://github.com/go-vgo/Mingw)

####For everything else:
```
GCC
    
X11 with the XTest extension (also known as the Xtst library)

Event:
    
xcb,xkb,libxkbcommon
``` 

#####Ubuntu:
```yml
sudo apt-get install gcc libc6-dev

sudo apt-get install libx11-dev
#sudo apt-get install libgtkglextmm-x11-dev
#sudo apt-get install libghc6-x11-dev
#sudo apt-get install libgl1-mesa-swx11-dev
sudo apt-get install xorg-dev

sudo apt-get install libxtst-dev libpng++-dev   

#Event:

sudo apt-get install xcb libxcb-xkb-dev x11-xkb-utils libx11-xcb-dev libxkbcommon-x11-dev
sudo apt-get install libxkbcommon-dev

```

##Installation:
```
go get github.com/go-vgo/robotgo
```
  It's that easy!

png.h: No such file or directory? Please see [issues/47](https://github.com/go-vgo/robotgo/issues/47).

##Update:
```
go get -u github.com/go-vgo/robotgo  
```

##[Examples:](https://github.com/go-vgo/robotgo/blob/master/example/main.go)

###[Mouse](https://github.com/go-vgo/robotgo/blob/master/example/main.go#L45)

```Go
package main

import (
	"github.com/go-vgo/robotgo"
)

func main() {
  robotgo.ScrollMouse(10, "up")
  robotgo.MouseClick("left", true)
  robotgo.MoveMouseSmooth(100, 200, 1.0, 100.0)
} 
``` 

###[Keyboard](https://github.com/go-vgo/robotgo/blob/master/example/main.go#L22)

```Go
package main

import (
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

###[Screen](https://github.com/go-vgo/robotgo/blob/master/example/main.go#L71)

```Go
package main

import (
	"fmt"

	"github.com/go-vgo/robotgo"
)

func main() {
  x, y := robotgo.GetMousePos()
  fmt.Println("pos:", x, y)
  color := robotgo.GetPixelColor(100, 200)
  fmt.Println("color----", color)
} 
```

###[Bitmap](https://github.com/go-vgo/robotgo/blob/master/example/main.go#L90)

```Go
package main

import (
	"fmt"

	"github.com/go-vgo/robotgo"
)

func main() {
  bitmap := robotgo.CaptureScreen(10, 20, 30, 40)
  fmt.Println("...", bitmap)

  fx, fy := robotgo.FindBitmap(bitmap)
  fmt.Println("FindBitmap------", fx, fy)

  robotgo.SaveBitmap(bitmap, "test.png")
} 
```

###[Event](https://github.com/go-vgo/robotgo/blob/master/example/main.go#L124)

```Go
package main

import (
	"fmt"

	"github.com/go-vgo/robotgo"
)

func main() {
  keve := robotgo.AddEvent("k")
  if keve == 0 {
    fmt.Println("you press...", "k")
  }

  mleft := robotgo.AddEvent("mleft")
  if mleft == 0 {
    fmt.Println("you press...", "mouse left button")
  }
} 
```

###[Window](https://github.com/go-vgo/robotgo/blob/master/example/main.go#L160)

```Go
package main

import (
	"fmt"

	"github.com/go-vgo/robotgo"
)

func main() {
  abool := robotgo.ShowAlert("test", "robotgo")
  if abool == 0 {
 	  fmt.Println("ok@@@", "ok")
  }

  title := robotgo.GetTitle()
  fmt.Println("title@@@", title)
} 
```

##Future
- Update Find an image on screen, read pixels from an image
- Update Window Handle
- Try support Android, maybe support IOS
- Remove zlib/libpng dependencies

##Contributors

- See [contributors page](https://github.com/go-vgo/robotgo/graphs/contributors) for full list of contributors.
- See [Contribution Guidelines](https://github.com/go-vgo/robotgo/blob/master/CONTRIBUTING.md).
