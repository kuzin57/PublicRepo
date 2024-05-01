package defaults

import (
	"context"

	"gitlab.com/slon/shad-go/gitfame/pkg/job"
)

type worker struct {
	qu job.IQueue
}

func NewWorker(qu job.IQueue) job.IWorker {
	return &worker{
		qu: qu,
	}
}

func (w *worker) Work(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			j := w.qu.Pop(ctx)
			if j != nil {
				if err := j.Run(ctx); err != nil {
					return err
				}
			}
		}
	}
}
