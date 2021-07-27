package picture

import "aws_rekognition_demo/internal/domain"

type TextDetectionOnImageDto struct {
    EventId       int64
    OriginalPath  string
    PreviewPath   string
    TextDetection []TextDetection
}

func NewImageTextDetectionDto(eventId int64, originalPath string, previewPath string, textDetections []TextDetection) *TextDetectionOnImageDto {
    t := new(TextDetectionOnImageDto)
    t.EventId = eventId
    t.OriginalPath = originalPath
    t.PreviewPath = previewPath
    t.TextDetection = textDetections
    return t
}

type SearchPictureDto struct {
    ParticipantNumber int
    Confidence        int
    Page              domain.Page
    EventId           int
}

func NewSearchPictureDto(participantNumber int, confidence int, eventId int, page domain.Page) SearchPictureDto {
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
    Pictures   []SearchPictureItem
}
