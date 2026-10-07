# Robotgo

<p align="center">
  <img src="../docs/robotgo-logo.svg" width="560" alt="RobotGo 標誌 — 帶滑鼠指標的機器人" />
</p>

[![Build Status](https://github.com/go-vgo/robotgo/workflows/Go/badge.svg)](https://github.com/go-vgo/robotgo/commits/master)
[![CircleCI Status](https://circleci.com/gh/go-vgo/robotgo.svg?style=shield)](https://circleci.com/gh/go-vgo/robotgo)
[![codecov](https://codecov.io/gh/go-vgo/robotgo/branch/master/graph/badge.svg)](https://codecov.io/gh/go-vgo/robotgo)
[![golangci-lint](https://github.com/go-vgo/robotgo/actions/workflows/lint.yml/badge.svg)](https://github.com/go-vgo/robotgo/actions/workflows/lint.yml)
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
- [無 Cgo 建置](#cgo-free-builds)
- [安裝](#installation)
- [更新](#update)
- [範例](#examples)
- [型別轉換與按鍵](https://github.com/go-vgo/robotgo/blob/master/docs/keys.md)
- [交叉編譯](https://github.com/go-vgo/robotgo/blob/master/docs/install.md#crosscompiling)
- [作者](#authors)
- [計畫](#plans)
- [授權](#license)

## Docs

- [GoDoc](https://godoc.org/github.com/go-vgo/robotgo) <br>
- [API 文件](https://github.com/go-vgo/robotgo/blob/master/docs/doc.md)（已棄用，不再更新）

## Binding

[ADB](https://github.com/vcaesar/adb)，封裝的 Android adb API。

## Requirements

**Go 1.26+**（見 [go.mod](../go.mod)）以及一個後端：

- **預設（Cgo）：** `CGO_ENABLED=1`、一個 C 編譯器，以及下列平台函式庫。
- **純 Go（實驗性）：** 無需 C 工具鏈；見[無 Cgo 建置](#cgo-free-builds)。

兩者都需要為你所用功能安裝對應的[功能相關依賴](#feature-specific-dependencies)。

### Default Cgo setup

#### macOS

Go 和 Xcode 命令列工具（Clang）：

```sh
brew install go
xcode-select --install
```

**權限（Cgo 與純 Go）：** 在 **系統設定 > 隱私權與安全性** 中，為應用程式或終端機
授予**輔助使用**（輸入）和**螢幕錄製**（擷取）權限。

#### Windows

```powershell
winget install GoLang.Go
```

此外還需**一個** C 工具鏈。[LLVM-MinGW](https://github.com/mstorsjo/llvm-mingw)：

```powershell
winget install MartinStorsjo.LLVM-MinGW.UCRT
$env:CC = "clang"
```

或 [MinGW-w64 (WinLibs)](https://winlibs.com/)：

```powershell
winget install BrechtSanders.WinLibs.POSIX.UCRT
$env:CC = "gcc"
```

工具鏈的 `bin` 目錄（例如 `C:\mingw64\bin`）必須位於 `PATH` 中且與 `GOARCH`
相符。其他任何相容 Cgo 的編譯器同樣可用，例如手動安裝的
[MinGW-w64](https://sourceforge.net/projects/mingw-w64/files)。

[Bitmap](https://github.com/vcaesar/bitmap) 僅為 MinGW-w64 隨附 libpng；
使用其他工具鏈時，需自行建置 libpng。

#### Linux (X11)

建置：GCC、libc 標頭檔、X11/XTest 標頭檔。執行：帶 XTEST 的 X 伺服器並設定
`DISPLAY`。原生 Wayland 請見[無 Cgo 建置](#cgo-free-builds)。

**Ubuntu / Debian：**

```sh
# Go（發行版自帶的套件過舊時用 Snap）
sudo snap install go --classic
# 或：sudo apt install golang

# 核心（Cgo）：編譯器、libc、X11、XTest
sudo apt install gcc libc6-dev libx11-dev xorg-dev libxtst-dev

# 剪貼簿（X11）
sudo apt install xsel xclip

# Bitmap：libpng
sudo apt install libpng++-dev

# GoHook：XCB、XKB
sudo apt install xcb libxcb-xkb-dev x11-xkb-utils libx11-xcb-dev libxkbcommon-x11-dev libxkbcommon-dev
```

**Fedora：**

```sh
# 核心（Cgo）：編譯器、libc、X11、XTest
sudo dnf install gcc glibc-devel libX11-devel libXtst-devel

# 剪貼簿（X11）
sudo dnf install xsel xclip

# Bitmap：libpng
sudo dnf install libpng-devel

# GoHook：XKB（Fedora < 34：用 xorg-x11-xkb-utils-devel 取代 xkbcomp-devel）
sudo dnf install libxkbcommon-devel libxkbcommon-x11-devel xkbcomp-devel
```

### Feature-specific dependencies

在核心環境之上，僅在使用相應功能時才需要：

- **Linux 剪貼簿（Cgo 與純 Go）：** X11 上需 `xclip` 或 `xsel`；
  Wayland 上需 `wl-clipboard`：

  ```sh
  sudo apt install wl-clipboard   # Ubuntu / Debian
  sudo dnf install wl-clipboard   # Fedora
  ```

- **[Bitmap](https://github.com/vcaesar/bitmap)：** libpng 標頭檔。
- **[GoHook](https://github.com/robotn/gohook)：** Linux 上的 XCB/XKB 標頭檔。

Bitmap 和 GoHook 是獨立的 Cgo 套件；純 Go 的 RobotGo 後端並不會消除它們對 C 的
依賴。

## Cgo-free Builds

**實驗性的純 Go 後端**可在 `CGO_ENABLED=0` 下建置與交叉編譯（無需 GCC、MinGW、
Xcode 或 X11 標頭檔）。匯入路徑和 API 不變；後端由**建置標籤**選擇 ——
僅設定 `CGO_ENABLED=0` 不會選擇任何後端。

| 目標 `GOOS`   | 建置標籤  | 套件                   | 執行時需求                                                       |
| ------------- | --------- | ---------------------- | ---------------------------------------------------------------- |
| `windows`     | `win`     | [win](../win/)         | 僅需 Win32，無額外執行時需求                                     |
| `darwin`      | `mac`     | [darwin](../darwin/)   | 系統框架與[隱私權限](#macos)                                     |
| `linux`       | `x11`     | [x11](../x11/)         | 帶 XTEST 的 X 伺服器，已設定 `DISPLAY`                           |
| `linux`       | `wayland` | [wayland](../wayland/) | wlroots 合成器，見 [Wayland](#wayland)                           |
| `linux`       | `libei`   | [libei](../libei/)     | RemoteDesktop portal（GNOME/KDE），見 [libei](#libei-gnome--kde) |

**`purego`** 是一個捷徑：macOS 上為 `mac`，Windows 上為 `win`，Linux 上為
`wayland`。在 Linux 上可追加 `x11` 或 `libei` 來覆寫（`-tags "purego,x11"`）。

只使用**一個**與 `GOOS` 相符的後端標籤；諸如 `x11,libei` 的組合無法編譯
（`purego` 加一個 Linux 覆寫標籤是唯一例外）。

### Build commands

```sh
# 當前作業系統
CGO_ENABLED=0 go build -tags purego .

# 交叉編譯
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -tags win .
CGO_ENABLED=0 GOOS=darwin  GOARCH=arm64 go build -tags mac .
CGO_ENABLED=0 GOOS=linux   GOARCH=amd64 go build -tags x11 .
CGO_ENABLED=0 GOOS=linux   GOARCH=amd64 go build -tags wayland .
CGO_ENABLED=0 GOOS=linux   GOARCH=amd64 go build -tags libei .
```

PowerShell：`$env:CGO_ENABLED = "0"; go build -tags purego .`

請建置模組根目錄（`.`）或你自己的套件，而不是 `./...`：`examples/` 和部分子套件
需要僅 Cgo 可用的 API。

### Limitations

上述執行時需求依然適用。僅 Cgo 可用的 API，如 `CaptureScreen` / `FreeBitmap`，
並不存在；請使用 `CaptureImg`。不支援的操作會回傳 `ErrNotSupported`（或空結果）。

- **`mac`：** 無視窗管理 —— `ActiveName` 回傳 `ErrNotSupported`，
  `GetTitle` 回傳 `""`，最小化/最大化/關閉為空操作。
- **`wayland`、`libei`：** `Location()` 回傳 RobotGo 最後注入的位置，而非實際
  游標位置。剪貼簿需要 `wl-clipboard`。

#### Wayland

需要一個 wlroots 合成器（Sway、Hyprland、Wayfire 等），已設定 `WAYLAND_DISPLAY`
和 `XDG_RUNTIME_DIR`，並公開：

| 功能     | 協定全域物件                       |
| -------- | ---------------------------------- |
| 滑鼠控制 | `zwlr_virtual_pointer_manager_v1`  |
| 鍵盤控制 | `zwp_virtual_keyboard_manager_v1`  |
| 螢幕擷取 | `zwlr_screencopy_manager_v1`       |
| 視窗管理 | `zwlr_foreign_toplevel_manager_v1` |

GNOME 和 KDE 不提供這些協定；請在其上使用 `libei`。

#### libei (GNOME / KDE)

輸入透過 `xdg-desktop-portal` 的 **RemoteDesktop** D-Bus 介面完成：需要工作階段
匯流排、`xdg-desktop-portal` 以及一個實作該介面的後端
（`xdg-desktop-portal-gnome`、`-kde` 或 `-wlr`）。首次執行會顯示授權對話框；
還原權杖會快取在 `$XDG_STATE_HOME/robotgo` 下。

絕對移動使用連結的 ScreenCast 串流（預設）；設定 `libei.LinkScreenCast = false`
則僅使用相對移動。螢幕擷取和視窗管理會回報 `ErrNotSupported`。

## Installation

```sh
go get github.com/go-vgo/robotgo
```

```go
import "github.com/go-vgo/robotgo"
```

使用 Bitmap 時出現 `png.h: No such file or directory`：請查看其 libpng 需求以及
[issues/47](https://github.com/go-vgo/robotgo/issues/47)。

## Update

```sh
go get -u github.com/go-vgo/robotgo
```

注意 go1.10.x 的 C 檔案編譯快取問題，[golang #24355](https://github.com/golang/go/issues/24355)。
`go mod vendor` 問題，[golang #26366](https://github.com/golang/go/issues/26366)。

## [Examples](https://github.com/go-vgo/robotgo/blob/master/examples)

#### [滑鼠](https://github.com/go-vgo/robotgo/blob/master/examples/mouse/main.go)

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

  // 檢查回傳的錯誤, 無效參數會回傳錯誤
  if err := robotgo.Click("left", "double"); err != nil {
    fmt.Println("robotgo.Click error:", err)
  }
  if err := robotgo.ScrollDir(10, "forward"); err != nil {
    fmt.Println("robotgo.ScrollDir error:", err)
  }

  // MoveSmoothRelative 將平滑移動失敗回報為 ErrSmoothMove
  err := robotgo.MoveSmoothRelative(10, -10)
  if errors.Is(err, robotgo.ErrSmoothMove) {
    fmt.Println("smooth move failed:", err)
  }

  // 純 Go 後端對不支援的操作回傳 ErrNotSupported
  err = robotgo.Move(100, 200)
  if errors.Is(err, robotgo.ErrNotSupported) {
    fmt.Println("robotgo.Move is not supported by this backend")
  }
}
```

#### [鍵盤](https://github.com/go-vgo/robotgo/blob/master/examples/key/main.go)

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

  // 未知的按鍵名稱會回傳錯誤, 而不是按下錯誤的按鍵
  if err := robotgo.KeyTap("notakey"); err != nil {
    fmt.Println("robotgo.KeyTap error:", err)
  }

  // TypeStr 回傳第一個輸入錯誤, Type 只回傳已輸入的字元數
  if err := robotgo.TypeStr("Hello, 世界"); err != nil {
    fmt.Println("robotgo.TypeStr error:", err)
  }

  // 按住修飾鍵, 然後一律放開它並檢查兩個錯誤
  if err := robotgo.KeyDown("shift"); err == nil {
    err = errors.Join(robotgo.KeyTap("a"), robotgo.KeyUp("shift"))
    if err != nil {
      fmt.Println("shift + a error:", err)
    }
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

- 更好的多螢幕支援
- 更新視窗控制代碼
- 嘗試支援 Android 和 iOS

## Contributors

- 完整的貢獻者列表請見[貢獻者頁面](https://github.com/go-vgo/robotgo/graphs/contributors)。
- 請參閱[貢獻指南](https://github.com/go-vgo/robotgo/blob/master/CONTRIBUTING.md)。

## License

Robotgo 主要依據「Apache License (Version 2.0)」的條款進行散布，部分內容受各類 BSD 風格授權條款約束。

詳見 [LICENSE-APACHE](http://www.apache.org/licenses/LICENSE-2.0)、[LICENSE](https://github.com/go-vgo/robotgo/blob/master/LICENSE)。
