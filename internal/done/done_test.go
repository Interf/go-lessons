package done

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestOr(t *testing.T) {
	ch1 := make(chan any)
	ch2 := make(chan any)
	ch3 := make(chan any)

	result := Or(ch1, ch2, ch3)

	requireNotClosed(t, result)

	close(ch2)

	requireClosed(t, result)
}

func TestOr2(t *testing.T) {
	ch1 := make(chan any)
	ch2 := make(chan any)
	ch3 := make(chan any)

	result := Or2(ch1, ch2, ch3)

	requireNotClosed(t, result)

	close(ch2)

	requireClosed(t, result)
}

func TestOrNil(t *testing.T) {
	result := Or()

	require.Nil(t, result)
}

func TestOr2Nil(t *testing.T) {
	result := Or2()

	require.Nil(t, result)
}

func TestOrWithOneChannel(t *testing.T) {
	ch := make(chan any)

	result := Or(ch)

	require.Equal(t, (<-chan any)(ch), result)
}

func TestOr2WithOneChannel(t *testing.T) {
	ch := make(chan any)

	result := Or2(ch)

	require.Equal(t, (<-chan any)(ch), result)
}

func requireClosed(t *testing.T, ch <-chan any) {
	t.Helper()

	select {
	case <-ch:
	case <-time.After(100 * time.Millisecond):
		require.Fail(t, "expected channel to be closed")
	}
}

func requireNotClosed(t *testing.T, ch <-chan any) {
	t.Helper()

	select {
	case <-ch:
		require.Fail(t, "expected channel to be closed")
	case <-time.After(10 * time.Millisecond):
		return
	}
}

func TestOr2WithFileChannels(t *testing.T) {
	ch1 := make(chan any)
	ch2 := make(chan any)
	ch3 := make(chan any)
	ch4 := make(chan any)
	ch5 := make(chan any)

	result := Or2(ch1, ch2, ch3, ch4, ch5)

	requireNotClosed(t, result)

	close(ch4)

	requireClosed(t, result)
}

func TestOr2CloseFirstChannel(t *testing.T) {
	ch1 := make(chan any)
	ch2 := make(chan any)
	ch3 := make(chan any)

	result := Or2(ch1, ch2, ch3)

	requireNotClosed(t, result)

	close(ch1)

	requireClosed(t, result)
}
