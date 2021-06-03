package hw05parallelexecution

import (
	"errors"
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/benbjohnson/clock"
	"github.com/stretchr/testify/require"
	"go.uber.org/goleak"
)

func TestRun(t *testing.T) {
	defer goleak.VerifyNone(t)

	t.Run("No tasks", func(t *testing.T) {
		tasks := make([]Task, 0)
		err := Run(tasks, 10, 1)

		require.Nil(t, err)
	})

	t.Run("Invalid (zero) goroutines count", func(t *testing.T) {
		tasks := make([]Task, 3)
		err := Run(tasks, 0, 2)

		require.Truef(t, errors.Is(err, ErrInvalidGoroutinesCount), "actual err - %v", err)
	})

	t.Run("Invalid (negative) goroutines count", func(t *testing.T) {
		tasks := make([]Task, 3)
		err := Run(tasks, -1, 2)

		require.Truef(t, errors.Is(err, ErrInvalidGoroutinesCount), "actual err - %v", err)
	})

	t.Run("Invalid (zero) errors count", func(t *testing.T) {
		tasks := make([]Task, 3)
		err := Run(tasks, 2, 0)

		require.Truef(t, errors.Is(err, ErrErrorsLimitExceeded), "actual err - %v", err)
	})

	t.Run("Invalid (negative) errors count", func(t *testing.T) {
		tasks := make([]Task, 3)
		err := Run(tasks, 2, -1)

		require.Truef(t, errors.Is(err, ErrErrorsLimitExceeded), "actual err - %v", err)
	})

	t.Run("if were errors in first M tasks, than finished not more N+M tasks, "+
		"workers cnt less than max errs cnt", func(t *testing.T) {
		tasksCount := 50
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32

		for i := 0; i < tasksCount; i++ {
			err := fmt.Errorf("error from task %d", i)

			tasks = append(tasks, func() error {
				time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
				atomic.AddInt32(&runTasksCount, 1)

				return err
			})
		}

		workersCount := 10
		maxErrorsCount := 23
		err := Run(tasks, workersCount, maxErrorsCount)

		require.Truef(t, errors.Is(err, ErrErrorsLimitExceeded), "actual err - %v", err)
		require.LessOrEqual(t, runTasksCount, int32(workersCount+maxErrorsCount), "extra tasks were started")
	})

	t.Run("if were errors in first M tasks, than finished not more N+M tasks, "+
		"workers cnt greater than max errs cnt", func(t *testing.T) {
		tasksCount := 50
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32

		for i := 0; i < tasksCount; i++ {
			err := fmt.Errorf("error from task %d", i)

			tasks = append(tasks, func() error {
				time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
				atomic.AddInt32(&runTasksCount, 1)

				return err
			})
		}

		workersCount := 10
		maxErrorsCount := 4
		err := Run(tasks, workersCount, maxErrorsCount)

		require.Truef(t, errors.Is(err, ErrErrorsLimitExceeded), "actual err - %v", err)
		require.LessOrEqual(t, runTasksCount, int32(workersCount+maxErrorsCount), "extra tasks were started")
	})

	t.Run("tasks without errors and sleeps", func(t *testing.T) {
		tasksCount := 8
		workersCount := 4
		maxErrorsCount := 1
		tasks := make([]Task, 0, tasksCount)
		tmpCh := make(chan struct{}, workersCount)

		var runTasksCount int32
		mock := clock.NewMock()

		wg := sync.WaitGroup{}
		defer wg.Wait()
		defer close(tmpCh)

		wg.Add(1)
		go func() {
			var internalCounter int32

			for range tmpCh {
				internalCounter++

				if internalCounter == int32(workersCount) {
					internalCounter = 0

					mock.Add(time.Duration(workersCount/2) * time.Millisecond)
				}
			}

			wg.Done()
		}()

		for i := 0; i < tasksCount; i++ {
			tasks = append(tasks, func() error {
				tmpCh <- struct{}{}
				mock.Sleep(time.Millisecond)
				atomic.AddInt32(&runTasksCount, 1)

				return nil
			})
		}

		require.Eventually(t, func() bool {
			err := Run(tasks, workersCount, maxErrorsCount)

			return err == nil && int32(tasksCount) == runTasksCount
		}, time.Second, time.Millisecond, "tasks were run sequentially?")
	})
}
