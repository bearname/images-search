package event

import (
	"photofinish/pkg/domain"
)

type Service interface {
	Create(event *CreateEventInputDto) (int, error)
	DeleteEvent(eventId int) error
	Search(page domain.Page) ([]Event, error)
}
