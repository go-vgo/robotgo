#Robotgo
  
  >Golang 跨平台自动系统，控制键盘鼠标和读取屏幕
  
RobotGo 支持 Mac, Windows, and Linux(X11).

这是一项正在完善中的工作.

提Issues请到[Github](https://github.com/go-vgo/robotgo),便于统一管理和即时更新


##[中文文档](http://git.oschina.net/veni0/robotgo/blob/master/zh_doc.md)&nbsp;&nbsp;&nbsp;[API Document](https://github.com/go-vgo/robotgo/blob/master/doc.md) 
 英文文档请点击API Document.



##安装:
    go get github.com/go-vgo/robotgo

  It's that easy!

###环境要求:

####ALL:  
    Golang
    //Gcc

####For Mac OS X:

    Xcode Command Line Tools

####For Windows:

    MinGW or other GCC
    
####For everything else(Linux等其他系统):
    GCC
    
    X11 with the XTest extension (also known as the Xtst library)


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
