package event

type CreateEventInputDto struct {
	Name     string `json:"name"`
	Location string `json:"location"`
	Date     string `json:"date"`
}

//
//func NewImageTextDetectionDto(eventId int64, originalPath string, previewPath string, textDetections []TextDetection) *TextDetectionOnImageDto {
//    t := new(TextDetectionOnImageDto)
//    t.EventId = eventId
//    t.OriginalPath = originalPath
//    t.PreviewPath = previewPath
//    t.TextDetection = textDetections
//    return t
//}
//
//type SearchPictureDto struct {
//    ParticipantNumber int
//    Confidence        int
//    Page              domain.Page
//    EventId           int
//}
//
//func NewSearchPictureDto(participantNumber int, confidence int, eventId int, page domain.Page) SearchPictureDto {
//    return SearchPictureDto{
//        ParticipantNumber: participantNumber,
//        Confidence:        confidence,
//        EventId:           eventId,
//        Page:              page,
//    }
//}
//
//type Event struct {
//    EventId   int64
//    EventName string
//}
//
//type TextDetectionDto struct {
//    TextDetection
//    Event
//}
//
//type SearchPictureItem struct {
//    PictureId      int
//    Path           string
//    TextDetections []TextDetectionDto
//}
//
//type SearchPictureResultDto struct {
//    Pictures []SearchPictureItem
//}
