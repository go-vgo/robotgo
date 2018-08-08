# CHANGELOG

<!--### RobotGo-->
# ## RobotGo v0.50.0, The Appalachian Mountains

## Add  

- [NEW] Add simple ocr support
- [NEW] Add max and min window api and win32.h file
- [NEW] Automatic free internal bitmap and add bitmapStr example
- [NEW] Update findBitmap and findColor default tolerance 0.5 to 0.01, [improve find accuracy and time]
- [NEW] Add more Window example
- [NEW] Add cross compile docs
- [NEW] Add free bitmap and tolerance godoc
- [NEW] Add GetForegroundWindow and FindWindow func support
- [NEW] Add bitmap to CBitmap func, Add ToCBitmap example to examples
- [NEW] Add get Scale and GetScaleSize func, get primary display DPI scale factor fix #129, #135
   Add Custom DPI Scaling support, 
   Add scale default value,
   Add scale example

## Update

- [NEW] Update README.md [add freeBitmap example]
- [NEW] Optimize findColor and uniform API with findBitmap
- [NEW] Update godoc, CI and README.md
- [NEW] Update CHANGELOG.md
- [NEW] Update examples [add freeBitmap and update findColor]
- [NEW] Optimize bitmap code, optimize args and not try [many methods]
- [NEW] Update getPid type to int32
- [NEW] Update var and other code style, fix non-standard code
   Update code and update some name
- [NEW] Update pkg to newest
- [NEW] Remove duplicate code and old useless code
- [NEW] Update mouse click and fix moveClick and movesClick args
- [NEW] Update code style use if not try
- [NEW] Update clipboard example
- [NEW] Update typestr use return not else
- [NEW] Update mouse toggle, keytap and savebitmap func args
- [NEW] Update examples remove duplicate code
- [NEW] Update bitmap and other examples
- [NEW] Simplify linux dependency installation commands
- [NEW] Update issue_template.md
-[NEW] Update pull_request_template.md
- [NEW] Move govendor to dep
- [NEW] Update robotgo ci to 1.10.3

## Fix

- [FIX] Update active pid to fix #140, fixed linux activePid
- [FIX] Fixed findBitmap and findPic memory leak
- [FIX] Add getPxColor destroyMMBitmap fix memory leak
- [FIX] Fix float args not float32
- [FIX] Fix windows clipboard memory leak
- [FIX] Update macos .a downgrade to 10.10 just warning not exit [fix #102, #128, #134]
- [FIX] use 10.10 to compile .a verifyed multi os
- [FIX] Fix #145 not assert
- [FIX] Fix some warning use supplemental code

See Commits for more details, after Apr 30.


## RobotGo v0.49.0, Olympus Mons

### Add  

- [NEW] Add get image size func
- [NEW] Add linux type string utf-8 support
- [NEW] Add scroll mouse support x, y
- [NEW] Add AddEvent() "esc" support fix #105
- [NEW] Add AddEvent "space" fix #110
- [NEW] Add clipboard choose primary mode on unix
- [NEW] Add move smooth return
- [NEW] Add more bitmap func and examples
- [NEW] Add MicroSleep func
- [NEW] Add find image by path


### Update

- [NEW] Update KeyToggle code
- [NEW] Update activePid allow Windows via hwnd
- [NEW] Update godoc and README.md
- [NEW] Update CHANGELOG.md
- [NEW] Update Kill() parameter and examples
- [NEW] Update examples and remove useless function
- [NEW] Update appveyor, circle and dockerfile
- [NEW] Update code style
- [NEW] Update and optimize func
- [NEW] Update travis support go 1.10
- [NEW] Update CI (use custom go image) and add func internalFindBitmap
- [NEW] Update godoc and deprecated GetBHandle
- [NEW] Optimize code func args and name


### Fix

- [FIX] Fix mac input method keytap not work
- [FIX] Fix clipboard golint
- [FIX] Update move smooth fix #96 (set mouse smooth speed)
- [FIX] Fix Getportion param to go type
- [FIX] Fix XFlush wait for events flushing

See Commits for more details, after Jan 25.

## RobotGo v0.48.0, Ben Nevis

### Add  

- [NEW] Add active window by name func ActiveName
- [NEW] Add type string utf-8 support

Add func CharCodeAt, UnicodeType, PasteStr and update TypeStr, TypeString

- [NEW] Add count of bitmap func CountBitmap
- [NEW] Add func SaveCapture and examples
- [NEW] Add time sleep func Sleep
- [NEW] Add more key listen
- [NEW] Add func PointInBounds and examples
- [NEW] Add func GetPxColor return C.MMRGBHex
- [NEW] Add FindColorCS param tolerance
- [NEW] Add func ToBitmap and examples
- [NEW] Add CBitmap type and examples
- [NEW] Add more examples
- [NEW] Add func ToMMBitmapRef
- [NEW] Add func BitmapClick and MovesClick
- [NEW] Add func ToMMRGBHex convert color hex
- [NEW] Add  func count bitmap color and CountColorCS
- [NEW] Add more color processing and conversion

Add func ToMMRGBHex, U32ToHex, U8ToHex, PadHex, HexToRgb, RgbToHex and examples

- [NEW] Add func tochar bitmap and gostring and fmt code


### Update
- [NEW] Remove robot and examples
- [NEW] Update vendor and appveyor.yml
- [NEW] Update keyboard code
- [NEW] Update godoc
- [NEW] Update CHANGELOG.md
- [NEW] Change TostringBitmap return string
- [NEW] Update C language code and other naming 
- [NEW] Update code and code style
- [NEW] Update move mouse smooth


### Fix

- [FIX] Fix mac set active and active by pid
- [FIX] Fix windows active by pid #101 
- [FIX] Fix FindColor param tolerance
- [FIX] Fix find bitmap float args
- [FIX] Fix some range error
- [FIX] Update doc fix #97
- [FIX] Update README.md fix link error

See Commits for more details, after Dec 13.

## RobotGo v0.47.0, Mount Cook

### Add  

- [NEW] Add windows 32bit and 64bit dependency
- [NEW] Add macOs dependency
- [NEW] Add pkg to vendor

Solve the problem of dependence, remove zlib/libpng dependencies

- [NEW] Add FindColorCS(x, y, w, h int, color CHex), CHex type and examples #84
- [NEW] Add kill the process
- [NEW] Add public event and update code
- [NEW] Add  Windows 32bit and 64bit Appveyor CI
 

### Update
- [NEW] Update png io
- [NEW] Update cgo link
- [NEW] Update .gitignore
- [NEW] Update README.md and godoc
- [NEW] Update CHANGELOG.md
- [NEW] Update circle to 2.0, add robotgo Dockerfile custom image
- [NEW] Update and fmt C code
- [NEW] Update GetTitle default value "null" to ""


### Fix

- [FIX] Fix FindColor inconvenient parameters
- [FIX] Fix installation requirements #72
- [FIX] Fix GetTitle `return address of local variable` in the higher gcc version. #81

See Commits for more details, after Nov 10.


## RobotGo v0.46.6, Pyrenees Mountains: Aneto Peak

## RobotGo v0.46.0, Pyrenees Mountains

### Add

- [NEW] Add ActivePID
- [NEW] Add FindBit
- [NEW] Add robot branch, where there is no zlib and libpng dependency

### Update

- [NEW] Update README.md
- [NEW] Update FindIds
- [NEW] Update examples
- [NEW] Update vendor
- [NEW] Update godoc and docs
- [NEW] Update and fix bitmap

### Fix

- [FIX] Fix MoveMouseSmooth args
- [FIX] Fix name err
- [FIX] Fix FindBitmap

## RobotGo v0.45.0, Mount Qomolangma

### Add
- [NEW] Add Process
- [NEW] Add TypeStr
- [NEW] Add DeepCopyBit
- [NEW] Add CopyBitpb
- [NEW] Add ReadBitmap
- [NEW] Add vendor.json
- [NEW] Add ReadAll: clipboard
- [NEW] Add WriteAll: clipboard
- [NEW] Add Pids : get the all process id
- [NEW] Add FindName: find the process name by the process id
- [NEW] Add FindNames: find the all process name
- [NEW] Add PidExists: determine whether the process exists
- [NEW] Add FindIds: find the process id by the process name
- [NEW] Add FreeBitmap and Update docs


### Update
- [NEW] Update docs
- [NEW] Update test
- [NEW] Update godoc
- [NEW] Update CHANGELOG.md
- [NEW] Update .gitignore
- [NEW] Update examples and docs
- [NEW] Update examples link
- [NEW] Update README.md and clipboard


### Fix

- [FIX] Fix release key
- [FIX] Fix godoc error


## RobotGo v0.44.0, Mount Kailash

### Add

- Add CHANGELOG.md
- Format some code	
- Add fedora dependencies

### Update

- Update test
- Update keys.md
- Update and Split example
- Update godoc and docs
- Update and Cleanup README.md
- Update CONTRIBUTING.md and issue_template.md

### Fix

- Fix typesetting and MD error
- Fix fedora dependencies #55 
- Fix doc.md and README.md

