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
		return nil, fmt.Errorf("empty key combination")
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
	case "backspace":
		return []uint32{0x08}, nil
	case "enter":
		return []uint32{0x0D}, nil
	case "printscreen":
		return []uint32{0x2C}, nil
	case "menu":
		return []uint32{0x5D}, nil
	case "left":
		return []uint32{0x25}, nil
	case "up":
		return []uint32{0x26}, nil
	case "right":
		return []uint32{0x27}, nil
	case "down":
		return []uint32{0x28}, nil
	case "num0":
		return []uint32{0x60}, nil
	case "num1":
		return []uint32{0x61}, nil
	case "num2":
		return []uint32{0x62}, nil
	case "num3":
		return []uint32{0x63}, nil
	case "num4":
		return []uint32{0x64}, nil
	case "num5":
		return []uint32{0x65}, nil
	case "num6":
		return []uint32{0x66}, nil
	case "num7":
		return []uint32{0x67}, nil
	case "num8":
		return []uint32{0x68}, nil
	case "num9":
		return []uint32{0x69}, nil
	case "multiply":
		return []uint32{0x6A}, nil
	case "add":
		return []uint32{0x6B}, nil
	case "subtract":
		return []uint32{0x6D}, nil
	case "decimal":
		return []uint32{0x6E}, nil
	case "divide":
		return []uint32{0x6F}, nil
	case "semicolon":
		return []uint32{0xBA}, nil
	case "equals":
		return []uint32{0xBB}, nil
	case "comma":
		return []uint32{0xBC}, nil
	case "minus":
		return []uint32{0xBD}, nil
	case "period":
		return []uint32{0xBE}, nil
	case "slash":
		return []uint32{0xBF}, nil
	case "backquote":
		return []uint32{0xC0}, nil
	case "lbracket":
		return []uint32{0xDB}, nil
	case "backslash":
		return []uint32{0xDC}, nil
	case "rbracket":
		return []uint32{0xDD}, nil
	case "quote":
		return []uint32{0xDE}, nil
	case "volumemute":
		return []uint32{0xAD}, nil
	case "volumedown":
		return []uint32{0xAE}, nil
	case "volumeup":
		return []uint32{0xAF}, nil
	case "medianext":
		return []uint32{0xB0}, nil
	case "mediaprev":
		return []uint32{0xB1}, nil
	case "mediastop":
		return []uint32{0xB2}, nil
	case "mediaplay":
		return []uint32{0xB3}, nil
	case "browserback":
		return []uint32{0xA6}, nil
	case "browserforward":
		return []uint32{0xA7}, nil
	case "browserrefresh":
		return []uint32{0xA8}, nil
	case "browserstop":
		return []uint32{0xA9}, nil
	case "browsersearch":
		return []uint32{0xAA}, nil
	case "browserfavorites":
		return []uint32{0xAB}, nil
	case "browserhome":
		return []uint32{0xAC}, nil
	case "launchmail":
		return []uint32{0xB4}, nil
	case "launchmedia":
		return []uint32{0xB5}, nil
	case "launchapp1":
		return []uint32{0xB6}, nil
	case "launchapp2":
		return []uint32{0xB7}, nil
	}
	if strings.HasPrefix(name, "vk0x") && len(name) == 6 {
		if n, err := strconv.ParseUint(name[4:], 16, 32); err == nil && n > 0 && n < 256 {
			return []uint32{uint32(n)}, nil
		}
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
	return nil, fmt.Errorf("unknown key in the combination: %q", name)
}

type comboDef struct {
	id     string
	groups [][]uint32
}

const vkEscape = 0x1B

const vkReturn = 0x0D

const vkUp = 0x26

const vkDown = 0x28

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
	handle   uintptr
	q        *evqueue.Queue[func()]
}

func (h *hotkeyHook) post(fn func()) {
	if fn == nil {
		return
	}
	if !h.q.Push(fn) {
		log.Printf("hotkey queue is full, event dropped (%d in total)", h.q.Dropped())
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
		h.mu.Lock()
		h.handle = hook
		h.mu.Unlock()
		errCh <- nil
		var m msgStruct
		for {
			procGetMessageW.Call(uintptr(unsafe.Pointer(&m)), 0, 0, 0)
		}
	}()
	return h, <-errCh
}

// release takes the low-level keyboard hook back off the system. Windows cleans
// it up on its own when the process ends, but a hook left hanging while the
// text services rewire themselves is what makes ctfmon.exe fall over.
func (h *hotkeyHook) release() {
	if h == nil {
		return
	}
	h.mu.Lock()
	hook := h.handle
	h.handle = 0
	h.mu.Unlock()
	if hook != 0 {
		procUnhookWindowsHookEx.Call(hook)
	}
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
				if down && askActive() && askKey(k.VkCode) {
					return 1
				}
				if down && tmActive() && tmKey(k.VkCode) {
					return 1
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
