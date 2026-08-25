package main

import (
	"log"
	"sync"
	"unsafe"

	"holdtotype/internal/plexfont"
)

var (
	procAddFontMemResourceEx = gdi32.NewProc("AddFontMemResourceEx")
	bundledFontsOnce         sync.Once
)

func registerBundledFonts() {
	bundledFontsOnce.Do(func() {
		for _, data := range [][]byte{plexfont.Regular, plexfont.SemiBold} {
			if len(data) == 0 {
				continue
			}
			var installed uint32
			h, _, _ := procAddFontMemResourceEx.Call(
				uintptr(unsafe.Pointer(&data[0])), uintptr(len(data)), 0,
				uintptr(unsafe.Pointer(&installed)))
			if h == 0 {
				log.Printf("шрифт %s не удалось подключить к процессу", plexfont.Family)
				return
			}
		}
		log.Printf("шрифт %s подключён к процессу", plexfont.Family)
	})
}
