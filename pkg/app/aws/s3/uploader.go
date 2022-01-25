package s3

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"io"
	"photofinish/pkg/common/util/uuid"
	"photofinish/pkg/domain"
	"strings"
)

type AwsS3Uploader struct {
	uploader  *manager.Uploader
	awsBucket string
}

func NewAwsS3Uploader(uploader *manager.Uploader, awsBucket string) *AwsS3Uploader {
	s := new(AwsS3Uploader)
	s.uploader = uploader
	s.awsBucket = awsBucket
	return s
}

func (s *AwsS3Uploader) Upload(filename string, file io.Reader, acl types.ObjectCannedACL) (*domain.UploadOutput, error) {
	index := strings.Index(filename, ".")
	fileName := filename[:index] + uuid.Generate().String() + filename[index+1:]

	uploadOutput, err := s.uploader.Upload(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(s.awsBucket),
		Key:    aws.String(fileName),
		Body:   file,
		ACL:    acl,
	})
	if err != nil {
		return nil, err
	}

	return &domain.UploadOutput{
		Location:  uploadOutput.Location,
		VersionID: uploadOutput.VersionID,
		UploadID:  uploadOutput.UploadID,
	}, err
}
