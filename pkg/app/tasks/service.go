package tasks

import (
	"github.com/pkg/errors"
	"photofinish/pkg/common/util"
	"photofinish/pkg/domain/domainerror"
	"photofinish/pkg/domain/dto"
	"photofinish/pkg/domain/tasks"
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
