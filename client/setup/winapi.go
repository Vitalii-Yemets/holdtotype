package main

import (
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
	user32   = windows.NewLazySystemDLL("user32.dll")
	kernel32 = windows.NewLazySystemDLL("kernel32.dll")
	dwmapi   = windows.NewLazySystemDLL("dwmapi.dll")
	shell32  = windows.NewLazySystemDLL("shell32.dll")

	procMessageBoxW           = user32.NewProc("MessageBoxW")
	procGetWindowLongPtrW     = user32.NewProc("GetWindowLongPtrW")
	procSetWindowLongPtrW     = user32.NewProc("SetWindowLongPtrW")
	procSetWindowPos          = user32.NewProc("SetWindowPos")
	procReleaseCapture        = user32.NewProc("ReleaseCapture")
	procSendMessageW          = user32.NewProc("SendMessageW")
	procShowWindow            = user32.NewProc("ShowWindow")
	procSetForegroundWnd      = user32.NewProc("SetForegroundWindow")
	procDwmSetWindowAttribute = dwmapi.NewProc("DwmSetWindowAttribute")
	procShellExecuteW         = shell32.NewProc("ShellExecuteW")
	procGetUserDefaultUILang  = kernel32.NewProc("GetUserDefaultUILanguage")

	ole32                    = windows.NewLazySystemDLL("ole32.dll")
	procSHBrowseForFolderW   = shell32.NewProc("SHBrowseForFolderW")
	procSHGetPathFromIDListW = shell32.NewProc("SHGetPathFromIDListW")
	procCoTaskMemFree        = ole32.NewProc("CoTaskMemFree")
)

type browseInfoW struct {
	Owner       uintptr
	Root        uintptr
	DisplayName *uint16
	Title       *uint16
	Flags       uint32
	Callback    uintptr
	LParam      uintptr
	Image       int32
}

func browseFolder(owner uintptr, title string) string {
	display := make([]uint16, 260)
	t, _ := windows.UTF16PtrFromString(title)
	bi := browseInfoW{
		Owner:       owner,
		DisplayName: &display[0],
		Title:       t,
		Flags:       0x00000041,
	}
	pidl, _, _ := procSHBrowseForFolderW.Call(uintptr(unsafe.Pointer(&bi)))
	if pidl == 0 {
		return ""
	}
	defer procCoTaskMemFree.Call(pidl)
	buf := make([]uint16, 300)
	if r, _, _ := procSHGetPathFromIDListW.Call(pidl, uintptr(unsafe.Pointer(&buf[0]))); r == 0 {
		return ""
	}
	return windows.UTF16ToString(buf)
}

func msgBox(title, text string) {
	t, _ := windows.UTF16PtrFromString(title)
	m, _ := windows.UTF16PtrFromString(text)
	procMessageBoxW.Call(0, uintptr(unsafe.Pointer(m)), uintptr(unsafe.Pointer(t)), 0x10)
}

func shellOpenURL(url string) {
	verb, _ := windows.UTF16PtrFromString("open")
	p, _ := windows.UTF16PtrFromString(url)
	procShellExecuteW.Call(0, uintptr(unsafe.Pointer(verb)), uintptr(unsafe.Pointer(p)), 0, 0, 1)
}

func applyDarkCaption(hwnd uintptr) {
	one := int32(1)
	procDwmSetWindowAttribute.Call(hwnd, 20, uintptr(unsafe.Pointer(&one)), 4)
	capColor := int32(0x0C0F0B)
	procDwmSetWindowAttribute.Call(hwnd, 35, uintptr(unsafe.Pointer(&capColor)), 4)
	txtColor := int32(0x3CFF6E)
	procDwmSetWindowAttribute.Call(hwnd, 36, uintptr(unsafe.Pointer(&txtColor)), 4)
	border := int32(0x2B4A1D)
	procDwmSetWindowAttribute.Call(hwnd, 34, uintptr(unsafe.Pointer(&border)), 4)
	corner := int32(2)
	procDwmSetWindowAttribute.Call(hwnd, 33, uintptr(unsafe.Pointer(&corner)), 4)
}

func makeBorderless(hwnd uintptr) {
	const gwlStyle = ^uintptr(15)
	style, _, _ := procGetWindowLongPtrW.Call(hwnd, gwlStyle)
	style &^= 0x00C00000 | 0x00040000
	procSetWindowLongPtrW.Call(hwnd, gwlStyle, style)
	procSetWindowPos.Call(hwnd, 0, 0, 0, 0, 0, 0x0001|0x0002|0x0004|0x0020)
}

func beginWindowDrag(hwnd uintptr) {
	procReleaseCapture.Call()
	procSendMessageW.Call(hwnd, 0x00A1, 2, 0)
}

func uiLang() string {
	id, _, _ := procGetUserDefaultUILang.Call()
	switch id & 0x3FF {
	case 0x19:
		return "ru"
	case 0x22:
		return "uk"
	}
	return "en"
}

func runHidden(name string, args ...string) error {
	return makeCmd(name, args...).Run()
}
