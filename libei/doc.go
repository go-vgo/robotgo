//go:build linux
// +build linux

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

// Package libei provides a pure-Go libei/portal implementation of the robotgo
// API for desktop automation on Wayland compositors that do NOT expose the
// wlroots virtual-input protocols — most importantly GNOME (Mutter) and
// KDE Plasma (KWin).
//
// Unlike the wayland backend (which speaks zwlr_virtual_pointer_v1 and
// zwp_virtual_keyboard_v1 directly to the compositor socket), this backend
// drives input through the freedesktop xdg-desktop-portal RemoteDesktop
// interface over D-Bus. The portal is the cross-desktop, security-reviewed
// path for emulated input and is the same mechanism libei is reached through.
//
// Requirements at runtime:
//   - A session D-Bus bus (DBUS_SESSION_BUS_ADDRESS).
//   - xdg-desktop-portal plus a desktop backend implementing
//     org.freedesktop.portal.RemoteDesktop (xdg-desktop-portal-gnome,
//     xdg-desktop-portal-kde, or xdg-desktop-portal-wlr).
//
// The first run shows a one-time "allow remote control" consent dialog. A
// restore token (persist_mode = persistent) is cached under
// $XDG_STATE_HOME/robotgo so subsequent runs do not re-prompt.
//
// Capability notes:
//   - Keyboard, relative pointer motion, buttons and scroll are supported.
//   - Absolute pointer motion (Move) requires a linked ScreenCast session and
//     is not yet wired up; Move falls back to a best-effort relative move and
//     screen/window helpers report ErrNotSupported.
//
// The transport is intentionally hidden behind the injector interface (see
// conn.go) so the real libei/EIS wire protocol (via RemoteDesktop.ConnectToEIS)
// can be slotted in later without changing this package's public API.
package libei
