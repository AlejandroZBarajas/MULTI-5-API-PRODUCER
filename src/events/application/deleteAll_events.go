package application

import "minimulti/src/events/domain"

type DeletEvents struct {
	repo domain.EventInterface
}

func NewDeletEvents(repo domain.EventInterface) *DeletEvents {
	return &DeletEvents{repo: repo}
}

func (de *DeletEvents) Run() error {
	return de.repo.DeleteAll()
}
