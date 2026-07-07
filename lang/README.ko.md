# Robotgo

[![Build Status](https://github.com/go-vgo/robotgo/workflows/Go/badge.svg)](https://github.com/go-vgo/robotgo/commits/master)
[![CircleCI Status](https://circleci.com/gh/go-vgo/robotgo.svg?style=shield)](https://circleci.com/gh/go-vgo/robotgo)
[![Go Report Card](https://goreportcard.com/badge/github.com/go-vgo/robotgo)](https://goreportcard.com/report/github.com/go-vgo/robotgo)
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
- [설치](#installation)
- [업데이트](#update)
- [Cgo 없는 빌드](#cgo-free-builds)
- [예제](#examples)
- [타입 변환과 키](https://github.com/go-vgo/robotgo/blob/master/docs/keys.md)
- [크로스 컴파일](https://github.com/go-vgo/robotgo/blob/master/docs/install.md#crosscompiling)
- [작성자](#authors)
- [계획](#plans)
- [라이선스](#license)

## Docs

- [GoDoc](https://godoc.org/github.com/go-vgo/robotgo) <br>
- [API 문서](https://github.com/go-vgo/robotgo/blob/master/docs/doc.md) (지원 중단, 더 이상 업데이트되지 않음)

## Binding:

[ADB](https://github.com/vcaesar/adb), Android adb API를 래핑한 패키지.

## Requirements:

이제 RobotGo를 설치하기 전에 `Golang, GCC`가 올바르게 설치되어 있는지 확인하세요.

### 전체 플랫폼:

```
Golang

GCC
```

#### MacOS:

```
brew install go
```

Xcode 명령줄 도구; <br>
그리고 개인정보 보호 설정에서 다음 위치에 "화면 기록"과 "손쉬운 사용" 권한을 추가하세요: <br>
`시스템 설정 > 개인정보 보호 및 보안 > 손쉬운 사용, 화면 및 시스템 오디오 기록`.

```
xcode-select --install
```

#### Windows:

```
winget install Golang.go
```

[llvm-mingw](https://github.com/mstorsjo/llvm-mingw)

```
winget install MartinStorsjo.LLVM-MinGW.UCRT
```

또는 [Mingw-w64](https://sourceforge.net/projects/mingw-w64/files)

```
winget install BrechtSanders.WinLibs.POSIX.UCRT
```

또는 [Mingw-w64](https://sourceforge.net/projects/mingw-w64/files)와 다른 gcc를 다운로드한 다음, `C:\mingw64\bin`과 같은 경로를 시스템 환경 변수 `Path`에 설정하세요.
[명령줄에서 GCC를 실행하도록 환경 변수 설정하기](https://www.youtube.com/results?search_query=Set+environment+variables+to+run+GCC+from+command+line).

`또는 다른 GCC 사용`(Mingw-w64를 제외하고, [bitmap](https://github.com/vcaesar/bitmap)을 사용할 때는 "libpng"를 직접 컴파일해야 합니다.)

#### 그 외 모든 플랫폼:

```
GCC

XTest 확장이 포함된 X11 (즉 Xtst 라이브러리)

"클립보드": xsel xclip

"비트맵": libpng (오직 "bitmap"에서만 사용.)

"이벤트-Gohook": xcb, xkb, libxkbcommon (오직 "hook"에서만 사용.)
```

##### Ubuntu:

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

##### Fedora:

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

Wayland 백엔드는 **순수 Go(Cgo 없음)** 구현이므로 시스템 C 라이브러리가 필요하지
않습니다. 다음 프로토콜을 지원하는 wlroots 기반 컴포지터(Sway, Hyprland, Wayfire 등)가
필요합니다:

```
zwlr_virtual_pointer_v1            (마우스 제어)
zwp_virtual_keyboard_v1            (키보드 제어)
zwlr_screencopy_v1                 (화면 캘처)
zwlr_foreign_toplevel_management_v1 (창 관리)
```

GNOME과 KDE는 이러한 프로토콜을 기본적으로 지원하지 **않습니다**.

#### libei(GNOME / KDE)

libei 백엔드도 **순수 Go(Cgo 없음)** 구현입니다. freedesktop의
`xdg-desktop-portal` RemoteDesktop 인터페이스를 통해 입력을 구동하므로, wlroots
Wayland 백엔드와 달리 GNOME과 KDE에서 작동합니다. 다음이 필요합니다:

```
xdg-desktop-portal               (portal D-Bus 서비스)
xdg-desktop-portal-gnome / -kde  (데스크톱의 portal 백엔드)
```

참고: libei 백엔드는 마우스와 키보드 입력만 처리합니다. 화면 캘처와 창 관리는
`ErrNotSupported`를 반환합니다.

## Installation:

Go 모듈을 지원하는 경우(Go 1.11+), import만 하면 됩니다:

```go
import "github.com/go-vgo/robotgo"
```

그렇지 않으면 다음 명령을 실행하여 robotgo 패키지를 설치하세요:

```
go get github.com/go-vgo/robotgo
```

png.h: No such file or directory? [issues/47](https://github.com/go-vgo/robotgo/issues/47)을 참조하세요.

## Update:

```
go get -u github.com/go-vgo/robotgo
```

go1.10.x의 C 파일 컴파일 캐시 문제에 주의하세요, [golang #24355](https://github.com/golang/go/issues/24355).
`go mod vendor` 문제, [golang #26366](https://github.com/golang/go/issues/26366).

## Cgo-free Builds:

RobotGo는 Windows, macOS, X11, Wayland, libei(Linux)용 **순수 Go(Cgo 없음)** 백엔드를 제공합니다. 현재는 실험적 기능입니다.
동일한 `robotgo` API를 노출하므로 코드를 변경할 필요 없이 빌드 태그만 지정하면 됩니다.
이 백엔드는 `CGO_ENABLED=0`으로 크로스 컴파일됩니다(GCC, MinGW, Xcode, X11 헤더 불필요).

| 백엔드                         | 빌드 태그 | Go 패키지                           |
| ------------------------------ | --------- | ----------------------------------- |
| Windows(Cgo 없음)              | `win`     | `github.com/go-vgo/robotgo/win`     |
| macOS(purego를 통한 Quartz)    | `mac`     | `github.com/go-vgo/robotgo/darwin`  |
| X11(Linux, 순수 Go X 프로토콜) | `x11`     | `github.com/go-vgo/robotgo/x11`     |
| Wayland(Linux, wlroots)        | `wayland` | `github.com/go-vgo/robotgo/wayland` |
| libei(Linux, GNOME/KDE portal) | `libei`   | `github.com/go-vgo/robotgo/libei`   |
| 순수 Go 기본값(모든 플랫폼)    | `purego`  | 위의 `mac`/`win`/`wayland` 선택     |

```sh
# 플랫폼별 순수 Go 기본 백엔드, 모든 대상에 하나의 태그 사용:
# macOS -> mac, Windows -> win, Linux -> wayland(x11/libei와 조합해 재정의 가능)
go build -tags purego ./...
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -tags "purego,x11" ./...

# Windows, Cgo / MinGW 불필요
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -tags win ./...

# macOS, purego를 통해 런타임에 Quartz/CoreGraphics 로드(Xcode 불필요)
CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -tags mac ./...

# X11, 순수 Go X 프로토콜(XTEST) — X11 헤더 불필요
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -tags x11 ./...

# Wayland, wlroots 기반 컴포지터(Sway, Hyprland, Wayfire 등)
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -tags wayland ./...

# libei, xdg-desktop-portal RemoteDesktop를 통한 GNOME/KDE 지원
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -tags libei ./...
```

`win` 태그에서는 기본 Cgo/Win32 백엔드가 제외되고 호출이 순수 Go `win` 패키지로 전달됩니다. `mac` 태그에서는 기본 Cgo/Quartz 백엔드가 제외되고 호출이 순수 Go `darwin` 패키지로 전달됩니다(창 관리는 `ErrNotSupported`를 반환합니다). `x11` 태그에서는 Cgo/X11 백엔드가 제외되고 호출이 순수 Go `x11` 패키지로 전달됩니다. `wayland` 태그에서는 Cgo/X11 백엔드가 제외되고 호출이 순수 Go `wayland` 패키지로 전달됩니다. `libei` 태그에서는 Cgo/X11과 wlroots Wayland 백엔드가 모두 제외되고 호출이 순수 Go `libei` 패키지로 전달됩니다.

`purego` 태그는 크로스 플랫폼 단축 태그입니다. 모든 플랫폼에서 Cgo 백엔드를 제외하고 대상 OS의 순수 Go 기본 백엔드(macOS는 `mac`, Windows는 `win`, Linux는 `wayland`)를 선택합니다. Linux에서는 `x11` 또는 `libei`와 조합하여(예: `-tags "purego,libei"`) 다른 순수 Go 백엔드를 선택할 수 있습니다.

## [Examples:](https://github.com/go-vgo/robotgo/blob/master/examples)

#### [마우스](https://github.com/go-vgo/robotgo/blob/master/examples/mouse/main.go)

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

#### [키보드](https://github.com/go-vgo/robotgo/blob/master/examples/key/main.go)

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

- 일부 C 코드를 Go로 리팩터링 (예: x11, windows)
- 더 나은 멀티 스크린 지원
- Wayland 지원
- 윈도우 핸들 업데이트
- Android 및 iOS 지원 시도

## Contributors

- 전체 기여자 목록은 [기여자 페이지](https://github.com/go-vgo/robotgo/graphs/contributors)를 참조하세요.
- [기여 가이드라인](https://github.com/go-vgo/robotgo/blob/master/CONTRIBUTING.md)을 참조하세요.

## License

Robotgo는 주로 "the Apache License (Version 2.0)" 조건에 따라 배포되며, 일부 내용은 다양한 BSD 계열 라이선스의 적용을 받습니다.

[LICENSE-APACHE](http://www.apache.org/licenses/LICENSE-2.0), [LICENSE](https://github.com/go-vgo/robotgo/blob/master/LICENSE)를 참조하세요.
