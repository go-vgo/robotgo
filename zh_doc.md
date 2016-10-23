#方法:

##键盘
    Keys
    SetKeyboardDelay
    KeyTap
    KeyToggle
    TypeString
    TypeStringDelayed
##鼠标
    SetMouseDelay
    MoveMouse
    MoveMouseSmooth
    MouseClick
    MouseToggle
    DragMouse
    GetMousePos
    ScrollMouse
##屏幕
    GetPixelColor
    GetScreenSize
    CaptureScreen
    GetXDisplayName(Linux)
    SetXDisplayName(Linux)
##位图
    This is a work in progress.(工作正在进行中)


##键盘
###.SetKeyboardDelay(ms)
    设置键盘延迟(在键盘一个事件后),单位ms,默认值10ms

    Sets the delay in milliseconds to sleep after a keyboard event. This is 10ms by default.

####参数:
    延迟时间,单位ms

    ms - Time to sleep in milliseconds.

###.KeyTap(key, modifier)
    模拟键盘按键

    Press a single key.

####参数:
    键盘值
    修饰值(可选类型, 字符串或者数组(数组类型正在添加中)) - 可选值: alt, command (win), control, and shift.

    key - See keys.
    modifier (optional, string or array) - Accepts alt, command (win), control, and shift.

###.KeyToggle(key, down, modifier)
    键盘切换,按住或释放一个键位

    Hold down or release a key.

####参数:

    key - See keys.
    down - Accepts 'down' or 'up'.
    modifier (optional, string or array) - Accepts alt, command (mac), control, and shift.

###.TypeString(string)

####参数:

    string - The string to send.

###.TypeStringDelayed(string, cpm)

####参数:

    string - The string to send.
    cpm - Characters per minute.



##鼠标
###.SetMouseDelay(ms)
    设置鼠标延迟(在一个鼠标事件后),单位ms,默认值10ms

    Sets the delay in milliseconds to sleep after a mouse event. This is 10ms by default.

####参数:

    ms - Time to sleep in milliseconds.

###.MoveMouse(x, y)
    移动鼠标

    Moves mouse to x, y instantly, with the mouse button up.

####参数:

    x,y

####例子:

```Go
//Move the mouse to 100, 100 on the screen. 
robotgo.MoveMouse(100, 100);
```

###.MoveMouseSmooth(x, y)
    模拟鼠标向X，Y平滑移动(像人类一样)，用鼠标按钮向上

    Moves mouse to x, y human like, with the mouse button up.

####参数:

    x,y

###.MouseClick(button, double)
    鼠标点击

    Clicks the mouse.

####参数:

    button (optional) - Accepts left, right, or middle. Defaults to left.
    double (optional) - Set to true to perform a double click. Defaults to false.

####例子:

```Go
    robogo.MouseClick();
```

###.mouseToggle(down, button)
    鼠标切换

    Toggles mouse button.

####参数:

    down (optional) - Accepts down or up. Defaults to down.
    button (optional) - Accepts left, right, or middle. Defaults to left.

####例子:

```Go
robotgo.MouseToggle("down")
```

###.DragMouse(x, y)
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

###.GetMousePos()
    获取鼠标的位置

    Gets the mouse coordinates.

####返回值:

    Returns an object with keys x and y.

####例子:

```Go
x,y := robotgo.GetMousePos()
fmt.Println("pos:", x, y)
```

###.ScrollMouse(magnitude, direction)
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


##屏幕
###.GetPixelColor(x, y)
    获取坐标为x,y位置处的颜色

    Gets the pixel color at x, y. This function is perfect for getting a pixel or two, but if you'll be reading large portions of the screen use screen.capture.

####参数:

    x,y

####返回值:

    Returns the hex color code of the pixel at x, y.

###.GetScreenSize()
    获取屏幕大小

    Gets the screen width and height.

####返回值:

    Returns an object with .width and .height.

###.CaptureScreen
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

    Returns a bitmap object.

##位图

    This is a work in progress.


