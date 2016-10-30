#方法:

##[键盘](#Keyboard)

#####[Keys](#keys)
#####[SetKeyboardDelay](#SetKeyboardDelay)
#####[KeyTap](#KeyTap)
#####[KeyToggle](#KeyToggle)
#####[TypeString](#TypeString)
#####[TypeStringDelayed](#TypeStringDelayed)

##[鼠标](#Mouse)

#####[SetMouseDelay](#SetMouseDelay)
#####[MoveMouse](#MoveMouse)
#####[MoveMouseSmooth](#MoveMouseSmooth)
#####[MouseClick](#MouseClick)
#####[MouseToggle](#MouseToggle)
#####[DragMouse](#DragMouse)
#####[GetMousePos](#GetMousePos)
#####[ScrollMouse](#ScrollMouse)

##[屏幕](#Screen)

#####[GetPixelColor](#GetPixelColor)
#####[GetScreenSize](#GetScreenSize)
#####[CaptureScreen](#CaptureScreen)
#####[GetXDisplayName(Linux)](#GetXDisplayName)
#####[SetXDisplayName(Linux)](#SetXDisplayName)

##[位图](#Bitmap)
    This is a work in progress.(工作正在进行中)

#####[FindBitmap](#FindBitmap)
#####[OpenBitmap](#OpenBitmap)
#####[SaveBitmap](#SaveBitmap)
#####[TostringBitmap](#TostringBitmap)
#####[GetPortion](#GetPortion)
#####[Convert](#Convert)

##[事件](#Event)

#####[LEvent](#LEvent)


##<h2 id="Keyboard">键盘</h2>
###<h3 id="SetKeyboardDelay">.SetKeyboardDelay(ms)</h3>
    设置键盘延迟(在键盘一个事件后),单位ms,默认值10ms

    Sets the delay in milliseconds to sleep after a keyboard event. This is 10ms by default.

####参数:
    延迟时间,单位ms

    ms - Time to sleep in milliseconds.

###<h3 id="KeyTap">.KeyTap(key, modifier)</h3>
    模拟键盘按键

    Press a single key.

####参数:
    键盘值
    修饰值(可选类型, 字符串或者数组(数组类型正在添加中)) - 可选值: alt, command (win), control, and shift.

    key - See keys.
    modifier (optional, string or array) - Accepts alt, command (win), control, and shift.

###<h3 id="KeyToggle">.KeyToggle(key, down, modifier)</h3>
    键盘切换,按住或释放一个键位

    Hold down or release a key.

####参数:

    key - See keys.
    down - Accepts 'down' or 'up'.
    modifier (optional, string or array) - Accepts alt, command (mac), control, and shift.

###<h3 id="TypeString">.TypeString(string)</h3>

####参数:

    string - The string to send.

###<h3 id="TypeStringDelayed">.TypeStringDelayed(string, cpm)</h3>

####参数:

    string - The string to send.
    cpm - Characters per minute.



##<h2 id="Mouse">鼠标</h2>
###<h3 id="SetMouseDelay">.SetMouseDelay(ms)</h3>
    设置鼠标延迟(在一个鼠标事件后),单位ms,默认值10ms

    Sets the delay in milliseconds to sleep after a mouse event. This is 10ms by default.

####参数:

    ms - Time to sleep in milliseconds.

###<h3 id="MoveMouse">.MoveMouse(x, y)</h3>
    移动鼠标

    Moves mouse to x, y instantly, with the mouse button up.

####参数:

    x,y

####例子:

```Go
//Move the mouse to 100, 100 on the screen. 
robotgo.MoveMouse(100, 100)
```

###<h3 id="MoveMouseSmooth">.MoveMouseSmooth(x, y)</h3>
    模拟鼠标向X，Y平滑移动(像人类一样)，用鼠标按钮向上

    Moves mouse to x, y human like, with the mouse button up.

####参数:

    x,y

###<h3 id="MouseClick">.MouseClick(button, double)</h3>
    鼠标点击

    Clicks the mouse.

####参数:

    button (optional) - Accepts left, right, or middle. Defaults to left.
    double (optional) - Set to true to perform a double click. Defaults to false.

####例子:

```Go
    robogo.MouseClick()
```

###<h3 id="MouseToggle">.MouseToggle(down, button)</h3>
    鼠标切换

    Toggles mouse button.

####参数:

    down (optional) - Accepts down or up. Defaults to down.
    button (optional) - Accepts left, right, or middle. Defaults to left.

####例子:

```Go
robotgo.MouseToggle("down")
```

###<h3 id="DragMouse">.DragMouse(x, y)</h3>
    拖动鼠标

    Moves mouse to x, y instantly, with the mouse button held down.

####参数:

    x,y

####例子:

```Go
//Mouse down at 0, 0 and then drag to 100, 100 and release. 
robotgo.MoveMouse(0, 0)
robotgo.MouseToggle("down")
robotgo.DragMouse(100, 100)
robotgo.MouseToggle("up")
```

###<h3 id="GetMousePos">.GetMousePos()</h3>
    获取鼠标的位置

    Gets the mouse coordinates.

####返回值:

    Returns an object with keys x and y.

####例子:

```Go
x,y := robotgo.GetMousePos()
fmt.Println("pos:", x, y)
```

###<h3 id="ScrollMouse">.ScrollMouse(magnitude, direction)</h3>
    滚动鼠标

    Scrolls the mouse either up or down.

####参数:
    滚动位置的大小
    滚动方向:up(向上滚动)  down(向下滚动)

    magnitude - The amount to scroll.
    direction - Accepts down or up.

####例子:

```Go
robotgo.ScrollMouse(50, "up")

robotgo.ScrollMouse(50, "down")
```


##<h2 id="Screen">屏幕</h2>
###<h3 id="GetPixelColor">.GetPixelColor(x, y)
    获取坐标为x,y位置处的颜色

    Gets the pixel color at x, y. This function is perfect for getting a pixel or two, but if you'll be reading large portions of the screen use screen.capture.

####参数:

    x,y

####返回值:

    Returns the hex color code of the pixel at x, y.

###<h3 id="GetScreenSize">.GetScreenSize()</h3>
    获取屏幕大小

    Gets the screen width and height.

####返回值:

    Returns an object with .width and .height.

###<h3 id="CaptureScreen">.CaptureScreen</h3>
    //ScreenCapture
    获取部分或者全部屏幕
    Gets part or all of the screen.

####参数:

    x (optional)
    y (optional)
    height (optional)
    width (optional)
    If no arguments are provided, screen.capture will get the full screen.

####返回值:

    返回一个bitmap object.

##<h2 id="Bitmap">位图</h2>

    This is a work in progress.

###<h3 id="FindBitmap">.FindBitmap</h3>

    查找bitmap.

####参数:

    bitmap;
    rect(可选参数): x, y, w, h

####Return:

    查找到,返回bitmap的x和y坐标;没有返回nil


###<h3 id="OpenBitmap">.OpenBitmap</h3>

    打开bitmap图片.

####参数:

    bitmap图片路径,
    MMImageType(可选)

####返回值:

     返回一个bitmap对象

###<h3 id="SaveBitmap">.SaveBitmap</h3>

    保存一个bitmap图片.

####参数:

    bitmap对象,
    保存路径,
    imagetype(int) 

####返回值:

    保存图片


###<h3 id="TostringBitmap">.TostringBitmap</h3>

     将一个bitmap对象转换为字符串对象.

####参数:

    bitmap对象 

####Return:

    返回一个bitmap字符串 

###<h3 id="GetPortion">.GetPortion</h3>

     bitmap from a portion

####参数:

    bitmap,
    rect: x, y, w, h 

####返回值:

    Returns new bitmap object created from a portion of another    

###<h3 id="Convert">.Convert(openpath, savepath,MMImageType)</h3>

    转换bitmap图片格式

####参数:

    openpath,
    savepath,
    MMImageType(可选)

####例子:

```Go
    robotgo.Convert("test.png", "test.tif")
```  
##<h2 id="Event">事件</h2> 

###<h3 id="LEvent">.LEvent(string)</h3>

    监听全局事件

####参数:

    string

    (鼠标参数:mleft mright wheelDown wheelUp wheelLeft wheelRight)

####返回值:

    监听成功返回0

####例子:

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