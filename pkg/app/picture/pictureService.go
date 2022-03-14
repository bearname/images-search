package picture

import (
	"encoding/json"
	rabbitmq "github.com/col3name/images-search/pkg/common/infrarstructure/amqp"
	"github.com/col3name/images-search/pkg/common/util/uuid"
	"github.com/col3name/images-search/pkg/domain"
	"github.com/col3name/images-search/pkg/domain/domainerror"
	"github.com/col3name/images-search/pkg/domain/pictures"
	"github.com/col3name/images-search/pkg/domain/tasks"
)

type ServiceImpl struct {
	pictureRepo       pictures.Repo
	amqpChannel       rabbitmq.BrokerService
	downloader        domain.Downloader
	addNewImagesTopic string
	tasksService      tasks.Service
}

func NewPictureService(pictureRepo pictures.Repo,
	channel rabbitmq.BrokerService,
	downloader domain.Downloader,
	brokerTopic string,
	tasksService tasks.Service,
) *ServiceImpl {
	s := new(ServiceImpl)
	s.pictureRepo = pictureRepo
	s.amqpChannel = channel
	s.downloader = downloader
	s.addNewImagesTopic = brokerTopic
	s.tasksService = tasksService
	return s
}

func (s *ServiceImpl) Create(image *pictures.TextDetectionOnImageDto) error {
	if image == nil {
		return domainerror.ErrNilObject
	}
	return s.pictureRepo.Store(image)
}

func (s *ServiceImpl) DetectImageFromUrl(dropboxPath string, eventId int) (*pictures.TaskResponse, error) {
	taskId := uuid.Generate().String()

	data, t, err := s.init(dropboxPath, taskId, eventId)
	if err != nil {
		return nil, err
	}

	err = s.tasksService.Store(t)
	if err != nil {
		return nil, err
	}
	err = s.amqpChannel.PublishToQueue(s.addNewImagesTopic, data)
	return &pictures.TaskResponse{TaskId: taskId}, err
}

func (s *ServiceImpl) init(dropboxPath, taskId string, eventId int) ([]byte, *tasks.AddImageDto, error) {
	task := tasks.Task{
		Id:          taskId,
		EventId:     eventId,
		DropboxPath: dropboxPath,
	}
	data, err := json.Marshal(tasks.Task{
		Id:          taskId,
		EventId:     eventId,
		DropboxPath: dropboxPath,
	})
	if err != nil {
		return nil, nil, err
	}
	return data, &tasks.AddImageDto{
		BrokerTopic: s.addNewImagesTopic,
		TaskData:    string(data),
		Task:        task,
	}, nil
}

func (s *ServiceImpl) Search(dto *pictures.SearchPictureDto) (*pictures.SearchPictureResultDto, error) {
	return s.pictureRepo.Search(dto)
}

func (s *ServiceImpl) GetDropboxFolders() ([]string, error) {
	return s.downloader.GetListFolder("")
}

func (s *ServiceImpl) Delete(pictureId string) error {
	//img, err := (*s.pictureRepo).FindPicture(pictureId)
	//if err != nil {
	//	log.Error(err)
	//	return pictures.ErrNotFound
	//}
	//var ids []types.ObjectIdentifier
	//if img.IsOriginalSaved {
	//	ids = append(ids, types.ObjectIdentifier{
	//		Key: aws.String(img.OriginalS3Id),
	//	})
	//}
	//if img.IsPreviewSaved {
	//	ids = append(ids, types.ObjectIdentifier{Key: aws.String(img.PreviewS3Id)})
	//}
	//
	//_, err = s.s3Client.DeleteObjects(context.Background(), &s3.DeleteObjectsInput{
	//	Bucket: aws.String(s.bucket),
	//	Delete: &types.Delete{
	//		Objects: ids,
	//		Quiet:   true,
	//	},
	//})
	//
	//if err != nil {
	//	return err
	//}

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
	return s.pictureRepo.Delete(pictureId)
}
