package main

import (
	"log"
	"runtime"
	"sync"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"

	"holdtotype/internal/appid"
)

func init() {
	runtime.LockOSThread()
}

const (
	wmTrayCallback = 0x0400 + 100
	wmTrayUpdate   = 0x0400 + 101

	nimAdd    = 0
	nimModify = 1
	nimDelete = 2

	nifMessage = 0x01
	nifIcon    = 0x02
	nifTip     = 0x04

	wmLButtonUp   = 0x0202
	wmRButtonUp   = 0x0205
	wmLButtonDbl  = 0x0203
	wmCommand     = 0x0111
	mfString      = 0x0000
	mfGrayed      = 0x0001
	mfSeparator   = 0x0800
	tpmReturnCmd  = 0x0100
	tpmRightAlign = 0x0008
)

const (
	cmdStatus = iota + 1
	cmdSettings
	cmdToggle
	cmdOpenConfig
	cmdOpenLog
	cmdAbout
	cmdQuit
	cmdLastCopy
)

var (
	procShellNotifyIconW         = shell32.NewProc("Shell_NotifyIconW")
	procCreatePopupMenu          = user32.NewProc("CreatePopupMenu")
	procDestroyMenu              = user32.NewProc("DestroyMenu")
	procAppendMenuW              = user32.NewProc("AppendMenuW")
	procTrackPopupMenu           = user32.NewProc("TrackPopupMenu")
	procGetCursorPos             = user32.NewProc("GetCursorPos")
	procCreateIconFromResourceEx = user32.NewProc("CreateIconFromResourceEx")
	procSetMenuInfo              = user32.NewProc("SetMenuInfo")
	procRegisterWindowMessageW   = user32.NewProc("RegisterWindowMessageW")
)

type notifyIconData struct {
	CbSize           uint32
	_                uint32
	HWnd             uintptr
	UID              uint32
	UFlags           uint32
	UCallbackMessage uint32
	_                uint32
	HIcon            uintptr
	SzTip            [128]uint16
	DwState          uint32
	DwStateMask      uint32
	SzInfo           [256]uint16
	UVersion         uint32
	SzInfoTitle      [64]uint16
	DwInfoFlags      uint32
	GuidItem         [16]byte
	HBalloonIcon     uintptr
}

type point struct{ X, Y int32 }

var (
	trayMu      sync.Mutex
	trayHwnd    uintptr
	trayReady   bool
	trayIconCur uintptr
	trayTipCur  string
	trayApp     *App

	trayIcons map[int]uintptr
)

const (
	trayIdle = iota
	trayRecording
	trayProcessing
	trayOff
	trayError
)

func hIconFromPNG(png []byte) uintptr {
	if len(png) == 0 {
		return 0
	}
	h, _, _ := procCreateIconFromResourceEx.Call(
		uintptr(unsafe.Pointer(&png[0])), uintptr(len(png)),
		1, 0x30000, 32, 32, 0)
	return h
}

func traySetIcon(state int) {
	trayMu.Lock()
	h := trayIcons[state]
	trayIconCur = h
	hwnd := trayHwnd
	ready := trayReady
	trayMu.Unlock()
	if ready {
		procPostMessageW.Call(hwnd, wmTrayUpdate, 0, 0)
	}
}

func traySetTooltip(tip string) {
	trayMu.Lock()
	trayTipCur = tip
	hwnd := trayHwnd
	ready := trayReady
	trayMu.Unlock()
	if ready {
		procPostMessageW.Call(hwnd, wmTrayUpdate, 0, 0)
	}
}

func trayNotifyData() notifyIconData {
	var nid notifyIconData
	nid.CbSize = uint32(unsafe.Sizeof(nid))
	nid.HWnd = trayHwnd
	nid.UID = 1
	nid.UFlags = nifMessage | nifIcon | nifTip
	nid.UCallbackMessage = wmTrayCallback
	nid.HIcon = trayIconCur
	tip, _ := windows.UTF16FromString(trayTipCur)
	copy(nid.SzTip[:], tip)
	return nid
}

var wmTaskbarCreated uintptr

func trayReadd() {
	trayMu.Lock()
	ready := trayReady
	nid := trayNotifyData()
	trayMu.Unlock()
	if !ready {
		return
	}
	procShellNotifyIconW.Call(nimAdd, uintptr(unsafe.Pointer(&nid)))
	log.Printf("tray: the taskbar restarted, the icon was added again")
}

func trayWndProc(hwnd, msg, wParam, lParam uintptr) uintptr {
	if wmTaskbarCreated != 0 && msg == wmTaskbarCreated {
		trayReadd()
		return 0
	}
	switch msg {
	case wmTrayCallback:
		switch lParam {
		case wmLButtonUp, wmLButtonDbl:
			go trayApp.openSettings("state")
		case wmRButtonUp:
			trayShowMenu(hwnd)
		}
		return 0
	case wmChromeRefresh:
		applyChrome(hwnd)
		return 0
	case wmTrayUpdate:
		trayMu.Lock()
		nid := trayNotifyData()
		trayMu.Unlock()
		procShellNotifyIconW.Call(nimModify, uintptr(unsafe.Pointer(&nid)))
		return 0
	case wmCommand:
		if wParam&0xFFFF == cmdQuit && trayApp != nil {
			go func() {
				trayApp.onExit()
				procPostMessageW.Call(hwnd, wmClose, 0, 0)
			}()
		}
		return 0
	case wmDestroy:
		procPostQuitMessage.Call(0)
		return 0
	}
	r, _, _ := procDefWindowProcW.Call(hwnd, msg, wParam, lParam)
	return r
}

func trayShowMenu(hwnd uintptr) {
	a := trayApp
	trayMu.Lock()
	status := trayTipCur
	trayMu.Unlock()
	a.mu.Lock()
	enabled := a.enabled
	a.mu.Unlock()

	toggleText := tr("menu.disable")
	if !enabled {
		toggleText = tr("menu.enable")
	}
	a.mu.Lock()
	last := a.lastResult
	a.mu.Unlock()
	items := []tmItem{
		{text: status, id: cmdStatus, grayed: true},
		{sep: true},
		{text: tr("menu.settings"), id: cmdSettings},
		{text: toggleText, id: cmdToggle},
		{text: tr("menu.lastcopy"), id: cmdLastCopy, grayed: last == ""},
		{sep: true},
		{text: tr("menu.open.config"), id: cmdOpenConfig},
		{text: tr("menu.open.log"), id: cmdOpenLog},
		{text: tr("menu.about"), id: cmdAbout},
		{sep: true},
		{text: tr("menu.quit"), id: cmdQuit},
	}
	cmd := showTrayMenu(items)

	switch cmd {
	case cmdLastCopy:
		ok, msg := a.copyLastResult()
		if a.snapshot().Overlay {
			if ok {
				overlaySet(ovFlashOK, tr("ov.copied"))
			} else {
				overlaySet(ovFlashErr, msg)
			}
		}
	case cmdSettings:
		go a.openSettings("state")
	case cmdToggle:
		a.toggleEnabled()
	case cmdOpenConfig:
		shellOpen("config.json")
	case cmdOpenLog:
		shellOpen(appid.LogFile)
	case cmdAbout:
		go a.openSettings("about")
	case cmdQuit:
		go func() {
			a.onExit()
			trayMu.Lock()
			h := trayHwnd
			trayMu.Unlock()
			procPostMessageW.Call(h, wmClose, 0, 0)
		}()
	}
}

func runTray(a *App) {
	trayApp = a
	trayIcons = loadTrayIcons()

	if name, err := windows.UTF16PtrFromString("TaskbarCreated"); err == nil {
		wmTaskbarCreated, _, _ = procRegisterWindowMessageW.Call(uintptr(unsafe.Pointer(name)))
	}

	className, _ := windows.UTF16PtrFromString(appid.Class("TrayWnd"))
	cb := syscall.NewCallback(trayWndProc)
	wc := wndClassExW{
		Size:      uint32(unsafe.Sizeof(wndClassExW{})),
		WndProc:   cb,
		ClassName: className,
	}
	procRegisterClassExW.Call(uintptr(unsafe.Pointer(&wc)))
	hwnd, _, _ := procCreateWindowExW.Call(0,
		uintptr(unsafe.Pointer(className)), 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0)
	if hwnd == 0 {
		msgBox(tr("err.title"), "tray window failed")
		return
	}

	trayMu.Lock()
	trayHwnd = hwnd
	if trayIconCur == 0 {
		trayIconCur = trayIcons[trayIdle]
	}
	nid := trayNotifyData()
	trayReady = true
	trayMu.Unlock()
	procShellNotifyIconW.Call(nimAdd, uintptr(unsafe.Pointer(&nid)))

	var m msgStruct
	for {
		r, _, _ := procGetMessageW.Call(uintptr(unsafe.Pointer(&m)), 0, 0, 0)
		if int32(r) <= 0 {
			break
		}
		procTranslateMessage.Call(uintptr(unsafe.Pointer(&m)))
		procDispatchMessageW.Call(uintptr(unsafe.Pointer(&m)))
	}
	trayMu.Lock()
	nid = trayNotifyData()
	trayMu.Unlock()
	procShellNotifyIconW.Call(nimDelete, uintptr(unsafe.Pointer(&nid)))
}

func trayReloadIcons(state int) {
	trayMu.Lock()
	trayIcons = loadTrayIcons()
	trayMu.Unlock()
	traySetIcon(state)
}
