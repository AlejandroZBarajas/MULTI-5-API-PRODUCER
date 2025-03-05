package application

import (
	"fmt"
	"minimulti/src/events/domain"
	evententity "minimulti/src/events/domain/event_entity"
)

type CreateEvent struct {
	repo domain.EventInterface
}

func NewCreateEvent(repo domain.EventInterface) *CreateEvent {
	return &CreateEvent{repo: repo}
}

func (ce *CreateEvent) Run(Device_name string) error {
	event := evententity.CreateEvent(Device_name)

	err := ce.repo.Create(event)
	if err != nil {
		return fmt.Errorf("error :'%w'", err)
	}
	return nil
}
