package common

import "gitlab.com/slon/shad-go/gitfame/internal/repositories"

type Repositories struct {
	ParticipantRepository repositories.IParticipantRepository
}

func NewRepositories() *Repositories {
	return &Repositories{
		ParticipantRepository: repositories.NewParticipantRepository(),
	}
}
