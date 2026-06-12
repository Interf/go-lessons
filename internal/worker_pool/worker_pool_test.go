package worker_pool

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewPool(t *testing.T) {
	t.Run("valid parameters", func(t *testing.T) {
		pool := NewPool(5, 3)
		assert.Equal(t, 5, pool.workersCount)
		assert.Equal(t, 3, pool.errorsLimit)
	})

	t.Run("zero workers becomes one", func(t *testing.T) {
		pool := NewPool(0, 3)
		assert.Equal(t, 1, pool.workersCount)
	})

	t.Run("negative workers becomes one", func(t *testing.T) {
		pool := NewPool(-5, 3)
		assert.Equal(t, 1, pool.workersCount)
	})

	t.Run("negative errors limit becomes zero", func(t *testing.T) {
		pool := NewPool(1, -10)
		assert.Equal(t, 0, pool.errorsLimit)
	})
}

func TestRunAllTasksComplete(t *testing.T) {
	pool := NewPool(3, 10)
	ctx := context.Background()

	var counter int32
	tasks := make([]Task, 10)

	for i := range tasks {
		tasks[i] = func() error {
			atomic.AddInt32(&counter, 1)
			return nil
		}
	}

	err := pool.Run(ctx, tasks)

	assert.NoError(t, err)
	assert.Equal(t, int32(10), atomic.LoadInt32(&counter))
}

func TestRunErrorLimitExceeded(t *testing.T) {
	pool := NewPool(3, 2)
	ctx := context.Background()

	tasks := make([]Task, 10)
	for i := range tasks {
		tasks[i] = func() error {
			return assert.AnError
		}
	}

	err := pool.Run(ctx, tasks)

	assert.ErrorIs(t, err, ErrErrorsLimitExceeded)
}

func TestRunNoErrors(t *testing.T) {
	pool := NewPool(3, 5)
	ctx := context.Background()

	tasks := make([]Task, 10)
	for i := range tasks {
		tasks[i] = func() error {
			return nil
		}
	}

	err := pool.Run(ctx, tasks)

	assert.NoError(t, err)
}

func TestRunPanicRecovery(t *testing.T) {
	pool := NewPool(3, 10)
	ctx := context.Background()

	var completed int32
	tasks := make([]Task, 5)

	tasks[0] = func() error {
		panic("unexpected error")
	}
	for i := 1; i < 5; i++ {
		tasks[i] = func() error {
			atomic.AddInt32(&completed, 1)
			return nil
		}
	}

	err := pool.Run(ctx, tasks)

	assert.NoError(t, err)
	assert.Equal(t, int32(4), atomic.LoadInt32(&completed))
}

func TestRunWithContextTimeout(t *testing.T) {
	pool := NewPool(1, 100)
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	tasks := make([]Task, 5)
	for i := range tasks {
		tasks[i] = func() error {
			time.Sleep(200 * time.Millisecond)
			return nil
		}
	}

	start := time.Now()
	err := pool.Run(ctx, tasks)
	elapsed := time.Since(start)

	assert.ErrorIs(t, err, context.DeadlineExceeded)
	assert.Less(t, elapsed, 1*time.Second)
}

func TestRunWithContextCancel(t *testing.T) {
	pool := NewPool(1, 100)
	ctx, cancel := context.WithCancel(context.Background())

	var completed int32
	tasks := make([]Task, 10)
	for i := range tasks {
		tasks[i] = func() error {
			time.Sleep(50 * time.Millisecond)
			atomic.AddInt32(&completed, 1)
			return nil
		}
	}

	go func() {
		time.Sleep(25 * time.Millisecond)
		cancel()
	}()

	pool.Run(ctx, tasks)

	assert.Less(t, atomic.LoadInt32(&completed), int32(10))
}

func TestRunConcurrentExecution(t *testing.T) {
	workersCount := 5
	pool := NewPool(workersCount, 100)
	ctx := context.Background()

	var maxConcurrent int32
	var currentConcurrent int32
	var mu sync.Mutex

	tasks := make([]Task, 20)
	for i := range tasks {
		tasks[i] = func() error {
			mu.Lock()
			currentConcurrent++
			if currentConcurrent > maxConcurrent {
				maxConcurrent = currentConcurrent
			}
			mu.Unlock()

			time.Sleep(10 * time.Millisecond)

			mu.Lock()
			currentConcurrent--
			mu.Unlock()
			return nil
		}
	}

	err := pool.Run(ctx, tasks)

	assert.NoError(t, err)
	assert.GreaterOrEqual(t, int(maxConcurrent), 1)
	assert.LessOrEqual(t, int(maxConcurrent), workersCount)
}

func TestRunZeroTasks(t *testing.T) {
	pool := NewPool(3, 5)
	ctx := context.Background()

	err := pool.Run(ctx, []Task{})

	assert.NoError(t, err)
}

func TestRunSingleWorker(t *testing.T) {
	pool := NewPool(1, 10)
	ctx := context.Background()

	var completed int32
	tasks := make([]Task, 5)
	for i := range tasks {
		tasks[i] = func() error {
			atomic.AddInt32(&completed, 1)
			return nil
		}
	}

	err := pool.Run(ctx, tasks)

	assert.NoError(t, err)
	assert.Equal(t, int32(5), atomic.LoadInt32(&completed))
}

func TestRunZeroErrorLimit(t *testing.T) {
	pool := NewPool(3, 0)
	ctx := context.Background()

	tasks := make([]Task, 5)
	for i := range tasks {
		tasks[i] = func() error {
			return assert.AnError
		}
	}

	err := pool.Run(ctx, tasks)

	assert.ErrorIs(t, err, ErrErrorsLimitExceeded)
}

func TestRunCancelledDuringTaskSending(t *testing.T) {
	pool := NewPool(1, 1)
	ctx := context.Background()

	var started int32
	tasks := make([]Task, 5)

	tasks[0] = func() error {
		atomic.AddInt32(&started, 1)
		return assert.AnError
	}
	for i := 1; i < 5; i++ {
		tasks[i] = func() error {
			return nil
		}
	}

	err := pool.Run(ctx, tasks)

	assert.ErrorIs(t, err, ErrErrorsLimitExceeded)
}

func TestRunDoneChannelDuringSending(t *testing.T) {
	pool := NewPool(1, 1)
	ctx := context.Background()

	tasks := make([]Task, 5)

	tasks[0] = func() error {
		return assert.AnError
	}
	for i := 1; i < 5; i++ {
		tasks[i] = func() error {
			return nil
		}
	}

	err := pool.Run(ctx, tasks)

	assert.ErrorIs(t, err, ErrErrorsLimitExceeded)
}

func TestRunPartialErrors(t *testing.T) {
	pool := NewPool(3, 10)
	ctx := context.Background()

	var successCount int32
	tasks := make([]Task, 10)

	for i := range tasks {
		if i%2 == 0 {
			tasks[i] = func() error {
				return assert.AnError
			}
		} else {
			tasks[i] = func() error {
				atomic.AddInt32(&successCount, 1)
				return nil
			}
		}
	}

	err := pool.Run(ctx, tasks)

	assert.NoError(t, err)
	assert.Equal(t, int32(5), atomic.LoadInt32(&successCount))
}
