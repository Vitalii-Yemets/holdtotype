package main

import (
	"math"
	"sync"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"

	"holdtotype/internal/appid"
)

var (
	procUpdateLayeredWindow = user32.NewProc("UpdateLayeredWindow")
	procCreateDIBSection    = gdi32.NewProc("CreateDIBSection")
)

const (
	glowPadDIP  = 26
	wsExTransp  = 0x00000020
	swpNoSize   = 0x0001
	swpNoMove   = 0x0002
	swpNoActive = 0x0010
	swpShow     = 0x0040
	swpNoRedraw = 0x0008
	ulwAlpha    = 0x00000002
)

type bitmapInfoHeader struct {
	Size          uint32
	Width         int32
	Height        int32
	Planes        uint16
	BitCount      uint16
	Compression   uint32
	SizeImage     uint32
	XPelsPerMeter int32
	YPelsPerMeter int32
	ClrUsed       uint32
	ClrImportant  uint32
}

type blendFunction struct {
	Op     byte
	Flags  byte
	Alpha  byte
	Format byte
}

type sizeL struct{ Cx, Cy int32 }

type glowState struct {
	mu      sync.Mutex
	hwnd    uintptr
	owner   uintptr
	w, h    int32
	pad     int32
	radius  int32
	colour  uintptr
	visible bool
}

var glow glowState

func glowColour() (uintptr, bool) {
	look := themeLook()
	if !look.Round || look.Palette.Halo == "" {
		return 0, false
	}
	return colorref(look.Palette.Halo), true
}

func glowEnsureWindow() uintptr {
	if glow.hwnd != 0 {
		return glow.hwnd
	}
	className, _ := windows.UTF16PtrFromString(appid.Class("GlowWnd"))
	wc := wndClassExW{
		Size:      uint32(unsafe.Sizeof(wndClassExW{})),
		WndProc:   syscall.NewCallback(glowWndProc),
		ClassName: className,
	}
	procRegisterClassExW.Call(uintptr(unsafe.Pointer(&wc)))
	hwnd, _, _ := procCreateWindowExW.Call(
		wsExLayered|wsExToolWindow|wsExNoActivate|wsExTransp,
		uintptr(unsafe.Pointer(className)), 0,
		wsPopup,
		0, 0, 0, 0,
		0, 0, 0, 0,
	)
	glow.hwnd = hwnd
	return hwnd
}

func glowWndProc(hwnd, msg, wp, lp uintptr) uintptr {
	r, _, _ := procDefWindowProcW.Call(hwnd, msg, wp, lp)
	return r
}

func glowPaint(hwnd uintptr, w, h, pad, radius int32, colour uintptr) {
	if w <= 0 || h <= 0 {
		return
	}
	screen, _, _ := procGetDC.Call(0)
	if screen == 0 {
		return
	}
	defer procReleaseDC.Call(0, screen)
	dc, _, _ := procCreateCompatibleDC.Call(screen)
	if dc == 0 {
		return
	}
	defer procDeleteDC.Call(dc)

	head := bitmapInfoHeader{
		Size:     uint32(unsafe.Sizeof(bitmapInfoHeader{})),
		Width:    w,
		Height:   -h,
		Planes:   1,
		BitCount: 32,
	}
	var bits unsafe.Pointer
	bmp, _, _ := procCreateDIBSection.Call(dc, uintptr(unsafe.Pointer(&head)), 0,
		uintptr(unsafe.Pointer(&bits)), 0, 0)
	if bmp == 0 || bits == nil {
		return
	}
	defer procDeleteObject.Call(bmp)
	old, _, _ := procSelectObject.Call(dc, bmp)
	defer procSelectObject.Call(dc, old)

	px := unsafe.Slice((*uint32)(bits), int(w)*int(h))
	br := float64(colour & 0xFF)
	bg := float64((colour >> 8) & 0xFF)
	bb := float64((colour >> 16) & 0xFF)
	inLeft, inTop := float64(pad), float64(pad)
	inRight, inBottom := float64(w-pad-1), float64(h-pad-1)
	rad := float64(radius)
	span := float64(pad)
	for y := int32(0); y < h; y++ {
		fy := float64(y)
		row := int(y) * int(w)
		deep := fy > inTop+rad && fy < inBottom-rad
		for x := int32(0); x < w; x++ {
			if deep && x == pad {
				for ; x < w-pad; x++ {
					px[row+int(x)] = 0
				}
				if x >= w {
					break
				}
			}
			fx := float64(x)
			dx := math.Max(math.Max(inLeft+rad-fx, fx-(inRight-rad)), 0)
			dy := math.Max(math.Max(inTop+rad-fy, fy-(inBottom-rad)), 0)
			d := math.Hypot(dx, dy) - rad
			var a float64
			if d > 0 && d < span {
				t := 1 - d/span
				a = t * t * 0.62
			}
			if a <= 0 {
				px[row+int(x)] = 0
				continue
			}
			alpha := uint32(a * 255)
			px[row+int(x)] = alpha<<24 |
				uint32(bb*a)<<16 | uint32(bg*a)<<8 | uint32(br*a)
		}
	}

	size := sizeL{Cx: w, Cy: h}
	src := pointL{}
	blend := blendFunction{Op: 0, Alpha: 255, Format: 1}
	procUpdateLayeredWindow.Call(hwnd, screen, 0,
		uintptr(unsafe.Pointer(&size)), dc, uintptr(unsafe.Pointer(&src)),
		0, uintptr(unsafe.Pointer(&blend)), ulwAlpha)
}

func glowSync(owner uintptr) {
	if owner == 0 {
		return
	}
	colour, want := glowColour()
	glow.mu.Lock()
	defer glow.mu.Unlock()
	if !want {
		if glow.hwnd != 0 && glow.visible {
			procShowWindow.Call(glow.hwnd, swHide)
			glow.visible = false
		}
		return
	}
	if vis, _, _ := procIsWindowVisible.Call(owner); vis == 0 {
		if glow.hwnd != 0 && glow.visible {
			procShowWindow.Call(glow.hwnd, swHide)
			glow.visible = false
		}
		return
	}
	if z, _, _ := procIsZoomed.Call(owner); z != 0 {
		if glow.hwnd != 0 && glow.visible {
			procShowWindow.Call(glow.hwnd, swHide)
			glow.visible = false
		}
		return
	}
	var rc rect
	procGetWindowRect.Call(owner, uintptr(unsafe.Pointer(&rc)))
	dpi := dpiFor(owner)
	pad := scaleDPI(glowPadDIP, dpi)
	radius := scaleDPI(themeLook().Radius, dpi)
	w := rc.Right - rc.Left + pad*2
	h := rc.Bottom - rc.Top + pad*2
	if w <= 0 || h <= 0 {
		return
	}
	hwnd := glowEnsureWindow()
	if hwnd == 0 {
		return
	}
	glow.owner = owner
	if glow.w != w || glow.h != h || glow.pad != pad || glow.radius != radius || glow.colour != colour {
		procSetWindowPos.Call(hwnd, 0, uintptr(uint32(rc.Left-pad)), uintptr(uint32(rc.Top-pad)),
			uintptr(w), uintptr(h), swpNoActive|swpNoRedraw)
		glowPaint(hwnd, w, h, pad, radius, colour)
		glow.w, glow.h, glow.pad, glow.radius, glow.colour = w, h, pad, radius, colour
	} else {
		procSetWindowPos.Call(hwnd, 0, uintptr(uint32(rc.Left-pad)), uintptr(uint32(rc.Top-pad)),
			0, 0, swpNoSize|swpNoActive|swpNoRedraw)
	}
	if !glow.visible {
		procShowWindow.Call(hwnd, swShowNoActivate)
		glow.visible = true
	}
	procSetWindowPos.Call(hwnd, owner, 0, 0, 0, 0, swpNoSize|swpNoMove|swpNoActive)
}

func glowHide() {
	glow.mu.Lock()
	defer glow.mu.Unlock()
	if glow.hwnd != 0 && glow.visible {
		procShowWindow.Call(glow.hwnd, swHide)
		glow.visible = false
	}
}

func glowDestroy() {
	glow.mu.Lock()
	defer glow.mu.Unlock()
	if glow.hwnd != 0 {
		procDestroyWindow.Call(glow.hwnd)
	}
	glow.hwnd, glow.owner, glow.visible = 0, 0, false
	glow.w, glow.h, glow.pad, glow.radius, glow.colour = 0, 0, 0, 0, 0
}
