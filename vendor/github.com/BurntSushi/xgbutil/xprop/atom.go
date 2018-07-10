package xprop

/*
xprop/atom.go contains functions related to interning atoms and retrieving
atom names from an atom identifier.

It also manages an atom cache so that once an atom is interned from the X
server, all future atom interns use that value. (So that one and only one
request is sent for interning each atom.)
*/

import (
	"fmt"

	"github.com/BurntSushi/xgb/xproto"

	"github.com/BurntSushi/xgbutil"
)

// Atm is a short alias for Atom in the common case of interning an atom.
// Namely, interning the atom always succeeds. (If the atom does not already
// exist, a new one is created.)
func Atm(xu *xgbutil.XUtil, name string) (xproto.Atom, error) {
	aid, err := Atom(xu, name, false)
	if err != nil {
		return 0, err
	}
	if aid == 0 {
		return 0, fmt.Errorf("Atm: '%s' returned an identifier of 0.", name)
	}

	return aid, err
}

// Atom interns an atom and panics if there is any error.
func Atom(xu *xgbutil.XUtil, name string,
	onlyIfExists bool) (xproto.Atom, error) {

	// Check the cache first
	if aid, ok := atomGet(xu, name); ok {
		return aid, nil
	}

	reply, err := xproto.InternAtom(xu.Conn(), onlyIfExists,
		uint16(len(name)), name).Reply()
	if err != nil {
		return 0, fmt.Errorf("Atom: Error interning atom '%s': %s", name, err)
	}

	// If we're here, it means we didn't have this atom cached. So cache it!
	cacheAtom(xu, name, reply.Atom)

	return reply.Atom, nil
}

// AtomName fetches a string representation of an ATOM given its integer id.
func AtomName(xu *xgbutil.XUtil, aid xproto.Atom) (string, error) {
	// Check the cache first
	if atomName, ok := atomNameGet(xu, aid); ok {
		return string(atomName), nil
	}

	reply, err := xproto.GetAtomName(xu.Conn(), aid).Reply()
	if err != nil {
		return "", fmt.Errorf("AtomName: Error fetching name for ATOM "+
			"id '%d': %s", aid, err)
	}

	// If we're here, it means we didn't have ths ATOM id cached. So cache it.
	atomName := string(reply.Name)
	cacheAtom(xu, atomName, aid)

	return atomName, nil
}

// atomGet retrieves an atom identifier from a cache if it exists.
func atomGet(xu *xgbutil.XUtil, name string) (xproto.Atom, bool) {
	xu.AtomsLck.RLock()
	defer xu.AtomsLck.RUnlock()

	aid, ok := xu.Atoms[name]
	return aid, ok
}

// atomNameGet retrieves an atom name from a cache if it exists.
func atomNameGet(xu *xgbutil.XUtil, aid xproto.Atom) (string, bool) {
	xu.AtomNamesLck.RLock()
	defer xu.AtomNamesLck.RUnlock()

	name, ok := xu.AtomNames[aid]
	return name, ok
}

// cacheAtom puts an atom into the cache.
func cacheAtom(xu *xgbutil.XUtil, name string, aid xproto.Atom) {
	xu.AtomsLck.Lock()
	xu.AtomNamesLck.Lock()
	defer xu.AtomsLck.Unlock()
	defer xu.AtomNamesLck.Unlock()

	xu.Atoms[name] = aid
	xu.AtomNames[aid] = name
}
