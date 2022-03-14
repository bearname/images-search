package tasks

import (
	"github.com/col3name/images-search/pkg/domain/dto"
)

type Service interface {
	Store(task *AddImageDto) error
	GetTaskStatistic(taskId string) (*TaskStats, error)
	GetTasks(page *dto.Page) (*[]TaskReturnDTO, error)
}
