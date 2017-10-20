# Robotgo
[![Build Status](https://travis-ci.org/go-vgo/robotgo.svg)](https://travis-ci.org/go-vgo/robotgo)
[![CircleCI Status](https://circleci.com/gh/go-vgo/robotgo.svg?style=shield)](https://circleci.com/gh/go-vgo/robotgo)
[![Go Report Card](https://goreportcard.com/badge/github.com/go-vgo/robotgo)](https://goreportcard.com/report/github.com/go-vgo/robotgo)
[![GoDoc](https://godoc.org/github.com/go-vgo/robotgo?status.svg)](https://godoc.org/github.com/go-vgo/robotgo)
[![Release](https://github-release-version.herokuapp.com/github/go-vgo/robotgo/release.svg?style=flat)](https://github.com/go-vgo/robotgo/releases/latest)
[![Join the chat at https://gitter.im/go-vgo/robotgo](https://badges.gitter.im/Join%20Chat.svg)](https://gitter.im/go-vgo/robotgo?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)
  
  >Golang 跨平台自动化系统，控制键盘鼠标位图和读取屏幕，窗口句柄以及全局事件监听
  
RobotGo 支持 Mac, Windows, and Linux(X11).

这是一项正在完善中的工作.

提 Issues 请到 [Github](https://github.com/go-vgo/robotgo), 便于统一管理和即时更新

QQ 群: 595877611

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
- [中文文档](https://github.com/go-vgo/robotgo/blob/master/docs/doc_zh.md)&nbsp;&nbsp;&nbsp;
- [English Docs](https://github.com/go-vgo/robotgo/blob/master/docs/doc.md) 
- [GoDoc](https://godoc.org/github.com/go-vgo/robotgo)

## Requirements:
环境要求:

在安装 RobotGo 之前, 请确保 Golang、GCC 被正确安装

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

#### For everything else (Linux 等其他系统):
```
GCC
    
X11 with the XTest extension (also known as the Xtst library)

事件:
    
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

#### [鼠标](https://github.com/go-vgo/robotgo/blob/master/examples/mouse/main.go)

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

#### [键盘](https://github.com/go-vgo/robotgo/blob/master/examples/key/main.go)

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

  robotgo.WriteAll("测试")
  text, err := robotgo.ReadAll()
  if err == nil {
    fmt.Println(text)
  }
} 
```

#### [屏幕](https://github.com/go-vgo/robotgo/blob/master/examples/screen/main.go)

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

<!-- #### [位图](https://github.com/go-vgo/robotgo/blob/master/examples/bitmap/mian.go)

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

#### [事件](https://github.com/go-vgo/robotgo/blob/master/examples/event/main.go)

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
#### [窗口句柄](https://github.com/go-vgo/robotgo/blob/master/examples/window/main.go)

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
- 更新 Find an image on screen, read pixels from an image
- 更新 Window Handle
- 支持 UTF-8, 过渡方案: TypeStr
- 尝试支持 Android, 也许支持 IOS
- 移除 zlib/libpng 依赖

## Donate

支持 robotgo, [buy me a coffee](https://github.com/go-vgo/buy-me-a-coffee).

#### Paypal

Donate money by [paypal](https://www.paypal.me/veni0/25) to my account [vzvway@gmail.com](vzvway@gmail.com)

## Contributors

- See [contributors page](https://github.com/go-vgo/robotgo/graphs/contributors) for full list of contributors.
- See [Contribution Guidelines](https://github.com/go-vgo/robotgo/blob/master/CONTRIBUTING.md).

## License

Robotgo is primarily distributed under the terms of both the MIT license and the Apache License (Version 2.0), with portions covered by various BSD-like licenses.

See [LICENSE-APACHE](http://www.apache.org/licenses/LICENSE-2.0), [LICENSE-MIT](https://github.com/go-vgo/robotgo/blob/master/LICENSE).