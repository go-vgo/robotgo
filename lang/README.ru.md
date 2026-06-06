# Robotgo

[![Build Status](https://github.com/go-vgo/robotgo/workflows/Go/badge.svg)](https://github.com/go-vgo/robotgo/commits/master)
[![CircleCI Status](https://circleci.com/gh/go-vgo/robotgo.svg?style=shield)](https://circleci.com/gh/go-vgo/robotgo)
[![Go Report Card](https://goreportcard.com/badge/github.com/go-vgo/robotgo)](https://goreportcard.com/report/github.com/go-vgo/robotgo)
[![GoDoc](https://pkg.go.dev/badge/github.com/go-vgo/robotgo?status.svg)](https://pkg.go.dev/github.com/go-vgo/robotgo?tab=doc)
[![GitHub release](https://img.shields.io/github/release/go-vgo/robotgo.svg)](https://github.com/go-vgo/robotgo/releases/latest)
<a href="https://discord.gg/npPb3NzE4A"><img src="https://img.shields.io/discord/1484658282777018551.svg?logo=discord&logoColor=white&label=Discord&color=5865F2" alt="Join the Discord chat at https://discord.gg/npPb3NzE4A"></a>

[English](../README.md) | [简体中文](README.zh.md) | [繁體中文](README.zht.md) | [日本語](README.ja.md) | [한국어](README.ko.md) | [Français](README.fr.md) | [Deutsch](README.de.md) | [Español](README.es.md) | Русский | [Português](README.pt.md)

> Автоматизация рабочего стола на Golang, автотестирование и управление компьютером с помощью ИИ (Computer Use). <br>
> Управление мышью и клавиатурой, чтение экрана, процессы, дескрипторы окон, изображения и битовые карты, а также глобальный перехват событий.

RobotGo поддерживает Mac, Windows и Linux; а также поддерживает архитектуры arm64 и x86-amd64.

Сейчас я создаю [Codg](https://github.com/vcaesar/codg) — простую и удобную рабочую систему ИИ-агентов (AI Agent): автоматизация, асинхронность, параллелизм, эффективность и высокая точность.

<p align="center">
<a href="https://github.com/vcaesar/codg" rel="nofollow">
<img width="800" alt="Codg Demo" src="https://github.com/vcaesar/codg/raw/main/demo/26-04-1.png" />
</a>
</p>

[RobotGo-Pro](https://github.com/vcaesar/robotgo-pro) предоставляет версии на JavaScript, Python, Lua и других языках, техническую поддержку, новые возможности, а также новейшую версию robotgo («сейчас нет версии с открытым исходным кодом»).

## Содержание

- [Документация](#docs)
- [Привязки](#binding)
- [Требования](#requirements)
- [Установка](#installation)
- [Обновление](#update)
- [Сборки без Cgo](#cgo-free-builds)
- [Примеры](#examples)
- [Преобразование типов и клавиши](https://github.com/go-vgo/robotgo/blob/master/docs/keys.md)
- [Кросс-компиляция](https://github.com/go-vgo/robotgo/blob/master/docs/install.md#crosscompiling)
- [Авторы](#authors)
- [Планы](#plans)
- [Лицензия](#license)

## Docs

- [GoDoc](https://godoc.org/github.com/go-vgo/robotgo) <br>
- [Документация API](https://github.com/go-vgo/robotgo/blob/master/docs/doc.md) (устарела, больше не обновляется)

## Binding:

[ADB](https://github.com/vcaesar/adb) — обёртка над Android adb API.

## Requirements:

Прежде чем устанавливать RobotGo, пожалуйста, убедитесь, что `Golang, GCC` установлены корректно.

### Все платформы:

```
Golang

GCC
```

#### Для MacOS:

```
brew install go
```

Инструменты командной строки Xcode; <br>
А также в настройках конфиденциальности добавьте разрешения «Запись экрана» и «Универсальный доступ» в разделе: <br>
`Системные настройки > Конфиденциальность и безопасность > Универсальный доступ, Запись экрана и системного звука`.

```
xcode-select --install
```

#### Для Windows:

```
winget install Golang.go
```

[llvm-mingw](https://github.com/mstorsjo/llvm-mingw)

```
winget install MartinStorsjo.LLVM-MinGW.UCRT
```

или [Mingw-w64](https://sourceforge.net/projects/mingw-w64/files)

```
winget install BrechtSanders.WinLibs.POSIX.UCRT
```

Или скачайте [Mingw-w64](https://sourceforge.net/projects/mingw-w64/files) и другие gcc, затем добавьте путь вида `C:\mingw64\bin` в системную переменную окружения `Path`.
[Настройка переменных окружения для запуска GCC из командной строки](https://www.youtube.com/results?search_query=Set+environment+variables+to+run+GCC+from+command+line).

`Или другие GCC` (кроме Mingw-w64; при использовании [bitmap](https://github.com/vcaesar/bitmap) библиотеку "libpng" вам придётся скомпилировать самостоятельно.)

#### Для всех остальных платформ:

```
GCC

X11 с расширением XTest (библиотека Xtst)

"Clipboard": xsel xclip

"Bitmap": libpng (Используется только в "bitmap".)

"Event-Gohook": xcb, xkb, libxkbcommon (Используется только в "hook".)
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

Бэкенд Wayland — это **чистая реализация на Go (без Cgo)**, поэтому системные
C-библиотеки не требуются. Он требует композитора на основе wlroots (Sway,
Hyprland, Wayfire и др.), поддерживающего следующие протоколы:

```
zwlr_virtual_pointer_v1            (управление мышью)
zwp_virtual_keyboard_v1            (управление клавиатурой)
zwlr_screencopy_v1                 (захват экрана)
zwlr_foreign_toplevel_management_v1 (управление окнами)
```

GNOME и KDE **не** поддерживают эти протоколы нативно.

#### libei (GNOME / KDE)

Бэкенд libei также является **чистой реализацией на Go (без Cgo)**. Он
управляет вводом через интерфейс RemoteDesktop из `xdg-desktop-portal`
freedesktop, поэтому работает на GNOME и KDE (в отличие от wlroots Wayland
бэкенда). Он требует:

```
xdg-desktop-portal               (D-Bus-служба portal)
xdg-desktop-portal-gnome / -kde  (portal-бэкенд вашего рабочего стола)
```

Примечание: бэкенд libei обрабатывает только ввод с мыши и клавиатуры. Захват
экрана и управление окнами возвращают `ErrNotSupported`.

## Installation:

При поддержке Go module (Go 1.11+) достаточно выполнить import:

```go
import "github.com/go-vgo/robotgo"
```

В противном случае установите пакет robotgo, выполнив команду:

```
go get github.com/go-vgo/robotgo
```

png.h: No such file or directory? См. [issues/47](https://github.com/go-vgo/robotgo/issues/47).

## Update:

```
go get -u github.com/go-vgo/robotgo
```

Обратите внимание на проблему кэширования компиляции C-файлов в go1.10.x, [golang #24355](https://github.com/golang/go/issues/24355).
Проблема `go mod vendor`, [golang #26366](https://github.com/golang/go/issues/26366).

## Cgo-free Builds:

RobotGo предоставляет **чистые Go-бэкенды (без Cgo)** для Windows, Wayland и
libei (Linux). Они предоставляют тот же API `robotgo`, поэтому ваш код не нуждается в
изменениях — достаточно тега сборки. Эти бэкенды кросс-компилируются с
`CGO_ENABLED=0` (без GCC, MinGW или заголовков X11).

| Бэкенд                          | Тег сборки | Go-пакет                            |
| ------------------------------- | ---------- | ----------------------------------- |
| Windows (без Cgo)               | `win`      | `github.com/go-vgo/robotgo/win`     |
| Wayland (Linux, wlroots)        | `wayland`  | `github.com/go-vgo/robotgo/wayland` |
| libei (Linux, portal GNOME/KDE) | `libei`    | `github.com/go-vgo/robotgo/libei`   |

```sh
# Windows, без Cgo / без MinGW
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -tags win ./...

# Wayland, композитор на основе wlroots (Sway, Hyprland, Wayfire и др.)
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -tags wayland ./...

# libei, GNOME/KDE через xdg-desktop-portal RemoteDesktop
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -tags libei ./...
```

С тегом `win` бэкенд Cgo/Win32 по умолчанию исключается, а вызовы
перенаправляются в чистый Go-пакет `win`; с тегом `wayland` исключается бэкенд
Cgo/X11, а вызовы перенаправляются в чистый Go-пакет `wayland`; с тегом `libei`
исключаются как бэкенд Cgo/X11, так и wlroots Wayland бэкенд, а вызовы
перенаправляются в чистый Go-пакет `libei`.

## [Examples:](https://github.com/go-vgo/robotgo/blob/master/examples)

#### [Мышь](https://github.com/go-vgo/robotgo/blob/master/examples/mouse/main.go)

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

#### [Клавиатура](https://github.com/go-vgo/robotgo/blob/master/examples/key/main.go)

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

#### [Экран](https://github.com/go-vgo/robotgo/blob/master/examples/screen/main.go)

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

#### [Битовая карта](https://github.com/vcaesar/bitmap/blob/main/examples/main.go)

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

#### [Событие](https://github.com/robotn/gohook/blob/master/examples/main.go)

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

#### [Окно](https://github.com/go-vgo/robotgo/blob/master/examples/window/main.go)

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

- [Автор — Evans](https://github.com/vcaesar)
- [Сопровождающие](https://github.com/orgs/go-vgo/people)

## Plans

- Переписать часть C-кода на Go (например, x11, windows)
- Улучшить поддержку нескольких экранов
- Поддержка Wayland
- Обновить работу с дескрипторами окон
- Попробовать добавить поддержку Android и iOS

## Contributors

- Полный список участников см. на [странице участников](https://github.com/go-vgo/robotgo/graphs/contributors).
- См. [Руководство для участников](https://github.com/go-vgo/robotgo/blob/master/CONTRIBUTING.md).

## License

Robotgo распространяется преимущественно на условиях «Apache License (Version 2.0)», при этом отдельные части подпадают под различные лицензии в стиле BSD.

См. [LICENSE-APACHE](http://www.apache.org/licenses/LICENSE-2.0), [LICENSE](https://github.com/go-vgo/robotgo/blob/master/LICENSE).
