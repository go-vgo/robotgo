# Robotgo

[![Build Status](https://github.com/go-vgo/robotgo/workflows/Go/badge.svg)](https://github.com/go-vgo/robotgo/commits/master)
[![CircleCI Status](https://circleci.com/gh/go-vgo/robotgo.svg?style=shield)](https://circleci.com/gh/go-vgo/robotgo)
[![Go Report Card](https://goreportcard.com/badge/github.com/go-vgo/robotgo)](https://goreportcard.com/report/github.com/go-vgo/robotgo)
[![GoDoc](https://pkg.go.dev/badge/github.com/go-vgo/robotgo?status.svg)](https://pkg.go.dev/github.com/go-vgo/robotgo?tab=doc)
[![GitHub release](https://img.shields.io/github/release/go-vgo/robotgo.svg)](https://github.com/go-vgo/robotgo/releases/latest)
<a href="https://discord.gg/npPb3NzE4A"><img src="https://img.shields.io/discord/1484658282777018551.svg?logo=discord&logoColor=white&label=Discord&color=5865F2" alt="Join the Discord chat at https://discord.gg/npPb3NzE4A"></a>

[English](../README.md) | [简体中文](README.zh.md) | 繁體中文 | [日本語](README.ja.md) | [한국어](README.ko.md) | [Français](README.fr.md) | [Deutsch](README.de.md) | [Español](README.es.md) | [Русский](README.ru.md) | [Português](README.pt.md)

> Golang 桌面自動化、自動測試以及 AI 電腦操作（Computer Use）。<br>
> 控制滑鼠、鍵盤，讀取螢幕，行程、視窗控制代碼、影像與點陣圖，以及全域事件監聽。

RobotGo 支援 Mac、Windows 和 Linux；並且支援 arm64 與 x86-amd64 架構。

我正在打造 [Codg](https://github.com/vcaesar/codg)，一個簡單易用的 AI 智慧代理（Agent）工作系統：自動化、非同步、並行、高效且高準確度。

<p align="center">
<a href="https://github.com/vcaesar/codg" rel="nofollow">
<img width="800" alt="Codg Demo" src="https://github.com/vcaesar/codg/raw/main/demo/26-04-1.png" />
</a>
</p>

[RobotGo-Pro](https://github.com/vcaesar/robotgo-pro) 提供 JavaScript、Python、Lua 等其他語言版本、技術支援、新功能以及最新的 robotgo 版本（「目前無開源版本」）。

## 目錄

- [文件](#docs)
- [綁定](#binding)
- [環境需求](#requirements)
- [安裝](#installation)
- [更新](#update)
- [無 Cgo 構建](#cgo-free-builds)
- [範例](#examples)
- [型別轉換與按鍵](https://github.com/go-vgo/robotgo/blob/master/docs/keys.md)
- [交叉編譯](https://github.com/go-vgo/robotgo/blob/master/docs/install.md#crosscompiling)
- [作者](#authors)
- [計畫](#plans)
- [授權](#license)

## Docs

- [GoDoc](https://godoc.org/github.com/go-vgo/robotgo) <br>
- [API 文件](https://github.com/go-vgo/robotgo/blob/master/docs/doc.md)（已棄用，不再更新）

## Binding:

[ADB](https://github.com/vcaesar/adb)，封裝的 Android adb API。

## Requirements:

現在，請在安裝 RobotGo 之前確保 `Golang、GCC` 已被正確安裝。

### 全部平台：

```
Golang

GCC
```

#### MacOS：

```
brew install go
```

Xcode 命令列工具；<br>
並在隱私設定中，於以下位置新增「螢幕錄製」和「輔助使用」權限：<br>
`系統設定 > 隱私權與安全性 > 輔助使用、螢幕與系統音訊錄製`。

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

或者下載 [Mingw-w64](https://sourceforge.net/projects/mingw-w64/files) 以及其他 gcc，然後將類似 `C:\mingw64\bin` 的路徑設定到系統環境變數 `Path` 中。
[設定環境變數以便從命令列執行 GCC](https://www.youtube.com/results?search_query=Set+environment+variables+to+run+GCC+from+command+line)。

`或者使用其他 GCC`（除 Mingw-w64 之外，使用 [bitmap](https://github.com/vcaesar/bitmap) 時你需要自行編譯「libpng」。）

#### 其他所有平台：

```
GCC

帶 XTest 擴充功能的 X11（即 Xtst 函式庫）

「剪貼簿」：xsel xclip

「點陣圖」：libpng（僅「bitmap」使用。）

「事件-Gohook」：xcb, xkb, libxkbcommon（僅「hook」使用。）
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

Wayland 後端是 **純 Go（無 Cgo）** 實作，因此無需任何系統 C 函式庫。它需要一個
基於 wlroots 的合成器（Sway、Hyprland、Wayfire 等），並支援以下協定：

```
zwlr_virtual_pointer_v1            (滑鼠控制)
zwp_virtual_keyboard_v1            (鍵盤控制)
zwlr_screencopy_v1                 (螢幕擷取)
zwlr_foreign_toplevel_management_v1 (視窗管理)
```

GNOME 和 KDE **不**原生支援這些協定。

#### libei（GNOME / KDE）

libei 後端同樣是 **純 Go（無 Cgo）** 實作。它透過 freedesktop 的
`xdg-desktop-portal` RemoteDesktop 介面驅動輸入，因此可在 GNOME 和 KDE 上運作
（不同於 wlroots Wayland 後端）。它需要：

```
xdg-desktop-portal               (portal D-Bus 服務)
xdg-desktop-portal-gnome / -kde  (你桌面的 portal 後端)
```

注意：libei 後端僅處理滑鼠和鍵盤輸入。螢幕擷取和視窗管理會返回 `ErrNotSupported`。

## Installation:

在支援 Go module 的情況下（Go 1.11+），只需 import：

```go
import "github.com/go-vgo/robotgo"
```

否則，執行以下命令安裝 robotgo 套件：

```
go get github.com/go-vgo/robotgo
```

png.h: No such file or directory？請參閱 [issues/47](https://github.com/go-vgo/robotgo/issues/47)。

## Update:

```
go get -u github.com/go-vgo/robotgo
```

注意 go1.10.x 的 C 檔案編譯快取問題，[golang #24355](https://github.com/golang/go/issues/24355)。
`go mod vendor` 問題，[golang #26366](https://github.com/golang/go/issues/26366)。

## Cgo-free Builds:

RobotGo 為 Windows、macOS、X11、Wayland 和 libei（Linux）提供了 **純 Go（無 Cgo）** 後端，目前為實驗性功能。
它們暴露相同的 `robotgo` API，因此你的程式碼無需改動 —— 只需一個構建標籤。
這些後端可在 `CGO_ENABLED=0` 下交叉編譯（無需 GCC、MinGW、Xcode 或 X11 頭檔案）。

| 後端                             | 構建標籤  | Go 套件                             |
| -------------------------------- | --------- | ----------------------------------- |
| Windows（無 Cgo）                | `win`     | `github.com/go-vgo/robotgo/win`     |
| macOS（透過 purego 呼叫 Quartz） | `mac`     | `github.com/go-vgo/robotgo/darwin`  |
| X11（Linux，純 Go X 協定）       | `x11`     | `github.com/go-vgo/robotgo/x11`     |
| Wayland（Linux，wlroots）        | `wayland` | `github.com/go-vgo/robotgo/wayland` |
| libei（Linux，GNOME/KDE portal） | `libei`   | `github.com/go-vgo/robotgo/libei`   |
| 純 Go 預設（所有平台）           | `purego`  | 選擇上面的 `mac`/`win`/`wayland`    |

```sh
# 每個平台的純 Go 預設後端，一個標籤適用於所有目標：
# macOS -> mac，Windows -> win，Linux -> wayland（可與 x11/libei 組合以覆蓋）
go build -tags purego .
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -tags "purego,x11" .

# Windows，無需 Cgo / 無需 MinGW
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -tags win .

# macOS，透過 purego 在執行時載入 Quartz/CoreGraphics（無需 Xcode）
CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -tags mac .

# X11，純 Go 實作的 X 協定（XTEST）—— 無需 X11 頭檔案
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -tags x11 .

# Wayland，基於 wlroots 的合成器（Sway、Hyprland、Wayfire 等）
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -tags wayland .

# libei，透過 xdg-desktop-portal RemoteDesktop 支援 GNOME/KDE
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -tags libei .
```

注意：上面的範例建置模組根目錄（`.`）而不是 `./...`，因為 `examples/` 和一些特定平台的子套件使用了僅在預設 Cgo 後端下可用的 API。

在 `win` 標籤下，預設的 Cgo/Win32 後端被排除，呼叫轉發到純 Go 的 `win` 套件；
在 `mac` 標籤下，預設的 Cgo/Quartz 後端被排除，呼叫轉發到純 Go 的 `darwin` 套件（視窗管理回報 `ErrNotSupported`）；
在 `x11` 標籤下，Cgo/X11 後端被排除，呼叫轉發到純 Go 的 `x11` 套件；
在 `wayland` 標籤下，Cgo/X11 後端被排除，呼叫轉發到純 Go 的 `wayland` 套件；
在 `libei` 標籤下，Cgo/X11 和 wlroots Wayland 後端均被排除，呼叫轉發到純 Go 的 `libei` 套件。

`purego` 標籤是一個跨平台快捷方式：它會在所有平台排除 Cgo 後端，並按目標 OS 選擇預設純 Go 後端 —— macOS 使用 `mac`，Windows 使用 `win`，Linux 使用 `wayland`。在 Linux 上，你可以將它與 `x11` 或 `libei` 組合（例如 `-tags "purego,libei"`）以選擇不同的純 Go 後端。

## [Examples:](https://github.com/go-vgo/robotgo/blob/master/examples)

#### [滑鼠](https://github.com/go-vgo/robotgo/blob/master/examples/mouse/main.go)

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

#### [鍵盤](https://github.com/go-vgo/robotgo/blob/master/examples/key/main.go)

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

#### [螢幕](https://github.com/go-vgo/robotgo/blob/master/examples/screen/main.go)

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

#### [點陣圖](https://github.com/vcaesar/bitmap/blob/main/examples/main.go)

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

#### [視窗](https://github.com/go-vgo/robotgo/blob/master/examples/window/main.go)

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
- [維護者](https://github.com/orgs/go-vgo/people)

## Plans

- 將部分 C 程式碼重構為 Go（例如 x11、windows）
- 更好的多螢幕支援
- Wayland 支援
- 更新視窗控制代碼
- 嘗試支援 Android 和 iOS

## Contributors

- 完整的貢獻者列表請見[貢獻者頁面](https://github.com/go-vgo/robotgo/graphs/contributors)。
- 請參閱[貢獻指南](https://github.com/go-vgo/robotgo/blob/master/CONTRIBUTING.md)。

## License

Robotgo 主要依據「Apache License (Version 2.0)」的條款進行散布，部分內容受各類 BSD 風格授權條款約束。

詳見 [LICENSE-APACHE](http://www.apache.org/licenses/LICENSE-2.0)、[LICENSE](https://github.com/go-vgo/robotgo/blob/master/LICENSE)。
