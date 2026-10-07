# Robotgo

<p align="center">
  <img src="../docs/robotgo-logo.svg" width="560" alt="RobotGo 标志 — 带鼠标指针的机器人" />
</p>

[![Build Status](https://github.com/go-vgo/robotgo/workflows/Go/badge.svg)](https://github.com/go-vgo/robotgo/commits/master)
[![CircleCI Status](https://circleci.com/gh/go-vgo/robotgo.svg?style=shield)](https://circleci.com/gh/go-vgo/robotgo)
[![codecov](https://codecov.io/gh/go-vgo/robotgo/branch/master/graph/badge.svg)](https://codecov.io/gh/go-vgo/robotgo)
[![golangci-lint](https://github.com/go-vgo/robotgo/actions/workflows/lint.yml/badge.svg)](https://github.com/go-vgo/robotgo/actions/workflows/lint.yml)
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
- [无 Cgo 构建](#cgo-free-builds)
- [安装](#installation)
- [更新](#update)
- [示例](#examples)
- [类型转换与按键](https://github.com/go-vgo/robotgo/blob/master/docs/keys.md)
- [交叉编译](https://github.com/go-vgo/robotgo/blob/master/docs/install.md#crosscompiling)
- [作者](#authors)
- [计划](#plans)
- [许可证](#license)

## Docs

- [GoDoc](https://godoc.org/github.com/go-vgo/robotgo) <br>
- [API 文档](https://github.com/go-vgo/robotgo/blob/master/docs/doc.md)（已弃用，不再更新）

## Binding

[ADB](https://github.com/vcaesar/adb)，封装的 Android adb API。

## Requirements

**Go 1.26+**（见 [go.mod](../go.mod)）以及一个后端：

- **默认（Cgo）：** `CGO_ENABLED=1`、一个 C 编译器，以及下列平台库。
- **纯 Go（实验性）：** 无需 C 工具链；见[无 Cgo 构建](#cgo-free-builds)。

两者都需要为你所用功能安装对应的[功能相关依赖](#feature-specific-dependencies)。

### Default Cgo setup

#### macOS

Go 和 Xcode 命令行工具（Clang）：

```sh
brew install go
xcode-select --install
```

**权限（Cgo 与纯 Go）：** 在 **系统设置 > 隐私与安全性** 中，为应用或终端授予
**辅助功能**（输入）和 **屏幕录制**（捕获）权限。

#### Windows

```powershell
winget install GoLang.Go
```

此外还需 **一个** C 工具链。[LLVM-MinGW](https://github.com/mstorsjo/llvm-mingw)：

```powershell
winget install MartinStorsjo.LLVM-MinGW.UCRT
$env:CC = "clang"
```

或 [MinGW-w64 (WinLibs)](https://winlibs.com/)：

```powershell
winget install BrechtSanders.WinLibs.POSIX.UCRT
$env:CC = "gcc"
```

工具链的 `bin` 目录（例如 `C:\mingw64\bin`）必须位于 `PATH` 中且与 `GOARCH`
匹配。其他任何兼容 Cgo 的编译器同样可用，例如手动安装的
[MinGW-w64](https://sourceforge.net/projects/mingw-w64/files)。

[Bitmap](https://github.com/vcaesar/bitmap) 仅为 MinGW-w64 附带 libpng；
使用其他工具链时，需自行构建 libpng。

#### Linux (X11)

构建：GCC、libc 头文件、X11/XTest 头文件。运行：带 XTEST 的 X 服务器并设置
`DISPLAY`。原生 Wayland 请见[无 Cgo 构建](#cgo-free-builds)。

**Ubuntu / Debian：**

```sh
# Go（发行版自带的包过旧时用 Snap）
sudo snap install go --classic
# 或：sudo apt install golang

# 核心（Cgo）：编译器、libc、X11、XTest
sudo apt install gcc libc6-dev libx11-dev xorg-dev libxtst-dev

# 剪贴板（X11）
sudo apt install xsel xclip

# Bitmap：libpng
sudo apt install libpng++-dev

# GoHook：XCB、XKB
sudo apt install xcb libxcb-xkb-dev x11-xkb-utils libx11-xcb-dev libxkbcommon-x11-dev libxkbcommon-dev
```

**Fedora：**

```sh
# 核心（Cgo）：编译器、libc、X11、XTest
sudo dnf install gcc glibc-devel libX11-devel libXtst-devel

# 剪贴板（X11）
sudo dnf install xsel xclip

# Bitmap：libpng
sudo dnf install libpng-devel

# GoHook：XKB（Fedora < 34：用 xorg-x11-xkb-utils-devel 代替 xkbcomp-devel）
sudo dnf install libxkbcommon-devel libxkbcommon-x11-devel xkbcomp-devel
```

### Feature-specific dependencies

在核心环境之上，仅在使用相应功能时才需要：

- **Linux 剪贴板（Cgo 与纯 Go）：** X11 上需 `xclip` 或 `xsel`；
  Wayland 上需 `wl-clipboard`：

  ```sh
  sudo apt install wl-clipboard   # Ubuntu / Debian
  sudo dnf install wl-clipboard   # Fedora
  ```

- **[Bitmap](https://github.com/vcaesar/bitmap)：** libpng 头文件。
- **[GoHook](https://github.com/robotn/gohook)：** Linux 上的 XCB/XKB 头文件。

Bitmap 和 GoHook 是独立的 Cgo 包；纯 Go 的 RobotGo 后端并不会消除它们对 C 的
依赖。

## Cgo-free Builds

**实验性的纯 Go 后端**可在 `CGO_ENABLED=0` 下构建与交叉编译（无需 GCC、MinGW、
Xcode 或 X11 头文件）。导入路径和 API 不变；后端由**构建标签**选择 ——
仅设置 `CGO_ENABLED=0` 不会选择任何后端。

| 目标 `GOOS`   | 构建标签  | 包                     | 运行时要求                                                       |
| ------------- | --------- | ---------------------- | ---------------------------------------------------------------- |
| `windows`     | `win`     | [win](../win/)         | 仅需 Win32，无额外运行时                                         |
| `darwin`      | `mac`     | [darwin](../darwin/)   | 系统框架与[隐私权限](#macos)                                     |
| `linux`       | `x11`     | [x11](../x11/)         | 带 XTEST 的 X 服务器，已设置 `DISPLAY`                           |
| `linux`       | `wayland` | [wayland](../wayland/) | wlroots 合成器，见 [Wayland](#wayland)                           |
| `linux`       | `libei`   | [libei](../libei/)     | RemoteDesktop portal（GNOME/KDE），见 [libei](#libei-gnome--kde) |

**`purego`** 是一个快捷方式：macOS 上为 `mac`，Windows 上为 `win`，Linux 上为
`wayland`。在 Linux 上可追加 `x11` 或 `libei` 来覆盖（`-tags "purego,x11"`）。

只使用**一个**与 `GOOS` 匹配的后端标签；诸如 `x11,libei` 的组合无法编译
（`purego` 加一个 Linux 覆盖标签是唯一例外）。

### Build commands

```sh
# 当前操作系统
CGO_ENABLED=0 go build -tags purego .

# 交叉编译
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -tags win .
CGO_ENABLED=0 GOOS=darwin  GOARCH=arm64 go build -tags mac .
CGO_ENABLED=0 GOOS=linux   GOARCH=amd64 go build -tags x11 .
CGO_ENABLED=0 GOOS=linux   GOARCH=amd64 go build -tags wayland .
CGO_ENABLED=0 GOOS=linux   GOARCH=amd64 go build -tags libei .
```

PowerShell：`$env:CGO_ENABLED = "0"; go build -tags purego .`

请构建模块根目录（`.`）或你自己的包，而不是 `./...`：`examples/` 和部分子包需要
仅 Cgo 可用的 API。

### Limitations

上述运行时要求依然适用。仅 Cgo 可用的 API，如 `CaptureScreen` / `FreeBitmap`，
并不存在；请使用 `CaptureImg`。不支持的操作返回 `ErrNotSupported`（或空结果）。

- **`mac`：** 无窗口管理 —— `ActiveName` 返回 `ErrNotSupported`，
  `GetTitle` 返回 `""`，最小化/最大化/关闭为空操作。
- **`wayland`、`libei`：** `Location()` 返回 RobotGo 最后注入的位置，而非物理
  光标位置。剪贴板需要 `wl-clipboard`。

#### Wayland

需要一个 wlroots 合成器（Sway、Hyprland、Wayfire 等），已设置 `WAYLAND_DISPLAY`
和 `XDG_RUNTIME_DIR`，并暴露：

| 功能     | 协议全局对象                       |
| -------- | ---------------------------------- |
| 鼠标控制 | `zwlr_virtual_pointer_manager_v1`  |
| 键盘控制 | `zwp_virtual_keyboard_manager_v1`  |
| 屏幕捕获 | `zwlr_screencopy_manager_v1`       |
| 窗口管理 | `zwlr_foreign_toplevel_manager_v1` |

GNOME 和 KDE 不提供这些协议；请在其上使用 `libei`。

#### libei (GNOME / KDE)

输入通过 `xdg-desktop-portal` 的 **RemoteDesktop** D-Bus 接口完成：需要会话总线、
`xdg-desktop-portal` 以及一个实现该接口的后端（`xdg-desktop-portal-gnome`、
`-kde` 或 `-wlr`）。首次运行会显示授权对话框；恢复令牌缓存在
`$XDG_STATE_HOME/robotgo` 下。

绝对移动使用关联的 ScreenCast 流（默认）；设置 `libei.LinkScreenCast = false`
则仅使用相对移动。屏幕捕获和窗口管理返回 `ErrNotSupported`。

## Installation

```sh
go get github.com/go-vgo/robotgo
```

```go
import "github.com/go-vgo/robotgo"
```

使用 Bitmap 时出现 `png.h: No such file or directory`：请查看其 libpng 要求以及
[issues/47](https://github.com/go-vgo/robotgo/issues/47)。

## Update

```sh
go get -u github.com/go-vgo/robotgo
```

注意 go1.10.x 的 C 文件编译缓存问题，[golang #24355](https://github.com/golang/go/issues/24355)。
`go mod vendor` 问题，[golang #26366](https://github.com/golang/go/issues/26366)。

## [Examples](https://github.com/go-vgo/robotgo/blob/master/examples)

#### [鼠标](https://github.com/go-vgo/robotgo/blob/master/examples/mouse/main.go)

```Go
package main

import (
  "errors"
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

  // 检查返回的错误, 无效参数会返回错误
  if err := robotgo.Click("left", "double"); err != nil {
    fmt.Println("robotgo.Click error:", err)
  }
  if err := robotgo.ScrollDir(10, "forward"); err != nil {
    fmt.Println("robotgo.ScrollDir error:", err)
  }

  // MoveSmoothRelative 将平滑移动失败报告为 ErrSmoothMove
  err := robotgo.MoveSmoothRelative(10, -10)
  if errors.Is(err, robotgo.ErrSmoothMove) {
    fmt.Println("smooth move failed:", err)
  }

  // 纯 Go 后端对不支持的操作返回 ErrNotSupported
  err = robotgo.Move(100, 200)
  if errors.Is(err, robotgo.ErrNotSupported) {
    fmt.Println("robotgo.Move is not supported by this backend")
  }
}
```

#### [键盘](https://github.com/go-vgo/robotgo/blob/master/examples/key/main.go)

```Go
package main

import (
  "errors"
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

  // 未知的按键名会返回错误, 而不是按下错误的按键
  if err := robotgo.KeyTap("notakey"); err != nil {
    fmt.Println("robotgo.KeyTap error:", err)
  }

  // TypeStr 返回第一个输入错误, Type 只返回已输入的字符数
  if err := robotgo.TypeStr("Hello, 世界"); err != nil {
    fmt.Println("robotgo.TypeStr error:", err)
  }

  // 按住修饰键, 然后始终释放它并检查两个错误
  if err := robotgo.KeyDown("shift"); err == nil {
    err = errors.Join(robotgo.KeyTap("a"), robotgo.KeyUp("shift"))
    if err != nil {
      fmt.Println("shift + a error:", err)
    }
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

- 更好的多屏支持
- 更新窗口句柄
- 尝试支持 Android 和 iOS

## Contributors

- 完整的贡献者列表请见[贡献者页面](https://github.com/go-vgo/robotgo/graphs/contributors)。
- 请参阅[贡献指南](https://github.com/go-vgo/robotgo/blob/master/CONTRIBUTING.md)。

## License

Robotgo 主要依据 “Apache License (Version 2.0)” 的条款进行分发，部分内容受各类 BSD 风格许可证约束。

详见 [LICENSE-APACHE](http://www.apache.org/licenses/LICENSE-2.0)、[LICENSE](https://github.com/go-vgo/robotgo/blob/master/LICENSE)。
