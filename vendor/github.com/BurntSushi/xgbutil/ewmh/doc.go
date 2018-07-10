/*
Package ewmh provides a comprehensive API to get and set properties specified
by the EWMH spec, as well as perform actions specified by the EWMH spec.

Since there are so many functions and they adhere to an existing spec,
this package file does not contain much documentation. Indeed, each
method has only a single comment associated with it: the EWMH property name.

The idea is to
provide a consistent interface to use all facilities described in the EWMH
spec: http://standards.freedesktop.org/wm-spec/wm-spec-latest.html.

Naming scheme

Using "_NET_ACTIVE_WINDOW" as an example,
functions "ActiveWindowGet" and "ActiveWindowSet" get and set the
property, respectively. Both of these functions exist for most EWMH
properties. Additionally, some EWMH properties support sending a client
message event to request the window manager to perform some action. In the
case of "_NET_ACTIVE_WINDOW", this request is used to set the active
window.

These sorts of functions end in "Req". So for "_NET_ACTIVE_WINDOW",
the method name is "ActiveWindowReq". Moreover, most requests include
various parameters that don't need to be changed often (like the source
indication). Thus, by default, functions ending in "Req" force these to
sensible defaults. If you need access to all of the parameters, use the
corresponding "ReqExtra" method. So for "_NET_ACTIVE_WINDOW", that would
be "ActiveWindowReqExtra". (If no "ReqExtra" method exists, then the
"Req" method covers all available parameters.)

This naming scheme has one exception: if a property's only use is through
sending an event (like "_NET_CLOSE_WINDOW"), then the name will be
"CloseWindow" for the short-hand version and "CloseWindowExtra"
for access to all of the parameters. (Since there is no "_NET_CLOSE_WINDOW"
property, there is no need for "CloseWindowGet" and "CloseWindowSet"
functions.)

For properties that store more than just a simple integer, name or list
of integers, structs have been created and exposed to organize the
information returned in a sensible manner. For example, the
"_NET_DESKTOP_GEOMETRY" property would typically return a slice of integers
of length 2, where the first integer is the width and the second is the
height. Xgbutil will wrap this in a struct with the obvious members. These
structs are documented.

Finally, functions ending in "*Set" are typically only used when setting
properties on clients *you've* created or when the window manager sets
properties. Thus, it's unlikely that you should use them unless you're
creating a top-level client or building a window manager.

Functions ending in "Get" or "Req[Extra]" are commonly used.

N.B. Not all properties have "*Req" functions.
*/
package ewmh
