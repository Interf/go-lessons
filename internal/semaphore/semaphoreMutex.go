package semaphore

import (
	"context"
	"errors"
	"sync"
	"time"
)

type Semaphore interface {
	Acquire(context.Context, int64) error
	TryAcquire(int64) bool
	Release(int64)
}

type SemaphoreMutex struct {
	count int64
	max   int64
	mu    sync.Mutex
}

var ErrParamLessZero = errors.New("param less zero")
var ErrReleaseParmMoreAvailable = errors.New("parm more then available")

func (s *SemaphoreMutex) Acquire(ctx context.Context, n int64) error {
	if n <= 0 {
		return ErrParamLessZero
	}

	for {
		s.mu.Lock()

		if s.count+n <= s.max {
			s.count += n
			s.mu.Unlock()
			return nil
		}

		s.mu.Unlock()

		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			time.Sleep(time.Millisecond)
		}

	}

}

func (s *SemaphoreMutex) TryAcquire(n int64) bool {
	if n <= 0 {
		return false
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if n > s.max {
		return false
	}

	if s.count+n <= s.max {
		s.count += n
		return true
	}

	return false
}

func (s *SemaphoreMutex) Release(n int64) {
	if n <= 0 {
		panic(ErrParamLessZero)
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if n > s.count {
		panic(ErrReleaseParmMoreAvailable)
	}

	s.count -= n
}
