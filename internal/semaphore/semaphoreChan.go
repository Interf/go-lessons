package semaphore

import (
	"context"
	"go-lessons/internal/errorsList"
)

type SemaphoreChan struct {
	ch chan struct{}
}

func NewSemaphoreChan(max int64) (*SemaphoreChan, error) {
	if max <= 0 {
		return nil, errorsList.ErrParamLessZero
	}

	return &SemaphoreChan{
		ch: make(chan struct{}, max),
	}, nil

}

func (c *SemaphoreChan) Acquire(ctx context.Context, n int64) error {

	if n <= 0 {
		return errorsList.ErrParamLessZero
	}

	if n > int64(cap(c.ch)) {
		return errorsList.ErrParmMoreAvailable
	}

	acquired := int64(0)

	for acquired < n {
		select {
		case c.ch <- struct{}{}:
			acquired++

		case <-ctx.Done():
			for acquired > 0 {
				<-c.ch
				acquired--
			}

			return ctx.Err()
		}
	}

	return nil
}

func (c *SemaphoreChan) TryAcquire(n int64) bool {
	if n <= 0 {
		return false
	}

	if n > int64(cap(c.ch)) {
		return false
	}

	acquired := int64(0)

	for acquired < n {
		select {
		case c.ch <- struct{}{}:
			acquired++

		default:
			for acquired > 0 {
				<-c.ch
				acquired--
			}

			return false
		}
	}

	return true
}

func (c *SemaphoreChan) Release(n int64) bool {
	if n <= 0 {
		return false
	}

	for i := int64(0); i < n; i++ {
		select {
		case <-c.ch:
		default:
			return false
		}
	}

	return true
}
