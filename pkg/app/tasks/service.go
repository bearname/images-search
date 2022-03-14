package tasks

import (
	"github.com/col3name/images-search/pkg/common/util"
	"github.com/col3name/images-search/pkg/domain/domainerror"
	"github.com/col3name/images-search/pkg/domain/dto"
	"github.com/col3name/images-search/pkg/domain/tasks"
	"github.com/pkg/errors"
)

type ServiceImpl struct {
	repo tasks.Repo
}

func NewService(eventRepo tasks.Repo) *ServiceImpl {
	s := new(ServiceImpl)
	s.repo = eventRepo
	return s
}

func (s *ServiceImpl) Store(task *tasks.AddImageDto) error {
	if task == nil {
		return domainerror.ErrNilObject
	}
	return s.repo.Store(task)
}

func (s *ServiceImpl) GetTaskStatistic(taskId string) (*tasks.TaskStats, error) {
	if !util.IsUUID(taskId) {
		return nil, errors.New("invalid task id")
	}

	return s.repo.GetStatsByTask(taskId)
}

func (s *ServiceImpl) GetTasks(page *dto.Page) (*[]tasks.TaskReturnDTO, error) {
	return s.repo.GetTaskList(page)
}
