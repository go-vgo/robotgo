# Robotgo

<p align="center">
  <img src="../docs/robotgo-logo.svg" width="560" alt="RobotGo 로고 — 마우스 포인터가 있는 로봇" />
</p>

[![Build Status](https://github.com/go-vgo/robotgo/workflows/Go/badge.svg)](https://github.com/go-vgo/robotgo/commits/master)
[![CircleCI Status](https://circleci.com/gh/go-vgo/robotgo.svg?style=shield)](https://circleci.com/gh/go-vgo/robotgo)
[![codecov](https://codecov.io/gh/go-vgo/robotgo/branch/master/graph/badge.svg)](https://codecov.io/gh/go-vgo/robotgo)
[![golangci-lint](https://github.com/go-vgo/robotgo/actions/workflows/lint.yml/badge.svg)](https://github.com/go-vgo/robotgo/actions/workflows/lint.yml)
[![GoDoc](https://pkg.go.dev/badge/github.com/go-vgo/robotgo?status.svg)](https://pkg.go.dev/github.com/go-vgo/robotgo?tab=doc)
[![GitHub release](https://img.shields.io/github/release/go-vgo/robotgo.svg)](https://github.com/go-vgo/robotgo/releases/latest)
<a href="https://discord.gg/npPb3NzE4A"><img src="https://img.shields.io/discord/1484658282777018551.svg?logo=discord&logoColor=white&label=Discord&color=5865F2" alt="Join the Discord chat at https://discord.gg/npPb3NzE4A"></a>

[English](../README.md) | [简体中文](README.zh.md) | [繁體中文](README.zht.md) | [日本語](README.ja.md) | 한국어 | [Français](README.fr.md) | [Deutsch](README.de.md) | [Español](README.es.md) | [Русский](README.ru.md) | [Português](README.pt.md)

> Golang 데스크톱 자동화, 자동 테스트 및 AI 컴퓨터 사용(Computer Use). <br>
> 마우스와 키보드 제어, 화면 읽기, 프로세스, 윈도우 핸들, 이미지와 비트맵, 그리고 전역 이벤트 리스너.

RobotGo는 Mac, Windows, Linux를 지원하며, arm64와 x86-amd64 아키텍처도 지원합니다.

저는 지금 [Codg](https://github.com/vcaesar/codg)를 만들고 있습니다. 간편하게 코딩하고 작업할 수 있는 AI 에이전트(Agent) 시스템으로, 자동화, 비동기, 동시성, 고효율 그리고 높은 정확도를 갖추고 있습니다.

<p align="center">
<a href="https://github.com/vcaesar/codg" rel="nofollow">
<img width="800" alt="Codg Demo" src="https://github.com/vcaesar/codg/raw/main/demo/26-04-1.png" />
</a>
</p>

[RobotGo-Pro](https://github.com/vcaesar/robotgo-pro)는 JavaScript, Python, Lua 등 다른 언어 버전과 기술 지원, 새로운 기능, 그리고 최신 robotgo 버전("현재 오픈소스 버전 없음")을 제공합니다.

## 목차

- [문서](#docs)
- [바인딩](#binding)
- [요구 사항](#requirements)
- [Cgo 없는 빌드](#cgo-free-builds)
- [설치](#installation)
- [업데이트](#update)
- [예제](#examples)
- [타입 변환과 키](https://github.com/go-vgo/robotgo/blob/master/docs/keys.md)
- [크로스 컴파일](https://github.com/go-vgo/robotgo/blob/master/docs/install.md#crosscompiling)
- [작성자](#authors)
- [계획](#plans)
- [라이선스](#license)

## Docs

- [GoDoc](https://godoc.org/github.com/go-vgo/robotgo) <br>
- [API 문서](https://github.com/go-vgo/robotgo/blob/master/docs/doc.md) (지원 중단, 더 이상 업데이트되지 않음)

## Binding

[ADB](https://github.com/vcaesar/adb), Android adb API를 래핑한 패키지.

## Requirements

**Go 1.26+**([go.mod](../go.mod) 참조)와 백엔드 하나가 필요합니다:

- **기본(Cgo):** `CGO_ENABLED=1`, C 컴파일러, 그리고 아래의 플랫폼 라이브러리.
- **순수 Go(실험적):** C 툴체인 불필요. [Cgo 없는 빌드](#cgo-free-builds) 참조.

두 방식 모두 사용하는 기능에 맞는 [기능별 의존성](#feature-specific-dependencies)이
필요합니다.

### Default Cgo setup

#### macOS

Go와 Xcode 명령줄 도구(Clang):

```sh
brew install go
xcode-select --install
```

**권한(Cgo와 순수 Go 공통):** **시스템 설정 > 개인정보 보호 및 보안**에서 앱 또는
터미널에 **손쉬운 사용**(입력)과 **화면 기록**(캡처) 권한을 부여하세요.

#### Windows

```powershell
winget install GoLang.Go
```

추가로 C 툴체인 **하나**가 필요합니다. [LLVM-MinGW](https://github.com/mstorsjo/llvm-mingw):

```powershell
winget install MartinStorsjo.LLVM-MinGW.UCRT
$env:CC = "clang"
```

또는 [MinGW-w64 (WinLibs)](https://winlibs.com/):

```powershell
winget install BrechtSanders.WinLibs.POSIX.UCRT
$env:CC = "gcc"
```

툴체인의 `bin` 경로(예: `C:\mingw64\bin`)가 `PATH`에 있어야 하고 `GOARCH`와 일치해야
합니다. 수동으로 설치한 [MinGW-w64](https://sourceforge.net/projects/mingw-w64/files)
처럼 Cgo와 호환되는 다른 컴파일러도 사용할 수 있습니다.

[Bitmap](https://github.com/vcaesar/bitmap)은 MinGW-w64용 libpng만 포함합니다.
다른 툴체인을 쓸 때는 libpng를 직접 빌드하세요.

#### Linux (X11)

빌드: GCC, libc 헤더, X11/XTest 헤더. 런타임: XTEST를 지원하는 X 서버와 `DISPLAY`
설정. 네이티브 Wayland는 [Cgo 없는 빌드](#cgo-free-builds)를 참조하세요.

**Ubuntu / Debian:**

```sh
# Go (배포판 패키지가 너무 오래되었으면 Snap 사용)
sudo snap install go --classic
# 또는: sudo apt install golang

# 핵심(Cgo): 컴파일러, libc, X11, XTest
sudo apt install gcc libc6-dev libx11-dev xorg-dev libxtst-dev

# 클립보드(X11)
sudo apt install xsel xclip

# Bitmap: libpng
sudo apt install libpng++-dev

# GoHook: XCB, XKB
sudo apt install xcb libxcb-xkb-dev x11-xkb-utils libx11-xcb-dev libxkbcommon-x11-dev libxkbcommon-dev
```

**Fedora:**

```sh
# 핵심(Cgo): 컴파일러, libc, X11, XTest
sudo dnf install gcc glibc-devel libX11-devel libXtst-devel

# 클립보드(X11)
sudo dnf install xsel xclip

# Bitmap: libpng
sudo dnf install libpng-devel

# GoHook: XKB (Fedora 34 미만: xkbcomp-devel 대신 xorg-x11-xkb-utils-devel)
sudo dnf install libxkbcommon-devel libxkbcommon-x11-devel xkbcomp-devel
```

### Feature-specific dependencies

핵심 설정에 더해, 사용하는 기능에만 필요합니다:

- **Linux 클립보드(Cgo와 순수 Go 공통):** X11에서는 `xclip` 또는 `xsel`,
  Wayland에서는 `wl-clipboard`:

  ```sh
  sudo apt install wl-clipboard   # Ubuntu / Debian
  sudo dnf install wl-clipboard   # Fedora
  ```

- **[Bitmap](https://github.com/vcaesar/bitmap):** libpng 헤더.
- **[GoHook](https://github.com/robotn/gohook):** Linux에서는 XCB/XKB 헤더.

Bitmap과 GoHook은 별도의 Cgo 패키지입니다. 순수 Go RobotGo 백엔드를 써도 이들의 C
의존성은 없어지지 않습니다.

## Cgo-free Builds

**실험적인 순수 Go 백엔드**는 `CGO_ENABLED=0`으로 빌드 및 크로스 컴파일됩니다
(GCC, MinGW, Xcode, X11 헤더 불필요). import 경로와 API는 동일하며 **빌드 태그**로
백엔드를 선택합니다 — `CGO_ENABLED=0`만으로는 아무것도 선택되지 않습니다.

| 대상 `GOOS`   | 빌드 태그 | 패키지                 | 런타임 요구 사항                                                        |
| ------------- | --------- | ---------------------- | ----------------------------------------------------------------------- |
| `windows`     | `win`     | [win](../win/)         | Win32만 사용, 추가 런타임 없음                                          |
| `darwin`      | `mac`     | [darwin](../darwin/)   | 시스템 프레임워크와 [개인정보 권한](#macos)                             |
| `linux`       | `x11`     | [x11](../x11/)         | XTEST를 지원하는 X 서버, `DISPLAY` 설정                                 |
| `linux`       | `wayland` | [wayland](../wayland/) | wlroots 컴포지터, [Wayland](#wayland) 참조                              |
| `linux`       | `libei`   | [libei](../libei/)     | RemoteDesktop portal(GNOME/KDE), [libei](#libei-gnome--kde) 참조        |

**`purego`**는 단축 태그입니다: macOS에서는 `mac`, Windows에서는 `win`, Linux에서는
`wayland`. Linux에서 재정의하려면 `x11` 또는 `libei`를 추가하세요(`-tags "purego,x11"`).

`GOOS`에 맞는 **백엔드 태그 하나**만 사용하세요. `x11,libei` 같은 조합은 컴파일에
실패합니다(`purego`와 Linux 재정의 하나를 함께 쓰는 경우만 예외).

### Build commands

```sh
# 현재 OS
CGO_ENABLED=0 go build -tags purego .

# 크로스 컴파일
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -tags win .
CGO_ENABLED=0 GOOS=darwin  GOARCH=arm64 go build -tags mac .
CGO_ENABLED=0 GOOS=linux   GOARCH=amd64 go build -tags x11 .
CGO_ENABLED=0 GOOS=linux   GOARCH=amd64 go build -tags wayland .
CGO_ENABLED=0 GOOS=linux   GOARCH=amd64 go build -tags libei .
```

PowerShell: `$env:CGO_ENABLED = "0"; go build -tags purego .`

`./...`가 아니라 모듈 루트(`.`) 또는 직접 만든 패키지를 빌드하세요. `examples/`와
일부 하위 패키지는 Cgo 전용 API가 필요합니다.

### Limitations

위의 런타임 요구 사항은 그대로 적용됩니다. `CaptureScreen` / `FreeBitmap` 같은 Cgo
전용 API는 존재하지 않으므로 `CaptureImg`를 사용하세요. 지원되지 않는 동작은
`ErrNotSupported`를 반환합니다(또는 빈 결과).

- **`mac`:** 창 관리 없음 — `ActiveName`은 `ErrNotSupported`를 반환하고,
  `GetTitle`은 `""`를 반환하며, 최소화/최대화/닫기는 아무 동작도 하지 않습니다.
- **`wayland`, `libei`:** `Location()`은 실제 커서가 아니라 RobotGo가 마지막으로
  주입한 위치를 반환합니다. 클립보드에는 `wl-clipboard`가 필요합니다.

#### Wayland

`WAYLAND_DISPLAY`와 `XDG_RUNTIME_DIR`가 설정되어 있고 다음을 노출하는 wlroots
컴포지터(Sway, Hyprland, Wayfire 등):

| 기능        | 프로토콜 global                    |
| ----------- | ---------------------------------- |
| 마우스 제어 | `zwlr_virtual_pointer_manager_v1`  |
| 키보드 제어 | `zwp_virtual_keyboard_manager_v1`  |
| 화면 캡처   | `zwlr_screencopy_manager_v1`       |
| 창 관리     | `zwlr_foreign_toplevel_manager_v1` |

GNOME과 KDE는 이러한 프로토콜을 제공하지 않으므로, 그곳에서는 `libei`를 사용하세요.

#### libei (GNOME / KDE)

입력은 `xdg-desktop-portal`의 **RemoteDesktop** D-Bus 인터페이스를 통해 전달됩니다.
세션 버스, `xdg-desktop-portal`, 그리고 이를 구현하는 백엔드
(`xdg-desktop-portal-gnome`, `-kde`, `-wlr`)가 필요합니다. 첫 실행 시 동의 대화상자가
표시되며, 복원 토큰은 `$XDG_STATE_HOME/robotgo`에 캐시됩니다.

절대 좌표 이동은 연결된 ScreenCast 스트림을 사용합니다(기본값). 상대 이동만 쓰려면
`libei.LinkScreenCast = false`로 설정하세요. 화면 캡처와 창 관리는
`ErrNotSupported`를 반환합니다.

## Installation

```sh
go get github.com/go-vgo/robotgo
```

```go
import "github.com/go-vgo/robotgo"
```

Bitmap에서 `png.h: No such file or directory`가 발생하면 해당 프로젝트의 libpng 요구
사항과 [issues/47](https://github.com/go-vgo/robotgo/issues/47)을 참조하세요.

## Update

```sh
go get -u github.com/go-vgo/robotgo
```

go1.10.x의 C 파일 컴파일 캐시 문제에 주의하세요, [golang #24355](https://github.com/golang/go/issues/24355).
`go mod vendor` 문제, [golang #26366](https://github.com/golang/go/issues/26366).

## [Examples](https://github.com/go-vgo/robotgo/blob/master/examples)

#### [마우스](https://github.com/go-vgo/robotgo/blob/master/examples/mouse/main.go)

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

  // 반환된 오류를 확인합니다. 잘못된 인자는 오류를 반환합니다
  if err := robotgo.Click("left", "double"); err != nil {
    fmt.Println("robotgo.Click error:", err)
  }
  if err := robotgo.ScrollDir(10, "forward"); err != nil {
    fmt.Println("robotgo.ScrollDir error:", err)
  }

  // MoveSmoothRelative는 부드러운 이동 실패를 ErrSmoothMove로 보고합니다
  err := robotgo.MoveSmoothRelative(10, -10)
  if errors.Is(err, robotgo.ErrSmoothMove) {
    fmt.Println("smooth move failed:", err)
  }

  // 순수 Go 백엔드는 지원하지 않는 작업에 ErrNotSupported를 반환합니다
  err = robotgo.Move(100, 200)
  if errors.Is(err, robotgo.ErrNotSupported) {
    fmt.Println("robotgo.Move is not supported by this backend")
  }
}
```

#### [키보드](https://github.com/go-vgo/robotgo/blob/master/examples/key/main.go)

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

  // 알 수 없는 키 이름은 잘못된 키를 누르지 않고 오류를 반환합니다
  if err := robotgo.KeyTap("notakey"); err != nil {
    fmt.Println("robotgo.KeyTap error:", err)
  }

  // TypeStr는 첫 번째 입력 오류를 반환하고, Type은 입력한 문자 수만 반환합니다
  if err := robotgo.TypeStr("Hello, 世界"); err != nil {
    fmt.Println("robotgo.TypeStr error:", err)
  }

  // 수정 키를 누른 뒤 항상 해제하고 두 오류를 모두 확인합니다
  if err := robotgo.KeyDown("shift"); err == nil {
    err = errors.Join(robotgo.KeyTap("a"), robotgo.KeyUp("shift"))
    if err != nil {
      fmt.Println("shift + a error:", err)
    }
  }
}
```

#### [화면](https://github.com/go-vgo/robotgo/blob/master/examples/screen/main.go)

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

#### [비트맵](https://github.com/vcaesar/bitmap/blob/main/examples/main.go)

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

#### [이벤트](https://github.com/robotn/gohook/blob/master/examples/main.go)

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

#### [윈도우](https://github.com/go-vgo/robotgo/blob/master/examples/window/main.go)

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

- [작성자 Evans](https://github.com/vcaesar)
- [관리자](https://github.com/orgs/go-vgo/people)

## Plans

- 더 나은 멀티 스크린 지원
- 윈도우 핸들 업데이트
- Android 및 iOS 지원 시도

## Contributors

- 전체 기여자 목록은 [기여자 페이지](https://github.com/go-vgo/robotgo/graphs/contributors)를 참조하세요.
- [기여 가이드라인](https://github.com/go-vgo/robotgo/blob/master/CONTRIBUTING.md)을 참조하세요.

## License

Robotgo는 주로 "the Apache License (Version 2.0)" 조건에 따라 배포되며, 일부 내용은 다양한 BSD 계열 라이선스의 적용을 받습니다.

[LICENSE-APACHE](http://www.apache.org/licenses/LICENSE-2.0), [LICENSE](https://github.com/go-vgo/robotgo/blob/master/LICENSE)를 참조하세요.
