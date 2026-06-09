package semaphore

import (
	"context"
	"go-lessons/internal/errorsList"
	"sync"
	"time"
)

type SemaphoreMutex struct {
	count int64
	max   int64
	mu    sync.Mutex
}

func (s *SemaphoreMutex) Acquire(ctx context.Context, n int64) error {
	if n <= 0 {
		return errorsList.ErrParamLessZero
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
		panic(errorsList.ErrParamLessZero)
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if n > s.count {
		panic(errorsList.ErrParmMoreAvailable)
	}

	s.count -= n
}
