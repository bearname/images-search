package tasks

import (
	"github.com/col3name/images-search/pkg/domain/dto"
)

type Repo interface {
	Store(task *AddImageDto) error
	GetStatsByTask(taskId string) (*TaskStats, error)
	GetTaskList(page *dto.Page) (*[]TaskReturnDTO, error)
}
