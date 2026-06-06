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

package libei

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/godbus/dbus/v5"
)

// Portal D-Bus addressing.
const (
	portalDest = "org.freedesktop.portal.Desktop"
	portalPath = "/org/freedesktop/portal/desktop"

	ifaceRemoteDesktop = "org.freedesktop.portal.RemoteDesktop"
	ifaceRequest       = "org.freedesktop.portal.Request"
	ifaceSession       = "org.freedesktop.portal.Session"
)

// RemoteDesktop device-type bitmask (SelectDevices "types").
const (
	deviceKeyboard = 1
	devicePointer  = 2
	// deviceTouchscreen = 4 // not used
)

// persist_mode values for SelectDevices.
const (
	persistNone       = 0
	persistTransient  = 1
	persistPersistent = 2
)

// Notify* state values (keyboard keycode/keysym, pointer button).
const (
	stateReleased uint32 = 0
	statePressed  uint32 = 1
)

// NotifyPointerAxisDiscrete axis values.
const (
	axisVertical   uint32 = 0
	axisHorizontal uint32 = 1
)

// Request.Response response codes.
const (
	respSuccess   uint32 = 0
	respCancelled uint32 = 1
	respEnded     uint32 = 2
)

// ErrNotSupported is returned when the running portal cannot satisfy a request
// (for example, screen capture or window management, which this backend does
// not implement, or absolute pointer motion without a linked ScreenCast).
var ErrNotSupported = errors.New("robotgo/libei: operation not supported by the RemoteDesktop portal")

// ErrNoConnection is returned when the RemoteDesktop portal session could not
// be established (no session bus, no portal backend, or the user denied access).
var ErrNoConnection = errors.New("robotgo/libei: RemoteDesktop portal session not established")

// injector abstracts the input transport. The default implementation drives
// the RemoteDesktop portal Notify* methods. A future implementation can speak
// the libei/EIS wire protocol obtained via RemoteDesktop.ConnectToEIS without
// changing this package's public API.
type injector interface {
	keyboardKeycode(code int32, state uint32) error
	keyboardKeysym(sym int32, state uint32) error
	pointerMotion(dx, dy float64) error
	pointerButton(button int32, state uint32) error
	pointerAxisDiscrete(axis uint32, steps int32) error
	pointerMotionAbsolute(stream uint32, x, y float64) error
}

// conn holds the singleton portal connection and negotiated session state.
type conn struct {
	mu sync.Mutex

	bus           *dbus.Conn
	obj           dbus.BusObject
	uniqueName    string
	sessionHandle dbus.ObjectPath

	devices uint32 // bitmask of granted device types
	streams []stream

	inj injector
}

// stream describes a ScreenCast PipeWire stream linked to the session. It is
// only populated when a ScreenCast source has been negotiated (not done yet),
// and is required for absolute pointer motion.
type stream struct {
	nodeID        uint32
	x, y          int32
	width, height int32
}

var (
	globalConn *conn
	connMu     sync.Mutex

	// tokenCounter makes handle/session tokens unique within a process.
	tokenCounter uint64
)

// ensureConn lazily establishes the portal session, re-establishing it if the
// previous connection was never created or has since been closed. Using a
// mutex (instead of sync.Once) keeps the backend recoverable after a failed
// handshake or an explicit Close.
func ensureConn() (*conn, error) {
	connMu.Lock()
	defer connMu.Unlock()

	if globalConn != nil && globalConn.healthy() {
		return globalConn, nil
	}

	c, err := newConn()
	if err != nil {
		globalConn = nil
		return nil, err
	}
	globalConn = c
	return globalConn, nil
}

// healthy reports whether the connection still has a live session bus.
func (c *conn) healthy() bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.bus != nil
}

// newConn connects to the session bus and runs the
// CreateSession -> SelectDevices -> Start handshake.
func newConn() (*conn, error) {
	bus, err := dbus.ConnectSessionBus()
	if err != nil {
		return nil, fmt.Errorf("%w: connect session bus: %v", ErrNoConnection, err)
	}

	c := &conn{
		bus:        bus,
		obj:        bus.Object(portalDest, dbus.ObjectPath(portalPath)),
		uniqueName: bus.Names()[0],
	}

	if err := c.createSession(); err != nil {
		_ = bus.Close()
		return nil, err
	}
	if err := c.selectDevices(); err != nil {
		_ = bus.Close()
		return nil, err
	}
	if err := c.start(); err != nil {
		_ = bus.Close()
		return nil, err
	}

	c.inj = &notifyInjector{c: c}
	return c, nil
}

// nextToken returns a process-unique token usable as a portal handle_token.
func nextToken(prefix string) string {
	n := atomic.AddUint64(&tokenCounter, 1)
	return fmt.Sprintf("robotgo_%s_%d", prefix, n)
}

// escapedSender converts the unique bus name (":1.42") into the form used in
// portal Request/Session object paths ("1_42").
func (c *conn) escapedSender() string {
	s := strings.TrimPrefix(c.uniqueName, ":")
	return strings.ReplaceAll(s, ".", "_")
}

// callPortal invokes a RemoteDesktop method that returns a Request handle and
// blocks until the corresponding Request.Response signal arrives, returning the
// results vardict. It subscribes to the predicted Response signal before making
// the call to avoid a race where the response fires first.
func (c *conn) callPortal(method string, args ...interface{}) (map[string]dbus.Variant, error) {
	token := nextToken("req")
	reqPath := dbus.ObjectPath(fmt.Sprintf(
		"/org/freedesktop/portal/desktop/request/%s/%s", c.escapedSender(), token))

	// Subscribe to the Response signal before issuing the call.
	if err := c.bus.AddMatchSignal(
		dbus.WithMatchObjectPath(reqPath),
		dbus.WithMatchInterface(ifaceRequest),
		dbus.WithMatchMember("Response"),
	); err != nil {
		return nil, fmt.Errorf("%w: add match: %v", ErrNoConnection, err)
	}
	defer c.bus.RemoveMatchSignal(
		dbus.WithMatchObjectPath(reqPath),
		dbus.WithMatchInterface(ifaceRequest),
		dbus.WithMatchMember("Response"),
	)

	sigCh := make(chan *dbus.Signal, 4)
	c.bus.Signal(sigCh)
	defer c.bus.RemoveSignal(sigCh)

	// The last options arg always carries our handle_token so the portal
	// uses the path we predicted above.
	if len(args) == 0 {
		return nil, errors.New("robotgo/libei: callPortal requires an options argument")
	}
	if opts, ok := args[len(args)-1].(map[string]dbus.Variant); ok {
		opts["handle_token"] = dbus.MakeVariant(token)
	}

	var returnedPath dbus.ObjectPath
	call := c.obj.Call(ifaceRemoteDesktop+"."+method, 0, args...)
	if call.Err != nil {
		return nil, fmt.Errorf("%w: %s: %v", ErrNoConnection, method, call.Err)
	}
	if err := call.Store(&returnedPath); err != nil {
		return nil, fmt.Errorf("%w: %s store handle: %v", ErrNoConnection, method, err)
	}

	timeout := time.NewTimer(60 * time.Second)
	defer timeout.Stop()

	for {
		select {
		case sig := <-sigCh:
			if sig == nil || sig.Path != reqPath || sig.Name != ifaceRequest+".Response" {
				continue
			}
			if len(sig.Body) < 2 {
				return nil, fmt.Errorf("robotgo/libei: %s malformed response", method)
			}
			code, _ := sig.Body[0].(uint32)
			results, _ := sig.Body[1].(map[string]dbus.Variant)
			switch code {
			case respSuccess:
				return results, nil
			case respCancelled:
				// The user denied or cancelled the portal dialog: the session
				// could not be established, not a missing capability.
				return nil, fmt.Errorf("%w: %s cancelled by user", ErrNoConnection, method)
			default:
				return nil, fmt.Errorf("%w: %s ended (code %d)", ErrNoConnection, method, code)
			}
		case <-timeout.C:
			return nil, fmt.Errorf("%w: %s timed out", ErrNoConnection, method)
		}
	}
}

// createSession creates the RemoteDesktop session and records its handle.
func (c *conn) createSession() error {
	sessionToken := nextToken("session")
	opts := map[string]dbus.Variant{
		"session_handle_token": dbus.MakeVariant(sessionToken),
	}
	results, err := c.callPortal("CreateSession", opts)
	if err != nil {
		return err
	}

	// Prefer the handle the portal reports; fall back to the predicted path.
	// Portals may return the handle either as a dbus.ObjectPath or a plain
	// string depending on the backend, so accept both.
	if v, ok := results["session_handle"]; ok {
		switch s := v.Value().(type) {
		case dbus.ObjectPath:
			if s != "" {
				c.sessionHandle = s
			}
		case string:
			if s != "" {
				c.sessionHandle = dbus.ObjectPath(s)
			}
		}
	}
	if c.sessionHandle == "" {
		c.sessionHandle = dbus.ObjectPath(fmt.Sprintf(
			"/org/freedesktop/portal/desktop/session/%s/%s", c.escapedSender(), sessionToken))
	}
	return nil
}

// selectDevices requests keyboard + pointer access, reusing a cached restore
// token (persistent) so the user is only prompted on the first run.
func (c *conn) selectDevices() error {
	opts := map[string]dbus.Variant{
		"types":        dbus.MakeVariant(uint32(deviceKeyboard | devicePointer)),
		"persist_mode": dbus.MakeVariant(uint32(persistPersistent)),
	}
	if tok := loadRestoreToken(); tok != "" {
		opts["restore_token"] = dbus.MakeVariant(tok)
	}
	_, err := c.callPortal("SelectDevices", c.sessionHandle, opts)
	return err
}

// start activates the session and records granted devices, any ScreenCast
// streams, and the refreshed (single-use) restore token.
func (c *conn) start() error {
	opts := map[string]dbus.Variant{}
	results, err := c.callPortal("Start", c.sessionHandle, "", opts)
	if err != nil {
		return err
	}

	if v, ok := results["devices"]; ok {
		if d, ok := v.Value().(uint32); ok {
			c.devices = d
		}
	}
	if v, ok := results["restore_token"]; ok {
		if tok, ok := v.Value().(string); ok && tok != "" {
			// single-use: always re-persist the latest. A persistence failure
			// is non-fatal (the user is re-prompted next run), but surface it
			// rather than dropping it silently.
			if werr := saveRestoreToken(tok); werr != nil {
				log.Printf("robotgo/libei: persist restore token: %v", werr)
			}
		}
	}
	c.streams = parseStreams(results["streams"])

	if c.devices == 0 {
		return fmt.Errorf("%w: no input devices were granted", ErrNotSupported)
	}
	return nil
}

// parseStreams decodes the Start "streams" result (a(ua{sv})).
func parseStreams(v dbus.Variant) []stream {
	raw, ok := v.Value().([][]interface{})
	if !ok {
		return nil
	}
	var out []stream
	for _, item := range raw {
		if len(item) < 2 {
			continue
		}
		node, _ := item[0].(uint32)
		props, _ := item[1].(map[string]dbus.Variant)
		s := stream{nodeID: node}
		if p, ok := props["position"]; ok {
			if xy, ok := p.Value().([]interface{}); ok && len(xy) == 2 {
				s.x, _ = xy[0].(int32)
				s.y, _ = xy[1].(int32)
			}
		}
		if p, ok := props["size"]; ok {
			if wh, ok := p.Value().([]interface{}); ok && len(wh) == 2 {
				s.width, _ = wh[0].(int32)
				s.height, _ = wh[1].(int32)
			}
		}
		out = append(out, s)
	}
	return out
}

// hasKeyboard reports whether keyboard input was granted.
func (c *conn) hasKeyboard() bool { return c.devices&deviceKeyboard != 0 }

// hasPointer reports whether pointer input was granted.
func (c *conn) hasPointer() bool { return c.devices&devicePointer != 0 }

// Close closes the portal session and the bus connection. After Close, a
// subsequent call into the backend re-establishes a fresh session.
func Close() {
	connMu.Lock()
	defer connMu.Unlock()
	if globalConn == nil {
		return
	}
	c := globalConn
	globalConn = nil

	c.mu.Lock()
	defer c.mu.Unlock()
	if c.bus != nil {
		// Best-effort close of the session object.
		_ = c.bus.Object(portalDest, c.sessionHandle).
			Call(ifaceSession+".Close", 0).Err
		_ = c.bus.Close()
		c.bus = nil
	}
}

// --- restore token persistence ---

func tokenPath() string {
	dir := os.Getenv("XDG_STATE_HOME")
	if dir == "" {
		if home, err := os.UserHomeDir(); err == nil {
			dir = filepath.Join(home, ".local", "state")
		}
	}
	if dir == "" {
		dir = os.TempDir()
	}
	return filepath.Join(dir, "robotgo", "portal_token")
}

func loadRestoreToken() string {
	data, err := os.ReadFile(tokenPath())
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(data))
}

func saveRestoreToken(tok string) error {
	p := tokenPath()
	if err := os.MkdirAll(filepath.Dir(p), 0o700); err != nil {
		return err
	}
	return os.WriteFile(p, []byte(tok), 0o600)
}
