# Robotgo

<p align="center">
  <img src="../docs/robotgo-logo.svg" width="560" alt="RobotGo-Logo — ein Roboter mit Mauszeiger" />
</p>

[![Build Status](https://github.com/go-vgo/robotgo/workflows/Go/badge.svg)](https://github.com/go-vgo/robotgo/commits/master)
[![CircleCI Status](https://circleci.com/gh/go-vgo/robotgo.svg?style=shield)](https://circleci.com/gh/go-vgo/robotgo)
[![golangci-lint](https://github.com/go-vgo/robotgo/actions/workflows/lint.yml/badge.svg)](https://github.com/go-vgo/robotgo/actions/workflows/lint.yml)
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
- [Cgo-freie Builds](#cgo-free-builds)
- [Installation](#installation)
- [Aktualisierung](#update)
- [Beispiele](#examples)
- [Typkonvertierung und Tasten](https://github.com/go-vgo/robotgo/blob/master/docs/keys.md)
- [Cross-Compiling](https://github.com/go-vgo/robotgo/blob/master/docs/install.md#crosscompiling)
- [Autoren](#authors)
- [Pläne](#plans)
- [Lizenz](#license)

## Docs

- [GoDoc](https://godoc.org/github.com/go-vgo/robotgo) <br>
- [API-Dokumentation](https://github.com/go-vgo/robotgo/blob/master/docs/doc.md) (Veraltet, nicht aktualisiert)

## Binding

[ADB](https://github.com/vcaesar/adb), kapselt die Android-adb-API.

## Requirements

**Go 1.26+** (siehe [go.mod](../go.mod)) und ein Backend:

- **Standard (Cgo):** `CGO_ENABLED=1`, ein C-Compiler und die unten genannten
  Plattformbibliotheken.
- **Reines Go (experimentell):** keine C-Toolchain; siehe [Cgo-freie Builds](#cgo-free-builds).

Beide benötigen die [funktionsspezifischen Abhängigkeiten](#feature-specific-dependencies)
für die jeweils genutzten Funktionen.

### Default Cgo setup

#### macOS

Go und die Xcode Command Line Tools (Clang):

```sh
brew install go
xcode-select --install
```

**Berechtigungen (Cgo und reines Go):** unter **Systemeinstellungen > Datenschutz
& Sicherheit** der App oder dem Terminal **Bedienungshilfen** (Eingabe) und
**Bildschirmaufnahme** (Erfassung) erteilen.

#### Windows

```powershell
winget install GoLang.Go
```

Dazu **eine** C-Toolchain. [LLVM-MinGW](https://github.com/mstorsjo/llvm-mingw):

```powershell
winget install MartinStorsjo.LLVM-MinGW.UCRT
$env:CC = "clang"
```

Oder [MinGW-w64 (WinLibs)](https://winlibs.com/):

```powershell
winget install BrechtSanders.WinLibs.POSIX.UCRT
$env:CC = "gcc"
```

Das `bin`-Verzeichnis der Toolchain (z. B. `C:\mingw64\bin`) muss in `PATH`
liegen und zu `GOARCH` passen. Jeder andere Cgo-kompatible Compiler, etwa ein
manuell installiertes [MinGW-w64](https://sourceforge.net/projects/mingw-w64/files),
funktioniert ebenfalls.

[Bitmap](https://github.com/vcaesar/bitmap) liefert libpng nur für MinGW-w64 mit;
mit einer anderen Toolchain muss libpng selbst gebaut werden.

#### Linux (X11)

Build: GCC, libc-Header, X11/XTest-Header. Laufzeit: ein X-Server mit XTEST und
gesetztem `DISPLAY`. Für natives Wayland siehe [Cgo-freie Builds](#cgo-free-builds).

**Ubuntu / Debian:**

```sh
# Go (Snap, falls das Distributionspaket zu alt ist)
sudo snap install go --classic
# oder: sudo apt install golang

# Kern (Cgo): Compiler, libc, X11, XTest
sudo apt install gcc libc6-dev libx11-dev xorg-dev libxtst-dev

# Zwischenablage (X11)
sudo apt install xsel xclip

# Bitmap: libpng
sudo apt install libpng++-dev

# GoHook: XCB, XKB
sudo apt install xcb libxcb-xkb-dev x11-xkb-utils libx11-xcb-dev libxkbcommon-x11-dev libxkbcommon-dev
```

**Fedora:**

```sh
# Kern (Cgo): Compiler, libc, X11, XTest
sudo dnf install gcc glibc-devel libX11-devel libXtst-devel

# Zwischenablage (X11)
sudo dnf install xsel xclip

# Bitmap: libpng
sudo dnf install libpng-devel

# GoHook: XKB (Fedora < 34: xorg-x11-xkb-utils-devel statt xkbcomp-devel)
sudo dnf install libxkbcommon-devel libxkbcommon-x11-devel xkbcomp-devel
```

### Feature-specific dependencies

Zusätzlich zur Grundinstallation, nur für die genutzten Funktionen nötig:

- **Linux-Zwischenablage (Cgo und reines Go):** `xclip` oder `xsel` unter X11,
  `wl-clipboard` unter Wayland:

  ```sh
  sudo apt install wl-clipboard   # Ubuntu / Debian
  sudo dnf install wl-clipboard   # Fedora
  ```

- **[Bitmap](https://github.com/vcaesar/bitmap):** libpng-Header.
- **[GoHook](https://github.com/robotn/gohook):** XCB/XKB-Header unter Linux.

Bitmap und GoHook sind eigene Cgo-Pakete; ein reines Go-Backend von RobotGo hebt
deren C-Anforderungen nicht auf.

## Cgo-free Builds

**Experimentelle reine Go-Backends** bauen und cross-kompilieren mit
`CGO_ENABLED=0` (ohne GCC, MinGW, Xcode oder X11-Header). Gleicher Importpfad und
gleiche API; ein **Build-Tag** wählt das Backend — `CGO_ENABLED=0` allein wählt
nichts aus.

| Ziel-`GOOS`   | Build-Tag | Paket                  | Laufzeitanforderungen                                                      |
| ------------- | --------- | ---------------------- | -------------------------------------------------------------------------- |
| `windows`     | `win`     | [win](../win/)         | nur Win32, keine zusätzliche Laufzeit                                      |
| `darwin`      | `mac`     | [darwin](../darwin/)   | System-Frameworks und [Datenschutzberechtigungen](#macos)                  |
| `linux`       | `x11`     | [x11](../x11/)         | X-Server mit XTEST, `DISPLAY` gesetzt                                      |
| `linux`       | `wayland` | [wayland](../wayland/) | wlroots-Compositor, siehe [Wayland](#wayland)                              |
| `linux`       | `libei`   | [libei](../libei/)     | RemoteDesktop-Portal (GNOME/KDE), siehe [libei](#libei-gnome--kde)         |

**`purego`** ist eine Abkürzung: `mac` unter macOS, `win` unter Windows,
`wayland` unter Linux. Mit zusätzlichem `x11` oder `libei` lässt sich das unter
Linux überschreiben (`-tags "purego,x11"`).

**Ein Backend-Tag** passend zu `GOOS` verwenden; Kombinationen wie `x11,libei`
lassen sich nicht kompilieren (`purego` plus ein Linux-Override ist die einzige
Ausnahme).

### Build commands

```sh
# Aktuelles Betriebssystem
CGO_ENABLED=0 go build -tags purego .

# Cross-Kompilierung
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -tags win .
CGO_ENABLED=0 GOOS=darwin  GOARCH=arm64 go build -tags mac .
CGO_ENABLED=0 GOOS=linux   GOARCH=amd64 go build -tags x11 .
CGO_ENABLED=0 GOOS=linux   GOARCH=amd64 go build -tags wayland .
CGO_ENABLED=0 GOOS=linux   GOARCH=amd64 go build -tags libei .
```

PowerShell: `$env:CGO_ENABLED = "0"; go build -tags purego .`

Das Modul-Stammverzeichnis (`.`) oder das eigene Paket bauen, nicht `./...`:
`examples/` und einige Unterpakete brauchen Cgo-only-APIs.

### Limitations

Die Laufzeitanforderungen von oben gelten weiterhin. Cgo-only-APIs wie
`CaptureScreen` / `FreeBitmap` gibt es nicht; stattdessen `CaptureImg` nutzen.
Nicht unterstützte Operationen geben `ErrNotSupported` zurück (oder ein leeres
Ergebnis).

- **`mac`:** keine Fensterverwaltung — `ActiveName` gibt `ErrNotSupported`
  zurück, `GetTitle` gibt `""` zurück, Minimieren/Maximieren/Schließen sind
  No-Ops.
- **`wayland`, `libei`:** `Location()` gibt die letzte von RobotGo injizierte
  Position zurück, nicht den physischen Cursor. Die Zwischenablage braucht
  `wl-clipboard`.

#### Wayland

Ein wlroots-Compositor (Sway, Hyprland, Wayfire, ...) mit gesetztem
`WAYLAND_DISPLAY` und `XDG_RUNTIME_DIR`, der Folgendes bereitstellt:

| Funktion            | Protokoll-Global                   |
| ------------------- | ---------------------------------- |
| Maussteuerung       | `zwlr_virtual_pointer_manager_v1`  |
| Tastatursteuerung   | `zwp_virtual_keyboard_manager_v1`  |
| Bildschirmaufnahme  | `zwlr_screencopy_manager_v1`       |
| Fensterverwaltung   | `zwlr_foreign_toplevel_manager_v1` |

GNOME und KDE stellen diese Protokolle nicht bereit; dort `libei` verwenden.

#### libei (GNOME / KDE)

Eingaben laufen über die **RemoteDesktop**-D-Bus-Schnittstelle von
`xdg-desktop-portal`: nötig sind ein Session-Bus, `xdg-desktop-portal` und ein
Backend, das sie implementiert (`xdg-desktop-portal-gnome`, `-kde` oder `-wlr`).
Beim ersten Start erscheint ein Zustimmungsdialog; das Restore-Token wird unter
`$XDG_STATE_HOME/robotgo` zwischengespeichert.

Absolute Bewegungen nutzen einen verknüpften ScreenCast-Stream (Standard); mit
`libei.LinkScreenCast = false` gibt es nur relative Bewegungen.
Bildschirmaufnahme und Fensterverwaltung melden `ErrNotSupported`.

## Installation

```sh
go get github.com/go-vgo/robotgo
```

```go
import "github.com/go-vgo/robotgo"
```

`png.h: No such file or directory` mit Bitmap: siehe dessen libpng-Anforderungen
und [issues/47](https://github.com/go-vgo/robotgo/issues/47).

## Update

```sh
go get -u github.com/go-vgo/robotgo
```

Beachten Sie das Problem mit dem Kompilierungs-Cache für C-Dateien in go1.10.x, [golang #24355](https://github.com/golang/go/issues/24355).
`go mod vendor`-Problem, [golang #26366](https://github.com/golang/go/issues/26366).

## [Examples](https://github.com/go-vgo/robotgo/blob/master/examples)

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

- Bessere Multiscreen-Unterstützung
- Fensterhandle aktualisieren
- Versuch, Android und iOS zu unterstützen

## Contributors

- Die vollständige Liste der Mitwirkenden finden Sie auf der [Mitwirkenden-Seite](https://github.com/go-vgo/robotgo/graphs/contributors).
- Siehe [Beitragsrichtlinien](https://github.com/go-vgo/robotgo/blob/master/CONTRIBUTING.md).

## License

Robotgo wird primär unter den Bedingungen „der Apache-Lizenz (Version 2.0)“ vertrieben, wobei Teile von verschiedenen BSD-ähnlichen Lizenzen abgedeckt sind.

Siehe [LICENSE-APACHE](http://www.apache.org/licenses/LICENSE-2.0), [LICENSE](https://github.com/go-vgo/robotgo/blob/master/LICENSE).
