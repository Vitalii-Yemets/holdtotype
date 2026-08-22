package evqueue

import (
	"sync"
	"testing"
	"time"
)

func TestOrderPreserved(t *testing.T) {
	q := New[int](100)
	for i := 0; i < 50; i++ {
		if !q.Push(i) {
			t.Fatalf("push %d отклонён", i)
		}
	}
	for i := 0; i < 50; i++ {
		v, ok := q.Pop()
		if !ok || v != i {
			t.Fatalf("получено %d (ok=%v), ожидалось %d", v, ok, i)
		}
	}
	if _, ok := q.Pop(); ok {
		t.Fatal("очередь должна быть пустой")
	}
}

func TestCapDropsAndCounts(t *testing.T) {
	q := New[int](3)
	for i := 0; i < 3; i++ {
		if !q.Push(i) {
			t.Fatalf("push %d отклонён до заполнения", i)
		}
	}
	if q.Push(99) {
		t.Fatal("push сверх предела должен вернуть false")
	}
	if q.Dropped() != 1 {
		t.Fatalf("Dropped = %d, ожидалось 1", q.Dropped())
	}
	if v, _ := q.Pop(); v != 0 {
		t.Fatal("отброшен должен быть новый элемент, а не старый")
	}
	if !q.Push(4) {
		t.Fatal("после освобождения места push должен пройти")
	}
}

func TestSignalWakesConsumer(t *testing.T) {
	q := New[int](10)
	got := make(chan int, 10)
	done := make(chan struct{})
	go func() {
		defer close(done)
		for range q.Signal() {
			for {
				v, ok := q.Pop()
				if !ok {
					break
				}
				got <- v
			}
		}
	}()
	for i := 1; i <= 5; i++ {
		q.Push(i)
		time.Sleep(time.Millisecond)
	}
	for i := 1; i <= 5; i++ {
		select {
		case v := <-got:
			if v != i {
				t.Fatalf("получено %d, ожидалось %d", v, i)
			}
		case <-time.After(2 * time.Second):
			t.Fatalf("потребитель не проснулся на элементе %d", i)
		}
	}
	q.Close()
	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("потребитель не завершился после Close")
	}
}

func TestNoSignalLostOnBurst(t *testing.T) {
	q := New[int](1000)
	var wg sync.WaitGroup
	for p := 0; p < 4; p++ {
		wg.Add(1)
		go func(base int) {
			defer wg.Done()
			for i := 0; i < 100; i++ {
				q.Push(base + i)
			}
		}(p * 1000)
	}
	wg.Wait()
	if q.Len() != 400 {
		t.Fatalf("в очереди %d элементов, ожидалось 400", q.Len())
	}
	seen := map[int]bool{}
	for {
		v, ok := q.Pop()
		if !ok {
			break
		}
		if seen[v] {
			t.Fatalf("элемент %d получен дважды", v)
		}
		seen[v] = true
	}
	if len(seen) != 400 {
		t.Fatalf("извлечено %d элементов, ожидалось 400", len(seen))
	}
}

func TestPushAfterCloseRejected(t *testing.T) {
	q := New[int](5)
	q.Close()
	if q.Push(1) {
		t.Fatal("push после Close должен вернуть false")
	}
	q.Close()
}
