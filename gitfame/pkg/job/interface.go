package job

import "context"

type IJob interface {
	GetMetadata() IJobMetadata
	Run(ctx context.Context) error
}

type IJobMetadata interface {
	GetName() string
	GetParallel() int
}

type IManager interface {
	RegisterJob(job IJobMetadata)
	Run(ctx context.Context) error
	EnqueueJob(ctx context.Context, job IJob) error
}

type IQueue interface {
	Pop(ctx context.Context) IJob
	Push(ctx context.Context, j IJob)
}

type IWorker interface {
	Work(ctx context.Context) error
}
