package pipeline

import (
	"context"
	"log"
)

type (
	In  = <-chan any
	Out = In

	Stage func(data any) (any, error)
)

type Pipeline interface {
	ExecutePipeline(ctx context.Context, in In, stages ...Stage) Out
}

type pipeline struct {
	logger *log.Logger
}

func NewPipeline(logger *log.Logger) Pipeline {
	if logger == nil {
		logger = log.Default()
	}

	return &pipeline{
		logger: logger,
	}
}

func (p *pipeline) ExecutePipeline(ctx context.Context, in In, stages ...Stage) Out {
	out := in

	for index, stage := range stages {
		out = p.runStage(ctx, index, out, stage)
	}

	return out
}

func (p *pipeline) runStage(ctx context.Context, stageIndex int, in In, stage Stage) Out {

	out := make(chan any)

	go func() {
		defer close(out)

		p.logger.Printf("Running stage %d", stageIndex)
		defer p.logger.Printf("Finished stage %d", stageIndex)

		for {
			select {
			case <-ctx.Done():
				p.logger.Printf("Stage %d. Context done", stageIndex)
				return

			case data, ok := <-in:
				if !ok {
					p.logger.Printf("Stage %d. Channel closed", stageIndex)
					return
				}

				result, err := stage(data)
				if err != nil {
					p.logger.Printf("Stage %d. Error: %s", stageIndex, err)
					continue
				}

				select {
				case <-ctx.Done():
					p.logger.Printf("Stage %d. Context done", stageIndex)
					return
				case out <- result:
					p.logger.Printf("Stage %d. Success", stageIndex)
				}

			}
		}

	}()

	return out
}
