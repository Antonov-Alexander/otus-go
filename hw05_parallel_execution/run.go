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

	var errCount atomic.Int32 //nolint:typecheck
	errLimit := int32(m)

	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for task := range ch {
				if errCount.Load() >= errLimit { //nolint:typecheck
					break
				}

				if err := task(); err != nil {
					errCount.Add(1)
				}
			}
		}()
	}

	for _, task := range tasks {
		if errCount.Load() >= errLimit { //nolint:typecheck
			break
		}

		ch <- task
	}

	close(ch)
	wg.Wait()

	if errCount.Load() >= errLimit { //nolint:typecheck
		return ErrErrorsLimitExceeded
	}

	return nil
}
