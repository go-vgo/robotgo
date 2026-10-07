# Robotgo

<p align="center">
  <img src="../docs/robotgo-logo.svg" width="560" alt="Логотип RobotGo — робот с указателем мыши" />
</p>

[![Build Status](https://github.com/go-vgo/robotgo/workflows/Go/badge.svg)](https://github.com/go-vgo/robotgo/commits/master)
[![CircleCI Status](https://circleci.com/gh/go-vgo/robotgo.svg?style=shield)](https://circleci.com/gh/go-vgo/robotgo)
[![codecov](https://codecov.io/gh/go-vgo/robotgo/branch/master/graph/badge.svg)](https://codecov.io/gh/go-vgo/robotgo)
[![golangci-lint](https://github.com/go-vgo/robotgo/actions/workflows/lint.yml/badge.svg)](https://github.com/go-vgo/robotgo/actions/workflows/lint.yml)
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
- [Сборки без Cgo](#cgo-free-builds)
- [Установка](#installation)
- [Обновление](#update)
- [Примеры](#examples)
- [Преобразование типов и клавиши](https://github.com/go-vgo/robotgo/blob/master/docs/keys.md)
- [Кросс-компиляция](https://github.com/go-vgo/robotgo/blob/master/docs/install.md#crosscompiling)
- [Авторы](#authors)
- [Планы](#plans)
- [Лицензия](#license)

## Docs

- [GoDoc](https://godoc.org/github.com/go-vgo/robotgo) <br>
- [Документация API](https://github.com/go-vgo/robotgo/blob/master/docs/doc.md) (устарела, больше не обновляется)

## Binding

[ADB](https://github.com/vcaesar/adb) — обёртка над Android adb API.

## Requirements

**Go 1.26+** (см. [go.mod](../go.mod)) и один из бэкендов:

- **По умолчанию (Cgo):** `CGO_ENABLED=1`, компилятор C и системные библиотеки,
  перечисленные ниже.
- **Чистый Go (экспериментально):** без набора инструментов C; см.
  [Сборки без Cgo](#cgo-free-builds).

Обоим нужны [зависимости возможностей](#feature-specific-dependencies) — только
для тех возможностей, которые вы используете.

### Default Cgo setup

#### macOS

Go и Xcode Command Line Tools (Clang):

```sh
brew install go
xcode-select --install
```

**Разрешения (Cgo и чистый Go):** в **Системных настройках > Конфиденциальность
и безопасность** выдайте приложению или терминалу **Универсальный доступ**
(ввод) и **Запись экрана** (захват).

#### Windows

```powershell
winget install GoLang.Go
```

Плюс **один** набор инструментов C. [LLVM-MinGW](https://github.com/mstorsjo/llvm-mingw):

```powershell
winget install MartinStorsjo.LLVM-MinGW.UCRT
$env:CC = "clang"
```

Или [MinGW-w64 (WinLibs)](https://winlibs.com/):

```powershell
winget install BrechtSanders.WinLibs.POSIX.UCRT
$env:CC = "gcc"
```

Каталог `bin` набора инструментов (например, `C:\mingw64\bin`) должен быть в
`PATH` и соответствовать `GOARCH`. Подойдёт и любой другой совместимый с Cgo
компилятор, например установленный вручную
[MinGW-w64](https://sourceforge.net/projects/mingw-w64/files).

[Bitmap](https://github.com/vcaesar/bitmap) поставляет libpng только для
MinGW-w64; с другим набором инструментов libpng придётся собрать самостоятельно.

#### Linux (X11)

Сборка: GCC, заголовки libc, заголовки X11/XTest. Выполнение: X-сервер с XTEST
и заданной переменной `DISPLAY`. Для нативного Wayland см.
[Сборки без Cgo](#cgo-free-builds).

**Ubuntu / Debian:**

```sh
# Go (Snap, если пакет дистрибутива слишком старый)
sudo snap install go --classic
# или: sudo apt install golang

# Основное (Cgo): компилятор, libc, X11, XTest
sudo apt install gcc libc6-dev libx11-dev xorg-dev libxtst-dev

# Буфер обмена (X11)
sudo apt install xsel xclip

# Bitmap: libpng
sudo apt install libpng++-dev

# GoHook: XCB, XKB
sudo apt install xcb libxcb-xkb-dev x11-xkb-utils libx11-xcb-dev libxkbcommon-x11-dev libxkbcommon-dev
```

**Fedora:**

```sh
# Основное (Cgo): компилятор, libc, X11, XTest
sudo dnf install gcc glibc-devel libX11-devel libXtst-devel

# Буфер обмена (X11)
sudo dnf install xsel xclip

# Bitmap: libpng
sudo dnf install libpng-devel

# GoHook: XKB (Fedora < 34: xorg-x11-xkb-utils-devel вместо xkbcomp-devel)
sudo dnf install libxkbcommon-devel libxkbcommon-x11-devel xkbcomp-devel
```

### Feature-specific dependencies

Нужны дополнительно к основной настройке и только для используемых
возможностей:

- **Буфер обмена в Linux (Cgo и чистый Go):** `xclip` или `xsel` в X11;
  `wl-clipboard` в Wayland:

  ```sh
  sudo apt install wl-clipboard   # Ubuntu / Debian
  sudo dnf install wl-clipboard   # Fedora
  ```

- **[Bitmap](https://github.com/vcaesar/bitmap):** заголовки libpng.
- **[GoHook](https://github.com/robotn/gohook):** заголовки XCB/XKB в Linux.

Bitmap и GoHook — отдельные Cgo-пакеты; чистый Go-бэкенд RobotGo не отменяет их
требований к C.

## Cgo-free Builds

**Экспериментальные чистые Go-бэкенды** собираются и кросс-компилируются с
`CGO_ENABLED=0` (без GCC, MinGW, Xcode и заголовков X11). Тот же путь импорта и
тот же API; бэкенд выбирается **тегом сборки** — сам по себе `CGO_ENABLED=0`
ничего не выбирает.

| Целевой `GOOS` | Тег сборки | Пакет                  | Требования времени выполнения                                    |
| -------------- | ---------- | ---------------------- | ---------------------------------------------------------------- |
| `windows`      | `win`      | [win](../win/)         | Только Win32, без дополнительных зависимостей                    |
| `darwin`       | `mac`      | [darwin](../darwin/)   | Системные фреймворки и [разрешения конфиденциальности](#macos)   |
| `linux`        | `x11`      | [x11](../x11/)         | X-сервер с XTEST, задана переменная `DISPLAY`                    |
| `linux`        | `wayland`  | [wayland](../wayland/) | Композитор wlroots, см. [Wayland](#wayland)                      |
| `linux`        | `libei`    | [libei](../libei/)     | Портал RemoteDesktop (GNOME/KDE), см. [libei](#libei-gnome--kde) |

**`purego`** — это сокращение: `mac` на macOS, `win` на Windows, `wayland` на
Linux. Добавьте `x11` или `libei`, чтобы переопределить выбор в Linux
(`-tags "purego,x11"`).

Используйте **один тег бэкенда**, соответствующий `GOOS`; комбинации вида
`x11,libei` не компилируются (единственное исключение — `purego` плюс одно
переопределение для Linux).

### Build commands

```sh
# Текущая ОС
CGO_ENABLED=0 go build -tags purego .

# Кросс-компиляция
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -tags win .
CGO_ENABLED=0 GOOS=darwin  GOARCH=arm64 go build -tags mac .
CGO_ENABLED=0 GOOS=linux   GOARCH=amd64 go build -tags x11 .
CGO_ENABLED=0 GOOS=linux   GOARCH=amd64 go build -tags wayland .
CGO_ENABLED=0 GOOS=linux   GOARCH=amd64 go build -tags libei .
```

PowerShell: `$env:CGO_ENABLED = "0"; go build -tags purego .`

Собирайте корень модуля (`.`) или свой собственный пакет, а не `./...`:
`examples/` и некоторые подпакеты используют API, доступные только с Cgo.

### Limitations

Приведённые выше требования времени выполнения остаются в силе. API, доступных
только с Cgo, таких как `CaptureScreen` / `FreeBitmap`, здесь нет; используйте
`CaptureImg`. Неподдерживаемые операции возвращают `ErrNotSupported` (или
пустой результат).

- **`mac`:** без управления окнами — `ActiveName` возвращает `ErrNotSupported`,
  `GetTitle` возвращает `""`, свернуть/развернуть/закрыть ничего не делают.
- **`wayland`, `libei`:** `Location()` возвращает последнюю позицию, заданную
  RobotGo, а не положение физического курсора. Для буфера обмена нужен
  `wl-clipboard`.

#### Wayland

Композитор wlroots (Sway, Hyprland, Wayfire, ...) с заданными
`WAYLAND_DISPLAY` и `XDG_RUNTIME_DIR`, предоставляющий:

| Возможность            | Глобальный объект протокола        |
| ---------------------- | ---------------------------------- |
| Управление мышью       | `zwlr_virtual_pointer_manager_v1`  |
| Управление клавиатурой | `zwp_virtual_keyboard_manager_v1`  |
| Захват экрана          | `zwlr_screencopy_manager_v1`       |
| Управление окнами      | `zwlr_foreign_toplevel_manager_v1` |

GNOME и KDE не предоставляют эти протоколы; там используйте `libei`.

#### libei (GNOME / KDE)

Ввод идёт через D-Bus-интерфейс **RemoteDesktop** из `xdg-desktop-portal`:
нужны сессионная шина, `xdg-desktop-portal` и реализующий его бэкенд
(`xdg-desktop-portal-gnome`, `-kde` или `-wlr`). При первом запуске появляется
диалог подтверждения; токен восстановления кэшируется в
`$XDG_STATE_HOME/robotgo`.

Абсолютное перемещение использует связанный поток ScreenCast (по умолчанию);
задайте `libei.LinkScreenCast = false`, чтобы оставить только относительное
перемещение. Захват экрана и управление окнами возвращают `ErrNotSupported`.

## Installation

```sh
go get github.com/go-vgo/robotgo
```

```go
import "github.com/go-vgo/robotgo"
```

`png.h: No such file or directory` при использовании Bitmap: см. его требования
к libpng и [issues/47](https://github.com/go-vgo/robotgo/issues/47).

## Update

```sh
go get -u github.com/go-vgo/robotgo
```

Обратите внимание на проблему кэширования компиляции C-файлов в go1.10.x, [golang #24355](https://github.com/golang/go/issues/24355).
Проблема `go mod vendor`, [golang #26366](https://github.com/golang/go/issues/26366).

## [Examples](https://github.com/go-vgo/robotgo/blob/master/examples)

#### [Мышь](https://github.com/go-vgo/robotgo/blob/master/examples/mouse/main.go)

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

  // Проверяйте возвращаемые ошибки, неверные аргументы возвращают ошибку
  if err := robotgo.Click("left", "double"); err != nil {
    fmt.Println("robotgo.Click error:", err)
  }
  if err := robotgo.ScrollDir(10, "forward"); err != nil {
    fmt.Println("robotgo.ScrollDir error:", err)
  }

  // MoveSmoothRelative сообщает о неудачном плавном перемещении как ErrSmoothMove
  err := robotgo.MoveSmoothRelative(10, -10)
  if errors.Is(err, robotgo.ErrSmoothMove) {
    fmt.Println("smooth move failed:", err)
  }

  // Бэкенды на чистом Go возвращают ErrNotSupported для неподдерживаемых операций
  err = robotgo.Move(100, 200)
  if errors.Is(err, robotgo.ErrNotSupported) {
    fmt.Println("robotgo.Move is not supported by this backend")
  }
}
```

#### [Клавиатура](https://github.com/go-vgo/robotgo/blob/master/examples/key/main.go)

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

  // Неизвестное имя клавиши возвращает ошибку вместо нажатия неверной клавиши
  if err := robotgo.KeyTap("notakey"); err != nil {
    fmt.Println("robotgo.KeyTap error:", err)
  }

  // TypeStr возвращает первую ошибку ввода, Type возвращает только количество
  if err := robotgo.TypeStr("Hello, 世界"); err != nil {
    fmt.Println("robotgo.TypeStr error:", err)
  }

  // Удерживайте модификатор, затем всегда отпускайте его и проверяйте обе ошибки
  if err := robotgo.KeyDown("shift"); err == nil {
    err = errors.Join(robotgo.KeyTap("a"), robotgo.KeyUp("shift"))
    if err != nil {
      fmt.Println("shift + a error:", err)
    }
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

- Улучшить поддержку нескольких экранов
- Обновить работу с дескрипторами окон
- Попробовать добавить поддержку Android и iOS

## Contributors

- Полный список участников см. на [странице участников](https://github.com/go-vgo/robotgo/graphs/contributors).
- См. [Руководство для участников](https://github.com/go-vgo/robotgo/blob/master/CONTRIBUTING.md).

## License

Robotgo распространяется преимущественно на условиях «Apache License (Version 2.0)», при этом отдельные части подпадают под различные лицензии в стиле BSD.

См. [LICENSE-APACHE](http://www.apache.org/licenses/LICENSE-2.0), [LICENSE](https://github.com/go-vgo/robotgo/blob/master/LICENSE).
