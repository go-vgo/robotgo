package xprop

import (
	"fmt"

	"github.com/BurntSushi/xgb"
	"github.com/BurntSushi/xgb/xproto"

	"github.com/BurntSushi/xgbutil"
)

// GetProperty abstracts the messiness of calling xgb.GetProperty.
func GetProperty(xu *xgbutil.XUtil, win xproto.Window, atom string) (
	*xproto.GetPropertyReply, error) {

	atomId, err := Atm(xu, atom)
	if err != nil {
		return nil, err
	}

	reply, err := xproto.GetProperty(xu.Conn(), false, win, atomId,
		xproto.GetPropertyTypeAny, 0, (1<<32)-1).Reply()

	if err != nil {
		return nil, fmt.Errorf("GetProperty: Error retrieving property '%s' "+
			"on window %x: %s", atom, win, err)
	}

	if reply.Format == 0 {
		return nil, fmt.Errorf("GetProperty: No such property '%s' on "+
			"window %x.", atom, win)
	}

	return reply, nil
}

// ChangeProperty abstracts the semi-nastiness of xgb.ChangeProperty.
func ChangeProp(xu *xgbutil.XUtil, win xproto.Window, format byte, prop string,
	typ string, data []byte) error {

	propAtom, err := Atm(xu, prop)
	if err != nil {
		return err
	}

	typAtom, err := Atm(xu, typ)
	if err != nil {
		return err
	}

	return xproto.ChangePropertyChecked(xu.Conn(), xproto.PropModeReplace, win,
		propAtom, typAtom, format,
		uint32(len(data)/(int(format)/8)), data).Check()
}

// ChangeProperty32 makes changing 32 bit formatted properties easier
// by constructing the raw X data for you.
func ChangeProp32(xu *xgbutil.XUtil, win xproto.Window, prop string, typ string,
	data ...uint) error {

	buf := make([]byte, len(data)*4)
	for i, datum := range data {
		xgb.Put32(buf[(i*4):], uint32(datum))
	}

	return ChangeProp(xu, win, 32, prop, typ, buf)
}

// WindowToUint is a covenience function for converting []xproto.Window
// to []uint.
func WindowToInt(ids []xproto.Window) []uint {
	ids32 := make([]uint, len(ids))
	for i, v := range ids {
		ids32[i] = uint(v)
	}
	return ids32
}

// AtomToInt is a covenience function for converting []xproto.Atom
// to []uint.
func AtomToUint(ids []xproto.Atom) []uint {
	ids32 := make([]uint, len(ids))
	for i, v := range ids {
		ids32[i] = uint(v)
	}
	return ids32
}

// StrToAtoms is a convenience function for converting
// []string to []uint32 atoms.
// NOTE: If an atom name in the list doesn't exist, it will be created.
func StrToAtoms(xu *xgbutil.XUtil, atomNames []string) ([]uint, error) {
	var err error
	atoms := make([]uint, len(atomNames))
	for i, atomName := range atomNames {
		a, err := Atom(xu, atomName, false)
		if err != nil {
			return nil, err
		}
		atoms[i] = uint(a)
	}
	return atoms, err
}

// PropValAtom transforms a GetPropertyReply struct into an ATOM name.
// The property reply must be in 32 bit format.
func PropValAtom(xu *xgbutil.XUtil, reply *xproto.GetPropertyReply,
	err error) (string, error) {

	if err != nil {
		return "", err
	}
	if reply.Format != 32 {
		return "", fmt.Errorf("PropValAtom: Expected format 32 but got %d",
			reply.Format)
	}

	return AtomName(xu, xproto.Atom(xgb.Get32(reply.Value)))
}

// PropValAtoms is the same as PropValAtom, except that it returns a slice
// of atom names. Also must be 32 bit format.
// This is a method of an XUtil struct, unlike the other 'PropVal...' functions.
func PropValAtoms(xu *xgbutil.XUtil, reply *xproto.GetPropertyReply,
	err error) ([]string, error) {

	if err != nil {
		return nil, err
	}
	if reply.Format != 32 {
		return nil, fmt.Errorf("PropValAtoms: Expected format 32 but got %d",
			reply.Format)
	}

	ids := make([]string, reply.ValueLen)
	vals := reply.Value
	for i := 0; len(vals) >= 4; i++ {
		ids[i], err = AtomName(xu, xproto.Atom(xgb.Get32(vals)))
		if err != nil {
			return nil, err
		}

		vals = vals[4:]
	}
	return ids, nil
}

// PropValWindow transforms a GetPropertyReply struct into an X resource
// window identifier.
// The property reply must be in 32 bit format.
func PropValWindow(reply *xproto.GetPropertyReply,
	err error) (xproto.Window, error) {

	if err != nil {
		return 0, err
	}
	if reply.Format != 32 {
		return 0, fmt.Errorf("PropValId: Expected format 32 but got %d",
			reply.Format)
	}
	return xproto.Window(xgb.Get32(reply.Value)), nil
}

// PropValWindows is the same as PropValWindow, except that it returns a slice
// of identifiers. Also must be 32 bit format.
func PropValWindows(reply *xproto.GetPropertyReply,
	err error) ([]xproto.Window, error) {

	if err != nil {
		return nil, err
	}
	if reply.Format != 32 {
		return nil, fmt.Errorf("PropValIds: Expected format 32 but got %d",
			reply.Format)
	}

	ids := make([]xproto.Window, reply.ValueLen)
	vals := reply.Value
	for i := 0; len(vals) >= 4; i++ {
		ids[i] = xproto.Window(xgb.Get32(vals))
		vals = vals[4:]
	}
	return ids, nil
}

// PropValNum transforms a GetPropertyReply struct into an unsigned
// integer. Useful when the property value is a single integer.
func PropValNum(reply *xproto.GetPropertyReply, err error) (uint, error) {
	if err != nil {
		return 0, err
	}
	if reply.Format != 32 {
		return 0, fmt.Errorf("PropValNum: Expected format 32 but got %d",
			reply.Format)
	}
	return uint(xgb.Get32(reply.Value)), nil
}

// PropValNums is the same as PropValNum, except that it returns a slice
// of integers. Also must be 32 bit format.
func PropValNums(reply *xproto.GetPropertyReply, err error) ([]uint, error) {
	if err != nil {
		return nil, err
	}
	if reply.Format != 32 {
		return nil, fmt.Errorf("PropValIds: Expected format 32 but got %d",
			reply.Format)
	}

	nums := make([]uint, reply.ValueLen)
	vals := reply.Value
	for i := 0; len(vals) >= 4; i++ {
		nums[i] = uint(xgb.Get32(vals))
		vals = vals[4:]
	}
	return nums, nil
}

// PropValNum64 transforms a GetPropertyReply struct into a 64 bit
// integer. Useful when the property value is a single integer.
func PropValNum64(reply *xproto.GetPropertyReply, err error) (int64, error) {
	if err != nil {
		return 0, err
	}
	if reply.Format != 32 {
		return 0, fmt.Errorf("PropValNum: Expected format 32 but got %d",
			reply.Format)
	}
	return int64(xgb.Get32(reply.Value)), nil
}

// PropValStr transforms a GetPropertyReply struct into a string.
// Useful when the property value is a null terminated string represented
// by integers. Also must be 8 bit format.
func PropValStr(reply *xproto.GetPropertyReply, err error) (string, error) {
	if err != nil {
		return "", err
	}
	if reply.Format != 8 {
		return "", fmt.Errorf("PropValStr: Expected format 8 but got %d",
			reply.Format)
	}
	return string(reply.Value), nil
}

// PropValStrs is the same as PropValStr, except that it returns a slice
// of strings. The raw byte string is a sequence of null terminated strings,
// which is translated into a slice of strings.
func PropValStrs(reply *xproto.GetPropertyReply, err error) ([]string, error) {
	if err != nil {
		return nil, err
	}
	if reply.Format != 8 {
		return nil, fmt.Errorf("PropValStrs: Expected format 8 but got %d",
			reply.Format)
	}

	var strs []string
	sstart := 0
	for i, c := range reply.Value {
		if c == 0 {
			strs = append(strs, string(reply.Value[sstart:i]))
			sstart = i + 1
		}
	}
	if sstart < int(reply.ValueLen) {
		strs = append(strs, string(reply.Value[sstart:]))
	}
	return strs, nil
}
