package domain

import "github.com/col3name/images-search/pkg/domain/pictures"

type TextDetector interface {
	DetectTextFromImage(imageBytes *[]byte, minConfidence int) ([]pictures.TextDetection, error)
}
