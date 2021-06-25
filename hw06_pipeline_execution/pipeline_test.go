package hw06pipelineexecution

import (
	"reflect"
	"strconv"
	"sync/atomic"
	"testing"
	"time"

	"github.com/benbjohnson/clock"
	"github.com/stretchr/testify/require"
	"go.uber.org/goleak"
)

const (
	sleepPerStage = time.Millisecond * 100
	fault         = sleepPerStage / 2
)

func contains(ar []uint32, num uint32) bool {
	for _, v := range ar {
		if v == num {
			return true
		}
	}

	return false
}

func TestPipelineEmpty(t *testing.T) {
	defer goleak.VerifyNone(t)
	mock := clock.NewMock()

	t.Run("empty slice case", func(t *testing.T) {
		in := make(Bi)
		var data []int

		// Stage generator
		g := func(_ string, f func(v interface{}) interface{}) Stage {
			return func(in In) Out {
				out := make(Bi)
				go func() {
					defer close(out)

					for v := range in {
						mock.Sleep(sleepPerStage)
						out <- f(v)
					}
				}()

				return out
			}
		}

		stages := []Stage{
			g("Dummy", func(v interface{}) interface{} { return v }),
			g("Multiplier (* 2)", func(v interface{}) interface{} { return v.(int) * 2 }),
			g("Adder (+ 100)", func(v interface{}) interface{} { return v.(int) + 100 }),
			g("Stringifier", func(v interface{}) interface{} {
				return strconv.Itoa(v.(int))
			}),
		}

		go func() {
			for _, v := range data {
				in <- v
			}

			close(in)
		}()

		result := make([]string, 0, 10)
		require.Eventually(t, func() bool {
			for s := range ExecutePipeline(in, nil, stages...) {
				result = append(result, s.(string))
			}

			return len(result) == 0
		}, fault, time.Millisecond)
	})
}

func TestPipelineSimple(t *testing.T) {
	// Without sleeps
	defer goleak.VerifyNone(t)
	mock := clock.NewMock()

	t.Run("simple case", func(t *testing.T) {
		in := make(Bi)
		data := []int{1, 2, 3, 4}
		rewindTimePoints := []uint32{1, 3, 6, 10, 13, 15, 16} // len should be equal len(stages) + len(data) - 1
		var stepToTheFuture uint32                            // step to check in backToTheFuture var
		syncCh := make(chan struct{}, 16)

		// Stage generator
		g := func(_ string, f func(v interface{}) interface{}) Stage {
			return func(in In) Out {
				out := make(Bi)
				go func() {
					defer close(out)

					for v := range in {
						syncCh <- struct{}{}
						mock.Sleep(sleepPerStage)
						out <- f(v)
					}
				}()

				return out
			}
		}

		stages := []Stage{
			g("Dummy", func(v interface{}) interface{} { return v }),
			g("Multiplier (* 2)", func(v interface{}) interface{} { return v.(int) * 2 }),
			g("Adder (+ 100)", func(v interface{}) interface{} { return v.(int) + 100 }),
			g("Stringifier", func(v interface{}) interface{} {
				return strconv.Itoa(v.(int))
			}),
		}

		go func() {
			for range syncCh {
				atomic.AddUint32(&stepToTheFuture, 1)
				if contains(rewindTimePoints, stepToTheFuture) {
					mock.Add(sleepPerStage)
				}
			}
		}()

		go func() {
			for _, v := range data {
				in <- v
			}

			close(in)
		}()

		result := make([]string, 0, 10)
		require.Eventually(t, func() bool {
			for s := range ExecutePipeline(in, nil, stages...) {
				result = append(result, s.(string))
			}
			close(syncCh)

			return reflect.DeepEqual(result, []string{"102", "104", "106", "108"})
		}, time.Duration(int64(sleepPerStage)*int64(len(stages)+len(data)-1)+int64(fault)), time.Millisecond)
	})
}

func TestPipelineDone(t *testing.T) {
	// Without sleeps
	defer goleak.VerifyNone(t)
	mock := clock.NewMock()

	t.Run("done case", func(t *testing.T) {
		in := make(Bi)
		done := make(Bi)
		data := []int{1, 2, 3, 4}
		syncCh := make(chan struct{}, 16)

		// Stage generator
		g := func(_ string, f func(v interface{}) interface{}) Stage {
			return func(in In) Out {
				out := make(Bi)
				go func() {
					defer close(out)

					for v := range in {
						syncCh <- struct{}{}
						mock.Sleep(sleepPerStage)
						out <- f(v)
					}
				}()

				return out
			}
		}

		stages := []Stage{
			g("Dummy", func(v interface{}) interface{} { return v }),
			g("Multiplier (* 2)", func(v interface{}) interface{} { return v.(int) * 2 }),
			g("Adder (+ 100)", func(v interface{}) interface{} { return v.(int) + 100 }),
			g("Stringifier", func(v interface{}) interface{} {
				return strconv.Itoa(v.(int))
			}),
		}

		// Abort after 200ms
		abortDur := sleepPerStage * 2
		go func() {
			<-mock.After(abortDur)
			close(done)
		}()

		rewindTimePoints := []uint32{1, 3, 6, 10, 13, 15, 16} // len should be equal len(stages) + len(data) - 1
		var stepToTheFuture uint32                            // step to check in backToTheFuture var

		go func() {
			for range syncCh {
				atomic.AddUint32(&stepToTheFuture, 1)
				if contains(rewindTimePoints, stepToTheFuture) {
					mock.Add(sleepPerStage)
				}
			}
		}()

		go func() {
			for _, v := range data {
				in <- v
			}
			close(in)
		}()

		result := make([]string, 0, 10)

		require.Eventually(t, func() bool {
			for s := range ExecutePipeline(in, done, stages...) {
				result = append(result, s.(string))
			}
			defer close(syncCh)

			return len(result) == 0
		}, time.Duration(int64(abortDur)+int64(fault)), time.Millisecond)
	})
}
