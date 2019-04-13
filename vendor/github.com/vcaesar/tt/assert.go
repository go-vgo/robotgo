package tt

import "testing"

// Assertions provides assertion methods around the
// TestingT interface.
type Assertions struct {
	t TestingT
}

// New makes a new Assertions object for the specified TestingT.
func New(t TestingT) *Assertions {
	return &Assertions{
		t: t,
	}
}

// BM func Benchmark1(b *testing.B, fn func())
func (at *Assertions) BM(b *testing.B, fn func()) {
	for i := 0; i < b.N; i++ {
		fn()
	}
}

// Equal asserts that two objects are equal.
func (at *Assertions) Equal(expect, actual interface{}, args ...int) bool {
	call := 5
	if len(args) > 0 {
		call = args[0]
	}

	return Equal(at.t, expect, actual, call)
}

// Expect asserts that string and objects are equal.
func (at *Assertions) Expect(expect string, actual interface{}, args ...int) bool {
	call := 4
	if len(args) > 0 {
		call = args[0]
	}

	return Expect(at.t, expect, actual, call)
}

// Nil asserts that nil and objects are equal.
func (at *Assertions) Nil(actual interface{}, args ...int) bool {
	call := 4
	if len(args) > 0 {
		call = args[0]
	}

	return Expect(at.t, "<nil>", actual, call)
}

// Empty asserts that empty and objects are equal.
func (at *Assertions) Empty(actual interface{}, args ...int) bool {
	call := 4
	if len(args) > 0 {
		call = args[0]
	}

	return Expect(at.t, "", actual, call)
}

// Bool asserts that true and objects are equal.
func (at *Assertions) Bool(actual interface{}, args ...int) bool {
	call := 4
	if len(args) > 0 {
		call = args[0]
	}

	return Expect(at.t, "true", actual, call)
}

// True asserts that true and objects are equal.
func (at *Assertions) True(actual interface{}, args ...int) bool {
	call := 4
	if len(args) > 0 {
		call = args[0]
	}

	return Expect(at.t, "true", actual, call)
}

// False asserts that flase and objects are equal.
func (at *Assertions) False(actual interface{}, args ...int) bool {
	call := 4
	if len(args) > 0 {
		call = args[0]
	}

	return Expect(at.t, "false", actual, call)
}

// Not asserts that two objects are not equal.
func (at *Assertions) Not(expect, actual interface{}, args ...int) bool {
	call := 5
	if len(args) > 0 {
		call = args[0]
	}

	return Not(at.t, expect, actual, call)
}

// NotEqual asserts that two objects are not equal.
func (at *Assertions) NotEqual(expect, actual interface{}, args ...int) bool {
	call := 5
	if len(args) > 0 {
		call = args[0]
	}

	return NotEqual(at.t, expect, actual, call)
}

// NotExpect asserts that string and objects are not equal.
func (at *Assertions) NotExpect(expect string, actual interface{}, args ...int) bool {
	call := 4
	if len(args) > 0 {
		call = args[0]
	}

	return NotExpect(at.t, expect, actual, call)
}
