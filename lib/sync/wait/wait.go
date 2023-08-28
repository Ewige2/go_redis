package wait

import (
	"sync"
	"time"
)

// Wait is similar with sync.WaitGroup which can wait with timeout
type Wait struct {
	wg sync.WaitGroup
}

// Add  将 delta（可能为负）添加到 WaitGroup 计数器。
func (w *Wait) Add(delta int) {
	w.wg.Add(delta)
}

// Done 将 WaitGroup 计数器减 1
func (w *Wait) Done() {
	w.wg.Done()
}

// Wait 阻塞直到 WaitGroup 计数器为零或超时
func (w *Wait) Wait() {
	w.wg.Wait()
}

// WaitWithTimeout 阻塞直到 WaitGroup 计数器为零或超时
// returns true if timeout
func (w *Wait) WaitWithTimeout(timeout time.Duration) bool {
	c := make(chan bool, 1)
	go func() {
		defer close(c)
		w.wg.Wait()
		c <- true
	}()
	select {
	case <-c:
		return false // completed normally
	case <-time.After(timeout):
		return true // timed out
	}
}
