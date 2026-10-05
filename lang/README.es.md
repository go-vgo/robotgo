# Robotgo

<p align="center">
  <img src="../docs/robotgo-logo.svg" width="560" alt="Logotipo de RobotGo — un robot con un puntero de ratón" />
</p>

[![Build Status](https://github.com/go-vgo/robotgo/workflows/Go/badge.svg)](https://github.com/go-vgo/robotgo/commits/master)
[![CircleCI Status](https://circleci.com/gh/go-vgo/robotgo.svg?style=shield)](https://circleci.com/gh/go-vgo/robotgo)
[![golangci-lint](https://github.com/go-vgo/robotgo/actions/workflows/lint.yml/badge.svg)](https://github.com/go-vgo/robotgo/actions/workflows/lint.yml)
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
- [Compilaciones sin Cgo](#cgo-free-builds)
- [Instalación](#installation)
- [Actualización](#update)
- [Ejemplos](#examples)
- [Conversión de tipos y teclas](https://github.com/go-vgo/robotgo/blob/master/docs/keys.md)
- [Compilación cruzada](https://github.com/go-vgo/robotgo/blob/master/docs/install.md#crosscompiling)
- [Autores](#authors)
- [Planes](#plans)
- [Licencia](#license)

## Docs

- [GoDoc](https://godoc.org/github.com/go-vgo/robotgo) <br>
- [Documentación de la API](https://github.com/go-vgo/robotgo/blob/master/docs/doc.md) (Obsoleta, sin actualizar)

## Binding

[ADB](https://github.com/vcaesar/adb), que empaqueta la API de adb de Android.

## Requirements

**Go 1.26+** (consulte [go.mod](../go.mod)) y un backend:

- **Predeterminado (Cgo):** `CGO_ENABLED=1`, un compilador de C y las bibliotecas
  de la plataforma indicadas abajo.
- **Go puro (experimental):** sin cadena de herramientas de C; consulte
  [Compilaciones sin Cgo](#cgo-free-builds).

Ambos necesitan las [dependencias por función](#feature-specific-dependencies)
de las funciones que utilice.

### Default Cgo setup

#### macOS

Go y las Xcode Command Line Tools (Clang):

```sh
brew install go
xcode-select --install
```

**Permisos (Cgo y Go puro):** en **Ajustes del Sistema > Privacidad y
seguridad**, conceda a la aplicación o al terminal **Accesibilidad** (entrada) y
**Grabación de pantalla** (captura).

#### Windows

```powershell
winget install GoLang.Go
```

Más **una** cadena de herramientas de C. [LLVM-MinGW](https://github.com/mstorsjo/llvm-mingw):

```powershell
winget install MartinStorsjo.LLVM-MinGW.UCRT
$env:CC = "clang"
```

O [MinGW-w64 (WinLibs)](https://winlibs.com/):

```powershell
winget install BrechtSanders.WinLibs.POSIX.UCRT
$env:CC = "gcc"
```

El directorio `bin` de la cadena de herramientas (p. ej. `C:\mingw64\bin`) debe
estar en `PATH` y coincidir con `GOARCH`. También funciona cualquier otro
compilador compatible con Cgo, como una instalación manual de
[MinGW-w64](https://sourceforge.net/projects/mingw-w64/files).

[Bitmap](https://github.com/vcaesar/bitmap) incluye libpng solo para MinGW-w64;
con otra cadena de herramientas, compile libpng usted mismo.

#### Linux (X11)

Compilación: GCC, encabezados de libc y de X11/XTest. Ejecución: un servidor X
con XTEST y `DISPLAY` definida. Para Wayland nativo, consulte
[Compilaciones sin Cgo](#cgo-free-builds).

**Ubuntu / Debian:**

```sh
# Go (Snap si el paquete de la distribución es demasiado antiguo)
sudo snap install go --classic
# o: sudo apt install golang

# Base (Cgo): compilador, libc, X11, XTest
sudo apt install gcc libc6-dev libx11-dev xorg-dev libxtst-dev

# Portapapeles (X11)
sudo apt install xsel xclip

# Bitmap: libpng
sudo apt install libpng++-dev

# GoHook: XCB, XKB
sudo apt install xcb libxcb-xkb-dev x11-xkb-utils libx11-xcb-dev libxkbcommon-x11-dev libxkbcommon-dev
```

**Fedora:**

```sh
# Base (Cgo): compilador, libc, X11, XTest
sudo dnf install gcc glibc-devel libX11-devel libXtst-devel

# Portapapeles (X11)
sudo dnf install xsel xclip

# Bitmap: libpng
sudo dnf install libpng-devel

# GoHook: XKB (Fedora < 34: xorg-x11-xkb-utils-devel en lugar de xkbcomp-devel)
sudo dnf install libxkbcommon-devel libxkbcommon-x11-devel xkbcomp-devel
```

### Feature-specific dependencies

Necesarias además de la configuración base, solo para las funciones que utilice:

- **Portapapeles en Linux (Cgo y Go puro):** `xclip` o `xsel` en X11;
  `wl-clipboard` en Wayland:

  ```sh
  sudo apt install wl-clipboard   # Ubuntu / Debian
  sudo dnf install wl-clipboard   # Fedora
  ```

- **[Bitmap](https://github.com/vcaesar/bitmap):** encabezados de libpng.
- **[GoHook](https://github.com/robotn/gohook):** encabezados de XCB/XKB en Linux.

Bitmap y GoHook son paquetes Cgo independientes; un backend de RobotGo en Go puro
no elimina sus requisitos de C.

## Cgo-free Builds

Los **backends experimentales en Go puro** se compilan y se compilan de forma
cruzada con `CGO_ENABLED=0` (sin GCC, MinGW, Xcode ni encabezados X11). La misma
ruta de importación y la misma API; una **etiqueta de compilación** selecciona el
backend: `CGO_ENABLED=0` por sí sola no selecciona nada.

| `GOOS` de destino | Etiqueta  | Paquete                | Requisitos de ejecución                                               |
| ----------------- | --------- | ---------------------- | --------------------------------------------------------------------- |
| `windows`         | `win`     | [win](../win/)         | Solo Win32, sin runtime adicional                                     |
| `darwin`          | `mac`     | [darwin](../darwin/)   | Frameworks del sistema y [permisos de privacidad](#macos)             |
| `linux`           | `x11`     | [x11](../x11/)         | Servidor X con XTEST y `DISPLAY` definida                             |
| `linux`           | `wayland` | [wayland](../wayland/) | Compositor wlroots, consulte [Wayland](#wayland)                      |
| `linux`           | `libei`   | [libei](../libei/)     | Portal RemoteDesktop (GNOME/KDE), consulte [libei](#libei-gnome--kde) |

**`purego`** es un atajo: `mac` en macOS, `win` en Windows, `wayland` en Linux.
Añada `x11` o `libei` para cambiarlo en Linux (`-tags "purego,x11"`).

Use **una sola etiqueta de backend** acorde con `GOOS`; combinaciones como
`x11,libei` no compilan (`purego` más una sustitución en Linux es la única
excepción).

### Build commands

```sh
# Sistema operativo actual
CGO_ENABLED=0 go build -tags purego .

# Compilación cruzada
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -tags win .
CGO_ENABLED=0 GOOS=darwin  GOARCH=arm64 go build -tags mac .
CGO_ENABLED=0 GOOS=linux   GOARCH=amd64 go build -tags x11 .
CGO_ENABLED=0 GOOS=linux   GOARCH=amd64 go build -tags wayland .
CGO_ENABLED=0 GOOS=linux   GOARCH=amd64 go build -tags libei .
```

PowerShell: `$env:CGO_ENABLED = "0"; go build -tags purego .`

Compile la raíz del módulo (`.`) o su propio paquete, no `./...`: `examples/` y
algunos subpaquetes necesitan APIs exclusivas de Cgo.

### Limitations

Los requisitos de ejecución anteriores siguen vigentes. Las APIs exclusivas de
Cgo, como `CaptureScreen` / `FreeBitmap`, no existen; use `CaptureImg`. Las
operaciones no soportadas devuelven `ErrNotSupported` (o un resultado vacío).

- **`mac`:** sin gestión de ventanas: `ActiveName` devuelve `ErrNotSupported`,
  `GetTitle` devuelve `""` y minimizar/maximizar/cerrar no hacen nada.
- **`wayland`, `libei`:** `Location()` devuelve la última posición inyectada por
  RobotGo, no la del cursor físico. El portapapeles necesita `wl-clipboard`.

#### Wayland

Un compositor wlroots (Sway, Hyprland, Wayfire, ...) con `WAYLAND_DISPLAY` y
`XDG_RUNTIME_DIR` definidas, que exponga:

| Función             | Global del protocolo               |
| ------------------- | ---------------------------------- |
| Control del ratón   | `zwlr_virtual_pointer_manager_v1`  |
| Control del teclado | `zwp_virtual_keyboard_manager_v1`  |
| Captura de pantalla | `zwlr_screencopy_manager_v1`       |
| Gestión de ventanas | `zwlr_foreign_toplevel_manager_v1` |

GNOME y KDE no ofrecen estos protocolos; en ese caso use `libei`.

#### libei (GNOME / KDE)

La entrada pasa por la interfaz D-Bus **RemoteDesktop** de `xdg-desktop-portal`:
requiere un bus de sesión, `xdg-desktop-portal` y un backend que la implemente
(`xdg-desktop-portal-gnome`, `-kde` o `-wlr`). La primera ejecución muestra un
diálogo de consentimiento; el token de restauración se guarda en
`$XDG_STATE_HOME/robotgo`.

El movimiento absoluto usa un flujo ScreenCast vinculado (predeterminado); defina
`libei.LinkScreenCast = false` para usar solo movimiento relativo. La captura de
pantalla y la gestión de ventanas devuelven `ErrNotSupported`.

## Installation

```sh
go get github.com/go-vgo/robotgo
```

```go
import "github.com/go-vgo/robotgo"
```

`png.h: No such file or directory` con Bitmap: consulte sus requisitos de libpng
y [issues/47](https://github.com/go-vgo/robotgo/issues/47).

## Update

```sh
go get -u github.com/go-vgo/robotgo
```

Tenga en cuenta el problema de la caché de compilación de archivos C en go1.10.x, [golang #24355](https://github.com/golang/go/issues/24355).
Problema de `go mod vendor`, [golang #26366](https://github.com/golang/go/issues/26366).

## [Examples](https://github.com/go-vgo/robotgo/blob/master/examples)

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

- Mejor soporte multipantalla
- Actualizar el manejador de ventanas
- Intentar dar soporte a Android e iOS

## Contributors

- Consulte la [página de colaboradores](https://github.com/go-vgo/robotgo/graphs/contributors) para ver la lista completa de colaboradores.
- Consulte las [Directrices de contribución](https://github.com/go-vgo/robotgo/blob/master/CONTRIBUTING.md).

## License

Robotgo se distribuye principalmente bajo los términos de «la Licencia Apache (Versión 2.0)», con partes cubiertas por varias licencias de tipo BSD.

Consulte [LICENSE-APACHE](http://www.apache.org/licenses/LICENSE-2.0), [LICENSE](https://github.com/go-vgo/robotgo/blob/master/LICENSE).
