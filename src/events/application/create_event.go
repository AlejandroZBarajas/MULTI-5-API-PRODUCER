package application

import (
	"minimulti/src/events/domain"
	evententity "minimulti/src/events/domain/event_entity"
)

type CreateEvent struct {
	repo domain.EventInterface
}

func NewCreateEvent(repo domain.EventInterface) *CreateEvent {
	return &CreateEvent{repo: repo}
}

func (ce *CreateEvent) Run(serie string, title string, description float32, emitter string, topic string) (int, string, error) {
	event := evententity.CreateEvent(serie, title, description, emitter, topic)

	id, createdAt, err := ce.repo.Create(event)
	if err != nil {
		return 0, "", err
	}
	return id, createdAt, nil
}
