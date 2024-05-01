package common

import (
	"gitlab.com/slon/shad-go/gitfame/pkg/job"
	"gitlab.com/slon/shad-go/gitfame/pkg/job/defaults"
	"gitlab.com/slon/shad-go/gitfame/tools"
)

type Resources struct {
	JobManager job.IManager
	Validator  *tools.Validator
	Parser     *tools.BlameParser
	Monitor    *tools.Monitor
}

func NewResources(params *Params) *Resources {
	return &Resources{
		JobManager: defaults.NewJobManager(),
		Validator:  tools.NewValidator(params.RestrictTo, params.Exclude),
		Parser:     tools.NewBlameParser(params.UseCommiter),
		Monitor:    tools.NewMonitor(),
	}
}
