# Robotgo

[![Build Status](https://github.com/go-vgo/robotgo/workflows/Go/badge.svg)](https://github.com/go-vgo/robotgo/commits/master)
[![CircleCI Status](https://circleci.com/gh/go-vgo/robotgo.svg?style=shield)](https://circleci.com/gh/go-vgo/robotgo)
[![Build Status](https://travis-ci.org/go-vgo/robotgo.svg)](https://travis-ci.org/go-vgo/robotgo)
![Appveyor](https://ci.appveyor.com/api/projects/status/github/go-vgo/robotgo?branch=master&svg=true)
[![Go Report Card](https://goreportcard.com/badge/github.com/go-vgo/robotgo)](https://goreportcard.com/report/github.com/go-vgo/robotgo)
[![GoDoc](https://godoc.org/github.com/go-vgo/robotgo?status.svg)](https://godoc.org/github.com/go-vgo/robotgo)
[![GitHub release](https://img.shields.io/github/release/go-vgo/robotgo.svg)](https://github.com/go-vgo/robotgo/releases/latest)
[![Join the chat at https://gitter.im/go-vgo/robotgo](https://badges.gitter.im/Join%20Chat.svg)](https://gitter.im/go-vgo/robotgo?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)

  >Golang 跨平台自动化系统，控制键盘鼠标位图和读取屏幕，窗口句柄以及全局事件监听

RobotGo 支持 Mac, Windows, and Linux(X11).

<br>

提 Issues 请到 [Github](https://github.com/go-vgo/robotgo), 便于统一管理和即时更新

QQ 讨论群: 595877611 (仅用于讨论)

## Contents
- [Docs](#docs)
- [Binding](#binding)
- [Requirements](#requirements)
- [Installation](#installation)
- [Update](#update)
- [Examples](#examples)
- [Cross-Compiling](#crosscompiling)
- [Authors](#authors)
- [Plans](#plans)
- [Donate](#donate)
- [Contributors](#contributors)
- [License](#license)

## Docs
- [GoDoc](https://godoc.org/github.com/go-vgo/robotgo)
- [中文文档](https://github.com/go-vgo/robotgo/blob/master/docs/doc_zh.md)&nbsp;&nbsp;&nbsp;
- [English Docs](https://github.com/go-vgo/robotgo/blob/master/docs/doc.md)

## Binding:

[Robotn](https://github.com/vcaesar/robotn), binding JavaScript and other, support more language.

## Requirements:
环境要求:

在安装 RobotGo 之前, 请确保 `Golang、GCC` 被正确安装

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
MinGW-w64 (推荐使用) or other GCC
```

#### For everything else (Linux 等其他系统):

```
GCC, libpng

X11 with the XTest extension (also known as the Xtst library)

事件:

xcb, xkb, libxkbcommon
```

##### Ubuntu:

```yml
sudo apt-get install gcc libc6-dev

sudo apt-get install libx11-dev xorg-dev libxtst-dev libpng++-dev

sudo apt-get install xcb libxcb-xkb-dev x11-xkb-utils libx11-xcb-dev libxkbcommon-x11-dev
sudo apt-get install libxkbcommon-dev

sudo apt-get install xsel xclip

```

#### Fedora:

```yml
sudo dnf install libxkbcommon-devel libXtst-devel libxkbcommon-x11-devel xorg-x11-xkb-utils-devel

sudo dnf install libpng-devel

sudo dnf install xsel xclip
```

## Installation:
```
go get github.com/go-vgo/robotgo
```
  It's that easy!

png.h: No such file or directory? Please see [issues/47](https://github.com/go-vgo/robotgo/issues/47).

## Update:
```
go get -u github.com/go-vgo/robotgo
```

注意 go1.10.x C 文件编译缓存问题, [golang #24355](https://github.com/golang/go/issues/24355).
`go mod vendor` problem, [golang #26366](https://github.com/golang/go/issues/26366).


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
  robotgo.TypeStr("Hello World")
  robotgo.TypeStr("测试")
  // robotgo.TypeString("测试")

  robotgo.TypeStr("山达尔星新星军团, galaxy. こんにちは世界.")
  robotgo.Sleep(1)

  // ustr := uint32(robotgo.CharCodeAt("测试", 0))
  // robotgo.UnicodeType(ustr)

  robotgo.KeyTap("enter")
  // robotgo.TypeString("en")
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
  fmt.Println("pos: ", x, y)

  color := robotgo.GetPixelColor(100, 200)
  fmt.Println("color----", color)
}
```

#### [位图](https://github.com/go-vgo/robotgo/blob/master/examples/bitmap/main.go)

```Go
package main

import (
	"fmt"

	"github.com/go-vgo/robotgo"
)

func main() {
  bitmap := robotgo.CaptureScreen(10, 20, 30, 40)
  // use `defer robotgo.FreeBitmap(bit)` to free the bitmap
  defer robotgo.FreeBitmap(bitmap)
  fmt.Println("...", bitmap)

  fx, fy := robotgo.FindBitmap(bitmap)
  fmt.Println("FindBitmap------", fx, fy)

  robotgo.SaveBitmap(bitmap, "test.png")
}
```

#### [事件](https://github.com/go-vgo/robotgo/blob/master/examples/event/main.go)

```Go
package main

import (
	"fmt"

	"github.com/go-vgo/robotgo"
)

func main() {
  ok := robotgo.AddEvents("q", "ctrl", "shift")
  if ok {
    fmt.Println("add events...")
  }

  keve := robotgo.AddEvent("k")
  if keve {
    fmt.Println("you press... ", "k")
  }

  mleft := robotgo.AddEvent("mleft")
  if mleft {
    fmt.Println("you press... ", "mouse left button")
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

    if len(fpid) > 0 {
      robotgo.ActivePID(fpid[0])

      robotgo.Kill(fpid[0])
    }
  }

  robotgo.ActiveName("chrome")

  isExist, err := robotgo.PidExists(100)
  if err == nil && isExist {
    fmt.Println("pid exists is", isExist)

    robotgo.Kill(100)
  }

  abool := robotgo.ShowAlert("test", "robotgo")
  if abool == 0 {
    fmt.Println("ok@@@ ", "ok")
  }

  title := robotgo.GetTitle()
  fmt.Println("title@@@ ", title)
}
```

## CrossCompiling

##### Windows64 to win32
```Go
SET CGO_ENABLED=1
SET GOARCH=386
go build main.go
```

#### Other to windows
```Go
GOOS=windows GOARCH=amd64 CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc CXX=x86_64-w64-mingw32-g++ go build -x ./
```
```
// CC=mingw-w64\x86_64-7.2.0-win32-seh-rt_v5-rev1\mingw64\bin\gcc.exe
// CXX=mingw-w64\x86_64-7.2.0-win32-seh-rt_v5-rev1\mingw64\bin\g++.exe
```

## Authors
* [The author is vz](https://github.com/vcaesar)
* [Maintainers](https://github.com/orgs/go-vgo/people)
* [Contributors](https://github.com/go-vgo/robotgo/graphs/contributors)

## Plans
- 更新 Find an image on screen, read pixels from an image
- 更新 Window Handle
- 尝试支持 Android, 也许支持 IOS

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