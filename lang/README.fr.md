# Robotgo

<p align="center">
  <img src="../docs/robotgo-logo.svg" width="560" alt="Logo RobotGo — un robot avec un pointeur de souris" />
</p>

[![Build Status](https://github.com/go-vgo/robotgo/workflows/Go/badge.svg)](https://github.com/go-vgo/robotgo/commits/master)
[![CircleCI Status](https://circleci.com/gh/go-vgo/robotgo.svg?style=shield)](https://circleci.com/gh/go-vgo/robotgo)
[![golangci-lint](https://github.com/go-vgo/robotgo/actions/workflows/lint.yml/badge.svg)](https://github.com/go-vgo/robotgo/actions/workflows/lint.yml)
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
- [Builds sans Cgo](#cgo-free-builds)
- [Installation](#installation)
- [Mise à jour](#update)
- [Exemples](#examples)
- [Conversion de types et touches](https://github.com/go-vgo/robotgo/blob/master/docs/keys.md)
- [Compilation croisée](https://github.com/go-vgo/robotgo/blob/master/docs/install.md#crosscompiling)
- [Auteurs](#authors)
- [Projets](#plans)
- [Licence](#license)

## Docs

- [GoDoc](https://godoc.org/github.com/go-vgo/robotgo) <br>
- [Documentation de l'API](https://github.com/go-vgo/robotgo/blob/master/docs/doc.md) (obsolète, plus mise à jour)

## Binding

[ADB](https://github.com/vcaesar/adb), encapsulation de l'API adb d'Android.

## Requirements

**Go 1.26+** (voir [go.mod](../go.mod)) et un backend :

- **Par défaut (Cgo) :** `CGO_ENABLED=1`, un compilateur C et les bibliothèques
  système ci-dessous.
- **Go pur (expérimental) :** aucune chaîne d'outils C ; voir [Builds sans Cgo](#cgo-free-builds).

Les deux nécessitent les [dépendances par fonctionnalité](#feature-specific-dependencies)
correspondant aux fonctionnalités utilisées.

### Default Cgo setup

#### macOS

Go et les Xcode Command Line Tools (Clang) :

```sh
brew install go
xcode-select --install
```

**Autorisations (Cgo et Go pur) :** dans **Réglages Système > Confidentialité et
sécurité**, accordez à l'application ou au terminal l'**Accessibilité** (saisie)
et l'**Enregistrement de l'écran** (capture).

#### Windows

```powershell
winget install GoLang.Go
```

Plus **une** chaîne d'outils C. [LLVM-MinGW](https://github.com/mstorsjo/llvm-mingw) :

```powershell
winget install MartinStorsjo.LLVM-MinGW.UCRT
$env:CC = "clang"
```

Ou [MinGW-w64 (WinLibs)](https://winlibs.com/) :

```powershell
winget install BrechtSanders.WinLibs.POSIX.UCRT
$env:CC = "gcc"
```

Le répertoire `bin` de la chaîne d'outils (par exemple `C:\mingw64\bin`) doit
figurer dans `PATH` et correspondre à `GOARCH`. Tout autre compilateur compatible
Cgo, comme une installation manuelle de
[MinGW-w64](https://sourceforge.net/projects/mingw-w64/files), fonctionne aussi.

[Bitmap](https://github.com/vcaesar/bitmap) ne fournit libpng que pour MinGW-w64 ;
avec une autre chaîne d'outils, compilez libpng vous-même.

#### Linux (X11)

Build : GCC, en-têtes libc, en-têtes X11/XTest. Exécution : un serveur X avec
XTEST et `DISPLAY` défini. Pour Wayland natif, voir [Builds sans Cgo](#cgo-free-builds).

**Ubuntu / Debian :**

```sh
# Go (Snap si le paquet de la distribution est trop ancien)
sudo snap install go --classic
# ou : sudo apt install golang

# Base (Cgo) : compilateur, libc, X11, XTest
sudo apt install gcc libc6-dev libx11-dev xorg-dev libxtst-dev

# Presse-papiers (X11)
sudo apt install xsel xclip

# Bitmap : libpng
sudo apt install libpng++-dev

# GoHook : XCB, XKB
sudo apt install xcb libxcb-xkb-dev x11-xkb-utils libx11-xcb-dev libxkbcommon-x11-dev libxkbcommon-dev
```

**Fedora :**

```sh
# Base (Cgo) : compilateur, libc, X11, XTest
sudo dnf install gcc glibc-devel libX11-devel libXtst-devel

# Presse-papiers (X11)
sudo dnf install xsel xclip

# Bitmap : libpng
sudo dnf install libpng-devel

# GoHook : XKB (Fedora < 34 : xorg-x11-xkb-utils-devel au lieu de xkbcomp-devel)
sudo dnf install libxkbcommon-devel libxkbcommon-x11-devel xkbcomp-devel
```

### Feature-specific dependencies

À ajouter à l'installation de base, uniquement pour les fonctionnalités utilisées :

- **Presse-papiers Linux (Cgo et Go pur) :** `xclip` ou `xsel` sous X11 ;
  `wl-clipboard` sous Wayland :

  ```sh
  sudo apt install wl-clipboard   # Ubuntu / Debian
  sudo dnf install wl-clipboard   # Fedora
  ```

- **[Bitmap](https://github.com/vcaesar/bitmap) :** en-têtes libpng.
- **[GoHook](https://github.com/robotn/gohook) :** en-têtes XCB/XKB sous Linux.

Bitmap et GoHook sont des paquets Cgo distincts ; un backend RobotGo en Go pur ne
supprime pas leurs dépendances C.

## Cgo-free Builds

Les **backends expérimentaux en Go pur** se compilent et se cross-compilent avec
`CGO_ENABLED=0` (sans GCC, MinGW, Xcode ni en-têtes X11). Même chemin d'import et
même API ; un **tag de build** sélectionne le backend — `CGO_ENABLED=0` seul ne
sélectionne rien.

| `GOOS` cible  | Tag de build | Paquet                 | Prérequis d'exécution                                                   |
| ------------- | ------------ | ---------------------- | ----------------------------------------------------------------------- |
| `windows`     | `win`        | [win](../win/)         | Win32 uniquement, aucun runtime supplémentaire                          |
| `darwin`      | `mac`        | [darwin](../darwin/)   | frameworks système et [autorisations de confidentialité](#macos)        |
| `linux`       | `x11`        | [x11](../x11/)         | serveur X avec XTEST, `DISPLAY` défini                                  |
| `linux`       | `wayland`    | [wayland](../wayland/) | compositeur wlroots, voir [Wayland](#wayland)                           |
| `linux`       | `libei`      | [libei](../libei/)     | portail RemoteDesktop (GNOME/KDE), voir [libei](#libei-gnome--kde)      |

**`purego`** est un raccourci : `mac` sur macOS, `win` sur Windows, `wayland` sur
Linux. Ajoutez `x11` ou `libei` pour le remplacer sous Linux (`-tags "purego,x11"`).

Utilisez **un seul tag de backend** correspondant à `GOOS` ; les combinaisons
telles que `x11,libei` ne compilent pas (`purego` plus un remplacement Linux est
la seule exception).

### Build commands

```sh
# OS courant
CGO_ENABLED=0 go build -tags purego .

# Compilation croisée
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -tags win .
CGO_ENABLED=0 GOOS=darwin  GOARCH=arm64 go build -tags mac .
CGO_ENABLED=0 GOOS=linux   GOARCH=amd64 go build -tags x11 .
CGO_ENABLED=0 GOOS=linux   GOARCH=amd64 go build -tags wayland .
CGO_ENABLED=0 GOOS=linux   GOARCH=amd64 go build -tags libei .
```

PowerShell : `$env:CGO_ENABLED = "0"; go build -tags purego .`

Compilez la racine du module (`.`) ou votre propre paquet, pas `./...` :
`examples/` et certains sous-paquets ont besoin d'API réservées à Cgo.

### Limitations

Les prérequis d'exécution ci-dessus s'appliquent toujours. Les API réservées à
Cgo comme `CaptureScreen` / `FreeBitmap` n'existent pas ; utilisez `CaptureImg`.
Les opérations non prises en charge renvoient `ErrNotSupported` (ou un résultat
vide).

- **`mac` :** pas de gestion des fenêtres — `ActiveName` renvoie
  `ErrNotSupported`, `GetTitle` renvoie `""`, et réduire/agrandir/fermer n'ont
  aucun effet.
- **`wayland`, `libei` :** `Location()` renvoie la dernière position injectée par
  RobotGo, pas le curseur physique. Le presse-papiers nécessite `wl-clipboard`.

#### Wayland

Un compositeur wlroots (Sway, Hyprland, Wayfire, ...) avec `WAYLAND_DISPLAY` et
`XDG_RUNTIME_DIR` définis, exposant :

| Fonctionnalité        | Global du protocole                |
| --------------------- | ---------------------------------- |
| Contrôle de la souris | `zwlr_virtual_pointer_manager_v1`  |
| Contrôle du clavier   | `zwp_virtual_keyboard_manager_v1`  |
| Capture d'écran       | `zwlr_screencopy_manager_v1`       |
| Gestion des fenêtres  | `zwlr_foreign_toplevel_manager_v1` |

GNOME et KDE ne fournissent pas ces protocoles ; utilisez `libei` dans ce cas.

#### libei (GNOME / KDE)

La saisie passe par l'interface D-Bus **RemoteDesktop** de
`xdg-desktop-portal` : il faut un bus de session, `xdg-desktop-portal` et un
backend qui l'implémente (`xdg-desktop-portal-gnome`, `-kde` ou `-wlr`). Au
premier lancement, une boîte de dialogue de consentement s'affiche ; le jeton de
restauration est mis en cache dans `$XDG_STATE_HOME/robotgo`.

Les déplacements absolus utilisent un flux ScreenCast lié (par défaut) ;
définissez `libei.LinkScreenCast = false` pour n'avoir que des déplacements
relatifs. La capture d'écran et la gestion des fenêtres renvoient
`ErrNotSupported`.

## Installation

```sh
go get github.com/go-vgo/robotgo
```

```go
import "github.com/go-vgo/robotgo"
```

`png.h: No such file or directory` avec Bitmap : consultez ses prérequis libpng
et [issues/47](https://github.com/go-vgo/robotgo/issues/47).

## Update

```sh
go get -u github.com/go-vgo/robotgo
```

Notez le problème de cache de compilation des fichiers C de go1.10.x, [golang #24355](https://github.com/golang/go/issues/24355).
Problème de `go mod vendor`, [golang #26366](https://github.com/golang/go/issues/26366).

## [Examples](https://github.com/go-vgo/robotgo/blob/master/examples)

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

- Meilleure prise en charge du multi-écran
- Mise à jour du handle de fenêtre
- Essayer de prendre en charge Android et iOS

## Contributors

- Consultez la [page des contributeurs](https://github.com/go-vgo/robotgo/graphs/contributors) pour la liste complète des contributeurs.
- Consultez les [directives de contribution](https://github.com/go-vgo/robotgo/blob/master/CONTRIBUTING.md).

## License

Robotgo est principalement distribué selon les termes de « the Apache License (Version 2.0) », certaines parties étant couvertes par diverses licences de type BSD.

Voir [LICENSE-APACHE](http://www.apache.org/licenses/LICENSE-2.0), [LICENSE](https://github.com/go-vgo/robotgo/blob/master/LICENSE).
