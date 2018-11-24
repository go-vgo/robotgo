module github.com/go-vgo/robotgo

require (
	github.com/BurntSushi/xgb v0.0.0-20160522181843-27f122750802
	github.com/BurntSushi/xgbutil v0.0.0-20160919175755-f7c97cef3b4e
	github.com/StackExchange/wmi v0.0.0-20180725035823-b12b22c5341f // indirect
	github.com/go-ole/go-ole v1.2.1 // indirect
	github.com/lxn/win v0.0.0-20181015143721-a7f87360b10e
	github.com/otiai10/gosseract v2.2.0+incompatible
	github.com/robotn/gohook v0.0.0-20181113164304-f27e2e52653b
	github.com/shirou/gopsutil v2.18.10+incompatible
	github.com/shirou/w32 v0.0.0-20160930032740-bb4de0191aa4 // indirect
	github.com/vcaesar/imgo v0.0.0-20181001170449-7a535c786a55
	golang.org/x/image v0.0.0-20181116024801-cd38e8056d9b // indirect
	golang.org/x/sys v0.0.0-20181122145206-62eef0e2fa9b // indirect
)

replace (
	golang.org/x/image/bmp v0.0.0-20181116024801-cd38e8056d9b => github.com/golang/image/bmp v0.0.0-20181116024801-cd38e8056d9b
	golang.org/x/sys v0.0.0-20181122145206-62eef0e2fa9b => github.com/golang/sys v0.0.0-20181122145206-62eef0e2fa9b
)
