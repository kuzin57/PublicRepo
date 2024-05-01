package parsefile

import (
	"context"
	"fmt"
	"os/exec"
	"strings"

	appCommon "gitlab.com/slon/shad-go/gitfame/internal/app/common"
	"gitlab.com/slon/shad-go/gitfame/internal/entities"
	"gitlab.com/slon/shad-go/gitfame/internal/jobs/common"
	srcJobs "gitlab.com/slon/shad-go/gitfame/pkg/job"
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
) srcJobs.IJob {
	return &job{
		path:      path,
		metadata:  metadata,
		resources: resources,
		services:  services,
		params:    params,
	}
}

func (j *job) GetMetadata() srcJobs.IJobMetadata {
	return j.metadata
}

func (j *job) Run(ctx context.Context) error {
	defer func() {
		j.resources.Monitor.UpdateParseFileJobs(-1)
	}()

	cmd := exec.CommandContext(ctx, "git", "blame", "--line-porcelain", j.params.Revision, j.path)
	cmd.Dir = j.params.RepositoryPath
	output, err := cmd.Output()
	if err != nil {
		return err
	}

	if len(output) == 0 {
		// to save defer
		err := j.parseFileAuthor(ctx)
		return err
	}

	return j.services.ParticipantService.Update(ctx, j.resources.Parser.Parse(j.path, output))
}

func (j *job) parseFileAuthor(ctx context.Context) error {
	cmd := exec.CommandContext(ctx, "git", "log", j.params.Revision, "--", j.path)
	cmd.Dir = j.params.RepositoryPath
	output, err := cmd.Output()
	if err != nil {
		return err
	}

	// Format:

	//0  commit 2de4941d15f39841e7a57aca86886b6246a56808
	//1  Author: kuzin57 <rkuzin.2003@gmail.com>
	//2  Date:   Sat Mar 23 21:54:53 2024 +0300
	//3
	//4  Add test file

	lines := strings.Split(string(output), "\n")
	commitHash := strings.Split(lines[0], " ")[1]

	author := strings.TrimPrefix(lines[1], "Author: ")
	emailStart := strings.Index(author, " <")
	author = author[:emailStart]

	checkBeforeCmd := exec.CommandContext(
		ctx,
		"git",
		"rev-list",
		fmt.Sprintf("%s..%s", commitHash[:len(j.params.Revision)], j.params.Revision),
	)
	checkBeforeCmd.Dir = j.params.RepositoryPath
	output, err = checkBeforeCmd.Output()
	if err != nil {
		return err
	}

	if len(output) == 0 && commitHash[:len(j.params.Revision)] != j.params.Revision {
		return nil
	}

	parsedInfo := map[string]*entities.Participant{
		author: {
			Name:    author,
			Lines:   0,
			Commits: 1,
			Files:   1,
			CommitsList: map[string]struct{}{
				commitHash: {},
			},
		},
	}
	return j.services.ParticipantService.Update(ctx, parsedInfo)
}
