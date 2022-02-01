package pictures

import (
	"photofinish/pkg/common/util/uuid"
	"photofinish/pkg/domain/dto"
)

type TextDetectionOnImageDto struct {
	EventId       int64
	OriginalPath  string
	PreviewPath   string
	TextDetection []TextDetection
}

type SearchPictureDto struct {
	ParticipantNumber int
	Confidence        int
	Page              dto.Page
	EventId           int
}

func NewSearchPictureDto(participantNumber int, confidence int, eventId int, page dto.Page) SearchPictureDto {
	return SearchPictureDto{
		ParticipantNumber: participantNumber,
		Confidence:        confidence,
		EventId:           eventId,
		Page:              page,
	}
}

type Event struct {
	EventId   int64
	EventName string
}

type TextDetectionDto struct {
	TextDetection
	Event
}

type SearchPictureItem struct {
	PictureId      string
	Path           string
	TextDetections []TextDetectionDto
}

type SearchPictureResultDto struct {
	CountAllItems int
	Pictures      []SearchPictureItem
}
type InitialDropboxImage struct {
	Images  []string
	EventId int
	Path    string
}

type InitialDropboxImageResult struct {
	ImagesId []uuid.UUID
	TaskId   uuid.UUID
}

type TaskResponse struct {
	TaskId          string
	CountAllImages  int
	CompletedImages int
}

type DropboxImage struct {
	Path    string
	EventId int
	Id      uuid.UUID
}
type DropboxImages struct {
	Images  []DropboxImage
	EventId int
}

type InitialImage struct {
	EventId     int64
	DropboxPath string
}

type PictureDTO struct {
	Id      string
	EventId int

	TaskId string

	DropboxPath  string
	OriginalPath string
	PreviewPath  string

	Attempts         int
	ProcessingStatus ProcessingStatus

	IsOriginalSaved  bool
	IsPreviewSaved   bool
	IsTextRecognized bool
}
