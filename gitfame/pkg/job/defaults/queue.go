package defaults

import (
	"context"
	"sync"

	"gitlab.com/slon/shad-go/gitfame/pkg/job"
)

type queue struct {
	mu       *sync.Mutex
	notEmpty chan struct{}
	buf      []job.IJob
}

func NewQueue[T any]() job.IQueue {
	mu := &sync.Mutex{}
	return &queue{
		mu:       mu,
		notEmpty: make(chan struct{}, 1),
	}
}

func (q *queue) wait(ctx context.Context) {
	q.mu.Unlock()
	select {
	case <-ctx.Done():
	case <-q.notEmpty:
	}
	q.mu.Lock()
}

func (q *queue) Pop(ctx context.Context) job.IJob {
	q.mu.Lock()
	defer q.mu.Unlock()

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			if len(q.buf) == 0 {
				q.wait(ctx)
				continue
			}

			first := q.buf[0]
			q.buf = q.buf[1:]
			return first
		}
	}
}

func (q *queue) Len() int {
	q.mu.Lock()
	defer q.mu.Unlock()

	return len(q.buf)
}

func (q *queue) wakeAll(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case q.notEmpty <- struct{}{}:
		default:
			<-q.notEmpty
			return
		}
	}
}

func (q *queue) Push(ctx context.Context, j job.IJob) {
	q.mu.Lock()
	defer q.mu.Unlock()

	select {
	case <-ctx.Done():
		return
	default:
		q.buf = append(q.buf, j)
	}

	q.wakeAll(ctx)
}
