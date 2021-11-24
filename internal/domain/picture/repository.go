package picture

type Repository interface {
	FindById(imageId string) error
	Store(imageTextDetectionDto *TextDetectionOnImageDto) error
	StoreAll(arr []*TextDetectionOnImageDto) error
	Search(dto SearchPictureDto) (SearchPictureResultDto, error)
	Delete(imageId string) error
}
