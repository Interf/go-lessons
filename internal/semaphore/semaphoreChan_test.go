package semaphore

import (
	"context"
	"go-lessons/internal/errorsList"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestSemaphoreChan(t *testing.T, max int64) *SemaphoreChan {
	t.Helper()

	sem, err := NewSemaphoreChan(max)
	require.NoError(t, err)

	return sem
}

func TestNewSemaphoreChan(t *testing.T) {
	tests := map[string]struct {
		max     int64
		wantErr error
	}{
		"less zero": {
			max:     -1,
			wantErr: errorsList.ErrParamLessZero,
		},
		"zero": {
			max:     0,
			wantErr: errorsList.ErrParamLessZero,
		},
		"correct": {
			max:     5,
			wantErr: nil,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			sem, err := NewSemaphoreChan(tc.max)

			if tc.wantErr != nil {
				require.ErrorIs(t, err, tc.wantErr)
				assert.Nil(t, sem)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, sem)
			assert.Equal(t, int(tc.max), cap(sem.ch))
			assert.Equal(t, 0, len(sem.ch))
		})
	}
}

func TestSemaphoreChan_Acquire(t *testing.T) {
	tests := map[string]struct {
		sem     *SemaphoreChan
		n       int64
		wantErr error
		wantLen int
	}{
		"less zero": {
			sem:     newTestSemaphoreChan(t, 5),
			n:       -1,
			wantErr: errorsList.ErrParamLessZero,
			wantLen: 0,
		},
		"zero": {
			sem:     newTestSemaphoreChan(t, 5),
			n:       0,
			wantErr: errorsList.ErrParamLessZero,
			wantLen: 0,
		},
		"more than max": {
			sem:     newTestSemaphoreChan(t, 5),
			n:       6,
			wantErr: errorsList.ErrParmMoreAvailable,
			wantLen: 0,
		},
		"correct acquire": {
			sem:     newTestSemaphoreChan(t, 5),
			n:       3,
			wantErr: nil,
			wantLen: 3,
		},
		"acquire all": {
			sem:     newTestSemaphoreChan(t, 5),
			n:       5,
			wantErr: nil,
			wantLen: 5,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			err := tc.sem.Acquire(context.Background(), tc.n)

			if tc.wantErr != nil {
				require.ErrorIs(t, err, tc.wantErr)
			} else {
				require.NoError(t, err)
			}

			assert.Equal(t, tc.wantLen, len(tc.sem.ch))
		})
	}
}

func TestSemaphoreChan_TryAcquire(t *testing.T) {
	tests := map[string]struct {
		sem     *SemaphoreChan
		n       int64
		want    bool
		wantLen int
	}{
		"less zero": {
			sem:     newTestSemaphoreChan(t, 5),
			n:       -1,
			want:    false,
			wantLen: 0,
		},
		"zero": {
			sem:     newTestSemaphoreChan(t, 5),
			n:       0,
			want:    false,
			wantLen: 0,
		},
		"more than max": {
			sem:     newTestSemaphoreChan(t, 5),
			n:       6,
			want:    false,
			wantLen: 0,
		},
		"correct acquire": {
			sem:     newTestSemaphoreChan(t, 5),
			n:       3,
			want:    true,
			wantLen: 3,
		},
		"acquire all": {
			sem:     newTestSemaphoreChan(t, 5),
			n:       5,
			want:    true,
			wantLen: 5,
		},
		"not enough available with rollback": {
			sem: func() *SemaphoreChan {
				sem := newTestSemaphoreChan(t, 5)

				for i := 0; i < 4; i++ {
					sem.ch <- struct{}{}
				}

				return sem
			}(),
			n:       2,
			want:    false,
			wantLen: 4,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := tc.sem.TryAcquire(tc.n)

			assert.Equal(t, tc.want, got)
			assert.Equal(t, tc.wantLen, len(tc.sem.ch))
		})
	}
}

func TestSemaphoreChan_Release(t *testing.T) {
	tests := map[string]struct {
		sem     *SemaphoreChan
		n       int64
		want    bool
		wantLen int
	}{
		"less zero": {
			sem: func() *SemaphoreChan {
				sem := newTestSemaphoreChan(t, 5)
				require.NoError(t, sem.Acquire(context.Background(), 3))
				return sem
			}(),
			n:       -1,
			want:    false,
			wantLen: 3,
		},
		"zero": {
			sem: func() *SemaphoreChan {
				sem := newTestSemaphoreChan(t, 5)
				require.NoError(t, sem.Acquire(context.Background(), 3))
				return sem
			}(),
			n:       0,
			want:    false,
			wantLen: 3,
		},
		"release more than acquired": {
			sem: func() *SemaphoreChan {
				sem := newTestSemaphoreChan(t, 5)
				require.NoError(t, sem.Acquire(context.Background(), 2))
				return sem
			}(),
			n:       3,
			want:    false,
			wantLen: 0,
		},
		"correct release": {
			sem: func() *SemaphoreChan {
				sem := newTestSemaphoreChan(t, 5)
				require.NoError(t, sem.Acquire(context.Background(), 4))
				return sem
			}(),
			n:       2,
			want:    true,
			wantLen: 2,
		},
		"release all": {
			sem: func() *SemaphoreChan {
				sem := newTestSemaphoreChan(t, 5)
				require.NoError(t, sem.Acquire(context.Background(), 5))
				return sem
			}(),
			n:       5,
			want:    true,
			wantLen: 0,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := tc.sem.Release(tc.n)

			assert.Equal(t, tc.want, got)
			assert.Equal(t, tc.wantLen, len(tc.sem.ch))
		})
	}
}

func TestSemaphoreChan_AcquireWithAlreadyCanceledContext(t *testing.T) {
	sem := newTestSemaphoreChan(t, 5)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := sem.Acquire(ctx, 1)

	require.Error(t, err)
	assert.ErrorIs(t, err, context.Canceled)
	assert.Equal(t, 0, len(sem.ch))
}

func TestSemaphoreChan_AcquireRollbackOnContextCancel(t *testing.T) {
	sem := newTestSemaphoreChan(t, 3)

	require.NoError(t, sem.Acquire(context.Background(), 1))
	assert.Equal(t, 1, len(sem.ch))

	ctx, cancel := context.WithCancel(context.Background())

	done := make(chan error, 1)

	go func() {
		done <- sem.Acquire(ctx, 3)
	}()

	require.Eventually(t, func() bool {
		return len(sem.ch) == 3
	}, 200*time.Millisecond, time.Millisecond)

	cancel()

	select {
	case err := <-done:
		require.Error(t, err)
		assert.ErrorIs(t, err, context.Canceled)
	case <-time.After(200 * time.Millisecond):
		t.Fatal("Acquire should finish after context cancel")
	}

	assert.Equal(t, 1, len(sem.ch))
}
