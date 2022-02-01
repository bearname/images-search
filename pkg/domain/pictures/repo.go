package pictures

type Repository interface {
	IsExists(pictureId string) error
	FindPicture(pictureId string) (*PictureDTO, error)
	Search(dto *SearchPictureDto) (*SearchPictureResultDto, error)
	SaveInitialPicture(image *InitialImage) (int, error)
	SaveInitialPictures(image *InitialDropboxImage) (*InitialDropboxImageResult, error)
	UpdateImageHandle(picture *Picture) error
	Store(imageTextDetectionDto *TextDetectionOnImageDto) error
	Delete(imageId string) error
}
