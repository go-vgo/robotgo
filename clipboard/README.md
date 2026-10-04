This directory based on [clipboard](https://github.com/atotto/clipboard).

[![Build Status](https://travis-ci.org/atotto/clipboard.svg?branch=master)](https://travis-ci.org/atotto/clipboard)
[![GoDoc](https://godoc.org/github.com/atotto/clipboard?status.svg)](http://godoc.org/github.com/atotto/clipboard)
<!--[![Build Status](https://drone.io/github.com/atotto/clipboard/status.png)](https://drone.io/github.com/atotto/clipboard/latest) -->

# Clipboard for Go

Provide copying and pasting to the Clipboard for Go.

<!--Download shell commands at https://drone.io/github.com/atotto/clipboard/files-->

Build:

    $ go get github.com/atotto/clipboard

Platforms:

* OSX
* Windows 7 (probably work on other Windows)
* Linux, Unix, the first available tool is used:
  * Wayland: `wl-clipboard` (`wl-copy` / `wl-paste`), when `WAYLAND_DISPLAY` is set
  * X11: `xclip` or `xsel`
  * Android Termux: `termux-clipboard-get` / `termux-clipboard-set` (Termux:API add-on)
  * WSL: the Windows clipboard via `clip.exe` / `powershell.exe`
  * tmux: the tmux paste buffer, when running inside tmux (`TMUX` is set)
* Plan 9 (`/dev/snarf`)

`clipboard.Primary = true` selects the primary selection on Wayland and X11,
other tools fall back to the clipboard.


Document: 

* http://godoc.org/github.com/atotto/clipboard

Notes:

* Text string only
* UTF-8 text encoding only (no conversion)

TODO:

* Clipboard watcher(?)

## Commands:

paste shell command:

    $ go get github.com/atotto/clipboard/cmd/gopaste
    $ # example:
    $ gopaste > document.txt

copy shell command:

    $ go get github.com/atotto/clipboard/cmd/gocopy
    $ # example:
    $ cat document.txt | gocopy



