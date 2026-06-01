package pi

import (
	"context"
	"math"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func BenchmarkCalculate1Worker(b *testing.B) {
	benchmarkCalculate(b, 1)
}

func BenchmarkCalculate4Worker(b *testing.B) {
	benchmarkCalculate(b, 4)
}

func BenchmarkCalculate100Worker(b *testing.B) {
	benchmarkCalculate(b, 100)
}

func benchmarkCalculate(b *testing.B, workers int) {
	for i := 0; i < b.N; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)

		_ = calculate(ctx, workers)

		cancel()
	}
}

func TestCalculate(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	result := calculate(ctx, 4)

	assert.InDelta(t, math.Pi, result.pi, 0.01)
	assert.Greater(t, result.iterations, uint64(0))
	assert.Greater(t, result.duration, time.Duration(0))
}
