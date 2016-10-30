#Robotgo
  
  >Golang 跨平台自动系统，控制键盘鼠标位图和读取屏幕以及全局事件监听
  
RobotGo 支持 Mac, Windows, and Linux(X11).

这是一项正在完善中的工作.

提Issues请到[Github](https://github.com/go-vgo/robotgo),便于统一管理和即时更新


##[中文文档](https://github.com/go-vgo/robotgo/blob/master/zh_doc.md)&nbsp;&nbsp;&nbsp;[API Document](https://github.com/go-vgo/robotgo/blob/master/doc.md) 
 英文文档请点击API Document.



##安装:
    go get github.com/go-vgo/robotgo

  It's that easy!

###环境要求:

####ALL:  
    Golang
    //Gcc
    zlib & libpng (bitmap)

####For Mac OS X:
    Xcode Command Line Tools

    brew install libpng
    brew install homebrew/dupes/zlib
####For Windows:
    MinGW or other GCC

#####[zlib&libpng Windows32 GCC教程](https://github.com/go-vgo/Mingw32)
#####[下载包含zlib和libpng的64位MinGW](https://github.com/go-vgo/Mingw)

####For everything else(Linux等其他系统):
    GCC
    
    X11 with the XTest extension (also known as the Xtst library)

#####Ubuntu:

      sudo apt-get install libx11-dev
      sudo apt-get install libgtkglextmm-x11-dev
      sudo apt-get install libghc6-x11-dev
      sudo apt-get install libgl1-mesa-swx11-dev
      sudo apt-get install xorg-dev

      sudo apt-get install libxtst-dev libpng++-dev

      事件:

      sudo apt-get install xcb libxcb-xkb-dev x11-xkb-utils libx11-xcb-dev libxkbcommon-x11-dev


##例子:

###鼠标

```Go
package main

import (
	//. "fmt"

	"github.com/go-vgo/robotgo"
)

func main() {
  robotgo.ScrollMouse(10, "up")
} 
``` 

###键盘

```Go
package main

import (
	//. "fmt"

	"github.com/go-vgo/robotgo"
)

func main() {
  robotgo.TypeString("Hello World")
  robotgo.KeyTap("enter")
  robotgo.TypeString("en")
} 
```

###屏幕

```Go
package main

import (
	. "fmt"

	"github.com/go-vgo/robotgo"
)

func main() {
  x, y := robotgo.GetMousePos()
  Println("pos:", x, y)
} 
```

###位图

```Go
package main

import (
	. "fmt"

	"github.com/go-vgo/robotgo"
)

func main() {
  bit_map := robotgo.CaptureScreen(10, 20, 30, 40)
  Println("...", bit_map)

  fx, fy := robotgo.FindBitmap(bit_map)
  Println("FindBitmap------", fx, fy)

  robotgo.SaveBitmap(bit_map, "test.png")
} 
```

###事件

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
