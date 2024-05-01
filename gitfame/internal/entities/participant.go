package entities

type Participant struct {
	CommitsList map[string]struct{} `json:"-" csv:"-" tab:"-"`
	Name        string              `json:"name" csv:"Name" tab:"Name"`
	Lines       int                 `json:"lines" csv:"Lines" tab:"Lines"`
	Commits     int                 `json:"commits" csv:"Commits" tab:"Commits"`
	Files       int                 `json:"files" csv:"Files" tab:"Files"`
}

func (p *Participant) Add(other *Participant) {
	p.Files += other.Files
	p.Lines += other.Lines

	for commit := range other.CommitsList {
		if _, ok := p.CommitsList[commit]; !ok {
			p.CommitsList[commit] = struct{}{}
			p.Commits++
		}
	}
}
