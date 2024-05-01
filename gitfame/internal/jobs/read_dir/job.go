package readdir

import (
	"context"
	"os/exec"
	"strings"

	appCommon "gitlab.com/slon/shad-go/gitfame/internal/app/common"
	"gitlab.com/slon/shad-go/gitfame/internal/jobs/common"
	"gitlab.com/slon/shad-go/gitfame/internal/jobs/parse_file"
	srcJob "gitlab.com/slon/shad-go/gitfame/pkg/job"
)

type job struct {
	path      string
	metadata  *common.JobMetadata
	resources *appCommon.Resources
	services  *appCommon.Services
	params    *appCommon.Params
}

func NewJob(
	path string,
	metadata *common.JobMetadata,
	resources *appCommon.Resources,
	services *appCommon.Services,
	params *appCommon.Params,
) srcJob.IJob {
	return &job{
		path:      path,
		metadata:  metadata,
		resources: resources,
		services:  services,
		params:    params,
	}
}

func (j *job) GetMetadata() srcJob.IJobMetadata {
	return j.metadata
}

func (j *job) Run(ctx context.Context) error {
	const (
		fileOrDir = 1
		filename  = 3
	)

	defer func() {
		j.resources.Monitor.UpdateReadDirJobs(-1)
	}()

	cmd := exec.CommandContext(ctx, "git", "ls-tree", j.params.Revision)
	cmd.Dir = j.path
	output, err := cmd.Output()
	if err != nil {
		return err
	}

	for _, line := range strings.Split(string(output), "\n") {
		if line == "" {
			continue
		}

		parts := strings.Fields(line)

		// to process filenames with spaces in name
		name := strings.Join(parts[filename:], " ")
		path := j.path + "/" + name

		if parts[fileOrDir] == "tree" {
			j.resources.Monitor.UpdateReadDirJobs(1)
			if err := j.resources.JobManager.EnqueueJob(
				ctx,
				NewJob(
					path,
					j.metadata,
					j.resources,
					j.services,
					j.params,
				),
			); err != nil {
				return err
			}

			continue
		}

		skip, err := j.resources.Validator.Validate(name, path)
		if err != nil {
			return err
		}

		if skip {
			continue
		}

		j.resources.Monitor.UpdateParseFileJobs(1)

		if err := j.resources.JobManager.EnqueueJob(
			ctx,
			parsefile.NewJob(
				path,
				j.metadata,
				j.resources,
				j.services,
				j.params,
			),
		); err != nil {
			return err
		}
	}

	return nil
}
