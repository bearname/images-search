package tasks

import (
	"photofinish/pkg/domain/dto"
)

type Repo interface {
	Store(task *AddImageDto) error
	GetStatsByTask(taskId string) (*TaskStats, error)
	GetTaskList(page *dto.Page) (*[]TaskReturnDTO, error)
}
