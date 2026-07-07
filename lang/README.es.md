# Robotgo

[![Build Status](https://github.com/go-vgo/robotgo/workflows/Go/badge.svg)](https://github.com/go-vgo/robotgo/commits/master)
[![CircleCI Status](https://circleci.com/gh/go-vgo/robotgo.svg?style=shield)](https://circleci.com/gh/go-vgo/robotgo)
[![Go Report Card](https://goreportcard.com/badge/github.com/go-vgo/robotgo)](https://goreportcard.com/report/github.com/go-vgo/robotgo)
[![GoDoc](https://pkg.go.dev/badge/github.com/go-vgo/robotgo?status.svg)](https://pkg.go.dev/github.com/go-vgo/robotgo?tab=doc)
[![GitHub release](https://img.shields.io/github/release/go-vgo/robotgo.svg)](https://github.com/go-vgo/robotgo/releases/latest)
<a href="https://discord.gg/npPb3NzE4A"><img src="https://img.shields.io/discord/1484658282777018551.svg?logo=discord&logoColor=white&label=Discord&color=5865F2" alt="Join the Discord chat at https://discord.gg/npPb3NzE4A"></a>

[English](../README.md) | [简体中文](README.zh.md) | [繁體中文](README.zht.md) | [日本語](README.ja.md) | [한국어](README.ko.md) | [Français](README.fr.md) | [Deutsch](README.de.md) | Español | [Русский](README.ru.md) | [Português](README.pt.md)

> Automatización de escritorio con Golang, pruebas automáticas y uso del ordenador con IA (AI Computer Use). <br>
> Controla el ratón, el teclado, lee la pantalla, procesos, manejadores de ventanas, imágenes y mapas de bits, y escucha de eventos globales.

RobotGo es compatible con Mac, Windows y Linux; y robotgo es compatible con arm64 y x86-amd64.

Ahora estoy construyendo [Codg](https://github.com/vcaesar/codg), un sistema de agentes de IA sencillo para programar y trabajar: automático, asíncrono, concurrente, eficiente y de alta precisión.

<p align="center">
<a href="https://github.com/vcaesar/codg" rel="nofollow">
<img width="800" alt="Codg Demo" src="https://github.com/vcaesar/codg/raw/main/demo/26-04-1.png" />
</a>
</p>

[RobotGo-Pro](https://github.com/vcaesar/robotgo-pro) ofrece las versiones en JavaScript, Python, Lua y otras, soporte técnico, nuevas funciones y la versión más reciente de robotgo («sin versión de código abierto por ahora»).

## Contenido

- [Documentación](#docs)
- [Bindings](#binding)
- [Requisitos](#requirements)
- [Instalación](#installation)
- [Actualización](#update)
- [Compilaciones sin Cgo](#cgo-free-builds)
- [Ejemplos](#examples)
- [Conversión de tipos y teclas](https://github.com/go-vgo/robotgo/blob/master/docs/keys.md)
- [Compilación cruzada](https://github.com/go-vgo/robotgo/blob/master/docs/install.md#crosscompiling)
- [Autores](#authors)
- [Planes](#plans)
- [Licencia](#license)

## Docs

- [GoDoc](https://godoc.org/github.com/go-vgo/robotgo) <br>
- [Documentación de la API](https://github.com/go-vgo/robotgo/blob/master/docs/doc.md) (Obsoleta, sin actualizar)

## Binding:

[ADB](https://github.com/vcaesar/adb), que empaqueta la API de adb de Android.

## Requirements:

Ahora, asegúrese de que `Golang, GCC` estén instalados correctamente antes de instalar RobotGo.

### TODOS:

```
Golang

GCC
```

#### Para MacOS:

```
brew install go
```

Herramientas de línea de comandos de Xcode; <br>
Y en los ajustes de privacidad, añada «Grabación de pantalla» y «Accesibilidad» en: <br>
`Ajustes del sistema > Privacidad y seguridad > Accesibilidad, Grabación de pantalla y audio del sistema`.

```
xcode-select --install
```

#### Para Windows:

```
winget install Golang.go
```

[llvm-mingw](https://github.com/mstorsjo/llvm-mingw)

```
winget install MartinStorsjo.LLVM-MinGW.UCRT
```

o [Mingw-w64](https://sourceforge.net/projects/mingw-w64/files)

```
winget install BrechtSanders.WinLibs.POSIX.UCRT
```

O descargue [Mingw-w64](https://sourceforge.net/projects/mingw-w64/files) y los demás gcc, luego configure variables de entorno del sistema como `C:\mingw64\bin` en la variable `Path`.
[Configurar variables de entorno para ejecutar GCC desde la línea de comandos](https://www.youtube.com/results?search_query=Set+environment+variables+to+run+GCC+from+command+line).

`O los demás GCC` (Excepto Mingw-w64, debe compilar «libpng» usted mismo cuando use el [bitmap](https://github.com/vcaesar/bitmap).)

#### Para todo lo demás:

```
GCC

X11 con la extensión XTest (la biblioteca Xtst)

"Clipboard": xsel xclip

"Bitmap": libpng (Usado solo por "bitmap".)

"Event-Gohook": xcb, xkb, libxkbcommon (Usado solo por "hook".)
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

El backend de Wayland es una implementación **100 % Go (sin Cgo)**, por lo que no
se requiere ninguna biblioteca C del sistema. Necesita un compositor basado en
wlroots (Sway, Hyprland, Wayfire, ...) que admita los siguientes protocolos:

```
zwlr_virtual_pointer_v1            (control del ratón)
zwp_virtual_keyboard_v1            (control del teclado)
zwlr_screencopy_v1                 (captura de pantalla)
zwlr_foreign_toplevel_management_v1 (gestión de ventanas)
```

GNOME y KDE **no** admiten estos protocolos de forma nativa.

#### libei (GNOME / KDE)

El backend de libei también es una implementación **100 % Go (sin Cgo)**. Controla
la entrada a través de la interfaz RemoteDesktop de `xdg-desktop-portal` de
freedesktop, por lo que funciona en GNOME y KDE (a diferencia del backend Wayland
wlroots). Requiere:

```
xdg-desktop-portal               (el servicio D-Bus del portal)
xdg-desktop-portal-gnome / -kde  (el backend de portal de tu escritorio)
```

Nota: el backend de libei solo gestiona la entrada de ratón y teclado. La captura
de pantalla y la gestión de ventanas devuelven `ErrNotSupported`.

## Installation:

Con soporte para módulos de Go (Go 1.11+), solo importe:

```go
import "github.com/go-vgo/robotgo"
```

De lo contrario, para instalar el paquete robotgo, ejecute el comando:

```
go get github.com/go-vgo/robotgo
```

png.h: No such file or directory? Consulte [issues/47](https://github.com/go-vgo/robotgo/issues/47).

## Update:

```
go get -u github.com/go-vgo/robotgo
```

Tenga en cuenta el problema de la caché de compilación de archivos C en go1.10.x, [golang #24355](https://github.com/golang/go/issues/24355).
Problema de `go mod vendor`, [golang #26366](https://github.com/golang/go/issues/26366).

## Cgo-free Builds:

RobotGo ofrece backends **100 % Go (sin Cgo)** para Windows, macOS, X11, Wayland y libei (Linux); es experimental.
Exponen la misma API `robotgo`, por lo que tu código no necesita cambios — solo una etiqueta de compilación.
Estos backends se compilan de forma cruzada con `CGO_ENABLED=0` (sin GCC, MinGW, Xcode ni encabezados X11).

| Backend                          | Etiqueta de build | Paquete Go                          |
| -------------------------------- | ----------------- | ----------------------------------- |
| Windows (sin Cgo)                | `win`             | `github.com/go-vgo/robotgo/win`     |
| macOS (Quartz vía purego)        | `mac`             | `github.com/go-vgo/robotgo/darwin`  |
| X11 (Linux, protocolo X en Go puro) | `x11`          | `github.com/go-vgo/robotgo/x11`     |
| Wayland (Linux, wlroots)         | `wayland`         | `github.com/go-vgo/robotgo/wayland` |
| libei (Linux, portal GNOME/KDE)  | `libei`           | `github.com/go-vgo/robotgo/libei`   |
| Go puro predeterminado (todas las plataformas) | `purego` | selecciona `mac`/`win`/`wayland` anteriores |

```sh
# Backend Go puro predeterminado por plataforma, una etiqueta para todos los destinos:
# macOS -> mac, Windows -> win, Linux -> wayland (combínala con x11/libei para sobrescribir)
go build -tags purego ./...
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -tags "purego,x11" ./...

# Windows, sin Cgo / sin MinGW
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -tags win ./...

# macOS, Quartz/CoreGraphics cargado en tiempo de ejecución vía purego (sin Xcode)
CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -tags mac ./...

# X11, protocolo X en Go puro (XTEST) — sin encabezados X11
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -tags x11 ./...

# Wayland, compositor basado en wlroots (Sway, Hyprland, Wayfire, ...)
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -tags wayland ./...

# libei, GNOME/KDE a través de la interfaz RemoteDesktop de xdg-desktop-portal
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -tags libei ./...
```

Con la etiqueta `win`, el backend Cgo/Win32 predeterminado se excluye y las llamadas se reenvían al paquete Go puro `win`; con la etiqueta `mac`, el backend Cgo/Quartz predeterminado se excluye y las llamadas se reenvían al paquete Go puro `darwin` (la gestión de ventanas devuelve `ErrNotSupported`); con la etiqueta `x11`, el backend Cgo/X11 se excluye y las llamadas se reenvían al paquete Go puro `x11`; con la etiqueta `wayland`, el backend Cgo/X11 se excluye y las llamadas se reenvían al paquete Go puro `wayland`; con la etiqueta `libei`, se excluyen tanto el backend Cgo/X11 como el backend Wayland wlroots y las llamadas se reenvían al paquete Go puro `libei`.

La etiqueta `purego` es un atajo multiplataforma: excluye el backend Cgo en todas partes y elige el backend Go puro predeterminado para el sistema operativo de destino — `mac` en macOS, `win` en Windows y `wayland` en Linux. En Linux puedes combinarla con `x11` o `libei` (por ejemplo, `-tags "purego,libei"`) para elegir otro backend Go puro.

## [Examples:](https://github.com/go-vgo/robotgo/blob/master/examples)

#### [Ratón](https://github.com/go-vgo/robotgo/blob/master/examples/mouse/main.go)

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

#### [Teclado](https://github.com/go-vgo/robotgo/blob/master/examples/key/main.go)

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

#### [Pantalla](https://github.com/go-vgo/robotgo/blob/master/examples/screen/main.go)

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

#### [Evento](https://github.com/robotn/gohook/blob/master/examples/main.go)

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

#### [Ventana](https://github.com/go-vgo/robotgo/blob/master/examples/window/main.go)

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

- [El autor es Evans](https://github.com/vcaesar)
- [Mantenedores](https://github.com/orgs/go-vgo/people)

## Plans

- Refactorizar parte del código C a Go (como x11, windows)
- Mejor soporte multipantalla
- Soporte para Wayland
- Actualizar el manejador de ventanas
- Intentar dar soporte a Android e iOS

## Contributors

- Consulte la [página de colaboradores](https://github.com/go-vgo/robotgo/graphs/contributors) para ver la lista completa de colaboradores.
- Consulte las [Directrices de contribución](https://github.com/go-vgo/robotgo/blob/master/CONTRIBUTING.md).

## License

Robotgo se distribuye principalmente bajo los términos de «la Licencia Apache (Versión 2.0)», con partes cubiertas por varias licencias de tipo BSD.

Consulte [LICENSE-APACHE](http://www.apache.org/licenses/LICENSE-2.0), [LICENSE](https://github.com/go-vgo/robotgo/blob/master/LICENSE).
