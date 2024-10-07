# Robotgo

<!--<img align="right" src="https://raw.githubusercontent.com/go-vgo/robotgo/master/logo.jpg">-->
<!--[![Build Status](https://travis-ci.org/go-vgo/robotgo.svg)](https://travis-ci.org/go-vgo/robotgo)
[![codecov](https://codecov.io/gh/go-vgo/robotgo/branch/master/graph/badge.svg)](https://codecov.io/gh/go-vgo/robotgo)-->
<!--<a href="https://circleci.com/gh/go-vgo/robotgo/tree/dev"><img src="https://img.shields.io/circleci/project/go-vgo/robotgo/dev.svg" alt="Build Status"></a>-->

[![Build Status](https://github.com/go-vgo/robotgo/workflows/Go/badge.svg)](https://github.com/go-vgo/robotgo/commits/master)
[![CircleCI Status](https://circleci.com/gh/go-vgo/robotgo.svg?style=shield)](https://circleci.com/gh/go-vgo/robotgo)
[![Build Status](https://travis-ci.org/go-vgo/robotgo.svg)](https://travis-ci.org/go-vgo/robotgo)
![Appveyor](https://ci.appveyor.com/api/projects/status/github/go-vgo/robotgo?branch=master&svg=true)
[![Go Report Card](https://goreportcard.com/badge/github.com/go-vgo/robotgo)](https://goreportcard.com/report/github.com/go-vgo/robotgo)
[![GoDoc](https://pkg.go.dev/badge/github.com/go-vgo/robotgo?status.svg)](https://pkg.go.dev/github.com/go-vgo/robotgo?tab=doc)
[![GitHub release](https://img.shields.io/github/release/go-vgo/robotgo.svg)](https://github.com/go-vgo/robotgo/releases/latest)
[![Join the chat at https://gitter.im/go-vgo/robotgo](https://badges.gitter.im/Join%20Chat.svg)](https://gitter.im/go-vgo/robotgo?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)

<!-- [![Release](https://github-release-version.herokuapp.com/github/go-vgo/robotgo/release.svg?style=flat)](https://github.com/go-vgo/robotgo/releases/latest) -->
<!-- <a href="https://github.com/go-vgo/robotgo/releases"><img src="https://img.shields.io/badge/%20version%20-%206.0.0%20-blue.svg?style=flat-square" alt="Releases"></a> -->

> Golang Desktop Automation. Control the mouse, keyboard, read the screen, process, Window Handle, image and bitmap and global event listener.

RobotGo supports Mac, Windows, and Linux(X11); and robotgo supports arm64 and x86-amd64.

## Contents

- [Docs](#docs)
- [Binding](#binding)
- [Requirements](#requirements)
- [Installation](#installation)
- [Update](#update)
- [Examples](#examples)
- [Type Conversion and keys](https://github.com/go-vgo/robotgo/blob/master/docs/keys.md)
- [Cross-Compiling](https://github.com/go-vgo/robotgo/blob/master/docs/install.md#crosscompiling)
- [Authors](#authors)
- [Plans](#plans)
- [Donate](#donate)
- [Contributors](#contributors)
- [License](#license)

## Docs

- [GoDoc](https://godoc.org/github.com/go-vgo/robotgo) <br>
- [API Docs](https://github.com/go-vgo/robotgo/blob/master/docs/doc.md) (Deprecated, no updated)

## Binding:

[ADB](https://github.com/vcaesar/adb), packaging android adb API.

[Robotn](https://github.com/vcaesar/robotn), binding JavaScript and other, support more language.

## Requirements:

Now, Please make sure `Golang, GCC` is installed correctly before installing RobotGo.

### ALL:

```
Golang

GCC
```

#### For MacOS:

Xcode Command Line Tools (And Privacy setting: [#277](https://github.com/go-vgo/robotgo/issues/277))

```
xcode-select --install
```

#### For Windows:

[MinGW-w64](https://sourceforge.net/projects/mingw-w64/files) (Use recommended) or others Mingw [llvm-mingw](https://github.com/mstorsjo/llvm-mingw);

Download the Mingw, then set system environment variables `C:\mingw64\bin` to the Path.
[Set environment variables to run GCC from command line](https://www.youtube.com/results?search_query=Set+environment+variables+to+run+GCC+from+command+line).

`Or the other GCC` (But you should compile the "libpng" with yourself when use the [bitmap](https://github.com/vcaesar/bitmap).)

#### For everything else:

```
GCC

X11 with the XTest extension (the Xtst library)

"Clipboard": xsel xclip


"Bitmap": libpng (Just used by the "bitmap".)

"Event-Gohook": xcb, xkb, libxkbcommon (Just used by the "hook".)

```

##### Ubuntu:

```yml
# gcc
sudo apt install gcc libc6-dev

# x11
sudo apt install libx11-dev xorg-dev libxtst-dev

# Clipboard
sudo apt install xsel xclip

#
# Bitmap
sudo apt install libpng++-dev

# GoHook
sudo apt install xcb libxcb-xkb-dev x11-xkb-utils libx11-xcb-dev libxkbcommon-x11-dev libxkbcommon-dev

```

##### Fedora:

```yml
# x11
sudo dnf install libXtst-devel

# Clipboard
sudo dnf install xsel xclip

#
# Bitmap
sudo dnf install libpng-devel

# GoHook
sudo dnf install libxkbcommon-devel libxkbcommon-x11-devel xorg-x11-xkb-utils-devel

```

## Installation:

With Go module support (Go 1.11+), just import:

```go
import "github.com/go-vgo/robotgo"
```

Otherwise, to install the robotgo package, run the command:

```
go get github.com/go-vgo/robotgo
```

png.h: No such file or directory? Please see [issues/47](https://github.com/go-vgo/robotgo/issues/47).

## Update:

```
go get -u github.com/go-vgo/robotgo
```

Note go1.10.x C file compilation cache problem, [golang #24355](https://github.com/golang/go/issues/24355).
`go mod vendor` problem, [golang #26366](https://github.com/golang/go/issues/26366).

## [Examples:](https://github.com/go-vgo/robotgo/blob/master/examples)

#### [Mouse](https://github.com/go-vgo/robotgo/blob/master/examples/mouse/main.go)

```Go
package main

import (
  "github.com/go-vgo/robotgo"
)

func main() {
  robotgo.MouseSleep = 100

  robotgo.ScrollDir(10, "up")
  robotgo.ScrollDir(20, "right")

  robotgo.Scroll(0, -10)
  robotgo.Scroll(100, 0)

  robotgo.MilliSleep(100)
  robotgo.ScrollSmooth(-10, 6)
  // robotgo.ScrollRelative(10, -100)

  robotgo.Move(10, 20)
  robotgo.MoveRelative(0, -10)
  robotgo.DragSmooth(10, 10)

  robotgo.Click("wheelRight")
  robotgo.Click("left", true)
  robotgo.MoveSmooth(100, 200, 1.0, 10.0)

  robotgo.Toggle("left")
  robotgo.Toggle("left", "up")
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
  robotgo.TypeStr("Hello World")
  robotgo.TypeStr("だんしゃり", 0, 1)
  // robotgo.TypeStr("テストする")

  robotgo.TypeStr("Hi, Seattle space needle, Golden gate bridge, One world trade center.")
  robotgo.TypeStr("Hi galaxy, hi stars, hi MT.Rainier, hi sea. こんにちは世界.")
  robotgo.Sleep(1)

  // ustr := uint32(robotgo.CharCodeAt("Test", 0))
  // robotgo.UnicodeType(ustr)

  robotgo.KeySleep = 100
  robotgo.KeyTap("enter")
  // robotgo.TypeStr("en")
  robotgo.KeyTap("i", "alt", "cmd")

  arr := []string{"alt", "cmd"}
  robotgo.KeyTap("i", arr)

  robotgo.MilliSleep(100)
  robotgo.KeyToggle("a")
  robotgo.KeyToggle("a", "up")

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
  "github.com/vcaesar/imgo"
)

func main() {
  x, y := robotgo.Location()
  fmt.Println("pos: ", x, y)

  color := robotgo.GetPixelColor(100, 200)
  fmt.Println("color---- ", color)

  sx, sy := robotgo.GetScreenSize()
  fmt.Println("get screen size: ", sx, sy)

  bit := robotgo.CaptureScreen(10, 10, 30, 30)
  defer robotgo.FreeBitmap(bit)

  img := robotgo.ToImage(bit)
  imgo.Save("test.png", img)

  num := robotgo.DisplaysNum()
  for i := 0; i < num; i++ {
    robotgo.DisplayID = i
    img1, _ := robotgo.CaptureImg()
    path1 := "save_" + strconv.Itoa(i)
    robotgo.Save(img1, path1+".png")
    robotgo.SaveJpeg(img1, path1+".jpeg", 50)

    img2, _ := robotgo.CaptureImg(10, 10, 20, 20)
    robotgo.Save(img2, "test_"+strconv.Itoa(i)+".png")

    x, y, w, h := robotgo.GetDisplayBounds(i)
    img3, err := robotgo.CaptureImg(x, y, w, h)
    fmt.Println("Capture error: ", err)
    robotgo.Save(img3, path1+"_1.png")
  }
}
```

#### [Bitmap](https://github.com/vcaesar/bitmap/blob/main/examples/main.go)

```Go
package main

import (
  "fmt"

  "github.com/go-vgo/robotgo"
  "github.com/vcaesar/bitmap"
)

func main() {
  bit := robotgo.CaptureScreen(10, 20, 30, 40)
  // use `defer robotgo.FreeBitmap(bit)` to free the bitmap
  defer robotgo.FreeBitmap(bit)

  fmt.Println("bitmap...", bit)
  img := robotgo.ToImage(bit)
  // robotgo.SavePng(img, "test_1.png")
  robotgo.Save(img, "test_1.png")

  bit2 := robotgo.ToCBitmap(robotgo.ImgToBitmap(img))
  fx, fy := bitmap.Find(bit2)
  fmt.Println("FindBitmap------ ", fx, fy)
  robotgo.Move(fx, fy)

  arr := bitmap.FindAll(bit2)
  fmt.Println("Find all bitmap: ", arr)

  fx, fy = bitmap.Find(bit)
  fmt.Println("FindBitmap------ ", fx, fy)

  bitmap.Save(bit, "test.png")
}
```

#### [OpenCV](https://github.com/vcaesar/gcv)

```Go
package main

import (
  "fmt"
  "math/rand"

  "github.com/go-vgo/robotgo"
  "github.com/vcaesar/gcv"
  "github.com/vcaesar/bitmap"
)

func main() {
  opencv()
}

func opencv() {
  name := "test.png"
  name1 := "test_001.png"
  robotgo.SaveCapture(name1, 10, 10, 30, 30)
  robotgo.SaveCapture(name)

  fmt.Print("gcv find image: ")
  fmt.Println(gcv.FindImgFile(name1, name))
  fmt.Println(gcv.FindAllImgFile(name1, name))

  bit := bitmap.Open(name1)
  defer robotgo.FreeBitmap(bit)
  fmt.Print("find bitmap: ")
  fmt.Println(bitmap.Find(bit))

  // bit0 := robotgo.CaptureScreen()
  // img := robotgo.ToImage(bit0)
  // bit1 := robotgo.CaptureScreen(10, 10, 30, 30)
  // img1 := robotgo.ToImage(bit1)
  // defer robotgo.FreeBitmapArr(bit0, bit1)
  img, _ := robotgo.CaptureImg()
  img1, _ := robotgo.CaptureImg(10, 10, 30, 30)

  fmt.Print("gcv find image: ")
  fmt.Println(gcv.FindImg(img1, img))
  fmt.Println()

  res := gcv.FindAllImg(img1, img)
  fmt.Println(res[0].TopLeft.Y, res[0].Rects.TopLeft.X, res)
  x, y := res[0].TopLeft.X, res[0].TopLeft.Y
  robotgo.Move(x, y-rand.Intn(5))
  robotgo.MilliSleep(100)
  robotgo.Click()

  res = gcv.FindAll(img1, img) // use find template and sift
  fmt.Println("find all: ", res)
  res1 := gcv.Find(img1, img)
  fmt.Println("find: ", res1)

  img2, _, _ := robotgo.DecodeImg("test_001.png")
  x, y = gcv.FindX(img2, img)
  fmt.Println(x, y)
}
```

#### [Event](https://github.com/robotn/gohook/blob/master/examples/main.go)

```Go
package main

import (
  "fmt"

  // "github.com/go-vgo/robotgo"
  hook "github.com/robotn/gohook"
)

func main() {
  add()
  low()
  event()
}

func add() {
  fmt.Println("--- Please press ctrl + shift + q to stop hook ---")
  hook.Register(hook.KeyDown, []string{"q", "ctrl", "shift"}, func(e hook.Event) {
    fmt.Println("ctrl-shift-q")
    hook.End()
  })

  fmt.Println("--- Please press w---")
  hook.Register(hook.KeyDown, []string{"w"}, func(e hook.Event) {
    fmt.Println("w")
  })

  s := hook.Start()
  <-hook.Process(s)
}

func low() {
	evChan := hook.Start()
	defer hook.End()

	for ev := range evChan {
		fmt.Println("hook: ", ev)
	}
}

func event() {
  ok := hook.AddEvents("q", "ctrl", "shift")
  if ok {
    fmt.Println("add events...")
  }

  keve := hook.AddEvent("k")
  if keve {
    fmt.Println("you press... ", "k")
  }

  mleft := hook.AddEvent("mleft")
  if mleft {
    fmt.Println("you press... ", "mouse left button")
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
    fmt.Println("pids... ", fpid)

    if len(fpid) > 0 {
      robotgo.TypeStr("Hi galaxy!", fpid[0])
      robotgo.KeyTap("a", fpid[0], "cmd")

      robotgo.KeyToggle("a", fpid[0])
      robotgo.KeyToggle("a", fpid[0], "up")

      robotgo.ActivePid(fpid[0])

      robotgo.Kill(fpid[0])
    }
  }

  robotgo.ActiveName("chrome")

  isExist, err := robotgo.PidExists(100)
  if err == nil && isExist {
    fmt.Println("pid exists is", isExist)

    robotgo.Kill(100)
  }

  abool := robotgo.Alert("test", "robotgo")
  if abool {
 	  fmt.Println("ok@@@ ", "ok")
  }

  title := robotgo.GetTitle()
  fmt.Println("title@@@ ", title)
}
```

## Authors

- [The author is vz](https://github.com/vcaesar)
- [Maintainers](https://github.com/orgs/go-vgo/people)
- [Contributors](https://github.com/go-vgo/robotgo/graphs/contributors)

## Plans

- Refactor some C code to Go (such as x11, windows)
- Better multiscreen support
- Wayland support
- Update Window Handle
- Try to support Android and IOS

## Contributors

- See [contributors page](https://github.com/go-vgo/robotgo/graphs/contributors) for full list of contributors.
- See [Contribution Guidelines](https://github.com/go-vgo/robotgo/blob/master/CONTRIBUTING.md).

## License

Robotgo is primarily distributed under the terms of "both the MIT license and the Apache License (Version 2.0)", with portions covered by various BSD-like licenses.

See [LICENSE-APACHE](http://www.apache.org/licenses/LICENSE-2.0), [LICENSE-MIT](https://github.com/go-vgo/robotgo/blob/master/LICENSE).
