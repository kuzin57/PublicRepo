package repositories

import (
	"context"

	"gitlab.com/slon/shad-go/gitfame/internal/entities"
)

type IParticipantRepository interface {
	FindAll() []entities.Participant
	UpdateOrAdd(context.Context, map[string]*entities.Participant) error
}
