package pictures

import (
	"photofinish/pkg/common/util/uuid"
	"time"
)

const ValueNotSet = 99999999

type ProcessingStatus int

const (
	Success ProcessingStatus = iota
	Processing
	Failed
	TooBig
	Deleted
)

type Picture struct {
	Id      uuid.UUID
	EventId int

	TaskId string

	DropboxPath  string
	OriginalS3Id string
	OriginalPath string
	PreviewS3Id  string
	PreviewPath  string

	Attempts         int
	ProcessingStatus ProcessingStatus
	ExecuteAfter     time.Time
	UpdatedAt        time.Time

	IsOriginalSaved  bool
	IsPreviewSaved   bool
	IsTextRecognized bool
	IsMobileSaved    bool

	DetectedTexts []TextDetection
}

type TextDetection struct {
	DetectedText string
	Confidence   float64
}

type SupportedImgType string

const (
	JPEG SupportedImgType = "jpeg"
	JPG  SupportedImgType = "jpg"
	PNG  SupportedImgType = "png"
)

func NewTextDetection(detectedText string, confidence float64) *TextDetection {
	t := new(TextDetection)
	t.DetectedText = detectedText
	t.Confidence = confidence
	return t
}
