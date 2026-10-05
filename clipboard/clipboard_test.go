// Copyright 2013 @atotto. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package clipboard_test

import (
	"strings"
	"testing"
	"time"

	"github.com/go-vgo/robotgo/clipboard"
)

func skipUnsupported(tb testing.TB) {
	tb.Helper()
	if clipboard.Unsupported {
		tb.Skip("clipboard utilities are not available")
	}
}

// roundTrip retries on mismatch: the system clipboard is shared with other
// test processes (go test ./... runs packages in parallel, and robotgo's
// TestClip/TestTypeStr write it too), so a write can be overwritten
// before it is read back.
func roundTrip(t *testing.T, expected string) {
	t.Helper()
	var actual string
	for range 5 {
		if err := clipboard.WriteAll(expected); err != nil {
			t.Fatal(err)
		}

		var err error
		actual, err = clipboard.ReadAll()
		if err != nil {
			t.Fatal(err)
		}
		if actual == expected {
			return
		}
		time.Sleep(100 * time.Millisecond)
	}
	t.Errorf("want %q (len %d), got %q (len %d)",
		short(expected), len(expected), short(actual), len(actual))
}

func short(s string) string {
	if len(s) > 64 {
		return s[:64] + "..."
	}
	return s
}

func TestCopyAndPaste(t *testing.T) {
	skipUnsupported(t)
	roundTrip(t, "日本語")
}

func TestMultiCopyAndPaste(t *testing.T) {
	skipUnsupported(t)
	roundTrip(t, "French: éèêëàùœç")
	roundTrip(t, "Weird UTF-8: 💩☃")
}

func TestCopyAndPasteEmpty(t *testing.T) {
	skipUnsupported(t)
	roundTrip(t, "")
}

func TestCopyAndPasteLarge(t *testing.T) {
	skipUnsupported(t)
	// larger than the old 1<<20 UTF-16 read limit on windows
	roundTrip(t, strings.Repeat("ab日", 1<<19))
}

func BenchmarkReadAll(b *testing.B) {
	skipUnsupported(b)
	for b.Loop() {
		if _, err := clipboard.ReadAll(); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkWriteAll(b *testing.B) {
	skipUnsupported(b)
	text := "いろはにほへと"
	for b.Loop() {
		if err := clipboard.WriteAll(text); err != nil {
			b.Fatal(err)
		}
	}
}
