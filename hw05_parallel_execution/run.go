package hw05parallelexecution

import (
	"errors"
	"sync"
)

var (
	ErrErrorsLimitExceeded    = errors.New("errors limit exceeded")
	ErrInvalidGoroutinesCount = errors.New("invalid goroutines count")
)

type Task func() error

type ErrorCounter struct {
	sync.Mutex
	val    int
	maxCnt int
}

func (ec *ErrorCounter) inc() {
	ec.Lock()
	defer ec.Unlock()

	ec.val++
}

func (ec *ErrorCounter) isLimitExceeded() bool {
	ec.Lock()
	defer ec.Unlock()

	return ec.val >= ec.maxCnt
}

func producer(ch chan<- Task, quitCh <-chan struct{}, tasks []Task) error {
	for _, t := range tasks {
		select {
		case <-quitCh:
			return ErrErrorsLimitExceeded
		default:
		}

		select {
		case ch <- t:
		case <-quitCh:
			return ErrErrorsLimitExceeded
		}
	}

	return nil
}

func consumer(ch <-chan Task, quitCh chan<- struct{}, ec *ErrorCounter) {
	for t := range ch {
		e := t()

		if e != nil {
			ec.inc()

			if ec.isLimitExceeded() {
				quitCh <- struct{}{}

				return
			}
		}
	}
}

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	if len(tasks) == 0 {
		return nil
	}

	if n <= 0 {
		return ErrInvalidGoroutinesCount
	}

	if m < 0 {
		m = 0
	}

	if m == 0 {
		return ErrErrorsLimitExceeded
	}

	ch := make(chan Task, n)
	quitCh := make(chan struct{}, 2*n-1)
	ec := ErrorCounter{maxCnt: m}
	wg := sync.WaitGroup{}

	defer close(quitCh)
	defer wg.Wait()
	defer close(ch)

	wg.Add(n)
	for i := 0; i < n; i++ {
		go func() {
			consumer(ch, quitCh, &ec)

			wg.Done()
		}()
	}

	e := producer(ch, quitCh, tasks)

	return e
}
