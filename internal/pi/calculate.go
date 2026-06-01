package pi

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

var ErrCountWorkers = errors.New("workers must be greater than 0")
var ErrPrecision = errors.New("precision must be greater than 0")

const batchSize = 10000

type CalculateResult struct {
	pi         float64
	iterations uint64
	duration   time.Duration
}

func Run() error {
	workers := flag.Int("workers", 1, "Number of workers")
	precision := flag.Int("precision", 1, "precision")
	timeout := flag.Duration("timeout", 0, "timeout")
	isVerbose := flag.Bool("verbose", false, "Verbose mode")

	flag.Parse()

	if *workers <= 0 {
		return ErrCountWorkers
	}

	if *precision <= 0 {
		return ErrPrecision
	}

	signalCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	ctx := signalCtx

	if *timeout > 0 {
		timeoutCtx, cancel := context.WithTimeout(ctx, *timeout)
		defer cancel()

		ctx = timeoutCtx
	}

	result := calculate(ctx, *workers)

	fmt.Printf("%.*f\n", *precision, result.pi)

	if *isVerbose {
		fmt.Printf("iterations: %d\n", result.iterations)
		fmt.Printf("duration: %s\n", result.duration)
	}

	return nil
}

func calculate(ctx context.Context, workers int) CalculateResult {
	start := time.Now()

	var wg sync.WaitGroup

	sums := make([]float64, workers)
	iterations := make([]uint64, workers)

	for id := 0; id < workers; id++ {
		wg.Add(1)

		go func(id int) {
			defer wg.Done()

			var localSum float64
			var localIterations uint64

			n := uint64(id)
			step := uint64(workers)

			for {
				select {
				case <-ctx.Done():
					sums[id] = localSum
					iterations[id] = localIterations
					return
				default:
				}

				for i := 0; i < batchSize; i++ {
					term := 1 / float64(2*n+1)

					if n%2 == 0 {
						localSum += term
					} else {
						localSum -= term
					}

					n += step
					localIterations++
				}

			}

		}(id)

	}

	wg.Wait()

	var totalSum float64
	var totalIterations uint64

	for i := 0; i < workers; i++ {
		totalSum += sums[i]
		totalIterations += iterations[i]
	}

	return CalculateResult{
		pi:         4 * totalSum,
		iterations: totalIterations,
		duration:   time.Since(start),
	}
}
