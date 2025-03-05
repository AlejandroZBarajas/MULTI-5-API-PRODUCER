package application

import (
	"minimulti/src/events/domain"
	evententity "minimulti/src/events/domain/event_entity"
)

type GetAllEvents struct {
	repo domain.EventInterface
}

func NewGetAllEvents(repo domain.EventInterface) *GetAllEvents {
	return &GetAllEvents{repo: repo}
}

func (gae *GetAllEvents) Execute() ([]*evententity.Event, error) {
	return gae.repo.GetAll()
}
