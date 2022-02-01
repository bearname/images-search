package tasks

import (
	"photofinish/pkg/domain/dto"
)

type Service interface {
	GetTaskStatistic(taskId string) (*TaskStats, error)
	GetTasks(page *dto.Page) (*[]TaskReturnDTO, error)
}
