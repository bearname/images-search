package domain

import "photofinish/pkg/domain/pictures"

type TextDetector interface {
	DetectTextFromImage(imageBytes *[]byte, minConfidence int) ([]pictures.TextDetection, error)
}
