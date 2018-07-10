/*
Package xprop provides a cache for interning atoms and helper functions for
dealing with GetProperty and ChangeProperty X requests.

Atoms

Atoms in X are interned, meaning that strings are assigned unique integer
identifiers. This minimizes the amount of data transmitted over an X connection.

Once atoms have been interned, they are never changed while the X server is
running. xgbutil takes advantage of this invariant and will only issue an
intern atom request once and cache the result.

To use the xprop package to intern an atom, use Atom:

	atom, err := xprop.Atom(XUtilValue, "THE_ATOM_NAME", false)
	if err == nil {
		println("The atom number: ", atom.Atom)
	}

The 'false' parameter corresponds to the 'only_if_exists' parameter of the
X InternAtom request. When it's false, the atom being interned always returns
a non-zero atom number---even if the string being interned hasn't been interned
before. If 'only_if_exists' is true, the atom being interned will return a 0
atom number if it hasn't already been interned.

The typical case is to set 'only_if_exists' to false. To this end, xprop.Atm is
an alias that always sets this value to false.

The reverse can also be done: getting an atom string if you have an atom
number. This can be done with the xprop.AtomName function.

Properties

The other facility of xprop is to help with the use of GetProperty and
ChangeProperty. Please see the source code of the ewmh package for plenty of
examples.

*/
package xprop
