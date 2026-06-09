package semaphore

import (
	"context"
	"go-lessons/internal/errorsList"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSemaphoreMutex_Acquire(t *testing.T) {

	tests := map[string]struct {
		sem   *SemaphoreMutex
		n     int64
		isErr bool
	}{
		"less zero": {
			sem:   &SemaphoreMutex{max: 5},
			n:     -5,
			isErr: true,
		},
		"correct value": {
			sem:   &SemaphoreMutex{max: 5},
			n:     5,
			isErr: false,
		},
		"more then max": {
			sem:   &SemaphoreMutex{max: 5},
			n:     9,
			isErr: true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {

			ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
			defer cancel()

			err := tc.sem.Acquire(ctx, tc.n)

			if tc.isErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

		})
	}

}

func TestSemaphoreMutex_TryAcquire(t *testing.T) {

	tests := map[string]struct {
		sem  *SemaphoreMutex
		n    int64
		want bool
	}{
		"less zero": {
			sem:  &SemaphoreMutex{max: 5, count: 5},
			n:    -5,
			want: false,
		},
		"correct value": {
			sem:  &SemaphoreMutex{max: 5, count: 0},
			n:    5,
			want: true,
		},
		"more then max": {
			sem:  &SemaphoreMutex{max: 5},
			n:    9,
			want: false,
		},
		"no space": {
			sem:  &SemaphoreMutex{max: 5, count: 4},
			n:    5,
			want: false,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {

			res := tc.sem.TryAcquire(tc.n)

			assert.Equal(t, tc.want, res)

		})
	}

}

func TestSemaphoreMutex_Release(t *testing.T) {
	tests := map[string]struct {
		sem       *SemaphoreMutex
		n         int64
		wantCount int64
		wantPanic any
	}{
		"less zero": {
			sem:       &SemaphoreMutex{max: 5, count: 3},
			n:         -1,
			wantCount: 3,
			wantPanic: errorsList.ErrParamLessZero,
		},
		"zero": {
			sem:       &SemaphoreMutex{max: 5, count: 3},
			n:         0,
			wantCount: 3,
			wantPanic: errorsList.ErrParamLessZero,
		},
		"release more than acquired": {
			sem:       &SemaphoreMutex{max: 5, count: 2},
			n:         3,
			wantCount: 2,
			wantPanic: errorsList.ErrParmMoreAvailable,
		},
		"correct release": {
			sem:       &SemaphoreMutex{max: 5, count: 4},
			n:         2,
			wantCount: 2,
			wantPanic: nil,
		},
		"release all": {
			sem:       &SemaphoreMutex{max: 5, count: 5},
			n:         5,
			wantCount: 0,
			wantPanic: nil,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			if tc.wantPanic != nil {
				assert.PanicsWithValue(t, tc.wantPanic, func() {
					tc.sem.Release(tc.n)
				})

				assert.Equal(t, tc.wantCount, tc.sem.count)
				return
			}

			assert.NotPanics(t, func() {
				tc.sem.Release(tc.n)
			})

			assert.Equal(t, tc.wantCount, tc.sem.count)
		})
	}
}
