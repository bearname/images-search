package picture

import (
	"encoding/json"
	"github.com/streadway/amqp"
	"photofinish/pkg/app/dropbox"
	rabbitmq "photofinish/pkg/common/infrarstructure/amqp"
	"photofinish/pkg/domain/broker"
	"photofinish/pkg/domain/pictures"
	"photofinish/pkg/domain/tasks"
)

type Processor interface {
	GetImages(t *tasks.Task) ([]string, error)
}

type ProcessorImpl struct {
	downloader           *dropbox.SDKDownloader
	pictureRepo          pictures.Repo
	amqpChannel          *amqp.Channel
	topicImageProcessing string
	outboxRepo           broker.Repo
}

func NewPictureProcessor(downloader *dropbox.SDKDownloader, pictureRepo pictures.Repo, amqpChannel *amqp.Channel, topicImageProcessing string, outboxRepo broker.Repo) *ProcessorImpl {
	p := new(ProcessorImpl)
	p.downloader = downloader
	p.pictureRepo = pictureRepo
	p.amqpChannel = amqpChannel
	p.topicImageProcessing = topicImageProcessing
	p.outboxRepo = outboxRepo
	return p
}

func (s *ProcessorImpl) PerformAddImagesToQueue(t *tasks.Task) error {
	dropboxImages, err := s.retrieveImages(t)
	if err != nil {
		return err
	}
	err = s.publishToQueue(dropboxImages)
	if err != nil {
		return err
	}
	err = s.outboxRepo.UpdateStatus(t.Id, broker.OutboxDone)

	return err

	//TODO
	// type TaskData struct {
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
}

func (s *ProcessorImpl) retrieveImages(t *tasks.Task) (*pictures.DropboxImages, error) {
	dropboxPath := t.DropboxPath
	images, err := s.downloader.GetListFolder(dropboxPath, true, true)
	if err != nil {
		return nil, err
	}
	eventId := t.EventId
	image := pictures.InitialDropboxImage{
		Images: images, EventId: eventId, Path: dropboxPath,
	}

	result, err := s.pictureRepo.SaveInitialPictures(&image)
	if err != nil {
		return nil, err
	}
	dropboxImages := s.fillImages(eventId, result, images)
	return &dropboxImages, nil
}

func (s *ProcessorImpl) fillImages(eventId int, result *pictures.InitialDropboxImageResult, images []string) pictures.DropboxImages {
	var dropboxImages pictures.DropboxImages
	dropboxImages.EventId = eventId
	id := result.ImagesId
	for i, img := range images {
		dropboxImages.Images = append(dropboxImages.Images, pictures.DropboxImage{
			Path: img,
			Id:   id[i],
		})
	}
	return dropboxImages
}

func (s *ProcessorImpl) publishToQueue(dropboxImages *pictures.DropboxImages) error {
	data, err := json.Marshal(&dropboxImages)
	if err != nil {
		return err
	}
	message := amqp.Publishing{
		ContentType: "text/plain",
		Body:        data,
	}

	err = rabbitmq.Publish(s.amqpChannel, s.topicImageProcessing, message)
	if err != nil {
		return err
	}
	return err
}
