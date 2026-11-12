package gpucard

import (
	"errors"
	"fmt"
	"strings"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
	dxcoreDLL         = windows.NewLazySystemDLL("dxcore.dll")
	procCreateFactory = dxcoreDLL.NewProc("DXCoreCreateAdapterFactory")

	iidAdapterFactory = windows.GUID{Data1: 0x78ee5945, Data2: 0xc36e, Data3: 0x4b13, Data4: [8]byte{0xa6, 0x69, 0x00, 0x5d, 0xd1, 0x1c, 0x0f, 0x06}}
	iidAdapterList    = windows.GUID{Data1: 0x526c7776, Data2: 0x40e9, Data3: 0x459b, Data4: [8]byte{0xb7, 0x11, 0xf3, 0x2a, 0xd7, 0x6d, 0xfc, 0x28}}
	iidAdapter        = windows.GUID{Data1: 0xf0db4c7f, Data2: 0xfe5a, Data3: 0x42a2, Data4: [8]byte{0xbd, 0x62, 0xf2, 0xa6, 0xcf, 0x6f, 0xc8, 0x3e}}
	attrCoreCompute   = windows.GUID{Data1: 0x248e2800, Data2: 0xa793, Data3: 0x4724, Data4: [8]byte{0xab, 0xaa, 0x23, 0xa6, 0xde, 0x1b, 0xe0, 0x90}}
)

const (
	propDriverDescription = 2
	propDedicatedMemory   = 7
	propIsHardware        = 11
	propIsIntegrated      = 12
)

func call(obj uintptr, index int, args ...uintptr) uintptr {
	vtbl := *(*uintptr)(unsafe.Pointer(obj))
	fn := *(*uintptr)(unsafe.Pointer(vtbl + uintptr(index)*unsafe.Sizeof(uintptr(0))))
	r, _, _ := syscall.SyscallN(fn, append([]uintptr{obj}, args...)...)
	return r
}

func release(obj uintptr) {
	if obj != 0 {
		call(obj, 2)
	}
}

func failed(hr uintptr) bool { return int32(hr) < 0 }

func Adapters() ([]Adapter, error) {
	if err := procCreateFactory.Find(); err != nil {
		return nil, errors.New("this Windows cannot describe its graphics adapters (no DXCore)")
	}
	var factory uintptr
	hr, _, _ := procCreateFactory.Call(uintptr(unsafe.Pointer(&iidAdapterFactory)), uintptr(unsafe.Pointer(&factory)))
	if failed(hr) || factory == 0 {
		return nil, fmt.Errorf("DXCore factory: 0x%08x", uint32(hr))
	}
	defer release(factory)
	var list uintptr
	hr = call(factory, 3, 1, uintptr(unsafe.Pointer(&attrCoreCompute)), uintptr(unsafe.Pointer(&iidAdapterList)), uintptr(unsafe.Pointer(&list)))
	if failed(hr) || list == 0 {
		return nil, fmt.Errorf("DXCore adapter list: 0x%08x", uint32(hr))
	}
	defer release(list)
	count := uint32(call(list, 4))
	var out []Adapter
	for i := uint32(0); i < count; i++ {
		var ad uintptr
		if hr := call(list, 3, uintptr(i), uintptr(unsafe.Pointer(&iidAdapter)), uintptr(unsafe.Pointer(&ad))); failed(hr) || ad == 0 {
			continue
		}
		out = append(out, describe(ad))
		release(ad)
	}
	return out, nil
}

func describe(ad uintptr) Adapter {
	var a Adapter
	var mem uint64
	if !failed(call(ad, 6, propDedicatedMemory, 8, uintptr(unsafe.Pointer(&mem)))) {
		a.DedicatedMB = int(mem / (1024 * 1024))
	}
	var flag byte
	if !failed(call(ad, 6, propIsHardware, 1, uintptr(unsafe.Pointer(&flag)))) {
		a.Hardware = flag != 0
	}
	flag = 1
	if !failed(call(ad, 6, propIsIntegrated, 1, uintptr(unsafe.Pointer(&flag)))) {
		a.Integrated = flag != 0
	} else {
		a.Integrated = true
	}
	var size uintptr
	if !failed(call(ad, 7, propDriverDescription, uintptr(unsafe.Pointer(&size)))) && size > 0 && size < 4096 {
		buf := make([]byte, size)
		if !failed(call(ad, 6, propDriverDescription, size, uintptr(unsafe.Pointer(&buf[0])))) {
			a.Name = strings.TrimSpace(strings.TrimRight(string(buf), "\x00"))
		}
	}
	return a
}
