# CHANGELOG

<!--### RobotGo-->
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

