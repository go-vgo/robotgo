# RobotGo

Go native cross-platform desktop automation: mouse, keyboard, screen, bitmap, process, window handle, clipboard, and global event listener. Supports macOS, Windows, Linux; amd64 and arm64.

Module: `github.com/go-vgo/robotgo` — `go.mod` declares `go 1.26.0` (GitHub Actions sets up Go 1.26.x).

## Build/Test/Lint Commands

Default Cgo backend prerequisites: `GCC` must be installed. `CGO_ENABLED=1` (default). On macOS, Xcode Command Line Tools + Accessibility/Screen Recording permissions. On Linux, X11 + XTest (`libx11-dev xorg-dev libxtst-dev`). The pure-Go backends support `CGO_ENABLED=0`; see the backend tags below.

- **Build**: `go build -v .`
- **Build all subpackages**: `go build -v ./...`
- **Fetch deps**: `go get -v -t -d ./...`
- **Test (Cgo CI smoke tests — queries screen/window state)**: `go test -v robot_info_test.go`
- **Test (pure-Go default backend)**: `go test -v -tags purego .`
- **Test (Linux pure-Go X11)**: `go test -v -tags "purego,x11" . ./x11`
- **Test (Linux pure-Go libei)**: `go test -v -tags "purego,libei" . ./libei`
- **Test (full)**: `go test -v ./...` (Linux CI wraps with `xvfb-run` — see `.circleci/config.yml`)
- **Single test**: `go test -v -run TestGetScreenSize .`
- **Format**: `gofmt -w .` (code uses tab indentation, standard `gofmt` style)
- **Vet**: `go vet ./...`
- **Run an example**: `cd examples/mouse && go run main.go`

There is no Makefile / Taskfile / linter config. `.github/workflows/go.yml` runs Go 1.26.x on macOS + Windows (Cgo build/smoke tests and pure-Go root tests) and Ubuntu (pure-Go X11/libei root + backend tests with `CGO_ENABLED=0`). `.circleci/config.yml` runs Linux package tests under xvfb; its container image is still `golang:1.25.0`, below the module's declared Go version. The old `appveyor.yml` has been removed.

## Architecture

Single Go package `robotgo` at repo root (flat layout) with platform-specific files and C-binding subpackages. The default backend is Cgo wrappers over C headers vendored in subdirectories; build tags split platform implementations. Five **pure-Go (no-Cgo) backends** live in their own packages: `win/` (`-tags win`, Win32 via tailscale/win), `darwin/` (`-tags mac`, Quartz via purego), `x11/` (`-tags x11`, X11 protocol), `wayland/` (`-tags wayland`, wlroots virtual-input protocols), and `libei/` (`-tags libei`, xdg-desktop-portal RemoteDesktop for GNOME/KDE). `-tags purego` selects macOS/Windows/Linux defaults (`mac`/`win`/`wayland`); Linux can override with `purego,x11` or `purego,libei`. `robotgo.go` excludes all these tags. Build the module root (`.`) for these backends; some examples/subpackages still require Cgo APIs.

```
robotgo/
├── robotgo.go              # default Cgo API; excludes wayland, win, libei, mac, x11, purego
├── robotgo_pub.go          # portable pkg vars: Version, MouseSleep, KeySleep, DisplayID, Scale...
├── doc.go                  # package doc
├── robotgo_mac.go          # macOS Cgo implementation; excludes mac/purego
├── robotgo_mac_unix.go     # non-Windows Cgo implementation; excludes pure-Go backend tags
├── robotgo_mac_win.go      # macOS/Windows Cgo implementation; excludes win/mac/purego
├── robotgo_win.go          # //go:build windows
├── robotgo_x11.go          # default non-macOS/non-Windows Cgo implementation
├── robotgo_android.go, robotgo_adb.go
├── robotgo_ocr.go          # //go:build ocr (gosseract OCR)
├── libei.go                # //go:build linux && libei — wires libei/ backend into robotgo pkg
├── wayland_n.go, windows_n.go, darwin.go, x11_n.go # root pure-Go backend adapters
├── key.go, keycode.go, screen.go, img.go, ps.go
├── robotgo_fn_v1.go        # deprecated v1 aliases (kept for compat)
├── robot_info_test.go      # default Cgo smoke tests (used by GitHub Actions)
├── robotgo_test.go         # interactive Cgo tests, macOS/Windows only
├── img_test.go, key_test.go, robot_mac_test.go # image, Cgo key and macOS pure-Go tests
├── base/       # C helpers (MMBitmap, rgb, microsleep, types, os, pubs, xdisplay)
├── mouse/      # Go pkg + C (mouse.h, mouse_c.h) with *_darwin.go/_windows.go/_x11.go
├── key/        # Go pkg + C (keycode.h, keycode_c.h, keypress.h, keypress_c.h, key_windows.go)
├── screen/     # Go pkg + C (goScreen.h, screen.go, screen_c.h, screengrab_c.h)
├── window/     # Go pkg + C (goWindow.h, window.h, alert_c.h, win_sys.h, pub.h)
├── clipboard/  # Go pkg (darwin/unix/windows variants) + cmd/gocopy, cmd/gopaste, example/
├── win/        # pure-Go Windows backend (no Cgo); //go:build windows
├── wayland/    # pure-Go wlroots Wayland backend; internal/protocols/ (wlr_*), libei/
├── libei/      # pure-Go libei/xdg-portal backend (GNOME/KDE); //go:build linux
├── mcp/        # MCP server package (mcp.go, currently a stub)
├── event/      # C headers for android/ios global hooks (event_c.h)
├── cv/         # OpenCV helper (gocv.go)
├── examples/   # main.go + mouse, key, screen, window, scale — runnable main.go
├── lang/       # translated READMEs (de, es, fr, ja, ko, pt, ru, zh, zht)
├── skills/     # SKILL.md (agent skill descriptor)
├── x11/, darwin/   # implemented pure-Go X11/macOS backends, with tests
├── docs/       # install.md, keys.md, CHANGELOG.md, README.md, archive/
└── .github/workflows/go.yml, .circleci/config.yml
```

Key subpackage relationships: the root `robotgo` package pulls C code from `screen/goScreen.h`, `mouse/mouse_c.h`, `window/goWindow.h`. The `key/` and `clipboard/` directories are importable Go packages; `base/` is header-only C support. The pure-Go packages expose platform-specific subsets of the robotgo API through root adapters; unsupported operations may return `ErrNotSupported` (for example, macOS pure-Go window management).

## Code Style

- **Copyright header**: every Go and C file starts with the 10-line `Copyright (c) 2016-2026 AtomAI...` block (see `CONTRIBUTING.md`). Preserve it verbatim when editing; add a second header only if authorship changes.
- **Indentation**: tabs (Go default). Run `gofmt` before committing.
- **Build tags**: use both forms together — `//go:build darwin` plus legacy `// +build darwin` — matching existing files.
- **Cgo**: `import "C"` immediately follows a `/* ... */` comment block containing `#cgo` directives and `#include`s. Keep LDFLAGS per-OS (`#cgo darwin LDFLAGS`, `#cgo linux LDFLAGS`, `#cgo windows LDFLAGS`).
- **Naming**: exported `CamelCase`; key constants prefixed `Key` (e.g. `KeyA`, `KeyEnter`); C-type aliases prefixed `C` (`CBitmap`, `CHex`). Mirror existing patterns.
- **Imports**: stdlib first, blank line, then third-party (`github.com/...`). Internal subpackage imported as `github.com/go-vgo/robotgo/clipboard` (full module path, not relative).
- **Comments**: godoc-style `// FuncName does X` on every exported symbol. Package doc lives in `doc.go` / top of `robotgo.go`.
- **Error handling**: return `error` as last result; use `errors.New` / `fmt.Errorf`. `Try(fn, handler)` helper wraps panics via `recover`. Do not swallow errors.
- **Types**: existing code uses `interface{}` (pre-generics) — match local style when editing that file, but prefer `any` in new code. Do not mass-rewrite; gopls emits hints, not errors.

## Testing

- Framework: stdlib `testing` + `github.com/vcaesar/tt` (`tt.Expect(t, want, got)`).
- Test files: `*_test.go` beside sources. Package declared as `robotgo_test` (external) for API-surface tests, or `robotgo` for internal.
- **Default Cgo smoke tests** live in `robot_info_test.go`, explicitly selected by GitHub Actions on macOS/Windows. They query screen/window state and are excluded from pure-Go package test runs. Backend-independent tests such as `img_test.go` run under the pure-Go root test jobs; match build tags to the APIs being tested.
- **Interactive / display-required tests** live in `robotgo_test.go`; its build tags restrict it to macOS/Windows Cgo, so CircleCI's Linux `xvfb-run go test -v ./...` does not include that file. Pure-Go backend tests live beside their implementations; macOS root adapter tests are in `robot_mac_test.go`.
- Run one test: `go test -v -run TestGetScreenSize .`
- No fixtures, snapshots, or golden files in use. Screenshots produced by examples are `.gitignore`d.

## Key Patterns

- **Platform and backend build tags are mandatory**. Keep Cgo implementations separate from pure-Go adapters. Any new cross-platform API needs appropriate implementations (or explicit unsupported stubs) for darwin, linux, windows — examine `mouse/mouse_darwin.go`, `mouse_windows.go`, `mouse_x11.go` for Cgo and the root backend adapters for pure-Go conventions.
- **Free C-allocated bitmaps**: every `CaptureScreen`, `ToCBitmap`, etc. must be paired with `defer robotgo.FreeBitmap(bit)` or `robotgo.FreeBitmapArr(...)`. Leaking is a memory bug on all platforms.
- **Global tunables** are package-level vars, not config structs: `MouseSleep`, `KeySleep`, `DisplayID`, `NotPid`, `Scale`. Callers mutate them directly (see README examples). Do not hide them behind getters.
- **`robotgo_fn_v1.go`** contains deprecated v1 aliases — do not add new APIs there, but do not delete existing ones (backwards compatibility).
- **Version string** lives in `robotgo_pub.go` as `const Version = "v2.00.0.1658, MT. Baker!"` (moved out of `robotgo.go`). Bump it when releasing; `TestGetVer` asserts it matches `GetVersion()`. The pure-Go backends carry their own versions in each package's `robotgo.go`: `v0.1.0-windows`, `v0.1.0-darwin`, `v0.1.0-x11`, `v0.1.0-wayland`, `v0.1.0-libei`.
- **Windows pid vs hwnd**: set `robotgo.NotPid = true` to pass window handles instead of pids into the window/key APIs on Windows.
- **macOS permissions**: most screen/input APIs silently fail without Accessibility + Screen Recording grants. When reproducing bugs on darwin, verify System Settings → Privacy & Security first.
- **Do not vendor**: `vendor/` is in `.gitignore`; also avoid `go mod vendor` (upstream note in README references golang/go#26366).
- **C artifacts** (`*.cgo1.go`, `*.cgo2.c`, `_cgo_*`, `*.o`, `*.a` except whitelisted `cdeps/...` libpng archives) are git-ignored — do not commit them.
- **Commit sign-off** is expected (see `CONTRIBUTING.md`); PRs require ≥2 maintainer review (LGTM).

## Dependencies

- `github.com/jezek/xgb`, `github.com/jezek/xgbutil` — X11 protocol on Linux.
- `github.com/vcaesar/go-wayland` — Wayland protocol client (used by the `wayland/` pure-Go backend).
- `github.com/godbus/dbus/v5` — D-Bus, used on Linux for Wayland/desktop ops and the `libei/` xdg-portal backend.
- `github.com/tailscale/win`, `github.com/dblohm7/wingoes`, `github.com/yusufpapurcu/wmi`, `github.com/go-ole/go-ole`, `golang.org/x/sys` — Windows system APIs (also used by the `win/` pure-Go backend).
- `github.com/ebitengine/purego` — direct dependency for runtime native-library loading (including the macOS backend); `github.com/gen2brain/shm` — indirect shared-memory helper.
- `github.com/vcaesar/keycode` — cross-platform keycode mapping (used by `key/`).
- `github.com/vcaesar/imgo`, `golang.org/x/image` — image encode/decode (PNG/JPEG save).
- `github.com/vcaesar/screenshot` — screenshot backend.
- `github.com/vcaesar/gops`, `github.com/shirou/gopsutil/v4` — process enumeration (`FindIds`, `PidExists`, `Kill`).
- `github.com/vcaesar/tt` — testing assertions.
- `github.com/otiai10/gosseract/v2` — OCR (used by `robotgo_ocr.go`; needs `libtesseract`).
- `github.com/godbus/dbus/v5` — used on Linux for Wayland/desktop ops.
- Companion repos (not in `go.mod`, referenced in README/examples): `github.com/vcaesar/bitmap`, `github.com/vcaesar/gcv` (OpenCV), `github.com/jezek/gohook` (global event hook).
