# tt
Simple and colorful test tools

[![CircleCI Status](https://circleci.com/gh/vcaesar/tt.svg?style=shield)](https://circleci.com/gh/vcaesar/tt)
![Appveyor](https://ci.appveyor.com/api/projects/status/github/vcaesar/tt?branch=master&svg=true)
[![codecov](https://codecov.io/gh/vcaesar/tt/branch/master/graph/badge.svg)](https://codecov.io/gh/vcaesar/tt)
[![Build Status](https://travis-ci.org/vcaesar/tt.svg)](https://travis-ci.org/vcaesar/tt)
[![Go Report Card](https://goreportcard.com/badge/github.com/vcaesar/tt)](https://goreportcard.com/report/github.com/vcaesar/tt)
[![GoDoc](https://godoc.org/github.com/vcaesar/tt?status.svg)](https://godoc.org/github.com/vcaesar/tt)
[![Release](https://github-release-version.herokuapp.com/github/vcaesar/tt/release.svg?style=flat)](https://github.com/vcaesar/tt/releases/latest)
[![Join the chat at https://gitter.im/go-ego/ego](https://badges.gitter.im/Join%20Chat.svg)](https://gitter.im/go-ego/ego?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)

## Installation/Update

```
go get -u github.com/vcaesar/tt
```

## Usage:

#### [Look at an example](/example/)

```go
package tt

import (
	"fmt"
	"testing"

	"github.com/vcaesar/tt"
	"github.com/vcaesar/tt/example"
)

func TestAdd(t *testing.T) {
	fmt.Println(add.Add(1, 1))

	tt.Expect(t, "1", add.Add(1, 1))
	tt.Expect(t, "2", add.Add(1, 1))

	tt.Equal(t, 1, add.Add(1, 1))
	tt.Equal(t, 2, add.Add(1, 1))

	at := tt.New(t)
	at.Expect("2", add.Add(1, 1))
	at.Equal(2, add.Add(1, 1))
}

func Benchmark1(b *testing.B) {
	at := tt.New(b)
	fn := func() {
		at.Equal(2, add.Add(1, 1))
	}

	tt.BM(b, fn)
	// at.BM(b, fn)
}

func Benchmark2(b *testing.B) {
	at := tt.New(b)
	for i := 0; i < b.N; i++ {
		at.Equal(2, Add(1, 1))
	}
}

```
## Thanks

[Testify](https://github.com/stretchr/testify), the code has some inspiration.