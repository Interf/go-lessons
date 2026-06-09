package semaphore

import "context"

type Semaphore interface {
	Acquire(context.Context, int64) error
	TryAcquire(int64) bool
	Release(int64)
}
