package domain

import evententity "minimulti/src/events/domain/event_entity"

type EventInterface interface {
	Create(event *evententity.Event) (int, error)
	GetAll() ([]*evententity.Event, error)
	DeleteAll() error
}
