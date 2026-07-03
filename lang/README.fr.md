# Robotgo

[![Build Status](https://github.com/go-vgo/robotgo/workflows/Go/badge.svg)](https://github.com/go-vgo/robotgo/commits/master)
[![CircleCI Status](https://circleci.com/gh/go-vgo/robotgo.svg?style=shield)](https://circleci.com/gh/go-vgo/robotgo)
[![Go Report Card](https://goreportcard.com/badge/github.com/go-vgo/robotgo)](https://goreportcard.com/report/github.com/go-vgo/robotgo)
[![GoDoc](https://pkg.go.dev/badge/github.com/go-vgo/robotgo?status.svg)](https://pkg.go.dev/github.com/go-vgo/robotgo?tab=doc)
[![GitHub release](https://img.shields.io/github/release/go-vgo/robotgo.svg)](https://github.com/go-vgo/robotgo/releases/latest)
<a href="https://discord.gg/npPb3NzE4A"><img src="https://img.shields.io/discord/1484658282777018551.svg?logo=discord&logoColor=white&label=Discord&color=5865F2" alt="Join the Discord chat at https://discord.gg/npPb3NzE4A"></a>

[English](../README.md) | [简体中文](README.zh.md) | [繁體中文](README.zht.md) | [日本語](README.ja.md) | [한국어](README.ko.md) | Français | [Deutsch](README.de.md) | [Español](README.es.md) | [Русский](README.ru.md) | [Português](README.pt.md)

> Automatisation de bureau en Golang, tests automatisés et utilisation de l'ordinateur par l'IA (Computer Use). <br>
> Contrôlez la souris et le clavier, lisez l'écran, gérez les processus, les handles de fenêtre, les images et les bitmaps, ainsi que l'écoute globale des événements.

RobotGo prend en charge Mac, Windows et Linux ; et robotgo prend en charge les architectures arm64 et x86-amd64.

Je développe actuellement [Codg](https://github.com/vcaesar/codg), un système de travail à base d'agents IA simple à utiliser : automatisé, asynchrone, concurrent, efficace et d'une grande précision.

<p align="center">
<a href="https://github.com/vcaesar/codg" rel="nofollow">
<img width="800" alt="Codg Demo" src="https://github.com/vcaesar/codg/raw/main/demo/26-04-1.png" />
</a>
</p>

[RobotGo-Pro](https://github.com/vcaesar/robotgo-pro) propose les versions JavaScript, Python, Lua et d'autres langages, le support technique, de nouvelles fonctionnalités ainsi que la toute dernière version de robotgo (« aucune version open source pour le moment »).

## Sommaire

- [Documentation](#docs)
- [Binding](#binding)
- [Prérequis](#requirements)
- [Installation](#installation)
- [Mise à jour](#update)
- [Builds sans Cgo](#cgo-free-builds)
- [Exemples](#examples)
- [Conversion de types et touches](https://github.com/go-vgo/robotgo/blob/master/docs/keys.md)
- [Compilation croisée](https://github.com/go-vgo/robotgo/blob/master/docs/install.md#crosscompiling)
- [Auteurs](#authors)
- [Projets](#plans)
- [Licence](#license)

## Docs

- [GoDoc](https://godoc.org/github.com/go-vgo/robotgo) <br>
- [Documentation de l'API](https://github.com/go-vgo/robotgo/blob/master/docs/doc.md) (obsolète, plus mise à jour)

## Binding:

[ADB](https://github.com/vcaesar/adb), encapsulation de l'API adb d'Android.

## Requirements:

Désormais, veuillez vous assurer que `Golang, GCC` sont correctement installés avant d'installer RobotGo.

### Toutes les plateformes :

```
Golang

GCC
```

#### Pour MacOS :

```
brew install go
```

Outils en ligne de commande Xcode ; <br>
Et dans les réglages de confidentialité, ajoutez « Enregistrement de l'écran » et « Accessibilité » dans : <br>
`Réglages Système > Confidentialité et sécurité > Accessibilité, Enregistrement de l'écran et de l'audio système`.

```
xcode-select --install
```

#### Pour Windows :

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

Ou téléchargez [Mingw-w64](https://sourceforge.net/projects/mingw-w64/files) ainsi que d'autres gcc, puis ajoutez un chemin tel que `C:\mingw64\bin` à la variable d'environnement système `Path`.
[Définir les variables d'environnement pour exécuter GCC depuis la ligne de commande](https://www.youtube.com/results?search_query=Set+environment+variables+to+run+GCC+from+command+line).

`Ou les autres GCC` (à l'exception de Mingw-w64, vous devez compiler vous-même la « libpng » lorsque vous utilisez [bitmap](https://github.com/vcaesar/bitmap).)

#### Pour tout le reste :

```
GCC

X11 avec l'extension XTest (la bibliothèque Xtst)

"Clipboard" : xsel xclip

"Bitmap" : libpng (Utilisé uniquement par "bitmap".)

"Event-Gohook" : xcb, xkb, libxkbcommon (Utilisé uniquement par "hook".)
```

##### Ubuntu :

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

##### Fedora :

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

Le backend Wayland est une implémentation **100 % Go (sans Cgo)**, aucune
bibliothèque C système n'est donc requise. Il nécessite un compositeur basé sur
wlroots (Sway, Hyprland, Wayfire, ...) prenant en charge les protocoles suivants :

```
zwlr_virtual_pointer_v1            (contrôle de la souris)
zwp_virtual_keyboard_v1            (contrôle du clavier)
zwlr_screencopy_v1                 (capture d'écran)
zwlr_foreign_toplevel_management_v1 (gestion des fenêtres)
```

GNOME et KDE ne prennent **pas** en charge ces protocoles nativement.

#### libei (GNOME / KDE)

Le backend libei est également une implémentation **100 % Go (sans Cgo)**. Il
pilote les entrées via l'interface RemoteDesktop de `xdg-desktop-portal` de
freedesktop, il fonctionne donc sur GNOME et KDE (contrairement au backend
Wayland wlroots). Il nécessite :

```
xdg-desktop-portal               (le service D-Bus du portail)
xdg-desktop-portal-gnome / -kde  (le backend de portail de votre bureau)
```

Remarque : le backend libei ne gère que les entrées souris et clavier. La capture
d'écran et la gestion des fenêtres renvoient `ErrNotSupported`.

## Installation:

Avec la prise en charge des modules Go (Go 1.11+), il suffit d'importer :

```go
import "github.com/go-vgo/robotgo"
```

Sinon, pour installer le paquet robotgo, exécutez la commande :

```
go get github.com/go-vgo/robotgo
```

png.h: No such file or directory ? Veuillez consulter [issues/47](https://github.com/go-vgo/robotgo/issues/47).

## Update:

```
go get -u github.com/go-vgo/robotgo
```

Notez le problème de cache de compilation des fichiers C de go1.10.x, [golang #24355](https://github.com/golang/go/issues/24355).
Problème de `go mod vendor`, [golang #26366](https://github.com/golang/go/issues/26366).

## Cgo-free Builds:

RobotGo fournit des backends **100 % Go (sans Cgo)** pour Windows, macOS, X11,
Wayland et libei (Linux). Ils exposent la même API `robotgo`, votre code n'a donc pas besoin
d'être modifié — il suffit d'un tag de build. Ces backends se compilent en
cross-compilation avec `CGO_ENABLED=0` (sans GCC, MinGW, Xcode ni en-têtes X11).

| Backend                          | Tag de build | Paquet Go                           |
| -------------------------------- | ------------ | ----------------------------------- |
| Windows (sans Cgo)               | `win`        | `github.com/go-vgo/robotgo/win`     |
| macOS (Quartz via purego)        | `mac`        | `github.com/go-vgo/robotgo/darwin`  |
| X11 (Linux, protocole X pur Go)  | `x11`        | `github.com/go-vgo/robotgo/x11`     |
| Wayland (Linux, wlroots)         | `wayland`    | `github.com/go-vgo/robotgo/wayland` |
| libei (Linux, portail GNOME/KDE) | `libei`      | `github.com/go-vgo/robotgo/libei`   |

```sh
# Windows, sans Cgo / sans MinGW
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -tags win ./...

# macOS, Quartz/CoreGraphics chargé à l'exécution via purego (sans Xcode)
CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -tags mac ./...

# X11, protocole X pur Go (XTEST) — sans en-têtes X11
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -tags x11 ./...

# Wayland, compositeur basé sur wlroots (Sway, Hyprland, Wayfire, ...)
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -tags wayland ./...

# libei, GNOME/KDE via l'interface RemoteDesktop de xdg-desktop-portal
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -tags libei ./...
```

Avec le tag `win`, le backend Cgo/Win32 par défaut est exclu et les appels sont
redirigés vers le paquet Go pur `win` ; avec le tag `mac`, le backend Cgo/Quartz
par défaut est exclu et les appels sont redirigés vers le paquet Go pur `darwin`
(la gestion des fenêtres renvoie `ErrNotSupported`) ; avec le tag `x11`, le
backend Cgo/X11 est exclu et les appels sont redirigés vers le paquet Go pur
`x11` ; avec le tag `wayland`, le backend Cgo/X11 est exclu et les appels sont
redirigés vers le paquet Go pur `wayland` ; avec le tag `libei`, les backends
Cgo/X11 et Wayland wlroots sont tous deux exclus et les appels sont redirigés
vers le paquet Go pur `libei`.

## [Examples:](https://github.com/go-vgo/robotgo/blob/master/examples)

#### [Souris](https://github.com/go-vgo/robotgo/blob/master/examples/mouse/main.go)

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

#### [Clavier](https://github.com/go-vgo/robotgo/blob/master/examples/key/main.go)

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

#### [Écran](https://github.com/go-vgo/robotgo/blob/master/examples/screen/main.go)

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

#### [Événement](https://github.com/robotn/gohook/blob/master/examples/main.go)

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

#### [Fenêtre](https://github.com/go-vgo/robotgo/blob/master/examples/window/main.go)

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

- [L'auteur est Evans](https://github.com/vcaesar)
- [Mainteneurs](https://github.com/orgs/go-vgo/people)

## Plans

- Réécrire une partie du code C en Go (comme x11, windows)
- Meilleure prise en charge du multi-écran
- Prise en charge de Wayland
- Mise à jour du handle de fenêtre
- Essayer de prendre en charge Android et iOS

## Contributors

- Consultez la [page des contributeurs](https://github.com/go-vgo/robotgo/graphs/contributors) pour la liste complète des contributeurs.
- Consultez les [directives de contribution](https://github.com/go-vgo/robotgo/blob/master/CONTRIBUTING.md).

## License

Robotgo est principalement distribué selon les termes de « the Apache License (Version 2.0) », certaines parties étant couvertes par diverses licences de type BSD.

Voir [LICENSE-APACHE](http://www.apache.org/licenses/LICENSE-2.0), [LICENSE](https://github.com/go-vgo/robotgo/blob/master/LICENSE).
