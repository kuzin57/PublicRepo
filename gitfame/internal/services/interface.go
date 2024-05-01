package services

import (
	"context"

	"gitlab.com/slon/shad-go/gitfame/internal/entities"
)

type IParticipantService interface {
	Show() error
	Update(context.Context, map[string]*entities.Participant) error
}
