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

package wlcursor

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"testing"

	"github.com/godbus/dbus/v5"
)

// shortDir returns a temp dir short enough for a unix socket path.
func shortDir(t *testing.T) string {
	t.Helper()
	dir, err := os.MkdirTemp("", "wlc")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.RemoveAll(dir) })
	return dir
}

// fakeHyprland serves one request on sock with reply and records the request.
func fakeHyprland(t *testing.T, sock, reply string) <-chan string {
	t.Helper()
	ln, err := net.Listen("unix", sock)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { ln.Close() })
	got := make(chan string, 1)
	go func() {
		c, err := ln.Accept()
		if err != nil {
			return
		}
		defer c.Close()
		buf := make([]byte, 64)
		n, _ := c.Read(buf)
		got <- string(buf[:n])
		fmt.Fprint(c, reply)
	}()
	return got
}

func TestHyprland(t *testing.T) {
	sock := filepath.Join(shortDir(t), "s.sock")
	got := fakeHyprland(t, sock, "\n{\n    \"x\": 1234,\n    \"y\": -56\n}\n")
	x, y, err := hyprland(sock)
	if err != nil || x != 1234 || y != -56 {
		t.Fatalf("hyprland = %d, %d, %v; want 1234, -56, nil", x, y, err)
	}
	if req := <-got; req != "j/cursorpos" {
		t.Fatalf("request = %q", req)
	}
}

func TestHyprlandErrors(t *testing.T) {
	dir := shortDir(t)
	if _, _, err := hyprland(filepath.Join(dir, "missing.sock")); err == nil {
		t.Fatal("missing socket: want error")
	}
	for i, reply := range []string{"unknown request", `{"x": 1}`} {
		sock := filepath.Join(dir, fmt.Sprintf("%d.sock", i))
		fakeHyprland(t, sock, reply)
		if _, _, err := hyprland(sock); err == nil {
			t.Fatalf("reply %q: want error", reply)
		}
	}
}

func TestPositionDispatch(t *testing.T) {
	t.Setenv("HYPRLAND_INSTANCE_SIGNATURE", "")
	t.Setenv("XDG_CURRENT_DESKTOP", "GNOME")
	if _, _, err := Position(); !errors.Is(err, ErrUnavailable) {
		t.Fatalf("GNOME: err = %v, want ErrUnavailable", err)
	}

	run := shortDir(t)
	if err := os.MkdirAll(filepath.Join(run, "hypr", "sig"), 0o700); err != nil {
		t.Fatal(err)
	}
	t.Setenv("XDG_RUNTIME_DIR", run)
	t.Setenv("HYPRLAND_INSTANCE_SIGNATURE", "sig")
	fakeHyprland(t, filepath.Join(run, "hypr", "sig", ".socket.sock"), `{"x": 7, "y": 8}`)
	if x, y, err := Position(); err != nil || x != 7 || y != 8 {
		t.Fatalf("Hyprland: Position = %d, %d, %v", x, y, err)
	}
}

// privateBus starts a dbus-daemon for the test and returns its address.
func privateBus(t *testing.T) string {
	t.Helper()
	bin, err := exec.LookPath("dbus-daemon")
	if err != nil {
		t.Skip("dbus-daemon not installed")
	}
	dir := shortDir(t)
	cfg := filepath.Join(dir, "bus.conf")
	conf := `<busconfig><type>session</type><listen>unix:tmpdir=` + dir + `</listen>
<auth>EXTERNAL</auth><policy context="default"><allow send_destination="*" eavesdrop="true"/>
<allow eavesdrop="true"/><allow own="*"/></policy></busconfig>`
	if err := os.WriteFile(cfg, []byte(conf), 0o600); err != nil {
		t.Fatal(err)
	}
	cmd := exec.Command(bin, "--config-file="+cfg, "--nofork", "--print-address")
	out, err := cmd.StdoutPipe()
	if err != nil {
		t.Fatal(err)
	}
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := cmd.Process.Kill(); err != nil {
			t.Error(err)
		}
		// Wait reports the kill as an *exec.ExitError.
		var exit *exec.ExitError
		if err := cmd.Wait(); err != nil && !errors.As(err, &exit) {
			t.Error(err)
		}
	})
	addr, err := bufio.NewReader(out).ReadString('\n')
	if err != nil {
		t.Fatal(err)
	}
	return strings.TrimSpace(addr)
}

func connect(t *testing.T, addr string) *dbus.Conn {
	t.Helper()
	c, err := dbus.Connect(addr)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { c.Close() })
	return c
}

var callRe = regexp.MustCompile(`callDBus\("([^"]+)", "([^"]+)", "([^"]+)", "result"`)

// fakeKWin mimics org.kde.KWin scripting: run answers the loaded script's
// callDBus target with reply.
type fakeKWin struct {
	bus      *dbus.Conn
	reply    string
	mu       sync.Mutex
	script   string
	unloaded []string
}

func (k *fakeKWin) export(t *testing.T, scriptPath dbus.ObjectPath) {
	t.Helper()
	if _, err := k.bus.RequestName(kwinDest, dbus.NameFlagDoNotQueue); err != nil {
		t.Fatal(err)
	}
	load := func(path, name string) (int32, *dbus.Error) {
		data, err := os.ReadFile(path)
		if err != nil {
			return 0, dbus.MakeFailedError(err)
		}
		k.mu.Lock()
		k.script = string(data)
		k.mu.Unlock()
		return 7, nil
	}
	unload := func(name string) (bool, *dbus.Error) {
		k.mu.Lock()
		k.unloaded = append(k.unloaded, name)
		k.mu.Unlock()
		return true, nil
	}
	run := func() *dbus.Error {
		k.mu.Lock()
		m := callRe.FindStringSubmatch(k.script)
		k.mu.Unlock()
		if m == nil {
			return dbus.MakeFailedError(errors.New("no callDBus in script"))
		}
		k.bus.Object(m[1], dbus.ObjectPath(m[2])).Go(m[3]+".result", 0, nil, k.reply)
		return nil
	}
	for _, err := range []error{
		k.bus.ExportMethodTable(map[string]any{"loadScript": load, "unloadScript": unload},
			"/Scripting", kwinScripting),
		k.bus.ExportMethodTable(map[string]any{"run": run}, scriptPath, kwinScript),
	} {
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestKWin(t *testing.T) {
	addr := privateBus(t)
	for _, tc := range []struct {
		plasma5 bool
		path    dbus.ObjectPath
	}{{false, "/Scripting/Script7"}, {true, "/7"}} {
		t.Run(string(tc.path), func(t *testing.T) {
			k := &fakeKWin{bus: connect(t, addr), reply: "640 -480"}
			k.export(t, tc.path)

			x, y, err := kwin(connect(t, addr), tc.plasma5)
			if err != nil || x != 640 || y != -480 {
				t.Fatalf("kwin = %d, %d, %v; want 640, -480, nil", x, y, err)
			}
			k.mu.Lock()
			defer k.mu.Unlock()
			if !strings.Contains(k.script, "workspace.cursorPos") {
				t.Fatalf("script does not read cursorPos:\n%s", k.script)
			}
			if len(k.unloaded) != 1 || !strings.HasPrefix(k.unloaded[0], "robotgo-cursor-") {
				t.Fatalf("unloaded = %v", k.unloaded)
			}
		})
	}
}

func TestKWinNoService(t *testing.T) {
	addr := privateBus(t)
	if _, _, err := kwin(connect(t, addr), false); err == nil {
		t.Fatal("no KWin on the bus: want error")
	}
}

func TestKWinBackoff(t *testing.T) {
	t.Setenv("HYPRLAND_INSTANCE_SIGNATURE", "")
	t.Setenv("XDG_CURRENT_DESKTOP", "KDE")
	t.Setenv("DBUS_SESSION_BUS_ADDRESS", "unix:path="+filepath.Join(shortDir(t), "none"))
	kwinRetry.Store(0)
	t.Cleanup(func() { kwinRetry.Store(0) })

	if _, _, err := Position(); err == nil || errors.Is(err, ErrUnavailable) {
		t.Fatalf("first call: err = %v, want connect error", err)
	}
	if _, _, err := Position(); !errors.Is(err, ErrUnavailable) {
		t.Fatalf("second call: err = %v, want ErrUnavailable (backoff)", err)
	}
}
