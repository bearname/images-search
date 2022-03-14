package event

import (
	"github.com/col3name/images-search/pkg/domain/dto"
)

type Service interface {
	Create(event *CreateEventInputDto) (int, error)
	DeleteEvent(eventId int) error
	Search(page *dto.Page) ([]Event, error)
}
