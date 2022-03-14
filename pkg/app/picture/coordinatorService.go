package picture

import (
	"bytes"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/col3name/images-search/pkg/domain"
	"github.com/col3name/images-search/pkg/domain/pictures"
	"github.com/col3name/images-search/pkg/infrastructure/aws/recognition"
	"github.com/col3name/images-search/pkg/infrastructure/postgres"
	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox/files"
	log "github.com/sirupsen/logrus"
	"strconv"
	"strings"
	"time"
)

const MaxSize = 14680064

type CoordinatorServiceImpl struct {
	maxAttemptsBeforeNotify int
	pictureRepo             *postgres.PictureRepositoryImpl
	downloader              domain.Downloader
	uploader                domain.Uploader
	textDetector            recognition.AmazonTextRecognition
	compressor              domain.ImageCompressor
	notifier                domain.NotifierService
	minConfidence           int
}

func NewCoordinatorServiceImpl(maxAttemptsBeforeNotify int, pictureRepo *postgres.PictureRepositoryImpl,
	downloader domain.Downloader,
	uploader domain.Uploader,
	textDetector *recognition.AmazonTextRecognition,
	compressor domain.ImageCompressor,
	notifier domain.NotifierService,
	minConfidence int) *CoordinatorServiceImpl {

	c := new(CoordinatorServiceImpl)

	c.maxAttemptsBeforeNotify = maxAttemptsBeforeNotify
	c.pictureRepo = pictureRepo
	c.downloader = downloader
	c.uploader = uploader
	c.textDetector = *textDetector
	c.compressor = compressor
	c.notifier = notifier
	c.minConfidence = minConfidence

	return c
}

func (c *CoordinatorServiceImpl) PerformAddImage(image *pictures.Picture) error {
	var originalData *[]byte
	var err error
	var metadata *files.FileMetadata
	isDownloaded := false
	if !image.IsOriginalSaved {
		metadata, originalData, err = c.downloadImage(image)
		if err != nil {
			return err
		}
		isDownloaded = true
		if metadata.Size <= MaxSize {
			name := c.getUploadFileName(image, "origin")
			uploadOutput, err := c.uploader.Upload(name, bytes.NewReader(*originalData), types.ObjectCannedACLBucketOwnerRead)
			if err != nil {
				c.handleError(image, err)
				return err
			}
			image.ProcessingStatus = pictures.Processing
			image.OriginalPath = uploadOutput.Location
			image.OriginalS3Id = uploadOutput.UploadID
			image.IsOriginalSaved = true
		} else {
			image.ProcessingStatus = pictures.TooBig
		}

		err = c.pictureRepo.UpdateImageHandle(image)
		if err != nil {
			c.handleError(image, err)
			return err
		}
		if image.ProcessingStatus == pictures.TooBig {
			return nil
		}
		log.Println("IsOriginalSaved")
	}
	if !image.IsPreviewSaved {
		if !isDownloaded {
			_, originalData, err = c.downloadImage(image)
			if err != nil {
				return err
			}
		}

		extension, err := c.getExtension(image.DropboxPath)
		if err != nil {
			c.handleError(image, err)
			return err
		}

		compressBuffer, ok := c.compressor.Compress(originalData, 90, 300, extension)
		if !ok {
			err = pictures.ErrFailedScale
			c.handleError(image, err)
			return err
		}
		name := c.getUploadFileName(image, "thumb")
		uploadOutput, err := c.uploader.Upload(name, compressBuffer, types.ObjectCannedACLPublicRead)
		if err != nil {
			c.handleError(image, err)
			return err
		}
		image.IsPreviewSaved = true
		image.PreviewPath = uploadOutput.Location
		image.PreviewS3Id = uploadOutput.UploadID

		image.ProcessingStatus = pictures.Processing

		err = c.pictureRepo.UpdateImageHandle(image)
		if err != nil {
			c.handleError(image, err)
			return err
		}
	}
	if !image.IsTextRecognized {
		if !isDownloaded {
			_, originalData, err = c.downloadImage(image)
			if err != nil {
				return err
			}
		}

		var detectedText []pictures.TextDetection
		detectedText, err = c.textDetector.DetectTextFromImage(originalData, c.minConfidence)
		if err != nil {
			c.handleError(image, err)
			return err
		}

		image.IsTextRecognized = true
		image.DetectedTexts = detectedText
		image.ProcessingStatus = pictures.Processing
		err = c.pictureRepo.UpdateImageHandle(image)
		if err != nil {
			c.handleError(image, err)
			return err
		}
	}

	image.ProcessingStatus = pictures.Success

	return c.pictureRepo.UpdateImageHandle(image)
}

func (c *CoordinatorServiceImpl) getExtension(path string) (pictures.SupportedImgType, error) {
	if strings.LastIndex(path, "."+string(pictures.PNG)) != -1 {
		return pictures.PNG, nil
	}
	if strings.LastIndex(path, "."+string(pictures.JPEG)) != -1 {
		return pictures.JPEG, nil
	}
	if strings.LastIndex(path, "."+string(pictures.JPG)) != -1 {
		return pictures.JPG, nil
	}

	return "", pictures.ErrUnsupportedType
}

func (c *CoordinatorServiceImpl) downloadImage(image *pictures.Picture) (*files.FileMetadata, *[]byte, error) {
	metadata, data, err := c.downloader.DownloadFile(image.DropboxPath)
	if err != nil {
		c.handleError(image, err)
		return nil, data, err
	}
	return metadata, data, nil
}

func (c *CoordinatorServiceImpl) getUploadFileName(image *pictures.Picture, typeSize string) string {
	return strconv.Itoa(image.EventId) + "/" + image.Id.String() + "-" + typeSize + ".jpg"
}

func (c *CoordinatorServiceImpl) handleError(image *pictures.Picture, err error) {
	now := time.Now()
	image.Attempts++
	if image.Attempts > c.maxAttemptsBeforeNotify {
		err = c.notifier.Notify(domain.Message{
			Message: "failed processing image: " + image.Id.String(),
		})
		if err != nil {
			log.Println("failed notify developer", err)
		}
	}

	image.ProcessingStatus = pictures.Failed
	image.ExecuteAfter = now.Add(time.Duration(image.Attempts*1) * time.Minute)
	e := c.pictureRepo.UpdateImageHandle(image)
	if e != nil {
		log.Println(err, image, "save")
	}
}
