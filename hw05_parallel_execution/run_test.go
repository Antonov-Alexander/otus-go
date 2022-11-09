package hw05parallelexecution

import (
	"errors"
	"fmt"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/goleak"
)

func TestRun(t *testing.T) {
	defer goleak.VerifyNone(t)

	t.Run("if were errors in first M tasks, than finished not more N+M tasks", func(t *testing.T) {
		tasksCount := 50

		stopCh := make(chan struct{}, tasksCount)
		defer close(stopCh)

		sleepTime := time.Millisecond * 10
		executionTime := sleepTime * time.Duration(tasksCount+1)

		var stopEventsCount int
		require.Eventually(t, func() bool {
			stopCh <- struct{}{}
			stopEventsCount++
			return stopEventsCount == tasksCount
		}, executionTime, sleepTime)

		var runTasksCount int32
		tasks := make([]Task, 0, tasksCount)

		for i := 0; i < tasksCount; i++ {
			err := fmt.Errorf("error from task %d", i)
			tasks = append(tasks, func() error {
				atomic.AddInt32(&runTasksCount, 1)
				for range stopCh {
					break
				}
				return err
			})
		}

		workersCount := 10
		maxErrorsCount := 23
		err := Run(tasks, workersCount, maxErrorsCount)

		require.Truef(t, errors.Is(err, ErrErrorsLimitExceeded), "actual err - %v", err)
		require.LessOrEqual(t, runTasksCount, int32(workersCount+maxErrorsCount), "extra tasks were started")
	})

	t.Run("tasks without errors", func(t *testing.T) {
		tasksCount := 50

		stopCh := make(chan struct{}, tasksCount)
		defer close(stopCh)

		sleepTime := time.Millisecond * 10
		executionTime := sleepTime * time.Duration(tasksCount+1)

		var stopEventsCount int
		require.Eventually(t, func() bool {
			stopCh <- struct{}{}
			stopEventsCount++
			return stopEventsCount == tasksCount
		}, executionTime, sleepTime)

		var runTasksCount int32
		tasks := make([]Task, 0, tasksCount)

		for i := 0; i < tasksCount; i++ {
			tasks = append(tasks, func() error {
				atomic.AddInt32(&runTasksCount, 1)
				for range stopCh {
					break
				}
				return nil
			})
		}

		workersCount := 5
		maxErrorsCount := 1

		err := Run(tasks, workersCount, maxErrorsCount)

		require.NoError(t, err)
		require.Equal(t, runTasksCount, int32(tasksCount), "not all tasks were completed")
	})

	t.Run("negative errors limit test", func(t *testing.T) {
		tests := []struct {
			errorsLimit   int
			expectedError error
		}{
			{errorsLimit: 0, expectedError: ErrErrorsLimitExceeded},
			{errorsLimit: -1, expectedError: ErrErrorsLimitExceeded},
		}

		tasks := make([]Task, 0)
		for _, test := range tests {
			err := Run(tasks, 10, test.errorsLimit)
			require.Truef(t, errors.Is(err, test.expectedError), "actual err - %v", err)
		}
	})

	t.Run("single worker deadlock test", func(t *testing.T) {
		tasks := []Task{
			func() error {
				return nil
			},
			func() error {
				return ErrErrorsLimitExceeded
			},
			func() error {
				return nil
			},
		}

		_ = Run(tasks, 1, 1)
	})
}
