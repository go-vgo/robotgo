# 方法:

##### [GetVersion](#GetVersion)

## [键盘](#Keyboard)

##### [Keys](https://github.com/go-vgo/robotgo/blob/master/docs/keys.md)
##### [SetKeyboardDelay](#SetKeyDelay) (相当于 SetKeyDelay, 废弃 API)
##### [SetKeyDelay](#SetKeyDelay)
##### [KeyTap](#KeyTap)
##### [KeyToggle](#KeyToggle)
##### [TypeString](#TypeString)
##### [TypeStringDelayed](#TypeStrDelay) (相当于 TypeStrDelay, 废弃 API)
##### [TypeStrDelay](#TypeStrDelay)
##### [TypeStr](#TypeStr)
##### [WriteAll](#WriteAll)
##### [ReadAll](#ReadAll)

## [鼠标](#Mouse)

##### [SetMouseDelay](#SetMouseDelay)
##### [MoveMouse](#MoveMouse)
##### [Move](#MoveMouse) (相当于 MoveMouse)
##### [MoveMouseSmooth](#MoveMouseSmooth)
##### [MoveSmooth](#MoveMouseSmooth) (相当与 MoveMouseSmooth)
##### [MouseClick](#MouseClick)
##### [Click](#MouseClick) (相当于 MouseClick)
##### [MoveClick](#MoveClick)
##### [MouseToggle](#MouseToggle)
##### [DragMouse](#DragMouse)
##### [Drag](#DragMouse) (相当于 DragMouse)
##### [GetMousePos](#GetMousePos)
##### [ScrollMouse](#ScrollMouse)

### <h3 id="GetVersion">.GetVersion()</h3>
    获取 robotgo 版本

## [屏幕](#Screen)

##### [GetPixelColor](#GetPixelColor)
##### [GetScreenSize](#GetScreenSize)
##### [CaptureScreen](#CaptureScreen)
##### [GetXDisplayName(Linux)](#GetXDisplayName)
##### [SetXDisplayName(Linux)](#SetXDisplayName)

## [位图](#Bitmap)
    This is a work in progress. (工作正在进行中)

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

## [事件](#Event)

##### [LEvent](#AddEvent) (相当于 AddEvent, 废弃 API)
##### [AddEvent](#AddEvent)
##### [StopEvent](#StopEvent)

## [窗口](#Window)
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


## <h2 id="Keyboard">键盘</h2>

### <h3 id="SetKeyDelay">.SetKeyDelay(ms)</h3>
    设置键盘延迟 (在键盘一个事件后), 单位 ms, 默认值 10ms

    Sets the delay in milliseconds to sleep after a keyboard event. This is 10ms by default.

#### 参数:
    延迟时间,单位 ms

    ms - Time to sleep in milliseconds.

### <h3 id="KeyTap">.KeyTap(key, modifier)</h3>
    模拟键盘按键

    Press a single key.

#### 参数:
    键盘值
    修饰值 (可选类型, 字符串或者数组) - 可选值: alt, command (win), control, and shift.

key - See [keys](https://github.com/go-vgo/robotgo/blob/master/docs/keys.md).  
modifier (optional, string or array) - Accepts alt, command (win), control, and shift.
#### 示例:

```Go
robotgo.KeyTap("h", "command")
robotgo.KeyTap("i", "alt", "command")
arr := []string{"alt", "command"}
robotgo.KeyTap("i", arr)
```

### <h3 id="KeyToggle">.KeyToggle(key, down, modifier)</h3>
    键盘切换, 按住或释放一个键位

    Hold down or release a key.

#### 参数:

key - See [keys](https://github.com/go-vgo/robotgo/blob/master/docs/keys.md).  
down - Accepts 'down' or 'up'.  
modifier (optional, string or array) - Accepts alt, command (mac), control, and shift.
### 返回值:

    返回 KeyToggle 状态

### <h3 id="TypeString">.TypeString(string)</h3>

#### 参数:

    string - The string to send.

### <h3 id="TypeStrDelay">.TypeStrDelay(string, cpm)</h3>

#### 参数:

    string - The string to send.
    cpm - Characters per minute.

### <h3 id="TypeStr">.TypeStr(string)</h3>

#### 参数:

    string - The string to send.

### <h3 id="WriteAll">.WriteAll(text string)</h3>

#### 参数:
    text string
#### 返回值:      

### <h3 id="ReadAll">.ReadAll()</h3>

#### 参数:

#### 返回值: 
    text,
    error 

## <h2 id="Mouse">鼠标</h2>
### <h3 id="SetMouseDelay">.SetMouseDelay(ms)</h3>
    设置鼠标延迟 (在一个鼠标事件后), 单位 ms, 默认值 10 ms

    Sets the delay in milliseconds to sleep after a mouse event. This is 10ms by default.

#### 参数:

    ms - Time to sleep in milliseconds.

### <h3 id="MoveMouse">.MoveMouse(x, y)</h3>
    移动鼠标

    Moves mouse to x, y instantly, with the mouse button up.

#### 参数:

    x, y

#### 示例:

```Go
// Move the mouse to 100, 100 on the screen. 
robotgo.MoveMouse(100, 100)
```

### <h3 id="MoveMouseSmooth">.MoveMouseSmooth(x, y)</h3>
    模拟鼠标向 X，Y 平滑移动(像人类一样)，用鼠标按钮向上

    Moves mouse to x, y human like, with the mouse button up.

#### 参数:

    x, y
    lowspeed, highspeed

#### 示例:

```Go
robotgo.MoveMouseSmooth(100, 200)
robotgo.MoveMouseSmooth(100, 200, 1.0, 100.0)
```       

### <h3 id="MouseClick">.MouseClick(button, double)</h3>
    鼠标点击

    Clicks the mouse.

#### 参数:

    button (optional) - Accepts left, right, or center. Defaults to left.
    double (optional) - Set to true to perform a double click. Defaults to false.

#### 示例:

```Go
robogo.MouseClick()
robogo.MouseClick("left", true)
```

### <h3 id="MoveClick">.MoveClick(x, y, button, double)</h3>

    移动并点击鼠标

#### 参数:
    x,
    y,

    button (optional) - Accepts "left", "right", or "center". Defaults to left.
    double (optional) - Set to true to perform a double click. Defaults to false.

#### 示例:

```Go
robogo.MoveClick(10, 20)
robogo.MoveClick(10, 20, "left", true)
```

### <h3 id="MouseToggle">.MouseToggle(down, button)</h3>
    鼠标切换

    Toggles mouse button.

#### 参数:

    down (optional) - Accepts down or up. Defaults to down.
    button (optional) - Accepts "left", "right", or "center". Defaults to left.

#### 示例:

```Go
robotgo.MouseToggle("down")
robotgo.MouseToggle("down", "right")
```

### <h3 id="DragMouse">.DragMouse(x, y)</h3>
    拖动鼠标

    Moves mouse to x, y instantly, with the mouse button held down.

#### 参数:

    x, y

#### 示例:

```Go
// Mouse down at 0, 0 and then drag to 100, 100 and release. 
robotgo.MoveMouse(0, 0)
robotgo.MouseToggle("down")
robotgo.DragMouse(100, 100)
robotgo.MouseToggle("up")
```

### <h3 id="GetMousePos">.GetMousePos()</h3>
    获取鼠标的位置

    Gets the mouse coordinates.

#### 返回值:

    Returns an object with keys x and y.

#### 示例:

```Go
x, y := robotgo.GetMousePos()
fmt.Println("pos:", x, y)
```

### <h3 id="ScrollMouse">.ScrollMouse(magnitude, direction)</h3>
    滚动鼠标

    Scrolls the mouse either up or down.

#### 参数:
    滚动位置的大小
    滚动方向: up (向上滚动)  down (向下滚动)

    magnitude - The amount to scroll.
    direction - Accepts down or up.

#### 示例:

```Go
robotgo.ScrollMouse(50, "up")

robotgo.ScrollMouse(50, "down")
```


## <h2 id="Screen">屏幕</h2>
### <h3 id="GetPixelColor">.GetPixelColor(x, y)
    获取坐标为 x, y 位置处的颜色

    Gets the pixel color at x, y. This function is perfect for getting a pixel or two,  
    but if you'll be reading large portions of the screen use screen.capture.

#### 参数:

    x, y

#### 返回值:

    Returns the hex color code of the pixel at x, y.

### <h3 id="GetScreenSize">.GetScreenSize()</h3>
    获取屏幕大小

    Gets the screen width and height.

#### 返回值:

    Returns an object with .width and .height.

### <h3 id="CaptureScreen">.CaptureScreen</h3>
    // CaptureScreen
    获取部分或者全部屏幕
    Gets part or all of the screen.

    GoCaptureScreen Returns a go struct
    Capture_Screen (废弃)

#### 参数:

    x (optional)
    y (optional)
    height (optional)
    width (optional)
    If no arguments are provided, capturescreen(screencapture) will get the full screen.

#### 返回值:

    返回一个 bitmap object.

## <h2 id="Bitmap">位图</h2>

    This is a work in progress.

### <h3 id="FindBitmap">.FindBitmap</h3>

    查找 bitmap.

#### 参数:

    bitmap;
    rect (可选参数): x, y, w, h

#### Return:

    查找到, 返回 bitmap 的 x 和y 坐标; 没有返回 nil


### <h3 id="OpenBitmap">.OpenBitmap</h3>

    打开 bitmap 图片.

#### 参数:

    bitmap 图片路径,
    MMImageType (可选)

#### 返回值:

     返回一个 bitmap 对象

### <h3 id="SaveBitmap">.SaveBitmap</h3>

    保存一个 bitmap 图片.

#### 参数:

    bitmap 对象,
    保存路径,
    imagetype (int) 

#### 返回值:

    保存图片, 返回保存状态


### <h3 id="TostringBitmap">.TostringBitmap</h3>

     将一个 bitmap 对象转换为字符串对象.

#### 参数:

    bitmap 对象 

#### Return:

    返回一个 bitmap 字符串 

### <h3 id="GetPortion">.GetPortion</h3>

     bitmap from a portion

#### 参数:

    bitmap,
    rect: x, y, w, h 

#### 返回值:

    Returns new bitmap object created from a portion of another    

### <h3 id="Convert">.Convert(openpath, savepath, MMImageType)</h3>

    转换 bitmap 图片格式

#### 参数:

    openpath,
    savepath,
    MMImageType (可选)

#### 示例:

```Go
robotgo.Convert("test.png", "test.tif")
```  

### <h3 id="FreeBitmap">.FreeBitmap(MMBitmapRef)</h3>

    FreeBitmap free and dealloc bitmap

#### 参数:

    MMBitmapRef

#### 示例:

```Go
robotgo.FreeBitmap(bitmap)
```    


### <h3 id="ReadBitmap">.ReadBitmap(MMBitmapRef)</h3>

    ReadBitmap returns false and sets error if |bitmap| is NULL

#### 参数:

    MMBitmapRef

#### 返回值:

    bool
    
#### 示例:

```Go
robotgo.ReadBitmap(bitmap)
```    


### <h3 id="CopyBitpb">.CopyBitpb(MMBitmapRef)</h3>

   CopyBitpb copy bitmap to pasteboard

#### 参数:

    MMBitmapRef

#### 返回值:

    bool

#### 示例:

```Go
robotgo.CopyBitpb(bitmap)
```    

### <h3 id="DeepCopyBit">.DeepCopyBit(MMBitmapRef)</h3>

   DeepCopyBit deep copy bitmap

#### 参数:

    MMBitmapRef

#### 返回值:

    MMBitmapRef

#### 示例:

```Go
robotgo.DeepCopyBit(bitmap)
```  

## <h2 id="Event">事件</h2> 

### <h3 id="AddEvent">.AddEvent(string)</h3>

    监听全局事件

#### 参数:

    string

    (鼠标参数: mleft, mright, wheelDown, wheelUp, wheelLeft, wheelRight)
 
#### 返回值:

    监听成功返回 0

#### 示例:

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
    停止事件监听

## <h2 id="Window">窗口</h2> 

### <h3 id="ShowAlert">.ShowAlert(title, msg, defaultButton, cancelButton string)</h3>

    Displays alert with the given attributes. If cancelButton is not given, only the defaultButton is displayed

#### 参数:

    title (string),
    msg (string),
    defaultButton (optional string),
    cancelButton (optional string)
           

#### 返回值:

    Returns 0 (True) if the default button was pressed, or 1 (False) if cancelled.
   
### <h3 id="CloseWindow">.CloseWindow()</h3>

    关闭窗口

#### 参数:
    无    

#### 返回值:
    无

### <h3 id="IsValid">.IsValid()</h3>

    Valid the Window

#### 参数:
    无      

#### 返回值:
    如果窗口 selected 返回 true 


### <h3 id="SetActive">.SetActive()</h3>

    设为当前窗口

#### 参数:
    hwnd  

#### 返回值:
    void


### <h3 id="GetActive">.GetActive()</h3>

    获取当前窗口

#### 参数:
    无       

#### 返回值:
    Returns hwnd

### <h3 id="SetHandle">.SetHandle()</h3>

    Set the Window Handle

#### 参数:
    int 

#### 返回值:
    bool    

### <h3 id="GetHandle">.GetHandle()</h3>

    获取窗口 Handle

#### 参数:
    无    

#### 返回值:
    Returns hwnd  

### <h3 id="GetTitle">.GetTitle()</h3>

    获取窗口标题

#### 参数:
           

#### 返回值:
    返回窗口标题  

### <h3 id="GetPID">.GetPID()</h3>

    获取进程 id

#### 参数:
    无   

#### 返回值:
    返回进程 id    

### <h3 id="Pids">.Pids()</h3>

    获取所有进程 id

#### 参数:
    无   

#### 返回值:
    返回进程 id   

### <h3 id="PidExists">.PidExists()</h3>

    判断进程 id 是否存在

#### 参数:
    pid  

#### 返回值:
    返回 bool 

### <h3 id="Process">.Process()</h3>

  Process get the all process struct

#### 参数:
    无  

#### 返回值:
    Returns []Nps, error

### <h3 id="FindName">.FindName()</h3>

    FindName find the process name by the process id

#### 参数:
    pid  

#### 返回值:
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