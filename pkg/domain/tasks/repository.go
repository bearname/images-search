package tasks

type Repository interface {
	GetStatsByTask(taskId string) (*TaskStats, error)
}
