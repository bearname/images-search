package picture

import (
	"bytes"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"photofinish/pkg/app/aws/recognition"
	"photofinish/pkg/domain"
	"photofinish/pkg/domain/pictures"
	"photofinish/pkg/infrastructure/postgres"
	"strconv"
	"strings"
	"time"
)

type CoordinatorServiceImpl struct {
	maxAttemptsBeforeNotify int
	pictureRepo             *postgres.PictureRepositoryImpl
	downloader              domain.Downloader
	uploader                domain.Uploader
	textDetector            recognition.AmazonTextRecognition
	compressor              domain.ImageCompressor

	minConfidence int
}

func NewCoordinatorServiceImpl(maxAttemptsBeforeNotify int, pictureRepo *postgres.PictureRepositoryImpl,
	downloader domain.Downloader,
	uploader domain.Uploader,
	textDetector *recognition.AmazonTextRecognition,
	compressor domain.ImageCompressor,
	minConfidence int) *CoordinatorServiceImpl {

	c := new(CoordinatorServiceImpl)

	c.maxAttemptsBeforeNotify = maxAttemptsBeforeNotify
	c.pictureRepo = pictureRepo
	c.downloader = downloader
	c.uploader = uploader
	c.textDetector = *textDetector
	c.compressor = compressor
	c.minConfidence = minConfidence

	return c
}

func (c *CoordinatorServiceImpl) PerformAddImage(image *pictures.Picture) error {
	var originalData *[]byte
	var err error
	isDownloaded := false
	if !image.IsOriginalSaved {
		originalData, err = c.downloadImage(image)
		if err != nil {
			return err
		}

		isDownloaded = true
		name := c.getUploadFileName(image, "origin")
		uploadOutput, err := c.uploader.Upload(name, bytes.NewReader(*originalData), types.ObjectCannedACLBucketOwnerRead)
		if err != nil {
			c.handleError(image, err)
			return err
		}
		image.ProcessingStatus = pictures.Processing
		image.IsOriginalSaved = true
		image.OriginalPath = uploadOutput.Location
		err = c.pictureRepo.UpdateImageHandle(image)
		if err != nil {
			c.handleError(image, err)
			return err
		}
		fmt.Println("IsOriginalSaved")
	}
	if !image.IsPreviewSaved {
		if !isDownloaded {
			originalData, err = c.downloadImage(image)
			if err != nil {
				return err
			}
		}

		extension, err := getExtension(image.DropboxPath)
		if err != nil {
			err = errors.New("Failed scale image")
			c.handleError(image, err)
			return err
		}

		compressBuffer, ok := c.compressor.Compress(originalData, 90, 300, extension)
		if !ok {
			err = errors.New("Failed scale image")
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
		image.ProcessingStatus = pictures.Processing
		err = c.pictureRepo.UpdateImageHandle(image)
		if err != nil {
			c.handleError(image, err)
			return err
		}
	}
	if !image.IsTextRecognized {
		if !isDownloaded {
			originalData, err = c.downloadImage(image)
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

		//log.Println(detectedText)
		image.IsTextRecognized = true
		image.DetectedTexts = detectedText
		image.ProcessingStatus = pictures.Processing
		err = c.pictureRepo.UpdateImageHandle(image)
		if err != nil {
			c.handleError(image, err)
			return err
		}
	}

	//if !image.IsMobileSaved {
	//    if !isDownloaded {
	//        originalData, err = c.downloadImage(image)
	//        if err != nil {
	//            return err
	//        }
	//    }
	//
	//    fmt.Println("IsMobileSaved")
	//}

	image.ProcessingStatus = pictures.Success

	return c.pictureRepo.UpdateImageHandle(image)
}

func getExtension(path string) (string, error) {
	if strings.LastIndex(path, ".png") != -1 {
		return "png", nil
	}
	if strings.LastIndex(path, ".jpeg") != -1 {
		return "jpeg", nil
	}
	if strings.LastIndex(path, ".jpg") != -1 {
		return "jpg", nil
	}

	return "", errors.New("unknown extension")

}

func (c *CoordinatorServiceImpl) downloadImage(image *pictures.Picture) (*[]byte, error) {
	_, data, err := c.downloader.DownloadFile(image.DropboxPath)
	if err != nil {
		c.handleError(image, err)
		return data, err
	}
	return data, nil
}

func (c *CoordinatorServiceImpl) getUploadFileName(image *pictures.Picture, typeSize string) string {
	return strconv.Itoa(image.EventId) + "/" + image.Id.String() + "-" + typeSize + ".jpg"
}

func (c *CoordinatorServiceImpl) handleError(image *pictures.Picture, err error) {
	now := time.Now()
	image.Attempts++
	if image.Attempts > c.maxAttemptsBeforeNotify {

		log.Println("Notify developer")
		//err := c.notifier.Notify(pictures)
		//if err != nil {
		//    log.Println("failed notify developer")
		//    log.Println(err, pictures, "err")
		//}
	}

	image.ProcessingStatus = pictures.Failed
	image.ExecuteAfter = now.Add(time.Duration(image.Attempts*1) * time.Minute)
	e := c.pictureRepo.UpdateImageHandle(image)
	if e != nil {
		log.Println(err, image, "save")
	}
}
