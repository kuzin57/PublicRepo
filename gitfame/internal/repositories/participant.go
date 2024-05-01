package repositories

import (
	"context"
	"sync"

	"gitlab.com/slon/shad-go/gitfame/internal/entities"
)

type participantRepository struct {
	mu  *sync.Mutex
	buf map[string]*entities.Participant
}

func NewParticipantRepository() IParticipantRepository {
	return &participantRepository{
		mu:  &sync.Mutex{},
		buf: make(map[string]*entities.Participant),
	}
}

func (r *participantRepository) UpdateOrAdd(ctx context.Context, stat map[string]*entities.Participant) error {
	for name, participant := range stat {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			func() {
				r.mu.Lock()
				defer r.mu.Unlock()

				if p, ok := r.buf[name]; ok {
					p.Add(participant)
					return
				}

				r.buf[name] = participant
			}()
		}
	}
	return nil
}

func (r *participantRepository) FindAll() []entities.Participant {
	result := make([]entities.Participant, 0, len(r.buf))
	for _, participant := range r.buf {
		result = append(result, *participant)
	}

	return result
}
