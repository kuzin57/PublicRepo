package defaults

import (
	"context"
	"errors"
	"sync"

	"gitlab.com/slon/shad-go/gitfame/pkg/awg"
	"gitlab.com/slon/shad-go/gitfame/pkg/job"
)

type jobManager struct {
	mu       *sync.Mutex
	jobsInfo map[string]job.IJobMetadata
	queues   map[string]job.IQueue
}

func NewJobManager() job.IManager {
	return &jobManager{
		mu:       &sync.Mutex{},
		queues:   make(map[string]job.IQueue),
		jobsInfo: make(map[string]job.IJobMetadata),
	}
}

func (m *jobManager) RegisterJob(jm job.IJobMetadata) {
	m.jobsInfo[jm.GetName()] = jm
	m.queues[jm.GetName()] = NewQueue[job.IJob]()
}

func (m *jobManager) EnqueueJob(ctx context.Context, j job.IJob) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.queues[j.GetMetadata().GetName()]; !ok {
		return ErrJobNotRegistered
	}

	m.queues[j.GetMetadata().GetName()].Push(ctx, j)
	return nil
}

func (m *jobManager) Run(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	wg := &awg.AdvancedWaitGroup{}

	for _, metadata := range m.jobsInfo {
		metadata := metadata
		wg.Add(func() error {
			return m.startWorkers(ctx, metadata)
		})
	}
	return errors.Join(wg.Start().SetStopOnError(true, cancel).GetAllErrors()...)
}

func (m *jobManager) startWorkers(ctx context.Context, metadata job.IJobMetadata) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	wg := awg.AdvancedWaitGroup{}

	q := m.queues[metadata.GetName()]
	for range metadata.GetParallel() {
		wg.Add(func() error {
			w := NewWorker(q)
			return w.Work(ctx)
		})
	}
	return errors.Join(wg.Start().SetStopOnError(true, cancel).GetAllErrors()...)
}
