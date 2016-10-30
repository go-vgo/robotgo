#Methods:

##[Keyboard](#Keyboard)

#####[Keys](#keys)
#####[SetKeyboardDelay](#SetKeyboardDelay)
#####[KeyTap](#KeyTap)
#####[KeyToggle](#KeyToggle)
#####[TypeString](#TypeString)
#####[TypeStringDelayed](#TypeStringDelayed)

##[Mouse](#Mouse)

#####[SetMouseDelay](#SetMouseDelay)
#####[MoveMouse](#MoveMouse)
#####[MoveMouseSmooth](#MoveMouseSmooth)
#####[MouseClick](#MouseClick)
#####[MouseToggle](#MouseToggle)
#####[DragMouse](#DragMouse)
#####[GetMousePos](#GetMousePos)
#####[ScrollMouse](#ScrollMouse)

##[Screen](#Screen)

#####[GetPixelColor](#GetPixelColor)
#####[GetScreenSize](#GetScreenSize)
#####[CaptureScreen](#CaptureScreen)
#####[GetXDisplayName(Linux)](#GetXDisplayName)
#####[SetXDisplayName(Linux)](#SetXDisplayName)

##[Bitmap](#Bitmap)
    This is a work in progress.

#####[FindBitmap](#FindBitmap)
#####[OpenBitmap](#OpenBitmap)
#####[SaveBitmap](#SaveBitmap)
#####[TostringBitmap](#TostringBitmap)
#####[GetPortion](#GetPortion)
#####[Convert](#Convert)

##[Event](#Event)

#####[LEvent](#LEvent)


##<h2 id="Keyboard">Keyboard</h2>

###<h3 id="SetKeyboardDelay">.SetKeyboardDelay(ms)</h3>

    Sets the delay in milliseconds to sleep after a keyboard event. This is 10ms by default.

####Arguments:

    ms - Time to sleep in milliseconds.

###<h3 id="KeyTap">.KeyTap(key, modifier)</h3>

    Press a single key.

####Arguments:

    key - See keys.
    modifier (optional, string or array) - Accepts alt, command (win), control, and shift.

###<h3 id="KeyToggle">.KeyToggle(key, down, modifier)</h3>

    Hold down or release a key.

####Arguments:

    key - See keys.
    down - Accepts 'down' or 'up'.
    modifier (optional, string or array) - Accepts alt, command (mac), control, and shift.

###<h3 id="TypeString">.TypeString(string)</h3>

####Arguments:

    string - The string to send.

###<h3 id="TypeStringDelayed">.TypeStringDelayed(string, cpm)</h3>

####Arguments:

    string - The string to send.
    cpm - Characters per minute.



##<h2 id="Mouse">Mouse</h2>

###<h3 id="SetMouseDelay">.SetMouseDelay(ms)</h3>

    Sets the delay in milliseconds to sleep after a mouse event. This is 10ms by default.

####Arguments:

    ms - Time to sleep in milliseconds.

###<h3 id="MoveMouse">.MoveMouse(x, y)</h3>

    Moves mouse to x, y instantly, with the mouse button up.

####Arguments:

    x,y

####Examples:

```Go
//Move the mouse to 100, 100 on the screen. 
robotgo.MoveMouse(100, 100)
```

###<h3 id="MoveMouseSmooth">.MoveMouseSmooth(x, y)</h3>

    Moves mouse to x, y human like, with the mouse button up.

####Arguments:

    x,y

###<h3 id="MouseClick">.MouseClick(button, double)</h3>

    Clicks the mouse.

####Arguments:

    button (optional) - Accepts left, right, or middle. Defaults to left.
    double (optional) - Set to true to perform a double click. Defaults to false.

####Examples:

```Go
    robogo.MouseClick()
```

###<h3 id="MouseToggle">.MouseToggle(down, button)</h3>

    Toggles mouse button.

####Arguments:

    down (optional) - Accepts down or up. Defaults to down.
    button (optional) - Accepts left, right, or middle. Defaults to left.

####Examples:

```Go
robotgo.MouseToggle("down")
```

###<h3 id="DragMouse">.DragMouse(x, y)</h3>

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

###<h3 id="GetMousePos">.GetMousePos()</h3>

    Gets the mouse coordinates.

####Return:

    Returns an object with keys x and y.

####Examples:

```Go
x,y := robotgo.GetMousePos()
fmt.Println("pos:", x, y)
```

###<h3 id="ScrollMouse">.ScrollMouse(magnitude, direction)</h3>

    Scrolls the mouse either up or down.

####Arguments:

    magnitude - The amount to scroll.
    direction - Accepts down or up.

####Examples:

```Go
robotgo.ScrollMouse(50, "up")

robotgo.ScrollMouse(50, "down")
```


##<h2 id="Screen">Screen</h2>

###<h3 id="GetPixelColor">.GetPixelColor(x, y)

    Gets the pixel color at x, y. This function is perfect for getting a pixel or two, but if you'll be reading large portions of the screen use screen.capture.

####Arguments:

    x,y

####Return:

    Returns the hex color code of the pixel at x, y.

###<h3 id="GetScreenSize">.GetScreenSize()</h3>

    Gets the screen width and height.

####Return:

    Returns an object with .width and .height.

###<h3 id="CaptureScreen">.CaptureScreen</h3>
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

##<h2 id="Bitmap">Bitmap</h2>

    This is a work in progress.

###<h3 id="FindBitmap">.FindBitmap</h3>

    find bitmap.

####Arguments:

    bitmap;
    rect(optional): x, y, w, h

####Return:

    Returns a position x and y


###<h3 id="OpenBitmap">.OpenBitmap</h3>

    open bitmap .

####Arguments:

    bitmap image path,
    MMImageType(optional) 

####Return:

    Returns a bitmap

###<h3 id="SaveBitmap">.SaveBitmap</h3>

    save a image with bitmap.

####Arguments:

    bitmap,
    path,
    imagetype(int) 

####Return:

    Return a imgage


###<h3 id="TostringBitmap">.TostringBitmap</h3>

     bitmap to string

####Arguments:

    bitmap 

####Return:

    Return a sting bitmap

###<h3 id="GetPortion">.GetPortion</h3>

     bitmap from a portion

####Arguments:

    bitmap,
    rect: x, y, w, h 

####Return:

    Returns new bitmap object created from a portion of another 

###<h3 id="Convert">.Convert(openpath, savepath,MMImageType)</h3>

    Convert the image format

####Arguments:

    openpath,
    savepath,
    MMImageType(optional)

####Examples:

```Go
    robotgo.Convert("test.png", "test.tif")
```             

##<h2 id="Event">Event</h2> 

###<h3 id="LEvent">.LEvent(string)</h3>

    Listening global event

####Arguments:

    string

####Return:

   if listened return 0

####Examples:

```Go
package main

import (
    . "fmt"

    "github.com/go-vgo/robotgo"
)

func main() {
  keve := robotgo.LEvent("k")
  if keve == 0 {
    Println("you press...", "k")
  }

  mleft := robotgo.LEvent("mleft")
  if mleft == 0 {
    Println("you press...", "mouse left button")
  }
} 
```    
