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

func NewPictureProcessor(downloader *dropbox.SDKDownloader,
	pictureRepo pictures.Repo,
	amqpChannel *amqp.Channel,
	topicImageProcessing string,
	outboxRepo broker.Repo) *ProcessorImpl {
	p := new(ProcessorImpl)
	p.downloader = downloader
	p.pictureRepo = pictureRepo
	p.amqpChannel = amqpChannel
	p.topicImageProcessing = topicImageProcessing
	p.outboxRepo = outboxRepo
	return p
}

func (s *ProcessorImpl) PerformAddImagesToQueue(t *tasks.Task) error {
	images, err := s.retrieveImages(t)
	if err != nil {
		return err
	}

	data, err := json.Marshal(images)
	if err != nil {
		return err
	}
	err = rabbitmq.PublishToQueue(s.amqpChannel, s.topicImageProcessing, data)
	if err != nil {
		return err
	}
	err = s.outboxRepo.UpdateStatus(t.Id, broker.OutboxDone)

	return err
}

func (s *ProcessorImpl) retrieveImages(t *tasks.Task) (*pictures.DropboxImages, error) {
	dropboxPath := t.DropboxPath
	images, err := s.downloader.GetListFiles(dropboxPath)
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
