package main

var (
	procCreateCompatibleDC     = gdi32.NewProc("CreateCompatibleDC")
	procCreateCompatibleBitmap = gdi32.NewProc("CreateCompatibleBitmap")
	procBitBlt                 = gdi32.NewProc("BitBlt")
	procDeleteDC               = gdi32.NewProc("DeleteDC")
)

const srcCopy = 0x00CC0020

type backBuf struct {
	dc   uintptr
	bmp  uintptr
	prev uintptr
	w    int32
	h    int32
}

func (b *backBuf) begin(hdc uintptr, w, h int32) uintptr {
	if w <= 0 || h <= 0 {
		return 0
	}
	if b.dc != 0 && b.w == w && b.h == h {
		return b.dc
	}
	b.release()
	dc, _, _ := procCreateCompatibleDC.Call(hdc)
	if dc == 0 {
		return 0
	}
	bmp, _, _ := procCreateCompatibleBitmap.Call(hdc, uintptr(w), uintptr(h))
	if bmp == 0 {
		procDeleteDC.Call(dc)
		return 0
	}
	prev, _, _ := procSelectObject.Call(dc, bmp)
	b.dc, b.bmp, b.prev, b.w, b.h = dc, bmp, prev, w, h
	return b.dc
}

func (b *backBuf) blit(hdc uintptr) {
	if b.dc == 0 {
		return
	}
	procBitBlt.Call(hdc, 0, 0, uintptr(b.w), uintptr(b.h), b.dc, 0, 0, srcCopy)
}

func (b *backBuf) release() {
	if b.dc != 0 {
		if b.prev != 0 {
			procSelectObject.Call(b.dc, b.prev)
		}
		procDeleteDC.Call(b.dc)
	}
	if b.bmp != 0 {
		procDeleteObject.Call(b.bmp)
	}
	b.dc, b.bmp, b.prev, b.w, b.h = 0, 0, 0, 0, 0
}
