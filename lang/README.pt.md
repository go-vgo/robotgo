# Robotgo

[![Build Status](https://github.com/go-vgo/robotgo/workflows/Go/badge.svg)](https://github.com/go-vgo/robotgo/commits/master)
[![CircleCI Status](https://circleci.com/gh/go-vgo/robotgo.svg?style=shield)](https://circleci.com/gh/go-vgo/robotgo)
[![Go Report Card](https://goreportcard.com/badge/github.com/go-vgo/robotgo)](https://goreportcard.com/report/github.com/go-vgo/robotgo)
[![GoDoc](https://pkg.go.dev/badge/github.com/go-vgo/robotgo?status.svg)](https://pkg.go.dev/github.com/go-vgo/robotgo?tab=doc)
[![GitHub release](https://img.shields.io/github/release/go-vgo/robotgo.svg)](https://github.com/go-vgo/robotgo/releases/latest)
<a href="https://discord.gg/npPb3NzE4A"><img src="https://img.shields.io/discord/1484658282777018551.svg?logo=discord&logoColor=white&label=Discord&color=5865F2" alt="Join the Discord chat at https://discord.gg/npPb3NzE4A"></a>

[English](../README.md) | [简体中文](README.zh.md) | [繁體中文](README.zht.md) | [日本語](README.ja.md) | [한국어](README.ko.md) | [Français](README.fr.md) | [Deutsch](README.de.md) | [Español](README.es.md) | [Русский](README.ru.md) | Português

> Automação de desktop em Golang, testes automatizados e uso de computador por IA (AI Computer Use). <br>
> Controle o mouse e o teclado, leia a tela, processos, identificadores de janela (Window Handle), imagens e bitmaps, e o ouvinte global de eventos.

O RobotGo é compatível com Mac, Windows e Linux; e o robotgo também suporta arm64 e x86-amd64.

Estou construindo o [Codg](https://github.com/vcaesar/codg) agora, um sistema de agentes de IA fácil de programar e usar: automático, assíncrono, concorrente, eficiente e de alta precisão.

<p align="center">
<a href="https://github.com/vcaesar/codg" rel="nofollow">
<img width="800" alt="Codg Demo" src="https://github.com/vcaesar/codg/raw/main/demo/26-04-1.png" />
</a>
</p>

O [RobotGo-Pro](https://github.com/vcaesar/robotgo-pro) oferece versões em JavaScript, Python, Lua e outras linguagens, suporte técnico, novos recursos e a versão mais recente do robotgo ("sem versão open-source no momento").

## Índice

- [Documentação](#docs)
- [Binding](#binding)
- [Requisitos](#requirements)
- [Instalação](#installation)
- [Atualização](#update)
- [Builds sem Cgo](#cgo-free-builds)
- [Exemplos](#examples)
- [Conversão de tipos e teclas](https://github.com/go-vgo/robotgo/blob/master/docs/keys.md)
- [Compilação cruzada](https://github.com/go-vgo/robotgo/blob/master/docs/install.md#crosscompiling)
- [Autores](#authors)
- [Planos](#plans)
- [Licença](#license)

## Docs

- [GoDoc](https://godoc.org/github.com/go-vgo/robotgo) <br>
- [Documentação da API](https://github.com/go-vgo/robotgo/blob/master/docs/doc.md) (Obsoleta, sem atualizações)

## Binding:

[ADB](https://github.com/vcaesar/adb), encapsulamento da API adb do Android.

## Requirements:

Agora, certifique-se de que o `Golang, GCC` esteja instalado corretamente antes de instalar o RobotGo.

### Todas as plataformas:

```
Golang

GCC
```

#### Para MacOS:

```
brew install go
```

Ferramentas de Linha de Comando do Xcode; <br>
E, nas configurações de Privacidade, adicione Gravação de Tela e Acessibilidade em: <br>
`Ajustes do Sistema > Privacidade e Segurança > Acessibilidade, Gravação de Tela e Áudio do Sistema`.

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

ou [Mingw-w64](https://sourceforge.net/projects/mingw-w64/files)

```
winget install BrechtSanders.WinLibs.POSIX.UCRT
```

Ou baixe o [Mingw-w64](https://sourceforge.net/projects/mingw-w64/files) e os outros gcc, depois defina variáveis de ambiente do sistema como `C:\mingw64\bin` na variável `Path`.
[Defina variáveis de ambiente para executar o GCC pela linha de comando](https://www.youtube.com/results?search_query=Set+environment+variables+to+run+GCC+from+command+line).

`Ou outros GCC` (Exceto o Mingw-w64, você precisa compilar a "libpng" por conta própria ao usar o [bitmap](https://github.com/vcaesar/bitmap).)

#### Para todas as outras plataformas:

```
GCC

X11 com a extensão XTest (a biblioteca Xtst)

"Área de transferência": xsel xclip

"Bitmap": libpng (Usado apenas pelo "bitmap".)

"Evento-Gohook": xcb, xkb, libxkbcommon (Usado apenas pelo "hook".)
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

O backend Wayland é uma implementação **100% Go (sem Cgo)**, portanto nenhuma
biblioteca C do sistema é necessária. Ele requer um compositor baseado em wlroots
(Sway, Hyprland, Wayfire, ...) que suporte os seguintes protocolos:

```
zwlr_virtual_pointer_v1            (controle do mouse)
zwp_virtual_keyboard_v1            (controle do teclado)
zwlr_screencopy_v1                 (captura de tela)
zwlr_foreign_toplevel_management_v1 (gerenciamento de janelas)
```

GNOME e KDE **não** suportam esses protocolos nativamente.

#### libei (GNOME / KDE)

O backend libei também é uma implementação **100% Go (sem Cgo)**. Ele controla a
entrada por meio da interface RemoteDesktop do `xdg-desktop-portal` do
freedesktop, portanto funciona no GNOME e no KDE (ao contrário do backend Wayland
wlroots). Ele requer:

```
xdg-desktop-portal               (o serviço D-Bus do portal)
xdg-desktop-portal-gnome / -kde  (o backend de portal do seu desktop)
```

Observação: o backend libei lida apenas com entrada de mouse e teclado. A captura
de tela e o gerenciamento de janelas retornam `ErrNotSupported`.

## Installation:

Com suporte a Go module (Go 1.11+), basta importar:

```go
import "github.com/go-vgo/robotgo"
```

Caso contrário, para instalar o pacote robotgo, execute o comando:

```
go get github.com/go-vgo/robotgo
```

png.h: No such file or directory? Consulte [issues/47](https://github.com/go-vgo/robotgo/issues/47).

## Update:

```
go get -u github.com/go-vgo/robotgo
```

Observe o problema de cache de compilação de arquivos C no go1.10.x, [golang #24355](https://github.com/golang/go/issues/24355).
Problema com `go mod vendor`, [golang #26366](https://github.com/golang/go/issues/26366).

## Cgo-free Builds:

O RobotGo fornece backends **100% Go (sem Cgo)** para Windows, macOS, X11, Wayland e libei (Linux); é experimental.
Eles expõem a mesma API `robotgo`, portanto seu código não precisa de alterações — apenas uma tag de build.
Esses backends fazem compilação cruzada com `CGO_ENABLED=0` (sem GCC, MinGW, Xcode ou cabeçalhos X11).

| Backend                          | Tag de build | Pacote Go                           |
| -------------------------------- | ------------ | ----------------------------------- |
| Windows (sem Cgo)                | `win`        | `github.com/go-vgo/robotgo/win`     |
| macOS (Quartz via purego)        | `mac`        | `github.com/go-vgo/robotgo/darwin`  |
| X11 (Linux, protocolo X em Go puro) | `x11`     | `github.com/go-vgo/robotgo/x11`     |
| Wayland (Linux, wlroots)         | `wayland`    | `github.com/go-vgo/robotgo/wayland` |
| libei (Linux, portal GNOME/KDE)  | `libei`      | `github.com/go-vgo/robotgo/libei`   |
| Go puro padrão (todas as plataformas) | `purego` | seleciona `mac`/`win`/`wayland` acima |

```sh
# Backend Go puro padrão por plataforma, uma tag para todos os alvos:
# macOS -> mac, Windows -> win, Linux -> wayland (combine com x11/libei para substituir)
go build -tags purego ./...
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -tags "purego,x11" ./...

# Windows, sem Cgo / sem MinGW
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -tags win ./...

# macOS, Quartz/CoreGraphics carregado em tempo de execução via purego (sem Xcode)
CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -tags mac ./...

# X11, protocolo X em Go puro (XTEST) — sem cabeçalhos X11
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -tags x11 ./...

# Wayland, compositor baseado em wlroots (Sway, Hyprland, Wayfire, ...)
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -tags wayland ./...

# libei, GNOME/KDE via interface RemoteDesktop do xdg-desktop-portal
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -tags libei ./...
```

Com a tag `win`, o backend Cgo/Win32 padrão é excluído e as chamadas são redirecionadas para o pacote Go puro `win`; com a tag `mac`, o backend Cgo/Quartz padrão é excluído e as chamadas são redirecionadas para o pacote Go puro `darwin` (o gerenciamento de janelas retorna `ErrNotSupported`); com a tag `x11`, o backend Cgo/X11 é excluído e as chamadas são redirecionadas para o pacote Go puro `x11`; com a tag `wayland`, o backend Cgo/X11 é excluído e as chamadas são redirecionadas para o pacote Go puro `wayland`; com a tag `libei`, tanto o backend Cgo/X11 quanto o backend Wayland wlroots são excluídos e as chamadas são redirecionadas para o pacote Go puro `libei`.

A tag `purego` é um atalho multiplataforma: ela exclui o backend Cgo em todos os lugares e escolhe o backend Go puro padrão para o SO de destino — `mac` no macOS, `win` no Windows e `wayland` no Linux. No Linux, você pode combiná-la com `x11` ou `libei` (por exemplo, `-tags "purego,libei"`) para escolher outro backend Go puro.

## [Examples:](https://github.com/go-vgo/robotgo/blob/master/examples)

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

#### [Tela](https://github.com/go-vgo/robotgo/blob/master/examples/screen/main.go)

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

#### [Janela](https://github.com/go-vgo/robotgo/blob/master/examples/window/main.go)

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

- [O autor é Evans](https://github.com/vcaesar)
- [Mantenedores](https://github.com/orgs/go-vgo/people)

## Plans

- Refatorar parte do código C para Go (como x11, windows)
- Melhor suporte a múltiplas telas
- Suporte a Wayland
- Atualizar o Window Handle
- Tentar oferecer suporte a Android e iOS

## Contributors

- Veja a [página de contribuidores](https://github.com/go-vgo/robotgo/graphs/contributors) para a lista completa de contribuidores.
- Veja as [Diretrizes de Contribuição](https://github.com/go-vgo/robotgo/blob/master/CONTRIBUTING.md).

## License

O Robotgo é distribuído principalmente sob os termos da "Apache License (Version 2.0)", com partes cobertas por diversas licenças no estilo BSD.

Veja [LICENSE-APACHE](http://www.apache.org/licenses/LICENSE-2.0), [LICENSE](https://github.com/go-vgo/robotgo/blob/master/LICENSE).
