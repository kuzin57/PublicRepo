package participant

import (
	"context"
	"sort"

	"gitlab.com/slon/shad-go/gitfame/internal/entities"
	"gitlab.com/slon/shad-go/gitfame/internal/repositories"
	"gitlab.com/slon/shad-go/gitfame/internal/services"
	"gitlab.com/slon/shad-go/gitfame/pkg/writer"
)

type service struct {
	repo         repositories.IParticipantRepository
	resultWriter writer.IWriter
	orderBy      string
}

func NewService(
	orderBy string,
	repo repositories.IParticipantRepository,
	resultWriter writer.IWriter,
) services.IParticipantService {
	return &service{
		repo:         repo,
		orderBy:      orderBy,
		resultWriter: resultWriter,
	}
}

func (s *service) Update(ctx context.Context, stat map[string]*entities.Participant) error {
	return s.repo.UpdateOrAdd(ctx, stat)
}

func cmp(intMetricsFirst, intMetricsSecond []int, nameFirst, nameSecond string) bool {
	for i := range len(intMetricsFirst) {
		if intMetricsFirst[i] > intMetricsSecond[i] {
			return true
		}

		if intMetricsFirst[i] < intMetricsSecond[i] {
			return false
		}
	}

	return nameFirst < nameSecond
}

func (s *service) Show() error {
	allParticipants := s.repo.FindAll()

	sort.Slice(allParticipants, func(i, j int) bool {
		switch s.orderBy {
		case "lines":
			return cmp(
				[]int{allParticipants[i].Lines, allParticipants[i].Commits, allParticipants[i].Files},
				[]int{allParticipants[j].Lines, allParticipants[j].Commits, allParticipants[j].Files},
				allParticipants[i].Name,
				allParticipants[j].Name,
			)
		case "commits":
			return cmp(
				[]int{allParticipants[i].Commits, allParticipants[i].Lines, allParticipants[i].Files},
				[]int{allParticipants[j].Commits, allParticipants[j].Lines, allParticipants[j].Files},
				allParticipants[i].Name,
				allParticipants[j].Name,
			)
		case "files":
			return cmp(
				[]int{allParticipants[i].Files, allParticipants[i].Lines, allParticipants[i].Commits},
				[]int{allParticipants[j].Files, allParticipants[j].Lines, allParticipants[j].Commits},
				allParticipants[i].Name,
				allParticipants[j].Name,
			)
		}
		return false
	})

	return s.resultWriter.Write(allParticipants)
}
