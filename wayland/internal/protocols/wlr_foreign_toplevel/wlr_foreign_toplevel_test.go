package wlr_foreign_toplevel

import "testing"

func TestWireStringBounds(t *testing.T) {
	// length 3 incl. NUL: "ab\0" + padding
	good := []byte{3, 0, 0, 0, 'a', 'b', 0, 0}
	if got := wireString(good, 0); got != "ab" {
		t.Errorf("good: got %q", got)
	}
	// Length claims more than the payload: must not panic.
	bad := []byte{200, 0, 0, 0, 'a', 'b', 0, 0}
	if got := wireString(bad, 0); got != "" {
		t.Errorf("overlong: got %q", got)
	}
	if got := wireString([]byte{1, 0}, 0); got != "" {
		t.Errorf("truncated header: got %q", got)
	}
	if got := wireString([]byte{0, 0, 0, 0}, 0); got != "" {
		t.Errorf("empty: got %q", got)
	}
}

func TestDispatchMalformedDoesNotPanic(t *testing.T) {
	h := &ZwlrForeignToplevelHandleV1{}
	var gotState []byte
	h.SetTitleHandler(func(ZwlrForeignToplevelHandleV1TitleEvent) {})
	h.SetStateHandler(func(e ZwlrForeignToplevelHandleV1StateEvent) { gotState = e.State })
	h.SetOutputEnterHandler(func(ZwlrForeignToplevelHandleV1OutputEnterEvent) {})
	h.Dispatch(0, -1, []byte{9, 0, 0, 0, 'x'})
	h.Dispatch(4, -1, []byte{16, 0, 0, 0, 2, 0, 0, 0})
	h.Dispatch(2, -1, []byte{1})
	if len(gotState) != 4 {
		t.Errorf("state truncated to payload: got %v", gotState)
	}
}
