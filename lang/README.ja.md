# Robotgo

[![Build Status](https://github.com/go-vgo/robotgo/workflows/Go/badge.svg)](https://github.com/go-vgo/robotgo/commits/master)
[![CircleCI Status](https://circleci.com/gh/go-vgo/robotgo.svg?style=shield)](https://circleci.com/gh/go-vgo/robotgo)
[![Go Report Card](https://goreportcard.com/badge/github.com/go-vgo/robotgo)](https://goreportcard.com/report/github.com/go-vgo/robotgo)
[![GoDoc](https://pkg.go.dev/badge/github.com/go-vgo/robotgo?status.svg)](https://pkg.go.dev/github.com/go-vgo/robotgo?tab=doc)
[![GitHub release](https://img.shields.io/github/release/go-vgo/robotgo.svg)](https://github.com/go-vgo/robotgo/releases/latest)
<a href="https://discord.gg/npPb3NzE4A"><img src="https://img.shields.io/discord/1484658282777018551.svg?logo=discord&logoColor=white&label=Discord&color=5865F2" alt="Join the Discord chat at https://discord.gg/npPb3NzE4A"></a>

[English](../README.md) | [简体中文](README.zh.md) | [繁體中文](README.zht.md) | 日本語 | [한국어](README.ko.md) | [Français](README.fr.md) | [Deutsch](README.de.md) | [Español](README.es.md) | [Русский](README.ru.md) | [Português](README.pt.md)

> Golang によるデスクトップ自動化、自動テスト、そして AI コンピュータ操作（Computer Use）。<br>
> マウスやキーボードの制御、画面の読み取り、プロセス、ウィンドウハンドル、画像とビットマップ、そしてグローバルイベントの監視を行えます。

RobotGo は Mac、Windows、Linux に対応しており、arm64 および x86-amd64 アーキテクチャをサポートしています。

現在 [Codg](https://github.com/vcaesar/codg) を開発しています。シンプルで使いやすい AI エージェント（Agent）作業システムで、自動化・非同期・並行処理・高効率・高精度を実現します。

<p align="center">
<a href="https://github.com/vcaesar/codg" rel="nofollow">
<img width="800" alt="Codg Demo" src="https://github.com/vcaesar/codg/raw/main/demo/26-04-1.png" />
</a>
</p>

[RobotGo-Pro](https://github.com/vcaesar/robotgo-pro) では、JavaScript、Python、Lua などの他言語版、テクニカルサポート、新機能、そして最新の robotgo バージョン（「現在オープンソース版はありません」）を入手できます。

## 目次

- [ドキュメント](#docs)
- [バインディング](#binding)
- [動作環境](#requirements)
- [インストール](#installation)
- [アップデート](#update)
- [Cgo 不要ビルド](#cgo-free-builds)
- [サンプル](#examples)
- [型変換とキー](https://github.com/go-vgo/robotgo/blob/master/docs/keys.md)
- [クロスコンパイル](https://github.com/go-vgo/robotgo/blob/master/docs/install.md#crosscompiling)
- [作者](#authors)
- [計画](#plans)
- [ライセンス](#license)

## Docs

- [GoDoc](https://godoc.org/github.com/go-vgo/robotgo) <br>
- [API ドキュメント](https://github.com/go-vgo/robotgo/blob/master/docs/doc.md)（非推奨、更新されていません）

## Binding:

[ADB](https://github.com/vcaesar/adb)、Android の adb API をラップしたものです。

## Requirements:

RobotGo をインストールする前に、まず `Golang、GCC` が正しくインストールされていることを確認してください。

### すべてのプラットフォーム：

```
Golang

GCC
```

#### MacOS：

```
brew install go
```

Xcode コマンドラインツール；<br>
さらにプライバシー設定で、以下の場所に「画面収録」と「アクセシビリティ」の権限を追加してください：<br>
`システム設定 > プライバシーとセキュリティ > アクセシビリティ、画面とシステムオーディオの収録`。

```
xcode-select --install
```

#### Windows：

```
winget install Golang.go
```

[llvm-mingw](https://github.com/mstorsjo/llvm-mingw)

```
winget install MartinStorsjo.LLVM-MinGW.UCRT
```

または [Mingw-w64](https://sourceforge.net/projects/mingw-w64/files)

```
winget install BrechtSanders.WinLibs.POSIX.UCRT
```

または [Mingw-w64](https://sourceforge.net/projects/mingw-w64/files) やその他の gcc をダウンロードし、`C:\mingw64\bin` のようなパスをシステム環境変数 `Path` に設定してください。
[コマンドラインから GCC を実行するための環境変数を設定する](https://www.youtube.com/results?search_query=Set+environment+variables+to+run+GCC+from+command+line)。

`または、その他の GCC`（Mingw-w64 を除き、[bitmap](https://github.com/vcaesar/bitmap) を使用する場合は「libpng」を自分でコンパイルする必要があります。）

#### その他すべてのプラットフォーム：

```
GCC

X11 拡張の XTest（Xtst ライブラリ）

"Clipboard": xsel xclip

"Bitmap": libpng（"bitmap" でのみ使用。）

"Event-Gohook": xcb, xkb, libxkbcommon（"hook" でのみ使用。）
```

##### Ubuntu：

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

##### Fedora：

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

Wayland バックエンドは **純粋 Go（Cgo 不要）** 実装なので、システムの C ライブラリは
一切必要ありません。以下のプロトコルをサポートする wlroots ベースのコンポジッタ
（Sway、Hyprland、Wayfire など）が必要です：

```
zwlr_virtual_pointer_v1            (マウス制御)
zwp_virtual_keyboard_v1            (キーボード制御)
zwlr_screencopy_v1                 (画面キャプチャ)
zwlr_foreign_toplevel_management_v1 (ウィンドウ管理)
```

GNOME と KDE はこれらのプロトコルをネイティブにはサポートしていません。

#### libei（GNOME / KDE）

libei バックエンドも **純粋 Go（Cgo 不要）** 実装です。freedesktop の
`xdg-desktop-portal` RemoteDesktop インターフェースを介して入力を駆動するため、
wlroots Wayland バックエンドとは異なり GNOME と KDE で動作します。以下が必要です：

```
xdg-desktop-portal               (portal D-Bus サービス)
xdg-desktop-portal-gnome / -kde  (お使いのデスクトップの portal バックエンド)
```

注意：libei バックエンドはマウスとキーボード入力のみを処理します。画面キャプチャと
ウィンドウ管理は `ErrNotSupported` を返します。

## Installation:

Go module に対応している場合（Go 1.11 以降）、import するだけです：

```go
import "github.com/go-vgo/robotgo"
```

そうでない場合は、次のコマンドを実行して robotgo パッケージをインストールしてください：

```
go get github.com/go-vgo/robotgo
```

png.h: No such file or directory？[issues/47](https://github.com/go-vgo/robotgo/issues/47) を参照してください。

## Update:

```
go get -u github.com/go-vgo/robotgo
```

go1.10.x の C ファイルコンパイルキャッシュ問題に注意してください。[golang #24355](https://github.com/golang/go/issues/24355)。
`go mod vendor` の問題、[golang #26366](https://github.com/golang/go/issues/26366)。

## Cgo-free Builds:

RobotGo は Windows、macOS、X11、Wayland、libei（Linux）向けに **純粋 Go（Cgo 不要）** バックエンドを提供します。現在は実験的な機能です。
これらは同じ `robotgo` API を公開するため、コードの変更は不要で、ビルドタグを指定するだけです。
これらのバックエンドは `CGO_ENABLED=0` でクロスコンパイルできます（GCC、MinGW、Xcode、X11 ヘッダー不要）。

| バックエンド                     | ビルドタグ | Go パッケージ                       |
| -------------------------------- | ---------- | ----------------------------------- |
| Windows（Cgo 不要）              | `win`      | `github.com/go-vgo/robotgo/win`     |
| macOS（purego 経由の Quartz）    | `mac`      | `github.com/go-vgo/robotgo/darwin`  |
| X11（Linux、純粋 Go X プロトコル）| `x11`      | `github.com/go-vgo/robotgo/x11`     |
| Wayland（Linux、wlroots）        | `wayland`  | `github.com/go-vgo/robotgo/wayland` |
| libei（Linux、GNOME/KDE portal） | `libei`    | `github.com/go-vgo/robotgo/libei`   |
| 純粋 Go デフォルト（全プラットフォーム） | `purego` | 上記の `mac`/`win`/`wayland` を選択 |

```sh
# プラットフォームごとの純粋 Go デフォルトバックエンド。すべてのターゲットに 1 つのタグで対応：
# macOS -> mac、Windows -> win、Linux -> wayland（x11/libei と組み合わせて上書き可能）
go build -tags purego ./...
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -tags "purego,x11" ./...

# Windows、Cgo / MinGW 不要
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -tags win ./...

# macOS、purego 経由で実行時に Quartz/CoreGraphics をロード（Xcode 不要）
CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -tags mac ./...

# X11、純粋 Go の X プロトコル（XTEST）—— X11 ヘッダー不要
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -tags x11 ./...

# Wayland、wlroots ベースのコンポジッタ（Sway、Hyprland、Wayfire など）
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -tags wayland ./...

# libei、xdg-desktop-portal RemoteDesktop 経由で GNOME/KDE に対応
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -tags libei ./...
```

`win` タグではデフォルトの Cgo/Win32 バックエンドが除外され、呼び出しは純粋 Go の `win` パッケージに転送されます。`mac` タグではデフォルトの Cgo/Quartz バックエンドが除外され、呼び出しは純粋 Go の `darwin` パッケージに転送されます（ウィンドウ管理は `ErrNotSupported` を返します）。`x11` タグでは Cgo/X11 バックエンドが除外され、呼び出しは純粋 Go の `x11` パッケージに転送されます。`wayland` タグでは Cgo/X11 バックエンドが除外され、呼び出しは純粋 Go の `wayland` パッケージに転送されます。`libei` タグでは Cgo/X11 と wlroots Wayland バックエンドの両方が除外され、呼び出しは純粋 Go の `libei` パッケージに転送されます。

`purego` タグはクロスプラットフォームのショートカットです。すべてのプラットフォームで Cgo バックエンドを除外し、対象 OS の純粋 Go デフォルトバックエンド（macOS は `mac`、Windows は `win`、Linux は `wayland`）を選択します。Linux では `x11` または `libei` と組み合わせて（例：`-tags "purego,libei"`）、別の純粋 Go バックエンドを選択できます。

## [Examples:](https://github.com/go-vgo/robotgo/blob/master/examples)

#### [マウス](https://github.com/go-vgo/robotgo/blob/master/examples/mouse/main.go)

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

#### [キーボード](https://github.com/go-vgo/robotgo/blob/master/examples/key/main.go)

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

#### [スクリーン](https://github.com/go-vgo/robotgo/blob/master/examples/screen/main.go)

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

#### [ビットマップ](https://github.com/vcaesar/bitmap/blob/main/examples/main.go)

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

#### [イベント](https://github.com/robotn/gohook/blob/master/examples/main.go)

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

#### [ウィンドウ](https://github.com/go-vgo/robotgo/blob/master/examples/window/main.go)

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

- [作者は Evans です](https://github.com/vcaesar)
- [メンテナー](https://github.com/orgs/go-vgo/people)

## Plans

- 一部の C コードを Go にリファクタリング（x11、windows など）
- マルチスクリーン対応の改善
- Wayland 対応
- ウィンドウハンドルの更新
- Android と iOS への対応を試みる

## Contributors

- 貢献者の完全な一覧は[コントリビューターページ](https://github.com/go-vgo/robotgo/graphs/contributors)をご覧ください。
- [コントリビューションガイドライン](https://github.com/go-vgo/robotgo/blob/master/CONTRIBUTING.md)を参照してください。

## License

Robotgo は主に「Apache License (Version 2.0)」の条件の下で配布されており、一部は各種の BSD 系ライセンスの対象となっています。

詳しくは [LICENSE-APACHE](http://www.apache.org/licenses/LICENSE-2.0)、[LICENSE](https://github.com/go-vgo/robotgo/blob/master/LICENSE) を参照してください。
