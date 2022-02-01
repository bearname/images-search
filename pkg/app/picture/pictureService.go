package picture

import (
	"encoding/json"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"photofinish/pkg/app/dropbox"
	rabbitmq "photofinish/pkg/common/infrarstructure/amqp"
	"photofinish/pkg/domain/pictures"
	"photofinish/pkg/infrastructure/postgres"
)

type ServiceImpl struct {
	pictureRepo *postgres.PictureRepositoryImpl
	amqpChannel *amqp.Channel
	downloader  *dropbox.SDKDownloader
	s3Client    *s3.Client
	bucket      string
}

func NewPictureService(pictureRepo *postgres.PictureRepositoryImpl, amqpChannel *amqp.Channel, downloader *dropbox.SDKDownloader, s3Client *s3.Client, bucket string) *ServiceImpl {
	s := new(ServiceImpl)
	s.pictureRepo = pictureRepo
	s.amqpChannel = amqpChannel
	s.downloader = downloader
	s.s3Client = s3Client
	s.bucket = bucket
	return s
}

func (s *ServiceImpl) Create(imageTextDetectionDto *pictures.TextDetectionOnImageDto) error {
	return (*s.pictureRepo).Store(imageTextDetectionDto)
}

func (s *ServiceImpl) DetectImageFromUrl(dropboxPath string, eventId int) (*pictures.TaskResponse, error) {
	images, err := s.downloader.GetListFolder(dropboxPath, true, true)
	if err != nil {
		return nil, err
	}

	image := pictures.InitialDropboxImage{
		Images: images, EventId: eventId, Path: dropboxPath,
	}
	result, err := (*s.pictureRepo).SaveInitialPictures(&image)

	if err != nil {
		return nil, err

	}
	var dropboxImages pictures.DropboxImages
	dropboxImages.EventId = eventId
	id := result.ImagesId
	for i, img := range images {
		dropboxImages.Images = append(dropboxImages.Images, pictures.DropboxImage{
			Path: img,
			Id:   id[i],
		})
	}
	//TODO
	// type Task struct {
	//  TaskId int
	//  DropboxPath string
	//  CountImage int
	//  Status ProcessingStatus
	// }
	//
	//TODO
	// type Picture struct {
	//    TaskId uuid
	//    DropboxPath string
	//    TaskId int
	// }
	// pictureRepo.SaveInitialPictures(image pictures.InitialDropboxImage) {
	//   begin
	//    insert into tasks values (dropboxPath, Processing, len(image.images)) returning id
	//    insert into pictures values (pictureData, taskId)
	//   commit
	// }
	// type TaskStatus struct {
	//   Status,
	//   Percent,
	// }
	//  long pool from frontend
	// .controller GET /api/v1/tasks/{taskId}
	//  GetProcessingStatus(taskId) TaskStatus {
	//    countCompletedTask = select count(id) from pictures where task_id = %taskId% and processing_status="Complete"
	//    countImagesInTask = select countImage from tasks where id = $task_id$
	//    percent = countCompletedTask / countImagesInTask * 100
	//    return TaskStatus { Percent: countCompletedTask / countImagesInTask * 100, Status: percent = 100 ? "success": "processing"}
	//  }
	//  .controller GET /api/v1/tasks
	//  GetTasks() []Tasks {
	//     return select * from tasks;
	//  }
	//  .controller DELETE /api/v1/tasks/{taskId}
	//  CancelTask(taskId) error {
	//    picturesInTask = select id, previewPath, originalPath from pictures where task_id = %{taskId}
	//    picturesInTaskDeleteImageFromS3(picturesInTask)
	//    for _, picture := range picturesInTask {
	//      err := this.Delete(pictureId)
	//      if err != nil {
	//         this.Delete(pictureId)
	//      }
	//    }
	//  }
	//  .
	//  .
	//  . вывод списка task на фронте по клике на кнопку "задачи"
	//

	marshal, err := json.Marshal(dropboxImages)
	if err != nil {
		return nil, err
	}
	message := amqp.Publishing{
		ContentType: "text/plain",
		Body:        marshal,
	}

	return &pictures.TaskResponse{
		TaskId:          result.TaskId.String(),
		CountAllImages:  len(images),
		CompletedImages: 0,
	}, rabbitmq.Publish(s.amqpChannel, rabbitmq.TargetQueue, message)
}

func (s *ServiceImpl) Search(dto *pictures.SearchPictureDto) (*pictures.SearchPictureResultDto, error) {
	return (*s.pictureRepo).Search(dto)
}

func (s *ServiceImpl) GetDropboxFolders() ([]string, error) {
	return s.downloader.GetListFolder("", false, false)
}

func (s *ServiceImpl) Delete(pictureId string) error {
	err := (*s.pictureRepo).IsExists(pictureId)
	if err != nil {
		log.Error(err)
		return pictures.ErrNotFound
	}

	// TODO
	//  deletePreviewImage
	//  deleteOriginalImage
	//   _, err := svc.DeleteObject(&s3.DeleteObjectInput{
	//        Bucket: bucket,
	//        Key:    item,
	//    })
	//    if err != nil {
	//        return err
	//    }
	//   .
	//    err = svc.WaitUntilObjectNotExists(&s3.HeadObjectInput{
	//        Bucket: bucket,
	//        Key:    item,
	//    })
	//    if err != nil {
	//        return err
	//    }
	//   .
	//    return nil
	return (*s.pictureRepo).Delete(pictureId)
}
