# Robotgo

[![Build Status](https://github.com/go-vgo/robotgo/workflows/Go/badge.svg)](https://github.com/go-vgo/robotgo/commits/master)
[![CircleCI Status](https://circleci.com/gh/go-vgo/robotgo.svg?style=shield)](https://circleci.com/gh/go-vgo/robotgo)
[![Go Report Card](https://goreportcard.com/badge/github.com/go-vgo/robotgo)](https://goreportcard.com/report/github.com/go-vgo/robotgo)
[![GoDoc](https://pkg.go.dev/badge/github.com/go-vgo/robotgo?status.svg)](https://pkg.go.dev/github.com/go-vgo/robotgo?tab=doc)
[![GitHub release](https://img.shields.io/github/release/go-vgo/robotgo.svg)](https://github.com/go-vgo/robotgo/releases/latest)
<a href="https://discord.gg/npPb3NzE4A"><img src="https://img.shields.io/discord/1484658282777018551.svg?logo=discord&logoColor=white&label=Discord&color=5865F2" alt="Join the Discord chat at https://discord.gg/npPb3NzE4A"></a>

[English](../README.md) | 简体中文 | [繁體中文](README.zht.md) | [日本語](README.ja.md) | [한국어](README.ko.md) | [Français](README.fr.md) | [Deutsch](README.de.md) | [Español](README.es.md) | [Русский](README.ru.md) | [Português](README.pt.md)

> Golang 桌面自动化、自动测试以及 AI 计算机操作（Computer Use）。<br>
> 控制鼠标、键盘，读取屏幕，进程、窗口句柄、图像与位图，以及全局事件监听。

RobotGo 支持 Mac、Windows 和 Linux；并且支持 arm64 与 x86-amd64 架构。

我正在打造 [Codg](https://github.com/vcaesar/codg)，一个简单易用的 AI 智能体（Agent）工作系统：自动化、异步、并发、高效且高准确度。

<p align="center">
<a href="https://github.com/vcaesar/codg" rel="nofollow">
<img width="800" alt="Codg Demo" src="https://github.com/vcaesar/codg/raw/main/demo/26-04-1.png" />
</a>
</p>

[RobotGo-Pro](https://github.com/vcaesar/robotgo-pro) 提供 JavaScript、Python、Lua 等其他语言版本、技术支持、新功能以及最新的 robotgo 版本（“目前无开源版本”）。

## 目录

- [文档](#docs)
- [绑定](#binding)
- [环境要求](#requirements)
- [安装](#installation)
- [更新](#update)
- [无 Cgo 构建](#cgo-free-builds)
- [示例](#examples)
- [类型转换与按键](https://github.com/go-vgo/robotgo/blob/master/docs/keys.md)
- [交叉编译](https://github.com/go-vgo/robotgo/blob/master/docs/install.md#crosscompiling)
- [作者](#authors)
- [计划](#plans)
- [许可证](#license)

## Docs

- [GoDoc](https://godoc.org/github.com/go-vgo/robotgo) <br>
- [API 文档](https://github.com/go-vgo/robotgo/blob/master/docs/doc.md)（已弃用，不再更新）

## Binding:

[ADB](https://github.com/vcaesar/adb)，封装的 Android adb API。

## Requirements:

现在，请在安装 RobotGo 之前确保 `Golang、GCC` 已被正确安装。

### 全部平台：

```
Golang

GCC
```

#### MacOS：

```
brew install go
```

Xcode 命令行工具；<br>
并在隐私设置中，于以下位置添加“屏幕录制”和“辅助功能”权限：<br>
`系统设置 > 隐私与安全性 > 辅助功能、屏幕与系统音频录制`。

```
xcode-select --install
```

#### Windows：

```
winget install Golang.go
```

[llvm-mingw](https://github.com/mstorsjo/llvm-mingw)

```
winget install MartinStorsjo.LLVM-MinGW.UCRT
```

或者 [Mingw-w64](https://sourceforge.net/projects/mingw-w64/files)

```
winget install BrechtSanders.WinLibs.POSIX.UCRT
```

或者下载 [Mingw-w64](https://sourceforge.net/projects/mingw-w64/files) 以及其他 gcc，然后将类似 `C:\mingw64\bin` 的路径设置到系统环境变量 `Path` 中。
[设置环境变量以便从命令行运行 GCC](https://www.youtube.com/results?search_query=Set+environment+variables+to+run+GCC+from+command+line)。

`或者使用其他 GCC`（除 Mingw-w64 之外，使用 [bitmap](https://github.com/vcaesar/bitmap) 时你需要自行编译 “libpng”。）

#### 其他所有平台：

```
GCC

带 XTest 扩展的 X11（即 Xtst 库）

“剪贴板”：xsel xclip

“位图”：libpng（仅 “bitmap” 使用。）

“事件-Gohook”：xcb, xkb, libxkbcommon（仅 “hook” 使用。）
```

##### Ubuntu：

```yml
# sudo apt install golang
sudo snap install go  --classic

# gcc
sudo apt install gcc libc6-dev

# x11
sudo apt install libx11-dev xorg-dev libxtst-dev

# Clipboard
sudo apt install xsel xclip

# Bitmap
sudo apt install libpng++-dev

# GoHook
sudo apt install xcb libxcb-xkb-dev x11-xkb-utils libx11-xcb-dev libxkbcommon-x11-dev libxkbcommon-dev
```

##### Fedora：

```yml
# x11
sudo dnf install libXtst-devel

# Clipboard
sudo dnf install xsel xclip

# Bitmap
sudo dnf install libpng-devel

# GoHook
sudo dnf install libxkbcommon-devel libxkbcommon-x11-devel xkbcomp-devel
xorg-x11-xkb-utils-devel (< Fedora 34)
```

#### Wayland

Wayland 后端是 **纯 Go（无 Cgo）** 实现，因此无需任何系统 C 库。它需要一个
基于 wlroots 的合成器（Sway、Hyprland、Wayfire 等），并支持以下协议：

```
zwlr_virtual_pointer_v1            (鼠标控制)
zwp_virtual_keyboard_v1            (键盘控制)
zwlr_screencopy_v1                 (屏幕捕获)
zwlr_foreign_toplevel_management_v1 (窗口管理)
```

GNOME 和 KDE **不**原生支持这些协议。

#### libei（GNOME / KDE）

libei 后端同样是 **纯 Go（无 Cgo）** 实现。它通过 freedesktop 的
`xdg-desktop-portal` RemoteDesktop 接口驱动输入，因此可在 GNOME 和 KDE 上工作
（不同于 wlroots Wayland 后端）。它需要：

```
xdg-desktop-portal               (portal D-Bus 服务)
xdg-desktop-portal-gnome / -kde  (你桌面的 portal 后端)
```

注意：libei 后端仅处理鼠标和键盘输入。屏幕捕获和窗口管理会返回 `ErrNotSupported`。

## Installation:

在支持 Go module 的情况下（Go 1.11+），只需 import：

```go
import "github.com/go-vgo/robotgo"
```

否则，运行以下命令安装 robotgo 包：

```
go get github.com/go-vgo/robotgo
```

png.h: No such file or directory？请参阅 [issues/47](https://github.com/go-vgo/robotgo/issues/47)。

## Update:

```
go get -u github.com/go-vgo/robotgo
```

注意 go1.10.x 的 C 文件编译缓存问题，[golang #24355](https://github.com/golang/go/issues/24355)。
`go mod vendor` 问题，[golang #26366](https://github.com/golang/go/issues/26366)。

## Cgo-free Builds:

RobotGo 为 Windows、macOS、X11、Wayland 和 libei（Linux）提供了 **纯 Go（无 Cgo）** 后端，当前为实验性功能。
它们暴露相同的 `robotgo` API，因此你的代码无需改动 —— 只需一个构建标签。
这些后端可在 `CGO_ENABLED=0` 下交叉编译（无需 GCC、MinGW、Xcode 或 X11 头文件）。

| 后端                             | 构建标签  | Go 包                               |
| -------------------------------- | --------- | ----------------------------------- |
| Windows（无 Cgo）                | `win`     | `github.com/go-vgo/robotgo/win`     |
| macOS（通过 purego 调用 Quartz） | `mac`     | `github.com/go-vgo/robotgo/darwin`  |
| X11（Linux，纯 Go X 协议）       | `x11`     | `github.com/go-vgo/robotgo/x11`     |
| Wayland（Linux，wlroots）        | `wayland` | `github.com/go-vgo/robotgo/wayland` |
| libei（Linux，GNOME/KDE portal） | `libei`   | `github.com/go-vgo/robotgo/libei`   |
| 纯 Go 默认（所有平台）           | `purego`  | 选择上面的 `mac`/`win`/`wayland`    |

```sh
# 每个平台的纯 Go 默认后端，一个标签适用于所有目标：
# macOS -> mac，Windows -> win，Linux -> wayland（可与 x11/libei 组合以覆盖）
go build -tags purego .
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -tags "purego,x11" .

# Windows，无需 Cgo / 无需 MinGW
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -tags win .

# macOS，通过 purego 在运行时加载 Quartz/CoreGraphics（无需 Xcode）
CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -tags mac .

# X11，纯 Go 实现的 X 协议（XTEST）—— 无需 X11 头文件
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -tags x11 .

# Wayland，基于 wlroots 的合成器（Sway、Hyprland、Wayfire 等）
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -tags wayland .

# libei，通过 xdg-desktop-portal RemoteDesktop 支持 GNOME/KDE
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -tags libei .
```

注意：上面的示例构建模块根目录（`.`）而不是 `./...`，因为 `examples/` 和一些特定平台的子包使用了仅在默认 Cgo 后端下可用的 API。

在 `win` 标签下，默认的 Cgo/Win32 后端被排除，调用转发到纯 Go 的 `win` 包；
在 `mac` 标签下，默认的 Cgo/Quartz 后端被排除，调用转发到纯 Go 的 `darwin` 包（窗口管理返回 `ErrNotSupported`）；
在 `x11` 标签下，Cgo/X11 后端被排除，调用转发到纯 Go 的 `x11` 包；
在 `wayland` 标签下，Cgo/X11 后端被排除，调用转发到纯 Go 的 `wayland` 包；
在 `libei` 标签下，Cgo/X11 和 wlroots Wayland 后端均被排除，调用转发到纯 Go 的 `libei` 包。

`purego` 标签是一个跨平台快捷方式：它会在所有平台排除 Cgo 后端，并按目标 OS 选择默认纯 Go 后端 —— macOS 使用 `mac`，Windows 使用 `win`，Linux 使用 `wayland`。在 Linux 上，你可以将它与 `x11` 或 `libei` 组合（例如 `-tags "purego,libei"`）以选择不同的纯 Go 后端。

## [Examples:](https://github.com/go-vgo/robotgo/blob/master/examples)

#### [鼠标](https://github.com/go-vgo/robotgo/blob/master/examples/mouse/main.go)

```Go
package main

import (
  "fmt"
  "github.com/go-vgo/robotgo"
)

func main() {
  robotgo.MouseSleep = 300

  robotgo.Move(100, 100)
  fmt.Println(robotgo.Location())
  robotgo.Move(100, -200) // multi screen supported
  robotgo.MoveSmooth(120, -150)
  fmt.Println(robotgo.Location())

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

#### [键盘](https://github.com/go-vgo/robotgo/blob/master/examples/key/main.go)

```Go
package main

import (
  "fmt"

  "github.com/go-vgo/robotgo"
)

func main() {
  robotgo.Type("Hello World")
  robotgo.Type("だんしゃり", 0, 1)
  // robotgo.Type("テストする")

  robotgo.Type("Hi, Seattle space needle, Golden gate bridge, One world trade center.")
  robotgo.Type("Hi galaxy, hi stars, hi MT.Rainier, hi sea. こんにちは世界.")
  robotgo.Sleep(1)

  // ustr := uint32(robotgo.CharCodeAt("Test", 0))
  // robotgo.UnicodeType(ustr)

  robotgo.KeySleep = 100
  robotgo.KeyTap("enter")
  // robotgo.Type("en")
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

#### [屏幕](https://github.com/go-vgo/robotgo/blob/master/examples/screen/main.go)

```Go
package main

import (
  "fmt"
  "strconv"

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

#### [位图](https://github.com/vcaesar/bitmap/blob/main/examples/main.go)

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

#### [事件](https://github.com/robotn/gohook/blob/master/examples/main.go)

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

#### [窗口](https://github.com/go-vgo/robotgo/blob/master/examples/window/main.go)

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
      robotgo.Type("Hi galaxy!", fpid[0])
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

- [作者 Evans](https://github.com/vcaesar)
- [维护者](https://github.com/orgs/go-vgo/people)

## Plans

- 将部分 C 代码重构为 Go（例如 x11、windows）
- 更好的多屏支持
- Wayland 支持
- 更新窗口句柄
- 尝试支持 Android 和 iOS

## Contributors

- 完整的贡献者列表请见[贡献者页面](https://github.com/go-vgo/robotgo/graphs/contributors)。
- 请参阅[贡献指南](https://github.com/go-vgo/robotgo/blob/master/CONTRIBUTING.md)。

## License

Robotgo 主要依据 “Apache License (Version 2.0)” 的条款进行分发，部分内容受各类 BSD 风格许可证约束。

详见 [LICENSE-APACHE](http://www.apache.org/licenses/LICENSE-2.0)、[LICENSE](https://github.com/go-vgo/robotgo/blob/master/LICENSE)。
