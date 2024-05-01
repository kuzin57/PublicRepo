package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"

	appCommon "gitlab.com/slon/shad-go/gitfame/internal/app/common"
	"gitlab.com/slon/shad-go/gitfame/internal/config"
	"gitlab.com/slon/shad-go/gitfame/internal/jobs/common"
	"gitlab.com/slon/shad-go/gitfame/internal/jobs/read_dir"
	"gitlab.com/slon/shad-go/gitfame/internal/repositories"
	"gitlab.com/slon/shad-go/gitfame/internal/services/participant"
	"gitlab.com/slon/shad-go/gitfame/pkg/writer"
	csvWriter "gitlab.com/slon/shad-go/gitfame/pkg/writer/csv"
	jsonWriter "gitlab.com/slon/shad-go/gitfame/pkg/writer/json"
	jsonLWriter "gitlab.com/slon/shad-go/gitfame/pkg/writer/json_lines"
	"gitlab.com/slon/shad-go/gitfame/pkg/writer/tab"
	"gopkg.in/yaml.v2"
)

const (
	configPath          = "../../configs/config.yaml"
	languagesConfigPath = "../../configs/language_extensions.json"
	defaultRevision     = "HEAD"
	cuttedRevLength     = 6
)

type App struct {
	cfg *config.Config

	params       *appCommon.Params
	resources    *appCommon.Resources
	repositories *appCommon.Repositories
	services     *appCommon.Services
}

type languagesConfig struct {
	Languages []struct {
		Name       string   `json:"name"`
		Extensions []string `json:"extensions"`
	} `json:"languages"`
}

func (c *languagesConfig) mapNamesToExt() map[string][]string {
	res := make(map[string][]string, len(c.Languages))
	for _, lang := range c.Languages {
		res[strings.ToLower(lang.Name)] = lang.Extensions
	}
	return res
}

func (a *App) getConfig() error {
	content, err := os.ReadFile(configPath)
	if err != nil {
		return err
	}

	a.cfg = &config.Config{}
	if e := yaml.Unmarshal(content, a.cfg); e != nil {
		return e
	}

	languagesContent, err := os.ReadFile(languagesConfigPath)
	if err != nil {
		return err
	}

	lCfg := &languagesConfig{}
	if err := json.Unmarshal(languagesContent, lCfg); err != nil {
		return err
	}

	nameToExt := lCfg.mapNamesToExt()
	for _, lang := range a.params.Languages {
		extensions, ok := nameToExt[lang]
		if !ok {
			continue
		}

		for _, ext := range extensions {
			a.params.RestrictTo = append(a.params.RestrictTo, "*"+ext)
		}
	}

	for _, ext := range a.params.Extensions {
		a.params.RestrictTo = append(a.params.RestrictTo, "*"+ext)
	}
	return nil
}

func (a *App) initRevision(ctx context.Context) error {
	cmd := exec.CommandContext(ctx, "git", "rev-parse", a.params.Revision)
	cmd.Dir = a.params.RepositoryPath
	output, err := cmd.Output()
	if err != nil {
		return err
	}

	a.params.Revision = string(output)[:cuttedRevLength]
	return nil
}

func (a *App) complementParams(ctx context.Context) error {
	if a.params.Revision == "" {
		a.params.Revision = defaultRevision
	}

	return a.initRevision(ctx)
}

func (a *App) initResources() {
	a.resources = appCommon.NewResources(a.params)
}

func (a *App) registerJobs() {
	for _, m := range a.cfg.Jobs {
		a.resources.JobManager.RegisterJob(m)
	}
}

func (a *App) initRepositories() {
	a.repositories = &appCommon.Repositories{
		ParticipantRepository: repositories.NewParticipantRepository(),
	}
}

func (a *App) initServices() error {
	var resultWriter writer.IWriter
	switch a.params.Format {
	case "json":
		resultWriter = jsonWriter.NewWriter(os.Stdout)
	case "csv":
		resultWriter = csvWriter.NewWriter(os.Stdout)
	case "json-lines":
		resultWriter = jsonLWriter.NewWriter(os.Stdout)
	case "tabular":
		resultWriter = tab.NewWriter(os.Stdout)
	default:
		return fmt.Errorf("unknown format type")
	}

	a.services = &appCommon.Services{
		ParticipantService: participant.NewService(
			a.params.OrderBy,
			a.repositories.ParticipantRepository,
			resultWriter,
		),
	}

	return nil
}

func (a *App) Init(ctx context.Context) error {
	if err := a.getConfig(); err != nil {
		return err
	}

	if err := a.complementParams(ctx); err != nil {
		return err
	}

	a.initResources()
	a.initRepositories()

	if err := a.initServices(); err != nil {
		return err
	}

	a.registerJobs()

	return nil
}

func (a *App) Run(ctx context.Context) error {
	var (
		wg = &sync.WaitGroup{}

		err error
	)

	defer func() {
		if e := a.services.ParticipantService.Show(); e != nil {
			panic(e)
		}
	}()

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	wg.Add(1)
	go func() {
		defer wg.Done()
		err = a.resources.JobManager.Run(ctx)
	}()

	a.resources.Monitor.UpdateReadDirJobs(1)
	if e := a.resources.JobManager.EnqueueJob(
		ctx,
		readdir.NewJob(
			a.params.RepositoryPath,
			a.cfg.Jobs[common.ReadDirJob],
			a.resources,
			a.services,
			a.params,
		),
	); e != nil {
		return e
	}

	<-a.resources.Monitor.Finish()
	cancel()

	wg.Wait()
	if errors.Is(err, ctx.Err()) {
		return nil
	}
	return err
}
