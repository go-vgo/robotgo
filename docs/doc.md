# Methods:

##### [GetVersion](#GetVersion)

## [Keyboard](#Keyboard)

##### [Keys](https://github.com/go-vgo/robotgo/blob/master/docs/keys.md)
##### [SetKeyboardDelay](#SetKeyDelay) (Equivalent to SetKeyDelay, Wno-deprecated)
##### [SetKeyDelay](#SetKeyDelay)
##### [KeyTap](#KeyTap)
##### [KeyToggle](#KeyToggle)
##### [TypeString](#TypeString)
##### [TypeStringDelayed](#TypeStrDelay) (Equivalent to TypeStrDelay, Wno-deprecated)
##### [TypeStrDelay](#TypeStrDelay)
##### [TypeStr](#TypeStr)
##### [WriteAll](#WriteAll)
##### [ReadAll](#ReadAll)


## [Mouse](#Mouse)

##### [SetMouseDelay](#SetMouseDelay)
##### [MoveMouse](#MoveMouse)
##### [Move](#MoveMouse) (Equivalent to MoveMouse)
##### [MoveMouseSmooth](#MoveMouseSmooth)
##### [MoveSmooth](#MoveMouseSmooth) (Equivalent to MoveMouseSmooth)
##### [MouseClick](#MouseClick)
##### [Click](#MouseClick) (Equivalent to MouseClick)
##### [MoveClick](#MoveClick)
##### [MouseToggle](#MouseToggle)
##### [DragMouse](#DragMouse)
##### [Drag](#DragMouse) (Equivalent to DragMouse)
##### [GetMousePos](#GetMousePos)
##### [ScrollMouse](#ScrollMouse)

## [Screen](#Screen)

##### [GetPixelColor](#GetPixelColor)
##### [GetScreenSize](#GetScreenSize)
##### [CaptureScreen](#CaptureScreen)
##### [GetXDisplayName(Linux)](#GetXDisplayName)
##### [SetXDisplayName(Linux)](#SetXDisplayName)

## [Bitmap](#Bitmap)
    This is a work in progress.

##### [FindBitmap](#FindBitmap)
##### [OpenBitmap](#OpenBitmap)
##### [SaveBitmap](#SaveBitmap)
##### [TostringBitmap](#TostringBitmap)
##### [GetPortion](#GetPortion)
##### [Convert](#Convert)
#### [FreeBitmap](#FreeBitmap)
#### [ReadBitmap](#ReadBitmap)
#### [CopyBitpb](#CopyBitpb)
#### [DeepCopyBit](#DeepCopyBit)

## [Event](#Event)

##### [LEvent](#LEvent) (Equivalent to AddEvent, Wno-deprecated)
##### [AddEvent](#AddEvent)
##### [StopEvent](#StopEvent)

## [Window](#Window)
    This is a work in progress.

##### [ShowAlert](#ShowAlert)
##### [CloseWindow](#CloseWindow)
##### [IsValid](#IsValid)
##### [SetActive](#SetActive)
##### [GetActive](#GetActive)
##### [SetHandle](#SetHandle)
##### [GetHandle](#GetHandle)
##### [GetBHandle](#GetHandle)
##### [GetTitle](#GetTitle)
##### [GetPID](#GetPID)
##### [Pids](#Pids)
##### [PidExists](#PidExists)
##### [Process](#Process)
##### [FindName](#FindName)
##### [FindNames](#FindNames)
##### [FindIds](#FindIds)
#### [ActivePID](#ActivePID)

### <h3 id="GetVersion">.GetVersion()</h3>
    Get robotgo version

## <h2 id="Keyboard">Keyboard</h2>

### <h3 id="SetKeyDelay">.SetKeyDelay(ms)</h3>

    Sets the delay in milliseconds to sleep after a keyboard event. This is 10ms by default.

#### Arguments:

    ms - Time to sleep in milliseconds.

### <h3 id="KeyTap">.KeyTap(key, modifier)</h3>

    Press a single key.

#### Arguments:

key - See [keys](https://github.com/go-vgo/robotgo/blob/master/docs/keys.md).  
modifier (optional, string or array) - Accepts alt, command (win), control, and shift.

#### Examples:

```Go
robotgo.KeyTap("h", "command")
robotgo.KeyTap("i", "alt", "command")
arr := []string{"alt", "command"}
robotgo.KeyTap("i", arr)
```

### <h3 id="KeyToggle">.KeyToggle(key, down, modifier)</h3>

    Hold down or release a key.

#### Arguments:

key - See [keys](https://github.com/go-vgo/robotgo/blob/master/docs/keys.md).  
down - Accepts 'down' or 'up'.  
modifier (optional, string or array) - Accepts alt, command (mac), control, and shift.

#### Return:
    return KeyToggle status

### <h3 id="TypeString">.TypeString(string)</h3>

#### Arguments:

    string - The string to send.

### <h3 id="TypeStrDelay">.TypeStrDelay(string, cpm)</h3>

#### Arguments:

    string - The string to send.
    cpm - Characters per minute.

### <h3 id="TypeStr">.TypeStr(string)</h3>

#### Arguments:

    string - The string to send.

### <h3 id="WriteAll">.WriteAll(text string)</h3>

#### Arguments:
    text string
#### Return:      

### <h3 id="ReadAll">.ReadAll()</h3>

#### Arguments:

#### Return: 
    text,
    error 


## <h2 id="Mouse">Mouse</h2>

### <h3 id="SetMouseDelay">.SetMouseDelay(ms)</h3>

    Sets the delay in milliseconds to sleep after a mouse event. This is 10ms by default.

#### Arguments:

    ms - Time to sleep in milliseconds.

### <h3 id="MoveMouse">.MoveMouse(x, y)</h3>

    Moves mouse to x, y instantly, with the mouse button up.

#### Arguments:

    x, y

#### Examples:

```Go
// Move the mouse to 100, 100 on the screen. 
robotgo.MoveMouse(100, 100)
```

### <h3 id="MoveMouseSmooth">.MoveMouseSmooth(x, y)</h3>

    Moves mouse to x, y human like, with the mouse button up.

#### Arguments:

    x, y
    lowspeed, highspeed

#### Examples:

```Go
robotgo.MoveMouseSmooth(100, 200)
robotgo.MoveMouseSmooth(100, 200, 1.0, 100.0)
```    

### <h3 id="MouseClick">.MouseClick(button, double)</h3>

    Clicks the mouse.

#### Arguments:

    button (optional) - Accepts "left", "right", or "center". Defaults to left.
    double (optional) - Set to true to perform a double click. Defaults to false.

#### Examples:

```Go
robogo.MouseClick()
robogo.MouseClick("left", true)
```

### <h3 id="MoveClick">.MoveClick(x, y, button, double)</h3>

    Move and click the mouse.

#### Arguments:
    x,
    y,

    button (optional) - Accepts "left", "right", or "center". Defaults to left.
    double (optional) - Set to true to perform a double click. Defaults to false.

#### Examples:

```Go
robogo.MoveClick(10, 20)
robogo.MoveClick(10, 20, "left", true)
```

### <h3 id="MouseToggle">.MouseToggle(down, button)</h3>

    Toggles mouse button.

#### Arguments:

    down (optional) - Accepts down or up. Defaults to down.
    button (optional) - Accepts "left", "right", or "center". Defaults to left.

#### Examples:

```Go
robotgo.MouseToggle("down")
robotgo.MouseToggle("down", "right")
```

### <h3 id="DragMouse">.DragMouse(x, y)</h3>

    Moves mouse to x, y instantly, with the mouse button held down.

#### Arguments:

    x, y

#### Examples:

```Go
// Mouse down at 0, 0 and then drag to 100, 100 and release. 
robotgo.MoveMouse(0, 0)
robotgo.MouseToggle("down")
robotgo.DragMouse(100, 100)
robotgo.MouseToggle("up")
```

### <h3 id="GetMousePos">.GetMousePos()</h3>

    Gets the mouse coordinates.

#### Return:

    Returns an object with keys x and y.

#### Examples:

```Go
x,y := robotgo.GetMousePos()
fmt.Println("pos:", x, y)
```

### <h3 id="ScrollMouse">.ScrollMouse(magnitude, direction)</h3>

    Scrolls the mouse either up or down.

#### Arguments:

    magnitude - The amount to scroll.
    direction - Accepts down or up.

#### Examples:

```Go
robotgo.ScrollMouse(50, "up")

robotgo.ScrollMouse(50, "down")
```


## <h2 id="Screen">Screen</h2>

### <h3 id="GetPixelColor">.GetPixelColor(x, y)

    Gets the pixel color at x, y. This function is perfect for getting a pixel or two, 
    but if you'll be reading large portions of the screen use screen.capture.

#### Arguments:

    x, y

#### Return:

    Returns the hex color code of the pixel at x, y.

### <h3 id="GetScreenSize">.GetScreenSize()</h3>

    Gets the screen width and height.

#### Return:

    Returns an object with .width and .height.

### <h3 id="CaptureScreen">.CaptureScreen</h3>
    // CaptureScreen

    Gets part or all of the screen.

    GoCaptureScreen Returns a go struct
    Capture_Screen (Drop support)

#### Arguments:

    x (optional)
    y (optional)
    height (optional)
    width (optional)
    If no arguments are provided, capturescreen(screencapture) will get the full screen.

#### Return:

    Returns a bitmap object.

## <h2 id="Bitmap">Bitmap</h2>

    This is a work in progress.

### <h3 id="FindBitmap">.FindBitmap</h3>

    find bitmap.

#### Arguments:

    bitmap;
    rect(optional): x, y, w, h

#### Return:

    Returns a position x and y


### <h3 id="OpenBitmap">.OpenBitmap</h3>

    open bitmap.

#### Arguments:

    bitmap image path,
    MMImageType (optional) 

#### Return:

    Returns a bitmap

### <h3 id="SaveBitmap">.SaveBitmap</h3>

    Save a image with bitmap.

#### Arguments:

    bitmap,
    path,
    imagetype (int) 

#### Return:

    return save image status


### <h3 id="TostringBitmap">.TostringBitmap</h3>

     bitmap to string

#### Arguments:

    bitmap 

#### Return:

    Return a sting bitmap

### <h3 id="GetPortion">.GetPortion</h3>

     bitmap from a portion

#### Arguments:

    bitmap,
    rect: x, y, w, h 

#### Return:

    Returns new bitmap object created from a portion of another 

### <h3 id="Convert">.Convert(openpath, savepath, MMImageType)</h3>

    Convert the image format

#### Arguments:

    openpath,
    savepath,
    MMImageType (optional)

#### Examples:

```Go
robotgo.Convert("test.png", "test.tif")
```             

### <h3 id="FreeBitmap">.FreeBitmap(MMBitmapRef)</h3>

    FreeBitmap free and dealloc bitmap

#### Arguments:

    MMBitmapRef

#### Examples:

```Go
robotgo.FreeBitmap(bitmap)
```    


### <h3 id="ReadBitmap">.ReadBitmap(MMBitmapRef)</h3>

    ReadBitmap returns false and sets error if |bitmap| is NULL

#### Arguments:

    MMBitmapRef

#### Return:

    bool

#### Examples:

```Go
robotgo.ReadBitmap(bitmap)
```    


### <h3 id="CopyBitpb">.CopyBitpb(MMBitmapRef)</h3>

   CopyBitpb copy bitmap to pasteboard

#### Arguments:

    MMBitmapRef

#### Return:

    bool

#### Examples:

```Go
robotgo.CopyBitpb(bitmap)
```    

### <h3 id="DeepCopyBit">.DeepCopyBit(MMBitmapRef)</h3>

   DeepCopyBit deep copy bitmap

#### Arguments:

    MMBitmapRef

#### Return:

    MMBitmapRef

#### Examples:

```Go
robotgo.DeepCopyBit(bitmap)
```  

## <h2 id="Event">Event</h2> 

### <h3 id="AddEvent">.AddEvent(string)</h3>

    Listening global event

#### Arguments:

    string

    (mouse arguments: mleft, mright, wheelDown, wheelUp, wheelLeft, wheelRight)

#### Return:

    if listened return 0

#### Examples:

```Go
package main

import (
    "fmt"

    "github.com/go-vgo/robotgo"
)

func main() {
  keve := robotgo.AddEvent("k")
  if keve {
    fmt.Println("you press...", "k")
  }

  mleft := robotgo.AddEvent("mleft")
  if mleft {
    fmt.Println("you press...", "mouse left button")
  }
} 
```  
### <h3 id="StopEvent">.StopEvent()</h3>  
    stop listen global event

## <h2 id="Window">Window</h2> 

### <h3 id="ShowAlert">.ShowAlert(title, msg, defaultButton, cancelButton string)</h3>

    Displays alert with the given attributes. If cancelButton is not given, only the defaultButton is displayed

#### Arguments:

    title (string),
    msg (string),
    defaultButton (optional string),
    cancelButton (optional string)
           

#### Return:

    Returns 0(True) if the default button was pressed, or 1(False) if cancelled. 

### <h3 id="CloseWindow">.CloseWindow()</h3>

    Close the Window

#### Arguments:
    None     

#### Return:
    None

### <h3 id="IsValid">.IsValid()</h3>

    Valid the Window

#### Arguments:
    None     

#### Return:
    Returns true if a window has been selected


### <h3 id="SetActive">.SetActive()</h3>

   Set the Active Window

#### Arguments:
    hwnd  

#### Return:
    void
    

### <h3 id="GetActive">.GetActive()</h3>

    Get the Active Window

#### Arguments:
    None   

#### Return:
    Returns hwnd

### <h3 id="SetHandle">.SetHandle()</h3>

    Set the Window Handle

#### Arguments:
    int 

#### Return:
    bool

### <h3 id="GetHandle">.GetHandle()</h3>

    Get the Window Handle

#### Arguments:
    None    

#### Return:
    Returns hwnd  

### <h3 id="GetTitle">.GetTitle()</h3>

    Get the Window Title

#### Arguments:
    None   

#### Return:
    Returns Window Title      

### <h3 id="GetPID">.GetPID()</h3>

    Get the process id

#### Arguments:
    None       

#### Return:
    Returns the process id         
     
### <h3 id="Pids">.Pids()</h3>

    Pids get the all process id

#### Arguments:
    None   

#### Return:
    Returns all process id   

### <h3 id="PidExists">.PidExists()</h3>

    PidExists determine whether the process exists

#### Arguments:
    pid  

#### Return:
    Returns bool 

### <h3 id="Process">.Process()</h3>

  Process get the all process struct

#### Arguments:
    none  

#### Return:
    Returns []Nps, error

### <h3 id="FindName">.FindName()</h3>

    FindName find the process name by the process id

#### Arguments:
    pid  

#### Return:
    Returns string, error       

### <h3 id="FindNames">.FindNames()</h3>

    FindNames find the all process name

#### Arguments:
    none  

#### Return:
    Returns []string, error  

### <h3 id="FindIds">.FindIds()</h3>

    FindIds find the process id by the process name

#### Arguments:
    name string  

#### Return:
    Returns []int32, error 


### <h3 id="ActivePID">.ActivePID()</h3>

    ActivePID window active by PID

#### Arguments:
    pid int32 

#### Return:
    none                  