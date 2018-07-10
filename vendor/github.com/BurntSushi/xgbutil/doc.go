/*
Package xgbutil is a utility library designed to make common tasks with the X
server easier. The central design choice that has driven development is to hide
the complexity of X wherever possible but expose it when necessary.

For example, the xevent package provides an implementation of an X event loop
that acts as a dispatcher to event handlers set up with the xevent, keybind and
mousebind packages. At the same time, the event queue is exposed and can be
modified using xevent.Peek and xevent.DequeueAt.

Sub-packages

The xgbutil package is considerably small, and only contains some type
definitions and the initial setup for an X connection. Much of the
functionality of xgbutil comes from its sub-packages. Each sub-package is
appropriately documented.

Installation

xgbutil is go-gettable:

	go get github.com/BurntSushi/xgbutil

Dependencies

XGB is the main dependency, and is required for all packages inside xgbutil.

graphics-go and freetype-go are also required if using the xgraphics package.

Quick Example

A quick example to demonstrate that xgbutil is working correctly:

	go get github.com/BurntSushi/xgbutil/examples/window-name-sizes
	GO/PATH/bin/window-name-sizes

The output will be a list of names of all top-level windows and their geometry
including window manager decorations. (Assuming your window manager supports
some basic EWMH properties.)

Examples

The examples directory contains a sizable number of examples demonstrating
common tasks with X. They are intended to demonstrate a single thing each,
although a few that require setup are necessarily long. Each example is
heavily documented.

The examples directory should be your first stop when learning how to use
xgbutil.

xgbutil is also used heavily throughout my window manager, Wingo. It may be
useful reference material.

Wingo project page: https://github.com/BurntSushi/wingo

Thread Safety

While I am fairly confident that XGB is thread safe, I am only somewhat
confident that xgbutil is thread safe. It simply has not been tested enough for
my confidence to be higher.

Note that the xevent package's X event loop is not concurrent. Namely,
designing a generally concurrent X event loop is extremely complex. Instead,
the onus is on you, the user, to design concurrent callback functions if
concurrency is desired.
*/
package xgbutil
