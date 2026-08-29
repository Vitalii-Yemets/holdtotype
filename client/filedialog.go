package main

import (
	"runtime"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
	comdlg32             = windows.NewLazySystemDLL("comdlg32.dll")
	procGetOpenFileNameW = comdlg32.NewProc("GetOpenFileNameW")
	procGetSaveFileNameW = comdlg32.NewProc("GetSaveFileNameW")
)

const (
	ofnHideReadOnly     = 0x00000004
	ofnOverwritePrompt  = 0x00000002
	ofnPathMustExist    = 0x00000800
	ofnFileMustExist    = 0x00001000
	ofnExplorer         = 0x00080000
	ofnNoChangeDir      = 0x00000008
	fileDialogPathChars = 520
)

type openFileNameW struct {
	StructSize    uint32
	Owner         uintptr
	Instance      uintptr
	Filter        *uint16
	CustomFilter  *uint16
	MaxCustFilter uint32
	FilterIndex   uint32
	File          *uint16
	MaxFile       uint32
	FileTitle     *uint16
	MaxFileTitle  uint32
	InitialDir    *uint16
	Title         *uint16
	Flags         uint32
	FileOffset    uint16
	FileExtension uint16
	DefExt        *uint16
	CustData      uintptr
	FnHook        uintptr
	TemplateName  *uint16
	PvReserved    uintptr
	DwReserved    uint32
	FlagsEx       uint32
}

func jsonFilter() []uint16 {
	parts := []string{"JSON", "*.json", "", ""}
	var buf []uint16
	for _, p := range parts {
		buf = append(buf, utf16Of(p)...)
	}
	return buf
}

func utf16Of(s string) []uint16 {
	u, err := windows.UTF16FromString(s)
	if err != nil {
		return []uint16{0}
	}
	return u
}

func askFilePath(save bool, title, suggested string) string {
	done := make(chan string, 1)
	go func() {
		runtime.LockOSThread()
		defer runtime.UnlockOSThread()
		if err := windows.CoInitializeEx(0, windows.COINIT_APARTMENTTHREADED); err == nil {
			defer windows.CoUninitialize()
		}
		done <- runFileDialog(save, title, suggested)
	}()
	return <-done
}

func runFileDialog(save bool, title, suggested string) string {
	buf := make([]uint16, fileDialogPathChars)
	copy(buf, utf16Of(suggested))
	filter := jsonFilter()
	ofn := openFileNameW{
		Filter:      &filter[0],
		FilterIndex: 1,
		File:        &buf[0],
		MaxFile:     uint32(len(buf)),
		Title:       syscall.StringToUTF16Ptr(title),
		DefExt:      syscall.StringToUTF16Ptr("json"),
		Flags:       ofnExplorer | ofnHideReadOnly | ofnPathMustExist | ofnNoChangeDir,
	}
	ofn.StructSize = uint32(unsafe.Sizeof(ofn))
	if save {
		ofn.Flags |= ofnOverwritePrompt
	} else {
		ofn.Flags |= ofnFileMustExist
	}
	proc := procGetOpenFileNameW
	if save {
		proc = procGetSaveFileNameW
	}
	ret, _, _ := proc.Call(uintptr(unsafe.Pointer(&ofn)))
	if ret == 0 {
		return ""
	}
	return windows.UTF16ToString(buf)
}
