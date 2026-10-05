# Robotgo

<p align="center">
  <img src="docs/robotgo-logo.svg" width="560" alt="RobotGo logo — a robot with a mouse pointer" />
</p>

<!-- <img align="right" src="https://raw.githubusercontent.com/go-vgo/robotgo/master/logo.jpg"> -->
<!-- [![codecov](https://codecov.io/gh/go-vgo/robotgo/branch/master/graph/badge.svg)](https://codecov.io/gh/go-vgo/robotgo) -->

[![Build Status](https://github.com/go-vgo/robotgo/workflows/Go/badge.svg)](https://github.com/go-vgo/robotgo/commits/master)
[![CircleCI Status](https://circleci.com/gh/go-vgo/robotgo.svg?style=shield)](https://circleci.com/gh/go-vgo/robotgo)
[![golangci-lint](https://github.com/go-vgo/robotgo/actions/workflows/lint.yml/badge.svg)](https://github.com/go-vgo/robotgo/actions/workflows/lint.yml)
[![GoDoc](https://pkg.go.dev/badge/github.com/go-vgo/robotgo?status.svg)](https://pkg.go.dev/github.com/go-vgo/robotgo?tab=doc)
[![GitHub release](https://img.shields.io/github/release/go-vgo/robotgo.svg)](https://github.com/go-vgo/robotgo/releases/latest)
<a href="https://discord.gg/npPb3NzE4A"><img src="https://img.shields.io/discord/1484658282777018551.svg?logo=discord&logoColor=white&label=Discord&color=5865F2" alt="Join the Discord chat at https://discord.gg/npPb3NzE4A"></a>

English | [简体中文](lang/README.zh.md) | [繁體中文](lang/README.zht.md) | [日本語](lang/README.ja.md) | [한국어](lang/README.ko.md) | [Français](lang/README.fr.md) | [Deutsch](lang/README.de.md) | [Español](lang/README.es.md) | [Русский](lang/README.ru.md) | [Português](lang/README.pt.md)

> Golang Desktop Automation, auto test and AI Computer Use. <br>
> Control the mouse, keyboard, read the screen, process, Window Handle, image and bitmap and global event listener.

RobotGo supports Mac, Windows, and Linux; and robotgo supports arm64 and x86-amd64.

I build [Codg](https://github.com/vcaesar/codg) now, Easy code and work AI agent system: auto, asynchronous, concurrency, efficiently and High accuracy

<p align="center">
<a href="https://github.com/vcaesar/codg" rel="nofollow">
<!-- <img width="800" alt="Codg Demo" src="https://github.com/vcaesar/codg/raw/main/demo/26-04.png" /> -->
<img width="800" alt="Codg Demo" src="https://github.com/vcaesar/codg/raw/main/demo/26-04-1.png" />
</a>
</p>

[RobotGo-Pro](https://github.com/vcaesar/robotgo-pro) get the JavaScript, Python, Lua and others version, tech supports, new features and newest robotgo version ("no open-source version now").

## Contents

- [Docs](#docs)
- [Binding](#binding)
- [Requirements](#requirements)
- [Cgo-free Builds](#cgo-free-builds)
- [Installation](#installation)
- [Update](#update)
- [Examples](#examples)
- [Type Conversion and keys](https://github.com/go-vgo/robotgo/blob/master/docs/keys.md)
- [Cross-Compiling](https://github.com/go-vgo/robotgo/blob/master/docs/install.md#crosscompiling)
- [Authors](#authors)
- [Plans](#plans)
- [License](#license)

## Docs

- [GoDoc](https://godoc.org/github.com/go-vgo/robotgo) <br>
- [API Docs](https://github.com/go-vgo/robotgo/blob/master/docs/doc.md) (Deprecated, no updated)

## Binding

[ADB](https://github.com/vcaesar/adb), packaging android adb API.

## Requirements

**Go 1.26+** (see [go.mod](go.mod)) and one backend:

- **Default (Cgo):** `CGO_ENABLED=1`, a C compiler and the platform libraries below.
- **Pure Go (experimental):** no C toolchain; see [Cgo-free Builds](#cgo-free-builds).

Both need the [feature-specific dependencies](#feature-specific-dependencies)
for the features you use.

### Default Cgo setup

#### macOS

Go and the Xcode Command Line Tools (Clang):

```sh
brew install go
xcode-select --install
```

**Permissions (Cgo and pure Go):** in **System Settings > Privacy & Security**,
grant the app or terminal **Accessibility** (input) and **Screen Recording**
(capture).

#### Windows

```powershell
winget install GoLang.Go
```

Plus **one** C toolchain. [LLVM-MinGW](https://github.com/mstorsjo/llvm-mingw):

```powershell
winget install MartinStorsjo.LLVM-MinGW.UCRT
$env:CC = "clang"
```

Or [MinGW-w64 (WinLibs)](https://winlibs.com/):

```powershell
winget install BrechtSanders.WinLibs.POSIX.UCRT
$env:CC = "gcc"
```

The toolchain's `bin` (e.g. `C:\mingw64\bin`) must be on `PATH` and match
`GOARCH`. Any other Cgo-compatible compiler, such as a manual
[MinGW-w64](https://sourceforge.net/projects/mingw-w64/files) install, also works.

[Bitmap](https://github.com/vcaesar/bitmap) ships libpng for MinGW-w64 only;
with another toolchain, build libpng yourself.

#### Linux (X11)

Build: GCC, libc headers, X11/XTest headers. Runtime: an X server with XTEST
and `DISPLAY` set. For native Wayland, see [Cgo-free Builds](#cgo-free-builds).

**Ubuntu / Debian:**

```sh
# Go (Snap if the distro package is too old)
sudo snap install go --classic
# or: sudo apt install golang

# Core (Cgo): compiler, libc, X11, XTest
sudo apt install gcc libc6-dev libx11-dev xorg-dev libxtst-dev

# Clipboard (X11)
sudo apt install xsel xclip

# Bitmap: libpng
sudo apt install libpng++-dev

# GoHook: XCB, XKB
sudo apt install xcb libxcb-xkb-dev x11-xkb-utils libx11-xcb-dev libxkbcommon-x11-dev libxkbcommon-dev
```

**Fedora:**

```sh
# Core (Cgo): compiler, libc, X11, XTest
sudo dnf install gcc glibc-devel libX11-devel libXtst-devel

# Clipboard (X11)
sudo dnf install xsel xclip

# Bitmap: libpng
sudo dnf install libpng-devel

# GoHook: XKB (Fedora < 34: xorg-x11-xkb-utils-devel instead of xkbcomp-devel)
sudo dnf install libxkbcommon-devel libxkbcommon-x11-devel xkbcomp-devel
```

### Feature-specific dependencies

Needed on top of the core setup, only for the features you use:

- **Linux clipboard (Cgo and pure Go):** `xclip` or `xsel` on X11;
  `wl-clipboard` on Wayland:

  ```sh
  sudo apt install wl-clipboard   # Ubuntu / Debian
  sudo dnf install wl-clipboard   # Fedora
  ```

- **[Bitmap](https://github.com/vcaesar/bitmap):** libpng headers.
- **[GoHook](https://github.com/robotn/gohook):** XCB/XKB headers on Linux.

Bitmap and GoHook are separate Cgo packages; a pure-Go RobotGo backend does not
remove their C requirements.

## Cgo-free Builds

**Experimental pure-Go backends** build and cross-compile with `CGO_ENABLED=0`
(no GCC, MinGW, Xcode or X11 headers). Same import path and API; a **build tag**
selects the backend — `CGO_ENABLED=0` alone selects nothing.

| Target `GOOS` | Build tag | Package             | Runtime requirements                                             |
| ------------- | --------- | ------------------- | ---------------------------------------------------------------- |
| `windows`     | `win`     | [win](win/)         | Win32 only, no extra runtime                                     |
| `darwin`      | `mac`     | [darwin](darwin/)   | System frameworks and [privacy permissions](#macos)              |
| `linux`       | `x11`     | [x11](x11/)         | X server with XTEST, `DISPLAY` set                               |
| `linux`       | `wayland` | [wayland](wayland/) | wlroots compositor, see [Wayland](#wayland)                      |
| `linux`       | `libei`   | [libei](libei/)     | RemoteDesktop portal (GNOME/KDE), see [libei](#libei-gnome--kde) |

**`purego`** is a shortcut: `mac` on macOS, `win` on Windows, `wayland` on
Linux. Add `x11` or `libei` to override on Linux (`-tags "purego,x11"`).

Use **one backend tag** matching `GOOS`; combinations such as `x11,libei` fail
to compile (`purego` plus one Linux override is the only exception).

### Build commands

```sh
# Current OS
CGO_ENABLED=0 go build -tags purego .

# Cross-compile
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -tags win .
CGO_ENABLED=0 GOOS=darwin  GOARCH=arm64 go build -tags mac .
CGO_ENABLED=0 GOOS=linux   GOARCH=amd64 go build -tags x11 .
CGO_ENABLED=0 GOOS=linux   GOARCH=amd64 go build -tags wayland .
CGO_ENABLED=0 GOOS=linux   GOARCH=amd64 go build -tags libei .
```

PowerShell: `$env:CGO_ENABLED = "0"; go build -tags purego .`

Build the module root (`.`) or your own package, not `./...`: `examples/` and
some subpackages need Cgo-only APIs.

### Limitations

The runtime requirements above still apply. Cgo-only APIs such as
`CaptureScreen` / `FreeBitmap` do not exist; use `CaptureImg`. Unsupported
operations return `ErrNotSupported` (or an empty result).

- **`mac`:** no window management — `ActiveName` returns `ErrNotSupported`,
  `GetTitle` returns `""`, minimize/maximize/close are no-ops.
- **`wayland`, `libei`:** `Location()` returns the last position injected by
  RobotGo, not the physical cursor. Clipboard needs `wl-clipboard`.

#### Wayland

A wlroots compositor (Sway, Hyprland, Wayfire, ...) with `WAYLAND_DISPLAY` and
`XDG_RUNTIME_DIR` set, exposing:

| Feature           | Protocol global                    |
| ----------------- | ---------------------------------- |
| Mouse control     | `zwlr_virtual_pointer_manager_v1`  |
| Keyboard control  | `zwp_virtual_keyboard_manager_v1`  |
| Screen capture    | `zwlr_screencopy_manager_v1`       |
| Window management | `zwlr_foreign_toplevel_manager_v1` |

GNOME and KDE do not provide these protocols; use `libei` there.

#### libei (GNOME / KDE)

Input goes through the `xdg-desktop-portal` **RemoteDesktop** D-Bus interface:
needs a session bus, `xdg-desktop-portal` and a backend implementing it
(`xdg-desktop-portal-gnome`, `-kde` or `-wlr`). The first run shows a consent
dialog; the restore token is cached under `$XDG_STATE_HOME/robotgo`.

Absolute motion uses a linked ScreenCast stream (default); set
`libei.LinkScreenCast = false` for relative motion only. Screen capture and
window management report `ErrNotSupported`.

## Installation

```sh
go get github.com/go-vgo/robotgo
```

```go
import "github.com/go-vgo/robotgo"
```

`png.h: No such file or directory` with Bitmap: see its libpng requirements and
[issues/47](https://github.com/go-vgo/robotgo/issues/47).

## Update

```sh
go get -u github.com/go-vgo/robotgo
```

Note go1.10.x C file compilation cache problem, [golang #24355](https://github.com/golang/go/issues/24355).
`go mod vendor` problem, [golang #26366](https://github.com/golang/go/issues/26366).

## [Examples](https://github.com/go-vgo/robotgo/blob/master/examples)

#### [Mouse](https://github.com/go-vgo/robotgo/blob/master/examples/mouse/main.go)

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

#### [Keyboard](https://github.com/go-vgo/robotgo/blob/master/examples/key/main.go)

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

#### [Screen](https://github.com/go-vgo/robotgo/blob/master/examples/screen/main.go)

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

- [The author is Evans](https://github.com/vcaesar)
- [Maintainers](https://github.com/orgs/go-vgo/people)

## Plans

- Better multiscreen support
- Update Window Handle
- Try to support Android and IOS

## Contributors

- See [contributors page](https://github.com/go-vgo/robotgo/graphs/contributors) for full list of contributors.
- See [Contribution Guidelines](https://github.com/go-vgo/robotgo/blob/master/CONTRIBUTING.md).

## License

Robotgo is primarily distributed under the terms of "the Apache License (Version 2.0)", with portions covered by various BSD-like licenses.

See [LICENSE-APACHE](http://www.apache.org/licenses/LICENSE-2.0), [LICENSE](https://github.com/go-vgo/robotgo/blob/master/LICENSE).
