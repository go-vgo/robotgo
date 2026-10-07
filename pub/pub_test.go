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

package pub

import "testing"

func TestSetDelay(t *testing.T) {
	k, m := KeySleep, MouseSleep
	t.Cleanup(func() { KeySleep, MouseSleep = k, m })

	SetDelay(25)
	if KeySleep != 25 || MouseSleep != 25 {
		t.Fatalf("SetDelay(25): KeySleep, MouseSleep = %d, %d", KeySleep, MouseSleep)
	}
	SetDelay()
	if KeySleep != 10 || MouseSleep != 10 {
		t.Fatalf("SetDelay(): KeySleep, MouseSleep = %d, %d, want 10", KeySleep, MouseSleep)
	}
}
