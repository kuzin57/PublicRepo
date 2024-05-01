package common

type Params struct {
	RepositoryPath string `names:"--repository, -r" usage:"path to git repo" default:"."`
	Revision       string `names:"--revision, -rev" usage:"revision hash" default:""`
	ExcludeArg     string `names:"--exclude, -e" usage:"files to exclude" default:""`
	RestrictToArg  string `names:"--restrict-to, -rt" usage:"restriction patterns on files" default:""`
	LanguagesArg   string `names:"--languages, -l" usage:"programming languages" default:""`
	ExtensionsArg  string `names:"--extensions, -ext" usage:"files extensions" default:""`
	OrderBy        string `names:"--order-by, -o" usage:"results sort key" default:"lines"`
	UseCommiter    bool   `names:"--use-committer, -uc" usage:"use commiter to calculate results" default:""`
	Format         string `names:"--format, -f" usage:"format of result output" default:"tabular"`
	RestrictTo     []string
	Exclude        []string
	Languages      []string
	Extensions     []string
}
