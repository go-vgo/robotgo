# RobotGo

Go native cross-platform desktop automation: mouse, keyboard, screen, bitmap, process, window handle, clipboard, and global event listener. Supports macOS, Windows, Linux; amd64 and arm64.

Module: `github.com/go-vgo/robotgo` — `go.mod` declares `go 1.26.0` (GitHub Actions sets up Go 1.26.x).

## Build/Test/Lint Commands

Prerequisites (default Cgo backend): `GCC` must be installed. `CGO_ENABLED=1` (default). On macOS, Xcode Command Line Tools + Accessibility/Screen Recording permissions. On Linux, X11 + XTest (`libx11-dev xorg-dev libxtst-dev`). The pure-Go backends (see Architecture) need none of these and build with `CGO_ENABLED=0`.

- **Build**: `go build -v .`
- **Build all subpackages**: `go build -v ./...`
- **Fetch deps**: `go get -v -t -d ./...`
- **Test (Cgo CI smoke tests — query screen/window state, need a display session)**: `go test -v robot_info_test.go`
- **Test (full)**: `go test -v ./...` (Linux CI wraps with `xvfb-run` — see `.circleci/config.yml`)
- **Single test**: `go test -v -run TestGetScreenSize .`
- **Format**: `gofmt -w .` (code uses tab indentation, standard `gofmt` style)
- **Vet**: `go vet ./...`
- **Run an example**: `cd examples/mouse && go run main.go`
- **Pure-Go backend test**: `go test -v -tags purego .` (picks `mac`/`win`/`wayland` per OS)
- **Pure-Go Linux tests**: `CGO_ENABLED=0 go test -v -tags "purego,x11" . ./x11` and `CGO_ENABLED=0 go test -v -tags "purego,libei" . ./libei`
- **Pure-Go cross build**: `CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -tags mac .`; `GOOS` must match the backend — `GOOS=windows` for `win`, `GOOS=linux` for `x11`/`wayland`/`libei` (a mismatched `GOOS` leaves the root package with no backend). Build the module root `.`, not `./...` — `examples/` and some subpackages need the Cgo backend.

There is no Makefile / Taskfile / linter config. CI:
- `.github/workflows/go.yml` (Go 1.26.x) — job `cgo` (macOS/Windows/Linux): `go build -v ./...`, `go vet .`, `go test -v robot_info_test.go` + `go test -v .` (Linux: apt X11/xvfb deps, `xvfb-run go test -v ./...`); job `purego` (`CGO_ENABLED=0`, vet + test of `.` and the backend pkg for every tag on its native OS: `mac`/`purego` on macOS, `win`/`purego` on Windows, `x11`, `purego,x11`, `wayland`, `purego`, `libei`, `purego,libei` on Linux); job `cross` (`CGO_ENABLED=0` cross-build of `.` for every backend tag × GOOS/GOARCH); job `fmt` (`gofmt -l .` must be empty).
- `.circleci/config.yml` — Linux Cgo full tests: `xvfb-run go test -v ./...`.

The old `appveyor.yml` has been removed.

## Architecture

Single Go package `robotgo` at repo root (flat layout) with platform-specific files and C-binding subpackages. The default backend is Cgo wrappers over C headers vendored in subdirectories; build tags split platform implementations. In addition, five experimental **pure-Go (no-Cgo) backends** live in their own packages and are selected via build tags on the root package:

| Tag | Package | Root wiring file | Notes |
| --- | --- | --- | --- |
| `win` | `win/` | `windows_n.go` | Win32 via tailscale/win |
| `mac` | `darwin/` | `darwin.go` | Quartz/CoreGraphics via ebitengine/purego; window mgmt unsupported: `ActiveName` returns `ErrNotSupported`, `GetTitle` returns `""`, `MinWindow`/`MaxWindow`/`CloseWindow` are no-ops |
| `x11` | `x11/` | `x11_n.go` | XTEST/EWMH via jezek/xgb + xgbutil |
| `wayland` | `wayland/` | `wayland_n.go` | wlroots virtual-input/screencopy protocols |
| `libei` | `libei/` | `libei.go` | xdg-desktop-portal RemoteDesktop (GNOME/KDE) |
| `purego` | — | — | shortcut: `mac` on darwin, `win` on windows, `wayland` on linux; combine with `x11` or `libei` on Linux to override |

Default Cgo files carry `//go:build !wayland && !win && !libei && !mac && !x11 && !purego` (`robotgo.go`, `key.go`, `robotgo_fn_v1.go`, `robot_info_test.go`, ...) so any pure-Go tag disables the Cgo backend. Select at most one backend tag per build (`purego` may be paired with one Linux override): the constraints do not enforce mutual exclusion, e.g. `-tags "x11,libei"` compiles both `x11_n.go` and `libei.go` and fails with duplicate declarations. When adding a new pure-Go tag, extend these exclusions everywhere.

```
robotgo/
├── robotgo.go              # default Cgo API + preamble; excluded by any pure-Go tag
├── robotgo_pub.go          # portable (untagged): Version, MouseSleep, KeySleep, DisplayID, Scale...
├── doc.go                  # package doc + supported key names
├── robotgo_mac.go          # darwin && !mac && !purego
├── robotgo_mac_unix.go     # !windows, Cgo only (darwin + Cgo X11)
├── robotgo_mac_win.go      # (darwin || windows) && !win && !mac && !purego
├── robotgo_win.go          # windows (shared by Cgo and win backends)
├── robotgo_x11.go          # Cgo X11: !darwin && !windows && no pure-Go tag
├── robotgo_android.go, robotgo_adb.go
├── robotgo_ocr.go          # //go:build ocr only — not in normal builds (gosseract OCR)
├── darwin.go               # darwin && (mac || purego) — wires darwin/
├── windows_n.go            # windows && (win || purego) — wires win/
├── x11_n.go                # linux && x11 — wires x11/
├── wayland_n.go            # linux && (wayland || purego) && !libei && !x11 — wires wayland/
├── libei.go                # linux && libei — wires libei/
├── key.go                  # Cgo only; keycode.go, screen.go, img.go, ps.go untagged
├── robotgo_fn_v1.go        # deprecated v1 aliases (kept for compat), Cgo only
├── robot_info_test.go      # Cgo smoke tests (used by GitHub Actions)
├── robot_mac_test.go       # darwin && (mac || purego) tests
├── key_c_test.go           # Cgo-only key/mouse-arg unit tests; img_test.go untagged
├── robotgo_test.go         # interactive tests, (darwin || windows) Cgo only
├── base/       # C helpers (MMBitmap, rgb, microsleep, types, os, pubs, xdisplay)
├── mouse/      # Go pkg + C (mouse.h, mouse_c.h) with *_darwin.go/_windows.go/_x11.go
├── key/        # Go pkg + C (keycode.h, keycode_c.h, keypress.h, keypress_c.h, key_windows.go)
├── screen/     # Go pkg + C (goScreen.h, screen.go, screen_c.h, screengrab_c.h)
├── window/     # Go pkg + C (goWindow.h, window.h, alert_c.h, win_sys.h, pub.h)
├── clipboard/  # Go pkg (darwin/unix/windows variants) + cmd/gocopy, cmd/gopaste, example/
├── win/        # pure-Go Windows backend (no Cgo); //go:build windows
├── darwin/     # pure-Go macOS backend (purego + objc, app.go/cg.go); //go:build darwin
├── x11/        # pure-Go X11 backend (xgb/xgbutil); //go:build linux
├── wayland/    # pure-Go wlroots Wayland backend; internal/protocols/wlr_{foreign_toplevel,screencopy,virtual_keyboard,virtual_pointer}; libei/ (empty)
├── libei/      # pure-Go libei/xdg-portal backend (GNOME/KDE); //go:build linux
├── mcp/        # MCP server package (mcp.go, a bare `package mcp` stub)
├── cuse/       # computer/, browser/, gui/ — empty placeholder dirs
├── event/      # C headers for android/ios global hooks (event_c.h)
├── cv/         # OpenCV helper (gocv.go)
├── examples/   # main.go + mouse, key, screen, window, scale — runnable main.go
├── lang/       # translated READMEs (de, es, fr, ja, ko, pt, ru, zh, zht)
├── skills/     # SKILL.md (agent skill descriptor)
├── docs/       # install.md, keys.md, CHANGELOG.md, README.md, robotgo-logo.svg, archive/
├── test/       # index.html; tmp/, logs/ are git-ignored scratch dirs
└── .github/workflows/go.yml, .circleci/config.yml
```

Key subpackage relationships: the root `robotgo` package pulls C code from `screen/goScreen.h`, `mouse/mouse_c.h`, `window/goWindow.h`. The `key/` and `clipboard/` directories are importable Go packages; `base/` is header-only C support. The pure-Go `win/`, `darwin/`, `x11/`, `wayland/` and `libei/` packages each mirror the robotgo API surface (mouse/keyboard/screen/window/process, plus `errors.go`/`ErrNotSupported` where applicable) so a backend can be swapped per platform with a build tag; the root wiring files forward `robotgo.*` calls to them.

## Code Style

- **Copyright header**: every Go and C file starts with the 10-line `Copyright (c) 2016-2026 AtomAI...` block (see `CONTRIBUTING.md`). Preserve it verbatim when editing; add a second header only if authorship changes.
- **Indentation**: tabs (Go default). Run `gofmt` before committing.
- **Build tags**: use both forms together — `//go:build darwin` plus legacy `// +build darwin` lines — matching existing files (e.g. `robotgo_mac_win.go` uses one `// +build` line per AND-ed term).
- **Cgo**: `import "C"` immediately follows a `/* ... */` comment block containing `#cgo` directives and `#include`s. Keep LDFLAGS per-OS (`#cgo darwin LDFLAGS`, `#cgo linux LDFLAGS`, `#cgo windows LDFLAGS`).
- **Naming**: exported `CamelCase`; key constants prefixed `Key` (e.g. `KeyA`, `KeyEnter`); C-type aliases prefixed `C` (`CBitmap`, `CHex`). Mirror existing patterns.
- **Imports**: stdlib first, blank line, then third-party (`github.com/...`). Internal subpackage imported as `github.com/go-vgo/robotgo/clipboard` (full module path, not relative).
- **Comments**: godoc-style `// FuncName does X` on every exported symbol. Package doc lives in `doc.go` / top of `robotgo.go`.
- **Error handling**: return `error` as last result; use `errors.New` / `fmt.Errorf`. `Try(fn, handler)` helper wraps panics via `recover`. Do not swallow errors.
- **Types**: existing code uses `interface{}` (pre-generics) — match local style when editing that file, but prefer `any` in new code. Do not mass-rewrite; gopls emits hints, not errors.

## Testing

- Framework: stdlib `testing` + `github.com/vcaesar/tt` (`tt.Expect(t, want, got)`).
- Test files: `*_test.go` beside sources. Package declared as `robotgo_test` (external) for API-surface tests, or `robotgo` for internal.
- **Cgo smoke tests** live in `robot_info_test.go` — explicitly selected by GitHub Actions on macOS/Windows. They query screen size/location/scale/window title (not truly headless) and are excluded from pure-Go test runs by build tags. Keep new lightweight Cgo tests here; match build tags to the APIs being tested.
- **Pure-Go tests**: root `robot_mac_test.go` (darwin `mac`/`purego`) plus each backend's `*/robotgo_test.go` (`win/`, `darwin/`, `x11/`, `wayland/` + `keyboard_wire_test.go`, `libei/`). GitHub Actions runs `-tags purego .` on macOS/Windows and the x11/libei suites on Linux with `CGO_ENABLED=0`.
- **Interactive / display-required tests** go in `robotgo_test.go`; its tags restrict it to darwin/windows Cgo, so CircleCI's Linux `xvfb-run go test -v ./...` does **not** include it (Linux root coverage there comes from `robot_info_test.go`, `key_c_test.go`, `img_test.go`).
- Other unit tests: `key_c_test.go` (Cgo), `img_test.go` (untagged, also runs in pure-Go root jobs), `clipboard/*_test.go`.
- Run one test: `go test -v -run TestGetScreenSize .`
- No fixtures, snapshots, or golden files in use. Screenshots produced by examples are `.gitignore`d.

## Key Patterns

- **Cgo + platform split is mandatory**. Any new OS-specific function must be gated by `//go:build` tags and have implementations (even stub) for darwin, linux, windows — examine `mouse/mouse_darwin.go`, `mouse_windows.go`, `mouse_x11.go` as the template.
- **Keep backends in sync**: a new public `robotgo` API added to the build-tagged Cgo surface (e.g. `robotgo.go`, `key.go`, `robotgo_mac*.go`) needs a forwarder in each pure-Go wiring file (`darwin.go`, `windows_n.go`, `x11_n.go`, `wayland_n.go`, `libei.go`) and an implementation (or `ErrNotSupported`) in the matching backend package, otherwise `-tags purego` builds break. APIs in untagged portable files (`robotgo_pub.go`, `ps.go`, `screen.go`, `img.go`, `keycode.go`) are already shared by every backend — do not redeclare them in wiring files.
- **Free C-allocated bitmaps**: every `CaptureScreen`, `ToCBitmap`, etc. must be paired with `defer robotgo.FreeBitmap(bit)` or `robotgo.FreeBitmapArr(...)`. Leaking is a memory bug on all platforms.
- **Global tunables** are package-level vars, not config structs: `MouseSleep`, `KeySleep`, `DisplayID`, `NotPid`, `Scale`. Callers mutate them directly (see README examples). Do not hide them behind getters.
- **`robotgo_fn_v1.go`** contains deprecated v1 aliases — do not add new APIs there, but do not delete existing ones (backwards compatibility).
- **Version string** lives in `robotgo_pub.go` as `const Version = "v2.00.0.1658, MT. Baker!"` (moved out of `robotgo.go`). Bump it when releasing; `TestGetVer` asserts it matches `GetVersion()`. The pure-Go backends carry their own `const Version` (`win/` `v0.1.0-windows`, `darwin/` `v0.1.0-darwin`, `x11/` `v0.1.0-x11`, `wayland/` `v0.1.0-wayland`, `libei/` `v0.1.0-libei`).
- **Windows pid vs hwnd**: set `robotgo.NotPid = true` to pass window handles instead of pids into the window/key APIs on Windows.
- **macOS permissions**: most screen/input APIs silently fail without Accessibility + Screen Recording grants (both the Cgo and pure-Go `darwin/` backends). When reproducing bugs on darwin, verify System Settings → Privacy & Security first.
- **Do not vendor**: `vendor/` is in `.gitignore`; also avoid `go mod vendor` (upstream note in README references golang/go#26366).
- **C artifacts** (`*.cgo1.go`, `*.cgo2.c`, `_cgo_*`, `*.o`, `*.a` except whitelisted `cdeps/...` libpng archives) are git-ignored — do not commit them.
- **Commit sign-off** is expected (see `CONTRIBUTING.md`); PRs require ≥2 maintainer review (LGTM).

## Dependencies

- `github.com/jezek/xgb`, `github.com/jezek/xgbutil` — X11 protocol on Linux (also the `x11/` pure-Go backend).
- `github.com/vcaesar/go-wayland` — Wayland protocol client (used by the `wayland/` pure-Go backend).
- `github.com/godbus/dbus/v5` — D-Bus, used on Linux for Wayland/desktop ops and the `libei/` xdg-portal backend.
- `github.com/tailscale/win`, `golang.org/x/sys` (direct); `github.com/dblohm7/wingoes`, `github.com/yusufpapurcu/wmi`, `github.com/go-ole/go-ole` (indirect) — Windows system APIs (also used by the `win/` pure-Go backend).
- `github.com/ebitengine/purego` — direct; dlopen + `objc` runtime for the `darwin/` pure-Go backend. `github.com/gen2brain/shm` — indirect shared-memory helper for screenshot paths.
- `github.com/vcaesar/keycode` — cross-platform keycode mapping (used by `key/`).
- `github.com/vcaesar/imgo`, `golang.org/x/image` — image encode/decode (PNG/JPEG save).
- `github.com/vcaesar/screenshot` — screenshot backend.
- `github.com/vcaesar/gops` (direct, imported by `ps.go`), `github.com/shirou/gopsutil/v4` (indirect) — process enumeration (`FindIds`, `PidExists`, `Kill`).
- `github.com/vcaesar/tt` — testing assertions.
- `github.com/otiai10/gosseract/v2` — OCR (used by `robotgo_ocr.go`; needs `libtesseract`).
- Companion repos (not in `go.mod`, referenced in README/examples): `github.com/vcaesar/bitmap`, `github.com/vcaesar/gcv` (OpenCV), `github.com/jezek/gohook` (global event hook).
