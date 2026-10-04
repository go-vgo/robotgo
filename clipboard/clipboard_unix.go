// Copyright 2013 @atotto. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build freebsd || linux || netbsd || openbsd || solaris || dragonfly
// +build freebsd linux netbsd openbsd solaris dragonfly

package clipboard

import (
	"bytes"
	"encoding/binary"
	"errors"
	"os"
	"os/exec"
	"strings"
	"unicode/utf16"
)

// tool describes a clipboard command line utility
type tool struct {
	paste, copy []string
	// pastePrimary and copyPrimary are nil when the tool has no
	// primary selection, the clipboard is used instead
	pastePrimary, copyPrimary []string

	// encode and decode convert the text when the tool does not speak UTF-8
	encode func(string) []byte
	decode func([]byte) string
}

var (
	// Primary choose primary mode on unix
	Primary bool

	current *tool

	// Wayland, wl-clipboard
	wlTool = &tool{
		paste:        []string{"wl-paste", "--no-newline"},
		copy:         []string{"wl-copy"},
		pastePrimary: []string{"wl-paste", "--no-newline", "--primary"},
		copyPrimary:  []string{"wl-copy", "--primary"},
	}
	// X11
	xclipTool = &tool{
		paste:        []string{"xclip", "-out", "-selection", "clipboard"},
		copy:         []string{"xclip", "-in", "-selection", "clipboard"},
		pastePrimary: []string{"xclip", "-out", "-selection", "primary"},
		copyPrimary:  []string{"xclip", "-in", "-selection", "primary"},
	}
	xselTool = &tool{
		paste:        []string{"xsel", "--output", "--clipboard"},
		copy:         []string{"xsel", "--input", "--clipboard"},
		pastePrimary: []string{"xsel", "--output", "--primary"},
		copyPrimary:  []string{"xsel", "--input", "--primary"},
	}
	// Android, Termux:API add-on
	termuxTool = &tool{
		paste: []string{"termux-clipboard-get"},
		copy:  []string{"termux-clipboard-set"},
	}
	// Windows Subsystem for Linux, uses the Windows clipboard
	wslTool = &tool{
		paste: []string{"powershell.exe", "-NoProfile", "-NonInteractive", "-Command",
			"[Console]::OutputEncoding = New-Object System.Text.UTF8Encoding $false; Get-Clipboard -Raw"},
		copy:   []string{"clip.exe"},
		encode: encodeUTF16LE,
		decode: decodeWSL,
	}
	// tmux paste buffer, for terminal sessions without a display server
	tmuxTool = &tool{
		paste: []string{"tmux", "save-buffer", "-"},
		copy:  []string{"tmux", "load-buffer", "-"},
	}

	errMissingCommands = errors.New("no clipboard utilities available, please install " +
		"wl-clipboard, xclip, xsel, Termux:API (termux-clipboard-get/set) or run inside WSL or tmux")
)

func init() {
	current = detect(os.Getenv, exec.LookPath)
	Unsupported = current == nil
}

// detect returns the first available clipboard tool for the current session
func detect(getenv func(string) string, lookPath func(string) (string, error)) *tool {
	var candidates []*tool
	if getenv("WAYLAND_DISPLAY") != "" {
		candidates = append(candidates, wlTool)
	}
	x11 := getenv("DISPLAY") != ""
	if x11 {
		candidates = append(candidates, xclipTool, xselTool)
	}
	candidates = append(candidates, termuxTool, wslTool)
	if getenv("TMUX") != "" {
		candidates = append(candidates, tmuxTool)
	}
	if !x11 {
		// the display may still be reachable, e.g. DISPLAY set later
		candidates = append(candidates, xclipTool, xselTool)
	}

	for _, t := range candidates {
		if found(lookPath, t.paste[0]) && found(lookPath, t.copy[0]) {
			return t
		}
	}
	return nil
}

func found(lookPath func(string) (string, error), name string) bool {
	_, err := lookPath(name)
	return err == nil
}

// encodeUTF16LE encodes text as UTF-16LE with a BOM,
// so clip.exe stores it as unicode instead of the ANSI code page
func encodeUTF16LE(text string) []byte {
	u := utf16.Encode([]rune(text))
	buf := make([]byte, 2+2*len(u))
	buf[0], buf[1] = 0xff, 0xfe
	for i, c := range u {
		binary.LittleEndian.PutUint16(buf[2+2*i:], c)
	}
	return buf
}

// decodeWSL drops the line terminator powershell appends to its output
func decodeWSL(out []byte) string {
	return string(bytes.TrimSuffix(out, []byte("\r\n")))
}

func getPasteCommand() *exec.Cmd {
	args := current.paste
	if Primary && current.pastePrimary != nil {
		args = current.pastePrimary
	}
	return exec.Command(args[0], args[1:]...)
}

func getCopyCommand() *exec.Cmd {
	args := current.copy
	if Primary && current.copyPrimary != nil {
		args = current.copyPrimary
	}
	return exec.Command(args[0], args[1:]...)
}

func readAll() (string, error) {
	if Unsupported || current == nil {
		return "", errMissingCommands
	}

	out, err := getPasteCommand().Output()
	if err != nil {
		return "", err
	}
	if current.decode != nil {
		return current.decode(out), nil
	}
	return string(out), nil
}

func writeAll(text string) error {
	if Unsupported || current == nil {
		return errMissingCommands
	}

	copyCmd := getCopyCommand()
	if current.encode != nil {
		copyCmd.Stdin = bytes.NewReader(current.encode(text))
	} else {
		copyCmd.Stdin = strings.NewReader(text)
	}
	return copyCmd.Run()
}
