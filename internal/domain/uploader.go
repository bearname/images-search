package domain

import (
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"io"
)

type Uploader interface {
	Upload(filename string, file io.Reader, acl types.ObjectCannedACL) (*UploadOutput, error)
}

type UploadOutput struct {
	// The URL where the object was uploaded to.
	Location string

	// The version of the object that was uploaded. Will only be populated if
	// the S3 Bucket is versioned. If the bucket is not versioned this field
	// will not be set.
	VersionID *string

	// The ID for a multipart upload to S3. In the case of an error the error
	// can be cast to the MultiUploadFailure interface to extract the upload ID.
	UploadID string
}