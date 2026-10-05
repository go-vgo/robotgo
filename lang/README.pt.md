# Robotgo

<p align="center">
  <img src="../docs/robotgo-logo.svg" width="560" alt="Logotipo do RobotGo — um robô com um ponteiro de mouse" />
</p>

[![Build Status](https://github.com/go-vgo/robotgo/workflows/Go/badge.svg)](https://github.com/go-vgo/robotgo/commits/master)
[![CircleCI Status](https://circleci.com/gh/go-vgo/robotgo.svg?style=shield)](https://circleci.com/gh/go-vgo/robotgo)
[![codecov](https://codecov.io/gh/go-vgo/robotgo/branch/master/graph/badge.svg)](https://codecov.io/gh/go-vgo/robotgo)
[![golangci-lint](https://github.com/go-vgo/robotgo/actions/workflows/lint.yml/badge.svg)](https://github.com/go-vgo/robotgo/actions/workflows/lint.yml)
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
- [Builds sem Cgo](#cgo-free-builds)
- [Instalação](#installation)
- [Atualização](#update)
- [Exemplos](#examples)
- [Conversão de tipos e teclas](https://github.com/go-vgo/robotgo/blob/master/docs/keys.md)
- [Compilação cruzada](https://github.com/go-vgo/robotgo/blob/master/docs/install.md#crosscompiling)
- [Autores](#authors)
- [Planos](#plans)
- [Licença](#license)

## Docs

- [GoDoc](https://godoc.org/github.com/go-vgo/robotgo) <br>
- [Documentação da API](https://github.com/go-vgo/robotgo/blob/master/docs/doc.md) (Obsoleta, sem atualizações)

## Binding

[ADB](https://github.com/vcaesar/adb), encapsulamento da API adb do Android.

## Requirements

**Go 1.26+** (veja [go.mod](../go.mod)) e um backend:

- **Padrão (Cgo):** `CGO_ENABLED=1`, um compilador C e as bibliotecas da
  plataforma listadas abaixo.
- **Go puro (experimental):** sem toolchain C; veja
  [Builds sem Cgo](#cgo-free-builds).

Ambos precisam das [dependências por recurso](#feature-specific-dependencies)
dos recursos que você usar.

### Default Cgo setup

#### macOS

Go e as Xcode Command Line Tools (Clang):

```sh
brew install go
xcode-select --install
```

**Permissões (Cgo e Go puro):** em **Ajustes do Sistema > Privacidade e
Segurança**, conceda ao app ou ao terminal **Acessibilidade** (entrada) e
**Gravação de Tela** (captura).

#### Windows

```powershell
winget install GoLang.Go
```

Mais **uma** toolchain C. [LLVM-MinGW](https://github.com/mstorsjo/llvm-mingw):

```powershell
winget install MartinStorsjo.LLVM-MinGW.UCRT
$env:CC = "clang"
```

Ou [MinGW-w64 (WinLibs)](https://winlibs.com/):

```powershell
winget install BrechtSanders.WinLibs.POSIX.UCRT
$env:CC = "gcc"
```

O diretório `bin` da toolchain (por exemplo, `C:\mingw64\bin`) precisa estar no
`PATH` e corresponder ao `GOARCH`. Qualquer outro compilador compatível com Cgo
também funciona, como uma instalação manual do
[MinGW-w64](https://sourceforge.net/projects/mingw-w64/files).

O [Bitmap](https://github.com/vcaesar/bitmap) inclui a libpng apenas para
MinGW-w64; com outra toolchain, compile a libpng por conta própria.

#### Linux (X11)

Build: GCC, cabeçalhos da libc e cabeçalhos X11/XTest. Execução: um servidor X
com XTEST e `DISPLAY` definido. Para Wayland nativo, veja
[Builds sem Cgo](#cgo-free-builds).

**Ubuntu / Debian:**

```sh
# Go (Snap, se o pacote da distribuição for muito antigo)
sudo snap install go --classic
# ou: sudo apt install golang

# Base (Cgo): compilador, libc, X11, XTest
sudo apt install gcc libc6-dev libx11-dev xorg-dev libxtst-dev

# Área de transferência (X11)
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

# Área de transferência (X11)
sudo dnf install xsel xclip

# Bitmap: libpng
sudo dnf install libpng-devel

# GoHook: XKB (Fedora < 34: xorg-x11-xkb-utils-devel em vez de xkbcomp-devel)
sudo dnf install libxkbcommon-devel libxkbcommon-x11-devel xkbcomp-devel
```

### Feature-specific dependencies

Necessárias além da configuração base, apenas para os recursos que você usar:

- **Área de transferência no Linux (Cgo e Go puro):** `xclip` ou `xsel` no X11;
  `wl-clipboard` no Wayland:

  ```sh
  sudo apt install wl-clipboard   # Ubuntu / Debian
  sudo dnf install wl-clipboard   # Fedora
  ```

- **[Bitmap](https://github.com/vcaesar/bitmap):** cabeçalhos da libpng.
- **[GoHook](https://github.com/robotn/gohook):** cabeçalhos XCB/XKB no Linux.

Bitmap e GoHook são pacotes Cgo separados; um backend do RobotGo em Go puro não
elimina os requisitos de C deles.

## Cgo-free Builds

Os **backends experimentais em Go puro** compilam e fazem compilação cruzada com
`CGO_ENABLED=0` (sem GCC, MinGW, Xcode ou cabeçalhos X11). Mesmo caminho de
importação e mesma API; uma **tag de build** seleciona o backend —
`CGO_ENABLED=0` sozinho não seleciona nada.

| `GOOS` de destino | Tag de build | Pacote                 | Requisitos de execução                                            |
| ----------------- | ------------ | ---------------------- | ----------------------------------------------------------------- |
| `windows`         | `win`        | [win](../win/)         | Apenas Win32, sem runtime extra                                   |
| `darwin`          | `mac`        | [darwin](../darwin/)   | Frameworks do sistema e [permissões de privacidade](#macos)       |
| `linux`           | `x11`        | [x11](../x11/)         | Servidor X com XTEST, `DISPLAY` definido                          |
| `linux`           | `wayland`    | [wayland](../wayland/) | Compositor wlroots, veja [Wayland](#wayland)                      |
| `linux`           | `libei`      | [libei](../libei/)     | Portal RemoteDesktop (GNOME/KDE), veja [libei](#libei-gnome--kde) |

**`purego`** é um atalho: `mac` no macOS, `win` no Windows, `wayland` no Linux.
Adicione `x11` ou `libei` para substituir no Linux (`-tags "purego,x11"`).

Use **uma única tag de backend** correspondente ao `GOOS`; combinações como
`x11,libei` não compilam (`purego` mais uma substituição no Linux é a única
exceção).

### Build commands

```sh
# Sistema operacional atual
CGO_ENABLED=0 go build -tags purego .

# Compilação cruzada
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -tags win .
CGO_ENABLED=0 GOOS=darwin  GOARCH=arm64 go build -tags mac .
CGO_ENABLED=0 GOOS=linux   GOARCH=amd64 go build -tags x11 .
CGO_ENABLED=0 GOOS=linux   GOARCH=amd64 go build -tags wayland .
CGO_ENABLED=0 GOOS=linux   GOARCH=amd64 go build -tags libei .
```

PowerShell: `$env:CGO_ENABLED = "0"; go build -tags purego .`

Compile a raiz do módulo (`.`) ou o seu próprio pacote, não `./...`: `examples/`
e alguns subpacotes precisam de APIs exclusivas do Cgo.

### Limitations

Os requisitos de execução acima continuam valendo. APIs exclusivas do Cgo, como
`CaptureScreen` / `FreeBitmap`, não existem; use `CaptureImg`. Operações sem
suporte retornam `ErrNotSupported` (ou um resultado vazio).

- **`mac`:** sem gerenciamento de janelas — `ActiveName` retorna
  `ErrNotSupported`, `GetTitle` retorna `""` e minimizar/maximizar/fechar não
  fazem nada.
- **`wayland`, `libei`:** `Location()` retorna a última posição injetada pelo
  RobotGo, não a do cursor físico. A área de transferência precisa do
  `wl-clipboard`.

#### Wayland

Um compositor wlroots (Sway, Hyprland, Wayfire, ...) com `WAYLAND_DISPLAY` e
`XDG_RUNTIME_DIR` definidos, expondo:

| Recurso                  | Global do protocolo                |
| ------------------------ | ---------------------------------- |
| Controle do mouse        | `zwlr_virtual_pointer_manager_v1`  |
| Controle do teclado      | `zwp_virtual_keyboard_manager_v1`  |
| Captura de tela          | `zwlr_screencopy_manager_v1`       |
| Gerenciamento de janelas | `zwlr_foreign_toplevel_manager_v1` |

GNOME e KDE não fornecem esses protocolos; neles, use `libei`.

#### libei (GNOME / KDE)

A entrada passa pela interface D-Bus **RemoteDesktop** do `xdg-desktop-portal`:
precisa de um barramento de sessão, do `xdg-desktop-portal` e de um backend que
o implemente (`xdg-desktop-portal-gnome`, `-kde` ou `-wlr`). A primeira execução
exibe um diálogo de consentimento; o token de restauração fica em cache em
`$XDG_STATE_HOME/robotgo`.

O movimento absoluto usa um stream ScreenCast vinculado (padrão); defina
`libei.LinkScreenCast = false` para apenas movimento relativo. Captura de tela e
gerenciamento de janelas retornam `ErrNotSupported`.

## Installation

```sh
go get github.com/go-vgo/robotgo
```

```go
import "github.com/go-vgo/robotgo"
```

`png.h: No such file or directory` com o Bitmap: veja os requisitos de libpng
dele e [issues/47](https://github.com/go-vgo/robotgo/issues/47).

## Update

```sh
go get -u github.com/go-vgo/robotgo
```

Observe o problema de cache de compilação de arquivos C no go1.10.x, [golang #24355](https://github.com/golang/go/issues/24355).
Problema com `go mod vendor`, [golang #26366](https://github.com/golang/go/issues/26366).

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

- Melhor suporte a múltiplas telas
- Atualizar o Window Handle
- Tentar oferecer suporte a Android e iOS

## Contributors

- Veja a [página de contribuidores](https://github.com/go-vgo/robotgo/graphs/contributors) para a lista completa de contribuidores.
- Veja as [Diretrizes de Contribuição](https://github.com/go-vgo/robotgo/blob/master/CONTRIBUTING.md).

## License

O Robotgo é distribuído principalmente sob os termos da "Apache License (Version 2.0)", com partes cobertas por diversas licenças no estilo BSD.

Veja [LICENSE-APACHE](http://www.apache.org/licenses/LICENSE-2.0), [LICENSE](https://github.com/go-vgo/robotgo/blob/master/LICENSE).
