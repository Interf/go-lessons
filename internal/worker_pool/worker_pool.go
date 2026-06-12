package worker_pool

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

type Pool interface {
	Run(ctx context.Context, tasks []Task) error
}

type WorkerPool struct {
	workersCount int
	errorsLimit  int
}

func NewPool(workersCount int, errorsLimit int) *WorkerPool {
	if workersCount <= 0 {
		workersCount = 1
	}
	if errorsLimit < 0 {
		errorsLimit = 0
	}
	return &WorkerPool{
		workersCount: workersCount,
		errorsLimit:  errorsLimit,
	}
}

func (p *WorkerPool) Run(ctx context.Context, tasks []Task) error {
	taskCh := make(chan Task, len(tasks))
	done := make(chan struct{})
	var once sync.Once

	closeDone := func() {
		once.Do(func() { close(done) })
	}

	var wg sync.WaitGroup

	var errCount int32
	var cancelled int32

	for i := 0; i < p.workersCount; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			for {
				select {
				case task, ok := <-taskCh:
					if !ok {
						return
					}

					if atomic.LoadInt32(&cancelled) == 1 {
						continue
					}

					err := safeCall(task)

					if err != nil {
						newCount := atomic.AddInt32(&errCount, 1)
						if int(newCount) >= p.errorsLimit {
							atomic.StoreInt32(&cancelled, 1)
							closeDone()
							return
						}
					}
				case <-done:
					return
				}
			}
		}()
	}

	go func() {
		<-ctx.Done()
		atomic.StoreInt32(&cancelled, 1)
		closeDone()
	}()

	for _, task := range tasks {
		if atomic.LoadInt32(&cancelled) == 1 {
			break
		}

		select {
		case taskCh <- task:
		case <-done:
			break
		}
	}

	close(taskCh)
	wg.Wait()

	if int(atomic.LoadInt32(&errCount)) >= p.errorsLimit {
		return ErrErrorsLimitExceeded
	}

	if ctx.Err() != nil {
		return ctx.Err()
	}

	return nil
}

func safeCall(task Task) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic recovered: %v", r)
		}
	}()
	return task()
}
