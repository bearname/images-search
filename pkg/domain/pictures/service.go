package pictures

type Service interface {
	Create(imageTextDetectionDto *TextDetectionOnImageDto) error
	DetectImageFromUrl(root string, eventId int) (*TaskResponse, error)
	Search(dto *SearchPictureDto) (*SearchPictureResultDto, error)
	Delete(imageId string) error
	GetDropboxFolders() ([]string, error)
}
