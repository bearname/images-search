package s3

import (
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"io"
	"photofinish/pkg/domain"
)

type MockUploader struct {
}

func NewMockUploader() *MockUploader {
	s := new(MockUploader)
	return s
}

func (s *MockUploader) Upload(filename string, file io.Reader, acl types.ObjectCannedACL) (*domain.UploadOutput, error) {
	//index := strings.Index(filename, ".")

	var id = "uploadOutput.VersionID"
	return &domain.UploadOutput{
		Location:  " uploadOutput.Location",
		VersionID: &id,
		UploadID:  "uploadOutput.UploadID",
	}, nil
}
