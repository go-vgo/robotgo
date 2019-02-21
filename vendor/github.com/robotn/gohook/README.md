# gohook

[![CircleCI Status](https://circleci.com/gh/robotn/gohook.svg?style=shield)](https://circleci.com/gh/robotn/gohook)
![Appveyor](https://ci.appveyor.com/api/projects/status/github/robotn/gohook?branch=master&svg=true)
[![Go Report Card](https://goreportcard.com/badge/github.com/robotn/gohook)](https://goreportcard.com/report/github.com/robotn/gohook)
[![GoDoc](https://godoc.org/github.com/robotn/gohook?status.svg)](https://godoc.org/github.com/robotn/gohook)
<!-- This is a work in progress. -->

```Go
package main

import (
	"fmt"

	"github.com/robotn/gohook"
)

func main() {
	EvChan := hook.Start()
	defer hook.End()
	
	for ev := range EvChan {
		fmt.Println(ev)
	}
}
```

Based on [libuiohook](https://github.com/kwhat/libuiohook).