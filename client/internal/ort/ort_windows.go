package ort

import (
	"errors"
	"fmt"
	"runtime"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

const apiVersion = 20

const (
	fnGetErrorMessage                  = 2
	fnCreateEnv                        = 3
	fnCreateSession                    = 7
	fnRun                              = 9
	fnDisableMemPattern                = 17
	fnDisableCpuMemArena               = 19
	fnCreateSessionOptions             = 10
	fnSetSessionGraphOptimizationLevel = 23
	fnSetIntraOpNumThreads             = 24
	fnSetInterOpNumThreads             = 25
	fnCreateTensorWithDataAsOrtValue   = 49
	fnGetTensorMutableData             = 51
	fnGetDimensionsCount               = 61
	fnGetDimensions                    = 62
	fnGetTensorTypeAndShape            = 65
	fnCreateCpuMemoryInfo              = 69
	fnReleaseEnv                       = 92
	fnReleaseStatus                    = 93
	fnReleaseMemoryInfo                = 94
	fnReleaseSession                   = 95
	fnReleaseValue                     = 96
	fnReleaseTensorTypeAndShapeInfo    = 99
	fnReleaseSessionOptions            = 100
)

const (
	typeInt64      = 7
	loggingError   = 3
	optimizeAll    = 99
	arenaAllocator = 1
	memTypeDefault = 0
)

type Session struct {
	api  uintptr
	env  uintptr
	sess uintptr
	mem  uintptr
}

func (s *Session) fn(i int) uintptr {
	return *(*uintptr)(unsafe.Pointer(s.api + uintptr(i)*unsafe.Sizeof(uintptr(0))))
}

//go:uintptrescapes
func (s *Session) call(i int, args ...uintptr) error {
	st, _, _ := syscall.SyscallN(s.fn(i), args...)
	if st == 0 {
		return nil
	}
	msg, _, _ := syscall.SyscallN(s.fn(fnGetErrorMessage), st)
	text := cString(msg)
	syscall.SyscallN(s.fn(fnReleaseStatus), st)
	return errors.New(text)
}

//go:uintptrescapes
func (s *Session) release(i int, p uintptr) {
	if p != 0 {
		syscall.SyscallN(s.fn(i), p)
	}
}

func cString(p uintptr) string {
	if p == 0 {
		return ""
	}
	var b []byte
	for i := uintptr(0); i < 4096; i++ {
		c := *(*byte)(unsafe.Pointer(p + i))
		if c == 0 {
			break
		}
		b = append(b, c)
	}
	return string(b)
}

func cBytes(s string) []byte { return append([]byte(s), 0) }

type Options struct {
	Threads   int
	NoArena   bool
	NoPattern bool
}

func Open(dll, model string, o Options) (*Session, error) {
	d, err := windows.LoadDLL(dll)
	if err != nil {
		return nil, fmt.Errorf("%s cannot be loaded: %w", dll, err)
	}
	p, err := d.FindProc("OrtGetApiBase")
	if err != nil {
		return nil, fmt.Errorf("%s has no OrtGetApiBase: %w", dll, err)
	}
	base, _, _ := p.Call()
	if base == 0 {
		return nil, errors.New("OrtGetApiBase gave nothing")
	}
	getAPI := *(*uintptr)(unsafe.Pointer(base))
	api, _, _ := syscall.SyscallN(getAPI, apiVersion)
	if api == 0 {
		return nil, fmt.Errorf("%s does not offer the API version %d", dll, apiVersion)
	}
	s := &Session{api: api}
	logid := cBytes("punct")
	if err := s.call(fnCreateEnv, loggingError, uintptr(unsafe.Pointer(&logid[0])), uintptr(unsafe.Pointer(&s.env))); err != nil {
		return nil, fmt.Errorf("the runtime did not start: %w", err)
	}
	var opts uintptr
	if err := s.call(fnCreateSessionOptions, uintptr(unsafe.Pointer(&opts))); err != nil {
		s.Close()
		return nil, err
	}
	defer s.release(fnReleaseSessionOptions, opts)
	threads := o.Threads
	if threads < 1 {
		threads = 1
	}
	for _, step := range []struct {
		fn  int
		arg uintptr
	}{{fnSetIntraOpNumThreads, uintptr(threads)}, {fnSetInterOpNumThreads, 1}, {fnSetSessionGraphOptimizationLevel, optimizeAll}} {
		if err := s.call(step.fn, opts, step.arg); err != nil {
			s.Close()
			return nil, err
		}
	}
	if o.NoArena {
		if err := s.call(fnDisableCpuMemArena, opts); err != nil {
			s.Close()
			return nil, err
		}
	}
	if o.NoPattern {
		if err := s.call(fnDisableMemPattern, opts); err != nil {
			s.Close()
			return nil, err
		}
	}
	path, err := windows.UTF16PtrFromString(model)
	if err != nil {
		s.Close()
		return nil, err
	}
	if err := s.call(fnCreateSession, s.env, uintptr(unsafe.Pointer(path)), opts, uintptr(unsafe.Pointer(&s.sess))); err != nil {
		s.Close()
		return nil, fmt.Errorf("the model %s cannot be opened: %w", model, err)
	}
	if err := s.call(fnCreateCpuMemoryInfo, arenaAllocator, memTypeDefault, uintptr(unsafe.Pointer(&s.mem))); err != nil {
		s.Close()
		return nil, err
	}
	return s, nil
}

func (s *Session) Close() {
	s.release(fnReleaseMemoryInfo, s.mem)
	s.release(fnReleaseSession, s.sess)
	s.release(fnReleaseEnv, s.env)
	s.mem, s.sess, s.env = 0, 0, 0
}

func (s *Session) tensor(data []int64, shape []int64) (uintptr, error) {
	var v uintptr
	err := s.call(fnCreateTensorWithDataAsOrtValue, s.mem, uintptr(unsafe.Pointer(&data[0])), uintptr(len(data)*8),
		uintptr(unsafe.Pointer(&shape[0])), uintptr(len(shape)), typeInt64, uintptr(unsafe.Pointer(&v)))
	return v, err
}

func (s *Session) Logits(ids []int64) ([]float32, []int64, error) {
	if len(ids) == 0 {
		return nil, nil, errors.New("nothing to read")
	}
	mask := make([]int64, len(ids))
	for i := range mask {
		mask[i] = 1
	}
	shape := []int64{1, int64(len(ids))}
	n0, n1, o0 := cBytes("input_ids"), cBytes("attention_mask"), cBytes("logits")
	var pin runtime.Pinner
	defer pin.Unpin()
	for _, p := range []unsafe.Pointer{unsafe.Pointer(&ids[0]), unsafe.Pointer(&mask[0]), unsafe.Pointer(&shape[0]), unsafe.Pointer(&n0[0]), unsafe.Pointer(&n1[0]), unsafe.Pointer(&o0[0])} {
		pin.Pin(p)
	}
	in0, err := s.tensor(ids, shape)
	if err != nil {
		return nil, nil, err
	}
	defer s.release(fnReleaseValue, in0)
	in1, err := s.tensor(mask, shape)
	if err != nil {
		return nil, nil, err
	}
	defer s.release(fnReleaseValue, in1)
	inNames := []uintptr{uintptr(unsafe.Pointer(&n0[0])), uintptr(unsafe.Pointer(&n1[0]))}
	outNames := []uintptr{uintptr(unsafe.Pointer(&o0[0]))}
	inVals := []uintptr{in0, in1}
	outVals := []uintptr{0}
	for _, p := range []unsafe.Pointer{unsafe.Pointer(&inNames[0]), unsafe.Pointer(&outNames[0]), unsafe.Pointer(&inVals[0]), unsafe.Pointer(&outVals[0])} {
		pin.Pin(p)
	}
	err = s.call(fnRun, s.sess, 0, uintptr(unsafe.Pointer(&inNames[0])), uintptr(unsafe.Pointer(&inVals[0])), 2,
		uintptr(unsafe.Pointer(&outNames[0])), 1, uintptr(unsafe.Pointer(&outVals[0])))
	if err != nil {
		return nil, nil, err
	}
	out := outVals[0]
	defer s.release(fnReleaseValue, out)
	var info uintptr
	if err := s.call(fnGetTensorTypeAndShape, out, uintptr(unsafe.Pointer(&info))); err != nil {
		return nil, nil, err
	}
	defer s.release(fnReleaseTensorTypeAndShapeInfo, info)
	var nd uintptr
	if err := s.call(fnGetDimensionsCount, info, uintptr(unsafe.Pointer(&nd))); err != nil {
		return nil, nil, err
	}
	if nd == 0 || nd > 8 {
		return nil, nil, fmt.Errorf("the answer has %d dimensions", nd)
	}
	dims := make([]int64, nd)
	if err := s.call(fnGetDimensions, info, uintptr(unsafe.Pointer(&dims[0])), nd); err != nil {
		return nil, nil, err
	}
	count := int64(1)
	for _, d := range dims {
		count *= d
	}
	var data uintptr
	if err := s.call(fnGetTensorMutableData, out, uintptr(unsafe.Pointer(&data))); err != nil {
		return nil, nil, err
	}
	res := make([]float32, count)
	copy(res, unsafe.Slice((*float32)(unsafe.Pointer(data)), count))
	return res, dims, nil
}
