# Robotgo
<!--<img align="right" src="https://raw.githubusercontent.com/go-vgo/robotgo/master/logo.jpg">-->
<!--[![Build Status](https://travis-ci.org/go-vgo/robotgo.svg)](https://travis-ci.org/go-vgo/robotgo)
[![codecov](https://codecov.io/gh/go-vgo/robotgo/branch/master/graph/badge.svg)](https://codecov.io/gh/go-vgo/robotgo)-->
<!--<a href="https://circleci.com/gh/go-vgo/robotgo/tree/dev"><img src="https://img.shields.io/circleci/project/go-vgo/robotgo/dev.svg" alt="Build Status"></a>-->
[![Build Status](https://travis-ci.org/go-vgo/robotgo.svg)](https://travis-ci.org/go-vgo/robotgo)
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

## Contents
- [Docs](#docs)
- [Requirements](#requirements)
- [Installation](#installation)
- [Update](#update)
- [Examples](#examples)
- [Plans](#plans)
- [Donate](#donate)
- [Contributors](#contributors)
- [License](#license)

## Docs
  - [API Docs](https://github.com/go-vgo/robotgo/blob/master/docs/doc.md) &nbsp;&nbsp;&nbsp;
  - [中文文档](https://github.com/go-vgo/robotgo/blob/master/docs/doc_zh.md)
  - [GoDoc](https://godoc.org/github.com/go-vgo/robotgo)

## Requirements:

Now, Please make sure Golang, GCC is installed correctly before installing RobotGo.

### ALL: 
``` 
Golang
GCC
```
#### For Mac OS X:
```
    Xcode Command Line Tools
```    
#### For Windows:
```
MinGW or other GCC

```

#### For everything else:
```
GCC
    
X11 with the XTest extension (also known as the Xtst library)

Event:
    
xcb, xkb, libxkbcommon
``` 

##### Ubuntu:
```yml
sudo apt-get install gcc libc6-dev

sudo apt-get install libx11-dev
sudo apt-get install xorg-dev  

sudo apt-get install xcb libxcb-xkb-dev x11-xkb-utils libx11-xcb-dev libxkbcommon-x11-dev
sudo apt-get install libxkbcommon-dev

sudo apt-get install xsel
sudo apt-get install xclip
```

#### Fedora:

```yml
sudo dnf install libxkbcommon-devel libXtst-devel libxkbcommon-x11-devel xorg-x11-xkb-utils-devel

sudo dnf install xsel
sudo dnf install xclip
```

## Installation:
```
go get github.com/go-vgo/robotgo
```
  It's that easy!

## Update:
```
go get -u github.com/go-vgo/robotgo  
```

## [Examples:](https://github.com/go-vgo/robotgo/blob/master/examples)

#### [Mouse](https://github.com/go-vgo/robotgo/blob/master/examples/mouse/main.go)

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

#### [Keyboard](https://github.com/go-vgo/robotgo/blob/master/examples/key/main.go)

```Go
package main

import (
  "fmt"

  "github.com/go-vgo/robotgo"
)

func main() {
  robotgo.TypeString("Hello World")
  robotgo.KeyTap("enter")
  robotgo.TypeString("en")
  robotgo.KeyTap("i", "alt", "command")
  arr := []string{"alt", "command"}
  robotgo.KeyTap("i", arr)

  robotgo.WriteAll("Test")
  text, err := robotgo.ReadAll()
  if err == nil {
    fmt.Println(text)
  }
} 
```

#### [Screen](https://github.com/go-vgo/robotgo/blob/master/examples/screen/main.go)

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

<!-- #### [Bitmap](https://github.com/go-vgo/robotgo/blob/master/examples/bitmap/mian.go)

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
``` -->

#### [Event](https://github.com/go-vgo/robotgo/blob/master/examples/event/main.go)

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

#### [Window](https://github.com/go-vgo/robotgo/blob/master/examples/window/main.go)

```Go
package main

import (
	"fmt"

	"github.com/go-vgo/robotgo"
)

func main() {
  fpid, err := robotgo.FindIds("Google")
  if err == nil {
    fmt.Println("pids...", fpid)
  }

  isExist, err := robotgo.PidExists(100)
  if err == nil {
    fmt.Println("pid exists is", isExist)
  }

  abool := robotgo.ShowAlert("test", "robotgo")
  if abool == 0 {
 	  fmt.Println("ok@@@", "ok")
  }

  title := robotgo.GetTitle()
  fmt.Println("title@@@", title)
} 
```

## Plans
- Update Find an image on screen, read pixels from an image
- Update Window Handle
- Support UTF-8, transitional plan: TypeStr
- Try support Android, maybe support IOS
- Remove zlib/libpng dependencies

## Donate

Supporting robotgo, [buy me a coffee](https://github.com/go-vgo/buy-me-a-coffee).

#### Paypal

Donate money by [paypal](https://www.paypal.me/veni0/25) to my account [vzvway@gmail.com](vzvway@gmail.com)


## Contributors

- See [contributors page](https://github.com/go-vgo/robotgo/graphs/contributors) for full list of contributors.
- See [Contribution Guidelines](https://github.com/go-vgo/robotgo/blob/master/CONTRIBUTING.md).

## License

Robotgo is primarily distributed under the terms of both the MIT license and the Apache License (Version 2.0), with portions covered by various BSD-like licenses.

See [LICENSE-APACHE](http://www.apache.org/licenses/LICENSE-2.0), [LICENSE-MIT](https://github.com/go-vgo/robotgo/blob/master/LICENSE).
