// Copyright (c) 2016-2026 AtomAI, All rights reserved.
//
// See the COPYRIGHT file at the top-level directory of this distribution and at
// https://github.com/go-vgo/robotgo/blob/master/LICENSE
//
// Licensed under the Apache License, Version 2.0 <LICENSE-APACHE or
// http://www.apache.org/licenses/LICENSE-2.0>
//
// This file may not be copied, modified, or distributed
// except according to those terms.

// Package wlcursor reads the global pointer position on Wayland through
// compositor-specific IPC, since no Wayland protocol or portal exposes it:
// Hyprland's request socket (cursorpos) and KWin's scripting D-Bus API
// (workspace.cursorPos). Coordinates are global logical layout coordinates.
package wlcursor

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"
	"time"

	"github.com/godbus/dbus/v5"
)

// ErrUnavailable is returned when the running compositor offers no cursor
// position query (GNOME, Sway and other wlroots compositors).
var ErrUnavailable = errors.New("wlcursor: compositor does not expose the cursor position")

// Timeout bounds one position query.
const Timeout = 2 * time.Second

// Backoff is how long the KWin query is skipped after it fails, so a broken
// session does not cost Timeout on every call.
const Backoff = 10 * time.Second

// kwinRetry is the UnixNano time before which the KWin query is skipped.
var kwinRetry atomic.Int64

// Position returns the global pointer position from the running compositor.
func Position() (int, int, error) {
	if sig := os.Getenv("HYPRLAND_INSTANCE_SIGNATURE"); sig != "" {
		return hyprland(hyprSocket(sig))
	}
	if strings.Contains(strings.ToUpper(os.Getenv("XDG_CURRENT_DESKTOP")), "KDE") {
		if time.Now().UnixNano() < kwinRetry.Load() {
			return 0, 0, fmt.Errorf("%w: kwin query failed recently", ErrUnavailable)
		}
		x, y, err := kwinSession()
		if err != nil {
			kwinRetry.Store(time.Now().Add(Backoff).UnixNano())
		}
		return x, y, err
	}
	return 0, 0, ErrUnavailable
}

// kwinSession runs the KWin query on a new session bus connection.
func kwinSession() (int, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), Timeout)
	defer cancel()
	bus, err := dbus.ConnectSessionBus(dbus.WithContext(ctx))
	if err != nil {
		return 0, 0, err
	}
	defer bus.Close()
	return kwin(bus, os.Getenv("KDE_SESSION_VERSION") == "5")
}

// hyprSocket returns the Hyprland request socket path; Hyprland < 0.40 used
// /tmp/hypr instead of $XDG_RUNTIME_DIR/hypr.
func hyprSocket(sig string) string {
	p := filepath.Join(os.Getenv("XDG_RUNTIME_DIR"), "hypr", sig, ".socket.sock")
	if _, err := os.Stat(p); err == nil {
		return p
	}
	return filepath.Join("/tmp/hypr", sig, ".socket.sock")
}

// hyprland sends "j/cursorpos" and decodes {"x": .., "y": ..}. Hyprland
// serves the socket synchronously, so the connection is closed right away.
func hyprland(sock string) (int, int, error) {
	c, err := net.DialTimeout("unix", sock, Timeout)
	if err != nil {
		return 0, 0, err
	}
	defer c.Close()
	if err := c.SetDeadline(time.Now().Add(Timeout)); err != nil {
		return 0, 0, err
	}
	if _, err := c.Write([]byte("j/cursorpos")); err != nil {
		return 0, 0, err
	}
	var p struct {
		X, Y *int
	}
	if err := json.NewDecoder(c).Decode(&p); err != nil {
		return 0, 0, fmt.Errorf("wlcursor: hyprland cursorpos: %w", err)
	}
	if p.X == nil || p.Y == nil {
		return 0, 0, errors.New("wlcursor: hyprland cursorpos: missing x/y")
	}
	return *p.X, *p.Y, nil
}

// KWin scripting D-Bus addressing.
const (
	kwinDest      = "org.kde.KWin"
	kwinScripting = "org.kde.kwin.Scripting"
	kwinScript    = "org.kde.kwin.Script"

	replyPath  = dbus.ObjectPath("/org/robotgo/Cursor")
	replyIface = "org.robotgo.Cursor"
)

var scriptSeq atomic.Uint64

// kwin loads a one-shot KWin script that calls back with workspace.cursorPos
// over D-Bus (the approach used by kdotool), then unloads it. plasma5 selects
// the Plasma 5 script object path.
func kwin(bus *dbus.Conn, plasma5 bool) (x, y int, err error) {
	type pos struct{ x, y int }
	ch := make(chan pos, 1)
	reply := func(s string) *dbus.Error {
		var p pos
		if _, err := fmt.Sscanf(s, "%d %d", &p.x, &p.y); err == nil {
			select {
			case ch <- p:
			default:
			}
		}
		return nil
	}
	if err := bus.ExportMethodTable(map[string]any{"result": reply}, replyPath, replyIface); err != nil {
		return 0, 0, err
	}
	defer func() {
		err = errors.Join(err, bus.ExportMethodTable(nil, replyPath, replyIface))
	}()

	name := fmt.Sprintf("robotgo-cursor-%d-%d", os.Getpid(), scriptSeq.Add(1))
	f, err := os.CreateTemp("", name+"-*.js")
	if err != nil {
		return 0, 0, err
	}
	defer os.Remove(f.Name())
	_, err = fmt.Fprintf(f, "const p = workspace.cursorPos;\n"+
		"callDBus(%q, %q, %q, \"result\", Math.round(p.x) + \" \" + Math.round(p.y));\n",
		bus.Names()[0], replyPath, replyIface)
	if cerr := f.Close(); err == nil {
		err = cerr
	}
	if err != nil {
		return 0, 0, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), Timeout)
	defer cancel()
	scripting := bus.Object(kwinDest, "/Scripting")
	var id int32
	if err := scripting.CallWithContext(ctx, kwinScripting+".loadScript", 0, f.Name(), name).Store(&id); err != nil {
		return 0, 0, fmt.Errorf("wlcursor: kwin loadScript: %w", err)
	}
	if id < 0 {
		return 0, 0, fmt.Errorf("wlcursor: kwin loadScript returned %d", id)
	}
	defer func() {
		uctx, ucancel := context.WithTimeout(context.Background(), Timeout)
		defer ucancel()
		uerr := scripting.CallWithContext(uctx, kwinScripting+".unloadScript", 0, name).Err
		if uerr != nil {
			err = errors.Join(err, fmt.Errorf("wlcursor: kwin unloadScript: %w", uerr))
		}
	}()

	path := dbus.ObjectPath(fmt.Sprintf("/Scripting/Script%d", id))
	if plasma5 {
		path = dbus.ObjectPath(fmt.Sprintf("/%d", id))
	}
	if err := bus.Object(kwinDest, path).CallWithContext(ctx, kwinScript+".run", 0).Err; err != nil {
		return 0, 0, fmt.Errorf("wlcursor: kwin run: %w", err)
	}
	select {
	case p := <-ch:
		return p.x, p.y, nil
	case <-ctx.Done():
		return 0, 0, fmt.Errorf("wlcursor: kwin cursorPos: %w", ctx.Err())
	}
}
