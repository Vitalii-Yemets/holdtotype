package main

import (
	"log"

	"holdtotype/internal/evqueue"

	"fmt"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"unsafe"
)

const (
	whKeyboardLL = 13
	wmKeyDown    = 0x0100
	wmKeyUp      = 0x0101
	wmSysKeyDown = 0x0104
	wmSysKeyUp   = 0x0105
)

type kbdllHookStruct struct {
	VkCode    uint32
	ScanCode  uint32
	Flags     uint32
	Time      uint32
	ExtraInfo uintptr
}

type msgStruct struct {
	Hwnd    uintptr
	Message uint32
	_       uint32
	WParam  uintptr
	LParam  uintptr
	Time    uint32
	PtX     int32
	PtY     int32
	_       uint32
}

func parseHotkey(s string) ([][]uint32, error) {
	var groups [][]uint32
	for _, part := range strings.Split(strings.ToLower(s), "+") {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		g, err := keyGroup(part)
		if err != nil {
			return nil, err
		}
		groups = append(groups, g)
	}
	if len(groups) == 0 {
		return nil, fmt.Errorf("пустое сочетание клавиш")
	}
	return groups, nil
}

func keyGroup(name string) ([]uint32, error) {
	switch name {
	case "ctrl", "control":
		return []uint32{0xA2, 0xA3, 0x11}, nil
	case "shift":
		return []uint32{0xA0, 0xA1, 0x10}, nil
	case "alt":
		return []uint32{0xA4, 0xA5, 0x12}, nil
	case "win", "super", "meta":
		return []uint32{0x5B, 0x5C}, nil
	case "space":
		return []uint32{0x20}, nil
	case "tab":
		return []uint32{0x09}, nil
	case "capslock", "caps":
		return []uint32{0x14}, nil
	case "esc", "escape":
		return []uint32{0x1B}, nil
	case "insert", "ins":
		return []uint32{0x2D}, nil
	case "delete", "del":
		return []uint32{0x2E}, nil
	case "home":
		return []uint32{0x24}, nil
	case "end":
		return []uint32{0x23}, nil
	case "pageup":
		return []uint32{0x21}, nil
	case "pagedown":
		return []uint32{0x22}, nil
	case "pause":
		return []uint32{0x13}, nil
	case "scrolllock":
		return []uint32{0x91}, nil
	case "numlock":
		return []uint32{0x90}, nil
	}
	if len(name) >= 2 && name[0] == 'f' {
		if n, err := strconv.Atoi(name[1:]); err == nil && n >= 1 && n <= 24 {
			return []uint32{uint32(0x70 + n - 1)}, nil
		}
	}
	if len(name) == 1 {
		c := name[0]
		if c >= 'a' && c <= 'z' {
			return []uint32{uint32(c) - 'a' + 0x41}, nil
		}
		if c >= '0' && c <= '9' {
			return []uint32{uint32(c) - '0' + 0x30}, nil
		}
	}
	return nil, fmt.Errorf("неизвестная клавиша в сочетании: %q", name)
}

type comboDef struct {
	id     string
	groups [][]uint32
}

const vkEscape = 0x1B

type hotkeyHook struct {
	mu       sync.Mutex
	combos   []comboDef
	pressed  map[uint32]bool
	activeID string
	active   bool
	armed    bool
	escFree  bool
	onDown   func(id string)
	onUp     func()
	onEsc    func() bool
	capture  *captureSession
	q        *evqueue.Queue[func()]
}

func (h *hotkeyHook) post(fn func()) {
	if fn == nil {
		return
	}
	if !h.q.Push(fn) {
		log.Printf("очередь хоткея переполнена, событие отброшено (всего %d)", h.q.Dropped())
	}
}

func (h *hotkeyHook) dispatch() {
	for range h.q.Signal() {
		for {
			fn, ok := h.q.Pop()
			if !ok {
				break
			}
			fn()
		}
	}
}

type captureSession struct {
	pressed  map[uint32]bool
	union    map[uint32]bool
	onUpdate func(live string)
	onDone   func(combo string, ok bool)
}

func (h *hotkeyHook) StartCapture(onUpdate func(string), onDone func(string, bool)) {
	h.mu.Lock()
	var fire func()
	if h.active {
		h.active = false
		h.activeID = ""
		fire = h.onUp
	}
	h.pressed = make(map[uint32]bool)
	h.capture = &captureSession{
		pressed:  make(map[uint32]bool),
		union:    make(map[uint32]bool),
		onUpdate: onUpdate,
		onDone:   onDone,
	}
	h.mu.Unlock()
	h.post(fire)
}

func (h *hotkeyHook) CancelCapture() {
	h.mu.Lock()
	h.capture = nil
	h.pressed = make(map[uint32]bool)
	h.active = false
	h.activeID = ""
	h.mu.Unlock()
}

func (h *hotkeyHook) captureEvent(vk uint32, down bool) bool {
	h.mu.Lock()
	c := h.capture
	if c == nil {
		h.mu.Unlock()
		return false
	}
	if down {
		c.pressed[vk] = true
		c.union[vk] = true
	} else {
		delete(c.pressed, vk)
	}
	var fire func()
	if len(c.pressed) == 0 && len(c.union) > 0 {
		combo := comboString(c.union)
		_, escOnly := c.union[0x1B]
		ok := combo != "" && !(escOnly && len(c.union) == 1)
		h.capture = nil
		h.pressed = make(map[uint32]bool)
		h.active = false
		h.activeID = ""
		fire = func() { c.onDone(combo, ok) }
	} else if down {
		live := comboString(c.union)
		if live == "" {
			live = "…"
		}
		fire = func() { c.onUpdate(live) }
	}
	h.mu.Unlock()
	h.post(fire)
	return true
}

func startHotkeyHook(combos []comboDef, onDown func(id string), onUp func(), onEsc func() bool) (*hotkeyHook, error) {
	h := &hotkeyHook{
		pressed: make(map[uint32]bool),
		armed:   true,
		onDown:  onDown,
		onUp:    onUp,
		onEsc:   onEsc,
		q:       evqueue.New[func()](256),
	}
	h.setCombosLocked(combos)
	go h.dispatch()
	errCh := make(chan error, 1)
	go func() {
		runtime.LockOSThread()
		cb := syscall.NewCallback(h.hookProc)
		hook, _, callErr := procSetWindowsHookExW.Call(whKeyboardLL, cb, 0, 0)
		if hook == 0 {
			errCh <- fmt.Errorf("SetWindowsHookEx: %v", callErr)
			return
		}
		errCh <- nil
		var m msgStruct
		for {
			procGetMessageW.Call(uintptr(unsafe.Pointer(&m)), 0, 0, 0)
		}
	}()
	return h, <-errCh
}

func (h *hotkeyHook) setCombosLocked(combos []comboDef) {
	sorted := make([]comboDef, len(combos))
	copy(sorted, combos)
	sort.SliceStable(sorted, func(i, j int) bool {
		return len(sorted[i].groups) > len(sorted[j].groups)
	})
	h.combos = sorted
	h.escFree = true
	for _, c := range sorted {
		for _, group := range c.groups {
			for _, vk := range group {
				if vk == vkEscape {
					h.escFree = false
				}
			}
		}
	}
}

func (h *hotkeyHook) escCancels() bool {
	h.mu.Lock()
	defer h.mu.Unlock()
	return h.escFree && h.capture == nil && h.onEsc != nil
}

func (h *hotkeyHook) SetCombos(combos []comboDef) {
	h.mu.Lock()
	h.setCombosLocked(combos)
	var fire func()
	if h.active && h.matchedCombo() == "" {
		h.active = false
		h.activeID = ""
		h.armed = false
		fire = h.onUp
	}
	h.mu.Unlock()
	h.post(fire)
}

func (h *hotkeyHook) comboPressed(id string) bool {
	for _, c := range h.combos {
		if c.id == id {
			return h.groupsPressed(c.groups)
		}
	}
	return false
}

func (h *hotkeyHook) matchedCombo() string {
	for _, c := range h.combos {
		if h.groupsPressed(c.groups) {
			return c.id
		}
	}
	return ""
}

func (h *hotkeyHook) groupsPressed(groups [][]uint32) bool {
	if len(groups) == 0 {
		return false
	}
	for _, group := range groups {
		ok := false
		for _, vk := range group {
			if h.pressed[vk] {
				ok = true
				break
			}
		}
		if !ok {
			return false
		}
	}
	return true
}

func (h *hotkeyHook) hookProc(nCode, wParam, lParam uintptr) uintptr {
	if int32(nCode) == 0 {
		k := (*kbdllHookStruct)(unsafe.Pointer(lParam))
		if k.ExtraInfo != injectedMarker {
			down := wParam == wmKeyDown || wParam == wmSysKeyDown
			up := wParam == wmKeyUp || wParam == wmSysKeyUp
			if down || up {
				if h.captureEvent(k.VkCode, down) {
					return 1
				}
				if k.VkCode == vkEscape && h.escCancels() {
					if down && h.onEsc() {
						return 1
					}
				}
				h.keyEvent(k.VkCode, down)
			}
		}
	}
	r, _, _ := procCallNextHookEx.Call(0, nCode, wParam, lParam)
	return r
}

func (h *hotkeyHook) keyEvent(vk uint32, down bool) {
	h.mu.Lock()
	if down {
		h.pressed[vk] = true
	} else {
		delete(h.pressed, vk)
	}
	var fire func()
	if h.active {
		if !h.comboPressed(h.activeID) {
			h.active = false
			h.activeID = ""
			h.armed = false
			fire = h.onUp
		}
	}
	if !h.active && fire == nil {
		if !h.armed {
			if h.matchedCombo() == "" {
				h.armed = true
			}
		} else if id := h.matchedCombo(); id != "" {
			h.active = true
			h.activeID = id
			down := h.onDown
			fire = func() { down(id) }
		}
	}
	h.mu.Unlock()
	h.post(fire)
}
