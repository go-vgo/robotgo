#Robotgo
[![CircleCI Status](https://circleci.com/gh/go-vgo/robotgo.svg?style=shield)](https://circleci.com/gh/go-vgo/robotgo)
[![Go Report Card](https://goreportcard.com/badge/github.com/go-vgo/robotgo)](https://goreportcard.com/report/github.com/go-vgo/robotgo)
[![GoDoc](https://godoc.org/github.com/go-vgo/robotgo?status.svg)](https://godoc.org/github.com/go-vgo/robotgo)
[![Release](https://github-release-version.herokuapp.com/github/go-vgo/robotgo/release.svg?style=flat)](https://github.com/go-vgo/robotgo/releases/latest)
[![Join the chat at https://gitter.im/go-vgo/robotgo](https://badges.gitter.im/Join%20Chat.svg)](https://gitter.im/go-vgo/robotgo?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)
  
  >Golang 跨平台自动系统，控制键盘鼠标位图和读取屏幕,窗口句柄以及全局事件监听
  
RobotGo 支持 Mac, Windows, and Linux(X11).

这是一项正在完善中的工作.

提Issues请到[Github](https://github.com/go-vgo/robotgo),便于统一管理和即时更新


##[中文文档](https://github.com/go-vgo/robotgo/blob/master/docs/doc_zh.md)&nbsp;&nbsp;&nbsp;[API Document](https://github.com/go-vgo/robotgo/blob/master/docs/doc.md) 
 英文文档请点击API Document.

- [Requirements](#requirements)
- [Installation](#installation)
- [Update](#update)
- [Examples](#examples)
- [Future](#future)
- [Contributors](#contributors)

###Requirements:
(环境要求)

####ALL:  
```
Golang
//Gcc
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
```
#####[zlib&libpng Windows32 GCC教程](https://github.com/go-vgo/Mingw32)
#####[下载包含zlib和libpng的64位MinGW](https://github.com/go-vgo/Mingw)

####For everything else(Linux等其他系统):
```
GCC
    
X11 with the XTest extension (also known as the Xtst library)

事件:
    
xcb,xkb,libxkbcommon
```
#####Ubuntu:

```yml

sudo apt-get install libx11-dev
#sudo apt-get install libgtkglextmm-x11-dev
#sudo apt-get install libghc6-x11-dev
#sudo apt-get install libgl1-mesa-swx11-dev
sudo apt-get install xorg-dev

sudo apt-get install libxtst-dev libpng++-dev   

#事件:

sudo apt-get install xcb libxcb-xkb-dev x11-xkb-utils libx11-xcb-dev libxkbcommon-x11-dev
sudo apt-get install libxkbcommon-dev

```

##Installation:
```
go get github.com/go-vgo/robotgo
```
  It's that easy!

##Update:
```
go get -u github.com/go-vgo/robotgo   
```

##[Examples:](https://github.com/go-vgo/robotgo/blob/master/example/main.go)

###鼠标

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

###键盘

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

###屏幕

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

###位图

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

###事件

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
  abool := robotgo.ShowAlert("test", "robotgo")
  if abool == 0 {
    Println("ok@@@", "ok")
  }

  title := robotgo.GetTitle()
  Println("title@@@", title)
} 
```

##Future
- Update Find an image on screen, read pixels from an image
- Update Window Handle
- Try support Android,maybe support IOS.

##Contributors

- See [contributors page](https://github.com/go-vgo/robotgo/graphs/contributors) for full list of contributors.
- See [Contribution Guidelines](https://github.com/go-vgo/robotgo/blob/master/CONTRIBUTING.md).