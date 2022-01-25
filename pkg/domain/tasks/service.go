package tasks

type Service interface {
	GetStatistics(taskId string) (*TaskStats, error)
}
