package event

import (
	"github.com/col3name/images-search/pkg/domain/dto"
)

type Repository interface {
	CheckExist(eventId int) error
	Store(imageTextDetectionDto *CreateEventInputDto) (int, error)
	Delete(eventId int) error
	FindAll(page *dto.Page) ([]Event, error)
}
