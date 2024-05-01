package config

import jobsCommon "gitlab.com/slon/shad-go/gitfame/internal/jobs/common"

type Config struct {
	Jobs map[jobsCommon.JobName]*jobsCommon.JobMetadata `yaml:"jobs"`
}
