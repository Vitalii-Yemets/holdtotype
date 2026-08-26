package main

var (
	procStretchBlt        = gdi32.NewProc("StretchBlt")
	procSetStretchBltMode = gdi32.NewProc("SetStretchBltMode")
	procSaveDC            = gdi32.NewProc("SaveDC")
	procRestoreDC         = gdi32.NewProc("RestoreDC")
	procIntersectClipRect = gdi32.NewProc("IntersectClipRect")
	procSetBrushOrgEx     = gdi32.NewProc("SetBrushOrgEx")
)

const (
	stretchHalftone = 4
	dotScale        = 4
)

type dotCanvas struct {
	dc   uintptr
	bmp  uintptr
	prev uintptr
	side int32
}

func (d *dotCanvas) begin(hdc uintptr, side int32) uintptr {
	if side <= 0 {
		return 0
	}
	if d.dc != 0 && d.side == side {
		return d.dc
	}
	d.release()
	dc, _, _ := procCreateCompatibleDC.Call(hdc)
	if dc == 0 {
		return 0
	}
	bmp, _, _ := procCreateCompatibleBitmap.Call(hdc, uintptr(side), uintptr(side))
	if bmp == 0 {
		procDeleteDC.Call(dc)
		return 0
	}
	prev, _, _ := procSelectObject.Call(dc, bmp)
	d.dc, d.bmp, d.prev, d.side = dc, bmp, prev, side
	return d.dc
}

func (d *dotCanvas) release() {
	if d.dc != 0 {
		if d.prev != 0 {
			procSelectObject.Call(d.dc, d.prev)
		}
		procDeleteDC.Call(d.dc)
	}
	if d.bmp != 0 {
		procDeleteObject.Call(d.bmp)
	}
	d.dc, d.bmp, d.prev, d.side = 0, 0, 0, 0
}

var ovDotCanvas dotCanvas

func fillCircle(dc uintptr, cx, cy, r int32, color uintptr) {
	if r < 0 {
		return
	}
	br, _, _ := procCreateSolidBrush.Call(color)
	oldBr, _, _ := procSelectObject.Call(dc, br)
	pen, _, _ := procGetStockObject.Call(nullPen)
	oldPen, _, _ := procSelectObject.Call(dc, pen)
	procEllipse.Call(dc, uintptr(cx-r), uintptr(cy-r), uintptr(cx+r+1), uintptr(cy+r+1))
	procSelectObject.Call(dc, oldPen)
	procSelectObject.Call(dc, oldBr)
	procDeleteObject.Call(br)
}

func drawSmoothDot(hdc uintptr, cx, cy, core, halo int32, col, bg uintptr, strength float64) {
	if core < 1 {
		core = 1
	}
	if halo < core {
		halo = core
	}
	box := halo*2 + 3
	dc := ovDotCanvas.begin(hdc, box*dotScale)
	if dc == 0 {
		if halo > core {
			for r := halo; r > core; r-- {
				t := float64(halo-r) / float64(halo-core)
				fillCircle(hdc, cx, cy, r, blendCol(bg, col, strength*t*t))
			}
		}
		fillCircle(hdc, cx, cy, core, col)
		return
	}
	x0, y0 := cx-box/2, cy-box/2
	procStretchBlt.Call(dc, 0, 0, uintptr(box*dotScale), uintptr(box*dotScale),
		hdc, uintptr(x0), uintptr(y0), uintptr(box), uintptr(box), srcCopy)
	c := box * dotScale / 2
	if halo > core {
		for r := halo * dotScale; r > core*dotScale; r-- {
			t := float64(halo*dotScale-r) / float64((halo-core)*dotScale)
			fillCircle(dc, c, c, r, blendCol(bg, col, strength*t*t))
		}
	}
	fillCircle(dc, c, c, core*dotScale, col)
	procSetStretchBltMode.Call(hdc, stretchHalftone)
	procSetBrushOrgEx.Call(hdc, 0, 0, 0)
	procStretchBlt.Call(hdc, uintptr(x0), uintptr(y0), uintptr(box), uintptr(box),
		dc, 0, 0, uintptr(box*dotScale), uintptr(box*dotScale), srcCopy)
}
