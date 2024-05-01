package tools

import (
	"strings"

	"gitlab.com/slon/shad-go/gitfame/internal/entities"
)

type lineType uint8

const (
	commit lineType = iota
	author
	authorMail
	authorTime
	authorTz
	committer
	committerMail
	committerTime
	committerTz
	summary
	boundaryOrPrevious
	filename
	code
)

const (
	blockSize = 13
)

type BlameParser struct {
	useCommiter bool
}

func NewBlameParser(useCommitter bool) *BlameParser {
	return &BlameParser{
		useCommiter: useCommitter,
	}
}

func (p *BlameParser) Parse(filename string, content []byte) map[string]*entities.Participant {
	var (
		lines = strings.Split(string(content), "\n")
		stat  = make(map[string]*entities.Participant)
	)

	var cursor int
	for cursor < len(lines) {
		if cursor == len(lines)-1 {
			break
		}
		var participantName string

		switch p.useCommiter {
		case true:
			participantName = strings.TrimPrefix(lines[cursor+int(committer)], "committer ")
		case false:
			participantName = strings.TrimPrefix(lines[cursor+int(author)], "author ")
		}

		var (
			participant *entities.Participant
			ok          bool
		)

		participant, ok = stat[participantName]
		if !ok {
			participant = &entities.Participant{
				CommitsList: make(map[string]struct{}),
				Name:        participantName,
				Files:       1,
			}
			stat[participantName] = participant
		}

		participant.Lines++

		commitHash := strings.Split(lines[cursor+int(commit)], " ")[0]
		if _, ok = participant.CommitsList[commitHash]; !ok {
			participant.Commits++
			participant.CommitsList[commitHash] = struct{}{}
		}

		firstWord := strings.Split(lines[cursor+int(boundaryOrPrevious)], " ")[0]
		if firstWord != "previous" && firstWord != "boundary" {
			cursor += blockSize - 1
			continue
		}
		cursor += blockSize
	}
	return stat
}
