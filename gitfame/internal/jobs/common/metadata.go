package common

type JobName string

const (
	ReadDirJob JobName = "read_dir"
	ParseJob   JobName = "parse_job"
)

type JobMetadata struct {
	Name     JobName `yaml:"name"`
	Parallel int     `yaml:"parallel"`
}

func (m *JobMetadata) GetName() string {
	return string(m.Name)
}

func (m *JobMetadata) GetParallel() int {
	return m.Parallel
}
