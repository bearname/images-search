package tasks

import (
	"photofinish/pkg/domain/pictures"
	"time"
)

type TaskStatsItem struct {
	Status pictures.ProcessingStatus
	Count  int
}

type TaskStats struct {
	Stats         []TaskStatsItem
	StartedAt     time.Time
	LastUpdatedAt time.Time
}

type TaskReturnDTO struct {
	Id           string
	IsCompleted  bool
	CountImages  int
	StartedAt    time.Time
	LastUpdateAt time.Time
}
