package tasks

import (
	"photofinish/pkg/domain/dto"
)

type Repository interface {
	GetStatsByTask(taskId string) (*TaskStats, error)
	GetTaskList(page *dto.Page) (*[]TaskReturnDTO, error)
}
