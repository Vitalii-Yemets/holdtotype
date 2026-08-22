package evqueue

import "sync"

type Queue[T any] struct {
	mu      sync.Mutex
	buf     []T
	sig     chan struct{}
	max     int
	dropped int
	closed  bool
}

func New[T any](max int) *Queue[T] {
	if max <= 0 {
		max = 1
	}
	return &Queue[T]{sig: make(chan struct{}, 1), max: max}
}

func (q *Queue[T]) Push(v T) bool {
	q.mu.Lock()
	if q.closed {
		q.mu.Unlock()
		return false
	}
	if len(q.buf) >= q.max {
		q.dropped++
		q.mu.Unlock()
		return false
	}
	q.buf = append(q.buf, v)
	q.mu.Unlock()
	select {
	case q.sig <- struct{}{}:
	default:
	}
	return true
}

func (q *Queue[T]) Pop() (T, bool) {
	var zero T
	q.mu.Lock()
	defer q.mu.Unlock()
	if len(q.buf) == 0 {
		return zero, false
	}
	v := q.buf[0]
	q.buf[0] = zero
	q.buf = q.buf[1:]
	return v, true
}

func (q *Queue[T]) Signal() <-chan struct{} { return q.sig }

func (q *Queue[T]) Len() int {
	q.mu.Lock()
	defer q.mu.Unlock()
	return len(q.buf)
}

func (q *Queue[T]) Dropped() int {
	q.mu.Lock()
	defer q.mu.Unlock()
	return q.dropped
}

func (q *Queue[T]) Close() {
	q.mu.Lock()
	if q.closed {
		q.mu.Unlock()
		return
	}
	q.closed = true
	q.mu.Unlock()
	close(q.sig)
}
