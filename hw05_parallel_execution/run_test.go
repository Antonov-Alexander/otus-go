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
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32
		taskSleep := time.Millisecond * 100

		for i := 0; i < tasksCount; i++ {
			err := fmt.Errorf("error from task %d", i)
			tasks = append(tasks, func() error {
				atomic.AddInt32(&runTasksCount, 1)
				require.Eventually(t, func() bool {
					return true
				}, taskSleep*2, taskSleep)

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
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32
		var sumTime time.Duration

		for i := 0; i < tasksCount; i++ {
			taskSleep := time.Millisecond * 100
			sumTime += taskSleep

			tasks = append(tasks, func() error {
				atomic.AddInt32(&runTasksCount, 1)

				require.Eventually(t, func() bool {
					return true
				}, taskSleep*2, taskSleep)

				return nil
			})
		}

		workersCount := 5
		maxErrorsCount := 1

		start := time.Now()
		err := Run(tasks, workersCount, maxErrorsCount)
		elapsedTime := time.Since(start)
		require.NoError(t, err)

		require.Equal(t, runTasksCount, int32(tasksCount), "not all tasks were completed")
		require.LessOrEqual(t, int64(elapsedTime), int64(sumTime/2), "tasks were run sequentially?")
	})

	t.Run("negative errors limit test", func(t *testing.T) {
		tasks := make([]Task, 0)
		for _, i := range []int{0, -1} {
			err := Run(tasks, 10, i)
			require.Truef(t, errors.Is(err, ErrErrorsLimitExceeded), "actual err - %v", err)
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

		// интересно, почему здесь не работает присвоение ":=", ведь ранее "_" не объявлялась
		_ = Run(tasks, 1, 1)
	})
}
