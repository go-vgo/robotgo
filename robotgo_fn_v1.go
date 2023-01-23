package robotgo

import "github.com/vcaesar/tt"

// Deprecated: use the Move(),
//
// MoveMouse move the mouse
func MoveMouse(x, y int) {
	Move(x, y)
}

// Deprecated: use the DragSmooth(),
//
// DragMouse drag the mouse to (x, y),
// It's same with the DragSmooth() now
func DragMouse(x, y int, args ...interface{}) {
	Toggle("left")
	MilliSleep(50)
	// Drag(x, y, args...)
	MoveSmooth(x, y, args...)
	Toggle("left", "up")
}

// Deprecated: use the MoveSmooth(),
//
// MoveMouseSmooth move the mouse smooth,
// moves mouse to x, y human like, with the mouse button up.
func MoveMouseSmooth(x, y int, args ...interface{}) bool {
	return MoveSmooth(x, y, args...)
}

// Deprecated: use the function Location()
//
// GetMousePos get the mouse's position return x, y
func GetMousePos() (int, int) {
	return Location()
}

// Deprecated: use the Click(),
//
// # MouseClick click the mouse
//
// robotgo.MouseClick(button string, double bool)
func MouseClick(args ...interface{}) {
	Click(args...)
}

// Deprecated: use the TypeStr(),
//
// # TypeStringDelayed type string delayed, Wno-deprecated
//
// This function will be removed in version v1.0.0
func TypeStringDelayed(str string, delay int) {
	tt.Drop("TypeStringDelayed", "TypeStrDelay")
	TypeStrDelay(str, delay)
}

// Deprecated: use the ScaledF(),
//
// Scale1 get the screen scale (only windows old), drop
func Scale1() int {
	dpi := map[int]int{
		0: 100,
		// DPI Scaling Level
		96:  100,
		120: 125,
		144: 150,
		168: 175,
		192: 200,
		216: 225,
		// Custom DPI
		240: 250,
		288: 300,
		384: 400,
		480: 500,
	}

	x := ScaleX()
	return dpi[x]
}

// Deprecated: use the ScaledF(),
//
// Scale0 return ScaleX() / 0.96, drop
func Scale0() int {
	return int(float64(ScaleX()) / 0.96)
}

// Deprecated: use the ScaledF(),
//
// Mul mul the scale, drop
func Mul(x int) int {
	s := Scale1()
	return x * s / 100
}
