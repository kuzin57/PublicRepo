package awg

// !NB! для курса shad-go
// взято из моей предыдущей работы

import (
	"context"
	"fmt"
	"gitlab.com/slon/shad-go/gitfame/pkg/errors"
	"runtime"
	"sync"
	"time"
)

const (
	// StatusIdle means that WG did not run yet
	StatusIdle int = iota
	// StatusSuccess means successful execution of all tasks
	StatusSuccess
	// StatusTimeout means that job was broken by timeout
	StatusTimeout
	// StatusError means that job was broken by error in one task (if stopOnError is true)
	StatusError

	errTimeoutMessage = "Wait group timeout after %v"
	stackBufferSize   = 1000
)

// ErrorTimeout error on timeout
type ErrorTimeout time.Duration

// Error implementation
func (e ErrorTimeout) Error() string {
	return fmt.Sprintf(errTimeoutMessage, time.Duration(e).String())
}

// WaitgroupFunc func
type WaitgroupFunc func() error

// AdvancedWaitGroup enhanced wait group struct
type AdvancedWaitGroup struct {
	waitGroupStatus
	stackBuffer []WaitgroupFunc
	receiver    chan struct{}
	sender      chan WaitgroupFunc
	capacity    uint32
	length      int
	timeout     *time.Duration
	ctx         context.Context
	done        func() <-chan struct{}
	stopOnError bool
	errors      []error
	finished    chan struct{}
	onError     func()
}

type waitGroupStatus struct {
	status     int
	statusLock sync.RWMutex
}

func done() <-chan struct{} {
	return nil
}

// SetTimeout defines timeout for all tasks
func (wg *AdvancedWaitGroup) SetTimeout(t time.Duration) *AdvancedWaitGroup {
	wg.timeout = &t
	return wg
}

// Do not run task after one that returns an error.
// You can pass context cancel() as second parameter to stop running tasks.
func (wg *AdvancedWaitGroup) SetStopOnError(b bool, onError func()) *AdvancedWaitGroup {
	wg.stopOnError = b
	wg.onError = onError
	return wg
}

// Add adds new task in waitgroup
func (wg *AdvancedWaitGroup) Add(f ...WaitgroupFunc) *AdvancedWaitGroup {
	wg.stackBuffer = append(wg.stackBuffer, f...)
	return wg
}

// AddSlice adds new tasks in waitgroup
func (wg *AdvancedWaitGroup) AddSlice(s []WaitgroupFunc) *AdvancedWaitGroup {
	return wg.Add(s...)
}

// WithContext make WaitGroup work with context timeout and Done
func (wg *AdvancedWaitGroup) WithContext(ctx context.Context) *AdvancedWaitGroup {
	wg.ctx = ctx
	wg.done = ctx.Done
	return wg
}

// SetCapacity defines tasks channel capacity
func (wg *AdvancedWaitGroup) SetCapacity(c int) *AdvancedWaitGroup {
	if c >= 0 {
		wg.capacity = uint32(c)
	}
	return wg
}

// GetCapacity defines tasks channel capacity
func (wg *AdvancedWaitGroup) GetCapacity() int {
	return int(wg.capacity)
}

func (wg *AdvancedWaitGroup) init() {
	wg.setStatus(StatusSuccess)
	wg.finished = make(chan struct{})
	if wg.done == nil {
		wg.done = done
	}

	wg.length = len(wg.stackBuffer)
	cap := wg.length
	if c := wg.GetCapacity(); c > 0 {
		cap = c
	}

	wg.receiver = make(chan struct{}, cap)
	wg.sender = make(chan WaitgroupFunc, wg.length)
	for _, f := range wg.stackBuffer {
		wg.sender <- f
	}
}

func (wg *AdvancedWaitGroup) Done() <-chan struct{} {
	return wg.finished
}

// Start runs tasks in separate goroutines
func (wg *AdvancedWaitGroup) Start() *AdvancedWaitGroup {
	if wg.CheckStatus(StatusSuccess) {
		return wg
	}

	wg.init()

	if wg.length > 0 {
		failed := make(chan error, wg.length)
		done := make(chan struct{}, wg.length)
		wgDone := make(chan struct{})

		var startTime time.Time
		var timer <-chan time.Time

		if wg.timeout != nil {
			if *wg.timeout != 0 {
				startTime = time.Now()
			}
			timer = time.After(*wg.timeout)
		}
		if wg.ctx != nil {
			startTime = time.Now()
		}

		w := sync.WaitGroup{}
	ForLoop:
		for wg.length > 0 {
			select {
			case f := <-wg.sender:
				wg.receiver <- struct{}{}
				w.Add(1)
				go func(f WaitgroupFunc, failed chan<- error, done chan<- struct{}) {
					defer func() {
						<-wg.receiver
						w.Done()
					}()
					if wg.stopOnError {
						wg.doIfSuccess(f, failed, done)
						return
					}

					wg.do(f, failed, done)

				}(f, failed, done)
			case err := <-failed:
				if wg.onError != nil {
					wg.onError()
				}
				wg.errors = append(wg.errors, err)
				wg.length--
				if wg.stopOnError {
					wg.setStatus(StatusError)
					break ForLoop
				}
			case <-done:
				wg.length--
			case <-wg.done():
				if deadlineTime, ok := wg.ctx.Deadline(); ok {
					wg.errors = append(wg.errors, ErrorTimeout(deadlineTime.Sub(startTime)))
					wg.setStatus(StatusTimeout)
					break ForLoop
				}
				if wg.ctx.Err() != nil {
					wg.errors = append(wg.errors, wg.ctx.Err())
					wg.setStatus(StatusError)
				}
				break ForLoop
			case t := <-timer:
				d := t.Sub(startTime)
				wg.errors = append(wg.errors, ErrorTimeout(d))
				wg.setStatus(StatusTimeout)
				break ForLoop
			}
		}

		close(wgDone)
		close(wg.sender)
		w.Wait()
	}

	close(wg.finished)

	return wg
}

func (wg *AdvancedWaitGroup) do(f WaitgroupFunc, failed chan<- error, done chan<- struct{}) {
	// Handle panic and pack it into stdlib error
	defer func() {
		if r := recover(); r != nil {
			buf := make([]byte, stackBufferSize)
			count := runtime.Stack(buf, false)
			failed <- errors.NewRecoveredPanicError(r, buf[:count])
		}
	}()

	if err := f(); err != nil {
		failed <- err
		return
	}

	done <- struct{}{}
}

func (wg *AdvancedWaitGroup) doIfSuccess(f WaitgroupFunc, failed chan<- error, done chan<- struct{}) {
	// Handle panic and pack it into stdlib error
	defer func() {
		if r := recover(); r != nil {
			buf := make([]byte, stackBufferSize)
			count := runtime.Stack(buf, false)
			failed <- fmt.Errorf("panic handeled\n%v\n%s", r, buf[:count])
		}
	}()

	// Check stop on error
	if !wg.CheckStatus(StatusSuccess) {
		// If some other goroutine get an error
		done <- struct{}{}
		return
	}

	if err := f(); err != nil {
		failed <- err
		return
	}

	done <- struct{}{}
}

// Reset performs cleanup task queue and reset state
func (wg *AdvancedWaitGroup) Reset() {
	wg.stackBuffer = []WaitgroupFunc{}
	close(wg.receiver)
	wg.sender = nil
	wg.timeout = nil
	wg.stopOnError = false
	wg.setStatus(StatusIdle)

	// pool
	wg.errors = []error{}
}

// GetLastError returns last error that caught by execution process
func (wg *AdvancedWaitGroup) GetLastError() error {
	if l := len(wg.errors); l > 0 {
		return wg.errors[l-1]
	}
	return nil
}

// GetAllErrors returns all errors that caught by execution process
func (wg *AdvancedWaitGroup) GetAllErrors() []error {
	return wg.errors
}

func (wg *AdvancedWaitGroup) setStatus(status int) {
	if status < StatusIdle || status > StatusError {
		return
	}

	wg.statusLock.Lock()
	wg.status = status
	wg.statusLock.Unlock()
}

// Status return result state string
func (wg *AdvancedWaitGroup) Status() int {
	wg.statusLock.RLock()
	defer wg.statusLock.RUnlock()

	return wg.status
}

// CheckStatus return result of status compare
func (wg *AdvancedWaitGroup) CheckStatus(status int) bool {
	if status < StatusIdle || status > StatusError {
		return false
	}

	wg.statusLock.RLock()
	defer wg.statusLock.RUnlock()

	return wg.status == status
}
