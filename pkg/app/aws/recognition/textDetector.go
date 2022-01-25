package recognition

import (
	"github.com/aws/aws-sdk-go/service/rekognition"
	//"photofinish/pkg/app/aws/recognition/rekognition"
	"photofinish/pkg/common/util"
	"photofinish/pkg/domain"
	"photofinish/pkg/domain/pictures"
)

type AmazonTextRecognition struct {
	awsRekognition *rekognition.Rekognition
}

func NewAmazonTextRecognition(awsRekognition *rekognition.Rekognition) *AmazonTextRecognition {
	s := new(AmazonTextRecognition)
	s.awsRekognition = awsRekognition
	return s
}

func (s *AmazonTextRecognition) DetectTextFromImage(imageBytes *[]byte, minConfidence int) ([]pictures.TextDetection, error) {
	var decodedImage []byte
	decodedImage, err := util.ImageBase64(*imageBytes)

	if err != nil {
		return []pictures.TextDetection{}, err
	}

	input := &rekognition.DetectTextInput{
		Image: &rekognition.Image{
			Bytes: decodedImage,
		},
	}

	result, err := s.awsRekognition.DetectText(input)
	if err != nil {
		return nil, err
	}

	set := domain.MakeSet()
	for _, detection := range result.TextDetections {
		if int(*detection.Confidence) >= minConfidence {
			text := *detection.DetectedText
			numbers := util.ExtractNumberFromString(text)
			for _, number := range numbers {
				set.Add(number, *detection.Confidence)
			}
		}
	}

	var arr []pictures.TextDetection

	for detectedText, confidence := range set.GetAll() {
		detection := pictures.NewTextDetection(detectedText, confidence)
		arr = append(arr, *detection)
	}

	return arr, nil
}
