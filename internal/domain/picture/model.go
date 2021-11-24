package picture

const ValueNotSetted = 99999999

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
