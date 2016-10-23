#Methods:

##Keyboard
    Keys
    SetKeyboardDelay
    KeyTap
    KeyToggle
    TypeString
    TypeStringDelayed
##Mouse
    SetMouseDelay
    MoveMouse
    MoveMouseSmooth
    MouseClick
    MouseToggle
    DragMouse
    GetMousePos
    ScrollMouse
##Screen
    GetPixelColor
    GetScreenSize
    CaptureScreen
    GetXDisplayName(Linux)
    SetXDisplayName(Linux)
##Bitmap
    This is a work in progress.


##Keyboard
###.SetKeyboardDelay(ms)

    Sets the delay in milliseconds to sleep after a keyboard event. This is 10ms by default.

####Arguments:

    ms - Time to sleep in milliseconds.

###.KeyTap(key, modifier)

    Press a single key.

####Arguments:

    key - See keys.
    modifier (optional, string or array) - Accepts alt, command (win), control, and shift.

###.KeyToggle(key, down, modifier)

    Hold down or release a key.

####Arguments:

    key - See keys.
    down - Accepts 'down' or 'up'.
    modifier (optional, string or array) - Accepts alt, command (mac), control, and shift.

###.TypeString(string)

####Arguments:

    string - The string to send.

###.TypeStringDelayed(string, cpm)

####Arguments:

    string - The string to send.
    cpm - Characters per minute.



##Mouse
###.SetMouseDelay(ms)

    Sets the delay in milliseconds to sleep after a mouse event. This is 10ms by default.

####Arguments:

    ms - Time to sleep in milliseconds.

###.MoveMouse(x, y)

    Moves mouse to x, y instantly, with the mouse button up.

####Arguments:

    x,y

####Examples:

```Go
//Move the mouse to 100, 100 on the screen. 
robotgo.MoveMouse(100, 100)
```

###.MoveMouseSmooth(x, y)

    Moves mouse to x, y human like, with the mouse button up.

####Arguments:

    x,y

###.MouseClick(button, double)

    Clicks the mouse.

####Arguments:

    button (optional) - Accepts left, right, or middle. Defaults to left.
    double (optional) - Set to true to perform a double click. Defaults to false.

####Examples:

```Go
    robogo.MouseClick()
```

###.mouseToggle(down, button)

    Toggles mouse button.

####Arguments:

    down (optional) - Accepts down or up. Defaults to down.
    button (optional) - Accepts left, right, or middle. Defaults to left.

####Examples:

```Go
robotgo.MouseToggle("down")
```

###.DragMouse(x, y)

    Moves mouse to x, y instantly, with the mouse button held down.

####Arguments:

    x,y

####Examples:

```Go
//Mouse down at 0, 0 and then drag to 100, 100 and release. 
robotgo.MoveMouse(0, 0)
robotgo.MouseToggle("down")
robotgo.DragMouse(100, 100)
robotgo.MouseToggle("up")
```

###.GetMousePos()

    Gets the mouse coordinates.

####Return:

    Returns an object with keys x and y.

####Examples:

```Go
x,y := robotgo.GetMousePos()
fmt.Println("pos:", x, y)
```

###.ScrollMouse(magnitude, direction)

    Scrolls the mouse either up or down.

####Arguments:

    magnitude - The amount to scroll.
    direction - Accepts down or up.

####Examples:

```Go
robotgo.ScrollMouse(50, "up")

robotgo.ScrollMouse(50, "down")
```


##Screen
###.GetPixelColor(x, y)

    Gets the pixel color at x, y. This function is perfect for getting a pixel or two, but if you'll be reading large portions of the screen use screen.capture.

####Arguments:

    x,y

####Return:

    Returns the hex color code of the pixel at x, y.

###.GetScreenSize()

    Gets the screen width and height.

####Return:

    Returns an object with .width and .height.

###.CaptureScreen
    //ScreenCapture

    Gets part or all of the screen.

####Arguments:

    x (optional)
    y (optional)
    height (optional)
    width (optional)
    If no arguments are provided, screencapture will get the full screen.

####Return:

    Returns a bitmap object.

##Bitmap

    This is a work in progress.
