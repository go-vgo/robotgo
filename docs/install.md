## CrossCompiling

##### Windows64 to windows32

```Go
SET CGO_ENABLED=1
SET GOARCH=386
go build main.go
```

#### Other to windows

Install Requirements (Ubuntu):

```bash
sudo apt install gcc-multilib
sudo apt install gcc-mingw-w64
# fix err: zlib.h: No such file or directory, Just used by bitmap.
sudo apt install libz-mingw-w64-dev
```

Build the binary:

```Go
GOOS=windows GOARCH=amd64 CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc CXX=x86_64-w64-mingw32-g++ go build -x ./
```

```
// CC=mingw-w64\x86_64-7.2.0-win32-seh-rt_v5-rev1\mingw64\bin\gcc.exe
// CXX=mingw-w64\x86_64-7.2.0-win32-seh-rt_v5-rev1\mingw64\bin\g++.exe
```

Some discussions and questions, please see [issues/228](https://github.com/go-vgo/robotgo/issues/228), [issues/143](https://github.com/go-vgo/robotgo/issues/143).
