package event

import "aws_rekognition_demo/internal/domain"

type Repository interface {
    CheckExist(eventId int) error
    Store(imageTextDetectionDto *CreateEventInputDto) (int, error)
    Delete(eventId int) error
    FindAll(page domain.Page) ([]Event, error)
    //StoreAll(arr []*TextDetectionOnImageDto) error
    //Search(dto SearchPictureDto) ([]SearchPictureItem, error)
}
