# Robotgo

[![Build Status](https://github.com/go-vgo/robotgo/workflows/Go/badge.svg)](https://github.com/go-vgo/robotgo/commits/master)
[![CircleCI Status](https://circleci.com/gh/go-vgo/robotgo.svg?style=shield)](https://circleci.com/gh/go-vgo/robotgo)
[![Go Report Card](https://goreportcard.com/badge/github.com/go-vgo/robotgo)](https://goreportcard.com/report/github.com/go-vgo/robotgo)
[![GoDoc](https://pkg.go.dev/badge/github.com/go-vgo/robotgo?status.svg)](https://pkg.go.dev/github.com/go-vgo/robotgo?tab=doc)
[![GitHub release](https://img.shields.io/github/release/go-vgo/robotgo.svg)](https://github.com/go-vgo/robotgo/releases/latest)
<a href="https://discord.gg/npPb3NzE4A"><img src="https://img.shields.io/discord/1484658282777018551.svg?logo=discord&logoColor=white&label=Discord&color=5865F2" alt="Join the Discord chat at https://discord.gg/npPb3NzE4A"></a>

[English](../README.md) | [简体中文](README.zh.md) | [繁體中文](README.zht.md) | [日本語](README.ja.md) | [한국어](README.ko.md) | [Français](README.fr.md) | Deutsch | [Español](README.es.md) | [Русский](README.ru.md) | [Português](README.pt.md)

> Golang Desktop-Automatisierung, automatisiertes Testen und KI-gestützte Computer-Bedienung (Computer Use). <br>
> Steuerung von Maus und Tastatur, Auslesen des Bildschirms, Prozesse, Fensterhandles, Bilder und Bitmaps sowie globales Event-Listening.

RobotGo unterstützt Mac, Windows und Linux; außerdem unterstützt RobotGo arm64 und x86-amd64.

Ich entwickle jetzt [Codg](https://github.com/vcaesar/codg), ein einfach zu bedienendes KI-Agentensystem zum Programmieren und Arbeiten: automatisch, asynchron, nebenläufig, effizient und mit hoher Genauigkeit.

<p align="center">
<a href="https://github.com/vcaesar/codg" rel="nofollow">
<img width="800" alt="Codg Demo" src="https://github.com/vcaesar/codg/raw/main/demo/26-04-1.png" />
</a>
</p>

[RobotGo-Pro](https://github.com/vcaesar/robotgo-pro) bietet die JavaScript-, Python-, Lua- und weitere Versionen, technischen Support, neue Funktionen und die neueste robotgo-Version („derzeit keine Open-Source-Version“).

## Inhalt

- [Dokumentation](#docs)
- [Bindings](#binding)
- [Voraussetzungen](#requirements)
- [Installation](#installation)
- [Aktualisierung](#update)
- [Cgo-freie Builds](#cgo-free-builds)
- [Beispiele](#examples)
- [Typkonvertierung und Tasten](https://github.com/go-vgo/robotgo/blob/master/docs/keys.md)
- [Cross-Compiling](https://github.com/go-vgo/robotgo/blob/master/docs/install.md#crosscompiling)
- [Autoren](#authors)
- [Pläne](#plans)
- [Lizenz](#license)

## Docs

- [GoDoc](https://godoc.org/github.com/go-vgo/robotgo) <br>
- [API-Dokumentation](https://github.com/go-vgo/robotgo/blob/master/docs/doc.md) (Veraltet, nicht aktualisiert)

## Binding:

[ADB](https://github.com/vcaesar/adb), kapselt die Android-adb-API.

## Requirements:

Bitte stellen Sie nun sicher, dass `Golang, GCC` korrekt installiert sind, bevor Sie RobotGo installieren.

### ALLE:

```
Golang

GCC
```

#### Für MacOS:

```
brew install go
```

Xcode Command Line Tools; <br>
Und in den Datenschutzeinstellungen „Bildschirmaufnahme“ und „Bedienungshilfen“ hinzufügen unter: <br>
`Systemeinstellungen > Datenschutz & Sicherheit > Bedienungshilfen, Bildschirm- & Systemaudioaufnahme`.

```
xcode-select --install
```

#### Für Windows:

```
winget install Golang.go
```

[llvm-mingw](https://github.com/mstorsjo/llvm-mingw)

```
winget install MartinStorsjo.LLVM-MinGW.UCRT
```

oder [Mingw-w64](https://sourceforge.net/projects/mingw-w64/files)

```
winget install BrechtSanders.WinLibs.POSIX.UCRT
```

Oder laden Sie [Mingw-w64](https://sourceforge.net/projects/mingw-w64/files) und die anderen gcc herunter und setzen Sie anschließend Systemumgebungsvariablen wie `C:\mingw64\bin` in die Umgebungsvariable `Path`.
[Umgebungsvariablen setzen, um GCC über die Befehlszeile auszuführen](https://www.youtube.com/results?search_query=Set+environment+variables+to+run+GCC+from+command+line).

`Oder die anderen GCC` (Außer Mingw-w64 müssen Sie „libpng“ selbst kompilieren, wenn Sie die [bitmap](https://github.com/vcaesar/bitmap) verwenden.)

#### Für alles andere:

```
GCC

X11 mit der XTest-Erweiterung (die Xtst-Bibliothek)

"Clipboard": xsel xclip

"Bitmap": libpng (Wird nur von "bitmap" verwendet.)

"Event-Gohook": xcb, xkb, libxkbcommon (Wird nur von "hook" verwendet.)
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

Das Wayland-Backend ist eine **reine Go-Implementierung (ohne Cgo)**, daher sind
keine System-C-Bibliotheken erforderlich. Es benötigt einen wlroots-basierten
Compositor (Sway, Hyprland, Wayfire, ...), der die folgenden Protokolle
unterstützt:

```
zwlr_virtual_pointer_v1            (Maussteuerung)
zwp_virtual_keyboard_v1            (Tastatursteuerung)
zwlr_screencopy_v1                 (Bildschirmaufnahme)
zwlr_foreign_toplevel_management_v1 (Fensterverwaltung)
```

GNOME und KDE unterstützen diese Protokolle **nicht** nativ.

#### libei (GNOME / KDE)

Das libei-Backend ist ebenfalls eine **reine Go-Implementierung (ohne Cgo)**. Es
steuert Eingaben über die RemoteDesktop-Schnittstelle von `xdg-desktop-portal`
von freedesktop und funktioniert daher auf GNOME und KDE (im Gegensatz zum
wlroots-Wayland-Backend). Es benötigt:

```
xdg-desktop-portal               (der Portal-D-Bus-Dienst)
xdg-desktop-portal-gnome / -kde  (das Portal-Backend deines Desktops)
```

Hinweis: Das libei-Backend verarbeitet nur Maus- und Tastatureingaben.
Bildschirmaufnahme und Fensterverwaltung geben `ErrNotSupported` zurück.

## Installation:

Mit Go-Modul-Unterstützung (Go 1.11+) einfach importieren:

```go
import "github.com/go-vgo/robotgo"
```

Andernfalls führen Sie zur Installation des robotgo-Pakets den folgenden Befehl aus:

```
go get github.com/go-vgo/robotgo
```

png.h: No such file or directory? Bitte siehe [issues/47](https://github.com/go-vgo/robotgo/issues/47).

## Update:

```
go get -u github.com/go-vgo/robotgo
```

Beachten Sie das Problem mit dem Kompilierungs-Cache für C-Dateien in go1.10.x, [golang #24355](https://github.com/golang/go/issues/24355).
`go mod vendor`-Problem, [golang #26366](https://github.com/golang/go/issues/26366).

## Cgo-free Builds:

RobotGo bietet **reine Go-Backends (ohne Cgo)** für Windows, macOS, X11,
Wayland und libei (Linux).
Sie stellen dieselbe `robotgo`-API bereit, sodass dein Code nicht geändert werden
muss — nur ein Build-Tag ist nötig. Diese Backends lassen sich mit
`CGO_ENABLED=0` cross-kompilieren (ohne GCC, MinGW, Xcode oder X11-Header).

| Backend                          | Build-Tag | Go-Paket                            |
| -------------------------------- | --------- | ----------------------------------- |
| Windows (ohne Cgo)               | `win`     | `github.com/go-vgo/robotgo/win`     |
| macOS (Quartz über purego)       | `mac`     | `github.com/go-vgo/robotgo/darwin`  |
| X11 (Linux, reines Go-X-Protokoll) | `x11`   | `github.com/go-vgo/robotgo/x11`     |
| Wayland (Linux, wlroots)         | `wayland` | `github.com/go-vgo/robotgo/wayland` |
| libei (Linux, GNOME/KDE-Portal)  | `libei`   | `github.com/go-vgo/robotgo/libei`   |

```sh
# Windows, ohne Cgo / ohne MinGW
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -tags win ./...

# macOS, Quartz/CoreGraphics zur Laufzeit über purego geladen (ohne Xcode)
CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -tags mac ./...

# X11, reines Go-X-Protokoll (XTEST) — ohne X11-Header
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -tags x11 ./...

# Wayland, wlroots-basierter Compositor (Sway, Hyprland, Wayfire, ...)
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -tags wayland ./...

# libei, GNOME/KDE über xdg-desktop-portal RemoteDesktop
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -tags libei ./...
```

Mit dem Tag `win` wird das standardmäßige Cgo/Win32-Backend ausgeschlossen und
Aufrufe werden an das reine Go-Paket `win` weitergeleitet; mit dem Tag `mac`
wird das standardmäßige Cgo/Quartz-Backend ausgeschlossen und Aufrufe werden an
das reine Go-Paket `darwin` weitergeleitet (Fensterverwaltung meldet
`ErrNotSupported`); mit dem Tag `x11` wird das Cgo/X11-Backend ausgeschlossen
und Aufrufe werden an das reine Go-Paket `x11` weitergeleitet; mit dem Tag
`wayland` wird das Cgo/X11-Backend ausgeschlossen und Aufrufe werden an das
reine Go-Paket `wayland` weitergeleitet; mit dem Tag `libei` werden sowohl das
Cgo/X11- als auch das wlroots-Wayland-Backend ausgeschlossen und Aufrufe werden
an das reine Go-Paket `libei` weitergeleitet.

## [Examples:](https://github.com/go-vgo/robotgo/blob/master/examples)

#### [Maus](https://github.com/go-vgo/robotgo/blob/master/examples/mouse/main.go)

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

#### [Tastatur](https://github.com/go-vgo/robotgo/blob/master/examples/key/main.go)

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

#### [Bildschirm](https://github.com/go-vgo/robotgo/blob/master/examples/screen/main.go)

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

#### [Fenster](https://github.com/go-vgo/robotgo/blob/master/examples/window/main.go)

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

- [Der Autor ist Evans](https://github.com/vcaesar)
- [Maintainer](https://github.com/orgs/go-vgo/people)

## Plans

- Teil des C-Codes nach Go umbauen (etwa x11, windows)
- Bessere Multiscreen-Unterstützung
- Wayland-Unterstützung
- Fensterhandle aktualisieren
- Versuch, Android und iOS zu unterstützen

## Contributors

- Die vollständige Liste der Mitwirkenden finden Sie auf der [Mitwirkenden-Seite](https://github.com/go-vgo/robotgo/graphs/contributors).
- Siehe [Beitragsrichtlinien](https://github.com/go-vgo/robotgo/blob/master/CONTRIBUTING.md).

## License

Robotgo wird primär unter den Bedingungen „der Apache-Lizenz (Version 2.0)“ vertrieben, wobei Teile von verschiedenen BSD-ähnlichen Lizenzen abgedeckt sind.

Siehe [LICENSE-APACHE](http://www.apache.org/licenses/LICENSE-2.0), [LICENSE](https://github.com/go-vgo/robotgo/blob/master/LICENSE).
