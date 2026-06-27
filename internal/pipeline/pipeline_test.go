package pipeline

import (
	"context"
	"errors"
	"log"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func collect(out Out, timeout time.Duration) ([]any, bool) {
	var result []any
	done := make(chan struct{})

	go func() {
		defer close(done)

		for v := range out {
			result = append(result, v)
		}
	}()

	select {
	case <-done:
		return result, true
	case <-time.After(timeout):
		return result, false
	}
}

func TestNewPipeline_NilLogger(t *testing.T) {
	p := NewPipeline(nil)
	assert.NotNil(t, p)
}

func TestNewPipeline_WithLogger(t *testing.T) {
	logger := log.Default()
	p := NewPipeline(logger)
	require.NotNil(t, p)
}

func TestExecutePipeline_NoStages(t *testing.T) {
	p := NewPipeline(nil)
	ctx := context.Background()
	in := make(chan any, 1)
	in <- 42
	close(in)

	out := p.ExecutePipeline(ctx, in)
	result, closed := collect(out, time.Second)
	require.True(t, closed)
	require.Equal(t, []any{42}, result)
}

func TestExecutePipeline_SingleStage(t *testing.T) {
	p := NewPipeline(nil)
	ctx := context.Background()

	in := make(chan any, 2)
	in <- 1
	in <- 2
	close(in)

	double := func(data any) (any, error) {
		return data.(int) * 2, nil
	}

	out := p.ExecutePipeline(ctx, in, double)
	result, closed := collect(out, time.Second)
	require.True(t, closed)
	require.Equal(t, []any{2, 4}, result)
}

func TestExecutePipeline_MultipleStages(t *testing.T) {
	p := NewPipeline(nil)
	ctx := context.Background()

	in := make(chan any, 1)
	in <- 3
	close(in)

	addOne := func(data any) (any, error) {
		return data.(int) + 1, nil
	}
	mulTen := func(data any) (any, error) {
		return data.(int) * 10, nil
	}

	out := p.ExecutePipeline(ctx, in, addOne, mulTen)
	result, closed := collect(out, time.Second)
	require.True(t, closed)
	require.Equal(t, []any{40}, result)
}

func TestExecutePipeline_StageError(t *testing.T) {
	p := NewPipeline(nil)
	ctx := context.Background()

	in := make(chan any, 3)
	in <- 1
	in <- 2
	in <- 3
	close(in)

	errStage := func(data any) (any, error) {
		if data.(int) == 2 {
			return nil, errors.New("skip")
		}
		return data.(int) * 10, nil
	}

	out := p.ExecutePipeline(ctx, in, errStage)
	result, closed := collect(out, time.Second)
	require.True(t, closed)
	require.Equal(t, []any{10, 30}, result)
}

func TestExecutePipeline_ContextCancel(t *testing.T) {
	p := NewPipeline(nil)
	ctx, cancel := context.WithCancel(context.Background())

	in := make(chan any, 1)
	in <- 1
	close(in)

	stageStarted := make(chan struct{})
	slow := func(data any) (any, error) {
		close(stageStarted)
		<-ctx.Done()
		return nil, ctx.Err()
	}

	out := p.ExecutePipeline(ctx, in, slow)
	<-stageStarted
	cancel()
	result, closed := collect(out, 500*time.Millisecond)
	require.True(t, closed)
	require.Empty(t, result)
}

func TestExecutePipeline_UpstreamClose(t *testing.T) {
	p := NewPipeline(nil)
	ctx := context.Background()

	in := make(chan any)

	stage := func(data any) (any, error) {
		return data.(int) + 1, nil
	}

	out := p.ExecutePipeline(ctx, in, stage)
	close(in)

	result, closed := collect(out, time.Second)
	require.True(t, closed)
	require.Empty(t, result)
}

func TestExecutePipeline_DataFlow(t *testing.T) {
	p := NewPipeline(nil)
	ctx := context.Background()

	in := make(chan any, 5)
	for i := 0; i < 5; i++ {
		in <- i
	}
	close(in)

	inc := func(data any) (any, error) {
		return data.(int) + 1, nil
	}

	out := p.ExecutePipeline(ctx, in, inc, inc, inc)
	result, closed := collect(out, time.Second)
	require.True(t, closed)
	require.Equal(t, []any{3, 4, 5, 6, 7}, result)
}

func TestExecutePipeline_ConcurrentWriters(t *testing.T) {
	p := NewPipeline(nil)
	ctx := context.Background()

	in := make(chan any)

	stage := func(data any) (any, error) {
		return data.(int) * 2, nil
	}

	out := p.ExecutePipeline(ctx, in, stage)

	var wg sync.WaitGroup
	n := 100
	resultCh := make(chan []any, 1)
	go func() {
		result, closed := collect(out, time.Second)
		require.True(t, closed)
		resultCh <- result
	}()

	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(v int) {
			defer wg.Done()
			in <- v
		}(i)
	}
	wg.Wait()
	close(in)

	result := <-resultCh
	assert.Len(t, result, n)
}

func TestExecutePipeline_ManyItems(t *testing.T) {
	p := NewPipeline(nil)
	ctx := context.Background()

	in := make(chan any, 1000)
	for i := 0; i < 1000; i++ {
		in <- i
	}
	close(in)

	identity := func(data any) (any, error) {
		return data, nil
	}

	out := p.ExecutePipeline(ctx, in, identity, identity)
	result, closed := collect(out, 5*time.Second)
	require.True(t, closed)
	require.Len(t, result, 1000)
}
