package picture

type Service interface {
	Create(imageTextDetectionDto *TextDetectionOnImageDto) error
	DetectImageFromArchive(root string, minConfidence int, eventId int64) error
	Search(dto SearchPictureDto) (SearchPictureResultDto, error)
	Delete(imageId string) error
}
