//go:build linux
// +build linux

// Copyright (c) 2016-2026 AtomAI, All rights reserved.
//
// See the COPYRIGHT file at the top-level directory of this distribution and at
// https://github.com/go-vgo/robotgo/blob/master/LICENSE
//
// Licensed under the Apache License, Version 2.0 <LICENSE-APACHE or
// http://www.apache.org/licenses/LICENSE-2.0>
//
// This file may not be copied, modified, or distributed
// except according to those terms.

package wayland

import (
	"fmt"
	"image"
	"image/color"
	"os"
	"sync"
	"time"

	"golang.org/x/sys/unix"

	"github.com/vcaesar/go-wayland/client"

	"github.com/go-vgo/robotgo/wayland/internal/protocols/wlr_screencopy"
)

// wl_shm 32-bit formats the capture path understands (drm_fourcc codes; the
// two legacy values are the wl_shm enum). Byte order in memory, little-endian.
const (
	shmARGB8888 = 0          // B G R A
	shmXRGB8888 = 1          // B G R X
	shmABGR8888 = 0x34324241 // R G B A
	shmXBGR8888 = 0x34324258 // R G B X
)

// captureTimeout bounds how long CaptureImg waits for the compositor.
const captureTimeout = 5 * time.Second

// GetScreenSize returns the primary output's width and height.
func GetScreenSize() (int, int) {
	c, err := ensureConn()
	if err != nil {
		return 0, 0
	}
	o, ok := c.output(0)
	if !ok {
		return 0, 0
	}
	return int(o.width), int(o.height)
}

// GetScaleSize returns the scaled screen size. On Wayland the screencopy
// buffer is already in physical pixels, so the scale factor is 1.0 and this
// returns the same value as GetScreenSize. Provided for robotgo API parity.
func GetScaleSize(displayId ...int) (int, int) {
	return GetScreenSize()
}

// GetScreenRect returns the screen rectangle.
func GetScreenRect(displayId ...int) Rect {
	c, err := ensureConn()
	if err != nil {
		return Rect{}
	}
	idx := 0
	if len(displayId) > 0 {
		idx = displayId[0]
	}
	o, ok := c.output(idx)
	if !ok {
		return Rect{}
	}
	return Rect{
		Point: Point{X: int(o.x), Y: int(o.y)},
		Size:  Size{W: int(o.width), H: int(o.height)},
	}
}

// DisplaysNum returns the number of displays.
func DisplaysNum() int {
	c, err := ensureConn()
	if err != nil {
		return 0
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	return len(c.outputs)
}

// GetPixelColor returns the pixel color at (x, y) as a hex string.
func GetPixelColor(x, y int, displayId ...int) string {
	img, err := CaptureImg(displayId...)
	if err != nil || img == nil {
		return "000000"
	}

	bounds := img.Bounds()
	if x < bounds.Min.X || x >= bounds.Max.X || y < bounds.Min.Y || y >= bounds.Max.Y {
		return "000000"
	}

	c := img.At(x, y)
	r, g, b, _ := c.RGBA()
	return fmt.Sprintf("%02x%02x%02x", r>>8, g>>8, b>>8)
}

// CaptureImg captures the screen and returns an image.Image.
// Optional args: x, y, w, h, displayId
func CaptureImg(args ...int) (image.Image, error) {
	c, err := ensureConn()
	if err != nil {
		return nil, err
	}
	if c.screencopyMgr == nil || c.shm == nil {
		return nil, ErrNotSupported
	}

	// Parse args
	displayIdx := 0
	var region *image.Rectangle
	if len(args) >= 4 {
		r := image.Rect(args[0], args[1], args[0]+args[2], args[1]+args[3])
		region = &r
	}
	if len(args) >= 5 {
		displayIdx = args[4]
	}
	output, ok := c.output(displayIdx)
	if !ok {
		return nil, ErrNotSupported
	}

	// All screencopy events fire from the single background dispatch goroutine
	// (see conn.dispatchLoop) while it holds c.wl, so the handlers below may
	// call proxies directly. We collect state here and use a channel to hand
	// off completion to the caller, so only one goroutine ever touches the
	// captured pixels — no data race.
	var (
		bufFormat, bufWidth, bufHeight, bufStride uint32
		yInvert                                   bool
		capturedData                              []byte
		captureErr                                error

		// Wayland resources kept alive until the copy completes.
		mmapData []byte
		shmPool  *client.ShmPool
		buffer   *client.Buffer
	)

	done := make(chan struct{})
	var finishOnce sync.Once
	finish := func(err error) {
		finishOnce.Do(func() {
			captureErr = err
			close(done)
		})
	}

	// Request a frame and install the handlers in one locked section so no
	// event can be dispatched before the handlers exist.
	var frame *wlr_screencopy.ZwlrScreencopyFrameV1
	err = c.do(func() error {
		var err error
		if region != nil {
			frame, err = c.screencopyMgr.CaptureOutputRegion(
				1, output.output,
				int32(region.Min.X), int32(region.Min.Y),
				int32(region.Dx()), int32(region.Dy()),
			)
		} else {
			frame, err = c.screencopyMgr.CaptureOutput(1, output.output)
		}
		if err != nil {
			return fmt.Errorf("robotgo: capture output: %w", err)
		}

		frame.SetBufferHandler(func(e wlr_screencopy.ZwlrScreencopyFrameV1BufferEvent) {
			bufFormat = e.Format
			bufWidth = e.Width
			bufHeight = e.Height
			bufStride = e.Stride
		})

		frame.SetFlagsHandler(func(e wlr_screencopy.ZwlrScreencopyFrameV1FlagsEvent) {
			yInvert = e.Flags&uint32(wlr_screencopy.ZwlrScreencopyFrameV1FlagsYInvert) != 0
		})

		frame.SetFailedHandler(func(_ wlr_screencopy.ZwlrScreencopyFrameV1FailedEvent) {
			finish(fmt.Errorf("robotgo: compositor reported screen capture failure"))
		})

		// buffer_done signals that all buffer events have been received and we
		// may create the SHM buffer the compositor will copy into.
		frame.SetBufferDoneHandler(func(_ wlr_screencopy.ZwlrScreencopyFrameV1BufferDoneEvent) {
			switch bufFormat {
			case shmARGB8888, shmXRGB8888, shmABGR8888, shmXBGR8888:
			default:
				finish(fmt.Errorf("robotgo: unsupported capture format %#x", bufFormat))
				return
			}
			size := int(bufStride * bufHeight)
			if size == 0 {
				finish(fmt.Errorf("robotgo: invalid capture buffer size"))
				return
			}

			dir := os.Getenv("XDG_RUNTIME_DIR")
			if dir == "" {
				dir = os.TempDir()
			}

			f, ferr := os.CreateTemp(dir, "robotgo-screencopy-*")
			if ferr != nil {
				finish(fmt.Errorf("robotgo: create capture temp file: %w", ferr))
				return
			}
			// CreatePool dups the fd via SCM_RIGHTS and mmap holds its own
			// reference, so the file can be closed and unlinked immediately.
			defer f.Close()
			defer os.Remove(f.Name())

			if ferr := f.Truncate(int64(size)); ferr != nil {
				finish(fmt.Errorf("robotgo: truncate capture file: %w", ferr))
				return
			}

			data, merr := unix.Mmap(int(f.Fd()), 0, size, unix.PROT_READ|unix.PROT_WRITE, unix.MAP_SHARED)
			if merr != nil {
				finish(fmt.Errorf("robotgo: mmap capture file: %w", merr))
				return
			}

			pool, serr := c.shm.CreatePool(int(f.Fd()), int32(size))
			if serr != nil {
				unix.Munmap(data)
				finish(fmt.Errorf("robotgo: create shm pool: %w", serr))
				return
			}

			buf, berr := pool.CreateBuffer(0, int32(bufWidth), int32(bufHeight), int32(bufStride), bufFormat)
			if berr != nil {
				unix.Munmap(data)
				pool.Destroy()
				finish(fmt.Errorf("robotgo: create buffer: %w", berr))
				return
			}

			mmapData = data
			shmPool = pool
			buffer = buf

			// Ask the compositor to copy the framebuffer into our buffer.
			// The ready (or failed) event arrives afterwards.
			if cerr := frame.Copy(buf); cerr != nil {
				finish(fmt.Errorf("robotgo: frame copy: %w", cerr))
			}
		})

		frame.SetReadyHandler(func(_ wlr_screencopy.ZwlrScreencopyFrameV1ReadyEvent) {
			if mmapData != nil {
				capturedData = make([]byte, len(mmapData))
				copy(capturedData, mmapData)
			}
			finish(nil)
		})
		return nil
	})
	if err != nil {
		return nil, err
	}

	// Block until the compositor signals ready or failed. The channel receive
	// establishes a happens-before edge with the handler writes above. Give
	// up if the connection dies or the compositor never answers.
	select {
	case <-done:
	case <-c.dispatchDone:
		finish(ErrNoConnection)
	case <-time.After(captureTimeout):
		finish(fmt.Errorf("robotgo: screen capture timed out"))
	}

	// Release Wayland buffer resources now that the copy is complete (or
	// abandoned). This runs under c.wl, so a late ready handler (which runs
	// under the same lock) sees mmapData cleared and never touches the
	// unmapped memory; finish is already a no-op for it.
	c.wl.Lock()
	data := mmapData
	mmapData = nil
	if buffer != nil {
		buffer.Destroy()
	}
	if shmPool != nil {
		shmPool.Destroy()
	}
	frame.Destroy()
	c.wl.Unlock()
	if data != nil {
		unix.Munmap(data)
	}

	if captureErr != nil {
		return nil, captureErr
	}
	if capturedData == nil {
		return nil, fmt.Errorf("robotgo: screen capture produced no data")
	}

	return decodeShm(capturedData, bufFormat, int(bufWidth), int(bufHeight), int(bufStride), yInvert), nil
}

// decodeShm converts a 32-bit wl_shm pixel buffer into an RGBA image,
// honouring the channel order of the format and the y_invert flag.
func decodeShm(data []byte, format uint32, w, h, stride int, yInvert bool) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, w, h))

	// ARGB/XRGB are stored B,G,R,A; ABGR/XBGR are R,G,B,A. The X variants
	// carry no alpha, so force opaque.
	swapRB := format == shmARGB8888 || format == shmXRGB8888
	hasAlpha := format == shmARGB8888 || format == shmABGR8888
	for y := 0; y < h; y++ {
		srcY := y
		if yInvert {
			srcY = h - 1 - y
		}
		for x := 0; x < w; x++ {
			off := srcY*stride + x*4
			if off+3 >= len(data) {
				break
			}
			r, g, b := data[off+0], data[off+1], data[off+2]
			if swapRB {
				r, b = b, r
			}
			a := byte(0xff)
			if hasAlpha {
				a = data[off+3]
			}
			img.SetRGBA(x, y, color.RGBA{R: r, G: g, B: b, A: a})
		}
	}
	return img
}

// Capture captures the screen and returns an *image.RGBA.
func Capture(args ...int) (*image.RGBA, error) {
	img, err := CaptureImg(args...)
	if err != nil {
		return nil, err
	}
	if rgba, ok := img.(*image.RGBA); ok {
		return rgba, nil
	}
	// Convert
	bounds := img.Bounds()
	rgba := image.NewRGBA(bounds)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			rgba.Set(x, y, img.At(x, y))
		}
	}
	return rgba, nil
}
