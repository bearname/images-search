package pictures

import (
	"errors"
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

var ErrNotFound = errors.New("pictures not exist")

type Picture struct {
	Id      uuid.UUID
	EventId int

	TaskId string

	DropboxPath  string
	OriginalPath string
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

func NewTextDetection(detectedText string, confidence float64) *TextDetection {
	t := new(TextDetection)
	t.DetectedText = detectedText
	t.Confidence = confidence
	return t
}
