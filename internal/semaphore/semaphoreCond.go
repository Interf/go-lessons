package semaphore

import (
	"context"
	"go-lessons/internal/errorsList"
	"sync"
)

type SemaphoreCond struct {
	count int64
	max   int64
	mu    sync.Mutex
	cond  *sync.Cond
}

func NewSemaphoreCond(max int64) *SemaphoreCond {
	s := &SemaphoreCond{
		max: max,
	}

	s.cond = sync.NewCond(&s.mu)

	return s
}

func (c *SemaphoreCond) Acquire(ctx context.Context, n int64) error {

	if n <= 0 {
		return errorsList.ErrParamLessZero
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	if n > c.max {
		return errorsList.ErrParmMoreAvailable
	}

	if err := ctx.Err(); err != nil {
		return err
	}

	if c.count+n <= c.max {
		c.count += n
		return nil
	}

	done := make(chan struct{})
	defer close(done)

	go func() {
		select {
		case <-ctx.Done():
			c.mu.Lock()
			c.cond.Broadcast()
			c.mu.Unlock()
		case <-done:
		}
	}()

	for c.count+n > c.max {
		c.cond.Wait()

		if err := ctx.Err(); err != nil {
			return err
		}
	}

	c.count += n

	return nil
}

func (c *SemaphoreCond) TryAcquire(n int64) bool {
	if n <= 0 {
		return false
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	if n > c.max {
		return false
	}

	if c.count+n <= c.max {
		c.count += n
		return true
	}

	return false
}

func (c *SemaphoreCond) Release(n int64) bool {
	if n <= 0 {
		return false
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	if n > c.count {
		return false
	}

	c.count -= n
	c.cond.Broadcast()

	return true
}
