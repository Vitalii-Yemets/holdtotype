package appname

import (
	"fmt"
	"unsafe"

	"golang.org/x/sys/windows"
)

func ReadVersion(path string) (description, product string) {
	size, err := windows.GetFileVersionInfoSize(path, nil)
	if err != nil || size == 0 {
		return "", ""
	}
	buf := make([]byte, size)
	if err := windows.GetFileVersionInfo(path, 0, size, unsafe.Pointer(&buf[0])); err != nil {
		return "", ""
	}
	block := unsafe.Pointer(&buf[0])
	tables := []string{}
	var p unsafe.Pointer
	var n uint32
	if windows.VerQueryValue(block, `\VarFileInfo\Translation`, unsafe.Pointer(&p), &n) == nil && n >= 4 && p != nil {
		pair := (*[2]uint16)(p)
		tables = append(tables, fmt.Sprintf("%04x%04x", pair[0], pair[1]))
	}
	tables = append(tables, "040904b0", "040904e4", "04090000")
	value := func(table, key string) string {
		var v unsafe.Pointer
		var l uint32
		if windows.VerQueryValue(block, `\StringFileInfo\`+table+`\`+key, unsafe.Pointer(&v), &l) != nil || l == 0 || v == nil {
			return ""
		}
		return windows.UTF16PtrToString((*uint16)(v))
	}
	for _, t := range tables {
		description, product = value(t, "FileDescription"), value(t, "ProductName")
		if description != "" || product != "" {
			return description, product
		}
	}
	return "", ""
}
