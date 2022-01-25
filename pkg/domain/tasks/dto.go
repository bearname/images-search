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
