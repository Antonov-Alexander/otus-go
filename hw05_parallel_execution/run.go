package hw05parallelexecution

import (
	"errors"
	"sync"
	"sync/atomic"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

func Run(tasks []Task, n, m int) error {
	if m <= 0 {
		return ErrErrorsLimitExceeded
	}

	wg := sync.WaitGroup{}
	ch := make(chan Task)

	var errorsCount int32
	errorsLimit := int32(m)

	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for task := range ch {
				if atomic.LoadInt32(&errorsCount) >= errorsLimit {
					continue
				}

				if err := task(); err != nil {
					atomic.AddInt32(&errorsCount, 1)
				}
			}
		}()
	}

	for _, task := range tasks {
		if atomic.LoadInt32(&errorsCount) < errorsLimit {
			ch <- task
		}
	}

	close(ch)
	wg.Wait()

	if atomic.LoadInt32(&errorsCount) >= errorsLimit {
		return ErrErrorsLimitExceeded
	}

	return nil
}
