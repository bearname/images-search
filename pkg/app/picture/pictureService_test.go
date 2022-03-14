package picture

import (
	"errors"
	"fmt"
	"github.com/col3name/images-search/pkg/domain/domainerror"
	"github.com/col3name/images-search/pkg/domain/dto"
	"github.com/col3name/images-search/pkg/domain/pictures"
	"github.com/col3name/images-search/pkg/domain/tasks"
	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox/files"
	"github.com/streadway/amqp"
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"

	"github.com/stretchr/testify/mock"
)

// smsServiceMock
type smsServiceMock struct {
	mock.Mock
}

type IsmsService interface {
	SendChargeNotification(value int) bool
}

// Our mocked smsService method
func (m *smsServiceMock) SendChargeNotification(value int) bool {
	fmt.Println("Mocked charge notification function")
	fmt.Printf("Value passed in: %d\n", value)
	// this records that the method was called and passes in the value
	// it was called with
	args := m.Called(value)
	// it then returns whatever we tell it to return
	// in this case true to simulate an SMS Service Notification
	// sent out
	return args.Bool(0)
}

// we need to satisfy our MessageService interface
// which sadly means we have to stub out every method
// defined in that interface
func (m *smsServiceMock) DummyFunc() {
	fmt.Println("Dummy")
}

type MyService struct {
	SmsService IsmsService
}

func (s *MyService) ChargeCustomer(val int) {
	s.SmsService.SendChargeNotification(val)
}

func TestChargeCustomer(t *testing.T) {
	smsService := new(smsServiceMock)
	smsService.On("SendChargeNotification", 100).Return(true)

	myService := MyService{smsService}
	myService.ChargeCustomer(100)
	smsService.AssertExpectations(t)
}

type MockAmqpService struct {
	mock.Mock
}

func (s *MockAmqpService) Consume(queueName string) (<-chan amqp.Delivery, error) {
	args := s.Called(queueName)
	var ch <-chan amqp.Delivery
	return ch, args.Error(0)
}

func (s *MockAmqpService) PublishToQueue(queueName string, data []byte) error {
	args := s.Called(queueName, data)
	return args.Error(0)
}

type MockDownloader struct {
	mock.Mock
}

func (s *MockDownloader) GetListFiles(path string) ([]string, error) {
	args := s.Called(path)
	return []string{}, args.Error(0)
}

func (s *MockDownloader) GetListFolder(path string) ([]string, error) {
	args := s.Called(path)
	return []string{}, args.Error(0)
}

func (s *MockDownloader) DownloadFile(path string) (*files.FileMetadata, *[]byte, error) {
	args := s.Called(path)
	return &files.FileMetadata{}, &[]byte{}, args.Error(0)
}

type MockPictureRepo struct {
	mock.Mock
}

func (s *MockPictureRepo) IsExists(pictureId string) error {
	args := s.Called(pictureId)
	return args.Error(0)
}

func (s *MockPictureRepo) FindPicture(pictureId string) (*pictures.PictureDTO, error) {
	args := s.Called(pictureId)
	return args.Get(0).(*pictures.PictureDTO), args.Error(1)
}

func (s *MockPictureRepo) Search(dto *pictures.SearchPictureDto) (*pictures.SearchPictureResultDto, error) {
	args := s.Called(dto)
	return args.Get(0).(*pictures.SearchPictureResultDto), args.Error(1)
}

func (s *MockPictureRepo) SaveInitialPicture(image *pictures.InitialImage) (int, error) {
	args := s.Called(image)
	return args.Int(0), args.Error(1)
}

func (s *MockPictureRepo) SaveInitialPictures(image *pictures.InitialDropboxImage) (*pictures.InitialDropboxImageResult, error) {
	args := s.Called(image)
	return args.Get(0).(*pictures.InitialDropboxImageResult), args.Error(1)
}

func (s *MockPictureRepo) UpdateImageHandle(image *pictures.Picture) error {
	args := s.Called(image)
	return args.Error(0)
}

func (s *MockPictureRepo) Store(image *pictures.TextDetectionOnImageDto) error {
	args := s.Called(image)
	return args.Error(0)
}

func (s *MockPictureRepo) Delete(imageId string) error {
	args := s.Called(imageId)
	return args.Error(0)
}

type MockTaskService struct {
	mock.Mock
}

func (s *MockTaskService) Store(task *tasks.AddImageDto) error {
	args := s.Called(task)
	err := args.Error(0)
	fmt.Println(err)
	return err
}

func (s *MockTaskService) GetTaskStatistic(taskId string) (*tasks.TaskStats, error) {
	args := s.Called(taskId)
	return args.Get(0).(*tasks.TaskStats), args.Error(1)
}

func (s *MockTaskService) GetTasks(page *dto.Page) (*[]tasks.TaskReturnDTO, error) {
	args := s.Called(page)
	return args.Get(0).(*[]tasks.TaskReturnDTO), args.Error(1)
}

var errDb = errors.New("db error")
var errQueue = errors.New("failed publish to queue")

func TestCreatePicture(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		pictureRepo, _, _, _, pictureService := getFixture()
		pictureRepo.On("Store", mock.Anything).Return(true, nil)

		err := pictureService.Create(nil)
		assert.Error(t, err, domainerror.ErrNilObject)
		assert.Equal(t, errors.Is(err, domainerror.ErrNilObject), true)
	})
	imgDto := pictures.TextDetectionOnImageDto{
		EventId:       0,
		OriginalPath:  "/images",
		PreviewPath:   "/images",
		TextDetection: []pictures.TextDetection{},
	}
	t.Run("not empty", func(t *testing.T) {
		pictureRepo, _, _, _, pictureService := getFixture()
		pictureRepo.On("Store", &imgDto).Return(nil)

		err := pictureService.Create(&imgDto)
		assert.Nil(t, err)
		pictureRepo.AssertExpectations(t)
	})
	t.Run("db error", func(t *testing.T) {
		pictureRepo, _, _, taskService, pictureService := getFixture()
		taskService.On("Store", mock.Anything).Return(errDb)
		err := pictureService.Create(&imgDto)
		pictureRepo.AssertExpectations(t)
		assert.Error(t, err, errDb)
	})
}

func getFixture() (*MockPictureRepo, *MockAmqpService, *MockDownloader, *MockTaskService, *ServiceImpl) {
	pictureRepo := new(MockPictureRepo)
	amqpChannel := new(MockAmqpService)
	downloader := new(MockDownloader)
	taskService := new(MockTaskService)
	pictureService := NewPictureService(pictureRepo, amqpChannel, downloader, "topic", taskService)
	return pictureRepo,
		amqpChannel,
		downloader,
		taskService,
		pictureService
}

func TestAddDropboxFolder(t *testing.T) {
	t.Run("db error", func(t *testing.T) {
		pictureRepo := new(MockPictureRepo)
		amqpChannel := new(MockAmqpService)
		downloader := new(MockDownloader)
		taskService := new(MockTaskService)

		taskService.On("Store", mock.AnythingOfType("*tasks.AddImageDto")).Return(errDb)
		called := amqpChannel.AssertNotCalled(t, "PublishToQueue", mock.Anything)
		pictureService := NewPictureService(pictureRepo, amqpChannel, downloader, "topic", taskService)

		res, err := pictureService.DetectImageFromUrl("/images", 1)
		assert.Nil(t, res)
		assert.Equal(t, true, called)
		pictureRepo.AssertExpectations(t)
		assert.Error(t, err, errDb)
	})
	t.Run("publish error", func(t *testing.T) {
		pictureRepo := new(MockPictureRepo)
		amqpChannel := new(MockAmqpService)
		downloader := new(MockDownloader)
		taskService := new(MockTaskService)

		taskService.On("Store", mock.AnythingOfType("*tasks.AddImageDto")).Return(nil)
		amqpChannel.On("PublishToQueue", mock.Anything, mock.Anything).Return(errQueue)
		pictureService := NewPictureService(pictureRepo, amqpChannel, downloader, "topic", taskService)

		_, err := pictureService.DetectImageFromUrl("/images", 1)
		assert.Nil(t, err)
		pictureRepo.AssertExpectations(t)
		taskService.AssertExpectations(t)
	})
	t.Run("ok", func(t *testing.T) {
		pictureRepo := new(MockPictureRepo)
		amqpChannel := new(MockAmqpService)
		downloader := new(MockDownloader)
		taskService := new(MockTaskService)
		pictureService := NewPictureService(pictureRepo, amqpChannel, downloader, "topic", taskService)

		taskService.On("Store", mock.AnythingOfType("*tasks.AddImageDto")).Return(nil)
		amqpChannel.On("PublishToQueue", mock.Anything, mock.Anything).Return(nil)

		_, err := pictureService.DetectImageFromUrl("/images", 1)
		assert.Nil(t, err)
		pictureRepo.AssertExpectations(t)
		taskService.AssertExpectations(t)
	})
}

func TestServiceImpl_Search(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		pictureRepo := new(MockPictureRepo)
		amqpChannel := new(MockAmqpService)
		downloader := new(MockDownloader)
		taskService := new(MockTaskService)
		pictureDto := pictures.SearchPictureDto{
			ParticipantNumber: 80,
			Confidence:        12,
			Page: dto.Page{
				Offset: 0,
				Limit:  5,
			},
			EventId: 1,
		}
		var in pictures.SearchPictureResultDto
		var items []pictures.SearchPictureItem
		var item pictures.SearchPictureItem
		var j pictures.TextDetectionDto
		for i := 0; i < 5; i++ {
			itoa := strconv.Itoa(i)
			item.Path = "/images" + itoa
			item.PictureId = itoa
			j.Confidence = 86
			j.EventId = int64(pictureDto.EventId)
			j.EventName = "event" + itoa
			j.DetectedText = itoa
			item.TextDetections = append(item.TextDetections)
			items = append(items, item)
		}
		in.CountAllItems = 5
		in.Pictures = items
		pictureRepo.On("Search", &pictureDto).Return(in, nil)
		pictureService := NewPictureService(pictureRepo, amqpChannel, downloader, "topic", taskService)
		res, err := pictureService.Search(&pictureDto)
		assert.Nil(t, err)
		assert.NotNil(t, *res)
		ress := *res
		assert.Equal(t, ress.CountAllItems, items)
		assert.Equal(t, len(ress.Pictures), len(in.Pictures))
		assert.Equal(t, len(ress.Pictures), len(in.Pictures))
		for i, picture := range ress.Pictures {
			s := in.Pictures[i]
			assert.Equal(t, picture.PictureId, s.PictureId)
			assert.Equal(t, picture.Path, s.Path)
			assert.Equal(t, len(picture.TextDetections), len(s.TextDetections))
			for k, it := range picture.TextDetections {
				jt := s.TextDetections[k]
				assert.Equal(t, it.EventId, jt.EventId)
				assert.Equal(t, it.EventName, jt.EventName)
				assert.Equal(t, it.DetectedText, jt.DetectedText)
				assert.Equal(t, it.Confidence, jt.Confidence)
			}
		}
	})
}
