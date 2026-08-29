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
		for _, face := range plexfont.Faces() {
			if len(face.Data) == 0 {
				continue
			}
			var installed uint32
			h, _, _ := procAddFontMemResourceEx.Call(
				uintptr(unsafe.Pointer(&face.Data[0])), uintptr(len(face.Data)), 0,
				uintptr(unsafe.Pointer(&installed)))
			if h == 0 {
				log.Printf("font %s could not be added to the process", face.Family)
				return
			}
		}
		log.Printf("fonts %s and %s added to the process", plexfont.Sans, plexfont.Mono)
	})
}
