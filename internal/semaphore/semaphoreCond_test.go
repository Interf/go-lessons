package semaphore

import (
	"context"
	"go-lessons/internal/errorsList"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewSemaphoreCond(t *testing.T) {

}

func TestSemaphoreCond_Acquire(t *testing.T) {
	test := map[string]struct {
		sem       *SemaphoreCond
		n         int64
		wantErr   error
		wantCount int64
	}{
		"less zero": {
			sem:       NewSemaphoreCond(5),
			n:         -1,
			wantErr:   errorsList.ErrParamLessZero,
			wantCount: 0,
		},
		"zero": {
			sem:       NewSemaphoreCond(5),
			n:         0,
			wantErr:   errorsList.ErrParamLessZero,
			wantCount: 0,
		},
		"more then max": {
			sem:       NewSemaphoreCond(5),
			n:         6,
			wantErr:   errorsList.ErrParmMoreAvailable,
			wantCount: 0,
		},
		"correct": {
			sem:       NewSemaphoreCond(5),
			n:         3,
			wantErr:   nil,
			wantCount: 3,
		},
	}

	for name, tc := range test {
		t.Run(name, func(t *testing.T) {
			err := tc.sem.Acquire(context.Background(), tc.n)

			if tc.wantErr != nil {
				assert.ErrorIs(t, err, tc.wantErr)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tc.wantCount, tc.sem.count)
		})
	}

}

func TestSemaphoreCond_AcquireWaitsUntilRelease(t *testing.T) {
	sem := NewSemaphoreCond(2)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	err := sem.Acquire(ctx, 2)
	require.NoError(t, err)

	done := make(chan error, 1)

	go func() {
		done <- sem.Acquire(ctx, 1)
	}()

	select {
	case err := <-done:
		t.Fatalf("Acquire should block, but returned: %v", err)
	case <-time.After(20 * time.Millisecond):
		// ok: Acquire is blocked
	}

	ok := sem.Release(1)
	require.True(t, ok)

	select {
	case err := <-done:
		require.NoError(t, err)
	case <-time.After(200 * time.Millisecond):
		t.Fatal("Acquire should finish after Release")
	}

	assert.Equal(t, int64(2), sem.count)
}

func TestSemaphoreCond_AcquireReturnsContextErrorWhileWaiting(t *testing.T) {
	sem := NewSemaphoreCond(1)

	err := sem.Acquire(context.Background(), 1)
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer cancel()

	err = sem.Acquire(ctx, 1)

	require.Error(t, err)
	assert.ErrorIs(t, err, context.DeadlineExceeded)
	assert.Equal(t, int64(1), sem.count)
}

func TestSemaphoreCond_AcquireReturnsErrorWhenContextAlreadyCanceled(t *testing.T) {
	sem := NewSemaphoreCond(5)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := sem.Acquire(ctx, 1)

	require.Error(t, err)
	assert.ErrorIs(t, err, context.Canceled)
	assert.Equal(t, int64(0), sem.count)
}

func TestSemaphoreCond_TryAcquire(t *testing.T) {
	tests := map[string]struct {
		sem       *SemaphoreCond
		n         int64
		want      bool
		wantCount int64
	}{
		"less zero": {
			sem:       NewSemaphoreCond(5),
			n:         -1,
			want:      false,
			wantCount: 0,
		},
		"zero": {
			sem:       NewSemaphoreCond(5),
			n:         0,
			want:      false,
			wantCount: 0,
		},
		"more than max": {
			sem:       NewSemaphoreCond(5),
			n:         6,
			want:      false,
			wantCount: 0,
		},
		"correct acquire": {
			sem:       NewSemaphoreCond(5),
			n:         3,
			want:      true,
			wantCount: 3,
		},
		"acquire all": {
			sem:       NewSemaphoreCond(5),
			n:         5,
			want:      true,
			wantCount: 5,
		},
		"not enough available": {
			sem: func() *SemaphoreCond {
				sem := NewSemaphoreCond(5)
				sem.count = 4
				return sem
			}(),
			n:         2,
			want:      false,
			wantCount: 4,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := tc.sem.TryAcquire(tc.n)

			assert.Equal(t, tc.want, got)
			assert.Equal(t, tc.wantCount, tc.sem.count)
		})
	}
}

func TestSemaphoreCond_Release(t *testing.T) {
	tests := map[string]struct {
		sem       *SemaphoreCond
		n         int64
		want      bool
		wantCount int64
	}{
		"less zero": {
			sem:       NewSemaphoreCond(5),
			n:         -1,
			want:      false,
			wantCount: 0,
		},
		"zero": {
			sem:       NewSemaphoreCond(5),
			n:         0,
			want:      false,
			wantCount: 0,
		},
		"release more than acquired": {
			sem: func() *SemaphoreCond {
				sem := NewSemaphoreCond(5)
				sem.count = 2
				return sem
			}(),
			n:         3,
			want:      false,
			wantCount: 2,
		},
		"correct release": {
			sem: func() *SemaphoreCond {
				sem := NewSemaphoreCond(5)
				sem.count = 4
				return sem
			}(),
			n:         2,
			want:      true,
			wantCount: 2,
		},
		"release all": {
			sem: func() *SemaphoreCond {
				sem := NewSemaphoreCond(5)
				sem.count = 5
				return sem
			}(),
			n:         5,
			want:      true,
			wantCount: 0,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := tc.sem.Release(tc.n)

			assert.Equal(t, tc.want, got)
			assert.Equal(t, tc.wantCount, tc.sem.count)
		})
	}
}
