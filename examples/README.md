# Robotgo examples

## Install:
```
go get -u github.com/go-vgo/robotgo  
```

## [Examples:](https://github.com/go-vgo/robotgo/blob/master/examples)

#### [Mouse](https://github.com/go-vgo/robotgo/blob/master/examples/mouse/main.go)

```Go
package main

import (
	"github.com/go-vgo/robotgo"
)

func main() {
  // robotgo.ScrollMouse(10, "up")
  robotgo.Scroll(0, 10)
  robotgo.MouseClick("left", true)
  robotgo.MoveSmooth(100, 200, 1.0, 100.0)
} 
``` 

#### [Keyboard](https://github.com/go-vgo/robotgo/blob/master/examples/key/main.go)

```Go
package main

import (
  "fmt"

  "github.com/go-vgo/robotgo"
)

func main() {
  robotgo.TypeStr("Hello World")
  // robotgo.TypeStr("だんしゃり")
  robotgo.TypeStr("だんしゃり")
  // ustr := uint32(robotgo.CharCodeAt("だんしゃり", 0))
  // robotgo.UnicodeType(ustr)

  robotgo.KeyTap("enter")
  robotgo.TypeStr("en")
  robotgo.KeyTap("i", "alt", "command")
  arr := []string{"alt", "command"}
  robotgo.KeyTap("i", arr)

  robotgo.WriteAll("Test")
  text, err := robotgo.ReadAll()
  if err == nil {
    fmt.Println(text)
  }
} 
```

#### [Screen](https://github.com/go-vgo/robotgo/blob/master/examples/screen/main.go)

```Go
package main

import (
	"fmt"

	"github.com/go-vgo/robotgo"
)

func main() {
  x, y := robotgo.Location()
  fmt.Println("pos:", x, y)
  color := robotgo.GetPixelColor(100, 200)
  fmt.Println("color----", color)
} 
```

#### [Bitmap](https://github.com/go-vgo/robotgo/blob/master/examples/bitmap/main.go)

```Go
package main

import (
	"fmt"

	"github.com/go-vgo/robotgo"
)

func main() {
  bitmap := robotgo.CaptureScreen(10, 20, 30, 40)
  // use `defer robotgo.FreeBitmap(bit)` to free the bitmap
  defer robotgo.FreeBitmap(bitmap)
  fmt.Println("...", bitmap)

  fx, fy := robotgo.FindBitmap(bitmap)
  fmt.Println("FindBitmap------", fx, fy)

  robotgo.SaveBitmap(bitmap, "test.png")
} 
```

#### [Event](https://github.com/go-vgo/robotgo/blob/master/examples/event/main.go)

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

#### [Window](https://github.com/go-vgo/robotgo/blob/master/examples/window/main.go)

```Go
package main

import (
	"fmt"

	"github.com/go-vgo/robotgo"
)

func main() {
  fpid, err := robotgo.FindIds("Google")
  if err == nil {
    fmt.Println("pids...", fpid)

    if len(fpid) > 0 {
      robotgo.ActivePID(fpid[0])

      robotgo.Kill(fpid[0])
    }
  }

  robotgo.ActiveName("chrome")

  isExist, err := robotgo.PidExists(100)
  if err == nil && isExist {
    fmt.Println("pid exists is", isExist)

    robotgo.Kill(100)
  }

  abool := robotgo.ShowAlert("test", "robotgo")
  if abool == 0 {
 	  fmt.Println("ok@@@", "ok")
  }

  title := robotgo.GetTitle()
  fmt.Println("title@@@", title)
} 
```