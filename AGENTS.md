# RobotGo

Go native cross-platform desktop automation: mouse, keyboard, screen, bitmap, process, window handle, clipboard, and global event listener. Supports macOS, Windows, Linux; amd64 and arm64.

Module: `github.com/go-vgo/robotgo` — `go.mod` declares `go 1.25.0` (CI sets up Go 1.26.x).

## Build/Test/Lint Commands

Prerequisites: `GCC` must be installed. `CGO_ENABLED=1` (default). On macOS, Xcode Command Line Tools + Accessibility/Screen Recording permissions. On Linux, X11 + XTest (`libx11-dev xorg-dev libxtst-dev`).

- **Build**: `go build -v .`
- **Build all subpackages**: `go build -v ./...`
- **Fetch deps**: `go get -v -t -d ./...`
- **Test (CI minimal — no display required)**: `go test -v robot_info_test.go`
- **Test (full)**: `go test -v ./...` (Linux CI wraps with `xvfb-run` — see `.circleci/config.yml`)
- **Single test**: `go test -v -run TestGetScreenSize .`
- **Format**: `gofmt -w .` (code uses tab indentation, standard `gofmt` style)
- **Vet**: `go vet ./...`
- **Run an example**: `cd examples/mouse && go run main.go`

There is no Makefile / Taskfile / linter config. CI is `.github/workflows/go.yml` (macOS + Windows, Go 1.26.x: `go build -v .` then `go test -v robot_info_test.go` only) and `.circleci/config.yml` (Linux full tests under xvfb). The old `appveyor.yml` has been removed.

## Architecture

Single Go package `robotgo` at repo root (flat layout) with platform-specific files and C-binding subpackages. The default backend is Cgo wrappers over C headers vendored in subdirectories; build tags split platform implementations. In addition, three **pure-Go (no-Cgo) backends** now live in their own packages and are selected via build tags on the root package: `win/` (`-tags win`, Win32 via tailscale/win), `wayland/` (`-tags wayland`, wlroots virtual-input protocols), and `libei/` (`-tags libei`, xdg-desktop-portal RemoteDesktop for GNOME/KDE). `robotgo.go` carries `//go:build !wayland && !win && !libei` so exactly one backend compiles.

```
robotgo/
├── robotgo.go              # default Cgo API + preamble; //go:build !wayland && !win && !libei
├── robotgo_pub.go          # portable pkg vars: Version, MouseSleep, KeySleep, DisplayID, Scale...
├── doc.go                  # package doc
├── robotgo_mac.go          # //go:build darwin
├── robotgo_mac_unix.go     # //go:build darwin || linux
├── robotgo_mac_win.go      # //go:build darwin || windows
├── robotgo_win.go          # //go:build windows
├── robotgo_x11.go          # //go:build linux
├── robotgo_android.go, robotgo_adb.go
├── robotgo_ocr.go          # //go:build ocr (gosseract OCR)
├── libei.go                # //go:build linux && libei — wires libei/ backend into robotgo pkg
├── wayland_n.go, windows_n.go
├── key.go, keycode.go, screen.go, img.go, ps.go
├── robotgo_fn_v1.go        # deprecated v1 aliases (kept for compat)
├── robot_info_test.go      # only portable test (used by GitHub Actions)
├── robotgo_test.go         # full interactive tests
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
├── x11/, darwin/   # currently empty placeholder dirs
├── docs/       # install.md, keys.md, CHANGELOG.md, README.md, archive/
└── .github/workflows/go.yml, .circleci/config.yml
```

Key subpackage relationships: the root `robotgo` package pulls C code from `screen/goScreen.h`, `mouse/mouse_c.h`, `window/goWindow.h`. The `key/` and `clipboard/` directories are importable Go packages; `base/` is header-only C support. The pure-Go `win/`, `wayland/` and `libei/` packages each mirror the robotgo API surface (mouse/keyboard/screen/window/process) so a backend can be swapped per platform with a build tag.

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
- **Portable tests** live in `robot_info_test.go` — the only file exercised by GitHub Actions. Keep new tests here if they must run headless on macOS/Windows CI.
- **Interactive / display-required tests** go in `robotgo_test.go` and run only under CircleCI's `xvfb-run go test -v ./...`.
- Run one test: `go test -v -run TestGetScreenSize .`
- No fixtures, snapshots, or golden files in use. Screenshots produced by examples are `.gitignore`d.

## Key Patterns

- **Cgo + platform split is mandatory**. Any new OS-specific function must be gated by `//go:build` tags and have implementations (even stub) for darwin, linux, windows — examine `mouse/mouse_darwin.go`, `mouse_windows.go`, `mouse_x11.go` as the template.
- **Free C-allocated bitmaps**: every `CaptureScreen`, `ToCBitmap`, etc. must be paired with `defer robotgo.FreeBitmap(bit)` or `robotgo.FreeBitmapArr(...)`. Leaking is a memory bug on all platforms.
- **Global tunables** are package-level vars, not config structs: `MouseSleep`, `KeySleep`, `DisplayID`, `NotPid`, `Scale`. Callers mutate them directly (see README examples). Do not hide them behind getters.
- **`robotgo_fn_v1.go`** contains deprecated v1 aliases — do not add new APIs there, but do not delete existing ones (backwards compatibility).
- **Version string** lives in `robotgo_pub.go` as `const Version = "v2.00.0.1658, MT. Baker!"` (moved out of `robotgo.go`). Bump it when releasing; `TestGetVer` asserts it matches `GetVersion()`. The pure-Go backends carry their own versions (`win/robotgo.go` `v0.1.0-windows`, `wayland/robotgo.go` `v0.1.0-wayland`, `libei/robotgo.go` `v0.1.0-libei`).
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
- `github.com/ebitengine/purego`, `github.com/gen2brain/shm` — indirect, dlopen/shared-memory helpers pulled in by the screenshot/wayland paths.
- `github.com/vcaesar/keycode` — cross-platform keycode mapping (used by `key/`).
- `github.com/vcaesar/imgo`, `golang.org/x/image` — image encode/decode (PNG/JPEG save).
- `github.com/vcaesar/screenshot` — screenshot backend.
- `github.com/vcaesar/gops`, `github.com/shirou/gopsutil/v4` — process enumeration (`FindIds`, `PidExists`, `Kill`).
- `github.com/vcaesar/tt` — testing assertions.
- `github.com/otiai10/gosseract/v2` — OCR (used by `robotgo_ocr.go`; needs `libtesseract`).
- `github.com/godbus/dbus/v5` — used on Linux for Wayland/desktop ops.
- Companion repos (not in `go.mod`, referenced in README/examples): `github.com/vcaesar/bitmap`, `github.com/vcaesar/gcv` (OpenCV), `github.com/jezek/gohook` (global event hook).
