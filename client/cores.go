package main

import (
	"encoding/binary"
	"errors"
	"fmt"
	"log"
	"os/exec"
	"sync"
	"unsafe"

	"golang.org/x/sys/windows"

	"holdtotype/internal/cores"
)

var (
	procGetLogicalProcessorInformationEx = kernel32.NewProc("GetLogicalProcessorInformationEx")
	procGetActiveProcessorGroupCount     = kernel32.NewProc("GetActiveProcessorGroupCount")
	procSetProcessAffinityMask           = kernel32.NewProc("SetProcessAffinityMask")
	procGetProcessAffinityMask           = kernel32.NewProc("GetProcessAffinityMask")
)

const relationProcessorCore = 0

func readCores() ([]cores.Core, error) {
	if n, _, _ := procGetActiveProcessorGroupCount.Call(); uint16(n) > 1 {
		return nil, fmt.Errorf("the processor has %d processor groups", uint16(n))
	}
	var size uint32
	procGetLogicalProcessorInformationEx.Call(relationProcessorCore, 0, uintptr(unsafe.Pointer(&size)))
	if size == 0 {
		return nil, errors.New("the processor layout cannot be read")
	}
	buf := make([]byte, size)
	r, _, err := procGetLogicalProcessorInformationEx.Call(relationProcessorCore, uintptr(unsafe.Pointer(&buf[0])), uintptr(unsafe.Pointer(&size)))
	if r == 0 {
		return nil, fmt.Errorf("the processor layout cannot be read: %v", err)
	}
	var out []cores.Core
	for off := 0; off+8 <= int(size); {
		rel := binary.LittleEndian.Uint32(buf[off:])
		n := int(binary.LittleEndian.Uint32(buf[off+4:]))
		if n <= 0 || off+n > int(size) {
			break
		}
		if rel == relationProcessorCore && n >= 48 {
			if groups := binary.LittleEndian.Uint16(buf[off+30:]); groups != 1 {
				return nil, fmt.Errorf("a core spans %d processor groups", groups)
			}
			if group := binary.LittleEndian.Uint16(buf[off+40:]); group != 0 {
				return nil, errors.New("the processor has more than one processor group")
			}
			out = append(out, cores.Core{Mask: binary.LittleEndian.Uint64(buf[off+32:]), Class: int(buf[off+9])})
		}
		off += n
	}
	if len(out) == 0 {
		return nil, errors.New("the processor layout is empty")
	}
	return out, nil
}

func systemMask() uint64 {
	var proc, sys uintptr
	r, _, _ := procGetProcessAffinityMask.Call(uintptr(windows.CurrentProcess()), uintptr(unsafe.Pointer(&proc)), uintptr(unsafe.Pointer(&sys)))
	if r == 0 {
		return 0
	}
	return uint64(sys)
}

func setAffinity(pid int, mask uint64) error {
	if pid <= 0 || mask == 0 {
		return nil
	}
	h, err := windows.OpenProcess(windows.PROCESS_SET_INFORMATION|windows.PROCESS_QUERY_LIMITED_INFORMATION, false, uint32(pid))
	if err != nil {
		return err
	}
	defer windows.CloseHandle(h)
	if r, _, e := procSetProcessAffinityMask.Call(uintptr(h), uintptr(mask)); r == 0 {
		return e
	}
	return nil
}

func cmdPid(cmd *exec.Cmd) int {
	if cmd == nil || cmd.Process == nil {
		return 0
	}
	return cmd.Process.Pid
}

func recognizerPid(rec recognizer) int {
	switch s := rec.(type) {
	case *whisperServer:
		return cmdPid(s.cmd)
	case *sherpaServer:
		return cmdPid(s.cmd)
	case *streamServer:
		return cmdPid(s.cmd)
	}
	return 0
}

func (s *llamaServer) pid() int { return cmdPid(s.cmd) }

type coreDivider struct {
	mu      sync.Mutex
	read    bool
	split   cores.Split
	why     string
	said    bool
	full    uint64
	divided bool
	closed  bool
	asrPid  int
	llmPid  int
	offWhy  string
}

func (c *coreDivider) apply(want bool, off string, threads, asr, llm int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.closed {
		return
	}
	if !want {
		c.releaseLocked()
		if off != "" && off != c.offWhy {
			log.Printf("subtitles: cores are not divided: %s", off)
		}
		c.offWhy = off
		return
	}
	c.offWhy = ""
	if !c.read {
		c.read = true
		c.full = systemMask()
		list, err := readCores()
		if err != nil {
			c.why = err.Error()
		} else {
			c.split, c.why = cores.Divide(list, threads)
		}
		if c.why == "" && c.full == 0 {
			c.why = "the full processor mask cannot be read"
		}
	}
	if c.why != "" {
		if !c.said {
			c.said = true
			log.Printf("subtitles: cores are not divided: %s", c.why)
		}
		return
	}
	if asr != 0 && asr != c.asrPid {
		if err := setAffinity(asr, c.split.Recognition); err != nil {
			c.giveUpLocked(fmt.Sprintf("the recognition engine cannot be pinned: %v", err))
			return
		}
		c.asrPid = asr
	}
	if llm != 0 && llm != c.llmPid {
		if err := setAffinity(llm, c.split.Translation); err != nil {
			c.giveUpLocked(fmt.Sprintf("llama-server cannot be pinned: %v", err))
			return
		}
		c.llmPid = llm
	}
	if !c.divided && (c.asrPid != 0 || c.llmPid != 0) {
		c.divided = true
		log.Printf("subtitles: cores divided: recognition on %d logical processors, translation on %d", cores.Count(c.split.Recognition), cores.Count(c.split.Translation))
	}
}

func (c *coreDivider) giveUpLocked(why string) {
	c.releaseLocked()
	c.why = why
	c.said = true
	log.Printf("subtitles: cores are not divided: %s", why)
}

func (c *coreDivider) releaseLocked() {
	if c.asrPid != 0 {
		_ = setAffinity(c.asrPid, c.full)
	}
	if c.llmPid != 0 {
		_ = setAffinity(c.llmPid, c.full)
	}
	was := c.divided
	c.asrPid, c.llmPid, c.divided = 0, 0, false
	if was {
		log.Printf("subtitles: cores are shared again")
	}
}

func (c *coreDivider) close() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.releaseLocked()
	c.closed = true
}

func (r *subRun) syncCores() {
	cfg := r.app.snapshot()
	t := r.target()
	translating := t != nil && translationTarget(cfg, t) != ""
	kind := translationSource(cfg).kind
	want := translating && kind == editorLocal
	off := ""
	switch {
	case translating && kind == editorAPI:
		off = "the user's post-processing server translates"
	case translating && kind == editorQuality:
		off = "the graphics card translates"
	}
	r.mu.Lock()
	asr := r.engPid
	r.mu.Unlock()
	llm := 0
	if s := r.app.editorEng.live(); s != nil {
		llm = s.pid()
	}
	r.cores.apply(want, off, cfg.Threads, asr, llm)
}
